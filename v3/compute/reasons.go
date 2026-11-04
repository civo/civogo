// Code generated from compute.openapi3.json (proto commit 60c7e49248e86df73908df1dde1c7e92ce59bdfd) by gencmd. DO NOT EDIT.

// Package compute holds the error reasons the Compute v3 API reports.
//
// ErrorReason names the reasons the compute services report in the
// google.rpc.ErrorInfo detail of a failed RPC. The value name is the reason
// string on the wire, and it is shared by ComputeService, SSHKeyService,
// SizeService, AdminQuotaService, InstanceSnapshotService and
// SnapshotScheduleService: a condition that means the same thing in two
// services carries one reason.
//
// A request that breaks a validation rule reports REQUEST_VALIDATION_FAILED,
// with one google.rpc.BadRequest field violation per offending field, whose
// field is the field's path in the request message. A caller sees one contract
// whether the rejection came from a validation rule or from a check written by
// hand, so the other values below name conditions beyond the request's own
// shape.
//
// Failures that are the service's own fault reach the client masked, as a
// generic internal error with reason INTERNAL, or UNKNOWN when the service did
// not classify the error.
package compute

// Reasons reported in Error.Reason; match them with civogo.HasReason.
// The API may add reasons, so never treat this set as exhaustive.
const (
	// ReasonInstanceNotFound — No instance with the requested identifier exists in the account.
	ReasonInstanceNotFound = "INSTANCE_NOT_FOUND"
	// ReasonNetworkNotFound — No network with the requested identifier exists in the account.
	ReasonNetworkNotFound = "NETWORK_NOT_FOUND"
	// ReasonDefaultNetworkNotFound — The account has no default network in the region to place the instance in.
	ReasonDefaultNetworkNotFound = "DEFAULT_NETWORK_NOT_FOUND"
	// ReasonFirewallNotFound — No firewall with the requested identifier exists in the account.
	ReasonFirewallNotFound = "FIREWALL_NOT_FOUND"
	// ReasonReservedIPNotFound — No reserved IP with the requested identifier exists in the account.
	ReasonReservedIPNotFound = "RESERVED_IP_NOT_FOUND"
	// ReasonVolumeNotFound — No volume with the requested identifier exists in the account.
	ReasonVolumeNotFound = "VOLUME_NOT_FOUND"
	// ReasonVolumeTypeNotFound — No volume type with the requested name exists in the region.
	ReasonVolumeTypeNotFound = "VOLUME_TYPE_NOT_FOUND"
	// ReasonDiskImageNotFound — No disk image with the requested identifier exists.
	ReasonDiskImageNotFound = "DISK_IMAGE_NOT_FOUND"
	// ReasonSSHKeyNotFound — No SSH key with the requested identifier exists in the account.
	ReasonSSHKeyNotFound = "SSH_KEY_NOT_FOUND"
	// ReasonSizeNotFound — No size with the requested name exists in the region.
	ReasonSizeNotFound = "SIZE_NOT_FOUND"
	// ReasonSSHKeyAlreadyExists — The account already has an SSH key with the requested name.
	ReasonSSHKeyAlreadyExists = "SSH_KEY_ALREADY_EXISTS"
	// ReasonRebootAlreadyInProgress — A reboot of the instance is already in progress.
	ReasonRebootAlreadyInProgress = "REBOOT_ALREADY_IN_PROGRESS"
	// ReasonInstanceNotResizable — The instance is in a state that does not allow resizing.
	ReasonInstanceNotResizable = "INSTANCE_NOT_RESIZABLE"
	// ReasonInstanceNotExpandable — The instance is in a state that does not allow root disk expansion.
	ReasonInstanceNotExpandable = "INSTANCE_NOT_EXPANDABLE"
	// ReasonInstanceReservedIPAlreadyAssigned — The instance already has a reserved IP assigned.
	ReasonInstanceReservedIPAlreadyAssigned = "INSTANCE_RESERVED_IP_ALREADY_ASSIGNED"
	// ReasonInstanceReservedIPNotAssigned — The instance has no reserved IP assigned, so none can be unassigned.
	ReasonInstanceReservedIPNotAssigned = "INSTANCE_RESERVED_IP_NOT_ASSIGNED"
	// ReasonReservedIPAlreadyAssigned — The reserved IP is assigned to a different resource than the one the
	// request names.
	ReasonReservedIPAlreadyAssigned = "RESERVED_IP_ALREADY_ASSIGNED"
	// ReasonPrivateIPAlreadyAllocated — The requested static private IP is already allocated in the network.
	ReasonPrivateIPAlreadyAllocated = "PRIVATE_IP_ALREADY_ALLOCATED"
	// ReasonNetworkAddressesExhausted — The network has no free addresses left to give the instance.
	ReasonNetworkAddressesExhausted = "NETWORK_ADDRESSES_EXHAUSTED"
	// ReasonFirewallNetworkMismatch — The firewall is not in the instance's network.
	ReasonFirewallNetworkMismatch = "FIREWALL_NETWORK_MISMATCH"
	// ReasonVolumeNotAvailable — The volume is not in a state that allows attaching it.
	ReasonVolumeNotAvailable = "VOLUME_NOT_AVAILABLE"
	// ReasonVolumeNetworkMismatch — The volume is not in the instance's network.
	ReasonVolumeNetworkMismatch = "VOLUME_NETWORK_MISMATCH"
	// ReasonGPUSizeReserved — The requested GPU size is reserved and the account is not authorised for it.
	ReasonGPUSizeReserved = "GPU_SIZE_RESERVED"
	// ReasonGPUNotAvailable — The region has too few GPU cards free for the requested size.
	ReasonGPUNotAvailable = "GPU_NOT_AVAILABLE"
	// ReasonRootDiskExpansionNotSupported — The instance's storage class does not allow root disk expansion.
	ReasonRootDiskExpansionNotSupported = "ROOT_DISK_EXPANSION_NOT_SUPPORTED"
	// ReasonRootDiskExpansionInProgress — A root disk expansion to a different size is already in progress.
	ReasonRootDiskExpansionInProgress = "ROOT_DISK_EXPANSION_IN_PROGRESS"
	// ReasonRootDiskStorageInsufficient — The region has too little storage free to expand the root disk.
	ReasonRootDiskStorageInsufficient = "ROOT_DISK_STORAGE_INSUFFICIENT"
	// ReasonRootDiskExpansionTimedOut — The root disk expansion did not complete in time.
	ReasonRootDiskExpansionTimedOut = "ROOT_DISK_EXPANSION_TIMED_OUT"
	// ReasonRegionUnavailable — The requested region is not serving requests.
	ReasonRegionUnavailable = "REGION_UNAVAILABLE"
	// ReasonRegionFeatureUnavailable — The requested region does not offer the feature the request needs.
	ReasonRegionFeatureUnavailable = "REGION_FEATURE_UNAVAILABLE"
	// ReasonAccountFeatureNotEnabled — The account is not enabled for the feature the request needs.
	ReasonAccountFeatureNotEnabled = "ACCOUNT_FEATURE_NOT_ENABLED"
	// ReasonInstanceSnapshotNotFound — No instance snapshot with the requested identifier exists in the account.
	ReasonInstanceSnapshotNotFound = "INSTANCE_SNAPSHOT_NOT_FOUND"
	// ReasonInstanceSnapshotAlreadyExists — The account already has an instance snapshot with the requested name.
	ReasonInstanceSnapshotAlreadyExists = "INSTANCE_SNAPSHOT_ALREADY_EXISTS"
	// ReasonInstanceSnapshotRestoreNotFound — No restore with the requested identifier exists for the snapshot.
	ReasonInstanceSnapshotRestoreNotFound = "INSTANCE_SNAPSHOT_RESTORE_NOT_FOUND"
	// ReasonInstanceSnapshotNotSupported — The instance does not support snapshots.
	ReasonInstanceSnapshotNotSupported = "INSTANCE_SNAPSHOT_NOT_SUPPORTED"
	// ReasonInstanceSnapshotNotComplete — The snapshot has not completed, so it cannot be restored.
	ReasonInstanceSnapshotNotComplete = "INSTANCE_SNAPSHOT_NOT_COMPLETE"
	// ReasonSnapshotScheduleNotFound — No snapshot schedule with the requested identifier exists in the account.
	ReasonSnapshotScheduleNotFound = "SNAPSHOT_SCHEDULE_NOT_FOUND"
	// ReasonSnapshotScheduleAlreadyExists — The account already has a snapshot schedule with the requested name.
	ReasonSnapshotScheduleAlreadyExists = "SNAPSHOT_SCHEDULE_ALREADY_EXISTS"
	// ReasonSnapshotScheduleNetworkMismatch — The instances in the schedule are not all in the same network.
	ReasonSnapshotScheduleNetworkMismatch = "SNAPSHOT_SCHEDULE_NETWORK_MISMATCH"
	// ReasonInstanceConsoleNotFound — The instance has no console session.
	ReasonInstanceConsoleNotFound = "INSTANCE_CONSOLE_NOT_FOUND"
	// ReasonInstanceConsoleNotAvailable — The instance is not in a state that allows a console session.
	ReasonInstanceConsoleNotAvailable = "INSTANCE_CONSOLE_NOT_AVAILABLE"
	// ReasonRecoveryImageNotAvailable — No recovery image is available for the instance's operating system.
	ReasonRecoveryImageNotAvailable = "RECOVERY_IMAGE_NOT_AVAILABLE"
	// ReasonRequestValidationFailed — The request broke a validation rule. The google.rpc.BadRequest detail
	// names each offending field.
	ReasonRequestValidationFailed = "REQUEST_VALIDATION_FAILED"
	// ReasonInternal — The service failed on its own; the cause is masked from the client.
	ReasonInternal = "INTERNAL"
	// ReasonUnknown — The service failed with an error it did not classify; the cause is masked
	// from the client.
	ReasonUnknown = "UNKNOWN"
)
