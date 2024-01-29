package civogo

import (
	"fmt"
	"strings"
)

// ListLoadBalancers implemented in a fake way for automated tests
func (c *FakeClient) ListLoadBalancers() ([]LoadBalancer, error) {
	return c.LoadBalancers, nil
}

// GetLoadBalancer implemented in a fake way for automated tests
func (c *FakeClient) GetLoadBalancer(id string) (*LoadBalancer, error) {
	for _, lb := range c.LoadBalancers {
		if lb.ID == id {
			return &lb, nil
		}
	}

	err := fmt.Errorf("unable to get load balancer %s", id)
	return nil, DatabaseLoadBalancerNotFoundError.wrap(err)
}

// FindLoadBalancer implemented in a fake way for automated tests
func (c *FakeClient) FindLoadBalancer(search string) (*LoadBalancer, error) {
	exactMatch := false
	partialMatchesCount := 0
	result := LoadBalancer{}

	for _, lb := range c.LoadBalancers {
		if lb.ID == search || lb.Name == search {
			exactMatch = true
			result = lb
		} else if strings.Contains(lb.Name, search) || strings.Contains(lb.ID, search) {
			if !exactMatch {
				result = lb
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

// CreateLoadBalancer implemented in a fake way for automated tests
func (c *FakeClient) CreateLoadBalancer(r *LoadBalancerConfig) (*LoadBalancer, error) {
	loadbalancer := LoadBalancer{
		ID:                           c.generateID(),
		Name:                         r.Name,
		Algorithm:                    r.Algorithm,
		ExternalTrafficPolicy:        r.ExternalTrafficPolicy,
		SessionAffinityConfigTimeout: r.SessionAffinityConfigTimeout,
		SessionAffinity:              r.SessionAffinity,
		EnableProxyProtocol:          r.EnableProxyProtocol,
		FirewallID:                   r.FirewallID,
		ClusterID:                    r.ClusterID,
	}

	if r.Algorithm == "" {
		loadbalancer.Algorithm = "round_robin"
	}
	if r.FirewallID == "" {
		loadbalancer.FirewallID = c.generateID()
	}
	if r.ExternalTrafficPolicy == "" {
		loadbalancer.ExternalTrafficPolicy = "Cluster"
	}

	backends := make([]LoadBalancerBackend, 0)
	for _, b := range r.Backends {
		backend := LoadBalancerBackend{
			IP:         b.IP,
			Protocol:   b.Protocol,
			SourcePort: b.SourcePort,
			TargetPort: b.TargetPort,
		}
		backends = append(backends, backend)
	}
	loadbalancer.Backends = backends
	loadbalancer.PublicIP = c.generatePublicIP()
	loadbalancer.State = "available"

	c.LoadBalancers = append(c.LoadBalancers, loadbalancer)
	return &loadbalancer, nil
}

// UpdateLoadBalancer implemented in a fake way for automated tests
func (c *FakeClient) UpdateLoadBalancer(id string, r *LoadBalancerUpdateConfig) (*LoadBalancer, error) {
	for _, lb := range c.LoadBalancers {
		if lb.ID == id {
			lb.Name = r.Name
			lb.Algorithm = r.Algorithm
			lb.EnableProxyProtocol = r.EnableProxyProtocol
			lb.ExternalTrafficPolicy = r.ExternalTrafficPolicy
			lb.SessionAffinity = r.SessionAffinity
			lb.SessionAffinityConfigTimeout = r.SessionAffinityConfigTimeout

			backends := make([]LoadBalancerBackend, len(r.Backends))
			for i, b := range r.Backends {
				backends[i].IP = b.IP
				backends[i].Protocol = b.Protocol
				backends[i].SourcePort = b.SourcePort
				backends[i].TargetPort = b.TargetPort
			}

			if r.ExternalTrafficPolicy == "" {
				lb.ExternalTrafficPolicy = "Cluster"
			}

			return &lb, nil
		}
	}

	err := fmt.Errorf("unable to find load balancer %s", id)
	return nil, DatabaseLoadBalancerNotFoundError.wrap(err)
}

// DeleteLoadBalancer implemented in a fake way for automated tests
func (c *FakeClient) DeleteLoadBalancer(id string) (*SimpleResponse, error) {
	for i, lb := range c.LoadBalancers {
		if lb.ID == id {
			c.LoadBalancers[len(c.LoadBalancers)-1], c.LoadBalancers[i] = c.LoadBalancers[i], c.LoadBalancers[len(c.LoadBalancers)-1]
			c.LoadBalancers = c.LoadBalancers[:len(c.LoadBalancers)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}
