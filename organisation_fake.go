package civogo

import "time"

// GetOrganisation implemented in a fake way for automated tests
func (c *FakeClient) GetOrganisation() (*Organisation, error) {
	return &c.Organisation, nil
}

// CreateOrganisation implemented in a fake way for automated tests
func (c *FakeClient) CreateOrganisation(name string) (*Organisation, error) {
	c.Organisation.ID = c.generateID()
	c.Organisation.Name = name
	return &c.Organisation, nil
}

// RenameOrganisation implemented in a fake way for automated tests
func (c *FakeClient) RenameOrganisation(name string) (*Organisation, error) {
	c.Organisation.Name = name
	return &c.Organisation, nil
}

// AddAccountToOrganisation implemented in a fake way for automated tests
func (c *FakeClient) AddAccountToOrganisation(accountID string) ([]Account, error) {
	c.OrganisationAccounts = append(c.OrganisationAccounts, Account{
		ID:        accountID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return c.ListAccountsInOrganisation()
}

// ListAccountsInOrganisation implemented in a fake way for automated tests
func (c *FakeClient) ListAccountsInOrganisation() ([]Account, error) {
	return c.OrganisationAccounts, nil
}
