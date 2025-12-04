package vmmv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixVirtualMachinesV4() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixVirtualMachinesV4Read,
		Schema: map[string]*schema.Schema{
			"page": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixVirtualMachineV4(),
			},
		},
	}
}

func DatasourceNutanixVirtualMachinesV4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).VmmAPI

	// initialize query params
	var filter, orderBy, selects *string
	var page, limit *int

	if pagef, ok := d.GetOk("page"); ok {
		page = utils.IntPtr(pagef.(int))
	} else {
		page = nil
	}
	if limitf, ok := d.GetOk("limit"); ok {
		limit = utils.IntPtr(limitf.(int))
	} else {
		limit = nil
	}
	if filterf, ok := d.GetOk("filter"); ok {
		filter = utils.StringPtr(filterf.(string))
	} else {
		filter = nil
	}
	if order, ok := d.GetOk("order_by"); ok {
		orderBy = utils.StringPtr(order.(string))
	} else {
		orderBy = nil
	}
	if selectf, ok := d.GetOk("select"); ok {
		selects = utils.StringPtr(selectf.(string))
	} else {
		selects = nil
	}
	resp, err := conn.VMAPIInstance.ListVms(page, limit, filter, orderBy, selects)
	if err != nil {
		return diag.Errorf("error while fetching vms : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("vms", make([]interface{}, 0)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.GenUUID())

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "🫙 No data found.",
			Detail:   "The API returned an empty list of virtual machines.",
		}}
	}

	getResp := resp.Data.GetValue().([]config.Vm)

	if err := d.Set("vms", flattenVMEntities(getResp)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenVMEntities(vms []config.Vm) []interface{} {
	if len(vms) > 0 {
		vmsList := make([]interface{}, len(vms))

		for k, v := range vms {
			vm := make(map[string]interface{})

			if v.ExtId != nil {
				vm["ext_id"] = v.ExtId
			}
			if v.Name != nil {
				vm["name"] = v.Name
			}
			if v.Description != nil {
				vm["description"] = v.Description
			}
			if v.CreateTime != nil {
				t := v.CreateTime
				vm["create_time"] = t.String()
			}
			if v.UpdateTime != nil {
				t := v.UpdateTime
				vm["update_time"] = t.String()
			}
			if v.Source != nil {
				vm["source"] = flattenVMSourceReference(v.Source)
			}
			if v.NumSockets != nil {
				vm["num_sockets"] = v.NumSockets
			}
			if v.NumCoresPerSocket != nil {
				vm["num_cores_per_socket"] = v.NumCoresPerSocket
			}
			if v.NumThreadsPerCore != nil {
				vm["num_threads_per_core"] = v.NumThreadsPerCore
			}
			if v.NumNumaNodes != nil {
				vm["num_numa_nodes"] = v.NumNumaNodes
			}
			if v.MemorySizeBytes != nil {
				vm["memory_size_bytes"] = v.MemorySizeBytes
			}
			if v.IsVcpuHardPinningEnabled != nil {
				vm["is_vcpu_hard_pinning_enabled"] = v.IsVcpuHardPinningEnabled
			}
			if v.IsCpuPassthroughEnabled != nil {
				vm["is_cpu_passthrough_enabled"] = v.IsCpuPassthroughEnabled
			}
			if v.EnabledCpuFeatures != nil {
				vm["enabled_cpu_features"] = flattenCPUFeature(v.EnabledCpuFeatures)
			}
			if v.IsMemoryOvercommitEnabled != nil {
				vm["is_memory_overcommit_enabled"] = v.IsMemoryOvercommitEnabled
			}
			if v.IsGpuConsoleEnabled != nil {
				vm["is_gpu_console_enabled"] = v.IsGpuConsoleEnabled
			}
			if v.IsCpuHotplugEnabled != nil {
				vm["is_cpu_hotplug_enabled"] = v.IsCpuHotplugEnabled
			}
			if v.IsScsiControllerEnabled != nil {
				vm["is_scsi_controller_enabled"] = v.IsScsiControllerEnabled
			}
			if v.GenerationUuid != nil {
				vm["generation_uuid"] = v.GenerationUuid
			}
			if v.BiosUuid != nil {
				vm["bios_uuid"] = v.BiosUuid
			}
			if v.Categories != nil {
				vm["categories"] = flattenCategoryReference(v.Categories)
			}
			if v.Project != nil {
				vm["project"] = flattenProjectReference(v.Project)
			}
			if v.OwnershipInfo != nil {
				vm["ownership_info"] = flattenOwnershipInfo(v.OwnershipInfo)
			}
			if v.Host != nil {
				vm["host"] = flattenHostReference(v.Host)
			}
			if v.Cluster != nil {
				vm["cluster"] = flattenClusterReference(v.Cluster)
			}
			if v.GuestCustomization != nil {
				vm["guest_customization"] = flattenGuestCustomizationParams(v.GuestCustomization)
			}
			if v.GuestTools != nil {
				vm["guest_tools"] = flattenGuestTools(v.GuestTools)
			}
			if v.HardwareClockTimezone != nil {
				vm["hardware_clock_timezone"] = v.HardwareClockTimezone
			}
			if v.IsBrandingEnabled != nil {
				vm["is_branding_enabled"] = v.IsBrandingEnabled
			}
			if v.BootConfig != nil {
				vm["boot_config"] = flattenOneOfVMBootConfig(v.BootConfig)
			}
			if v.IsVgaConsoleEnabled != nil {
				vm["is_vga_console_enabled"] = v.IsVgaConsoleEnabled
			}
			if v.MachineType != nil {
				vm["machine_type"] = flattenMachineType(v.MachineType)
			}
			if v.PowerState != nil {
				vm["power_state"] = flattenPowerState(v.PowerState)
			}
			if v.VtpmConfig != nil {
				vm["vtpm_config"] = flattenVtpmConfig(v.VtpmConfig)
			}
			if v.IsAgentVm != nil {
				vm["is_agent_vm"] = v.IsAgentVm
			}
			if v.ApcConfig != nil {
				vm["apc_config"] = flattenApcConfig(v.ApcConfig)
			}
			if v.StorageConfig != nil {
				vm["storage_config"] = flattenADSFVmStorageConfig(v.StorageConfig)
			}
			if v.Disks != nil {
				vm["disks"] = flattenDisk(v.Disks)
			}
			if v.CdRoms != nil {
				vm["cd_roms"] = flattenCdRom(v.CdRoms)
			}
			if v.Nics != nil {
				vm["nics"] = flattenNic(v.Nics)
			}
			if v.Gpus != nil {
				vm["gpus"] = flattenGpu(v.Gpus)
			}
			if v.SerialPorts != nil {
				vm["serial_ports"] = flattenSerialPort(v.SerialPorts)
			}
			if v.ProtectionType != nil {
				vm["protection_type"] = flattenProtectionType(v.ProtectionType)
			}
			if v.ProtectionPolicyState != nil {
				vm["protection_policy_state"] = flattenProtectionPolicyState(v.ProtectionPolicyState)
			}

			vmsList[k] = vm
		}
		return vmsList
	}
	return nil
}
