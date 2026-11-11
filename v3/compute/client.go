// Code generated from compute.openapi3.json (proto commit 60c7e49248e86df73908df1dde1c7e92ce59bdfd) by gencmd. DO NOT EDIT.

// Package compute is the client for the Civo Compute v3 API.
package compute

import "github.com/civo/civogo/v3/api"

// Client groups the Compute services. Use civogo.New to build one.
type Client struct {
	Instances *InstancesService
}

// NewClient is called by the root civogo package.
func NewClient(t *api.Transport) *Client {
	return &Client{
		Instances: &InstancesService{t: t},
	}
}
