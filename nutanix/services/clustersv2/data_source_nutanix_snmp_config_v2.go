package clustersv2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// DatasourceNutanixSnmpConfigV2 is a singular datasource for fetching the full
// SNMP configuration of the cluster identified by {cluster_ext_id}. It exposes
// the SNMP enable status together with all configured transports, traps and
// users in a single read.
func DatasourceNutanixSnmpConfigV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixSnmpConfigV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the UUID of a cluster.",
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
			"is_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "SNMP status.",
			},
			"transports": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SNMP transport details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"traps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SNMP trap details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A globally unique identifier of an instance that is suitable for external consumption.",
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
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
													Type:     schema.TypeString,
													Computed: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Computed: true,
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
													Type:     schema.TypeString,
													Computed: true,
												},
												"prefix_length": {
													Type:     schema.TypeInt,
													Computed: true,
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
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "SNMP user information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ext_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
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
						"priv_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SNMP user encryption type.",
						},
					},
				},
			},
		},
	}
}

func DatasourceNutanixSnmpConfigV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	clusterExtID := d.Get("cluster_ext_id").(string)

	cfg, err := fetchSnmpConfig(ctx, conn, clusterExtID)
	if err != nil {
		return diag.Errorf("error while fetching SNMP config for cluster (%s): %v", clusterExtID, err)
	}

	aJSON, _ := json.MarshalIndent(cfg, "", "  ")
	log.Printf("[DEBUG] Read SNMP Config: %s", string(aJSON))

	if err := d.Set("ext_id", utils.StringValue(cfg.ExtId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", utils.StringValue(cfg.TenantId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", flattenLinks(cfg.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", utils.BoolValue(cfg.IsEnabled)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("transports", flattenSnmpTransports(cfg.Transports)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("traps", flattenSnmpTraps(cfg.Traps)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("users", flattenSnmpUsers(cfg.Users)); err != nil {
		return diag.FromErr(err)
	}

	// SnmpConfig is a per-cluster singleton. Use the cluster ext_id as the
	// Terraform resource ID so the data source is stable across reads even
	// when the upstream config has no extId of its own.
	if extID := utils.StringValue(cfg.ExtId); extID != "" {
		d.SetId(extID)
	} else {
		d.SetId(clusterExtID)
	}
	return nil
}

func flattenSnmpTransports(in []config.SnmpTransport) []map[string]interface{} {
	if len(in) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(in))
	for _, t := range in {
		out = append(out, map[string]interface{}{
			"port":     utils.IntValue(t.Port),
			"protocol": common.FlattenPtrEnum(t.Protocol),
		})
	}
	return out
}

func flattenSnmpTraps(in []config.SnmpTrap) []map[string]interface{} {
	if len(in) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(in))
	for _, t := range in {
		out = append(out, map[string]interface{}{
			"ext_id":           utils.StringValue(t.ExtId),
			"tenant_id":        utils.StringValue(t.TenantId),
			"links":            flattenLinks(t.Links),
			"address":          flattenIPAddress(t.Address),
			"community_string": utils.StringValue(t.CommunityString),
			"engine_id":        utils.StringValue(t.EngineId),
			"port":             utils.IntValue(t.Port),
			"protocol":         common.FlattenPtrEnum(t.Protocol),
			"receiver_name":    utils.StringValue(t.RecieverName),
			"should_inform":    utils.BoolValue(t.ShouldInform),
			"username":         utils.StringValue(t.Username),
			"version":          common.FlattenPtrEnum(t.Version),
		})
	}
	return out
}

func flattenSnmpUsers(in []config.SnmpUser) []map[string]interface{} {
	if len(in) == 0 {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(in))
	for _, u := range in {
		out = append(out, map[string]interface{}{
			"ext_id":    utils.StringValue(u.ExtId),
			"tenant_id": utils.StringValue(u.TenantId),
			"links":     flattenLinks(u.Links),
			"username":  utils.StringValue(u.Username),
			"auth_type": common.FlattenPtrEnum(u.AuthType),
			"priv_type": common.FlattenPtrEnum(u.PrivType),
		})
	}
	return out
}
