package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	clustermgmtPrism "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixClusterCategoriesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterCategoriesV2Create,
		ReadContext:   ResourceNutanixClusterCategoriesV2Read,
		UpdateContext: ResourceNutanixClusterCategoriesV2Update,
		DeleteContext: ResourceNutanixClusterCategoriesV2Delete,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"categories": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      common.HashStringItem,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ResourceNutanixClusterCategoriesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id")

	body := config.CategoryEntityReferences{}

	if categories, ok := d.GetOk("categories"); ok {
		body.Categories = common.ExpandListOfString(common.InterfaceToSlice(categories))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Associate Categories to Cluster Request Body: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.AssociateCategoriesToCluster(utils.StringPtr(clusterExtID.(string)), &body)
	if err != nil {
		return diag.Errorf("error while associating categories to cluster : %v", err)
	}

	TaskRef := resp.Data.GetValue().(clustermgmtPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the categories to be associated to the cluster
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for categories to be associated to the cluster (%s) : %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get UUID from TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching associate categories to cluster task : %v", err)
	}

	aJSON, _ = json.Marshal(taskResp)
	log.Printf("[DEBUG] associate categories to cluster task details: %s", string(aJSON))

	d.SetId(resource.UniqueId())
	return nil
}

func ResourceNutanixClusterCategoriesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterCategoriesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterCategoriesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	body := &config.CategoryEntityReferences{}
	clusterExtID := d.Get("cluster_ext_id")

	if categories, ok := d.GetOk("categories"); ok {
		body.Categories = common.ExpandListOfString(common.InterfaceToSlice(categories))
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Disassociate Categories from Cluster Request Body: %s", string(aJSON))

	resp, err := conn.ClusterEntityAPI.DisassociateCategoriesFromCluster(utils.StringPtr(clusterExtID.(string)), body)
	if err != nil {
		return diag.Errorf("error while Disassociating Categories from Cluster : %v", err)
	}

	TaskRef := resp.Data.GetValue().(clustermgmtPrism.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the node to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		resourceUUID, _ := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		rUUID := resourceUUID.Data.GetValue().(import2.Task)
		aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
		log.Printf("Error Disassociate Categories from Cluster Task Details : %s", string(aJSON))
		return diag.Errorf("error waiting for categories to be disassociated from cluster (%s) : %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details TASK API
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching disassociate categories from cluster task : %v", err)
	}

	aJSON, _ = json.MarshalIndent(taskResp, "", "  ")
	log.Printf("disassociate categories from cluster task details : %s", string(aJSON))
	return nil
}
