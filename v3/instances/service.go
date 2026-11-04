// Code generated from compute.openapi3.json (proto commit 60c7e49248e86df73908df1dde1c7e92ce59bdfd) by gencmd. DO NOT EDIT.

package instances

import (
	"context"
	"fmt"
	"github.com/civo/civogo/v3/api"
	"iter"
	"net/url"
	"strings"
)

// Service provides access to the instances API.
type Service struct {
	t *api.Transport
}

// NewService is called by the root civogo package; use civogo.New.
func NewService(t *api.Transport) *Service { return &Service{t: t} }

// AssignIP calls POST /v3alpha1/compute/instances/{instance_id}/assign-ip.
// The API returns no payload for this call.
func (s *Service) AssignIP(ctx context.Context, instanceID string, req AssignIPToInstanceRequest) error {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/assign-ip", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, nil); err != nil {
		return err
	}
	return nil
}

// AttachVolume calls POST /v3alpha1/compute/instances/{instance_id}/attach-volume.
// The API returns no payload for this call.
func (s *Service) AttachVolume(ctx context.Context, instanceID string, req AttachVolumeToInstanceRequest) error {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/attach-volume", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, nil); err != nil {
		return err
	}
	return nil
}

// Create calls POST /v3alpha1/compute/instances.
func (s *Service) Create(ctx context.Context, req CreateInstanceRequest) ([]Instance, error) {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	var out struct {
		V []Instance `json:"instances"`
	}
	if err := s.t.Do(ctx, "POST", "/v3alpha1/compute/instances", nil, req, &out); err != nil {
		return nil, err
	}
	return out.V, nil
}

// CreateConsole calls PUT /v3alpha1/compute/instances/{instance_id}/console.
func (s *Service) CreateConsole(ctx context.Context, instanceID string, req CreateInstanceConsoleRequest) (*InstanceConsole, error) {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	var out struct {
		V InstanceConsole `json:"console"`
	}
	if err := s.t.Do(ctx, "PUT", strings.Replace("/v3alpha1/compute/instances/{instance_id}/console", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// Delete calls DELETE /v3alpha1/compute/instances/{instance_id}.
func (s *Service) Delete(ctx context.Context, instanceID string, opts ...api.CallOption) (*Instance, error) {
	q := api.ResolveCallParams(s.t, opts...).Query()
	var out struct {
		V Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "DELETE", strings.Replace("/v3alpha1/compute/instances/{instance_id}", "{instance_id}", url.PathEscape(instanceID), 1), q, nil, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// DeleteConsole calls DELETE /v3alpha1/compute/instances/{instance_id}/console.
// The API returns no payload for this call.
func (s *Service) DeleteConsole(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	q := api.ResolveCallParams(s.t, opts...).Query()
	if err := s.t.Do(ctx, "DELETE", strings.Replace("/v3alpha1/compute/instances/{instance_id}/console", "{instance_id}", url.PathEscape(instanceID), 1), q, nil, nil); err != nil {
		return err
	}
	return nil
}

// DetachVolume calls POST /v3alpha1/compute/instances/{instance_id}/detach-volume.
// The API returns no payload for this call.
func (s *Service) DetachVolume(ctx context.Context, instanceID string, req DetachVolumeFromInstanceRequest) error {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/detach-volume", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, nil); err != nil {
		return err
	}
	return nil
}

// DisableRecoveryMode calls DELETE /v3alpha1/compute/instances/{instance_id}/recovery.
func (s *Service) DisableRecoveryMode(ctx context.Context, instanceID string, opts ...api.CallOption) (*InstanceRecoveryMode, error) {
	q := api.ResolveCallParams(s.t, opts...).Query()
	var out struct {
		V InstanceRecoveryMode `json:"recovery_mode"`
	}
	if err := s.t.Do(ctx, "DELETE", strings.Replace("/v3alpha1/compute/instances/{instance_id}/recovery", "{instance_id}", url.PathEscape(instanceID), 1), q, nil, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// EnableRecoveryMode calls PUT /v3alpha1/compute/instances/{instance_id}/recovery.
func (s *Service) EnableRecoveryMode(ctx context.Context, instanceID string, req EnableInstanceRecoveryModeRequest) (*InstanceRecoveryMode, error) {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	var out struct {
		V InstanceRecoveryMode `json:"recovery_mode"`
	}
	if err := s.t.Do(ctx, "PUT", strings.Replace("/v3alpha1/compute/instances/{instance_id}/recovery", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// ExpandRootDisk calls POST /v3alpha1/compute/instances/{instance_id}/expand-root-disk.
func (s *Service) ExpandRootDisk(ctx context.Context, instanceID string, req ExpandInstanceRootDiskRequest) (*Instance, error) {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	var out struct {
		V Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/expand-root-disk", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// Get calls GET /v3alpha1/compute/instances/{instance_id}.
func (s *Service) Get(ctx context.Context, instanceID string, opts ...api.CallOption) (*Instance, error) {
	q := api.ResolveCallParams(s.t, opts...).Query()
	var out struct {
		V Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "GET", strings.Replace("/v3alpha1/compute/instances/{instance_id}", "{instance_id}", url.PathEscape(instanceID), 1), q, nil, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// GetConsole calls GET /v3alpha1/compute/instances/{instance_id}/console.
func (s *Service) GetConsole(ctx context.Context, instanceID string, opts ...api.CallOption) (*InstanceConsole, error) {
	q := api.ResolveCallParams(s.t, opts...).Query()
	var out struct {
		V InstanceConsole `json:"console"`
	}
	if err := s.t.Do(ctx, "GET", strings.Replace("/v3alpha1/compute/instances/{instance_id}/console", "{instance_id}", url.PathEscape(instanceID), 1), q, nil, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// GetRecoveryMode calls GET /v3alpha1/compute/instances/{instance_id}/recovery.
func (s *Service) GetRecoveryMode(ctx context.Context, instanceID string, opts ...api.CallOption) (*InstanceRecoveryMode, error) {
	q := api.ResolveCallParams(s.t, opts...).Query()
	var out struct {
		V InstanceRecoveryMode `json:"recovery_mode"`
	}
	if err := s.t.Do(ctx, "GET", strings.Replace("/v3alpha1/compute/instances/{instance_id}/recovery", "{instance_id}", url.PathEscape(instanceID), 1), q, nil, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// HardReboot calls POST /v3alpha1/compute/instances/{instance_id}/reboot_hard.
// The API returns no payload for this call.
func (s *Service) HardReboot(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	body := api.ResolveCallParams(s.t, opts...)
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/reboot_hard", "{instance_id}", url.PathEscape(instanceID), 1), nil, body, nil); err != nil {
		return err
	}
	return nil
}

// ListRequest selects what List returns.
type ListRequest struct {
	//
	AccountID string
	// Regions to list instances from. At least one region is required.
	Regions []string
	// Page selects one page; the zero value uses server defaults.
	Page api.PageParams
}

// List calls GET /v3alpha1/compute/instances, returning one page.
func (s *Service) List(ctx context.Context, req ListRequest) (*ListInstancesResponse, error) {
	q := url.Values{}
	regions := req.Regions
	if len(regions) == 0 && s.t.DefaultRegion() != "" {
		regions = []string{s.t.DefaultRegion()}
	}
	q["regions"] = regions
	if req.AccountID != "" {
		q.Set("account_id", req.AccountID)
	}
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
func (s *Service) All(ctx context.Context, req ListRequest) iter.Seq2[Instance, error] {
	return api.All(func(p api.PageParams) ([]Instance, *api.PageInfo, error) {
		req.Page = p
		resp, err := s.List(ctx, req)
		if err != nil {
			return nil, nil, err
		}
		return resp.Instances, &resp.PageInfo, nil
	})
}

// Resize calls POST /v3alpha1/compute/instances/{instance_id}/resize.
func (s *Service) Resize(ctx context.Context, instanceID string, req ResizeInstanceRequest) (*Instance, error) {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	var out struct {
		V Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/resize", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}

// SoftReboot calls POST /v3alpha1/compute/instances/{instance_id}/reboot_soft.
// The API returns no payload for this call.
func (s *Service) SoftReboot(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	body := api.ResolveCallParams(s.t, opts...)
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/reboot_soft", "{instance_id}", url.PathEscape(instanceID), 1), nil, body, nil); err != nil {
		return err
	}
	return nil
}

// Start calls POST /v3alpha1/compute/instances/{instance_id}/start.
// The API returns no payload for this call.
func (s *Service) Start(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	body := api.ResolveCallParams(s.t, opts...)
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/start", "{instance_id}", url.PathEscape(instanceID), 1), nil, body, nil); err != nil {
		return err
	}
	return nil
}

// Stop calls POST /v3alpha1/compute/instances/{instance_id}/stop.
// The API returns no payload for this call.
func (s *Service) Stop(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	body := api.ResolveCallParams(s.t, opts...)
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/stop", "{instance_id}", url.PathEscape(instanceID), 1), nil, body, nil); err != nil {
		return err
	}
	return nil
}

// UnassignIP calls POST /v3alpha1/compute/instances/{instance_id}/unassign-ip.
// The API returns no payload for this call.
func (s *Service) UnassignIP(ctx context.Context, instanceID string, opts ...api.CallOption) error {
	q := api.ResolveCallParams(s.t, opts...).Query()
	if err := s.t.Do(ctx, "POST", strings.Replace("/v3alpha1/compute/instances/{instance_id}/unassign-ip", "{instance_id}", url.PathEscape(instanceID), 1), q, nil, nil); err != nil {
		return err
	}
	return nil
}

// Update calls PUT /v3alpha1/compute/instances/{instance_id}.
func (s *Service) Update(ctx context.Context, instanceID string, req UpdateInstanceRequest) (*Instance, error) {
	if req.Region == "" {
		req.Region = s.t.DefaultRegion()
	}
	var out struct {
		V Instance `json:"instance"`
	}
	if err := s.t.Do(ctx, "PUT", strings.Replace("/v3alpha1/compute/instances/{instance_id}", "{instance_id}", url.PathEscape(instanceID), 1), nil, req, &out); err != nil {
		return nil, err
	}
	return &out.V, nil
}
