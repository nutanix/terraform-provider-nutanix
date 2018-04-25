package nutanix

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/terraform-providers/terraform-provider-nutanix/client/v3"

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
	request := v3.VMIntentInput{
		APIVersion: version,
	}

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
		r := v3.Reference{
			Kind: a["kind"].(string),
			UUID: a["uuid"].(string),
		}
		if v, ok := a["name"]; ok {
			r.Name = v.(string)
		}
		request.Spec.AvailabilityZoneReference = r
	}
	if descok {
		request.Spec.Description = desc.(string)
	}
	if crok {
		a := cr.(map[string]interface{})
		r := v3.Reference{
			Kind: a["kind"].(string),
			UUID: a["uuid"].(string),
		}
		if v, ok := a["name"]; ok {
			r.Name = v.(string)
		}
		request.Spec.ClusterReference = r
	}

	request.Spec.Name = n.(string)

	metad := m.(map[string]interface{})
	metadata := v3.VMMetadata{
		Kind: metad["kind"].(string),
	}
	if v, ok := metad["uuid"]; ok {
		metadata.UUID = v.(string)
	}
	if v, ok := metad["spec_version"]; ok {
		metadata.SpecVersion = int64(v.(int))
	}
	if v, ok := metad["spec_hash"]; ok {
		metadata.SpecHash = v.(string)
	}
	if v, ok := metad["name"]; ok {
		metadata.Name = v.(string)
	}
	if v, ok := metad["categories"]; ok {
		metadata.Categories = v.(map[string]string)
	}
	if v, ok := metad["project_reference"]; ok {
		pr := v.(map[string]interface{})
		r := v3.Reference{
			Kind: pr["kind"].(string),
			UUID: pr["uuid"].(string),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = v1.(string)
		}
		metadata.ProjectReference = r
	}
	if v, ok := metad["owner_reference"]; ok {
		pr := v.(map[string]interface{})
		r := v3.Reference{
			Kind: pr["kind"].(string),
			UUID: pr["uuid"].(string),
		}
		if v1, ok1 := pr["name"]; ok1 {
			r.Name = v1.(string)
		}
		metadata.OwnerReference = r
	}

	request.Metadata = metadata
	request.Spec.Resources = setVMResources(r)

	// Make request to the API
	resp, err := conn.V3.CreateVM(request)
	if err != nil {
		return err
	}

	uuid := resp.Metadata.UUID

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
	if resp.Spec.Resources.NicList != nil && resp.Spec.Resources.PowerState == "ON" {
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

	// Set vm values
	// set availability zone reference values
	availabilityZoneReference := make(map[string]interface{})
	availabilityZoneReference["kind"] = resp.Status.AvailabilityZoneReference.Kind
	availabilityZoneReference["name"] = resp.Status.AvailabilityZoneReference.Name
	availabilityZoneReference["uuid"] = resp.Status.AvailabilityZoneReference.UUID

	// set message list values
	messages := make([]map[string]interface{}, len(resp.Status.MessageList))
	for k, v := range resp.Status.MessageList {
		message := make(map[string]interface{})

		message["message"] = v.Message
		message["reason"] = v.Reason
		message["details"] = v.Details

		messages[k] = message
	}

	// set cluster reference values
	clusterReference := make(map[string]interface{})
	clusterReference["kind"] = resp.Status.ClusterReference.Kind
	clusterReference["name"] = resp.Status.ClusterReference.Name
	clusterReference["uuid"] = resp.Status.ClusterReference.UUID

	// set resources values
	resouces := make(map[string]interface{})

	vnumaConfig := make(map[string]interface{})
	vnumaConfig["num_vnuma_nodes"] = resp.Status.Resources.VnumaConfig.NumVnumaNodes

	resouces["vnuma_config"] = vnumaConfig

	nics := resp.Status.Resources.NicList

	nicLists := make([]map[string]interface{}, len(nics))
	for k, v := range nics {
		nic := make(map[string]interface{})
		// simple firts
		nic["nic_type"] = v.NicType
		nic["uuid"] = v.UUID
		nic["floating_ip"] = v.FloatingIP
		nic["network_function_nic_type"] = v.NetworkFunctionNicType
		nic["mac_address"] = v.MacAddress
		nic["model"] = v.Model

		ipEndpointList := make([]map[string]interface{}, len(v.IPEndpointList))
		for k1, v1 := range v.IPEndpointList {
			ipEndpoint := make(map[string]interface{})
			ipEndpoint["ip"] = v1.IP
			ipEndpoint["type"] = v1.Type
			ipEndpointList[k1] = ipEndpoint
		}
		nic["ip_endpoint_list"] = ipEndpointList

		netFnChainRef := make(map[string]interface{})
		netFnChainRef["kind"] = v.NetworkFunctionChainReference.Kind
		netFnChainRef["name"] = v.NetworkFunctionChainReference.Name
		netFnChainRef["uuid"] = v.NetworkFunctionChainReference.UUID

		nic["network_function_chain_reference"] = netFnChainRef

		subtnetRef := make(map[string]interface{})
		subtnetRef["kind"] = v.SubnetReference.Kind
		subtnetRef["name"] = v.SubnetReference.Name
		subtnetRef["uuid"] = v.SubnetReference.UUID

		nic["subnet_reference"] = subtnetRef

		nicLists[k] = nic
	}

	resouces["nic_list"] = nicLists
	hostRef := make(map[string]interface{})
	hostRef["kind"] = resp.Status.Resources.HostReference.Kind
	hostRef["name"] = resp.Status.Resources.HostReference.Name
	hostRef["uuid"] = resp.Status.Resources.HostReference.UUID

	resouces["host_reference"] = hostRef

	guestTools := make(map[string]interface{})

	tools := resp.Status.Resources.GuestTools.NutanixGuestTools
	nutanixGuestTools := make(map[string]interface{})
	nutanixGuestTools["available_version"] = tools.AvailableVersion
	nutanixGuestTools["iso_mount_state"] = tools.IsoMountState
	nutanixGuestTools["state"] = tools.State
	nutanixGuestTools["version"] = tools.Version
	nutanixGuestTools["guest_os_version"] = tools.GuestOsVersion

	capList := make([]string, len(tools.EnabledCapabilityList))
	for k, v := range tools.EnabledCapabilityList {
		capList[k] = v
	}
	nutanixGuestTools["enabled_capability_list"] = capList
	nutanixGuestTools["vss_snapshot_capable"] = tools.VSSSnapshotCapable
	nutanixGuestTools["is_reachable"] = tools.IsReachable
	nutanixGuestTools["vm_mobility_drivers_installed"] = tools.VMMobilityDriversInstalled

	guestTools["nutanix_guest_tools"] = nutanixGuestTools

	resouces["guest_tools"] = guestTools

	gpuList := make([]map[string]interface{}, len(resp.Status.Resources.GpuList))
	for k, v := range resp.Status.Resources.GpuList {
		gpu := make(map[string]interface{})
		gpu["frame_buffer_size_mib"] = v.FrameBufferSizeMib
		gpu["vendor"] = v.Vendor
		gpu["uuid"] = v.UUID
		gpu["name"] = v.Name
		gpu["pci_address"] = v.PCIAddress
		gpu["fraction"] = v.Fraction
		gpu["mode"] = v.Mode
		gpu["num_virtual_display_heads"] = v.NumVirtualDisplayHeads
		gpu["guest_driver_version"] = v.GuestDriverVersion
		gpu["device_id"] = v.DeviceID

		gpuList[k] = gpu
	}

	resouces["gpu_list"] = gpuList

	parentRef := make(map[string]interface{})
	parentRef["kind"] = resp.Status.Resources.ParentReference.Kind
	parentRef["name"] = resp.Status.Resources.ParentReference.Name
	parentRef["uuid"] = resp.Status.Resources.ParentReference.UUID

	resouces["parent_reference"] = parentRef

	bootConfig := make(map[string]interface{})
	boots := make([]string, len(resp.Status.Resources.BootConfig.BootDeviceOrderList))
	for k, v := range resp.Status.Resources.BootConfig.BootDeviceOrderList {
		boots[k] = v
	}
	bootDevice := make(map[string]interface{})
	diskAddress := make(map[string]interface{})
	diskAddress["device_index"] = resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex
	diskAddress["adapter_type"] = resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType
	bootDevice["disk_address"] = diskAddress
	bootDevice["mac_address"] = resp.Status.Resources.BootConfig.BootDevice.MacAddress

	bootConfig["boot_device"] = bootDevice
	bootConfig["boot_device_order_list"] = boots

	resouces["boot_config"] = bootConfig

	guestCustom := make(map[string]interface{})
	cloudInit := make(map[string]interface{})
	cloudInit["meta_data"] = resp.Status.Resources.GuestCustomization.CloudInit.MetaData
	cloudInit["user_data"] = resp.Status.Resources.GuestCustomization.CloudInit.UserData
	cloudInit["custom_key_values"] = resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues

	guestCustom["cloud_init"] = cloudInit
	guestCustom["is_overridable"] = resp.Status.Resources.GuestCustomization.IsOverridable

	sysprep := make(map[string]interface{})
	sysprep["install_type"] = resp.Status.Resources.GuestCustomization.Sysprep.InstallType
	sysprep["unattend_xml"] = resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML
	sysprep["custom_key_values"] = resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues

	guestCustom["sysprep"] = sysprep

	resouces["guest_customization"] = guestCustom

	powerStateMechanism := make(map[string]interface{})
	powerStateMechanism["mechanism"] = resp.Status.Resources.PowerStateMechanism.Mechanism

	guestTransition := make(map[string]interface{})
	guestTransition["should_fail_on_script_failure"] = resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure
	guestTransition["enable_script_exec"] = resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec

	powerStateMechanism["guest_transition_config"] = guestTransition

	resouces["power_state_mechanism"] = powerStateMechanism

	diskList := make([]map[string]interface{}, len(resp.Status.Resources.DiskList))
	for k, v := range resp.Status.Resources.DiskList {
		disk := make(map[string]interface{})
		disk["uuid"] = v.UUID
		disk["disk_size_bytes"] = v.DiskSizeBytes
		disk["disk_size_mib"] = v.DiskSizeMib

		dsourceRef := make(map[string]interface{})
		dsourceRef["kind"] = v.DataSourceReference.Kind
		dsourceRef["name"] = v.DataSourceReference.Name
		dsourceRef["uuid"] = v.DataSourceReference.UUID

		disk["data_source_reference"] = dsourceRef

		volumeRef := make(map[string]interface{})
		volumeRef["kind"] = v.VolumeGroupReference.Kind
		volumeRef["name"] = v.VolumeGroupReference.Name
		volumeRef["uuid"] = v.VolumeGroupReference.UUID

		disk["volume_group_reference"] = volumeRef

		deviceProps := make(map[string]interface{})
		deviceProps["device_type"] = v.DeviceProperties.DeviceType

		diskAddress := make(map[string]interface{})
		diskAddress["device_index"] = v.DeviceProperties.DiskAddress.DeviceIndex
		diskAddress["adapter_type"] = v.DeviceProperties.DiskAddress.AdapterType

		deviceProps["disk_address"] = diskAddress

		disk["device_properties"] = deviceProps

		diskList[k] = disk
	}

	resouces["disk_list"] = diskList

	// set metadata values
	metadata := make(map[string]interface{})
	metadata["last_update_time"] = resp.Metadata.LastUpdateTime
	metadata["kind"] = resp.Metadata.Kind
	metadata["uuid"] = resp.Metadata.UUID
	metadata["creation_time"] = resp.Metadata.CreationTime
	metadata["spec_version"] = resp.Metadata.SpecVersion
	metadata["spec_hash"] = resp.Metadata.SpecHash
	metadata["categories"] = resp.Metadata.Categories
	metadata["name"] = resp.Metadata.Name

	pr := make(map[string]interface{})
	pr["kind"] = resp.Metadata.ProjectReference.Kind
	pr["name"] = resp.Metadata.ProjectReference.Name
	pr["uuid"] = resp.Metadata.ProjectReference.UUID

	or := make(map[string]interface{})
	or["kind"] = resp.Metadata.OwnerReference.Kind
	or["name"] = resp.Metadata.OwnerReference.Name
	or["uuid"] = resp.Metadata.OwnerReference.UUID

	metadata["project_reference"] = pr
	metadata["owner_reference"] = or

	// Simple first
	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}
	if err := d.Set("name", resp.Status.Name); err != nil {
		return err
	}
	if err := d.Set("state", resp.Status.State); err != nil {
		return err
	}
	if err := d.Set("description", resp.Status.Description); err != nil {
		return err
	}
	if err := d.Set("availability_zone_reference", availabilityZoneReference); err != nil {
		return err
	}
	if err := d.Set("message_list", messages); err != nil {
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
	vmSpec := setVMResources(spec)

	log.Printf("[DEBUG] Updating Virtual Machine: %s, %s", name, uuid)

	d.Partial(true)

	if d.HasChange("name") || d.HasChange("resources") || d.HasChange("metadata") {
		request := v3.VMIntentInput{}
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

	getEntitiesRequest := v3.VMListMetadata{}

	resp, err := conn.V3.ListVM(getEntitiesRequest)
	if err != nil {
		return false, err
	}

	for i := range resp.Entities {
		if resp.Entities[i].Metadata.UUID == d.Id() {
			return true, nil
		}
	}
	return false, nil
}

func setVMResources(m interface{}) v3.VMResources {

	vm := v3.VMResources{}

	resources := m.(map[string]interface{})

	if v, ok := resources["vnuma_config"]; ok {
		vm.VMVnumaConfig.NumVnumaNodes = int64(v.(int))
	}

	if v, ok := resources["nic_list"]; ok {
		var nics []v3.VMNic

		for _, val := range v.([]map[string]interface{}) {
			nic := v3.VMNic{}

			if value, ok := val["nic_type"]; ok {
				nic.NicType = value.(string)
			}
			if value, ok := val["uuid"]; ok {
				nic.UUID = value.(string)
			}
			if value, ok := val["network_function_nic_type"]; ok {
				nic.NetworkFunctionNicType = value.(string)
			}
			if value, ok := val["mac_address"]; ok {
				nic.MacAddress = value.(string)
			}
			if value, ok := val["model"]; ok {
				nic.Model = value.(string)
			}
			if value, ok := val["ip_endpoint_list"]; ok {
				var ip []v3.IPAddress
				for _, v := range value.([]map[string]interface{}) {
					ip = append(ip, v3.IPAddress{IP: v["ip"].(string), Type: v["type"].(string)})
				}
				nic.IPEndpointList = ip
			}
			if value, ok := val["network_function_chain_reference"]; ok {
				v := value.(map[string]string)
				nic.NetworkFunctionChainReference.Kind = v["kind"]
				nic.NetworkFunctionChainReference.UUID = v["uuid"]
				if j, ok1 := v["name"]; ok1 {
					nic.NetworkFunctionChainReference.Name = j
				}
			}
			if value, ok := val["subnet_reference"]; ok {
				v := value.(map[string]string)
				nic.SubnetReference.Kind = v["kind"]
				nic.SubnetReference.UUID = v["uuid"]
				if j, ok1 := v["name"]; ok1 {
					nic.SubnetReference.Name = j
				}
			}

			nics = append(nics, nic)
		}

		vm.NicList = nics
	}
	if v, ok := resources["guest_tools"]; ok {
		ngt := v.(map[string]interface{})
		if k, ok1 := ngt["nutanix_guest_tools"]; ok1 {
			ngts := k.(map[string]interface{})
			if val, ok2 := ngts["iso_mount_state"]; ok2 {
				vm.GuestTools.NutanixGuestTools.IsoMountState = val.(string)
			}
			if val, ok2 := ngts["state"]; ok2 {
				vm.GuestTools.NutanixGuestTools.State = val.(string)
			}
			if val, ok2 := ngts["enabled_capability_list"]; ok2 {
				var l []string
				for _, list := range val.([]interface{}) {
					l = append(l, list.(string))
				}
				vm.GuestTools.NutanixGuestTools.EnabledCapabilityList = l
			}
		}
	}
	if v, ok := resources["gpu_list"]; ok {
		var gpl []v3.VMGpu
		for _, val := range v.([]map[string]interface{}) {
			gpu := v3.VMGpu{}
			if value, ok1 := val["vendor"]; ok1 {
				gpu.Vendor = value.(string)
			}
			if value, ok1 := val["device_id"]; ok1 {
				gpu.DeviceID = int64(value.(int))
			}
			if value, ok1 := val["mode"]; ok1 {
				gpu.Mode = value.(string)
			}
			gpl = append(gpl, gpu)
		}
		vm.GpuList = gpl
	}
	if v, ok := resources["parent_reference"]; ok {
		val := v.(map[string]string)
		vm.ParentReference.Kind = val["kind"]
		vm.ParentReference.UUID = val["uuid"]
		if j, ok1 := val["name"]; ok1 {
			vm.ParentReference.Name = j
		}
	}
	if v, ok := resources["boot_config"]; ok {
		val := v.(map[string]interface{})
		if value1, ok1 := val["boot_device_order_list"]; ok1 {
			var b []string
			for _, boot := range value1.([]interface{}) {
				b = append(b, boot.(string))
			}
			vm.BootConfig.BootDeviceOrderList = b
		}
		if value1, ok1 := val["boot_device"]; ok1 {
			bdi := value1.(map[string]interface{})
			bd := v3.VMBootDevice{}
			if value2, ok2 := bdi["disk_address"]; ok2 {
				dai := value2.(map[string]interface{})
				da := v3.DiskAddress{}
				if value3, ok3 := dai["device_index"]; ok3 {
					da.DeviceIndex = int64(value3.(int))
				}
				if value3, ok3 := dai["adapter_type"]; ok3 {
					da.AdapterType = value3.(string)
				}
				bd.DiskAddress = da
			}
			if value2, ok2 := bdi["mac_address"]; ok2 {
				bd.MacAddress = value2.(string)
			}
			vm.BootConfig.BootDevice = bd
		}
	}

	if v, ok := resources["guest_customization"]; ok {
		gci := v.(map[string]interface{})
		gc := v3.GuestCustomization{}

		if v1, ok1 := gci["cloud_init"]; ok1 {
			cii := v1.(map[string]interface{})
			if v2, ok2 := cii["meta_data"]; ok2 {
				gc.CloudInit.MetaData = v2.(string)
			}
			if v2, ok2 := cii["user_data"]; ok2 {
				gc.CloudInit.UserData = v2.(string)
			}
			if v2, ok2 := cii["custom_key_values"]; ok2 {
				gc.CloudInit.CustomKeyValues = v2.(map[string]string)
			}
		}
		if v1, ok1 := gci["sysprep"]; ok1 {
			spi := v1.(map[string]interface{})
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
		if v1, ok1 := gci["is_overridable"]; ok1 {
			gc.IsOverridable = v1.(bool)
		}

		vm.GuestCustomization = gc
	}

	return vm
}

func waitForVMProcess(conn *v3.Client, uuid string) (bool, error) {
	for {
		resp, err := conn.V3.GetVM(uuid)
		if err != nil {
			return false, err
		}

		if resp.Status.State == "COMPLETE" {
			return true, nil
		} else if resp.Status.State == "ERROR" {
			return false, fmt.Errorf("Error while waiting for resource to be up")
		}
		time.Sleep(3000 * time.Millisecond)
	}
	return false, nil
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
					if ip := resp.Status.Resources.NicList[i].IPEndpointList[0].IP; ip != "" {
						// TODO set ip address
						d.Set("ip_address", ip)
						return nil
					}
				}
			}
		}
		time.Sleep(3000 * time.Millisecond)
	}
	return nil
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
					"categories": &schema.Schema{
						Type:     schema.TypeMap,
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
						Type:     schema.TypeMap,
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
									Type:     schema.TypeString,
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
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
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
					"boot_config": &schema.Schema{
						Type:     schema.TypeMap,
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
									Type:     schema.TypeMap,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"disk_address": &schema.Schema{
												Type:     schema.TypeMap,
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
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cloud_init": &schema.Schema{
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
								"is_overridable": &schema.Schema{
									Type:     schema.TypeBool,
									Optional: true,
									Computed: true,
								},
								"sysprep": &schema.Schema{
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
							},
						},
					},
					"power_state_mechanism": &schema.Schema{
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"guest_transition_config": &schema.Schema{
									Type:     schema.TypeMap,
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
									Type:     schema.TypeMap,
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
												Type:     schema.TypeMap,
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
								"disk_size_mib": &schema.Schema{
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
								},
								"volume_group_reference": &schema.Schema{
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
				},
			},
		},
	}
}
