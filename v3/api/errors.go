package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Canonical google.rpc.Code names reported in Error.Status.
const (
	StatusInvalidArgument   = "INVALID_ARGUMENT"
	StatusNotFound          = "NOT_FOUND"
	StatusAlreadyExists     = "ALREADY_EXISTS"
	StatusPermissionDenied  = "PERMISSION_DENIED"
	StatusResourceExhausted = "RESOURCE_EXHAUSTED"
	StatusUnauthenticated   = "UNAUTHENTICATED"
)

// ReasonValidationFailed is the reason of every request rejected for
// breaking a validation rule; FieldViolations names the offending fields.
const ReasonValidationFailed = "REQUEST_VALIDATION_FAILED"

// Error is the single error type returned for any v3 API failure.
// Branch on Reason: it is the stable, machine-readable condition.
type Error struct {
	HTTPStatus int
	// Status is the canonical gRPC code name, e.g. "NOT_FOUND".
	Status  string
	Message string

	// Reason is the google.rpc.ErrorInfo reason, e.g. "INSTANCE_NOT_FOUND".
	Reason string
	// Domain is the service that failed, e.g. "compute.civo.com".
	Domain    string
	Metadata  map[string]string
	RequestID string

	FieldViolations []FieldViolation
	QuotaViolations []QuotaViolation

	// Details holds every detail object of the response, including any
	// the fields above do not model.
	Details []json.RawMessage
}

// FieldViolation is one entry of a google.rpc.BadRequest detail.
type FieldViolation struct {
	// Field is the path of the offending field in the request, e.g. "account_id".
	Field       string `json:"field"`
	Description string `json:"description"`
}

// QuotaViolation is one entry of a google.rpc.QuotaFailure detail.
type QuotaViolation struct {
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

func (e *Error) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "civo: %s (http %d", e.Message, e.HTTPStatus)
	if e.Reason != "" {
		fmt.Fprintf(&b, ", reason %s", e.Reason)
	}
	if e.RequestID != "" {
		fmt.Fprintf(&b, ", request %s", e.RequestID)
	}
	b.WriteString(")")
	for _, v := range e.FieldViolations {
		fmt.Fprintf(&b, "; %s: %s", v.Field, v.Description)
	}
	return b.String()
}

// HasReason reports whether err is a v3 API error with the given reason.
func HasReason(err error, reason string) bool {
	var apiErr *Error
	return errors.As(err, &apiErr) && apiErr.Reason == reason
}

// IsNotFound reports whether err is a v3 API "not found" error.
func IsNotFound(err error) bool { return hasStatus(err, StatusNotFound, 404) }

// IsAlreadyExists reports whether err is a v3 API "already exists" error.
func IsAlreadyExists(err error) bool { return hasStatus(err, StatusAlreadyExists, 409) }

// IsAuthFailure reports whether err is an authentication or authorization failure.
func IsAuthFailure(err error) bool {
	return hasStatus(err, StatusUnauthenticated, 401) || hasStatus(err, StatusPermissionDenied, 403)
}

// IsRateLimited reports whether err is a rate-limit / quota-exhausted error.
func IsRateLimited(err error) bool { return hasStatus(err, StatusResourceExhausted, 429) }

// IsInvalidArgument reports whether err is a request-validation error.
func IsInvalidArgument(err error) bool { return hasStatus(err, StatusInvalidArgument, 400) }

func hasStatus(err error, status string, httpStatus int) bool {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.Status == status || apiErr.HTTPStatus == httpStatus
}

// errorEnvelope is the AIP-193 error body the v3 gateway writes.
type errorEnvelope struct {
	Error *struct {
		Code    int               `json:"code"` // HTTP status, not the gRPC code
		Status  string            `json:"status"`
		Message string            `json:"message"`
		Details []json.RawMessage `json:"details"`
	} `json:"error"`
}

// detail is the union of the google.rpc detail types, discriminated by @type.
type detail struct {
	Type            string            `json:"@type"`
	Reason          string            `json:"reason"`
	Domain          string            `json:"domain"`
	Metadata        map[string]string `json:"metadata"`
	FieldViolations []FieldViolation  `json:"fieldViolations"`
	Violations      []QuotaViolation  `json:"violations"`
}

const typePrefix = "type.googleapis.com/google.rpc."

func decodeError(httpStatus int, body []byte) *Error {
	var env errorEnvelope
	if err := json.Unmarshal(body, &env); err != nil || env.Error == nil {
		// Not a structured envelope (e.g. a proxy error page):
		// keep the raw body as the message so nothing is lost.
		return &Error{HTTPStatus: httpStatus, Message: strings.TrimSpace(string(body))}
	}

	e := &Error{
		HTTPStatus: httpStatus,
		Status:     env.Error.Status,
		Message:    env.Error.Message,
		Details:    env.Error.Details,
	}
	for _, raw := range env.Error.Details {
		var d detail
		if json.Unmarshal(raw, &d) != nil {
			continue
		}
		switch d.Type {
		case typePrefix + "ErrorInfo":
			e.Reason, e.Domain, e.Metadata = d.Reason, d.Domain, d.Metadata
			e.RequestID = d.Metadata["request_id"]
		case typePrefix + "BadRequest":
			e.FieldViolations = append(e.FieldViolations, d.FieldViolations...)
		case typePrefix + "QuotaFailure":
			e.QuotaViolations = append(e.QuotaViolations, d.Violations...)
		}
	}
	return e
}
