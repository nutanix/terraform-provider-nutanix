package clustersv2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import1 "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/common/v1/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func DatasourceNutanixHostEntityV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: DatasourceNutanixHostEntityV2Read,
		Schema: map[string]*schema.Schema{
			"ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_ext_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": common.LinksSchema(),
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"number_of_vms": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acropolis_connection_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"cluster": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"controller_vm": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"external_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"backplane_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"rdma_backplane_address": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"nat_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"nat_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"maintenance_mode": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"disk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size_in_bytes": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"serial_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_tier": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"is_degraded": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_secure_booted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_hardware_virtualized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"has_csr": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"key_management_device_to_cert_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_management_server_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_certificate_present": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"number_of_cpu_cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"number_of_cpu_threads": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"number_of_cpu_sockets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_capacity_hz": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_frequency_hz": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_driver_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gpu_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_vhd_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_vhd_container_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_vm_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_vm_container_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_reboot_pending": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"failover_cluster_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failover_cluster_node_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"boot_time_usecs": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"block_serial": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipmi": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     common.SchemaForIPList(false),
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"rackable_unit_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DatasourceNutanixHostEntityV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.Client).ClusterAPI

	extID := d.Get("ext_id")
	clsID := d.Get("cluster_ext_id")
	resp, err := conn.ClusterEntityAPI.GetHostById(utils.StringPtr(clsID.(string)), utils.StringPtr(extID.(string)))
	if err != nil {
		return diag.Errorf("error while fetching host entity : %v", err)
	}

	getResp := resp.Data.GetValue().(import1.Host)

	if err := d.Set("tenant_id", getResp.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("links", common.FlattenLinks(getResp.Links)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_name", getResp.HostName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host_type", flattenHostTypeEnum(getResp.HostType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("hypervisor", flattenHypervisorReference(getResp.Hypervisor)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cluster", flattenClusterReference(getResp.Cluster)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("controller_vm", flattenControllerVMReference(getResp.ControllerVm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk", flattenDiskReference(getResp.Disk)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_degraded", getResp.IsDegraded); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_secure_booted", getResp.IsSecureBooted); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_hardware_virtualized", getResp.IsHardwareVirtualized); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("has_csr", getResp.HasCsr); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key_management_device_to_cert_status", flattenKeyManagementDeviceToCertStatusInfo(getResp.KeyManagementDeviceToCertStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("number_of_cpu_cores", getResp.NumberOfCpuCores); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("number_of_cpu_threads", getResp.NumberOfCpuThreads); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("number_of_cpu_sockets", getResp.NumberOfCpuSockets); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_capacity_hz", getResp.CpuCapacityHz); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_frequency_hz", getResp.CpuFrequencyHz); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_model", getResp.CpuModel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpu_driver_version", getResp.GpuDriverVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("gpu_list", getResp.GpuList); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_vhd_location", getResp.DefaultVhdLocation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_vhd_container_uuid", getResp.DefaultVhdContainerUuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_vm_location", getResp.DefaultVmLocation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_vm_container_uuid", getResp.DefaultVmContainerUuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_reboot_pending", getResp.IsRebootPending); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("failover_cluster_fqdn", getResp.FailoverClusterFqdn); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("failover_cluster_node_status", getResp.FailoverClusterNodeStatus); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("boot_time_usecs", getResp.BootTimeUsecs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("memory_size_bytes", getResp.MemorySizeBytes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("block_serial", getResp.BlockSerial); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("block_model", getResp.BlockModel); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("maintenance_state", getResp.MaintenanceState); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("node_status", flattenNodeStatus(getResp.NodeStatus)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ipmi", flattenIpmiReference(getResp.Ipmi)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rackable_unit_uuid", getResp.RackableUnitUuid); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*getResp.ExtId)
	return nil
}

func flattenNodeStatus(nodeStatus *import1.NodeStatus) string {
	return common.FlattenPtrEnum(nodeStatus)
}

func flattenHostTypeEnum(hostTypeEnum *import1.HostTypeEnum) string {
	return common.FlattenPtrEnum(hostTypeEnum)
}

func flattenHypervisorReference(pr *import1.HypervisorReference) []map[string]interface{} {
	if pr != nil {
		hypervisorRef := make([]map[string]interface{}, 0)

		hyper := make(map[string]interface{})

		hyper["external_address"] = flattenIPAddress(pr.ExternalAddress)
		hyper["user_name"] = pr.UserName
		hyper["full_name"] = pr.FullName
		hyper["type"] = flattenHostHypervisorType(pr.Type)
		hyper["number_of_vms"] = pr.NumberOfVms
		hyper["state"] = flattenHypervisorState(pr.State)
		hyper["acropolis_connection_state"] = flattenAcropolisConnectionState(pr.AcropolisConnectionState)

		hypervisorRef = append(hypervisorRef, hyper)
		return hypervisorRef
	}
	return nil
}

func flattenHypervisorState(hypervisorState *import1.HypervisorState) string {
	return common.FlattenPtrEnum(hypervisorState)
}

func flattenHostHypervisorType(hypervisorType *import1.HypervisorType) string {
	return common.FlattenPtrEnum(hypervisorType)
}

func flattenAcropolisConnectionState(acropolisConnState *import1.AcropolisConnectionState) string {
	return common.FlattenPtrEnum(acropolisConnState)
}

func flattenClusterReference(pr *import1.ClusterReference) []map[string]interface{} {
	if pr != nil {
		clsRef := make([]map[string]interface{}, 0)

		cls := make(map[string]interface{})

		cls["uuid"] = pr.Uuid
		cls["name"] = pr.Name

		clsRef = append(clsRef, cls)
		return clsRef
	}
	return nil
}

func flattenControllerVMReference(pr *import1.ControllerVmReference) []map[string]interface{} {
	if pr != nil {
		cvmRef := make([]map[string]interface{}, 0)

		cvm := make(map[string]interface{})

		cvm["external_address"] = flattenIPAddress(pr.ExternalAddress)
		cvm["backplane_address"] = flattenIPAddress(pr.BackplaneAddress)
		cvm["rdma_backplane_address"] = flattenIPAddressList(pr.RdmaBackplaneAddress)
		cvm["nat_ip"] = flattenIPAddress(pr.NatIp)
		cvm["nat_port"] = pr.NatPort
		cvm["maintenance_mode"] = pr.IsInMaintenanceMode

		cvmRef = append(cvmRef, cvm)
		return cvmRef
	}
	return nil
}

func flattenIPAddressList(pr []config.IPAddress) []map[string]interface{} {
	if len(pr) > 0 {
		ips := make([]map[string]interface{}, len(pr))

		for k, v := range pr {
			ip := make(map[string]interface{})

			ip["ipv4"] = flattenIPv4Address(v.Ipv4)
			ip["ipv6"] = flattenIPv6Address(v.Ipv6)

			ips[k] = ip
		}
		return ips
	}
	return nil
}

func flattenIpmiReference(pr *import1.IpmiReference) []map[string]interface{} {
	if pr != nil {
		ipmiRef := make([]map[string]interface{}, 0)
		ipmi := make(map[string]interface{})

		ipmi["ip"] = flattenIPAddress(pr.Ip)
		ipmi["username"] = pr.Username

		ipmiRef = append(ipmiRef, ipmi)
		return ipmiRef
	}
	return nil
}

func flattenDiskReference(pr []import1.DiskReference) []interface{} {
	if len(pr) > 0 {
		diskRef := make([]interface{}, len(pr))

		for k, v := range pr {
			disk := make(map[string]interface{})

			disk["uuid"] = v.Uuid
			disk["mount_path"] = v.MountPath
			disk["size_in_bytes"] = v.SizeInBytes
			disk["serial_id"] = v.SerialId
			disk["storage_tier"] = flattenStorageTierReference(v.StorageTier)

			diskRef[k] = disk
		}
		return diskRef
	}
	return nil
}

func flattenKeyManagementDeviceToCertStatusInfo(pr []import1.KeyManagementDeviceToCertStatusInfo) []interface{} {
	if len(pr) > 0 {
		keymgmInfo := make([]interface{}, len(pr))

		for k, v := range pr {
			mgm := make(map[string]interface{})

			mgm["key_management_server_name"] = v.KeyManagementServerName
			mgm["is_certificate_present"] = v.IsCertificatePresent

			keymgmInfo[k] = mgm
		}
		return keymgmInfo
	}
	return nil
}

func flattenStorageTierReference(storageTierReference *import1.StorageTierReference) string {
	return common.FlattenPtrEnum(storageTierReference)
}
