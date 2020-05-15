package nutanix

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNutanixProtectionRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNutanixProtectionRuleRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"protection_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type: schema.TypeString,
						},
						"uuid": {
							Type: schema.TypeString,
						},
						"name": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type: schema.TypeString,
						},
						"uuid": {
							Type: schema.TypeString,
						},
						"name": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_connectivity_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_availability_zone_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_availability_zone_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"snapshot_schedule_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"recovery_point_objective_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"local_snapshot_retention_policy": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_snapshots": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"rollup_retention_policy_multiple": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"rollup_retention_policy_snapshot_interval_type": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"auto_suspend_timeout_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshot_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"remote_snapshot_retention_policy": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"num_snapshots": {
													Type: schema.TypeInt,
												},
												"rollup_retention_policy_multiple": {
													Type: schema.TypeInt,
												},
												"rollup_retention_policy_snapshot_interval_type": {
													Type: schema.TypeInt,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"ordered_availability_zone_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"category_filter": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"params": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      filterParamsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixProtectionRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	protectionRuleID := d.Get("protection_rule_id").(string)
	resp, err := conn.V3.GetProtectionRule(protectionRuleID)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("name", resp.Spec.Name); err != nil {
		return err
	}
	if err := d.Set("start_time", resp.Spec.Resources.StartTime); err != nil {
		return err
	}
	if err := d.Set("category_filter", flattenCategoriesFilter(resp.Spec.Resources.CategoryFilter)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_connectivity_list",
		flattenAvailabilityZoneConnectivityList(resp.Spec.Resources.AvailabilityZoneConnectivityList)); err != nil {
		return err
	}
	if err := d.Set("ordered_availability_zone_list",
		flattenOrderAvailibilityList(resp.Spec.Resources.OrderedAvailabilityZoneList)); err != nil {
		return err
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}
