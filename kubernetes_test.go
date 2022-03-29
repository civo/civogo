package civogo

import (
	"reflect"
	"testing"
	"time"
)

func TestListKubernetesClusters(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters": `{"page":1,"per_page":20,"pages":1,"items":[{
		  "id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
		  "name": "your-cluster-name",
		  "version": "2",
		  "status": "ACTIVE",
		  "ready": true,
		  "num_target_nodes": 1,
		  "target_nodes_size": "g2.xsmall",
		  "built_at": "2019-09-23T13:04:23.000+01:00",
		  "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		  "kubernetes_version": "0.8.1",
		  "api_endpoint": "https://your.cluster.ip.address:6443",
		  "master_ip": "your.cluster.ip.address",
		  "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		  "tags": [],
		  "created_at": "2019-09-23T13:02:59.000+01:00",
		  "firewall_id": "42118911-44c2-4cab-ad77-bcae062815b3",
		  "instances": [{
			"hostname": "kube-master-HEXDIGITS",
			"size": "g2.xsmall",
			"region": "lon1",
			"created_at": "2019-09-23T13:03:00.000+01:00",
			"status": "ACTIVE",
			"firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			"public_ip": "your.cluster.ip.address",
			"tags": ["civo-kubernetes:installed", "civo-kubernetes:master"]
		  }],
			"required_pools": [
        {"id": "d6e52d2a-9199-4b04-9118-0559e4d0ce63", "size": "g3.k3s.xsmall","count": 1},
        {"id": "fc432d8c-c3ab-4716-8ab0-d2164b932da7", "size": "g4s.kube.medium", "count": 2}
    ],
		  "installed_applications": [{
			"application": "Traefik",
			"title": null,
			"version": "(default)",
			"dependencies": null,
			"maintainer": "@Rancher_Labs",
			"description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			"post_install": "Some documentation here\n",
			"installed": true,
			"url": "https://traefik.io",
			"category": "architecture",
			"updated_at": "2019-09-23T13:02:59.000+01:00",
			"image_url": "https://api.civo.com/k3s-marketplace/traefik.png",
			"plan": null,
			"configuration": {}
		  }],
			"cni_plugin": "flannel",
			"ccm_installed": "true"
		}]}`,
	})
	defer server.Close()
	got, err := client.ListKubernetesClusters()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	buildAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:04:23.000+01:00")
	createAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")
	createAtInstance, _ := time.Parse(time.RFC3339, "2019-09-23T13:03:00.000+01:00")
	updateAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")

	expected := &PaginatedKubernetesClusters{
		Page:    1,
		PerPage: 20,
		Pages:   1,
		Items: []KubernetesCluster{
			{
				ID:                "69a23478-a89e-41d2-97b1-6f4c341cee70",
				Name:              "your-cluster-name",
				Version:           "2",
				Status:            "ACTIVE",
				BuiltAt:           buildAt,
				Ready:             true,
				NumTargetNode:     1,
				TargetNodeSize:    "g2.xsmall",
				KubeConfig:        "YAML_VERSION_OF_KUBECONFIG_HERE\n",
				KubernetesVersion: "0.8.1",
				APIEndPoint:       "https://your.cluster.ip.address:6443",
				MasterIP:          "your.cluster.ip.address",
				DNSEntry:          "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
				CreatedAt:         createAt,
				Tags:              []string{},
				FirewallID:        "42118911-44c2-4cab-ad77-bcae062815b3",
				Instances: []KubernetesInstance{{
					Hostname:   "kube-master-HEXDIGITS",
					Size:       "g2.xsmall",
					Region:     "lon1",
					CreatedAt:  createAtInstance,
					Status:     "ACTIVE",
					FirewallID: "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
					PublicIP:   "your.cluster.ip.address",
					Tags:       []string{"civo-kubernetes:installed", "civo-kubernetes:master"},
				}},
				RequiredPools: []RequiredPools{
					{
						ID:    "d6e52d2a-9199-4b04-9118-0559e4d0ce63",
						Size:  "g3.k3s.xsmall",
						Count: 1,
					},
					{
						ID:    "fc432d8c-c3ab-4716-8ab0-d2164b932da7",
						Size:  "g4s.kube.medium",
						Count: 2,
					},
				},
				InstalledApplications: []KubernetesInstalledApplication{{
					Application:   "Traefik",
					Version:       "(default)",
					Maintainer:    "@Rancher_Labs",
					Description:   "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
					PostInstall:   "Some documentation here\n",
					URL:           "https://traefik.io",
					UpdatedAt:     updateAt,
					Installed:     true,
					Category:      "architecture",
					ImageURL:      "https://api.civo.com/k3s-marketplace/traefik.png",
					Configuration: map[string]ApplicationConfiguration{},
				}},
				CNIPlugin:    "flannel",
				CCMInstalled: "true",
			},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestFindKubernetesCluster(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters": `{"page":1,"per_page":20,"pages":1,"items":[
			{ "id": "69a23478-a89e-41d2-97b1-6f4c341cee70", "name": "your-first-cluster-name", "version": "2", "status": "ACTIVE", "ready": true, "num_target_nodes": 1, "target_nodes_size": "g2.xsmall", "built_at": "2019-09-23T13:04:23.000+01:00", "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n", "kubernetes_version": "0.8.1", "api_endpoint": "https://your.cluster.ip.address:6443", "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com", "tags": [], "created_at": "2019-09-23T13:02:59.000+01:00", "firewall_id": "42118911-44c2-4cab-ad77-bcae062815b3", "instances": [{ "hostname": "kube-master-HEXDIGITS", "size": "g2.xsmall", "region": "lon1", "created_at": "2019-09-23T13:03:00.000+01:00", "status": "ACTIVE", "firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64", "public_ip": "your.cluster.ip.address", "tags": ["civo-kubernetes:installed", "civo-kubernetes:master"] }], "installed_applications": [{ "application": "Traefik", "title": null, "version": "(default)", "dependencies": null, "maintainer": "@Rancher_Labs", "description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.", "post_install": "Some documentation here\n", "installed": true, "url": "https://traefik.io", "category": "architecture", "updated_at": "2019-09-23T13:02:59.000+01:00", "image_url": "https://api.civo.com/k3s-marketplace/traefik.png", "plan": null, "configuration": {} }] },
			{ "id": "d1cd0b71-5da1-492e-9d0d-a46ccdaae2fa", "name": "your-second-cluster-name", "version": "2", "status": "ACTIVE", "ready": true, "num_target_nodes": 1, "target_nodes_size": "g2.xsmall", "built_at": "2019-09-23T13:04:23.000+01:00", "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n", "kubernetes_version": "0.8.1", "api_endpoint": "https://your.cluster.ip.address:6443", "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com", "tags": [], "created_at": "2019-09-23T13:02:59.000+01:00", "firewall_id": "42118911-44c2-4cab-ad77-bcae062815b3", "instances": [{ "hostname": "kube-master-HEXDIGITS", "size": "g2.xsmall", "region": "lon1", "created_at": "2019-09-23T13:03:00.000+01:00", "status": "ACTIVE", "firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64", "public_ip": "your.cluster.ip.address", "tags": ["civo-kubernetes:installed", "civo-kubernetes:master"] }], "installed_applications": [{ "application": "Traefik", "title": null, "version": "(default)", "dependencies": null, "maintainer": "@Rancher_Labs", "description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.", "post_install": "Some documentation here\n", "installed": true, "url": "https://traefik.io", "category": "architecture", "updated_at": "2019-09-23T13:02:59.000+01:00", "image_url": "https://api.civo.com/k3s-marketplace/traefik.png", "plan": null, "configuration": {} }] }
		]}`,
	})
	defer server.Close()

	got, _ := client.FindKubernetesCluster("69a23478")
	if got.ID != "69a23478-a89e-41d2-97b1-6f4c341cee70" {
		t.Errorf("Expected %s, got %s", "69a23478-a89e-41d2-97b1-6f4c341cee70", got.ID)
	}

	got, _ = client.FindKubernetesCluster("d1cd0b71")
	if got.ID != "d1cd0b71-5da1-492e-9d0d-a46ccdaae2fa" {
		t.Errorf("Expected %s, got %s", "d1cd0b71-5da1-492e-9d0d-a46ccdaae2fa", got.ID)
	}

	got, _ = client.FindKubernetesCluster("first")
	if got.ID != "69a23478-a89e-41d2-97b1-6f4c341cee70" {
		t.Errorf("Expected %s, got %s", "69a23478-a89e-41d2-97b1-6f4c341cee70", got.ID)
	}

	got, _ = client.FindKubernetesCluster("YOUR-FIRST-CLUSTER-NAME")
	if got.ID != "69a23478-a89e-41d2-97b1-6f4c341cee70" {
		t.Errorf("Expected %s, got %s", "69a23478-a89e-41d2-97b1-6f4c341cee70", got.ID)
	}

	got, _ = client.FindKubernetesCluster("second")
	if got.ID != "d1cd0b71-5da1-492e-9d0d-a46ccdaae2fa" {
		t.Errorf("Expected %s, got %s", "d1cd0b71-5da1-492e-9d0d-a46ccdaae2fa", got.ID)
	}

	_, err := client.FindKubernetesCluster("cluster")
	if err.Error() != "MultipleMatchesError: unable to find cluster because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find cluster because there were multiple matches", err.Error())
	}

	_, err = client.FindKubernetesCluster("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}

func TestNewKubernetesClusters(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters": `{
		  "id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
		  "name": "your-cluster-name",
		  "version": "2",
		  "status": "ACTIVE",
		  "ready": true,
		  "num_target_nodes": 1,
		  "target_nodes_size": "g2.xsmall",
		  "built_at": "2019-09-23T13:04:23.000+01:00",
		  "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		  "kubernetes_version": "0.8.1",
		  "api_endpoint": "https://your.cluster.ip.address:6443",
		  "master_ip": "your.cluster.ip.address",
		  "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		  "tags": [],
		  "created_at": "2019-09-23T13:02:59.000+01:00",
		  "firewall_id": "42118911-44c2-4cab-ad77-bcae062815b3",
		  "instances": [{
			"hostname": "kube-master-HEXDIGITS",
			"size": "g2.xsmall",
			"region": "lon1",
			"created_at": "2019-09-23T13:03:00.000+01:00",
			"status": "ACTIVE",
			"firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			"public_ip": "your.cluster.ip.address",
			"tags": ["civo-kubernetes:installed", "civo-kubernetes:master"]
		  }],
		  "installed_applications": [{
			"application": "Traefik",
			"title": null,
			"version": "(default)",
			"dependencies": null,
			"maintainer": "@Rancher_Labs",
			"description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			"post_install": "Some documentation here\n",
			"installed": true,
			"url": "https://traefik.io",
			"category": "architecture",
			"updated_at": "2019-09-23T13:02:59.000+01:00",
			"image_url": "https://api.civo.com/k3s-marketplace/traefik.png",
			"plan": null,
			"configuration": {}
		  }],
			"cni_plugin": "flannel"
		}`,
	})
	defer server.Close()

	cfg := &KubernetesClusterConfig{
		Name:              "your-cluster-name",
		Tags:              "",
		KubernetesVersion: "0.8.1",
		NumTargetNodes:    3,
		TargetNodesSize:   "g2.xsmall",
		Applications:      "traefik",
	}
	got, err := client.NewKubernetesClusters(cfg)

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	buildAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:04:23.000+01:00")
	createAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")
	createAtInstance, _ := time.Parse(time.RFC3339, "2019-09-23T13:03:00.000+01:00")
	updateAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")

	expected := &KubernetesCluster{
		ID:                "69a23478-a89e-41d2-97b1-6f4c341cee70",
		Name:              "your-cluster-name",
		Version:           "2",
		Status:            "ACTIVE",
		BuiltAt:           buildAt,
		Ready:             true,
		NumTargetNode:     1,
		TargetNodeSize:    "g2.xsmall",
		KubeConfig:        "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		KubernetesVersion: "0.8.1",
		APIEndPoint:       "https://your.cluster.ip.address:6443",
		MasterIP:          "your.cluster.ip.address",
		DNSEntry:          "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		CreatedAt:         createAt,
		Tags:              []string{},
		FirewallID:        "42118911-44c2-4cab-ad77-bcae062815b3",
		Instances: []KubernetesInstance{{
			Hostname:   "kube-master-HEXDIGITS",
			Size:       "g2.xsmall",
			Region:     "lon1",
			CreatedAt:  createAtInstance,
			Status:     "ACTIVE",
			FirewallID: "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			PublicIP:   "your.cluster.ip.address",
			Tags:       []string{"civo-kubernetes:installed", "civo-kubernetes:master"},
		}},
		InstalledApplications: []KubernetesInstalledApplication{{
			Application:   "Traefik",
			Version:       "(default)",
			Maintainer:    "@Rancher_Labs",
			Description:   "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			PostInstall:   "Some documentation here\n",
			URL:           "https://traefik.io",
			UpdatedAt:     updateAt,
			Installed:     true,
			Category:      "architecture",
			ImageURL:      "https://api.civo.com/k3s-marketplace/traefik.png",
			Configuration: map[string]ApplicationConfiguration{},
		}},
		CNIPlugin: "flannel",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetKubernetesClusters(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/69a23478-a89e-41d2-97b1-6f4c341cee70": `{
		  "id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
		  "name": "your-cluster-name",
		  "version": "2",
		  "status": "ACTIVE",
		  "ready": true,
		  "num_target_nodes": 1,
		  "target_nodes_size": "g2.xsmall",
		  "built_at": "2019-09-23T13:04:23.000+01:00",
		  "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		  "kubernetes_version": "0.8.1",
		  "api_endpoint": "https://your.cluster.ip.address:6443",
		  "master_ip": "your.cluster.ip.address",
		  "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		  "tags": [],
		  "created_at": "2019-09-23T13:02:59.000+01:00",
		  "firewall_id": "42118911-44c2-4cab-ad77-bcae062815b3",
		  "instances": [{
			"hostname": "kube-master-HEXDIGITS",
			"size": "g2.xsmall",
			"region": "lon1",
			"created_at": "2019-09-23T13:03:00.000+01:00",
			"status": "ACTIVE",
			"firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			"public_ip": "your.cluster.ip.address",
			"tags": ["civo-kubernetes:installed", "civo-kubernetes:master"]
		  }],
		  "installed_applications": [{
			"application": "Traefik",
			"title": null,
			"version": "(default)",
			"dependencies": null,
			"maintainer": "@Rancher_Labs",
			"description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			"post_install": "Some documentation here\n",
			"installed": true,
			"url": "https://traefik.io",
			"category": "architecture",
			"updated_at": "2019-09-23T13:02:59.000+01:00",
			"image_url": "https://api.civo.com/k3s-marketplace/traefik.png",
			"plan": null,
			"configuration": {}
		  }],
			"cni_plugin": "flannel",
			"ccm_installed": "false"
		}`,
	})
	defer server.Close()

	got, err := client.GetKubernetesCluster("69a23478-a89e-41d2-97b1-6f4c341cee70")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	buildAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:04:23.000+01:00")
	createAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")
	createAtInstance, _ := time.Parse(time.RFC3339, "2019-09-23T13:03:00.000+01:00")
	updateAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")

	expected := &KubernetesCluster{
		ID:                "69a23478-a89e-41d2-97b1-6f4c341cee70",
		Name:              "your-cluster-name",
		Version:           "2",
		Status:            "ACTIVE",
		BuiltAt:           buildAt,
		Ready:             true,
		NumTargetNode:     1,
		TargetNodeSize:    "g2.xsmall",
		KubeConfig:        "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		KubernetesVersion: "0.8.1",
		APIEndPoint:       "https://your.cluster.ip.address:6443",
		MasterIP:          "your.cluster.ip.address",
		DNSEntry:          "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		CreatedAt:         createAt,
		Tags:              []string{},
		FirewallID:        "42118911-44c2-4cab-ad77-bcae062815b3",
		Instances: []KubernetesInstance{{
			Hostname:   "kube-master-HEXDIGITS",
			Size:       "g2.xsmall",
			Region:     "lon1",
			CreatedAt:  createAtInstance,
			Status:     "ACTIVE",
			FirewallID: "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			PublicIP:   "your.cluster.ip.address",
			Tags:       []string{"civo-kubernetes:installed", "civo-kubernetes:master"},
		}},
		InstalledApplications: []KubernetesInstalledApplication{{
			Application:   "Traefik",
			Version:       "(default)",
			Maintainer:    "@Rancher_Labs",
			Description:   "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			PostInstall:   "Some documentation here\n",
			URL:           "https://traefik.io",
			UpdatedAt:     updateAt,
			Installed:     true,
			Category:      "architecture",
			ImageURL:      "https://api.civo.com/k3s-marketplace/traefik.png",
			Configuration: map[string]ApplicationConfiguration{},
		}},
		CNIPlugin:    "flannel",
		CCMInstalled: "false",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpdateKubernetesClusters(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/69a23478-a89e-41d2-97b1-6f4c341cee70": `{
		  "id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
		  "name": "cluster-name",
		  "version": "2",
		  "status": "ACTIVE",
		  "ready": true,
		  "num_target_nodes": 6,
		  "target_nodes_size": "g2.xsmall",
		  "built_at": "2019-09-23T13:04:23.000+01:00",
		  "kubeconfig": "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		  "kubernetes_version": "0.8.1",
		  "api_endpoint": "https://your.cluster.ip.address:6443",
		  "master_ip": "your.cluster.ip.address",
		  "dns_entry": "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		  "tags": [],
		  "created_at": "2019-09-23T13:02:59.000+01:00",
		  "firewall_id": "42118911-44c2-4cab-ad77-bcae062815b3",
		  "instances": [{
			"hostname": "kube-master-HEXDIGITS",
			"size": "g2.xsmall",
			"region": "lon1",
			"created_at": "2019-09-23T13:03:00.000+01:00",
			"status": "ACTIVE",
			"firewall_id": "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			"public_ip": "your.cluster.ip.address",
			"tags": ["civo-kubernetes:installed", "civo-kubernetes:master"]
		  }],
		  "installed_applications": [{
			"application": "Traefik",
			"title": null,
			"version": "(default)",
			"dependencies": null,
			"maintainer": "@Rancher_Labs",
			"description": "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			"post_install": "Some documentation here\n",
			"installed": true,
			"url": "https://traefik.io",
			"category": "architecture",
			"updated_at": "2019-09-23T13:02:59.000+01:00",
			"image_url": "https://api.civo.com/k3s-marketplace/traefik.png",
			"plan": null,
			"configuration": {}
		  }],
			"cni_plugin": "flannel"
		}`,
	})
	defer server.Close()

	cfg := &KubernetesClusterConfig{
		Name:           "cluster-name",
		NumTargetNodes: 6,
	}

	got, err := client.UpdateKubernetesCluster("69a23478-a89e-41d2-97b1-6f4c341cee70", cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	buildAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:04:23.000+01:00")
	createAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")
	createAtInstance, _ := time.Parse(time.RFC3339, "2019-09-23T13:03:00.000+01:00")
	updateAt, _ := time.Parse(time.RFC3339, "2019-09-23T13:02:59.000+01:00")

	expected := &KubernetesCluster{
		ID:                "69a23478-a89e-41d2-97b1-6f4c341cee70",
		Name:              "cluster-name",
		Version:           "2",
		Status:            "ACTIVE",
		BuiltAt:           buildAt,
		Ready:             true,
		NumTargetNode:     6,
		TargetNodeSize:    "g2.xsmall",
		KubeConfig:        "YAML_VERSION_OF_KUBECONFIG_HERE\n",
		KubernetesVersion: "0.8.1",
		APIEndPoint:       "https://your.cluster.ip.address:6443",
		MasterIP:          "your.cluster.ip.address",
		DNSEntry:          "69a23478-a89e-41d2-97b1-6f4c341cee70.k8s.civo.com",
		CreatedAt:         createAt,
		Tags:              []string{},
		FirewallID:        "42118911-44c2-4cab-ad77-bcae062815b3",
		Instances: []KubernetesInstance{{
			Hostname:   "kube-master-HEXDIGITS",
			Size:       "g2.xsmall",
			Region:     "lon1",
			CreatedAt:  createAtInstance,
			Status:     "ACTIVE",
			FirewallID: "5f0ba9ed-5ca7-4e14-9a09-449a84196d64",
			PublicIP:   "your.cluster.ip.address",
			Tags:       []string{"civo-kubernetes:installed", "civo-kubernetes:master"},
		}},
		InstalledApplications: []KubernetesInstalledApplication{{
			Application:   "Traefik",
			Version:       "(default)",
			Maintainer:    "@Rancher_Labs",
			Description:   "A reverse proxy/load-balancer that's easy, dynamic, automatic, fast, full-featured, open source, production proven and provides metrics.",
			PostInstall:   "Some documentation here\n",
			URL:           "https://traefik.io",
			UpdatedAt:     updateAt,
			Installed:     true,
			Category:      "architecture",
			ImageURL:      "https://api.civo.com/k3s-marketplace/traefik.png",
			Configuration: map[string]ApplicationConfiguration{},
		}},
		CNIPlugin: "flannel",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestKubernetesMarketplaceApplication(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/applications": `[{
		  "name": "MariaDB",
		  "title": null,
		  "version": "10.4.7",
		  "default": null,
		  "dependencies": ["Longhorn"],
		  "maintainer": "hello@civo.com",
		  "description": "MariaDB is a community-developed fork of MySQL intended to remain free under the GNU GPL.",
		  "post_install": "Instructions go here\n",
		  "url": "https://mariadb.com",
		  "category": "database",
		  "image_url": "https://api.civo.com/k3s-marketplace/mariadb.png",
		  "plans": [{
			"label": "5GB",
			"configuration": {
			  "VOLUME_SIZE": {
				"value": "5Gi"
			  }
			}
		  }, {
			"label": "10GB",
			"configuration": {
			  "VOLUME_SIZE": {
				"value": "10Gi"
			  }
			}
		  }]
		}]`,
	})
	defer server.Close()
	got, err := client.ListKubernetesMarketplaceApplications()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := []KubernetesMarketplaceApplication{{
		Name:         "MariaDB",
		Version:      "10.4.7",
		Dependencies: []string{"Longhorn"},
		Maintainer:   "hello@civo.com",
		Description:  "MariaDB is a community-developed fork of MySQL intended to remain free under the GNU GPL.",
		PostInstall:  "Instructions go here\n",
		URL:          "https://mariadb.com",
		Category:     "database",
		Plans: []KubernetesMarketplacePlan{{
			Label:         "5GB",
			Configuration: map[string]KubernetesPlanConfiguration{"VOLUME_SIZE": {Value: "5Gi"}},
		}, {
			Label:         "10GB",
			Configuration: map[string]KubernetesPlanConfiguration{"VOLUME_SIZE": {Value: "10Gi"}},
		}},
	}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestDeleteKubernetesClusters(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/12346": `{"result": "success"}`,
	})
	defer server.Close()
	got, err := client.DeleteKubernetesCluster("12346")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestRecycleKubernetesClusters(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/clusters/12346": `{"result": "success"}`,
	})
	defer server.Close()

	got, err := client.RecycleKubernetesCluster("12346", "test-hostname")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &SimpleResponse{Result: "success"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListAvailableKubernetesVersions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/kubernetes/versions": `[
		  {
			"version": "0.9.1",
			"type": "stable",
			"default": true
		  },
		  {
			"version": "0.8.1",
			"type": "legacy"
		  }
		]`,
	})
	defer server.Close()
	got, err := client.ListAvailableKubernetesVersions()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := []KubernetesVersion{{
		Default: true,
		Type:    "stable",
		Version: "0.9.1",
	}, {
		Type:    "legacy",
		Version: "0.8.1",
	}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
