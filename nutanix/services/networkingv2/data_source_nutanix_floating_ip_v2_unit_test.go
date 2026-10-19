package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/common/v1/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenFloatingIP(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenFloatingIP(&import1.FloatingIPAddress{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenFloatingIP(&import1.FloatingIPAddress{})
		_ = got
	})

}

func TestUnitFlattenFloatingIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenFloatingIPv4Address(&import1.FloatingIPv4Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenFloatingIPv4Address(&import1.FloatingIPv4Address{})
		_ = got
	})

}

func TestUnitFlattenFloatingIPv6Address(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenFloatingIPv6Address(&import1.FloatingIPv6Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenFloatingIPv6Address(&import1.FloatingIPv6Address{})
		_ = got
	})

}

func TestUnitFlattenVpc(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVpc(&import1.Vpc{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVpc(&import1.Vpc{})
		_ = got
	})

}

func TestUnitFlattenVMNic(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenVMNic(&import1.VmNic{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenVMNic(&import1.VmNic{})
		_ = got
	})

}

func TestUnitFlattenExternalSubnet(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenExternalSubnet(&import1.Subnet{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenExternalSubnet(&import1.Subnet{})
		_ = got
	})

}

func TestUnitFlattenAssociation(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenAssociation(nil))
	})
}

func TestUnitFlattenIPAddress(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPAddress(&config.IPAddress{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPAddress(&config.IPAddress{})
		_ = got
	})

}

func TestUnitFlattenIPv4Address(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Address(&config.IPv4Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Address(&config.IPv4Address{})
		_ = got
	})

}

func TestUnitFlattenIPv6Address(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv6Address(&config.IPv6Address{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv6Address(&config.IPv6Address{})
		_ = got
	})

}
