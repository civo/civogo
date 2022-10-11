package civogo

import (
	"context"
	"time"

	// "errors"
	// "fmt"
	"net/http"
)

const instanceBasePath = "/v2/instances"

// NetworkService is an interface for interfacing with the network
type InstancesService interface {
	List(ctx context.Context) (*PaginatedInstancesList, *Metadata, error)
	// Default(ctx context.Context) (*Network, *Metadata, error)
	// GetByID(ctx context.Context, networkID string) (*Network, *Metadata, error)
	// Find(ctx context.Context, value string) (*Network, *Metadata, error)
	// Create(ctx context.Context, createRequest *NetworkCreateRequest) (*SimpleResponse, *Metadata, error)
	// Update(ctx context.Context, networkID string, updateRequest *NetworkUpdateRequest) (*SimpleResponse, *Metadata, error)
	// Delete(ctx context.Context, networkID string) (*SimpleResponse, *Metadata, error)
}

// NetworkServiceOp Service used for communicating with the API
type InstancesServiceOp struct {
	client  *Client
	network string
}

type InstancesGetter interface {
	Instances(network string) InstancesService
}

// newSSHKey returns a SSHKey
func newInstances(c *Client, network string) *InstancesServiceOp {
	return &InstancesServiceOp{
		client:  c,
		network: network,
	}
}

// Instance represents a virtual server within Civo's infrastructure
type Instance struct {
	ID                       string    `json:"id,omitempty"`
	Hostname                 string    `json:"hostname,omitempty"`
	ReverseDNS               string    `json:"reverse_dns,omitempty"`
	Size                     string    `json:"size,omitempty"`
	Region                   string    `json:"region,omitempty"`
	NetworkID                string    `json:"network_id,omitempty"`
	PrivateIP                string    `json:"private_ip,omitempty"`
	PublicIP                 string    `json:"public_ip,omitempty"`
	PseudoIP                 string    `json:"pseudo_ip,omitempty"`
	TemplateID               string    `json:"template_id,omitempty"`
	SourceType               string    `json:"source_type,omitempty"`
	SourceID                 string    `json:"source_id,omitempty"`
	SnapshotID               string    `json:"snapshot_id,omitempty"`
	InitialUser              string    `json:"initial_user,omitempty"`
	InitialPassword          string    `json:"initial_password,omitempty"`
	SSHKey                   string    `json:"ssh_key,omitempty"`
	SSHKeyID                 string    `json:"ssh_key_id,omitempty"`
	Status                   string    `json:"status,omitempty"`
	Notes                    string    `json:"notes,omitempty"`
	FirewallID               string    `json:"firewall_id,omitempty"`
	Tags                     []string  `json:"tags,omitempty"`
	CivostatsdToken          string    `json:"civostatsd_token,omitempty"`
	CivostatsdStats          string    `json:"civostatsd_stats,omitempty"`
	CivostatsdStatsPerMinute []string  `json:"civostatsd_stats_per_minute,omitempty"`
	CivostatsdStatsPerHour   []string  `json:"civostatsd_stats_per_hour,omitempty"`
	VolumeBacked             bool      `json:"volume_backed,omitempty"`
	CPUCores                 int       `json:"cpu_cores,omitempty"`
	RAMMegabytes             int       `json:"ram_mb,omitempty"`
	DiskGigabytes            int       `json:"disk_gb,omitempty"`
	Script                   string    `json:"script,omitempty"`
	CreatedAt                time.Time `json:"created_at,omitempty"`
	ReservedIPID             string    `json:"reserved_ip_id,omitempty"`
	ReservedIPName           string    `json:"reserved_ip_name,omitempty"`
	ReservedIP               string    `json:"reserved_ip,omitempty"`
}

// PaginatedInstanceList returns a paginated list of Instance object
type PaginatedInstancesList struct {
	Page    int        `json:"page"`
	PerPage int        `json:"per_page"`
	Pages   int        `json:"pages"`
	Items   []Instance `json:"items"`
}

// InstanceConfig describes the parameters for a new instance
// none of the fields are mandatory and will be automatically
// set with default values
type InstanceCreateRequest struct {
	Count            int      `json:"count"`
	Hostname         string   `json:"hostname"`
	ReverseDNS       string   `json:"reverse_dns"`
	Size             string   `json:"size"`
	Region           string   `json:"region"`
	PublicIPRequired string   `json:"public_ip"`
	NetworkID        string   `json:"network_id"`
	TemplateID       string   `json:"template_id"`
	SourceType       string   `json:"source_type"`
	SourceID         string   `json:"source_id"`
	SnapshotID       string   `json:"snapshot_id"`
	InitialUser      string   `json:"initial_user"`
	SSHKeyID         string   `json:"ssh_key_id"`
	Script           string   `json:"script"`
	Tags             []string `json:"-"`
	TagsList         string   `json:"tags"`
	FirewallID       string   `json:"firewall_id"`
}

// List function returns a list of all instances you can pass the network to filter the results
func (s *InstancesServiceOp) List(ctx context.Context) (*PaginatedInstancesList, *Metadata, error) {
	path := instanceBasePath

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(PaginatedInstancesList)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if s.network != "" && len(root.Items) > 0 {
		network, meta, err := s.client.Network().Find(ctx, s.network)
		if err != nil {
			return nil, meta, err
		}

		for _, items := range root.Items {
			if items.NetworkID != network.ID {
				// remove the instance from the list
				root.Items = removeInstance(root.Items, items)
			}
		}
	}

	return root, resp, err
}

// removeInstance removes an instance from the list
func removeInstance(s []Instance, r Instance) []Instance {
	for i, v := range s {
		if v.ID == r.ID {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
