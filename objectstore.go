package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ObjectStore is the struct for the ObjectStorage model
type ObjectStorage struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	GeneratedName string `json:"generated_name"`
	//default:1000
	MaxObjects int `json:"max_objects"`
	//default:500G
	MaxSize         string `json:"max_size"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	BucketURL       string `json:"bucket_URL"`
	//Status can be one of - 1. Ready, 2.Creating and 3. failed
	Status string `json:"status"`
}

// CreateObjectStorageRequest holds the request to create a new object storage
type CreateObjectStorageRequest struct {
	Name       string `json:"name" validate:"required"`
	MaxSizeGB  int    `json:"max_size_gb" validate:"required"`
	MaxObjects int    `json:"max_objects"`
}

// UpdateObjectStorageRequest holds the request to update a specified object storage's details
type UpdateObjectStorageRequest struct {
	MaxSizeGB  int `json:"max_size_gb"`
	MaxObjects int `json:"max_objects"`
}

// ListObjectStores returns all objectstores in that specific region
func (c *Client) ListObjectStores() ([]ObjectStorage, error) {
	resp, err := c.SendGetRequest("/v2/objectstorage/buckets")
	if err != nil {
		return nil, decodeError(err)
	}

	var objectstores = make([]ObjectStorage, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&objectstores); err != nil {
		return nil, err
	}

	return objectstores, nil
}

// GetObjectStore finds an objectstore by the full ID
func (c *Client) GetObjectStore(id string) (*ObjectStorage, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/objectstorage/buckets/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	var os = ObjectStorage{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&os); err != nil {
		return nil, err
	}

	return &os, nil
}

// FindObjectStore finds an objectstore by name or by accesskeyID
func (c *Client) FindObjectStore(search string) (*ObjectStorage, error) {
	objectstores, err := c.ListObjectStores()
	if err != nil {
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := ObjectStorage{}

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
func (c *Client) NewObjectStore(v *CreateObjectStorageRequest) (*ObjectStorage, error) {
	body, err := c.SendPostRequest("/v2/objectstorage/buckets", v)
	if err != nil {
		return nil, decodeError(err)
	}

	var result = &ObjectStorage{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateObjectStore updates an objectstore
func (c *Client) UpdateObjectStore(id string, v *UpdateObjectStorageRequest) (*ObjectStorage, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/objectstorage/buckets/%s", id), v)
	if err != nil {
		return nil, decodeError(err)
	}

	var result = &ObjectStorage{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteObjectStore deletes an objectstore
func (c *Client) DeleteObjectStore(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/objectstorage/buckets/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}
