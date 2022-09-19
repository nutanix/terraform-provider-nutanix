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
	RefID       *string `json:"ref_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Value       *string `json:"value,omitempty"`
	Secure      bool    `json:"secure,omitempty"`
	Description *string `json:"description,omitempty"`
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

type ListDatabaseTypesResponse map[string]DatabaseTypeProperties

type DatabaseTypeProperties struct {
	Databasetype                  string `json:"databaseType"`
	Stagingdriveautotunesupported bool   `json:"stagingDriveAutoTuneSupported"`
	Defaultstagingdriveautotune   bool   `json:"defaultStagingDriveAutoTune"`
	Logdriveautotunesupported     bool   `json:"logDriveAutoTuneSupported"`
	Defaultlogdriveautotune       bool   `json:"defaultLogDriveAutoTune"`
}

// ProvisionDatabaseRequestStructs
type ProvisionDatabaseRequest struct {
	Databasetype             string            `json:"databaseType,omitempty"`
	Name                     string            `json:"name,omitempty"`
	Databasedescription      string            `json:"databaseDescription,omitempty"`
	DatabaseServerID         string            `json:"dbserverId,omitempty"`
	Softwareprofileid        string            `json:"softwareProfileId,omitempty"`
	Softwareprofileversionid string            `json:"softwareProfileVersionId,omitempty"`
	Computeprofileid         string            `json:"computeProfileId,omitempty"`
	Networkprofileid         string            `json:"networkProfileId,omitempty"`
	Dbparameterprofileid     string            `json:"dbParameterProfileId,omitempty"`
	Newdbservertimezone      string            `json:"newDbServerTimeZone,omitempty"`
	Timemachineinfo          Timemachineinfo   `json:"timeMachineInfo,omitempty"`
	Actionarguments          []Actionarguments `json:"actionArguments,omitempty"`
	Createdbserver           bool              `json:"createDbserver,omitempty"`
	Nodecount                int               `json:"nodeCount,omitempty"`
	Nxclusterid              string            `json:"nxClusterId,omitempty"`
	Sshpublickey             string            `json:"sshPublicKey,omitempty"`
	Clustered                bool              `json:"clustered,omitempty"`
	Nodes                    []Nodes           `json:"nodes,omitempty"`
	Autotunestagingdrive     bool              `json:"autoTuneStagingDrive,omitempty"`
}

type Snapshottimeofday struct {
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
	Seconds int `json:"seconds"`
}

type Continuousschedule struct {
	Enabled           bool `json:"enabled"`
	Logbackupinterval int  `json:"logBackupInterval"`
	Snapshotsperday   int  `json:"snapshotsPerDay"`
}

type Weeklyschedule struct {
	Enabled   bool   `json:"enabled"`
	Dayofweek string `json:"dayOfWeek"`
}

type Monthlyschedule struct {
	Enabled    bool `json:"enabled"`
	Dayofmonth int  `json:"dayOfMonth"`
}

type Quartelyschedule struct {
	Enabled    bool   `json:"enabled"`
	Startmonth string `json:"startMonth"`
	Dayofmonth int    `json:"dayOfMonth"`
}

type Yearlyschedule struct {
	Enabled    bool   `json:"enabled"`
	Dayofmonth int    `json:"dayOfMonth"`
	Month      string `json:"month"`
}

type Dailyschedule struct {
	Enabled bool `json:"enabled"`
}

type Schedule struct {
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	UniqueName         string              `json:"uniqueName"`
	OwnerID            string              `json:"ownerId"`
	SystemPolicy       bool                `json:"systemPolicy"`
	GlobalPolicy       bool                `json:"globalPolicy"`
	Datecreated        string              `json:"dateCreated"`
	Datemodified       string              `json:"dateModified"`
	Snapshottimeofday  *Snapshottimeofday  `json:"snapshotTimeOfDay"`
	Continuousschedule *Continuousschedule `json:"continuousSchedule"`
	Weeklyschedule     *Weeklyschedule     `json:"weeklySchedule"`
	Dailyschedule      *Dailyschedule      `json:"dailySchedule"`
	Monthlyschedule    *Monthlyschedule    `json:"monthlySchedule"`
	Quartelyschedule   *Quartelyschedule   `json:"quartelySchedule"`
	Yearlyschedule     *Yearlyschedule     `json:"yearlySchedule"`
	ReferenceCount     int                 `json:"referenceCount"`
	StartTime          string              `json:"startTime"`
	TimeZone           string              `json:"timeZone"`
}

type Timemachineinfo struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Slaid       string        `json:"slaId"`
	Schedule    Schedule      `json:"schedule"`
	Tags        []interface{} `json:"tags"`

	Autotunelogdrive bool `json:"autoTuneLogDrive"`
}

type Actionarguments struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Nodes struct {
	Properties       []interface{} `json:"properties"`
	Vmname           string        `json:"vmName,omitempty"`
	Networkprofileid string        `json:"networkProfileId,omitempty"`
	DatabaseServerID string        `json:"dbserverId,omitempty"`
}

// ProvisionDatabaseResponse structs
type ProvisionDatabaseResponse struct {
	Name                 string      `json:"name"`
	Workid               string      `json:"workId"`
	Operationid          string      `json:"operationId"`
	Dbserverid           string      `json:"dbserverId"`
	Message              interface{} `json:"message"`
	Entityid             string      `json:"entityId"`
	Entityname           string      `json:"entityName"`
	Entitytype           string      `json:"entityType"`
	Status               string      `json:"status"`
	Associatedoperations interface{} `json:"associatedOperations"`
	Dependencyreport     interface{} `json:"dependencyReport"`
}

// ListDatabaseParamsResponse structs
type ListDatabaseParamsResponse struct {
	Properties []DatabaseProperties `json:"properties"`
}
type DatabaseProperties struct {
	RefID                    string `json:"ref_id"`
	Name                     string `json:"name"`
	Type                     string `json:"type"`
	ValueType                string `json:"value_type"`
	Category                 string `json:"category"`
	Regex                    string `json:"regex"`
	Secure                   string `json:"secure"`
	Required                 string `json:"required"`
	Custom1                  string `json:"custom1"`
	Custom2                  string `json:"custom2"`
	Custom3                  string `json:"custom3"`
	DefaultValue             string `json:"default_value"`
	Sensitive                string `json:"sensitive"`
	DisplayName              string `json:"display_name"`
	Description              string `json:"description"`
	Index                    int    `json:"index"`
	Alias                    string `json:"alias"`
	ParameterizedDisplayName string `json:"parameterized_display_name"`
	ParameterizedDescription string `json:"parameterized_description"`
	Isduplicable             string `json:"isDuplicable"`
}

// ListDatabaseInstancesResponse structs
type ListDatabaseInstancesResponse struct {
	Databases []DatabaseInstance `json:"databases"`
}

type DatabaseInstance struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListDatabaseServerVMResponse structs
type ListDatabaseServerVMResponse struct {
	Dbserverclusters []interface{} `json:"dbserverClusters"`
	Dbservers        []Dbservers   `json:"dbservers"`
}
type DatabaseServerProperties struct {
	RefID       string      `json:"ref_id"`
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Secure      bool        `json:"secure"`
	Description interface{} `json:"description"`
}
type Metadata struct {
	Secureinfo                      interface{} `json:"secureInfo"`
	Info                            interface{} `json:"info"`
	Deregisterinfo                  interface{} `json:"deregisterInfo"`
	Databasetype                    string      `json:"databaseType"`
	Physicaleradrive                bool        `json:"physicalEraDrive"`
	Clustered                       bool        `json:"clustered"`
	Singleinstance                  bool        `json:"singleInstance"`
	Eradriveinitialised             bool        `json:"eraDriveInitialised"`
	Provisionoperationid            string      `json:"provisionOperationId"`
	Markedfordeletion               bool        `json:"markedForDeletion"`
	Associatedtimemachines          interface{} `json:"associatedTimeMachines"`
	Softwaresnaphotinterval         int         `json:"softwareSnaphotInterval"`
	Protectiondomainmigrationstatus interface{} `json:"protectionDomainMigrationStatus"`
	Lastclocksyncalerttime          interface{} `json:"lastClockSyncAlertTime"`
}
type Dbservers struct {
	ID                       string                     `json:"id"`
	Name                     string                     `json:"name"`
	Description              string                     `json:"description"`
	Ownerid                  string                     `json:"ownerId"`
	Datecreated              string                     `json:"dateCreated"`
	Datemodified             string                     `json:"dateModified"`
	Properties               []DatabaseServerProperties `json:"properties"`
	Tags                     []interface{}              `json:"tags"`
	Eracreated               bool                       `json:"eraCreated"`
	Internal                 bool                       `json:"internal"`
	Dbserverclusterid        interface{}                `json:"dbserverClusterId"`
	Vmclustername            string                     `json:"vmClusterName"`
	Vmclusteruuid            string                     `json:"vmClusterUuid"`
	Ipaddresses              []string                   `json:"ipAddresses"`
	Fqdns                    interface{}                `json:"fqdns"`
	Macaddresses             []string                   `json:"macAddresses"`
	Type                     string                     `json:"type"`
	Placeholder              bool                       `json:"placeholder"`
	Status                   string                     `json:"status"`
	Clientid                 string                     `json:"clientId"`
	Nxclusterid              string                     `json:"nxClusterId"`
	Eradriveid               string                     `json:"eraDriveId"`
	Eraversion               string                     `json:"eraVersion"`
	Vmtimezone               string                     `json:"vmTimeZone"`
	Vminfo                   interface{}                `json:"vmInfo"`
	Info                     interface{}                `json:"info"`
	Metadata                 Metadata                   `json:"metadata"`
	Metric                   interface{}                `json:"metric"`
	Lcmconfig                interface{}                `json:"lcmConfig"`
	Clustered                bool                       `json:"clustered"`
	Requestedversion         interface{}                `json:"requestedVersion"`
	IsServerDriven           bool                       `json:"is_server_driven"`
	AssociatedTimeMachineID  interface{}                `json:"associated_time_machine_id"`
	TimeMachineInfo          interface{}                `json:"time_machine_info"`
	Eradrive                 interface{}                `json:"eraDrive"`
	Databases                interface{}                `json:"databases"`
	Clones                   interface{}                `json:"clones"`
	Accesskey                interface{}                `json:"accessKey"`
	Softwareinstallations    interface{}                `json:"softwareInstallations"`
	Protectiondomainid       string                     `json:"protectionDomainId"`
	Protectiondomain         interface{}                `json:"protectionDomain"`
	Databasetype             string                     `json:"databaseType"`
	Accesskeyid              string                     `json:"accessKeyId"`
	Associatedtimemachineids interface{}                `json:"associatedTimeMachineIds"`
	Dbserverinvalideastate   bool                       `json:"dbserverInValidEaState"`
	Workingdirectory         string                     `json:"workingDirectory"`
}

// GetOperationRequest struct
type GetOperationRequest struct {
	OperationID string `json:"operation_id"`
}

// GetOperationResponse struct
type GetOperationResponse struct {
	Entityname              string            `json:"entityName"`
	Work                    interface{}       `json:"work"`
	Stepgenenabled          bool              `json:"stepGenEnabled"`
	Setstarttime            bool              `json:"setStartTime"`
	Timezone                string            `json:"timeZone"`
	ID                      string            `json:"id"`
	Name                    string            `json:"name"`
	Uniquename              interface{}       `json:"uniqueName"`
	Type                    string            `json:"type"`
	Starttime               string            `json:"startTime"`
	Timeout                 int               `json:"timeout"`
	Endtime                 string            `json:"endTime"`
	Instanceid              interface{}       `json:"instanceId"`
	Ownerid                 string            `json:"ownerId"`
	Status                  string            `json:"status"`
	Percentagecomplete      string            `json:"percentageComplete"`
	Steps                   []Steps           `json:"steps"`
	Properties              []interface{}     `json:"properties"`
	Parentid                interface{}       `json:"parentId"`
	Parentstep              int               `json:"parentStep"`
	Message                 interface{}       `json:"message"`
	Metadata                OperationMetadata `json:"metadata"`
	Entityid                string            `json:"entityId"`
	Entitytype              string            `json:"entityType"`
	Systemtriggered         bool              `json:"systemTriggered"`
	Uservisible             bool              `json:"userVisible"`
	Dbserverid              string            `json:"dbserverId"`
	Datesubmitted           string            `json:"dateSubmitted"`
	Deferredby              interface{}       `json:"deferredBy"`
	Scheduletime            interface{}       `json:"scheduleTime"`
	Isinternal              bool              `json:"isInternal"`
	Nxclusterid             string            `json:"nxClusterId"`
	Dbserverstatus          string            `json:"dbserverStatus"`
	Childoperations         []interface{}     `json:"childOperations"`
	Userrequestedaction     string            `json:"userRequestedAction"`
	Userrequestedactiontime interface{}       `json:"userRequestedActionTime"`
}
type Steps struct {
	Stepgenenabled     bool        `json:"stepGenEnabled"`
	Setstarttimevalue  bool        `json:"setStartTimeValue"`
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Uniquename         interface{} `json:"uniqueName"`
	Definitionid       string      `json:"definitionId"`
	Starttime          string      `json:"startTime"`
	Endtime            string      `json:"endTime"`
	Instanceid         interface{} `json:"instanceId"`
	Parentid           interface{} `json:"parentId"`
	Level              string      `json:"level"`
	Status             string      `json:"status"`
	Fileid             interface{} `json:"fileId"`
	Percentagecomplete string      `json:"percentageComplete"`
	Message            interface{} `json:"message"`
	Sequencenumber     int         `json:"sequenceNumber"`
	Childsteps         interface{} `json:"childSteps"`
	Weightage          int         `json:"weightage"`
}
type Executioncontext struct {
	Affecteddbservers         []string `json:"affectedDBServers"`
	Extendedaffecteddbservers []string `json:"extendedAffectedDBServers"`
	Applicationtype           string   `json:"applicationType"`
}
type OperationMetadata struct {
	Linkedoperations             interface{}      `json:"linkedOperations"`
	Associatedentities           interface{}      `json:"associatedEntities"`
	Oldstatus                    interface{}      `json:"oldStatus"`
	Userrequestedaction          string           `json:"userRequestedAction"`
	Userrequestedactiontimestamp interface{}      `json:"userRequestedActionTimestamp"`
	Controlmessage               interface{}      `json:"controlMessage"`
	Executioncontext             Executioncontext `json:"executionContext"`
	Scheduletime                 interface{}      `json:"scheduleTime"`
	Scheduledby                  string           `json:"scheduledBy"`
	Scheduledon                  string           `json:"scheduledOn"`
	Retryparentid                interface{}      `json:"retryParentId"`
	Retryimmediateparentid       interface{}      `json:"retryImmediateParentId"`
	Retriedoperations            interface{}      `json:"retriedOperations"`
	Switcheddbservers            interface{}      `json:"switchedDbservers"`
	Linkedoperationsdescription  string           `json:"linkedOperationsDescription"`
}

// Common Error response

type ErrorResponse struct {
	Errorcode            string        `json:"errorCode"`
	Reason               string        `json:"Reason"`
	Remedy               string        `json:"remedy"`
	Message              string        `json:"message"`
	Stacktrace           []interface{} `json:"stackTrace"`
	Suppressedexceptions []interface{} `json:"suppressedExceptions"`
}

// DeleteDatabase models

type DeleteDatabaseRequest struct {
	Delete               bool `json:"delete"`
	Remove               bool `json:"remove"`
	Softremove           bool `json:"softRemove"`
	Forced               bool `json:"forced"`
	Deletetimemachine    bool `json:"deleteTimeMachine"`
	Deletelogicalcluster bool `json:"deleteLogicalCluster"`
}

type DeleteDatabaseResponse struct {
	Name                 string      `json:"name"`
	Workid               string      `json:"workId"`
	Operationid          string      `json:"operationId"`
	Dbserverid           string      `json:"dbserverId"`
	Message              interface{} `json:"message"`
	Entityid             string      `json:"entityId"`
	Entityname           string      `json:"entityName"`
	Entitytype           string      `json:"entityType"`
	Status               string      `json:"status"`
	Associatedoperations interface{} `json:"associatedOperations"`
	Dependencyreport     interface{} `json:"dependencyReport"`
}

// UpdateDatabase models
type UpdateDatabaseRequest struct {
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	Tags             []interface{} `json:"tags"`
	Resetname        bool          `json:"resetName"`
	Resetdescription bool          `json:"resetDescription"`
	Resettags        bool          `json:"resetTags"`
}

type UpdateDatabaseResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EraDbExpiryDetails struct {
	RemindBeforeInDays int    `json:"remindBeforeInDays,omitempty"`
	EffectiveTimestamp string `json:"effectiveTimestamp,omitempty"`
	ExpiryTimestamp    string `json:"expiryTimestamp,omitempty"`
	ExpiryDateTimezone string `json:"expiryDateTimezone,omitempty"`
	UserCreated        bool   `json:"userCreated,omitempty"`
	ExpireInDays       int    `json:"expireInDays,omitempty"`
	DeleteDatabase     bool   `json:"deleteDatabase,omitempty"`
	DeleteTimeMachine  bool   `json:"deleteTimeMachine,omitempty"`
	DeleteVM           bool   `json:"deleteVM,omitempty"`
}

type EraDbRefreshDetails struct {
	RefreshInDays       int    `json:"refreshInDays,omitempty"`
	RefreshInHours      int    `json:"refreshInHours,omitempty"`
	RefreshInMonths     int    `json:"refreshInMonths,omitempty"`
	LastRefreshDate     string `json:"lastRefreshDate,omitempty"`
	NextRefreshDate     string `json:"nextRefreshDate,omitempty"`
	RefreshTime         string `json:"refreshTime,omitempty"`
	RefreshDateTimezone string `json:"refreshDateTimezone,omitempty"`
}

type EraDbPrePostDeleteCommand struct {
	Command string `json:"command,omitempty"`
}

type EraDbPostDeleteCommand struct{}

type LcmConfig struct {
	ExpiryDetails     *EraDbExpiryDetails        `json:"expiryDetails,omitempty"`
	RefreshDetails    *EraDbRefreshDetails       `json:"refreshDetails,omitempty"`
	PreDeleteCommand  *EraDbPrePostDeleteCommand `json:"preDeleteCommand,omitempty"`
	PostDeleteCommand *EraDbPrePostDeleteCommand `json:"postDeleteCommand,omitempty"`
}

type GetDatabaseResponse struct {
	ID                       string                 `json:"id"`
	Name                     string                 `json:"name"`
	Description              string                 `json:"description"`
	Ownerid                  string                 `json:"ownerId"`
	Datecreated              string                 `json:"dateCreated"`
	Datemodified             string                 `json:"dateModified"`
	AccessLevel              interface{}            `json:"accessLevel"`
	Properties               []DBInstanceProperties `json:"properties"`
	Tags                     []EraTags              `json:"tags"`
	Clustered                bool                   `json:"clustered"`
	Clone                    bool                   `json:"clone"`
	Eracreated               bool                   `json:"eraCreated"`
	Internal                 bool                   `json:"internal"`
	Placeholder              bool                   `json:"placeholder"`
	Databasename             string                 `json:"databaseName"`
	Type                     string                 `json:"type"`
	Databaseclustertype      interface{}            `json:"databaseClusterType"`
	Status                   string                 `json:"status"`
	Databasestatus           string                 `json:"databaseStatus"`
	Dbserverlogicalclusterid interface{}            `json:"dbserverLogicalClusterId"`
	Timemachineid            string                 `json:"timeMachineId"`
	Parenttimemachineid      interface{}            `json:"parentTimeMachineId"`
	Timezone                 string                 `json:"timeZone"`
	Info                     *Info                  `json:"info"`
	GroupInfo                interface{}            `json:"groupInfo"`
	Metadata                 *DBInstanceMetadata    `json:"metadata"`
	Metric                   interface{}            `json:"metric"`
	Category                 string                 `json:"category"`
	ParentDatabaseId         interface{}            `json:"parentDatabaseId,omitempty"`
	ParentSourceDatabaseId   interface{}            `json:"parentSourceDatabaseId,omitempty"`
	Lcmconfig                *LcmConfig             `json:"lcmConfig"`
	TimeMachine              *TimeMachine           `json:"timeMachine"`
	Dbserverlogicalcluster   interface{}            `json:"dbserverlogicalCluster"`
	Databasenodes            []Databasenodes        `json:"databaseNodes"`
	Linkeddatabases          []Linkeddatabases      `json:"linkedDatabases"`
	Databases                interface{}            `json:"databases,omitempty"`
	DatabaseGroupStateInfo   interface{}            `json:"databaseGroupStateInfo"`
}

type DBInstanceProperties struct {
	RefID       string      `json:"ref_id"`
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Secure      bool        `json:"secure"`
	Description interface{} `json:"description"`
}

type Secureinfo struct {
}

type DataDisks struct {
	Count float64 `json:"count"`
}

type LogDisks struct {
	Count float64 `json:"count"`
	Size  float64 `json:"size"`
}
type ArchiveStorage struct {
	Size float64 `json:"size"`
}
type Storage struct {
	DataDisks      *DataDisks      `json:"data_disks"`
	LogDisks       *LogDisks       `json:"log_disks"`
	ArchiveStorage *ArchiveStorage `json:"archive_storage"`
}
type VMProperties struct {
	NrHugepages             float64 `json:"nr_hugepages"`
	OvercommitMemory        float64 `json:"overcommit_memory"`
	DirtyBackgroundRatio    float64 `json:"dirty_background_ratio"`
	DirtyRatio              float64 `json:"dirty_ratio"`
	DirtyExpireCentisecs    float64 `json:"dirty_expire_centisecs"`
	DirtyWritebackCentisecs float64 `json:"dirty_writeback_centisecs"`
	Swappiness              float64 `json:"swappiness"`
}
type BpgDbParam struct {
	SharedBuffers               string `json:"shared_buffers"`
	MaintenanceWorkMem          string `json:"maintenance_work_mem"`
	WorkMem                     string `json:"work_mem"`
	EffectiveCacheSize          string `json:"effective_cache_size"`
	MaxWorkerProcesses          string `json:"max_worker_processes"`
	MaxParallelWorkersPerGather string `json:"max_parallel_workers_per_gather"`
}
type BpgConfigs struct {
	Storage      *Storage      `json:"storage"`
	VMProperties *VMProperties `json:"vm_properties"`
	BpgDbParam   *BpgDbParam   `json:"bpg_db_param"`
}
type InfoBpgConfig struct {
	BpgConfigs *BpgConfigs `json:"bpg_configs"`
}
type Info struct {
	Secureinfo interface{}    `json:"secureInfo"`
	Info       *InfoBpgConfig `json:"info"`
}
type DBInstanceMetadata struct {
	Secureinfo                          interface{} `json:"secureInfo,omitempty"`
	Info                                interface{} `json:"info,omitempty"`
	Deregisterinfo                      interface{} `json:"deregisterInfo,omitempty"`
	Tmactivateoperationid               string      `json:"tmActivateOperationId,omitempty"`
	Createddbservers                    interface{} `json:"createdDbservers,omitempty"`
	Registereddbservers                 interface{} `json:"registeredDbservers,omitempty"`
	Lastrefreshtimestamp                interface{} `json:"lastRefreshTimestamp,omitempty"`
	Lastrequestedrefreshtimestamp       interface{} `json:"lastRequestedRefreshTimestamp,omitempty"`
	CapabilityResetTime                 interface{} `json:"capabilityResetTime,omitempty"`
	Statebeforerefresh                  interface{} `json:"stateBeforeRefresh,omitempty"`
	Statebeforerestore                  interface{} `json:"stateBeforeRestore,omitempty"`
	Statebeforescaling                  interface{} `json:"stateBeforeScaling,omitempty"`
	Logcatchupforrestoredispatched      bool        `json:"logCatchUpForRestoreDispatched,omitempty"`
	Lastlogcatchupforrestoreoperationid interface{} `json:"lastLogCatchUpForRestoreOperationId,omitempty"`
	BaseSizeComputed                    bool        `json:"baseSizeComputed,omitempty"`
	Originaldatabasename                interface{} `json:"originalDatabaseName,omitempty"`
	ProvisionOperationId                string      `json:"provisionOperationId,omitempty"`
	SourceSnapshotId                    interface{} `json:"sourceSnapshotId,omitempty"`
	PitrBased                           bool        `json:"pitrBased,omitempty"`
	Sanitised                           bool        `json:"sanitised,omitempty"`
	RefreshBlockerInfo                  interface{} `json:"refreshBlockerInfo,omitempty"`
	DeregisteredWithDeleteTimeMachine   bool        `json:"deregisteredWithDeleteTimeMachine,omitempty"`
}

type DbserverMetadata struct {
	Secureinfo                      interface{} `json:"secureInfo"`
	Info                            interface{} `json:"info"`
	Deregisterinfo                  interface{} `json:"deregisterInfo"`
	Databasetype                    string      `json:"databaseType"`
	Physicaleradrive                bool        `json:"physicalEraDrive"`
	Clustered                       bool        `json:"clustered"`
	Singleinstance                  bool        `json:"singleInstance"`
	Eradriveinitialised             bool        `json:"eraDriveInitialised"`
	Provisionoperationid            string      `json:"provisionOperationId"`
	Markedfordeletion               bool        `json:"markedForDeletion"`
	Associatedtimemachines          interface{} `json:"associatedTimeMachines"`
	Softwaresnaphotinterval         int         `json:"softwareSnaphotInterval"`
	Protectiondomainmigrationstatus interface{} `json:"protectionDomainMigrationStatus"`
	Lastclocksyncalerttime          interface{} `json:"lastClockSyncAlertTime"`
}

type Dbserver struct {
	ID                       string           `json:"id"`
	Name                     string           `json:"name"`
	Description              string           `json:"description"`
	Ownerid                  string           `json:"ownerId"`
	Datecreated              string           `json:"dateCreated"`
	Datemodified             string           `json:"dateModified"`
	Properties               []Properties     `json:"properties"`
	Tags                     []interface{}    `json:"tags"`
	Eracreated               bool             `json:"eraCreated"`
	Internal                 bool             `json:"internal"`
	Dbserverclusterid        interface{}      `json:"dbserverClusterId"`
	Vmclustername            string           `json:"vmClusterName"`
	Vmclusteruuid            string           `json:"vmClusterUuid"`
	Ipaddresses              []string         `json:"ipAddresses"`
	Fqdns                    interface{}      `json:"fqdns"`
	Macaddresses             []string         `json:"macAddresses"`
	Type                     string           `json:"type"`
	Placeholder              bool             `json:"placeholder"`
	Status                   string           `json:"status"`
	Clientid                 string           `json:"clientId"`
	Nxclusterid              string           `json:"nxClusterId"`
	Eradriveid               string           `json:"eraDriveId"`
	Eraversion               string           `json:"eraVersion"`
	Vmtimezone               string           `json:"vmTimeZone"`
	Vminfo                   interface{}      `json:"vmInfo"`
	Info                     interface{}      `json:"info"`
	Metadata                 DbserverMetadata `json:"metadata"`
	Metric                   interface{}      `json:"metric"`
	Lcmconfig                interface{}      `json:"lcmConfig"`
	Clustered                bool             `json:"clustered"`
	Requestedversion         interface{}      `json:"requestedVersion"`
	IsServerDriven           bool             `json:"is_server_driven"`
	AssociatedTimeMachineID  interface{}      `json:"associated_time_machine_id"`
	TimeMachineInfo          interface{}      `json:"time_machine_info"`
	Eradrive                 interface{}      `json:"eraDrive"`
	Databases                interface{}      `json:"databases"`
	Clones                   interface{}      `json:"clones"`
	Accesskey                interface{}      `json:"accessKey"`
	Softwareinstallations    interface{}      `json:"softwareInstallations"`
	Protectiondomainid       string           `json:"protectionDomainId"`
	Protectiondomain         interface{}      `json:"protectionDomain"`
	Databasetype             string           `json:"databaseType"`
	Accesskeyid              string           `json:"accessKeyId"`
	Associatedtimemachineids interface{}      `json:"associatedTimeMachineIds"`
	Dbserverinvalideastate   bool             `json:"dbserverInValidEaState"`
	Workingdirectory         string           `json:"workingDirectory"`
}

type EraTags struct {
	TagId      string      `json:"tagId,omitempty"`
	EntityId   string      `json:"entityId,omitempty"`
	EntityType interface{} `json:"entityType,omitempty"`
	Value      string      `json:"value,omitempty"`
	TagName    string      `json:"tagName,omitempty"`
}

type Protectiondomain struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Eracreated    bool                   `json:"eraCreated"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	Cloudid       string                 `json:"cloudId"`
	Datecreated   string                 `json:"dateCreated"`
	Datemodified  string                 `json:"dateModified"`
	Ownerid       string                 `json:"ownerId"`
	Status        string                 `json:"status"`
	PrimaryHost   string                 `json:"primaryHost,omitempty"`
	Properties    []DBInstanceProperties `json:"properties"`
	Tags          []EraTags              `json:"tags,omitempty"`
	AssocEntities []string               `json:"assocEntities,omitempty"`
}
type Databasenodes struct {
	ID                     string            `json:"id"`
	Name                   string            `json:"name"`
	Description            string            `json:"description"`
	Ownerid                string            `json:"ownerId"`
	Datecreated            string            `json:"dateCreated"`
	Datemodified           string            `json:"dateModified"`
	AccessLevel            interface{}       `json:"accessLevel,omitempty"`
	Properties             []interface{}     `json:"properties"`
	Tags                   []EraTags         `json:"tags"`
	Databaseid             string            `json:"databaseId"`
	Status                 string            `json:"status"`
	Databasestatus         string            `json:"databaseStatus"`
	Primary                bool              `json:"primary"`
	Dbserverid             string            `json:"dbserverId"`
	Softwareinstallationid string            `json:"softwareInstallationId"`
	Protectiondomainid     string            `json:"protectionDomainId"`
	Info                   Info              `json:"info"`
	Metadata               interface{}       `json:"metadata"`
	Protectiondomain       *Protectiondomain `json:"protectionDomain"`
	// Valideastate           bool             `json:"validEaState"`
}

type Linkeddatabases struct {
	ID                     string      `json:"id"`
	Name                   string      `json:"name"`
	DatabaseName           string      `json:"databaseName,omitempty"`
	Description            string      `json:"description"`
	Status                 string      `json:"status"`
	Databasestatus         string      `json:"databaseStatus"`
	ParentDatabaseId       string      `json:"parentDatabaseId"`
	ParentLinkedDatabaseId string      `json:"parentLinkedDatabaseId"`
	Ownerid                string      `json:"ownerId"`
	Datecreated            string      `json:"dateCreated"`
	Datemodified           string      `json:"dateModified"`
	TimeZone               string      `json:"timeZone"`
	Info                   Info        `json:"info"`
	Metadata               interface{} `json:"metadata"`
	Metric                 interface{} `json:"metric"`
	SnapshotId             string      `json:"snapshotId"`
}

type TimeMachine struct {
	ID                  string                 `json:"id,omitempty"`
	Name                string                 `json:"name,omitempty"`
	Description         string                 `json:"description,omitempty"`
	OwnerID             string                 `json:"ownerId,omitempty"`
	DateCreated         string                 `json:"dateCreated,omitempty"`
	DateModified        string                 `json:"dateModified,omitempty"`
	AccessLevel         interface{}            `json:"accessLevel,omitempty"`
	Properties          []DBInstanceProperties `json:"properties,omitempty"`
	Tags                []EraTags              `json:"tags,omitempty"`
	Clustered           bool                   `json:"clustered,omitempty"`
	Clone               bool                   `json:"clone,omitempty"`
	Internal            bool                   `json:"internal,omitempty"`
	DatabaseID          string                 `json:"databaseId,omitempty"`
	Type                string                 `json:"type,omitempty"`
	Category            string                 `json:"category,omitempty"`
	Status              string                 `json:"status,omitempty"`
	EaStatus            string                 `json:"eaStatus,omitempty"`
	Scope               string                 `json:"scope,omitempty"`
	SLAID               string                 `json:"slaId,omitempty"`
	ScheduleID          string                 `json:"scheduleId,omitempty"`
	Info                interface{}            `json:"info,omitempty"`
	Metadata            *TimeMachineMetadata   `json:"metadata,omitempty"`
	Metric              interface{}            `json:"metric,omitempty"`
	SLAUpdateMetadata   interface{}            `json:"slaUpdateMetadata,omitempty"`
	SLA                 *ListSLAResponse       `json:"sla,omitempty"`
	Schedule            *Schedule              `json:"schedule,omitempty"`
	Database            interface{}            `json:"database,omitempty"`
	Clones              interface{}            `json:"clones,omitempty"`
	AssociatedClusters  interface{}            `json:"associatedClusters,omitempty"`
	SourceNxClusters    []string               `json:"sourceNxClusters,omitempty"`
	SLAUpdateInProgress bool                   `json:"slaUpdateInProgress,omitempty"`
}

type TimeMachineMetadata struct {
	SecureInfo                                          interface{} `json:"secureInfo,omitempty"`
	Info                                                interface{} `json:"info,omitempty"`
	DeregisterInfo                                      interface{} `json:"deregisterInfo,omitempty"`
	CapabilityResetTime                                 interface{} `json:"capabilityResetTime,omitempty"`
	AutoHeal                                            bool        `json:"autoHeal,omitempty"`
	AutoHealRetryCount                                  *int        `json:"autoHealRetryCount,omitempty"`
	AutoHealSnapshotCount                               *int        `json:"autoHealSnapshotCount,omitempty"`
	AutoHealLogCatchupCount                             *int        `json:"autoHealLogCatchupCount,omitempty"`
	FirstSnapshotCaptured                               bool        `json:"firstSnapshotCaptured,omitempty"`
	FirstSnapshotDispatched                             bool        `json:"firstSnapshotDispatched,omitempty"`
	LastSnapshotTime                                    string      `json:"lastSnapshotTime,omitempty"`
	LastAutoSnapshotTime                                string      `json:"lastAutoSnapshotTime,omitempty"`
	LastSnapshotOperationID                             string      `json:"lastSnapshotOperationId,omitempty"`
	LastAutoSnapshotOperationID                         string      `json:"lastAutoSnapshotOperationId,omitempty"`
	LastSuccessfulSnapshotOperationID                   string      `json:"lastSuccessfulSnapshotOperationId,omitempty"`
	SnapshotSuccessiveFailureCount                      *int        `json:"snapshotSuccessiveFailureCount,omitempty"`
	LastHealSnapshotOperation                           interface{} `json:"lastHealSnapshotOperation,omitempty"`
	DatabasesFirstSnapshotInfo                          interface{} `json:"databasesFirstSnapshotInfo,omitempty"`
	DispatchOnboardingSnapshot                          bool        `json:"dispatchOnboardingSnapshot,omitempty"`
	OnboardingSnapshotProperties                        interface{} `json:"onboardingSnapshotProperties,omitempty"`
	LastNonExtraAutoSnapshotTime                        string      `json:"lastNonExtraAutoSnapshotTime,omitempty"`
	SnapshotCapturedForTheDay                           bool        `json:"snapshotCapturedForTheDay,omitempty"`
	FirstSnapshotRetryCount                             *int        `json:"firstSnapshotRetryCount,omitempty"`
	LastLogCatchupTime                                  interface{} `json:"lastLogCatchupTime,omitempty"`
	LastSuccessfulLogCatchupOperationID                 interface{} `json:"lastSuccessfulLogCatchupOperationId,omitempty"`
	LastLogCatchupOperationID                           interface{} `json:"lastLogCatchupOperationId,omitempty"`
	LogCatchupSuccessiveFailureCount                    *int        `json:"logCatchupSuccessiveFailureCount,omitempty"`
	LastLogCatchupSkipped                               bool        `json:"lastLogCatchupSkipped,omitempty"`
	LastSuccessfulLogCatchupPostHealWithResetCapability interface{} `json:"lastSuccessfulLogCatchupPostHealWithResetCapability,omitempty"`
	LastPauseTime                                       interface{} `json:"lastPauseTime,omitempty"`
	LastPauseByForce                                    bool        `json:"lastPauseByForce,omitempty"`
	LastResumeTime                                      interface{} `json:"lastResumeTime,omitempty"`
	LastPauseReason                                     interface{} `json:"lastPauseReason,omitempty"`
	StateBeforeRestore                                  interface{} `json:"stateBeforeRestore,omitempty"`
	LastHealthAlertedTime                               interface{} `json:"lastHealthAlertedTime,omitempty"`
	ImplicitResumeCount                                 *int        `json:"implicitResumeCount,omitempty"`
	LastImplicitResumeTime                              interface{} `json:"lastImplicitResumeTime,omitempty"`
	StorageLimitExhausted                               bool        `json:"storageLimitExhausted,omitempty"`
	AbsoluteThresholdExhausted                          bool        `json:"absoluteThresholdExhausted,omitempty"`
	RequiredSpace                                       *float64    `json:"requiredSpace,omitempty"`
	LastEaBreakdownTime                                 interface{} `json:"lastEaBreakdownTime,omitempty"`
	AutoSnapshotRetryInfo                               interface{} `json:"autoSnapshotRetryInfo,omitempty"`
	AuthorizedDbservers                                 interface{} `json:"authorizedDbservers,omitempty"`
	LastHealTime                                        interface{} `json:"lastHealTime,omitempty"`
	LastHealSystemTriggered                             bool        `json:"lastHealSystemTriggered,omitempty"`
}
