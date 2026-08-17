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

// ResourceNutanixSuspendTenantV2 suspends a tenant, temporarily blocking all of
// its users without losing data. This is a standalone mutating action resource.
func ResourceNutanixSuspendTenantV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSuspendTenantV2Create,
		ReadContext:   ResourceNutanixSuspendTenantV2Read,
		DeleteContext: ResourceNutanixSuspendTenantV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// ResourceNutanixSuspendTenantV2Create issues the suspend action for a tenant.
func ResourceNutanixSuspendTenantV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	extID := d.Get("ext_id").(string)
	suspendRequest := import4.SuspendTenantByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	_, err := conn.TenantManagementAPIInstance.SuspendTenantById(ctx, &suspendRequest)
	if err != nil {
		return diag.Errorf("error while suspending tenant: %v", err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

// ResourceNutanixSuspendTenantV2Read is a no-op; the suspend action has no read API.
func ResourceNutanixSuspendTenantV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixSuspendTenantV2Delete removes the action resource from state.
func ResourceNutanixSuspendTenantV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
