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
	c.ListIPs()
}

// TestIPs is a test for the IPs method.
func TestIPs(t *testing.T) {
	g := NewWithT(t)

	client, err := NewFakeClient()
	g.Expect(err).To(BeNil())

	config := &CreateIPRequest{
		Name: "test-ip",
	}

	expected := &IP{
		Name: "test-ip",
	}

	ip, err := client.NewIP(config)
	g.Expect(err).To(BeNil())
	expected.ID = ip.ID
	g.Expect(ip.Name).To(Equal(expected.Name))

	ip, err = client.GetIP(ip.ID)
	g.Expect(err).To(BeNil())
	g.Expect(ip.Name).To(Equal(expected.Name))

	ips, err := client.ListIPs()
	g.Expect(err).To(BeNil())
	g.Expect(len(ips.Items)).To(Equal(1))

	ip, err = client.FindIP(ip.ID)
	g.Expect(err).To(BeNil())
	g.Expect(ip.Name).To(Equal(expected.Name))

	resp, err := client.DeleteIP(ip.ID)
	g.Expect(err).To(BeNil())
	g.Expect(resp).To(Equal(&SimpleResponse{Result: "success"}))
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

// TestKubernetesClustersPools is a test for the KubernetesClustersPools method.
func TestKubernetesClustersPools(t *testing.T) {
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
				{
					ID:       "aa2de9f9-c26e-4faf-8af5-26d7c9e6facf",
					Hostname: "bar",
				},
			},
			Pools: []KubernetesPool{
				{
					ID:    "33de5de2-14fd-44ba-a621-f6efbeeb9639",
					Count: 2,
					Size:  "small",
					InstanceNames: []string{
						"foo",
						"bar",
					},
					Instances: []KubernetesInstance{
						{
							ID:       "ad0dbf3f-4036-47f5-b33b-6822cf90799c0",
							Hostname: "foo",
						},
						{
							ID:       "aa2de9f9-c26e-4faf-8af5-26d7c9e6facf",
							Hostname: "bar",
						},
					},
				},
			},
		},
	}

	pools, err := client.ListKubernetesClusterPools("9c89d8b9-463d-45f2-8928-455eb3f3726")
	g.Expect(err).To(BeNil())
	g.Expect(len(pools)).To(Equal(1))

	pool, err := client.GetKubernetesClusterPool("9c89d8b9-463d-45f2-8928-455eb3f3726", "33de5de2-14fd-44ba-a621-f6efbeeb9639")
	g.Expect(err).To(BeNil())
	g.Expect(pool.ID).To(Equal("33de5de2-14fd-44ba-a621-f6efbeeb9639"))
	g.Expect(len(pool.InstanceNames)).To(Equal(2))

	pool, err = client.FindKubernetesClusterPool("9c89d8b9-463d-45f2-8928-455eb3f3726", "33de5de2-14fd-44ba-a621-f6efbeeb9639")
	g.Expect(err).To(BeNil())
	g.Expect(pool.ID).To(Equal("33de5de2-14fd-44ba-a621-f6efbeeb9639"))
	g.Expect(len(pool.InstanceNames)).To(Equal(2))

	result, err := client.DeleteKubernetesClusterPoolInstance("9c89d8b9-463d-45f2-8928-455eb3f3726", "33de5de2-14fd-44ba-a621-f6efbeeb9639", "ad0dbf3f-4036-47f5-b33b-6822cf90799c0")
	g.Expect(err).To(BeNil())
	g.Expect(string(result.Result)).To(Equal("success"))

	pc := KubernetesClusterPoolUpdateConfig{
		Count: 4,
	}
	pool, err = client.UpdateKubernetesClusterPool("9c89d8b9-463d-45f2-8928-455eb3f3726", "33de5de2-14fd-44ba-a621-f6efbeeb9639", &pc)
	g.Expect(err).To(BeNil())
	g.Expect(pool.Count).To(Equal(4))
}
