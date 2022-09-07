package era

type Clusteravailability struct {
	Nxclusterid  *string `json:"nxClusterId,omitempty"`
	Datecreated  *string `json:"dateCreated,omitempty"`
	Datemodified *string `json:"dateModified,omitempty"`
	Ownerid      *string `json:"ownerId,omitempty"`
	Status       *string `json:"status,omitempty"`
	Profileid    *string `json:"profileId,omitempty"`
}

// ListProfile response
type ListProfileResponse struct {
	ID                  *string                `json:"id,omitempty"`
	Name                *string                `json:"name,omitempty"`
	Description         *string                `json:"description,omitempty"`
	Status              *string                `json:"status,omitempty"`
	Datecreated         *string                `json:"dateCreated,omitempty"`
	Datemodified        *string                `json:"dateModified,omitempty"`
	Owner               *string                `json:"owner,omitempty"`
	Enginetype          *string                `json:"engineType,omitempty"`
	Type                *string                `json:"type,omitempty"`
	Topology            *string                `json:"topology,omitempty"`
	Dbversion           *string                `json:"dbVersion,omitempty"`
	Systemprofile       bool                   `json:"systemProfile,omitempty"`
	Latestversion       *string                `json:"latestVersion,omitempty"`
	Latestversionid     *string                `json:"latestVersionId,omitempty"`
	Versions            []*Versions            `json:"versions,omitempty"`
	Assocdbservers      []interface{}          `json:"assocDbServers,omitempty"`
	Assocdatabases      []*string              `json:"assocDatabases,omitempty"`
	Nxclusterid         *string                `json:"nxClusterId,omitempty"`
	Clusteravailability []*Clusteravailability `json:"clusterAvailability,omitempty"`
}

type ProfileListResponse []ListProfileResponse

type Propertiesmap struct {
	DefaultContainer *string `json:"DEFAULT_CONTAINER"`
	MaxVdiskSize     *string `json:"MAX_VDISK_SIZE"`
}

type Properties struct {
	RefID       interface{} `json:"ref_id,omitempty"`
	Name        *string     `json:"name,omitempty"`
	Value       *string     `json:"value,omitempty"`
	Secure      bool        `json:"secure,omitempty"`
	Description interface{} `json:"description,omitempty"`
}

type VersionClusterAssociation struct {
	NxClusterID              *string       `json:"nxClusterId,omitempty"`
	DateCreated              *string       `json:"dateCreated,omitempty"`
	DateModified             *string       `json:"dateModified,omitempty"`
	OwnerID                  *string       `json:"ownerId,omitempty"`
	Status                   *string       `json:"status,omitempty"`
	ProfileVersionID         *string       `json:"profileVersionId,omitempty"`
	Properties               []*Properties `json:"properties,omitempty"`
	OptimizedForProvisioning bool          `json:"optimizedForProvisioning,omitempty"`
}

type Versions struct {
	ID                        *string                      `json:"id,omitempty"`
	Name                      *string                      `json:"name,omitempty"`
	Description               *string                      `json:"description,omitempty"`
	Status                    *string                      `json:"status,omitempty"`
	Datecreated               *string                      `json:"dateCreated,omitempty"`
	Datemodified              *string                      `json:"dateModified,omitempty"`
	Owner                     *string                      `json:"owner,omitempty"`
	Enginetype                *string                      `json:"engineType,omitempty"`
	Type                      *string                      `json:"type,omitempty"`
	Topology                  *string                      `json:"topology,omitempty"`
	Dbversion                 *string                      `json:"dbVersion,omitempty"`
	Version                   *string                      `json:"version,omitempty"`
	Profileid                 *string                      `json:"profileId,omitempty"`
	Published                 bool                         `json:"published,omitempty"`
	Deprecated                bool                         `json:"deprecated,omitempty"`
	Systemprofile             bool                         `json:"systemProfile,omitempty"`
	Propertiesmap             map[string]interface{}       `json:"propertiesMap,omitempty"`
	Properties                []*Properties                `json:"properties,omitempty"`
	VersionClusterAssociation []*VersionClusterAssociation `json:"versionClusterAssociation,omitempty"`
}

// ListClustersResponse structs
type ListClusterResponse struct {
	ID                   *string         `json:"id,omitempty"`
	Name                 *string         `json:"name,omitempty"`
	Uniquename           *string         `json:"uniqueName,omitempty"`
	Ipaddresses          []*string       `json:"ipAddresses,omitempty"`
	Fqdns                interface{}     `json:"fqdns,omitempty"`
	Nxclusteruuid        *string         `json:"nxClusterUUID,omitempty"`
	Description          *string         `json:"description,omitempty"`
	Cloudtype            *string         `json:"cloudType,omitempty"`
	Datecreated          *string         `json:"dateCreated,omitempty"`
	Datemodified         *string         `json:"dateModified,omitempty"`
	Ownerid              *string         `json:"ownerId,omitempty"`
	Status               *string         `json:"status,omitempty"`
	Version              *string         `json:"version,omitempty"`
	Hypervisortype       *string         `json:"hypervisorType,omitempty"`
	Hypervisorversion    *string         `json:"hypervisorVersion,omitempty"`
	Properties           []*Properties   `json:"properties,omitempty"`
	Referencecount       int             `json:"referenceCount,omitempty"`
	Username             interface{}     `json:"username,omitempty"`
	Password             interface{}     `json:"password,omitempty"`
	Cloudinfo            interface{}     `json:"cloudInfo,omitempty"`
	Resourceconfig       *Resourceconfig `json:"resourceConfig,omitempty"`
	Managementserverinfo interface{}     `json:"managementServerInfo,omitempty"`
	Entitycounts         interface{}     `json:"entityCounts,omitempty"`
	Healthy              bool            `json:"healthy,omitempty"`
}

type ClusterListResponse []ListClusterResponse

type Resourceconfig struct {
	Storagethresholdpercentage float64 `json:"storageThresholdPercentage,omitempty"`
	Memorythresholdpercentage  float64 `json:"memoryThresholdPercentage,omitempty"`
}

// ListSLAResponse structs
type ListSLAResponse struct {
	ID                     *string `json:"id,omitempty"`
	Name                   *string `json:"name,omitempty"`
	Uniquename             *string `json:"uniqueName,omitempty"`
	Description            *string `json:"description,omitempty"`
	Ownerid                *string `json:"ownerId,omitempty"`
	Datecreated            *string `json:"dateCreated,omitempty"`
	Datemodified           *string `json:"dateModified,omitempty"`
	CurrentActiveFrequency *string `json:"currentActiveFrequency,omitempty"`
	Continuousretention    int     `json:"continuousRetention,omitempty"`
	Dailyretention         int     `json:"dailyRetention,omitempty"`
	Weeklyretention        int     `json:"weeklyRetention,omitempty"`
	Monthlyretention       int     `json:"monthlyRetention,omitempty"`
	Quarterlyretention     int     `json:"quarterlyRetention,omitempty"`
	Yearlyretention        int     `json:"yearlyRetention,omitempty"`
	Referencecount         int     `json:"referenceCount,omitempty"`
	PitrEnabled            bool    `json:"pitrEnabled,omitempty"`
	Systemsla              bool    `json:"systemSla,omitempty"`
}

type SLAResponse []ListSLAResponse
