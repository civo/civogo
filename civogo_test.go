package civogo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	ctx    = context.TODO()
	client *Client
	server *httptest.Server
)

func initServer() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewClient("test", "TEST")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func downServer() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

// toJSON converts a struct to JSON
func toJSON(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Error marshalling JSON: %v", err)
	}
	return string(b)
}
