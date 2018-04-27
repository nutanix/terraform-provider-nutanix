package nutanix

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixVirtualMachineCreate,
		Read:   resourceNutanixVirtualMachineRead,
		Update: resourceNutanixVirtualMachineUpdate,
		Delete: resourceNutanixVirtualMachineDelete,
		Exists: resourceNutanixVirtualMachineExists,

		Schema: getVMSchema(),
	}
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*NutanixClient).API

	var version string
	if v, ok := d.GetOk("api_version"); ok {
		version = v.(string)
	} else {
		version = Version
	}

	// Prepare request
	request := &v3.VMIntentInput{
		APIVersion: &version,
	}

	spec := &v3.VM{}

	// Read Arguments and set request values
	m, mok := d.GetOk("metadata")
	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")
	azr, azrok := d.GetOk("availability_zone_reference")
	cr, crok := d.GetOk("cluster_reference")
	r, rok := d.GetOk("resources")

	if !mok && !rok && !nok {
		return fmt.Errorf("Please provide the required attributes metadata, name and resources")
	}

	if azrok {
		a := azr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
		}
		if v, ok := a["name"]; ok {
			r.Name = utils.String(v.(string))
		}
		spec.AvailabilityZoneReference = r
	}
	if descok {
		spec.Description = utils.String(desc.(string))
	}
	if crok {
		a := cr.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
		}
		if v, ok := a["name"]; ok {
			r.Name = utils.String(v.(string))
		}
		spec.ClusterReference = r
	}

	spec.Name = utils.String(n.(string))

	metad := m.([]interface{})[0].(map[string]interface{})

	metadata := &v3.VMMetadata{
		Kind: utils.String(metad["kind"].(string)),
	}
	if v, ok := metad["uuid"]; ok && v != "" {
		metadata.UUID = utils.String(v.(string))
	}
	if v, ok := metad["spec_version"]; ok && v != 0 {
		metadata.SpecVersion = utils.Int64(int64(v.(int)))
	}
	if v, ok := metad["spec_hash"]; ok && v != "" {
		metadata.SpecHash = utils.String(v.(string))
	}
	if v, ok := metad["name"]; ok {
		metadata.Name = utils.String(v.(string))
	}
	if v, ok := metad["categories"]; ok {
		p := v.([]interface{})
		if len(p) > 0 {
			c := p[0].(map[string]interface{})
			labels := map[string]string{}

			for k, v := range c {
				labels[k] = v.(string)
			}
			metadata.Categories = labels
		}

	}
	if v, ok := metad["project_reference"]; ok {
		p := v.([]interface{})

		if len(p) > 0 {
			pr := p[0].(map[string]interface{})

			r := &v3.Reference{
				Kind: utils.String(pr["kind"].(string)),
				UUID: utils.String(pr["uuid"].(string)),
			}
			if v1, ok1 := pr["name"]; ok1 {
				r.Name = utils.String(v1.(string))
			}
			metadata.ProjectReference = r
		}
	}
	if v, ok := metad["owner_reference"]; ok {
		p := v.([]interface{})

		if len(p) > 0 {
			pr := p[0].(map[string]interface{})
			r := &v3.Reference{
				Kind: utils.String(pr["kind"].(string)),
				UUID: utils.String(pr["uuid"].(string)),
			}
			if v1, ok1 := pr["name"]; ok1 {
				r.Name = utils.String(v1.(string))
			}
			metadata.OwnerReference = r
		}
	}

	res, err := setVMResources(r)

	if err != nil {
		return err
	}

	spec.Resources = res

	request.Metadata = metadata
	request.Spec = spec

	utils.PrintToJSON(request, "REQUEST VM")

	// Make request to the API
	resp, err := conn.V3.CreateVM(request)
	if err != nil {
		return err
	}

	uuid := *resp.Metadata.UUID

	// Wait for the VM to be available
	status, err := waitForVMProcess(conn, uuid)
	for status != true {
		return err
	}

	// Set terraform state id
	d.SetId(uuid)
	d.Partial(true)
	d.Set("ip_address", "")
	d.Partial(false)

	// Read the ip
	if resp.Spec.Resources.NicList != nil && *resp.Spec.Resources.PowerState == "ON" {
		log.Printf("[DEBUG] Polling for IP\n")
		err := waitForIP(conn, uuid, d)
		if err != nil {
			return err
		}
	}

	return resourceNutanixVirtualMachineRead(d, meta)
}

func resourceNutanixVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*NutanixClient).API

	// Make request to the API
	resp, err := conn.V3.GetVM(d.Id())
	if err != nil {
		return err
	}

	utils.PrintToJSON(resp, "VMRead Response")

	// Set vm values
	// set availability zone reference values
	availabilityZoneReference := make(map[string]interface{})

	if resp.Status.AvailabilityZoneReference != nil {
		availabilityZoneReference["kind"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Kind)
		availabilityZoneReference["name"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Name)
		availabilityZoneReference["uuid"] = utils.StringValue(resp.Status.AvailabilityZoneReference.UUID)
	}

	// set message list values
	if resp.Status.MessageList != nil {
		messages := make([]map[string]interface{}, len(resp.Status.MessageList))
		for k, v := range resp.Status.MessageList {
			message := make(map[string]interface{})

			message["message"] = utils.StringValue(v.Message)
			message["reason"] = utils.StringValue(v.Reason)
			message["details"] = v.Details

			messages[k] = message
		}

		if err := d.Set("message_list", messages); err != nil {
			return err
		}
	}

	// set cluster reference values
	clusterReference := make(map[string]interface{})
	clusterReference["kind"] = utils.StringValue(resp.Status.ClusterReference.Kind)
	clusterReference["name"] = utils.StringValue(resp.Status.ClusterReference.Name)
	clusterReference["uuid"] = utils.StringValue(resp.Status.ClusterReference.UUID)

	// set resources values
	resouces := make(map[string]interface{})

	vnumaConfig := make(map[string]interface{})
	vnumaConfig["num_vnuma_nodes"] = utils.Int64Value(resp.Status.Resources.VnumaConfig.NumVnumaNodes)

	resouces["vnuma_config"] = vnumaConfig

	nics := resp.Status.Resources.NicList
	if nics != nil {

		nicLists := make([]map[string]interface{}, len(nics))
		for k, v := range nics {
			nic := make(map[string]interface{})
			// simple firts
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["floating_ip"] = utils.StringValue(v.FloatingIP)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)

			ipEndpointList := make([]map[string]interface{}, len(v.IPEndpointList))
			for k1, v1 := range v.IPEndpointList {
				ipEndpoint := make(map[string]interface{})
				ipEndpoint["ip"] = utils.StringValue(v1.IP)
				ipEndpoint["type"] = utils.StringValue(v1.Type)
				ipEndpointList[k1] = ipEndpoint
			}
			nic["ip_endpoint_list"] = ipEndpointList

			netFnChainRef := make(map[string]interface{})
			netFnChainRef["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
			netFnChainRef["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
			netFnChainRef["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)

			nic["network_function_chain_reference"] = netFnChainRef

			subtnetRef := make(map[string]interface{})
			subtnetRef["kind"] = utils.StringValue(v.SubnetReference.Kind)
			subtnetRef["name"] = utils.StringValue(v.SubnetReference.Name)
			subtnetRef["uuid"] = utils.StringValue(v.SubnetReference.UUID)

			nic["subnet_reference"] = subtnetRef

			nicLists[k] = nic
		}
		resouces["nic_list"] = nicLists
	}

	hostRef := make(map[string]interface{})
	hostRef["kind"] = utils.StringValue(resp.Status.Resources.HostReference.Kind)
	hostRef["name"] = utils.StringValue(resp.Status.Resources.HostReference.Name)
	hostRef["uuid"] = utils.StringValue(resp.Status.Resources.HostReference.UUID)

	resouces["host_reference"] = hostRef

	guestTools := make(map[string]interface{})

	if resp.Status.Resources.GuestTools != nil {
		tools := resp.Status.Resources.GuestTools.NutanixGuestTools
		nutanixGuestTools := make(map[string]interface{})
		nutanixGuestTools["available_version"] = utils.StringValue(tools.AvailableVersion)
		nutanixGuestTools["iso_mount_state"] = utils.StringValue(tools.IsoMountState)
		nutanixGuestTools["state"] = utils.StringValue(tools.State)
		nutanixGuestTools["version"] = utils.StringValue(tools.Version)
		nutanixGuestTools["guest_os_version"] = utils.StringValue(tools.GuestOsVersion)

		capList := make([]string, len(tools.EnabledCapabilityList))
		for k, v := range tools.EnabledCapabilityList {
			capList[k] = *v
		}
		nutanixGuestTools["enabled_capability_list"] = capList
		nutanixGuestTools["vss_snapshot_capable"] = utils.BoolValue(tools.VSSSnapshotCapable)
		nutanixGuestTools["is_reachable"] = utils.BoolValue(tools.IsReachable)
		nutanixGuestTools["vm_mobility_drivers_installed"] = utils.BoolValue(tools.VMMobilityDriversInstalled)

		guestTools["nutanix_guest_tools"] = nutanixGuestTools

		resouces["guest_tools"] = guestTools
	}

	gpuList := make([]map[string]interface{}, len(resp.Status.Resources.GpuList))

	if resp.Status.Resources.GpuList != nil {
		for k, v := range resp.Status.Resources.GpuList {
			gpu := make(map[string]interface{})
			gpu["frame_buffer_size_mib"] = utils.Int64Value(v.FrameBufferSizeMib)
			gpu["vendor"] = utils.StringValue(v.Vendor)
			gpu["uuid"] = utils.StringValue(v.UUID)
			gpu["name"] = utils.StringValue(v.Name)
			gpu["pci_address"] = utils.StringValue(v.PCIAddress)
			gpu["fraction"] = utils.Int64Value(v.Fraction)
			gpu["mode"] = utils.StringValue(v.Mode)
			gpu["num_virtual_display_heads"] = utils.Int64Value(v.NumVirtualDisplayHeads)
			gpu["guest_driver_version"] = utils.StringValue(v.GuestDriverVersion)
			gpu["device_id"] = utils.Int64Value(v.DeviceID)

			gpuList[k] = gpu
		}
		resouces["gpu_list"] = gpuList
	}

	if resp.Status.Resources.ParentReference != nil {
		parentRef := make(map[string]interface{})
		parentRef["kind"] = utils.StringValue(resp.Status.Resources.ParentReference.Kind)
		parentRef["name"] = utils.StringValue(resp.Status.Resources.ParentReference.Name)
		parentRef["uuid"] = utils.StringValue(resp.Status.Resources.ParentReference.UUID)

		resouces["parent_reference"] = parentRef
	}

	if resp.Status.Resources.BootConfig != nil {
		bootConfig := make(map[string]interface{})
		boots := make([]string, len(resp.Status.Resources.BootConfig.BootDeviceOrderList))
		for k, v := range resp.Status.Resources.BootConfig.BootDeviceOrderList {
			boots[k] = utils.StringValue(v)
		}
		bootDevice := make(map[string]interface{})
		diskAddress := make(map[string]interface{})
		diskAddress["device_index"] = utils.Int64Value(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)
		diskAddress["adapter_type"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
		bootDevice["disk_address"] = diskAddress
		bootDevice["mac_address"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.MacAddress)

		bootConfig["boot_device"] = bootDevice
		bootConfig["boot_device_order_list"] = boots

		resouces["boot_config"] = bootConfig
	}

	if resp.Status.Resources.GuestCustomization != nil {
		guestCustom := make(map[string]interface{})
		cloudInit := make(map[string]interface{})
		cloudInit["meta_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.MetaData)
		cloudInit["user_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.UserData)
		cloudInit["custom_key_values"] = resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues

		guestCustom["cloud_init"] = cloudInit
		guestCustom["is_overridable"] = utils.BoolValue(resp.Status.Resources.GuestCustomization.IsOverridable)

		sysprep := make(map[string]interface{})
		sysprep["install_type"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.InstallType)
		sysprep["unattend_xml"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
		sysprep["custom_key_values"] = resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues

		guestCustom["sysprep"] = sysprep

		resouces["guest_customization"] = guestCustom
	}

	powerStateMechanism := make([]map[string]interface{}, 1)
	psm := make(map[string]interface{})

	psm["mechanism"] = utils.StringValue(resp.Status.Resources.PowerStateMechanism.Mechanism)

	guestTransition := make(map[string]interface{})
	guestTransition["should_fail_on_script_failure"] = utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure)
	guestTransition["enable_script_exec"] = utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec)

	psm["guest_transition_config"] = guestTransition
	powerStateMechanism[0] = psm
	resouces["power_state_mechanism"] = powerStateMechanism

	diskList := make([]map[string]interface{}, len(resp.Status.Resources.DiskList))

	if resp.Status.Resources.DiskList != nil {
		for k, v := range resp.Status.Resources.DiskList {
			disk := make(map[string]interface{})
			disk["uuid"] = *v.UUID
			disk["disk_size_bytes"] = *v.DiskSizeBytes
			disk["disk_size_mib"] = *v.DiskSizeMib

			dsourceRef := make(map[string]interface{})
			dsourceRef["kind"] = utils.StringValue(v.DataSourceReference.Kind)
			dsourceRef["name"] = utils.StringValue(v.DataSourceReference.Name)
			dsourceRef["uuid"] = utils.StringValue(v.DataSourceReference.UUID)

			disk["data_source_reference"] = dsourceRef

			volumeRef := make(map[string]interface{})
			volumeRef["kind"] = utils.StringValue(v.VolumeGroupReference.Kind)
			volumeRef["name"] = utils.StringValue(v.VolumeGroupReference.Name)
			volumeRef["uuid"] = utils.StringValue(v.VolumeGroupReference.UUID)

			disk["volume_group_reference"] = volumeRef

			deviceProps := make(map[string]interface{})
			deviceProps["device_type"] = utils.StringValue(v.DeviceProperties.DeviceType)

			diskAddress := make(map[string]interface{})
			diskAddress["device_index"] = utils.Int64Value(v.DeviceProperties.DiskAddress.DeviceIndex)
			diskAddress["adapter_type"] = utils.StringValue(v.DeviceProperties.DiskAddress.AdapterType)

			deviceProps["disk_address"] = diskAddress

			disk["device_properties"] = deviceProps

			diskList[k] = disk
		}

		resouces["disk_list"] = diskList
	}

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = utils.TimeValue(resp.Metadata.LastUpdateTime)
	metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
	metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
	metadata["creation_time"] = utils.TimeValue(resp.Metadata.CreationTime)
	metadata["spec_version"] = utils.Int64Value(resp.Metadata.SpecVersion)
	metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
	metadata["categories"] = resp.Metadata.Categories
	metadata["name"] = utils.StringValue(resp.Metadata.Name)

	pr := make(map[string]interface{})
	pr["kind"] = utils.StringValue(resp.Metadata.ProjectReference.Kind)
	pr["name"] = utils.StringValue(resp.Metadata.ProjectReference.Name)
	pr["uuid"] = utils.StringValue(resp.Metadata.ProjectReference.UUID)

	or := make(map[string]interface{})
	or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
	or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
	or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)

	metadata["project_reference"] = pr
	metadata["owner_reference"] = or

	// Simple first
	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}
	if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
		return err
	}
	if err := d.Set("state", utils.StringValue(resp.Status.State)); err != nil {
		return err
	}
	if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_reference", availabilityZoneReference); err != nil {
		return err
	}

	if err := d.Set("cluster_reference", clusterReference); err != nil {
		return err
	}
	if err := d.Set("resources", resouces); err != nil {
		return err
	}
	if err := d.Set("metadata", metadata); err != nil {
		return err
	}
	d.SetId(resource.UniqueId())

	return nil
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API

	// get state
	uuid := d.Id()
	name := d.Get("name").(string)
	spec := d.Get("resources").(map[string]interface{})

	res, err := setVMResources(spec)

	if err != nil {
		return err
	}

	vmSpec := res

	log.Printf("[DEBUG] Updating Virtual Machine: %s, %s", name, uuid)

	d.Partial(true)

	if d.HasChange("name") || d.HasChange("resources") || d.HasChange("metadata") {
		request := &v3.VMIntentInput{}
		request.Spec.Resources = vmSpec

		_, err := conn.V3.UpdateVM(uuid, request)
		if err != nil {
			return err
		}
		d.SetPartial("resources")
		d.SetPartial("metadata")
		d.Set("ip_address", "")
	}

	d.Partial(false)

	status, err := waitForVMProcess(conn, uuid)
	for status != true {
		return err
	}

	return resourceNutanixVirtualMachineRead(d, meta)
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API
	uuid := d.Id()

	if err := conn.V3.DeleteVM(uuid); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceNutanixVirtualMachineExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*NutanixClient).API

	getEntitiesRequest := &v3.VMListMetadata{}

	resp, err := conn.V3.ListVM(getEntitiesRequest)
	if err != nil {
		return false, err
	}

	for i := range resp.Entities {
		if *resp.Entities[i].Metadata.UUID == d.Id() {
			return true, nil
		}
	}
	return false, nil
}

func setVMResources(m interface{}) (*v3.VMResources, error) {

	vm := &v3.VMResources{}

	resources := m.(map[string]interface{})

	if v, ok := resources["vnuma_config"]; ok {
		vc := v.([]interface{})
		if len(vc) > 0 {
			vnuma := vc[0].(map[string]interface{})
			vm.VMVnumaConfig.NumVnumaNodes = utils.Int64(vnuma["num_vnuma_nodes"].(int64))
		}
	}
	if v, ok := resources["nic_list"]; ok {
		n := v.([]interface{})
		if len(n) > 0 {
			var nics []v3.VMNic

			for _, nc := range n {
				val := nc.(map[string]interface{})
				nic := v3.VMNic{}

				if value, ok := val["nic_type"]; ok {
					nic.NicType = utils.String(value.(string))
				}
				if value, ok := val["uuid"]; ok {
					nic.UUID = utils.String(value.(string))
				}
				if value, ok := val["network_function_nic_type"]; ok {
					nic.NetworkFunctionNicType = utils.String(value.(string))
				}
				if value, ok := val["mac_address"]; ok {
					nic.MacAddress = utils.String(value.(string))
				}
				if value, ok := val["model"]; ok {
					nic.Model = utils.String(value.(string))
				}
				if value, ok := val["ip_endpoint_list"]; ok {
					ipl := value.([]interface{})
					if len(ipl) > 0 {
						var ip []*v3.IPAddress
						for _, i := range ipl {
							v := i.(map[string]interface{})
							ip = append(ip, &v3.IPAddress{IP: utils.String(v["ip"].(string)), Type: utils.String(v["type"].(string))})
						}
						nic.IPEndpointList = ip
					}
				}
				if value, ok := val["network_function_chain_reference"]; ok {
					nfcr := value.([]interface{})
					if len(nfcr) > 0 {
						v := nfcr[0].(map[string]string)
						nic.NetworkFunctionChainReference.Kind = utils.String(v["kind"])
						nic.NetworkFunctionChainReference.UUID = utils.String(v["uuid"])
						if j, ok1 := v["name"]; ok1 {
							nic.NetworkFunctionChainReference.Name = utils.String(j)
						}
					}
				}
				if value, ok := val["subnet_reference"]; ok {
					sr := value.([]interface{})
					if len(sr) > 0 {
						v := sr[0].(map[string]string)
						nic.SubnetReference.Kind = utils.String(v["kind"])
						nic.SubnetReference.UUID = utils.String(v["uuid"])
						if j, ok1 := v["name"]; ok1 {
							nic.SubnetReference.Name = utils.String(j)
						}
					}
				}

				nics = append(nics, nic)
			}

			vm.NicList = nics
		}
	}
	if v, ok := resources["guest_os_id"]; ok {
		vm.GuestOsID = utils.String(v.(string))
	}
	if v, ok := resources["power_state"]; ok {
		vm.PowerState = utils.String(v.(string))
	}
	if v, ok := resources["guest_tools"]; ok {
		gt := v.([]interface{})
		if len(gt) > 0 {
			ngt := gt[0].(map[string]interface{})
			if k, ok1 := ngt["nutanix_guest_tools"]; ok1 {
				ngtsi := k.([]interface{})
				if len(ngtsi) > 0 {
					ngts := ngtsi[0].(map[string]interface{})
					if val, ok2 := ngts["iso_mount_state"]; ok2 {
						vm.GuestTools.NutanixGuestTools.IsoMountState = utils.String(val.(string))
					}
					if val, ok2 := ngts["state"]; ok2 {
						vm.GuestTools.NutanixGuestTools.State = utils.String(val.(string))
					}
					if val, ok2 := ngts["enabled_capability_list"]; ok2 {
						var l []*string
						for _, list := range val.([]interface{}) {
							l = append(l, utils.String(list.(string)))
						}
						vm.GuestTools.NutanixGuestTools.EnabledCapabilityList = l
					}
				}
			}
		}
	}
	if v, ok := resources["num_vcpus_per_socket"]; ok {
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			return nil, err
		}
		vm.NumVcpusPerSocket = utils.Int64(int64(i))
	}
	if v, ok := resources["num_sockets"]; ok {
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			return nil, err
		}
		vm.NumSockets = utils.Int64(int64(i))
	}
	if v, ok := resources["gpu_list"]; ok {
		var gpl []*v3.VMGpu
		for _, va := range v.([]interface{}) {
			val := va.(map[string]interface{})
			gpu := &v3.VMGpu{}
			if value, ok1 := val["vendor"]; ok1 {
				gpu.Vendor = utils.String(value.(string))
			}
			if value, ok1 := val["device_id"]; ok1 {
				gpu.DeviceID = utils.Int64(int64(value.(int)))
			}
			if value, ok1 := val["mode"]; ok1 {
				gpu.Mode = utils.String(value.(string))
			}
			gpl = append(gpl, gpu)
		}
		vm.GpuList = gpl
	}
	if v, ok := resources["parent_reference"]; ok {
		pr := v.([]interface{})
		if len(pr) > 0 {
			val := pr[0].(map[string]string)
			vm.ParentReference.Kind = utils.String(val["kind"])
			vm.ParentReference.UUID = utils.String(val["uuid"])
			if j, ok1 := val["name"]; ok1 {
				vm.ParentReference.Name = utils.String(j)
			}
		}
	}
	if v, ok := resources["memory_size_mib"]; ok {
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			return nil, err
		}
		vm.MemorySizeMib = utils.Int64(int64(i))
	}
	if v, ok := resources["boot_config"]; ok {
		btc := v.([]interface{})
		if len(btc) > 0 {
			val := btc[0].(map[string]interface{})
			if value1, ok1 := val["boot_device_order_list"]; ok1 {
				var b []*string
				for _, boot := range value1.([]interface{}) {
					b = append(b, utils.String(boot.(string)))
				}
				vm.BootConfig.BootDeviceOrderList = b
			}
			if value1, ok1 := val["boot_device"]; ok1 {
				btd := value1.([]interface{})
				if len(btd) > 0 {
					bdi := btd[0].(map[string]interface{})
					bd := &v3.VMBootDevice{}
					if value2, ok2 := bdi["disk_address"]; ok2 {
						dka := value2.([]interface{})
						if len(dka) > 0 {
							dai := dka[0].(map[string]interface{})
							da := &v3.DiskAddress{}
							if value3, ok3 := dai["device_index"]; ok3 {
								da.DeviceIndex = utils.Int64(int64(value3.(int)))
							}
							if value3, ok3 := dai["adapter_type"]; ok3 {
								da.AdapterType = utils.String(value3.(string))
							}
							bd.DiskAddress = da
						}
					}
					if value2, ok2 := bdi["mac_address"]; ok2 {
						bd.MacAddress = utils.String(value2.(string))
					}
					vm.BootConfig.BootDevice = bd
				}
			}
		}
	}
	if v, ok := resources["hardware_clock_timezone"]; ok {
		vm.HardwareClockTimezone = utils.String(v.(string))
	}
	if v, ok := resources["guest_customization"]; ok {
		gst := v.([]interface{})
		if len(gst) > 0 {
			gci := gst[0].(map[string]interface{})
			gc := &v3.GuestCustomization{}

			if v1, ok1 := gci["cloud_init"]; ok1 {
				cld := v1.([]interface{})
				if len(cld) > 0 {
					cii := cld[0].(map[string]interface{})
					if v2, ok2 := cii["meta_data"]; ok2 {
						gc.CloudInit.MetaData = utils.String(v2.(string))
					}
					if v2, ok2 := cii["user_data"]; ok2 {
						gc.CloudInit.UserData = utils.String(v2.(string))
					}
					if v2, ok2 := cii["custom_key_values"]; ok2 {
						gc.CloudInit.CustomKeyValues = v2.(map[string]string)
					}
				}
			}
			if v1, ok1 := gci["sysprep"]; ok1 {
				sys := v1.([]interface{})
				if len(sys) > 0 {
					spi := sys[0].(map[string]interface{})
					if v2, ok2 := spi["install_type"]; ok2 {
						gc.Sysprep["install_type"] = v2.(string)
					}
					if v2, ok2 := spi["unattend_xml"]; ok2 {
						gc.Sysprep["unattend_xml"] = v2.(string)
					}
					if v2, ok2 := spi["custom_key_values"]; ok2 {
						gc.Sysprep["custom_key_values"] = v2.(map[string]string)
					}
				}
			}
			if v1, ok1 := gci["is_overridable"]; ok1 {
				gc.IsOverridable = utils.Bool(v1.(bool))
			}

			vm.GuestCustomization = gc
		}
	}
	if v, ok := resources["vga_console_enabled"]; ok {
		vm.VgaConsoleEnabled = utils.Bool(v.(bool))
	}
	if v, ok := resources["power_state_mechanism"]; ok {
		ps := v.([]interface{})
		if len(ps) > 0 {
			p := ps[0].(map[string]interface{})
			psm := &v3.VMPowerStateMechanism{}

			if v1, ok1 := p["mechanism"]; ok1 {
				psm.Mechanism = utils.String(v1.(string))
			}
			if v1, ok1 := p["guest_transition_config"]; ok1 {
				gst := v1.([]interface{})
				if len(gst) > 0 {
					g := gst[0].(map[string]interface{})
					gtc := &v3.VMGuestPowerStateTransitionConfig{}
					if v2, ok2 := g["should_fail_on_script_failure"]; ok2 {
						gtc.ShouldFailOnScriptFailure = utils.Bool(v2.(bool))
					}
					if v2, ok2 := g["enable_script_exec"]; ok2 {
						gtc.EnableScriptExec = utils.Bool(v2.(bool))
					}
					psm.GuestTransitionConfig = gtc
				}
			}

			vm.PowerStateMechanism = psm
		}
	}
	if v, ok := resources["disk_list"]; ok {
		dsk := v.([]interface{})
		if len(dsk) > 0 {
			dls := make([]*v3.VMDisk, len(dsk))

			for k, val := range dsk {
				v := val.(map[string]interface{})
				dl := &v3.VMDisk{}
				if v1, ok1 := v["uuid"]; ok1 {
					dl.UUID = utils.String(v1.(string))
				}
				if v1, ok1 := v["disk_size_bytes"]; ok1 {
					dl.DiskSizeBytes = utils.Int64(int64(v1.(int)))
				}
				if v1, ok1 := v["device_properties"]; ok1 {
					dvp := v1.([]interface{})
					if len(dvp) > 0 {
						d := dvp[0].(map[string]interface{})
						dp := &v3.VMDiskDeviceProperties{}
						if v, ok := d["device_type"]; ok {
							dp.DeviceType = utils.String(v.(string))
						}
						if v, ok := d["disk_address"]; ok {
							da := v.([]interface{})[0].(map[string]interface{})
							dp.DiskAddress = v3.DiskAddress{
								DeviceIndex: utils.Int64(int64(da["device_index"].(int))),
								AdapterType: utils.String(da["adapter_type"].(string)),
							}
						}
						dl.DeviceProperties = dp
					}
				}
				if v1, ok := v["data_source_reference"]; ok {
					dsref := v1.([]interface{})
					if len(dsref) > 0 {
						dsri := dsref[0].(map[string]interface{})
						dsr := &v3.Reference{
							Kind: utils.String(dsri["kind"].(string)),
							UUID: utils.String(dsri["uuid"].(string)),
						}
						if v2, ok2 := dsri["name"]; ok2 {
							dsr.Name = utils.String(v2.(string))
						}
						dl.DataSourceReference = dsr
					}
				}
				if v1, ok := v["volume_group_reference"]; ok {
					volgr := v1.([]interface{})
					if len(volgr) > 0 {
						dsri := volgr[0].(map[string]interface{})
						dsr := &v3.Reference{
							Kind: utils.String(dsri["kind"].(string)),
							UUID: utils.String(dsri["uuid"].(string)),
						}
						if v2, ok2 := dsri["name"]; ok2 {
							dsr.Name = utils.String(v2.(string))
						}
						dl.VolumeGroupReference = dsr
					}
				}
				if v1, ok := v["disk_size_mib"]; ok {
					dl.DiskSizeMib = utils.Int64(int64(v1.(int)))
				}
				dls[k] = dl
			}

			vm.DiskList = dls
		}
	}

	return vm, nil
}

func waitForVMProcess(conn *v3.Client, uuid string) (bool, error) {
	for {
		resp, err := conn.V3.GetVM(uuid)
		if err != nil {
			return false, err
		}

		if utils.StringValue(resp.Status.State) == "COMPLETE" {
			return true, nil
		} else if utils.StringValue(resp.Status.State) == "ERROR" {
			return false, fmt.Errorf("Error while waiting for resource to be up")
		}
		time.Sleep(3000 * time.Millisecond)
	}
	// return false, nil
}

func waitForIP(conn *v3.Client, uuid string, d *schema.ResourceData) error {
	for {
		resp, err := conn.V3.GetVM(uuid)
		if err != nil {
			return err
		}

		if len(resp.Status.Resources.NicList) != 0 {
			for i := range resp.Status.Resources.NicList {
				if len(resp.Status.Resources.NicList[i].IPEndpointList) != 0 {
					if ip := resp.Status.Resources.NicList[i].IPEndpointList[0].IP; ip != nil {
						// TODO set ip address
						d.Set("ip_address", *ip)
						return nil
					}
				}
			}
		}
		time.Sleep(3000 * time.Millisecond)
	}
	// return nil
}

func getVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"project_reference": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"creation_time": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"spec_version": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"spec_hash": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"owner_reference": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"categories": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeMap},
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"availability_zone_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"message_list": &schema.Schema{
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"message": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"reason": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"details": &schema.Schema{
						Type:     schema.TypeMap,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"name": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"resources": &schema.Schema{
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vnuma_config": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"num_vnuma_nodes": &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"nic_list": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"nic_type": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"ip_endpoint_list": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"ip": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"type": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
										},
									},
								},
								"network_function_chain_reference": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"kind": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
											"name": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
								"floating_ip": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"network_function_nic_type": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"mac_address": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"subnet_reference": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"kind": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
											"name": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
								"model": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"host_reference": &schema.Schema{
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
					"guest_os_id": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"power_state": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"guest_tools": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"nutanix_guest_tools": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"available_version": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"iso_mount_state": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"state": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"version": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"guest_os_version": &schema.Schema{
												Type:     schema.TypeString,
												Computed: true,
											},
											"enabled_capability_list": &schema.Schema{
												Type:     schema.TypeList,
												Optional: true,
												Computed: true,
												Elem:     &schema.Schema{Type: schema.TypeString},
											},
											"vss_snapshot_capable": &schema.Schema{
												Type:     schema.TypeBool,
												Computed: true,
											},
											"is_reachable": &schema.Schema{
												Type:     schema.TypeBool,
												Computed: true,
											},
											"vm_mobility_drivers_installed": &schema.Schema{
												Type:     schema.TypeBool,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"hypervisor_type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"num_vcpus_per_socket": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"num_sockets": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"gpu_list": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"frame_buffer_size_mib": &schema.Schema{
									Type:     schema.TypeInt,
									Computed: true,
								},
								"vendor": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"pci_address": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"fraction": &schema.Schema{
									Type:     schema.TypeInt,
									Computed: true,
								},
								"mode": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"num_virtual_display_heads": &schema.Schema{
									Type:     schema.TypeInt,
									Computed: true,
								},
								"guest_driver_version": &schema.Schema{
									Type:     schema.TypeString,
									Computed: true,
								},
								"device_id": &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"parent_reference": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"memory_size_mib": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"boot_config": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"boot_device_order_list": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem:     &schema.Schema{Type: schema.TypeString},
								},
								"boot_device": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"disk_address": &schema.Schema{
												Type:     schema.TypeList,
												Optional: true,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"device_index": &schema.Schema{
															Type:     schema.TypeInt,
															Optional: true,
															Computed: true,
														},
														"adapter_type": &schema.Schema{
															Type:     schema.TypeString,
															Optional: true,
															Computed: true,
														},
													},
												},
											},
											"mac_address": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"hardware_clock_timezone": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"guest_customization": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cloud_init": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"meta_data": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"user_data": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"custom_key_values": &schema.Schema{
												Type:     schema.TypeMap,
												Optional: true,
												Computed: true,
											},
										},
									},
								},
								"is_overridable": &schema.Schema{
									Type:     schema.TypeBool,
									Optional: true,
									Computed: true,
								},
								"sysprep": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"install_type": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"unattend_xml": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"custom_key_values": &schema.Schema{
												Type:     schema.TypeMap,
												Optional: true,
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
					"power_state_mechanism": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"guest_transition_config": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"should_fail_on_script_failure": &schema.Schema{
												Type:     schema.TypeBool,
												Optional: true,
												Computed: true,
											},
											"enable_script_exec": &schema.Schema{
												Type:     schema.TypeBool,
												Optional: true,
												Computed: true,
											},
										},
									},
								},
								"mechanism": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"vga_console_enabled": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
						Computed: true,
					},
					"disk_list": &schema.Schema{
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"uuid": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"disk_size_bytes": &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"device_properties": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_type": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"disk_address": &schema.Schema{
												Type:     schema.TypeList,
												Optional: true,
												Computed: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														"device_index": &schema.Schema{
															Type:     schema.TypeInt,
															Required: true,
														},
														"adapter_type": &schema.Schema{
															Type:     schema.TypeString,
															Required: true,
														},
													},
												},
											},
										},
									},
								},
								"data_source_reference": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"kind": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
											"name": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
								"disk_size_mib": &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"volume_group_reference": &schema.Schema{
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"kind": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
											"name": &schema.Schema{
												Type:     schema.TypeString,
												Optional: true,
												Computed: true,
											},
											"uuid": &schema.Schema{
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
