package civogo

import (
	"fmt"
	"strings"
)

func (c *FakeClient) ListObjectStores() (*PaginatedObjectstores, error) {
	return &PaginatedObjectstores{
		Page:    1,
		PerPage: 20,
		Pages:   100,
		Items:   c.ObjectStore,
	}, nil
}

func (c *FakeClient) GetObjectStore(id string) (*ObjectStore, error) {
	for _, objectstore := range c.ObjectStore {
		if strings.Contains(objectstore.ID, id) {
			return &objectstore, nil
		}
	}

	err := fmt.Errorf("unable to find the object store, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

func (c *FakeClient) FindObjectStore(search string) (*ObjectStore, error) {
	for _, objectstore := range c.ObjectStore {
		if strings.Contains(objectstore.ID, search) || strings.Contains(objectstore.Name, search) {
			return &objectstore, nil
		}
	}

	err := fmt.Errorf("unable to find the object store, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

func (c *FakeClient) NewObjectStore(v *CreateObjectStoreRequest) (*ObjectStore, error) {
	objectstore := ObjectStore{
		ID:        c.generateID(),
		Name:      v.Name,
		MaxSize:   int(v.MaxSizeGB),
		OwnerInfo: BucketOwner{},
		BucketURL: fmt.Sprintf("https://objectstorage.%s.civo.com/%s", v.Region, v.Name),
		Status:    "ready",
	}

	c.ObjectStore = append(c.ObjectStore, objectstore)

	return &objectstore, nil
}

func (c *FakeClient) UpdateObjectStore(id string, v *UpdateObjectStoreRequest) (*ObjectStore, error) {
	for i, objectstore := range c.ObjectStore {
		if objectstore.ID == id {
			c.ObjectStore[i].MaxSize = int(v.MaxSizeGB)
			return &c.ObjectStore[i], nil
		}
	}

	err := fmt.Errorf("unable to find the object store, zero matches")
	return nil, ZeroMatchesError.wrap(err)
}

func (c *FakeClient) DeleteObjectStore(id string) (*SimpleResponse, error) {
	for i, objectstore := range c.ObjectStore {
		if objectstore.ID == id {
			c.ObjectStore[len(c.ObjectStore)-1], c.ObjectStore[i] = c.ObjectStore[i], c.ObjectStore[len(c.ObjectStore)-1]
			c.ObjectStore = c.ObjectStore[:len(c.ObjectStore)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}

func (c *FakeClient) GetObjectStoreStats(id string) (*ObjectStoreStats, error) {
	return &ObjectStoreStats{
		SizeKBUtilised: 2000000,
		MaxSizeKB:      10000000,
		NumObjects:     1000,
	}, nil
}
