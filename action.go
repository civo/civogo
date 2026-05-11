package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-querystring/query"
)

// PaginateActionList is a struct for a page of actions
type PaginateActionList struct {
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Pages   int      `json:"pages"`
	Items   []Action `json:"items"`
}

// Action is a struct for an individual action within the database and when serialized
type Action struct {
	ID          int       `json:"id" gorm:"autoIncrement"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	AccountID   string    `json:"account_id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Details     string    `json:"details,omitempty"`
	RelatedID   string    `json:"related_id,omitempty"`
	RelatedType string    `json:"related_type,omitempty"`
	Debug       bool      `json:"debug"`
}

// ActionListRequest is a struct for the request to list actions
type ActionListRequest struct {
	PerPage      int    `json:"per_page,omitempty" url:"per_page,omitempty"`
	Page         int    `json:"page,omitempty" url:"page,omitempty"`
	IncludeDebug bool   `json:"include_debug,omitempty" url:"include_debug,omitempty"`
	ResourceID   string `json:"resource_id,omitempty" url:"resource_id,omitempty"`
	Details      string `json:"details,omitempty" url:"details,omitempty"`
	RelatedID    string `json:"related_id,omitempty" url:"related_id,omitempty"`
	ResourceType string `json:"resource_type,omitempty" url:"resource_type,omitempty"`
	ActionType   string `json:"action_type,omitempty" url:"action_type,omitempty"`
	CreatedAt    string `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	UserID       string `json:"user_id,omitempty" url:"user_id,omitempty"`
}

// ListActions returns a page of actions. Callers that want to enumerate
// every action matching the supplied filters should use ListAllActions
// instead — without explicit Page/PerPage values this function returns
// only the server-side first page (default 20 items).
func (c *Client) ListActions(listRequest *ActionListRequest) (*PaginateActionList, error) {
	url := "/v2/actions"

	vals, err := query.Values(listRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.SendGetRequest(fmt.Sprintf("%s?%s", url, vals.Encode()))
	if err != nil {
		return nil, decodeError(err)
	}

	paginateActionList := PaginateActionList{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&paginateActionList)
	return &paginateActionList, err
}

// ListAllActions returns every action matching the supplied filters by
// iterating server-side pagination internally. Any Page/PerPage values
// set on filters are ignored (the SDK chooses these to maximise iteration
// safety; see paginateAll in pagination.go).
//
// Pass nil for filters to enumerate all actions on the account with no
// filter criteria.
func (c *Client) ListAllActions(filters *ActionListRequest) ([]Action, error) {
	return paginateAll("actions", func(page, perPage int) ([]Action, int, error) {
		req := ActionListRequest{}
		if filters != nil {
			req = *filters
		}
		req.Page = page
		req.PerPage = perPage

		vals, err := query.Values(&req)
		if err != nil {
			return nil, 0, err
		}
		resp, err := c.SendGetRequest(fmt.Sprintf("/v2/actions?%s", vals.Encode()))
		if err != nil {
			return nil, 0, decodeError(err)
		}
		pg := PaginateActionList{}
		if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&pg); err != nil {
			return nil, 0, err
		}
		return pg.Items, pg.Pages, nil
	})
}
