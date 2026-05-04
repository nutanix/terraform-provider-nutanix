package cluster_managementv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	prismConfig "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/prism/v4/config"
	prismConfigV2 "github.com/nutanix/ntnx-api-golang-clients/prism-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func ResourceNutanixAddSnmpTransportV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixAddSnmpTransportV2Create,
		ReadContext:   resourceNutanixAddSnmpTransportV2Read,
		DeleteContext: resourceNutanixAddSnmpTransportV2Delete,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Indicates the UUID of a cluster.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "SNMP port.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "SNMP transport protocol.",
				ValidateFunc: validation.StringInSlice([]string{"UDP", "UDP6", "TCP", "TCP6"}, false),
			},
		},
	}
}

func resourceNutanixAddSnmpTransportV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Cluster_managementAPI

	body := &config.SnmpTransport{}

	clusterExtID := d.Get("cluster_ext_id").(string)

	if port, ok := d.GetOk("port"); ok {
		body.Port = utils.IntPtr(port.(int))
	}
	if protocol, ok := d.GetOk("protocol"); ok {
		p := common.ExpandEnum[config.SnmpProtocol](protocol)
		body.Protocol = p
	}

	resp, err := conn.ClustersAPI.AddSnmpTransport(utils.StringPtr(clusterExtID), body)
	if err != nil {
		return diag.Errorf("error while adding SNMP transport: %v", err)
	}

	TaskRef := resp.Data.GetValue().(prismConfig.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for SNMP transport task (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	resourceUUID, err := taskconn.TaskRefAPI.GetTaskById(taskUUID, nil)
	if err != nil {
		return diag.Errorf("error while fetching SNMP transport task details: %v", err)
	}
	rUUID := resourceUUID.Data.GetValue().(prismConfigV2.Task)

	uuid := rUUID.EntitiesAffected[0].ExtId
	d.SetId(*uuid)

	return resourceNutanixAddSnmpTransportV2Read(ctx, d, meta)
}

func resourceNutanixAddSnmpTransportV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNutanixAddSnmpTransportV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
