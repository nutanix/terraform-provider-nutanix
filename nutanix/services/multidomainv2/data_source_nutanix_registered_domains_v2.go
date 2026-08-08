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

func DatasourceNutanixRegisteredDomainsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRegisteredDomainsV2Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expand": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"registered_domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixRegisteredDomainV2(),
			},
		},
	}
}

func DatasourceNutanixRegisteredDomainsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	listReq := import1.ListRegisteredDomainsRequest{}
	if v, ok := d.GetOk("page"); ok {
		listReq.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listReq.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listReq.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listReq.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("expand"); ok {
		listReq.Expand_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listReq.Select_ = utils.StringPtr(v.(string))
	}
	resp, err := conn.RegisteredDomains.ListRegisteredDomains(ctx, &listReq)
	if err != nil {
		return diag.Errorf("error while listing registered domains: %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("registered_domains", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "No data found.",
			Detail:   "The API returned an empty list of registered domains.",
		}}
	}

	registeredDomains := resp.Data.GetValue().([]config.RegisteredDomain)
	if err := d.Set("registered_domains", flattenRegisteredDomains(registeredDomains)); err != nil {
		return diag.Errorf("error setting registered_domains: %v", err)
	}

	d.SetId(utils.GenUUID())
	return nil
}

func flattenRegisteredDomains(registeredDomains []config.RegisteredDomain) []map[string]interface{} {
	if len(registeredDomains) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(registeredDomains))
	for _, rd := range registeredDomains {
		m := map[string]interface{}{
			"ext_id":                utils.StringValue(rd.ExtId),
			"name":                  utils.StringValue(rd.Name),
			"dns_name":              utils.StringValue(rd.DnsName),
			"fqdn":                  utils.StringValue(rd.Fqdn),
			"domain_ext_id":         utils.StringValue(rd.DomainExtId),
			"location_ext_id":       utils.StringValue(rd.LocationExtId),
			"credential_issuer":     utils.StringValue(rd.CredentialIssuer),
			"created_by":            utils.StringValue(rd.CreatedBy),
			"api_credential_status": flattenApiCredentialStatus(rd.ApiCredentialStatus),
			"registration_state":    flattenRegistrationState(rd.RegistrationState),
			"registered_time":       utils.TimeStringValue(rd.RegisteredTime),
			"platform_data":         flattenRegisteredDomainPlatformData(rd.PlatformData),
			"registration_config":   flattenRegisteredDomainRegistrationConfig(rd.RegistrationConfig),
			"tunnel":                flattenRegisteredDomainTunnel(rd.Tunnel),
			"tenant_id":             utils.StringValue(rd.TenantId),
			"links":                 flattenLinks(rd.Links),
		}
		out = append(out, m)
	}
	return out
}
