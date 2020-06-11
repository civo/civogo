package civogo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Errors raised by package civogo
var (
	ResponseDecodeFailedError = constError("Failed to decode the simple response expected from the API")
	ListInstanceSizesError    = constError("Failed to get the instances size list from the API")
	KubernetesClusterError    = constError("KubernetesClusterError")

	// Kubernetes Errors
	DatabaseKubernetesClusterInvalidError         = constError("DatabaseKubernetesClusterInvalid")
	DatabaseKubernetesApplicationNotFoundError    = constError("DatabaseKubernetesApplicationNotFound")
	DatabaseKubernetesApplicationInvalidPlanError = constError("DatabaseKubernetesApplicationInvalidPlan")
	DatabaseKubernetesClusterDuplicateError       = constError("DatabaseKubernetesClusterDuplicate")
	DatabaseKubernetesClusterNotFoundError        = constError("DatabaseKubernetesClusterNotFound")
	DatabaseKubernetesNodeNotFoundError           = constError("DatabaseKubernetesNodeNotFound")

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
		panic(err)
	}

	if _, ok := dat["reason"]; ok {
		msg.WriteString(dat["reason"].(string))
		if _, ok := dat["details"]; ok {
			msg.WriteString(", " + dat["details"].(string))
		}
	}

	switch dat["code"] {
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
	}

	return UnknowError.wrap(err)
}
