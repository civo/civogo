package civogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// This file covers the /v3alpha1 database API. The v3alpha1 gateway is a
// protojson surface: responses are camelCase and wrapped in an envelope
// ({"database": {...}}, {"backups": [...]}) and request bodies accept the
// proto field names (snake_case). Everything is decoded into the same public
// types the /v2 methods return, so callers only pick the method, not the
// shape.

// DatabaseReadEndpoint is the load-balanced read-only endpoint served by the
// standby replicas of a database with three or more nodes.
type DatabaseReadEndpoint struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DNSEntry string `json:"dns_entry,omitempty"`
}

// DatabaseExtension is an entry of the extension catalogue a region offers.
type DatabaseExtension struct {
	Name            string   `json:"name"`
	Description     string   `json:"description,omitempty"`
	Versions        []string `json:"versions,omitempty"`
	RequiresRestart bool     `json:"requires_restart"`
}

// DatabaseBackupDestinationRegion is a region that can hold database backups.
type DatabaseBackupDestinationRegion struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	ObjectStore bool   `json:"object_store"`
	IsCurrent   bool   `json:"is_current"`
}

// DatabaseBackupObjectStoreDestination stores backups in a Civo Object Store,
// in this region (leave Region empty) or in another one.
type DatabaseBackupObjectStoreDestination struct {
	Region        string `json:"region,omitempty"`
	ObjectStoreID string `json:"object_store_id,omitempty"`
}

// DatabaseBackupExternalS3Destination stores backups in any S3-compatible
// bucket outside Civo.
type DatabaseBackupExternalS3Destination struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region,omitempty"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	PathPrefix      string `json:"path_prefix,omitempty"`
}

// DatabaseBackupDestination selects where a backup is written. Exactly one of
// the two fields is set.
type DatabaseBackupDestination struct {
	ObjectStore *DatabaseBackupObjectStoreDestination `json:"object_store,omitempty"`
	ExternalS3  *DatabaseBackupExternalS3Destination  `json:"external_s3,omitempty"`
}

// ForkDatabaseRequest creates a new, independent database from an existing
// one. The source is the database at its latest backup by default; set
// SourceBackup to fork a specific backup, or SourceTargetTime to fork the
// database as it was at that moment (point-in-time, within the restore
// window reported by Database.EarliestRestoreTime).
type ForkDatabaseRequest struct {
	Name             string     `json:"name" validate:"required"`
	Size             string     `json:"size" validate:"required"`
	Nodes            int        `json:"nodes" validate:"required"`
	SourceDatabaseID string     `json:"source_database_id" validate:"required"`
	SourceBackup     string     `json:"source_backup,omitempty"`
	SourceTargetTime *time.Time `json:"source_target_time,omitempty"`
	Software         string     `json:"software,omitempty"`
	SoftwareVersion  string     `json:"software_version,omitempty"`
	NetworkID        string     `json:"network_id,omitempty"`
	FirewallID       string     `json:"firewall_id,omitempty"`
	FirewallRules    string     `json:"firewall_rules,omitempty"`
	StorageGB        *int       `json:"storage_gb,omitempty"`
	Extensions       []string   `json:"extensions,omitempty"`
	Region           string     `json:"region,omitempty"`
}

// RestoreDatabaseToPointInTimeRequest restores a database in place to the
// state it had at TargetTime. Data written after that moment is lost.
type RestoreDatabaseToPointInTimeRequest struct {
	Name       string    `json:"name" validate:"required"`
	TargetTime time.Time `json:"target_time" validate:"required"`
	Region     string    `json:"region,omitempty"`
}

// UpdateDatabaseBackupV3alpha1Request is the partial update of the scheduled
// backup accepted by the v3alpha1 API. Unset pointers leave the field
// untouched; an explicit empty string would clear the name or schedule.
type UpdateDatabaseBackupV3alpha1Request struct {
	Name        *string                    `json:"name,omitempty"`
	Schedule    *string                    `json:"schedule,omitempty"`
	Destination *DatabaseBackupDestination `json:"destination,omitempty"`
	Region      string                     `json:"region,omitempty"`
}

// UpdateDatabaseV3alpha1Request is the partial update accepted by the
// v3alpha1 API. Unset pointers leave the field untouched. Extensions and
// DatabaseParameters are only applied when the matching Update* flag is set,
// which is how an empty list ("remove everything") is told apart from "not
// sent".
type UpdateDatabaseV3alpha1Request struct {
	Name                     *string           `json:"name,omitempty"`
	Nodes                    *int              `json:"nodes,omitempty"`
	FirewallID               *string           `json:"firewall_id,omitempty"`
	StorageGB                *int              `json:"storage_gb,omitempty"`
	Extensions               []string          `json:"extensions,omitempty"`
	UpdateExtensions         *bool             `json:"update_extensions,omitempty"`
	DatabaseParameters       map[string]string `json:"database_parameters,omitempty"`
	UpdateDatabaseParameters *bool             `json:"update_database_parameters,omitempty"`
	Region                   string            `json:"region,omitempty"`
}

// v3alpha1Database is the protojson rendering of a database.
type v3alpha1Database struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Nodes            int                `json:"nodes"`
	Size             string             `json:"size"`
	Software         string             `json:"software"`
	SoftwareVersion  string             `json:"softwareVersion"`
	PublicIPv4       string             `json:"publicIpv4"`
	PrivateIPv4      string             `json:"privateIpv4"`
	Port             int                `json:"port"`
	Username         string             `json:"username"`
	Password         string             `json:"password"`
	DatabaseUserInfo []DatabaseUserInfo `json:"databaseUserInfo"`
	NetworkID        string             `json:"networkId"`
	FirewallID       string             `json:"firewallId"`
	BackupConfigured bool               `json:"backupConfigured"`
	Status           string             `json:"status"`
	CreatedAt        *time.Time         `json:"createdAt,omitempty"`
	DNSEntry         string             `json:"dnsEntry"`
	Region           string             `json:"region"`
	ReadEndpoint     *struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		DNSEntry string `json:"dnsEntry"`
	} `json:"readEndpoint,omitempty"`
	StorageGB              int               `json:"storageGb"`
	BaseStorageGB          int               `json:"baseStorageGb"`
	Extensions             []string          `json:"extensions"`
	DatabaseParameters     map[string]string `json:"databaseParameters"`
	EarliestRestoreTime    *time.Time        `json:"earliestRestoreTime,omitempty"`
	EngineVersion          string            `json:"engineVersion"`
	UpgradeAvailableTo     string            `json:"upgradeAvailableTo"`
	MajorUpgradesAvailable []string          `json:"majorUpgradesAvailable"`
	VersionUpgrade         *struct {
		Phase         string     `json:"phase"`
		FromVersion   string     `json:"fromVersion"`
		TargetVersion string     `json:"targetVersion"`
		BackupName    string     `json:"backupName"`
		Message       string     `json:"message"`
		StartedAt     *time.Time `json:"startedAt,omitempty"`
		FinishedAt    *time.Time `json:"finishedAt,omitempty"`
	} `json:"versionUpgrade,omitempty"`
}

func (d *v3alpha1Database) toDatabase() Database {
	db := Database{
		ID:                     d.ID,
		Name:                   d.Name,
		Nodes:                  d.Nodes,
		Size:                   d.Size,
		Software:               d.Software,
		SoftwareVersion:        d.SoftwareVersion,
		PublicIPv4:             d.PublicIPv4,
		PrivateIPv4:            d.PrivateIPv4,
		NetworkID:              d.NetworkID,
		FirewallID:             d.FirewallID,
		Port:                   d.Port,
		Username:               d.Username,
		Password:               d.Password,
		DatabaseUserInfo:       d.DatabaseUserInfo,
		DNSEntry:               d.DNSEntry,
		Status:                 d.Status,
		Region:                 d.Region,
		BackupConfigured:       d.BackupConfigured,
		CreatedAt:              d.CreatedAt,
		StorageGB:              d.StorageGB,
		BaseStorageGB:          d.BaseStorageGB,
		Extensions:             d.Extensions,
		DatabaseParameters:     d.DatabaseParameters,
		EarliestRestoreTime:    d.EarliestRestoreTime,
		EngineVersion:          d.EngineVersion,
		UpgradeAvailableTo:     d.UpgradeAvailableTo,
		MajorUpgradesAvailable: d.MajorUpgradesAvailable,
	}
	if u := d.VersionUpgrade; u != nil {
		db.VersionUpgrade = &DatabaseVersionUpgrade{
			Phase:         u.Phase,
			FromVersion:   u.FromVersion,
			TargetVersion: u.TargetVersion,
			BackupName:    u.BackupName,
			Message:       u.Message,
			StartedAt:     u.StartedAt,
			FinishedAt:    u.FinishedAt,
		}
	}
	if d.ReadEndpoint != nil {
		db.ReadEndpoint = &DatabaseReadEndpoint{Host: d.ReadEndpoint.Host, Port: d.ReadEndpoint.Port, DNSEntry: d.ReadEndpoint.DNSEntry}
	}
	return db
}

// v3alpha1DatabaseBackup is the protojson rendering of a backup.
type v3alpha1DatabaseBackup struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Software     string     `json:"software"`
	State        string     `json:"state"`
	Schedule     string     `json:"schedule"`
	DatabaseName string     `json:"databaseName"`
	DatabaseID   string     `json:"databaseId"`
	IsScheduled  bool       `json:"isScheduled"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	Region       string     `json:"region"`
	Destination  *struct {
		ObjectStore *struct {
			Region        string `json:"region"`
			ObjectStoreID string `json:"objectStoreId"`
		} `json:"objectStore,omitempty"`
		ExternalS3 *struct {
			Endpoint        string `json:"endpoint"`
			Region          string `json:"region"`
			Bucket          string `json:"bucket"`
			AccessKeyID     string `json:"accessKeyId"`
			SecretAccessKey string `json:"secretAccessKey"`
			PathPrefix      string `json:"pathPrefix"`
		} `json:"externalS3,omitempty"`
	} `json:"destination,omitempty"`
}

func (b *v3alpha1DatabaseBackup) toDatabaseBackup() DatabaseBackup {
	bk := DatabaseBackup{
		ID:           b.ID,
		Name:         b.Name,
		Software:     b.Software,
		Status:       databaseBackupStateToStatus(b.State),
		Schedule:     b.Schedule,
		DatabaseName: b.DatabaseName,
		DatabaseID:   b.DatabaseID,
		IsScheduled:  b.IsScheduled,
		Region:       b.Region,
	}
	if b.CreatedAt != nil {
		bk.CreatedAt = *b.CreatedAt
	}
	if b.Destination != nil {
		bk.Destination = &DatabaseBackupDestination{}
		if b.Destination.ObjectStore != nil {
			bk.Destination.ObjectStore = &DatabaseBackupObjectStoreDestination{
				Region:        b.Destination.ObjectStore.Region,
				ObjectStoreID: b.Destination.ObjectStore.ObjectStoreID,
			}
		}
		if b.Destination.ExternalS3 != nil {
			bk.Destination.ExternalS3 = &DatabaseBackupExternalS3Destination{
				Endpoint:        b.Destination.ExternalS3.Endpoint,
				Region:          b.Destination.ExternalS3.Region,
				Bucket:          b.Destination.ExternalS3.Bucket,
				AccessKeyID:     b.Destination.ExternalS3.AccessKeyID,
				SecretAccessKey: b.Destination.ExternalS3.SecretAccessKey,
				PathPrefix:      b.Destination.ExternalS3.PathPrefix,
			}
		}
	}
	return bk
}

// databaseBackupStateToStatus turns the proto enum name
// (DATABASE_BACKUP_STATE_READY) into the status word the /v2 API uses (Ready).
func databaseBackupStateToStatus(state string) string {
	s := strings.TrimPrefix(state, "DATABASE_BACKUP_STATE_")
	if s == "" || s == "UNSPECIFIED" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

type v3alpha1PageInfo struct {
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
	PerPage    int `json:"perPage"`
}

// ListDatabasesV3alpha1 lists the databases of the client's region through
// the v3alpha1 API, which also reports the read-only endpoint, storage,
// extensions, parameters and restore window of each database.
func (c *Client) ListDatabasesV3alpha1() (*PaginatedDatabases, error) {
	items, err := paginateAll("databases", func(page, perPage int) ([]Database, int, error) {
		resp, err := c.SendGetRequest(fmt.Sprintf("/v3alpha1/databases?page=%d&per_page=%d", page, perPage))
		if err != nil {
			return nil, 0, decodeError(err)
		}
		pg := struct {
			Databases []v3alpha1Database `json:"databases"`
			PageInfo  *v3alpha1PageInfo  `json:"pageInfo,omitempty"`
		}{}
		if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&pg); err != nil {
			return nil, 0, err
		}
		dbs := make([]Database, 0, len(pg.Databases))
		for i := range pg.Databases {
			dbs = append(dbs, pg.Databases[i].toDatabase())
		}
		totalPages := 1
		if pg.PageInfo != nil && pg.PageInfo.TotalPages > 0 {
			totalPages = pg.PageInfo.TotalPages
		}
		return dbs, totalPages, nil
	})
	if err != nil {
		return nil, err
	}
	return &PaginatedDatabases{Page: 1, Pages: 1, PerPage: len(items), Items: items}, nil
}

// GetDatabaseV3alpha1 fetches a database through the v3alpha1 API, including
// the read-only endpoint, storage, extensions, parameters and the earliest
// point in time it can be restored or forked to.
func (c *Client) GetDatabaseV3alpha1(id string) (*Database, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v3alpha1/databases/%s", id))
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeV3alpha1Database(resp)
}

// v3alpha1Region returns the region code the v3alpha1 gateway expects: the
// request's own when set, else the client's, always lowercase since the
// gateway matches region codes case-sensitively.
func (c *Client) v3alpha1Region(region string) string {
	if region == "" {
		region = c.Region
	}
	return strings.ToLower(region)
}

// ForkDatabase creates a new database from the backups of an existing one.
// Software and SoftwareVersion should be set to the source's: the API does
// not derive them from the source and falls back to its default engine.
func (c *Client) ForkDatabase(v *ForkDatabaseRequest) (*Database, error) {
	if v == nil || v.SourceDatabaseID == "" {
		return nil, fmt.Errorf("a source database id is required to fork a database")
	}
	if v.Name == "" || v.Size == "" || v.Nodes < 1 {
		return nil, fmt.Errorf("a name, a size and at least one node are required to fork a database")
	}
	v.Region = c.v3alpha1Region(v.Region)
	body, err := c.SendPostRequest("/v3alpha1/databases", v)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeV3alpha1Database(body)
}

// RestoreDatabaseToPointInTime restores a database in place to the given
// moment. The moment must not be earlier than Database.EarliestRestoreTime.
func (c *Client) RestoreDatabaseToPointInTime(id string, v *RestoreDatabaseToPointInTimeRequest) (*Database, error) {
	if v == nil || v.TargetTime.IsZero() {
		return nil, fmt.Errorf("a target time is required to restore a database to a point in time")
	}
	if v.Name == "" {
		return nil, fmt.Errorf("a name is required to restore a database")
	}
	v.Region = c.v3alpha1Region(v.Region)
	body, err := c.SendPostRequest(fmt.Sprintf("/v3alpha1/databases/%s/restore", id), v)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeV3alpha1Database(body)
}

// UpdateDatabaseV3alpha1 applies a partial update through the v3alpha1 API.
func (c *Client) UpdateDatabaseV3alpha1(id string, v *UpdateDatabaseV3alpha1Request) (*Database, error) {
	if v == nil {
		return nil, fmt.Errorf("an update request is required")
	}
	v.Region = c.v3alpha1Region(v.Region)
	body, err := c.SendPatchRequest(fmt.Sprintf("/v3alpha1/databases/%s", id), v)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeV3alpha1Database(body)
}

// ResizeDatabaseStorage grows the storage of every node of a database to
// storageGB. Storage can only be increased; the GB beyond the plan's base
// storage are billed as an add-on.
func (c *Client) ResizeDatabaseStorage(id string, storageGB int) (*Database, error) {
	return c.UpdateDatabaseV3alpha1(id, &UpdateDatabaseV3alpha1Request{StorageGB: &storageGB, Region: c.Region})
}

// UpdateDatabaseExtensions replaces the set of enabled extensions. Pass an
// empty slice to disable them all. Extensions that need it trigger a brief
// rolling restart.
func (c *Client) UpdateDatabaseExtensions(id string, extensions []string) (*Database, error) {
	if extensions == nil {
		extensions = []string{}
	}
	update := true
	return c.UpdateDatabaseV3alpha1(id, &UpdateDatabaseV3alpha1Request{
		Extensions:       extensions,
		UpdateExtensions: &update,
		Region:           c.Region,
	})
}

// UpdateDatabaseParameters replaces the custom database parameters. Pass an
// empty map to reset them all to their defaults.
func (c *Client) UpdateDatabaseParameters(id string, parameters map[string]string) (*Database, error) {
	if parameters == nil {
		parameters = map[string]string{}
	}
	update := true
	return c.UpdateDatabaseV3alpha1(id, &UpdateDatabaseV3alpha1Request{
		DatabaseParameters:       parameters,
		UpdateDatabaseParameters: &update,
		Region:                   c.Region,
	})
}

// ListDatabaseExtensions returns the extension catalogue of the client's
// region.
func (c *Client) ListDatabaseExtensions() ([]DatabaseExtension, error) {
	resp, err := c.SendGetRequest("/v3alpha1/databases/extensions")
	if err != nil {
		return nil, decodeError(err)
	}
	out := struct {
		Extensions []struct {
			Name            string   `json:"name"`
			Description     string   `json:"description"`
			Versions        []string `json:"versions"`
			RequiresRestart bool     `json:"requiresRestart"`
		} `json:"extensions"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&out); err != nil {
		return nil, err
	}
	exts := make([]DatabaseExtension, 0, len(out.Extensions))
	for _, e := range out.Extensions {
		exts = append(exts, DatabaseExtension{Name: e.Name, Description: e.Description, Versions: e.Versions, RequiresRestart: e.RequiresRestart})
	}
	return exts, nil
}

// ListDatabaseBackupDestinationRegions returns the regions whose Civo Object
// Store can receive backups from the client's region.
func (c *Client) ListDatabaseBackupDestinationRegions() ([]DatabaseBackupDestinationRegion, error) {
	resp, err := c.SendGetRequest("/v3alpha1/databases/backup-destination-regions")
	if err != nil {
		return nil, decodeError(err)
	}
	out := struct {
		Regions []struct {
			Code        string `json:"code"`
			Name        string `json:"name"`
			ObjectStore bool   `json:"objectStore"`
			IsCurrent   bool   `json:"isCurrent"`
		} `json:"regions"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&out); err != nil {
		return nil, err
	}
	regions := make([]DatabaseBackupDestinationRegion, 0, len(out.Regions))
	for _, r := range out.Regions {
		regions = append(regions, DatabaseBackupDestinationRegion{Code: r.Code, Name: r.Name, ObjectStore: r.ObjectStore, IsCurrent: r.IsCurrent})
	}
	return regions, nil
}

// ListDatabaseBackupsV3alpha1 lists the backups of a database through the
// v3alpha1 API, which also reports each backup's destination.
func (c *Client) ListDatabaseBackupsV3alpha1(did string) (*PaginatedDatabaseBackup, error) {
	items, err := paginateAll("databases/"+did+"/backups", func(page, perPage int) ([]DatabaseBackup, int, error) {
		resp, err := c.SendGetRequest(fmt.Sprintf("/v3alpha1/databases/%s/backups?page=%d&per_page=%d", did, page, perPage))
		if err != nil {
			return nil, 0, decodeError(err)
		}
		pg := struct {
			Backups  []v3alpha1DatabaseBackup `json:"backups"`
			PageInfo *v3alpha1PageInfo        `json:"pageInfo,omitempty"`
		}{}
		if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&pg); err != nil {
			return nil, 0, err
		}
		backups := make([]DatabaseBackup, 0, len(pg.Backups))
		for i := range pg.Backups {
			backups = append(backups, pg.Backups[i].toDatabaseBackup())
		}
		totalPages := 1
		if pg.PageInfo != nil && pg.PageInfo.TotalPages > 0 {
			totalPages = pg.PageInfo.TotalPages
		}
		return backups, totalPages, nil
	})
	if err != nil {
		return nil, err
	}
	return &PaginatedDatabaseBackup{Page: 1, Pages: 1, PerPage: len(items), Items: items}, nil
}

// GetDatabaseBackupV3alpha1 fetches one backup through the v3alpha1 API. id is
// the manual backup UUID or the literal "scheduled".
func (c *Client) GetDatabaseBackupV3alpha1(did, id string) (*DatabaseBackup, error) {
	resp, err := c.SendGetRequest(fmt.Sprintf("/v3alpha1/databases/%s/backups/%s", did, id))
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeV3alpha1DatabaseBackup(resp)
}

// CreateDatabaseBackupV3alpha1 creates a manual or scheduled backup through
// the v3alpha1 API, optionally in a Civo Object Store of another region or an
// external S3 bucket (see DatabaseBackupCreateRequest.Destination).
func (c *Client) CreateDatabaseBackupV3alpha1(did string, v *DatabaseBackupCreateRequest) (*DatabaseBackup, error) {
	if v == nil {
		return nil, fmt.Errorf("a backup request is required")
	}
	v.Region = c.v3alpha1Region(v.Region)
	body, err := c.SendPostRequest(fmt.Sprintf("/v3alpha1/databases/%s/backups", did), v)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeV3alpha1DatabaseBackup(body)
}

// UpdateDatabaseBackupV3alpha1 updates the scheduled backup of a database
// through the v3alpha1 API, including its destination. Only the fields set
// in the request are changed.
func (c *Client) UpdateDatabaseBackupV3alpha1(did string, v *UpdateDatabaseBackupV3alpha1Request) (*DatabaseBackup, error) {
	if v == nil {
		return nil, fmt.Errorf("a backup update request is required")
	}
	v.Region = c.v3alpha1Region(v.Region)
	body, err := c.SendPatchRequest(fmt.Sprintf("/v3alpha1/databases/%s/backups", did), v)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeV3alpha1DatabaseBackup(body)
}

// DeleteDatabaseBackupV3alpha1 deletes a backup through the v3alpha1 API. id
// is the manual backup UUID or the literal "scheduled". The API answers with
// an empty body, so there is no response to return.
func (c *Client) DeleteDatabaseBackupV3alpha1(did, id string) error {
	_, err := c.SendDeleteRequest(fmt.Sprintf("/v3alpha1/databases/%s/backups/%s", did, id))
	if err != nil {
		return decodeError(err)
	}
	return nil
}

func decodeV3alpha1Database(body []byte) (*Database, error) {
	out := struct {
		Database *v3alpha1Database `json:"database"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&out); err != nil {
		return nil, err
	}
	if out.Database == nil {
		return nil, fmt.Errorf("v3alpha1 response carried no database")
	}
	db := out.Database.toDatabase()
	return &db, nil
}

func decodeV3alpha1DatabaseBackup(body []byte) (*DatabaseBackup, error) {
	out := struct {
		Backup *v3alpha1DatabaseBackup `json:"backup"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&out); err != nil {
		return nil, err
	}
	if out.Backup == nil {
		return nil, fmt.Errorf("v3alpha1 response carried no backup")
	}
	bk := out.Backup.toDatabaseBackup()
	return &bk, nil
}

// UpgradeDatabaseRequest upgrades a database to a newer engine release.
type UpgradeDatabaseRequest struct {
	// TargetVersion is the engine release to upgrade to; empty means the
	// newest available.
	TargetVersion string `json:"target_version,omitempty"`
	// ValidateOnly only reports what would be upgraded.
	ValidateOnly bool   `json:"validate_only,omitempty"`
	Region       string `json:"region,omitempty"`
}

// DatabaseUpgrade describes an engine upgrade.
type DatabaseUpgrade struct {
	FromVersion string `json:"from_version"`
	ToVersion   string `json:"to_version"`
	// Databases are every database upgraded: the requested one and the others
	// on its network, which share the engine.
	Databases []DatabaseUpgradeTarget `json:"databases"`
	// Started is false for validate-only requests.
	Started bool `json:"started"`
}

// DatabaseUpgradeTarget is a database included in an upgrade.
type DatabaseUpgradeTarget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpgradeDatabase upgrades a database, and the other databases on its network,
// to a newer engine release. Each database restarts while it moves; until
// this is called nothing restarts. Use ValidateOnly to list the databases that
// would be upgraded first.
func (c *Client) UpgradeDatabase(id string, v *UpgradeDatabaseRequest) (*DatabaseUpgrade, error) {
	if v == nil {
		v = &UpgradeDatabaseRequest{}
	}
	v.Region = c.v3alpha1Region(v.Region)
	body, err := c.SendPostRequest(fmt.Sprintf("/v3alpha1/databases/%s/upgrade", id), v)
	if err != nil {
		return nil, decodeError(err)
	}
	out := struct {
		FromVersion string `json:"fromVersion"`
		ToVersion   string `json:"toVersion"`
		Databases   []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"databases"`
		Started bool `json:"started"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&out); err != nil {
		return nil, err
	}
	result := &DatabaseUpgrade{FromVersion: out.FromVersion, ToVersion: out.ToVersion, Started: out.Started, Databases: []DatabaseUpgradeTarget{}}
	for _, d := range out.Databases {
		result.Databases = append(result.Databases, DatabaseUpgradeTarget{ID: d.ID, Name: d.Name})
	}
	return result, nil
}

// UpgradeDatabaseVersionRequest upgrades a PostgreSQL database in place to a
// newer major version.
type UpgradeDatabaseVersionRequest struct {
	// TargetVersion is the PostgreSQL major to upgrade to, e.g. "18".
	TargetVersion string `json:"target_version"`
	// ValidateOnly only runs the pre-flight checks.
	ValidateOnly bool   `json:"validate_only,omitempty"`
	Region       string `json:"region,omitempty"`
}

// DatabaseVersionUpgradeCheck is one pre-flight check of a major version
// upgrade.
type DatabaseVersionUpgradeCheck struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// UpgradeDatabaseVersionResponse is the result of a major version upgrade
// request.
type UpgradeDatabaseVersionResponse struct {
	FromVersion string `json:"from_version"`
	ToVersion   string `json:"to_version"`
	// Checks are every pre-flight check, passed or not.
	Checks []DatabaseVersionUpgradeCheck `json:"checks"`
	// Started is false for validate-only requests.
	Started bool `json:"started"`
	// BackupName is the backup taken before the upgrade.
	BackupName string `json:"backup_name,omitempty"`
}

// UpgradeDatabaseVersion upgrades a PostgreSQL database in place to a newer
// major version. A backup is taken first and kept; the database then stops
// for a few minutes while its data is upgraded, and its replicas re-sync
// afterwards. There is no downgrade. Use ValidateOnly to run the pre-flight
// checks without changing anything; when a check fails on a real request the
// API returns an error listing the failed checks.
func (c *Client) UpgradeDatabaseVersion(id string, v *UpgradeDatabaseVersionRequest) (*UpgradeDatabaseVersionResponse, error) {
	if v == nil {
		v = &UpgradeDatabaseVersionRequest{}
	}
	v.Region = c.v3alpha1Region(v.Region)
	body, err := c.SendPostRequest(fmt.Sprintf("/v3alpha1/databases/%s/upgrade-version", id), v)
	if err != nil {
		return nil, decodeError(err)
	}
	out := struct {
		FromVersion string `json:"fromVersion"`
		ToVersion   string `json:"toVersion"`
		Checks      []struct {
			Name    string `json:"name"`
			Passed  bool   `json:"passed"`
			Message string `json:"message"`
		} `json:"checks"`
		Started    bool   `json:"started"`
		BackupName string `json:"backupName"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&out); err != nil {
		return nil, err
	}
	result := &UpgradeDatabaseVersionResponse{
		FromVersion: out.FromVersion,
		ToVersion:   out.ToVersion,
		Started:     out.Started,
		BackupName:  out.BackupName,
		Checks:      []DatabaseVersionUpgradeCheck{},
	}
	for _, ch := range out.Checks {
		result.Checks = append(result.Checks, DatabaseVersionUpgradeCheck{Name: ch.Name, Passed: ch.Passed, Message: ch.Message})
	}
	return result, nil
}

// ListDatabaseVersionsForNetwork returns the database versions a new database
// on a network can use in the client's region. Databases on a network share an
// engine release, which may not run every version until it is upgraded. An
// empty networkID returns the region's offering.
func (c *Client) ListDatabaseVersionsForNetwork(networkID string) (map[string][]SupportedSoftwareVersion, error) {
	path := "/v3alpha1/databases/versions"
	if networkID != "" {
		path += "?network_id=" + url.QueryEscape(networkID)
	}
	resp, err := c.SendGetRequest(path)
	if err != nil {
		return nil, decodeError(err)
	}
	out := struct {
		Versions map[string]struct {
			Versions []struct {
				SoftwareVersion string `json:"softwareVersion"`
				Default         bool   `json:"default"`
			} `json:"versions"`
		} `json:"versions"`
	}{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&out); err != nil {
		return nil, err
	}
	versions := make(map[string][]SupportedSoftwareVersion, len(out.Versions))
	for software, list := range out.Versions {
		for _, v := range list.Versions {
			versions[software] = append(versions[software], SupportedSoftwareVersion{SoftwareVersion: v.SoftwareVersion, Default: v.Default})
		}
	}
	return versions, nil
}
