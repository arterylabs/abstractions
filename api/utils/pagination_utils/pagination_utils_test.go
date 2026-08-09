package pagination_utils_test

import (
	"testing"

	"abstractions/utils/pagination_utils"

	"github.com/stretchr/testify/assert"
)

func TestParseSortOrder(t *testing.T) {
	assert.Equal(t, pagination_utils.Asc, pagination_utils.ParseSortOrder("asc"))
	assert.Equal(t, pagination_utils.Desc, pagination_utils.ParseSortOrder("desc"))
	assert.Equal(t, pagination_utils.Desc, pagination_utils.ParseSortOrder("invalid"))
	assert.Equal(t, pagination_utils.Desc, pagination_utils.ParseSortOrder(""))
	assert.Equal(t, "asc", pagination_utils.Asc.String())
	assert.Equal(t, "desc", pagination_utils.Desc.String())
}

func TestPaginationParamsValidate(t *testing.T) {
	// Defaults when missing/invalid.
	p := pagination_utils.PaginationParams{Page: 0, PerPage: 0, SortOrder: "sideways"}
	p.Validate()
	assert.Equal(t, pagination_utils.DefaultPage, p.Page)
	assert.Equal(t, pagination_utils.DefaultPerPage, p.PerPage)
	assert.Equal(t, pagination_utils.DefaultSortOrder, p.SortOrder)

	// Clamp perPage to max.
	p2 := pagination_utils.PaginationParams{Page: 2, PerPage: 500, SortOrder: "asc"}
	p2.Validate()
	assert.Equal(t, 2, p2.Page)
	assert.Equal(t, pagination_utils.MaxPerPage, p2.PerPage)
	assert.Equal(t, "asc", p2.SortOrder)

	// Offset/limit helpers.
	p3 := pagination_utils.PaginationParams{Page: 3, PerPage: 25}
	p3.Validate()
	assert.Equal(t, 50, p3.GetOffset())
	assert.Equal(t, 25, p3.GetLimit())
}
