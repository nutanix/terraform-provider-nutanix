package microsegv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/common/v1/config"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/prism/v4/config"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	commonUtils "github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixEntityGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixEntityGroupV2Create,
		ReadContext:   ResourceNutanixEntityGroupV2Read,
		UpdateContext: ResourceNutanixEntityGroupV2Update,
		DeleteContext: ResourceNutanixEntityGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_ext_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     schemaAllowedConfigResource(),
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": schemaForLinks(),
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceNutanixEntityGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	bodySpec := import1.NewEntityGroup()

	if name, ok := d.GetOk("name"); ok {
		bodySpec.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		bodySpec.Description = utils.StringPtr(description.(string))
	}
	if ownerExtId, ok := d.GetOk("owner_ext_id"); ok {
		bodySpec.OwnerExtId = utils.StringPtr(ownerExtId.(string))
	}
	if policyExtIds, ok := d.GetOk("policy_ext_ids"); ok {
		bodySpec.PolicyExtIds = commonUtils.ExpandListOfString(policyExtIds.([]interface{}))
	}
	if allowedConfig, ok := d.GetOk("allowed_config"); ok {
		bodySpec.AllowedConfig = expandAllowedConfig(allowedConfig.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(bodySpec, "", "  ")
	log.Printf("[DEBUG] Create Entity Group Body Spec: %s", string(aJSON))

	resp, err := conn.EntityGroupsAPIInstance.CreateEntityGroup(bodySpec)
	if err != nil {
		return diag.Errorf("error while creating Entity Group: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the task to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Entity Group (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Entity Group Task : %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Create Entity Group Task Response Details: %s", string(aJSON))

	uuid := taskDetails.CompletionDetails[0].Value.GetValue().(string)

	d.SetId(uuid)

	return ResourceNutanixEntityGroupV2Read(ctx, d, meta)
}

func ResourceNutanixEntityGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	extID := d.Id()

	resp, err := conn.EntityGroupsAPIInstance.GetEntityGroupById(utils.StringPtr(extID))
	if err != nil {
		return diag.Errorf("error while fetching Entity Group: %s", err)
	}

	getResp := resp.Data.GetValue().(import1.EntityGroup)

	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] Read Entity Group Response Details: %s", string(aJSON))

	if err := d.Set("ext_id", getResp.ExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", getResp.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", getResp.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("owner_ext_id", getResp.OwnerExtId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("policy_ext_ids", getResp.PolicyExtIds); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allowed_config", flattenAllowedConfig(getResp.AllowedConfig)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixEntityGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	readResp, err := conn.EntityGroupsAPIInstance.GetEntityGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while fetching Entity Group: %v", err)
	}
	// extract e-tag
	args := make(map[string]interface{})
	etag := conn.EntityGroupsAPIInstance.ApiClient.GetEtag(readResp)
	args["If-Match"] = utils.StringPtr(etag)

	updateSpec := import1.NewEntityGroup()

	if name, ok := d.GetOk("name"); ok {
		updateSpec.Name = utils.StringPtr(name.(string))
	}
	if description, ok := d.GetOk("description"); ok {
		updateSpec.Description = utils.StringPtr(description.(string))
	}
	if ownerExtId, ok := d.GetOk("owner_ext_id"); ok {
		updateSpec.OwnerExtId = utils.StringPtr(ownerExtId.(string))
	}
	if policyExtIds, ok := d.GetOk("policy_ext_ids"); ok {
		updateSpec.PolicyExtIds = commonUtils.ExpandListOfString(policyExtIds.([]interface{}))
	}
	if allowedConfig, ok := d.GetOk("allowed_config"); ok {
		updateSpec.AllowedConfig = expandAllowedConfig(allowedConfig.([]interface{}))
	}

	resp, err := conn.EntityGroupsAPIInstance.UpdateEntityGroupById(utils.StringPtr(d.Id()), updateSpec, args)
	if err != nil {
		return diag.Errorf("error while updating Entity Group: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the task to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Entity Group (%s) to update: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Entity Group Task : %v", err)
	}

	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Update Entity Group Task Response Details: %s", string(aJSON))

	return ResourceNutanixEntityGroupV2Read(ctx, d, meta)
}

func ResourceNutanixEntityGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).MicroSegAPI

	resp, err := conn.EntityGroupsAPIInstance.DeleteEntityGroupById(utils.StringPtr(d.Id()))
	if err != nil {
		return diag.Errorf("error while deleting Entity Group: %v", err)
	}
	TaskRef := resp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the task to complete
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: commonUtils.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Entity Group (%s) to delete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while deleting Entity Group Task : %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(prismConfig.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Delete Entity Group Task Response Details: %s", string(aJSON))

	return nil
}

// schemas funcs for resource
func schemaAllowedConfigResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"entities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     schemaAllowedEntityResource(),
			},
		},
	}
}

func schemaAllowedEntityResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kube_entities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"reference_ext_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"select_by": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     schemaAllowedSelectByResource(),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func schemaAllowedSelectByResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Add fields based on AllowedSelectBy struct if needed
			// This is a placeholder as the exact structure may vary
		},
	}
}

// expander funcs
func expandAllowedConfig(allowedConfig []interface{}) *import1.AllowedConfig {
	if len(allowedConfig) == 0 {
		return nil
	}

	allowedConfigVal := allowedConfig[0].(map[string]interface{})
	allowedConfigSpec := import1.NewAllowedConfig()

	if entities, ok := allowedConfigVal["entities"]; ok {
		allowedConfigSpec.Entities = expandAllowedEntities(entities.([]interface{}))
	}

	return allowedConfigSpec
}

func expandAllowedEntities(entities []interface{}) []import1.AllowedEntity {
	if len(entities) == 0 {
		return nil
	}

	entitiesSpec := make([]import1.AllowedEntity, 0)

	for _, entity := range entities {
		entityVal := entity.(map[string]interface{})

		entitySpec := import1.AllowedEntity{}
		if kubeEntities, ok := entityVal["kube_entities"]; ok {
			entitySpec.KubeEntities = commonUtils.ExpandListOfString(kubeEntities.([]interface{}))
		}
		if referenceExtIds, ok := entityVal["reference_ext_ids"]; ok {
			entitySpec.ReferenceExtIds = commonUtils.ExpandListOfString(referenceExtIds.([]interface{}))
		}
		if selectBy, ok := entityVal["select_by"]; ok {
			entitySpec.SelectBy = expandAllowedSelectBy(selectBy.([]interface{}))
		}
		if entityType, ok := entityVal["type"]; ok {
			entitySpec.Type = expandAllowedType(entityType.(string))
		}

		entitiesSpec = append(entitiesSpec, entitySpec)
	}

	return entitiesSpec
}

func expandAllowedSelectBy(selectBy []interface{}) *import1.AllowedSelectBy {
	if len(selectBy) == 0 {
		return nil
	}

	selectByVal := selectBy[0].(map[string]interface{})
	selectBySpec := import1.NewAllowedSelectBy()

	// Add fields based on AllowedSelectBy struct if needed
	_ = selectByVal

	return selectBySpec
}

func expandAllowedType(entityType string) *import1.AllowedType {
	if entityType == "" {
		return nil
	}

	// Map string to AllowedType enum
	// This is a placeholder - adjust based on actual enum values
	var allowedType import1.AllowedType
	// Add enum mapping logic here based on actual AllowedType enum values

	return &allowedType
}
