package civogo

import (
	"reflect"
	"testing"
	"time"
)

func TestClienterDiskImage(t *testing.T) {
	var c Clienter

	c, _ = NewClient("foo", "NYC1")
	c, _ = NewFakeClient()
	_, _ = c.ListDiskImages()
}

func TestGetDiskImage(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/disk_images/b82168fe-66f6-4b38-a3b8-5283542d5475": `{
			"id": "b82168fe-66f6-4b38-a3b8-5283542d5475",
			"name": "centos-7",
			"version": "7",
			"state": "available",
			"initial_user": "centos",
			"distribution": "centos",
			"os": "linux",
			"description": "CentOS 7 disk image",
			"label": "centos",
			"disk_image_url": "https://example.com/centos-7.img",
			"disk_image_size_bytes": 1073741824,
			"logo_url": "https://example.com/centos-logo.png",
			"created_at": "2023-01-01T00:00:00Z",
			"created_by": "system",
			"distribution_default": true
		}`})
	defer server.Close()

	results, err := client.GetDiskImage("b82168fe-66f6-4b38-a3b8-5283542d5475")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &DiskImage{
		ID:                  "b82168fe-66f6-4b38-a3b8-5283542d5475",
		Name:                "centos-7",
		Version:             "7",
		State:               "available",
		InitialUser:         "centos",
		Distribution:        "centos",
		OS:                  "linux",
		Description:         "CentOS 7 disk image",
		Label:               "centos",
		DiskImageURL:        "https://example.com/centos-7.img",
		DiskImageSizeBytes:  1073741824,
		LogoURL:             "https://example.com/centos-logo.png",
		CreatedAt:           time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedBy:           "system",
		DistributionDefault: true,
	}

	if !reflect.DeepEqual(results, expected) {
		t.Errorf("Expected %+v, got %+v", expected, results)
	}
}

func TestFindDiskImage(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/disk_images": `[{
			"id": "b82168fe-66f6-4b38-a3b8-52835428965",
			"name": "debian-10",
			"version": "10",
			"state": "available",
			"initial_user": "debian",
			"distribution": "debian",
			"os": "linux",
			"description": "Debian 10 disk image",
			"label": "debian",
			"disk_image_url": "https://example.com/debian-10.img",
			"disk_image_size_bytes": 2147483648,
			"logo_url": "https://example.com/debian-logo.png",
			"created_at": "2023-01-01T00:00:00Z",
			"created_by": "system",
			"distribution_default": true
		}]`})
	defer server.Close()

	results, err := client.FindDiskImage("debian-10")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &DiskImage{
		ID:                  "b82168fe-66f6-4b38-a3b8-52835428965",
		Name:                "debian-10",
		Version:             "10",
		State:               "available",
		InitialUser:         "debian",
		Distribution:        "debian",
		OS:                  "linux",
		Description:         "Debian 10 disk image",
		Label:               "debian",
		DiskImageURL:        "https://example.com/debian-10.img",
		DiskImageSizeBytes:  2147483648,
		LogoURL:             "https://example.com/debian-logo.png",
		CreatedAt:           time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedBy:           "system",
		DistributionDefault: true,
	}

	if !reflect.DeepEqual(results, expected) {
		t.Errorf("Expected %+v, got %+v", expected, results)
	}
}

func TestListDiskImages(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/disk_images": `[{"id":"22552dcf-aea3-4403-ae62-93651932bbd7","name":"centos-7","version":"7","state":"available","initial_user":"centos","distribution":"centos","os":"linux","description":"","label":"","disk_image_url":"https://example.com/centos-7.img","disk_image_size_bytes":1073741824,"logo_url":"https://example.com/centos-logo.png","created_at":"2023-01-01T00:00:00Z","created_by":"system","distribution_default":true},{"id":"4204229c-510c-4ba4-ab07-522e2aaa2cf8","name":"debian-10","version":"10","state":"available","initial_user":"","distribution":"debian","os":"","description":"","label":"","disk_image_url":"","disk_image_size_bytes":0,"logo_url":"","created_at":"0001-01-01T00:00:00Z","created_by":"","distribution_default":false},{"id":"cddce6c9-f84e-4e4f-ab8d-7a33cab85158","name":"debian-9","version":"9","state":"available","initial_user":"","distribution":"debian","os":"","description":"","label":"","disk_image_url":"","disk_image_size_bytes":0,"logo_url":"","created_at":"0001-01-01T00:00:00Z","created_by":"","distribution_default":false},{"id":"c3b28d45-c161-4abc-bdda-4facac38f2b1","name":"ubuntu-bionic","version":"18.04","state":"available","initial_user":"","distribution":"ubuntu","os":"","description":"","label":"","disk_image_url":"","disk_image_size_bytes":0,"logo_url":"","created_at":"0001-01-01T00:00:00Z","created_by":"","distribution_default":false},{"id":"8eb48e20-e5db-49fe-9cdf-cc8f381c61c6","name":"ubuntu-focal","version":"20.04","state":"available","initial_user":"","distribution":"ubuntu","os":"","description":"","label":"","disk_image_url":"","disk_image_size_bytes":0,"logo_url":"","created_at":"0001-01-01T00:00:00Z","created_by":"","distribution_default":false}]`,
	})
	defer server.Close()
	got, err := client.ListDiskImages()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []DiskImage{
		{
			ID:                  "22552dcf-aea3-4403-ae62-93651932bbd7",
			Name:                "centos-7",
			Version:             "7",
			State:               "available",
			InitialUser:         "centos",
			Distribution:        "centos",
			OS:                  "linux",
			Description:         "",
			Label:               "",
			DiskImageURL:        "https://example.com/centos-7.img",
			DiskImageSizeBytes:  1073741824,
			LogoURL:             "https://example.com/centos-logo.png",
			CreatedAt:           time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedBy:           "system",
			DistributionDefault: true,
		},
		{
			ID:           "4204229c-510c-4ba4-ab07-522e2aaa2cf8",
			Name:         "debian-10",
			Version:      "10",
			State:        "available",
			Distribution: "debian",
			Description:  "",
			Label:        "",
		},
		{
			ID:           "cddce6c9-f84e-4e4f-ab8d-7a33cab85158",
			Name:         "debian-9",
			Version:      "9",
			State:        "available",
			Distribution: "debian",
			Description:  "",
			Label:        "",
		},
		{
			ID:           "c3b28d45-c161-4abc-bdda-4facac38f2b1",
			Name:         "ubuntu-bionic",
			Version:      "18.04",
			State:        "available",
			Distribution: "ubuntu",
			Description:  "",
			Label:        "",
		},
		{
			ID:           "8eb48e20-e5db-49fe-9cdf-cc8f381c61c6",
			Name:         "ubuntu-focal",
			Version:      "20.04",
			State:        "available",
			Distribution: "ubuntu",
			Description:  "",
			Label:        "",
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetDiskImageByName(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/disk_images": `[{ "ID": "329d473e-f110-4852-b2fa-fe65aa6bff4a", "Name": "ubuntu-bionic", "Version": "18.04", "State": "available", "Distribution": "ubuntu", "Description": "", "Label": "bionic" }, { "ID": "77bea4dd-bfd4-492c-823d-f92eb6dd962d", "Name": "ubuntu-focal", "Version": "20.04", "State": "available", "Distribution": "ubuntu", "Description": "", "Label": "focal" }]`,
	})
	defer server.Close()

	got, err := client.GetDiskImageByName("ubuntu-bionic")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.ID != "329d473e-f110-4852-b2fa-fe65aa6bff4a" {
		t.Errorf("Expected %s, got %s", "329d473e-f110-4852-b2fa-fe65aa6bff4a", got.ID)
	}
}

func TestGetMostRecentDistro(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/disk_images": `[{ "ID": "329d473e-f110-4852-b2fa-fe65aa6bff4a", "Name": "ubuntu-bionic", "Version": "18.04", "State": "available", "Distribution": "ubuntu", "Description": "", "Label": "bionic" }, { "ID": "77bea4dd-bfd4-492c-823d-f92eb6dd962d", "Name": "ubuntu-focal", "Version": "20.04", "State": "available", "Distribution": "ubuntu", "Description": "", "Label": "focal" }]`,
	})
	defer server.Close()

	got, err := client.GetMostRecentDistro("ubuntu")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Name != "ubuntu-focal" {
		t.Errorf("Expected %s, got %s", "ubuntu-focal", got.Name)
	}
}
