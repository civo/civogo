package civogo

import "time"

// ListCharges implemented in a fake way for automated tests
func (c *FakeClient) ListCharges(from, to time.Time) ([]Charge, error) {
	return []Charge{}, nil
}
