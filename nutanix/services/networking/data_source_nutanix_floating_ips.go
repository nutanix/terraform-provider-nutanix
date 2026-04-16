package networking

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DataSourceNutanixFloatingIPs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFloatingIPsRead,
		Schema: map[string]*schema.Schema{
			// COMPUTED RESOURCES
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sort_order": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"offset": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sort_attribute": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"total_matches": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceNutanixFloatingIPsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).API

	req := &v3.DSMetadata{}

	metadata, filtersOk := d.GetOk("metadata")

	if filtersOk {
		req = buildDataSourceListMetadata(metadata.(*schema.Set))
	}

	resp, err := conn.V3.ListAllFloatingIPs(ctx, utils.StringValue(req.Filter))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("metadata", flattenFIPsMetadata(resp.Metadata)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("entities", flattenFIPsEntities(resp.Entities)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenFIPsEntities(ent []*v3.FloatingIPsIntentResponse) []map[string]interface{} {
	if len(ent) > 0 {
		entList := make([]map[string]interface{}, len(ent))
		for k, v := range ent {
			ents := make(map[string]interface{})

			ents["status"] = flattenStatusFIP(v.Status)
			ents["spec"] = flattenSpecFIP(v.Spec)
			m, _ := setRSEntityMetadata(v.Metadata)
			ents["metadata"] = m

			entList[k] = ents
		}
		return entList
	}
	return nil
}

func flattenFIPsMetadata(met *v3.ListMetadataOutput) []interface{} {
	metList := make([]interface{}, 0)

	if met != nil {
		mets := make(map[string]interface{})

		mets["total_matches"] = met.TotalMatches
		mets["kind"] = met.Kind
		mets["length"] = met.Length
		mets["offset"] = met.Offset
		mets["filter"] = met.Filter

		metList = append(metList, mets)
	}
	return metList
}
