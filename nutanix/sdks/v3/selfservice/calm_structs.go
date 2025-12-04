package selfservice

import "encoding/json"

type BlueprintProvisionInput struct {
	Spec BPspec `json:"spec,omitempty"`
}

type BPspec struct {
	AppName             string              `json:"app_name,omitempty"`
	AppDesc             string              `json:"app_description,omitempty"`
	RuntimeEditables    *RuntimeEditables   `json:"runtime_editables,omitempty"`
	AppProfileReference AppProfileReference `json:"app_profile_reference,omitempty"`
}

type AppProfileReference struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type AppProvisionTaskOutput struct {
	AppName string          `json:"app_name"`
	AppDesc string          `json:"app_description"`
	Spec    json.RawMessage `json:"spec"`
	Status  BpRespStatus    `json:"status"`
}

type BlueprintResponse struct {
	APIVersion *string `json:"api_version" mapstructure:"api_version"`

	Spec json.RawMessage `json:"spec,omitempty" mapstructure:"spec,omitempty"`

	Status BpRespStatus `json:"status,omitempty" mapstructure:"status,omitempty"`

	Metadata json.RawMessage `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
}

type BpRespStatus struct {
	RequestID string `json:"request_id,omitempty" mapstructure:"request_id,omitempty"`
}

type PollResponse struct {
	Status     PollStatus `json:"status,omitempty"`
	APIVersion *string    `json:"api_version,omitempty"`
}

type PollStatus struct {
	AppUUID       *string `json:"application_uuid,omitempty"`
	AppName       *string `json:"app_name,omitempty"`
	State         *string `json:"state,omitempty"`
	BlueprintUUID *string `json:"bp_uuid,omitempty"`
	BpName        *string `json:"bp_name,omitempty"`
}

type DeleteAppResp struct {
	Status     DeleteAppStatus `json:"status"`
	APIVersion string          `json:"api_version"`
}

type DeleteAppStatus struct {
	RunlogUUID    string `json:"runlog_uuid"`
	ErgonTaskUUID string `json:"ergon_task_uuid"`
}

type AppResponse struct {
	Status     json.RawMessage `json:"status"`
	Spec       json.RawMessage `json:"spec"`
	APIVersion string          `json:"api_version"`
	Metadata   json.RawMessage `json:"metadata"`
}

type ActionSpec struct {
	Name string `json:"name"`
}

type AppActionResponse struct {
	RunlogUUID string `json:"runlog_uuid"`
}

type AppRunlogsResponse struct {
	APIVersion *string          `json:"api_version"`
	Status     *AppRunlogStatus `json:"status"`
	OutputList []*OutputList    `json:"output_list"`
}

type AppRunlogStatus struct {
	RunlogState *string `json:"runlog_state"`
	ExitCode    *int    `json:"exit_code"`
}

type OutputList struct {
	Output *string `json:"output"`
	Script *string `json:"script"`
}

type ListInput struct {
	Filter string `json:"filter"`
}

type BlueprintListInput = ListInput
type ApplicationListInput = ListInput
type AccountsListInput = ListInput

type ListResponse struct {
	Entities json.RawMessage `json:"entities"`
}

type BlueprintListResponse = ListResponse
type ApplicationListResponse = ListResponse

type RuntimeEditablesResponse struct {
	Resources []*Resources `json:"resources"`
}

type Resources struct {
	AppProfileReference *AppProfileReference `json:"app_profile_reference"`
	RuntimeEditables    *RuntimeEditables    `json:"runtime_editables"`
}

type RuntimeEditables struct {
	ActionList         []*RuntimeSpec `json:"action_list,omitempty"`
	ServiceList        []*RuntimeSpec `json:"service_list,omitempty"`
	CredentialList     []*RuntimeSpec `json:"credential_list,omitempty"`
	SubstrateList      []*RuntimeSpec `json:"substrate_list,omitempty"`
	PackageList        []*RuntimeSpec `json:"package_list,omitempty"`
	SnapshotConfifList []*RuntimeSpec `json:"snapshot_config_list,omitempty"`
	AppProfile         *RuntimeSpec   `json:"app_profile,omitempty"`
	TaskList           []*RuntimeSpec `json:"task_list,omitempty"`
	RestoreConfigList  []*RuntimeSpec `json:"restore_config_list,omitempty"`
	VariableList       []*RuntimeSpec `json:"variable_list,omitempty"`
	DeploymentList     []*RuntimeSpec `json:"deployment_list,omitempty"`
}

type RuntimeSpec struct {
	Description *string          `json:"description,omitempty"`
	Value       *json.RawMessage `json:"value,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Context     *string          `json:"context,omitempty"`
	Type        *string          `json:"type,omitempty"`
	UUID        *string          `json:"uuid,omitempty"`
}
type PatchInput struct {
	Spec       PatchSpec              `json:"spec"`
	APIVersion string                 `json:"api_version"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type PatchSpec struct {
	Args       ArgsSpec `json:"args"`
	TargetUUID string   `json:"target_uuid"`
	TargetKind string   `json:"target_kind"`
}

type ArgsSpec struct {
	Variables []*VariableList        `json:"variables"`
	Patch     map[string]interface{} `json:"patch"`
}

type VariableList struct {
	TaskUUID string `json:"task_uuid,omitempty"`
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
}

type AppTaskResponse struct {
	Status     ActionRunStatus `json:"status"`
	Spec       json.RawMessage `json:"spec"`
	APIVersion string          `json:"api_version"`
	Metadata   json.RawMessage `json:"metadata"`
}

type ActionRunStatus struct {
	RunlogUUID string `json:"runlog_uuid"`
}

type ActionInput struct {
	Spec       TaskSpec               `json:"spec"`
	APIVersion string                 `json:"api_version"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type TaskSpec struct {
	Args       []*VariableList `json:"args"`
	TargetUUID string          `json:"target_uuid"`
	TargetKind string          `json:"target_kind"`
}

type ActionResponse struct {
	Status ActionRunStatus `json:"status,omitempty"`
}

type PolicyListInput struct {
	Length int    `json:"length,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type PolicyListResponse struct {
	Entities []map[string]interface{} `json:"entities,omitempty"`
}

type RecoveryPointsListInput struct {
	Filter string `json:"filter,omitempty"`
	Length int    `json:"length,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type ListResponseData struct {
	APIVersion string                   `json:"api_version,omitempty"`
	Metadata   map[string]interface{}   `json:"metadata,omitempty"`
	Entities   []map[string]interface{} `json:"entities,omitempty"`
}
type RecoveryPointsListResponse = ListResponseData
type AccountsListResponse = ListResponseData
