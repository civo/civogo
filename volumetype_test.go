package civogo

import (
	"reflect"
	"testing"
)

func TestListVolumeTypes(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/volumetypes": `[{
			"name": "my-volume-type",
			"description": "a volume type",
			"enabled": true,
			"labels": ["label"]
		}]`,
	})
	defer server.Close()

	got, err := client.ListVolumeTypes()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := []VolumeType{{
		Name:        "my-volume-type",
		Description: "a volume type",
		Enabled:     true,
		Labels:      []string{"label"},
	}}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}
