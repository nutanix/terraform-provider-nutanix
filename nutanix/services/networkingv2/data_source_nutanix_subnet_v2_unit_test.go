package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/networking-go-client/v4/models/networking/v4/config"
)

func TestUnitFlattenLinks(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenLinks(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenLinks([]import2.ApiLink{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenLinks([]import2.ApiLink{import2.ApiLink{}})
		_ = got
	})

}

func TestUnitFlattenDhcpOptions(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenDhcpOptions(&import1.DhcpOptions{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenDhcpOptions(&import1.DhcpOptions{})
		_ = got
	})

}

func TestUnitFlattenNtpServer(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenNtpServer(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNtpServer([]config.IPAddress{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNtpServer([]config.IPAddress{config.IPAddress{}})
		_ = got
	})

}

func TestUnitFlattenIPv4(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4(&config.IPv4Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4(&config.IPv4Address{})
		_ = got
	})

}

func TestUnitFlattenIPv6(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv6(&config.IPv6Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv6(&config.IPv6Address{})
		_ = got
	})

}

func TestUnitFlattenIPConfig(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIPConfig(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPConfig([]import1.IPConfig{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPConfig([]import1.IPConfig{import1.IPConfig{}})
		_ = got
	})

}

func TestUnitFlattenIPv4Subnet(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Subnet(&import1.IPv4Subnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Subnet(&import1.IPv4Subnet{})
		_ = got
	})

}

func TestUnitFlattenIPv6Subnet(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv6Subnet(&import1.IPv6Subnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv6Subnet(&import1.IPv6Subnet{})
		_ = got
	})

}

func TestUnitFlattenPoolListIPv4(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenPoolListIPv4(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPoolListIPv4([]import1.IPv4Pool{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPoolListIPv4([]import1.IPv4Pool{import1.IPv4Pool{}})
		_ = got
	})

}

func TestUnitFlattenPoolListIPv6(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenPoolListIPv6(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenPoolListIPv6([]import1.IPv6Pool{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenPoolListIPv6([]import1.IPv6Pool{import1.IPv6Pool{}})
		_ = got
	})

}

func TestUnitFlattenIpv4Config(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv4Config(&import1.IPv4Config{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv4Config(&import1.IPv4Config{})
		_ = got
	})

}

func TestUnitFlattenIpv6Config(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv6Config(&import1.IPv6Config{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIpv6Config(&import1.IPv6Config{})
		_ = got
	})

}

func TestUnitFlattenSubnetType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenSubnetType(func() *import1.SubnetType { var v import1.SubnetType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenSubnetType(func() *import1.SubnetType { v := import1.SubnetType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenReservedIPAddresses(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenReservedIPAddresses(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenReservedIPAddresses([]config.IPAddress{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenReservedIPAddresses([]config.IPAddress{config.IPAddress{}})
		_ = got
	})

}

func TestUnitFlattenVirtualSwitch(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualSwitch(&import1.VirtualSwitch{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVirtualSwitch(&import1.VirtualSwitch{})
		_ = got
	})

}

func TestUnitFlattenBondMode(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenBondMode(func() *import1.BondModeType { var v import1.BondModeType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenBondMode(func() *import1.BondModeType { v := import1.BondModeType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenMetadata(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenMetadata(&config.Metadata{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenMetadata(&config.Metadata{})
		_ = got
	})

}

func TestUnitFlattenClusters(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenClusters(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenClusters([]import1.Cluster{})
		require.Nil(t, got)
	})

}

func TestUnitFlattenHosts(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenHosts(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenHosts([]import1.Host{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenHosts([]import1.Host{import1.Host{}})
		_ = got
	})

}

func TestUnitFlattenVPC(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVPC(&import1.Vpc{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVPC(&import1.Vpc{})
		_ = got
	})

}

func TestUnitFlattenCommonDhcpOptions(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenCommonDhcpOptions(&import1.VpcDhcpOptions{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCommonDhcpOptions(&import1.VpcDhcpOptions{})
		_ = got
	})

}

func TestUnitFlattenExternalSubnets(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenExternalSubnets(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenExternalSubnets([]import1.ExternalSubnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenExternalSubnets([]import1.ExternalSubnet{import1.ExternalSubnet{}})
		_ = got
	})

}

func TestUnitFlattenActiveGatewayNode(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenActiveGatewayNode(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenActiveGatewayNode([]import1.GatewayNodeReference{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenActiveGatewayNode([]import1.GatewayNodeReference{import1.GatewayNodeReference{}})
		_ = got
	})

}

func TestUnitFlattenNodeIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenNodeIPAddress(&config.IPAddress{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenNodeIPAddress(&config.IPAddress{})
		_ = got
	})

}

func TestUnitFlattenExternallyRoutablePrefixes(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenExternallyRoutablePrefixes(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenExternallyRoutablePrefixes([]import1.IPSubnet{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenExternallyRoutablePrefixes([]import1.IPSubnet{import1.IPSubnet{}})
		_ = got
	})

}

func TestUnitFlattenIPUsage(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPUsage(&import1.IPUsage{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPUsage(&import1.IPUsage{})
		_ = got
	})

}

func TestUnitFlattenIPPoolUsages(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIPPoolUsages(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPPoolUsages([]import1.IPPoolUsage{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPPoolUsages([]import1.IPPoolUsage{import1.IPPoolUsage{}})
		_ = got
	})

}

func TestUnitFlattenIPv4Pool(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Pool(&import1.IPv4Pool{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Pool(&import1.IPv4Pool{})
		_ = got
	})

}

func TestUnitFlattenMigrationState(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenMigrationState(func() *import1.MigrationState { var v import1.MigrationState; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenMigrationState(func() *import1.MigrationState { v := import1.MigrationState(1); return &v }())
		_ = got
	})

}
