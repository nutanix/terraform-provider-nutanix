package dataprotectionv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/dataprotection-go-client/v4/models/dataprotection/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixRecoveryPointsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRecoveryPointsV2Read,
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
			"apply": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AOS",
			},
			"recovery_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixRecoveryPointV2(),
			},
		},
	}
}

func DatasourceNutanixRecoveryPointsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).DataProtectionAPI

	// initialize query params
	var filter, orderBy, selectQ *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectQy, ok := d.GetOk("select"); ok {
		selectQ = utils.StringPtr(selectQy.(string))
	} else {
		selectQ = nil
	}

	clusterID := d.Get("cluster_id").(string)

	resp, err := conn.RecoveryPoint.ListRecoveryPoints(&clusterID, page, limit, filter, orderBy, selectQ)
	if err != nil {
		return diag.Errorf("error while fetching Recovery Points : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("recovery_points", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of recovery points.",
		}}
	}

	getResp := resp.Data.GetValue().([]config.RecoveryPoint)

	if err := d.Set("recovery_points", flattenRecoveryPoints(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenRecoveryPoints(recoveryPoints []config.RecoveryPoint) []interface{} {
	if recoveryPoints == nil {
		return make([]interface{}, 0)
	}

	result := make([]interface{}, len(recoveryPoints))

	for i, recoveryPoint := range recoveryPoints {
		result[i] = map[string]interface{}{
			"ext_id":                       recoveryPoint.ExtId,
			"tenant_id":                    recoveryPoint.TenantId,
			"links":                        flattenLinks(recoveryPoint.Links),
			"location_agnostic_id":         recoveryPoint.LocationAgnosticId,
			"name":                         recoveryPoint.Name,
			"creation_time":                flattenTime(recoveryPoint.CreationTime),
			"expiration_time":              flattenTime(recoveryPoint.ExpirationTime),
			"status":                       flattenStatus(recoveryPoint.Status),
			"recovery_point_type":          flattenRecoveryPointType(recoveryPoint.RecoveryPointType),
			"owner_ext_id":                 recoveryPoint.OwnerExtId,
			"location_references":          flattenLocationReferences(recoveryPoint.LocationReferences),
			"vm_recovery_points":           flattenVMRecoveryPoints(recoveryPoint.VmRecoveryPoints),
			"volume_group_recovery_points": flattenVolumeGroupRecoveryPoints(recoveryPoint.VolumeGroupRecoveryPoints),
		}
	}
	return result
}
