package civogo

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_AdvancedClientForTesting(t *testing.T) {
	g := NewGomegaWithT(t)

	serverValue := []ConfigAdvanceClientForTesting{
		{
			Method: "GET",
			Value: []ValueAdvanceClientForTesting{
				{
					RequestBody: "",
					URL:         "/v2/dns",
					ResponseBody: `[
						{"id": "12345", "account_id": "1", "name": "example.com"},
						{"id": "12346", "account_id": "1", "name": "example.net"}
						]`,
				},
				{
					RequestBody:  "",
					URL:          "/v2/dns/12345/records",
					ResponseBody: `[{"id": "1", "domain_id":"12345", "account_id": "1", "name": "txt", "type": "A", "value": "target", "ttl": 600}]`,
				},
			},
		},
	}

	client, server, _ := NewAdvancedClientForTesting(serverValue)
	defer server.Close()
	g.Expect(client.UserAgent).To(Equal("civogo/dev"))

	// Update the UserAgent
	client.SetUserAgent("civogo/test")
	g.Expect(client.UserAgent).To(Equal("civogo/test civogo/dev"))

	// Check the records for the domain
	records, err := client.ListDNSRecords("12345")
	g.Expect(err).To(BeNil())
	g.Expect(len(records)).To(Equal(1))

	// Check the doamins
	domains, err := client.ListDNSDomains()
	g.Expect(err).To(BeNil())
	g.Expect(len(domains)).To(Equal(2))

}
