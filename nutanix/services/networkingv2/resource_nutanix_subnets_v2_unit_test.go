package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitExpandDhcpOptions(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandDhcpOptions(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandDhcpOptions([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandDhcpOptions([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIPAddress(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddress([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddress([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4Address(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Address([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv6Address(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv6Address(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Address([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandVirtualSwitch(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVirtualSwitch([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandCluster(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandCluster(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandCluster([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandCluster([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandHost(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandHost(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandHost([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandHost([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4Subnet(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4Subnet(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Subnet([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4SubnetFromMap(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4SubnetFromMap(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4SubnetFromMap(map[string]interface{}{"prefix_length": 24})
		_ = got
	})

}

func TestUnitExpandIPv6Subnet(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv6Subnet(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Subnet([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv6SubnetFromMap(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6SubnetFromMap(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6SubnetFromMap(map[string]interface{}{"prefix_length": 24})
		_ = got
	})

}

func TestUnitExpandVpc(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVpc([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandVpcDhcpOptions(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVpcDhcpOptions([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandExternalSubnet(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandExternalSubnet(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandExternalSubnet([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandExternalSubnet([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandGatewayNodeReference(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandGatewayNodeReference([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPAddressMap(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPAddressMap([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4AddressMap(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv4AddressMap(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressMap([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv6AddressMap(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPv6AddressMap(nil))
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6AddressMap([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPSubnet(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIPSubnet(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPSubnet([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPSubnet([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPUsage(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPUsage([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIPConfig(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPConfig([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPConfig([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4Config(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Config([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv6Config(t *testing.T) {
	t.Parallel()

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Config([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4Pool(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIPv4Pool(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Pool([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4Pool([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv6Pool(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = expandIPv6Pool(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Pool([]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv6Pool([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}
