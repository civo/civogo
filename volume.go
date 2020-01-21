package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Volume struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	InstanceID string    `json:"instance_id"`
	MountPoint string    `json:"mountpoint"`
	SizeGB     int       `json:"size_gb"`
	Bootable   bool      `json:"bootable"`
	CreatedAt  time.Time `json:"created_at"`
}

type VolumeResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

type VolumeConfig struct {
	Name     string `form:"name"`
	SizeGB   int    `form:"size_gb"`
	Bootable bool   `form:"bootable"`
}

// ListVolumes returns all volumes owned by the calling API account
func (c *Client) ListVolumes() ([]Volume, error) {
	resp, err := c.SendGetRequest("/v2/volumes")
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var volumes = make([]Volume, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&volumes); err != nil {
		return nil, err
	}

	return volumes, nil
}

// NewVolume creates a new volume
func (c *Client) NewVolume(v *VolumeConfig) (*VolumeResult, error) {
	body, err := c.SendPostRequest("/v2/volumes/", v)
	if err != nil {
		return nil, err
	}

	var result = &VolumeResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// ResizeVolume resizes a volume
func (c *Client) ResizeVolume(id string, size int) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/resize", id), map[string]int{
		"size_gb": size,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// AttachVolume attaches a volume to an intance
func (c *Client) AttachVolume(id string, instance string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/attach", id), map[string]string{
		"instance_id": instance,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// DetachVolume attach volume from any instances
func (c *Client) DetachVolume(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/volumes/%s/detach", id), "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// DeleteVolume deletes a volumes
func (c *Client) DeleteVolume(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/volumes/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
