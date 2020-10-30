package civogo

import (
	"testing"
)

func TestListInstanceSizes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sizes": `[
			{
				"id": "d6b170f2-d2b3-4205-84c4-61898622393d",
				"name": "micro",
				"nice_name": "Micro",
				"cpu_cores": 1,
				"ram_mb": 1024,
				"disk_gb": 25,
				"description": "Micro - 1GB RAM, 1 CPU Core, 25GB SSD Disk",
				"selectable": true
			}
		]
		`,
	})
	defer server.Close()

	got, err := client.ListInstanceSizes()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].ID != "d6b170f2-d2b3-4205-84c4-61898622393d" {
		t.Errorf("Expected %s, got %s", "d6b170f2-d2b3-4205-84c4-61898622393d", got[0].ID)
	}
	if got[0].Name != "micro" {
		t.Errorf("Expected %s, got %s", "micro", got[0].Name)
	}
	if got[0].NiceName != "Micro" {
		t.Errorf("Expected %s, got %s", "Micro", got[0].NiceName)
	}
	if got[0].NiceName != "Micro" {
		t.Errorf("Expected %s, got %s", "Micro", got[0].NiceName)
	}
	if !got[0].Selectable {
		t.Errorf("Expected first result to be selectable")
	}
	if got[0].CPUCores != 1 {
		t.Errorf("Expected %d, got %d", 1, got[0].CPUCores)
	}
	if got[0].RAMMegabytes != 1024 {
		t.Errorf("Expected %d, got %d", 1024, got[0].RAMMegabytes)
	}
	if got[0].DiskGigabytes != 25 {
		t.Errorf("Expected %d, got %d", 25, got[0].DiskGigabytes)
	}
	if got[0].Description != "Micro - 1GB RAM, 1 CPU Core, 25GB SSD Disk" {
		t.Errorf("Expected %s, got %s", "Micro - 1GB RAM, 1 CPU Core, 25GB SSD Disk", got[0].Description)
	}
}

func TestFindInstanceSizes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sizes": `[
			{
				"id": "d6b170f2-d2b3-4205-84c4-61898622393d",
				"name": "debian-10",
				"nice_name": "Debian 10",
				"cpu_cores": 1,
				"ram_mb": 1024,
				"disk_gb": 25,
				"description": "Micro - 1GB RAM, 1 CPU Core, 25GB SSD Disk",
				"selectable": true
			},
			{
				"id": "456780f2-d3b3-345-84c4-12345622393d",
				"name": "debian-9",
				"nice_name": "Debian 9",
				"cpu_cores": 1,
				"ram_mb": 2024,
				"disk_gb": 50,
				"description": "Debian - 2GB RAM, 1 CPU Core, 50GB SSD Disk",
				"selectable": true
			}
		]
		`,
	})
	defer server.Close()

	got, _ := client.FindInstanceSizes("456780f2")
	if got.ID != "456780f2-d3b3-345-84c4-12345622393d" {
		t.Errorf("Expected %s, got %s", "456780f2-d3b3-345-84c4-12345622393d", got.ID)
	}

	got, _ = client.FindInstanceSizes("d6b170")
	if got.ID != "d6b170f2-d2b3-4205-84c4-61898622393d" {
		t.Errorf("Expected %s, got %s", "d6b170f2-d2b3-4205-84c4-61898622393d", got.ID)
	}

	got, _ = client.FindInstanceSizes("debian-1")
	if got.ID != "d6b170f2-d2b3-4205-84c4-61898622393d" {
		t.Errorf("Expected %s, got %s", "d6b170f2-d2b3-4205-84c4-61898622393d", got.ID)
	}

	got, _ = client.FindInstanceSizes("debian-9")
	if got.ID != "456780f2-d3b3-345-84c4-12345622393d" {
		t.Errorf("Expected %s, got %s", "456780f2-d3b3-345-84c4-12345622393d", got.ID)
	}

	_, err := client.FindInstanceSizes("debian")
	if err.Error() != "MultipleMatchesError: unable to find debian because there were multiple matches" {
		t.Errorf("Expected %s, got %s", "unable to find com because there were multiple matches", err.Error())
	}

	_, err = client.FindInstanceSizes("missing")
	if err.Error() != "ZeroMatchesError: unable to find missing, zero matches" {
		t.Errorf("Expected %s, got %s", "unable to find missing, zero matches", err.Error())
	}
}
