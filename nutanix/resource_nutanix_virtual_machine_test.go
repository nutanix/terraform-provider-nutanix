package nutanix

// import (
// 	"bufio"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	flag "terraform-provider-nutanix/testflg"
// 	"testing"

// 	"github.com/hashicorp/terraform/helper/resource"
// 	"github.com/hashicorp/terraform/terraform"
// )

// // Base setup function to check that api_version is set
// func testBasicPreCheck(t *testing.T) {
// 	testAccPreCheck(t)

// 	if flag.NutanixName == "" {
// 		t.Fatal("env variable NUTANIX_NAME must be set for acceptance tests")
// 	}
// }

// type TemplateBasicBodyVars struct {
// 	name                string
// 	numVCPUs            string
// 	numSockets          string
// 	memorySizeMb        string
// 	powerState          string
// 	nicType             string
// 	nicKind             string
// 	nicUUID             string
// 	networkFunctionType string
// 	project             string
// }

// func (body TemplateBasicBodyVars) testSprintfTemplateBodyWithoutNic(template string) string {
// 	return fmt.Sprintf(
// 		template,
// 		body.name,
// 		body.numSockets,
// 		body.numVCPUs,
// 		body.memorySizeMb,
// 		body.powerState,
// 		body.project,
// 	)
// }

// func (body TemplateBasicBodyVars) testSprintfTemplateBody(template string) string {
// 	return fmt.Sprintf(
// 		template,
// 		body.name,
// 		body.numSockets,
// 		body.numVCPUs,
// 		body.memorySizeMb,
// 		body.powerState,
// 		body.nicType,
// 		body.nicKind,
// 		body.nicUUID,
// 		body.networkFunctionType,
// 		body.project,
// 	)
// }

// func (body TemplateBasicBodyVars) testSprintfTemplateBodyUpdateName(template string) string {
// 	return fmt.Sprintf(
// 		template,
// 		flag.NutanixUpdateName,
// 		body.numSockets,
// 		body.numVCPUs,
// 		body.memorySizeMb,
// 		body.powerState,
// 		body.nicType,
// 		body.nicKind,
// 		body.nicUUID,
// 		body.networkFunctionType,
// 		body.project,
// 	)
// }

// // setups variables used by fixed ip tests
// func setupTemplateBasicBodyVars() TemplateBasicBodyVars {
// 	data := TemplateBasicBodyVars{
// 		name:                flag.NutanixName,
// 		numSockets:          flag.NutanixNumSockets,
// 		numVCPUs:            flag.NutanixNumVCPUs,
// 		memorySizeMb:        flag.NutanixMemorySize,
// 		powerState:          flag.NutanixPowerState,
// 		nicType:             flag.NutanixNicType,
// 		nicKind:             flag.NutanixNicKind,
// 		nicUUID:             flag.NutanixNicUUID,
// 		networkFunctionType: flag.NutanixNetworkFunctionType,
// 		project:             flag.NutanixProject,
// 	}
// 	return data
// }

// // Basic data to create series of testing functions
// type TestFuncData struct {
// 	vm           nutanixV3.VmIntentInput
// 	vmName       string
// 	name         string
// 	numVCPUs     string
// 	numSockets   string
// 	memorySizeMb string
// 	powerState   string
// }

// // returns TestCheckFunc's that will be used in most of our tests
// // numVCPUs, numSockets defaults to 1
// // APIVersion defaults to 3.0 specVersion 0 and memorySizeMb tp 1024
// // kind defaults to "vm" and powerState to "ON", vmName to "nutanix_virtual_machine"
// func (test TestFuncData) testCheckFuncBasic() (resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc) {
// 	vmName := test.vmName
// 	if vmName == "" {
// 		vmName = "nutanix_virtual_machine.my-machine"
// 	}
// 	name := test.name
// 	if name == "" {
// 		name = flag.NutanixName
// 	}
// 	memorySize := test.memorySizeMb
// 	if memorySize == "" {
// 		memorySize = flag.NutanixMemorySize
// 	}

// 	return testAccCheckNutanixVirtualMachineExists(vmName, &test.vm),
// 		resource.TestCheckResourceAttr(vmName, "spec.#", "1"),
// 		resource.TestCheckResourceAttr(vmName, "spec.0.resources.#", "1"),
// 		resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.power_state", flag.NutanixPowerState),
// 		resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.memory_size_mb", memorySize),
// 		resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.num_sockets", flag.NutanixNumSockets),
// 		resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.num_vcpus_per_socket", flag.NutanixNumVCPUs),
// 		resource.TestCheckResourceAttr(vmName, "name", name)
// }

// const testAccCheckNutanixVirtualMachineConfigReallyBasic = `
// resource "nutanix_virtual_machine" "my-machine" {
// 	name = "%s"
// ` + testAccTemplateBasicBodyWithEnd

// const testAccCheckNutanixVirtualMachineConfigMostBasic = `
// resource "nutanix_virtual_machine" "my-machine" {
// 	name = "%s"
// ` + testAccTemplateMostBasicBody + `
// }`

// const testAccTemplateSpecBody = `
// spec = {
// `
// const testAccTemplateResourcesBody = testAccTemplateBasicResourcesBody + nicListBody
// const testAccTemplateBasicResourcesBody = `
// 		resources = {
// 			num_sockets = %s
// 			num_vcpus_per_socket = %s
// 			memory_size_mb = %s
// 			power_state = "%s"
// `
// const nicListBody = `
// 			nic_list = [
// 				{
// 					nic_type = "%s"
// 					subnet_reference = {
// 						kind = "%s"
// 						uuid = "%s"
// 					}
// 					network_function_nic_type = "%s"
// 				}
// 			]
// `
// const testAccTemplateMetadata = `
// 	metadata = {
// 		categories = {
// 			"Project" = "%s"
// 		}
// `
// const testAccTemplateBasicBody = testAccTemplateSpecBody +
// 	testAccTemplateResourcesBody + `
// 		}
// 	}
// ` +
// 	testAccTemplateMetadata + `
// 	}
// `
// const testAccTemplateMostBasicBody = testAccTemplateSpecBody +
// 	testAccTemplateBasicResourcesBody + `
// 		}
// 	}
// ` +
// 	testAccTemplateMetadata + `
// 	}
// `
// const testAccTemplateBasicBodyWithEnd = testAccTemplateBasicBody + `
// }`

// func TestAccNutanixVirtualMachine_basic(t *testing.T) {
// 	var vm nutanixV3.VmIntentInput
// 	basicVars := setupTemplateBasicBodyVars()
// 	config := basicVars.testSprintfTemplateBodyWithoutNic(testAccCheckNutanixVirtualMachineConfigMostBasic)

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testBasicPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: config,
// 				Check: resource.ComposeTestCheckFunc(
// 					TestFuncData{vm: vm, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
// 				),
// 			},
// 		},
// 	})
// }

// // testing vms with nic_list config
// func TestAccNutanixVirtualMachine_nicList(t *testing.T) {
// 	var vm nutanixV3.VmIntentInput
// 	basicVars := setupTemplateBasicBodyVars()
// 	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)

// 	log.Printf("[DEBUG] template config= %s", config)

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testBasicPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: config,
// 				Check: resource.ComposeTestCheckFunc(
// 					TestFuncData{vm: vm, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
// 				),
// 			},
// 		},
// 	})
// }

// const diskTemplate = `
// 				{
// 					data_source_reference = {
// 						kind = "%s"
// 						name = "%s"
// 						uuid = "%s"
// 					}
// 					device_properties = {
// 						device_type = "%s"
// 					}
// 					disk_size_mib = %s
// 				}`

// func diskSet() string {
// 	diskList := `
// 			disk_list = [
// 	`
// 	diskNo, _ := strconv.Atoi(flag.NutanixDiskNo)
// 	for i := 0; i < diskNo-1; i++ {
// 		disk := fmt.Sprintf(diskTemplate, flag.NutanixDiskKind[i], flag.NutanixDiskName[i], flag.NutanixDiskUUID[i], flag.NutanixDiskDeviceType[i], flag.NutanixDiskSize[i]) + ","
// 		diskList = diskList + disk
// 	}
// 	i := diskNo - 1
// 	if diskNo > 0 {
// 		disk := fmt.Sprintf(diskTemplate, flag.NutanixDiskKind[i], flag.NutanixDiskName[i], flag.NutanixDiskUUID[i], flag.NutanixDiskDeviceType[i], flag.NutanixDiskSize[i])
// 		diskList = diskList + disk
// 	}
// 	diskList = diskList + `
// 			]
// 	`
// 	return diskList
// }

// // testing vms with disk list
// func TestAccNutanixVirtualMachine_diskList(t *testing.T) {
// 	var vm nutanixV3.VmIntentInput
// 	basicVars := setupTemplateBasicBodyVars()
// 	diskList := diskSet()
// 	vmName := "nutanix_virtual_machine.my-machine"
// 	testAccTemplateDiskBody := testAccTemplateSpecBody +
// 		testAccTemplateResourcesBody + diskList + `
// 		}
// 	}
// 	` +
// 		testAccTemplateMetadata + `
// 	}
// 	`
// 	testAccCheckNutanixVirtualMachineConfigDisk := `
// resource "nutanix_virtual_machine" "my-machine" {
// 	name = "%s"
// ` + testAccTemplateDiskBody + `
// }`

// 	basicVars.powerState = powerOFF
// 	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigDisk)
// 	log.Printf("[DEBUG] template config= %s", config)

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testBasicPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: config,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.#", "2"),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.0.data_source_reference.#", "1"),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.0.data_source_reference.0.kind", flag.NutanixDiskKind[0]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.0.data_source_reference.0.name", flag.NutanixDiskName[0]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.0.data_source_reference.0.uuid", flag.NutanixDiskUUID[0]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.0.device_properties.#", "1"),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.0.device_properties.0.device_type", flag.NutanixDiskDeviceType[0]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.0.disk_size_mib", flag.NutanixDiskSize[0]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.1.data_source_reference.#", "1"),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.1.data_source_reference.0.kind", flag.NutanixDiskKind[1]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.1.data_source_reference.0.name", flag.NutanixDiskName[1]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.1.data_source_reference.0.uuid", flag.NutanixDiskUUID[1]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.1.device_properties.#", "1"),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.1.device_properties.0.device_type", flag.NutanixDiskDeviceType[1]),
// 					resource.TestCheckResourceAttr(vmName, "spec.0.resources.0.disk_list.1.disk_size_mib", flag.NutanixDiskSize[1]),
// 				),
// 			},
// 		},
// 	})
// }

// // testing update memory in vm
// func TestAccNutanixVirtualMachine_updateMemory(t *testing.T) {
// 	var vm nutanixV3.VmIntentInput
// 	basicVars := setupTemplateBasicBodyVars()
// 	basicVars.powerState = powerOFF
// 	configOFF := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)
// 	basicVars.memorySizeMb = flag.NutanixUpdateMemorySize
// 	basicVars.powerState = flag.NutanixPowerState
// 	configON := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)
// 	log.Printf("[DEBUG] template configOFF= %s", configOFF)
// 	log.Printf("[DEBUG] template configON= %s", configON)

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testBasicPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: configOFF,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
// 				),
// 			},
// 			resource.TestStep{
// 				Config: configON,
// 				Check: resource.ComposeTestCheckFunc(
// 					TestFuncData{vm: vm, memorySizeMb: flag.NutanixUpdateMemorySize, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
// 				),
// 			},
// 		},
// 	})
// }

// // testing update name of the vm
// func TestAccNutanixVirtualMachine_updateName(t *testing.T) {
// 	var vm nutanixV3.VmIntentInput
// 	basicVars := setupTemplateBasicBodyVars()
// 	basicVars.powerState = powerOFF
// 	configOFF := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)
// 	basicVars.powerState = flag.NutanixPowerState
// 	configON := basicVars.testSprintfTemplateBodyUpdateName(testAccCheckNutanixVirtualMachineConfigReallyBasic)
// 	log.Printf("[DEBUG] template config= %s", configOFF)

// 	log.Printf("[DEBUG] template configON= %s", configON)

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testBasicPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: configOFF,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
// 				),
// 			},
// 			resource.TestStep{
// 				Config: configON,
// 				Check: resource.ComposeTestCheckFunc(
// 					TestFuncData{vm: vm, name: flag.NutanixUpdateName, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {

// 	for i := range s.RootModule().Resources {
// 		if s.RootModule().Resources[i].Type != "nutanix_virtual_machine" {
// 			continue
// 		}
// 		id := string(s.RootModule().Resources[i].Primary.ID)
// 		if id == "" {
// 			err := errors.New("ID is already set to the null string")
// 			return err
// 		}
// 		return nil
// 	}
// 	return nil
// }

// func testAccCheckNutanixVirtualMachineExists(n string, vm *nutanixV3.VmIntentInput) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		f, _ := os.Create("/tmp/check")
// 		w := bufio.NewWriter(f)
// 		defer f.Close()
// 		defer w.Flush()
// 		fmt.Fprintf(w, "%+v \n %+v\n %+v", s.Serial, s, *s)
// 		if n == "" {
// 			return fmt.Errorf("No vm name passed in")
// 		}
// 		if vm == nil {
// 			return fmt.Errorf("No vm obj passed in")
// 		}
// 		rs, ok := s.RootModule().Resources[n]
// 		if !ok {
// 			return fmt.Errorf("Not found: %s", n)
// 		}
// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("No ID is set")
// 		}
// 		ipAddress := string(rs.Primary.Attributes["ip_address"])
// 		if ipAddress == "" {
// 			fmt.Errorf("ip_address is not defined")
// 		}

// 		return nil
// 	}
// }
