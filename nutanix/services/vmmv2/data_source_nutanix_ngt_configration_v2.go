package vmmv2

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vmmAhvConfig "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixNGTConfigurationV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixNGTConfigurationV4Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_iso_inserted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"capablities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"available_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guest_os_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_reachable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_vss_snapshot_capable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_vm_mobility_drivers_installed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixNGTConfigurationV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	extID := d.Get("ext_id")
	resp, err := conn.VMAPIInstance.GetGuestToolsById(utils.StringPtr(extID.(string)))
	log.Printf("[DEBUG] GetGuestToolsById : %v", resp)
	if err != nil {
		return diag.Errorf("error while fetching Gest Tool : %v", err)
	}
	getResp := resp.Data.GetValue().(vmmAhvConfig.GuestTools)

	if err := d.Set("version", getResp.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_installed", getResp.IsInstalled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_enabled", getResp.IsEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_iso_inserted", getResp.IsIsoInserted); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("capablities", flattenCapabilities(getResp.Capabilities)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("available_version", getResp.AvailableVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guest_os_version", getResp.GuestOsVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_reachable", getResp.IsReachable); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_vss_snapshot_capable", getResp.IsVssSnapshotCapable); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_vm_mobility_drivers_installed", getResp.IsVmMobilityDriversInstalled); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenCapabilities(capabilities []vmmAhvConfig.NgtCapability) []string {
	if len(capabilities) > 0 {
		c := make([]string, len(capabilities))
		for i, v := range capabilities {
			c[i] = v.GetName()
		}
		sort.Strings(c)
		log.Printf("[DEBUG] Read NGT Capabilities : %v", c)
		return c
	}
	log.Printf("[DEBUG] Read NGT Capabilities : nil")
	return nil
}
