package vmm

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/api"
	vmm "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v16/client"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	ImagesAPIInstance          *api.ImagesApi
	TemplatesAPIInstance       *api.TemplatesApi
	VMAPIInstance              *api.VmApi
	ImagesPlacementAPIInstance *api.ImagePlacementPoliciesApi
}

func NewVmmClient(credentials client.Credentials) (*Client, error) {
	var baseClient *vmm.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := vmm.NewApiClient()

		pcClient.Host = credentials.Endpoint
		pcClient.Password = credentials.Password
		pcClient.Username = credentials.Username
		pcClient.Port = 9440
		pcClient.VerifySSL = false

		baseClient = pcClient
	}

	f := &Client{
		ImagesAPIInstance:          api.NewImagesApi(baseClient),
		TemplatesAPIInstance:       api.NewTemplatesApi(baseClient),
		VMAPIInstance:              api.NewVmApi(baseClient),
		ImagesPlacementAPIInstance: api.NewImagePlacementPoliciesApi(baseClient),
	}

	return f, nil
}
