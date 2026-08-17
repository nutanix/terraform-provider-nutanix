package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixScheduledDeleteTenantV2 schedules deletion of a tenant after a
// configurable retention period. This is a standalone mutating action resource.
func ResourceNutanixScheduledDeleteTenantV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixScheduledDeleteTenantV2Create,
		ReadContext:   ResourceNutanixScheduledDeleteTenantV2Read,
		DeleteContext: ResourceNutanixScheduledDeleteTenantV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_seconds": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// ResourceNutanixScheduledDeleteTenantV2Create issues the scheduled-delete action for a tenant.
func ResourceNutanixScheduledDeleteTenantV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	extID := d.Get("ext_id").(string)
	body := import1.TenantScheduledDeleteSpec{
		RetentionSeconds: utils.Int64Ptr(int64(d.Get("retention_seconds").(int))),
	}
	scheduledDeleteRequest := import4.ScheduledDeleteTenantByIdRequest{
		ExtId: utils.StringPtr(extID),
		Body:  &body,
	}
	_, err := conn.TenantManagementAPIInstance.ScheduledDeleteTenantById(ctx, &scheduledDeleteRequest)
	if err != nil {
		return diag.Errorf("error while scheduling deletion of tenant: %v", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// ResourceNutanixScheduledDeleteTenantV2Read is a no-op; the scheduled-delete action has no read API.
func ResourceNutanixScheduledDeleteTenantV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixScheduledDeleteTenantV2Delete removes the action resource from state.
func ResourceNutanixScheduledDeleteTenantV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
