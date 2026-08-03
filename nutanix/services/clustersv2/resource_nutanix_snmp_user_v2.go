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

// ResourceNutanixSnmpUserV2 manages SNMP user configuration for a cluster.
// CRUD: CreateSnmpUser, GetSnmpUserById, UpdateSnmpUserById, DeleteSnmpUserById.
func ResourceNutanixSnmpUserV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixSnmpUserV2Create,
		ReadContext:   ResourceNutanixSnmpUserV2Read,
		UpdateContext: ResourceNutanixSnmpUserV2Update,
		DeleteContext: ResourceNutanixSnmpUserV2Delete,
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
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SNMP username. For SNMP trap v3 version, SNMP username is required parameter.",
			},
			"auth_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(SnmpAuthTypeStrings, false),
				Description:  "SNMP user authentication type. One of MD5, SHA, SHA224, SHA256, SHA384, SHA512.",
			},
			"auth_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "SNMP user authentication key.",
			},
			"priv_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(SnmpPrivTypeStrings, false),
				Description:  "SNMP user encryption type. One of DES, AES, AES192, AES256.",
			},
			"priv_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "SNMP user encryption key.",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier of an instance that is suitable for external consumption.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
			},
			"links": schemaForLinks(),
		},
	}
}

func ResourceNutanixSnmpUserV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	username := d.Get("username").(string)

	body := cmgmtConfig.NewSnmpUser()
	body.Username = utils.StringPtr(username)
	body.AuthKey = utils.StringPtr(d.Get("auth_key").(string))
	body.AuthType = common.ExpandEnum[cmgmtConfig.SnmpAuthType](d.Get("auth_type").(string))

	if v, ok := d.GetOk("priv_type"); ok && v.(string) != "" {
		body.PrivType = common.ExpandEnum[cmgmtConfig.SnmpPrivType](v.(string))
	}
	if v, ok := d.GetOk("priv_key"); ok && v.(string) != "" {
		body.PrivKey = utils.StringPtr(v.(string))
	}

	bJSON, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("[DEBUG] Create SNMP User body: %s", string(bJSON))

	req := cmgmtRequest.CreateSnmpUserRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		Body:         body,
	}
	resp, err := conn.ClusterEntityAPI.CreateSnmpUser(ctx, &req)
	if err != nil {
		return diag.Errorf("error while creating SNMP user for cluster (%s): %v", clusterExtID, err)
	}

	taskRef, ok := resp.Data.GetValue().(prismCfg.TaskReference)
	if !ok {
		return diag.Errorf("unexpected create SNMP user response data type")
	}
	taskUUID := utils.StringValue(taskRef.ExtId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, taskUUID),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for SNMP user create task (%s): %v", taskUUID, errWait)
	}

	// CreateSnmpUser returns only a TaskReference; resolve the new user's extId by
	// querying the SNMP config and matching on the unique username.
	extID, errLookup := lookupSnmpUserExtIDByUsername(ctx, conn, clusterExtID, username)
	if errLookup != nil {
		return diag.Errorf("error resolving created SNMP user ext_id: %v", errLookup)
	}
	d.SetId(extID)

	return ResourceNutanixSnmpUserV2Read(ctx, d, meta)
}

func ResourceNutanixSnmpUserV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	req := cmgmtRequest.GetSnmpUserByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	resp, err := conn.ClusterEntityAPI.GetSnmpUserById(ctx, &req)
	if err != nil {
		return diag.Errorf("error while reading SNMP user (%s) for cluster (%s): %v", extID, clusterExtID, err)
	}

	user, ok := resp.Data.GetValue().(cmgmtConfig.SnmpUser)
	if !ok {
		return diag.Errorf("unexpected response data type when reading SNMP user")
	}

	if err := d.Set("ext_id", utils.StringValue(user.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(user.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(user.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("username", utils.StringValue(user.Username)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("auth_type", common.FlattenPtrEnum(user.AuthType)); err != nil {
		return diag.FromErr(err)
	}
	if user.AuthKey != nil {
		if err := d.Set("auth_key", utils.StringValue(user.AuthKey)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("priv_type", common.FlattenPtrEnum(user.PrivType)); err != nil {
		return diag.FromErr(err)
	}
	if user.PrivKey != nil {
		if err := d.Set("priv_key", utils.StringValue(user.PrivKey)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func ResourceNutanixSnmpUserV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	// Read current state for ETag.
	getReq := cmgmtRequest.GetSnmpUserByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	getResp, err := conn.ClusterEntityAPI.GetSnmpUserById(ctx, &getReq)
	if err != nil {
		return diag.Errorf("error while reading SNMP user before update (%s): %v", extID, err)
	}
	args := make(map[string]interface{})
	etag := conn.ClusterEntityAPI.ApiClient.GetEtag(getResp)
	args["If-Match"] = utils.StringPtr(etag)

	body := cmgmtConfig.NewSnmpUser()
	body.ExtId = utils.StringPtr(extID)
	body.Username = utils.StringPtr(d.Get("username").(string))
	body.AuthKey = utils.StringPtr(d.Get("auth_key").(string))
	body.AuthType = common.ExpandEnum[cmgmtConfig.SnmpAuthType](d.Get("auth_type").(string))

	if v, ok := d.GetOk("priv_type"); ok && v.(string) != "" {
		body.PrivType = common.ExpandEnum[cmgmtConfig.SnmpPrivType](v.(string))
	}
	if v, ok := d.GetOk("priv_key"); ok && v.(string) != "" {
		body.PrivKey = utils.StringPtr(v.(string))
	}

	updReq := cmgmtRequest.UpdateSnmpUserByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
		Body:         body,
	}
	resp, err := conn.ClusterEntityAPI.UpdateSnmpUserById(ctx, &updReq, args)
	if err != nil {
		return diag.Errorf("error while updating SNMP user (%s): %v", extID, err)
	}

	taskRef, ok := resp.Data.GetValue().(prismCfg.TaskReference)
	if !ok {
		return diag.Errorf("unexpected update SNMP user response data type")
	}
	taskUUID := utils.StringValue(taskRef.ExtId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, taskUUID),
		Timeout: d.Timeout(schema.TimeoutUpdate),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for SNMP user update task (%s): %v", taskUUID, errWait)
	}

	return ResourceNutanixSnmpUserV2Read(ctx, d, meta)
}

func ResourceNutanixSnmpUserV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI
	taskconn := meta.(*conns.Client).PrismAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Id()

	delReq := cmgmtRequest.DeleteSnmpUserByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	resp, err := conn.ClusterEntityAPI.DeleteSnmpUserById(ctx, &delReq)
	if err != nil {
		return diag.Errorf("error while deleting SNMP user (%s): %v", extID, err)
	}

	taskRef, ok := resp.Data.GetValue().(prismCfg.TaskReference)
	if !ok {
		return diag.Errorf("unexpected delete SNMP user response data type")
	}
	taskUUID := utils.StringValue(taskRef.ExtId)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, taskUUID),
		Timeout: d.Timeout(schema.TimeoutDelete),
	}
	if _, errWait := stateConf.WaitForStateContext(ctx); errWait != nil {
		return diag.Errorf("error waiting for SNMP user delete task (%s): %v", taskUUID, errWait)
	}

	return nil
}
