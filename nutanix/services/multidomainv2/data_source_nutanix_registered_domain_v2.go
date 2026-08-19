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

func DatasourceNutanixRegisteredDomainV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRegisteredDomainV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credential_issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_credential_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registration_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registered_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"platform_data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaDatasourceRegisteredDomainPlatformData(),
			},
			"registration_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaDatasourceRegisteredDomainRegistrationConfig(),
			},
			"tunnel": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaDatasourceRegisteredDomainTunnel(),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
		},
	}
}

func schemaDatasourceRegisteredDomainPlatformData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connectivity_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pc_env_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pc_provider_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pc_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func schemaDatasourceRegisteredDomainRegistrationConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"key_id": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func schemaDatasourceRegisteredDomainTunnel() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixRegisteredDomainV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NCDomainManagementAPI

	extID := d.Get("ext_id").(string)
	getReq := import1.GetRegisteredDomainByIdRequest{
		ExtId: utils.StringPtr(extID),
	}
	resp, err := conn.RegisteredDomains.GetRegisteredDomainById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching registered domain: %v", err)
	}

	rd := resp.Data.GetValue().(config.RegisteredDomain)
	if err := setRegisteredDomainToResourceData(d, &rd); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(rd.ExtId))
	return nil
}
