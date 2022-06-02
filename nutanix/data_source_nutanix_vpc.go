package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNutanixVPC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixVPCRead,
		Schema: map[string]*schema.Schema{
			"vpc_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},

			// COMPUTED RESOURCES
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"categories": categoriesSchema(),
			"external_subnet_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_subnet_reference": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"externally_routable_prefix_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prefix_length": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"common_domain_name_server_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_backup": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"external_subnet_list_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_subnet_reference": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"external_ip_list": {
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
									"host_reference": {
										Type:     schema.TypeMap,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"ip_address": {
										Type:     schema.TypeString,
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

func dataSourceNutanixVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	vpcUUID, vpcOk := d.GetOk("vpc_uuid")

	if !vpcOk {
		return diag.Errorf("please provide `vpc_uuid`")
	}

	resp, err := conn.V3.GetVPC(ctx, vpcUUID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", resp.Spec.Name); err != nil {
		return diag.FromErr(err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("categories", c); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("external_subnet_list", flattenExtSubnetList(resp.Spec.Resources.ExternalSubnetList)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("externally_routable_prefix_list", flattenExtRoutableList(resp.Spec.Resources.ExternallyRoutablePrefixList)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("common_domain_name_server_ip_list", flattenCommonDNSIPList(resp.Spec.Resources.CommonDomainNameServerIPList)); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("external_subnet_list_status", flattenExtSubnetListStatus(resp.Status.Resources.ExternalSubnetList)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}
