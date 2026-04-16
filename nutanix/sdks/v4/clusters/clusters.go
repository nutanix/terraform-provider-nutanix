package clusters

import (
	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/api"
	cluster "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ClusterEntityAPI     *api.ClustersApi
	StorageContainersAPI *api.StorageContainersApi
	PasswordManagerAPI   *api.PasswordManagerApi
	ClusterProfilesAPI   *api.ClusterProfilesApi
	SSLCertificateAPI    *api.SSLCertificateApi
}

func NewClustersClient(credentials client.Credentials) (*Client, error) {
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
		ClusterEntityAPI:     api.NewClustersApi(baseClient),
		StorageContainersAPI: api.NewStorageContainersApi(baseClient),
		PasswordManagerAPI:   api.NewPasswordManagerApi(baseClient),
		ClusterProfilesAPI:   api.NewClusterProfilesApi(baseClient),
		SSLCertificateAPI:    api.NewSSLCertificateApi(baseClient),
	}, nil
}
