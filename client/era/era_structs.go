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
	Enabled    bool   `json:"enabled"`
	Dayofmonth string `json:"dayOfMonth"`
}

type Quartelyschedule struct {
	Enabled    bool   `json:"enabled"`
	Startmonth string `json:"startMonth"`
	Dayofmonth string `json:"dayOfMonth"`
}

type Yearlyschedule struct {
	Enabled    bool   `json:"enabled"`
	Dayofmonth int    `json:"dayOfMonth"`
	Month      string `json:"month"`
}

type Schedule struct {
	Snapshottimeofday  Snapshottimeofday  `json:"snapshotTimeOfDay"`
	Continuousschedule Continuousschedule `json:"continuousSchedule"`
	Weeklyschedule     Weeklyschedule     `json:"weeklySchedule"`
	Monthlyschedule    Monthlyschedule    `json:"monthlySchedule"`
	Quartelyschedule   Quartelyschedule   `json:"quartelySchedule"`
	Yearlyschedule     Yearlyschedule     `json:"yearlySchedule"`
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

// Properties are redeclared in the block with additional information along with ListClusterResponse
// type Properties struct {
// 	Name   string `json:"name"`
// 	Value  string `json:"value"`
// 	Secure bool   `json:"secure"`
// }

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

/*
type UpdateDatabaseResponse struct {
	ID                       string          `json:"id"`
	Name                     string          `json:"name"`
	Description              string          `json:"description"`
	Ownerid                  string          `json:"ownerId"`
	Datecreated              string          `json:"dateCreated"`
	Datemodified             string          `json:"dateModified"`
	Properties               []Properties    `json:"properties"`
	Tags                     []interface{}   `json:"tags"`
	Clustered                bool            `json:"clustered"`
	Clone                    bool            `json:"clone"`
	Eracreated               bool            `json:"eraCreated"`
	Internal                 bool            `json:"internal"`
	Placeholder              bool            `json:"placeholder"`
	Databasename             string          `json:"databaseName"`
	Type                     string          `json:"type"`
	Databaseclustertype      interface{}     `json:"databaseClusterType"`
	Status                   string          `json:"status"`
	Databasestatus           string          `json:"databaseStatus"`
	Dbserverlogicalclusterid interface{}     `json:"dbserverLogicalClusterId"`
	Timemachineid            string          `json:"timeMachineId"`
	Parenttimemachineid      interface{}     `json:"parentTimeMachineId"`
	Timezone                 string          `json:"timeZone"`
	Info                     Info            `json:"info"`
	Metadata                 Metadata        `json:"metadata"`
	Metric                   interface{}     `json:"metric"`
	Category                 string          `json:"category"`
	Lcmconfig                interface{}     `json:"lcmConfig"`
	Timemachine              interface{}     `json:"timeMachine"`
	Dbserverlogicalcluster   interface{}     `json:"dbserverlogicalCluster"`
	Databasenodes            []Databasenodes `json:"databaseNodes"`
	Linkeddatabases          interface{}     `json:"linkedDatabases"`
}
type Properties struct {
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
	DataDisks      DataDisks      `json:"data_disks"`
	LogDisks       LogDisks       `json:"log_disks"`
	ArchiveStorage ArchiveStorage `json:"archive_storage"`
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
	Storage      Storage      `json:"storage"`
	VMProperties VMProperties `json:"vm_properties"`
	BpgDbParam   BpgDbParam   `json:"bpg_db_param"`
}
type Info struct {
	BpgConfigs BpgConfigs `json:"bpg_configs"`
}
type Info struct {
	Secureinfo Secureinfo `json:"secureInfo"`
	Info       Info       `json:"info"`
}
type Metadata struct {
	Secureinfo                          interface{} `json:"secureInfo"`
	Info                                interface{} `json:"info"`
	Deregisterinfo                      interface{} `json:"deregisterInfo"`
	Tmactivateoperationid               string      `json:"tmActivateOperationId"`
	Createddbservers                    interface{} `json:"createdDbservers"`
	Registereddbservers                 interface{} `json:"registeredDbservers"`
	Lastrefreshtimestamp                interface{} `json:"lastRefreshTimestamp"`
	Lastrequestedrefreshtimestamp       interface{} `json:"lastRequestedRefreshTimestamp"`
	Statebeforerefresh                  interface{} `json:"stateBeforeRefresh"`
	Statebeforerestore                  interface{} `json:"stateBeforeRestore"`
	Statebeforescaling                  interface{} `json:"stateBeforeScaling"`
	Logcatchupforrestoredispatched      bool        `json:"logCatchUpForRestoreDispatched"`
	Lastlogcatchupforrestoreoperationid interface{} `json:"lastLogCatchUpForRestoreOperationId"`
	Originaldatabasename                interface{} `json:"originalDatabaseName"`
}
type Info struct {
}
type Info struct {
	Secureinfo interface{} `json:"secureInfo"`
	Info       Info        `json:"info"`
}
type Databasenodes struct {
	ID                     string        `json:"id"`
	Name                   string        `json:"name"`
	Description            string        `json:"description"`
	Ownerid                string        `json:"ownerId"`
	Datecreated            string        `json:"dateCreated"`
	Datemodified           string        `json:"dateModified"`
	Properties             []interface{} `json:"properties"`
	Tags                   []interface{} `json:"tags"`
	Databaseid             string        `json:"databaseId"`
	Status                 string        `json:"status"`
	Databasestatus         string        `json:"databaseStatus"`
	Primary                bool          `json:"primary"`
	Dbserverid             string        `json:"dbserverId"`
	Softwareinstallationid string        `json:"softwareInstallationId"`
	Protectiondomainid     string        `json:"protectionDomainId"`
	Info                   Info          `json:"info"`
	Metadata               interface{}   `json:"metadata"`
	Dbserver               interface{}   `json:"dbserver"`
	Protectiondomain       interface{}   `json:"protectionDomain"`
	Valideastate           bool          `json:"validEaState"`
}
*/

// TODO: Change some fields like metadata etc to map[string]interface{} if necessary

type GetDatabaseResponse struct {
	ID                       string                 `json:"id"`
	Name                     string                 `json:"name"`
	Description              string                 `json:"description"`
	Ownerid                  string                 `json:"ownerId"`
	Datecreated              string                 `json:"dateCreated"`
	Datemodified             string                 `json:"dateModified"`
	Properties               []DBInstanceProperties `json:"properties"`
	Tags                     []interface{}          `json:"tags"`
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
	Info                     Info                   `json:"info"`
	Metadata                 DBInstanceMetadata     `json:"metadata"`
	Metric                   interface{}            `json:"metric"`
	Category                 string                 `json:"category"`
	Lcmconfig                interface{}            `json:"lcmConfig"`
	Dbserverlogicalcluster   interface{}            `json:"dbserverlogicalCluster"`
	Databasenodes            []Databasenodes        `json:"databaseNodes"`
	Linkeddatabases          interface{}            `json:"linkedDatabases"`
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
	Count int `json:"count"`
}

type LogDisks struct {
	Count int `json:"count"`
	Size  int `json:"size"`
}
type ArchiveStorage struct {
	Size int `json:"size"`
}
type Storage struct {
	DataDisks      DataDisks      `json:"data_disks"`
	LogDisks       LogDisks       `json:"log_disks"`
	ArchiveStorage ArchiveStorage `json:"archive_storage"`
}
type VMProperties struct {
	NrHugepages             int `json:"nr_hugepages"`
	OvercommitMemory        int `json:"overcommit_memory"`
	DirtyBackgroundRatio    int `json:"dirty_background_ratio"`
	DirtyRatio              int `json:"dirty_ratio"`
	DirtyExpireCentisecs    int `json:"dirty_expire_centisecs"`
	DirtyWritebackCentisecs int `json:"dirty_writeback_centisecs"`
	Swappiness              int `json:"swappiness"`
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
	Storage      Storage      `json:"storage"`
	VMProperties VMProperties `json:"vm_properties"`
	BpgDbParam   BpgDbParam   `json:"bpg_db_param"`
}

type Info struct {
	Secureinfo Secureinfo  `json:"secureInfo"`
	Info       interface{} `json:"info"`
	BpgConfigs BpgConfigs  `json:"bpg_configs"`
}
type DBInstanceMetadata struct {
	Secureinfo                          interface{} `json:"secureInfo"`
	Info                                interface{} `json:"info"`
	Deregisterinfo                      interface{} `json:"deregisterInfo"`
	Tmactivateoperationid               string      `json:"tmActivateOperationId"`
	Createddbservers                    interface{} `json:"createdDbservers"`
	Registereddbservers                 interface{} `json:"registeredDbservers"`
	Lastrefreshtimestamp                interface{} `json:"lastRefreshTimestamp"`
	Lastrequestedrefreshtimestamp       interface{} `json:"lastRequestedRefreshTimestamp"`
	Statebeforerefresh                  interface{} `json:"stateBeforeRefresh"`
	Statebeforerestore                  interface{} `json:"stateBeforeRestore"`
	Statebeforescaling                  interface{} `json:"stateBeforeScaling"`
	Logcatchupforrestoredispatched      bool        `json:"logCatchUpForRestoreDispatched"`
	Lastlogcatchupforrestoreoperationid interface{} `json:"lastLogCatchUpForRestoreOperationId"`
	Originaldatabasename                interface{} `json:"originalDatabaseName"`
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
type Protectiondomain struct {
	ID                                 string        `json:"id"`
	Name                               string        `json:"name"`
	Eracreated                         bool          `json:"eraCreated"`
	Description                        string        `json:"description"`
	Type                               string        `json:"type"`
	Status                             string        `json:"status"`
	Cloudid                            string        `json:"cloudId"`
	Parentprotectiondomainid           interface{}   `json:"parentProtectionDomainId"`
	Ownerid                            string        `json:"ownerId"`
	Datecreated                        string        `json:"dateCreated"`
	Datemodified                       string        `json:"dateModified"`
	Info                               interface{}   `json:"info"`
	Metadata                           interface{}   `json:"metadata"`
	Replicatedprotectiondomains        []interface{} `json:"replicatedProtectionDomains"`
	Assocentities                      interface{}   `json:"assocEntities"`
	Referencecount                     int           `json:"referenceCount"`
	Referencingtimemachines            interface{}   `json:"referencingTimeMachines"`
	Referencingtimemachinesanystatus   interface{}   `json:"referencingTimeMachinesAnyStatus"`
	Sourcepd                           bool          `json:"sourcePD"`
	Replicatedpd                       bool          `json:"replicatedPD"`
	Timemachinereferencecount          int           `json:"timeMachineReferenceCount"`
	Timemachineanystatusreferencecount int           `json:"timeMachineAnyStatusReferenceCount"`
}
type Databasenodes struct {
	ID                     string           `json:"id"`
	Name                   string           `json:"name"`
	Description            string           `json:"description"`
	Ownerid                string           `json:"ownerId"`
	Datecreated            string           `json:"dateCreated"`
	Datemodified           string           `json:"dateModified"`
	Properties             []interface{}    `json:"properties"`
	Tags                   []interface{}    `json:"tags"`
	Databaseid             string           `json:"databaseId"`
	Status                 string           `json:"status"`
	Databasestatus         string           `json:"databaseStatus"`
	Primary                bool             `json:"primary"`
	Dbserverid             string           `json:"dbserverId"`
	Softwareinstallationid string           `json:"softwareInstallationId"`
	Protectiondomainid     string           `json:"protectionDomainId"`
	Info                   Info             `json:"info"`
	Metadata               interface{}      `json:"metadata"`
	Dbserver               Dbserver         `json:"dbserver"`
	Protectiondomain       Protectiondomain `json:"protectionDomain"`
	Valideastate           bool             `json:"validEaState"`
}
