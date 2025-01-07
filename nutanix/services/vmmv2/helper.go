package vmmv2

import (
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/vmm"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func getEtagHeader(resp interface{}, conn *vmm.Client) *string {
	// Extract E-Tag Header
	etagValue := conn.VMAPIInstance.ApiClient.GetEtag(resp)
	return utils.StringPtr(etagValue)
}
