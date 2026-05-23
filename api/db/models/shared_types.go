package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringArray is a custom type for PostgreSQL text[] columns.
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "{}", nil
	}
	result := "{"
	for i, v := range s {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%q", v)
	}
	result += "}"
	return result, nil
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	// Parse PostgreSQL array literal {a,b,c}
	if len(str) < 2 || str[0] != '{' || str[len(str)-1] != '}' {
		*s = StringArray{}
		return nil
	}
	inner := str[1 : len(str)-1]
	if inner == "" {
		*s = StringArray{}
		return nil
	}
	var result []string
	if err := json.Unmarshal([]byte(`[`+inner+`]`), &result); err != nil {
		// fallback: simple comma split
		result = splitPGArray(inner)
	}
	*s = result
	return nil
}

// splitPGArray splits a raw PostgreSQL array interior on commas, stripping quotes.
func splitPGArray(s string) []string {
	var result []string
	current := ""
	inQuote := false
	for _, c := range s {
		switch {
		case c == '"':
			inQuote = !inQuote
		case c == ',' && !inQuote:
			result = append(result, current)
			current = ""
		default:
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// JSONB is a custom type for PostgreSQL jsonb columns backed by json.RawMessage.
// It stores arbitrary JSON as raw bytes and preserves the original encoding on round-trips.
type JSONB json.RawMessage

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = JSONB(v)
	default:
		return fmt.Errorf("unsupported type for JSONB: %T", value)
	}
	return nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	*j = append((*j)[0:0], data...)
	return nil
}

// JSONBRaw stores any JSON value (object, array, scalar) as raw bytes in a JSONB column.
// Unlike JSONB, a nil JSONBRaw serialises to the SQL literal "null" rather than SQL NULL.
type JSONBRaw json.RawMessage

func (j JSONBRaw) Value() (driver.Value, error) {
	if j == nil {
		return "null", nil
	}
	return string(j), nil
}

func (j *JSONBRaw) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("JSONBRaw: unsupported type %T", value)
	}
	raw := make(JSONBRaw, len(b))
	copy(raw, b)
	*j = raw
	return nil
}

// MarshalJSON returns the raw JSON bytes.
func (j JSONBRaw) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// JSONBField is a custom type to handle JSONB columns whose value is a JSON object
// (key→value map). Use JSONB or JSONBRaw when the column may hold arrays or scalars.
type JSONBField map[string]interface{}

// Value implements the driver.Valuer interface for database storage.
func (j JSONBField) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval.
func (j *JSONBField) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}
