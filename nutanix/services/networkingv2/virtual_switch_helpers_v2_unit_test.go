package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	commonConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/common/v1/config"
	networkingConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitExpandVirtualSwitchClusters(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVirtualSwitchClusters(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVirtualSwitchClusters([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVirtualSwitchClusters([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandVirtualSwitchHosts(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandVirtualSwitchHosts(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandVirtualSwitchHosts([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandVirtualSwitchHosts([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandIPv4AddressVS(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressVS(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4AddressVS(map[string]interface{}{"value": "UNKNOWN", "prefix_length": 24})
		_ = got
	})

}

func TestUnitExpandIPv4SubnetVS(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4SubnetVS(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIPv4SubnetVS(map[string]interface{}{"prefix_length": 24})
		_ = got
	})

}

func TestUnitExpandIgmpSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandIgmpSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandIgmpSpec([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandIgmpSpec([]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}})
		_ = got
	})

}

func TestUnitExpandQuerierSpec(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandQuerierSpec(map[string]interface{}{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandQuerierSpec(map[string]interface{}{"is_querier_enabled": true})
		_ = got
	})

}

func TestUnitExpandStringList(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandStringList(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := expandStringList([]interface{}{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := expandStringList([]interface{}{"a"})
		_ = got
	})

}

func TestUnitFlattenVirtualSwitchClusters(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVirtualSwitchClusters(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualSwitchClusters([]networkingConfig.Cluster{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualSwitchClusters([]networkingConfig.Cluster{networkingConfig.Cluster{}})
		_ = got
	})

}

func TestUnitFlattenVirtualSwitchHosts(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenVirtualSwitchHosts(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualSwitchHosts([]networkingConfig.Host{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualSwitchHosts([]networkingConfig.Host{networkingConfig.Host{}})
		_ = got
	})

}

func TestUnitFlattenIPv4AddressVS(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIPv4AddressVS(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressVS(&commonConfig.IPv4Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressVS(&commonConfig.IPv4Address{})
		_ = got
	})

}

func TestUnitFlattenIPv4SubnetVS(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIPv4SubnetVS(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4SubnetVS(&networkingConfig.IPv4Subnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4SubnetVS(&networkingConfig.IPv4Subnet{})
		_ = got
	})

}

func TestUnitFlattenIgmpSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIgmpSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIgmpSpec(&networkingConfig.IgmpSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIgmpSpec(&networkingConfig.IgmpSpec{})
		_ = got
	})

}

func TestUnitFlattenQuerierSpec(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenQuerierSpec(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenQuerierSpec(&networkingConfig.QuerierSpec{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenQuerierSpec(&networkingConfig.QuerierSpec{})
		_ = got
	})

}
