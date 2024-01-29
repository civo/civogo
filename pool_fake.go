package civogo

import (
	"fmt"
	"strings"
)

// ListKubernetesClusterPools implemented in a fake way for automated tests
func (c *FakeClient) ListKubernetesClusterPools(cid string) ([]KubernetesPool, error) {
	pools := []KubernetesPool{}
	found := false

	for _, cs := range c.Kubernetes {
		if cs.ID == cid {
			found = true
			pools = cs.Pools
			break
		}
	}

	if found {
		return pools, nil
	}

	err := fmt.Errorf("unable to get kubernetes cluster %s", cid)
	return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
}

// GetKubernetesClusterPool implemented in a fake way for automated tests
func (c *FakeClient) GetKubernetesClusterPool(cid, pid string) (*KubernetesPool, error) {
	pool := &KubernetesPool{}
	clusterFound := false
	poolFound := false

	for _, cs := range c.Kubernetes {
		if cs.ID == cid {
			clusterFound = true
			for _, p := range cs.Pools {
				if p.ID == pid {
					poolFound = true
					pool = &p
					break
				}
			}
		}
	}

	if !clusterFound {
		err := fmt.Errorf("unable to get kubernetes cluster %s", cid)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	if !poolFound {
		err := fmt.Errorf("unable to get kubernetes pool %s", pid)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	return pool, nil
}

// FindKubernetesClusterPool implemented in a fake way for automated tests
func (c *FakeClient) FindKubernetesClusterPool(cid, search string) (*KubernetesPool, error) {
	pool := &KubernetesPool{}
	clusterFound := false
	poolFound := false

	for _, cs := range c.Kubernetes {
		if cs.ID == cid {
			clusterFound = true
			for _, p := range cs.Pools {
				if p.ID == search || strings.Contains(p.ID, search) {
					poolFound = true
					pool = &p
					break
				}
			}
		}
	}

	if !clusterFound {
		err := fmt.Errorf("unable to get kubernetes cluster %s", cid)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	if !poolFound {
		err := fmt.Errorf("unable to get kubernetes pool %s", search)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	return pool, nil
}

// DeleteKubernetesClusterPoolInstance implemented in a fake way for automated tests
func (c *FakeClient) DeleteKubernetesClusterPoolInstance(cid, pid, id string) (*SimpleResponse, error) {
	clusterFound := false
	poolFound := false
	instanceFound := false

	for ci, cs := range c.Kubernetes {
		if cs.ID == cid {
			clusterFound = true
			for pi, p := range cs.Pools {
				if p.ID == pid {
					poolFound = true
					for i, in := range p.Instances {
						if in.ID == id {
							instanceFound = true
							p.Instances = append(p.Instances[:i], p.Instances[i+1:]...)

							instanceNames := []string{}
							for _, in := range p.Instances {
								instanceNames = append(instanceNames, in.Hostname)
							}
							p.InstanceNames = instanceNames
							c.Kubernetes[ci].Pools[pi] = p
							break
						}
					}
				}
			}
		}
	}

	if !clusterFound {
		err := fmt.Errorf("unable to get kubernetes cluster %s", cid)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	if !poolFound {
		err := fmt.Errorf("unable to get kubernetes pool %s", pid)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	if !instanceFound {
		err := fmt.Errorf("unable to get kubernetes pool instance %s", id)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	return &SimpleResponse{
		Result: "success",
	}, nil
}

// UpdateKubernetesClusterPool implemented in a fake way for automated tests
func (c *FakeClient) UpdateKubernetesClusterPool(cid, pid string, config *KubernetesClusterPoolUpdateConfig) (*KubernetesPool, error) {
	clusterFound := false
	poolFound := false

	pool := KubernetesPool{}
	for _, cs := range c.Kubernetes {
		if cs.ID == cid {
			clusterFound = true
			for _, p := range cs.Pools {
				if p.ID == pid {
					poolFound = true
					p.Count = config.Count
					pool = p
				}
			}
		}
	}

	if !clusterFound {
		err := fmt.Errorf("unable to get kubernetes cluster %s", cid)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	if !poolFound {
		err := fmt.Errorf("unable to get kubernetes pool %s", pid)
		return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
	}

	return &pool, nil
}
