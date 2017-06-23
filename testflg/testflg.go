package testflg

/*
Uses Viper to read values from ENV variable or commandline flags transparently.

a flag with name "x-y" can be set in CLI as <binary> --x-y
if the same flag has to be set in ENV, it has tobe set as X_Y

conflicts & resolution order in the descending order of precedence
    flag
    env
*/

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

// NutanixUsername username for api call
var NutanixUsername string

// NutanixPassword password for api call
var NutanixPassword string

// NutanixEndpoint endpoint must be set
var NutanixEndpoint string

// NutanixInsecure insecure flag must set true to allow provider to perform insecure SSL requests.
var NutanixInsecure bool

// NutanixPort port for api call
var NutanixPort string

// NutanixNumSockets is num_sockets for the testcase vm
var NutanixNumSockets string

// NutanixNumVCPUs is num_vcpus for the testcase vm
var NutanixNumVCPUs string

// NutanixMemorySize is the memory_size_mb for testcase vm
var NutanixMemorySize string

// NutanixPowerState is power_state for testcase vm
var NutanixPowerState string

// NutanixDiskNo is the number of disks attached to the disktestcase vm
var NutanixDiskNo string

// NutanixDiskKind is slice of Kind of disks
var NutanixDiskKind []string

// NutanixDiskName is slice of disk names
var NutanixDiskName []string

// NutanixDiskUUID is slice of uuid of disks
var NutanixDiskUUID []string

// NutanixDiskSize is slice of size of disks
var NutanixDiskSize []string

// NutanixDiskDeviceType is slice of device types
var NutanixDiskDeviceType []string

// NutanixName is the name of the vm
var NutanixName string

// NutanixUpdateMemorySize is the memory size to which vm gets upgraded in updateMemory testcase
var NutanixUpdateMemorySize string

// NutanixUpdateName is the updated name of the vm in updateName testcase
var NutanixUpdateName string

// NutanixNicType is the nic_type of network adapter
var NutanixNicType string

// NutanixNicKind is the kind of network adapter
var NutanixNicKind string

// NutanixNicUUID is the nic_uuid of network adapter
var NutanixNicUUID string

// NutanixNetworkFunctionType is the network_function_type of network adapter
var NutanixNetworkFunctionType string

// NutanixProject is name any of project inside metadata categories.
var NutanixProject string

func init() {
	var diskKind1, diskKind2, diskName1, diskName2, diskUUID1, diskUUID2 string
	var diskDeviceType1, diskDeviceType2, diskSize1, diskSize2 string
	flag.StringVar(&NutanixUsername, "username", "", "username for api call")
	flag.StringVar(&NutanixPassword, "password", "", "password for api call")
	flag.StringVar(&NutanixEndpoint, "endpoint", "", "endpoint must be set")
	flag.BoolVar(&NutanixInsecure, "insecure", false, "insecure flag must set true to allow provider to perform insecure SSL requests. ")
	flag.StringVar(&NutanixPort, "port", "9440", "port for api call")
	flag.StringVar(&NutanixNumSockets, "num-sockets", "1", "This is num_sockets for the testcase vm.")
	flag.StringVar(&NutanixNumVCPUs, "num-vcpus", "1", "This is num_vcpus for the testcase vm.")
	flag.StringVar(&NutanixMemorySize, "memory-size", "1024", "This is the memory_size_mb for testcase vm.")
	flag.StringVar(&NutanixPowerState, "power-state", "ON", "This is power_state for testcase vm.")
	flag.StringVar(&NutanixName, "name", "vm_test1", "This is the name of the vm.")
	flag.StringVar(&NutanixUpdateName, "update-name", "vm_test2", "This is the updated name of the vm in updateName testcase.")
	flag.StringVar(&NutanixUpdateMemorySize, "update-memory-size", "2048", "This is the memory size to which vm gets upgraded in updateMemory testcase.")
	flag.StringVar(&NutanixDiskNo, "diskNo", "2", "This is the number of disks attached to the disktestcase vm.")
	flag.StringVar(&diskKind1, "disk-kind-1", "image", "This is Kind field for the first disk.")
	flag.StringVar(&diskName1, "disk-name-1", "Centos7", "This is disk name of first disk.")
	flag.StringVar(&diskUUID1, "disk-uuid-1", "9eabbb39-1baf-4872-beaf-adedcb612a0b", "This is UUID of first disk.")
	flag.StringVar(&diskSize1, "disk-size-1", "1", "This is size of the first disk")
	flag.StringVar(&diskDeviceType1, "disk-device-type-1", "DISK", "This is device type for the first disk.")
	flag.StringVar(&diskKind2, "disk-kind-2", "image", "This is Kind field for the second disk.")
	flag.StringVar(&diskName2, "disk-name-2", "Centos7", "This is disk name of second disk.")
	flag.StringVar(&diskUUID2, "disk-uuid-2", "9eabbb39-1baf-4872-beaf-adedcb612a0b", "This is UUID of second disk.")
	flag.StringVar(&diskSize2, "disk-size-2", "1", "This is size of the second disk")
	flag.StringVar(&diskDeviceType2, "disk-device-type-2", "DISK", "This is device type for the second disk.")
	flag.StringVar(&NutanixNicType, "nic-type", "NORMAL_NIC", "This is the nic_type of network adapter.")
	flag.StringVar(&NutanixNicUUID, "nic-uuid", "c03ecf8f-aa1c-4a07-af43-9f2f198713c0", "This is the nic_uuid of network adapter.")
	flag.StringVar(&NutanixNicKind, "nic-kind", "subnet", "This is the kind of network adapter.")
	flag.StringVar(&NutanixNetworkFunctionType, "network-function-nic-type", "INGRESS", "This is the network_function_type of network adapter.")
	flag.StringVar(&NutanixProject, "project", "nucalm", "Name any of project inside metadata categories.")

	//pflag configuration
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()
	pflag.Visit(func(f *pflag.Flag) {
		fmt.Printf("FlagValue %s overridden: %s -> %s\n", f.Name, f.DefValue, f.Value)
	})

	//Env configuration
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	//Config Init
	NutanixUsername = viper.GetString("username")
	NutanixPassword = viper.GetString("password")
	NutanixEndpoint = viper.GetString("endpoint")
	NutanixInsecure = viper.GetBool("insecure")
	NutanixPort = viper.GetString("port")
	NutanixNumSockets = viper.GetString("num-sockets")
	NutanixNumVCPUs = viper.GetString("num-vcpus")
	NutanixMemorySize = viper.GetString("memory-size")
	NutanixPowerState = viper.GetString("power-state")
	NutanixName = viper.GetString("name")
	NutanixUpdateName = viper.GetString("update-name")
	NutanixUpdateMemorySize = viper.GetString("update-memory-size")
	NutanixDiskNo = viper.GetString("diskNo")
	diskKind1 = viper.GetString("disk-kind-1")
	diskName1 = viper.GetString("disk-name-1")
	diskUUID1 = viper.GetString("disk-uuid-1")
	diskSize1 = viper.GetString("disk-size-1")
	diskDeviceType1 = viper.GetString("disk-device-type-1")
	diskKind2 = viper.GetString("disk-kind-2")
	diskName2 = viper.GetString("disk-name-2")
	diskUUID2 = viper.GetString("disk-uuid-2")
	diskSize2 = viper.GetString("disk-size-2")
	diskDeviceType2 = viper.GetString("disk-device-type-2")
	NutanixNicType = viper.GetString("nic-type")
	NutanixNicUUID = viper.GetString("nic-uuid")
	NutanixNicKind = viper.GetString("nic-kind")
	NutanixNetworkFunctionType = viper.GetString("network-function-nic-type")
	NutanixProject = viper.GetString("project")

	// Appending to the Disk List
	NutanixDiskKind = append(NutanixDiskKind, diskKind1)
	NutanixDiskName = append(NutanixDiskName, diskName1)
	NutanixDiskUUID = append(NutanixDiskUUID, diskUUID1)
	NutanixDiskSize = append(NutanixDiskSize, diskSize1)
	NutanixDiskDeviceType = append(NutanixDiskDeviceType, diskDeviceType1)
	NutanixDiskKind = append(NutanixDiskKind, diskKind2)
	NutanixDiskName = append(NutanixDiskName, diskName2)
	NutanixDiskUUID = append(NutanixDiskUUID, diskUUID2)
	NutanixDiskSize = append(NutanixDiskSize, diskSize2)
	NutanixDiskDeviceType = append(NutanixDiskDeviceType, diskDeviceType2)
}
