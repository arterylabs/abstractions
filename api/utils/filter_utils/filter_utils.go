package utils

import (
	"abstractions/utils/transformer_utils"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Filter struct {
	OriginalKey string
	Type        string
}

func BuildFilterString(filters map[string][]string, allowedFilters map[string]Filter) (string, []interface{}, error) {
	var inputArgs []interface{}
	var filterArrays []string

	// filter key is the url name of the filter used as the lookup for the allowed filters list
	for filterKey, filterValList := range filters {
		if realFilter, ok := allowedFilters[filterKey]; ok {
			if len(filterValList) == 0 {
				continue
			}

			if realFilter.Type == "" || realFilter.Type == "string" {
				filterArrays = append(filterArrays, fmt.Sprintf("%s IN (?)", realFilter.OriginalKey))
			} else if realFilter.Type == "array" {
				filterArrays = append(filterArrays, fmt.Sprintf("array[?] @> %s", realFilter.OriginalKey))
			} else if realFilter.Type == "exists" {
				if len(filterValList) > 0 && filterValList[0] == "true" {
					filterArrays = append(filterArrays, fmt.Sprintf("%s IS NOT NULL", realFilter.OriginalKey))
				} else if len(filterValList) > 0 && filterValList[0] == "false" {
					filterArrays = append(filterArrays, fmt.Sprintf("%s IS NULL", realFilter.OriginalKey))
				}
			} else if realFilter.Type == "boolean" {
				filterArrays = append(filterArrays, fmt.Sprintf("%s = ?", realFilter.OriginalKey))
			}

			if realFilter.Type != "exists" {
				inputArgs = append(inputArgs, filterValList)
			}
		}
	}

	filterString := strings.Join(filterArrays, " AND ")
	// template the where clause into the original query and then expand the IN clauses with sqlx
	query, args, err := sqlx.In(filterString, inputArgs...)
	if err != nil {
		return "", nil, err
	}

	// using postgres means we need to rebind the ? bindvars that sqlx.IN creates by default to $ bindvars
	// you can omit this if you are using mysql
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	return query, args, nil
}

func BuildFilterStringV2(filters map[string]interface{}, allowedFilters map[string]Filter) (string, []interface{}, error) {
	var inputArgs []interface{}
	var filterArrays []string

	// filter key is the url name of the filter used as the lookup for the allowed filters list
	for filterKey, filterVal := range filters {
		if realFilter, ok := allowedFilters[filterKey]; ok {
			if filterVal == nil {
				continue
			}

			if realFilter.Type == "" || realFilter.Type == "string" {
				filterArrays = append(filterArrays, fmt.Sprintf("%s::text=any(array[?])", realFilter.OriginalKey))
				inputArgs = append(inputArgs, transformer_utils.GetInterfaceArrayFromInterface(filterVal))
			} else if realFilter.Type == "array" {
				filterArrays = append(filterArrays, fmt.Sprintf("%s && array[?]", realFilter.OriginalKey))
				inputArgs = append(inputArgs, filterVal)
			} else if realFilter.Type == "json_array" {
				filterArrays = append(filterArrays, fmt.Sprintf(`%s @> array_to_json(array[?])::jsonb`, realFilter.OriginalKey))
				inputArgs = append(inputArgs, filterVal)
			} else if realFilter.Type == "exists" {
				if filterVal != nil && transformer_utils.GetBooleanFromInterface(filterVal) {
					filterArrays = append(filterArrays, fmt.Sprintf("%s IS NOT NULL", realFilter.OriginalKey))
				} else if filterVal != nil && !transformer_utils.GetBooleanFromInterface(filterVal) {
					filterArrays = append(filterArrays, fmt.Sprintf("%s IS NULL", realFilter.OriginalKey))
				}
			} else if realFilter.Type == "boolean" {
				filterArrays = append(filterArrays, fmt.Sprintf("%s = ?", realFilter.OriginalKey))
				inputArgs = append(inputArgs, filterVal)
			}
		}
	}

	filterString := strings.Join(filterArrays, " AND ")
	// template the where clause into the original query and then expand the IN clauses with sqlx
	query, args, err := sqlx.In(filterString, inputArgs...)
	if err != nil {
		return "", nil, err
	}

	// using postgres means we need to rebind the ? bindvars that sqlx.IN creates by default to $ bindvars
	// you can omit this if you are using mysql
	// query = sqlx.Rebind(sqlx.DOLLAR, query)
	return query, args, nil
}

type FilterInput struct {
	Field    string
	Operator string
	Value    []interface{}
}

func GetFilterDetailsByField(field string, allowedFilters map[string]Filter) *Filter {
	for key, filter := range allowedFilters {
		if key == field {
			return &filter
		}
	}
	return nil
}

func BuildFilterStringV3(filters []FilterInput, allowedFilters map[string]Filter) (string, []interface{}, error) {
	var inputArgs []interface{}
	var filterArrays []string

	// filter key is the url name of the filter used as the lookup for the allowed filters list
	for _, filterObj := range filters {
		mappedFilter := GetFilterDetailsByField(filterObj.Field, allowedFilters)
		if mappedFilter != nil {
			if filterObj.Value == nil {
				continue
			}

			if mappedFilter.Type == "" || mappedFilter.Type == "string" {
				if filterObj.Operator == "eq" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s::text= ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "not_eq" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s::text != ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "contains" {
					newFilters := []string{}
					args := []interface{}{}
					for _, val := range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("(%s::text like ? or %s=?)", mappedFilter.OriginalKey, mappedFilter.OriginalKey))
						args = append(args, "%"+transformer_utils.GetStringFromInterface(val)+"%", val)
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, args...)
				} else if filterObj.Operator == "not_contains" {
					newFilters := []string{}
					args := []interface{}{}
					for _, val := range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("(%s::text not like ? and %s!=?)", mappedFilter.OriginalKey, mappedFilter.OriginalKey))
						args = append(args, "%"+transformer_utils.GetStringFromInterface(val)+"%", val)
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, args...)
				} else if filterObj.Operator == "starts_with" {
					newFilters := []string{}
					args := []interface{}{}
					for _, val := range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("(%s::text like ?)", mappedFilter.OriginalKey))
						args = append(args, transformer_utils.GetStringFromInterface(val)+"%")
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, args...)
				} else if filterObj.Operator == "ends_with" {
					newFilters := []string{}
					args := []interface{}{}
					for _, val := range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("(%s::text like ?)", mappedFilter.OriginalKey))
						args = append(args, "%"+transformer_utils.GetStringFromInterface(val))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, args...)
				}
			} else if mappedFilter.Type == "number" {
				if filterObj.Operator == "eq" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s = ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")

					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "not_eq" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s != ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")

					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "gt" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s > ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")

					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "lt" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s < ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")

					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "gt_or_eq" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s >= ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")

					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "lt_or_eq" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s <= ?", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")

					inputArgs = append(inputArgs, filterObj.Value...)
				}
			} else if mappedFilter.Type == "array" {
				if filterObj.Operator == "contains" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("%s && array[?]", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, filterObj.Value...)
				} else if filterObj.Operator == "not_contains" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf("(%s && array[?]) is distinct from true", mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, filterObj.Value...)
				}
			} else if mappedFilter.Type == "json_array" {
				if filterObj.Operator == "contains" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf(`%s @> array_to_json(array[?])::jsonb`, mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, filterObj.Value)
				} else if filterObj.Operator == "not_contains" {
					newFilters := []string{}
					for range filterObj.Value {
						newFilters = append(newFilters, fmt.Sprintf(`(%s @> array_to_json(array[?])::jsonb) is distinct from true`, mappedFilter.OriginalKey))
					}
					filterArrays = append(filterArrays, "("+strings.Join(newFilters, " OR ")+")")
					inputArgs = append(inputArgs, filterObj.Value)
				}
			} else if mappedFilter.Type == "exists" {
				if filterObj.Value != nil && transformer_utils.GetBooleanFromInterface(filterObj.Value) {
					filterArrays = append(filterArrays, fmt.Sprintf("%s IS NOT NULL", mappedFilter.OriginalKey))
				} else if filterObj.Value != nil && !transformer_utils.GetBooleanFromInterface(filterObj.Value) {
					filterArrays = append(filterArrays, fmt.Sprintf("%s IS NULL", mappedFilter.OriginalKey))
				}
			} else if mappedFilter.Type == "boolean" {
				filterArrays = append(filterArrays, fmt.Sprintf("%s = ?", mappedFilter.OriginalKey))
				inputArgs = append(inputArgs, filterObj.Value)
			}
		}
	}

	filterString := strings.Join(filterArrays, " AND ")
	// template the where clause into the original query and then expand the IN clauses with sqlx
	query, args, err := sqlx.In(filterString, inputArgs...)
	if err != nil {
		return "", nil, err
	}

	// using postgres means we need to rebind the ? bindvars that sqlx.IN creates by default to $ bindvars
	// you can omit this if you are using mysql
	// query = sqlx.Rebind(sqlx.DOLLAR, query)
	return query, args, nil
}
