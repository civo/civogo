package civogo

import "fmt"

// Internal pagination defaults and safety caps used by the no-arg List*
// convenience functions. These exist so callers writing
// `c.ListIPs()` do not silently get a truncated result when the account
// owns more than the server-side default per_page.
//
// Values picked deliberately:
//   - paginationPerPage:  large enough that most accounts complete in one
//     round trip, small enough that the server is not asked to marshal a
//     huge response. 100 matches the upper bound used by api-go's existing
//     ReadPaginationOptions defaults in practice.
//   - paginationMaxPages: bounded so a runaway loop on a misbehaving server
//     cannot spin forever. With per_page=100 this caps a List call at
//     10,000 items.
//
// If the bound is reached, paginateAll returns an explicit error rather
// than silently truncating again — the whole point of this package.
const (
	paginationPerPage  = 100
	paginationMaxPages = 100
	paginationMaxItems = paginationPerPage * paginationMaxPages
)

// ErrPaginationCapExceeded is returned by paginateAll when the iteration
// would exceed paginationMaxPages or paginationMaxItems. Callers seeing
// this should switch to an explicit-page API (e.g. ListInstances) and
// process results in chunks.
type ErrPaginationCapExceeded struct {
	Resource   string
	PagesSeen  int
	ItemsSeen  int
	TotalPages int
}

func (e ErrPaginationCapExceeded) Error() string {
	return fmt.Sprintf(
		"civogo: pagination cap exceeded for %s (saw %d pages / %d items; server reports %d total pages; max %d pages / %d items)",
		e.Resource, e.PagesSeen, e.ItemsSeen, e.TotalPages, paginationMaxPages, paginationMaxItems,
	)
}

// pageFetcher fetches a single page of a paginated v2 list endpoint.
// It must return the items on that page, the server's reported total
// number of pages, and any error.
type pageFetcher[T any] func(page, perPage int) (items []T, totalPages int, err error)

// paginateAll iterates a paginated v2 list endpoint sequentially,
// starting at page 1, until all pages reported by the server have been
// consumed. It returns the merged Items slice in server order.
//
// Iteration is sequential (never parallel) because the server's total
// page count is not known until page 1 returns. A hard cap on pages and
// items guards against runaway responses.
func paginateAll[T any](resource string, fetch pageFetcher[T]) ([]T, error) {
	var all []T
	page := 1
	for {
		items, totalPages, err := fetch(page, paginationPerPage)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if len(all) > paginationMaxItems {
			return nil, ErrPaginationCapExceeded{Resource: resource, PagesSeen: page, ItemsSeen: len(all), TotalPages: totalPages}
		}
		if totalPages <= page {
			break
		}
		if page >= paginationMaxPages {
			return nil, ErrPaginationCapExceeded{Resource: resource, PagesSeen: page, ItemsSeen: len(all), TotalPages: totalPages}
		}
		page++
	}
	return all, nil
}
