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
	GetApp(ctx context.Context, appUUID string) (*AppResponse, error)
	PerformAction(ctx context.Context, appUUID string, spec *ActionSpec) (*AppActionResponse, error)
	AppRunlogs(ctx context.Context, appUUID, runlogUUID string) (*AppRunlogsResponse, error)
	ListBlueprint(ctx context.Context, filter *BlueprintListInput) (*BlueprintListResponse, error)
	GetRuntimeEditables(ctx context.Context, bpUUID string) (*RuntimeEditablesResponse, error)
	PatchApp(ctx context.Context, appUUID string, patchUUID string, input *PatchInput) (*AppPatchResponse, error)
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

func (op Operations) PatchApp(ctx context.Context, appUUID string, patchUUID string, input *PatchInput) (*AppPatchResponse, error) {
	path := fmt.Sprintf("/apps/%s/patch/%s/run", appUUID, patchUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppPatchResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}
