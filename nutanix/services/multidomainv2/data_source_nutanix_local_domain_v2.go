package multidomainv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/management"
	localdomainReq "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/request/localdomain"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLocalDomainV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLocalDomainV2Read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaLocalDomainConnectionStatus(),
			},
			"registration_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaLocalDomainRegistrationStatus(),
			},
			"location": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaLocalDomainLocation(),
			},
		},
	}
}

func DatasourceNutanixLocalDomainV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MultidomainAPI

	getReq := localdomainReq.GetLocalDomainRequest{}
	resp, err := conn.LocalDomain.GetLocalDomain(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while fetching LocalDomain: %v", err)
	}

	domain := resp.Data.GetValue().(management.LocalDomain)
	if err := d.Set("name", utils.StringValue(domain.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_url", utils.StringValue(domain.TargetUrl)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connection_status", flattenLocalDomainConnectionStatus(domain.ConnectionStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("registration_status", flattenLocalDomainRegistrationStatus(domain.RegistrationStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location", flattenLocalDomainLocation(domain.Location)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
