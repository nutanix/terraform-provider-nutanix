package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixFloatingIPsList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFloatingIPsListRead,
		Schema: map[string]*schema.Schema{
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},

			// COMPUTED RESOURCES
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_matches": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"offset": {
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

func dataSourceNutanixFloatingIPsListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	req := &v3.DSMetadata{}
	length, lok := d.GetOk("length")
	if lok {
		req.Length = utils.Int64Ptr(int64(length.(int)))
	}

	offset, ok := d.GetOk("offset")
	if ok {
		req.Offset = utils.Int64Ptr(int64(offset.(int)))
	}

	resp, err := conn.V3.ListFloatingIPs(ctx, req)
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

		metList = append(metList, mets)
	}
	return metList
}
