package iam

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/api"
	iam "github.com/nutanix-core/ntnx-api-golang-sdk-internal/iam-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	APIClientInstance           *iam.ApiClient
	DirectoryServiceAPIInstance *api.DirectoryServicesServiceApi
	SamlIdentityAPIInstance     *api.SAMLIdentityProvidersServiceApi
	UsersAPIInstance            *api.UsersServiceApi
	UserGroupsAPIInstance       *api.UserGroupsServiceApi
	RolesAPIInstance            *api.RolesServiceApi
	OperationsAPIInstance       *api.OperationsServiceApi
	AuthAPIInstance             *api.AuthorizationPoliciesServiceApi
	RoleMembershipAPIInstance   *api.RoleMembershipServiceApi
	ScopeTemplatesAPIInstance   *api.ScopeTemplatesServiceApi
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

	return &Client{
		DirectoryServiceAPIInstance: api.NewDirectoryServicesServiceApi(baseClient),
		SamlIdentityAPIInstance:     api.NewSAMLIdentityProvidersServiceApi(baseClient),
		UserGroupsAPIInstance:       api.NewUserGroupsServiceApi(baseClient),
		RolesAPIInstance:            api.NewRolesServiceApi(baseClient),
		OperationsAPIInstance:       api.NewOperationsServiceApi(baseClient),
		UsersAPIInstance:            api.NewUsersServiceApi(baseClient),
		AuthAPIInstance:             api.NewAuthorizationPoliciesServiceApi(baseClient),
		RoleMembershipAPIInstance:   api.NewRoleMembershipServiceApi(baseClient),
		ScopeTemplatesAPIInstance:   api.NewScopeTemplatesServiceApi(baseClient),
		APIClientInstance:           iam.NewApiClient(),
	}, nil
}
