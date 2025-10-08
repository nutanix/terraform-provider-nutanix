package clustersv2

import (
	"strings"

	sdk "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/models/clustermgmt/v4/config"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/clusters"

	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func getEtagHeader(resp interface{}, conn *clusters.Client) map[string]interface{} {
	// Extract E-Tag Header
	etagValue := conn.ClusterEntityAPI.ApiClient.GetEtag(resp)

	// Extract E-Tag Header
	args := make(map[string]interface{})
	args["If-Match"] = utils.StringPtr(etagValue)

	return args
}

var allowedConfigTypeNames = buildConfigTypeNames()

// buildConfigTypeNames inspects sdk.ConfigType values and returns a deduped list
// of human-facing enum names (skips names that start with '$').
func buildConfigTypeNames() []string {
	seen := map[string]bool{}
	out := make([]string, 0, 16)

	// safe upper bound to discover enum values; adjust if you expect >256.
	const maxTry = 256
	for i := 0; i < maxTry; i++ {
		name := sdk.ConfigType(i).GetName()
		// skip unknown/internal markers that start with $
		if name == "" || strings.HasPrefix(name, "$") {
			continue
		}
		if !seen[name] {
			seen[name] = true
			out = append(out, name)
		}
	}

	return out
}
