package lcmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixLcmConfigrationV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixLcmConfigrationV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"links": schemaForLinks(),
		},
	}
}

func DatasourceNutanixLcmConfigrationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	domainManagerExtID := d.Get("domain_manager_ext_id").(string)
	backupTargetExtID := d.Get("ext_id").(string)

	resp, err := conn.DomainManagerBackupsAPIInstance.GetBackupTargetById(utils.StringPtr(domainManagerExtID), utils.StringPtr(backupTargetExtID), nil)

	if err != nil {
		return diag.Errorf("error while fetching Backup Target: %s", err)
	}

	backupTarget := resp.Data.GetValue().(management.BackupTarget)

	if err := d.Set("tenant_id", backupTarget.TenantId); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("ext_id", backupTarget.ExtId); err != nil {
		return diag.Errorf("error setting ext_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(backupTarget.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("last_sync_time", flattenTime(backupTarget.LastSyncTime)); err != nil {
		return diag.Errorf("error setting last_sync_time: %s", err)
	}
	if err := d.Set("is_backup_paused", backupTarget.IsBackupPaused); err != nil {
		return diag.Errorf("error setting is_backup_paused: %s", err)
	}
	if err := d.Set("backup_pause_reason", backupTarget.BackupPauseReason); err != nil {
		return diag.Errorf("error setting backup_pause_reason: %s", err)
	}
	if err := d.Set("location", flattenBackupTargetLocation(backupTarget.Location)); err != nil {
		return diag.Errorf("error setting location: %s", err)
	}

	return nil
}

func schemaForLinks() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rel": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"href": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
