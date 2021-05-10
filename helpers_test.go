package civogo

import (
	"testing"
)

// EnsureSuccessfulSimpleResponse simply takes a simple response and the client's last
// error and ensure the last request was successful
func EnsureSuccessfulSimpleResponse(t *testing.T, got *SimpleResponse, err error) {
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Result == "" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}
