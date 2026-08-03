package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	cmgmtConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	cmgmtRequest "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/request/clusters"
	prismCfg "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixSnmpTrapV2 manages SNMP trap configuration for a cluster.
// CRUD: CreateSnmpTrap, GetSnmpTrapById, UpdateSnmpTrapById, DeleteSnmpTrapById.
//
// Identity model mirrors the SNMP user resource: traps are children of a
// cluster's SNMP config, each with their own ext_id. Because CreateSnmpTrap
// returns only a TaskReference we resolve the new ext_id after task success
// by matching version + address + port + protocol via lookupSnmpTrapExtIDByAttrs.
func ResourceNutanixSnmpTrapV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSnmpTrapV2Create,
		ReadContext:   ResourceNutanixSnmpTrapV2Read,
		UpdateContext: ResourceNutanixSnmpTrapV2Update,
		DeleteContext: ResourceNutanixSnmpTrapV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Indicates the UUID of a cluster.",
			},
			"address": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Address of the SNMP trap receiver.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IPv4 address of the host.",
									},
									"prefix_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The prefix length of the network to which this host IPv4 address belongs.",
									},
								},
							},
						},
						"ipv6": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IPv6 address of the host.",
									},
									"prefix_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The prefix length of the network to which this host IPv6 address belongs.",
									},
								},
							},
						},
					},
				},
			},
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(SnmpTrapVersionStrings, false),
				Description:  "SNMP version. One of V2 or V3. Switching versions is not supported in place; changing this attribute forces resource replacement.",
			},
			"community_string": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Community string(plaintext) for SNMP version 2.0.",
			},
			"engine_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SNMP engine Id.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "SNMP port.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(SnmpProtocolStrings, false),
				Description:  "SNMP protocol. One of TCP, TCP6, UDP, UDP6.",
			},
			"receiver_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SNMP receiver name.",
			},
			"should_inform": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "SNMP information status.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SNMP username. Required when version is V3 and references an existing nutanix_snmp_user_v2 on the same cluster.",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity.",
			},
			"links": schemaForLinks(),
		},
	}
}

// buildSnmpTrapBody assembles a *cmgmtConfig.SnmpTrap from the resource state.
// Used by both Create and Update to keep the marshalling logic in one place.
func buildSnmpTrapBody(d *schema.ResourceData) *cmgmtConfig.SnmpTrap {
	body := cmgmtConfig.NewSnmpTrap()
	body.Version = common.ExpandEnum[cmgmtConfig.SnmpTrapVersion](d.Get("version").(string))
	body.Address = expandIPAddress(d.Get("address"))
	body.Port = utils.IntPtr(d.Get("port").(int))
	body.Protocol = common.ExpandEnum[cmgmtConfig.SnmpProtocol](d.Get("protocol").(string))

	if v, ok := d.GetOk("community_string"); ok && v.(string) != "" {
		body.CommunityString = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("engine_id"); ok && v.(string) != "" {
		body.EngineId = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("receiver_name"); ok && v.(string) != "" {
		body.RecieverName = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOkExists("should_inform"); ok {
		body.ShouldInform = utils.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOk("username"); ok && v.(string) != "" {
		body.Username = utils.StringPtr(v.(string))
	}
	return body
}

func ResourceNutanixSnmpTrapV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	body := buildSnmpTrapBody(d)

	bJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create SNMP Trap body: %s", string(bJSON))

	req := cmgmtRequest.CreateSnmpTrapRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		Body:         body,
	}
	resp, err := conn.ClusterEntityAPI.CreateSnmpTrap(ctx, &req)
	if err != nil {
		return diag.Errorf("error while creating SNMP trap on cluster (%s): %v", clusterExtID, err)
	}

	taskRef, ok := resp.Data.GetValue().(prismCfg.TaskReference)
	if !ok {
		return diag.Errorf("unexpected create SNMP trap response data type")
	}
	taskUUID := utils.StringValue(taskRef.ExtId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, taskUUID),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for SNMP trap create task (%s): %v", taskUUID, errWait)
	}

	// Resolve ext_id by re-reading the SNMP config and matching on the
	// natural identity tuple (version, ipv4|ipv6, port, protocol).
	match := snmpTrapMatch{
		Version:  d.Get("version").(string),
		Port:     d.Get("port").(int),
		Protocol: d.Get("protocol").(string),
	}
	if body.Address != nil {
		if body.Address.Ipv4 != nil {
			match.IPv4 = utils.StringValue(body.Address.Ipv4.Value)
		}
		if body.Address.Ipv6 != nil {
			match.IPv6 = utils.StringValue(body.Address.Ipv6.Value)
		}
	}
	extID, errLookup := lookupSnmpTrapExtIDByAttrs(ctx, conn, clusterExtID, match)
	if errLookup != nil {
		return diag.Errorf("error resolving created SNMP trap ext_id: %v", errLookup)
	}
	d.SetId(extID)

	return ResourceNutanixSnmpTrapV2Read(ctx, d, meta)
}

func ResourceNutanixSnmpTrapV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	req := cmgmtRequest.GetSnmpTrapByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	resp, err := conn.ClusterEntityAPI.GetSnmpTrapById(ctx, &req)
	if err != nil {
		return diag.Errorf("error while reading SNMP trap (%s) on cluster (%s): %v", extID, clusterExtID, err)
	}

	trap, ok := resp.Data.GetValue().(cmgmtConfig.SnmpTrap)
	if !ok {
		return diag.Errorf("unexpected response data type when reading SNMP trap")
	}

	if err := d.Set("ext_id", utils.StringValue(trap.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(trap.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(trap.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("address", flattenIPAddress(trap.Address)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("community_string", utils.StringValue(trap.CommunityString)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("engine_id", utils.StringValue(trap.EngineId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", utils.IntValue(trap.Port)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("protocol", common.FlattenPtrEnum(trap.Protocol)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("receiver_name", utils.StringValue(trap.RecieverName)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("should_inform", utils.BoolValue(trap.ShouldInform)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username", utils.StringValue(trap.Username)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", common.FlattenPtrEnum(trap.Version)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceNutanixSnmpTrapV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	// Read current state to capture ETag for the optimistic-concurrency PUT.
	getReq := cmgmtRequest.GetSnmpTrapByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	getResp, err := conn.ClusterEntityAPI.GetSnmpTrapById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while reading SNMP trap before update (%s): %v", extID, err)
	}
	args := make(map[string]interface{})
	etag := conn.ClusterEntityAPI.ApiClient.GetEtag(getResp)
	args["If-Match"] = utils.StringPtr(etag)

	body := buildSnmpTrapBody(d)
	body.ExtId = utils.StringPtr(extID)

	updReq := cmgmtRequest.UpdateSnmpTrapByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
		Body:         body,
	}
	resp, err := conn.ClusterEntityAPI.UpdateSnmpTrapById(ctx, &updReq, args)
	if err != nil {
		return diag.Errorf("error while updating SNMP trap (%s): %v", extID, err)
	}

	taskRef, ok := resp.Data.GetValue().(prismCfg.TaskReference)
	if !ok {
		return diag.Errorf("unexpected update SNMP trap response data type")
	}
	taskUUID := utils.StringValue(taskRef.ExtId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, taskUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for SNMP trap update task (%s): %v", taskUUID, errWait)
	}

	return ResourceNutanixSnmpTrapV2Read(ctx, d, meta)
}

func ResourceNutanixSnmpTrapV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	delReq := cmgmtRequest.DeleteSnmpTrapByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	resp, err := conn.ClusterEntityAPI.DeleteSnmpTrapById(ctx, &delReq)
	if err != nil {
		return diag.Errorf("error while deleting SNMP trap (%s): %v", extID, err)
	}

	taskRef, ok := resp.Data.GetValue().(prismCfg.TaskReference)
	if !ok {
		return diag.Errorf("unexpected delete SNMP trap response data type")
	}
	taskUUID := utils.StringValue(taskRef.ExtId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, taskUUID),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for SNMP trap delete task (%s): %v", taskUUID, errWait)
	}
	return nil
}
