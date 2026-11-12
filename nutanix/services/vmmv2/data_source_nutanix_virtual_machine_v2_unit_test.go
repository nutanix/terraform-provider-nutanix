package vmmv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	import4 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/config"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/common/v1/response"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitFlattenVMSourceReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMSourceReference(&config.VmSourceReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMSourceReference(&config.VmSourceReference{})
		_ = got
	})

}

func TestUnitFlattenVMSourceReferenceEntityType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMSourceReferenceEntityType(func() *config.VmSourceReferenceEntityType { var v config.VmSourceReferenceEntityType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMSourceReferenceEntityType(func() *config.VmSourceReferenceEntityType { v := config.VmSourceReferenceEntityType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenCPUFeature(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenCPUFeature(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCPUFeature([]config.CpuFeature{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCPUFeature([]config.CpuFeature{config.CpuFeature(1)})
		_ = got
	})

}

func TestUnitFlattenCategoryReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenCategoryReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryReference([]config.CategoryReference{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategoryReference([]config.CategoryReference{config.CategoryReference{}})
		_ = got
	})

}

func TestUnitFlattenOwnershipInfo(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenOwnershipInfo(&config.OwnershipInfo{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenOwnershipInfo(&config.OwnershipInfo{})
		_ = got
	})

}

func TestUnitFlattenOwnerReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenOwnerReference(&config.OwnerReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenOwnerReference(&config.OwnerReference{})
		_ = got
	})

}

func TestUnitFlattenHostReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenHostReference(&config.HostReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenHostReference(&config.HostReference{})
		_ = got
	})

}

func TestUnitFlattenClusterReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenClusterReference(&config.ClusterReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenClusterReference(&config.ClusterReference{})
		_ = got
	})

}

func TestUnitFlattenAvailabilityZoneReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenAvailabilityZoneReference(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenAvailabilityZoneReference(&config.AvailabilityZoneReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenAvailabilityZoneReference(&config.AvailabilityZoneReference{})
		_ = got
	})

}

func TestUnitFlattenGuestCustomizationParams(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGuestCustomizationParams(&config.GuestCustomizationParams{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGuestCustomizationParams(&config.GuestCustomizationParams{})
		_ = got
	})

}

func TestUnitFlattenOneOfGuestCustomizationParamsConfig(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfGuestCustomizationParamsConfig(nil))
	})
}

func TestUnitFlattenInstallType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenInstallType(func() *config.InstallType { var v config.InstallType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenInstallType(func() *config.InstallType { v := config.InstallType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenOneOfSysprepSysprepScript(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfSysprepSysprepScript(nil))
	})
}

func TestUnitFlattenCloudInitDataSourceType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCloudInitDataSourceType(func() *config.CloudInitDataSourceType { var v config.CloudInitDataSourceType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCloudInitDataSourceType(func() *config.CloudInitDataSourceType { v := config.CloudInitDataSourceType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenOneOfCloudInitCloudInitScript(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfCloudInitCloudInitScript(nil))
	})
}

func TestUnitFlattenGuestTools(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGuestTools(&config.GuestTools{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGuestTools(&config.GuestTools{})
		_ = got
	})

}

func TestUnitFlattenNgtCapability(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenNgtCapability(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNgtCapability([]config.NgtCapability{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNgtCapability([]config.NgtCapability{config.NgtCapability(1)})
		_ = got
	})

}

func TestUnitFlattenOneOfVMBootConfig(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfVMBootConfig(nil))
	})
}

func TestUnitFlattenOneOfUefiBootBootDevice(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfUefiBootBootDevice(nil))
	})
}

func TestUnitFlattenOneOfLegacyBootBootDevice(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfLegacyBootBootDevice(nil))
	})
}

func TestUnitFlattenDiskAddress(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDiskAddress(&config.DiskAddress{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDiskAddress(&config.DiskAddress{})
		_ = got
	})

}

func TestUnitFlattenDiskBusType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDiskBusType(func() *config.DiskBusType { var v config.DiskBusType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDiskBusType(func() *config.DiskBusType { v := config.DiskBusType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenBootDeviceType(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenBootDeviceType(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenBootDeviceType([]config.BootDeviceType{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenBootDeviceType([]config.BootDeviceType{config.BootDeviceType(1)})
		_ = got
	})

}

func TestUnitFlattenNvramDevice(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNvramDevice(&config.NvramDevice{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNvramDevice(&config.NvramDevice{})
		_ = got
	})

}

func TestUnitFlattenVMDisk(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMDisk(&config.VmDisk{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMDisk(&config.VmDisk{})
		_ = got
	})

}

func TestUnitFlattenVMDiskContainerReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMDiskContainerReference(&config.VmDiskContainerReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMDiskContainerReference(&config.VmDiskContainerReference{})
		_ = got
	})

}

func TestUnitFlattenVMDiskStorageConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMDiskStorageConfig(&config.VmDiskStorageConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMDiskStorageConfig(&config.VmDiskStorageConfig{})
		_ = got
	})

}

func TestUnitFlattenDataSource(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDataSource(&config.DataSource{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDataSource(&config.DataSource{})
		_ = got
	})

}

func TestUnitFlattenOneOfDataSourceReference(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfDataSourceReference(nil))
	})
}

func TestUnitFlattenVMReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMReference(&config.VmReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMReference(&config.VmReference{})
		_ = got
	})

}

func TestUnitFlattenMachineType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenMachineType(func() *config.MachineType { var v config.MachineType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenMachineType(func() *config.MachineType { v := config.MachineType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenVtpmConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVtpmConfig(&config.VtpmConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVtpmConfig(&config.VtpmConfig{})
		_ = got
	})

}

func TestUnitFlattenApcConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenApcConfig(&config.ApcConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenApcConfig(&config.ApcConfig{})
		_ = got
	})

}

func TestUnitFlattenCPUModelReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCPUModelReference(&config.CpuModelReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCPUModelReference(&config.CpuModelReference{})
		_ = got
	})

}

func TestUnitFlattenADSFVmStorageConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenADSFVmStorageConfig(&config.ADSFVmStorageConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenADSFVmStorageConfig(&config.ADSFVmStorageConfig{})
		_ = got
	})

}

func TestUnitFlattenQosConfig(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenQosConfig(&config.QosConfig{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenQosConfig(&config.QosConfig{})
		_ = got
	})

}

func TestUnitFlattenDisk(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenDisk(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenDisk([]config.Disk{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDisk([]config.Disk{{ExtId: utils.StringPtr("disk-1")}})
		require.Len(t, got, 1)
	})
}

func TestUnitFlattenOneOfDiskBackingInfo(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfDiskBackingInfo(nil))
	})
}

func TestUnitFlattenCdRom(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenCdRom(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCdRom([]config.CdRom{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCdRom([]config.CdRom{config.CdRom{}})
		_ = got
	})

}

func TestUnitFlattenIsoType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIsoType(func() *config.IsoType { var v config.IsoType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIsoType(func() *config.IsoType { v := config.IsoType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenCdRomAddress(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCdRomAddress(&config.CdRomAddress{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCdRomAddress(&config.CdRomAddress{})
		_ = got
	})

}

func TestUnitFlattenCdRomBusType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCdRomBusType(func() *config.CdRomBusType { var v config.CdRomBusType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCdRomBusType(func() *config.CdRomBusType { v := config.CdRomBusType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenNic(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenNic(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNic([]config.Nic{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNic([]config.Nic{config.Nic{}})
		_ = got
	})

}

func TestUnitFlattenEmulatedNic(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenEmulatedNic(&config.EmulatedNic{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenEmulatedNic(&config.EmulatedNic{})
		_ = got
	})

}

func TestUnitFlattenEmulatedNicModel(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenEmulatedNicModel(func() *config.EmulatedNicModel { var v config.EmulatedNicModel; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenEmulatedNicModel(func() *config.EmulatedNicModel { v := config.EmulatedNicModel(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenNicNetworkInfo(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicNetworkInfo(&config.NicNetworkInfo{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNicNetworkInfo(&config.NicNetworkInfo{})
		_ = got
	})

}

func TestUnitFlattenNicType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNicType(func() *config.NicType { var v config.NicType; return &v }())
		_ = got
	})

}

func TestUnitFlattenNetworkFunctionChainReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkFunctionChainReference(&config.NetworkFunctionChainReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkFunctionChainReference(&config.NetworkFunctionChainReference{})
		_ = got
	})

}

func TestUnitFlattenNetworkFunctionNicType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkFunctionNicType(func() *config.NetworkFunctionNicType { var v config.NetworkFunctionNicType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNetworkFunctionNicType(func() *config.NetworkFunctionNicType { v := config.NetworkFunctionNicType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenSubnetReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSubnetReference(&config.SubnetReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSubnetReference(&config.SubnetReference{})
		_ = got
	})

}

func TestUnitFlattenVlanMode(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVlanMode(func() *config.VlanMode { var v config.VlanMode; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVlanMode(func() *config.VlanMode { v := config.VlanMode(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenIpv4Config(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv4Config(&config.Ipv4Config{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv4Config(&config.Ipv4Config{})
		_ = got
	})

}

func TestUnitFlattenIpv4Info(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv4Info(&config.Ipv4Info{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv4Info(&config.Ipv4Info{})
		_ = got
	})

}

func TestUnitFlattenIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Address(&import4.IPv4Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Address(&import4.IPv4Address{})
		_ = got
	})

}

func TestUnitFlattenIPv4AddressList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenIPv4AddressList(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressList([]import4.IPv4Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressList([]import4.IPv4Address{import4.IPv4Address{}})
		_ = got
	})

}

func TestUnitFlattenGpu(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenGpu(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGpu([]config.Gpu{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGpu([]config.Gpu{config.Gpu{}})
		_ = got
	})

}

func TestUnitFlattenGpuMode(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGpuMode(func() *config.GpuMode { var v config.GpuMode; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGpuMode(func() *config.GpuMode { v := config.GpuMode(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenGpuVendor(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenGpuVendor(func() *config.GpuVendor { var v config.GpuVendor; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenGpuVendor(func() *config.GpuVendor { v := config.GpuVendor(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenSBDF(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSBDF(&config.SBDF{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSBDF(&config.SBDF{})
		_ = got
	})

}

func TestUnitFlattenSerialPort(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenSerialPort(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSerialPort([]config.SerialPort{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSerialPort([]config.SerialPort{config.SerialPort{}})
		_ = got
	})

}

func TestUnitFlattenProtectionType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenProtectionType(func() *config.ProtectionType { var v config.ProtectionType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenProtectionType(func() *config.ProtectionType { v := config.ProtectionType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenProtectionPolicyState(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenProtectionPolicyState(&config.ProtectionPolicyState{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenProtectionPolicyState(&config.ProtectionPolicyState{})
		_ = got
	})

}

func TestUnitFlattenPolicyReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicyReference(&config.PolicyReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPolicyReference(&config.PolicyReference{})
		_ = got
	})

}

func TestUnitFlattenAPILink(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenAPILink(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenAPILink([]response.ApiLink{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenAPILink([]response.ApiLink{response.ApiLink{}})
		_ = got
	})

}

func TestUnitFlattenProjectReference(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenProjectReference(&config.ProjectReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenProjectReference(&config.ProjectReference{})
		_ = got
	})

}

const testDiskSerialID = "NFS_6_0_382_a36dcd9e_d389_4394_6557_be092fda879d"

func newTestVMDisk() *config.VmDisk {
	vmDisk := config.NewVmDisk()
	vmDisk.DiskExtId = utils.StringPtr("a36dcd9e-d389-4394-6557-be092fda879d")
	vmDisk.DiskSizeBytes = utils.Int64Ptr(10737418240)
	vmDisk.SerialId = utils.StringPtr(testDiskSerialID)
	return vmDisk
}

func newTestDisk(t *testing.T) config.Disk {
	t.Helper()
	backingInfo := config.NewOneOfDiskBackingInfo()
	if err := backingInfo.SetValue(*newTestVMDisk()); err != nil {
		t.Fatalf("failed to build OneOfDiskBackingInfo: %v", err)
	}
	disk := config.NewDisk()
	disk.BackingInfo = backingInfo
	return *disk
}

func TestFlattenVMDisk_SerialID(t *testing.T) {
	t.Parallel()

	got := flattenVMDisk(newTestVMDisk())

	if len(got) != 1 {
		t.Fatalf("expected exactly one element, got %d", len(got))
	}
	serialID, ok := got[0]["serial_id"].(*string)
	if !ok || serialID == nil || *serialID != testDiskSerialID {
		t.Fatalf("unexpected serial_id: %v", got[0]["serial_id"])
	}
}

func TestFlattenVMDisk_NoSerialID(t *testing.T) {
	t.Parallel()

	vmDisk := newTestVMDisk()
	vmDisk.SerialId = nil
	got := flattenVMDisk(vmDisk)

	if _, ok := got[0]["serial_id"]; ok {
		t.Fatalf("expected serial_id to be omitted when the API returns none, got %v", got[0]["serial_id"])
	}
}

func TestFlattenOneOfDiskBackingInfo_SerialID(t *testing.T) {
	t.Parallel()

	disk := newTestDisk(t)
	got := flattenOneOfDiskBackingInfo(disk.BackingInfo)

	if len(got) != 1 {
		t.Fatalf("expected exactly one element, got %d", len(got))
	}
	vmDisks, ok := got[0]["vm_disk"].([]map[string]interface{})
	if !ok || len(vmDisks) != 1 {
		t.Fatalf("expected one vm_disk, got %v", got[0]["vm_disk"])
	}
	serialID, ok := vmDisks[0]["serial_id"].(*string)
	if !ok || serialID == nil || *serialID != testDiskSerialID {
		t.Fatalf("unexpected serial_id: %v", vmDisks[0]["serial_id"])
	}
}

// Every vm_disk schema that receives flattenVMDisk or flattenOneOfDiskBackingInfo
// output must declare serial_id, otherwise d.Set fails for that resource.
func TestVMDiskSchemas_HaveSerialID(t *testing.T) {
	t.Parallel()

	resources := map[string]*schema.Resource{
		"nutanix_virtual_machine_v2":       ResourceNutanixVirtualMachineV2(),
		"data.nutanix_virtual_machine_v2":  DatasourceNutanixVirtualMachineV4(),
		"data.nutanix_virtual_machines_v2": DatasourceNutanixVirtualMachinesV4(),
		"nutanix_vm_clone_v2":              ResourceNutanixVMCloneV2(),
		"nutanix_template_v2":              ResourceNutanixTemplatesV2(),
		"nutanix_ova_vm_deploy_v2":         ResourceNutanixOvaVMDeploymentV2(),
		"nutanix_vm_cdrom_insert_eject_v2": ResourceNutanixVmsCdRomsInsertEjectV2(),
		"data.nutanix_template_v2":         DatasourceNutanixTemplateV2(),
		"data.nutanix_templates_v2":        DatasourceNutanixTemplatesV2(),
	}

	for name, r := range resources {
		found := 0
		walkVMDiskSchemas(r.Schema, func(path string, s map[string]*schema.Schema) {
			found++
			serialID, ok := s["serial_id"]
			if !ok {
				t.Errorf("%s: %s has no serial_id", name, path)
				return
			}
			if serialID.Type != schema.TypeString || !serialID.Computed || serialID.Optional || serialID.Required {
				t.Errorf("%s: %s.serial_id must be a computed-only string", name, path)
			}
		})
		if found == 0 {
			t.Errorf("%s: no vm_disk schema found", name)
		}
	}
}

// walkVMDiskSchemas calls fn for every nested schema shaped like a VmDisk,
// identified by its is_migration_in_progress attribute.
func walkVMDiskSchemas(s map[string]*schema.Schema, fn func(string, map[string]*schema.Schema)) {
	var walk func(string, map[string]*schema.Schema)
	walk = func(path string, s map[string]*schema.Schema) {
		if _, ok := s["is_migration_in_progress"]; ok {
			fn(path, s)
		}
		for k, v := range s {
			if r, ok := v.Elem.(*schema.Resource); ok {
				walk(path+"."+k, r.Schema)
			}
		}
	}
	walk("", s)
}

func TestSetDisksAndCdRoms_SerialID(t *testing.T) {
	t.Parallel()

	cdRom := config.NewCdRom()
	cdRom.BackingInfo = newTestVMDisk()
	disks := flattenDisk([]config.Disk{newTestDisk(t)})
	cdRoms := flattenCdRom([]config.CdRom{*cdRom})

	resources := map[string]*schema.Resource{
		"nutanix_virtual_machine_v2":      ResourceNutanixVirtualMachineV2(),
		"data.nutanix_virtual_machine_v2": DatasourceNutanixVirtualMachineV4(),
		"nutanix_vm_clone_v2":             ResourceNutanixVMCloneV2(),
	}

	for name, r := range resources {
		d := r.TestResourceData()
		if err := d.Set("disks", disks); err != nil {
			t.Fatalf("%s: setting disks: %v", name, err)
		}
		if err := d.Set("cd_roms", cdRoms); err != nil {
			t.Fatalf("%s: setting cd_roms: %v", name, err)
		}
		if got := d.Get("disks.0.backing_info.0.vm_disk.0.serial_id"); got != testDiskSerialID {
			t.Errorf("%s: disks serial_id = %q", name, got)
		}
		if got := d.Get("cd_roms.0.backing_info.0.serial_id"); got != testDiskSerialID {
			t.Errorf("%s: cd_roms serial_id = %q", name, got)
		}
	}

	d := ResourceNutanixVmsCdRomsInsertEjectV2().TestResourceData()
	if err := d.Set("backing_info", flattenVMDisk(newTestVMDisk())); err != nil {
		t.Fatalf("nutanix_vm_cdrom_insert_eject_v2: setting backing_info: %v", err)
	}
	if got := d.Get("backing_info.0.serial_id"); got != testDiskSerialID {
		t.Errorf("nutanix_vm_cdrom_insert_eject_v2: backing_info serial_id = %q", got)
	}
}
