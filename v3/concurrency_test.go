package civogo_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	civogo "github.com/civo/civogo/v3"
	"github.com/civo/civogo/v3/api"
	"github.com/civo/civogo/v3/compute"
)

func TestClientIsConcurrencySafe(t *testing.T) {
	client := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"instance": {"id": "i-1"}, "instances": [], "page_info": {"total_pages": 1}}`)
	}, civogo.WithDefaultRegion("lon1"))

	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctx := context.Background()
			switch i % 3 {
			case 0:
				if _, err := client.Compute.Instances.Get(ctx, "i-1", api.WithRegion("fra1")); err != nil {
					t.Error(err)
				}
			case 1:
				if _, err := client.Compute.Instances.List(ctx, compute.ListInstancesRequest{}); err != nil {
					t.Error(err)
				}
			case 2:
				if err := client.Compute.Instances.Stop(ctx, "i-1"); err != nil {
					t.Error(err)
				}
			}
		}(i)
	}
	wg.Wait()
}
