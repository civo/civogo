// Code generated from compute.openapi3.json (proto commit 60c7e49248e86df73908df1dde1c7e92ce59bdfd) by gencmd. DO NOT EDIT.

package compute

import (
	"github.com/civo/civogo/v3/api"
	"time"
)

// AssignIPToInstanceRequest — Request to assign a reserved IP to a compute instance.
type AssignIPToInstanceRequest struct {
	// Identifier for the reserved IP to assign. Accepts UUID, name, or label.
	ReservedIPID string `json:"reserved_ip_id,omitempty"`
}

// AttachVolumeToInstanceRequest — Request to attach a volume to a compute instance.
type AttachVolumeToInstanceRequest struct {
	// Whether to attach the volume at boot time.
	AttachAtBoot bool `json:"attach_at_boot,omitempty"`
	// Identifier for the volume to attach.
	VolumeID string `json:"volume_id,omitempty"`
}

// CreateInstanceConsoleRequest — Request to start a console session on a compute instance. An instance has at
// most one session: if one is already running, it is returned unchanged.
type CreateInstanceConsoleRequest struct {
	// How long the session lasts, between 5m and 24h. Defaults to 1h.
	Duration string `json:"duration,omitempty"`
}

// DetachVolumeFromInstanceRequest — Request to detach a volume from a compute instance.
type DetachVolumeFromInstanceRequest struct {
	// Identifier for the volume to detach.
	VolumeID string `json:"volume_id,omitempty"`
}

// EnableInstanceRecoveryModeRequest — Request to enable recovery mode on a compute instance.
type EnableInstanceRecoveryModeRequest struct {
	// Name of the recovery disk image to boot from. Defaults to the best match
	// for the instance's operating system.
	RecoveryImage *string `json:"recovery_image,omitempty"`
}

// ExpandInstanceRootDiskRequest — Request to expand the root disk of a compute instance.
type ExpandInstanceRootDiskRequest struct {
	// The desired root disk size in GB; the target is absolute, not a delta. Set
	// on the returned Instance as root_disk.target_size_gb; progress and outcome
	// are reported via root_disk.expansion. Expansion is grow-only: this must be
	// greater than the instance's current root_disk.current_size_gb (enforced by
	// the service/admission webhook, not by this proto). The upper bound below
	// is a sanity ceiling only; the authoritative maximum is dynamic and
	// enforced downstream.
	TargetSizeGB int `json:"target_size_gb,omitempty"`
}

// ResizeInstanceRequest — Request to resize a compute instance.
type ResizeInstanceRequest struct {
	// The target size (e.g. "g4s.large"). Must be a valid, selectable instance size.
	Size string `json:"size,omitempty"`
}

// UpdateInstanceRequest — Request to update an existing compute instance.
type UpdateInstanceRequest struct {
	// New allowed IPs for the instance.
	AllowedIps []string `json:"allowed_ips,omitempty"`
	// New firewall ID to apply.
	FirewallID *string `json:"firewall_id,omitempty"`
	// New hostname for the instance.
	Hostname *string `json:"hostname,omitempty"`
	// New network bandwidth limit.
	NetworkBandwidthLimit *string `json:"network_bandwidth_limit,omitempty"`
	// New notes for the instance.
	Notes         *string        `json:"notes,omitempty"`
	PlacementRule *PlacementRule `json:"placement_rule,omitempty"`
	// Whether a public IPv4 address is required.
	PublicIPRequired *bool `json:"public_ip_required,omitempty"`
	// Reserved IPv4 to assign (UUID, IP, or name). Empty string to unassign.
	ReservedIPv4 *string `json:"reserved_ipv4,omitempty"`
	// New reverse DNS entry for the instance.
	ReverseDNS *string `json:"reverse_dns,omitempty"`
	// New tags for the instance.
	Tags []string `json:"tags,omitempty"`
	// Set to true to apply the allowed_ips field (allows clearing with an empty list).
	UpdateAllowedIps bool `json:"update_allowed_ips,omitempty"`
	// Set to true to apply the tags field (allows clearing tags with an empty list).
	UpdateTags bool `json:"update_tags,omitempty"`
}

// AffinityRule — A scheduling rule that controls whether instances sharing tags should be
// co-located on the same host ("affinity") or separated across hosts
// ("anti-affinity").
type AffinityRule struct {
	// If true, the rule is a hard requirement; if false, it is treated as a preference.
	Exclusive bool `json:"exclusive,omitempty"`
	// Tags identifying the instances this rule applies to.
	Tags []string `json:"tags,omitempty"`
	// Rule type: "affinity" to co-locate, "anti-affinity" to separate.
	Type string `json:"type,omitempty"`
}

// AssignIPToInstanceResponse — Response for assigning a reserved IP to a compute instance.
type AssignIPToInstanceResponse struct {
}

// AttachVolumeToInstanceResponse — Response for attaching a volume to a compute instance.
type AttachVolumeToInstanceResponse struct {
}

// AttachedVolume — A volume currently attached to an instance.
type AttachedVolume struct {
	// Identifier of the attached volume.
	ID string `json:"id,omitempty"`
}

// CreateInstanceConsoleResponse — Response containing the console session.
type CreateInstanceConsoleResponse struct {
	Console *InstanceConsole `json:"console,omitempty"`
}

// CreateInstanceRequest — Request to create a new compute instance.
type CreateInstanceRequest struct {
	// Allowed IPs for the instance (requires advanced instance networking region feature).
	AllowedIps []string `json:"allowed_ips,omitempty"`
	// Cloud-init configuration to apply to the instance.
	CloudInit string `json:"cloud_init,omitempty"`
	// Number of instances to create (1-10). Optional; defaults to 1 when unset.
	// When count > 1, hostnames are suffixed: hostname-1, hostname-2, etc.
	// Reserved IP (reserved_ip_id) is not supported when count > 1.
	Count *int `json:"count,omitempty"`
	// The firewall ID to apply. Uses default firewall if empty.
	FirewallID string `json:"firewall_id,omitempty"`
	// The hostname for the instance.
	Hostname string `json:"hostname,omitempty"`
	// The initial password.
	InitialPassword string `json:"initial_password,omitempty"`
	// The initial SSH user.
	InitialUser string `json:"initial_user,omitempty"`
	// Whether an IPv6 address is required.
	IPv6Required bool `json:"ipv6_required,omitempty"`
	// Network bandwidth limit for the instance (requires advanced instance networking region feature).
	NetworkBandwidthLimit int `json:"network_bandwidth_limit,omitempty"`
	// The network ID to place the instance in.
	NetworkID string `json:"network_id,omitempty"`
	// Notes for the instance.
	Notes         string         `json:"notes,omitempty"`
	PlacementRule *PlacementRule `json:"placement_rule,omitempty"`
	// Whether a public IPv4 address is required.
	PublicIPRequired bool `json:"public_ip_required,omitempty"`
	// Reserved IP ID to assign.
	ReservedIPID string `json:"reserved_ip_id,omitempty"`
	// Reverse DNS entry.
	ReverseDNS string `json:"reverse_dns,omitempty"`
	// The size of the instance (e.g. "g4s.medium"). Defaults to "g4s.medium" if empty.
	Size string `json:"size,omitempty"`
	// The source image ID.
	SourceID string `json:"source_id,omitempty"`
	// The source type (e.g. "diskimage").
	SourceType string `json:"source_type,omitempty"`
	// The SSH public key content.
	SSHKey string `json:"ssh_key,omitempty"`
	// The SSH key ID reference.
	SSHKeyID string `json:"ssh_key_id,omitempty"`
	// Static private IPv4 address within the network CIDR.
	StaticPrivateIPv4 string `json:"static_private_ipv4,omitempty"`
	// Tags to apply to the instance.
	Tags []string `json:"tags,omitempty"`
	// The disk image identifier used to define the instance.
	TemplateID string `json:"template_id,omitempty"`
	// A user script to execute after cloud-init.
	UserScript string `json:"user_script,omitempty"`
	// CivoVolume IDs to attach at startup.
	VolumeIds []string `json:"volume_ids,omitempty"`
	// The volume type name (e.g. "standard", "encrypted"). Uses default resolution if empty.
	VolumeType string `json:"volume_type,omitempty"`
}

// CreateInstanceResponse — Response containing the created compute instance(s).
type CreateInstanceResponse struct {
	Instances []Instance `json:"instances,omitempty"`
}

// DeleteInstanceConsoleResponse — Response for an ended console session.
type DeleteInstanceConsoleResponse struct {
}

// DeleteInstanceResponse — Response containing the deleted compute instance.
type DeleteInstanceResponse struct {
	Instance *Instance `json:"instance,omitempty"`
}

// DetachVolumeFromInstanceResponse — Response for detaching a volume from a compute instance.
type DetachVolumeFromInstanceResponse struct {
}

// DisableInstanceRecoveryModeResponse — Response containing the recovery mode.
type DisableInstanceRecoveryModeResponse struct {
	RecoveryMode *InstanceRecoveryMode `json:"recovery_mode,omitempty"`
}

// EnableInstanceRecoveryModeResponse — Response containing the recovery mode.
type EnableInstanceRecoveryModeResponse struct {
	RecoveryMode *InstanceRecoveryMode `json:"recovery_mode,omitempty"`
}

// ExpandInstanceRootDiskResponse — Response containing the compute instance after a root-disk expansion was
// accepted. The expansion runs asynchronously: root_disk.target_size_gb is
// updated immediately; poll GetInstance and observe root_disk.expansion
// (status/reason) for progress and outcome. root_disk.current_size_gb
// converges to the target only on success.
type ExpandInstanceRootDiskResponse struct {
	Instance *Instance `json:"instance,omitempty"`
}

// GetInstanceConsoleResponse — Response containing the console session.
type GetInstanceConsoleResponse struct {
	Console *InstanceConsole `json:"console,omitempty"`
}

// GetInstanceRecoveryModeResponse — Response containing the recovery mode.
type GetInstanceRecoveryModeResponse struct {
	RecoveryMode *InstanceRecoveryMode `json:"recovery_mode,omitempty"`
}

// GetInstanceResponse — Response containing a single compute instance.
type GetInstanceResponse struct {
	Instance *Instance `json:"instance,omitempty"`
}

// HardRebootInstanceResponse — Response for an hard rebooted instance.
type HardRebootInstanceResponse struct {
}

// Instance — A compute instance and its current state.
type Instance struct {
	// Account that owns the instance.
	AccountID string `json:"account_id,omitempty"`
	// Source IPs or CIDR ranges the instance is permitted to send traffic from.
	AllowedIps []string `json:"allowed_ips,omitempty"`
	// Volumes currently attached to the instance.
	AttachedVolumes []AttachedVolume `json:"attached_volumes,omitempty"`
	// Latest metric snapshot reported by the instance, encoded as a comma-separated list.
	CivoStatsdStats string `json:"civo_statsd_stats,omitempty"`
	// Per-hour metric snapshots, ordered from oldest to newest.
	CivoStatsdStatsPerHour []string `json:"civo_statsd_stats_per_hour,omitempty"`
	// Per-minute metric snapshots, ordered from oldest to newest.
	CivoStatsdStatsPerMinute []string `json:"civo_statsd_stats_per_minute,omitempty"`
	// Per-instance token used to authenticate metric pushes from the instance's civostatsd agent. Server-generated at creation; treat as read-only.
	CivoStatsdToken string `json:"civo_statsd_token,omitempty"`
	// Identifier of the Kubernetes cluster the instance belongs to; only set for cluster nodes.
	ClusterID string `json:"cluster_id,omitempty"`
	// Number of vCPU cores allocated to the instance, derived from size.
	CPUCores int `json:"cpu_cores,omitempty"`
	// Time at which the instance was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Base root disk capacity allocated by the instance's size, in gigabytes; derived from size.
	// See root_disk.current_size_gb for the actual provisioned size, which may exceed this once the root disk has been expanded.
	DiskGB int `json:"disk_gb,omitempty"`
	// Distribution of the disk image used to provision the instance (e.g. "ubuntu").
	Distribution string `json:"distribution,omitempty"`
	// Firewall applied to the instance.
	FirewallID string `json:"firewall_id,omitempty"`
	// Label of the firewall applied to the instance.
	FirewallLabel string `json:"firewall_label,omitempty"`
	// Name of the firewall applied to the instance.
	FirewallName string `json:"firewall_name,omitempty"`
	// Number of GPUs allocated to the instance; 0 when none.
	GPUCount int `json:"gpu_count,omitempty"`
	// Model of the GPUs allocated to the instance, when gpu_count > 0.
	GPUType string `json:"gpu_type,omitempty"`
	// User-visible hostname assigned to the instance.
	Hostname string `json:"hostname,omitempty"`
	// Unique identifier of the instance.
	ID string `json:"id,omitempty"`
	// Password set for the initial user. Empty when authentication is via SSH key only.
	InitialPassword string `json:"initial_password,omitempty"`
	// Username of the initial account created on the instance.
	InitialUser string `json:"initial_user,omitempty"`
	// IPv6 address assigned to the instance, if IPv6 is enabled for the network.
	IPv6 string `json:"ipv6,omitempty"`
	// System-generated resource name for the instance.
	Name string `json:"name,omitempty"`
	// Identifier of the isolation boundary the instance lives in.
	NamespaceID string `json:"namespace_id,omitempty"`
	// Egress bandwidth cap in Mbps; 0 means uncapped.
	NetworkBandwidthLimit string `json:"network_bandwidth_limit,omitempty"`
	// Network the instance is attached to.
	NetworkID string `json:"network_id,omitempty"`
	// Label of the network the instance is attached to.
	NetworkLabel string `json:"network_label,omitempty"`
	// Name of the network the instance is attached to.
	NetworkName string `json:"network_name,omitempty"`
	// Free-form notes attached to the instance by the user.
	Notes         string         `json:"notes,omitempty"`
	PlacementRule *PlacementRule `json:"placement_rule,omitempty"`
	// Private IPv4 address assigned within the instance's network.
	PrivateIP string `json:"private_ip,omitempty"`
	// Public IPv4 address assigned to the instance, if any.
	PublicIP string `json:"public_ip,omitempty"`
	// Whether the instance was requested with a public IPv4 address.
	PublicIPRequired bool `json:"public_ip_required,omitempty"`
	// RAM allocated to the instance, in megabytes; derived from size.
	RAMMB int `json:"ram_mb,omitempty"`
	// Region the instance is located in.
	Region string `json:"region,omitempty"`
	// The reserved IP address currently assigned to the instance, if any.
	ReservedIP string `json:"reserved_ip,omitempty"`
	// Identifier of the reserved IP currently assigned to the instance, if any.
	ReservedIPID string `json:"reserved_ip_id,omitempty"`
	// Name of the reserved IP currently assigned to the instance, if any.
	ReservedIPName string `json:"reserved_ip_name,omitempty"`
	// Reverse-DNS hostname configured for the instance's public IP.
	ReverseDNS string    `json:"reverse_dns,omitempty"`
	RootDisk   *RootDisk `json:"root_disk,omitempty"`
	// Instance size identifier (e.g. "g4s.medium"); determines CPU, RAM and disk allocation.
	Size string `json:"size,omitempty"`
	// Human-readable label for the instance size (e.g. "Medium").
	SizeLabel string `json:"size_label,omitempty"`
	// Identifier of the source the instance was provisioned from, scoped to source_type.
	SourceID string `json:"source_id,omitempty"`
	// Source category for the boot image (e.g. "diskimage").
	SourceType string `json:"source_type,omitempty"`
	// SSH public key contents installed for the initial user.
	SSHKey string `json:"ssh_key,omitempty"`
	// Identifier of the SSH key resource installed on the instance, if one was supplied.
	SSHKeyID string         `json:"ssh_key_id,omitempty"`
	Status   InstanceStatus `json:"status,omitempty"`
	// Subnets the instance is connected to within its network.
	Subnets []Subnet `json:"subnets,omitempty"`
	// User-defined tags attached to the instance.
	Tags []string `json:"tags,omitempty"`
	// Identifier of the disk image used to provision the instance.
	TemplateID string `json:"template_id,omitempty"`
	// Classification of the instance: "compute" for standalone VMs, "cluster-master" or "cluster-worker" for Kubernetes cluster nodes.
	Type string `json:"type,omitempty"`
	// User-supplied script executed once on first boot.
	UserScript string `json:"user_script,omitempty"`
	// Storage class used for the instance's root disk.
	VolumeType string `json:"volume_type,omitempty"`
}

// InstanceConsole — A console (VNC) session on a compute instance.
type InstanceConsole struct {
	// When the session ends.
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	// URL to open the console in a browser.
	URI string `json:"uri,omitempty"`
}

// InstanceRecoveryMode — Recovery mode of a compute instance. While enabled, the instance boots from
// the recovery image instead of its root disk.
type InstanceRecoveryMode struct {
	Enabled bool `json:"enabled,omitempty"`
	// Name of the disk image the instance boots from in recovery mode.
	RecoveryImage string `json:"recovery_image,omitempty"`
}

// InstanceStatus — Lifecycle status of a compute instance.
//
//   - INSTANCE_STATUS_UNSPECIFIED: Default zero value; the status has not been set.
//   - INSTANCE_STATUS_ACTIVE: Instance is running and reachable.
//   - INSTANCE_STATUS_BUILDING: Instance is being provisioned for the first time.
//   - INSTANCE_STATUS_SHUTOFF: Instance has been powered off but its resources are retained and it can be started again.
//   - INSTANCE_STATUS_ERROR: Instance is in an error state; inspect the resource for details.
//   - INSTANCE_STATUS_REBOOTING: A reboot (soft or hard) is in progress.
//   - INSTANCE_STATUS_STARTING: Instance is powering up from SHUTOFF.
//   - INSTANCE_STATUS_RESCUE: Instance is in rescue mode for recovery/troubleshooting.
//   - INSTANCE_STATUS_UPGRADING: A resize is in progress (CPU/RAM/disk are being changed).
//   - INSTANCE_STATUS_DELETED: Instance has been deleted; its resources have been or are being released.
//   - INSTANCE_STATUS_MIGRATING: Instance is being live-migrated between hosts.
type InstanceStatus string

const (
	InstanceStatusUnspecified InstanceStatus = "INSTANCE_STATUS_UNSPECIFIED"
	InstanceStatusActive      InstanceStatus = "INSTANCE_STATUS_ACTIVE"
	InstanceStatusBuilding    InstanceStatus = "INSTANCE_STATUS_BUILDING"
	InstanceStatusShutoff     InstanceStatus = "INSTANCE_STATUS_SHUTOFF"
	InstanceStatusError       InstanceStatus = "INSTANCE_STATUS_ERROR"
	InstanceStatusRebooting   InstanceStatus = "INSTANCE_STATUS_REBOOTING"
	InstanceStatusStarting    InstanceStatus = "INSTANCE_STATUS_STARTING"
	InstanceStatusRescue      InstanceStatus = "INSTANCE_STATUS_RESCUE"
	InstanceStatusUpgrading   InstanceStatus = "INSTANCE_STATUS_UPGRADING"
	InstanceStatusDeleted     InstanceStatus = "INSTANCE_STATUS_DELETED"
	InstanceStatusMigrating   InstanceStatus = "INSTANCE_STATUS_MIGRATING"
)

// ListInstancesResponse — Response containing a list of compute instances.
type ListInstancesResponse struct {
	Instances []Instance   `json:"instances,omitempty"`
	PageInfo  api.PageInfo `json:"page_info,omitempty"`
}

// PlacementRule — Placement constraints applied when scheduling an instance onto a host.
type PlacementRule struct {
	// Affinity / anti-affinity rules to apply to this instance.
	AffinityRules []AffinityRule `json:"affinity_rules,omitempty"`
	// Required host attributes expressed as key/value pairs; an instance is
	// only scheduled onto hosts that match every entry.
	NodeSelector map[string]string `json:"node_selector,omitempty"`
}

// ResizeInstanceResponse — Response containing the resized compute instance.
type ResizeInstanceResponse struct {
	Instance *Instance `json:"instance,omitempty"`
}

// RootDisk describes the instance's root disk: its provisioned size, any
// requested expansion target, and the state of the most recent expansion.
//
// Always populated: every instance has a root disk, and every Instance
// returned by this service sets this field. An absent root_disk indicates a
// stale or faulty producer and must be treated as "no data", never as "never
// expanded". "Never expanded" is asserted only by a present root_disk whose
// expansion field is absent.
//
// The size-tier base allocation is Instance.disk_gb, which predates this
// message and remains top-level; current_size_gb may exceed disk_gb once the
// root disk has been expanded. Future root-disk attributes belong in this
// message.
//
// The instance's status remains INSTANCE_STATUS_ACTIVE throughout an
// expansion; expansion progress is reported here, not via Instance.status.
//
// expansion is the single authoritative progress and outcome signal; the size
// fields are plain data. Never infer expansion progress from a size
// comparison.
type RootDisk struct {
	// The provisioned size of the root disk in GB. Grows only when an expansion
	// completes; progress and outcome are reported via expansion.
	CurrentSizeGB int                `json:"current_size_gb,omitempty"`
	Expansion     *RootDiskExpansion `json:"expansion,omitempty"`
	// The last requested root disk size in GB. Not rolled back if an expansion
	// is rejected or times out, so it may exceed current_size_gb indefinitely
	// after the expansion settled.
	TargetSizeGB int `json:"target_size_gb,omitempty"`
}

// RootDiskExpansion reports the state of an online (non-disruptive) root-disk
// expansion. It is absent (root_disk.expansion unset) on an instance whose root
// disk has never been expanded. A present value carrying UNSPECIFIED
// status/reason means "no signal" and must not be read as success. The instance
// status stays INSTANCE_STATUS_ACTIVE throughout.
type RootDiskExpansion struct {
	// The last time status or reason changed.
	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`
	// Human-readable detail about the current status. Customer-facing: it never
	// carries internal identifiers (node names, hostnames, driver errors);
	// reason is the authoritative machine-readable signal.
	Message string                  `json:"message,omitempty"`
	Reason  RootDiskExpansionReason `json:"reason,omitempty"`
	Status  RootDiskExpansionStatus `json:"status,omitempty"`
	// The target size in GB this record pertains to. Compare with
	// root_disk.target_size_gb to detect a stale record immediately after a
	// re-target.
	TargetSizeGB int `json:"target_size_gb,omitempty"`
}

// RootDiskExpansionReason — Reason for the current root-disk expansion status. Clients must tolerate
// values not present in their generated set and classify by rule (status ==
// ROOT_DISK_EXPANSION_STATUS_SETTLED with reason != COMPLETED indicates a
// failure), not an exhaustive allow-list; new reasons may be added in this alpha.
//
//   - ROOT_DISK_EXPANSION_REASON_UNSPECIFIED: Default zero value; no reason reported. Never emitted on a populated
//
// RootDiskExpansion; treat as in-doubt, never as success.
//   - ROOT_DISK_EXPANSION_REASON_STARTED: Expansion has started and the disk is growing toward the target.
//
// (status IN_PROGRESS)
//   - ROOT_DISK_EXPANSION_REASON_TIMED_OUT: Expansion exceeded the operator's deadline and is no longer being driven.
//
// The underlying storage request cannot be withdrawn, so a later out-of-band
// completion flips reason to COMPLETED while status stays SETTLED.
// (status SETTLED)
//   - ROOT_DISK_EXPANSION_REASON_RETRYING: A transient error occurred and the operator is retrying. (status IN_PROGRESS)
//   - ROOT_DISK_EXPANSION_REASON_COMPLETED: Expansion completed successfully; the disk reached the target. (status SETTLED)
//   - ROOT_DISK_EXPANSION_REASON_EXCEEDS_MAX_SIZE: Rejected: the requested size exceeds the maximum provisionable size.
//
// (status SETTLED)
//   - ROOT_DISK_EXPANSION_REASON_NOT_PERMITTED_BY_STORAGE_CLASS: Rejected: the storage class does not permit expansion. (status SETTLED)
type RootDiskExpansionReason string

const (
	RootDiskExpansionReasonUnspecified                RootDiskExpansionReason = "ROOT_DISK_EXPANSION_REASON_UNSPECIFIED"
	RootDiskExpansionReasonStarted                    RootDiskExpansionReason = "ROOT_DISK_EXPANSION_REASON_STARTED"
	RootDiskExpansionReasonTimedOut                   RootDiskExpansionReason = "ROOT_DISK_EXPANSION_REASON_TIMED_OUT"
	RootDiskExpansionReasonRetrying                   RootDiskExpansionReason = "ROOT_DISK_EXPANSION_REASON_RETRYING"
	RootDiskExpansionReasonCompleted                  RootDiskExpansionReason = "ROOT_DISK_EXPANSION_REASON_COMPLETED"
	RootDiskExpansionReasonExceedsMaxSize             RootDiskExpansionReason = "ROOT_DISK_EXPANSION_REASON_EXCEEDS_MAX_SIZE"
	RootDiskExpansionReasonNotPermittedByStorageClass RootDiskExpansionReason = "ROOT_DISK_EXPANSION_REASON_NOT_PERMITTED_BY_STORAGE_CLASS"
)

// RootDiskExpansionStatus — Coarse state of an online root-disk expansion. A client can switch on this
// without knowing the full reason taxonomy; see RootDiskExpansionReason for detail.
//
//   - ROOT_DISK_EXPANSION_STATUS_UNSPECIFIED: Default zero value; no expansion state has been reported. Treat as "no
//
// signal" (in-doubt); never read as success.
//   - ROOT_DISK_EXPANSION_STATUS_IN_PROGRESS: An expansion is outstanding and the operator is driving the disk toward
//
// the requested size. See reason.
//   - ROOT_DISK_EXPANSION_STATUS_SETTLED: The operator is no longer driving the expansion: it completed successfully,
//
// was rejected, or timed out. See reason to distinguish success from failure.
type RootDiskExpansionStatus string

const (
	RootDiskExpansionStatusUnspecified RootDiskExpansionStatus = "ROOT_DISK_EXPANSION_STATUS_UNSPECIFIED"
	RootDiskExpansionStatusInProgress  RootDiskExpansionStatus = "ROOT_DISK_EXPANSION_STATUS_IN_PROGRESS"
	RootDiskExpansionStatusSettled     RootDiskExpansionStatus = "ROOT_DISK_EXPANSION_STATUS_SETTLED"
)

// SoftRebootInstanceResponse — Response for a soft rebooted instance.
type SoftRebootInstanceResponse struct {
}

// StartInstanceResponse — Response for a started instance.
type StartInstanceResponse struct {
}

// StopInstanceResponse — Response for a stopped instance.
type StopInstanceResponse struct {
}

// Subnet — A subnet allocated within a network and attached to an instance.
type Subnet struct {
	// System-generated subnet name.
	GeneratedName string `json:"generated_name,omitempty"`
	// Unique identifier of the subnet.
	ID string `json:"id,omitempty"`
	// User-provided label for the subnet.
	Name string `json:"name,omitempty"`
	// Identifier of the network this subnet belongs to.
	NetworkID string `json:"network_id,omitempty"`
	// Current status of the subnet (e.g. "active").
	Status string `json:"status,omitempty"`
	// Size of the subnet expressed as a CIDR suffix (e.g. "/24").
	SubnetSize string `json:"subnet_size,omitempty"`
}

// UnassignIPFromInstanceResponse — Response for unassigning a reserved IP from a compute instance.
type UnassignIPFromInstanceResponse struct {
}

// UpdateInstanceResponse — Response containing the updated compute instance.
type UpdateInstanceResponse struct {
	Instance *Instance `json:"instance,omitempty"`
}
