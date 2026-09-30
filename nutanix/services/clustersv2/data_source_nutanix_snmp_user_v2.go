package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	cmgmtRequest "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/request/clusters"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixSnmpUserV2 is a singular datasource for fetching an SNMP user
// configuration identified by {extId} associated with the cluster identified by {clusterExtId}.
func DatasourceNutanixSnmpUserV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSnmpUserV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the UUID of a cluster.",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SNMP user UUID.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
			},
			"links": schemaForLinks(),
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP username. For SNMP trap v3 version, SNMP username is required parameter.",
			},
			"auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP user authentication type.",
			},
			"auth_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "SNMP user authentication key.",
			},
			"priv_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP user encryption type.",
			},
			"priv_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "SNMP user encryption key.",
			},
		},
	}
}

func DatasourceNutanixSnmpUserV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Get("ext_id").(string)

	req := cmgmtRequest.GetSnmpUserByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	resp, err := conn.ClusterEntityAPI.GetSnmpUserById(ctx, &req)
	if err != nil {
		return diag.Errorf("error while fetching SNMP user (%s) for cluster (%s): %v", extID, clusterExtID, err)
	}

	user, ok := resp.Data.GetValue().(config.SnmpUser)
	if !ok {
		return diag.Errorf("unexpected response data type when fetching SNMP user")
	}

	aJSON, _ := json.MarshalIndent(user, "", "  ")
	log.Printf("[DEBUG] Read SNMP User: %s", string(aJSON))

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
	if err := d.Set("auth_key", utils.StringValue(user.AuthKey)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("priv_type", common.FlattenPtrEnum(user.PrivType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("priv_key", utils.StringValue(user.PrivKey)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.StringValue(user.ExtId))
	return nil
}
