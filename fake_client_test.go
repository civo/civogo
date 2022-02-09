package civogo

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestClienter(t *testing.T) {
	var c Clienter

	c, _ = NewClient("foo", "NYC1")
	c, _ = NewFakeClient()
	_, _ = c.ListAllInstances()
}

func TestInstances(t *testing.T) {
	client, _ := NewFakeClient()

	config := &InstanceConfig{
		Count:    1,
		Hostname: "foo.example.com",
	}
	client.CreateInstance(config)

	results, err := client.ListInstances(1, 10)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if results.Page != 1 {
		t.Errorf("Expected %+v, got %+v", 1, results.Page)
		return
	}
	if results.Pages != 1 {
		t.Errorf("Expected %+v, got %+v", 1, results.Pages)
		return
	}
	if results.PerPage != 10 {
		t.Errorf("Expected %+v, got %+v", 10, results.PerPage)
		return
	}
	if len(results.Items) != 1 {
		t.Errorf("Expected %+v, got %+v", 1, len(results.Items))
		return
	}
}

// TestLoadBalancers is a test for the LoadBalancers method.
func TestLoadBalancers(t *testing.T) {
	g := NewWithT(t)

	client, err := NewFakeClient()
	g.Expect(err).To(BeNil())

	backendConfig := []LoadBalancerBackendConfig{
		{
			IP:         "192.168.1.3",
			Protocol:   "TCP",
			SourcePort: 80,
			TargetPort: 31579,
		},
	}
	config := &LoadBalancerConfig{
		Name:      "foo",
		Algorithm: "round_robin",
		Backends:  backendConfig,
	}

	backends := []LoadBalancerBackend{
		{
			IP:         "192.168.1.3",
			Protocol:   "TCP",
			SourcePort: 80,
			TargetPort: 31579,
		},
	}

	expected := &LoadBalancer{
		Name:                  "foo",
		Algorithm:             "round_robin",
		Backends:              []LoadBalancerBackend(backends),
		ExternalTrafficPolicy: "Cluster",
	}

	loadbalancer, err := client.CreateLoadBalancer(config)
	g.Expect(err).To(BeNil())
	expected.ID = loadbalancer.ID
	g.Expect(loadbalancer.Name).To(Equal(expected.Name))
	g.Expect(loadbalancer.Algorithm).To(Equal(expected.Algorithm))
	g.Expect(loadbalancer.Backends).To(Equal(expected.Backends))
	g.Expect(loadbalancer.ExternalTrafficPolicy).To(Equal(expected.ExternalTrafficPolicy))

	loadbalancer, err = client.GetLoadBalancer(loadbalancer.ID)
	g.Expect(err).To(BeNil())
	g.Expect(loadbalancer.Name).To(Equal(expected.Name))
	g.Expect(loadbalancer.Algorithm).To(Equal(expected.Algorithm))
	g.Expect(loadbalancer.Backends).To(Equal(expected.Backends))
	g.Expect(loadbalancer.ExternalTrafficPolicy).To(Equal(expected.ExternalTrafficPolicy))

	loadbalancers, err := client.ListLoadBalancers()
	g.Expect(err).To(BeNil())
	g.Expect(len(loadbalancers)).To(Equal(1))

	loadbalancer, err = client.FindLoadBalancer(loadbalancer.ID)
	g.Expect(err).To(BeNil())
	g.Expect(loadbalancer.Name).To(Equal(expected.Name))
	g.Expect(loadbalancer.Algorithm).To(Equal(expected.Algorithm))
	g.Expect(loadbalancer.Backends).To(Equal(expected.Backends))
	g.Expect(loadbalancer.ExternalTrafficPolicy).To(Equal(expected.ExternalTrafficPolicy))

	resp, err := client.DeleteLoadBalancer(loadbalancer.ID)
	g.Expect(err).To(BeNil())
	g.Expect(resp).To(Equal(&SimpleResponse{Result: "success"}))
}

// TestKubernetesClustersInstances is a test for the KubernetesClustersInstances method.
func TestKubernetesClustersInstances(t *testing.T) {
	g := NewWithT(t)

	client, err := NewFakeClient()
	g.Expect(err).To(BeNil())

	client.Clusters = []KubernetesCluster{
		{
			ID:   "9c89d8b9-463d-45f2-8928-455eb3f3726",
			Name: "foo-cluster",
			Instances: []KubernetesInstance{
				{
					ID:       "ad0dbf3f-4036-47f5-b33b-6822cf90799c0",
					Hostname: "foo",
				},
			},
		},
	}
	client.Instances = []Instance{
		{
			ID:       "ad0dbf3f-4036-47f5-b33b-6822cf90799c0",
			Hostname: "foo",
		},
	}
	instances, err := client.ListKubernetesClusterInstances("9c89d8b9-463d-45f2-8928-455eb3f3726")
	g.Expect(err).To(BeNil())
	g.Expect(len(instances)).To(Equal(1))

	instance, err := client.FindKubernetesClusterInstance("9c89d8b9-463d-45f2-8928-455eb3f3726", "ad0dbf3f-4036-47f5-b33b-6822cf90799c0")
	g.Expect(err).To(BeNil())
	g.Expect(instance.ID).To(Equal("ad0dbf3f-4036-47f5-b33b-6822cf90799c0"))
	g.Expect(instance.Hostname).To(Equal("foo"))
}
