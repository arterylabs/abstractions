package requests

import "abstractions/utils/pagination_utils"

type PaginationParams struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	OrderBy string `query:"order_by"`
	Search  string `query:"search"`
}

func (p *PaginationParams) Validate() error {
	if p.Page < 1 {
		p.Page = pagination_utils.DefaultPage
	}
	if p.PerPage < 1 {
		p.PerPage = pagination_utils.DefaultPerPage
	}
	if p.PerPage > pagination_utils.MaxPerPage {
		p.PerPage = pagination_utils.MaxPerPage
	}
	return nil
}
