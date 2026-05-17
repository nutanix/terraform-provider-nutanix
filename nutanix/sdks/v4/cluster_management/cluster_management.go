package cluster_management

import (
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/api"
	cluster "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ClustersServiceAPI *api.ClustersServiceApi
}

func NewClusterManagementClient(credentials client.Credentials) (*Client, error) {
	var baseClient *cluster.ApiClient

	pcClient := cluster.NewApiClient()
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
		ClustersServiceAPI: api.NewClustersServiceApi(baseClient),
	}, nil
}
