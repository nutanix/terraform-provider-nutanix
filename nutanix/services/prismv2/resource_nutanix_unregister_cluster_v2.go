package prismv2

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/management"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	"log"
)

func ResourceNutanixUnregisterClusterV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixUnregisterClusterV2Create,
		ReadContext:   ResourceNutanixUnregisterClusterV2Read,
		UpdateContext: ResourceNutanixUnregisterClusterV2Update,
		DeleteContext: ResourceNutanixUnregisterClusterV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ResourceNutanixUnregisterClusterV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clsConn := meta.(*conns.Client).ClusterAPI
	conn := meta.(*conns.Client).PrismAPI
	extID := d.Get("ext_id")

	readClsResp, err := clsConn.ClusterEntityAPI.GetClusterById(utils.StringPtr(extID.(string)), nil)
	if err != nil {
		return diag.Errorf("error while fetching cluster entity : %v", err)
	}

	args := make(map[string]interface{})
	eTag := clsConn.ClusterEntityAPI.ApiClient.GetEtag(readClsResp)
	args["If-Match"] = utils.StringPtr(eTag)

	body := management.ClusterReference{
		ExtId: utils.StringPtr(extID.(string)),
	}

	resp, err := conn.DomainManagerAPIInstance.Unregister(utils.StringPtr(extID.(string)), &body, args)

	if err != nil {
		return diag.Errorf("error while unregistering cluster : %v", err)
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
		return diag.Errorf("error waiting for Cluster Unregister Task to complete: %s", err)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching task details : %v", err)
	}

	rUUID := resourceUUID.Data.GetValue().(config.Task)
	aJSON, _ := json.MarshalIndent(rUUID, "", "  ")
	log.Printf("[DEBUG] Unregister Cluster Task Details: %s", string(aJSON))

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	return ResourceNutanixUnregisterClusterV2Read(ctx, d, meta)
}

func ResourceNutanixUnregisterClusterV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func ResourceNutanixUnregisterClusterV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceNutanixUnregisterClusterV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
