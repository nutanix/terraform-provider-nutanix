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
	NxClusterID              *string              `json:"nxClusterId,omitempty"`
	DateCreated              *string              `json:"dateCreated,omitempty"`
	DateModified             *string              `json:"dateModified,omitempty"`
	OwnerID                  *string              `json:"ownerId,omitempty"`
	Status                   *string              `json:"status,omitempty"`
	ProfileVersionID         *string              `json:"profileVersionId,omitempty"`
	Properties               []*ProfileProperties `json:"properties,omitempty"`
	OptimizedForProvisioning bool                 `json:"optimizedForProvisioning,omitempty"`
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
	Properties                []*ProfileProperties         `json:"properties,omitempty"`
	VersionClusterAssociation []*VersionClusterAssociation `json:"versionClusterAssociation,omitempty"`
}

type ProfilesEntity struct {
	WindowsDomain     *int `json:"WindowsDomain,omitempty"`
	Software          *int `json:"Software,omitempty"`
	Compute           *int `json:"Compute,omitempty"`
	Network           *int `json:"Network,omitempty"`
	Storage           *int `json:"Storage,omitempty"`
	DatabaseParameter *int `json:"Database_Parameter,omitempty"`
}

type ProfileTimeMachinesCount struct {
	Profiles     *ProfilesEntity `json:"profiles,omitempty"`
	TimeMachines *int            `json:"timeMachines,omitempty"`
}

type EngineCounts struct {
	OracleDatabase    *ProfileTimeMachinesCount `json:"oracle_database,omitempty"`
	PostgresDatabase  *ProfileTimeMachinesCount `json:"postgres_database,omitempty"`
	MongodbDatabase   *ProfileTimeMachinesCount `json:"mongodb_database,omitempty"`
	SqlserverDatabase *ProfileTimeMachinesCount `json:"sqlserver_database,omitempty"`
	SaphanaDatabase   *ProfileTimeMachinesCount `json:"saphana_database,omitempty"`
	MariadbDatabase   *ProfileTimeMachinesCount `json:"mariadb_database,omitempty"`
	MySQLDatabase     *ProfileTimeMachinesCount `json:"mysql_database,omitempty"`
}

type EntityCounts struct {
	DBServers    *int          `json:"dbServers,omitempty"`
	EngineCounts *EngineCounts `json:"engineCounts,omitempty"`
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
	EntityCounts         *EntityCounts   `json:"entityCounts,omitempty"`
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

type PrePostCommand struct {
	PreCommand  *string `json:"preCommand,omitempty"`
	PostCommand *string `json:"postCommand,omitempty"`
}

type Payload struct {
	PrePostCommand *PrePostCommand `json:"prePostCommand,omitempty"`
}

type Tasks struct {
	TaskType *string  `json:"taskType,omitempty"`
	Payload  *Payload `json:"payload,omitempty"`
}

type MaintenanceTasks struct {
	MaintenanceWindowID *string  `json:"maintenanceWindowId,omitempty"`
	Tasks               []*Tasks `json:"tasks,omitempty"`
}

type ClusterIPInfos struct {
	NxClusterID *string    `json:"nxClusterId,omitempty"`
	IPInfos     []*IPInfos `json:"ipInfos,omitempty"`
}

type ClusterInfo struct {
	ClusterIPInfos []*ClusterIPInfos `json:"clusterIpInfos,omitempty"`
}

// ProvisionDatabaseRequestStructs
type ProvisionDatabaseRequest struct {
	Createdbserver           bool               `json:"createDbserver,omitempty"`
	Clustered                bool               `json:"clustered,omitempty"`
	Autotunestagingdrive     bool               `json:"autoTuneStagingDrive,omitempty"`
	Nodecount                *int               `json:"nodeCount,omitempty"`
	Databasetype             *string            `json:"databaseType,omitempty"`
	Name                     *string            `json:"name,omitempty"`
	Databasedescription      *string            `json:"databaseDescription,omitempty"`
	DatabaseServerID         *string            `json:"dbserverId,omitempty"`
	Softwareprofileid        *string            `json:"softwareProfileId,omitempty"`
	Softwareprofileversionid *string            `json:"softwareProfileVersionId,omitempty"`
	Computeprofileid         *string            `json:"computeProfileId,omitempty"`
	Networkprofileid         *string            `json:"networkProfileId,omitempty"`
	Dbparameterprofileid     *string            `json:"dbParameterProfileId,omitempty"`
	Newdbservertimezone      *string            `json:"newFVMbServerTimeZone,omitempty"`
	Nxclusterid              *string            `json:"nxClusterId,omitempty"`
	Sshpublickey             *string            `json:"sshPublicKey,omitempty"`
	VMPassword               *string            `json:"vmPassword,omitempty"`
	Timemachineinfo          *Timemachineinfo   `json:"timeMachineInfo,omitempty"`
	Actionarguments          []*Actionarguments `json:"actionArguments,omitempty"`
	Nodes                    []*Nodes           `json:"nodes,omitempty"`
	Tags                     []*Tags            `json:"tags,omitempty"`
	MaintenanceTasks         *MaintenanceTasks  `json:"maintenanceTasks,omitempty"`
	ClusterInfo              *ClusterInfo       `json:"clusterInfo,omitempty"`
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
	ID                 *string             `json:"id,omitempty"`
	Name               *string             `json:"name,omitempty"`
	Description        *string             `json:"description,omitempty"`
	UniqueName         *string             `json:"uniqueName,omitempty"`
	OwnerID            *string             `json:"ownerId,omitempty"`
	StartTime          *string             `json:"startTime,omitempty"`
	TimeZone           *string             `json:"timeZone,omitempty"`
	Datecreated        *string             `json:"dateCreated,omitempty"`
	Datemodified       *string             `json:"dateModified,omitempty"`
	ReferenceCount     *int                `json:"referenceCount,omitempty"`
	SystemPolicy       bool                `json:"systemPolicy,omitempty"`
	GlobalPolicy       bool                `json:"globalPolicy,omitempty"`
	Snapshottimeofday  *Snapshottimeofday  `json:"snapshotTimeOfDay,omitempty"`
	Continuousschedule *Continuousschedule `json:"continuousSchedule,omitempty"`
	Weeklyschedule     *Weeklyschedule     `json:"weeklySchedule,omitempty"`
	Dailyschedule      *Dailyschedule      `json:"dailySchedule,omitempty"`
	Monthlyschedule    *Monthlyschedule    `json:"monthlySchedule,omitempty"`
	Quartelyschedule   *Quartelyschedule   `json:"quartelySchedule,omitempty"`
	Yearlyschedule     *Yearlyschedule     `json:"yearlySchedule,omitempty"`
}

type PrimarySLA struct {
	SLAID        *string   `json:"slaId,omitempty"`
	NxClusterIds []*string `json:"nxClusterIds,omitempty"`
}

type SLADetails struct {
	PrimarySLA *PrimarySLA `json:"primarySla,omitempty"`
}

type Timemachineinfo struct {
	Name             string      `json:"name,omitempty"`
	Description      string      `json:"description,omitempty"`
	Slaid            string      `json:"slaId,omitempty"`
	Schedule         Schedule    `json:"schedule,omitempty"`
	Tags             []*Tags     `json:"tags,omitempty"`
	Autotunelogdrive bool        `json:"autoTuneLogDrive,omitempty"`
	SLADetails       *SLADetails `json:"slaDetails,omitempty"`
}

type Actionarguments struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type NodesProperties struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type IPInfos struct {
	IPType      *string   `json:"ipType,omitempty"`
	IPAddresses []*string `json:"ipAddresses,omitempty"`
}

type Nodes struct {
	Properties          []*NodesProperties `json:"properties"`
	Vmname              *string            `json:"vmName,omitempty"`
	Networkprofileid    *string            `json:"networkProfileId,omitempty"`
	DatabaseServerID    *string            `json:"dbserverId,omitempty"`
	NxClusterID         *string            `json:"nxClusterId,omitempty"`
	ComputeProfileID    *string            `json:"computeProfileId,omitempty"`
	NewDBServerTimeZone *string            `json:"newDbServerTimeZone,omitempty"`
	IPInfos             []*IPInfos         `json:"ipInfos,omitempty"`
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
type DBServerMetadata struct {
	Physicaleradrive        bool            `json:"physicalEraDrive"`
	Clustered               bool            `json:"clustered"`
	Singleinstance          bool            `json:"singleInstance"`
	Eradriveinitialised     bool            `json:"eraDriveInitialised"`
	Markedfordeletion       bool            `json:"markedForDeletion"`
	Softwaresnaphotinterval int             `json:"softwareSnaphotInterval"`
	Databasetype            *string         `json:"databaseType"`
	Provisionoperationid    *string         `json:"provisionOperationId"`
	Associatedtimemachines  []*string       `json:"associatedTimeMachines"`
	Secureinfo              *Secureinfo     `json:"secureInfo"`
	Info                    *Info           `json:"info"`
	Deregisterinfo          *DeregisterInfo `json:"deregisterInfo"`
	// Protectiondomainmigrationstatus interface{}     `json:"protectionDomainMigrationStatus"`
	// Lastclocksyncalerttime          interface{}     `json:"lastClockSyncAlertTime"`
}
type Dbservers struct {
	Placeholder              bool                        `json:"placeholder"`
	Clustered                bool                        `json:"clustered"`
	Eracreated               bool                        `json:"eraCreated"`
	Internal                 bool                        `json:"internal"`
	IsServerDriven           bool                        `json:"is_server_driven"`
	Dbserverinvalideastate   bool                        `json:"dbserverInValidEaState"`
	ID                       *string                     `json:"id"`
	Name                     *string                     `json:"name"`
	Description              *string                     `json:"description"`
	Ownerid                  *string                     `json:"ownerId"`
	Datecreated              *string                     `json:"dateCreated"`
	Datemodified             *string                     `json:"dateModified"`
	Dbserverclusterid        *string                     `json:"dbserverClusterId"`
	Vmclustername            *string                     `json:"vmClusterName"`
	Vmclusteruuid            *string                     `json:"vmClusterUuid"`
	Type                     *string                     `json:"type"`
	Status                   *string                     `json:"status"`
	Clientid                 *string                     `json:"clientId"`
	Nxclusterid              *string                     `json:"nxClusterId"`
	Eradriveid               *string                     `json:"eraDriveId"`
	Eraversion               *string                     `json:"eraVersion"`
	Vmtimezone               *string                     `json:"vmTimeZone"`
	Accesskey                *string                     `json:"accessKey"`
	Protectiondomainid       *string                     `json:"protectionDomainId"`
	Databasetype             *string                     `json:"databaseType"`
	Accesskeyid              *string                     `json:"accessKeyId"`
	Requestedversion         *string                     `json:"requestedVersion"`
	AssociatedTimeMachineID  *string                     `json:"associated_time_machine_id"`
	Workingdirectory         *string                     `json:"workingDirectory"`
	Ipaddresses              []*string                   `json:"ipAddresses"`
	Fqdns                    []*string                   `json:"fqdns"`
	Macaddresses             []*string                   `json:"macAddresses"`
	Associatedtimemachineids []*string                   `json:"associatedTimeMachineIds"`
	Properties               []*DatabaseServerProperties `json:"properties"`
	Tags                     []*Tags                     `json:"tags"`
	Vminfo                   *VMInfo                     `json:"vmInfo"`
	Info                     *Info                       `json:"info"`
	Metadata                 *DBServerMetadata           `json:"metadata"`
	Metric                   *Metric                     `json:"metric"`
	Lcmconfig                *LcmConfig                  `json:"lcmConfig"`
	TimeMachineInfo          []*Properties               `json:"time_machine_info"`
	Eradrive                 interface{}                 `json:"eraDrive"`
	Databases                interface{}                 `json:"databases"`
	Clones                   interface{}                 `json:"clones"`
	Softwareinstallations    interface{}                 `json:"softwareInstallations"`
	// Protectiondomain         interface{}                 `json:"protectionDomain"`
}

// GetOperationRequest struct
type GetOperationRequest struct {
	OperationID string `json:"operation_id"`
}

// GetOperationResponse struct
type GetOperationResponse struct {
	Stepgenenabled          bool               `json:"stepGenEnabled"`
	Setstarttime            bool               `json:"setStartTime"`
	Systemtriggered         bool               `json:"systemTriggered"`
	Uservisible             bool               `json:"userVisible"`
	Isinternal              bool               `json:"isInternal"`
	Timeout                 int                `json:"timeout"`
	Parentstep              int                `json:"parentStep"`
	Entityname              *string            `json:"entityName"`
	Timezone                *string            `json:"timeZone"`
	ID                      *string            `json:"id"`
	Name                    *string            `json:"name"`
	Uniquename              *string            `json:"uniqueName"`
	Type                    *string            `json:"type"`
	Starttime               *string            `json:"startTime"`
	Endtime                 *string            `json:"endTime"`
	Instanceid              *string            `json:"instanceId"`
	Ownerid                 *string            `json:"ownerId"`
	Status                  *string            `json:"status"`
	Percentagecomplete      *string            `json:"percentageComplete"`
	Parentid                *string            `json:"parentId"`
	Message                 *string            `json:"message"`
	Scheduletime            *string            `json:"scheduleTime"`
	Nxclusterid             *string            `json:"nxClusterId"`
	Dbserverstatus          *string            `json:"dbserverStatus"`
	Userrequestedaction     *string            `json:"userRequestedAction"`
	Userrequestedactiontime *string            `json:"userRequestedActionTime"`
	Entityid                *string            `json:"entityId"`
	Entitytype              *string            `json:"entityType"`
	Dbserverid              *string            `json:"dbserverId"`
	Datesubmitted           *string            `json:"dateSubmitted"`
	Deferredby              *string            `json:"deferredBy"`
	DeferredByOpIDs         []*string          `json:"deferredByOpIds"`
	Steps                   []*Steps           `json:"steps"`
	Properties              []*Properties      `json:"properties"`
	Metadata                *OperationMetadata `json:"metadata"`
	Work                    interface{}        `json:"work"`
	Childoperations         []interface{}      `json:"childOperations"`
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
	Name             string  `json:"name,omitempty"`
	Description      string  `json:"description,omitempty"`
	Tags             []*Tags `json:"tags,omitempty"`
	Resetname        bool    `json:"resetName,omitempty"`
	Resetdescription bool    `json:"resetDescription,omitempty"`
	Resettags        bool    `json:"resetTags,omitempty"`
}

type UpdateDatabaseResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DBExpiryDetails struct {
	EffectiveTimestamp *string `json:"effectiveTimestamp,omitempty"`
	ExpiryTimestamp    *string `json:"expiryTimestamp,omitempty"`
	ExpiryDateTimezone *string `json:"expiryDateTimezone,omitempty"`
	RemindBeforeInDays *int    `json:"remindBeforeInDays,omitempty"`
	ExpireInDays       *int    `json:"expireInDays,omitempty"`
	DeleteDatabase     bool    `json:"deleteDatabase,omitempty"`
	DeleteTimeMachine  bool    `json:"deleteTimeMachine,omitempty"`
	DeleteVM           bool    `json:"deleteVM,omitempty"`
	UserCreated        bool    `json:"userCreated,omitempty"`
}

type DBRefreshDetails struct {
	RefreshInDays       int    `json:"refreshInDays,omitempty"`
	RefreshInHours      int    `json:"refreshInHours,omitempty"`
	RefreshInMonths     int    `json:"refreshInMonths,omitempty"`
	LastRefreshDate     string `json:"lastRefreshDate,omitempty"`
	NextRefreshDate     string `json:"nextRefreshDate,omitempty"`
	RefreshTime         string `json:"refreshTime,omitempty"`
	RefreshDateTimezone string `json:"refreshDateTimezone,omitempty"`
}

type DBPrePostDeleteCommand struct {
	Command string `json:"command,omitempty"`
}

type DBPostDeleteCommand struct{}

type LcmConfig struct {
	ExpiryDetails     *DBExpiryDetails        `json:"expiryDetails,omitempty"`
	RefreshDetails    *DBRefreshDetails       `json:"refreshDetails,omitempty"`
	PreDeleteCommand  *DBPrePostDeleteCommand `json:"preDeleteCommand,omitempty"`
	PostDeleteCommand *DBPrePostDeleteCommand `json:"postDeleteCommand,omitempty"`
}

type ListDatabaseInstance []GetDatabaseResponse

type GetDatabaseResponse struct {
	ID                       string                  `json:"id"`
	Name                     string                  `json:"name"`
	Description              string                  `json:"description"`
	Ownerid                  string                  `json:"ownerId"`
	Datecreated              string                  `json:"dateCreated"`
	Datemodified             string                  `json:"dateModified"`
	AccessLevel              interface{}             `json:"accessLevel"`
	Properties               []*DBInstanceProperties `json:"properties"`
	Tags                     []*Tags                 `json:"tags"`
	Clustered                bool                    `json:"clustered"`
	Clone                    bool                    `json:"clone"`
	Eracreated               bool                    `json:"eraCreated"`
	Internal                 bool                    `json:"internal"`
	Placeholder              bool                    `json:"placeholder"`
	Databasename             string                  `json:"databaseName"`
	Type                     string                  `json:"type"`
	Databaseclustertype      interface{}             `json:"databaseClusterType"`
	Status                   string                  `json:"status"`
	Databasestatus           string                  `json:"databaseStatus"`
	Dbserverlogicalclusterid interface{}             `json:"dbserverLogicalClusterId"`
	Timemachineid            string                  `json:"timeMachineId"`
	Parenttimemachineid      interface{}             `json:"parentTimeMachineId"`
	Timezone                 string                  `json:"timeZone"`
	Info                     *Info                   `json:"info"`
	GroupInfo                interface{}             `json:"groupInfo"`
	Metadata                 *DBInstanceMetadata     `json:"metadata"`
	Metric                   interface{}             `json:"metric"`
	Category                 string                  `json:"category"`
	ParentDatabaseID         interface{}             `json:"parentDatabaseId,omitempty"`
	ParentSourceDatabaseID   interface{}             `json:"parentSourceDatabaseId,omitempty"`
	Lcmconfig                *LcmConfig              `json:"lcmConfig"`
	TimeMachine              *TimeMachine            `json:"timeMachine"`
	Dbserverlogicalcluster   interface{}             `json:"dbserverlogicalCluster"`
	Databasenodes            []Databasenodes         `json:"databaseNodes"`
	Linkeddatabases          []Linkeddatabases       `json:"linkedDatabases"`
	Databases                interface{}             `json:"databases,omitempty"`
	DatabaseGroupStateInfo   interface{}             `json:"databaseGroupStateInfo"`
}

type DBInstanceProperties struct {
	RefID       string      `json:"ref_id"`
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Secure      bool        `json:"secure"`
	Description interface{} `json:"description"`
}

type Secureinfo struct{}

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
type BpgDBParam struct {
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
	BpgDBParam   *BpgDBParam   `json:"bpg_db_param"`
}
type InfoBpgConfig struct {
	CreatedBy  *string     `json:"created_by,omitempty"`
	BpgConfigs *BpgConfigs `json:"bpg_configs"`
}
type Info struct {
	Secureinfo interface{}    `json:"secureInfo"`
	Info       *InfoBpgConfig `json:"info"`
	CreatedBy  *string        `json:"created_by,omitempty"`
}
type DBInstanceMetadata struct {
	Logcatchupforrestoredispatched      bool            `json:"logCatchUpForRestoreDispatched,omitempty"`
	BaseSizeComputed                    bool            `json:"baseSizeComputed,omitempty"`
	PitrBased                           bool            `json:"pitrBased,omitempty"`
	DeregisteredWithDeleteTimeMachine   bool            `json:"deregisteredWithDeleteTimeMachine,omitempty"`
	Lastrefreshtimestamp                *string         `json:"lastRefreshTimestamp,omitempty"`
	Lastrequestedrefreshtimestamp       *string         `json:"lastRequestedRefreshTimestamp,omitempty"`
	Statebeforerefresh                  *string         `json:"stateBeforeRefresh,omitempty"`
	Statebeforerestore                  *string         `json:"stateBeforeRestore,omitempty"`
	Statebeforescaling                  *string         `json:"stateBeforeScaling,omitempty"`
	Lastlogcatchupforrestoreoperationid *string         `json:"lastLogCatchUpForRestoreOperationId,omitempty"`
	ProvisionOperationID                *string         `json:"provisionOperationId,omitempty"`
	SourceSnapshotID                    *string         `json:"sourceSnapshotId,omitempty"`
	Tmactivateoperationid               *string         `json:"tmActivateOperationId,omitempty"`
	Createddbservers                    []*string       `json:"createdDbservers,omitempty"`
	Secureinfo                          *Secureinfo     `json:"secureInfo,omitempty"`
	Info                                *Info           `json:"info,omitempty"`
	Deregisterinfo                      *DeregisterInfo `json:"deregisterInfo,omitempty"`
	Registereddbservers                 interface{}     `json:"registeredDbservers,omitempty"`
	CapabilityResetTime                 interface{}     `json:"capabilityResetTime,omitempty"`
	Originaldatabasename                interface{}     `json:"originalDatabaseName,omitempty"`
	RefreshBlockerInfo                  interface{}     `json:"refreshBlockerInfo,omitempty"`
}

type DbserverMetadata struct {
	Physicaleradrive        bool            `json:"physicalEraDrive"`
	Clustered               bool            `json:"clustered"`
	Singleinstance          bool            `json:"singleInstance"`
	Eradriveinitialised     bool            `json:"eraDriveInitialised"`
	Markedfordeletion       bool            `json:"markedForDeletion"`
	Softwaresnaphotinterval int             `json:"softwareSnaphotInterval"`
	Databasetype            *string         `json:"databaseType"`
	Provisionoperationid    *string         `json:"provisionOperationId"`
	Associatedtimemachines  []*string       `json:"associatedTimeMachines"`
	Secureinfo              *Secureinfo     `json:"secureInfo"`
	Info                    *Info           `json:"info"`
	Deregisterinfo          *DeregisterInfo `json:"deregisterInfo"`
	// Protectiondomainmigrationstatus interface{}     `json:"protectionDomainMigrationStatus"`
	// Lastclocksyncalerttime          interface{}     `json:"lastClockSyncAlertTime"`
}

// type VMInfo struct {
// 	OsType       *string `json:"osType,omitempty"`
// 	OsVersion    *string `json:"osVersion,omitempty"`
// 	Distribution *string `json:"distribution,omitempty"`
// }

type MetricVMInfo struct {
	NumVCPUs              *int    `json:"numVCPUs,omitempty"`
	NumCoresPerVCPU       *int    `json:"numCoresPerVCPU,omitempty"`
	HypervisorCPUUsagePpm []*int  `json:"hypervisorCpuUsagePpm,omitempty"`
	LastUpdatedTimeInUTC  *string `json:"lastUpdatedTimeInUTC,omitempty"`
}

type MetricMemoryInfo struct {
	LastUpdatedTimeInUTC *string `json:"lastUpdatedTimeInUTC,omitempty"`
	Memory               *int    `json:"memory,omitempty"`
	MemoryUsagePpm       []*int  `json:"memoryUsagePpm,omitempty"`
	Unit                 *string `json:"unit,omitempty"`
}

type MetricStorageInfo struct {
	LastUpdatedTimeInUTC        interface{} `json:"lastUpdatedTimeInUTC,omitempty"`
	ControllerNumIops           []*int      `json:"controllerNumIops,omitempty"`
	ControllerAvgIoLatencyUsecs []*int      `json:"controllerAvgIoLatencyUsecs,omitempty"`
	Size                        interface{} `json:"size,omitempty"`
	AllocatedSize               interface{} `json:"allocatedSize,omitempty"`
	UsedSize                    interface{} `json:"usedSize,omitempty"`
	Unit                        interface{} `json:"unit,omitempty"`
}

type Metric struct {
	LastUpdatedTimeInUTC *string            `json:"lastUpdatedTimeInUTC,omitempty"`
	Compute              *MetricVMInfo      `json:"compute,omitempty"`
	Memory               *MetricMemoryInfo  `json:"memory,omitempty"`
	Storage              *MetricStorageInfo `json:"storage,omitempty"`
}

type Dbserver struct {
	Placeholder              bool              `json:"placeholder,omitempty"`
	Eracreated               bool              `json:"eraCreated,omitempty"`
	Internal                 bool              `json:"internal,omitempty"`
	IsServerDriven           bool              `json:"is_server_driven,omitempty"`
	Clustered                bool              `json:"clustered,omitempty"`
	Dbserverinvalideastate   bool              `json:"dbserverInValidEaState,omitempty"`
	ID                       *string           `json:"id,omitempty"`
	Name                     *string           `json:"name,omitempty"`
	Description              *string           `json:"description,omitempty"`
	Ownerid                  *string           `json:"ownerId,omitempty"`
	Datecreated              *string           `json:"dateCreated,omitempty"`
	Datemodified             *string           `json:"dateModified,omitempty"`
	Vmclustername            *string           `json:"vmClusterName,omitempty"`
	Vmclusteruuid            *string           `json:"vmClusterUuid,omitempty"`
	Type                     *string           `json:"type,omitempty"`
	Status                   *string           `json:"status,omitempty"`
	Clientid                 *string           `json:"clientId,omitempty"`
	Nxclusterid              *string           `json:"nxClusterId,omitempty"`
	Eradriveid               *string           `json:"eraDriveId,omitempty"`
	Eraversion               *string           `json:"eraVersion,omitempty"`
	Vmtimezone               *string           `json:"vmTimeZone,omitempty"`
	Requestedversion         *string           `json:"requestedVersion,omitempty"`
	AssociatedTimeMachineID  *string           `json:"associated_time_machine_id,omitempty"`
	Accesskey                *string           `json:"accessKey,omitempty"`
	Protectiondomainid       *string           `json:"protectionDomainId,omitempty"`
	Databasetype             *string           `json:"databaseType,omitempty"`
	Accesskeyid              *string           `json:"accessKeyId,omitempty"`
	Associatedtimemachineids []*string         `json:"associatedTimeMachineIds,omitempty"`
	Workingdirectory         *string           `json:"workingDirectory,omitempty"`
	Ipaddresses              []*string         `json:"ipAddresses,omitempty"`
	Macaddresses             []*string         `json:"macAddresses,omitempty"`
	Vminfo                   *VMInfo           `json:"vmInfo,omitempty"`
	Info                     *Info             `json:"info,omitempty"`
	Metadata                 *DbserverMetadata `json:"metadata,omitempty"`
	Metric                   *Metric           `json:"metric,omitempty"`
	Lcmconfig                *LcmConfig        `json:"lcmConfig,omitempty"`
	TimeMachineInfo          []*Properties     `json:"time_machine_info"`
	Properties               []*Properties     `json:"properties,omitempty"`
	Eradrive                 interface{}       `json:"eraDrive,omitempty"`
	Databases                interface{}       `json:"databases,omitempty"`
	Clones                   interface{}       `json:"clones,omitempty"`
	Softwareinstallations    interface{}       `json:"softwareInstallations,omitempty"`
	Protectiondomain         interface{}       `json:"protectionDomain,omitempty"`
	Dbserverclusterid        interface{}       `json:"dbserverClusterId,omitempty"`
	Fqdns                    interface{}       `json:"fqdns,omitempty"`
	Tags                     []interface{}     `json:"tags,omitempty"`
}

type Tags struct {
	TagID      string      `json:"tagId,omitempty"`
	EntityID   string      `json:"entityId,omitempty"`
	EntityType interface{} `json:"entityType,omitempty"`
	Value      string      `json:"value,omitempty"`
	TagName    string      `json:"tagName,omitempty"`
}

type Protectiondomain struct {
	ID            string                  `json:"id"`
	Name          string                  `json:"name"`
	Eracreated    bool                    `json:"eraCreated"`
	Description   string                  `json:"description"`
	Type          string                  `json:"type"`
	Cloudid       string                  `json:"cloudId"`
	Datecreated   string                  `json:"dateCreated"`
	Datemodified  string                  `json:"dateModified"`
	Ownerid       string                  `json:"ownerId"`
	Status        string                  `json:"status"`
	PrimaryHost   string                  `json:"primaryHost,omitempty"`
	Properties    []*DBInstanceProperties `json:"properties"`
	Tags          []*Tags                 `json:"tags,omitempty"`
	AssocEntities []string                `json:"assocEntities,omitempty"`
}
type Databasenodes struct {
	ID                     string                  `json:"id"`
	Name                   string                  `json:"name"`
	Description            string                  `json:"description"`
	Ownerid                string                  `json:"ownerId"`
	Datecreated            string                  `json:"dateCreated"`
	Datemodified           string                  `json:"dateModified"`
	AccessLevel            interface{}             `json:"accessLevel,omitempty"`
	Properties             []*DBInstanceProperties `json:"properties"`
	Tags                   []*Tags                 `json:"tags"`
	Databaseid             string                  `json:"databaseId"`
	Status                 string                  `json:"status"`
	Databasestatus         string                  `json:"databaseStatus"`
	Primary                bool                    `json:"primary"`
	Dbserverid             string                  `json:"dbserverId"`
	Softwareinstallationid string                  `json:"softwareInstallationId"`
	Protectiondomainid     string                  `json:"protectionDomainId"`
	Info                   Info                    `json:"info"`
	Metadata               interface{}             `json:"metadata"`
	Protectiondomain       *Protectiondomain       `json:"protectionDomain"`
	// Valideastate           bool             `json:"validEaState"`
}

type Linkeddatabases struct {
	ID                     string      `json:"id"`
	Name                   string      `json:"name"`
	DatabaseName           string      `json:"databaseName,omitempty"`
	Description            string      `json:"description"`
	Status                 string      `json:"status"`
	Databasestatus         string      `json:"databaseStatus"`
	ParentDatabaseID       string      `json:"parentDatabaseId"`
	ParentLinkedDatabaseID string      `json:"parentLinkedDatabaseId"`
	Ownerid                string      `json:"ownerId"`
	Datecreated            string      `json:"dateCreated"`
	Datemodified           string      `json:"dateModified"`
	TimeZone               string      `json:"timeZone"`
	Info                   Info        `json:"info"`
	Metadata               interface{} `json:"metadata"`
	Metric                 interface{} `json:"metric"`
	SnapshotID             string      `json:"snapshotId"`
}

type TimeMachine struct {
	SLAUpdateInProgress bool                    `json:"slaUpdateInProgress,omitempty"`
	Clustered           bool                    `json:"clustered,omitempty"`
	Clone               bool                    `json:"clone,omitempty"`
	Internal            bool                    `json:"internal,omitempty"`
	ID                  *string                 `json:"id,omitempty"`
	Name                *string                 `json:"name,omitempty"`
	Description         *string                 `json:"description,omitempty"`
	OwnerID             *string                 `json:"ownerId,omitempty"`
	DateCreated         *string                 `json:"dateCreated,omitempty"`
	DateModified        *string                 `json:"dateModified,omitempty"`
	DatabaseID          *string                 `json:"databaseId,omitempty"`
	Type                *string                 `json:"type,omitempty"`
	Category            *string                 `json:"category,omitempty"`
	Status              *string                 `json:"status,omitempty"`
	EaStatus            *string                 `json:"eaStatus,omitempty"`
	Scope               *string                 `json:"scope,omitempty"`
	SLAID               *string                 `json:"slaId,omitempty"`
	ScheduleID          *string                 `json:"scheduleId,omitempty"`
	SourceNxClusters    []*string               `json:"sourceNxClusters,omitempty"`
	Properties          []*DBInstanceProperties `json:"properties,omitempty"`
	Tags                []*Tags                 `json:"tags,omitempty"`
	Info                *Info                   `json:"info,omitempty"`
	Metadata            *TimeMachineMetadata    `json:"metadata,omitempty"`
	SLA                 *ListSLAResponse        `json:"sla,omitempty"`
	Schedule            *Schedule               `json:"schedule,omitempty"`
	Database            *DatabaseInstance       `json:"database,omitempty"`
	Clones              interface{}             `json:"clones,omitempty"`
	AccessLevel         interface{}             `json:"accessLevel,omitempty"`
	Metric              interface{}             `json:"metric,omitempty"`
	//AssociatedClusters  interface{}             `json:"associatedClusters,omitempty"`
	// SLAUpdateMetadata   interface{}             `json:"slaUpdateMetadata,omitempty"`
}

type DeregisterInfo struct {
	Message    *string   `json:"message,omitempty"`
	Operations []*string `json:"operations,omitempty"`
}

type TimeMachineMetadata struct {
	LastHealSystemTriggered                             bool            `json:"lastHealSystemTriggered,omitempty"`
	AutoHeal                                            bool            `json:"autoHeal,omitempty"`
	DispatchOnboardingSnapshot                          bool            `json:"dispatchOnboardingSnapshot,omitempty"`
	LastLogCatchupSkipped                               bool            `json:"lastLogCatchupSkipped,omitempty"`
	FirstSnapshotCaptured                               bool            `json:"firstSnapshotCaptured,omitempty"`
	FirstSnapshotDispatched                             bool            `json:"firstSnapshotDispatched,omitempty"`
	StorageLimitExhausted                               bool            `json:"storageLimitExhausted,omitempty"`
	AbsoluteThresholdExhausted                          bool            `json:"absoluteThresholdExhausted,omitempty"`
	SnapshotCapturedForTheDay                           bool            `json:"snapshotCapturedForTheDay,omitempty"`
	LastPauseByForce                                    bool            `json:"lastPauseByForce,omitempty"`
	AutoHealRetryCount                                  *int            `json:"autoHealRetryCount,omitempty"`
	AutoHealSnapshotCount                               *int            `json:"autoHealSnapshotCount,omitempty"`
	AutoHealLogCatchupCount                             *int            `json:"autoHealLogCatchupCount,omitempty"`
	SnapshotSuccessiveFailureCount                      *int            `json:"snapshotSuccessiveFailureCount,omitempty"`
	FirstSnapshotRetryCount                             *int            `json:"firstSnapshotRetryCount,omitempty"`
	LogCatchupSuccessiveFailureCount                    *int            `json:"logCatchupSuccessiveFailureCount,omitempty"`
	ImplicitResumeCount                                 *int            `json:"implicitResumeCount,omitempty"`
	RequiredSpace                                       *float64        `json:"requiredSpace,omitempty"`
	CapabilityResetTime                                 *string         `json:"capabilityResetTime,omitempty"`
	LastSnapshotTime                                    *string         `json:"lastSnapshotTime,omitempty"`
	LastAutoSnapshotTime                                *string         `json:"lastAutoSnapshotTime,omitempty"`
	LastSnapshotOperationID                             *string         `json:"lastSnapshotOperationId,omitempty"`
	LastAutoSnapshotOperationID                         *string         `json:"lastAutoSnapshotOperationId,omitempty"`
	LastSuccessfulSnapshotOperationID                   *string         `json:"lastSuccessfulSnapshotOperationId,omitempty"`
	LastHealSnapshotOperation                           *string         `json:"lastHealSnapshotOperation,omitempty"`
	LastNonExtraAutoSnapshotTime                        *string         `json:"lastNonExtraAutoSnapshotTime,omitempty"`
	LastLogCatchupTime                                  *string         `json:"lastLogCatchupTime,omitempty"`
	LastSuccessfulLogCatchupOperationID                 *string         `json:"lastSuccessfulLogCatchupOperationId,omitempty"`
	LastLogCatchupOperationID                           *string         `json:"lastLogCatchupOperationId,omitempty"`
	LastPauseTime                                       *string         `json:"lastPauseTime,omitempty"`
	LastResumeTime                                      *string         `json:"lastResumeTime,omitempty"`
	LastPauseReason                                     *string         `json:"lastPauseReason,omitempty"`
	StateBeforeRestore                                  *string         `json:"stateBeforeRestore,omitempty"`
	LastHealthAlertedTime                               *string         `json:"lastHealthAlertedTime,omitempty"`
	LastImplicitResumeTime                              *string         `json:"lastImplicitResumeTime,omitempty"`
	LastEaBreakdownTime                                 *string         `json:"lastEaBreakdownTime,omitempty"`
	LastHealTime                                        *string         `json:"lastHealTime,omitempty"`
	AuthorizedDbservers                                 []string        `json:"authorizedDbservers,omitempty"`
	DeregisterInfo                                      *DeregisterInfo `json:"deregisterInfo,omitempty"`
	SecureInfo                                          interface{}     `json:"secureInfo,omitempty"`
	Info                                                interface{}     `json:"info,omitempty"`
	DatabasesFirstSnapshotInfo                          interface{}     `json:"databasesFirstSnapshotInfo,omitempty"`
	OnboardingSnapshotProperties                        interface{}     `json:"onboardingSnapshotProperties,omitempty"`
	LastSuccessfulLogCatchupPostHealWithResetCapability interface{}     `json:"lastSuccessfulLogCatchupPostHealWithResetCapability,omitempty"`
	AutoSnapshotRetryInfo                               interface{}     `json:"autoSnapshotRetryInfo,omitempty"`
}

type SLAIntentInput struct {
	Name                *string `json:"name,omitempty"`
	Description         *string `json:"description,omitempty"`
	ContinuousRetention *int    `json:"continuousRetention,omitempty"`
	DailyRetention      *int    `json:"dailyRetention,omitempty"`
	WeeklyRetention     *int    `json:"weeklyRetention,omitempty"`
	MonthlyRetention    *int    `json:"monthlyRetention,omitempty"`
	QuarterlyRetention  *int    `json:"quarterlyRetention,omitempty"`
	ID                  *string `json:"id,omitempty"`
}

type SLADeleteResponse struct {
	Status *string `json:"status,omitempty"`
}

type DatabaseRestoreRequest struct {
	SnapshotID        *string            `json:"snapshotId,omitempty"`
	LatestSnapshot    *string            `json:"latestSnapshot,omitempty"`
	UserPitrTimestamp *string            `json:"userPitrTimestamp,omitempty"`
	TimeZone          *string            `json:"timeZone,omitempty"`
	ActionArguments   []*Actionarguments `json:"actionArguments,omitempty"`
}

type LogCatchUpRequest struct {
	ForRestore      bool               `json:"for_restore,omitempty"`
	Actionarguments []*Actionarguments `json:"actionArguments,omitempty"`
}

type DatabaseScale struct {
	ApplicationType *string            `json:"applicationType,omitempty"`
	Actionarguments []*Actionarguments `json:"actionArguments,omitempty"`
}

type ProfileProperties struct {
	Name        *string `json:"name,omitempty"`
	Value       *string `json:"value,omitempty"`
	Secure      bool    `json:"secure"`
	Description *string `json:"description,omitempty"`
}

type ProfileRequest struct {
	EngineType                *string                      `json:"engineType,omitempty"`
	Type                      *string                      `json:"type,omitempty"`
	Topology                  *string                      `json:"topology,omitempty"`
	DBVersion                 *string                      `json:"dbVersion,omitempty"`
	Name                      *string                      `json:"name,omitempty"`
	Description               *string                      `json:"description,omitempty"`
	AvailableClusterIds       []*string                    `json:"availableClusterIds,omitempty"`
	SystemProfile             bool                         `json:"systemProfile,omitempty"`
	Published                 bool                         `json:"published"`
	Deprecated                bool                         `json:"deprecated"`
	Properties                []*ProfileProperties         `json:"properties,omitempty"`
	VersionClusterAssociation []*VersionClusterAssociation `json:"versionClusterAssociation,omitempty"`
}

type SoftwareProfileResponse struct {
	Name        *string `json:"name,omitempty"`
	WorkID      *string `json:"workId,omitempty"`
	OperationID *string `json:"operationId,omitempty"`
	DbserverID  *string `json:"dbserverId,omitempty"`
	EntityID    *string `json:"entityId,omitempty"`
	EntityName  *string `json:"entityName,omitempty"`
	EntityType  *string `json:"entityType,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type UpdateProfileRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ProfileFilter struct {
	Engine      string `json:"engine,omitempty"`
	ProfileType string `json:"profile_type,omitempty"`
	ProfileID   string `json:"profile_id,omitempty"`
	ProfileName string `json:"profile_name,omitempty"`
}

type RegisterDBInputRequest struct {
	NxClusterID                 *string            `json:"nxClusterId,omitempty"`
	DatabaseType                *string            `json:"databaseType,omitempty"`
	DatabaseName                *string            `json:"databaseName,omitempty"`
	Description                 *string            `json:"description,omitempty"`
	Category                    *string            `json:"category,omitempty"`
	VMIP                        *string            `json:"vmIp,omitempty"`
	WorkingDirectory            *string            `json:"workingDirectory,omitempty"`
	VMUsername                  *string            `json:"vmUsername,omitempty"`
	VMPassword                  *string            `json:"vmPassword,omitempty"`
	VMSshkey                    *string            `json:"vmSshkey,omitempty"`
	VMDescription               *string            `json:"vmDescription,omitempty"`
	ResetDescriptionInNxCluster bool               `json:"resetDescriptionInNxCluster,omitempty"`
	AutoTuneStagingDrive        bool               `json:"autoTuneStagingDrive,omitempty"`
	Clustered                   bool               `json:"clustered,omitempty"`
	ForcedInstall               bool               `json:"forcedInstall,omitempty"`
	TimeMachineInfo             *Timemachineinfo   `json:"timeMachineInfo,omitempty"`
	Tags                        []*Tags            `json:"tags,omitempty"`
	Actionarguments             []*Actionarguments `json:"actionArguments,omitempty"`
	MaintenanceTasks            *MaintenanceTasks  `json:"maintenanceTasks,omitempty"`
}

type UnRegisterDatabaseRequest struct {
	SoftRemove        bool `json:"softRemove,omitempty"`
	Remove            bool `json:"remove,omitempty"`
	Delete            bool `json:"delete,omitempty"`
	DeleteTimeMachine bool `json:"deleteTimeMachine,omitempty"`
}
type DatabaseSnapshotRequest struct {
	Name                *string            `json:"name,omitempty"`
	LcmConfig           *LCMConfigSnapshot `json:"lcmConfig,omitempty"`
	ReplicateToClusters []*string          `json:"replicateToClusterIds,omitempty"`
}

type LCMConfigSnapshot struct {
	SnapshotLCMConfig *SnapshotLCMConfig `json:"snapshotLCMConfig,omitempty"`
}

type SnapshotLCMConfig struct {
	ExpiryDetails *DBExpiryDetails `json:"expiryDetails,omitempty"`
}

type ListTimeMachines []*TimeMachine

type CloneLCMConfig struct {
	DatabaseLCMConfig *DatabaseLCMConfig `json:"databaseLCMConfig,omitempty"`
}

type DatabaseLCMConfig struct {
	ExpiryDetails  *DBExpiryDetails  `json:"expiryDetails,omitempty"`
	RefreshDetails *DBRefreshDetails `json:"refreshDetails,omitempty"`
}

type CloneRequest struct {
	Name                       *string            `json:"name,omitempty"`
	Description                *string            `json:"description,omitempty"`
	NxClusterID                *string            `json:"nxClusterId,omitempty"`
	SSHPublicKey               *string            `json:"sshPublicKey,omitempty"`
	DbserverID                 *string            `json:"dbserverId,omitempty"`
	DbserverClusterID          *string            `json:"dbserverClusterId,omitempty"`
	DbserverLogicalClusterID   *string            `json:"dbserverLogicalClusterId,omitempty"`
	TimeMachineID              *string            `json:"timeMachineId,omitempty"`
	SnapshotID                 *string            `json:"snapshotId,omitempty"`
	UserPitrTimestamp          *string            `json:"userPitrTimestamp,omitempty"`
	TimeZone                   *string            `json:"timeZone,omitempty"`
	VMPassword                 *string            `json:"vmPassword,omitempty"`
	ComputeProfileID           *string            `json:"computeProfileId,omitempty"`
	NetworkProfileID           *string            `json:"networkProfileId,omitempty"`
	DatabaseParameterProfileID *string            `json:"databaseParameterProfileId,omitempty"`
	NodeCount                  *int               `json:"nodeCount,omitempty"`
	Nodes                      []*Nodes           `json:"nodes,omitempty"`
	ActionArguments            []*Actionarguments `json:"actionArguments,omitempty"`
	Tags                       []*Tags            `json:"tags,omitempty"`
	LatestSnapshot             bool               `json:"latestSnapshot,omitempty"`
	CreateDbserver             bool               `json:"createDbserver,omitempty"`
	Clustered                  bool               `json:"clustered,omitempty"`
	LcmConfig                  *CloneLCMConfig    `json:"lcmConfig,omitempty"`
}

type AuthorizeDBServerResponse struct {
	ErrorCode *int    `json:"errorCode,omitempty"`
	Info      *string `json:"info,omitempty"`
	Message   *string `json:"message,omitempty"`
	Status    *string `json:"status,omitempty"`
}

type FilterParams struct {
	Detailed                      string `json:"detailed,omitempty"`
	AnyStatus                     string `json:"any-status,omitempty"`
	LoadDBServerCluster           string `json:"load-dbserver-cluster"`
	TimeZone                      string `json:"time-zone,omitempty"`
	OrderByDBServerCluster        string `json:"order-by-dbserver-cluster,omitempty"`
	OrderByDBServerLogicalCluster string `json:"order-by-dbserver-logical-cluster,omitempty"`
	LoadReplicatedChildSnapshots  string `json:"load-replicated-child-snapshots,omitempty"`
}

type UpdateSnapshotRequest struct {
	Name      *string `json:"name,omitempty"`
	ResetName bool    `json:"resetName,omitempty"`
}

type ListSnapshots []SnapshotResponse

type SnapshotResponse struct {
	ID                             *string                 `json:"id,omitempty"`
	Name                           *string                 `json:"name,omitempty"`
	Description                    *string                 `json:"description,omitempty"`
	OwnerID                        *string                 `json:"ownerId,omitempty"`
	DateCreated                    *string                 `json:"dateCreated,omitempty"`
	DateModified                   *string                 `json:"dateModified,omitempty"`
	SnapshotID                     *string                 `json:"snapshotId,omitempty"`
	SnapshotUUID                   *string                 `json:"snapshotUuid,omitempty"`
	NxClusterID                    *string                 `json:"nxClusterId,omitempty"`
	ProtectionDomainID             *string                 `json:"protectionDomainId,omitempty"`
	ParentSnapshotID               *string                 `json:"parentSnapshotId,omitempty"`
	TimeMachineID                  *string                 `json:"timeMachineId,omitempty"`
	DatabaseNodeID                 *string                 `json:"databaseNodeId,omitempty"`
	AppInfoVersion                 *string                 `json:"appInfoVersion,omitempty"`
	Status                         *string                 `json:"status,omitempty"`
	Type                           *string                 `json:"type,omitempty"`
	SnapshotTimeStamp              *string                 `json:"snapshotTimeStamp,omitempty"`
	TimeZone                       *string                 `json:"timeZone,omitempty"`
	SoftwareSnapshotID             *string                 `json:"softwareSnapshotId,omitempty"`
	FromTimeStamp                  *string                 `json:"fromTimeStamp,omitempty"`
	ToTimeStamp                    *string                 `json:"toTimeStamp,omitempty"`
	ApplicableTypes                []*string               `json:"applicableTypes,omitempty"`
	DBServerStorageMetadataVersion *int                    `json:"dbServerStorageMetadataVersion,omitempty"`
	SnapshotTimeStampDate          *int64                  `json:"snapshotTimeStampDate,omitempty"`
	SnapshotSize                   *float64                `json:"snapshotSize,omitempty"`
	ParentSnapshot                 *bool                   `json:"parentSnapshot,omitempty"`
	SoftwareDatabaseSnapshot       bool                    `json:"softwareDatabaseSnapshot,omitempty"`
	Processed                      bool                    `json:"processed,omitempty"`
	DatabaseSnapshot               bool                    `json:"databaseSnapshot,omitempty"`
	Sanitized                      bool                    `json:"sanitised,omitempty"` //nolint:all
	Properties                     []*DBInstanceProperties `json:"properties"`
	Tags                           []*Tags                 `json:"tags"`
	Info                           *CloneInfo              `json:"info,omitempty"`
	Metadata                       *ClonedMetadata         `json:"metadata,omitempty"`
	Metric                         *Metric                 `json:"metric,omitempty"`
	LcmConfig                      *LcmConfig              `json:"lcmConfig,omitempty"`
	SanitizedFromSnapshotID        interface{}             `json:"sanitisedFromSnapshotId,omitempty"`
	AccessLevel                    interface{}             `json:"accessLevel"`
	DbserverID                     interface{}             `json:"dbserverId,omitempty"`
	DbserverName                   interface{}             `json:"dbserverName,omitempty"`
	DbserverIP                     interface{}             `json:"dbserverIp,omitempty"`
	ReplicatedSnapshots            interface{}             `json:"replicatedSnapshots,omitempty"`
	SoftwareSnapshot               interface{}             `json:"softwareSnapshot,omitempty"`
	SanitizedSnapshots             interface{}             `json:"sanitisedSnapshots,omitempty"`
	SnapshotFamily                 interface{}             `json:"snapshotFamily,omitempty"`
}

type LinkedDBInfo struct {
	Info *Info `json:"info,omitempty"`
}

type CloneLinkedDBInfo struct {
	ID            *string       `json:"id,omitempty"`
	DatabaseName  *string       `json:"databaseName,omitempty"`
	Status        *string       `json:"status,omitempty"`
	Info          *LinkedDBInfo `json:"info,omitempty"`
	AppConsistent bool          `json:"appConsistent,omitempty"`
	Clone         bool          `json:"clone,omitempty"`
	Message       interface{}   `json:"message,omitempty"`
}

type CloneInfo struct {
	SecureInfo         interface{}          `json:"secureInfo,omitempty"`
	Info               interface{}          `json:"info,omitempty"`
	LinkedDatabases    []*CloneLinkedDBInfo `json:"linkedDatabases,omitempty"`
	Databases          interface{}          `json:"databases,omitempty"`
	DatabaseGroupID    interface{}          `json:"databaseGroupId,omitempty"`
	MissingDatabases   interface{}          `json:"missingDatabases,omitempty"`
	ReplicationHistory interface{}          `json:"replicationHistory,omitempty"`
}

type ClonedMetadata struct {
	SecureInfo                           interface{}   `json:"secureInfo,omitempty"`
	Info                                 interface{}   `json:"info,omitempty"`
	DeregisterInfo                       interface{}   `json:"deregisterInfo,omitempty"`
	FromTimeStamp                        string        `json:"fromTimeStamp,omitempty"`
	ToTimeStamp                          string        `json:"toTimeStamp,omitempty"`
	ReplicationRetryCount                int           `json:"replicationRetryCount,omitempty"`
	LastReplicationRetryTimestamp        interface{}   `json:"lastReplicationRetryTimestamp,omitempty"`
	LastReplicationRetrySourceSnapshotID interface{}   `json:"lastReplicationRetrySourceSnapshotId,omitempty"`
	Async                                bool          `json:"async,omitempty"`
	Standby                              bool          `json:"standby,omitempty"`
	CurationRetryCount                   int           `json:"curationRetryCount,omitempty"`
	OperationsUsingSnapshot              []interface{} `json:"operationsUsingSnapshot,omitempty"`
}

type Capability struct {
	Mode                      *string           `json:"mode,omitempty"`
	From                      *string           `json:"from,omitempty"`
	To                        *string           `json:"to,omitempty"`
	TimeUnit                  *string           `json:"timeUnit,omitempty"`
	TimeUnitNumber            *string           `json:"timeUnitNumber,omitempty"`
	DatabaseIds               []*string         `json:"databaseIds,omitempty"`
	Snapshots                 *ListSnapshots    `json:"snapshots,omitempty"`
	ContinuousRegion          *ContinuousRegion `json:"continuousRegion,omitempty"`
	DatabasesContinuousRegion interface{}       `json:"databasesContinuousRegion,omitempty"`
}

type TimeMachineCapability struct {
	TimeMachineID                 *string                 `json:"timeMachineId,omitempty"`
	OutputTimeZone                *string                 `json:"outputTimeZone,omitempty"`
	Type                          *string                 `json:"type,omitempty"`
	NxClusterID                   *string                 `json:"nxClusterId,omitempty"`
	NxClusterAssociationType      *string                 `json:"nxClusterAssociationType,omitempty"`
	SLAID                         *string                 `json:"slaId,omitempty"`
	CapabilityResetTime           *string                 `json:"capabilityResetTime,omitempty"`
	LastContinuousSnapshotTime    *string                 `json:"lastContinuousSnapshotTime,omitempty"`
	LogCatchupStartTime           *string                 `json:"logCatchupStartTime,omitempty"`
	DatabaseIds                   []*string               `json:"databaseIds,omitempty"`
	HealWithResetCapability       bool                    `json:"healWithResetCapability,omitempty"`
	Source                        bool                    `json:"source,omitempty"`
	Capability                    []*Capability           `json:"capability,omitempty"`
	LogTimeInfo                   map[string]interface{}  `json:"logTimeInfo,omitempty"`
	LastDBLog                     *DBLogs                 `json:"lastDbLog,omitempty"`
	LastContinuousSnapshot        *LastContinuousSnapshot `json:"lastContinuousSnapshot,omitempty"`
	OverallContinuousRangeEndTime interface{}             `json:"overallContinuousRangeEndTime,omitempty"`
}

type ProcessedRanges struct {
	First  string `json:"first,omitempty"`
	Second string `json:"second,omitempty"`
}

type DBLogsInfo struct {
	SecureInfo       interface{} `json:"secureInfo,omitempty"`
	Info             interface{} `json:"info,omitempty"`
	UnknownTimeRange bool        `json:"unknownTimeRange,omitempty"`
}

type DBLogsMetadata struct {
	SecureInfo         interface{}     `json:"secureInfo,omitempty"`
	Info               interface{}     `json:"info,omitempty"`
	DeregisterInfo     *DeregisterInfo `json:"deregisterInfo,omitempty"`
	CurationRetryCount int             `json:"curationRetryCount,omitempty"`
	CreatedDirectly    bool            `json:"createdDirectly,omitempty"`
	UpdatedDirectly    bool            `json:"updatedDirectly,omitempty"`
}

type DBLogs struct {
	ID                 string          `json:"id,omitempty"`
	Name               string          `json:"name,omitempty"`
	EraLogDriveID      string          `json:"eraLogDriveId,omitempty"`
	DatabaseNodeID     string          `json:"databaseNodeId,omitempty"`
	FromTime           string          `json:"fromTime,omitempty"`
	ToTime             string          `json:"toTime,omitempty"`
	Status             string          `json:"status,omitempty"`
	Size               int             `json:"size,omitempty"`
	Info               *DBLogsInfo     `json:"info,omitempty"`
	Metadata           *DBLogsMetadata `json:"metadata,omitempty"`
	DateCreated        string          `json:"dateCreated,omitempty"`
	DateModified       string          `json:"dateModified,omitempty"`
	OwnerID            string          `json:"ownerId,omitempty"`
	DatabaseID         interface{}     `json:"databaseId,omitempty"`
	Message            interface{}     `json:"message,omitempty"`
	Unprocessed        bool            `json:"unprocessed,omitempty"`
	LogCopyOperationID interface{}     `json:"logCopyOperationId,omitempty"`
}

type ContinuousRegion struct {
	FromTime              string             `json:"fromTime,omitempty"`
	ToTime                string             `json:"toTime,omitempty"`
	TimeZone              string             `json:"timeZone,omitempty"`
	SnapshotIds           []string           `json:"snapshotIds,omitempty"`
	PartialRanges         bool               `json:"partialRanges,omitempty"`
	SubRange              bool               `json:"subRange,omitempty"`
	Message               interface{}        `json:"message,omitempty"`
	UnknownTimeRanges     interface{}        `json:"unknownTimeRanges,omitempty"`
	TimeRangeAndDatabases interface{}        `json:"timeRangeAndDatabases,omitempty"`
	Snapshots             interface{}        `json:"snapshots,omitempty"`
	DBLogs                []*DBLogs          `json:"dbLogs,omitempty"`
	ProcessedRanges       []*ProcessedRanges `json:"processedRanges,omitempty"`
	UnprocessedRanges     []*ProcessedRanges `json:"unprocessedRanges,omitempty"`
}

type LastContinuousSnapshotMetadata struct {
	FromTimeStamp                        string        `json:"fromTimeStamp,omitempty"`
	ToTimeStamp                          string        `json:"toTimeStamp,omitempty"`
	ReplicationRetryCount                int           `json:"replicationRetryCount,omitempty"`
	CurationRetryCount                   int           `json:"curationRetryCount,omitempty"`
	Async                                bool          `json:"async,omitempty"`
	Standby                              bool          `json:"standby,omitempty"`
	SecureInfo                           interface{}   `json:"secureInfo,omitempty"`
	Info                                 interface{}   `json:"info,omitempty"`
	DeregisterInfo                       interface{}   `json:"deregisterInfo,omitempty"`
	LastReplicationRetryTimestamp        interface{}   `json:"lastReplicationRetryTimestamp,omitempty"`
	LastReplicationRetrySourceSnapshotID interface{}   `json:"lastReplicationRetrySourceSnapshotId,omitempty"`
	OperationsUsingSnapshot              []interface{} `json:"operationsUsingSnapshot,omitempty"`
}

type LastContinuousSnapshot struct {
	ID                             string                          `json:"id,omitempty"`
	Name                           string                          `json:"name,omitempty"`
	OwnerID                        string                          `json:"ownerId,omitempty"`
	DateCreated                    string                          `json:"dateCreated,omitempty"`
	DateModified                   string                          `json:"dateModified,omitempty"`
	SnapshotID                     string                          `json:"snapshotId,omitempty"`
	SnapshotUUID                   string                          `json:"snapshotUuid,omitempty"`
	NxClusterID                    string                          `json:"nxClusterId,omitempty"`
	ProtectionDomainID             string                          `json:"protectionDomainId,omitempty"`
	TimeMachineID                  string                          `json:"timeMachineId,omitempty"`
	DatabaseNodeID                 string                          `json:"databaseNodeId,omitempty"`
	AppInfoVersion                 string                          `json:"appInfoVersion,omitempty"`
	Status                         string                          `json:"status,omitempty"`
	Type                           string                          `json:"type,omitempty"`
	SnapshotTimeStamp              string                          `json:"snapshotTimeStamp,omitempty"`
	SoftwareSnapshotID             string                          `json:"softwareSnapshotId,omitempty"`
	TimeZone                       string                          `json:"timeZone,omitempty"`
	FromTimeStamp                  string                          `json:"fromTimeStamp,omitempty"`
	ToTimeStamp                    string                          `json:"toTimeStamp,omitempty"`
	ApplicableTypes                []string                        `json:"applicableTypes,omitempty"`
	SoftwareDatabaseSnapshot       bool                            `json:"softwareDatabaseSnapshot,omitempty"`
	Processed                      bool                            `json:"processed,omitempty"`
	DatabaseSnapshot               bool                            `json:"databaseSnapshot,omitempty"`
	ParentSnapshot                 bool                            `json:"parentSnapshot,omitempty"`
	DBServerStorageMetadataVersion int                             `json:"dbServerStorageMetadataVersion,omitempty"`
	SnapshotTimeStampDate          int64                           `json:"snapshotTimeStampDate,omitempty"`
	SnapshotSize                   float64                         `json:"snapshotSize,omitempty"`
	AccessLevel                    interface{}                     `json:"accessLevel,omitempty"`
	Metric                         interface{}                     `json:"metric,omitempty"`
	SanitizedFromSnapshotID        interface{}                     `json:"sanitisedFromSnapshotId,omitempty"`
	DBserverID                     interface{}                     `json:"dbserverId,omitempty"`
	DBserverName                   interface{}                     `json:"dbserverName,omitempty"`
	DBserverIP                     interface{}                     `json:"dbserverIp,omitempty"`
	ReplicatedSnapshots            interface{}                     `json:"replicatedSnapshots,omitempty"`
	SoftwareSnapshot               interface{}                     `json:"softwareSnapshot,omitempty"`
	SanitizedSnapshots             interface{}                     `json:"sanitisedSnapshots,omitempty"`
	Description                    interface{}                     `json:"description,omitempty"`
	SnapshotFamily                 interface{}                     `json:"snapshotFamily,omitempty"`
	ParentSnapshotID               interface{}                     `json:"parentSnapshotId,omitempty"`
	Properties                     []*DBInstanceProperties         `json:"properties,omitempty"`
	Tags                           []*Tags                         `json:"tags,omitempty"`
	Info                           *CloneInfo                      `json:"info,omitempty"`
	Metadata                       *LastContinuousSnapshotMetadata `json:"metadata,omitempty"`
	LcmConfig                      *LcmConfig                      `json:"lcmConfig,omitempty"`
}

type LinkedDatabases struct {
	DatabaseName *string `json:"databaseName,omitempty"`
}

type CreateLinkedDatabasesRequest struct {
	Databases []*LinkedDatabases `json:"databases,omitempty"`
}

type DeleteLinkedDatabaseRequest struct {
	Delete bool `json:"delete,omitempty"`
	Forced bool `json:"forced,omitempty"`
}

type MaintenaceSchedule struct {
	Recurrence  *string     `json:"recurrence,omitempty"`
	StartTime   *string     `json:"startTime,omitempty"`
	DayOfWeek   *string     `json:"dayOfWeek,omitempty"`
	WeekOfMonth *int        `json:"weekOfMonth,omitempty"`
	Duration    *int        `json:"duration,omitempty"`
	Threshold   interface{} `json:"threshold,omitempty"`
	Hour        *int        `json:"hour,omitempty"`
	Minute      *int        `json:"minute,omitempty"`
	TimeZone    *string     `json:"timeZone,omitempty"`
}

type MaintenanceWindowInput struct {
	Name             *string             `json:"name,omitempty"`
	Description      *string             `json:"description,omitempty"`
	Timezone         *string             `json:"timezone,omitempty"`
	Schedule         *MaintenaceSchedule `json:"schedule,omitempty"`
	ResetSchedule    *bool               `json:"resetSchedule,omitempty"`
	ResetDescription *bool               `json:"resetDescription,omitempty"`
	ResetName        *bool               `json:"resetName,omitempty"`
}

type MaintenaceWindowResponse struct {
	ID              *string                     `json:"id,omitempty"`
	Name            *string                     `json:"name,omitempty"`
	Description     *string                     `json:"description,omitempty"`
	OwnerID         *string                     `json:"ownerId,omitempty"`
	DateCreated     *string                     `json:"dateCreated,omitempty"`
	DateModified    *string                     `json:"dateModified,omitempty"`
	AccessLevel     interface{}                 `json:"accessLevel,omitempty"`
	Properties      []*Properties               `json:"properties,omitempty"`
	Tags            []*Tags                     `json:"tags,omitempty"`
	Schedule        *MaintenaceSchedule         `json:"schedule,omitempty"`
	Status          *string                     `json:"status,omitempty"`
	NextRunTime     *string                     `json:"nextRunTime,omitempty"`
	EntityTaskAssoc []*MaintenanceTasksResponse `json:"entityTaskAssoc,omitempty"`
	Timezone        *string                     `json:"timezone,omitempty"`
}

type ListMaintenanceWindowResponse []MaintenaceWindowResponse
type MaintenanceEntities struct {
	EraDBServer        []*string `json:"ERA_DBSERVER,omitempty"`
	EraDBServerCluster []*string `json:"ERA_DBSERVER_CLUSTER,omitempty"`
}

type MaintenanceTasksInput struct {
	Entities            *MaintenanceEntities `json:"entities,omitempty"`
	MaintenanceWindowID *string              `json:"maintenanceWindowId,omitempty"`
	Tasks               []*Tasks             `json:"tasks"`
}

type MaintenanceTasksResponse struct {
	ID                       *string       `json:"id,omitempty"`
	Name                     *string       `json:"name,omitempty"`
	Description              *string       `json:"description,omitempty"`
	OwnerID                  *string       `json:"ownerId,omitempty"`
	DateCreated              *string       `json:"dateCreated,omitempty"`
	DateModified             *string       `json:"dateModified,omitempty"`
	AccessLevel              *string       `json:"accessLevel,omitempty"`
	Properties               []*Properties `json:"properties,omitempty"`
	Tags                     []*Tags       `json:"tags,omitempty"`
	MaintenanceWindowID      *string       `json:"maintenanceWindowId,omitempty"`
	MaintenanceWindowOwnerID *string       `json:"maintenanceWindowOwnerId,omitempty"`
	EntityID                 *string       `json:"entityId,omitempty"`
	EntityType               *string       `json:"entityType,omitempty"`
	Status                   *string       `json:"status,omitempty"`
	TaskType                 *string       `json:"taskType,omitempty"`
	Payload                  *Payload      `json:"payload,omitempty"`
	Entity                   interface{}   `json:"entity,omitempty"`
}

type ListMaintenanceTasksResponse []MaintenanceTasksResponse
type TmsClusterIntentInput struct {
	NxClusterID *string `json:"nxClusterId,omitempty"`
	Type        *string `json:"type,omitempty"`
	SLAID       *string `json:"slaId,omitempty"`
	ResetSLAID  *bool   `json:"resetSlaId,omitempty"`
}

type TmsClusterResponse struct {
	TimeMachineID               *string     `json:"timeMachineId,omitempty"`
	NxClusterID                 *string     `json:"nxClusterId,omitempty"`
	LogDriveID                  *string     `json:"logDriveId,omitempty"`
	LogDriveStatus              *string     `json:"logDriveStatus,omitempty"`
	Type                        *string     `json:"type,omitempty"`
	Description                 *string     `json:"description,omitempty"`
	Status                      *string     `json:"status,omitempty"`
	SLAID                       *string     `json:"slaId,omitempty"`
	ScheduleID                  *string     `json:"scheduleId,omitempty"`
	OwnerID                     *string     `json:"ownerId,omitempty"`
	DateCreated                 *string     `json:"dateCreated,omitempty"`
	DateModified                *string     `json:"dateModified,omitempty"`
	Info                        interface{} `json:"info,omitempty"`
	Metadata                    interface{} `json:"metadata,omitempty"`
	NxCluster                   interface{} `json:"nxCluster,omitempty"`
	LogDrive                    interface{} `json:"logDrive,omitempty"`
	SLA                         interface{} `json:"sla,omitempty"`
	Schedule                    interface{} `json:"schedule,omitempty"`
	SourceClusters              []*string   `json:"sourceClusters,omitempty"`
	ResetSLAID                  *bool       `json:"resetSlaId,omitempty"`
	ResetDescription            *bool       `json:"resetDescription,omitempty"`
	ResetType                   *bool       `json:"resetType,omitempty"`
	SubmitActivateTimeMachineOp *bool       `json:"submitActivateTimeMachineOp,omitempty"`
	UpdateOperationSummary      interface{} `json:"updateOperationSummary,omitempty"`
	StorageResourceID           *string     `json:"storageResourceId,omitempty"`
	ForceVGBasedLogDrive        *bool       `json:"forceVGBasedLogDrive,omitempty"`
	Source                      *bool       `json:"source,omitempty"`
}

type DeleteTmsClusterInput struct {
	DeleteReplicatedSnapshots         *bool `json:"deleteReplicatedSnapshots,omitempty"`
	DeleteReplicatedProtectionDomains *bool `json:"deleteReplicatedProtectionDomains,omitempty"`
}
type CreateTagsInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	EntityType  *string `json:"entityType,omitempty"`
	Required    *bool   `json:"required,omitempty"`
}

type TagsIntentResponse struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Owner       *string `json:"owner,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Status      *string `json:"status,omitempty"`
	EntityType  *string `json:"entityType,omitempty"`
	Values      int     `json:"values,omitempty"`
}

type GetTagsResponse struct {
	ID           *string `json:"id,omitempty"`
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	DateCreated  *string `json:"dateCreated,omitempty"`
	DateModified *string `json:"dateModified,omitempty"`
	Owner        *string `json:"owner,omitempty"`
	Status       *string `json:"status,omitempty"`
	EntityType   *string `json:"entityType,omitempty"`
	Required     *bool   `json:"required,omitempty"`
	Values       *int    `json:"values,omitempty"`
}

type ListTagsResponse []*GetTagsResponse
type IPAddresses struct {
	IP           *string `json:"ip,omitempty"`
	Status       *string `json:"status,omitempty"`
	DBServerID   *string `json:"dbserverID,omitempty"`
	DBServerName *string `json:"dbserverName,omitempty"`
}

type IPPools struct {
	StartIP     *string        `json:"startIP,omitempty"`
	EndIP       *string        `json:"endIP,omitempty"`
	ID          *string        `json:"id,omitempty"`
	ModifiedBy  *string        `json:"modifiedBy,omitempty"`
	IPAddresses []*IPAddresses `json:"ipAddresses,omitempty"`
}

type NetworkIntentInput struct {
	Name       *string       `json:"name,omitempty"`
	Type       *string       `json:"type,omitempty"`
	Properties []*Properties `json:"properties,omitempty"`
	ClusterID  *string       `json:"clusterId,omitempty"`
	IPPools    []*IPPools    `json:"ipPools,omitempty"`
}

type NetworkPropertiesmap struct {
	VLANSecondaryDNS *string `json:"VLAN_SECONDARY_DNS,omitempty"`
	VLANSubnetMask   *string `json:"VLAN_SUBNET_MASK,omitempty"`
	VLANPrimaryDNS   *string `json:"VLAN_PRIMARY_DNS,omitempty"`
	VLANGateway      *string `json:"VLAN_GATEWAY,omitempty"`
}

type NetworkIntentResponse struct {
	ID              *string               `json:"id,omitempty"`
	Name            *string               `json:"name,omitempty"`
	Type            *string               `json:"type,omitempty"`
	ClusterID       *string               `json:"clusterId,omitempty"`
	Managed         *bool                 `json:"managed,omitempty"`
	Properties      []*Properties         `json:"properties,omitempty"`
	PropertiesMap   *NetworkPropertiesmap `json:"propertiesMap,omitempty"`
	StretchedVlanID *string               `json:"stretchedVlanId,omitempty"`
	IPPools         []*IPPools            `json:"ipPools,omitempty"`
	IPAddresses     []*IPAddresses        `json:"ipAddresses,omitempty"`
}

type ListNetworkResponse []*NetworkIntentResponse

type DBServerInputRequest struct {
	DatabaseType             *string            `json:"databaseType,omitempty"`
	SoftwareProfileID        *string            `json:"softwareProfileId,omitempty"`
	SoftwareProfileVersionID *string            `json:"softwareProfileVersionId,omitempty"`
	NetworkProfileID         *string            `json:"networkProfileId,omitempty"`
	ComputeProfileID         *string            `json:"computeProfileId,omitempty"`
	VMPassword               *string            `json:"vmPassword,omitempty"`
	NxClusterID              *string            `json:"nxClusterId,omitempty"`
	LatestSnapshot           bool               `json:"latestSnapshot,omitempty"`
	ActionArguments          []*Actionarguments `json:"actionArguments,omitempty"`
	Description              *string            `json:"description,omitempty"`
	TimeMachineID            *string            `json:"timeMachineId,omitempty"`
	SnapshotID               *string            `json:"snapshotId,omitempty"`
	TimeZone                 *string            `json:"timeZone,omitempty"`
	MaintenanceTasks         *MaintenanceTasks  `json:"maintenanceTasks,omitempty"`
}

type DeleteDBServerVMRequest struct {
	SoftRemove        bool `json:"softRemove,omitempty"`
	Remove            bool `json:"remove,omitempty"`
	Delete            bool `json:"delete,omitempty"`
	DeleteVgs         bool `json:"deleteVgs,omitempty"`
	DeleteVMSnapshots bool `json:"deleteVmSnapshots,omitempty"`
}

type VMCredentials struct {
	Username *string     `json:"username,omitempty"`
	Password *string     `json:"password,omitempty"`
	Label    interface{} `json:"label,omitempty"`
}

type UpdateDBServerVMRequest struct {
	Name                        *string          `json:"name,omitempty"`
	Description                 *string          `json:"description,omitempty"`
	ResetNameInNxCluster        *bool            `json:"resetNameInNxCluster,omitempty"`
	ResetDescriptionInNxCluster *bool            `json:"resetDescriptionInNxCluster,omitempty"`
	ResetCredential             *bool            `json:"resetCredential,omitempty"`
	ResetTags                   *bool            `json:"resetTags,omitempty"`
	ResetName                   *bool            `json:"resetName,omitempty"`
	ResetDescription            *bool            `json:"resetDescription,omitempty"`
	Tags                        []*Tags          `json:"tags,omitempty"`
	Credentials                 []*VMCredentials `json:"credentials,omitempty"`
}

type DiskList struct {
	DeviceName    *string `json:"device_name,omitempty"`
	LocalMapping  *string `json:"local_mapping,omitempty"`
	DiskIndex     *string `json:"disk_index,omitempty"`
	Path          *string `json:"path,omitempty"`
	DiskID        *string `json:"disk_id,omitempty"`
	Hypervisor    *string `json:"hypervisor,omitempty"`
	MountPoint    *string `json:"mount_point,omitempty"`
	DiskUUID      *string `json:"disk_uuid,omitempty"`
	DiskUser      *string `json:"disk_user,omitempty"`
	DiskGroup     *string `json:"disk_group,omitempty"`
	DiskMode      *string `json:"disk_mode,omitempty"`
	DiskType      *string `json:"disk_type,omitempty"`
	VirtualDiskID *string `json:"virtual_disk_id,omitempty"`
	FsType        *string `json:"fs_type,omitempty"`
	Size          *string `json:"size,omitempty"`
	DateCreated   *string `json:"date_created,omitempty"`
	IsEncrypted   bool    `json:"is_encrypted,omitempty"`
}

type StorageProfileVGList struct {
	Name          *string     `json:"name,omitempty"`
	VgID          *string     `json:"vg_id,omitempty"`
	VgType        *string     `json:"vg_type,omitempty"`
	VgIscsiTarget *string     `json:"vg_iscsi_target,omitempty"`
	DiskList      []*DiskList `json:"disk_list,omitempty"`
}

type StorageProfile struct {
	HostOsType      *string                 `json:"host_os_type,omitempty"`
	Hypervisor      *string                 `json:"hypervisor,omitempty"`
	IsEraDriveOnEsx *bool                   `json:"is_era_drive_on_esx,omitempty"`
	LvmBased        *bool                   `json:"lvm_based,omitempty"`
	LvPath          []*string               `json:"lv_path,omitempty"`
	DiskList        []interface{}           `json:"disk_list,omitempty"`
	VgList          []*StorageProfileVGList `json:"vg_list,omitempty"`
}

type DriveSoftware struct {
	StorageProfile *StorageProfile `json:"storage_profile,omitempty"`
}

type DriveStorageInfo struct {
	AttachedVM *string        `json:"attachedVm,omitempty"`
	VgName     *string        `json:"vgName,omitempty"`
	VgUUID     *string        `json:"vgUuid,omitempty"`
	PdName     *string        `json:"pdName,omitempty"`
	Software   *DriveSoftware `json:"software,omitempty"`
}

type DriveInfo struct {
	StorageInfo   *DriveStorageInfo `json:"storage_info,omitempty"`
	SourceEraPath *string           `json:"source_era_path,omitempty"`
}

type Disks struct {
	ID            *string     `json:"id,omitempty"`
	VdiskUUID     *string     `json:"vdiskUuid,omitempty"`
	TimeMachineID interface{} `json:"timeMachineId,omitempty"`
	EraDriveID    *string     `json:"eraDriveId,omitempty"`
	EraCreated    *string     `json:"eraCreated,omitempty"`
	Status        *string     `json:"status,omitempty"`
	Type          *string     `json:"type,omitempty"`
	TotalSize     int         `json:"totalSize,omitempty"`
	UsedSize      int         `json:"usedSize,omitempty"`
	Info          interface{} `json:"info,omitempty"`
	DateCreated   *string     `json:"dateCreated,omitempty"`
	DateModified  *string     `json:"dateModified,omitempty"`
	OwnerID       *string     `json:"ownerId,omitempty"`
	Message       interface{} `json:"message,omitempty"`
}

type Drive struct {
	ID                 *string     `json:"id,omitempty"`
	Path               *string     `json:"path,omitempty"`
	HostID             *string     `json:"hostId,omitempty"`
	VgUUID             *string     `json:"vgUuid,omitempty"`
	ClusterID          string      `json:"clusterId,omitempty"`
	ProtectionDomainID *string     `json:"protectionDomainId,omitempty"`
	EraCreated         bool        `json:"eraCreated,omitempty"`
	Status             *string     `json:"status,omitempty"`
	TotalSize          int         `json:"totalSize,omitempty"`
	UsedSize           int         `json:"usedSize,omitempty"`
	Info               *DriveInfo  `json:"info,omitempty"`
	DateCreated        *string     `json:"dateCreated,omitempty"`
	DateModified       *string     `json:"dateModified,omitempty"`
	OwnerID            *string     `json:"ownerId,omitempty"`
	Metadata           interface{} `json:"metadata,omitempty"`
	EraDisks           []*Disks    `json:"eraDisks,omitempty"`
	ProtectionDomain   interface{} `json:"protectionDomain,omitempty"`
	Message            interface{} `json:"message,omitempty"`
}

type SoftwareInstallationsInfo struct {
	Owner *string `json:"owner,omitempty"`
}
type SoftwareInstallations struct {
	ID                       *string                    `json:"id,omitempty"`
	Name                     *string                    `json:"name,omitempty"`
	EraCreated               bool                       `json:"eraCreated,omitempty"`
	Type                     *string                    `json:"type,omitempty"`
	DbserverID               *string                    `json:"dbserverId,omitempty"`
	SoftwareProfileID        *string                    `json:"softwareProfileId,omitempty"`
	SoftwareProfileVersionID *string                    `json:"softwareProfileVersionId,omitempty"`
	Version                  *string                    `json:"version,omitempty"`
	OwnerID                  *string                    `json:"ownerId,omitempty"`
	Info                     *SoftwareInstallationsInfo `json:"info,omitempty"`
	Metadata                 interface{}                `json:"metadata,omitempty"`
	DateCreated              *string                    `json:"dateCreated,omitempty"`
	DateModified             *string                    `json:"dateModified,omitempty"`
}

type DBServerVMResponse struct {
	ID                         *string                     `json:"id,omitempty"`
	Name                       *string                     `json:"name,omitempty"`
	Description                *string                     `json:"description,omitempty"`
	OwnerID                    *string                     `json:"ownerId,omitempty"`
	DateCreated                *string                     `json:"dateCreated,omitempty"`
	DateModified               *string                     `json:"dateModified,omitempty"`
	DbserverClusterID          *string                     `json:"dbserverClusterId,omitempty"`
	VMClusterName              *string                     `json:"vmClusterName,omitempty"`
	VMClusterUUID              *string                     `json:"vmClusterUuid,omitempty"`
	Status                     *string                     `json:"status,omitempty"`
	ClientID                   *string                     `json:"clientId,omitempty"`
	NxClusterID                *string                     `json:"nxClusterId,omitempty"`
	EraDriveID                 *string                     `json:"eraDriveId,omitempty"`
	EraVersion                 *string                     `json:"eraVersion,omitempty"`
	VMTimeZone                 *string                     `json:"vmTimeZone,omitempty"`
	WorkingDirectory           *string                     `json:"workingDirectory,omitempty"`
	Type                       *string                     `json:"type,omitempty"`
	DatabaseType               *string                     `json:"databaseType,omitempty"`
	AccessKeyID                *string                     `json:"accessKeyId,omitempty"`
	DbserverInValidEaState     *bool                       `json:"dbserverInValidEaState,omitempty"`
	ValidDiagnosticBundleState *bool                       `json:"validDiagnosticBundleState,omitempty"`
	WindowsDBServer            *bool                       `json:"windowsDBServer,omitempty"`
	EraCreated                 *bool                       `json:"eraCreated,omitempty"`
	Internal                   *bool                       `json:"internal,omitempty"`
	Placeholder                *bool                       `json:"placeholder,omitempty"`
	Clustered                  *bool                       `json:"clustered,omitempty"`
	IsServerDriven             *bool                       `json:"is_server_driven,omitempty"`
	QueryCount                 *int                        `json:"queryCount,omitempty"`
	IPAddresses                []*string                   `json:"ipAddresses,omitempty"`
	AssociatedTimeMachineIds   []*string                   `json:"associatedTimeMachineIds,omitempty"`
	MacAddresses               []*string                   `json:"macAddresses,omitempty"`
	VMInfo                     *VMInfo                     `json:"vmInfo,omitempty"`
	Metadata                   *DBServerMetadata           `json:"metadata,omitempty"`
	Metric                     *Metric                     `json:"metric,omitempty"`
	LcmConfig                  *LcmConfig                  `json:"lcmConfig,omitempty"`
	EraDrive                   *Drive                      `json:"eraDrive,omitempty"`
	SoftwareInstallations      []*SoftwareInstallations    `json:"softwareInstallations,omitempty"`
	Properties                 []*DatabaseServerProperties `json:"properties,omitempty"`
	Tags                       []*Tags                     `json:"tags,omitempty"`
	Databases                  interface{}                 `json:"databases,omitempty"`
	Clones                     interface{}                 `json:"clones,omitempty"`
	AccessKey                  interface{}                 `json:"accessKey,omitempty"`
	ProtectionDomainID         interface{}                 `json:"protectionDomainId,omitempty"`
	AccessLevel                interface{}                 `json:"accessLevel,omitempty"`
	Fqdns                      interface{}                 `json:"fqdns,omitempty"`
	Info                       interface{}                 `json:"info,omitempty"`
	RequestedVersion           interface{}                 `json:"requestedVersion,omitempty"`
	AssociatedTimeMachineID    interface{}                 `json:"associated_time_machine_id,omitempty"`
	TimeMachineInfo            interface{}                 `json:"time_machine_info,omitempty"`
	ProtectionDomain           interface{}                 `json:"protectionDomain,omitempty"`
}

type ListDBServerVMResponse []DBServerVMResponse

type AccessInfo struct {
	AccessType        *string `json:"accessType,omitempty"`
	DestinationSubnet *string `json:"destinationSubnet,omitempty"`
}
type NetworkInfo struct {
	VlanName             *string       `json:"vlanName,omitempty"`
	VlanUUID             *string       `json:"vlanUuid,omitempty"`
	VlanType             *string       `json:"vlanType,omitempty"`
	Gateway              *string       `json:"gateway,omitempty"`
	SubnetMask           *string       `json:"subnetMask,omitempty"`
	Hostname             *string       `json:"hostname,omitempty"`
	DeviceName           *string       `json:"deviceName,omitempty"`
	MacAddress           *string       `json:"macAddress,omitempty"`
	Flags                *string       `json:"flags,omitempty"`
	Mtu                  *string       `json:"mtu,omitempty"`
	DefaultGatewayDevice *bool         `json:"defaultGatewayDevice,omitempty"`
	EraConfigured        *bool         `json:"eraConfigured,omitempty"`
	IPAddresses          []*string     `json:"ipAddresses,omitempty"`
	AccessInfo           []*AccessInfo `json:"accessInfo,omitempty"`
}

type VMInfo struct {
	OsType         *string         `json:"osType,omitempty"`
	OsVersion      *string         `json:"osVersion,omitempty"`
	Distribution   *string         `json:"distribution,omitempty"`
	SecureInfo     interface{}     `json:"secureInfo,omitempty"`
	Info           interface{}     `json:"info,omitempty"`
	DeregisterInfo *DeregisterInfo `json:"deregisterInfo,omitempty"`
	NetworkInfo    []*NetworkInfo  `json:"networkInfo,omitempty"`
}

type DBServerRegisterInput struct {
	DatabaseType     *string            `json:"databaseType,omitempty"`
	VMIP             *string            `json:"vmIp,omitempty"`
	NxClusterUUID    *string            `json:"nxClusterUuid,omitempty"`
	ForcedInstall    *bool              `json:"forcedInstall,omitempty"`
	WorkingDirectory *string            `json:"workingDirectory,omitempty"`
	Username         *string            `json:"username,omitempty"`
	Password         *string            `json:"password,omitempty"`
	SSHPrivateKey    *string            `json:"sshPrivateKey,omitempty"`
	ActionArguments  []*Actionarguments `json:"actionArguments,omitempty"`
}

type DBServerFilterRequest struct {
	ID                *string `json:"id,omitempty"`
	Name              *string `json:"name,omitempty"`
	IP                *string `json:"ip,omitempty"`
	VMClusterName     *string `json:"vm-cluster-name,omitempty"`
	VMClusterID       *string `json:"vm-cluster-uuid,omitempty"`
	DBServerClusterID *string `json:"dbserver-cluster-id,omitempty"`
	NxClusterID       *string `json:"nx-cluster-id,omitempty"`
}

type StretchedVlanMetadata struct {
	Gateway    *string `json:"gateway,omitempty"`
	SubnetMask *string `json:"subnetMask,omitempty"`
}

type StretchedVlansInput struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Type        *string                `json:"type,omitempty"`
	Metadata    *StretchedVlanMetadata `json:"metadata,omitempty"`
	VlanIDs     []*string              `json:"vlanIds,omitempty"`
}

type StretchedVlanResponse struct {
	ID          *string                  `json:"id,omitempty"`
	Name        *string                  `json:"name,omitempty"`
	Type        *string                  `json:"type,omitempty"`
	Description *string                  `json:"description,omitempty"`
	Metadata    *StretchedVlanMetadata   `json:"metadata,omitempty"`
	Vlans       []*NetworkIntentResponse `json:"vlans,omitempty"`
	VlanIDs     []*string                `json:"vlanIds,omitempty"`
}

type CloneRefreshInput struct {
	SnapshotID        *string `json:"snapshotId,omitempty"`
	UserPitrTimestamp *string `json:"userPitrTimestamp,omitempty"`
	Timezone          *string `json:"timeZone,omitempty"`
}

type NameValueParams struct {
	Name  *string `json:"name,omitempty"`
	Value *string `json:"value,omitempty"`
}

type NetworksInfo struct {
	Type        *string            `json:"type,omitempty"`
	NetworkInfo []*NameValueParams `json:"networkInfo,omitempty"`
	AccessType  []*string          `json:"accessType,omitempty"`
}

type ClusterIntentInput struct {
	ClusterName        *string            `json:"clusterName,omitempty"`
	ClusterDescription *string            `json:"clusterDescription,omitempty"`
	ClusterIP          *string            `json:"clusterIP,omitempty"`
	StorageContainer   *string            `json:"storageContainer,omitempty"`
	AgentVMPrefix      *string            `json:"agentVMPrefix,omitempty"`
	Port               *int               `json:"port,omitempty"`
	Protocol           *string            `json:"protocol,omitempty"`
	ClusterType        *string            `json:"clusterType,omitempty"`
	Version            *string            `json:"version,omitempty"`
	CredentialsInfo    []*NameValueParams `json:"credentialsInfo,omitempty"`
	AgentNetworkInfo   []*NameValueParams `json:"agentNetworkInfo,omitempty"`
	NetworksInfo       []*NetworksInfo    `json:"networksInfo,omitempty"`
}

type DeleteClusterInput struct {
	DeleteRemoteSites bool `json:"deleteRemoteSites,omitempty"`
}

type ClusterUpdateInput struct {
	Username    *string   `json:"username,omitempty"`
	Password    *string   `json:"password,omitempty"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	IPAddresses []*string `json:"ipAddresses,omitempty"`
}

type GetNetworkAvailableIPs []struct {
	ID           *string   `json:"id,omitempty"`
	Name         *string   `json:"name,omitempty"`
	PropertyName *string   `json:"propertyName,omitempty"`
	Type         *string   `json:"type,omitempty"`
	ClusterID    *string   `json:"clusterId,omitempty"`
	ClusterName  *string   `json:"clusterName,omitempty"`
	IPAddresses  []*string `json:"ipAddresses,omitempty"`
	Managed      bool      `json:"managed,omitempty"`
}
