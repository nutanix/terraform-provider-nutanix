package dataprotectionv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v17/models/dataprotection/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/dataprotection-go-client/v17/models/dataprotection/v4/request/recoverypoints"
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

	clusterID := d.Get("cluster_id").(string)

	listRecoveryPointsRequest := import1.ListRecoveryPointsRequest{
		XClusterId: utils.StringPtr(clusterID),
	}

	if v, ok := d.GetOk("page"); ok {
		listRecoveryPointsRequest.Page_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("limit"); ok {
		listRecoveryPointsRequest.Limit_ = utils.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("filter"); ok {
		listRecoveryPointsRequest.Filter_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("order_by"); ok {
		listRecoveryPointsRequest.Orderby_ = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("select"); ok {
		listRecoveryPointsRequest.Select_ = utils.StringPtr(v.(string))
	}

	resp, err := conn.RecoveryPoint.ListRecoveryPoints(ctx, &listRecoveryPointsRequest)
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
			"project_ext_id":               recoveryPoint.ProjectExtId,
		}
	}
	return result
}
