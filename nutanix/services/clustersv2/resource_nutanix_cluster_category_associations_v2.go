package clustersv2

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	import5 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/request/clusters"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/request/tasks"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixClusterCategoriesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixClusterCategoriesV2Create,
		ReadContext:   resourceNutanixClusterCategoriesV2Read,
		UpdateContext: resourceNutanixClusterCategoriesV2Update,
		DeleteContext: resourceNutanixClusterCategoriesV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNutanixClusterCategoriesV2Import,
		},
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"categories": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceNutanixClusterCategoriesV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	categories := common.InterfaceToSlice(d.Get("categories"))

	body := config.CategoryEntityReferences{
		Categories: common.ExpandListOfString(categories),
	}
	associateCategoriesToClusterRequest := import5.AssociateCategoriesToClusterRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		Body:         &body,
	}
	resp, err := conn.ClusterEntityAPI.AssociateCategoriesToCluster(ctx, &associateCategoriesToClusterRequest)
	if err != nil {
		return diag.Errorf("error while associating categories to cluster: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// Wait for the categories to be associated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for categories to be associated to the cluster (%s): %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details
	getTaskByIdRequest := import3.GetTaskByIdRequest{
		ExtId: taskUUID,
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching associate categories to cluster task: %v", err)
	}

	aJSON, _ := json.Marshal(taskResp)
	log.Printf("[DEBUG] Associate categories to cluster task details: %s", string(aJSON))

	d.SetId(clusterExtID)
	return nil
}

func resourceNutanixClusterCategoriesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	conn := meta.(*conns.Client).ClusterAPI
	clusterExtID := d.Get("cluster_ext_id").(string)

	getClusterByIdRequest := import5.GetClusterByIdRequest{
		ExtId:   utils.StringPtr(clusterExtID),
		Expand_: nil,
	}

	resp, err := conn.ClusterEntityAPI.GetClusterById(ctx, &getClusterByIdRequest)
	if err != nil {
		return diag.Errorf("error fetching cluster categories: %v", err)
	}

	categories := resp.Data.GetValue().(config.Cluster).Categories

	d.Set("categories", categories)
	return nil

}

func resourceNutanixClusterCategoriesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterExtID := d.Get("cluster_ext_id").(string)

	log.Printf("[DEBUG] Handling category association changes for cluster: %s", clusterExtID)

	// Get old and new category values
	oldCategoriesRaw, newCategoriesRaw := d.GetChange("categories")

	// Use shared function to handle category updates
	return UpdateClusterCategories(ctx, d, meta, clusterExtID, oldCategoriesRaw, newCategoriesRaw)
}

func resourceNutanixClusterCategoriesV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI
	clusterExtID := d.Get("cluster_ext_id").(string)

	categoriesToDisassociate := common.ExpandListOfString(common.InterfaceToSlice(d.Get("categories")))

	if len(categoriesToDisassociate) == 0 {
		d.SetId("")
		log.Printf("[DEBUG] No categories to disassociate from cluster, resource destroyed: %s", clusterExtID)
		return nil
	}
	body := &config.CategoryEntityReferences{
		Categories: categoriesToDisassociate,
	}
	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Disassociate Categories from Cluster Request Body: %s", string(aJSON))

	disassociateCategoriesFromClusterRequest := import5.DisassociateCategoriesFromClusterRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		Body:         body,
	}
	resp, err := conn.ClusterEntityAPI.DisassociateCategoriesFromCluster(ctx, &disassociateCategoriesFromClusterRequest)
	if err != nil {
		return diag.Errorf("error while disassociating categories from cluster: %v", err)
	}

	TaskRef := resp.Data.GetValue().(import1.TaskReference)
	taskUUID := TaskRef.ExtId

	// Wait for the categories to be disassociated
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		getTaskByIdRequest := import3.GetTaskByIdRequest{
			ExtId: taskUUID,
		}
		resourceUUID, _ := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
		if resourceUUID != nil {
			rUUID := resourceUUID.Data.GetValue().(import2.Task)
			aJSON, _ = json.MarshalIndent(rUUID, "", "  ")
			log.Printf("[DEBUG] Error Disassociate Categories from Cluster Task Details: %s", string(aJSON))
		}
		return diag.Errorf("error waiting for categories to be disassociated from cluster (%s): %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get task details
	getTaskByIdRequest := import3.GetTaskByIdRequest{
		ExtId: taskUUID,
	}
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(ctx, &getTaskByIdRequest)
	if err != nil {
		return diag.Errorf("error while fetching disassociate categories from cluster task: %v", err)
	}

	aJSON, _ = json.MarshalIndent(taskResp, "", "  ")
	log.Printf("[DEBUG] Disassociate categories from cluster task details: %s", string(aJSON))
	d.SetId("")
	return nil
}

func resourceNutanixClusterCategoriesV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*conns.Client).ClusterAPI
	// During StateContext import, Terraform populates the resource ID (d.Id()) but not
	// necessarily the required attributes. Use d.Id() as the source of truth.
	clusterExtID := d.Id()
	if clusterExtID == "" {
		clusterExtID = d.Get("cluster_ext_id").(string)
	}
	if err := d.Set("cluster_ext_id", clusterExtID); err != nil {
		return nil, err
	}

	getClusterByIdRequest := import5.GetClusterByIdRequest{
		ExtId:   utils.StringPtr(clusterExtID),
		Expand_: nil,
	}

	resp, err := conn.ClusterEntityAPI.GetClusterById(ctx, &getClusterByIdRequest)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster for cluster categories import: %v", err)
	}

	categories := resp.Data.GetValue().(config.Cluster).Categories

	d.Set("categories", categories)

	d.SetId(clusterExtID)
	return []*schema.ResourceData{d}, nil
}
