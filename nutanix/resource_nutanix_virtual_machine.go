package nutanix

import (
	"fmt"
	"log"
	"reflect"
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
	conn := meta.(*Client).API

	// Prepare request
	request := &v3.VMIntentInput{}
	spec := &v3.VM{}
	metadata := &v3.Metadata{}
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
		return fmt.Errorf("please provide the required name attribute")
	}
	if err := getMetadataAttributes(d, metadata, "vm"); err != nil {
		return err
	}
	if descok {
		spec.Description = utils.String(desc.(string))
	}
	if azrok {
		a := azr.(map[string]interface{})
		spec.AvailabilityZoneReference = validateRef(a)
	}
	if crok {
		a := cr.(map[string]interface{})
		spec.ClusterReference = validateRef(a)
	}

	if err := getVMResources(d, res); err != nil {
		return err
	}

	spec.Name = utils.String(n.(string))
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

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
		return fmt.Errorf("error waiting for vm (%s) to create: %s", d.Id(), err)
	}

	return resourceNutanixVirtualMachineRead(d, meta)
}

func resourceNutanixVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetVM(d.Id())
	if err != nil {
		return err
	}

	m, c := setRSEntityMetadata(resp.Metadata)
	n := setNicList(resp.Status.Resources.NicList)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", getReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", getReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_reference", getReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	if err := d.Set("cluster_reference", getClusterReferenceValues(resp.Status.ClusterReference)); err != nil {
		return err
	}
	if err := d.Set("nic_list", n); err != nil {
		return err
	}
	if err := d.Set("host_reference", getReferenceValues(resp.Status.Resources.HostReference)); err != nil {
		return err
	}
	if err := d.Set("nutanix_guest_tools", setNutanixGuestTools(resp.Status.Resources.GuestTools)); err != nil {
		return err
	}
	if err := d.Set("gpu_list", setGPUList(resp.Status.Resources.GpuList)); err != nil {
		return err
	}
	if err := d.Set("parent_reference", getReferenceValues(resp.Status.Resources.ParentReference)); err != nil {
		return err
	}

	diskAddress := make(map[string]interface{})
	mac := ""
	b := make([]string, 0)

	if resp.Status.Resources.BootConfig != nil {
		if resp.Status.Resources.BootConfig.BootDevice.DiskAddress != nil {
			i := strconv.Itoa(int(utils.Int64Value(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.DeviceIndex)))
			diskAddress["device_index"] = i
			diskAddress["adapter_type"] = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.DiskAddress.AdapterType)
		}
		if resp.Status.Resources.BootConfig.BootDeviceOrderList != nil {
			b = utils.StringValueSlice(resp.Status.Resources.BootConfig.BootDeviceOrderList)
		}
		mac = utils.StringValue(resp.Status.Resources.BootConfig.BootDevice.MacAddress)
	}

	d.Set("boot_device_order_list", b)
	d.Set("boot_device_disk_address", diskAddress)
	d.Set("boot_device_mac_address", mac)

	sysprep := make(map[string]interface{})
	sysprepCV := make(map[string]string)
	cloudInit := make(map[string]interface{})
	cloudInitCV := make(map[string]string)
	isOv := false
	if resp.Status.Resources.GuestCustomization != nil {
		isOv = utils.BoolValue(resp.Status.Resources.GuestCustomization.IsOverridable)
		if resp.Status.Resources.GuestCustomization.CloudInit != nil {
			cloudInit["meta_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.MetaData)
			cloudInit["user_data"] = utils.StringValue(resp.Status.Resources.GuestCustomization.CloudInit.UserData)
			if resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues != nil {
				for k, v := range resp.Status.Resources.GuestCustomization.CloudInit.CustomKeyValues {
					cloudInitCV[k] = v
				}
			}
		}
		if resp.Status.Resources.GuestCustomization.Sysprep != nil {
			sysprep["install_type"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.InstallType)
			sysprep["unattend_xml"] = utils.StringValue(resp.Status.Resources.GuestCustomization.Sysprep.UnattendXML)
			if resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues != nil {
				for k, v := range resp.Status.Resources.GuestCustomization.Sysprep.CustomKeyValues {
					sysprepCV[k] = v
				}
			}
		}
	}
	if err := d.Set("guest_customization_cloud_init_custom_key_values", cloudInitCV); err != nil {
		return err
	}
	if err := d.Set("guest_customization_sysprep_custom_key_values", sysprepCV); err != nil {
		return err
	}
	if err := d.Set("guest_customization_sysprep", sysprep); err != nil {
		return err
	}
	if err := d.Set("guest_customization_cloud_init", cloudInit); err != nil {
		return err
	}

	d.Set("hardware_clock_timezone", utils.StringValue(resp.Status.Resources.HardwareClockTimezone))
	d.Set("cluster_reference_name", utils.StringValue(resp.Status.ClusterReference.Name))
	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Status.Name))
	d.Set("description", utils.StringValue(resp.Status.Description))
	d.Set("state", utils.StringValue(resp.Status.State))
	d.Set("num_vnuma_nodes", utils.Int64Value(resp.Status.Resources.VnumaConfig.NumVnumaNodes))
	d.Set("guest_os_id", utils.StringValue(resp.Status.Resources.GuestOsID))
	d.Set("power_state", utils.StringValue(resp.Status.Resources.PowerState))
	d.Set("num_vcpus_per_socket", utils.Int64Value(resp.Status.Resources.NumVcpusPerSocket))
	d.Set("num_sockets", utils.Int64Value(resp.Status.Resources.NumSockets))
	d.Set("memory_size_mib", utils.Int64Value(resp.Status.Resources.MemorySizeMib))
	d.Set("guest_customization_is_overridable", isOv)
	d.Set("should_fail_on_script_failure", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.ShouldFailOnScriptFailure))
	d.Set("enable_script_exec", utils.BoolValue(resp.Status.Resources.PowerStateMechanism.GuestTransitionConfig.EnableScriptExec))
	d.Set("power_state_mechanism", utils.StringValue(resp.Status.Resources.PowerStateMechanism.Mechanism))
	d.Set("vga_console_enabled", utils.BoolValue(resp.Status.Resources.VgaConsoleEnabled))
	d.SetId(*resp.Metadata.UUID)

	return d.Set("disk_list", setDiskList(resp.Status.Resources.DiskList, resp.Status.Resources.GuestCustomization))
}

func resourceNutanixVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	log.Printf("Updating VM values %s", d.Id())
	fmt.Printf("Updating VM values %s", d.Id())

	request := &v3.VMIntentInput{}
	metadata := &v3.Metadata{}
	res := &v3.VMResources{}
	spec := &v3.VM{}
	guest := &v3.GuestCustomization{}
	guestTool := &v3.GuestToolsSpec{}
	boot := &v3.VMBootConfig{}
	pw := &v3.VMPowerStateMechanism{}

	response, err := conn.V3.GetVM(d.Id())
	preFillResUpdateRequest(res, response)
	preFillGTUpdateRequest(guestTool, response)
	preFillGUpdateRequest(guest, response)
	preFillPWUpdateRequest(pw, response)

	if err != nil {
		return err
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if d.HasChange("categories") {
		catl := d.Get("categories").([]interface{})

		if len(catl) > 0 && catl != nil {
			cl := make(map[string]string)
			for _, v := range catl {
				item := v.(map[string]interface{})

				if i, ok := item["name"]; ok && i.(string) != "" {
					if k, kok := item["value"]; kok && k.(string) != "" {
						cl[i.(string)] = k.(string)
					}
				}
			}
			metadata.Categories = cl
		} else {
			metadata.Categories = nil
		}
	}
	metadata.OwnerReference = response.Metadata.OwnerReference
	if d.HasChange("owner_reference") {
		_, n := d.GetChange("owner_reference")
		metadata.OwnerReference = validateRef(n.(map[string]interface{}))
	}
	metadata.ProjectReference = response.Metadata.ProjectReference
	if d.HasChange("project_reference") {
		_, n := d.GetChange("project_reference")
		metadata.ProjectReference = validateRef(n.(map[string]interface{}))
	}
	spec.Name = response.Status.Name
	if d.HasChange("name") {
		_, n := d.GetChange("name")
		spec.Name = utils.String(n.(string))
	}
	spec.Description = response.Status.Description
	if d.HasChange("description") {
		_, n := d.GetChange("description")
		spec.Description = utils.String(n.(string))
	}
	spec.AvailabilityZoneReference = response.Status.AvailabilityZoneReference
	if d.HasChange("availability_zone_reference") {
		_, n := d.GetChange("availability_zone_reference")
		spec.AvailabilityZoneReference = validateRef(n.(map[string]interface{}))
	}
	spec.ClusterReference = response.Status.ClusterReference
	if d.HasChange("cluster_reference") {
		_, n := d.GetChange("cluster_reference")
		spec.ClusterReference = validateRef(n.(map[string]interface{}))
	}
	if d.HasChange("parent_reference") {
		_, n := d.GetChange("parent_reference")
		res.ParentReference = validateRef(n.(map[string]interface{}))
	}
	if d.HasChange("num_vnuma_nodes") {
		_, n := d.GetChange("num_vnuma_nodes")
		res.VMVnumaConfig = &v3.VMVnumaConfig{
			NumVnumaNodes: utils.Int64(int64(n.(int))),
		}
	}
	if d.HasChange("guest_os_id") {
		_, n := d.GetChange("guest_os_id")
		res.GuestOsID = utils.String(n.(string))
	}
	if d.HasChange("power_state") {
		_, n := d.GetChange("power_state")

		res.PowerState = utils.String(n.(string))
	}
	if d.HasChange("num_vcpus_per_socket") {
		_, n := d.GetChange("num_vcpus_per_socket")
		res.NumVcpusPerSocket = utils.Int64(int64(n.(int)))
	}
	if d.HasChange("num_sockets") {
		_, n := d.GetChange("num_sockets")
		res.NumSockets = utils.Int64(int64(n.(int)))
	}
	if d.HasChange("memory_size_mib") {
		_, n := d.GetChange("memory_size_mib")
		res.MemorySizeMib = utils.Int64(int64(n.(int)))
	}
	if d.HasChange("hardware_clock_timezone") {
		_, n := d.GetChange("hardware_clock_timezone")
		res.HardwareClockTimezone = utils.String(n.(string))
	}
	if d.HasChange("vga_console_enabled") {
		_, n := d.GetChange("vga_console_enabled")
		res.VgaConsoleEnabled = utils.Bool(n.(bool))
	}
	if d.HasChange("guest_customization_is_overridable") {
		_, n := d.GetChange("guest_customization_is_overridable")
		guest.IsOverridable = utils.Bool(n.(bool))
	}
	if d.HasChange("power_state_mechanism") {
		_, n := d.GetChange("power_state_mechanism")
		pw.Mechanism = utils.String(n.(string))
	}
	if d.HasChange("power_state_guest_transition_config") {
		_, n := d.GetChange("power_state_guest_transition_config")
		val := n.(map[string]interface{})

		p := &v3.VMGuestPowerStateTransitionConfig{}
		if v, ok := val["enable_script_exec"]; ok {
			p.EnableScriptExec = utils.Bool(v.(bool))
		}
		if v, ok := val["should_fail_on_script_failure"]; ok {
			p.ShouldFailOnScriptFailure = utils.Bool(v.(bool))
		}
		pw.GuestTransitionConfig = p
	}
	if d.HasChange("guest_customization_cloud_init") {
		_, n := d.GetChange("guest_customization_cloud_init")
		a := n.(map[string]interface{})

		guest.CloudInit = &v3.GuestCustomizationCloudInit{
			MetaData: validateMapStringValue(a, "meta_data"),
			UserData: validateMapStringValue(a, "user_data"),
		}
	}
	if d.HasChange("guest_customization_cloud_init_custom_key_values") {
		if guest.CloudInit == nil {
			guest.CloudInit = &v3.GuestCustomizationCloudInit{}
		}
		_, n := d.GetChange("guest_customization_cloud_init_custom_key_values")
		guest.CloudInit.CustomKeyValues = n.(map[string]string)
	}
	if d.HasChange("guest_customization_sysprep") {
		_, n := d.GetChange("guest_customization_sysprep")
		a := n.(map[string]interface{})

		guest.Sysprep = &v3.GuestCustomizationSysprep{
			InstallType: validateMapStringValue(a, "install_type"),
			UnattendXML: validateMapStringValue(a, "unattend_xml"),
		}
	}
	if d.HasChange("guest_customization_sysprep_custom_key_values") {
		if guest.Sysprep == nil {
			guest.Sysprep = &v3.GuestCustomizationSysprep{}
		}
		_, n := d.GetChange("guest_customization_sysprep_custom_key_values")
		guest.Sysprep.CustomKeyValues = n.(map[string]string)
	}
	if d.HasChange("nic_list") {
		_, ne := d.GetChange("nic_list")
		n := ne.([]interface{})
		if len(n) > 0 && n != nil {
			nics := make([]*v3.VMNic, len(n))

			for k, nc := range n {
				val := nc.(map[string]interface{})
				net := val["network_function_chain_reference"].(map[string]interface{})
				sub := val["subnet_reference"].(map[string]interface{})

				nic := &v3.VMNic{
					NicType: validateMapStringValue(val, "nic_type"),
					UUID:    validateMapStringValue(val, "uuid"),
					NetworkFunctionNicType: validateMapStringValue(val, "network_function_nic_type"),
					MacAddress:             validateMapStringValue(val, "mac_address"),
					Model:                  validateMapStringValue(val, "model"),
				}

				nic.NetworkFunctionChainReference = validateRef(net)
				nic.SubnetReference = validateRef(sub)

				if value, ok := val["ip_endpoint_list"]; ok && len(value.([]interface{})) > 0 && value.([]interface{}) != nil {
					ip := make([]*v3.IPAddress, len(value.([]interface{})))
					for k, i := range value.([]interface{}) {
						v := i.(map[string]interface{})
						ip[k] = &v3.IPAddress{
							IP:   validateMapStringValue(v, "ip"),
							Type: validateMapStringValue(v, "type"),
						}
					}
					nic.IPEndpointList = ip
				}

				nics[k] = nic
			}
			res.NicList = nics
		}
	}
	if d.HasChange("nutanix_guest_tools") {
		_, n := d.GetChange("nutanix_guest_tools")
		ngt := n.(map[string]interface{})

		tool := &v3.NutanixGuestToolsSpec{}
		tool.IsoMountState = validateMapStringValue(ngt, "iso_mount_state")
		tool.State = validateMapStringValue(ngt, "state")

		if val, ok2 := ngt["enabled_capability_list"]; ok2 && val.([]interface{}) != nil {
			var l []*string
			for _, list := range val.([]interface{}) {
				l = append(l, utils.String(list.(string)))
			}
			tool.EnabledCapabilityList = l
		}
		guestTool.NutanixGuestTools = tool
	}
	if d.HasChange("gpu_list") {
		_, n := d.GetChange("gpu_list")

		if v, ok := n.([]interface{}); ok && v != nil {
			gpl := make([]*v3.VMGpu, len(v))

			for k, va := range v {
				val := va.(map[string]interface{})
				gpl[k] = &v3.VMGpu{
					Vendor:   validateMapStringValue(val, "vendor"),
					Mode:     validateMapStringValue(val, "mode"),
					DeviceID: validateMapIntValue(val, "device_id"),
				}
			}
			res.GpuList = gpl
		}
	}
	if d.HasChange("boot_device_order_list") {
		_, n := d.GetChange("boot_device_order_list")
		var b []*string
		for _, boot := range n.([]interface{}) {
			b = append(b, utils.String(boot.(string)))
		}
		boot.BootDeviceOrderList = b
	}

	bd := &v3.VMBootDevice{}
	dska := &v3.DiskAddress{}
	if d.HasChange("boot_device_disk_address") {
		_, n := d.GetChange("boot_device_disk_address")
		dai := n.(map[string]interface{})
		dska = &v3.DiskAddress{
			DeviceIndex: validateMapIntValue(dai, "device_index"),
			AdapterType: validateMapStringValue(dai, "adapter_type"),
		}
	}
	if d.HasChange("boot_device_mac_address") {
		_, n := d.GetChange("boot_device_mac_address")
		bd.MacAddress = utils.String(n.(string))
	}
	if d.HasChange("disk_list") {
		_, n := d.GetChange("disk_list")
		if dsk, ok := n.([]interface{}); ok && len(dsk) > 0 && dsk != nil {
			dls := make([]*v3.VMDisk, len(dsk))

			for k, val := range dsk {

				if response.Status.Resources.GuestCustomization.CloudInit != nil && val.(map[string]interface{})["device_properties"].([]interface{})[0].(map[string]interface{})["device_type"].(string) == "CDROM" {
					continue
				}

				v := val.(map[string]interface{})
				dl := &v3.VMDisk{
					UUID:          validateMapStringValue(v, "uuid"),
					DiskSizeBytes: validateMapIntValue(v, "disk_size_bytes"),
					DiskSizeMib:   validateMapIntValue(v, "disk_size_mib"),
				}
				if v1, ok1 := v["device_properties"]; ok1 && v1.([]interface{}) != nil && len(v1.([]interface{})) > 0 {
					dvp := v1.([]interface{})
					d := dvp[0].(map[string]interface{})
					dp := &v3.VMDiskDeviceProperties{
						DeviceType: validateMapStringValue(d, "device_type"),
					}
					if v2, ok := d["disk_address"]; ok && len(v2.([]interface{})) > 0 && v2.([]interface{})[0] != nil {
						da := v2.([]interface{})[0].(map[string]interface{})
						dp.DiskAddress = &v3.DiskAddress{
							AdapterType: validateMapStringValue(da, "adapter_type"),
							DeviceIndex: validateMapIntValue(da, "device_index"),
						}
					}
					dl.DeviceProperties = dp
				}
				if v1, ok := v["data_source_reference"]; ok {
					dsref := v1.([]interface{})
					if len(dsref) > 0 && dsref[0] != nil {
						dsri := dsref[0].(map[string]interface{})
						dl.DataSourceReference = validateRef(dsri)
					}
				}
				if v1, ok := v["volume_group_reference"]; ok {
					volgr := v1.([]interface{})
					if len(volgr) > 0 && volgr[0] != nil {
						dsri := volgr[0].(map[string]interface{})
						dl.VolumeGroupReference = validateRef(dsri)
					}
				}
				dls[k] = dl
			}
			res.DiskList = dls
		}
	}
	boot.BootDevice = bd

	if dska.AdapterType == nil && dska.DeviceIndex == nil && bd.MacAddress == nil {
		boot = nil
	}

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

	_, err2 := conn.V3.UpdateVM(d.Id(), request)
	if err2 != nil {
		return err2
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
			"error waiting for vm (%s) to update: %s", d.Id(), err)
	}

	return resourceNutanixVirtualMachineRead(d, meta)
}

func resourceNutanixVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

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
			"error waiting for vm (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixVirtualMachineExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*Client).API

	getEntitiesRequest := &v3.DSMetadata{}
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
					nic.NetworkFunctionChainReference = validateRef(v)
				}
				if value, ok := val["subnet_reference"]; ok {
					v := value.(map[string]interface{})
					nic.SubnetReference = validateRef(v)
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
		val := v.(map[string]interface{})
		vm.ParentReference = validateRef(val)
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

	guestCustom := &v3.GuestCustomization{}

	if v, ok := d.GetOk("guest_customization_cloud_init"); ok {
		guestCustom.CloudInit = &v3.GuestCustomizationCloudInit{}
		cii := v.(map[string]interface{})
		if v2, ok2 := cii["meta_data"]; ok2 {
			guestCustom.CloudInit.MetaData = utils.String(v2.(string))
		}
		if v2, ok2 := cii["user_data"]; ok2 {
			guestCustom.CloudInit.UserData = utils.String(v2.(string))
		} else if *guestCustom.CloudInit.UserData != "" {
			guestCustom.CloudInit.MetaData = utils.String("")
		}
	}

	if v, ok := d.GetOk("guest_customization_cloud_init_custom_key_values"); ok {
		if guestCustom.CloudInit == nil {
			guestCustom.CloudInit = &v3.GuestCustomizationCloudInit{}
		}
		guestCustom.CloudInit.CustomKeyValues = v.(map[string]string)
	}
	if v, ok := d.GetOk("guest_customization_is_overridable"); ok {
		guestCustom.IsOverridable = utils.Bool(v.(bool))
	}
	if v, ok := d.GetOk("guest_customization_sysprep"); ok {
		guestCustom.Sysprep = &v3.GuestCustomizationSysprep{}
		spi := v.(map[string]interface{})
		if v2, ok2 := spi["install_type"]; ok2 {
			guestCustom.Sysprep.InstallType = utils.String(v2.(string))
		}
		if v2, ok2 := spi["unattend_xml"]; ok2 {
			guestCustom.Sysprep.UnattendXML = utils.String(v2.(string))
		}
	}

	if v, ok := d.GetOk("guest_customization_sysprep_custom_key_values"); ok {
		if guestCustom.CloudInit == nil {
			guestCustom.CloudInit = &v3.GuestCustomizationCloudInit{}
		}
		guestCustom.Sysprep.CustomKeyValues = v.(map[string]string)
	}

	if !reflect.DeepEqual(*guestCustom, (v3.GuestCustomization{})) {
		vm.GuestCustomization = guestCustom
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
						if v1, ok := d["device_type"]; ok {
							dp.DeviceType = utils.String(v1.(string))
						}
						if v2, ok := d["disk_address"]; ok {
							if len(v2.([]interface{})) > 0 {
								da := v2.([]interface{})[0].(map[string]interface{})
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
						dl.DataSourceReference = validateRef(dsri)
					}
				}
				if v1, ok := v["volume_group_reference"]; ok {
					volgr := v1.([]interface{})
					if len(volgr) > 0 {
						dsri := volgr[0].(map[string]interface{})
						dl.VolumeGroupReference = validateRef(dsri)
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

func preFillResUpdateRequest(res *v3.VMResources, response *v3.VMIntentResponse) {
	res.ParentReference = response.Status.Resources.ParentReference
	res.VMVnumaConfig = &v3.VMVnumaConfig{NumVnumaNodes: response.Status.Resources.VnumaConfig.NumVnumaNodes}
	res.GuestOsID = response.Status.Resources.GuestOsID
	res.PowerState = response.Status.Resources.PowerState
	res.NumVcpusPerSocket = response.Status.Resources.NumVcpusPerSocket
	res.NumSockets = response.Status.Resources.NumSockets
	res.MemorySizeMib = response.Status.Resources.MemorySizeMib
	res.HardwareClockTimezone = response.Status.Resources.HardwareClockTimezone
	res.VgaConsoleEnabled = response.Status.Resources.VgaConsoleEnabled

	nold := make([]*v3.VMNic, len(response.Status.Resources.NicList))
	if len(response.Status.Resources.NicList) > 0 {
		for k, v := range response.Status.Resources.NicList {
			nold[k] = &v3.VMNic{
				NicType: v.NicType,
				UUID:    v.UUID,
				NetworkFunctionNicType: v.NetworkFunctionNicType,
				MacAddress:             v.MacAddress,
				Model:                  v.Model,
				NetworkFunctionChainReference: v.NetworkFunctionChainReference,
				SubnetReference:               v.SubnetReference,
				IPEndpointList:                v.IPEndpointList,
			}
		}
	} else {
		nold = nil
	}
	res.NicList = nold

	gold := make([]*v3.VMGpu, len(response.Status.Resources.GpuList))
	if len(response.Status.Resources.GpuList) > 0 {
		for k, v := range response.Status.Resources.GpuList {
			gold[k] = &v3.VMGpu{
				Vendor:   v.Vendor,
				Mode:     v.Mode,
				DeviceID: v.DeviceID,
			}
		}
	} else {
		gold = nil
	}
	res.GpuList = gold
	if response.Status.Resources.BootConfig != nil {
		res.BootConfig = response.Status.Resources.BootConfig
	} else {
		res.BootConfig = nil
	}
}

func preFillGTUpdateRequest(guestTool *v3.GuestToolsSpec, response *v3.VMIntentResponse) {
	if response.Status.Resources.GuestTools != nil {
		guestTool.NutanixGuestTools = &v3.NutanixGuestToolsSpec{
			EnabledCapabilityList: response.Status.Resources.GuestTools.NutanixGuestTools.EnabledCapabilityList,
			IsoMountState:         response.Status.Resources.GuestTools.NutanixGuestTools.IsoMountState,
			State:                 response.Status.Resources.GuestTools.NutanixGuestTools.State,
		}
	} else {
		guestTool = nil
	}
}

func preFillGUpdateRequest(guest *v3.GuestCustomization, response *v3.VMIntentResponse) {
	if response.Status.Resources.GuestCustomization != nil {
		guest.CloudInit = response.Status.Resources.GuestCustomization.CloudInit
		guest.Sysprep = response.Status.Resources.GuestCustomization.Sysprep
		guest.IsOverridable = response.Status.Resources.GuestCustomization.IsOverridable
	} else {
		guest = nil
	}
}

func preFillPWUpdateRequest(pw *v3.VMPowerStateMechanism, response *v3.VMIntentResponse) {
	if response.Status.Resources.PowerStateMechanism != nil {
		pw.Mechanism = response.Status.Resources.PowerStateMechanism.Mechanism
		pw.GuestTransitionConfig = response.Status.Resources.PowerStateMechanism.GuestTransitionConfig
	} else {
		pw = nil
	}
}

func getVMSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metadata": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"last_update_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_time": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"spec_hash": {
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
		"categories": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"project_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"owner_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"api_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"availability_zone_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"cluster_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"host_reference": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"hypervisor_type": {
			Type:     schema.TypeString,
			Computed: true,
		},

		// RESOURCES ARGUMENTS

		"num_vnuma_nodes": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"nic_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nic_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"floating_ip": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"model": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"network_function_nic_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"mac_address": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"ip_endpoint_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"ip": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"type": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"network_function_chain_reference": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"name": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"subnet_reference": {
						Type:     schema.TypeMap,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Required: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"subnet_reference_name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"guest_os_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"power_state": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nutanix_guest_tools": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"available_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"iso_mount_state": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"guest_os_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"enabled_capability_list": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"vss_snapshot_capable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"is_reachable": {
						Type:     schema.TypeBool,
						Computed: true,
					},
					"vm_mobility_drivers_installed": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"num_vcpus_per_socket": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"num_sockets": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"gpu_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"frame_buffer_size_mib": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"vendor": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"pci_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"fraction": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"mode": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"num_virtual_display_heads": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"guest_driver_version": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"device_id": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"parent_reference": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kind": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"uuid": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"memory_size_mib": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"boot_device_order_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"boot_device_disk_address": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"device_index": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"adapter_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"boot_device_mac_address": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"hardware_clock_timezone": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"guest_customization_cloud_init": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"meta_data": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"user_data": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"guest_customization_cloud_init_custom_key_values": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"guest_customization_is_overridable": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"guest_customization_sysprep": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"install_type": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"unattend_xml": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"guest_customization_sysprep_custom_key_values": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"should_fail_on_script_failure": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"enable_script_exec": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"power_state_mechanism": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"vga_console_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"disk_list": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"uuid": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"disk_size_bytes": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"disk_size_mib": {
						Type:     schema.TypeInt,
						Optional: true,
						Computed: true,
					},
					"device_properties": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"device_type": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"disk_address": {
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"device_index": {
												Type:     schema.TypeInt,
												Optional: true,
												Computed: true,
											},
											"adapter_type": {
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
					"data_source_reference": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
							},
						},
					},

					"volume_group_reference": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"kind": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Optional: true,
									Computed: true,
								},
								"uuid": {
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
	}
}
