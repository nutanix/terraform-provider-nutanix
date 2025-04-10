package calm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Operations struct {
	client *client.Client
}

type Service interface {
	ProvisionBlueprint(ctx context.Context, uuid string, input *BlueprintProvisionInput) (*AppProvisionTaskOutput, error)
	GetBlueprint(ctx context.Context, uuid string) (*BlueprintResponse, error)
	TaskPoll(ctx context.Context, bpUUID string, launchID string) (*PollResponse, error)
	DeleteApp(ctx context.Context, appUUID string) (*DeleteAppResp, error)
	SoftDeleteApp(ctx context.Context, appUUID string) (*DeleteAppResp, error)
	GetApp(ctx context.Context, appUUID string) (*AppResponse, error)
	PerformAction(ctx context.Context, appUUID string, spec *ActionSpec) (*AppActionResponse, error)
	AppRunlogs(ctx context.Context, appUUID, runlogUUID string) (*AppRunlogsResponse, error)
	ListBlueprint(ctx context.Context, filter *BlueprintListInput) (*BlueprintListResponse, error)
	GetRuntimeEditables(ctx context.Context, bpUUID string) (*RuntimeEditablesResponse, error)
	PatchApp(ctx context.Context, appUUID string, patchUUID string, input *PatchInput) (*AppTaskResponse, error)
	PerformActionUuid(ctx context.Context, appUUID string, actionUUID string, input *ActionInput) (*AppTaskResponse, error)
	ExecuteRunbook(ctx context.Context, rbUUID string, input *RunbookProvisionInput) (*RunbookResponse, error)
	ListRunbook(ctx context.Context, filter *RunbookListInput) (*RunbookListResponse, error)
	GetRunbook(ctx context.Context, rbUUID string) (*RunbookResponse, error)
	RbRunlogs(ctx context.Context, runlogUUID string) (*RbRunlogsResponse, error)
	GetAppProtectionPolicyList(ctx context.Context, bpUUID string, appUUID string, configUUID string, policyListInput *PolicyListInput) (*PolicyListResponse, error)
	RecoveryPointsList(ctx context.Context, appUUID string, input *RecoveryPointsListInput) (*RecoveryPointsListResponse, error)
	CreateBlueprint(ctx context.Context, input CreateBlueprintResponse) (*CreateBlueprintResponse, error)
	UpdateBlueprint(ctx context.Context, bpUUID string, input CreateBlueprintResponse) (*CreateBlueprintResponse, error)
	RecoveryPointsDelete(ctx context.Context, appUUID string, input *ActionInput) (*AppTaskResponse, error)
	RunbookImport(ctx context.Context, input *RunbookImportInput) (*RunbookImportResponse, error)
	DeleteRunbook(ctx context.Context, RbUUID string) (*DeleteRbResp, error)
	CreateEndpoint(ctx context.Context, input *EndpointCreateInput) (*EndpointCreateResponse, error)
}

func (op Operations) ProvisionBlueprint(ctx context.Context, bpUUID string, input *BlueprintProvisionInput) (*AppProvisionTaskOutput, error) {
	path := fmt.Sprintf("/blueprints/%s/simple_launch", bpUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppProvisionTaskOutput)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) GetBlueprint(ctx context.Context, bpUUID string) (*BlueprintResponse, error) {
	path := fmt.Sprintf("/blueprints/%s", bpUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(BlueprintResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) TaskPoll(ctx context.Context, bpUUID string, launchID string) (*PollResponse, error) {
	path := fmt.Sprintf("/blueprints/%s/pending_launches/%s", bpUUID, launchID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(PollResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) DeleteApp(ctx context.Context, id string) (*DeleteAppResp, error) {
	httpReq, err := op.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/apps/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(DeleteAppResp)
	return res, op.client.Do(ctx, httpReq, res)
}

func (op Operations) SoftDeleteApp(ctx context.Context, id string) (*DeleteAppResp, error) {
	httpReq, err := op.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/apps/%s?type=soft", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(DeleteAppResp)
	return res, op.client.Do(ctx, httpReq, res)
}

func (op Operations) GetApp(ctx context.Context, id string) (*AppResponse, error) {
	httpReq, err := op.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%s", id), nil)
	if err != nil {
		return nil, err
	}
	res := new(AppResponse)
	return res, op.client.Do(ctx, httpReq, res)
}

func (op Operations) PerformAction(ctx context.Context, appUUID string, input *ActionSpec) (*AppActionResponse, error) {
	path := fmt.Sprintf("/apps/%s/actions/run", appUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppActionResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) AppRunlogs(ctx context.Context, appUUID string, runlogUUID string) (*AppRunlogsResponse, error) {
	path := fmt.Sprintf("/apps/%s/app_runlogs/%s/output", appUUID, runlogUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(AppRunlogsResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) ListBlueprint(ctx context.Context, filter *BlueprintListInput) (*BlueprintListResponse, error) {
	path := "/blueprints/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, filter)

	appResponse := new(BlueprintListResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) GetRuntimeEditables(ctx context.Context, bpUUID string) (*RuntimeEditablesResponse, error) {
	path := fmt.Sprintf("/blueprints/%s/runtime_editables", bpUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(RuntimeEditablesResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) PatchApp(ctx context.Context, appUUID string, patchUUID string, input *PatchInput) (*AppTaskResponse, error) {
	path := fmt.Sprintf("/apps/%s/patch/%s/run", appUUID, patchUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppTaskResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) PerformActionUuid(ctx context.Context, appUUID string, actionUUID string, input *ActionInput) (*AppTaskResponse, error) {
	path := fmt.Sprintf("/apps/%s/actions/%s/run", appUUID, actionUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppTaskResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) ExecuteRunbook(ctx context.Context, rbUUID string, input *RunbookProvisionInput) (*RunbookResponse, error) {
	path := fmt.Sprintf("/runbooks/%s/execute", rbUUID)
	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	rbResponse := new(RunbookResponse)

	if err != nil {
		return nil, err
	}

	return rbResponse, op.client.Do(ctx, req, rbResponse)
}

func (op Operations) ListRunbook(ctx context.Context, filter *RunbookListInput) (*RunbookListResponse, error) {
	path := "/runbooks/list"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, filter)

	rbResponse := new(RunbookListResponse)

	if err != nil {
		return nil, err
	}

	return rbResponse, op.client.Do(ctx, req, rbResponse)
}

func (op Operations) GetRunbook(ctx context.Context, rbUUID string) (*RunbookResponse, error) {
	path := fmt.Sprintf("/runbooks/%s", rbUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(RunbookResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) RbRunlogs(ctx context.Context, runlogUUID string) (*RbRunlogsResponse, error) {
	path := fmt.Sprintf("/runbooks/runlogs/%s", runlogUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	rbResponse := new(RbRunlogsResponse)

	if err != nil {
		return nil, err
	}

	return rbResponse, op.client.Do(ctx, req, rbResponse)
}

func (op Operations) GetAppProtectionPolicyList(ctx context.Context, bpUUID string, appUUID string, configUUID string, policyListInput *PolicyListInput) (*PolicyListResponse, error) {
	path := fmt.Sprintf("/blueprints/%s/app_profile/%s/config_spec/%s/app_protection_policies/list", bpUUID, appUUID, configUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, policyListInput)

	plResponse := new(PolicyListResponse)

	if err != nil {
		return nil, err
	}

	return plResponse, op.client.Do(ctx, req, plResponse)
}

func (op Operations) RecoveryPointsList(ctx context.Context, appUUID string, input *RecoveryPointsListInput) (*RecoveryPointsListResponse, error) {
	path := fmt.Sprintf("/apps/%s/recovery_groups/list", appUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	listResponse := new(RecoveryPointsListResponse)

	if err != nil {
		return nil, err
	}

	return listResponse, op.client.Do(ctx, req, listResponse)
}

func (op Operations) CreateBlueprint(ctx context.Context, input CreateBlueprintResponse) (*CreateBlueprintResponse, error) {
	path := "/blueprints/import_json"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	bpResponse := new(CreateBlueprintResponse)

	if err != nil {
		return nil, err
	}

	return bpResponse, op.client.Do(ctx, req, bpResponse)
}

func (op Operations) RunbookImport(ctx context.Context, input *RunbookImportInput) (*RunbookImportResponse, error) {
	path := "/runbooks/import_json"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	RbImportResponse := new(RunbookImportResponse)

	if err != nil {
		return nil, err
	}

	return RbImportResponse, op.client.Do(ctx, req, RbImportResponse)
}

func (op Operations) UpdateBlueprint(ctx context.Context, bpUUID string, input CreateBlueprintResponse) (*CreateBlueprintResponse, error) {
	path := fmt.Sprintf("/blueprints/%s", bpUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPut, path, input)

	bpResponse := new(CreateBlueprintResponse)

	if err != nil {
		return nil, err
	}

	return bpResponse, op.client.Do(ctx, req, bpResponse)
}

func (op Operations) DeleteRunbook(ctx context.Context, RbUUID string) (*DeleteRbResp, error) {
	httpReq, err := op.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf("/runbooks/%s", RbUUID), nil)
	if err != nil {
		return nil, err
	}
	res := new(DeleteRbResp)
	return res, op.client.Do(ctx, httpReq, res)
}

func (op Operations) RecoveryPointsDelete(ctx context.Context, appUUID string, input *ActionInput) (*AppTaskResponse, error) {
	path := fmt.Sprintf("/apps/%s/recovery_group_delete", appUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppTaskResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) CreateEndpoint(ctx context.Context, input *EndpointCreateInput) (*EndpointCreateResponse, error) {
	path := "/endpoints"

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	endpointResponse := new(EndpointCreateResponse)

	if err != nil {
		return nil, err
	}

	return endpointResponse, op.client.Do(ctx, req, endpointResponse)
}
