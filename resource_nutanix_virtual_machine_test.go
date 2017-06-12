package nutanix

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Base setup function to check that api_version is set
func testBasicPreCheck(t *testing.T) {
	testAccPreCheck(t)

	if v := os.Getenv("NUTANIX_API_VERSION"); v == "" {
		log.Printf("yo")
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
		name:         os.Getenv("NUTANIX_NAME"),
		numSockets:   os.Getenv("NUTANIX_NUM_SOCKETS"),
		numVCPUs:     os.Getenv("NUTANIX_NUM_VCPUS"),
		memorySizeMb: os.Getenv("NUTANIX_MEMORY_SIZE_MB"),
		powerState:   os.Getenv("NUTANIX_POWER_STATE"),
		kind:         os.Getenv("NUTANIX_KIND"),
		specVersion:  os.Getenv("NUTANIX_SPEC_VERSION"),
		APIVersion:   os.Getenv("NUTANIX_API_VERSION"),
	}
	return data
}

// Basic data to create series of testing functions
type TestFuncData struct {
	vm           Machine
	vmName       string
	name         string
	numVCPUs     string
	numSockets   string
	memorySizeMb string
	powerState   string
	APIversion   string
	kind         string
	specVersion  string
}

func hashmapKey(s, t string) string {
	words := strings.Fields(terraformState)
	prefix := s + "."
	suffix := "." + t
	for _, word := range words {
		if (word == strings.TrimPrefix(word, prefix+"#")) && (word != strings.TrimPrefix(word, prefix)) {
			str1 := strings.TrimPrefix(word, prefix)
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
// APIversion defaults to 3.0 specVersion 0 and memorySizeMb tp 1024
// kind defaults to "vm" and powerState to "POWERED_ON", vmName to "nutanix_virtual_machine"
func (test TestFuncData) testCheckFuncBasic() (resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc, resource.TestCheckFunc) {
	vmName := test.vmName
	if vmName == "" {
		vmName = "nutanix_virtual_machine.my-machine"
	}
	kind := test.kind
	if kind == "" {
		kind = "vm"
	}
	powerState := test.powerState
	if powerState == "" {
		powerState = "POWERED_ON"
	}
	APIversion := test.APIversion
	if APIversion == "" {
		APIversion = "3.0"
	}
	numSockets := test.numSockets
	if numSockets == "" {
		numSockets = "1"
	}
	numVCPUs := test.numVCPUs
	if numVCPUs == "" {
		numVCPUs = "1"
	}
	name := test.name
	if name == "" {
		name = "kritagya_test1"
	}
	memorySizeMb := test.memorySizeMb
	if memorySizeMb == "" {
		memorySizeMb = "1024"
	}
	specVersion := test.specVersion
	if specVersion == "" {
		specVersion = "0"
	}
	return testAccCheckNutanixVirtualMachineExists(vmName, &test.vm),
		resource.TestCheckResourceAttr(vmName, "api_version", APIversion),
		resource.TestCheckResourceAttr(vmName, "spec.#", "1"),
		resource.TestCheckResourceAttr(vmName, specKey+".resources.#", "1"),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".power_state", powerState),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".memory_size_mb", memorySizeMb),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".num_sockets", numSockets),
		resource.TestCheckResourceAttr(vmName, specResourcesKey+".num_vcpus_per_socket", numVCPUs),
		resource.TestCheckResourceAttr(vmName, "metadata.#", "1"),
		resource.TestCheckResourceAttr(vmName, metadataKey+".kind", kind),
		resource.TestCheckResourceAttr(vmName, metadataKey+".spec_version", specVersion),
		resource.TestCheckResourceAttr(vmName, "name", name)

}

const testAccCheckNutanixVirtualMachineConfigReallyBasic = `
resource "nutanix_virtual_machine" "my-machine" {
	name = "kritagya_test1"
` + testAccTemplateBasicBodyWithEnd

const testAccTemplateSpecBody = `
spec = {
	name = "%s"
`
const testAccTemplateResourcesBody = `
		resources = {
			num_sockets = %s
			num_vcpus_per_socket = %s
			memory_size_mb = %s
			power_state = "%s"
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
const testAccTemplateBasicBodyWithEnd = testAccTemplateBasicBody + `
}`

// testing vms with basic config
func TestAccNutanixVirtualMachine_basic1(t *testing.T) {
	var vm Machine
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

// testing vms with basic config
func TestAccNutanixVirtualMachine_basic2(t *testing.T) {
	var vm Machine
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
	diskNo, _ := strconv.Atoi(os.Getenv("NUTANIX_DISK_NO"))
	for i := 0; i < diskNo-1; i++ {
		disk := fmt.Sprintf(diskTemplate, os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_KIND"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_NAME"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_UUID"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_DEVICETYPE"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_SIZE")) + ","
		diskList = diskList + disk
	}
	i := diskNo - 1
	if diskNo > 0 {
		disk := fmt.Sprintf(diskTemplate, os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_KIND"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_NAME"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_UUID"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_DEVICETYPE"), os.Getenv("NUTANIX_DISK_"+strconv.Itoa(i)+"_SIZE"))
		diskList = diskList + disk
	}
	diskList = diskList + `
			]
	`
	return diskList
}

// testing vms with disk list
func TestAccNutanixVirtualMachine_diskList1(t *testing.T) {
	var vm Machine
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
	name = "kritagya_test1"
` + testAccTemplateDiskBody + `
}`

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
	var vm Machine
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
	name = "kritagya_test1"
` + testAccTemplateDiskBody + `
}`

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
					resource.TestCheckResourceAttr(vmName, diskSourceReference0Key+".kind", os.Getenv("NUTANIX_DISK_0_KIND")),
					resource.TestCheckResourceAttr(vmName, diskSourceReference0Key+".name", os.Getenv("NUTANIX_DISK_0_NAME")),
					resource.TestCheckResourceAttr(vmName, diskSourceReference0Key+".uuid", os.Getenv("NUTANIX_DISK_0_UUID")),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.0.device_properties.#", "1"),
					resource.TestCheckResourceAttr(vmName, deviceProperties0Key+".device_type", os.Getenv("NUTANIX_DISK_0_DEVICETYPE")),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.0.disk_size_mib", os.Getenv("NUTANIX_DISK_0_SIZE")),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.1.data_source_reference.#", "1"),
					resource.TestCheckResourceAttr(vmName, diskSourceReference1Key+".kind", os.Getenv("NUTANIX_DISK_1_KIND")),
					resource.TestCheckResourceAttr(vmName, diskSourceReference1Key+".name", os.Getenv("NUTANIX_DISK_1_NAME")),
					resource.TestCheckResourceAttr(vmName, diskSourceReference1Key+".uuid", os.Getenv("NUTANIX_DISK_1_UUID")),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.1.device_properties.#", "1"),
					resource.TestCheckResourceAttr(vmName, deviceProperties1Key+".device_type", os.Getenv("NUTANIX_DISK_1_DEVICETYPE")),
					resource.TestCheckResourceAttr(vmName, specResourcesKey+".disk_list.1.disk_size_mib", os.Getenv("NUTANIX_DISK_1_SIZE")),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine" {
			continue
		}
		id := string(rs.Primary.ID)
		if id == "" {
			err := errors.New("ID is already set to the null string")
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckNutanixVirtualMachineExists(n string, vm *Machine) resource.TestCheckFunc {
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
