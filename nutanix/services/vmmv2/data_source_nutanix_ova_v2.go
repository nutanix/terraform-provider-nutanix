package vmmv2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/content"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixOvaV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNutanixOvaV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ova_sha1_checksum": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ova_sha256_checksum": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hex_digest": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ova_url_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"basic_auth": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"username": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"password": {
													Type:      schema.TypeString,
													Computed:  true,
													Sensitive: true,
												},
											},
										},
									},
								},
							},
						},
						"ova_vm_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_ext_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_file_format": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"object_lite_source": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"created_by": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idp_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"middle_initial": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locale": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_force_reset_password_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"additional_attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": schemaForValue(),
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buckets_access_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ext_id": {
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
									"access_key_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"secret_access_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"user_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"created_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"last_login_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster_location_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parent_vm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ResourceNutanixVirtualMachineV2(),
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// datasourceNutanixOvaV2Read implements the ReadContext function for the OVA data source.
func datasourceNutanixOvaV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id").(string)
	resp, err := conn.OvasAPIInstance.GetOvaById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error reading OVA: %v", err)
	}

	ova := resp.Data.GetValue().(import1.Ova)
	diags := setOva(d, ova)
	if diags.HasError() {
		return diags
	}
	d.SetId(*ova.ExtId)
	return nil
}

func setOva(d *schema.ResourceData, ova import1.Ova) diag.Diagnostics {
	if err := d.Set("ext_id", utils.StringValue(ova.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenAPILink(ova.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(ova.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", utils.StringValue(ova.Name)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("checksum", flattenOneOfOvaChecksum(ova.Checksum)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("size_bytes", int(*ova.SizeBytes)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by", flattenCreatedBy(ova.CreatedBy)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster_location_ext_ids", ova.ClusterLocationExtIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("parent_vm", utils.StringValue(ova.ParentVm)); err != nil {
		return diag.FromErr(err)
	}

	// Set the VM config
	fields, diags := extractVMConfigFields(*ova.VmConfig)
	if diags.HasError() {
		return diags
	}
	if err := d.Set("vm_config", []interface{}{fields}); err != nil {
		return diag.FromErr(fmt.Errorf("failed setting vm_config: %w", err))
	}

	if err := d.Set("disk_format", flattenOvaDiskFormat(ova.DiskFormat)); err != nil {
		return diag.FromErr(err)
	}

	if ova.CreateTime != nil {
		t := ova.CreateTime
		if err := d.Set("create_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	if ova.LastUpdateTime != nil {
		t := ova.LastUpdateTime
		if err := d.Set("last_update_time", t.String()); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
