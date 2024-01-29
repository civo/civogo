package civogo

import "fmt"

func (c *FakeClient) ListObjectStoreCredentials() (*PaginatedObjectStoreCredentials, error) {
	return &PaginatedObjectStoreCredentials{
		Page:    1,
		PerPage: 20,
		Pages:   100,
		Items:   c.ObjectStoreCredential,
	}, nil
}

func (c *FakeClient) GetObjectStoreCredential(id string) (*ObjectStoreCredential, error) {
	for _, objectstorecredential := range c.ObjectStoreCredential {
		if objectstorecredential.ID == id {
			return &objectstorecredential, nil
		}
	}

	err := fmt.Errorf("unable to find the object store credential, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

func (c *FakeClient) FindObjectStoreCredential(search string) (*ObjectStoreCredential, error) {
	for _, objectstorecredential := range c.ObjectStoreCredential {
		if objectstorecredential.ID == search || objectstorecredential.Name == search {
			return &objectstorecredential, nil
		}
	}

	err := fmt.Errorf("unable to find the object store credential, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

func (c *FakeClient) NewObjectStoreCredential(v *CreateObjectStoreCredentialRequest) (*ObjectStoreCredential, error) {
	return nil, nil
}

func (c *FakeClient) UpdateObjectStoreCredential(id string, v *UpdateObjectStoreCredentialRequest) (*ObjectStoreCredential, error) {
	return nil, nil
}

func (c *FakeClient) DeleteObjectStoreCredential(id string) (*SimpleResponse, error) {
	return nil, nil
}
