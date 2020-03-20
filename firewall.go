package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Firewall represents list of rule in Civo's infrastructure
type Firewall struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	RulesCount     int    `json:"rules_count"`
	InstancesCount int    `json:"instances_count"`
	Region         string `json:"region"`
}

// FirewallResult is the response from the Civo Firewall APIs
type FirewallResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

// FirewallRule represents a single rule for a given firewall, regarding
// which ports to open and which protocol, to which CIDR
type FirewallRule struct {
	ID         string   `json:"id"`
	FirewallID string   `json:"firewall_id"`
	Protocol   string   `json:"protocol"`
	StartPort  string   `json:"start_port"`
	EndPort    string   `json:"end_port"`
	Cidr       []string `json:"cidr"`
	Direction  string   `json:"direction"`
	Label      string   `json:"label,omitempty"`
}

// FirewallRuleConfig is how you specify the details when creating a new rule
type FirewallRuleConfig struct {
	FirewallID string   `json:"firewall_id"`
	Protocol   string   `json:"protocol"`
	StartPort  string   `json:"start_port"`
	EndPort    string   `json:"end_port"`
	Cidr       []string `json:"cidr"`
	Direction  string   `json:"direction"`
	Label      string   `json:"label,omitempty"`
}

type firewallConfig struct {
	Name string `json:"name"`
}

// ListFirewalls returns all firewall owned by the calling API account
func (c *Client) ListFirewalls() ([]Firewall, error) {
	resp, err := c.SendGetRequest("/v2/firewalls")
	if err != nil {
		return nil, err
	}

	firewall := make([]Firewall, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&firewall); err != nil {
		return nil, err
	}

	return firewall, nil
}

// FindFirewall finds a firewall by either part of the ID or part of the name
func (c *Client) FindFirewall(search string) (*Firewall, error) {
	firewalls, err := c.ListFirewalls()
	if err != nil {
		return nil, err
	}

	found := -1

	for i, firewall := range firewalls {
		if strings.Contains(firewall.ID, search) || strings.Contains(firewall.Name, search) {
			if found != -1 {
				return nil, fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
			found = i
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("unable to find %s, zero matches", search)
	}

	return &firewalls[found], nil
}

// NewFirewall creates a new firewall record
func (c *Client) NewFirewall(name string) (*FirewallResult, error) {
	fw := firewallConfig{Name: name}
	body, err := c.SendPostRequest("/v2/firewalls", fw)
	if err != nil {
		return nil, err
	}

	result := &FirewallResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// RenameFirewall rename firewall
func (c *Client) RenameFirewall(id string, name string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/firewalls/%s", id), map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// DeleteFirewall deletes an firewall
func (c *Client) DeleteFirewall(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest("/v2/firewalls/" + id)
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// NewFirewallRule creates a new rule within a firewall
func (c *Client) NewFirewallRule(r *FirewallRuleConfig) (*FirewallRule, error) {
	if len(r.FirewallID) == 0 {
		return nil, fmt.Errorf("the firewall ID is empty")
	}

	resp, err := c.SendPostRequest(fmt.Sprintf("/v2/firewalls/%s/rules", r.FirewallID), r)
	if err != nil {
		return nil, err
	}

	rule := &FirewallRule{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(rule); err != nil {
		return nil, err
	}

	return rule, nil
}

// ListFirewallRules get all rules for a firewall
func (c *Client) ListFirewallRules(id string) ([]FirewallRule, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/firewalls/%s/rules", id))
	if err != nil {
		return nil, err
	}

	firewallRule := make([]FirewallRule, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&firewallRule); err != nil {
		return nil, err
	}

	return firewallRule, nil
}

// FindFirewallRule finds a firewall Rule by ID or part of the same
func (c *Client) FindFirewallRule(firewallID string, search string) (*FirewallRule, error) {
	firewallsRules, err := c.ListFirewallRules(firewallID)
	if err != nil {
		return nil, err
	}

	found := -1

	for i, firewallRule := range firewallsRules {
		if strings.Contains(firewallRule.ID, search) {
			if found != -1 {
				return nil, fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
			found = i
		}
	}

	if found == -1 {
		return nil, fmt.Errorf("unable to find %s, zero matches", search)
	}

	return &firewallsRules[found], nil
}

// DeleteFirewallRule deletes an firewall
func (c *Client) DeleteFirewallRule(id string, ruleID string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/firewalls/%s/rules/%s", id, ruleID))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
