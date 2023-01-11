package civogo

import (
	"testing"
)

func TestGetQuota(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/quota": `{
			"id": "44aab548-61ca-11e5-860e-5cf9389be614",
			"default_user_id": "ca04ddda-06e1-469a-ad63-27ac3298c42c",
			"default_user_email_address": "johnsmith@example.com",
			"instance_count_limit": 16,
			"instance_count_usage": 6,
			"cpu_core_limit": 10,
			"cpu_core_usage": 3,
			"ram_mb_limit": 5120,
			"ram_mb_usage": 1536,
			"disk_gb_limit": 250,
			"disk_gb_usage": 75,
			"disk_volume_count_limit": 16,
			"disk_volume_count_usage": 6,
			"disk_snapshot_count_limit": 30,
			"disk_snapshot_count_usage": 0,
			"public_ip_address_limit": 16,
			"public_ip_address_usage": 6,
			"subnet_count_limit": 10,
			"subnet_count_usage": 1,
			"network_count_limit": 10,
			"network_count_usage": 1,
			"security_group_limit": 16,
			"security_group_usage": 5,
			"security_group_rule_limit": 160,
			"security_group_rule_usage": 24,
			"port_count_limit": 32,
			"port_count_usage": 7,
			"loadbalancer_count_limit": 16,
			"loadbalancer_count_usage": 1,
			"objectstore_gb_limit": 1000,
			"objectstore_gb_usage": 0,
			"database_count_limit": 4,
			"database_count_usage": 0,
			"database_snapshot_count_limit": 20,
			"database_snapshot_count_usage": 0,
			"database_cpu_core_limit": 120,
			"database_cpu_core_usage": 0,
			"database_ram_mb_limit": 786432,
			"database_ram_mb_usage": 0,
			"database_disk_gb_limit": 7680,
			"database_disk_gb_usage": 0
		}`,
	})
	defer server.Close()

	got, err := client.GetQuota()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.ID != "44aab548-61ca-11e5-860e-5cf9389be614" {
		t.Errorf("Expected %s, got %s", "44aab548-61ca-11e5-860e-5cf9389be614", got.ID)
	}
	if got.DefaultUserID != "ca04ddda-06e1-469a-ad63-27ac3298c42c" {
		t.Errorf("Expected %s, got %s", "ca04ddda-06e1-469a-ad63-27ac3298c42c", got.DefaultUserID)
	}
	if got.DefaultUserEmailAddress != "johnsmith@example.com" {
		t.Errorf("Expected %s, got %s", "johnsmith@example.com", got.DefaultUserEmailAddress)
	}
	if got.InstanceCountLimit != 16 {
		t.Errorf("Expected %d, got %d", 16, got.InstanceCountLimit)
	}
	if got.InstanceCountUsage != 6 {
		t.Errorf("Expected %d, got %d", 6, got.InstanceCountUsage)
	}
	if got.CPUCoreLimit != 10 {
		t.Errorf("Expected %d, got %d", 10, got.CPUCoreLimit)
	}
	if got.CPUCoreUsage != 3 {
		t.Errorf("Expected %d, got %d", 3, got.CPUCoreUsage)
	}
	if got.RAMMegabytesLimit != 5120 {
		t.Errorf("Expected %d, got %d", 5120, got.RAMMegabytesLimit)
	}
	if got.RAMMegabytesUsage != 1536 {
		t.Errorf("Expected %d, got %d", 1536, got.RAMMegabytesUsage)
	}
	if got.DiskGigabytesLimit != 250 {
		t.Errorf("Expected %d, got %d", 250, got.DiskGigabytesLimit)
	}
	if got.DiskGigabytesUsage != 75 {
		t.Errorf("Expected %d, got %d", 75, got.DiskGigabytesUsage)
	}
	if got.DiskVolumeCountLimit != 16 {
		t.Errorf("Expected %d, got %d", 16, got.DiskVolumeCountLimit)
	}
	if got.DiskVolumeCountUsage != 6 {
		t.Errorf("Expected %d, got %d", 6, got.DiskVolumeCountUsage)
	}
	if got.DiskSnapshotCountLimit != 30 {
		t.Errorf("Expected %d, got %d", 30, got.DiskSnapshotCountLimit)
	}
	if got.DiskSnapshotCountUsage != 0 {
		t.Errorf("Expected %d, got %d", 0, got.DiskSnapshotCountUsage)
	}
	if got.PublicIPAddressLimit != 16 {
		t.Errorf("Expected %d, got %d", 16, got.PublicIPAddressLimit)
	}
	if got.PublicIPAddressUsage != 6 {
		t.Errorf("Expected %d, got %d", 6, got.PublicIPAddressUsage)
	}
	if got.SubnetCountLimit != 10 {
		t.Errorf("Expected %d, got %d", 10, got.SubnetCountLimit)
	}
	if got.SubnetCountUsage != 1 {
		t.Errorf("Expected %d, got %d", 1, got.SubnetCountUsage)
	}
	if got.NetworkCountLimit != 10 {
		t.Errorf("Expected %d, got %d", 10, got.NetworkCountLimit)
	}
	if got.NetworkCountUsage != 1 {
		t.Errorf("Expected %d, got %d", 1, got.NetworkCountUsage)
	}
	if got.SecurityGroupLimit != 16 {
		t.Errorf("Expected %d, got %d", 16, got.SecurityGroupLimit)
	}
	if got.SecurityGroupUsage != 5 {
		t.Errorf("Expected %d, got %d", 5, got.SecurityGroupUsage)
	}
	if got.SecurityGroupRuleLimit != 160 {
		t.Errorf("Expected %d, got %d", 160, got.SecurityGroupRuleLimit)
	}
	if got.SecurityGroupRuleUsage != 24 {
		t.Errorf("Expected %d, got %d", 24, got.SecurityGroupRuleUsage)
	}
	if got.PortCountLimit != 32 {
		t.Errorf("Expected %d, got %d", 32, got.PortCountLimit)
	}
	if got.PortCountUsage != 7 {
		t.Errorf("Expected %d, got %d", 7, got.PortCountUsage)
	}
	if got.LoadBalancerCountLimit != 16 {
		t.Errorf("Expected %d, got %d", 16, got.LoadBalancerCountLimit)
	}
	if got.LoadBalancerCountUsage != 1 {
		t.Errorf("Expected %d, got %d", 1, got.LoadBalancerCountUsage)
	}
	if got.ObjectStoreGigabytesLimit != 1000 {
		t.Errorf("Expected %d, got %d", 1000, got.ObjectStoreGigabytesLimit)
	}
	if got.ObjectStoreGigabytesUsage != 0 {
		t.Errorf("Expected %d, got %d", 0, got.ObjectStoreGigabytesUsage)
	}
	if got.DatabaseCountLimit != 4 {
		t.Errorf("Expected %d, got %d", 4, got.DatabaseCountLimit)
	}
	if got.DatabaseCountUsage != 0 {
		t.Errorf("Expected %d, got %d", 0, got.DatabaseCountUsage)
	}
	if got.DatabaseSnapshotCountLimit != 20 {
		t.Errorf("Expected %d, got %d", 20, got.DatabaseSnapshotCountLimit)
	}
	if got.DatabaseSnapshotCountUsage != 0 {
		t.Errorf("Expected %d, got %d", 0, got.DatabaseSnapshotCountUsage)
	}
	if got.DatabaseCPUCoreLimit != 120 {
		t.Errorf("Expected %d, got %d", 120, got.DatabaseCPUCoreLimit)
	}
	if got.DatabaseCPUCoreUsage != 0 {
		t.Errorf("Expected %d, got %d", 0, got.DatabaseCPUCoreUsage)
	}
	if got.DatabaseRAMMegabytesLimit != 786432 {
		t.Errorf("Expected %d, got %d", 786432, got.DatabaseRAMMegabytesLimit)
	}
	if got.DatabaseRAMMegabytesUsage != 0 {
		t.Errorf("Expected %d, got %d", 0, got.DatabaseRAMMegabytesUsage)
	}
	if got.DatabaseDiskGigabytesLimit != 7680 {
		t.Errorf("Expected %d, got %d", 7680, got.DatabaseDiskGigabytesLimit)
	}
	if got.DatabaseDiskGigabytesUsage != 0 {
		t.Errorf("Expected %d, got %d", 0, got.DatabaseDiskGigabytesUsage)
	}
}
