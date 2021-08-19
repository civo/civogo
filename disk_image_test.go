package civogo

import (
	"reflect"
	"testing"
)

func TestClienterDiskImage(t *testing.T) {
	var c Clienter

	c, _ = NewClient("foo", "NYC1")
	c, _ = NewFakeClient()
	_, _ = c.ListDiskImages()
}

func TestListDiskImages(t *testing.T) {
	client, _ := NewFakeClient()

	results, err := client.ListDiskImages()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if len(results) != 4 {
		t.Errorf("Expected %+v, got %+v", 4, len(results))
		return
	}

}

func TestGetDiskImage(t *testing.T) {
	client, _ := NewFakeClient()

	results, err := client.GetDiskImage("b82168fe-66f6-4b38-a3b8-5283542d5475")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if results.ID != "b82168fe-66f6-4b38-a3b8-5283542d5475" {
		t.Errorf("Expected %+v, got %+v", "b82168fe-66f6-4b38-a3b8-5283542d5475", results.ID)
		return
	}

	if results.Name != "centos-7" {
		t.Errorf("Expected %+v, got %+v", "centos-7", results.Name)
		return
	}

}

func TestFindDiskImage(t *testing.T) {
	client, _ := NewFakeClient()

	results, err := client.FindDiskImage("debian-10")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if results.ID != "b82168fe-66f6-4b38-a3b8-52835428965" {
		t.Errorf("Expected %+v, got %+v", "b82168fe-66f6-4b38-a3b8-52835428965", results.ID)
		return
	}

	if results.Name != "debian-10" {
		t.Errorf("Expected %+v, got %+v", "debian-10", results.Name)
		return
	}

}

func TestListDiskImagesV2(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/disk_images": `[{"id":"22552dcf-aea3-4403-ae62-93651932bbd7","name":"centos-7","version":"7","state":"available","distribution":"centos","description":null,"label":null},{"id":"4204229c-510c-4ba4-ab07-522e2aaa2cf8","name":"debian-10","version":"10","state":"available","distribution":"debian","description":null,"label":null},{"id":"cddce6c9-f84e-4e4f-ab8d-7a33cab85158","name":"debian-9","version":"9","state":"available","distribution":"debian","description":null,"label":null},{"id":"c3b28d45-c161-4abc-bdda-4facac38f2b1","name":"ubuntu-bionic","version":"18.04","state":"available","distribution":"ubuntu","description":null,"label":null},{"id":"8eb48e20-e5db-49fe-9cdf-cc8f381c61c6","name":"ubuntu-focal","version":"20.04","state":"available","distribution":"ubuntu","description":null,"label":null}]`,
	})
	defer server.Close()
	got, err := client.ListDiskImages()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	expected := []DiskImage{
		{
			ID:           "22552dcf-aea3-4403-ae62-93651932bbd7",
			Name:         "centos-7",
			Version:      "7",
			State:        "available",
			Distribution: "centos",
			Description:  "",
			Label:        "",
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
