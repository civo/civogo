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
