package networkingv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	config "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/config"
	import2 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/common/v1/response"
	import1 "github.com/nutanix/ntnx-api-golang-clients/microseg-go-client/v4/models/microseg/v4/config"
)

func TestUnitFlattenLinksMicroSeg(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenLinksMicroSeg(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenLinksMicroSeg([]import2.ApiLink{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenLinksMicroSeg([]import2.ApiLink{import2.ApiLink{}})
		_ = got
	})

}

func TestUnitFlattenListofString(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenListofString(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenListofString([]string{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenListofString([]string{"a"})
		_ = got
	})

}

func TestUnitFlattenIPv4Range(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIPv4Range(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Range([]import1.IPv4Range{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4Range([]import1.IPv4Range{import1.IPv4Range{}})
		_ = got
	})

}

func TestUnitFlattenIPv4AddressMicroSeg(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenIPv4AddressMicroSeg(nil))
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressMicroSeg([]config.IPv4Address{})
		require.Nil(t, got)
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenIPv4AddressMicroSeg([]config.IPv4Address{config.IPv4Address{}})
		_ = got
	})

}
