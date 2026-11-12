package api

// PageInfo is the pagination metadata returned by every v3 list endpoint.
type PageInfo struct {
	Page            int `json:"page"`
	PerPage         int `json:"per_page"`
	ItemsCount      int `json:"items_count"`
	TotalItemsCount int `json:"total_items_count"`
	TotalPages      int `json:"total_pages"`
}

// PageParams selects a page of a list endpoint. The zero value asks for
// the server defaults.
type PageParams struct {
	Page    int
	PerPage int
}
