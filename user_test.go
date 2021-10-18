package civogo

import (
	"testing"
)

func TestGetUserEverything(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/users/123/everything": `{
			"user": {
				"id": "be6de2eb-540f-4d19-ab54-239b8333fafe",
				"first_name": "Test",
				"last_name": "User",
				"created_at": "2021-10-18T07:44:37+01:00",
				"updated_at": "2021-10-18T07:44:37+01:00",
				"deleted_at": null,
				"company_name": "Civo Ltd",
				"email_address": "test.user@example.com",
				"status": null,
				"flags": "group/live",
				"token": "c92b0064eb9245fb91f1e8a207cbb55e",
				"marketing_allowed": 1,
				"default_account_id": "4f229791-6088-42dd-8fbe-3f64ec47567f",
				"password_digest": "",
				"partner": null,
				"partner_user_id": null,
				"referral_id": null,
				"last_chosen_region": ""
			},
			"accounts": [
				{
					"id": "4f229791-6088-42dd-8fbe-3f64ec47567f",
					"username": null,
					"password": "",
					"api_key": "75e521bc74d34b21b42827fb58fcd590",
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00",
					"deleted_at": null,
					"email_address": "team@civo.com",
					"label": "team testing",
					"flags": null,
					"user_username": "",
					"user_password": null,
					"signup_number": 0,
					"salt": null,
					"timezone": "Europe/London",
					"partner": null,
					"partner_user_id": null,
					"status": null,
					"token": null,
					"email_confirmed": false,
					"credit_card_added": false,
					"staff_fraud_comments": null,
					"staff_fraud_comments_user_id": null,
					"organisation_id": "63bd4fe4-eeff-421b-aa24-1518decc5464",
					"enabled": false
				}
			],
			"organisations": [
				{
					"id": "63bd4fe4-eeff-421b-aa24-1518decc5464",
					"name": "Our-group Inc",
					"token": "7c6a21dc4f894b7282854fd23df4ea56",
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00"
				}
			],
			"teams": [
				{
					"id": "dd832bff-e86b-4c8f-b024-dc133212496e",
					"name": "Kube department",
					"organisation_id": "",
					"account_id": "4f229791-6088-42dd-8fbe-3f64ec47567f",
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00"
				}
			],
			"roles": [
				{
					"id": "04e8c8cb-1592-483d-8b10-6909cd0fd957",
					"name": "Billing administrator",
					"permissions": "billing.*",
					"organisation_id": "63bd4fe4-eeff-421b-aa24-1518decc5464",
					"account_id": "",
					"built_in": false,
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00"
				},
				{
					"id": "26010610-7015-4983-b10d-9e3affa4a710",
					"name": "Compute administrator",
					"permissions": "compute.*",
					"organisation_id": "",
					"account_id": "4f229791-6088-42dd-8fbe-3f64ec47567f",
					"built_in": false,
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00"
				},
				{
					"id": "6e7c6845-949e-4c94-8533-c3e27846e365",
					"name": "Compute viewer",
					"permissions": "compute.viewer",
					"organisation_id": "",
					"account_id": "",
					"built_in": true,
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00"
				},
				{
					"id": "152e7f6d-560a-49ca-b168-8053369f8ffb",
					"name": "Kubernetes administrator",
					"permissions": "kubernetes.*",
					"organisation_id": "",
					"account_id": "",
					"built_in": true,
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00"
				},
				{
					"id": "ebccdc4b-0356-472e-bf29-a6b1361d3b77",
					"name": "Owner",
					"permissions": "*.*",
					"organisation_id": "",
					"account_id": "",
					"built_in": true,
					"created_at": "2021-10-18T07:44:37+01:00",
					"updated_at": "2021-10-18T07:44:37+01:00"
				}
			]
		}`,
	})
	defer server.Close()
	got, err := client.GetUserEverything("123")

	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.User.ID != "be6de2eb-540f-4d19-ab54-239b8333fafe" {
		t.Errorf("Expected User ID %s, got %s", "be6de2eb-540f-4d19-ab54-239b8333fafe", got.User.ID)
	}
	if got.Accounts[0].ID != "4f229791-6088-42dd-8fbe-3f64ec47567f" {
		t.Errorf("Expected User ID %s, got %s", "4f229791-6088-42dd-8fbe-3f64ec47567f", got.Accounts[0].ID)
	}
	if got.Organisations[0].ID != "63bd4fe4-eeff-421b-aa24-1518decc5464" {
		t.Errorf("Expected User ID %s, got %s", "63bd4fe4-eeff-421b-aa24-1518decc5464", got.Organisations[0].ID)
	}
	if got.Teams[0].ID != "dd832bff-e86b-4c8f-b024-dc133212496e" {
		t.Errorf("Expected User ID %s, got %s", "dd832bff-e86b-4c8f-b024-dc133212496e", got.Teams[0].ID)
	}
	if got.Roles[1].ID != "26010610-7015-4983-b10d-9e3affa4a710" {
		t.Errorf("Expected User ID %s, got %s", "26010610-7015-4983-b10d-9e3affa4a710", got.Roles[1].ID)
	}
}
