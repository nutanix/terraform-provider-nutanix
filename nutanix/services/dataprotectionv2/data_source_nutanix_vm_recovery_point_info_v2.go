package dataprotectionv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVMRecoveryPointInfoV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVMRecoveryPointInfoV2Read,
		Schema: map[string]*schema.Schema{
			"recovery_point_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": SchemaForLinks(),
			"consistency_group_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_agnostic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_recovery_points": SchemaForDiskRecoveryPoints(),
			"vm_ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"application_consistent_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"should_include_writers": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"writers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"should_store_vss_metadata": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"object_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixVMRecoveryPointInfoV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DatasourceNutanixVMRecoveryPointInfoV2Read \n")

	conn := meta.(*conns.Client).DataProtectionAPI

	recoveryPointExtID := d.Get("recovery_point_ext_id").(string)
	extID := d.Get("ext_id").(string)

	resp, err := conn.RecoveryPoint.GetVmRecoveryPointById(&recoveryPointExtID, &extID)
	if err != nil {
		return diag.Errorf("error while fetching vm recovery point: %v", err)
	}

	getResp := resp.Data.GetValue().(config.VmRecoveryPoint)

	aJSON, _ := json.Marshal(getResp)
	log.Printf("[DEBUG] DatasourceNutanixVMRecoveryPointInfoV2Read response: \n%v\n", string(aJSON))

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("consistency_group_ext_id", getResp.ConsistencyGroupExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("location_agnostic_id", getResp.LocationAgnosticId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_recovery_points", flattenDiskRecoveryPoints(getResp.DiskRecoveryPoints)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_ext_id", getResp.VmExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_categories", getResp.VmCategories); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("application_consistent_properties", flattenApplicationConsistentProperties(getResp.ApplicationConsistentProperties)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(getResp.ExtId))
	return nil
}
