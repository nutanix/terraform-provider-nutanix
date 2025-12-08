package networkingv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixVPCv2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixVPCv2Read,
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
					Schema: DatasourceMetadataSchemaV2(),
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"common_dhcp_options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
					},
				},
			},
			"vpc_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"snat_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": SchemaForValuePrefixLength(),
						"ipv6": SchemaForValuePrefixLength(),
					},
				},
			},
			"external_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_reference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4": SchemaForValuePrefixLength(),
									"ipv6": SchemaForValuePrefixLength(),
								},
							},
						},
						"gateway_nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"active_gateway_node": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_ip_address": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4": SchemaForValuePrefixLength(),
												"ipv6": SchemaForValuePrefixLength(),
											},
										},
									},
								},
							},
						},
						"active_gateway_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"external_routing_domain_reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"externally_routable_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": SchemaForValuePrefixLength(),
									"prefix_length": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"ipv6": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": SchemaForValuePrefixLength(),
									"prefix_length": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixVPCv2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	extID := d.Get("ext_id")
	resp, err := conn.VpcAPIInstance.GetVpcById(utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching vpc : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Vpc)

	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpc_type", getResp.VpcType.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata", flattenMetadata(getResp.Metadata)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("common_dhcp_options", flattenCommonDhcpOptions(getResp.CommonDhcpOptions)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("snat_ips", flattenNtpServer(getResp.SnatIps)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("external_subnets", flattenExternalSubnets(getResp.ExternalSubnets)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("external_routing_domain_reference", getResp.ExternalRoutingDomainReference); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("externally_routable_prefixes", flattenExternallyRoutablePrefixes(getResp.ExternallyRoutablePrefixes)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}
