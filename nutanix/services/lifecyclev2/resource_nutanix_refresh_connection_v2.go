package lifecyclev2

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	taskRef "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixRefreshConnectionV2 is an action resource that refreshes nodes or
// resources (IP/MAC/server-identity pools) from a hardware provider connection. This
// is a state-changing, one-shot action with no lifecycle pair.
func ResourceNutanixRefreshConnectionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixRefreshConnectionV2Create,
		ReadContext:   ResourceNutanixRefreshConnectionV2Read,
		DeleteContext: ResourceNutanixRefreshConnectionV2Delete,
		Schema: map[string]*schema.Schema{
			"hardware_provider_ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "External ID of the hardware provider",
			},
			"ext_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "External ID of the hardware provider connection to refresh",
			},
			"refresh_resources_spec": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Refresh resources (IP/MAC/server-identity pools) from the connection",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"should_refresh_ip_pools": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"should_refresh_mac_pools": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"should_refresh_server_identity_pools": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"discover_nodes_spec": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Discover nodes from the connection",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func ResourceNutanixRefreshConnectionV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	hpExtID := d.Get("hardware_provider_ext_id").(string)
	extID := d.Get("ext_id").(string)

	body := import1.NewRefreshConnectionSpec()
	if config := expandRefreshConnectionConfiguration(d); config != nil {
		body.Configuration = config
	}

	aJSON, _ := json.MarshalIndent(body, "", " ")
	log.Printf("[DEBUG] Refresh connection payload : %s", string(aJSON))

	resp, err := conn.HardwareProvidersAPIInstance.RefreshConnection(utils.StringPtr(hpExtID), utils.StringPtr(extID), body)
	if err != nil {
		return diag.Errorf("error while refreshing connection : %v", err)
	}

	TaskRef := resp.Data.GetValue().(taskRef.TaskReference)
	taskUUID := TaskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"QUEUED", "RUNNING", "PENDING"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for connection (%s) to refresh: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func ResourceNutanixRefreshConnectionV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Refresh is a one-shot action with no persistent entity to read back.
	return nil
}

func ResourceNutanixRefreshConnectionV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Action resource: nothing to delete. The framework clears the ID after this returns.
	return nil
}

// expandRefreshConnectionConfiguration builds the OneOfRefreshConnectionSpecConfiguration
// wrapper from the refresh_resources_spec / discover_nodes_spec schema blocks.
func expandRefreshConnectionConfiguration(d *schema.ResourceData) *import1.OneOfRefreshConnectionSpecConfiguration {
	if v, ok := d.GetOk("refresh_resources_spec"); ok && len(v.([]interface{})) > 0 {
		spec := import1.NewRefreshResourcesSpec()
		m := v.([]interface{})[0].(map[string]interface{})
		if gids, ok := m["group_ids"]; ok {
			spec.GroupIds = common.ExpandListOfString(gids.([]interface{}))
		}
		if val, ok := m["should_refresh_ip_pools"]; ok {
			spec.ShouldRefreshIpPools = utils.BoolPtr(val.(bool))
		}
		if val, ok := m["should_refresh_mac_pools"]; ok {
			spec.ShouldRefreshMacPools = utils.BoolPtr(val.(bool))
		}
		if val, ok := m["should_refresh_server_identity_pools"]; ok {
			spec.ShouldRefreshServerIdentityPools = utils.BoolPtr(val.(bool))
		}
		config := import1.NewOneOfRefreshConnectionSpecConfiguration()
		if err := config.SetValue(*spec); err == nil {
			return config
		}
	}

	if v, ok := d.GetOk("discover_nodes_spec"); ok && len(v.([]interface{})) > 0 {
		spec := import1.NewDiscoverNodesSpec()
		m := v.([]interface{})[0].(map[string]interface{})
		if gids, ok := m["group_ids"]; ok {
			spec.GroupIds = common.ExpandListOfString(gids.([]interface{}))
		}
		config := import1.NewOneOfRefreshConnectionSpecConfiguration()
		if err := config.SetValue(*spec); err == nil {
			return config
		}
	}

	return nil
}
