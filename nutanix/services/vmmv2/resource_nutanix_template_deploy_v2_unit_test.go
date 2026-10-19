package vmmv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestUnitExpandVMConfigOverride(t *testing.T) {
	t.Parallel()
	schemaMap := ResourceNutanixTemplateDeployV2().Schema
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		got := expandVMConfigOverride([]interface{}{}, d)
		require.Nil(t, got)
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{
			"ext_id": "tmpl-1",
		})
		got := expandVMConfigOverride([]interface{}{map[string]interface{}{"name": "vm-1"}}, d)
		require.NotNil(t, got)
	})
}

func TestUnitExpandNic(t *testing.T) {
	t.Parallel()
	schemaMap := ResourceNutanixVirtualMachineV2().Schema
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.Nil(t, expandNic(nil, d, "nics"))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{})
		require.Nil(t, expandNic([]interface{}{}, d, "nics"))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{
			"nics": []interface{}{map[string]interface{}{"ext_id": "nic-1"}},
		})
		got := expandNic([]interface{}{map[string]interface{}{"ext_id": "nic-1"}}, d, "nics")
		require.NotNil(t, got)
		require.Equal(t, "nic-1", utils.StringValue(got[0].ExtId))
	})
}

func TestUnitExpandHostPcieDeviceReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandHostPcieDeviceReference(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandHostPcieDeviceReference([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandNicProfileReference(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandNicProfileReference(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNicProfileReference([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandEmulatedNic(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandEmulatedNic([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandNicNetworkInfo(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNicNetworkInfo([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4Info(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Info([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandGuestCustomizationParams(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandGuestCustomizationParams([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandNetworkFunctionChainReference(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandNetworkFunctionChainReference([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandSubnetReference(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandSubnetReference([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIpv4Config(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIpv4Config([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Address([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4AddressList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIPv4AddressList(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressList([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressList([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandOneOfGuestCustomizationParamsConfig(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfGuestCustomizationParamsConfig([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandOneOfSysprepSysprepScript(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfSysprepSysprepScript([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandOneOfCloudInitCloudInitScript(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandOneOfCloudInitCloudInitScript([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
