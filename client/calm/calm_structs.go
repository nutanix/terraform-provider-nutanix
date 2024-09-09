package calm

import "encoding/json"

type BlueprintProvisionInput struct {
	Spec BPspec `json:"spec,omitempty"`
}

type BPspec struct {
	AppName             string                 `json:"app_name,omitempty"`
	AppDesc             string                 `json:"app_description,omitempty"`
	RuntimeEditables    map[string]interface{} `json:"runtime_editables,omitempty"`
	AppProfileReference AppProfileReference    `json:"app_profile_reference,omitempty"`
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

type BlueprintListResponse struct {
	Entities json.RawMessage `json:"entities"`
}

// type ActionPayload struct {
// 	TargetUUID string `json:"target_uuid"`
// 	TargetKind string `json:"target_kind"`
// }
