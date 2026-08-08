package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import7 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/common"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/registereddomains"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixRegisteredDomainChangeStateV2 exposes the sync-registration
// action which updates the registration state of a registered domain.
func ResourceNutanixRegisteredDomainChangeStateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRegisteredDomainChangeStateV2Create,
		ReadContext:   ResourceNutanixRegisteredDomainChangeStateV2Read,
		DeleteContext: ResourceNutanixRegisteredDomainChangeStateV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"registration_state": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"UNREGISTERED", "REGISTERING", "REGISTERED", "UNREGISTERING",
					"REGISTRATION_ERROR", "UNREGISTRATION_ERROR", "DRAFT",
				}, false),
			},
		},
	}
}

func ResourceNutanixRegisteredDomainChangeStateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	extID := d.Get("ext_id").(string)

	spec := config.NewRegisteredDomainChangeStateSpec()
	spec.RegistrationState = commonUtils.ExpandEnum[import7.RegistrationState](d.Get("registration_state").(string))

	req := import1.ChangeStateRequest{
		ExtId: utils.StringPtr(extID),
		Body:  spec,
	}
	// sync-registration is a synchronous action: the response carries a list of
	// AppMessage entries rather than an async task reference.
	if _, err := conn.RegisteredDomains.ChangeState(ctx, &req); err != nil {
		return diag.Errorf("error changing registered domain state: %v", err)
	}

	d.SetId(extID)
	return ResourceNutanixRegisteredDomainChangeStateV2Read(ctx, d, meta)
}

func ResourceNutanixRegisteredDomainChangeStateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRegisteredDomainChangeStateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
