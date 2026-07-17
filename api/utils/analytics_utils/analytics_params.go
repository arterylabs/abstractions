package analytics_utils

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// DateParam is a custom type that can parse date strings in YYYY-MM-DD format
type DateParam struct {
	time.Time
}

// UnmarshalParam implements echo's param unmarshaling for query parameters
func (dp *DateParam) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}

	// Try parsing as YYYY-MM-DD format first (what frontend sends)
	t, err := time.Parse("2006-01-02", param)
	if err == nil {
		dp.Time = t
		return nil
	}

	// Fallback to RFC3339 format for backward compatibility
	t, err = time.Parse(time.RFC3339, param)
	if err == nil {
		dp.Time = t
		return nil
	}

	return fmt.Errorf("unable to parse date '%s': expected format YYYY-MM-DD or RFC3339", param)
}

// Helper methods for DateParam to make it compatible with time operations

// Before returns true if dp is before the given time
func (dp DateParam) Before(t time.Time) bool {
	return dp.Time.Before(t)
}

// Equal returns true if dp equals the given time
func (dp DateParam) Equal(t time.Time) bool {
	return dp.Time.Equal(t)
}

// Add returns a new DateParam with the duration added
func (dp DateParam) Add(d time.Duration) DateParam {
	return DateParam{Time: dp.Time.Add(d)}
}

// AddDate returns a new DateParam with the date components added
func (dp DateParam) AddDate(years, months, days int) DateParam {
	return DateParam{Time: dp.Time.AddDate(years, months, days)}
}

// NewDateParam creates a DateParam from a time.Time
func NewDateParam(t time.Time) DateParam {
	return DateParam{Time: t}
}

// Granularity represents the time granularity for analytics data
type Granularity string

const (
	GranularityMinute Granularity = "minute"
	GranularityHour   Granularity = "hour"
	GranularityDay    Granularity = "day"
	GranularityWeek   Granularity = "week"
	GranularityMonth  Granularity = "month"
	GranularityYear   Granularity = "year"
)

// Dimension represents what dimension to group by
type Dimension string

const (
	DimensionDate   Dimension = "date"
	DimensionHour   Dimension = "hour"
	DimensionWeek   Dimension = "week"
	DimensionMonth  Dimension = "month"
	DimensionStatus Dimension = "status"
	DimensionType   Dimension = "type"
)

// Metric represents what we're measuring
type Metric string

const (
	MetricUsersCreated        Metric = "users_created"
	MetricWorkspacesCreated   Metric = "workspaces_created"
	MetricUsersActive         Metric = "users_active"
	MetricInvitationsSent     Metric = "invitations_sent"
	MetricRepositoriesCreated Metric = "repositories_created"
	MetricNewslettersCreated  Metric = "newsletters_created"
)

// AnalyticsParams contains parameters for analytics queries
type AnalyticsParams struct {
	FromDate    DateParam   `json:"fromDate" query:"fromDate"`
	ToDate      DateParam   `json:"toDate" query:"toDate"`
	Metric      Metric      `json:"metric" query:"metric"`
	Dimension   Dimension   `json:"dimension" query:"dimension"`
	Granularity Granularity `json:"granularity" query:"granularity"`
}

// Validate validates the analytics parameters
func (ap *AnalyticsParams) Validate() error {
	return validation.ValidateStruct(ap,
		validation.Field(&ap.FromDate, validation.Required),
		validation.Field(&ap.ToDate, validation.Required, validation.By(func(value interface{}) error {
			toDate := value.(DateParam)
			if toDate.Before(ap.FromDate.Time) {
				return fmt.Errorf("toDate must be after fromDate")
			}
			return nil
		})),
		validation.Field(&ap.Metric, validation.Required, validation.In(
			MetricUsersCreated,
			MetricWorkspacesCreated,
			MetricUsersActive,
			MetricInvitationsSent,
			MetricRepositoriesCreated,
			MetricNewslettersCreated,
		)),
		validation.Field(&ap.Dimension, validation.Required, validation.In(
			DimensionDate,
			DimensionHour,
			DimensionWeek,
			DimensionMonth,
			DimensionStatus,
			DimensionType,
		)),
		validation.Field(&ap.Granularity, validation.Required, validation.In(
			GranularityMinute,
			GranularityHour,
			GranularityDay,
			GranularityWeek,
			GranularityMonth,
			GranularityYear,
		)),
	)
}

// SetDefaults sets default values for the analytics parameters
func (ap *AnalyticsParams) SetDefaults() {
	if ap.ToDate.IsZero() {
		ap.ToDate = NewDateParam(time.Now())
	}
	if ap.FromDate.IsZero() {
		ap.FromDate = ap.ToDate.AddDate(0, 0, -30) // Default to 30 days back
	}
	if ap.Dimension == "" {
		ap.Dimension = DimensionDate
	}
	if ap.Granularity == "" {
		ap.Granularity = GranularityDay
	}
}

// GetDateFormat returns the SQL date format string based on granularity
func (ap *AnalyticsParams) GetDateFormat() string {
	switch ap.Granularity {
	case GranularityMinute:
		return "%Y-%m-%d %H:%i:00"
	case GranularityHour:
		return "%Y-%m-%d %H:00:00"
	case GranularityDay:
		return "%Y-%m-%d"
	case GranularityWeek:
		return "%Y-%u" // Year-Week
	case GranularityMonth:
		return "%Y-%m"
	case GranularityYear:
		return "%Y"
	default:
		return "%Y-%m-%d"
	}
}

// GetDateTruncFormat returns the date truncation format for PostgreSQL
func (ap *AnalyticsParams) GetDateTruncFormat() string {
	switch ap.Granularity {
	case GranularityMinute:
		return "minute"
	case GranularityHour:
		return "hour"
	case GranularityDay:
		return "day"
	case GranularityWeek:
		return "week"
	case GranularityMonth:
		return "month"
	case GranularityYear:
		return "year"
	default:
		return "day"
	}
}

// ToDateExclusive returns the day after ToDate for use in SQL queries with < operator,
// ensuring records created at any time on ToDate are included.
func (ap *AnalyticsParams) ToDateExclusive() time.Time {
	return ap.ToDate.AddDate(0, 0, 1).Time
}

// GenerateDateRange generates all dates between FromDate and ToDate based on granularity
func (ap *AnalyticsParams) GenerateDateRange() []time.Time {
	var dates []time.Time
	current := ap.FromDate

	for current.Before(ap.ToDate.Time) || current.Equal(ap.ToDate.Time) {
		dates = append(dates, current.Time)

		switch ap.Granularity {
		case GranularityMinute:
			current = current.Add(time.Minute)
		case GranularityHour:
			current = current.Add(time.Hour)
		case GranularityDay:
			current = current.AddDate(0, 0, 1)
		case GranularityWeek:
			current = current.AddDate(0, 0, 7)
		case GranularityMonth:
			current = current.AddDate(0, 1, 0)
		case GranularityYear:
			current = current.AddDate(1, 0, 0)
		default:
			current = current.AddDate(0, 0, 1)
		}
	}

	return dates
}

// FormatDateForResponse formats a date for the response based on granularity
func (ap *AnalyticsParams) FormatDateForResponse(date time.Time) string {
	switch ap.Granularity {
	case GranularityMinute:
		return date.Format("2006-01-02 15:04")
	case GranularityHour:
		return date.Format("2006-01-02 15:00")
	case GranularityDay:
		return date.Format("2006-01-02")
	case GranularityWeek:
		year, week := date.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case GranularityMonth:
		return date.Format("2006-01")
	case GranularityYear:
		return date.Format("2006")
	default:
		return date.Format("2006-01-02")
	}
}
