package civogo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Errors raised by package civogo
var (
	ResponseDecodeFailedError = constError("ResponseDecodeFailed")
	DisabledServiceError      = constError("DisabledServiceError")

	CivoStatsdRecordFailedError = constError("CivoStatsdRecordFailedError")
	AuthenticationFailedError   = constError("AuthenticationFailedError")

	// Volume Error
	CannotRescueNewVolumeError              = constError("CannotRescueNewVolumeError")
	CannotRestoreNewVolumeError             = constError("CannotRestoreNewVolumeError")
	CannotScaleAlreadyRescalingClusterError = constError("CannotScaleAlreadyRescalingClusterError")
	VolumeInvalidSizeError                  = constError("VolumeInvalidSizeError")

	DatabaseAccountDestroyError      = constError("DatabaseAccountDestroyError")
	DatabaseAccountNotFoundError     = constError("DatabaseAccountNotFoundError")
	DatabaseAccountAccessDeniedError = constError("DatabaseAccountAccessDeniedError")
	DatabaseCreatingAccountError     = constError("DatabaseCreatingAccountError")

	DatabaseUpdatingAccountError = constError("DatabaseUpdatingAccountError")
	DatabaseAccountStatsError    = constError("DatabaseAccountStatsError")
	DatabaseActionListingError   = constError("DatabaseActionListingError")
	DatabaseAPIKeyCreateError    = constError("DatabaseApiKeyCreateError")
	DatabaseAPIKeyDuplicateError = constError("DatabaseApiKeyDuplicateError")
	DatabaseAPIKeyNotFoundError  = constError("DatabaseApiKeyNotFoundError")

	DatabaseAPIkeyDestroyError           = constError("DatabaseAPIkeyDestroyError")
	DatabaseAuditLogListingError         = constError("DatabaseAuditLogListingError")
	DatabaseBlueprintNotFoundError       = constError("DatabaseBlueprintNotFoundError")
	DatabaseBlueprintDeleteFailedError   = constError("DatabaseBlueprintDeleteFailedError")
	DatabaseBlueprintCreateError         = constError("DatabaseBlueprintCreateError")
	DatabaseBlueprintUpdateError         = constError("DatabaseBlueprintUpdateError")
	ParameterEmptyVolumeIDError          = constError("ParameterEmptyVolumeIDError")
	ParameterEmptyOpenstackVolumeIDError = constError("ParameterEmptyOpenstackVolumeIDError")
	DatabaseChangeAPIKeyError            = constError("DatabaseChangeAPIKeyError")
	DatabaseChargeListingError           = constError("DatabaseChargeListingError")
	DatabaseConnectionFailedError        = constError("DatabaseConnectionFailedError")

	DatabaseDNSDomainCreateError        = constError("DatabaseDnsDomainCreateError")
	DatabaseDNSDomainUpdateError        = constError("DatabaseDnsDomainUpdateError")
	DatabaseDNSDomainDuplicateNameError = constError("DatabaseDnsDomainDuplicateNameError")
	DatabaseDNSDomainNotFoundError      = constError("DatabaseDNSDomainNotFoundError")
	DatabaseDNSRecordCreateError        = constError("DatabaseDNSRecordCreateError")
	DatabaseDNSRecordNotFoundError      = constError("DatabaseDNSRecordNotFoundError")
	DatabaseDNSRecordUpdateError        = constError("DatabaseDNSRecordUpdateError")
	DatabaseListingDNSDomainsError      = constError("DatabaseListingDNSDomainsError")

	DatabaseFirewallCreateError           = constError("DatabaseFirewallCreateError")
	DatabaseFirewallDuplicateNameError    = constError("DatabaseFirewallDuplicateNameError")
	DatabaseFirewallMismatchError         = constError("DatabaseFirewallMismatchError")
	DatabaseFirewallNotFoundError         = constError("DatabaseFirewallNotFoundError")
	DatabaseFirewallSaveFailedError       = constError("DatabaseFirewallSaveFailedError")
	DatabaseFirewallDeleteFailedError     = constError("DatabaseFirewallDeleteFailedError")
	DatabaseFirewallRuleCreateError       = constError("DatabaseFirewallRuleCreateError")
	DatabaseFirewallRuleDeleteFailedError = constError("DatabaseFirewallRuleDeleteFailedError")
	DatabaseFirewallRulesFindError        = constError("DatabaseFirewallRulesFindError")
	DatabaseListingFirewallsError         = constError("DatabaseListingFirewallsError")
	FirewallDuplicateError                = constError("FirewallDuplicateError")

	// Instances Errors
	DatabaseInstanceAlreadyinRescueStateError              = constError("DatabaseInstanceAlreadyinRescueStateError")
	DatabaseInstanceBuildError                             = constError("DatabaseInstanceBuildError")
	DatabaseInstanceBuildMultipleWithExistingPublicIPError = constError("DatabaseInstanceBuildMultipleWithExistingPublicIPError")
	DatabaseInstanceCreateError                            = constError("DatabaseInstanceCreateError")
	DatabaseInstanceSnapshotTooBigError                    = constError("DatabaseInstanceSnapshotTooBigError")
	DatabaseInstanceDuplicateError                         = constError("DatabaseInstanceDuplicateError")
	DatabaseInstanceDuplicateNameError                     = constError("DatabaseInstanceDuplicateNameError")
	DatabaseInstanceListError                              = constError("DatabaseInstanceListError")
	DatabaseInstanceNotFoundError                          = constError("DatabaseInstanceNotFoundError")
	DatabaseInstanceNotInOpenStackError                    = constError("DatabaseInstanceNotInOpenStackError")
	DatabaseCannotManageClusterInstanceError               = constError("DatabaseCannotManageClusterInstanceError")
	DatabaseOldInstanceFindError                           = constError("DatabaseOldInstanceFindError")
	DatabaseCannotMoveIPError                              = constError("DatabaseCannotMoveIPError")
	DatabaseIPFindError                                    = constError("DatabaseIPFindError")

	// Kubernetes Errors
	DatabaseKubernetesClusterInvalidError         = constError("DatabaseKubernetesClusterInvalid")
	DatabaseKubernetesApplicationNotFoundError    = constError("DatabaseKubernetesApplicationNotFound")
	DatabaseKubernetesApplicationInvalidPlanError = constError("DatabaseKubernetesApplicationInvalidPlan")
	DatabaseKubernetesClusterDuplicateError       = constError("DatabaseKubernetesClusterDuplicate")
	DatabaseKubernetesClusterNotFoundError        = constError("DatabaseKubernetesClusterNotFound")
	DatabaseKubernetesNodeNotFoundError           = constError("DatabaseKubernetesNodeNotFound")

	DatabaseListingAccountsError              = constError("DatabaseListingAccountsError")
	DatabaseListingMembershipsError           = constError("DatabaseListingMembershipsError")
	DatabaseMembershipCannotDeleteError       = constError("DatabaseMembershipCannotDeleteError")
	DatabaseMembershipsGrantAccessError       = constError("DatabaseMembershipsGrantAccessError")
	DatabaseMembershipsInvalidInvitationError = constError("DatabaseMembershipsInvalidInvitationError")
	DatabaseMembershipsInvalidStatusError     = constError("DatabaseMembershipsInvalidStatusError")
	DatabaseMembershipsNotFoundError          = constError("DatabaseMembershipsNotFoundError")
	DatabaseMembershipsSuspendedError         = constError("DatabaseMembershipsSuspendedError")

	DatabaseLoadbalancerDuplicateError = constError("DatabaseLoadbalancerDuplicateError")
	DatabaseLoadbalancerInvalidError   = constError("DatabaseLoadbalancerInvalidError")
	DatabaseLoadbalancerNotFoundError  = constError("DatabaseLoadbalancerNotFoundError")

	DatabaseNetworksListError              = constError("DatabaseNetworksListError")
	DatabaseNetworkCreateError             = constError("DatabaseNetworkCreateError")
	DatabaseNetworkExistsError             = constError("DatabaseNetworkExistsError")
	DatabaseNetworkDeleteLastError         = constError("DatabaseNetworkDeleteLastError")
	DatabaseNetworkDeleteWithInstanceError = constError("DatabaseNetworkDeleteWithInstanceError")
	DatabaseNetworkDuplicateNameError      = constError("DatabaseNetworkDuplicateNameError")
	DatabaseNetworkLookupError             = constError("DatabaseNetworkLookupError")
	DatabaseNetworkNotFoundError           = constError("DatabaseNetworkNotFoundError")
	DatabaseNetworkSaveError               = constError("DatabaseNetworkSaveError")

	DatabasePrivateIPFromPublicIPError = constError("DatabasePrivateIPFromPublicIPError")

	DatabaseQuotaNotFoundError   = constError("DatabaseQuotaNotFoundError")
	DatabaseQuotaUpdateError     = constError("DatabaseQuotaUpdateError")
	QuotaLimitReachedError       = constError("QuotaLimitReachedError")
	DatabaseServiceNotFoundError = constError("DatabaseServiceNotFoundError")
	DatabaseSizesListError       = constError("DatabaseSizesListError")

	DatabaseSnapshotCannotDeleteInUseError      = constError("DatabaseSnapshotCannotDeleteInUseError")
	DatabaseSnapshotCannotReplaceError          = constError("DatabaseSnapshotCannotReplaceError")
	DatabaseSnapshotCreateError                 = constError("DatabaseSnapshotCreateError")
	DatabaseSnapshotCreateInstanceNotFoundError = constError("DatabaseSnapshotCreateInstanceNotFoundError")
	DatabaseSnapshotCreateAlreadyInProcessError = constError("DatabaseSnapshotCreateAlreadyInProcessError")
	DatabaseSnapshotNotFoundError               = constError("DatabaseSnapshotNotFoundError")
	DatabaseSnapshotsListError                  = constError("DatabaseSnapshotsListError")

	DatabaseSSHKeyDestroyError       = constError("DatabaseSSHKeyDestroyError")
	DatabaseSSHKeyCreateError        = constError("DatabaseSSHKeyCreateError")
	DatabaseSSHKeyUpdateError        = constError("DatabaseSSHKeyUpdateError")
	DatabaseSSHKeyDuplicateNameError = constError("DatabaseSSHKeyDuplicateNameError")
	DatabaseSSHKeyNotFoundError      = constError("DatabaseSSHKeyNotFoundError")
	SSHKeyDuplicateError             = constError("SSHKeyDuplicateError")

	DatabaseTeamCannotDeleteError      = constError("DatabaseTeamCannotDeleteError")
	DatabaseTeamCreateError            = constError("DatabaseTeamCreateError")
	DatabaseTeamListingError           = constError("DatabaseTeamListingError")
	DatabaseTeamMembershipCreateError  = constError("DatabaseTeamMembershipCreateError")
	DatabaseTeamNotFoundError          = constError("DatabaseTeamNotFoundError")
	DatabaseTemplateDestroyError       = constError("DatabaseTemplateDestroyError")
	DatabaseTemplateNotFoundError      = constError("DatabaseTemplateNotFoundError")
	DatabaseTemplateUpdateError        = constError("DatabaseTemplateUpdateError")
	DatabaseTemplateWouldConflictError = constError("DatabaseTemplateWouldConflictError")
	DatabaseImageIDInvalidError        = constError("DatabaseImageIDInvalidError")

	DatabaseUserAlreadyExistsError          = constError("DatabaseUserAlreadyExistsError")
	DatabaseUserNewError                    = constError("DatabaseUserNewError")
	DatabaseUserConfirmedError              = constError("DatabaseUserConfirmedError")
	DatabaseUserSuspendedError              = constError("DatabaseUserSuspendedError")
	DatabaseUserLoginFailedError            = constError("DatabaseUserLoginFailedError")
	DatabaseUserNoChangeStatusError         = constError("DatabaseUserNoChangeStatusError")
	DatabaseUserNotFoundError               = constError("DatabaseUserNotFoundError")
	DatabaseUserPasswordInvalidError        = constError("DatabaseUserPasswordInvalidError")
	DatabaseUserPasswordSecuringFailedError = constError("DatabaseUserPasswordSecuringFailedError")
	DatabaseUserUpdateError                 = constError("DatabaseUserUpdateError")
	DatabaseCreatingUserError               = constError("DatabaseCreatingUserError")

	DatabaseVolumeIDInvalidError                 = constError("DatabaseVolumeIDInvalidError")
	DatabaseVolumeDuplicateNameError             = constError("DatabaseVolumeDuplicateNameError")
	DatabaseVolumeCannotMultipleAttachError      = constError("DatabaseVolumeCannotMultipleAttachError")
	DatabaseVolumeStillAttachedCannotResizeError = constError("DatabaseVolumeStillAttachedCannotResizeError")
	DatabaseVolumeNotAttachedError               = constError("DatabaseVolumeNotAttachedError")
	DatabaseVolumeNotFoundError                  = constError("DatabaseVolumeNotFoundError")
	DatabaseVolumeDeleteFailedError              = constError("DatabaseVolumeDeleteFailedError")

	DatabaseWebhookDestroyError       = constError("DatabaseWebhookDestroyError")
	DatabaseWebhookNotFoundError      = constError("DatabaseWebhookNotFoundError")
	DatabaseWebhookUpdateError        = constError("DatabaseWebhookUpdateError")
	DatabaseWebhookWouldConflictError = constError("DatabaseWebhookWouldConflictError")

	OpenstackConnectionFailedError        = constError("OpenstackConnectionFailedError")
	OpenstackCreatingProjectError         = constError("OpenstackCreatingProjectError")
	OpenstackCreatingUserError            = constError("OpenstackCreatingUserError")
	OpenstackFirewallCreateError          = constError("OpenstackFirewallCreateError")
	OpenstackFirewallDestroyError         = constError("OpenstackFirewallDestroyError")
	OpenstackFirewallRuleDestroyError     = constError("OpenstackFirewallRuleDestroyError")
	OpenstackInstanceCreateError          = constError("OpenstackInstanceCreateError")
	OpenstackInstanceDestroyError         = constError("OpenstackInstanceDestroyError")
	OpenstackInstanceFindError            = constError("OpenstackInstanceFindError")
	OpenstackInstanceRebootError          = constError("OpenstackInstanceRebootError")
	OpenstackInstanceRebuildError         = constError("OpenstackInstanceRebuildError")
	OpenstackInstanceResizeError          = constError("OpenstackInstanceResizeError")
	OpenstackInstanceRestoreError         = constError("OpenstackInstanceRestoreError")
	OpenstackInstanceSetFirewallError     = constError("OpenstackInstanceSetFirewallError")
	OpenstackInstanceStartError           = constError("OpenstackInstanceStartError")
	OpenstackInstanceStopError            = constError("OpenstackInstanceStopError")
	OpenstackIPCreateError                = constError("OpenstackIPCreateError")
	OpenstackNetworkCreateFailedError     = constError("OpenstackNetworkCreateFailedError")
	OpenstackNnetworkDestroyFailedError   = constError("OpenstackNnetworkDestroyFailedError")
	OpenstackNetworkEnsureConfiguredError = constError("OpenstackNetworkEnsureConfiguredError")
	OpenstackPublicIPConnectError         = constError("OpenstackPublicIPConnectError")
	OpenstackQuotaApplyError              = constError("OpenstackQuotaApplyError")
	OpenstackSnapshotDestroyError         = constError("OpenstackSnapshotDestroyError")
	OpenstackSSHKeyUploadError            = constError("OpenstackSSHKeyUploadError")
	OpenstackProjectDestroyError          = constError("OpenstackProjectDestroyError")
	OpenstackProjectFindoError            = constError("OpenstackProjectFindoError")
	OpenstackUserDestroyError             = constError("OpenstackUserDestroyError")
	OpenstackURLGlanceError               = constError("OpenstackUrlGlanceError")
	OpenstackURLNovaError                 = constError("OpenstackURLNovaError")
	AuthenticationInvalidKeyError         = constError("AuthenticationInvalidKeyError")
	AuthenticationAccessDeniedError       = constError("AuthenticationAccessDeniedError")

	InstanceStateMustBeActiveOrShutoffError = constError("InstanceStateMustBeActiveOrShutoffError")
	MarshalingObjectsToJSONError            = constError("MarshalingObjectsToJsonError")
	NetworkCreateDefaultError               = constError("NetworkCreateDefaultError")
	NetworkDeleteDefaultError               = constError("NetworkDeleteDefaultError")
	ParameterTimeValueError                 = constError("ParameterTimeValueError")
	ParameterDateRangeTooLongError          = constError("ParameterDateRangeTooLongError")
	ParameterDNSRecordTypeError             = constError("ParameterDnsRecordTypeError")
	ParameterDNSRecordCnameApexError        = constError("ParameterDNSRecordCnameApexError")
	ParameterPublicKeyEmptyError            = constError("ParameterPublicKeyEmptyError")
	ParameterDateRangeError                 = constError("ParameterDateRangeError")
	ParameterIDMissingError                 = constError("ParameterIDMissingError")
	ParameterIDToIntegerError               = constError("ParameterIDToIntegerError")
	ParameterImageAndVolumeIDMissingError   = constError("ParameterImageAndVolumeIDMissingError")
	ParameterLabelInvalidError              = constError("ParameterLabelInvalidError")
	ParameterNameInvalidError               = constError("ParameterNameInvalidError")
	ParameterPrivateIPMissingError          = constError("ParameterPrivateIPMissingError")
	ParameterPublicIPMissingError           = constError("ParameterPublicIPMissingError")
	ParameterSizeMissingError               = constError("ParameterSizeMissingError")
	ParameterVolumeSizeIncorrectError       = constError("ParameterVolumeSizeIncorrectError")
	ParameterVolumeSizeMustIncreaseError    = constError("ParameterVolumeSizeMustIncreaseError")
	ParameterSnapshotMissingError           = constError("ParameterSnapshotMissingError")
	ParameterSnapshotIncorrectFormatError   = constError("ParameterSnapshotIncorrectFormatError")
	ParameterStartPortMissingError          = constError("ParameterStartPortMissingError")
	DatabaseTemplateParseRequestError       = constError("DatabaseTemplateParseRequestError")
	ParameterValueMissingError              = constError("ParameterValueMissingError")

	UnknowError = constError("Unknow Error")
)

type constError string

func (err constError) Error() string {
	return string(err)
}

func (err constError) Is(target error) bool {
	ts := target.Error()
	es := string(err)
	return ts == es || strings.HasPrefix(ts, es+": ")
}

func (err constError) wrap(inner error) error {
	return wrapError{msg: string(err), err: inner}
}

type wrapError struct {
	err error
	msg string
}

func (err wrapError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.msg, err.err)
	}
	return err.msg
}

func (err wrapError) Unwrap() error {
	return err.err
}

func (err wrapError) Is(target error) bool {
	return constError(err.msg).Is(target)
}

func decodeERROR(err error) error {
	errorData := err.(HTTPError)
	byt := []byte(errorData.Reason)

	var dat map[string]interface{}
	var msg strings.Builder

	if err := json.Unmarshal(byt, &dat); err != nil {
		err := errors.New("Failed to decode the response expected from the API")
		return ResponseDecodeFailedError.wrap(err)
	}

	if _, ok := dat["reason"]; ok {
		msg.WriteString(dat["reason"].(string))
		if _, ok := dat["details"]; ok {
			msg.WriteString(", " + dat["details"].(string))
		}
	}

	switch dat["code"] {
	// Kuberenetes Error
	case "database_kubernetes_cluster_invalid":
		err := errors.New(msg.String())
		return DatabaseKubernetesClusterInvalidError.wrap(err)
	case "database_kubernetes_application_not_found":
		err := errors.New(msg.String())
		return DatabaseKubernetesApplicationNotFoundError.wrap(err)
	case "database_kubernetes_application_invalid_plan":
		err := errors.New(msg.String())
		return DatabaseKubernetesApplicationInvalidPlanError.wrap(err)
	case "database_kubernetes_cluster_duplicate":
		err := errors.New(msg.String())
		return DatabaseKubernetesClusterDuplicateError.wrap(err)
	case "database_kubernetes_cluster_not_found":
		err := errors.New(msg.String())
		return DatabaseKubernetesClusterNotFoundError.wrap(err)
	case "database_kubernetes_node_not_found":
		err := errors.New(msg.String())
		return DatabaseKubernetesNodeNotFoundError.wrap(err)
	// Instances Error
	case "database_instance_already_in_rescue_state":
		err := errors.New(msg.String())
		return DatabaseInstanceAlreadyinRescueStateError.wrap(err)
	case "database_instance_build":
		err := errors.New(msg.String())
		return DatabaseInstanceBuildError.wrap(err)
	case "database_instance_build_multiple_with_existing_public_ip":
		err := errors.New(msg.String())
		return DatabaseInstanceBuildMultipleWithExistingPublicIPError.wrap(err)
	case "database_instance_create":
		err := errors.New(msg.String())
		return DatabaseInstanceCreateError.wrap(err)
	case "database_instance_snapshot_too_big":
		err := errors.New(msg.String())
		return DatabaseInstanceSnapshotTooBigError.wrap(err)
	case "instance_duplicate":
		err := errors.New(msg.String())
		return DatabaseInstanceDuplicateError.wrap(err)
	case "database_instance_duplicate_name":
		err := errors.New(msg.String())
		return DatabaseInstanceDuplicateNameError.wrap(err)
	case "database_instance_list":
		err := errors.New(msg.String())
		return DatabaseInstanceListError.wrap(err)
	case "database_instance_find":
		err := errors.New(msg.String())
		return DatabaseInstanceNotFoundError.wrap(err)
	case "database_instance_not_in_openstack":
		err := errors.New(msg.String())
		return DatabaseInstanceNotInOpenStackError.wrap(err)
	default:
		return UnknowError
	}
}
