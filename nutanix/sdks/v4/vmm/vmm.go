package vmm

import (
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/api"
	vmm "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ImagesAPIInstance                       *api.ImagesServiceApi
	TemplatesAPIInstance                    *api.TemplatesServiceApi
	VMAPIInstance                           *api.VmServiceApi
	ImagesPlacementAPIInstance              *api.ImagePlacementPoliciesServiceApi
	OvasAPIInstance                         *api.OvasServiceApi
	VMAntiAffinityPolicyAPIInstance         *api.VmAntiAffinityPoliciesServiceApi
	VMHostAffinityPolicyAPIInstance         *api.VmHostAffinityPoliciesServiceApi
	TemplatePlacementPoliciesAPIInstance    *api.TemplatePlacementPoliciesServiceApi
	ImageRateLimitPoliciesAPIInstance       *api.ImageRateLimitPoliciesServiceApi
	VmStartupPoliciesAPIInstance            *api.VmStartupPoliciesServiceApi
	VmGuestCustomizationProfilesAPIInstance *api.VmGuestCustomizationProfilesServiceApi
}

func NewVmmClient(credentials client.Credentials) (*Client, error) {
	var baseClient *vmm.ApiClient

	pcClient := vmm.NewApiClient()
	if cfg := sdkconfig.ConfigureV4Client(credentials, pcClient); cfg != nil {
		pcClient.Host = cfg.Host
		pcClient.Port = cfg.Port
		pcClient.Username = cfg.Username
		pcClient.Password = cfg.Password
		pcClient.VerifySSL = cfg.VerifySSL
		pcClient.AllowVersionNegotiation = cfg.AllowVersionNegotiation
		baseClient = pcClient
	}
	f := &Client{
		ImagesAPIInstance:                       api.NewImagesServiceApi(baseClient),
		TemplatesAPIInstance:                    api.NewTemplatesServiceApi(baseClient),
		VMAPIInstance:                           api.NewVmServiceApi(baseClient),
		ImagesPlacementAPIInstance:              api.NewImagePlacementPoliciesServiceApi(baseClient),
		OvasAPIInstance:                         api.NewOvasServiceApi(baseClient),
		VMAntiAffinityPolicyAPIInstance:         api.NewVmAntiAffinityPoliciesServiceApi(baseClient),
		VMHostAffinityPolicyAPIInstance:         api.NewVmHostAffinityPoliciesServiceApi(baseClient),
		TemplatePlacementPoliciesAPIInstance:    api.NewTemplatePlacementPoliciesServiceApi(baseClient),
		ImageRateLimitPoliciesAPIInstance:       api.NewImageRateLimitPoliciesServiceApi(baseClient),
		VmStartupPoliciesAPIInstance:            api.NewVmStartupPoliciesServiceApi(baseClient),
		VmGuestCustomizationProfilesAPIInstance: api.NewVmGuestCustomizationProfilesServiceApi(baseClient),
	}

	return f, nil
}
