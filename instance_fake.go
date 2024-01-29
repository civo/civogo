package civogo

import (
	"fmt"
	"strings"
)

// ListInstances implemented in a fake way for automated tests
func (c *FakeClient) ListInstances(page int, perPage int) (*PaginatedInstanceList, error) {
	return &PaginatedInstanceList{
		Items:   c.Instances,
		Page:    page,
		PerPage: perPage,
		Pages:   page,
	}, nil
}

// ListAllInstances implemented in a fake way for automated tests
func (c *FakeClient) ListAllInstances() ([]Instance, error) {
	return c.Instances, nil
}

// FindInstance implemented in a fake way for automated tests
func (c *FakeClient) FindInstance(search string) (*Instance, error) {
	for _, instance := range c.Instances {
		if strings.Contains(instance.Hostname, search) {
			return &instance, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// GetInstance implemented in a fake way for automated tests
func (c *FakeClient) GetInstance(id string) (*Instance, error) {
	for _, instance := range c.Instances {
		if instance.ID == id {
			return &instance, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// NewInstanceConfig implemented in a fake way for automated tests
func (c *FakeClient) NewInstanceConfig() (*InstanceConfig, error) {
	return &InstanceConfig{}, nil
}

// CreateInstance implemented in a fake way for automated tests
func (c *FakeClient) CreateInstance(config *InstanceConfig) (*Instance, error) {
	instance := Instance{
		ID:          c.generateID(),
		Hostname:    config.Hostname,
		Size:        config.Size,
		Region:      config.Region,
		TemplateID:  config.TemplateID,
		InitialUser: config.InitialUser,
		SSHKey:      config.SSHKeyID,
		Tags:        config.Tags,
		PublicIP:    c.generatePublicIP(),
	}
	c.Instances = append(c.Instances, instance)
	return &instance, nil
}

// SetInstanceTags implemented in a fake way for automated tests
func (c *FakeClient) SetInstanceTags(i *Instance, tags string) (*SimpleResponse, error) {
	for idx, instance := range c.Instances {
		if instance.ID == i.ID {
			c.Instances[idx].Tags = strings.Split(tags, " ")
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}

// UpdateInstance implemented in a fake way for automated tests
func (c *FakeClient) UpdateInstance(i *Instance) (*SimpleResponse, error) {
	for idx, instance := range c.Instances {
		if instance.ID == i.ID {
			c.Instances[idx] = *i
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}

// DeleteInstance implemented in a fake way for automated tests
func (c *FakeClient) DeleteInstance(id string) (*SimpleResponse, error) {
	for i, instance := range c.Instances {
		if instance.ID == id {
			c.Instances[len(c.Instances)-1], c.Instances[i] = c.Instances[i], c.Instances[len(c.Instances)-1]
			c.Instances = c.Instances[:len(c.Instances)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}

// RebootInstance implemented in a fake way for automated tests
func (c *FakeClient) RebootInstance(id string) (*SimpleResponse, error) {
	return &SimpleResponse{Result: "success"}, nil
}

// HardRebootInstance implemented in a fake way for automated tests
func (c *FakeClient) HardRebootInstance(id string) (*SimpleResponse, error) {
	return &SimpleResponse{Result: "success"}, nil
}

// SoftRebootInstance implemented in a fake way for automated tests
func (c *FakeClient) SoftRebootInstance(id string) (*SimpleResponse, error) {
	return &SimpleResponse{Result: "success"}, nil
}

// StopInstance implemented in a fake way for automated tests
func (c *FakeClient) StopInstance(id string) (*SimpleResponse, error) {
	return &SimpleResponse{Result: "success"}, nil
}

// StartInstance implemented in a fake way for automated tests
func (c *FakeClient) StartInstance(id string) (*SimpleResponse, error) {
	return &SimpleResponse{Result: "success"}, nil
}

// GetInstanceConsoleURL implemented in a fake way for automated tests
func (c *FakeClient) GetInstanceConsoleURL(id string) (string, error) {
	return fmt.Sprintf("https://console.example.com/%s", id), nil
}

// UpgradeInstance implemented in a fake way for automated tests
func (c *FakeClient) UpgradeInstance(id, newSize string) (*SimpleResponse, error) {
	for idx, instance := range c.Instances {
		if instance.ID == id {
			c.Instances[idx].Size = newSize
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}

// MovePublicIPToInstance implemented in a fake way for automated tests
func (c *FakeClient) MovePublicIPToInstance(id, ipAddress string) (*SimpleResponse, error) {
	oldIndex := -1
	for idx, instance := range c.Instances {
		if instance.PublicIP == ipAddress {
			oldIndex = idx
		}
	}

	newIndex := -1
	for idx, instance := range c.Instances {
		if instance.ID == id {
			newIndex = idx
		}
	}

	if oldIndex == -1 || newIndex == -1 {
		return &SimpleResponse{Result: "failed"}, nil
	}

	c.Instances[newIndex].PublicIP = c.Instances[oldIndex].PublicIP
	c.Instances[oldIndex].PublicIP = ""

	return &SimpleResponse{Result: "success"}, nil
}

// SetInstanceFirewall implemented in a fake way for automated tests
func (c *FakeClient) SetInstanceFirewall(id, firewallID string) (*SimpleResponse, error) {
	for idx, instance := range c.Instances {
		if instance.ID == id {
			c.Instances[idx].FirewallID = firewallID
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}
