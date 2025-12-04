package volumesv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	taskPoll "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/common/v1/config"
	volumesPrism "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/prism/v4/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixAssociateCategoryToVolumeGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixAssociateCategoryToVolumeGroupV2Create,
		ReadContext:   ResourceNutanixAssociateCategoryToVolumeGroupV2Read,
		UpdateContext: ResourceNutanixAssociateCategoryToVolumeGroupV2Update,
		DeleteContext: ResourceNutanixAssociateCategoryToVolumeGroupV2Delete,

		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"categories": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uris": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"entity_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "CATEGORY",
							ValidateFunc: validation.StringInSlice([]string{"CATEGORY"}, false),
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixAssociateCategoryToVolumeGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VolumeAPI

	extID := d.Get("ext_id")

	body := volumesClient.NewCategoryEntityReferences()

	if categories, ok := d.GetOk("categories"); ok {
		body.Categories = expandCategoryEntityReference(categories.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Payload to associate categories to Volume Group Body: %s", string(aJSON))

	resp, err := conn.VolumeAPIInstance.AssociateCategory(utils.StringPtr(extID.(string)), body)
	if err != nil {
		return diag.Errorf("error while associating categories to Volume Group : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the category to be associated to the Volume Group
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for associate categories task (%s) to finish: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Associate Category to Volume Group Task : %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(taskPoll.Task)

	aJSON, _ = json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Associate Category to Volume Group Task Details: %s", string(aJSON))

	d.SetId(utils.GenUUID())

	return nil
}

func ResourceNutanixAssociateCategoryToVolumeGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixAssociateCategoryToVolumeGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceNutanixAssociateCategoryToVolumeGroupV2Create(ctx, d, meta)
}

func ResourceNutanixAssociateCategoryToVolumeGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Dissociate Category from Volume Group
	conn := meta.(*conns.Client).VolumeAPI

	extID := d.Get("ext_id")

	body := volumesClient.NewCategoryEntityReferences()

	if categories, ok := d.GetOk("categories"); ok {
		body.Categories = expandCategoryEntityReference(categories.([]interface{}))
	}

	aJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Payload for disassociating category from Volume Group: %s", string(aJSON))

	resp, err := conn.VolumeAPIInstance.DisassociateCategory(utils.StringPtr(extID.(string)), body)
	if err != nil {
		return diag.Errorf("error while Dissociating Category from Volume Group : %v", err)
	}

	TaskRef := resp.Data.GetValue().(volumesPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: taskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for Category task (%s) to Dissociate from Volume Group: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching Dissociate Category from Volume Group Task : %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(taskPoll.Task)

	aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Dissociate Category from Volume Group Task Details: %s", string(aJSON))

	return nil
}

func expandCategoryEntityReference(categoryEntityReference interface{}) []config.EntityReference {
	if categoryEntityReference == nil {
		return nil
	}
	entityReferenceList := categoryEntityReference.([]interface{})
	entityReferenceListExpanded := make([]config.EntityReference, 0)
	for _, entityReference := range entityReferenceList {
		entityReferenceMap := entityReference.(map[string]interface{})
		entityReferenceExpanded := config.EntityReference{
			ExtId:      utils.StringPtr(entityReferenceMap["ext_id"].(string)),
			Name:       utils.StringPtr(entityReferenceMap["name"].(string)),
			EntityType: expandCategoryEntityType(entityReferenceMap["entity_type"].(string)),
			Uris:       expandListOfString(entityReferenceMap["uris"].([]interface{})),
		}
		entityReferenceListExpanded = append(entityReferenceListExpanded, entityReferenceExpanded)
	}

	return entityReferenceListExpanded
}

func expandCategoryEntityType(entityType string) *config.EntityType {
	if entityType == "" {
		return nil
	}

	if entityType == "CATEGORY" {
		p := config.ENTITYTYPE_CATEGORY
		return &p
	}
	return nil
}

func expandListOfString(list []interface{}) []string {
	stringListStr := make([]string, len(list))
	for i, v := range list {
		stringListStr[i] = v.(string)
	}
	return stringListStr
}
