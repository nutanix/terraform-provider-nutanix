package tenancyv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/config"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/tenancy-go-client/v17/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixAddDomainV2 adds an allowed domain to Nutanix Central tenancy.
// The tenancy API does not currently provide update or delete operations for domains.
func ResourceNutanixAddDomainV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixAddDomainV2Create,
		ReadContext:   ResourceNutanixAddDomainV2Read,
		DeleteContext: ResourceNutanixAddDomainV2Delete,
		Schema: map[string]*schema.Schema{
			"domain_ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
		},
	}
}

// ResourceNutanixAddDomainV2Create adds an allowed domain.
func ResourceNutanixAddDomainV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	body := import1.Domain{}
	if v, ok := d.GetOk("domain_ext_id"); ok {
		body.DomainExtId = utils.StringPtr(v.(string))
	}

	request := import4.AddDomainRequest{
		Body: &body,
	}
	resp, err := conn.TenantManagementAPIInstance.AddDomain(ctx, &request)
	if err != nil {
		return diag.Errorf("error while adding domain: %v", err)
	}
	if resp == nil || resp.Data == nil {
		return diag.Errorf("error while adding domain: API returned an empty response")
	}

	domain, ok := resp.Data.GetValue().(import1.Domain)
	if !ok {
		return diag.Errorf("error while adding domain: unexpected response type %T", resp.Data.GetValue())
	}
	if domain.ExtId == nil || *domain.ExtId == "" {
		return diag.Errorf("error while adding domain: response did not include an ext_id")
	}

	d.SetId(utils.StringValue(domain.ExtId))
	return ResourceNutanixAddDomainV2Read(ctx, d, meta)
}

// ResourceNutanixAddDomainV2Read refreshes domain state from the list endpoint.
func ResourceNutanixAddDomainV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ServiceProviderManagementAPI

	resp, err := conn.TenantManagementAPIInstance.ListDomains(ctx, &import4.ListDomainsRequest{})
	if err != nil {
		return diag.Errorf("error while fetching domains: %v", err)
	}
	if resp == nil || resp.Data == nil {
		return diag.Errorf("error while fetching domain %q: API returned an empty response", d.Id())
	}

	domains, ok := resp.Data.GetValue().([]import1.Domain)
	if !ok {
		return diag.Errorf("error while fetching domains: unexpected response type %T", resp.Data.GetValue())
	}
	for _, domain := range domains {
		if utils.StringValue(domain.ExtId) == d.Id() {
			return setDomainResourceData(d, domain)
		}
	}

	return diag.Errorf("error while fetching domain %q: domain was not returned by the list API", d.Id())
}

// ResourceNutanixAddDomainV2Delete removes the domain resource from Terraform state.
// The tenancy API does not expose a corresponding remote delete operation.
func ResourceNutanixAddDomainV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func setDomainResourceData(d *schema.ResourceData, domain import1.Domain) diag.Diagnostics {
	if err := d.Set("domain_ext_id", domain.DomainExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ext_id", domain.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", domain.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(domain.Links)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
