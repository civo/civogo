package civogo

// GetQuota implemented in a fake way for automated tests
func (c *FakeClient) GetQuota() (*Quota, error) {
	return &c.Quota, nil
}
