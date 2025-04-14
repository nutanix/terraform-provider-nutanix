package selfservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/calm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
)

func DatsourceNutanixCalmRuntimeEditables() *schema.Resource {
	return &schema.Resource{
		ReadContext: datsourceNutanixCalmRuntimeEditablesRead,
		Schema: map[string]*schema.Schema{
			"bp_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bp_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runtime_editables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"service_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"credential_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"substrate_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"package_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"snapshot_config_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"app_profile": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"task_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"restore_config_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"variable_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
						"deployment_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: RuntimeSpecDS(),
							},
						},
					},
				},
			},
		},
	}
}

func datsourceNutanixCalmRuntimeEditablesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).Calm

	var bpUUID string
	// fetch bp_uuid from bp_name
	bpName := d.Get("bp_name").(string)

	bpFilter := &calm.BlueprintListInput{}

	bpFilter.Filter = fmt.Sprintf("name==%s;state!=DELETED", bpName)

	bpNameResp, err := conn.Service.ListBlueprint(ctx, bpFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	var BpNameStatus []interface{}
	if err := json.Unmarshal([]byte(bpNameResp.Entities), &BpNameStatus); err != nil {
		fmt.Println("Error unmarshalling BPName:", err)
	}

	entities := BpNameStatus[0].(map[string]interface{})

	if entity, ok := entities["metadata"].(map[string]interface{}); ok {
		bpUUID = entity["uuid"].(string)
	}

	if bpUUID, ok := d.GetOk("bp_uuid"); ok {
		bpUUID = bpUUID.(string)
	}

	getRuntime, err := conn.Service.GetRuntimeEditables(ctx, bpUUID)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Println("Runtime Editables:", getRuntime)

	if err := d.Set("runtime_editables", flattenRuntimeEditables(getRuntime.Resources[0].RuntimeEditables)); err != nil {
		return diag.FromErr(err)
	}
	fmt.Println("res:", getRuntime.Resources[0].RuntimeEditables)

	d.SetId(bpUUID)
	return nil
}

func flattenRuntimeEditables(resource *calm.RuntimeEditables) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"action_list":     flattenRuntimeSpec(resource.ActionList),
			"service_list":    flattenRuntimeSpec(resource.ServiceList),
			"credential_list": flattenRuntimeSpec(resource.CredentialList),
			"substrate_list":  flattenRuntimeSpec(resource.SubstrateList),
			"package_list":    flattenRuntimeSpec(resource.PackageList),
			// "snapshot_config_list": flattenRuntimeSpec(resource.SnapshotConfigList),
			// "app_profile":          flattenRuntimeSpec(resource.AppProfile),
			"task_list":           flattenRuntimeSpec(resource.TaskList),
			"restore_config_list": flattenRuntimeSpec(resource.RestoreConfigList),
			"variable_list":       flattenRuntimeSpec(resource.VariableList),
			"deployment_list":     flattenRuntimeSpec(resource.DeploymentList),
		},
	}
}

func flattenRuntimeSpec(pr []*calm.RuntimeSpec) []interface{} {
	if pr == nil {
		return nil
	}

	var runtimeSpec []interface{}
	for _, r := range pr {
		runtimeSpec = append(runtimeSpec, map[string]interface{}{
			"uuid":        r.UUID,
			"name":        r.Name,
			"description": r.Description,
			"type":        r.Type,
			"context":     r.Context,
			"value": func() string {
				data := json.RawMessage(*r.Value)
				return string(data)
			}(),
		})
	}
	return runtimeSpec
}

func RuntimeSpecDS() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"uuid": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"context": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}
