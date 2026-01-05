package volumes

import (
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/api"
	prism "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	VolumeAPIInstance      *api.VolumeGroupsApi
	IscsiClientAPIInstance *api.IscsiClientsApi
}

func NewVolumeClient(credentials client.Credentials) (*Client, error) {
	var baseClient *prism.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := prism.NewApiClient()

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
		VolumeAPIInstance:      api.NewVolumeGroupsApi(baseClient),
		IscsiClientAPIInstance: api.NewIscsiClientsApi(baseClient),
	}

	return f, nil
}
