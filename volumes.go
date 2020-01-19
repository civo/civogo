package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Volumes struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	InstanceID string    `json:"instance_id"`
	MountPoint string    `json:"mountpoint"`
	SizeGB     int       `json:"size_gb"`
	Bootable   bool      `json:"bootable"`
	CreatedAt  time.Time `json:"created_at"`
}

type VolumesResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

type VolumesConfig struct {
	Name     string `json:"name"`
	SizeGB   int    `json:"size_gb"`
	Bootable bool   `json:"bootable"`
}

// ListVolumes returns all volumes owned by the calling API account
func (c *Client) ListVolumes() ([]Volumes, error) {
	resp, err := c.SendGetRequest("/v2/volumes")
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var volumes = make([]Volumes, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&volumes); err != nil {
		return nil, err
	}

	return volumes, nil
}

// NewVolumes creates a new volume
func (c *Client) NewVolumes(r *VolumesConfig) (*VolumesResult, error) {
	body, err := c.SendPostRequest("/v2/volumes/", r)
	if err != nil {
		return nil, err
	}

	var result = &VolumesResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// ResizeVolumes resize a volume
func (c *Client) ResizeVolumes(id string, size int) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/resize", id), map[string]int{
		"size_gb": size,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// AttachVolumes attach volume to a intances
func (c *Client) AttachVolumes(id string, instance string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/attach", id), map[string]string{
		"instance_id": instance,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// DetachVolumes attach volume to a intances
func (c *Client) DetachVolumes(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/detach", id), "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// DeleteVolumes deletes an volumes
func (c *Client) DeleteVolumes(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/volumes/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
