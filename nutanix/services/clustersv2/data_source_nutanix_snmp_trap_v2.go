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

// DatasourceNutanixSnmpTrapV2 is a singular datasource for fetching an SNMP trap configuration
// identified by {extId} associated with the cluster identified by {clusterExtId}.
func DatasourceNutanixSnmpTrapV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSnmpTrapV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the UUID of a cluster.",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SNMP trap UUID.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).",
			},
			"links": schemaForLinks(),
			"address": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Address of the SNMP trap receiver.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv4": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IPv4 address of the host.",
									},
									"prefix_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The prefix length of the network to which this host IPv4 address belongs.",
									},
								},
							},
						},
						"ipv6": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IPv6 address of the host.",
									},
									"prefix_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The prefix length of the network to which this host IPv6 address belongs.",
									},
								},
							},
						},
					},
				},
			},
			"community_string": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Community string(plaintext) for SNMP version 2.0.",
			},
			"engine_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP engine Id.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "SNMP port.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP protocol.",
			},
			"receiver_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP receiver name.",
			},
			"should_inform": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "SNMP information status.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP username. For SNMP trap v3 version, SNMP username is required parameter.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNMP version.",
			},
		},
	}
}

func DatasourceNutanixSnmpTrapV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)
	extID := d.Get("ext_id").(string)

	req := cmgmtRequest.GetSnmpTrapByIdRequest{
		ClusterExtId: utils.StringPtr(clusterExtID),
		ExtId:        utils.StringPtr(extID),
	}
	resp, err := conn.ClusterEntityAPI.GetSnmpTrapById(ctx, &req)
	if err != nil {
		return diag.Errorf("error while fetching SNMP trap (%s) for cluster (%s): %v", extID, clusterExtID, err)
	}

	trap, ok := resp.Data.GetValue().(config.SnmpTrap)
	if !ok {
		return diag.Errorf("unexpected response data type when fetching SNMP trap")
	}

	aJSON, _ := json.MarshalIndent(trap, "", "  ")
	log.Printf("[DEBUG] Read SNMP Trap: %s", string(aJSON))

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

	d.SetId(utils.StringValue(trap.ExtId))
	return nil
}
