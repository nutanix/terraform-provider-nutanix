package iam

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v16/api"
	iam "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v16/client"

	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
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
		pcClient.Port = 9440
		pcClient.VerifySSL = false

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
