package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type KubernetesInstances struct {
	Hostname   string    `json:"hostname"`
	Size       string    `json:"size"`
	Region     string    `json:"region"`
	CreatedAt  time.Time `json:"created_at"`
	Status     string    `json:"status"`
	FirewallID string    `json:"firewall_id"`
	PublicIP   string    `json:"public_ip"`
	Tags       []string  `json:"tags"`
}

type KubernetesApplications struct {
	Application   string            `json:"application"`
	Title         string            `json:"title,omitempty"`
	Version       string            `json:"version"`
	Dependencies  []string          `json:"dependencies,omitempty"`
	Maintainer    string            `json:"maintainer"`
	Description   string            `json:"description"`
	PostInstall   string            `json:"post_install"`
	Installed     bool              `json:"installed"`
	Url           string            `json:"url"`
	Category      string            `json:"category"`
	UpdatedAt     time.Time         `json:"updated_at"`
	ImageUrl      string            `json:"image_url"`
	Plan          string            `json:"plan,omitempty"`
	Configuration map[string]string `json:"configuration,omitempty"`
}

type KubernetesCluster struct {
	ID                    string                   `json:"id"`
	Name                  string                   `json:"name"`
	Version               string                   `json:"version"`
	Status                string                   `json:"status"`
	Ready                 bool                     `json:"ready"`
	NumTargetNode         int                      `json:"num_target_nodes"`
	TargetNodeSize        string                   `json:"target_nodes_size"`
	BuiltAt               time.Time                `json:"built_at"`
	KubeConfig            string                   `json:"kubeconfig"`
	KubernetesVersion     string                   `json:"kubernetes_version"`
	ApiEndPoint           string                   `json:"api_endpoint"`
	DnsEntry              string                   `json:"dns_entry"`
	Tags                  []string                 `json:"tags"`
	CreatedAt             time.Time                `json:"created_at"`
	Instances             []KubernetesInstances    `json:"instances"`
	InstalledApplications []KubernetesApplications `json:"installed_applications"`
}

type KubernetesClusterConfig struct {
	Name              string `form:"name"`
	NumTargetNodes    int    `form:"num_target_nodes"`
	TargetNodesSize   string `form:"target_nodes_size"`
	KubernetesVersion string `form:"kubernetes_version"`
	Tags              string `form:"tags"`
	Applications      string `form:"applications"`
}

type KubernetesPlanConfiguration struct {
	Value string `json:"value"`
}

type KubernetesMarketplacePlan struct {
	Label         string                                 `json:"label"`
	Configuration map[string]KubernetesPlanConfiguration `json:"configuration"`
}

type KubernetesMarketplaceApplication struct {
	Name         string                      `json:"name"`
	Title        string                      `json:"title,omitempty"`
	Version      string                      `json:"version"`
	Default      string                      `json:"default,omitempty"`
	Dependencies []string                    `json:"dependencies,omitempty"`
	Maintainer   string                      `json:"maintainer"`
	Description  string                      `json:"description"`
	PostInstall  string                      `json:"post_install"`
	Url          string                      `json:"url"`
	Category     string                      `json:"category"`
	Plans        []KubernetesMarketplacePlan `json:"plans"`
}

type KubernetesVersion struct {
	Version string `json:"version"`
	Type    string `json:"type"`
	Default bool   `json:"default,omitempty"`
}

// ListKubernetesCluster returns all cluster of kubernetes in the account
func (c *Client) ListKubernetesClusters() ([]KubernetesCluster, error) {
	resp, err := c.SendGetRequest("/v2/kubernetes/clusters")
	if err != nil {
		return nil, err
	}

	kubernetes := make([]KubernetesCluster, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}

// NewKubernetesCluster create a new cluster of kubernetes
func (c *Client) NewKubernetesClusters(kc *KubernetesClusterConfig) (*KubernetesCluster, error) {
	body, err := c.SendPostRequest("/v2/kubernetes/clusters", kc)
	if err != nil {
		return nil, err
	}

	kubernetes := &KubernetesCluster{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}

// GetKubernetesCluster returns a single kubernetes cluster by its full ID
func (c *Client) GetKubernetesClusters(id string) (*KubernetesCluster, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id))
	if err != nil {
		return nil, err
	}

	kubernetes := &KubernetesCluster{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(kubernetes)
	return kubernetes, nil
}

// UpdateKubernetesCluster update a single kubernetes cluster by its full ID
func (c *Client) UpdateKubernetesClusters(id string, i *KubernetesClusterConfig) (*KubernetesCluster, error) {
	params := map[string]interface{}{
		"name":             i.Name,
		"num_target_nodes": i.NumTargetNodes,
		"version":          i.KubernetesVersion,
		"applications":     i.Applications,
	}

	resp, err := c.SendPutRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id), params)
	if err != nil {
		return nil, err
	}

	kubernetes := &KubernetesCluster{}
	err = json.NewDecoder(bytes.NewReader(resp)).Decode(kubernetes)
	return kubernetes, nil
}

//ListKubernetesMarketplace returns all application inside marketplace
func (c *Client) ListKubernetesMarketplaceApplications() ([]KubernetesMarketplaceApplication, error) {
	resp, err := c.SendGetRequest("/v2/kubernetes/applications")
	if err != nil {
		return nil, err
	}

	kubernetes := make([]KubernetesMarketplaceApplication, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}

// DeleteKubernetesCluster deletes a cluster
func (c *Client) DeleteKubernetesClusters(id string) (*SimpleResponse, error) {
	resp, err := c.SendDeleteRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id))
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(resp)
}

// RecycleKubernetesCluster create a new cluster of kubernetes
func (c *Client) RecycleKubernetesClusters(id string, hostname string) (*SimpleResponse, error) {
	body, err := c.SendPostRequest(fmt.Sprintf("/v2/kubernetes/clusters/%s", id), map[string]string{
		"hostname": hostname,
	})
	if err != nil {
		return nil, err
	}

	return c.DecodeSimpleResponse(body)
}

// VersionKubernetesCluster returns all version of kubernetes in the cloud
func (c *Client) ListAvailableKubernetesVersions() ([]KubernetesVersion, error) {
	resp, err := c.SendGetRequest("/v2/kubernetes/versions")
	if err != nil {
		return nil, err
	}

	kubernetes := make([]KubernetesVersion, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&kubernetes); err != nil {
		return nil, err
	}

	return kubernetes, nil
}
