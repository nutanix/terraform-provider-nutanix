package nutanix

import (
	"bufio"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	//st "github.com/ideadevice/terraform-ahv-provider-plugin/virtualmachinestruct"
	"log"
	"os"
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

type TemplateBasicBodyVars struct {
	name         string
	numVCPUs     string
	numSockets   string
	memorySizeMb string
	powerState   string
	APIversion   string
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

	log.Printf("%+v", test.vm)
	return testAccCheckNutanixVirtualMachineExists(vmName, &test.vm),
		resource.TestCheckResourceAttr(vmName, "api_version", APIversion),
		resource.TestCheckResourceAttr(vmName, "spec.#", "1"),
		resource.TestCheckResourceAttr(vmName, "spec.1373794133.resources.#", "1"),
		resource.TestCheckResourceAttr(vmName, "spec.1373794133.resources.3082566069.power_state", powerState),
		resource.TestCheckResourceAttr(vmName, "spec.1373794133.resources.3082566069.memory_size_mb", memorySizeMb),
		resource.TestCheckResourceAttr(vmName, "spec.1373794133.resources.3082566069.num_sockets", numSockets),
		resource.TestCheckResourceAttr(vmName, "spec.1373794133.resources.3082566069.num_vcpus_per_socket", numVCPUs),
		resource.TestCheckResourceAttr(vmName, "metadata.#", "1"),
		resource.TestCheckResourceAttr(vmName, "metadata.341341886.kind", kind),
		resource.TestCheckResourceAttr(vmName, "metadata.341341886.spec_version", specVersion),
		resource.TestCheckResourceAttr(vmName, "name", name)

}

const testAccCheckNutanixVirtualMachineConfigReallyBasic = `
resource "nutanix_virtual_machine" "my-machine" {
	name = "kritagya_test1"
` + testAccTemplateBasicBodyWithEnd

const testAccTemplateBasicBody = `
	spec = {
		name = "%s"
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
		}
	}
	api_version = "%s"
	metadata = {
		kind = "%s"
		spec_version = %s
		name = "%s"
		categories = {
			"Project" = "nucalm"
		}
	}
`
const testAccTemplateBasicBodyWithEnd = testAccTemplateBasicBody + `
}`

func TestAccNutanixVirtualMachine_basic(t *testing.T) {
	var vm Machine
	basicVars := setupTemplateBasicBodyVars()
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigReallyBasic)

	log.Printf("[DEBUG] template= %s", testAccCheckNutanixVirtualMachineConfigReallyBasic)
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

const testAccCheckNutanixVirtualMachineConfigDebug = `
provider "nutanix"{
	username = "admin"
	password = "Nutanix123#"
	endpoint = "10.5.68.6"
}

` + testAccCheckNutanixVirtualMachineConfigReallyBasic

func TestAccNutanixVirtualMachine_client_debug(t *testing.T) {
	var vm Machine
	basicVars := setupTemplateBasicBodyVars()
	config := basicVars.testSprintfTemplateBody(testAccCheckNutanixVirtualMachineConfigDebug)

	log.Printf("[DEBUG] template= %s", testAccCheckNutanixVirtualMachineConfigDebug)
	log.Printf("[DEBUG] template config= %s", config)

	testExists, testAPIVersion, testSpec, testResources, testPowerState, testMemorySizeMb, testNumSockets, testNumVcpus, testMetadata, testKind, testSpecVersion, testName := TestFuncData{vm: vm}.testCheckFuncBasic()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testBasicPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixVirtualMachineDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testExists, testMetadata, testName, testSpec, testResources, testSpecVersion, testKind, testNumVcpus, testNumSockets, testMemorySizeMb, testPowerState, testAPIVersion, testAccCheckDebugExists(),
				),
			},
		},
	})
}

func testAccCheckNutanixVirtualMachineDestroy(s *terraform.State) error {
	//client := testAccProvider.Meta().(*MyClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine" {
			continue
		}
		/*	machine := Machine{
				Spec: &st.SpecStruct{
					Name: string(rs.Primary.Attributes["name"]),
				},
			}
		*/
		return nil
	}
	return nil
}

func testAccCheckDebugExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return nil
	}
}

func testAccCheckNutanixVirtualMachineExists(n string, vm *Machine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f, _ := os.Create("/tmp/check")
		w := bufio.NewWriter(f)
		defer f.Close()
		defer w.Flush()
		fmt.Fprintf(w, "%+v \n %+v\n %+v", s.Serial, s, *s)
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
		//client := testAccProvider.Meta().(*MyClient)

		return nil
	}
}
