package civogo_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/civo/civogo/v3"
	"github.com/civo/civogo/v3/instances"
)

func TestUpdateSendsExplicitFalse(t *testing.T) {
	var body string
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		io.WriteString(w, `{"instance":{"id":"i-1"}}`)
	})
	if _, err := client.Instances.Update(context.Background(), "i-1", instances.UpdateInstanceRequest{
		PublicIPRequired: civogo.Ptr(false),
		Notes:            civogo.Ptr(""),
	}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"public_ip_required":false`, `"notes":""`} {
		if !strings.Contains(body, want) {
			t.Errorf("body %s missing %s", body, want)
		}
	}
	if strings.Contains(body, "hostname") {
		t.Errorf("unset optional field sent: %s", body)
	}
}
