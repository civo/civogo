package civogo

import (
	"reflect"
	"testing"
)

func TestListApplications(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications": `{"page":1,"per_page":20,"pages":1,"items":[{
		  "id": "69a23478-a89e-41d2-97b1-6f4c341cee70",
		  "name": "your-app-name",
		  "status": "ACTIVE",
		  "account_id": "12345",
		  "network_id": "34567",
		  "description": "this is a test app",
		  "process_info": [
			  	{
					"processType": "web",
					"processCount": 1
				  }],
			"domains": [
				"your-app-name.example.com"
			],
			"ssh_key_ids": [
				"12345"
			],
			"config": [
				{
					"name": "PORT",
					"value": "80"
				}]
				}]}`,
	})

	defer server.Close()
	got, err := client.ListApplications()
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}

	expected := &PaginatedApplications{
		Page:    1,
		PerPage: 20,
		Pages:   1,
		Items: []Application{
			{
				ID:          "69a23478-a89e-41d2-97b1-6f4c341cee70",
				Name:        "your-app-name",
				Status:      "ACTIVE",
				NetworkID:   "34567",
				Description: "this is a test app",
				ProcessInfo: []ProcessInfo{
					{
						ProcessType:  "web",
						ProcessCount: 1,
					},
				},
				Domains:   []string{"your-app-name.example.com"},
				SSHKeyIDs: []string{"12345"},
				Config: []EnvVar{
					{
						Name:  "PORT",
						Value: "80",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestCreateApplication(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications": `{"name":"test-app"}`,
	})
	defer server.Close()

	cfg := &ApplicationConfig{
		Name:        "test-app",
		Description: "test app",
		SSHKeyIDs:   []string{"12345"},
		Size:        "small",
	}

	got, err := client.CreateApplication(cfg)
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Name != "test-app" {
		t.Errorf("Expected %s, got %s", "test-app", got.Name)
	}
}

func TestDeleteApplication(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v2/applications/12345": `{"result":"success"}`,
	})
	defer server.Close()

	got, err := client.DeleteApplication("12345")
	if err != nil {
		t.Errorf("Request returned an error: %s", err)
		return
	}
	if got.Result != "success" {
		t.Errorf("Expected %s, got %s", "success", got.Result)
	}
}
