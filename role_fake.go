package civogo

import "fmt"

// ListRoles implemented in a fake way for automated tests
func (c *FakeClient) ListRoles() ([]Role, error) {
	return c.OrganisationRoles, nil
}

// CreateRole implemented in a fake way for automated tests
func (c *FakeClient) CreateRole(name, permissions string) (*Role, error) {
	role := Role{
		ID:          c.generateID(),
		Name:        name,
		Permissions: permissions,
	}
	c.OrganisationRoles = append(c.OrganisationRoles, role)
	return &role, nil
}

// DeleteRole implemented in a fake way for automated tests
func (c *FakeClient) DeleteRole(id string) (*SimpleResponse, error) {
	for i, role := range c.OrganisationRoles {
		if role.ID == id {
			c.OrganisationRoles[len(c.OrganisationRoles)-1], c.OrganisationRoles[i] = c.OrganisationRoles[i], c.OrganisationRoles[len(c.OrganisationRoles)-1]
			c.OrganisationRoles = c.OrganisationRoles[:len(c.OrganisationRoles)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, fmt.Errorf("unable to find that role")
}
