package lifecyclev2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import3 "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/lifecycle/v4/config"
	prismConfigLc "github.com/nutanix/ntnx-api-golang-clients/lifecycle-go-client/v4/models/prism/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// ResourceNutanixImageNodeV2 installs AOS or hypervisor or host OS on the node.
// This is a standalone mutating action resource (no lifecycle Get/List pair).
func ResourceNutanixImageNodeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceNutanixImageNodeV2Create,
		ReadContext:   ResourceNutanixImageNodeV2Read,
		DeleteContext: ResourceNutanixImageNodeV2Delete,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// host_installation: install a patched hypervisor / host OS image.
			"host_installation": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"patched_image_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"patched_image_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// aos_installation: install AOS on the node.
			"aos_installation": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aos_image_ext_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cvm_memory_gb": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"management_network": imageNetworkDetailsSchema(),
					},
				},
			},
		},
	}
}

func imageNetworkDetailsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"gateway": ipAddressSchema(false),
				"ip":      ipAddressSchema(false),
				"vlan_id": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	}
}

func ResourceNutanixImageNodeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).LifecycleAPI

	extID := d.Get("ext_id").(string)

	spec, err := expandImageNodeSpec(d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := conn.NodesAPIInstance.ImageNode(utils.StringPtr(extID), spec)
	if err != nil {
		return diag.Errorf("error while imaging node : %v", err)
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
		return diag.Errorf("error waiting for node image (%s) to complete: %s", extID, errWaitTask)
	}

	d.SetId(extID)
	return ResourceNutanixImageNodeV2Read(ctx, d, meta)
}

// ResourceNutanixImageNodeV2Read is a no-op; an imaging action has no fetchable state of its own.
func ResourceNutanixImageNodeV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// ResourceNutanixImageNodeV2Delete clears the resource from state; imaging cannot be undone here.
func ResourceNutanixImageNodeV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func expandImageNodeSpec(d *schema.ResourceData) (*import3.ImageNodeSpec, error) {
	spec := import3.NewImageNodeSpec()

	if v, ok := d.GetOk("host_installation"); ok && len(v.([]interface{})) > 0 {
		details := expandHostInstallationDetails(v.([]interface{}))
		if err := spec.Configuration.SetValue(*details); err != nil {
			return nil, err
		}
		return spec, nil
	}
	if v, ok := d.GetOk("aos_installation"); ok && len(v.([]interface{})) > 0 {
		details := expandAosInstallationDetails(v.([]interface{}))
		if err := spec.Configuration.SetValue(*details); err != nil {
			return nil, err
		}
		return spec, nil
	}

	return nil, fmt.Errorf("exactly one configuration block (host_installation, aos_installation) must be set")
}

func expandHostInstallationDetails(pr []interface{}) *import3.HostInstallationDetails {
	val := pr[0].(map[string]interface{})
	details := import3.NewHostInstallationDetails()
	if v, ok := val["patched_image_ext_id"]; ok && v.(string) != "" {
		details.PatchedImageExtId = utils.StringPtr(v.(string))
	}
	if v, ok := val["patched_image_url"]; ok && v.(string) != "" {
		details.PatchedImageUrl = utils.StringPtr(v.(string))
	}
	return details
}

func expandAosInstallationDetails(pr []interface{}) *import3.AosInstallationDetails {
	val := pr[0].(map[string]interface{})
	details := import3.NewAosInstallationDetails()

	if v, ok := val["aos_image_ext_id"]; ok && v.(string) != "" {
		imageDetails := import3.NewLocalAOSImageDetails()
		imageDetails.ExtId = utils.StringPtr(v.(string))
		_ = details.AosImageDetails.SetValue(*imageDetails)
	}
	if v, ok := val["cvm_memory_gb"]; ok && v.(int) != 0 {
		details.CvmMemoryGB = utils.Int64Ptr(int64(v.(int)))
	}
	if v, ok := val["management_network"]; ok && len(v.([]interface{})) > 0 {
		networking := import3.NewAosNetworking()
		networking.Management = expandNetworkDetails(v.([]interface{}))
		details.NetworkDetails = networking
	}
	return details
}
