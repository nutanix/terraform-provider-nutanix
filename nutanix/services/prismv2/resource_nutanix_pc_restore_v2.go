package prismv2

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1*time.Hour + 30*time.Minute),
		},
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
						"config": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     schemaForPcConfig(),
						},
						"is_registered_with_hosting_cluster": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     schemaForPcNetwork(),
						},
						"hosting_cluster_ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"should_enable_high_availability": {
							Type:     schema.TypeBool,
							Optional: true,
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

	restorePcBody.Config = expandPCConfig(domainManagerConfig["config"])
	restorePcBody.Network = expandPCNetwork(domainManagerConfig["network"])
	restorePcBody.ShouldEnableHighAvailability = utils.BoolPtr(domainManagerConfig["should_enable_high_availability"].(bool))

	restoreSpec := management.NewRestoreSpec()

	restoreSpec.DomainManager = restorePcBody

	aJSON, _ := json.MarshalIndent(restoreSpec, "", "  ")
	log.Printf("[DEBUG] Restore PC Body: %s", string(aJSON))

	resp, err := conn.DomainManagerBackupsAPIInstance.Restore(utils.StringPtr(restoreSourceExtID), utils.StringPtr(restorableDomainManagerExtID), utils.StringPtr(restorePointExtID), restoreSpec)
	if err != nil {
		return diag.Errorf("error while restoring PC: %s", err)
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
		return diag.Errorf("error waiting for restore PC task to complete: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching task details : %v", err)
	}

	taskDetails := resourceUUID.Data.GetValue().(config.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Restore PC task details: %s", string(aJSON))

	entityAffected := taskDetails.EntitiesAffected

	// extract the PC UUID from the task response
	for _, entity := range entityAffected {
		if utils.StringValue(entity.Name) == "prism_central" && utils.StringValue(entity.Rel) == "prism:config:domain_manager" {
			uuid := entity.ExtId
			d.SetId(*uuid)
			return nil
		}
	}

	return diag.Errorf("error while fetching restored PC UUID From task response: %s", err)
}

func ResourceNutanixRestorePcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRestorePcUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixRestorePcDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
