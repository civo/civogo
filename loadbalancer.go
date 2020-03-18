package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// LoadBalancerBackend represents a backend instance being load-balanced
type LoadBalancerBackend struct {
	InstanceID string `json:"instance_id"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
}

// LoadBalancerBackendConfig is the configuration for creating backends
type LoadBalancerBackendConfig struct {
	InstanceID string `json:"instance_id"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
}

// LoadBalancer represents a load balancer configuration within Civo
type LoadBalancer struct {
	ID                      string `json:"id"`
	DefaultHostname         bool   `json:"default_hostname"`
	Hostname                string `json:"hostname"`
	Protocol                string `json:"protocol"`
	Port                    int    `json:"port"`
	MaxRequestSize          int    `json:"max_request_size"`
	TLSCertificate          string `json:"tls_certificate"`
	TLSKey                  string `json:"tls_key"`
	Policy                  string `json:"policy"`
	HealthCheckPath         string `json:"health_check_path"`
	FailTimeout             int    `json:"fail_timeout"`
	MaxConns                int    `json:"max_conns"`
	IgnoreInvalidBackendTLS bool   `json:"ignore_invalid_backend_tls"`
	Backends                []LoadBalancerBackend
}

// LoadBalancerConfig represents a load balancer to be created
type LoadBalancerConfig struct {
	Hostname                string                      `json:"hostname"`
	Protocol                string                      `json:"protocol"`
	TLSCertificate          string                      `json:"tls_certificate"`
	TLSKey                  string                      `json:"tls_key"`
	Policy                  string                      `json:"policy"`
	Port                    int                         `json:"port"`
	MaxRequestSize          int                         `json:"max_request_size"`
	HealthCheckPath         string                      `json:"health_check_path"`
	FailTimeout             int                         `json:"fail_timeout"`
	MaxConns                int                         `json:"max_conns"`
	IgnoreInvalidBackendTLS bool                        `json:"ignore_invalid_backend_tls"`
	Backends                []LoadBalancerBackendConfig `json:"backends"`
}

// ListLoadBalancers returns all load balancers owned by the calling API account
func (c *Client) ListLoadBalancers() ([]LoadBalancer, error) {
	resp, err := c.SendGetRequest("/v2/loadbalancers")
	if err != nil {
		return nil, err
	}

	loadbalancer := make([]LoadBalancer, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&loadbalancer); err != nil {
		return nil, err
	}

	return loadbalancer, nil
}

// FindLoadBalancer finds a load balancer by either part of the ID or part of the name
func (c *Client) FindLoadBalancer(search string) (*LoadBalancer, error) {
	lbs, err := c.ListLoadBalancers()
	if err != nil {
		return nil, err
	}

	found := -1

	for i, lb := range lbs {
		if strings.Contains(lb.ID, search) || strings.Contains(lb.Hostname, search) {
			if found != -1 {
				return nil, fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
			found = i
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("unable to find %s, zero matches", search)
	}

	return &lbs[found], nil
}

// CreateLoadBalancer creates a new load balancer
func (c *Client) CreateLoadBalancer(r *LoadBalancerConfig) (*LoadBalancer, error) {
	body, err := c.SendPostRequest("/v2/loadbalancers", r)
	if err != nil {
		return nil, err
	}

	loadbalancer := &LoadBalancer{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(loadbalancer); err != nil {
		return nil, err
	}

	return loadbalancer, nil
}

// UpdateLoadBalancer updates a load balancer
func (c *Client) UpdateLoadBalancer(id string, r *LoadBalancerConfig) (*LoadBalancer, error) {
	body, err := c.SendPutRequest(fmt.Sprintf("/v2/loadbalancers/%s", id), r)
	if err != nil {
		return nil, err
	}

	loadbalancer := &LoadBalancer{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(loadbalancer); err != nil {
		return nil, err
	}

	return loadbalancer, nil
}

// DeleteLoadBalancer deletes a load balancer
func (c *Client) DeleteLoadBalancer(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/loadbalancers/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
