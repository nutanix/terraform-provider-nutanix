package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import4 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"

	import3 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixClusterProfileAssociationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixClusterProfileAssociationV2Create,
		ReadContext:   ResourceNutanixClusterProfileAssociationV2Read,
		UpdateContext: ResourceNutanixClusterProfileAssociationV2Update,
		DeleteContext: ResourceNutanixClusterProfileAssociationV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dryrun": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"clusters": {
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

func ResourceNutanixClusterProfileAssociationV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	extID := d.Get("ext_id").(string)
	dryrun := d.Get("dryrun").(bool)
	clustersUUIDs := common.ExpandListOfString(common.InterfaceToSlice(d.Get("clusters")))

	clustersRef := make([]import4.ClusterReference, 0)

	for _, clusterUUID := range clustersUUIDs {
		clusterRef := import4.ClusterReference{
			Uuid: utils.StringPtr(clusterUUID),
		}
		clustersRef = append(clustersRef, clusterRef)
	}

	ClusterReferenceListSpec := &import4.ClusterReferenceListSpec{
		Clusters: clustersRef,
	}

	associateResp, associateErr := conn.ClusterProfilesAPI.ApplyClusterProfile(utils.StringPtr(extID), ClusterReferenceListSpec, utils.BoolPtr(dryrun))
	if associateErr != nil {
		return diag.FromErr(associateErr)
	}

	TaskRef := associateResp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster profile (%s) to associate: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get Task Details
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster profile association task UUID : %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Associate Cluster Profile Task Details: %s", string(aJSON))

	d.SetId(resource.UniqueId())

	return nil
}

func ResourceNutanixClusterProfileAssociationV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixClusterProfileAssociationV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	extID := d.Get("ext_id").(string)
	dryrun := d.Get("dryrun").(bool)

	// Get old and new cluster sets
	oldClustersRaw, newClustersRaw := d.GetChange("clusters")
	oldClusters := common.ExpandListOfString(common.InterfaceToSlice(oldClustersRaw))
	newClusters := common.ExpandListOfString(common.InterfaceToSlice(newClustersRaw))

	// Find clusters to associate and deassociate
	clustersToAssociate, clustersToDeassociate := common.DiffStringSets(oldClusters, newClusters)

	// Handle dryrun changes: if dryrun changed from true to false and clusters haven't changed,
	// we need to actually associate the clusters (since previous dryrun=true didn't actually associate them)
	hasDryrunChange := d.HasChange("dryrun")
	if hasDryrunChange {
		oldDryrunRaw, newDryrunRaw := d.GetChange("dryrun")
		oldDryrun := oldDryrunRaw.(bool)
		newDryrun := newDryrunRaw.(bool)

		// If dryrun changed from true to false, and clusters are the same, we need to associate
		if oldDryrun && !newDryrun && len(clustersToAssociate) == 0 && len(clustersToDeassociate) == 0 {
			log.Printf("[DEBUG] dryrun changed from true to false, associating existing clusters")
			clustersToAssociate = newClusters
		}
		// If dryrun changed from false to true, we'll just do a dry run
		if !oldDryrun && newDryrun {
			log.Printf("[DEBUG] dryrun changed from false to true, will perform dry run for cluster operations")
		}
	}

	// Associate new clusters
	if len(clustersToAssociate) > 0 {
		clustersRef := make([]import4.ClusterReference, 0)
		for _, clusterUUID := range clustersToAssociate {
			clusterRef := import4.ClusterReference{
				Uuid: utils.StringPtr(clusterUUID),
			}
			clustersRef = append(clustersRef, clusterRef)
		}

		ClusterReferenceListSpec := &import4.ClusterReferenceListSpec{
			Clusters: clustersRef,
		}

		associateResp, associateErr := conn.ClusterProfilesAPI.ApplyClusterProfile(utils.StringPtr(extID), ClusterReferenceListSpec, utils.BoolPtr(dryrun))
		if associateErr != nil {
			return diag.Errorf("error associating clusters to cluster profile: %v", associateErr)
		}

		TaskRef := associateResp.Data.GetValue().(import3.TaskReference)
		taskUUID := TaskRef.ExtId

		// Wait for the association task to complete
		stateConf := &resource.StateChangeConf{
			Pending: []string{"QUEUED", "RUNNING", "PENDING"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for cluster profile (%s) to associate: %s", utils.StringValue(taskUUID), errWaitTask)
		}

		// Get Task Details
		taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return diag.Errorf("error while fetching cluster profile association task UUID : %v", err)
		}
		taskDetails := taskResp.Data.GetValue().(import2.Task)
		aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
		log.Printf("[DEBUG] Associate Cluster Profile Task Details: %s", string(aJSON))
	}

	// Deassociate removed clusters
	if len(clustersToDeassociate) > 0 {
		clustersRef := make([]import4.ClusterReference, 0)
		for _, clusterUUID := range clustersToDeassociate {
			clusterRef := import4.ClusterReference{
				Uuid: utils.StringPtr(clusterUUID),
			}
			clustersRef = append(clustersRef, clusterRef)
		}

		ClusterReferenceListSpec := &import4.ClusterReferenceListSpec{
			Clusters: clustersRef,
		}

		disassociateResp, disassociateErr := conn.ClusterProfilesAPI.DisassociateClusterFromClusterProfile(utils.StringPtr(extID), ClusterReferenceListSpec)
		if disassociateErr != nil {
			return diag.Errorf("error deassociating clusters from cluster profile: %v", disassociateErr)
		}

		TaskRef := disassociateResp.Data.GetValue().(import3.TaskReference)
		taskUUID := TaskRef.ExtId

		// Wait for the deassociation task to complete
		stateConf := &resource.StateChangeConf{
			Pending: []string{"QUEUED", "RUNNING", "PENDING"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
			return diag.Errorf("error waiting for cluster profile (%s) to disassociate: %s", utils.StringValue(taskUUID), errWaitTask)
		}

		// Get Task Details
		taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
		if err != nil {
			return diag.Errorf("error while fetching cluster profile disassociation task UUID : %v", err)
		}
		taskDetails := taskResp.Data.GetValue().(import2.Task)
		aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
		log.Printf("[DEBUG] Disassociate Cluster Profile Task Details: %s", string(aJSON))
	}

	return nil
}

func ResourceNutanixClusterProfileAssociationV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	extID := d.Get("ext_id").(string)
	clustersUUIDs := common.ExpandListOfString(common.InterfaceToSlice(d.Get("clusters")))

	clustersRef := make([]import4.ClusterReference, 0)

	for _, clusterUUID := range clustersUUIDs {
		clusterRef := import4.ClusterReference{
			Uuid: utils.StringPtr(clusterUUID),
		}
		clustersRef = append(clustersRef, clusterRef)
	}

	ClusterReferenceListSpec := &import4.ClusterReferenceListSpec{
		Clusters: clustersRef,
	}

	disassociateResp, disassociateErr := conn.ClusterProfilesAPI.DisassociateClusterFromClusterProfile(utils.StringPtr(extID), ClusterReferenceListSpec)
	if disassociateErr != nil {
		return diag.FromErr(disassociateErr)
	}

	TaskRef := disassociateResp.Data.GetValue().(import3.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	// Wait for the cluster to be available
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for cluster profile (%s) to disassociate: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	// Get Task Details
	taskResp, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster profile disassociation task UUID : %v", err)
	}
	taskDetails := taskResp.Data.GetValue().(import2.Task)
	aJSON, _ := json.MarshalIndent(taskDetails, "", "  ")
	log.Printf("[DEBUG] Disassociate Cluster Profile Task Details: %s", string(aJSON))

	d.SetId("")

	return nil
}
