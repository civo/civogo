package civogo

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

const v3alpha1DatabaseJSON = `{
  "id": "12345",
  "name": "test-db",
  "nodes": 3,
  "size": "g3.db.small",
  "software": "PostgreSQL",
  "softwareVersion": "17",
  "publicIpv4": "1.2.3.4",
  "privateIpv4": "10.0.0.4",
  "port": 5432,
  "username": "civo",
  "password": "secret",
  "databaseUserInfo": [{"username": "civo", "password": "secret", "port": 5432}],
  "networkId": "net-1",
  "firewallId": "fw-1",
  "backupConfigured": true,
  "status": "Ready",
  "createdAt": "2026-10-20T08:55:00Z",
  "dnsEntry": "12345.db.civo.com",
  "region": "LON1",
  "readEndpoint": {"host": "5.6.7.8", "port": 5433, "dnsEntry": "12345-ro.db.civo.com"},
  "storageGb": 120,
  "baseStorageGb": 40,
  "extensions": ["pg_cron", "postgis"],
  "databaseParameters": {"max_connections": "200"},
  "earliestRestoreTime": "2026-10-21T00:00:00Z"
}`

func v3alpha1ExpectedDatabase() *Database {
	created := time.Date(2026, 10, 20, 8, 55, 0, 0, time.UTC)
	earliest := time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)
	return &Database{
		ID:                  "12345",
		Name:                "test-db",
		Nodes:               3,
		Size:                "g3.db.small",
		Software:            "PostgreSQL",
		SoftwareVersion:     "17",
		PublicIPv4:          "1.2.3.4",
		PrivateIPv4:         "10.0.0.4",
		NetworkID:           "net-1",
		FirewallID:          "fw-1",
		Port:                5432,
		Username:            "civo",
		Password:            "secret",
		DatabaseUserInfo:    []DatabaseUserInfo{{Username: "civo", Password: "secret", Port: 5432}},
		DNSEntry:            "12345.db.civo.com",
		Status:              "Ready",
		Region:              "LON1",
		BackupConfigured:    true,
		CreatedAt:           &created,
		ReadEndpoint:        &DatabaseReadEndpoint{Host: "5.6.7.8", Port: 5433, DNSEntry: "12345-ro.db.civo.com"},
		StorageGB:           120,
		BaseStorageGB:       40,
		Extensions:          []string{"pg_cron", "postgis"},
		DatabaseParameters:  map[string]string{"max_connections": "200"},
		EarliestRestoreTime: &earliest,
	}
}

func TestGetDatabaseV3alpha1(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/12345": `{"database": ` + v3alpha1DatabaseJSON + `}`,
	})
	defer server.Close()

	got, err := client.GetDatabaseV3alpha1("12345")
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if !reflect.DeepEqual(got, v3alpha1ExpectedDatabase()) {
		t.Errorf("Expected %+v, got %+v", v3alpha1ExpectedDatabase(), got)
	}
}

func TestGetDatabaseV3alpha1WithoutOptionalFields(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/12345": `{"database": {"id": "12345", "name": "test-db", "nodes": 1, "status": "Pending", "extensions": [], "databaseParameters": {}}}`,
	})
	defer server.Close()

	got, err := client.GetDatabaseV3alpha1("12345")
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.ReadEndpoint != nil || got.EarliestRestoreTime != nil || got.CreatedAt != nil {
		t.Errorf("Expected nil optional fields, got %+v", got)
	}
	if len(got.Extensions) != 0 || len(got.DatabaseParameters) != 0 {
		t.Errorf("Expected empty extensions and parameters, got %+v", got)
	}
}

func TestListDatabasesV3alpha1(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases": `{"databases": [` + v3alpha1DatabaseJSON + `], "pageInfo": {"page": 1, "totalPages": 1, "perPage": 100, "itemsCount": 1, "totalItemsCount": 1}}`,
	})
	defer server.Close()

	got, err := client.ListDatabasesV3alpha1()
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	expected := &PaginatedDatabases{Page: 1, Pages: 1, PerPage: 1, Items: []Database{*v3alpha1ExpectedDatabase()}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestForkDatabase(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases",
					RequestBody:  `{"name":"test-db-fork","size":"g3.db.small","nodes":1,"source_database_id":"12345","source_target_time":"2026-10-21T10:00:00Z","region":"lon1"}`,
					ResponseBody: `{"database": {"id": "67890", "name": "test-db-fork", "nodes": 1, "status": "Pending"}}`,
				},
			},
		},
	})
	defer server.Close()

	target := time.Date(2026, 10, 21, 10, 0, 0, 0, time.UTC)
	got, err := client.ForkDatabase(&ForkDatabaseRequest{
		Name:             "test-db-fork",
		Size:             "g3.db.small",
		Nodes:            1,
		SourceDatabaseID: "12345",
		SourceTargetTime: &target,
		Region:           "LON1",
	})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.ID != "67890" || got.Name != "test-db-fork" || got.Status != "Pending" {
		t.Errorf("Unexpected database %+v", got)
	}
}

func TestRestoreDatabaseToPointInTime(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345/restore",
					RequestBody:  `{"name":"test-db-pitr","target_time":"2026-10-21T10:00:00Z","region":"lon1"}`,
					ResponseBody: `{"database": {"id": "12345", "name": "test-db", "status": "Restoring"}}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.RestoreDatabaseToPointInTime("12345", &RestoreDatabaseToPointInTimeRequest{
		Name:       "test-db-pitr",
		TargetTime: time.Date(2026, 10, 21, 10, 0, 0, 0, time.UTC),
		Region:     "LON1",
	})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.Status != "Restoring" {
		t.Errorf("Expected status Restoring, got %+v", got)
	}
}

func TestResizeDatabaseStorage(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PATCH",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345",
					RequestBody:  `{"storage_gb":120,"region":"test"}`,
					ResponseBody: `{"database": {"id": "12345", "storageGb": 120, "baseStorageGb": 40}}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.ResizeDatabaseStorage("12345", 120)
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.StorageGB != 120 || got.BaseStorageGB != 40 {
		t.Errorf("Unexpected storage %+v", got)
	}
}

func TestUpdateDatabaseExtensions(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PATCH",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345",
					RequestBody:  `{"update_extensions":true,"region":"test"}`,
					ResponseBody: `{"database": {"id": "12345", "extensions": []}}`,
				},
			},
		},
	})
	defer server.Close()

	// An empty list must still send update_extensions so the server clears them.
	got, err := client.UpdateDatabaseExtensions("12345", nil)
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if len(got.Extensions) != 0 {
		t.Errorf("Expected no extensions, got %+v", got.Extensions)
	}
}

func TestUpdateDatabaseParameters(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PATCH",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345",
					RequestBody:  `{"database_parameters":{"max_connections":"200"},"update_database_parameters":true,"region":"test"}`,
					ResponseBody: `{"database": {"id": "12345", "databaseParameters": {"max_connections": "200"}}}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.UpdateDatabaseParameters("12345", map[string]string{"max_connections": "200"})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.DatabaseParameters["max_connections"] != "200" {
		t.Errorf("Unexpected parameters %+v", got.DatabaseParameters)
	}
}

func TestListDatabaseExtensions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/extensions": `{"extensions": [{"name": "pg_cron", "description": "Job scheduler", "versions": ["1.6"], "requiresRestart": true}, {"name": "postgis", "description": "GIS", "versions": [], "requiresRestart": false}]}`,
	})
	defer server.Close()

	got, err := client.ListDatabaseExtensions()
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	expected := []DatabaseExtension{
		{Name: "pg_cron", Description: "Job scheduler", Versions: []string{"1.6"}, RequiresRestart: true},
		{Name: "postgis", Description: "GIS", Versions: []string{}, RequiresRestart: false},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListDatabaseBackupDestinationRegions(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/backup-destination-regions": `{"regions": [{"code": "LON1", "name": "London 1", "objectStore": true, "isCurrent": true}, {"code": "FRA1", "name": "Frankfurt 1", "objectStore": true, "isCurrent": false}]}`,
	})
	defer server.Close()

	got, err := client.ListDatabaseBackupDestinationRegions()
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	expected := []DatabaseBackupDestinationRegion{
		{Code: "LON1", Name: "London 1", ObjectStore: true, IsCurrent: true},
		{Code: "FRA1", Name: "Frankfurt 1", ObjectStore: true, IsCurrent: false},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListDatabaseBackupsV3alpha1(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/12345/backups": `{"backups": [
			{"id": "b-1", "name": "nightly", "software": "PostgreSQL", "state": "DATABASE_BACKUP_STATE_READY", "schedule": "0 0 * * *", "databaseName": "test-db", "databaseId": "12345", "isScheduled": true, "createdAt": "2026-10-20T09:00:00Z", "region": "LON1", "destination": {"objectStore": {"region": "FRA1", "objectStoreId": "os-1"}}},
			{"id": "b-2", "name": "manual-1", "state": "DATABASE_BACKUP_STATE_RUNNING", "databaseId": "12345", "createdAt": "2026-10-21T09:00:00Z", "region": "LON1", "destination": {"externalS3": {"endpoint": "https://s3.example.com", "region": "eu-1", "bucket": "bk", "accessKeyId": "AK", "secretAccessKey": "", "pathPrefix": "db/"}}}
		], "pageInfo": {"page": 1, "totalPages": 1}}`,
	})
	defer server.Close()

	got, err := client.ListDatabaseBackupsV3alpha1("12345")
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	expected := &PaginatedDatabaseBackup{Page: 1, Pages: 1, PerPage: 2, Items: []DatabaseBackup{
		{
			ID: "b-1", Name: "nightly", Software: "PostgreSQL", Status: "Ready", Schedule: "0 0 * * *",
			DatabaseName: "test-db", DatabaseID: "12345", IsScheduled: true,
			CreatedAt: time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC), Region: "LON1",
			Destination: &DatabaseBackupDestination{ObjectStore: &DatabaseBackupObjectStoreDestination{Region: "FRA1", ObjectStoreID: "os-1"}},
		},
		{
			ID: "b-2", Name: "manual-1", Status: "Running", DatabaseID: "12345",
			CreatedAt: time.Date(2026, 10, 21, 9, 0, 0, 0, time.UTC), Region: "LON1",
			Destination: &DatabaseBackupDestination{ExternalS3: &DatabaseBackupExternalS3Destination{Endpoint: "https://s3.example.com", Region: "eu-1", Bucket: "bk", AccessKeyID: "AK", PathPrefix: "db/"}},
		},
	}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestCreateDatabaseBackupV3alpha1WithDestination(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345/backups",
					RequestBody:  `{"name":"offsite","schedule":"0 0 * * *","type":"scheduled","region":"lon1","destination":{"object_store":{"region":"FRA1"}}}`,
					ResponseBody: `{"backup": {"id": "b-3", "name": "offsite", "state": "DATABASE_BACKUP_STATE_READY", "isScheduled": true, "destination": {"objectStore": {"region": "FRA1", "objectStoreId": "os-9"}}}}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.CreateDatabaseBackupV3alpha1("12345", &DatabaseBackupCreateRequest{
		Name:        "offsite",
		Schedule:    "0 0 * * *",
		Type:        "scheduled",
		Region:      "LON1",
		Destination: &DatabaseBackupDestination{ObjectStore: &DatabaseBackupObjectStoreDestination{Region: "FRA1"}},
	})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.ID != "b-3" || got.Status != "Ready" || got.Destination == nil || got.Destination.ObjectStore.ObjectStoreID != "os-9" {
		t.Errorf("Unexpected backup %+v", got)
	}
}

func TestDatabaseBackupStateToStatus(t *testing.T) {
	cases := map[string]string{
		"DATABASE_BACKUP_STATE_READY":       "Ready",
		"DATABASE_BACKUP_STATE_RUNNING":     "Running",
		"DATABASE_BACKUP_STATE_FAILED":      "Failed",
		"DATABASE_BACKUP_STATE_DELETING":    "Deleting",
		"DATABASE_BACKUP_STATE_UNSPECIFIED": "",
		"":                                  "",
	}
	for in, want := range cases {
		if got := databaseBackupStateToStatus(in); got != want {
			t.Errorf("%q: expected %q, got %q", in, want, got)
		}
	}
}

func TestUpdateDatabaseBackupV3alpha1SendsOnlyWhatIsSet(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "PATCH",
			Value: []ValueAdvanceClientForTesting{
				{
					URL: "/v3alpha1/databases/12345/backups",
					// No name and no schedule: an explicit "" would clear them on the server.
					RequestBody:  `{"destination":{"object_store":{"region":"FRA1"}},"region":"lon1"}`,
					ResponseBody: `{"backup": {"id": "scheduled", "name": "nightly", "schedule": "0 0 * * *", "state": "DATABASE_BACKUP_STATE_READY", "isScheduled": true}}`,
				},
			},
		},
	})
	defer server.Close()

	got, err := client.UpdateDatabaseBackupV3alpha1("12345", &UpdateDatabaseBackupV3alpha1Request{
		Region:      "LON1",
		Destination: &DatabaseBackupDestination{ObjectStore: &DatabaseBackupObjectStoreDestination{Region: "FRA1"}},
	})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.Name != "nightly" || got.Schedule != "0 0 * * *" {
		t.Errorf("Unexpected backup %+v", got)
	}
}

func TestV3alpha1RequestGuards(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{})
	defer server.Close()

	if _, err := client.ForkDatabase(&ForkDatabaseRequest{Name: "x", Size: "g3.db.small", Nodes: 1}); err == nil {
		t.Errorf("expected an error without a source database id")
	}
	if _, err := client.ForkDatabase(&ForkDatabaseRequest{Name: "x", Size: "g3.db.small", SourceDatabaseID: "1"}); err == nil {
		t.Errorf("expected an error without nodes")
	}
	if _, err := client.RestoreDatabaseToPointInTime("1", &RestoreDatabaseToPointInTimeRequest{Name: "x"}); err == nil {
		t.Errorf("expected an error without a target time")
	}
}

func TestDecodeGatewayError(t *testing.T) {
	cases := []struct {
		body string
		want error
	}{
		{`{"error":{"code":404,"status":"NOT_FOUND","message":"database not found","details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"DATABASE_NOT_FOUND","domain":"database.civo.com"}]}}`, DBaaSDatabaseNotFoundError},
		{`{"error":{"code":404,"status":"NOT_FOUND","message":"backup not found","details":[{"reason":"BACKUP_NOT_FOUND"}]}}`, DBaaSBackupNotFoundError},
		{`{"error":{"code":400,"status":"INVALID_ARGUMENT","message":"invalid request body: unknown field \"target_time\""}}`, CommonError},
		{`{"error":{"code":401,"status":"UNAUTHENTICATED","message":"invalid token"}}`, AuthenticationError},
		{`{"error":{"code":429,"status":"RESOURCE_EXHAUSTED","message":"quota exceeded"}}`, QuotaLimitReachedError},
	}
	for _, c := range cases {
		err := decodeError(HTTPError{Code: 400, Status: "400 Bad Request", Reason: c.body})
		if !errors.Is(err, c.want) {
			t.Errorf("%s: expected %v, got %v", c.body, c.want, err)
		}
	}
	err := decodeError(HTTPError{Code: 400, Status: "400 Bad Request", Reason: `{"error":{"code":400,"status":"INVALID_ARGUMENT","message":"storage can only grow"}}`})
	if err == nil || !strings.Contains(err.Error(), "storage can only grow") {
		t.Errorf("expected the gateway message to be kept, got %v", err)
	}
}

func TestV3alpha1RegionIsLowercased(t *testing.T) {
	seen := ""
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/12345": `{"database": {"id": "12345"}}`,
	})
	defer server.Close()
	client.Region = "LON1"
	// Wrap the test server handler to capture the query the SDK sent.
	orig := server.Config.Handler
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Query().Get("region")
		orig.ServeHTTP(w, r)
	})

	if _, err := client.GetDatabaseV3alpha1("12345"); err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if seen != "lon1" {
		t.Errorf("expected the region query to be lowercased for v3alpha1, got %q", seen)
	}
}

func TestUpgradeDatabase(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345/upgrade",
					RequestBody:  `{"validate_only":true,"region":"lon1"}`,
					ResponseBody: `{"fromVersion": "2.7.0", "toVersion": "2.9.0", "databases": [{"id": "12345", "name": "app"}, {"id": "67890", "name": "analytics"}], "started": false}`,
				},
			},
		},
	})
	defer server.Close()
	client.Region = "LON1"

	got, err := client.UpgradeDatabase("12345", &UpgradeDatabaseRequest{ValidateOnly: true})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	expected := &DatabaseUpgrade{FromVersion: "2.7.0", ToVersion: "2.9.0", Started: false, Databases: []DatabaseUpgradeTarget{{ID: "12345", Name: "app"}, {ID: "67890", Name: "analytics"}}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestListDatabaseVersionsForNetwork(t *testing.T) {
	seen := ""
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/versions": `{"versions": {"PostgreSQL": {"versions": [{"softwareVersion": "16", "default": false}, {"softwareVersion": "17", "default": true}]}}}`,
	})
	defer server.Close()
	orig := server.Config.Handler
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = r.URL.Query().Get("network_id")
		orig.ServeHTTP(w, r)
	})

	got, err := client.ListDatabaseVersionsForNetwork("net-1")
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if seen != "net-1" {
		t.Errorf("expected network_id=net-1, got %q", seen)
	}
	expected := map[string][]SupportedSoftwareVersion{"PostgreSQL": {{SoftwareVersion: "16"}, {SoftwareVersion: "17", Default: true}}}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestGetDatabaseV3alpha1EngineFields(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/12345": `{"database": {"id": "12345", "engineVersion": "2.7.0", "upgradeAvailableTo": "2.9.0"}}`,
	})
	defer server.Close()

	got, err := client.GetDatabaseV3alpha1("12345")
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if got.EngineVersion != "2.7.0" || got.UpgradeAvailableTo != "2.9.0" {
		t.Errorf("unexpected engine fields %+v", got)
	}
}

func TestUpgradeDatabaseVersion(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345/upgrade-version",
					RequestBody:  `{"target_version":"18","validate_only":true,"region":"lon1"}`,
					ResponseBody: `{"fromVersion": "17", "toVersion": "18", "checks": [{"name": "database_ready", "passed": true}, {"name": "storage_headroom", "passed": false, "message": "the data volume is 91% full; grow the storage first"}], "started": false}`,
				},
			},
		},
	})
	defer server.Close()
	client.Region = "LON1"

	got, err := client.UpgradeDatabaseVersion("12345", &UpgradeDatabaseVersionRequest{TargetVersion: "18", ValidateOnly: true})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	expected := &UpgradeDatabaseVersionResponse{
		FromVersion: "17",
		ToVersion:   "18",
		Checks: []DatabaseVersionUpgradeCheck{
			{Name: "database_ready", Passed: true},
			{Name: "storage_headroom", Passed: false, Message: "the data volume is 91% full; grow the storage first"},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Expected %+v, got %+v", expected, got)
	}
}

func TestUpgradeDatabaseVersionStarted(t *testing.T) {
	client, server, _ := NewAdvancedClientForTesting([]ConfigAdvanceClientForTesting{
		{
			Method: "POST",
			Value: []ValueAdvanceClientForTesting{
				{
					URL:          "/v3alpha1/databases/12345/upgrade-version",
					RequestBody:  `{"target_version":"18","region":"lon1"}`,
					ResponseBody: `{"fromVersion": "17", "toVersion": "18", "checks": [], "started": true, "backupName": "pre-upgrade-pg18"}`,
				},
			},
		},
	})
	defer server.Close()
	client.Region = "LON1"

	got, err := client.UpgradeDatabaseVersion("12345", &UpgradeDatabaseVersionRequest{TargetVersion: "18"})
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if !got.Started || got.BackupName != "pre-upgrade-pg18" || len(got.Checks) != 0 {
		t.Errorf("unexpected response %+v", got)
	}
}

func TestGetDatabaseV3alpha1VersionUpgradeFields(t *testing.T) {
	client, server, _ := NewClientForTesting(map[string]string{
		"/v3alpha1/databases/12345": `{"database": {"id": "12345", "majorUpgradesAvailable": ["18"], "versionUpgrade": {"phase": "Upgrading", "fromVersion": "16", "targetVersion": "17", "backupName": "pre-upgrade-pg17", "startedAt": "2026-10-23T15:00:00Z"}}}`,
	})
	defer server.Close()

	got, err := client.GetDatabaseV3alpha1("12345")
	if err != nil {
		t.Fatalf("Request returned an error: %s", err)
	}
	if !reflect.DeepEqual(got.MajorUpgradesAvailable, []string{"18"}) {
		t.Errorf("unexpected major upgrades %v", got.MajorUpgradesAvailable)
	}
	u := got.VersionUpgrade
	started := time.Date(2026, 10, 23, 15, 0, 0, 0, time.UTC)
	if u == nil || u.Phase != "Upgrading" || u.FromVersion != "16" || u.TargetVersion != "17" || u.BackupName != "pre-upgrade-pg17" ||
		u.StartedAt == nil || !u.StartedAt.Equal(started) || u.FinishedAt != nil {
		t.Errorf("unexpected version upgrade %+v", u)
	}
}
