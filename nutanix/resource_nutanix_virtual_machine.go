package nutanix

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNutanixVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixVirtualMachineCreate,
		Read:   resourceNutanixVirtualMachineRead,
		Update: resourceNutanixVirtualMachineUpdate,
		Delete: resourceNutanixVirtualMachineDelete,
		Exists: resourceNutanixVirtualMachineExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: getVMSchema(),
	}
}

func resourceNutanixVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*NutanixClient).API

	// Prepare request
	request := &v3.VMIntentInput{}
	spec := &v3.VM{}
	metadata := &v3.VMMetadata{}
	res := &v3.VMResources{}

	// Read Arguments and set request values
	n, nok := d.GetOk("name")
	desc, descok := d.GetOk("description")
	azr, azrok := d.GetOk("availability_zone_reference")
	cr, crok := d.GetOk("cluster_reference")

	if v, ok := d.GetOk("api_version"); ok {
		request.APIVersion = utils.String(v.(string))
	}
	if !nok {
		return fmt.Errorf("Please provide the required name attribute")
	}
	if err := getMetadaAttributes(d, metadata); err != nil {
		return err
	}
	if descok {
		spec.Description = utils.String(desc.(string))
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

	if err := getVMResources(d, res); err != nil {
		return err
	}

	spec.Name = utils.String(n.(string))
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

	// Set terraform state id
	d.SetId(uuid)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    vmStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for vm (%s) to create: %s", d.Id(), err)
	}

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

	log.Printf("Reading VM values %s", d.Id())
	fmt.Printf("Reading VM values %s", d.Id())

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = resp.Metadata.LastUpdateTime.String()
	metadata["kind"] = utils.StringValue(resp.Metadata.Kind)
	metadata["uuid"] = utils.StringValue(resp.Metadata.UUID)
	metadata["creation_time"] = resp.Metadata.CreationTime.String()
	metadata["spec_version"] = strconv.Itoa(int(utils.Int64Value(resp.Metadata.SpecVersion)))
	metadata["spec_hash"] = utils.StringValue(resp.Metadata.SpecHash)
	metadata["name"] = utils.StringValue(resp.Metadata.Name)
	if err := d.Set("metadata", metadata); err != nil {
		return err
	}
	if err := d.Set("categories", resp.Metadata.Categories); err != nil {
		return err
	}
	pr := make(map[string]interface{})
	if resp.Metadata.ProjectReference != nil {
		pr["kind"] = utils.StringValue(resp.Metadata.ProjectReference.Kind)
		pr["name"] = utils.StringValue(resp.Metadata.ProjectReference.Name)
		pr["uuid"] = utils.StringValue(resp.Metadata.ProjectReference.UUID)

	}
	if err := d.Set("project_reference", pr); err != nil {
		return err
	}
	or := make(map[string]interface{})
	if resp.Metadata.OwnerReference != nil {
		or["kind"] = utils.StringValue(resp.Metadata.OwnerReference.Kind)
		or["name"] = utils.StringValue(resp.Metadata.OwnerReference.Name)
		or["uuid"] = utils.StringValue(resp.Metadata.OwnerReference.UUID)

	}
	if err := d.Set("owner_reference", or); err != nil {
		return err
	}
	if err := d.Set("api_version", utils.StringValue(resp.APIVersion)); err != nil {
		return err
	}
	if err := d.Set("name", utils.StringValue(resp.Status.Name)); err != nil {
		return err
	}
	if err := d.Set("description", utils.StringValue(resp.Status.Description)); err != nil {
		return err
	}
	// set availability zone reference values
	availabilityZoneReference := make(map[string]interface{})
	if resp.Status.AvailabilityZoneReference != nil {
		availabilityZoneReference["kind"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Kind)
		availabilityZoneReference["name"] = utils.StringValue(resp.Status.AvailabilityZoneReference.Name)
		availabilityZoneReference["uuid"] = utils.StringValue(resp.Status.AvailabilityZoneReference.UUID)
	}
	if err := d.Set("availability_zone_reference", availabilityZoneReference); err != nil {
		return err
	}
	// set cluster reference values
	clusterReference := make(map[string]interface{})
	if resp.Status.ClusterReference != nil {
		clusterReference["kind"] = utils.StringValue(resp.Status.ClusterReference.Kind)
		clusterReference["name"] = utils.StringValue(resp.Status.ClusterReference.Name)
		clusterReference["uuid"] = utils.StringValue(resp.Status.ClusterReference.UUID)
	}

	utils.PrintToJSON(clusterReference, "CLUSTER READ =>")

	if err := d.Set("cluster_reference", clusterReference); err != nil {
		return err
	}
	// set state value
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}
	// set vnuma_config value
	if err := d.Set("num_vnuma_nodes", utils.Int64Value(resp.Status.Resources.VnumaConfig.NumVnumaNodes)); err != nil {
		return err
	}
	// set nic list value
	nics := resp.Status.Resources.NicList
	nicLists := make([]map[string]interface{}, 0)
	if nics != nil {
		nicLists = make([]map[string]interface{}, len(nics))
		for k, v := range nics {
			nic := make(map[string]interface{})
			// simple first
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["floating_ip"] = utils.StringValue(v.FloatingIP)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)

			// set ip lists value
			ipEndpointList := make([]map[string]interface{}, len(v.IPEndpointList))
			for k1, v1 := range v.IPEndpointList {
				ipEndpoint := make(map[string]interface{})
				ipEndpoint["ip"] = utils.StringValue(v1.IP)
				ipEndpoint["type"] = utils.StringValue(v1.Type)
				ipEndpointList[k1] = ipEndpoint
			}
			nic["ip_endpoint_list"] = ipEndpointList

			// set network_function_chain_reference value
			netFnChainRef := make(map[string]interface{})
			if v.NetworkFunctionChainReference != nil {
				netFnChainRef["kind"] = utils.StringValue(v.NetworkFunctionChainReference.Kind)
				netFnChainRef["name"] = utils.StringValue(v.NetworkFunctionChainReference.Name)
				netFnChainRef["uuid"] = utils.StringValue(v.NetworkFunctionChainReference.UUID)
			}
			nic["network_function_chain_reference"] = netFnChainRef

			// set subnet_reference value
			subtnetRef := make(map[string]interface{})
			if v.SubnetReference != nil {
				subtnetRef["kind"] = utils.StringValue(v.SubnetReference.Kind)
				subtnetRef["name"] = utils.StringValue(v.SubnetReference.Name)
				subtnetRef["uuid"] = utils.StringValue(v.SubnetReference.UUID)
			}
			nic["subnet_reference"] = subtnetRef

			nicLists[k] = nic
		}

	}
	if err := d.Set("nic_list", nicLists); err != nil {
		return err
	}
	// set host_reference value
	hostRef := make(map[string]interface{})
	if resp.Status.Resources.HostReference != nil {
		hostRef["kind"] = utils.StringValue(resp.Status.Resources.HostReference.Kind)
		hostRef["name"] = utils.StringValue(resp.Status.Resources.HostReference.Name)
		hostRef["uuid"] = utils.StringValue(resp.Status.Resources.HostReference.UUID)

	}
	if err := d.Set("host_reference", hostRef); err != nil {
		return err
	}
	// set guest_os_id value
	if err := d.Set("guest_os_id", resp.Status.Resources.GuestOsID); err != nil {
		return err
	}
	// set power_state value
	if err := d.Set("power_state", resp.Status.Resources.PowerState); err != nil {
		return err
	}

	nutanixGuestTools := make(map[string]interface{})
	if resp.Status.Resources.GuestTools != nil {
		tools := resp.Status.Resources.GuestTools.NutanixGuestTools
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
	}
	if err := d.Set("nutanix_guest_tools", nutanixGuestTools); err != nil {
		return err
	}
	// set num_vcpus_per_socket value
	if err := d.Set("num_vcpus_per_socket", resp.Status.Resources.NumVcpusPerSocket); err != nil {
		return err
	}
	// set num_sockets value
	if err := d.Set("num_sockets", resp.Status.Resources.NumSockets); err != nil {
		return err
	}
	gpuList := make([]map[string]interface{}, 0)
	if resp.Status.Resources.GpuList != nil {
		gpuList = make([]map[string]interface{}, len(resp.Status.Resources.GpuList))
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
	}
	if err := d.Set("gpu_list", gpuList); err != nil {
		return err
	}

	parentRef := make(map[string]interface{})
	if resp.Status.Resources.ParentReference != nil {
		parentRef["kind"] = utils.StringValue(resp.Status.Resources.ParentReference.Kind)
		parentRef["name"] = utils.StringValue(resp.Status.Resources.ParentReference.Name)
		parentRef["uuid"] = utils.StringValue(resp.Status.Resources.ParentReference.UUID)
	}
	if err := d.Set("parent_reference", parentRef); err != nil {
		return err
	}

	// set memory_size_mib value
	if err := d.Set("memory_size_mib", utils.Int64Value(resp.Status.Resources.MemorySizeMib)); err != nil {
		return err
	}

	boots := make([]string, 0)
	diskAddress := make(map[string]interface{})
	mac := ""

	if resp.Status.Resources.BootConfig != nil {
		boots = make([]string, len(resp.Status.Resources.BootConfig.BootDeviceOrderList))
		for k, v := range resp.Status.Resources.BootConfig.BootDeviceOrderList {
			boots[k] = utils.StringValue(v)
		}

		if resp.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
			i := strconv.Itoa(int(utils.Int64Value(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)))

			diskAddress["device_index"] = i
			diskAddress["adapter_type"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
		}

		mac = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.MacAddress)
	}
	if err := d.Set("boot_device_order_list", boots); err != nil {
		return err
	}
	if err := d.Set("boot_device_disk_address", diskAddress); err != nil {
		return err
	}
	if err := d.Set("boot_device_mac_address", mac); err != nil {
		return err
	}
	// set hardware_clock_timezone value
	if err := d.Set("hardware_clock_timezone", utils.StringValue(resp.Status.Resources.HardwareClockTimezone)); err != nil {
		return err
	}
	cloudInit := make(map[string]interface{})
	sysprep := make(map[string]interface{})
	isOverride := false
	if resp.Status.Resources.GuestCustomization != nil {
		if resp.Status.Resources.GuestCustomization.CloudInit != nil {
			cloudInit["meta_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.MetaData)
			cloudInit["user_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.UserData)
			cloudInit["custom_key_values"] = resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues
		}
		isOverride = utils.BoolValue(resp.Status.Resources.GuestCustomization.IsOverridable)
		if resp.Status.Resources.GuestCustomization.Sysprep != nil {
			sysprep["install_type"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.InstallType)
			sysprep["unattend_xml"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
			sysprep["custom_key_values"] = resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues
		}
	}
	if err := d.Set("guest_customization_cloud_init", cloudInit); err != nil {
		return err
	}
	if err := d.Set("guest_customization_is_overridable", isOverride); err != nil {
		return err
	}
	if err := d.Set("guest_customization_sysprep", sysprep); err != nil {
		return err
	}
	// set power_state_guest_transition_config value
	if err := d.Set("should_fail_on_script_failure", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure)); err != nil {
		return err
	}
	if err := d.Set("enable_script_exec", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec)); err != nil {
		return err
	}
	// set power_state_mechanism value
	if err := d.Set("power_state_mechanism", utils.StringValue(resp.Status.Resources.PowerStateMechanism.Mechanism)); err != nil {
		return err
	}
	// set vga_console_enabled value
	if err := d.Set("vga_console_enabled", utils.BoolValue(resp.Status.Resources.VgaConsoleEnabled)); err != nil {
		return err
	}
	diskList := make([]map[string]interface{}, 0)
	if resp.Status.Resources.DiskList != nil {
		diskList := make([]map[string]interface{}, len(resp.Status.Resources.DiskList))
		for k, v := range resp.Status.Resources.DiskList {
			disk := make(map[string]interface{})
			disk["uuid"] = *v.UUID
			disk["disk_size_bytes"] = *v.DiskSizeBytes
			disk["disk_size_mib"] = *v.DiskSizeMib

			ds := make([]map[string]interface{}, 1)
			dsourceRef := make(map[string]interface{})
			if v.DataSourceReference != nil {
				dsourceRef["kind"] = utils.StringValue(v.DataSourceReference.Kind)
				dsourceRef["name"] = utils.StringValue(v.DataSourceReference.Name)
				dsourceRef["uuid"] = utils.StringValue(v.DataSourceReference.UUID)
			}
			ds[0] = dsourceRef

			disk["data_source_reference"] = ds

			vr := make([]map[string]interface{}, 1)
			volumeRef := make(map[string]interface{})
			if v.VolumeGroupReference != nil {
				volumeRef["kind"] = utils.StringValue(v.VolumeGroupReference.Kind)
				volumeRef["name"] = utils.StringValue(v.VolumeGroupReference.Name)
				volumeRef["uuid"] = utils.StringValue(v.VolumeGroupReference.UUID)
			}
			vr[0] = volumeRef

			disk["volume_group_reference"] = vr

			dp := make([]map[string]interface{}, 1)
			deviceProps := make(map[string]interface{})
			deviceProps["device_type"] = utils.StringValue(v.DeviceProperties.DeviceType)
			dp[0] = deviceProps

			da := make([]map[string]interface{}, 1)
			diskAddress := make(map[string]interface{})
			if v.DeviceProperties.DiskAddress != nil {
				diskAddress["device_index"] = utils.Int64Value(v.DeviceProperties.DiskAddress.DeviceIndex)
				diskAddress["adapter_type"] = utils.StringValue(v.DeviceProperties.DiskAddress.AdapterType)
			}
			da[0] = diskAddress
			deviceProps["disk_address"] = da

			disk["device_properties"] = dp

			diskList[k] = disk
		}
	}
	if err := d.Set("disk_list", diskList); err != nil {
		return err
	}

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API

	log.Printf("Updating VM values %s", d.Id())
	fmt.Printf("Updating VM values %s", d.Id())

	request := &v3.VMIntentInput{}
	metadata := &v3.VMMetadata{}
	res := &v3.VMResources{}
	spec := &v3.VM{}
	guest := &v3.GuestCustomization{}
	guestTool := &v3.GuestToolsSpec{}
	boot := &v3.VMBootConfig{}
	pw := &v3.VMPowerStateMechanism{}

	// get state
	if d.HasChange("metadata") {
		m := d.Get("metadata")
		metad := m.(map[string]interface{})
		if v, ok := metad["uuid"]; ok && v != "" {
			metadata.UUID = utils.String(v.(string))
		}
		if v, ok := metad["spec_version"]; ok && v != 0 {
			i, err := strconv.Atoi(v.(string))
			if err != nil {
				return err
			}
			metadata.SpecVersion = utils.Int64(int64(i))
		}
		if v, ok := metad["spec_hash"]; ok && v != "" {
			metadata.SpecHash = utils.String(v.(string))
		}
		if v, ok := metad["name"]; ok {
			metadata.Name = utils.String(v.(string))
		}
	}

	if d.HasChange("categories") {
		p := d.Get("categories").(map[string]interface{})
		labels := map[string]string{}
		for k, v := range p {
			labels[k] = v.(string)
		}
		metadata.Categories = labels
	}
	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(or["kind"].(string)),
			UUID: utils.String(or["uuid"].(string)),
			Name: utils.String(or["name"].(string)),
		}
		metadata.OwnerReference = r
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
			Name: utils.String(pr["name"].(string)),
		}
		metadata.ProjectReference = r
	}
	if d.HasChange("name") {
		spec.Name = utils.String(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.String(d.Get("description").(string))
	}
	if d.HasChange("availability_zone_reference") {
		a := d.Get("availability_zone_reference").(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
			Name: utils.String(a["name"].(string)),
		}
		spec.AvailabilityZoneReference = r
	}
	if d.HasChange("cluster_reference") {
		a := d.Get("cluster_reference").(map[string]interface{})

		utils.PrintToJSON(a, "CLUSTER VALUES =>")

		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
			Name: utils.String(a["name"].(string)),
		}
		spec.ClusterReference = r
	}
	if d.HasChange("parent_reference") {
		a := d.Get("parent_reference").(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(a["kind"].(string)),
			UUID: utils.String(a["uuid"].(string)),
			Name: utils.String(a["name"].(string)),
		}
		res.ParentReference = r
	}

	if d.HasChange("num_vnuma_nodes") {
		res.VMVnumaConfig = &v3.VMVnumaConfig{
			NumVnumaNodes: utils.Int64(int64(d.Get("num_vnuma_nodes").(int))),
		}
	}
	if d.HasChange("guest_os_id") {
		res.GuestOsID = utils.String(d.Get("guest_os_id").(string))
	}
	if d.HasChange("power_state") {
		res.PowerState = utils.String(d.Get("power_state").(string))
	}
	if d.HasChange("num_vcpus_per_socket") {
		res.NumVcpusPerSocket = utils.Int64(int64(d.Get("num_vcpus_per_socket").(int)))
	}
	if d.HasChange("num_sockets") {
		res.NumSockets = utils.Int64(int64(d.Get("num_sockets").(int)))
	}
	if d.HasChange("memory_size_mib") {
		res.MemorySizeMib = utils.Int64(int64(d.Get("memory_size_mib").(int)))
	}
	if d.HasChange("hardware_clock_timezone") {
		res.HardwareClockTimezone = utils.String(d.Get("hardware_clock_timezone").(string))
	}
	if d.HasChange("vga_console_enabled") {
		res.VgaConsoleEnabled = utils.Bool(d.Get("vga_console_enabled").(bool))
	}
	if d.HasChange("guest_customization_is_overridable") {
		guest.IsOverridable = utils.Bool(d.Get("guest_customization_is_overridable").(bool))
	}
	if d.HasChange("power_state_mechanism") {
		pw.Mechanism = utils.String(d.Get("power_state_mechanism").(string))
	}
	if d.HasChange("power_state_guest_transition_config") {
		val := d.Get("power_state_guest_transition_config").(map[string]interface{})
		pw.GuestTransitionConfig = &v3.VMGuestPowerStateTransitionConfig{
			EnableScriptExec:          utils.Bool(val["enable_script_exec"].(bool)),
			ShouldFailOnScriptFailure: utils.Bool(val["should_fail_on_script_failure"].(bool)),
		}
	}
	if d.HasChange("guest_customization_cloud_init") {
		a := d.Get("guest_customization_cloud_init").(map[string]interface{})
		r := &v3.GuestCustomizationCloudInit{
			MetaData:        utils.String(a["meta_data"].(string)),
			UserData:        utils.String(a["user_data"].(string)),
			CustomKeyValues: a["custom_key_values"].(map[string]string),
		}

		guest.CloudInit = r
	}
	if d.HasChange("guest_customization_sysprep") {
		a := d.Get("guest_customization_sysprep").(map[string]interface{})
		r := &v3.GuestCustomizationSysprep{
			InstallType:     utils.String(a["install_type"].(string)),
			UnattendXML:     utils.String(a["unattend_xml"].(string)),
			CustomKeyValues: a["custom_key_values"].(map[string]string),
		}

		guest.Sysprep = r
	}
	if d.HasChange("nic_list") {
		n := d.Get("nic_list").([]interface{})
		if len(n) > 0 {
			nics := make([]*v3.VMNic, len(n))

			for k, nc := range n {
				val := nc.(map[string]interface{})
				net := val["network_function_chain_reference"].(map[string]interface{})
				sub := val["subnet_reference"].(map[string]interface{})

				nic := &v3.VMNic{
					NicType: utils.String(val["nic_type"].(string)),
					UUID:    utils.String(val["uuid"].(string)),
					NetworkFunctionNicType: utils.String(val["network_function_nic_type"].(string)),
					MacAddress:             utils.String(val["mac_address"].(string)),
					Model:                  utils.String(val["model"].(string)),
					NetworkFunctionChainReference: &v3.Reference{
						Kind: utils.String(net["kind"].(string)),
						UUID: utils.String(net["uuid"].(string)),
						Name: utils.String(net["name"].(string)),
					},
					SubnetReference: &v3.Reference{
						Kind: utils.String(sub["kind"].(string)),
						UUID: utils.String(sub["uuid"].(string)),
						Name: utils.String(sub["name"].(string)),
					},
				}

				if value, ok := val["ip_endpoint_list"]; ok {
					ipl := value.([]interface{})
					if len(ipl) > 0 {
						ip := make([]*v3.IPAddress, len(ipl))
						for k, i := range ipl {
							v := i.(map[string]interface{})
							v3ip := &v3.IPAddress{
								IP:   utils.String(v["ip"].(string)),
								Type: utils.String(v["type"].(string)),
							}
							ip[k] = v3ip
						}
						nic.IPEndpointList = ip
					}
				}

				nics[k] = nic
			}
			res.NicList = nics
		}
	}
	if d.HasChange("nic_list") {
		ngt := d.Get("nutanix_guest_tools").(map[string]interface{})

		tool := &v3.NutanixGuestToolsSpec{
			IsoMountState: utils.String(ngt["iso_mount_state"].(string)),
			State:         utils.String(ngt["state"].(string)),
		}

		if val, ok2 := ngt["enabled_capability_list"]; ok2 {
			var l []*string
			for _, list := range val.([]interface{}) {
				l = append(l, utils.String(list.(string)))
			}
			tool.EnabledCapabilityList = l
		}
		guestTool.NutanixGuestTools = tool
	}
	if d.HasChange("gpu_list") {
		if v, ok := d.GetOk("gpu_list"); ok {
			gpl := make([]*v3.VMGpu, len(v.([]interface{})))

			for k, va := range v.([]interface{}) {
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
				gpl[k] = gpu
			}
			res.GpuList = gpl
		}
	}
	if d.HasChange("boot_device_order_list") {
		var b []*string
		for _, boot := range d.Get("boot_device_order_list").([]interface{}) {
			b = append(b, utils.String(boot.(string)))
		}
		boot.BootDeviceOrderList = b
	}

	bd := &v3.VMBootDevice{}
	if d.HasChange("boot_device_disk_address") {
		dai := d.Get("boot_device_disk_address").(map[string]interface{})
		da := &v3.DiskAddress{}
		if value3, ok3 := dai["device_index"]; ok3 {
			da.DeviceIndex = utils.Int64(int64(value3.(int)))
		}
		if value3, ok3 := dai["adapter_type"]; ok3 {
			da.AdapterType = utils.String(value3.(string))
		}
		bd.DiskAddress = da
	}

	if d.HasChange("boot_device_mac_address") {
		v := d.Get("boot_device_mac_address").(string)
		bd.MacAddress = utils.String(v)
	}

	if d.HasChange("disk_list") {
		if v, ok := d.GetOk("disk_list"); ok {
			dsk := v.([]interface{})
			if len(dsk) > 0 {
				dls := make([]*v3.VMDisk, len(dsk))

				for k, val := range dsk {
					v := val.(map[string]interface{})
					dl := &v3.VMDisk{}
					if v1, ok1 := v["uuid"]; ok1 && v1.(string) != "" {
						dl.UUID = utils.String(v1.(string))
					}
					if v1, ok1 := v["disk_size_bytes"]; ok1 && v1.(int) != 0 {
						dl.DiskSizeBytes = utils.Int64(int64(v1.(int)))
					}
					if v1, ok := v["disk_size_mib"]; ok && v1.(int) != 0 {
						dl.DiskSizeMib = utils.Int64(int64(v1.(int)))
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
								if len(v.([]interface{})) > 0 {
									da := v.([]interface{})[0].(map[string]interface{})
									v3disk := &v3.DiskAddress{}
									if di, diok := da["device_index"]; diok {
										v3disk.DeviceIndex = utils.Int64(int64(di.(int)))
									}
									if di, diok := da["adapter_type"]; diok {
										v3disk.AdapterType = utils.String(di.(string))
									}
									dp.DiskAddress = v3disk
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
					dls[k] = dl
				}
				res.DiskList = dls
			}
		}
	}

	boot.BootDevice = bd
	res.PowerStateMechanism = pw
	res.BootConfig = boot
	res.GuestTools = guestTool
	res.GuestCustomization = guest
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	log.Printf("[DEBUG] Updating Virtual Machine: %s, %s", d.Get("name").(string), d.Id())
	fmt.Printf("[DEBUG] Updating Virtual Machine: %s, %s", d.Get("name").(string), d.Id())

	utils.PrintToJSON(request, "UPDATE")
	_, err := conn.V3.UpdateVM(d.Id(), request)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING"},
		Target:     []string{"COMPLETE"},
		Refresh:    vmStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for vm (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixVirtualMachineRead(d, meta)
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NutanixClient).API

	log.Printf("[DEBUG] Deleting Virtual Machine: %s, %s", d.Get("name").(string), d.Id())
	if err := conn.V3.DeleteVM(d.Id()); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING", "RUNNING", "DELETE_IN_PROGRESS", "COMPLETE"},
		Target:     []string{"DELETED"},
		Refresh:    vmStateRefreshFunc(conn, d.Id()),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for vm (%s) to delete: %s", d.Id(), err)
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

func getMetadaAttributes(d *schema.ResourceData, metadata *v3.VMMetadata) error {
	m, mok := d.GetOk("metadata")
	metad := m.(map[string]interface{})

	if !mok {
		return fmt.Errorf("please provide metadata required attributes")
	}

	metadata.Kind = utils.String(metad["kind"].(string))

	if v, ok := metad["uuid"]; ok && v != "" {
		metadata.UUID = utils.String(v.(string))
	}
	if v, ok := metad["spec_version"]; ok && v != 0 {
		i, err := strconv.Atoi(v.(string))
		if err != nil {
			return err
		}
		metadata.SpecVersion = utils.Int64(int64(i))
	}
	if v, ok := metad["spec_hash"]; ok && v != "" {
		metadata.SpecHash = utils.String(v.(string))
	}
	if v, ok := metad["name"]; ok {
		metadata.Name = utils.String(v.(string))
	}
	if v, ok := d.GetOk("categories"); ok {
		c := v.(map[string]interface{})
		labels := map[string]string{}

		for k, v := range c {
			labels[k] = v.(string)
		}
		metadata.Categories = labels
	}
	if p, ok := d.GetOk("project_reference"); ok {
		pr := p.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(pr["kind"].(string)),
			UUID: utils.String(pr["uuid"].(string)),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.ProjectReference = r
	}
	if o, ok := metad["owner_reference"]; ok {
		or := o.(map[string]interface{})
		r := &v3.Reference{
			Kind: utils.String(or["kind"].(string)),
			UUID: utils.String(or["uuid"].(string)),
		}
		if v1, ok1 := or["name"]; ok1 {
			r.Name = utils.String(v1.(string))
		}
		metadata.OwnerReference = r
	}

	return nil
}

func getVMResources(d *schema.ResourceData, vm *v3.VMResources) error {
	if v, ok := d.GetOk("num_vnuma_nodes"); ok {
		vm.VMVnumaConfig.NumVnumaNodes = utils.Int64(v.(int64))
	}
	if v, ok := d.GetOk("nic_list"); ok {
		n := v.([]interface{})
		if len(n) > 0 {
			nics := make([]*v3.VMNic, len(n))

			for k, nc := range n {
				val := nc.(map[string]interface{})
				nic := &v3.VMNic{}

				if value, ok := val["nic_type"]; ok && value.(string) != "" {
					nic.NicType = utils.String(value.(string))
				}
				if value, ok := val["uuid"]; ok && value.(string) != "" {
					nic.UUID = utils.String(value.(string))
				}
				if value, ok := val["network_function_nic_type"]; ok && value.(string) != "" {
					nic.NetworkFunctionNicType = utils.String(value.(string))
				}
				if value, ok := val["mac_address"]; ok && value.(string) != "" {
					nic.MacAddress = utils.String(value.(string))
				}
				if value, ok := val["model"]; ok && value.(string) != "" {
					nic.Model = utils.String(value.(string))
				}
				if value, ok := val["ip_endpoint_list"]; ok {
					ipl := value.([]interface{})
					if len(ipl) > 0 {
						ip := make([]*v3.IPAddress, len(ipl))
						for k, i := range ipl {
							v := i.(map[string]interface{})
							v3ip := &v3.IPAddress{}

							if ipset, ipsetok := v["ip"]; ipsetok {
								v3ip.IP = utils.String(ipset.(string))
							}
							if iptype, iptypeok := v["type"]; iptypeok {
								v3ip.Type = utils.String(iptype.(string))
							}
							ip[k] = v3ip
						}
						nic.IPEndpointList = ip
					}
				}
				if value, ok := val["network_function_chain_reference"]; ok && len(value.(map[string]interface{})) != 0 {
					v := value.(map[string]interface{})
					ref := &v3.Reference{}
					if j, ok1 := v["kind"]; ok1 {
						ref.Kind = utils.String(j.(string))
					}
					if j, ok1 := v["uuid"]; ok1 {
						ref.UUID = utils.String(j.(string))
					}
					if j, ok1 := v["name"]; ok1 {
						ref.Name = utils.String(j.(string))
					}
					nic.NetworkFunctionChainReference = ref
				}
				if value, ok := val["subnet_reference"]; ok {
					v := value.(map[string]interface{})
					ref := &v3.Reference{}

					if j, ok1 := v["kind"]; ok1 {
						ref.Kind = utils.String(j.(string))
					}
					if j, ok1 := v["uuid"]; ok1 {
						ref.UUID = utils.String(j.(string))
					}
					if j, ok1 := v["name"]; ok1 {
						ref.Name = utils.String(j.(string))
					}
					nic.SubnetReference = ref
				}
				nics[k] = nic
			}
			vm.NicList = nics
		}
	}
	if v, ok := d.GetOk("guest_os_id"); ok {
		vm.GuestOsID = utils.String(v.(string))
	}
	if v, ok := d.GetOk("power_state"); ok {
		vm.PowerState = utils.String(v.(string))
	}
	if v, ok := d.GetOk("nutanix_guest_tools"); ok {
		ngt := v.(map[string]interface{})

		if val, ok2 := ngt["iso_mount_state"]; ok2 {
			vm.GuestTools.NutanixGuestTools.IsoMountState = utils.String(val.(string))
		}
		if val, ok2 := ngt["state"]; ok2 {
			vm.GuestTools.NutanixGuestTools.State = utils.String(val.(string))
		}
		if val, ok2 := ngt["enabled_capability_list"]; ok2 {
			var l []*string
			for _, list := range val.([]interface{}) {
				l = append(l, utils.String(list.(string)))
			}
			vm.GuestTools.NutanixGuestTools.EnabledCapabilityList = l
		}
	}
	if v, ok := d.GetOk("num_vcpus_per_socket"); ok {
		vm.NumVcpusPerSocket = utils.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("num_sockets"); ok {
		vm.NumSockets = utils.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("gpu_list"); ok {
		gpl := make([]*v3.VMGpu, len(v.([]interface{})))

		for k, va := range v.([]interface{}) {
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
			gpl[k] = gpu
		}
		vm.GpuList = gpl
	}
	if v, ok := d.GetOk("parent_reference"); ok {
		val := v.(map[string]string)
		vm.ParentReference.Kind = utils.String(val["kind"])
		vm.ParentReference.UUID = utils.String(val["uuid"])
		if j, ok1 := val["name"]; ok1 {
			vm.ParentReference.Name = utils.String(j)
		}
	}
	if v, ok := d.GetOk("memory_size_mib"); ok {
		vm.MemorySizeMib = utils.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("boot_device_order_list"); ok {
		var b []*string
		for _, boot := range v.([]interface{}) {
			b = append(b, utils.String(boot.(string)))
		}
		vm.BootConfig.BootDeviceOrderList = b
	}

	bd := &v3.VMBootDevice{}
	da := &v3.DiskAddress{}
	if v, ok := d.GetOk("boot_device_disk_address"); ok {
		dai := v.(map[string]interface{})

		if value3, ok3 := dai["device_index"]; ok3 {
			da.DeviceIndex = utils.Int64(int64(value3.(int)))
		}
		if value3, ok3 := dai["adapter_type"]; ok3 {
			da.AdapterType = utils.String(value3.(string))
		}
		bd.DiskAddress = da
		vm.BootConfig.BootDevice = bd
	}

	if v, ok := d.GetOk("boot_device_mac_address"); ok {
		bdi := v.(string)
		bd.MacAddress = utils.String(bdi)
		vm.BootConfig.BootDevice = bd
	}

	if v, ok := d.GetOk("hardware_clock_timezone"); ok {
		vm.HardwareClockTimezone = utils.String(v.(string))
	}
	if v, ok := d.GetOk("guest_customization_cloud_init"); ok {
		cii := v.(map[string]interface{})
		if v2, ok2 := cii["meta_data"]; ok2 {
			vm.GuestCustomization.CloudInit.MetaData = utils.String(v2.(string))
		}
		if v2, ok2 := cii["user_data"]; ok2 {
			vm.GuestCustomization.CloudInit.UserData = utils.String(v2.(string))
		}
		if v2, ok2 := cii["custom_key_values"]; ok2 {
			vm.GuestCustomization.CloudInit.CustomKeyValues = v2.(map[string]string)
		}
	}
	if v, ok := d.GetOk("guest_customization_is_overridable"); ok {
		vm.GuestCustomization.IsOverridable = utils.Bool(v.(bool))
	}
	if v, ok := d.GetOk("guest_customization_sysprep"); ok {
		spi := v.(map[string]interface{})
		if v2, ok2 := spi["install_type"]; ok2 {
			vm.GuestCustomization.Sysprep.InstallType = utils.String(v2.(string))
		}
		if v2, ok2 := spi["unattend_xml"]; ok2 {
			vm.GuestCustomization.Sysprep.UnattendXML = utils.String(v2.(string))
		}
		if v2, ok2 := spi["custom_key_values"]; ok2 {
			vm.GuestCustomization.Sysprep.CustomKeyValues = v2.(map[string]string)
		}
	}
	if v, ok := d.GetOk("vga_console_enabled"); ok {
		vm.VgaConsoleEnabled = utils.Bool(v.(bool))
	}
	if v, ok := d.GetOk("power_state_mechanism"); ok {
		vm.PowerStateMechanism.Mechanism = utils.String(v.(string))
	}
	if v, ok := d.GetOk("should_fail_on_script_failure"); ok {
		vm.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure = utils.Bool(v.(bool))
	}
	if v, ok := d.GetOk("enable_script_exec"); ok {
		vm.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec = utils.Bool(v.(bool))
	}
	if v, ok := d.GetOk("disk_list"); ok {
		dsk := v.([]interface{})
		if len(dsk) > 0 {
			dls := make([]*v3.VMDisk, len(dsk))

			for k, val := range dsk {
				v := val.(map[string]interface{})
				dl := &v3.VMDisk{}
				if v1, ok1 := v["uuid"]; ok1 && v1.(string) != "" {
					dl.UUID = utils.String(v1.(string))
				}
				if v1, ok1 := v["disk_size_bytes"]; ok1 && v1.(int) != 0 {
					dl.DiskSizeBytes = utils.Int64(int64(v1.(int)))
				}
				if v1, ok := v["disk_size_mib"]; ok && v1.(int) != 0 {
					dl.DiskSizeMib = utils.Int64(int64(v1.(int)))
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
							if len(v.([]interface{})) > 0 {
								da := v.([]interface{})[0].(map[string]interface{})
								v3disk := &v3.DiskAddress{}
								if di, diok := da["device_index"]; diok {
									v3disk.DeviceIndex = utils.Int64(int64(di.(int)))
								}
								if di, diok := da["adapter_type"]; diok {
									v3disk.AdapterType = utils.String(di.(string))
								}
								dp.DiskAddress = v3disk
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
				dls[k] = dl
			}
			vm.DiskList = dls
		}
	}

	return nil
}

func vmStateRefreshFunc(client *v3.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := client.V3.GetVM(uuid)

		if err != nil {
			if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
				return v, "DELETED", nil
			}
			log.Printf("ERROR %s", err)
			return nil, "", err
		}

		return v, *v.Status.State, nil
	}
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
			Type:     schema.TypeMap,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"kind": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": &schema.Schema{
						Type:     schema.TypeString,
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
		"categories": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"project_reference": &schema.Schema{
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
		"owner_reference": &schema.Schema{
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
		"api_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
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

		// COMPUTED
		"state": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_reference": &schema.Schema{
			Type:     schema.TypeMap,
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
		"hypervisor_type": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},

		// RESOURCES ARGUMENTS

		"num_vnuma_nodes": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
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
					"floating_ip": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
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
						Type:     schema.TypeMap,
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
					"subnet_reference": &schema.Schema{
						Type:     schema.TypeMap,
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
		"nutanix_guest_tools": &schema.Schema{
			Type:     schema.TypeMap,
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
			Type:     schema.TypeMap,
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
		"boot_device_order_list": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"boot_device_disk_address": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_index": &schema.Schema{
						Type:     schema.TypeString,
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
		"boot_device_mac_address": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"hardware_clock_timezone": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"guest_customization_cloud_init": &schema.Schema{
			Type:     schema.TypeMap,
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
		"guest_customization_is_overridable": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"guest_customization_sysprep": &schema.Schema{
			Type:     schema.TypeMap,
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
		"power_state_mechanism": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
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
					"disk_size_mib": &schema.Schema{
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
	}
}
