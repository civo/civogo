package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Firewall represents list of rule in Civo's infrastructure
type Firewall struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	RulesCount     string `json:"rules_count"`
	InstancesCount string `json:"instances_count"`
	Region         string `json:"region"`
}

type FirewallResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

type FirewallConfig struct {
	Name string `from:"name"`
}

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

type FirewallRuleConfig struct {
	FirewallID string   `from:"firewall_id"`
	Protocol   string   `from:"protocol"`
	StartPort  string   `from:"start_port"`
	EndPort    string   `from:"end_port"`
	Cidr       []string `from:"cidr"`
	Direction  string   `from:"direction"`
	Label      string   `json:"label,omitempty"`
}

// ListFirewall returns all firewall owned by the calling API account
func (c *Client) ListFirewall() ([]Firewall, error) {
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

// NewFirewall creates a new firewall record
func (c *Client) NewFirewall(r *FirewallConfig) (*FirewallResult, error) {
	body, err := c.SendPostRequest("/v2/firewalls/", r)
	if err != nil {
		return nil, err
	}

	result := &FirewallResult{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteFirewall deletes an firewall
func (c *Client) DeleteFirewall(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest("/v2/firewalls/" + id)
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// NewRecord creates a new DNS record
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

// ListFirewallRule get all rules for a firewall
func (c *Client) ListFirewallRule(id string) ([]FirewallRule, error) {
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

// DeleteFirewallRule deletes an firewall
func (c *Client) DeleteFirewallRule(id string, id_rule string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/firewalls/%s/rules/%s", id, id_rule))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}
