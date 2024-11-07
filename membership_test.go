package civogo

import "testing"

func TestListMemberships(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/memberships": `{
		"accounts": [
				{
					"id": "4f229791-6088-42dd-8fbe-3f64ec47567f",
					"api_key": "75e521bc74d34b21b42827fb58fcd590",
					"email_address": "team@civo.com",
					"label": "team testing",
					"organisation_id": "63bd4fe4-eeff-421b-aa24-1518decc5464"
				}
			],

				"organisations": [
				{
					"id": "63bd4fe4-eeff-421b-aa24-1518decc5464",
					"name": "Our-group Inc"
				}
			]
			
		}`,
	})
	defer server.Close()
	got, err := client.ListMemberships()

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	if got.Accounts[0].ID != "4f229791-6088-42dd-8fbe-3f64ec47567f" {
		t.Errorf("Expected User ID %s, got %s", "4f229791-6088-42dd-8fbe-3f64ec47567f", got.Accounts[0].ID)
	}

	if got.Accounts[0].OrganisationID != "63bd4fe4-eeff-421b-aa24-1518decc5464" {
		t.Errorf("Expected User ID %s, got %s", "63bd4fe4-eeff-421b-aa24-1518decc5464", got.Accounts[0].OrganisationID)
	}

	if got.Organisations[0].ID != "63bd4fe4-eeff-421b-aa24-1518decc5464" {
		t.Errorf("Expected User ID %s, got %s", "63bd4fe4-eeff-421b-aa24-1518decc5464", got.Organisations[0].ID)
	}

}
