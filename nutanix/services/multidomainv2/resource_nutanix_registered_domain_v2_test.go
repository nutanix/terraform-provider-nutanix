package multidomainv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	name := fmt.Sprintf("tf-test-rd-%d", r)
	updateName := fmt.Sprintf("tf-test-rd-%d-update", r)
	dataSourceGet := "data.nutanix_registered_domain_v2.get"
	dataSourceList := "data.nutanix_registered_domains_v2.list"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testRegisteredDomainV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisteredDomainV2Config(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRegisteredDomainV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "name", name),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "fqdn", "pc.example.com"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "dns_name", "pc.example.com"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "domain_ext_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "location_ext_id", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "credential_issuer", "nutanix-central"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "created_by", "admin"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "api_credential_status", "VALID"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "registration_state", "REGISTERED"),
					resource.TestCheckResourceAttrSet(resourceNameRegisteredDomainV2, "tenant_id"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "platform_data.#", "1"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "platform_data.0.pc_version", "pc.7.6"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "platform_data.0.connectivity_status", "REACHABLE"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "platform_data.0.cluster_count", "1"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "registration_config.#", "1"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "registration_config.0.target_url", "https://nc.example.com"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "tunnel.#", "1"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "tunnel.0.endpoint", "https://tunnel.example.com"),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "tunnel.0.provider", "IKAT"),
					resource.TestCheckResourceAttrPair(dataSourceGet, "name", resourceNameRegisteredDomainV2, "name"),
					resource.TestCheckResourceAttrPair(dataSourceGet, "fqdn", resourceNameRegisteredDomainV2, "fqdn"),
					resource.TestCheckResourceAttrPair(dataSourceGet, "registration_state", resourceNameRegisteredDomainV2, "registration_state"),
					resource.TestCheckResourceAttrPair(dataSourceGet, "platform_data.0.pc_version", resourceNameRegisteredDomainV2, "platform_data.0.pc_version"),
					resource.TestCheckResourceAttrPair(dataSourceGet, "registration_config.0.target_url", resourceNameRegisteredDomainV2, "registration_config.0.target_url"),
					resource.TestCheckResourceAttrPair(dataSourceGet, "tunnel.0.provider", resourceNameRegisteredDomainV2, "tunnel.0.provider"),
				),
			},
			{
				Config: testAccRegisteredDomainV2UpdateConfig(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "name", updateName),
					resource.TestCheckResourceAttr(resourceNameRegisteredDomainV2, "dns_name", "pc-updated.example.com"),
				),
			},
			{
				Config: testAccRegisteredDomainV2ListConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceList, "registered_domains.#"),
				),
			},
			{
				ResourceName:            resourceNameRegisteredDomainV2,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"registration_config.0.credentials.0.api_key"},
			},
		},
	})
}

func TestAccV2NutanixRegisteredDomainsDatasource_List(t *testing.T) {
	dataSourceList := "data.nutanix_registered_domains_v2.list"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
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
		PreCheck:  func() { acc.TestAccPreCheck(t) },
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

func TestAccV2NutanixRegisteredDomainResource_InvalidRegistrationState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccRegisteredDomainV2InvalidEnumConfig("registration_state", "INVALID_STATE"),
				ExpectError: regexp.MustCompile(`expected registration_state to be one of`),
			},
		},
	})
}

func TestAccV2NutanixRegisteredDomainResource_InvalidApiCredentialStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccRegisteredDomainV2InvalidEnumConfig("api_credential_status", "INVALID_STATUS"),
				ExpectError: regexp.MustCompile(`expected api_credential_status to be one of`),
			},
		},
	})
}

func TestAccV2NutanixRegisteredDomainResource_InvalidTunnelProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccRegisteredDomainV2InvalidTunnelProviderConfig(),
				ExpectError: regexp.MustCompile(`expected tunnel.0.provider to be one of`),
			},
		},
	})
}

func TestAccV2NutanixRegisteredDomainChangeState_InvalidState(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccRegisteredDomainChangeStateInvalidConfig(),
				ExpectError: regexp.MustCompile(`expected registration_state to be one of`),
			},
		},
	})
}

func testAccRegisteredDomainV2Config(name string) string {
	return fmt.Sprintf(`
resource "nutanix_registered_domain_v2" "test" {
  name                  = "%[1]s"
  fqdn                  = "pc.example.com"
  dns_name              = "pc.example.com"
  domain_ext_id         = "11111111-1111-1111-1111-111111111111"
  location_ext_id       = "22222222-2222-2222-2222-222222222222"
  credential_issuer     = "nutanix-central"
  created_by            = "admin"
  api_credential_status = "VALID"
  registration_state    = "REGISTERED"

  platform_data {
    pc_version          = "pc.7.6"
    connectivity_status = "REACHABLE"
    pc_env_type         = "ONPREM"
    pc_provider_type    = "VM"
    cluster_count       = 1
  }

  registration_config {
    domain_ext_id = "11111111-1111-1111-1111-111111111111"
    target_url    = "https://nc.example.com"
    tenant_ext_id = "33333333-3333-3333-3333-333333333333"
    credentials {
      api_key = "example-api-key"
      key_id  = "example-key-id"
    }
  }

  tunnel {
    endpoint = "https://tunnel.example.com"
    provider = "IKAT"
  }
}

data "nutanix_registered_domain_v2" "get" {
  ext_id = nutanix_registered_domain_v2.test.id
}

data "nutanix_registered_domains_v2" "list" {
  depends_on = [nutanix_registered_domain_v2.test]
}
`, name)
}

func testAccRegisteredDomainV2UpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_registered_domain_v2" "test" {
  name               = "%[1]s"
  fqdn               = "pc.example.com"
  dns_name           = "pc-updated.example.com"
  registration_state = "REGISTERED"

  platform_data {
    pc_version          = "pc.7.6"
    connectivity_status = "REACHABLE"
    cluster_count       = 2
  }
}
`, name)
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

func testAccRegisteredDomainV2InvalidEnumConfig(attr, value string) string {
	return fmt.Sprintf(`
resource "nutanix_registered_domain_v2" "test" {
  name = "tf-test-invalid"
  %s   = "%s"
}
`, attr, value)
}

func testAccRegisteredDomainV2InvalidTunnelProviderConfig() string {
	return `
resource "nutanix_registered_domain_v2" "test" {
  name = "tf-test-invalid-tunnel"
  tunnel {
    endpoint = "https://tunnel.example.com"
    provider = "INVALID_PROVIDER"
  }
}
`
}

func testAccRegisteredDomainChangeStateInvalidConfig() string {
	return `
resource "nutanix_registered_domain_change_state_v2" "test" {
  ext_id             = "00000000-0000-0000-0000-000000000000"
  registration_state = "INVALID_STATE"
}
`
}
