package vmmv2

import (
	"testing"

	"github.com/stretchr/testify/require"

	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/content"
)

func TestUnitFlattenImageType(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenImageType(func() *import5.ImageType { var v import5.ImageType; return &v }())
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenImageType(func() *import5.ImageType { v := import5.ImageType(1); return &v }())
		_ = got
	})

}

func TestUnitFlattenOneOfImageChecksum(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfImageChecksum(nil))
	})

}

func TestUnitFlattenOneOfImageSource(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenOneOfImageSource(nil))
	})
}

func TestUnitFlattenURLBasicAuth(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenURLBasicAuth(&import5.UrlBasicAuth{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenURLBasicAuth(&import5.UrlBasicAuth{})
		_ = got
	})

}

func TestUnitFlattenImagePlacementStatus(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenImagePlacementStatus(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenImagePlacementStatus([]import5.ImagePlacementStatus{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenImagePlacementStatus([]import5.ImagePlacementStatus{import5.ImagePlacementStatus{}})
		_ = got
	})

}
