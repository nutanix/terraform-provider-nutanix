package calm

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
	ApiVersion string          `json:"api_version"`
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
	// Spec *ActionPayload `json:"spec"`
	// ApiVersion string                 `json:"api_version"`
	// Metadata   map[string]interface{} `json:"metadata"`
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

type BlueprintListInput struct {
	Filter string `json:"filter"`
}

type ApplicationListInput struct {
	Filter string `json:"filter"`
}

type BlueprintListResponse struct {
	Entities json.RawMessage `json:"entities"`
}

type ApplicationListResponse struct {
	Entities json.RawMessage `json:"entities"`
}

// type ActionPayload struct {
// 	TargetUUID string `json:"target_uuid"`
// 	TargetKind string `json:"target_kind"`
// }

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
	Spec       TaskSpec               `json:"spec,omitempty"`
	APIVersion string                 `json:"api_version,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type TaskSpec struct {
	Args       []*VariableList `json:"args,omitempty"`
	TargetUUID string          `json:"target_uuid,omitempty"`
	TargetKind string          `json:"target_kind,omitempty"`
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

//	type AppCustomActionResponse struct {
//		Status     ActionRunStatus `json:"status"`
//		Spec       json.RawMessage `json:"spec"`
//		APIVersion string          `json:"api_version"`
//		Metadata   json.RawMessage `json:"metadata"`
//	}
type RunbookProvisionInput struct {
	Spec         RunbookProvisionSpec `json:"spec,omitempty"`
	VariableList json.RawMessage      `json:"variable_list,omitempty"`
}

type RunbookProvisionSpec struct {
	Args []RunbookArgs `json:"args,omitempty"`
}

type RunbookArgs struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type RBspec struct {
	RbName string `json:"rb_name,omitempty"`
}

type RunbookResponse struct {
	Spec   json.RawMessage `json:"spec"`
	Status RbRespStatus    `json:"status"`
}

type RunbookListInput struct {
	Filter string `json:"filter"`
}

type RunbookListResponse struct {
	Entities json.RawMessage `json:"entities"`
}

type RbRespStatus struct {
	RunlogUUID string `json:"runlog_uuid,omitempty"`
}

type RbRunlogsResponse struct {
	Status *RbRunlogStatus `json:"status"`
}

type RbRunlogStatus struct {
	State              *string             `json:"state"`
	OutputVariableList []*RbOutputVariable `json:"output_variable_list,omitempty"`
}

type RbOutputVariable struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type RecoveryPointsListInput struct {
	Filter string `json:"filter,omitempty"`
	Length int    `json:"length,omitempty"`
	Offset int    `json:"offset,omitempty"`
}

type RecoveryPointsListResponse struct {
	APIVersion string                   `json:"api_version,omitempty"`
	Metadata   map[string]interface{}   `json:"metadata,omitempty"`
	Entities   []map[string]interface{} `json:"entities,omitempty"`
}

type CreateBlueprintResponse struct {
	APIVersion string                 `json:"api_version"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       map[string]interface{} `json:"spec"`
}

type RefObject struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
}

type TaskDef struct {
	Name             string                 `json:"name"`
	Type             string                 `json:"type"`
	Attrs            map[string]interface{} `json:"attrs"`
	ChildTaskRefList []RefObject            `json:"child_tasks_local_reference_list"`
	StatusMapList    []interface{}          `json:"status_map_list"`
	VariableList     []interface{}          `json:"variable_list"`
	Retries          string                 `json:"retries"`
	Timeout          string                 `json:"timeout_secs"`
}

type RunbookDefinition struct {
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	MainTaskLocalRef   RefObject     `json:"main_task_local_reference"`
	TaskDefList        []TaskDef     `json:"task_definition_list"`
	VariableList       []interface{} `json:"variable_list"`
	OutputVariableList []interface{} `json:"output_variable_list"`
}

type RunbookResources struct {
	Runbook           RunbookDefinition      `json:"runbook"`
	EndpointDefList   []interface{}          `json:"endpoint_definition_list"`
	CredentialDefList []interface{}          `json:"credential_definition_list"`
	DefaultTargetRef  RefObject              `json:"default_target_reference"`
	ClientAttrs       map[string]interface{} `json:"client_attrs"`
}

type RunbookSpec struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Resources   RunbookResources `json:"resources"`
}

type RunbookImportInput struct {
	Spec       RunbookSpec            `json:"spec"`
	APIVersion string                 `json:"api_version"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type RunbookImportResponse struct {
	Status     map[string]interface{} `json:"status"`
	Spec       map[string]interface{} `json:"spec"`
	Metadata   map[string]interface{} `json:"metadata"`
	APIVersion string                 `json:"api_version"`
}

type DeleteRbResp struct {
	Description string `json:"description,omitempty"`
}

type EndpointCreateInput struct {
	Spec       map[string]interface{} `json:"spec"`
	Metadata   map[string]interface{} `json:"metadata"`
	APIVersion string                 `json:"api_version"`
}

type EndpointCreateResponse struct {
	Status     map[string]interface{} `json:"status"`
	Spec       map[string]interface{} `json:"spec"`
	Metadata   map[string]interface{} `json:"metadata"`
	APIVersion string                 `json:"api_version"`
}
