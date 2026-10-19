package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/registereddomains"
	multidomainPrism "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixRegisteredDomainRefreshApiCredentialsV2 exposes the
// refresh-api-credentials action which rotates the API credentials used to
// communicate with a registered domain.
func ResourceNutanixRegisteredDomainRefreshApiCredentialsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Create,
		ReadContext:   ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Read,
		DeleteContext: ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NCDomainManagementAPI

	extID := d.Get("ext_id").(string)

	req := import1.RefreshApiCredentialsRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.RegisteredDomains.RefreshApiCredentials(ctx, &req)
	if err != nil {
		return diag.Errorf("error refreshing registered domain API credentials: %v", err)
	}

	taskRef, ok := resp.Data.GetValue().(multidomainPrism.TaskReference)
	if !ok {
		return diag.Errorf("refresh API credentials response did not contain task reference")
	}
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for registered domain refresh API credentials task (%s): %s", utils.StringValue(taskUUID), errWait)
	}

	d.SetId(extID)
	return ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Read(ctx, d, meta)
}

func ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRegisteredDomainRefreshApiCredentialsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
