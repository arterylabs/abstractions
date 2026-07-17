package responses

import "abstractions/utils/pagination_utils"

type PageInfo struct {
	Page         int  `json:"page"`
	PerPage      int  `json:"perPage"`
	TotalPages   int  `json:"totalPages"`
	HasNextPage  bool `json:"hasNextPage"`
	TotalRecords int  `json:"totalRecords"`
	PageRecords  int  `json:"pageRecords"`
}

type PaginatedResponse struct {
	Results  interface{} `json:"results"`
	PageInfo PageInfo    `json:"pageInfo"`
}

func GetPageInfo(page, perPage, resultsCount, totalRecords int) PageInfo {
	// Validate and set defaults to prevent division by zero
	if perPage <= 0 {
		perPage = pagination_utils.DefaultPerPage // Default to 10 items per page
	}
	if page <= 0 {
		page = pagination_utils.DefaultPage // Default to first page
	}

	totalPages := (totalRecords + perPage - 1) / perPage
	hasNextPage := page < totalPages
	return PageInfo{
		Page:         page,
		PerPage:      perPage,
		TotalPages:   totalPages,
		HasNextPage:  hasNextPage,
		TotalRecords: totalRecords,
		PageRecords:  resultsCount,
	}
}
