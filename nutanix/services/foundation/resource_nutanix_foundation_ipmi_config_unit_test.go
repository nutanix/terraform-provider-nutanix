package foundation

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUnitFoundationIPMIConfigSchema(t *testing.T) {
	t.Parallel()

	r := ResourceNutanixFoundationIPMIConfig()
	require.NoError(t, r.InternalValidate(r.Schema, true))
}

func TestUnitGetRequiredErrorMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		attr string
		want string
	}{
		{name: "empty", attr: "", want: "please provide the value of required attribute: "},
		{name: "populated", attr: "ipmi_ip", want: "please provide the value of required attribute: ipmi_ip"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, getRequiredErrorMessage(tt.attr))
		})
	}
}

func TestUnitFoundationIPMIConfigReadAndDelete(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, ResourceNutanixFoundationIPMIConfig().Schema, map[string]interface{}{
		"ipmi_ip": "10.0.0.1",
	})
	require.Nil(t, resourceFoundationIPMIConfigRead(t.Context(), d, nil))
	require.Nil(t, resourceFoundationIPMIConfigDelete(t.Context(), d, nil))
}
