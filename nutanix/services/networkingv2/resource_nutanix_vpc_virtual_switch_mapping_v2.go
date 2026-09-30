package networkingv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkingConfig "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
	networkingMappingReq "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/request/vpcvirtualswitchmappings"
	import4 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixVpcVirtualSwitchMappingV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixVpcVirtualSwitchMappingV2Create,
		ReadContext:   resourceNutanixVpcVirtualSwitchMappingV2Read,
		DeleteContext: resourceNutanixVpcVirtualSwitchMappingV2Delete,
		Schema: map[string]*schema.Schema{
			"mappings": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "VPC to virtual switch mapping entries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_switch_uuid": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "UUID of the virtual switch.",
						},
						"cluster_uuids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "UUID of the cluster.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"is_all_traffic_permitted": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: "Whether to permit all traffic through virtual switch or only the ICMP and statistics collection requests.",
						},
						"metadata": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: DatasourceMetadataSchemaV2(),
							},
						},
					},
				},
			},
		},
	}
}

func resourceNutanixVpcVirtualSwitchMappingV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).NetworkingAPI

	mappings := d.Get("mappings").([]interface{})
	body := make([]networkingConfig.VpcVirtualSwitchMapping, len(mappings))

	for i, m := range mappings {
		mMap := m.(map[string]interface{})
		mapping := networkingConfig.VpcVirtualSwitchMapping{
			VirtualSwitchUuid: utils.StringPtr(mMap["virtual_switch_uuid"].(string)),
		}
		if clusterUuids, ok := mMap["cluster_uuids"]; ok {
			mapping.ClusterUuids = expandStringList(clusterUuids.([]interface{}))
		}
		if isAllTraffic, ok := mMap["is_all_traffic_permitted"]; ok {
			mapping.IsAllTrafficPermitted = utils.BoolPtr(isAllTraffic.(bool))
		}
		if metadataList, ok := mMap["metadata"].([]interface{}); ok && len(metadataList) > 0 {
			mapping.Metadata = expandMetadata(metadataList)
		}
		body[i] = mapping
	}

	createReq := networkingMappingReq.CreateVpcVirtualSwitchMappingRequest{
		Body: &body,
	}

	aJSON, _ := json.MarshalIndent(createReq, "", " ")
	log.Printf("[DEBUG] VpcVirtualSwitchMapping create payload: %s", string(aJSON))

	resp, err := conn.VpcVirtualSwitchMappingsAPI.CreateVpcVirtualSwitchMapping(ctx, &createReq)
	if err != nil {
		return diag.Errorf("error while creating VPC Virtual Switch Mapping: %v", err)
	}

	taskRef := resp.Data.GetValue().(import4.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for VPC Virtual Switch Mapping (%s) to create: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(resource.UniqueId())
	return resourceNutanixVpcVirtualSwitchMappingV2Read(ctx, d, meta)
}

func resourceNutanixVpcVirtualSwitchMappingV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixVpcVirtualSwitchMappingV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
