package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRouteTableV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRouteTableV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: DatasourceMetadataSchemaV4(),
				},
			},
			"vpc_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_routing_domain_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixRouteTableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	extID := d.Get("ext_id")
	resp, err := conn.RoutesTable.GetRouteTableById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching route table : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.RouteTable)

	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_reference", getResp.VpcReference); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("external_routing_domain_reference", getResp.ExternalRoutingDomainReference); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
