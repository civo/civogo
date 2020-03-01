package civogo

import (
	"testing"
	"time"
)

func TestListCharges(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/charges": `[
			{
				"code": "instance-g1.small",
				"label": "furry-apple.example.com",
				"from": "2016-03-18T10:46:06Z",
				"to": "2016-03-25T10:46:06Z",
				"num_hours": 168,
				"size_gb": 200
			}
		]
		`,
	})
	defer server.Close()

	from, _ := time.Parse(time.RFC3339, "2016-03-01T00:00:00Z")
	to, _ := time.Parse(time.RFC3339, "2016-03-31T23:59:59Z")

	got, err := client.ListCharges(from, to)

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got[0].Code != "instance-g1.small" {
		t.Errorf("Expected %s, got %s", "instance-g1.small", got[0].Code)
	}

	if got[0].Label != "furry-apple.example.com" {
		t.Errorf("Expected %s, got %s", "furry-apple.example.com", got[0].Label)
	}

	test, _ := time.Parse(time.RFC3339, "2016-03-18T10:46:06Z")
	if got[0].From != test {
		t.Errorf("Expected %v, got %v", test, got[0].From)
	}

	test, _ = time.Parse(time.RFC3339, "2016-03-25T10:46:06Z")
	if got[0].To != test {
		t.Errorf("Expected %v, got %v", test, got[0].To)
	}

	if got[0].Label != "furry-apple.example.com" {
		t.Errorf("Expected %s, got %s", "furry-apple.example.com", got[0].Label)
	}

	if got[0].NumHours != 168 {
		t.Errorf("Expected %d, got %d", 168, got[0].NumHours)
	}

	if got[0].SizeGigabytes != 200 {
		t.Errorf("Expected %d, got %d", 200, got[0].SizeGigabytes)
	}
}
