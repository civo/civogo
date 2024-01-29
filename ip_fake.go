package civogo

// ListIPs returns a list of fake IPs
func (c *FakeClient) ListIPs() (*PaginatedIPs, error) {
	return &PaginatedIPs{
		Page:    1,
		PerPage: 20,
		Pages:   100,
		Items: []IP{
			{
				ID:   c.generateID(),
				Name: "test-ip",
				IP:   c.generatePublicIP(),
			},
		},
	}, nil
}

// GetIP returns a fake IP
func (c *FakeClient) GetIP(id string) (*IP, error) {
	return &IP{
		ID:   c.generateID(),
		Name: "test-ip",
		IP:   c.generatePublicIP(),
	}, nil
}

// FindIP finds a fake IP
func (c *FakeClient) FindIP(search string) (*IP, error) {
	return &IP{
		ID:   c.generateID(),
		Name: "test-ip",
		IP:   c.generatePublicIP(),
	}, nil
}

// NewIP creates a fake IP
func (c *FakeClient) NewIP(v *CreateIPRequest) (*IP, error) {
	return &IP{
		ID:   c.generateID(),
		Name: "test-ip",
		IP:   c.generatePublicIP(),
	}, nil
}

// UpdateIP updates a fake IP
func (c *FakeClient) UpdateIP(id string, v *UpdateIPRequest) (*IP, error) {
	return &IP{
		ID:   c.generateID(),
		Name: v.Name,
		IP:   c.generatePublicIP(),
	}, nil
}

// DeleteIP deletes a fake IP
func (c *FakeClient) DeleteIP(id string) (*SimpleResponse, error) {
	return &SimpleResponse{
		Result: "success",
	}, nil
}

// AssignIP assigns a fake IP
func (c *FakeClient) AssignIP(id, resourceID, resourceType, region string) (*SimpleResponse, error) {
	return &SimpleResponse{
		Result: "success",
	}, nil
}

// UnassignIP unassigns a fake IP
func (c *FakeClient) UnassignIP(id, region string) (*SimpleResponse, error) {
	return &SimpleResponse{
		Result: "success",
	}, nil
}
