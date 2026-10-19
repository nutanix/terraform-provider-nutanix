package networkingv2

import (
	"testing"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/networking-go-client/v17/models/networking/v4/config"
)

func TestUnitFlattenRoutes(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_ = flattenRoutes(nil)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		got := flattenRoutes([]config.Route{})
		_ = got
	})

	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenRoutes([]config.Route{config.Route{}})
		_ = got
	})

}
