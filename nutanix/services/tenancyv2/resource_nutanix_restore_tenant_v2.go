package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixRestoreTenantV2 restores a suspended or pending-deletion tenant
// back to the active state. This is a standalone mutating action resource.
func ResourceNutanixRestoreTenantV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRestoreTenantV2Create,
		ReadContext:   ResourceNutanixRestoreTenantV2Read,
		DeleteContext: ResourceNutanixRestoreTenantV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// ResourceNutanixRestoreTenantV2Create issues the restore action for a tenant.
func ResourceNutanixRestoreTenantV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	extID := d.Get("ext_id").(string)
	restoreRequest := import4.RestoreTenantByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	_, err := conn.TenantManagementAPIInstance.RestoreTenantById(ctx, &restoreRequest)
	if err != nil {
		return diag.Errorf("error while restoring tenant: %v", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// ResourceNutanixRestoreTenantV2Read is a no-op; the restore action has no read API.
func ResourceNutanixRestoreTenantV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixRestoreTenantV2Delete removes the action resource from state.
func ResourceNutanixRestoreTenantV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
