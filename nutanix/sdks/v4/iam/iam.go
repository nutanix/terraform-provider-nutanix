package iam

import (
	"strconv"

	"github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/api"
	iam "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	APIClientInstance           *iam.ApiClient
	DirectoryServiceAPIInstance *api.DirectoryServicesApi
	SamlIdentityAPIInstance     *api.SAMLIdentityProvidersApi
	UsersAPIInstance            *api.UsersApi
	UserGroupsAPIInstance       *api.UserGroupsApi
	RolesAPIInstance            *api.RolesApi
	OperationsAPIInstance       *api.OperationsApi
	AuthAPIInstance             *api.AuthorizationPoliciesApi
}

func NewIamClient(credentials client.Credentials) (*Client, error) {
	var baseClient *iam.ApiClient

	// check if all required fields are present. Else create an empty client
	if credentials.Username != "" && credentials.Password != "" && credentials.Endpoint != "" {
		pcClient := iam.NewApiClient()

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
		DirectoryServiceAPIInstance: api.NewDirectoryServicesApi(baseClient),
		SamlIdentityAPIInstance:     api.NewSAMLIdentityProvidersApi(baseClient),
		UserGroupsAPIInstance:       api.NewUserGroupsApi(baseClient),
		RolesAPIInstance:            api.NewRolesApi(baseClient),
		OperationsAPIInstance:       api.NewOperationsApi(baseClient),
		UsersAPIInstance:            api.NewUsersApi(baseClient),
		AuthAPIInstance:             api.NewAuthorizationPoliciesApi(baseClient),
		APIClientInstance:           iam.NewApiClient(),
	}

	return f, nil
}
