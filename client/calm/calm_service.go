package calm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Operations implements Service interface
type Operations struct {
	client *client.Client
}

type Service interface {
	CreateProjectQuota(ctx context.Context, request *ProjectQuotaIntentInput) (*ProjectQuotaIntentResponse, error)
	UpdateProjectQuota(ctx context.Context, quotaID string, request *ProjectQuotaIntentInput) (*ProjectQuotaIntentResponse, error)
	EnableProjectQuota(ctx context.Context, request *EnableProjectQuotaInput) (*ProjectQuotaIntentResponse, error)
}

func (op Operations) CreateProjectQuota(ctx context.Context, request *ProjectQuotaIntentInput) (*ProjectQuotaIntentResponse, error) {
	req, err := op.client.NewRequest(ctx, http.MethodPost, "/quotas", request)
	if err != nil {
		return nil, err
	}

	projectResponse := new(ProjectQuotaIntentResponse)

	return projectResponse, op.client.Do(ctx, req, projectResponse)
}

func (op Operations) UpdateProjectQuota(ctx context.Context, quotaID string, request *ProjectQuotaIntentInput) (*ProjectQuotaIntentResponse, error) {
	path := fmt.Sprintf("/quotas/%s", quotaID)
	req, err := op.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return nil, err
	}

	projectResponse := new(ProjectQuotaIntentResponse)

	return projectResponse, op.client.Do(ctx, req, projectResponse)
}

func (op Operations) EnableProjectQuota(ctx context.Context, request *EnableProjectQuotaInput) (*ProjectQuotaIntentResponse, error) {
	req, err := op.client.NewRequest(ctx, http.MethodPut, "/quotas/update/state", request)
	if err != nil {
		return nil, err
	}

	quotaEnableResponse := new(ProjectQuotaIntentResponse)

	return quotaEnableResponse, op.client.Do(ctx, req, quotaEnableResponse)
}
