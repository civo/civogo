package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ObjectStore is the struct for the ObjectStore model
type ObjectStore struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	GeneratedName string `json:"generated_name"`
	//default:1000
	MaxObjects int `json:"max_objects"`
	//default:500G
	MaxSize             string `json:"max_size"`
	AccessKeyID         string `json:"access_key_id"`
	SecretAccessKey     string `json:"secret_access_key"`
	ObjectStoreEndpoint string `json:"objectstore_endpoint"`
	//Status can be one of - 1.ready, 2.creating and 3.failed
	Status string `json:"status"`
}

// CreateObjectStoreRequest holds the request to create a new object storage
type CreateObjectStoreRequest struct {
	//Name            string `json:"name,omitempty" schema:"name"`
	Name            string `json:"-"`
	MaxSizeGB       int    `json:"max_size_gb" validate:"required"`
	MaxObjects      int    `json:"max_objects"`
	Prefix          string `json:"prefix,omitempty"`
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	Region          string `json:"region"`
}

// UpdateObjectStoreRequest holds the request to update a specified object storage's details
type UpdateObjectStoreRequest struct {
	MaxSizeGB       int    `json:"max_size_gb"`
	MaxObjects      int    `json:"max_objects"`
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	Region          string `json:"region"`
}

// ListObjectStores returns all objectstores in that specific region
func (c *Client) ListObjectStores() ([]ObjectStore, error) {
	resp, err := c.SendGetRequest("/v2/objectstores")
	if err != nil {
		return nil, decodeError(err)
	}

	var stores = make([]ObjectStore, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&stores); err != nil {
		return nil, err
	}

	return stores, nil
}

// GetObjectStore finds an objectstore by the full ID
func (c *Client) GetObjectStore(id string) (*ObjectStore, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/objectstores/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	var os = ObjectStore{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&os); err != nil {
		return nil, err
	}

	return &os, nil
}

// FindObjectStore finds an objectstore by name or by accesskeyID
func (c *Client) FindObjectStore(search string) (*ObjectStore, error) {
	objectstores, err := c.ListObjectStores()
	if err != nil {
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := ObjectStore{}

	for _, value := range objectstores {
		if value.AccessKeyID == search || value.Name == search || value.ID == search {
			exactMatch = true
			result = value
		} else if strings.Contains(value.AccessKeyID, search) || strings.Contains(value.Name, search) || strings.Contains(value.ID, search) {
			if !exactMatch {
				result = value
				partialMatchesCount++
			}
		}
	}

	if exactMatch || partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}

// NewObjectStore creates a new objectstore
func (c *Client) NewObjectStore(v *CreateObjectStoreRequest) (*ObjectStore, error) {
	body, err := c.SendPostRequest("/v2/objectstores", v)
	if err != nil {
		return nil, decodeError(err)
	}

	var result = &ObjectStore{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateObjectStore updates an objectstore
func (c *Client) UpdateObjectStore(id string, v *UpdateObjectStoreRequest) (*ObjectStore, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/objectstores/%s", id), v)
	if err != nil {
		return nil, decodeError(err)
	}

	var result = &ObjectStore{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteObjectStore deletes an objectstore
func (c *Client) DeleteObjectStore(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/objectstores/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}
