package civogo

import (
	"testing"
)

func TestListInstanceSizes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sizes": `[
			{
        "type": "Instance",
        "name": "g4s.xsmall",
        "nice_name": "xSmall - Standard",
        "cpu_cores": 1,
        "gpu_count": 0,
        "gpu_type": "",
        "ram_mb": 1024,
        "disk_gb": 25,
        "transfer_tb": 1,
        "description": "xSmall - Standard",
        "selectable": true,
        "price_monthly": 5.00,
        "price_hourly": 0.00684
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
	if got[0].Name != "g4s.xsmall" {
		t.Errorf("Expected %s, got %s", "g4s.xsmall", got[0].Name)
	}
	if got[0].Type != "Instance" {
		t.Errorf("Expected %s, got %s", "Instance", got[0].Type)
	}
	if got[0].NiceName != "xSmall - Standard" {
		t.Errorf("Expected %s, got %s", "xSmall - Standard", got[0].NiceName)
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
	if got[0].Description != "xSmall - Standard" {
		t.Errorf("Expected %s, got %s", "xSmall - Standard", got[0].Description)
	}
	if got[0].PriceMonthly != 5.00 {
		t.Errorf("Expected monthly price %f, got %f", 5.00, got[0].PriceMonthly)
	}
	if got[0].PriceHourly != 0.00684 {
		t.Errorf("Expected hourly price %f, got %f", 0.00684, got[0].PriceHourly)
	}
}

func TestFindInstanceSizes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/sizes": `[
			{
        "type": "Instance",
        "name": "g3.xsmall",
        "nice_name": "Extra Small",
        "cpu_cores": 1,
        "gpu_count": 0,
        "gpu_type": "",
        "ram_mb": 1024,
        "disk_gb": 25,
        "transfer_tb": 1,
        "description": "Extra Small",
        "selectable": true,
        "price_monthly": 5.00,
        "price_hourly": 0.00684
    },
    {
        "type": "Instance",
        "name": "g3.small",
        "nice_name": "Small",
        "cpu_cores": 1,
        "gpu_count": 0,
        "gpu_type": "",
        "ram_mb": 2048,
        "disk_gb": 25,
        "transfer_tb": 2,
        "description": "Small",
        "selectable": true,
        "price_monthly": 10.00,
        "price_hourly": 0.01369
    }
		]
		`,
	})
	defer server.Close()

	got, _ := client.FindInstanceSizes("g3.small")
	if got.Name != "g3.small" {
		t.Errorf("Expected %s, got %s", "g3.small", got.Name)
	}

	got, _ = client.FindInstanceSizes("xsmall")
	if got.Name != "g3.xsmall" {
		t.Errorf("Expected %s, got %s", "g3.xsmall", got.Name)
	}
}
