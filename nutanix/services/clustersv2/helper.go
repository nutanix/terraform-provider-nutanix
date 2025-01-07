package clustersv2

import (
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
