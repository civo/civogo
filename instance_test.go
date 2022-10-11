package civogo

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"testing"
	. "github.com/onsi/gomega"
)

func TestInstances_List(t *testing.T) {
	g = NewGomegaWithT(t)
	initServer()
	defer downServer()

	mux.HandleFunc("/v2/instances", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Add("Content-Type", "application/json")
		allInstances := PaginatedInstancesList{
			Page:    1,
			PerPage: 10,
			Pages:   2,
			Items:   []Instance{
				{
					ID:                       "3e80c90e-1ca2-4221-8820-5aaf2c048378",
					Hostname:                 "test.com",
					NetworkID:                "4f9f972e-6031-4802-83f7-460b0d840ba0",
				},
				{
					ID:                       "c1f43c41-8630-4f7f-a946-b78d421a87bc",
					Hostname:                 "test1.com",
					NetworkID:                "4f9f972e-6031-4802-83f7-460b0d840ba0",
				},
				{
					ID:                       "6326019c-2e49-4638-b52a-c00690f60f1b",
					Hostname:                 "test2.com",
					NetworkID:                "b214410d-9ddf-4880-8a4f-7b62b03f32ca",
				},
			},
		}
		value := toJSON(t, allInstances)
		fmt.Fprint(w, value)
	})

	keys, meta, err := client.Instances("").List(ctx)
	g.Expect(err).To(BeNil())

	expectedKeys := PaginatedInstancesList{
		Page:    1,
		PerPage: 10,
		Pages:   2,
		Items:   []Instance{
			{
				ID:                       "3e80c90e-1ca2-4221-8820-5aaf2c048378",
				Hostname:                 "test.com",
				NetworkID:                "4f9f972e-6031-4802-83f7-460b0d840ba0",
			},
			{
				ID:                       "c1f43c41-8630-4f7f-a946-b78d421a87bc",
				Hostname:                 "test1.com",
				NetworkID:                "4f9f972e-6031-4802-83f7-460b0d840ba0",
			},
			{
				ID:                       "6326019c-2e49-4638-b52a-c00690f60f1b",
				Hostname:                 "test2.com",
				NetworkID:                "b214410d-9ddf-4880-8a4f-7b62b03f32ca",
			},
		},
	}
	g.Expect(keys).To(Equal(expectedKeys))
	g.Expect(meta.StatusCode).To(Equal(http.StatusOK))
}
