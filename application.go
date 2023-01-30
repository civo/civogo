package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/civo/civogo/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProcInfo is the struct for process information
type ProcInfo struct {
	ProcessType  string `json:"process_type,omitempty"`
	ProcessCount int    `json:"process_count,omitempty"`

	// Currently, command and Args are only applicable for when container image is provided, not when git repo is provided
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// Application is the struct for the Application model
type Application struct {
	Name               string         `json:"name" schema:"name"`
	ID                 string         `json:"id"`
	NetworkID          string         `json:"network_id" validate:"required" schema:"network_id"`
	FirewallID         string         `json:"firewall_id" schema:"firewall_id"`
	Image              *string        `json:"image,omitempty"  schema:"image"`
	Size               string         `json:"size"  schema:"size"`
	ProcessInfo        []ProcInfo     `json:"process_info,omitempty"`
	GitInfo            *GitInfo       `json:"git_info,omitempty" schema:"git_info"`
	Config             ObservedConfig `json:"config,omitempty"`
	AppIP              string         `json:"app_ip,omitempty"`
	Domains            []string       `json:"domains,omitempty"`
	PublicIPv4Required bool           `json:"public_ipv4_required" schema:"public_ipv4_required"`
	Status             string         `json:"status"`
}

// GitInfo holds the git information for the application
type GitInfo struct {
	GitURL         string          `json:"url" schema:"url"`
	GitToken       string          `json:"token" schema:"token"`
	PullPreference *PullPreference `json:"pull_preferences,omitempty" schema:"pull_preferences"`
}

// PullPreference determines which tag/branch should the image be built from.
// If both are specified, the tag will be used(Tags>Branches)
// If neither is specified, main/master will be used (main > master)
type PullPreference struct {
	Tag    *string `json:"tag,omitempty" schema:"tag"`
	Branch *string `json:"branch,omitempty" schema:"branch"`
}

// ObservedConfig defines the observed state of EnvVar
type ObservedConfig struct {
	Env          []EnvVar    `json:"env_vars,omitempty"`
	LastSyncedAt metav1.Time `json:"last_sycned_at,omitempty"`
}

// EnvVar holds key-value pairs for an application
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RegistryAuth holds the registry auth for an application
type RegistryAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ApplicationConfig describes the parameters for a new CivoApp
type ApplicationConfig struct {
	Name               string   `json:"name"`
	NetworkID          string   `json:"network_id" validate:"required"`
	Size               string   `json:"size"`
	Image              *string  `json:"image,omitempty"`
	GitInfo            *GitInfo `json:"git_info,omitempty"`
	PublicIPv4Required bool     `json:"public_ipv4_required" schema:"public_ipv4_required"`
	Region             string   `json:"region"`
}

// PaginatedApplications returns a paginated list of Application object
type PaginatedApplications struct {
	Page    int           `json:"page"`
	PerPage int           `json:"per_page"`
	Pages   int           `json:"pages"`
	Items   []Application `json:"items"`
}

// UpdateApplicationRequest is the struct for the UpdateApplication request
type UpdateApplicationRequest struct {
	Name         string        `json:"name"  schema:"name"`
	Size         string        `json:"size" schema:"size"`
	ProcessInfo  []ProcInfo    `json:"process_info" schema:"process_info"`
	GitInfo      *GitInfo      `json:"git_info,omitempty" schema:"git_info"`
	EnvVars      []EnvVar      `json:"env_vars,omitempty" schema:"env_vars"`
	FirewallID   string        `json:"firewall_id" schema:"firewall_id"`
	RegistryAuth *RegistryAuth `json:"registry_auth" schema:"registry_auth"`
	Image        *string       `json:"image,omitempty"  schema:"image"`
	Region       string        `json:"region"`
}

// ErrAppDomainNotFound is returned when the domain is not found
var ErrAppDomainNotFound = fmt.Errorf("domain not found")

// ListApplications returns all applications in that specific region
func (c *Client) ListApplications() (*PaginatedApplications, error) {
	resp, err := c.SendGetRequest("/v2/applications")
	if err != nil {
		return nil, decodeError(err)
	}

	application := &PaginatedApplications{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&application); err != nil {
		return nil, decodeError(err)
	}

	return application, nil
}

// GetApplication returns an application by ID
func (c *Client) GetApplication(id string) (*Application, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/applications/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	application := &Application{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&application); err != nil {
		return nil, decodeError(err)
	}

	return application, nil
}

// NewApplicationConfig returns an initialized config for a new application
func (c *Client) NewApplicationConfig() (*ApplicationConfig, error) {
	network, err := c.GetDefaultNetwork()
	if err != nil {
		return nil, decodeError(err)
	}

	return &ApplicationConfig{
		Name:      utils.RandomName(),
		NetworkID: network.ID,
		Size:      "small",
	}, nil
}

// FindApplication finds an application by either part of the ID or part of the name
func (c *Client) FindApplication(search string) (*Application, error) {
	apps, err := c.ListApplications()
	if err != nil {
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := Application{}

	for _, value := range apps.Items {
		if value.Name == search || value.ID == search {
			exactMatch = true
			result = value
		} else if strings.Contains(value.Name, search) || strings.Contains(value.ID, search) {
			if !exactMatch {
				result = value
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

// CreateApplication creates a new application
func (c *Client) CreateApplication(config *ApplicationConfig) (*Application, error) {
	body, err := c.SendPostRequest("/v2/applications", config)
	if err != nil {
		return nil, decodeError(err)
	}

	var application Application
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&application); err != nil {
		return nil, err
	}

	return &application, nil
}

// UpdateApplication updates an application
func (c *Client) UpdateApplication(id string, application *UpdateApplicationRequest) (*Application, error) {
	body, err := c.SendPutRequest(fmt.Sprintf("/v2/applications/%s", id), application)
	if err != nil {
		return nil, decodeError(err)
	}

	updatedApplication := &Application{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(updatedApplication); err != nil {
		return nil, err
	}

	return updatedApplication, nil
}

// DeleteApplication deletes an application
func (c *Client) DeleteApplication(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/applications/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}

	return c.DecodeSimpleResponse(resp)
}
