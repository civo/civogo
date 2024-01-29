package civogo

// ListPermissions implemented in a fake way for automated tests
func (c *FakeClient) ListPermissions() ([]Permission, error) {
	return []Permission{
		{
			Name:        "instance.create",
			Description: "Create Compute instances",
		},
		{
			Name:        "kubernetes.*",
			Description: "Manage Civo Kubernetes clusters",
		},
	}, nil
}
