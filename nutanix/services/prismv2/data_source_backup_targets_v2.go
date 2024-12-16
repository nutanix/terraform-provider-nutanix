package prismv2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixBackupTargetsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixBackupTargetsV2Read,
		Schema: map[string]*schema.Schema{
			"domain_manager_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_targets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixBackupTargetV2(),
			},
		},
	}
}

func DatasourceNutanixBackupTargetsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	domainManagerExtID := d.Get("domain_manager_ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.ListBackupTargets(utils.StringPtr(domainManagerExtID))

	if err != nil {
		return diag.Errorf("error while Listing Backup Targets for : %s err: %s", domainManagerExtID, err)
	}

	if resp.Data == nil {
		if err := d.Set("backup_targets", []interface{}{}); err != nil {
			return diag.Errorf("error setting backup_targets: %s", err)
		}
	}

	if err := d.Set("backup_targets", flattenBackupTargets(resp.Data.GetValue().([]management.BackupTarget))); err != nil {
		return diag.Errorf("error setting backup_targets: %s", err)
	}
	return nil
}

func flattenBackupTargets(backupTargets []management.BackupTarget) []map[string]interface{} {
	if backupTargets == nil || len(backupTargets) == 0 {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0)
	for _, backupTarget := range backupTargets {
		backupTargetMap := map[string]interface{}{
			"ext_id":              backupTarget.ExtId,
			"tenant_id":           backupTarget.TenantId,
			"links":               flattenLinks(backupTarget.Links),
			"location":            flattenBackupTargetLocation(backupTarget.Location),
			"last_sync_time":      flattenTime(backupTarget.LastSyncTime),
			"is_backup_paused":    backupTarget.IsBackupPaused,
			"backup_pause_reason": backupTarget.BackupPauseReason,
		}
		result = append(result, backupTargetMap)
	}
	return result
}
