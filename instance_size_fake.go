package civogo

import (
	"fmt"
	"strings"
)

// ListInstanceSizes implemented in a fake way for automated tests
func (c *FakeClient) ListInstanceSizes() ([]InstanceSize, error) {
	return c.InstanceSizes, nil
}

// FindInstanceSizes implemented in a fake way for automated tests
func (c *FakeClient) FindInstanceSizes(search string) (*InstanceSize, error) {
	for _, size := range c.InstanceSizes {
		if strings.Contains(size.Name, search) {
			return &size, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}
