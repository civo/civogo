package civogo

import (
	"fmt"
	"strings"
)

// ListFirewalls implemented in a fake way for automated tests
func (c *FakeClient) ListFirewalls() ([]Firewall, error) {
	return c.Firewalls, nil
}

// FindFirewall implemented in a fake way for automated tests
func (c *FakeClient) FindFirewall(search string) (*Firewall, error) {
	for _, firewall := range c.Firewalls {
		if strings.Contains(firewall.Name, search) {
			return &firewall, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// NewFirewall implemented in a fake way for automated tests
func (c *FakeClient) NewFirewall(*FirewallConfig) (*FirewallResult, error) {
	firewall := Firewall{
		ID:   c.generateID(),
		Name: "fw-name",
	}
	c.Firewalls = append(c.Firewalls, firewall)

	return &FirewallResult{
		ID:     firewall.ID,
		Name:   firewall.Name,
		Result: "success",
	}, nil
}

// RenameFirewall implemented in a fake way for automated tests
func (c *FakeClient) RenameFirewall(id string, f *FirewallConfig) (*SimpleResponse, error) {
	for i, firewall := range c.Firewalls {
		if firewall.ID == id {
			c.Firewalls[i].Name = f.Name
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// DeleteFirewall implemented in a fake way for automated tests
func (c *FakeClient) DeleteFirewall(id string) (*SimpleResponse, error) {
	for i, firewall := range c.Firewalls {
		if firewall.ID == id {
			c.Firewalls[len(c.Firewalls)-1], c.Firewalls[i] = c.Firewalls[i], c.Firewalls[len(c.Firewalls)-1]
			c.Firewalls = c.Firewalls[:len(c.Firewalls)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}

// NewFirewallRule implemented in a fake way for automated tests
func (c *FakeClient) NewFirewallRule(r *FirewallRuleConfig) (*FirewallRule, error) {
	rule := FirewallRule{
		ID:        c.generateID(),
		Protocol:  r.Protocol,
		StartPort: r.StartPort,
		EndPort:   r.EndPort,
		Cidr:      r.Cidr,
		Label:     r.Label,
	}
	c.FirewallRules = append(c.FirewallRules, rule)
	return &rule, nil
}

// ListFirewallRules implemented in a fake way for automated tests
func (c *FakeClient) ListFirewallRules(id string) ([]FirewallRule, error) {
	return c.FirewallRules, nil
}

// FindFirewallRule implemented in a fake way for automated tests
func (c *FakeClient) FindFirewallRule(firewallID string, search string) (*FirewallRule, error) {
	for _, rule := range c.FirewallRules {
		if rule.FirewallID == firewallID && strings.Contains(rule.Label, search) {
			return &rule, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// DeleteFirewallRule implemented in a fake way for automated tests
func (c *FakeClient) DeleteFirewallRule(id string, ruleID string) (*SimpleResponse, error) {
	for i, rule := range c.FirewallRules {
		if rule.ID == ruleID {
			c.FirewallRules[len(c.FirewallRules)-1], c.FirewallRules[i] = c.FirewallRules[i], c.FirewallRules[len(c.FirewallRules)-1]
			c.FirewallRules = c.FirewallRules[:len(c.FirewallRules)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}
