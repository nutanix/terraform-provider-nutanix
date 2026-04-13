package clusters

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/api"
	cluster "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	ClusterEntityAPI     *api.ClustersServiceApi
	StorageContainersAPI *api.StorageContainersServiceApi
	PasswordManagerAPI   *api.PasswordManagerServiceApi
	ClusterProfilesAPI   *api.ClusterProfilesServiceApi
	SSLCertificateAPI    *api.SSLCertificateServiceApi
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

	f := &Client{
		ClusterEntityAPI:     api.NewClustersServiceApi(baseClient),
		StorageContainersAPI: api.NewStorageContainersServiceApi(baseClient),
		PasswordManagerAPI:   api.NewPasswordManagerServiceApi(baseClient),
		ClusterProfilesAPI:   api.NewClusterProfilesServiceApi(baseClient),
		SSLCertificateAPI:    api.NewSSLCertificateServiceApi(baseClient),
	}

	return f, nil
}
