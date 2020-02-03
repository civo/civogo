package civogo

import (
	"bytes"
	"encoding/json"
	"time"
)

type KubeInstances struct {
	Hostname   string    `json:"hostname"`
	Size       string    `json:"size"`
	Region     string    `json:"region"`
	CreatedAt  time.Time `json:"created_at"`
	Status     string    `json:"status"`
	FirewallID string    `json:"firewall_id"`
	PublicIP   string    `json:"public_ip"`
	Tags       []string  `json:"tags"`
}

type KubeApplications struct {
	Application   string    `json:"application"`
	Title         string    `json:"title,omitempty"`
	Version       string    `json:"version"`
	Dependencies  string    `json:"dependencies,omitempty"`
	Maintainer    string    `json:"maintainer"`
	Description   string    `json:"description"`
	PostInstall   string    `json:"post_install"`
	Installed     bool      `json:"installed"`
	Url           string    `json:"url"`
	Category      string    `json:"category"`
	UpdatedAt     time.Time `json:"updated_at"`
	ImageUrl      string    `json:"image_url"`
	Plan          string    `json:"plan,omitempty"`
	Configuration string    `json:"plan,omitempty"`
}

type Kubernetes struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Version               string    `json:"version"`
	Status                string    `json:"status"`
	Ready                 bool      `json:"ready"`
	NumTargetNode         bool      `json:"num_target_nodes"`
	TargetNodeSize        bool      `json:"target_nodes_size"`
	BuiltAt               time.Time `json:"built_at"`
	KubeConfig            string    `json:"kubeconfig"`
	KubernetesVersion     string    `json:"kubernetes_version"`
	ApiEndPoint           string    `json:"api_endpoint"`
	DnsEntry              string    `json:"dns_entry"`
	Tags                  []string  `json:"tags"`
	CreatedAt             time.Time `json:"created_at"`
	Instances             []KubeInstances
	InstalledApplications []KubeApplications
}

type KubernetesConfig struct {
	Name              string `form:"name"`
	NumTargetNodes    string `form:"num_target_nodes"`
	TargetNodesSize   string `form:"target_nodes_size"`
	KubernetesVersion string `form:"kubernetes_version"`
	Tags              string `form:"tags"`
	Applications      string `form:"applications"`
}

// ListCluster returns all cluster of kubernetes in the account
func (c *Client) ListCluster() ([]Kubernetes, error) {
	resp, err := c.SendGetRequest("/v2/kubernetes/clusters")
	if err != nil {
		return nil, err
	}

	kubernetes := make([]Kubernetes, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}
