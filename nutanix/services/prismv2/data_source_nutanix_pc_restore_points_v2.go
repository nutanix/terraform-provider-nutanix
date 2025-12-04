package prismv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixFetchRestorePointsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixRestorePointsV2Read,
		Schema: map[string]*schema.Schema{
			"restore_source_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restorable_domain_manager_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50, //nolint:gomnd
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixFetchRestorePointV2(),
			},
		},
	}
}

func DatasourceNutanixRestorePointsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	restoreSourceExtID := utils.StringPtr(d.Get("restore_source_ext_id").(string))
	restorableDomainManagerExtID := utils.StringPtr(d.Get("restorable_domain_manager_ext_id").(string))

	var selects, orderBy, filter *string
	var page, limit *int

	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}
	if orderByf, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(orderByf.(string))
	} else {
		orderBy = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
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

	resp, err := conn.DomainManagerBackupsAPIInstance.ListRestorePoints(restoreSourceExtID, restorableDomainManagerExtID, page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching Domain Manager Restore Point Detail: %s", err)
	}

	if resp.Data == nil {
		if err := d.Set("restore_points", make([]interface{}, 0)); err != nil {
			return diag.Errorf("Error setting restore_points: %v", err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of restore points.",
		}}
	}

	restorePoints := resp.Data.GetValue().([]management.RestorePoint)

	if err := d.Set("restore_points", flattenRestorePoints(restorePoints)); err != nil {
		return diag.Errorf("Error setting restore_points: %v", err)
	}

	d.SetId(utils.GenUUID())

	return nil
}

func flattenRestorePoints(restorePoints []management.RestorePoint) []map[string]interface{} {
	restorePointsList := make([]map[string]interface{}, 0)
	for _, restorePoint := range restorePoints {
		restorePointMap := map[string]interface{}{
			"tenant_id":      utils.StringValue(restorePoint.TenantId),
			"ext_id":         utils.StringValue(restorePoint.ExtId),
			"links":          flattenLinks(restorePoint.Links),
			"creation_time":  flattenTime(restorePoint.CreationTime),
			"domain_manager": flattenDomainManager(restorePoint.DomainManager),
		}
		restorePointsList = append(restorePointsList, restorePointMap)
	}
	return restorePointsList
}
