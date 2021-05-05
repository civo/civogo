package civogo

import (
	"testing"
)

func TestListRegions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/regions": `[{"code":"NYC1","name":"New York 1","type":"civostack","out_of_capacity":false,"country":"us","country_name":"United States","features":{"iaas":false,"kubernetes":true}},{"code":"SVG1","name":"Stevenage 1","default":true,"type":"openstack","out_of_capacity":true,"country":"uk","country_name":"United Kingdom","features":{"iaas":true,"kubernetes":true}}]`,
	})
	defer server.Close()

	got, err := client.ListRegions()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got[0].Code != "NYC1" {
		t.Errorf("Expected %s, got %s", "NYC1", got[0].Code)
	}
	if got[0].Name != "New York 1" {
		t.Errorf("Expected %s, got %s", "New York 1", got[0].Name)
	}
}

func TestFindRegions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/regions": `[{"code":"NYC1","name":"New York 1","type":"civostack","out_of_capacity":false,"country":"us","country_name":"United States","features":{"iaas":false,"kubernetes":true}},{"code":"SVG1","name":"Stevenage 1","default":true,"type":"openstack","out_of_capacity":true,"country":"uk","country_name":"United Kingdom","features":{"iaas":true,"kubernetes":true}}]`,
	})
	defer server.Close()

	got, err := client.FindRegion("nyc1")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Code != "NYC1" {
		t.Errorf("Expected %s, got %s", "NYC1", got.Code)
	}
	if got.Name != "New York 1" {
		t.Errorf("Expected %s, got %s", "New York 1", got.Name)
	}
}
