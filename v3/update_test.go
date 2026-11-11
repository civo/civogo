package civogo_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/civo/civogo/v3"
	"github.com/civo/civogo/v3/compute"
)

func TestUpdateSendsExplicitFalse(t *testing.T) {
	var body string
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		io.WriteString(w, `{"instance":{"id":"i-1"}}`)
	})
	if _, err := client.Compute.Instances.Update(context.Background(), "i-1", compute.UpdateInstanceRequest{
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
	// Unset fields must stay off the wire: an empty placement_rule would
	// make the server wipe the instance's placement constraints.
	for _, unset := range []string{"hostname", "placement_rule"} {
		if strings.Contains(body, unset) {
			t.Errorf("unset field %s sent: %s", unset, body)
		}
	}
}
