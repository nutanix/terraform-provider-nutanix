package service_provider

import (
	"github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/api"
	serviceProviderClient "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/sdkconfig"
)

type ServiceProviderManagementClient struct {
	TenantMembershipAPIInstance *api.MembershipServiceApi
	TenantManagementAPIInstance *api.TenantsServiceApi
	APIClientInstance           *serviceProviderClient.ApiClient
}

func NewServiceProviderManagementClient(credentials client.Credentials) (*ServiceProviderManagementClient, error) {
	var baseClient *serviceProviderClient.ApiClient

	pcClient := serviceProviderClient.NewApiClient()

	if cfg := sdkconfig.ConfigureV4ClientForSP(credentials, pcClient); cfg != nil {
		pcClient.Host = cfg.Host
		pcClient.Port = cfg.Port
		pcClient.Username = cfg.Username
		pcClient.Password = cfg.Password
		pcClient.VerifySSL = cfg.VerifySSL
		pcClient.AllowVersionNegotiation = false
		baseClient = pcClient
	}

	f := &ServiceProviderManagementClient{
		TenantMembershipAPIInstance: api.NewMembershipServiceApi(baseClient),
		TenantManagementAPIInstance: api.NewTenantsServiceApi(baseClient),
		APIClientInstance:           baseClient,
	}

	return f, nil
}
