package civogo

// ListRegions implemented in a fake way for automated tests
func (c *FakeClient) ListRegions() ([]Region, error) {
	return []Region{
		{
			Code:    "FAKE1",
			Name:    "Fake testing region",
			Default: true,
		},
	}, nil
}
