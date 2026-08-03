package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/request/vpcs"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixVPCsv2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixVPCsv2Read,
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
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DataSourceNutanixVPCv2(),
			},
		},
	}
}

func dataSourceNutanixVPCsv2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	listVpcsRequest := import2.ListVpcsRequest{}

	if v, ok := d.GetOk("page"); ok {
		listVpcsRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listVpcsRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listVpcsRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listVpcsRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listVpcsRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.VpcAPIInstance.ListVpcs(ctx, &listVpcsRequest)
	if err != nil {
		return diag.Errorf("error while fetching vpcs : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("vpcs", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of VPCs.",
		}}
	}

	getResp := resp.Data.GetValue().([]import1.Vpc)

	if err := d.Set("vpcs", flattenVPCsEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVPCsEntities(pr []import1.Vpc) []map[string]interface{} {
	if len(pr) > 0 {
		vpcs := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			vpc := make(map[string]interface{})

			vpc["tenant_id"] = utils.StringValue(v.TenantId)
			vpc["ext_id"] = utils.StringValue(v.ExtId)
			vpc["links"] = flattenLinks(v.Links)
			vpc["metadata"] = flattenMetadata(v.Metadata)
			vpc["name"] = utils.StringValue(v.Name)
			vpc["description"] = utils.StringValue(v.Description)
			vpc["common_dhcp_options"] = flattenCommonDhcpOptions(v.CommonDhcpOptions)
			vpc["vpc_type"] = v.VpcType.GetName()
			vpc["snat_ips"] = flattenNtpServer(v.SnatIps)
			vpc["external_subnets"] = flattenExternalSubnets(v.ExternalSubnets)
			vpc["external_routing_domain_reference"] = v.ExternalRoutingDomainReference
			vpc["externally_routable_prefixes"] = flattenExternallyRoutablePrefixes(v.ExternallyRoutablePrefixes)
			vpc["project_ext_id"] = v.ProjectExtId
			vpc["shared_with_projects"] = v.SharedWithProjects
			vpc["should_advertise_connected_subnets"] = v.ShouldAdvertiseConnectedSubnets
			if v.SupportedMultipleExternalSubnetType != nil {
				vpc["supported_multiple_external_subnet_type"] = flattenSupportedMultipleExternalSubnetType(v.SupportedMultipleExternalSubnetType)
			}
			if v.Scope != nil {
				vpc["scope"] = flattenVpcScope(v.Scope)
			}
			vpc["kubernetes_clusters"] = flattenKubernetesClusters(v.KubernetesClusters)

			vpcs[k] = vpc
		}
		return vpcs
	}
	return nil
}
