package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/registereddomains"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixRegisteredDomainEstablishTrustV2 exposes the establish-trust
// action which verifies a registered domain by exchanging a client-signed
// request with Nutanix Central.
func ResourceNutanixRegisteredDomainEstablishTrustV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRegisteredDomainEstablishTrustV2Create,
		ReadContext:   ResourceNutanixRegisteredDomainEstablishTrustV2Read,
		DeleteContext: ResourceNutanixRegisteredDomainEstablishTrustV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_signing_request": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func ResourceNutanixRegisteredDomainEstablishTrustV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	extID := d.Get("ext_id").(string)

	spec := config.NewRegisteredDomainTrustSpec()
	spec.ClientSigningRequest = utils.StringPtr(d.Get("client_signing_request").(string))

	req := import1.EstablishTrustRequest{
		ExtId: utils.StringPtr(extID),
		Body:  spec,
	}
	// establish-trust is a synchronous action: the response carries a list of
	// AppMessage entries rather than an async task reference.
	if _, err := conn.RegisteredDomains.EstablishTrust(ctx, &req); err != nil {
		return diag.Errorf("error establishing trust for registered domain: %v", err)
	}

	d.SetId(extID)
	return ResourceNutanixRegisteredDomainEstablishTrustV2Read(ctx, d, meta)
}

func ResourceNutanixRegisteredDomainEstablishTrustV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRegisteredDomainEstablishTrustV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
