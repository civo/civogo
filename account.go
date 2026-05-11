package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PaginatedAccounts returns a paginated list of Account object
type PaginatedAccounts struct {
	Page    int       `json:"page"`
	PerPage int       `json:"per_page"`
	Pages   int       `json:"pages"`
	Items   []Account `json:"items"`
}

// ListAccounts lists all accounts
func (c *Client) ListAccounts() (*PaginatedAccounts, error) {
	items, err := paginateAll("accounts", func(page, perPage int) ([]Account, int, error) {
		resp, err := c.SendGetRequest(fmt.Sprintf("/v2/accounts?page=%d&per_page=%d", page, perPage))
		if err != nil {
			return nil, 0, decodeError(err)
		}
		pg := PaginatedAccounts{}
		if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&pg); err != nil {
			return nil, 0, decodeError(err)
		}
		return pg.Items, pg.Pages, nil
	})
	if err != nil {
		return nil, err
	}
	return &PaginatedAccounts{Page: 1, Pages: 1, PerPage: len(items), Items: items}, nil
}

// GetAccountID returns the account ID
func (c *Client) GetAccountID() string {
	accounts, err := c.ListAccounts()
	if err != nil {
		return ""
	}

	if len(accounts.Items) == 0 {
		return "No account found"
	}

	return accounts.Items[0].ID
}
