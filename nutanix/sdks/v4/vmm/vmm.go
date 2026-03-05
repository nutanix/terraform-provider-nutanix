package vmm

import (
	"strconv"

	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/api"
	vmm "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ImagesAPIInstance          *api.ImagesServiceApi
	TemplatesAPIInstance       *api.TemplatesServiceApi
	VMAPIInstance              *api.VmServiceApi
	ImagesPlacementAPIInstance *api.ImagePlacementPoliciesServiceApi
	OvasAPIInstance            *api.OvasServiceApi
}

func NewVmmClient(credentials client.Credentials) (*Client, error) {
	var baseClient *vmm.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := vmm.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = sdkconfig.DefaultPort
		if credentials.Port != "" {
			if p, err := strconv.Atoi(credentials.Port); err == nil {
				pcClient.Port = p
			}
		}
		pcClient.VerifySSL = false
		pcClient.AllowVersionNegotiation = sdkconfig.AllowVersionNegotiation
		baseClient = pcClient
	}

	f := &Client{
		ImagesAPIInstance:          api.NewImagesServiceApi(baseClient),
		TemplatesAPIInstance:       api.NewTemplatesServiceApi(baseClient),
		VMAPIInstance:              api.NewVmServiceApi(baseClient),
		ImagesPlacementAPIInstance: api.NewImagePlacementPoliciesServiceApi(baseClient),
		OvasAPIInstance:            api.NewOvasServiceApi(baseClient),
	}

	return f, nil
}
