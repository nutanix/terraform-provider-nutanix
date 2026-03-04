package clusters

import (
	"strconv"

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

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := cluster.NewApiClient()

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
		ClusterEntityAPI:     api.NewClustersServiceApi(baseClient),
		StorageContainersAPI: api.NewStorageContainersServiceApi(baseClient),
		PasswordManagerAPI:   api.NewPasswordManagerServiceApi(baseClient),
		ClusterProfilesAPI:   api.NewClusterProfilesServiceApi(baseClient),
		SSLCertificateAPI:    api.NewSSLCertificateServiceApi(baseClient),
	}

	return f, nil
}
