package Era

import "time"

// ListProfile response
type ListProfileResponse struct {
	ID                  string                `json:"id"`
	Name                string                `json:"name"`
	Description         string                `json:"description"`
	Status              string                `json:"status"`
	Datecreated         time.Time             `json:"dateCreated"`
	Datemodified        time.Time             `json:"dateModified"`
	Owner               string                `json:"owner"`
	Enginetype          string                `json:"engineType"`
	Type                string                `json:"type"`
	Topology            string                `json:"topology"`
	Dbversion           string                `json:"dbVersion"`
	Systemprofile       bool                  `json:"systemProfile"`
	Latestversion       string                `json:"latestVersion"`
	Latestversionid     string                `json:"latestVersionId"`
	Versions            []Versions            `json:"versions"`
	Assocdbservers      []interface{}         `json:"assocDbServers,omitempty"`
	Assocdatabases      []string              `json:"assocDatabases,omitempty"`
	Nxclusterid         string                `json:"nxClusterId,omitempty"`
	Clusteravailability []Clusteravailability `json:"clusterAvailability,omitempty"`
}

type ProfileResponse *ListProfileResponse

type ProfileListResponse []*ListProfileResponse

type Propertiesmap struct {
	DefaultContainer string `json:"DEFAULT_CONTAINER"`
	MaxVdiskSize     string `json:"MAX_VDISK_SIZE"`
}

type VersionClusterAssociation struct {
	NxClusterID              *string      `json:"nxClusterId,omitempty"`
	DateCreated              *string      `json:"dateCreated,omitempty"`
	DateModified             *string      `json:"dateModified,omitempty"`
	OwnerID                  *string      `json:"ownerId,omitempty"`
	Status                   *string      `json:"status,omitempty"`
	ProfileVersionID         *string      `json:"profileVersionId,omitempty"`
	Properties               []Properties `json:"properties,omitempty"`
	OptimizedForProvisioning bool         `json:"optimizedForProvisioning,omitempty"`
}

type Versions struct {
	ID                        string                      `json:"id"`
	Name                      string                      `json:"name"`
	Description               string                      `json:"description"`
	Status                    string                      `json:"status"`
	Datecreated               string                      `json:"dateCreated"`
	Datemodified              string                      `json:"dateModified"`
	Owner                     string                      `json:"owner"`
	Enginetype                string                      `json:"engineType"`
	Type                      string                      `json:"type"`
	Topology                  string                      `json:"topology"`
	Dbversion                 string                      `json:"dbVersion"`
	Systemprofile             bool                        `json:"systemProfile"`
	Version                   string                      `json:"version"`
	Profileid                 string                      `json:"profileId"`
	Published                 bool                        `json:"published"`
	Deprecated                bool                        `json:"deprecated"`
	Properties                []Properties                `json:"properties"`
	Propertiesmap             map[string]interface{}      `json:"propertiesMap"`
	VersionClusterAssociation []VersionClusterAssociation `json:"versionClusterAssociation"`
}

type Clusteravailability struct {
	Nxclusterid  string `json:"nxClusterId"`
	Datecreated  string `json:"dateCreated"`
	Datemodified string `json:"dateModified"`
	Ownerid      string `json:"ownerId"`
	Status       string `json:"status"`
	Profileid    string `json:"profileId"`
}

// ListClustersResponse structs
type ListClusterResponse []struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	Uniquename           string          `json:"uniqueName"`
	Ipaddresses          []string        `json:"ipAddresses"`
	Fqdns                interface{}     `json:"fqdns"`
	Nxclusteruuid        string          `json:"nxClusterUUID"`
	Description          string          `json:"description"`
	Cloudtype            string          `json:"cloudType"`
	Datecreated          string          `json:"dateCreated"`
	Datemodified         string          `json:"dateModified"`
	Ownerid              string          `json:"ownerId"`
	Status               string          `json:"status"`
	Version              string          `json:"version"`
	Hypervisortype       string          `json:"hypervisorType"`
	Hypervisorversion    string          `json:"hypervisorVersion"`
	Properties           []*Properties   `json:"properties"`
	Referencecount       int             `json:"referenceCount"`
	Username             interface{}     `json:"username"`
	Password             interface{}     `json:"password"`
	Cloudinfo            interface{}     `json:"cloudInfo"`
	Resourceconfig       *Resourceconfig `json:"resourceConfig"`
	Managementserverinfo interface{}     `json:"managementServerInfo"`
	Entitycounts         interface{}     `json:"entityCounts"`
	Healthy              bool            `json:"healthy"`
}

type Properties struct {
	RefID       interface{} `json:"ref_id"`
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Secure      bool        `json:"secure"`
	Description interface{} `json:"description"`
}

type Resourceconfig struct {
	Storagethresholdpercentage float64 `json:"storageThresholdPercentage"`
	Memorythresholdpercentage  float64 `json:"memoryThresholdPercentage"`
}

// ListSLAResponse structs
type ListSLAResponse struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Uniquename             string `json:"uniqueName"`
	Description            string `json:"description"`
	Ownerid                string `json:"ownerId"`
	Systemsla              bool   `json:"systemSla"`
	Datecreated            string `json:"dateCreated"`
	Datemodified           string `json:"dateModified"`
	Continuousretention    int    `json:"continuousRetention"`
	Dailyretention         int    `json:"dailyRetention"`
	Weeklyretention        int    `json:"weeklyRetention"`
	Monthlyretention       int    `json:"monthlyRetention"`
	Quarterlyretention     int    `json:"quarterlyRetention"`
	Yearlyretention        int    `json:"yearlyRetention"`
	Referencecount         int    `json:"referenceCount"`
	PitrEnabled            bool   `json:"pitrEnabled,omitempty"`
	CurrentActiveFrequency string `json:"currentActiveFrequency,omitempty"`
}

type SLAResponse []*ListSLAResponse
