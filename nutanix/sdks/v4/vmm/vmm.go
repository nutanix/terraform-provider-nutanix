package vmm

import (
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/api"
	vmm "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ImagesAPIInstance               *api.ImagesApi
	TemplatesAPIInstance            *api.TemplatesApi
	VMAPIInstance                   *api.VmApi
	ImagesPlacementAPIInstance      *api.ImagePlacementPoliciesApi
	OvasAPIInstance                 *api.OvasApi
	VMAntiAffinityPolicyAPIInstance *api.VmAntiAffinityPoliciesApi
	VMHostAffinityPolicyAPIInstance *api.VmHostAffinityPoliciesApi
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
		ImagesAPIInstance:               api.NewImagesApi(baseClient),
		TemplatesAPIInstance:            api.NewTemplatesApi(baseClient),
		VMAPIInstance:                   api.NewVmApi(baseClient),
		ImagesPlacementAPIInstance:      api.NewImagePlacementPoliciesApi(baseClient),
		OvasAPIInstance:                 api.NewOvasApi(baseClient),
		VMAntiAffinityPolicyAPIInstance: api.NewVmAntiAffinityPoliciesApi(baseClient),
		VMHostAffinityPolicyAPIInstance: api.NewVmHostAffinityPoliciesApi(baseClient),
	}

	return f, nil
}
