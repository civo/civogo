package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/civo/civogo/utils"
)

const (
	// DefaultInstanceUser is the default username used in newly created instances.
	DefaultInstanceUser string = "civo"
)

// Instance represents a virtual server within Civo's infrastructure
type Instance struct {
	ID                       string           `json:"id,omitempty"`
	OpenstackServerID        string           `json:"openstack_server_id,omitempty"`
	Hostname                 string           `json:"hostname,omitempty"`
	ReverseDNS               string           `json:"reverse_dns,omitempty"`
	Size                     string           `json:"size,omitempty"`
	Region                   string           `json:"region,omitempty"`
	NetworkID                string           `json:"network_id,omitempty"`
	PrivateIP                string           `json:"private_ip,omitempty"`
	PublicIP                 string           `json:"public_ip,omitempty"`
	IPv6                     string           `json:"ipv6,omitempty"`
	PseudoIP                 string           `json:"pseudo_ip,omitempty"`
	TemplateID               string           `json:"template_id,omitempty"`
	SourceType               string           `json:"source_type,omitempty"`
	SourceID                 string           `json:"source_id,omitempty"`
	SnapshotID               string           `json:"snapshot_id,omitempty"`
	InitialUser              string           `json:"initial_user,omitempty"`
	InitialPassword          string           `json:"initial_password,omitempty"`
	SSHKey                   string           `json:"ssh_key,omitempty"`
	SSHKeyID                 string           `json:"ssh_key_id,omitempty"`
	Status                   string           `json:"status,omitempty"`
	Notes                    string           `json:"notes,omitempty"`
	FirewallID               string           `json:"firewall_id,omitempty"`
	Tags                     []string         `json:"tags,omitempty"`
	CivostatsdToken          string           `json:"civostatsd_token,omitempty"`
	CivostatsdStats          string           `json:"civostatsd_stats,omitempty"`
	CivostatsdStatsPerMinute []string         `json:"civostatsd_stats_per_minute,omitempty"`
	CivostatsdStatsPerHour   []string         `json:"civostatsd_stats_per_hour,omitempty"`
	OpenstackImageID         string           `json:"openstack_image_id,omitempty"`
	RescuePassword           string           `json:"rescue_password,omitempty"`
	VolumeBacked             bool             `json:"volume_backed,omitempty"`
	CPUCores                 int              `json:"cpu_cores,omitempty"`
	RAMMegabytes             int              `json:"ram_mb,omitempty"`
	DiskGigabytes            int              `json:"disk_gb,omitempty"`
	GPUCount                 int              `json:"gpu_count,omitempty"`
	GPUType                  string           `json:"gpu_type,omitempty"`
	Script                   string           `json:"script,omitempty"`
	CreatedAt                time.Time        `json:"created_at,omitempty"`
	ReservedIPID             string           `json:"reserved_ip_id,omitempty"`
	ReservedIPName           string           `json:"reserved_ip_name,omitempty"`
	ReservedIP               string           `json:"reserved_ip,omitempty"`
	VolumeType               string           `json:"volume_type,omitempty"`
	Subnets                  []Subnet         `json:"subnets,omitempty"`
	AttachedVolumes          []AttachedVolume `json:"attached_volumes,omitempty"`
	PlacementRule            PlacementRule    `json:"placement_rule,omitempty"`
	NetworkBandwidthLimit    int              `json:"network_bandwidth_limit,omitempty"`
	AllowedIPs               []string         `json:"allowed_ips,omitempty"`
}

// InstanceVnc represents VNC information for an instances
type InstanceVnc struct {
	URI        string `json:"uri,omitempty"`
	Expiration string `json:"expiration,omitempty"`
}

// CreateInstanceVncResp represents VNC information for a new instance console
type CreateInstanceVncResp struct {
	URI      string `json:"uri,omitempty"`
	Duration string `json:"duration,omitempty"`
}

// PaginatedInstanceList returns a paginated list of Instance object
type PaginatedInstanceList struct {
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
	Pages   int        `json:"pages"`
	Items   []Instance `json:"items"`
}

// AttachedVolume disk information
type AttachedVolume struct {
	// ID of the volume to attach
	ID string `json:"id"`
}

// InstanceConfig describes the parameters for a new instance
// none of the fields are mandatory and will be automatically
// set with default values
type InstanceConfig struct {
	Count                 int              `json:"count"`
	Hostname              string           `json:"hostname"`
	ReverseDNS            string           `json:"reverse_dns"`
	Size                  string           `json:"size"`
	Region                string           `json:"region"`
	PublicIPRequired      string           `json:"public_ip"`
	ReservedIPv4          string           `json:"reserved_ipv4"`
	PrivateIPv4           string           `json:"private_ipv4"`
	NetworkID             string           `json:"network_id"`
	TemplateID            string           `json:"template_id"`
	SourceType            string           `json:"source_type"`
	SourceID              string           `json:"source_id"`
	SnapshotID            string           `json:"snapshot_id"`
	Subnets               []string         `json:"subnets,omitempty"`
	InitialUser           string           `json:"initial_user"`
	SSHKeyID              string           `json:"ssh_key_id"`
	Script                string           `json:"script"`
	Tags                  []string         `json:"-"`
	TagsList              string           `json:"tags"`
	FirewallID            string           `json:"firewall_id"`
	VolumeType            string           `json:"volume_type,omitempty"`
	AttachedVolumes       []AttachedVolume `json:"attached_volumes"`
	PlacementRule         PlacementRule    `json:"placement_rule"`
	NetworkBandwidthLimit int              `json:"network_bandwidth_limit,omitempty"`
	AllowedIPs            []string         `json:"allowed_ips,omitempty"`
}

// AffinityRule represents a affinity rule
type AffinityRule struct {
	Type      string   `json:"type"`
	Exclusive bool     `json:"exclusive"`
	Tags      []string `json:"tags"`
}

// PlacementRule represents a placement rule
type PlacementRule struct {
	AffinityRules []AffinityRule    `json:"affinity_rules,omitempty"`
	NodeSelector  map[string]string `json:"node_selector,omitempty"`
}

// ListInstances returns a paginated list of instances owned by the calling API account.
// This method supports pagination to handle large numbers of instances efficiently.
// Use page=0 and perPage=0 to get default pagination settings.
//
// Parameters:
//   - page: The page number to retrieve (1-based indexing)
//   - perPage: Number of instances per page (0 for default)
//
// Returns:
//   - *PaginatedInstanceList: Paginated list with instances and metadata
//   - error: Any error that occurred during the API request
func (c *Client) ListInstances(page int, perPage int) (*PaginatedInstanceList, error) {
	url := "/v2/instances"
	if page != 0 && perPage != 0 {
		url = url + fmt.Sprintf("?page=%d&per_page=%d", page, perPage)
	}

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, decodeError(err)
	}

	PaginatedInstances := PaginatedInstanceList{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&PaginatedInstances)
	return &PaginatedInstances, err
}

// ListAllInstances returns all instances owned by the calling API account.
// This is a convenience method that fetches all instances by requesting
// a very large page size. Use with caution for accounts with many instances.
//
// Returns:
//   - []Instance: A slice containing all instances
//   - error: Any error that occurred during the API request
func (c *Client) ListAllInstances() ([]Instance, error) {
	instances, err := c.ListInstances(1, 99999999)
	if err != nil {
		return []Instance{}, decodeError(err)
	}

	return instances.Items, nil
}

// FindInstance searches for an instance by partial ID or hostname match.
// This method provides a convenient way to locate instances when you only
// know part of the identifier. Returns an error if no matches or multiple
// partial matches are found.
//
// Parameters:
//   - search: Partial instance ID or hostname to search for
//
// Returns:
//   - *Instance: The found instance (exact match preferred over partial)
//   - error: MultipleMatchesError, ZeroMatchesError, or API errors
func (c *Client) FindInstance(search string) (*Instance, error) {
	instances, err := c.ListAllInstances()
	if err != nil {
		return nil, decodeError(err)
	}

	partialMatchesCount := 0
	result := Instance{}

	for _, value := range instances {
		if value.Hostname == search || value.ID == search {
			return &value, nil
		} else if strings.Contains(value.Hostname, search) || strings.Contains(value.ID, search) {
			partialMatchesCount++
			result = value
		}
	}

	if partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}

// GetInstance retrieves detailed information about a specific instance by its full ID.
// This method returns complete instance details including status, configuration,
// network settings, and resource allocations.
//
// Parameters:
//   - id: The complete instance ID (UUID)
//
// Returns:
//   - *Instance: Complete instance information
//   - error: Any error that occurred during retrieval
func (c *Client) GetInstance(id string) (*Instance, error) {
	resp, err := c.SendGetRequest("/v2/instances/" + id)
	if err != nil {
		return nil, decodeError(err)
	}

	instance := Instance{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&instance)
	return &instance, err
}

// NewInstanceConfig returns an initialized configuration for creating a new instance.
// This method sets up sensible defaults including a random hostname, default network,
// and standard configuration options. Modify the returned config as needed.
//
// Returns:
//   - *InstanceConfig: Pre-configured instance settings with defaults
//   - error: Any error that occurred while fetching default network
func (c *Client) NewInstanceConfig() (*InstanceConfig, error) {
	network, err := c.GetDefaultNetwork()
	if err != nil {
		return nil, decodeError(err)
	}

	return &InstanceConfig{
		Count:            1,
		Hostname:         utils.RandomName(),
		ReverseDNS:       "",
		Region:           c.Region,
		PublicIPRequired: "true",
		NetworkID:        network.ID,
		SnapshotID:       "",
		InitialUser:      DefaultInstanceUser,
		SSHKeyID:         "",
		Script:           "",
		Tags:             []string{""},
		FirewallID:       "",
	}, nil
}

// CreateInstance creates a new virtual server instance with the specified configuration.
// This provisions a new instance in the account with the given settings. The instance
// will be created asynchronously and its status can be monitored using GetInstance.
//
// Parameters:
//   - config: Complete instance configuration including size, image, network, etc.
//
// Returns:
//   - *Instance: The newly created instance information
//   - error: Any error that occurred during instance creation
func (c *Client) CreateInstance(config *InstanceConfig) (*Instance, error) {
	config.TagsList = strings.Join(config.Tags, " ")
	body, err := c.SendPostRequest("/v2/instances", config)
	if err != nil {
		return nil, decodeError(err)
	}

	var instance Instance
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&instance); err != nil {
		return nil, err
	}

	return &instance, nil
}

// SetInstanceTags updates the tags associated with the specified instance.
// Tags are space-separated labels that help organize and categorize instances.
// This operation replaces all existing tags with the new ones.
//
// Parameters:
//   - i: The instance to update (only ID is used)
//   - tags: Space-separated string of tags to apply
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during the tag update
func (c *Client) SetInstanceTags(i *Instance, tags string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/tags", i.ID), map[string]string{
		"tags":   tags,
		"region": c.Region,
	})
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// UpdateInstance modifies an existing instance's configuration.
// This method can update hostname, reverse DNS, notes, public IP, and subnet assignments.
// Use empty notes field with notes_delete=true to remove existing notes.
//
// Parameters:
//   - i: Instance object with updated fields (ID must be set)
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during the update
func (c *Client) UpdateInstance(i *Instance) (*SimpleResponse, error) {
	params := map[string]interface{}{
		"hostname":    i.Hostname,
		"reverse_dns": i.ReverseDNS,
		"notes":       i.Notes,
		"region":      c.Region,
		"public_ip":   i.PublicIP,
		"subnets":     i.Subnets,
	}

	if i.Notes == "" {
		params["notes_delete"] = "true"
	}

	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s", i.ID), params)
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// GetInstanceVnc enables and retrieves VNC access information for an instance.
// This creates a VNC session that allows remote desktop access to the instance.
// The session has a configurable duration and will expire automatically.
//
// Parameters:
//   - id: The instance ID to enable VNC for
//   - duration: Optional session duration (e.g., "30m", "1h", "24h"). If empty, uses default
//
// Returns:
//   - CreateInstanceVncResp: VNC connection details including URI and duration
//   - error: Any error that occurred during VNC session creation
func (c *Client) GetInstanceVnc(id string, duration ...string) (CreateInstanceVncResp, error) {
	url := fmt.Sprintf("/v2/instances/%s/vnc", id)
	if len(duration) > 0 && duration[0] != "" {
		url = fmt.Sprintf("%s?duration=%s", url, duration[0])
	}

	resp, err := c.SendPutRequest(url, map[string]string{
		"region": c.Region,
	})
	vnc := CreateInstanceVncResp{}

	if err != nil {
		return vnc, decodeError(err)
	}

	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&vnc)
	return vnc, err
}

// GetInstanceVncStatus returns the current VNC session status for an instance.
// This method checks if there's an active VNC session and provides connection
// details if one exists, including expiration time.
//
// Parameters:
//   - id: The instance ID to check VNC status for
//
// Returns:
//   - *InstanceVnc: Current VNC session information (URI and expiration)
//   - error: Any error that occurred during status retrieval
func (c *Client) GetInstanceVncStatus(id string) (*InstanceVnc, error) {
	url := fmt.Sprintf("/v2/instances/%s/vnc", id)
	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, decodeError(err)
	}

	vnc := InstanceVnc{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(&vnc)
	return &vnc, err

}

// DeleteInstanceVncSession terminates the active VNC session for an instance.
// This immediately closes any active VNC connections and prevents further
// remote desktop access until a new session is created.
//
// Parameters:
//   - id: The instance ID to terminate VNC session for
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during session termination
func (c *Client) DeleteInstanceVncSession(id string) (*SimpleResponse, error) {
	url := fmt.Sprintf("/v2/instances/%s/vnc", id)
	resp, err := c.SendDeleteRequest(url)
	if err != nil {
		return nil, decodeError(err)
	}
	return c.DecodeSimpleResponse(resp)
}

// DeleteInstance permanently removes an instance and frees all its resources.
// This action is irreversible and will destroy the instance, its local storage,
// and any data that hasn't been backed up. Network resources may be freed for reuse.
//
// Parameters:
//   - id: The instance ID to delete
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during deletion
func (c *Client) DeleteInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest("/v2/instances/" + id)
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// RebootInstance restarts an instance using a hard reboot method.
// This is equivalent to calling HardRebootInstance and forces an immediate
// restart similar to pressing a physical power button.
//
// Parameters:
//   - id: The instance ID to reboot
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during reboot initiation
func (c *Client) RebootInstance(id string) (*SimpleResponse, error) {
	return c.HardRebootInstance(id)
}

// HardRebootInstance performs a forced restart of an instance.
// This is equivalent to pulling the power cord and reconnecting it - the instance
// is immediately powered off and restarted without graceful shutdown procedures.
//
// Parameters:
//   - id: The instance ID to hard reboot
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during hard reboot initiation
func (c *Client) HardRebootInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPostRequest(fmt.Sprintf("/v2/instances/%s/hard_reboots", id), map[string]string{
		"region": c.Region,
	})
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// SoftRebootInstance performs a graceful restart of an instance.
// This sends a shutdown signal to the operating system, allowing it to
// close applications and services cleanly before restarting.
//
// Parameters:
//   - id: The instance ID to soft reboot
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during soft reboot initiation
func (c *Client) SoftRebootInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPostRequest(fmt.Sprintf("/v2/instances/%s/soft_reboots", id), map[string]string{
		"region": c.Region,
	})
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// StopInstance gracefully shuts down an instance without destroying it.
// The instance remains allocated and can be restarted later using StartInstance.
// This is similar to shutting down a computer - it's powered off but not deleted.
//
// Parameters:
//   - id: The instance ID to stop
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during shutdown initiation
func (c *Client) StopInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/stop", id), map[string]string{
		"region": c.Region,
	})
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// StartInstance powers on a previously stopped instance.
// This boots the instance from its shutdown state, resuming normal operation.
// The instance retains all its configuration and data from before it was stopped.
//
// Parameters:
//   - id: The instance ID to start
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during startup initiation
func (c *Client) StartInstance(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/start", id), map[string]string{
		"region": c.Region,
	})
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// UpgradeInstance resizes an instance to a larger specification.
// This increases the instance's resources (CPU, RAM, disk) to the specified size.
// Note: Instances can only be upgraded to larger sizes, not downgraded to smaller ones.
//
// Parameters:
//   - id: The instance ID to upgrade
//   - newSize: The target size identifier (e.g., "g3.medium", "g3.large")
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during upgrade initiation
func (c *Client) UpgradeInstance(id, newSize string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/resize", id), map[string]string{
		"size":   newSize,
		"region": c.Region,
	})
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// MovePublicIPToInstance transfers a public IP address to the specified instance.
// This operation moves an existing public IP from another instance or the IP pool
// to the target instance, updating network routing accordingly.
//
// Parameters:
//   - id: The instance ID to receive the public IP
//   - ipAddress: The public IP address to move
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during IP address migration
func (c *Client) MovePublicIPToInstance(id, ipAddress string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/ip/%s", id, ipAddress), "")
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// SetInstanceFirewall assigns a firewall to an instance.
// This changes the network security rules applied to the instance by associating
// it with a different firewall. The change takes effect immediately.
//
// Parameters:
//   - id: The instance ID to update
//   - firewallID: The ID of the firewall to assign
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during firewall assignment
func (c *Client) SetInstanceFirewall(id, firewallID string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/firewall", id), map[string]string{
		"firewall_id": firewallID,
		"region":      c.Region,
	})
	if err != nil {
		return nil, decodeError(err)
	}

	response, err := c.DecodeSimpleResponse(resp)
	return response, err
}

// EnableRecoveryMode activates recovery mode for an instance.
// Recovery mode allows access to an instance that may be in an unbootable state,
// providing emergency access for troubleshooting and data recovery.
//
// Parameters:
//   - id: The instance ID to enable recovery mode for
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during recovery mode activation
func (c *Client) EnableRecoveryMode(id string) (*SimpleResponse, error) {
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/recovery?region=%s", id, c.Region), nil)
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}

// DisableRecoveryMode deactivates recovery mode for an instance.
// This returns the instance to normal operation mode, removing special
// recovery access and restoring standard boot procedures.
//
// Parameters:
//   - id: The instance ID to disable recovery mode for
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during recovery mode deactivation
func (c *Client) DisableRecoveryMode(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/instances/%s/recovery?region=%s", id, c.Region))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}

// GetRecoveryStatus retrieves the current recovery mode status for an instance.
// This method checks whether recovery mode is currently active and provides
// relevant status information about the recovery state.
//
// Parameters:
//   - id: The instance ID to check recovery status for
//
// Returns:
//   - *SimpleResponse: Recovery status information
//   - error: Any error that occurred during status retrieval
func (c *Client) GetRecoveryStatus(id string) (*SimpleResponse, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/instances/%s/recovery", id))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}

// UpdateInstanceAllowedIPs sets the list of IP addresses that an instance is allowed to communicate with.
// This configures network access control by restricting the instance's outbound connections
// to only the specified IP addresses or CIDR blocks.
//
// Parameters:
//   - id: The instance ID to update
//   - allowedIPs: List of IP addresses or CIDR blocks to allow
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during the update
func (c *Client) UpdateInstanceAllowedIPs(id string, allowedIPs []string) (*SimpleResponse, error) {
	// Create a map to match the expected JSON structure
	payload := map[string][]string{
		"allowed_ips": allowedIPs,
	}
	// Send the payload map instead of the raw allowedIPs slice
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/allowed_ips", id), payload)
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}

// UpdateInstanceBandwidth updates the network bandwidth limit for a specific instance.
// This sets the maximum network bandwidth in MB/s that the instance can use.
//
// Parameters:
//   - id: The ID of the instance to update
//   - bandwidthLimit: The network bandwidth limit in MB/s
//
// Returns:
//   - *SimpleResponse: Response indicating success or failure
//   - error: Any error that occurred during the operation
func (c *Client) UpdateInstanceBandwidth(id string, bandwidthLimit int) (*SimpleResponse, error) {
	payload := map[string]int{
		"network_bandwidth_limit": bandwidthLimit,
	}
	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/instances/%s/network_bandwidth_limit", id), payload)
	if err != nil {
		return nil, decodeError(err)
	}
	return c.DecodeSimpleResponse(resp)
}
