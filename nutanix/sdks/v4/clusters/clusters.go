package clusters

import (
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/api"
	cluster "github.com/nutanix/ntnx-api-golang-clients/clustermgmt-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
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

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := cluster.NewApiClient()

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
		// pcClient.AllowVersionNegotiation = false
		baseClient = pcClient
	}

	f := &Client{
		ClusterEntityAPI:     api.NewClustersApi(baseClient),
		StorageContainersAPI: api.NewStorageContainersApi(baseClient),
		PasswordManagerAPI:   api.NewPasswordManagerApi(baseClient),
		ClusterProfilesAPI:   api.NewClusterProfilesApi(baseClient),
		SSLCertificateAPI:    api.NewSSLCertificateApi(baseClient),
	}

	return f, nil
}
