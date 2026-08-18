package multidomain

import (
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/api"
	multidomainClient "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type Client struct {
	Projects          *api.ProjectsServiceApi
	ResourceGroups    *api.ResourceGroupsServiceApi
	APIClientInstance *multidomainClient.ApiClient
}

func NewMultidomainClient(credentials client.Credentials) (*Client, error) {
	var baseClient *multidomainClient.ApiClient

	pcClient := multidomainClient.NewApiClient()
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
		EntityAPIInstance:           api.NewEntitiesServiceApi(baseClient),
		APIClientInstance:           baseClient,
	}, nil

}
