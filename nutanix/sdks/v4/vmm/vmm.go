package vmm

import (
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/api"
	vmm "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
)

type Client struct {
	ImagesAPIInstance          *api.ImagesApi
	TemplatesAPIInstance       *api.TemplatesApi
	VMAPIInstance              *api.VmApi
	ImagesPlacementAPIInstance *api.ImagePlacementPoliciesApi
	OvasAPIInstance            *api.OvasApi
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
		if credentials.Port != "" {
			if p, err := strconv.Atoi(credentials.Port); err == nil {
				pcClient.Port = p
			}
		}
		pcClient.VerifySSL = false
		pcClient.AllowVersionNegotiation = false
		baseClient = pcClient
	}

	f := &Client{
		ImagesAPIInstance:          api.NewImagesApi(baseClient),
		TemplatesAPIInstance:       api.NewTemplatesApi(baseClient),
		VMAPIInstance:              api.NewVmApi(baseClient),
		ImagesPlacementAPIInstance: api.NewImagePlacementPoliciesApi(baseClient),
		OvasAPIInstance:            api.NewOvasApi(baseClient),
	}

	return f, nil
}
