package repository

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int `json:"page" form:"page" binding:"omitempty,min=1"`
	Limit    int `json:"limit" form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortOrder string `json:"sort_order" form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// GetOffset calculates the offset for pagination
func (p *PaginationParams) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit returns the limit with a default value
func (p *PaginationParams) GetLimit() int {
	if p.Limit <= 0 {
		return 10 // Default limit
	}
	if p.Limit > 100 {
		return 100 // Maximum limit
	}
	return p.Limit
}

// GetSortOrder returns the sort order with a default value
func (p *PaginationParams) GetSortOrder() string {
	if p.SortOrder == "" {
		return "desc"
	}
	return p.SortOrder
}

// PaginationResult represents paginated results
type PaginationResult struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int `json:"total_pages"`
}

// CalculateTotalPages calculates the total pages based on total records and per page
func (p *PaginationResult) CalculateTotalPages() {
	if p.PerPage > 0 {
		p.TotalPages = int((p.Total + int64(p.PerPage) - 1) / int64(p.PerPage))
	}
}

// TodoFilters represents filters for todo queries
type TodoFilters struct {
	Search     string `json:"search" form:"search"`
	Completed  *bool  `json:"completed" form:"completed"`
	CategoryID *uint  `json:"category_id" form:"category_id"`
	Priority   string `json:"priority" form:"priority" binding:"omitempty,oneof=low medium high"`
}

// CategoryFilters represents filters for category queries
type CategoryFilters struct {
	Search string `json:"search" form:"search"`
}