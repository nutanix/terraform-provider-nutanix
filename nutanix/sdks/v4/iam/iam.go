package iam

import (
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
	EntityAPIInstance           *api.EntitiesApi
}

func NewIamClient(credentials client.Credentials) (*Client, error) {
	var baseClient *iam.ApiClient

	pcClient := iam.NewApiClient()
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
		DirectoryServiceAPIInstance: api.NewDirectoryServicesApi(baseClient),
		SamlIdentityAPIInstance:     api.NewSAMLIdentityProvidersApi(baseClient),
		UserGroupsAPIInstance:       api.NewUserGroupsApi(baseClient),
		RolesAPIInstance:            api.NewRolesApi(baseClient),
		OperationsAPIInstance:       api.NewOperationsApi(baseClient),
		UsersAPIInstance:            api.NewUsersApi(baseClient),
		AuthAPIInstance:             api.NewAuthorizationPoliciesApi(baseClient),
		EntityAPIInstance:           api.NewEntitiesApi(baseClient),
		APIClientInstance:           baseClient,
	}

	return f, nil
}
