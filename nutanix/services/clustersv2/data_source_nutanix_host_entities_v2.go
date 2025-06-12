package clustersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixHostEntitiesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixHostEntitiesV2Read,
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
			"apply": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     DatasourceNutanixHostEntityV2(),
			},
		},
	}
}

func DatasourceNutanixHostEntitiesV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	// initialize query params
	var filter, orderBy, apply, selectQ *string
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
	if applyf, ok := d.GetOk("apply"); ok {
		apply = utils.StringPtr(applyf.(string))
	} else {
		apply = nil
	}
	if selectQy, ok := d.GetOk("apply"); ok {
		selectQ = utils.StringPtr(selectQy.(string))
	} else {
		selectQ = nil
	}

	resp, err := conn.ClusterEntityAPI.ListHosts(page, limit, filter, orderBy, apply, selectQ)
	if err != nil {
		return diag.Errorf("error while fetching host entities : %v", err)
	}

	if resp.Data == nil {
		if err := d.Set("host_entities", []map[string]interface{}{}); err != nil {
			return diag.FromErr(err)
		}
	} else {
		getResp := resp.Data.GetValue().([]import1.Host)

		if err := d.Set("host_entities", flattenHostEntities(getResp)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenHostEntities(pr []import1.Host) []interface{} {
	if len(pr) > 0 {
		hostList := make([]interface{}, len(pr))

		for k, v := range pr {
			host := make(map[string]interface{})

			host["ext_id"] = v.ExtId
			host["tenant_id"] = v.TenantId
			host["links"] = flattenLinks(v.Links)
			host["host_name"] = v.HostName
			host["host_type"] = flattenHostTypeEnum(v.HostType)
			host["hypervisor"] = flattenHypervisorReference(v.Hypervisor)
			host["cluster"] = flattenClusterReference(v.Cluster)
			host["controller_vm"] = flattenControllerVMReference(v.ControllerVm)
			host["disk"] = flattenDiskReference(v.Disk)
			host["is_degraded"] = v.IsDegraded
			host["is_secure_booted"] = v.IsSecureBooted
			host["is_hardware_virtualized"] = v.IsHardwareVirtualized
			host["has_csr"] = v.HasCsr
			host["key_management_device_to_cert_status"] = flattenKeyManagementDeviceToCertStatusInfo(v.KeyManagementDeviceToCertStatus)
			host["number_of_cpu_cores"] = v.NumberOfCpuCores
			host["number_of_cpu_threads"] = v.NumberOfCpuThreads
			host["number_of_cpu_sockets"] = v.NumberOfCpuSockets
			host["cpu_capacity_hz"] = v.CpuCapacityHz
			host["cpu_frequency_hz"] = v.CpuFrequencyHz
			host["cpu_model"] = v.CpuModel
			host["gpu_driver_version"] = v.GpuDriverVersion
			host["gpu_list"] = v.GpuList
			host["default_vhd_location"] = v.DefaultVhdLocation
			host["default_vhd_container_uuid"] = v.DefaultVhdContainerUuid
			host["default_vm_location"] = v.DefaultVmLocation
			host["default_vm_container_uuid"] = v.DefaultVmContainerUuid
			host["is_reboot_pending"] = v.IsRebootPending
			host["failover_cluster_fqdn"] = v.FailoverClusterFqdn
			host["failover_cluster_node_status"] = v.FailoverClusterNodeStatus
			host["boot_time_usecs"] = v.BootTimeUsecs
			host["memory_size_bytes"] = v.MemorySizeBytes
			host["block_serial"] = v.BlockSerial
			host["block_model"] = v.BlockModel
			host["maintenance_state"] = v.MaintenanceState
			host["node_status"] = flattenNodeStatus(v.NodeStatus)

			hostList[k] = host
		}
		return hostList
	}
	return nil
}
