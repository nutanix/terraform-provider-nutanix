package nutanix

import (
	"flag"
	"github.com/hashicorp/terraform/builtin/providers/template"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"strconv"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccTemplateProvider *schema.Provider
var terraformState string
var NutanixUsername string
var NutanixPassword string
var NutanixEndpoint string
var NutanixInsecure bool
var NutanixPort string
var NutanixNumSockets string
var NutanixNumVCPUs string
var NutanixKind string
var NutanixMemorySize string
var NutanixPowerState string
var NutanixSpecVersion string
var NutanixAPIVersion string
var NutanixDiskNo string
var NutanixDiskKind []string
var NutanixDiskName []string
var NutanixDiskUUID []string
var NutanixDiskSize []string
var NutanixDiskDeviceType []string
var NutanixName string
var NutanixUpdateMemorySize string
var NutanixUpdateName string

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccTemplateProvider = template.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"nutanix":  testAccProvider,
		"template": testAccTemplateProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func init() {
	var temp string
	os.Setenv("TF_ACC", "1")
	flag.StringVar(&NutanixUsername, "username", "", "username for api call")
	flag.StringVar(&NutanixPassword, "password", "", "password for api call")
	flag.StringVar(&NutanixEndpoint, "endpoint", "", "endpoint must be set")
	flag.BoolVar(&NutanixInsecure, "insecure", false, "insecure flag must set true to allow provider to perform insecure SSL requests.")
	flag.StringVar(&NutanixPort, "port", "9440", "port for api call")
	flag.StringVar(&NutanixNumSockets, "numSockets", "1", "num_sockets")
	flag.StringVar(&NutanixNumVCPUs, "numVCPUs", "1", "num_vcpus")
	flag.StringVar(&NutanixKind, "kind", "vm", "kind")
	flag.StringVar(&NutanixMemorySize, "memorySize", "1024", "memory_size_mb")
	flag.StringVar(&NutanixPowerState, "powerState", "POWERED_ON", "power_state")
	flag.StringVar(&NutanixSpecVersion, "specVersion", "0", "spec_version")
	flag.StringVar(&NutanixAPIVersion, "apiVersion", "3.0", "api_version")
	flag.StringVar(&NutanixName, "name", "vm_test1", "name")
	flag.StringVar(&NutanixUpdateName, "updateName", "vm_test2", "update_name")
	flag.StringVar(&NutanixUpdateMemorySize, "updateMemorySize", "2048", "update_memory_size_name")
	flag.StringVar(&NutanixDiskNo, "diskNo", "2", "disk_No")
	flag.StringVar(&temp, "diskKind1", "image", "disk_kind_1")
	NutanixDiskKind = append(NutanixDiskKind, temp)
	flag.StringVar(&temp, "diskName1", "Centos7", "disk_name_1")
	NutanixDiskName = append(NutanixDiskName, temp)
	flag.StringVar(&temp, "diskUUID1", "9eabbb39-1baf-4872-beaf-adedcb612a0b", "disk_uuid_1")
	NutanixDiskUUID = append(NutanixDiskUUID, temp)
	flag.StringVar(&temp, "diskSize1", "1", "disk_size_1")
	NutanixDiskSize = append(NutanixDiskSize, temp)
	flag.StringVar(&temp, "diskDeviceType1", "DISK", "disk_device_type_1")
	NutanixDiskDeviceType = append(NutanixDiskDeviceType, temp)
	flag.StringVar(&temp, "diskKind2", "image", "disk_kind_2")
	NutanixDiskKind = append(NutanixDiskKind, temp)
	flag.StringVar(&temp, "diskName2", "Centos7", "disk_name_2")
	NutanixDiskName = append(NutanixDiskName, temp)
	flag.StringVar(&temp, "diskUUID2", "9eabbb39-1baf-4872-beaf-adedcb612a0b", "disk_uuid_2")
	NutanixDiskUUID = append(NutanixDiskUUID, temp)
	flag.StringVar(&temp, "diskSize2", "1", "disk_size_2")
	NutanixDiskSize = append(NutanixDiskSize, temp)
	flag.StringVar(&temp, "diskDeviceType2", "DISK", "disk_device_type_2")
	NutanixDiskDeviceType = append(NutanixDiskDeviceType, temp)

}
func main() {
	flag.Parse()
}

func testAccPreCheck(t *testing.T) {
	os.Setenv("NUTANIX_USERNAME", NutanixUsername)
	os.Setenv("NUTANIX_PASSWORD", NutanixPassword)
	os.Setenv("NUTANIX_ENDPOINT", NutanixEndpoint)
	os.Setenv("NUTANIX_PORT", NutanixPort)
	os.Setenv("NUTANIX_INSECURE", strconv.FormatBool(NutanixInsecure))
	if NutanixUsername == "" {
		t.Fatal("username flag must be set for acceptance tests")
	}
	if NutanixPassword == "" {
		t.Fatal("password must be set for acceptance tests")
	}
	if NutanixEndpoint == "" {
		t.Fatal("endpoint flag must be set for acceptance tests")
	}
	if NutanixInsecure == false {
		t.Fatal("insecure flag must be set true for acceptance tests")
	}
	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}
