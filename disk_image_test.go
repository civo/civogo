package civogo

import (
	"testing"
)

func TestClienterDiskImage(t *testing.T) {
	var c Clienter

	c, _ = NewClient("foo", "NYC1")
	c, _ = NewFakeClient()
	_, _ = c.ListDiskImages()
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
