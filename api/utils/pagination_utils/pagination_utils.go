package pagination_utils

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

// Default pagination constants
const (
	DefaultPage      = 1
	DefaultPerPage   = 10
	MaxPerPage       = 100
	DefaultSortOrder = "desc"
)

func (s SortOrder) IsValid() bool {
	return s == Asc || s == Desc
}

func (s SortOrder) String() string {
	return string(s)
}

func ParseSortOrder(order string) SortOrder {
	switch order {
	case "asc":
		return Asc
	case "desc":
		return Desc
	default:
		return Desc // Default to Desc if invalid
	}
}

// PaginationParams defines the parameters for paginated requests
// It includes fields for pagination, sorting, and searching.
type PaginationParams struct {
	Page      int    `query:"page"`      // Current page number
	PerPage   int    `query:"perPage"`   // Number of items per page
	Sort      string `query:"sort"`      // Column to order by (name, email, created_at)
	SortOrder string `query:"sortOrder"` // Order direction (asc, desc)
	Query     string `query:"query"`     // Search query for name or email
}

func (p *PaginationParams) Validate() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PerPage < 1 {
		p.PerPage = DefaultPerPage
	}
	if p.PerPage > MaxPerPage {
		p.PerPage = MaxPerPage
	}

	if p.SortOrder != "asc" && p.SortOrder != "desc" {
		p.SortOrder = DefaultSortOrder
	}
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

func (p *PaginationParams) GetLimit() int {
	return p.PerPage
}
