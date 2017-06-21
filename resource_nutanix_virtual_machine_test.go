package nutanix

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	vmdefn "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachine"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Base setup function to check that api_version is set
func testBasicPreCheck(t *testing.T) {
	testAccPreCheck(t)

	if NutanixAPIVersion == "" {
		t.Fatal("env variable NUTANIX_API_VERSION must be set for acceptance tests")
	}
}

var specKey string
var specResourcesKey string
var metadataKey string
var diskSourceReference0Key string
var diskSourceReference1Key string
var deviceProperties0Key string
var deviceProperties1Key string

type TemplateBasicBodyVars struct {
	name         string
	numVCPUs     string
	numSockets   string
	memorySizeMb string
	powerState   string
	kind         string
	specVersion  string
	APIVersion   string
}

func (body TemplateBasicBodyVars) testSprintfTemplateBody(template string) string {
	return fmt.Sprintf(
		template,
		body.name,
		body.name,
		body.numSockets,
		body.numVCPUs,
		body.memorySizeMb,
		body.powerState,
		body.APIVersion,
		body.kind,
		body.specVersion,
		body.name,
	)
}

func (body TemplateBasicBodyVars) testSprintfTemplateBodyUpdateName(template string) string {
	return fmt.Sprintf(
		template,
		NutanixUpdateName,
		body.name,
		body.numSockets,
		body.numVCPUs,
		body.memorySizeMb,
		body.powerState,
		body.APIVersion,
		body.kind,
		body.specVersion,
		body.name,
	)
}

// setups variables used by fixed ip tests
func setupTemplateBasicBodyVars() TemplateBasicBodyVars {
	data := TemplateBasicBodyVars{
		name:         NutanixName,
		numSockets:   NutanixNumSockets,
		numVCPUs:     NutanixNumVCPUs,
		memorySizeMb: NutanixMemorySize,
		powerState:   NutanixPowerState,
		kind:         NutanixKind,
		specVersion:  NutanixSpecVersion,
		APIVersion:   NutanixAPIVersion,
	}
	return data
}

// Basic data to create series of testing functions
type TestFuncData struct {
	vm           vmdefn.VirtualMachine
	vmName       string
	name         string
	numVCPUs     string
	numSockets   string
	memorySizeMb string
	powerState   string
	APIVersion   string
	kind         string
	specVersion  string
}

func hashmapKey(s, t string) string {
	words := strings.Fields(terraformState)
	prefix := s + "."
	suffix := "." + t
	for i := range words {
		if (words[i] == strings.TrimPrefix(words[i], prefix+"#")) && (words[i] != strings.TrimPrefix(words[i], prefix)) {
			str1 := strings.TrimPrefix(words[i], prefix)
			str2 := strings.TrimSuffix(str1, suffix)
			str3 := strings.TrimSuffix(str1, suffix+".#")
			if str2 != str1 {
				return str2
			} else if str3 != str1 {
				return str3
			}
		}
	}
	return ""
}

// returns TestCheckFunc's that will be used in most of our tests
// numVCPUs, numSockets defaults to 1
// APIVersion defaults to 3.0 specVersion 0 and memorySizeMb tp 1024
// kind defaults to "vm" and powerState to "POWERED_ON", vmName to "nutanix_virtual_machine"
func (test TestFuncData) testCheckFuncBasic() (resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc) {
	vmName := test.vmName
	if vmName == "" {
		vmName = "nutanix_virtual_machine.my-machine"
	}
	name := test.name
	if name == "" {
		name = NutanixName
	}
	return testAccCheckNutanixVirtualMachineExists(vmName, &test.vm),
		resource.TestCheckResourceAttr(vmName, "api_version", NutanixAPIVersion),
		resource.TestCheckResourceAttr(vmName, "spec.#", "1"),
		resource.TestCheckResourceAttr(vmName, specKey+".resources.#", "1"),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".power_state", NutanixPowerState),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".memory_size_mb", NutanixMemorySize),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".num_sockets", NutanixNumSockets),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".num_vcpus_per_socket", NutanixNumVCPUs),
		resource.TestCheckResourceAttr(vmName, "metadata.#", "1"),
		resource.TestCheckResourceAttr(vmName, metadataKey+".kind", NutanixKind),
		resource.TestCheckResourceAttr(vmName, metadataKey+".spec_version", NutanixSpecVersion),
		resource.TestCheckResourceAttr(vmName, "name", name)

}

const testAccCheckNutanixVirtualMachineConfigReallyBasic = `
resource "nutanix_virtual_machine" "my-machine" {
	name = "%s"
` + testAccTemplateBasicBodyWithEnd

const testAccCheckNutanixVirtualMachineConfigMostBasic = `
resource "nutanix_virtual_machine" "my-machine" {
	name = "%s"
` + testAccTemplateMostBasicBody + `
}`

const testAccTemplateSpecBody = `
spec = {
	name = "%s"
`
const testAccTemplateResourcesBody = testAccTemplateBasicResourcesBody + nicListBody
const testAccTemplateBasicResourcesBody = `
		resources = {
			num_sockets = %s
			num_vcpus_per_socket = %s
			memory_size_mb = %s
			power_state = "%s"
`
const nicListBody = `
			nic_list = [
				{
					nic_type = "NORMAL_NIC"
					subnet_reference = {
						kind = "subnet"
						uuid = "c03ecf8f-aa1c-4a07-af43-9f2f198713c0"
					}
					network_function_nic_type = "INGRESS"
				}
			]
`
const testAccTemplateMetadata = `
	metadata = {
		kind = "%s"
		spec_version = %s
		name = "%s"
		categories = {
			"Project" = "nucalm"
		}
`
const testAccTemplateBasicBody = testAccTemplateSpecBody +
	testAccTemplateResourcesBody + `
		}
	}
	api_version = "%s"
` +
	testAccTemplateMetadata + `
	}
`
const testAccTemplateMostBasicBody = testAccTemplateSpecBody +
	testAccTemplateBasicResourcesBody + `
		}
	}
	api_version = "%s"
` +
	testAccTemplateMetadata + `
	}
`
const testAccTemplateBasicBodyWithEnd = testAccTemplateBasicBody + `
}`

// testing vms with basic config
func TestAccNutanixVirtualMachine_basic1(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigMostBasic)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
				),
			},
		},
	})
}

// testing vms with basic config
func TestAccNutanixVirtualMachine_basic2(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigMostBasic)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					TestFuncData{vm: vm, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
				),
			},
		},
	})
}

// testing vms with nic_list config
func TestAccNutanixVirtualMachine_nicList1(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
				),
			},
		},
	})
}

// testing vms with nic_list config
func TestAccNutanixVirtualMachine_nicList2(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)

	log.Printf("[DEBUG] template config= %s", config)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					TestFuncData{vm: vm, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
				),
			},
		},
	})
}

const diskTemplate = `
				{
					data_source_reference = {
						kind = "%s"
						name = "%s"
						uuid = "%s"
					}
					device_properties = {
						device_type = "%s"
					}
					disk_size_mib = %s
				}`

func diskSet() string {
	diskList := `
			disk_list = [
	`
	diskNo, _ := strconv.Atoi(NutanixDiskNo)
	for i := 0; i < diskNo-1; i++ {
		disk := fmt.Sprintf(diskTemplate, NutanixDiskKind[i], NutanixDiskName[i], NutanixDiskUUID[i], NutanixDiskDeviceType[i], NutanixDiskSize[i]) + ","
		diskList = diskList + disk
	}
	i := diskNo - 1
	if diskNo > 0 {
		disk := fmt.Sprintf(diskTemplate, NutanixDiskKind[i], NutanixDiskName[i], NutanixDiskUUID[i], NutanixDiskDeviceType[i], NutanixDiskSize[i])
		diskList = diskList + disk
	}
	diskList = diskList + `
			]
	`
	return diskList
}

// testing vms with disk list
func TestAccNutanixVirtualMachine_diskList1(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	diskList := diskSet()
	testAccTemplateDiskBody := testAccTemplateSpecBody +
		testAccTemplateResourcesBody + diskList + `
		}
	}
	api_version = "%s"
	` +
		testAccTemplateMetadata + `
	}
	`
	testAccCheckNutanixVirtualMachineConfigDisk := `
resource "nutanix_virtual_machine" "my-machine" {
	name = "%s"
` + testAccTemplateDiskBody + `
}`

	basicVars.powerState = "POWERED_OFF"
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigDisk)
	log.Printf("[DEBUG] template config= %s", config)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
				),
			},
		},
	})
}

// testing vms with disk list
func TestAccNutanixVirtualMachine_diskList2(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	diskList := diskSet()
	vmName := "nutanix_virtual_machine.my-machine"
	testAccTemplateDiskBody := testAccTemplateSpecBody +
		testAccTemplateResourcesBody + diskList + `
		}
	}
	api_version = "%s"
	` +
		testAccTemplateMetadata + `
	}
	`
	testAccCheckNutanixVirtualMachineConfigDisk := `
resource "nutanix_virtual_machine" "my-machine" {
	name = "%s"
` + testAccTemplateDiskBody + `
}`

	basicVars.powerState = "POWERED_OFF"
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigDisk)
	log.Printf("[DEBUG] template config= %s", config)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.#", "2"),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.0.data_source_reference.#", "1"),
					resource.TestCheckResourceAttr(vmName, diskSourceReference0Key+".kind", NutanixDiskKind[0]),
					resource.TestCheckResourceAttr(vmName, diskSourceReference0Key+".name", NutanixDiskName[0]),
					resource.TestCheckResourceAttr(vmName, diskSourceReference0Key+".uuid", NutanixDiskUUID[0]),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.0.device_properties.#", "1"),
					resource.TestCheckResourceAttr(vmName, deviceProperties0Key+".device_type", NutanixDiskDeviceType[0]),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.0.disk_size_mib", NutanixDiskSize[0]),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.1.data_source_reference.#", "1"),
					resource.TestCheckResourceAttr(vmName, diskSourceReference1Key+".kind", NutanixDiskKind[1]),
					resource.TestCheckResourceAttr(vmName, diskSourceReference1Key+".name", NutanixDiskName[1]),
					resource.TestCheckResourceAttr(vmName, diskSourceReference1Key+".uuid", NutanixDiskUUID[1]),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.1.device_properties.#", "1"),
					resource.TestCheckResourceAttr(vmName, deviceProperties1Key+".device_type", NutanixDiskDeviceType[1]),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.1.disk_size_mib", NutanixDiskSize[1]),
				),
			},
		},
	})
}

// testing update memory in vm
func TestAccNutanixVirtualMachine_updateMemory1(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	basicVars.memorySizeMb = NutanixUpdateMemorySize
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)
	log.Printf("[DEBUG] template config= %s", config)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
				),
			},
		},
	})
}

// testing update memory in vm
func TestAccNutanixVirtualMachine_updateMemory2(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	basicVars.powerState = "POWERED_OFF"
	configOFF := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)
	basicVars.memorySizeMb = NutanixUpdateMemorySize
	basicVars.powerState = NutanixPowerState
	configON := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)
	log.Printf("[DEBUG] template configOFF= %s", configOFF)
	log.Printf("[DEBUG] template configON= %s", configON)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configOFF,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
				),
			},
			resource.TestStep{
				Config: configON,
				Check: resource.ComposeTestCheckFunc(
					TestFuncData{vm: vm, memorySizeMb: NutanixUpdateMemorySize, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
				),
			},
		},
	})
}

// testing update name of the vm
func TestAccNutanixVirtualMachine_updateName1(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	config := basicVars.testSprintfTemplateBodyUpdateName(testAccCheckNutanixVirtualMachineConfigReallyBasic)
	log.Printf("[DEBUG] template config= %s", config)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
				),
			},
		},
	})
}

// testing update name of the vm
func TestAccNutanixVirtualMachine_updateName2(t *testing.T) {
	var vm vmdefn.VirtualMachine
	basicVars := setupTemplateBasicBodyVars()
	basicVars.powerState = "POWERED_OFF"
	configOFF := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)
	basicVars.powerState = NutanixPowerState
	configON := basicVars.testSprintfTemplateBodyUpdateName(testAccCheckNutanixVirtualMachineConfigReallyBasic)
	log.Printf("[DEBUG] template config= %s", configOFF)

	log.Printf("[DEBUG] template configON= %s", configON)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: configOFF,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixVirtualMachineExists("nutanix_virtual_machine.my-machine", &vm),
				),
			},
			resource.TestStep{
				Config: configON,
				Check: resource.ComposeTestCheckFunc(
					TestFuncData{vm: vm, name: NutanixUpdateName, vmName: "nutanix_virtual_machine.my-machine"}.testCheckFuncBasic(),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {

	for i := range s.RootModule().Resources {
		if s.RootModule().Resources[i].Type != "nutanix_virtual_machine" {
			continue
		}
		id := string(s.RootModule().Resources[i].Primary.ID)
		if id == "" {
			err := errors.New("ID is already set to the null string")
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckNutanixVirtualMachineExists(n string, vm *vmdefn.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f, _ := os.Create("/tmp/check")
		w := bufio.NewWriter(f)
		defer f.Close()
		defer w.Flush()
		fmt.Fprintf(w, "%+v \n %+v\n %+v", s.Serial, s, *s)
		terraformState = fmt.Sprintf("%+v", s)
		specKey = "spec." + hashmapKey("spec", "resources")
		specResourcesKey = specKey + ".resources." + hashmapKey(specKey+".resources", "power_state")
		metadataKey = "metadata." + hashmapKey("metadata", "kind")
		diskSourceReference0Key = specResourcesKey + ".disk_list.0.data_source_reference." + hashmapKey(specResourcesKey+".disk_list.0.data_source_reference", "uuid")
		diskSourceReference1Key = specResourcesKey + ".disk_list.1.data_source_reference." + hashmapKey(specResourcesKey+".disk_list.1.data_source_reference", "uuid")
		deviceProperties0Key = specResourcesKey + ".disk_list.0.device_properties." + hashmapKey(specResourcesKey+".disk_list.0.device_properties", "device_type")
		deviceProperties1Key = specResourcesKey + ".disk_list.1.device_properties." + hashmapKey(specResourcesKey+".disk_list.1.device_properties", "device_type")
		if n == "" {
			return fmt.Errorf("No vm name passed in")
		}
		if vm == nil {
			return fmt.Errorf("No vm obj passed in")
		}
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		ipAddress := string(rs.Primary.Attributes["ip_address"])
		if ipAddress == "" {
			fmt.Errorf("ip_address is not defined")
		}

		return nil
	}
}
