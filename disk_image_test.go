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
