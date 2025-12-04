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
				Elem:     schemaForBackupTarget(),
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
		if err := d.Set("backup_targets", make([]interface{}, 0)); err != nil {
			return diag.Errorf("error setting backup_targets: %s", err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of backup targets.",
		}}
	}

	getResp := resp.Data.GetValue().([]management.BackupTarget)

	if err := d.Set("backup_targets", flattenBackupTargets(getResp)); err != nil {
		return diag.Errorf("error setting backup_targets: %s", err)
	}

	d.SetId(domainManagerExtID)

	return nil
}

func flattenBackupTargets(backupTargets []management.BackupTarget) []map[string]interface{} {
	if len(backupTargets) == 0 {
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

// schema for backup target
func schemaForBackupTarget() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"location": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_location": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ext_id": {
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
								},
							},
						},
						"object_store_location": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"credentials": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"access_key_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"secret_access_key": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"backup_policy": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rpo_in_minutes": {
													Type:     schema.TypeInt,
													Computed: true,
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
			"last_sync_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_backup_paused": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backup_pause_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
