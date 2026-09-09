package lifecyclev2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	import3 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	prismConfigLc "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixConfigureNodeV2 configures a node in the hardware provider with the given settings.
// This is a standalone mutating action resource (no lifecycle Get/List pair).
func ResourceNutanixConfigureNodeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixConfigureNodeV2Create,
		ReadContext:   ResourceNutanixConfigureNodeV2Read,
		DeleteContext: ResourceNutanixConfigureNodeV2Delete,
		Schema: map[string]*schema.Schema{
			// configure_server: configures a single server in the hardware provider.
			"configure_server": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server_settings_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_adaptor_fec_mode": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"CL91", "CL74", "OFF"}, false),
									},
									"server_identity_pool_ext_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			// un_configure_server: removes the configuration from the given node.
			"un_configure_server": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			// pre_cluster_config: pre-cluster operations on a set of nodes.
			"pre_cluster_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": nodeInfoListSchema(),
					},
				},
			},
			// post_cluster_config: post-cluster operations on a set of nodes.
			"post_cluster_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"nodes": nodeInfoListSchema(),
					},
				},
			},
			"ext_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func nodeInfoListSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ext_id": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func ResourceNutanixConfigureNodeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	spec, err := expandConfigureNodeSpec(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.NodesAPIInstance.ConfigureNode(spec)
	if err != nil {
		return diag.Errorf("error while configuring node : %v", err)
	}

	taskRef := resp.Data.GetValue().(prismConfigLc.TaskReference)
	taskUUID := taskRef.ExtId

	taskconn := meta.(*conns.Client).PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskconn, utils.StringValue(taskUUID)),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}
	if _, errWaitTask := stateConf.WaitForStateContext(ctx); errWaitTask != nil {
		return diag.Errorf("error waiting for node configure (%s) to complete: %s", utils.StringValue(taskUUID), errWaitTask)
	}

	d.SetId(utils.StringValue(taskUUID))
	if err := d.Set("ext_id", utils.StringValue(taskUUID)); err != nil {
		return diag.FromErr(err)
	}
	return ResourceNutanixConfigureNodeV2Read(ctx, d, meta)
}

// ResourceNutanixConfigureNodeV2Read is a no-op; a configure action has no fetchable state of its own.
func ResourceNutanixConfigureNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixConfigureNodeV2Delete clears the resource from state; a configure action cannot be undone here.
func ResourceNutanixConfigureNodeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func expandConfigureNodeSpec(d *schema.ResourceData) (*import3.ConfigureNodeSpec, error) {
	spec := import3.NewConfigureNodeSpec()

	if v, ok := d.GetOk("configure_server"); ok && len(v.([]interface{})) > 0 {
		details := expandConfigureServerDetails(v.([]interface{}))
		if err := spec.Configuration.SetValue(*details); err != nil {
			return nil, err
		}
		return spec, nil
	}
	if v, ok := d.GetOk("un_configure_server"); ok && len(v.([]interface{})) > 0 {
		details := expandUnConfigureServerDetails(v.([]interface{}))
		if err := spec.Configuration.SetValue(*details); err != nil {
			return nil, err
		}
		return spec, nil
	}
	if v, ok := d.GetOk("pre_cluster_config"); ok && len(v.([]interface{})) > 0 {
		details := expandPreClusterConfigDetails(v.([]interface{}))
		if err := spec.Configuration.SetValue(*details); err != nil {
			return nil, err
		}
		return spec, nil
	}
	if v, ok := d.GetOk("post_cluster_config"); ok && len(v.([]interface{})) > 0 {
		details := expandPostClusterConfigDetails(v.([]interface{}))
		if err := spec.Configuration.SetValue(*details); err != nil {
			return nil, err
		}
		return spec, nil
	}

	return nil, fmt.Errorf("exactly one configuration block (configure_server, un_configure_server, pre_cluster_config, post_cluster_config) must be set")
}

func expandConfigureServerDetails(pr []interface{}) *import3.ConfigureServerDetails {
	val := pr[0].(map[string]interface{})
	details := import3.NewConfigureServerDetails()
	if v, ok := val["node_ext_id"]; ok && v.(string) != "" {
		details.NodeExtId = utils.StringPtr(v.(string))
	}
	if v, ok := val["group_id"]; ok && v.(string) != "" {
		details.GroupId = utils.StringPtr(v.(string))
	}
	if v, ok := val["server_settings_config"]; ok && len(v.([]interface{})) > 0 {
		details.ServerSettingsConfig = expandServerSettingsConfig(v.([]interface{}))
	}
	return details
}

func expandServerSettingsConfig(pr []interface{}) *import3.ServerSettingsConfig {
	if len(pr) == 0 || pr[0] == nil {
		return nil
	}
	val := pr[0].(map[string]interface{})
	cfg := import3.NewServerSettingsConfig()
	if v, ok := val["network_adaptor_fec_mode"]; ok && v.(string) != "" {
		cfg.NetworkAdaptorFecMode = common.ExpandEnum[import3.NetworkAdaptorFecMode](v.(string))
	}
	if v, ok := val["server_identity_pool_ext_id"]; ok && v.(string) != "" {
		cfg.ServerIdentityPoolExtId = utils.StringPtr(v.(string))
	}
	return cfg
}

func expandUnConfigureServerDetails(pr []interface{}) *import3.UnConfigureServerDetails {
	val := pr[0].(map[string]interface{})
	details := import3.NewUnConfigureServerDetails()
	if v, ok := val["node_ext_id"]; ok && v.(string) != "" {
		details.NodeExtId = utils.StringPtr(v.(string))
	}
	return details
}

func expandPreClusterConfigDetails(pr []interface{}) *import3.PreClusterConfigDetails {
	val := pr[0].(map[string]interface{})
	details := import3.NewPreClusterConfigDetails()
	if v, ok := val["nodes"]; ok {
		details.Nodes = expandNodeInfoList(v.([]interface{}))
	}
	return details
}

func expandPostClusterConfigDetails(pr []interface{}) *import3.PostClusterConfigDetails {
	val := pr[0].(map[string]interface{})
	details := import3.NewPostClusterConfigDetails()
	if v, ok := val["cluster_ext_id"]; ok && v.(string) != "" {
		details.ClusterExtId = utils.StringPtr(v.(string))
	}
	if v, ok := val["nodes"]; ok {
		details.Nodes = expandNodeInfoList(v.([]interface{}))
	}
	return details
}

func expandNodeInfoList(pr []interface{}) []import3.NodeInfo {
	if len(pr) == 0 {
		return nil
	}
	out := make([]import3.NodeInfo, 0, len(pr))
	for _, item := range pr {
		if item == nil {
			continue
		}
		val := item.(map[string]interface{})
		info := import3.NewNodeInfo()
		if v, ok := val["ext_id"]; ok && v.(string) != "" {
			info.ExtId = utils.StringPtr(v.(string))
		}
		out = append(out, *info)
	}
	return out
}
