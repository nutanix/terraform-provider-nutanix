package prismv2

import (
	"context"
	"encoding/json"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
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
						"config":  schemaForPcConfig(),
						"network": schemaForPcNetwork(),
						"should_enable_high_availability": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
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
			"is_registered_with_hosting_cluster": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"hosting_cluster_ext_id": {
				Type:     schema.TypeString,
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
	}
}

func ResourceNutanixRestorePcCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).PrismAPI

	restoreSourceExtID := d.Get("restore_source_ext_id").(string)
	restoreableDomainManagerExtID := d.Get("restorable_domain_manager_ext_id").(string)
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

	resp, err := conn.DomainManagerBackupsAPIInstance.Restore(utils.StringPtr(restoreSourceExtID), utils.StringPtr(restoreableDomainManagerExtID), utils.StringPtr(restorePointExtID), restoreSpec)
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

	resp, err := conn.DomainManagerAPIInstance.GetDomainManagerById(utils.StringPtr(d.Id()))

	if err != nil {
		return diag.Errorf("error while fetching Domain Manager: %s", err)
	}

	deployPcBody := resp.Data.GetValue().(config.DomainManager)

	if err := d.Set("tenant_id", utils.StringValue(deployPcBody.TenantId)); err != nil {
		return diag.Errorf("error setting tenant_id: %s", err)
	}
	if err := d.Set("ext_id", utils.StringValue(deployPcBody.ExtId)); err != nil {
		return diag.Errorf("error setting ext_id: %s", err)
	}
	if err := d.Set("links", flattenLinks(deployPcBody.Links)); err != nil {
		return diag.Errorf("error setting links: %s", err)
	}
	if err := d.Set("config", flattenPCConfig(deployPcBody.Config)); err != nil {
		return diag.Errorf("error setting config: %s", err)
	}
	if err := d.Set("is_registered_with_hosting_cluster", utils.BoolValue(deployPcBody.IsRegisteredWithHostingCluster)); err != nil {
		return diag.Errorf("error setting is_registered_with_hosting_cluster: %s", err)
	}
	if err := d.Set("network", flattenPCNetwork(deployPcBody.Network)); err != nil {
		return diag.Errorf("error setting network: %s", err)
	}
	if err := d.Set("hosting_cluster_ext_id", utils.StringValue(deployPcBody.HostingClusterExtId)); err != nil {
		return diag.Errorf("error setting hosting_cluster_ext_id: %s", err)
	}
	if err := d.Set("should_enable_high_availability", utils.BoolValue(deployPcBody.ShouldEnableHighAvailability)); err != nil {
		return diag.Errorf("error setting should_enable_high_availability: %s", err)
	}
	if err := d.Set("node_ext_ids", deployPcBody.NodeExtIds); err != nil {
		return diag.Errorf("error setting node_ext_ids: %s", err)
	}
	return nil
}

func ResourceNutanixRestorePcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRestorePcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
