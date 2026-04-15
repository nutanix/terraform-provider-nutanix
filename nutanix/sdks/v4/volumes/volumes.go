package volumes

import (
	"github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/api"
	volumes "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	VolumeAPIInstance      *api.VolumeGroupsApi
	IscsiClientAPIInstance *api.IscsiClientsApi
}

func NewVolumeClient(credentials client.Credentials) (*Client, error) {
	var baseClient *volumes.ApiClient

	pcClient := volumes.NewApiClient()
	if cfg := sdkconfig.ConfigureV4Client(credentials, pcClient); cfg != nil {
		pcClient.Host = cfg.Host
		pcClient.Port = cfg.Port
		pcClient.Username = cfg.Username
		pcClient.Password = cfg.Password
		pcClient.VerifySSL = cfg.VerifySSL
		pcClient.AllowVersionNegotiation = cfg.AllowVersionNegotiation
		baseClient = pcClient
	}

	return &Client{
		VolumeAPIInstance:      api.NewVolumeGroupsApi(baseClient),
		IscsiClientAPIInstance: api.NewIscsiClientsApi(baseClient),
	}, nil
}
