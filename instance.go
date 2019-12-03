package civogo

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/civo/civogo/utils"
)

// Instance represents a virtual server within Civo's infrastructure
type Instance struct {
	ID                       string    `json:"id"`
	OpenstackServerID        string    `json:"openstack_server_id"`
	Hostname                 string    `json:"hostname"`
	ReverseDNS               string    `json:"reverse_dns"`
	Size                     string    `json:"size"`
	Region                   string    `json:"region"`
	NetworkID                string    `json:"network_id"`
	PrivateIP                string    `json:"private_ip"`
	PublicIP                 string    `json:"public_ip"`
	PseudoIP                 string    `json:"pseudo_ip"`
	TemplateID               string    `json:"template_id"`
	SnapshotID               string    `json:"snapshot_id"`
	InitialUser              string    `json:"initial_user"`
	InitialPassword          string    `json:"initial_password"`
	SSHKey                   string    `json:"ssh_key"`
	Status                   string    `json:"status"`
	Notes                    string    `json:"notes"`
	FirewallID               string    `json:"firewall_id"`
	Tags                     []string  `json:"tags"`
	CivostatsdToken          string    `json:"civostatsd_token"`
	CivostatsdStats          string    `json:"civostatsd_stats"`
	CivostatsdStatsPerMinute []string  `json:"civostatsd_stats_per_minute"`
	CivostatsdStatsPerHour   []string  `json:"civostatsd_stats_per_hour"`
	OpenstackImageID         string    `json:"openstack_image_id"`
	RescuePassword           string    `json:"rescue_password"`
	VolumeBacked             bool      `json:"volume_backed"`
	Script                   string    `json:"script"`
	CreatedAt                time.Time `json:"created_at"`
}

// PaginatedInstanceList returns a paginated list of Instance object
type PaginatedInstanceList struct {
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
	Pages   int        `json:"pages"`
	Items   []Instance `json:"items"`
}

// InstanceConfig describes the parameters for a new instance
// none of the fields are mandatory and will be automatically
// set with default values
type InstanceConfig struct {
	Count            int    `form:"count"`
	Hostname         string `form:"hostname"`
	ReverseDNS       string `form:"reverse_dns"`
	Size             string `form:"size"`
	Region           string `form:"region"`
	PublicIPRequired bool   `form:"public_ip_required"`
	NetworkID        string `form:"network_id"`
	TemplateID       string `form:"template_id"`
	SnapshotID       string `form:"snapshot_id"`
	InitialUser      string `form:"initial_user"`
	SSHKeyID         string `form:"ssh_key_id"`
	Script           string `form:"script"`
	Tags             string `form:"tags"`
}

// ListInstances returns a list of Instances owned by the calling API account
func (c *Client) ListInstances() (*PaginatedInstanceList, error) {
	resp, err := c.SendGetRequest("/v2/instances")
	if err != nil {
		return nil, err
	}

	PaginatedInstances := PaginatedInstanceList{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&PaginatedInstances)
	return &PaginatedInstances, err
}

// GetInstance returns a single Instance by its full ID
func (c *Client) GetInstance(id string) (*Instance, error) {
	resp, err := c.SendGetRequest("/v2/instances/" + id)
	if err != nil {
		return nil, err
	}

	instance := Instance{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&instance)
	return &instance, err
}

// NewInstanceConfig returns an initialized config for a new instance
func (c *Client) NewInstanceConfig() (*InstanceConfig, error) {
	var sshKeyID string
	sshKey, err := c.GetDefaultSSHKey()
	if err == nil {
		sshKeyID = sshKey.ID
	} else {
		sshKeyID = ""
	}

	network, err := c.GetDefaultNetwork()
	if err != nil {
		return nil, err
	}

	template, err := c.GetTemplateByCode("ubuntu-18.04")
	if err != nil {
		return nil, err
	}

	return &InstanceConfig{
		Count:            1,
		Hostname:         utils.RandomName(),
		ReverseDNS:       "",
		Size:             "g2.xsmall",
		Region:           "lon1",
		PublicIPRequired: true,
		NetworkID:        network.ID,
		TemplateID:       template.ID,
		SnapshotID:       "",
		InitialUser:      "civo",
		SSHKeyID:         sshKeyID,
		Script:           "",
		Tags:             "",
	}, nil
}

// CreateInstance creates a new instance in the account
func (c *Client) CreateInstance(config *InstanceConfig) (*Instance, error) {
	body, err := c.SendPostRequest("/v2/instances", config)
	if err != nil {
		return nil, err
	}

	var instance Instance
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&instance)

	return &instance, nil
}
