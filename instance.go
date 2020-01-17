package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// ListInstances returns a page of Instances owned by the calling API account
func (c *Client) ListInstances(page int, perPage int) (*PaginatedInstanceList, error) {
	url := "/v2/instances"
	if page != 0 && perPage != 0 {
		url = url + fmt.Sprintf("?page=%d&per_page=%d", page, perPage)
	}

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	PaginatedInstances := PaginatedInstanceList{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&PaginatedInstances)
	return &PaginatedInstances, err
}

// ListAllInstances returns all (well, upto 99,999,999 instances) Instances owned by the calling API account
func (c *Client) ListAllInstances() ([]Instance, error) {
	instances, err := c.ListInstances(1, 99999999)
	if err != nil {
		return []Instance{}, err
	}

	return instances.Items, nil
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

// SetInstanceTags sets the tags for the specified instance
func (c *Client) SetInstanceTags(id, tags string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest("/v2/instances/"+id+"/tags", map[string]string{
		"tags": tags,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// UpdateInstance updates an Instance's hostname, reverse DNS or notes
func (c *Client) UpdateInstance(i *Instance) (*SimpleResponse, error) {
	params := map[string]string{
		"hostname":    i.Hostname,
		"reverse_dns": i.ReverseDNS,
		"notes":       i.Notes,
	}

	if i.Notes == "" {
		params["notes_delete"] = "true"
	}

	resp, err := c.SendPutRequest("/v2/instances/"+i.ID, params)
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// DeleteInstance deletes an instance and frees its resources
func (c *Client) DeleteInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest("/v2/instances/" + id)
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// RebootInstance reboots an instance (short version of HardRebootInstance)
func (c *Client) RebootInstance(id string) (*SimpleResponse, error) {
	return c.HardRebootInstance(id)
}

// HardRebootInstance harshly reboots an instance (like shutting the power off and booting it again)
func (c *Client) HardRebootInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPostRequest("/v2/instances/"+id+"/hard_reboots", "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// SoftRebootInstance requests the VM to shut down nicely
func (c *Client) SoftRebootInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPostRequest("/v2/instances/"+id+"/soft_reboots", "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// StopInstance shuts the power down to the instance
func (c *Client) StopInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest("/v2/instances/"+id+"/stop", "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// StartInstance starts the instance booting from the shutdown state
func (c *Client) StartInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest("/v2/instances/"+id+"/start", "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// UpgradeInstance resizes the instance up to the new specification
// it's not possible to resize the instance to a smaller size
func (c *Client) UpgradeInstance(id, newSize string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest("/v2/instances/"+id+"/resize", map[string]string{
		"size": newSize,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// MovePublicIPToInstance moves a public IP to the specified instance
func (c *Client) MovePublicIPToInstance(id, ipAddress string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest("/v2/instances/"+id+"/ip/"+ipAddress, "")
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// SetInstanceFirewall changes the current firewall for an instance
func (c *Client) SetInstanceFirewall(id, firewallID string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest("/v2/instances/"+id+"/firewall", map[string]string{
		"firewall_id": firewallID,
	})
	if err != nil {
		return nil, err
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}
