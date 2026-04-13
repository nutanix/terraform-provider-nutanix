package volumes

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v17/api"
	volumes "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	VolumeAPIInstance      *api.VolumeGroupsServiceApi
	IscsiClientAPIInstance *api.IscsiClientsServiceApi
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

	f := &Client{
		VolumeAPIInstance:      api.NewVolumeGroupsServiceApi(baseClient),
		IscsiClientAPIInstance: api.NewIscsiClientsServiceApi(baseClient),
	}

	return f, nil
}
