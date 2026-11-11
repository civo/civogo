// Code generated from compute.openapi3.json (proto commit 60c7e49248e86df73908df1dde1c7e92ce59bdfd) by gencmd. DO NOT EDIT.

package compute

import (
	"context"
	"errors"
	"fmt"
	"github.com/civo/civogo/v3/api"
	"github.com/civo/civogo/v3/internal/transport"
	"iter"
	"net/url"
	"strings"
)

// InstancesService is the Compute instances API.
type InstancesService struct {
	t *transport.Transport
}

// AssignIP calls POST /v3alpha1/compute/instances/{instance_id}/assign-ip.
// The API returns no payload for this call.
func (s *InstancesService) AssignIP(ctx context.Context, instanceID string, req AssignIPToInstanceRequest, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/assign-ip", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		AssignIPToInstanceRequest
		transport.Scope
	}{req, sc}, nil)
}

// AttachVolume calls POST /v3alpha1/compute/instances/{instance_id}/attach-volume.
// The API returns no payload for this call.
func (s *InstancesService) AttachVolume(ctx context.Context, instanceID string, req AttachVolumeToInstanceRequest, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/attach-volume", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		AttachVolumeToInstanceRequest
		transport.Scope
	}{req, sc}, nil)
}

// Create calls POST /v3alpha1/compute/instances.
func (s *InstancesService) Create(ctx context.Context, req CreateInstanceRequest, opts ...api.CallOption) ([]Instance, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V []Instance `json:"instances"`
	}
	if err := s.t.Do(ctx, "POST", "/v3alpha1/compute/instances", nil, struct {
		CreateInstanceRequest
		transport.Scope
	}{req, sc}, &out); err != nil {
		return nil, err
	}
	return out.V, nil
}

// CreateConsole calls PUT /v3alpha1/compute/instances/{instance_id}/console.
func (s *InstancesService) CreateConsole(ctx context.Context, instanceID string, req CreateInstanceConsoleRequest, opts ...api.CallOption) (*InstanceConsole, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *InstanceConsole `json:"console"`
	}
	if err := s.t.Do(ctx, "PUT", strings.Replace("/v3alpha1/compute/instances/{instance_id}/console", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		CreateInstanceConsoleRequest
		transport.Scope
	}{req, sc}, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: CreateConsole response has no \"console\" field")
	}
	return out.V, nil
}

// Delete calls DELETE /v3alpha1/compute/instances/{instance_id}.
func (s *InstancesService) Delete(ctx context.Context, instanceID string, opts ...api.CallOption) (*Instance, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "DELETE", strings.Replace("/v3alpha1/compute/instances/{instance_id}", "{instance_id}", url.PathEscape(instanceID), 1), sc.Query(), nil, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: Delete response has no \"instance\" field")
	}
	return out.V, nil
}

// DeleteConsole calls DELETE /v3alpha1/compute/instances/{instance_id}/console.
// The API returns no payload for this call.
func (s *InstancesService) DeleteConsole(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "DELETE", strings.Replace("/v3alpha1/compute/instances/{instance_id}/console", "{instance_id}", url.PathEscape(instanceID), 1), sc.Query(), nil, nil)
}

// DetachVolume calls POST /v3alpha1/compute/instances/{instance_id}/detach-volume.
// The API returns no payload for this call.
func (s *InstancesService) DetachVolume(ctx context.Context, instanceID string, req DetachVolumeFromInstanceRequest, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/detach-volume", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		DetachVolumeFromInstanceRequest
		transport.Scope
	}{req, sc}, nil)
}

// DisableRecoveryMode calls DELETE /v3alpha1/compute/instances/{instance_id}/recovery.
func (s *InstancesService) DisableRecoveryMode(ctx context.Context, instanceID string, opts ...api.CallOption) (*InstanceRecoveryMode, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *InstanceRecoveryMode `json:"recovery_mode"`
	}
	if err := s.t.Do(ctx, "DELETE", strings.Replace("/v3alpha1/compute/instances/{instance_id}/recovery", "{instance_id}", url.PathEscape(instanceID), 1), sc.Query(), nil, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: DisableRecoveryMode response has no \"recovery_mode\" field")
	}
	return out.V, nil
}

// EnableRecoveryMode calls PUT /v3alpha1/compute/instances/{instance_id}/recovery.
func (s *InstancesService) EnableRecoveryMode(ctx context.Context, instanceID string, req EnableInstanceRecoveryModeRequest, opts ...api.CallOption) (*InstanceRecoveryMode, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *InstanceRecoveryMode `json:"recovery_mode"`
	}
	if err := s.t.Do(ctx, "PUT", strings.Replace("/v3alpha1/compute/instances/{instance_id}/recovery", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		EnableInstanceRecoveryModeRequest
		transport.Scope
	}{req, sc}, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: EnableRecoveryMode response has no \"recovery_mode\" field")
	}
	return out.V, nil
}

// ExpandRootDisk calls POST /v3alpha1/compute/instances/{instance_id}/expand-root-disk.
func (s *InstancesService) ExpandRootDisk(ctx context.Context, instanceID string, req ExpandInstanceRootDiskRequest, opts ...api.CallOption) (*Instance, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/expand-root-disk", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		ExpandInstanceRootDiskRequest
		transport.Scope
	}{req, sc}, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: ExpandRootDisk response has no \"instance\" field")
	}
	return out.V, nil
}

// Get calls GET /v3alpha1/compute/instances/{instance_id}.
func (s *InstancesService) Get(ctx context.Context, instanceID string, opts ...api.CallOption) (*Instance, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "GET", strings.Replace("/v3alpha1/compute/instances/{instance_id}", "{instance_id}", url.PathEscape(instanceID), 1), sc.Query(), nil, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: Get response has no \"instance\" field")
	}
	return out.V, nil
}

// GetConsole calls GET /v3alpha1/compute/instances/{instance_id}/console.
func (s *InstancesService) GetConsole(ctx context.Context, instanceID string, opts ...api.CallOption) (*InstanceConsole, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *InstanceConsole `json:"console"`
	}
	if err := s.t.Do(ctx, "GET", strings.Replace("/v3alpha1/compute/instances/{instance_id}/console", "{instance_id}", url.PathEscape(instanceID), 1), sc.Query(), nil, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: GetConsole response has no \"console\" field")
	}
	return out.V, nil
}

// GetRecoveryMode calls GET /v3alpha1/compute/instances/{instance_id}/recovery.
func (s *InstancesService) GetRecoveryMode(ctx context.Context, instanceID string, opts ...api.CallOption) (*InstanceRecoveryMode, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *InstanceRecoveryMode `json:"recovery_mode"`
	}
	if err := s.t.Do(ctx, "GET", strings.Replace("/v3alpha1/compute/instances/{instance_id}/recovery", "{instance_id}", url.PathEscape(instanceID), 1), sc.Query(), nil, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: GetRecoveryMode response has no \"recovery_mode\" field")
	}
	return out.V, nil
}

// HardReboot calls POST /v3alpha1/compute/instances/{instance_id}/reboot_hard.
// The API returns no payload for this call.
func (s *InstancesService) HardReboot(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/reboot_hard", "{instance_id}", url.PathEscape(instanceID), 1), nil, sc, nil)
}

// ListInstancesRequest selects what List returns.
type ListInstancesRequest struct {
	// Page selects one page; the zero value uses server defaults.
	Page api.PageParams
}

// List calls GET /v3alpha1/compute/instances, returning one page. Pass api.WithRegion once per
// region to list from several; with none, the client's default is used.
func (s *InstancesService) List(ctx context.Context, req ListInstancesRequest, opts ...api.CallOption) (*ListInstancesResponse, error) {
	q := transport.ListQuery(s.t, opts)
	if req.Page.Page > 0 {
		q.Set("pagination.page", fmt.Sprint(req.Page.Page))
	}
	if req.Page.PerPage > 0 {
		q.Set("pagination.per_page", fmt.Sprint(req.Page.PerPage))
	}
	var out ListInstancesResponse
	if err := s.t.Do(ctx, "GET", "/v3alpha1/compute/instances", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// All lazily iterates every Instance matching req across all pages.
func (s *InstancesService) All(ctx context.Context, req ListInstancesRequest, opts ...api.CallOption) iter.Seq2[Instance, error] {
	return transport.All(func(p api.PageParams) ([]Instance, *api.PageInfo, error) {
		r := req
		r.Page = p
		resp, err := s.List(ctx, r, opts...)
		if err != nil {
			return nil, nil, err
		}
		return resp.Instances, &resp.PageInfo, nil
	})
}

// Resize calls POST /v3alpha1/compute/instances/{instance_id}/resize.
func (s *InstancesService) Resize(ctx context.Context, instanceID string, req ResizeInstanceRequest, opts ...api.CallOption) (*Instance, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/resize", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		ResizeInstanceRequest
		transport.Scope
	}{req, sc}, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: Resize response has no \"instance\" field")
	}
	return out.V, nil
}

// SoftReboot calls POST /v3alpha1/compute/instances/{instance_id}/reboot_soft.
// The API returns no payload for this call.
func (s *InstancesService) SoftReboot(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/reboot_soft", "{instance_id}", url.PathEscape(instanceID), 1), nil, sc, nil)
}

// Start calls POST /v3alpha1/compute/instances/{instance_id}/start.
// The API returns no payload for this call.
func (s *InstancesService) Start(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/start", "{instance_id}", url.PathEscape(instanceID), 1), nil, sc, nil)
}

// Stop calls POST /v3alpha1/compute/instances/{instance_id}/stop.
// The API returns no payload for this call.
func (s *InstancesService) Stop(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/stop", "{instance_id}", url.PathEscape(instanceID), 1), nil, sc, nil)
}

// UnassignIP calls POST /v3alpha1/compute/instances/{instance_id}/unassign-ip.
// The API returns no payload for this call.
func (s *InstancesService) UnassignIP(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return err
	}
	return s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/unassign-ip", "{instance_id}", url.PathEscape(instanceID), 1), sc.Query(), nil, nil)
}

// Update calls PUT /v3alpha1/compute/instances/{instance_id}.
func (s *InstancesService) Update(ctx context.Context, instanceID string, req UpdateInstanceRequest, opts ...api.CallOption) (*Instance, error) {
	sc, err := transport.Resolve(s.t, opts)
	if err != nil {
		return nil, err
	}
	var out struct {
		V *Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "PUT", strings.Replace("/v3alpha1/compute/instances/{instance_id}", "{instance_id}", url.PathEscape(instanceID), 1), nil, struct {
		UpdateInstanceRequest
		transport.Scope
	}{req, sc}, &out); err != nil {
		return nil, err
	}
	if out.V == nil {
		return nil, errors.New("civo: Update response has no \"instance\" field")
	}
	return out.V, nil
}
