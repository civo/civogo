package civogo

import (
	"fmt"
	"strings"
)

// ListVolumes implemented in a fake way for automated tests
func (c *FakeClient) ListVolumes() ([]Volume, error) {
	return c.Volumes, nil
}

// GetVolume implemented in a fake way for automated tests
func (c *FakeClient) GetVolume(id string) (*Volume, error) {
	for _, volume := range c.Volumes {
		if volume.ID == id {
			return &volume, nil
		}
	}

	err := fmt.Errorf("unable to get volume %s", id)
	return nil, ZeroMatchesError.wrap(err)
}

// FindVolume implemented in a fake way for automated tests
func (c *FakeClient) FindVolume(search string) (*Volume, error) {
	for _, volume := range c.Volumes {
		if strings.Contains(volume.Name, search) || strings.Contains(volume.ID, search) {
			return &volume, nil
		}
	}

	err := fmt.Errorf("unable to find volume %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// NewVolume implemented in a fake way for automated tests
func (c *FakeClient) NewVolume(v *VolumeConfig) (*VolumeResult, error) {
	volume := Volume{
		ID:            c.generateID(),
		Name:          v.Name,
		SizeGigabytes: v.SizeGigabytes,
		Status:        "available",
	}
	c.Volumes = append(c.Volumes, volume)

	return &VolumeResult{
		ID:     volume.ID,
		Name:   volume.Name,
		Result: "success",
	}, nil
}

// ResizeVolume implemented in a fake way for automated tests
func (c *FakeClient) ResizeVolume(id string, size int) (*SimpleResponse, error) {
	for i, volume := range c.Volumes {
		if volume.ID == id {
			c.Volumes[i].SizeGigabytes = size
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	err := fmt.Errorf("unable to find volume %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// AttachVolume implemented in a fake way for automated tests
func (c *FakeClient) AttachVolume(id string, instance string) (*SimpleResponse, error) {
	for i, volume := range c.Volumes {
		if volume.ID == id {
			c.Volumes[i].InstanceID = instance
			c.Volumes[i].Status = "attached"
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	err := fmt.Errorf("unable to find volume %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// DetachVolume implemented in a fake way for automated tests
func (c *FakeClient) DetachVolume(id string) (*SimpleResponse, error) {
	for i, volume := range c.Volumes {
		if volume.ID == id {
			c.Volumes[i].InstanceID = ""
			c.Volumes[i].Status = "available"
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	err := fmt.Errorf("unable to find volume %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// DeleteVolume implemented in a fake way for automated tests
func (c *FakeClient) DeleteVolume(id string) (*SimpleResponse, error) {
	for i, volume := range c.Volumes {
		if volume.ID == id {
			c.Volumes[len(c.Volumes)-1], c.Volumes[i] = c.Volumes[i], c.Volumes[len(c.Volumes)-1]
			c.Volumes = c.Volumes[:len(c.Volumes)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}
