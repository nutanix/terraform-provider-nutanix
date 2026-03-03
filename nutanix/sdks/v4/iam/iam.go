package iam

import (
	"strconv"

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
		DirectoryServiceAPIInstance: api.NewDirectoryServicesServiceApi(baseClient),
		SamlIdentityAPIInstance:     api.NewSAMLIdentityProvidersServiceApi(baseClient),
		UserGroupsAPIInstance:       api.NewUserGroupsServiceApi(baseClient),
		RolesAPIInstance:            api.NewRolesServiceApi(baseClient),
		OperationsAPIInstance:       api.NewOperationsServiceApi(baseClient),
		UsersAPIInstance:            api.NewUsersServiceApi(baseClient),
		AuthAPIInstance:             api.NewAuthorizationPoliciesServiceApi(baseClient),
		APIClientInstance:           iam.NewApiClient(),
	}

	return f, nil
}
