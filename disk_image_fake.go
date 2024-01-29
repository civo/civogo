package civogo

import (
	"fmt"
	"strings"
)

// ListDiskImages implemented in a fake way for automated tests
func (c *FakeClient) ListDiskImages() ([]DiskImage, error) {
	return c.DiskImage, nil
}

// GetDiskImage implemented in a fake way for automated tests
func (c *FakeClient) GetDiskImage(id string) (*DiskImage, error) {
	for k, v := range c.DiskImage {
		if v.ID == id {
			return &c.DiskImage[k], nil
		}
	}

	err := fmt.Errorf("unable to find disk image %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// FindDiskImage implemented in a fake way for automated tests
func (c *FakeClient) FindDiskImage(search string) (*DiskImage, error) {
	for _, diskimage := range c.DiskImage {
		if strings.Contains(diskimage.Name, search) || strings.Contains(diskimage.ID, search) {
			return &diskimage, nil
		}
	}

	err := fmt.Errorf("unable to find volume %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}
