package civogo

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// FakeClient is a temporary storage structure for use when you don't want to communicate with a real Civo API server
type FakeClient struct {
	LastID                  int64
	Charges                 []Charge
	Domains                 []DNSDomain
	DomainRecords           []DNSRecord
	Firewalls               []Firewall
	FirewallRules           []FirewallRule
	InstanceSizes           []InstanceSize
	Instances               []Instance
	Kubernetes              []KubernetesCluster
	IP                      []IP
	Networks                []Network
	Volumes                 []Volume
	SSHKeys                 []SSHKey
	Webhooks                []Webhook
	DiskImage               []DiskImage
	Quota                   Quota
	Organisation            Organisation
	OrganisationAccounts    []Account
	OrganisationRoles       []Role
	OrganisationTeams       []Team
	OrganisationTeamMembers map[string][]TeamMember
	LoadBalancers           []LoadBalancer
	Pools                   []KubernetesPool
	ObjectStore             []ObjectStore
	ObjectStoreCredential   []ObjectStoreCredential
	// Snapshots            []Snapshot
	// Templates            []Template
}

// Clienter is the interface the real civogo.Client and civogo.FakeClient implement
type Clienter interface {
	// Charges
	ListCharges(from, to time.Time) ([]Charge, error)

	// DNS
	ListDNSDomains() ([]DNSDomain, error)
	FindDNSDomain(search string) (*DNSDomain, error)
	CreateDNSDomain(name string) (*DNSDomain, error)
	GetDNSDomain(name string) (*DNSDomain, error)
	UpdateDNSDomain(d *DNSDomain, name string) (*DNSDomain, error)
	DeleteDNSDomain(d *DNSDomain) (*SimpleResponse, error)
	CreateDNSRecord(domainID string, r *DNSRecordConfig) (*DNSRecord, error)
	ListDNSRecords(dnsDomainID string) ([]DNSRecord, error)
	GetDNSRecord(domainID, domainRecordID string) (*DNSRecord, error)
	UpdateDNSRecord(r *DNSRecord, rc *DNSRecordConfig) (*DNSRecord, error)
	DeleteDNSRecord(r *DNSRecord) (*SimpleResponse, error)

	// Firewalls
	ListFirewalls() ([]Firewall, error)
	FindFirewall(search string) (*Firewall, error)
	NewFirewall(*FirewallConfig) (*FirewallResult, error)
	RenameFirewall(id string, f *FirewallConfig) (*SimpleResponse, error)
	DeleteFirewall(id string) (*SimpleResponse, error)
	NewFirewallRule(r *FirewallRuleConfig) (*FirewallRule, error)
	ListFirewallRules(id string) ([]FirewallRule, error)
	FindFirewallRule(firewallID string, search string) (*FirewallRule, error)
	DeleteFirewallRule(id string, ruleID string) (*SimpleResponse, error)

	// Instances
	ListInstances(page int, perPage int) (*PaginatedInstanceList, error)
	ListAllInstances() ([]Instance, error)
	FindInstance(search string) (*Instance, error)
	GetInstance(id string) (*Instance, error)
	NewInstanceConfig() (*InstanceConfig, error)
	CreateInstance(config *InstanceConfig) (*Instance, error)
	SetInstanceTags(i *Instance, tags string) (*SimpleResponse, error)
	UpdateInstance(i *Instance) (*SimpleResponse, error)
	DeleteInstance(id string) (*SimpleResponse, error)
	RebootInstance(id string) (*SimpleResponse, error)
	HardRebootInstance(id string) (*SimpleResponse, error)
	SoftRebootInstance(id string) (*SimpleResponse, error)
	StopInstance(id string) (*SimpleResponse, error)
	StartInstance(id string) (*SimpleResponse, error)
	GetInstanceConsoleURL(id string) (string, error)
	UpgradeInstance(id, newSize string) (*SimpleResponse, error)
	MovePublicIPToInstance(id, ipAddress string) (*SimpleResponse, error)
	SetInstanceFirewall(id, firewallID string) (*SimpleResponse, error)

	// Instance sizes
	ListInstanceSizes() ([]InstanceSize, error)
	FindInstanceSizes(search string) (*InstanceSize, error)

	// Clusters
	ListKubernetesClusters() (*PaginatedKubernetesClusters, error)
	FindKubernetesCluster(search string) (*KubernetesCluster, error)
	NewKubernetesClusters(kc *KubernetesClusterConfig) (*KubernetesCluster, error)
	GetKubernetesCluster(id string) (*KubernetesCluster, error)
	UpdateKubernetesCluster(id string, i *KubernetesClusterConfig) (*KubernetesCluster, error)
	ListKubernetesMarketplaceApplications() ([]KubernetesMarketplaceApplication, error)
	DeleteKubernetesCluster(id string) (*SimpleResponse, error)
	RecycleKubernetesCluster(id string, hostname string) (*SimpleResponse, error)
	ListAvailableKubernetesVersions() ([]KubernetesVersion, error)
	ListKubernetesClusterInstances(id string) ([]Instance, error)
	FindKubernetesClusterInstance(clusterID, search string) (*Instance, error)

	//Pools
	ListKubernetesClusterPools(cid string) ([]KubernetesPool, error)
	GetKubernetesClusterPool(cid, pid string) (*KubernetesPool, error)
	FindKubernetesClusterPool(cid, search string) (*KubernetesPool, error)
	DeleteKubernetesClusterPoolInstance(cid, pid, id string) (*SimpleResponse, error)
	UpdateKubernetesClusterPool(cid, pid string, config *KubernetesClusterPoolUpdateConfig) (*KubernetesPool, error)

	// Networks
	GetDefaultNetwork() (*Network, error)
	NewNetwork(label string) (*NetworkResult, error)
	ListNetworks() ([]Network, error)
	FindNetwork(search string) (*Network, error)
	RenameNetwork(label, id string) (*NetworkResult, error)
	DeleteNetwork(id string) (*SimpleResponse, error)

	// Quota
	GetQuota() (*Quota, error)

	// Regions
	ListRegions() ([]Region, error)

	// SSHKeys
	ListSSHKeys() ([]SSHKey, error)
	NewSSHKey(name string, publicKey string) (*SimpleResponse, error)
	UpdateSSHKey(name string, sshKeyID string) (*SSHKey, error)
	FindSSHKey(search string) (*SSHKey, error)
	DeleteSSHKey(id string) (*SimpleResponse, error)

	// DiskImages
	ListDiskImages() ([]DiskImage, error)
	GetDiskImage(id string) (*DiskImage, error)
	FindDiskImage(search string) (*DiskImage, error)

	// Volumes
	ListVolumes() ([]Volume, error)
	GetVolume(id string) (*Volume, error)
	FindVolume(search string) (*Volume, error)
	NewVolume(v *VolumeConfig) (*VolumeResult, error)
	ResizeVolume(id string, size int) (*SimpleResponse, error)
	AttachVolume(id string, instance string) (*SimpleResponse, error)
	DetachVolume(id string) (*SimpleResponse, error)
	DeleteVolume(id string) (*SimpleResponse, error)

	// Webhooks
	CreateWebhook(r *WebhookConfig) (*Webhook, error)
	ListWebhooks() ([]Webhook, error)
	FindWebhook(search string) (*Webhook, error)
	UpdateWebhook(id string, r *WebhookConfig) (*Webhook, error)
	DeleteWebhook(id string) (*SimpleResponse, error)

	// Reserved IPs
	ListIPs() (*PaginatedIPs, error)
	FindIP(search string) (*IP, error)
	GetIP(id string) (*IP, error)
	NewIP(v *CreateIPRequest) (*IP, error)
	UpdateIP(id string, v *UpdateIPRequest) (*IP, error)
	DeleteIP(id string) (*SimpleResponse, error)
	AssignIP(id, resourceID, resourceType, region string) (*SimpleResponse, error)
	UnassignIP(id, region string) (*SimpleResponse, error)

	// LoadBalancer
	ListLoadBalancers() ([]LoadBalancer, error)
	GetLoadBalancer(id string) (*LoadBalancer, error)
	FindLoadBalancer(search string) (*LoadBalancer, error)
	CreateLoadBalancer(r *LoadBalancerConfig) (*LoadBalancer, error)
	UpdateLoadBalancer(id string, r *LoadBalancerUpdateConfig) (*LoadBalancer, error)
	DeleteLoadBalancer(id string) (*SimpleResponse, error)

	// ObjectStore
	ListObjectStores() (*PaginatedObjectstores, error)
	GetObjectStore(id string) (*ObjectStore, error)
	FindObjectStore(search string) (*ObjectStore, error)
	NewObjectStore(v *CreateObjectStoreRequest) (*ObjectStore, error)
	UpdateObjectStore(id string, v *UpdateObjectStoreRequest) (*ObjectStore, error)
	DeleteObjectStore(id string) (*SimpleResponse, error)
	GetObjectStoreStats(id string) (*ObjectStoreStats, error)

	// ObjectStoreCredentials
	ListObjectStoreCredentials() (*PaginatedObjectStoreCredentials, error)
	GetObjectStoreCredential(id string) (*ObjectStoreCredential, error)
	FindObjectStoreCredential(search string) (*ObjectStoreCredential, error)
	NewObjectStoreCredential(v *CreateObjectStoreCredentialRequest) (*ObjectStoreCredential, error)
	UpdateObjectStoreCredential(id string, v *UpdateObjectStoreCredentialRequest) (*ObjectStoreCredential, error)
	DeleteObjectStoreCredential(id string) (*SimpleResponse, error)
}

// NewFakeClient initializes a Client that doesn't attach to a
func NewFakeClient() (*FakeClient, error) {
	return &FakeClient{
		Quota: Quota{
			CPUCoreLimit:           10,
			InstanceCountLimit:     10,
			RAMMegabytesLimit:      100,
			DiskGigabytesLimit:     100,
			DiskVolumeCountLimit:   10,
			DiskSnapshotCountLimit: 10,
			PublicIPAddressLimit:   10,
			NetworkCountLimit:      10,
			SecurityGroupLimit:     10,
			SecurityGroupRuleLimit: 10,
		},
		InstanceSizes: []InstanceSize{
			{
				Name:          "g3.xsmall",
				CPUCores:      1,
				RAMMegabytes:  1024,
				DiskGigabytes: 10,
			},
			{
				Name:          "g3.small",
				CPUCores:      2,
				RAMMegabytes:  2048,
				DiskGigabytes: 20,
			},
			{
				Name:          "g3.medium",
				CPUCores:      4,
				RAMMegabytes:  4096,
				DiskGigabytes: 40,
			},
		},
		DiskImage: []DiskImage{
			{
				ID:           "b82168fe-66f6-4b38-a3b8-5283542d5475",
				Name:         "centos-7",
				Version:      "7",
				State:        "available",
				Distribution: "centos",
				Description:  "",
				Label:        "",
			},
			{
				ID:           "b82168fe-66f6-4b38-a3b8-52835425895",
				Name:         "debian-9",
				Version:      "9",
				State:        "available",
				Distribution: "debian",
				Description:  "",
				Label:        "",
			},
			{
				ID:           "b82168fe-66f6-4b38-a3b8-52835428965",
				Name:         "debian-10",
				Version:      "10",
				State:        "available",
				Distribution: "debian",
				Description:  "",
				Label:        "",
			},
			{
				ID:           "b82168fe-66f6-4b38-a3b8-528354282548",
				Name:         "ubuntu-20-4",
				Version:      "20.4",
				State:        "available",
				Distribution: "ubuntu",
				Description:  "",
				Label:        "",
			},
		},
	}, nil
}

// generateID generates a unique ID for use in the fake client
func (c *FakeClient) generateID() string {
	c.LastID++
	return strconv.FormatInt(c.LastID, 10)
}

// generatePublicIP generates a random IP address for use in the fake client
func (c *FakeClient) generatePublicIP() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return fmt.Sprintf("%v.%v.%v.%v", r.Intn(256), r.Intn(256), r.Intn(256), r.Intn(256))
}
