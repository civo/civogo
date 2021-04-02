package civogo

import (
	"testing"
)

func TestClienter(t *testing.T) {
	var c Clienter

	c, _ = NewClient("foo", "NYC1")
	c, _ = NewFakeClient()
	_, _ = c.ListAllInstances()
}

func TestInstances(t *testing.T) {
	client, _ := NewFakeClient()

	config := &InstanceConfig{
		Count:    1,
		Hostname: "foo.example.com",
	}
	client.CreateInstance(config)

	results, err := client.ListInstances(1, 10)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if results.Page != 1 {
		t.Errorf("Expected %+v, got %+v", 1, results.Page)
		return
	}
	if results.Pages != 1 {
		t.Errorf("Expected %+v, got %+v", 1, results.Pages)
		return
	}
	if results.PerPage != 10 {
		t.Errorf("Expected %+v, got %+v", 10, results.PerPage)
		return
	}
	if len(results.Items) != 1 {
		t.Errorf("Expected %+v, got %+v", 1, len(results.Items))
		return
	}
}
