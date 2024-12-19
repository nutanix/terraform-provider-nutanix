package prismv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"

	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixRestorePcV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRestorePcCreate,
		ReadContext:   ResourceNutanixRestorePcRead,
		UpdateContext: ResourceNutanixRestorePcUpdate,
		DeleteContext: ResourceNutanixRestorePcDelete,
		Schema: map[string]*schema.Schema{
			"restore_source_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"restorable_domain_manager_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_manager": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config": schemaForPcConfig(),
						"is_registered_with_hosting_cluster": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network": schemaForPcNetwork(),
						"hosting_cluster_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"should_enable_high_availability": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"node_ext_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			// read schema
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixRestorePcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	restoreSourceExtID := d.Get("restore_source_ext_id").(string)
	restorableDomainManagerExtID := d.Get("restorable_domain_manager_ext_id").(string)
	restorePointExtID := d.Get("ext_id").(string)

	restorePcBody := config.NewDomainManager()

	domainManagerConfigI := d.Get("domain_manager").([]interface{})[0]
	domainManagerConfig := domainManagerConfigI.(map[string]interface{})

	restorePcBody.Config = expandPCConfig(domainManagerConfig["config"].(map[string]interface{}))
	restorePcBody.Network = expandPCNetwork(domainManagerConfig["network"].(map[string]interface{}))
	restorePcBody.ShouldEnableHighAvailability = utils.BoolPtr(domainManagerConfig["should_enable_high_availability"].(bool))

	restoreSpec := management.NewRestoreSpec()

	restoreSpec.DomainManager = restorePcBody

	aJSON, _ := json.MarshalIndent(restoreSpec, "", "  ")
	log.Printf("[DEBUG] Restore PC Body: %s", string(aJSON))

	resp, err := conn.DomainManagerBackupsAPIInstance.Restore(utils.StringPtr(restoreSourceExtID), utils.StringPtr(restorableDomainManagerExtID), utils.StringPtr(restorePointExtID), restoreSpec)
	if err != nil {
		return diag.Errorf("error while restoring Domain Manager: %s", err)
	}

	TaskRef := resp.Data.GetValue().(config.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for Restore Domain Manager Task to complete: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching task details : %v", err)
	}

	rUUID := resourceUUID.Data.GetValue().(config.Task)
	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Restore Domain Manager Task Details: %s", string(aJSON))

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	return ResourceNutanixRestorePcRead(ctx, d, meta)
}

func ResourceNutanixRestorePcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	restoreSourceExtID := utils.StringPtr(d.Get("restore_source_ext_id").(string))
	restorableDomainManagerExtID := utils.StringPtr(d.Get("restorable_domain_manager_ext_id").(string))

	resp, err := conn.DomainManagerBackupsAPIInstance.GetRestorePointById(restoreSourceExtID, restorableDomainManagerExtID, utils.StringPtr(d.Id()))

	if err != nil {
		return diag.Errorf("error while fetching Domain Manager Restore Point Detail: %s", err)
	}

	restorePoint := resp.Data.GetValue().(management.RestorePoint)

	if err := d.Set("tenant_id", utils.StringValue(restorePoint.TenantId)); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("ext_id", utils.StringValue(restorePoint.ExtId)); err != nil {
		return diag.Errorf("error setting ext_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(restorePoint.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("creation_time", flattenTime(restorePoint.CreationTime)); err != nil {
		return diag.Errorf("error setting creation_time: %s", err)
	}
	if err := d.Set("domain_manager", flattenDomainManager(restorePoint.DomainManager)); err != nil {
		return diag.Errorf("error setting domain_manager: %s", err)
	}
	return nil
}

func ResourceNutanixRestorePcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRestorePcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
