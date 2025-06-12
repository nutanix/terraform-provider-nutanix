package selfservice

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
	ListApplication(ctx context.Context, filter *ApplicationListInput) (*ApplicationListResponse, error)
	GetRuntimeEditables(ctx context.Context, bpUUID string) (*RuntimeEditablesResponse, error)
	PatchApp(ctx context.Context, appUUID string, patchUUID string, input *PatchInput) (*AppTaskResponse, error)
	PerformActionUUID(ctx context.Context, appUUID string, actionUUID string, input *ActionInput) (*AppTaskResponse, error)
	RecoveryPointsList(ctx context.Context, appUUID string, input *RecoveryPointsListInput) (*RecoveryPointsListResponse, error)
	GetAppProtectionPolicyList(ctx context.Context, bpUUID string, appUUID string, configUUID string, policyListInput *PolicyListInput) (*PolicyListResponse, error)
	RecoveryPointsDelete(ctx context.Context, appUUID string, input *ActionInput) (*AppTaskResponse, error)
	ListAccounts(ctx context.Context, filter *AccountsListInput) (*AccountsListResponse, error)
}

func (op Operations) ProvisionBlueprint(ctx context.Context, bpUUID string, input *BlueprintProvisionInput) (*AppProvisionTaskOutput, error) {
	path := fmt.Sprintf(launchBlueprintAPI, bpUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppProvisionTaskOutput)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) GetBlueprint(ctx context.Context, bpUUID string) (*BlueprintResponse, error) {
	path := fmt.Sprintf(getBlueprintAPI, bpUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(BlueprintResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) TaskPoll(ctx context.Context, bpUUID string, launchID string) (*PollResponse, error) {
	path := fmt.Sprintf(pendingLaunchBlueprintAPI, bpUUID, launchID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(PollResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) DeleteApp(ctx context.Context, id string) (*DeleteAppResp, error) {
	httpReq, err := op.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf(getApplicationAPI, id), nil)
	if err != nil {
		return nil, err
	}
	res := new(DeleteAppResp)
	return res, op.client.Do(ctx, httpReq, res)
}

func (op Operations) SoftDeleteApp(ctx context.Context, id string) (*DeleteAppResp, error) {
	httpReq, err := op.client.NewRequest(ctx, http.MethodDelete, fmt.Sprintf(softDeleteApplicationAPI, id), nil)
	if err != nil {
		return nil, err
	}
	res := new(DeleteAppResp)
	return res, op.client.Do(ctx, httpReq, res)
}

func (op Operations) GetApp(ctx context.Context, id string) (*AppResponse, error) {
	httpReq, err := op.client.NewRequest(ctx, http.MethodGet, fmt.Sprintf(getApplicationAPI, id), nil)
	if err != nil {
		return nil, err
	}
	res := new(AppResponse)
	return res, op.client.Do(ctx, httpReq, res)
}

func (op Operations) PerformAction(ctx context.Context, appUUID string, input *ActionSpec) (*AppActionResponse, error) {
	path := fmt.Sprintf(runApplicationSystemActionAPI, appUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppActionResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) AppRunlogs(ctx context.Context, appUUID string, runlogUUID string) (*AppRunlogsResponse, error) {
	path := fmt.Sprintf(getAppRunlogOutputAPI, appUUID, runlogUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(AppRunlogsResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) ListBlueprint(ctx context.Context, filter *BlueprintListInput) (*BlueprintListResponse, error) {
	path := listBlueprintAPI

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, filter)

	appResponse := new(BlueprintListResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) ListApplication(ctx context.Context, filter *ApplicationListInput) (*ApplicationListResponse, error) {
	path := listApplicationAPI

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, filter)

	appResponse := new(ApplicationListResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) GetRuntimeEditables(ctx context.Context, bpUUID string) (*RuntimeEditablesResponse, error) {
	path := fmt.Sprintf(getBlueprintRuntimeEditables, bpUUID)

	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)

	appResponse := new(RuntimeEditablesResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) PatchApp(ctx context.Context, appUUID string, patchUUID string, input *PatchInput) (*AppTaskResponse, error) {
	path := fmt.Sprintf(runPatchActionAPI, appUUID, patchUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppTaskResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) PerformActionUUID(ctx context.Context, appUUID string, actionUUID string, input *ActionInput) (*AppTaskResponse, error) {
	path := fmt.Sprintf(runApplicationCustomActionAPI, appUUID, actionUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppTaskResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

func (op Operations) GetAppProtectionPolicyList(ctx context.Context, bpUUID string, appUUID string, configUUID string, policyListInput *PolicyListInput) (*PolicyListResponse, error) {
	path := fmt.Sprintf(listAppProtectionPolicyAPI, bpUUID, appUUID, configUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, policyListInput)

	plResponse := new(PolicyListResponse)

	if err != nil {
		return nil, err
	}

	return plResponse, op.client.Do(ctx, req, plResponse)
}

func (op Operations) RecoveryPointsList(ctx context.Context, appUUID string, input *RecoveryPointsListInput) (*RecoveryPointsListResponse, error) {
	path := fmt.Sprintf(listAppRecoveryPointsAPI, appUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	listResponse := new(RecoveryPointsListResponse)

	if err != nil {
		return nil, err
	}

	return listResponse, op.client.Do(ctx, req, listResponse)
}

func (op Operations) RecoveryPointsDelete(ctx context.Context, appUUID string, input *ActionInput) (*AppTaskResponse, error) {
	path := fmt.Sprintf(deleteRecoveryPointsAPI, appUUID)

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, input)

	appResponse := new(AppTaskResponse)

	if err != nil {
		return nil, err
	}

	return appResponse, op.client.Do(ctx, req, appResponse)
}

// ListAccounts lists the accounts available in the Nutanix environment.
func (op Operations) ListAccounts(ctx context.Context, filter *AccountsListInput) (*AccountsListResponse, error) {
	path := listAccountsAPI

	req, err := op.client.NewRequest(ctx, http.MethodPost, path, filter)

	accResponse := new(AccountsListResponse)

	if err != nil {
		return nil, err
	}

	return accResponse, op.client.Do(ctx, req, accResponse)
}
