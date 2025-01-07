package networking

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
)

func DataSourceNutanixFloatingIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFloatingIPRead,
		Schema: map[string]*schema.Schema{
			"floating_ip_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			// COMPUTED
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"external_subnet_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"floating_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vm_nic_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"vpc_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"execution_context": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_uuid": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"external_subnet_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"vm_nic_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"vpc_reference": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceNutanixFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	fUUID, ok := d.GetOk("floating_ip_uuid")
	if !ok {
		return diag.Errorf("please provide `floating_ip_uuid`")
	}

	resp, err := conn.V3.GetFloatingIPs(ctx, fUUID.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	m, _ := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", flattenStatusFIP(resp.Status)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("spec", flattenSpecFIP(resp.Spec)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*resp.Metadata.UUID)
	return nil
}

func flattenStatusFIP(stat *v3.FloatingIPDefStatus) []interface{} {
	statList := make([]interface{}, 0)

	if stat != nil {
		stats := make(map[string]interface{})
		stats["state"] = stat.State
		stats["resources"] = flattenFIPResource(stat.Resource)
		stats["execution_context"] = flattenExecutionContext(stat.ExecutionContext)

		statList = append(statList, stats)

		return statList
	}
	return nil
}

func flattenSpecFIP(spec *v3.FloatingIPSpec) []interface{} {
	specList := make([]interface{}, 0)

	if spec != nil {
		specs := make(map[string]interface{})

		specs["resources"] = flattenFIPResource(spec.Resource)

		specList = append(specList, specs)
		return specList
	}
	return nil
}

func flattenFIPResource(res *v3.FIPResource) []interface{} {
	resList := make([]interface{}, 0)
	if res != nil {
		ress := make(map[string]interface{})

		ress["external_subnet_reference"] = flattenReferenceValues(res.ExternalSubnetReference)
		if res.FloatingIP != nil {
			ress["floating_ip"] = res.FloatingIP
		}
		ress["vm_nic_reference"] = flattenReferenceValues(res.VMNICReference)
		ress["vpc_reference"] = flattenReferenceValues(res.VPCReference)

		resList = append(resList, ress)
		return resList
	}
	return nil
}
