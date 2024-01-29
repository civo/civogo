package civogo

import (
	"fmt"
	"strings"
)

// ListKubernetesClusters implemented in a fake way for automated tests
func (c *FakeClient) ListKubernetesClusters() (*PaginatedKubernetesClusters, error) {
	return &PaginatedKubernetesClusters{
		Items:   c.Kubernetes,
		Page:    1,
		PerPage: 10,
		Pages:   1,
	}, nil
}

// FindKubernetesCluster implemented in a fake way for automated tests
func (c *FakeClient) FindKubernetesCluster(search string) (*KubernetesCluster, error) {
	for _, cluster := range c.Kubernetes {
		if strings.Contains(cluster.Name, search) || cluster.ID == search {
			return &cluster, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", search)
	return nil, ZeroMatchesError.wrap(err)
}

// ListKubernetesClusterInstances implemented in a fake way for automated tests
func (c *FakeClient) ListKubernetesClusterInstances(id string) ([]Instance, error) {
	for _, cluster := range c.Kubernetes {
		if cluster.ID == id {
			instaces := make([]Instance, 0)
			for _, kins := range cluster.Instances {
				for _, instance := range c.Instances {
					if instance.ID == kins.ID {
						instaces = append(instaces, instance)
					}
				}
			}
			return instaces, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", id)
	return nil, DatabaseKubernetesClusterNotFoundError.wrap(err)
}

// FindKubernetesClusterInstance implemented in a fake way for automated tests
func (c *FakeClient) FindKubernetesClusterInstance(clusterID, search string) (*Instance, error) {
	instances, err := c.ListKubernetesClusterInstances(clusterID)
	if err != nil {
		return nil, decodeError(err)
	}

	exactMatch := false
	partialMatchesCount := 0
	result := Instance{}

	for _, instance := range instances {
		if instance.Hostname == search || instance.ID == search {
			exactMatch = true
			result = instance
		} else if strings.Contains(instance.Hostname, search) || strings.Contains(instance.ID, search) {
			if !exactMatch {
				result = instance
				partialMatchesCount++
			}
		}
	}

	if exactMatch || partialMatchesCount == 1 {
		return &result, nil
	} else if partialMatchesCount > 1 {
		err := fmt.Errorf("unable to find %s because there were multiple matches", search)
		return nil, MultipleMatchesError.wrap(err)
	} else {
		err := fmt.Errorf("unable to find %s, zero matches", search)
		return nil, ZeroMatchesError.wrap(err)
	}
}

// NewKubernetesClusters implemented in a fake way for automated tests
func (c *FakeClient) NewKubernetesClusters(kc *KubernetesClusterConfig) (*KubernetesCluster, error) {
	cluster := KubernetesCluster{
		ID:             c.generateID(),
		Name:           kc.Name,
		MasterIP:       c.generatePublicIP(),
		NumTargetNode:  kc.NumTargetNodes,
		TargetNodeSize: kc.TargetNodesSize,
		Ready:          true,
		Status:         "ACTIVE",
		Instances:      make([]KubernetesInstance, 0),
		Pools:          make([]KubernetesPool, 0),
	}
	pool := KubernetesPool{
		Instances: make([]KubernetesInstance, 0),
	}
	for i := 0; i < kc.NumTargetNodes; i++ {
		instance := KubernetesInstance{
			ID:       c.generateID(),
			Hostname: fmt.Sprintf("%s_pool_%d", kc.Name, i),
		}
		pool.Instances = append(pool.Instances, instance)
		cluster.Instances = append(pool.Instances, instance)
	}

	cluster.Pools = append(cluster.Pools, pool)
	c.Kubernetes = append(c.Kubernetes, cluster)
	return &cluster, nil
}

// GetKubernetesCluster implemented in a fake way for automated tests
func (c *FakeClient) GetKubernetesCluster(id string) (*KubernetesCluster, error) {
	for _, cluster := range c.Kubernetes {
		if cluster.ID == id {
			return &cluster, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// UpdateKubernetesCluster implemented in a fake way for automated tests
func (c *FakeClient) UpdateKubernetesCluster(id string, kc *KubernetesClusterConfig) (*KubernetesCluster, error) {
	for i, cluster := range c.Kubernetes {
		if cluster.ID == id {
			c.Kubernetes[i].Name = kc.Name
			c.Kubernetes[i].NumTargetNode = kc.NumTargetNodes
			c.Kubernetes[i].TargetNodeSize = kc.TargetNodesSize
			return &cluster, nil
		}
	}

	err := fmt.Errorf("unable to find %s, zero matches", id)
	return nil, ZeroMatchesError.wrap(err)
}

// ListKubernetesMarketplaceApplications implemented in a fake way for automated tests
func (c *FakeClient) ListKubernetesMarketplaceApplications() ([]KubernetesMarketplaceApplication, error) {
	return []KubernetesMarketplaceApplication{}, nil
}

// DeleteKubernetesCluster implemented in a fake way for automated tests
func (c *FakeClient) DeleteKubernetesCluster(id string) (*SimpleResponse, error) {
	for i, cluster := range c.Kubernetes {
		if cluster.ID == id {
			c.Kubernetes[len(c.Kubernetes)-1], c.Kubernetes[i] = c.Kubernetes[i], c.Kubernetes[len(c.Kubernetes)-1]
			c.Kubernetes = c.Kubernetes[:len(c.Kubernetes)-1]
			return &SimpleResponse{Result: "success"}, nil
		}
	}

	return &SimpleResponse{Result: "failed"}, nil
}

// RecycleKubernetesCluster implemented in a fake way for automated tests
func (c *FakeClient) RecycleKubernetesCluster(id string, hostname string) (*SimpleResponse, error) {
	return &SimpleResponse{Result: "success"}, nil
}

// ListAvailableKubernetesVersions implemented in a fake way for automated tests
func (c *FakeClient) ListAvailableKubernetesVersions() ([]KubernetesVersion, error) {
	return []KubernetesVersion{
		{
			Version: "1.20+k3s1",
			Type:    "stable",
		},
	}, nil
}
