package api

import "iter"

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

const allPerPage = 100

// FetchPage fetches one page of a list endpoint, returning the items on
// that page plus the server's pagination metadata.
type FetchPage[T any] func(p PageParams) (items []T, info *PageInfo, err error)

// All lazily iterates every item of a paginated list endpoint, fetching
// pages as needed. It stops at the first error, yielded with a zero item.
func All[T any](fetch FetchPage[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for page := 1; ; page++ {
			items, info, err := fetch(PageParams{Page: page, PerPage: allPerPage})
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			if info == nil || page >= info.TotalPages {
				return
			}
		}
	}
}
