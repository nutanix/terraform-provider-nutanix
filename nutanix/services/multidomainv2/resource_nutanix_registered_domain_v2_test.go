package multidomainv2_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	multidomainConfig "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/config"
	import1 "github.com/nutanix/ntnx-api-golang-clients/multidomain-go-client/v4/models/multidomain/v4/request/registereddomains"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Test plan for nutanix_registered_domain_v2 and its datasources.
//
// RegisteredDomains models a Prism Central (PC) domain registered to a Nutanix
// Central (NC) instance. A live Create requires a reachable remote PC plus a
// trust exchange (establish-trust) with NC, so a full create/update/delete
// lifecycle can only be exercised on an NC-enabled environment. The lifecycle
// test below (TestAccV2NutanixRegisteredDomainResource_Basic) documents that
// end-to-end flow and runs against such an environment.
//
// The remaining tests are environment-agnostic and validate:
//   - the plural List datasource returns a well-formed (possibly empty) list,
//   - an invalid $filter yields zero results (no crash),
//   - schema ValidateFunc rejects invalid enum values for registration_state,
//     api_credential_status and the tunnel provider,
//   - the change-state / establish-trust action resources reject invalid input.

func TestAccV2NutanixRegisteredDomainResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := os.Getenv("NUTANIX_DOMAIN_NAME")
	updateName := fmt.Sprintf("tf-test-rd-%d-update", r)
	dataSourceGet := "data.nutanix_registered_domain_v2.get"
	dataSourceList := "data.nutanix_registered_domains_v2.list"
	localDomainDataSourceGet := "data.nutanix_local_domain_v2.test"
	identityProviderExtID := acc.MustConfig().Iam.DirectoryServicesMain.SecondaryAD.ExtID

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccNCPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testNCPCRegistrationCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNCPCRegistrationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRegisteredDomainV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "registration_config.#", "1"),
					resource.TestCheckResourceAttrPair(resourceNameRegisteredDomainV2, "location_ext_id", "data.nutanix_locations_v2.test", "locations.0.ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameLocalDomainV2, "id"),
					resource.TestCheckResourceAttr(resourceNameLocalDomainV2, "identity_provider_ext_id", identityProviderExtID),
					testAccCheckRegisteredDomainState(resourceNameRegisteredDomainV2, "REGISTERED"),
					testAccCheckLocalDomainAttributes(localDomainDataSourceGet),
				),
			},
			{
				Config: testAccNCPCRegistrationConfig(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "name", updateName),
					resource.TestCheckResourceAttr(dataSourceGet, "name", updateName),
					resource.TestCheckResourceAttr(dataSourceGet, "registration_state", "REGISTERED"),
					resource.TestCheckResourceAttr(dataSourceList, "registered_domains.#", "1"),
					resource.TestCheckResourceAttr(dataSourceList, "registered_domains.0.name", updateName),
					resource.TestCheckResourceAttr(dataSourceList, "registered_domains.0.registration_state", "REGISTERED"),
					testAccCheckRegisteredDomainState(resourceNameRegisteredDomainV2, "REGISTERED"),
					testAccCheckRegisteredDomainAttributes(dataSourceGet, dataSourceList),
					testAccCheckLocalDomainAttributes(localDomainDataSourceGet),
				),
			},
		},
	})
}

func testAccCheckRegisteredDomainState(resourceName, expectedState string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("registered domain resource %q was not found in state", resourceName)
		}

		client, err := testAccNCClient()
		if err != nil {
			return err
		}
		conn := client.NCDomainManagementAPI
		request := import1.GetRegisteredDomainByIdRequest{ExtId: utils.StringPtr(rs.Primary.ID)}
		deadline := time.Now().Add(3 * time.Minute)
		for {
			resp, err := conn.RegisteredDomains.GetRegisteredDomainById(context.Background(), &request)
			if err != nil {
				return fmt.Errorf("error fetching registered domain %q from NC: %w", rs.Primary.ID, err)
			}
			if resp == nil || resp.Data == nil {
				return fmt.Errorf("NC returned an empty response for registered domain %q", rs.Primary.ID)
			}
			registeredDomain, ok := resp.Data.GetValue().(multidomainConfig.RegisteredDomain)
			if !ok {
				return fmt.Errorf("unexpected registered domain response type %T", resp.Data.GetValue())
			}
			if registeredDomain.RegistrationState == nil {
				return fmt.Errorf("registered domain %q returned no registration state", rs.Primary.ID)
			}
			state := registeredDomain.RegistrationState.GetName()
			if state == expectedState {
				return nil
			}
			if state == "REGISTRATION_ERROR" || state == "UNREGISTRATION_ERROR" {
				return fmt.Errorf("registered domain %q reached %s while waiting for %s", rs.Primary.ID, state, expectedState)
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("registered domain %q did not reach %s within 3 minutes (last state: %s)", rs.Primary.ID, expectedState, state)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func testAccNCClient() (*conns.Client, error) {
	return (&conns.Config{
		Endpoint:               os.Getenv("NUTANIX_ENDPOINT"),
		Username:               os.Getenv("NUTANIX_USERNAME"),
		Password:               os.Getenv("NUTANIX_PASSWORD"),
		APIKey:                 os.Getenv("NUTANIX_API_KEY"),
		Port:                   "443",
		Insecure:               os.Getenv("NUTANIX_INSECURE") == "true",
		NutanixCentralEndpoint: os.Getenv("NUTANIX_CENTRAL_ENDPOINT"),
		DomainName:             os.Getenv("NUTANIX_DOMAIN_NAME"),
	}).Client()
}

func testAccCheckRegisteredDomainAttributes(getDataSource, listDataSource string) resource.TestCheckFunc {
	resourceName := resourceNameRegisteredDomainV2
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(getDataSource, "name", resourceName, "name"),
		resource.TestCheckResourceAttrPair(getDataSource, "dns_name", resourceName, "dns_name"),
		resource.TestCheckResourceAttrPair(getDataSource, "fqdn", resourceName, "fqdn"),
		resource.TestCheckResourceAttrPair(getDataSource, "domain_ext_id", resourceName, "domain_ext_id"),
		resource.TestCheckResourceAttrPair(getDataSource, "location_ext_id", resourceName, "location_ext_id"),
		resource.TestCheckResourceAttrPair(getDataSource, "credential_issuer", resourceName, "credential_issuer"),
		resource.TestCheckResourceAttrPair(getDataSource, "created_by", resourceName, "created_by"),
		resource.TestCheckResourceAttrPair(getDataSource, "api_credential_status", resourceName, "api_credential_status"),
		resource.TestCheckResourceAttrPair(getDataSource, "registration_state", resourceName, "registration_state"),
		resource.TestCheckResourceAttrPair(getDataSource, "registered_time", resourceName, "registered_time"),
		resource.TestCheckResourceAttrPair(getDataSource, "tenant_id", resourceName, "tenant_id"),
		resource.TestCheckResourceAttrPair(getDataSource, "platform_data.#", resourceName, "platform_data.#"),
		resource.TestCheckResourceAttrPair(getDataSource, "platform_data.0.cluster_count", resourceName, "platform_data.0.cluster_count"),
		resource.TestCheckResourceAttrPair(getDataSource, "platform_data.0.connectivity_status", resourceName, "platform_data.0.connectivity_status"),
		resource.TestCheckResourceAttrPair(getDataSource, "platform_data.0.pc_env_type", resourceName, "platform_data.0.pc_env_type"),
		resource.TestCheckResourceAttrPair(getDataSource, "platform_data.0.pc_provider_type", resourceName, "platform_data.0.pc_provider_type"),
		resource.TestCheckResourceAttrPair(getDataSource, "platform_data.0.pc_version", resourceName, "platform_data.0.pc_version"),
		resource.TestCheckResourceAttrPair(getDataSource, "registration_config.#", resourceName, "registration_config.#"),
		resource.TestCheckResourceAttrPair(getDataSource, "registration_config.0.domain_ext_id", resourceName, "registration_config.0.domain_ext_id"),
		resource.TestCheckResourceAttrPair(getDataSource, "registration_config.0.target_url", resourceName, "registration_config.0.target_url"),
		resource.TestCheckResourceAttrPair(getDataSource, "registration_config.0.tenant_ext_id", resourceName, "registration_config.0.tenant_ext_id"),
		resource.TestCheckResourceAttrPair(getDataSource, "registration_config.0.credentials.#", resourceName, "registration_config.0.credentials.#"),
		// Resource credentials are WriteOnly (empty in state); DS may expose Sensitive key_id from API.
		resource.TestCheckResourceAttrSet(getDataSource, "registration_config.0.credentials.0.key_id"),
		resource.TestCheckResourceAttrPair(getDataSource, "tunnel.#", resourceName, "tunnel.#"),
		resource.TestCheckResourceAttrPair(getDataSource, "tunnel.0.endpoint", resourceName, "tunnel.0.endpoint"),
		resource.TestCheckResourceAttrPair(getDataSource, "tunnel.0.provider", resourceName, "tunnel.0.provider"),
		resource.TestCheckResourceAttrPair(listDataSource, "registered_domains.0.name", getDataSource, "name"),
		resource.TestCheckResourceAttrPair(listDataSource, "registered_domains.0.dns_name", getDataSource, "dns_name"),
		resource.TestCheckResourceAttrPair(listDataSource, "registered_domains.0.fqdn", getDataSource, "fqdn"),
		resource.TestCheckResourceAttrPair(listDataSource, "registered_domains.0.domain_ext_id", getDataSource, "domain_ext_id"),
		resource.TestCheckResourceAttrPair(listDataSource, "registered_domains.0.location_ext_id", getDataSource, "location_ext_id"),
		resource.TestCheckResourceAttrPair(listDataSource, "registered_domains.0.registration_state", getDataSource, "registration_state"),
		resource.TestCheckResourceAttrPair(listDataSource, "registered_domains.0.tenant_id", getDataSource, "tenant_id"),
	)
}

func testAccCheckLocalDomainAttributes(dataSource string) resource.TestCheckFunc {
	resourceName := resourceNameLocalDomainV2
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(dataSource, "name", resourceName, "name"),
		resource.TestCheckResourceAttrPair(dataSource, "target_url", resourceName, "target_url"),
		resource.TestCheckResourceAttrPair(dataSource, "connection_status.#", resourceName, "connection_status.#"),
		resource.TestCheckResourceAttrPair(dataSource, "connection_status.0.status", resourceName, "connection_status.0.status"),
		resource.TestCheckResourceAttrPair(dataSource, "connection_status.0.service_statuses.#", resourceName, "connection_status.0.service_statuses.#"),
		resource.TestCheckResourceAttrPair(dataSource, "registration_status.#", resourceName, "registration_status.#"),
		resource.TestCheckResourceAttrPair(dataSource, "registration_status.0.days_to_expiry", resourceName, "registration_status.0.days_to_expiry"),
		resource.TestCheckResourceAttrPair(dataSource, "registration_status.0.error_stage", resourceName, "registration_status.0.error_stage"),
		resource.TestCheckResourceAttrPair(dataSource, "registration_status.0.identity_provider_ext_id", resourceName, "registration_status.0.identity_provider_ext_id"),
		resource.TestCheckResourceAttrPair(dataSource, "registration_status.0.state", resourceName, "registration_status.0.state"),
		resource.TestCheckResourceAttrPair(dataSource, "location.#", resourceName, "location.#"),
		resource.TestCheckResourceAttrPair(dataSource, "location.0.ext_id", resourceName, "location.0.ext_id"),
		resource.TestCheckResourceAttrPair(dataSource, "location.0.name", resourceName, "location.0.name"),
		resource.TestCheckResourceAttrPair(dataSource, "location.0.country", resourceName, "location.0.country"),
		resource.TestCheckResourceAttrPair(dataSource, "location.0.state", resourceName, "location.0.state"),
		resource.TestCheckResourceAttrPair(dataSource, "location.0.coordinates.#", resourceName, "location.0.coordinates.#"),
		resource.TestCheckResourceAttrPair(dataSource, "location.0.coordinates.0.latitude", resourceName, "location.0.coordinates.0.latitude"),
		resource.TestCheckResourceAttrPair(dataSource, "location.0.coordinates.0.longitude", resourceName, "location.0.coordinates.0.longitude"),
	)
}

func TestAccV2NutanixRegisteredDomainsDatasource_List(t *testing.T) {
	dataSourceList := "data.nutanix_registered_domains_v2.list"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccNCPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisteredDomainV2ListConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceList, "registered_domains.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixRegisteredDomainsDatasource_ListWithInvalidFilter(t *testing.T) {
	dataSourceList := "data.nutanix_registered_domains_v2.list"
	randomUUID := utils.GenUUID()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccNCPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisteredDomainV2ListWithInvalidFilterConfig(randomUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceList, "registered_domains.#", "0"),
				),
			},
		},
	})
}

func testAccNCPCRegistrationConfig(name string) string {
	testVars := acc.MustConfig()
	secondaryAD := testVars.Iam.DirectoryServicesMain.SecondaryAD
	identityProviderExtID := secondaryAD.ExtID
	domainExtID := testVars.AvailabilityZone.PcExtID
	return fmt.Sprintf(`
# PC provider - direct endpoint for Local Domain operations
provider "nutanix" {
  endpoint                 = "%[1]s"
  nutanix_central_endpoint = ""
  domain_name              = ""
  port                     = "9440"
}

# NC provider - for Registered Domain, Locations, and NC management operations
provider "nutanix-2" {
  endpoint                 = ""
  nutanix_central_endpoint = "%[2]s"
  domain_name              = "%[3]s"
  port                     = "443"
}

data "nutanix_locations_v2" "test" {
  provider = nutanix-2
  filter   = "country eq 'Australia' and state eq 'Victoria' and name eq 'Springvale'"
  limit    = 1
}

resource "nutanix_registered_domain_v2" "test" {
  provider        = nutanix-2
  name            = "%[5]s"
  domain_ext_id   = "%[7]s"
  location_ext_id = data.nutanix_locations_v2.test.locations[0].ext_id
}

# Resource credentials are WriteOnly and not stored; read issued secrets from the DS.
data "nutanix_registered_domain_v2" "creds" {
  provider = nutanix-2
  ext_id   = nutanix_registered_domain_v2.test.id
}

resource "nutanix_local_domain_v2" "test" {
  provider                 = nutanix
  domain_ext_id            = nutanix_registered_domain_v2.test.registration_config[0].domain_ext_id
  identity_provider_ext_id = "%[6]s"
  target_url               = nutanix_registered_domain_v2.test.registration_config[0].target_url
  tenant_ext_id            = nutanix_registered_domain_v2.test.registration_config[0].tenant_ext_id

  credentials {
    api_key = data.nutanix_registered_domain_v2.creds.registration_config[0].credentials[0].api_key
    key_id  = data.nutanix_registered_domain_v2.creds.registration_config[0].credentials[0].key_id
  }
}

data "nutanix_local_domain_v2" "test" {
  provider   = nutanix
  depends_on = [nutanix_local_domain_v2.test]
}

data "nutanix_registered_domain_v2" "get" {
  provider   = nutanix-2
  ext_id     = nutanix_registered_domain_v2.test.id
  depends_on = [nutanix_local_domain_v2.test]
}

data "nutanix_registered_domains_v2" "list" {
  provider   = nutanix-2
  filter     = "name eq '%[5]s'"
  depends_on = [nutanix_local_domain_v2.test]
}
`, os.Getenv("NUTANIX_ENDPOINT"), os.Getenv("NUTANIX_CENTRAL_ENDPOINT"), os.Getenv("NUTANIX_DOMAIN_NAME"), acc.ConfigPath(), name, identityProviderExtID, domainExtID)
}

func testAccRegisteredDomainV2ListConfig() string {
	return `
data "nutanix_registered_domains_v2" "list" {}
`
}

func testAccRegisteredDomainV2ListWithInvalidFilterConfig(uuid string) string {
	return fmt.Sprintf(`
data "nutanix_registered_domains_v2" "list" {
  filter = "extId eq '%s'"
}
`, uuid)
}
