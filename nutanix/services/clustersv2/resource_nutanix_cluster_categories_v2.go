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
		Importer: &schema.ResourceImporter{
			StateContext: ResourceNutanixClusterCategoriesV2Import,
		},
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
	conn := meta.(*conns.Client).ClusterAPI

	// Get cluster_ext_id from resource configuration
	clusterExtID, ok := d.Get("cluster_ext_id").(string)
	if !ok || clusterExtID == "" {
		return diag.Errorf("cluster_ext_id is required and cannot be empty")
	}

	log.Printf("[DEBUG] Reading cluster categories for cluster_ext_id: %s", clusterExtID)

	// Call GetClusterById API using the cluster_ext_id
	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(clusterExtID), nil)
	if err != nil {
		// If cluster is not found, mark resource as removed by not setting state
		// This will cause Terraform to detect the resource needs to be recreated
		log.Printf("[DEBUG] Error fetching cluster by ID %s: %v", clusterExtID, err)
		d.SetId("")
		return diag.Errorf("error while fetching cluster by ID %s: %v. The cluster may have been deleted or does not exist", clusterExtID, err)
	}

	// Extract the cluster response
	getResp := resp.Data.GetValue().(config.Cluster)
	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] GetClusterById Response Details: %s", string(aJSON))

	// Extract the categories field from the cluster response
	categories := getResp.Categories

	// Convert API response format ([]string) to Terraform schema format (schema.Set)
	// Convert []string to []interface{}
	categoriesList := make([]interface{}, 0)
	if len(categories) > 0 {
		for _, category := range categories {
			if category != "" {
				categoriesList = append(categoriesList, category)
			}
		}
	}

	// Create a schema.Set from the list
	categoriesSet := schema.NewSet(common.HashStringItem, categoriesList)

	// Set cluster_ext_id in state to ensure consistency
	if err := d.Set("cluster_ext_id", clusterExtID); err != nil {
		return diag.FromErr(err)
	}

	// Set categories in state with current category associations
	if err := d.Set("categories", categoriesSet); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully read cluster categories. cluster_ext_id: %s, categories count: %d", clusterExtID, len(categoriesList))

	return nil
}

func ResourceNutanixClusterCategoriesV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterExtID := d.Get("cluster_ext_id").(string)

	// Check if categories have changed
	if !d.HasChange("categories") {
		log.Printf("[DEBUG] No changes detected in categories, skipping update")
		return ResourceNutanixClusterCategoriesV2Read(ctx, d, meta)
	}

	// Get old and new category values
	oldCategoriesRaw, newCategoriesRaw := d.GetChange("categories")

	// Use shared function to handle category updates
	if diags := UpdateClusterCategories(ctx, d, meta, clusterExtID, oldCategoriesRaw, newCategoriesRaw); diags.HasError() {
		return diags
	}

	// Refresh state by calling Read function
	return ResourceNutanixClusterCategoriesV2Read(ctx, d, meta)
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

func ResourceNutanixClusterCategoriesV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*conns.Client).ClusterAPI

	// The import ID is the cluster_ext_id
	clusterExtID := d.Id()
	if clusterExtID == "" {
		return nil, fmt.Errorf("cluster_ext_id cannot be empty")
	}

	log.Printf("[DEBUG] Importing cluster categories for cluster_ext_id: %s", clusterExtID)

	// Call GetClusterById API using the cluster_ext_id
	resp, err := conn.ClusterEntityAPI.GetClusterById(utils.StringPtr(clusterExtID), nil)
	if err != nil {
		return nil, fmt.Errorf("error while fetching cluster by ID %s: %v", clusterExtID, err)
	}

	// Extract the cluster response
	getResp := resp.Data.GetValue().(config.Cluster)
	aJSON, _ := json.MarshalIndent(getResp, "", "  ")
	log.Printf("[DEBUG] GetClusterById Response Details: %s", string(aJSON))

	// Extract the categories field from the cluster response
	categories := getResp.Categories

	// Convert API response format ([]string) to Terraform schema format (schema.Set)
	// Convert []string to []interface{}
	categoriesList := make([]interface{}, 0)
	if len(categories) > 0 {
		for _, category := range categories {
			if category != "" {
				categoriesList = append(categoriesList, category)
			}
		}
	}

	// Create a schema.Set from the list
	categoriesSet := schema.NewSet(common.HashStringItem, categoriesList)

	// Set cluster_ext_id in state
	if err := d.Set("cluster_ext_id", clusterExtID); err != nil {
		return nil, fmt.Errorf("error setting cluster_ext_id: %v", err)
	}

	// Set categories in state
	if err := d.Set("categories", categoriesSet); err != nil {
		return nil, fmt.Errorf("error setting categories: %v", err)
	}

	// Set the resource ID (using a unique ID similar to Create)
	d.SetId(resource.UniqueId())

	log.Printf("[DEBUG] Successfully imported cluster categories. cluster_ext_id: %s, categories count: %d", clusterExtID, len(categoriesList))

	return []*schema.ResourceData{d}, nil
}
