package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type LoadBalancerBackend struct {
	InstanceID string `json:"instance_id"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
}

type LoadBalancerBackendConfig struct {
	InstanceID string `from:"instance_id"`
	Protocol   string `from:"protocol"`
	Port       int    `from:"port"`
}

type LoadBalancer struct {
	ID                      string `json:"id"`
	DefaultHostname         bool   `json:"default_hostname"`
	Hostname                string `json:"hostname"`
	Protocol                string `json:"protocol"`
	Port                    string `json:"port"`
	MaxRequestSize          int    `json:"max_request_size"`
	TlsCertificate          string `json:"tls_certificate"`
	TlsKey                  string `json:"tls_key"`
	Policy                  string `json:"policy"`
	HealthCheckPath         string `json:"health_check_path"`
	FailTimeout             int    `json:"fail_timeout"`
	MaxConns                int    `json:"max_conns"`
	IgnoreInvalidBackendTls bool   `json:"ignore_invalid_backend_tls"`
	Backends                []LoadBalancerBackend
}

type LoadBalancerConfig struct {
	Hostname                string `from:"hostname"`
	Protocol                string `from:"protocol"`
	TlsCertificate          string `from:"tls_certificate"`
	TlsKey                  string `from:"tls_key"`
	Policy                  string `json:"policy"`
	Port                    string `from:"port"`
	MaxRequestSize          int    `from:"max_request_size"`
	HealthCheckPath         string `from:"health_check_path"`
	FailTimeout             int    `from:"fail_timeout"`
	MaxConns                int    `from:"max_conns"`
	IgnoreInvalidBackendTls bool   `from:"ignore_invalid_backend_tls"`
	Backends                []LoadBalancerBackendConfig
}

// ListLoadBalance returns all load balance owned by the calling API account
func (c *Client) ListLoadBalancer() ([]LoadBalancer, error) {
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

// NewLoadBalance creates a new load balance
func (c *Client) NewLoadBalancer(r *LoadBalancerConfig) (*LoadBalancer, error) {
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

// UpdateLoadBalance update a load balance
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

// DeleteLoadBalance deletes a load balance
func (c *Client) DeleteLoadBalancer(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/loadbalancers/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
