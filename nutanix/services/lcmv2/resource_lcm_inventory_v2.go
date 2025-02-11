package lcmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	lcmInventoryResp "github.com/nutanix/ntnx-api-golang-clients/lcm-go-client/v4/models/lcm/v4/operations"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func ResourceLcmInventoryV2() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ntnx_request_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"x_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceLcmPerformInventoryV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LcmAPI
	clusterId := d.Get("x_cluster_id").(string)
	ntnxRequestId, ok := d.Get("ntnx_request_id").(string)
	if !ok || ntnxRequestId == "" {
		return diag.Errorf("ntnx_request_id is required and cannot be null or empty")
	}

	args := make(map[string]interface{})
	args["X-Cluster-Id"] = clusterId
	args["NTNX-Request-Id"] = ntnxRequestId

	resp, err := conn.LcmInventoryAPIInstance.Inventory(args)
	if err != nil {
		return diag.Errorf("error while performing the inventory: %v", err)
	}
	getResp := resp.Data.GetValue().(lcmInventoryResp.InventoryApiResponse)
	TaskRef := getResp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	d.Set("ext_id", taskUUID)
	return nil
}
