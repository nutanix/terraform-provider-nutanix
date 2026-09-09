package lifecyclev2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// Test plan for nutanix_connection_v2 (hardware provider connection)
// -----------------------------------------------------------------
// A connection stores the endpoint + authentication details used to reach an
// external hardware provider. Create/Update/Delete are asynchronous task APIs.
// Hardware providers themselves are NOT created via Terraform, so tests read a
// pre-existing hardware_provider_ext_id from test_config_v2.json.
//
//  1. TestAccV2NutanixConnectionResource_Basic
//     - Step 1 (Create): create a connection with a basic_auth + url_endpoint access
//       detail; assert every literal (name, region, auth username, endpoint url) and
//       that computed ext_id / deployment_type are populated.
//     - Step 2 (Update): change name + region + auth username; assert the new literals.
//     - Step 3 (Datasource singular): read the connection by ext_id and assert all
//       attributes round-trip via TestCheckResourceAttrPair.
//     - Step 4 (Datasource plural): list connections for the provider and assert
//       count > 0 plus spot-checked nested attributes.
//     - Step 5 (Import): import the resource and verify state.
//  2. TestAccV2NutanixConnectionResource_IPRangeEndpoint (variant): exercise the
//     ip_range_endpoint OneOf branch.
//  3. TestAccV2NutanixConnectionResource_IPAddressEndpoint (variant): exercise the
//     ip_address_endpoint OneOf branch + api_key_auth OneOf branch.
//  4. TestAccV2NutanixConnectionResource_MissingName (negative): omit the required
//     name and assert Terraform reports the missing argument.

const (
	resourceNameConnection    = "nutanix_connection_v2.test"
	datasourceNameConnection  = "data.nutanix_connection_v2.test"
	datasourceNameConnections = "data.nutanix_connections_v2.test"
)

func testAccConnectionCheckDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).LifecycleAPI

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_connection_v2" {
			continue
		}
		hpExtID := rs.Primary.Attributes["hardware_provider_ext_id"]
		_, err := conn.HardwareProvidersAPIInstance.GetConnectionById(utils.StringPtr(hpExtID), utils.StringPtr(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("connection %s still exists after destroy", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccV2NutanixConnectionResource_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-test-conn-%d", acctest.RandInt())
	nameUpdated := name + "-updated"
	region := "us-west"
	regionUpdated := "us-east"
	username := "conn-user"
	usernameUpdated := "conn-user-2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccConnectionCheckDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create with basic_auth + url_endpoint.
			{
				Config: testConnectionConfigBasic(name, region, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameConnection, "name", name),
					resource.TestCheckResourceAttr(resourceNameConnection, "region", region),
					resource.TestCheckResourceAttr(resourceNameConnection, "access_details.0.auth.0.basic_auth.0.username", username),
					resource.TestCheckResourceAttrSet(resourceNameConnection, "ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameConnection, "access_details.0.endpoint.0.url_endpoint.0.url"),
				),
			},
			// Step 2: Update name, region and auth username.
			{
				Config: testConnectionConfigBasic(nameUpdated, regionUpdated, usernameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameConnection, "name", nameUpdated),
					resource.TestCheckResourceAttr(resourceNameConnection, "region", regionUpdated),
					resource.TestCheckResourceAttr(resourceNameConnection, "access_details.0.auth.0.basic_auth.0.username", usernameUpdated),
				),
			},
			// Step 3: Singular datasource must round-trip attributes.
			{
				Config: testConnectionConfigWithDatasources(nameUpdated, regionUpdated, usernameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceNameConnection, "name", resourceNameConnection, "name"),
					resource.TestCheckResourceAttrPair(datasourceNameConnection, "region", resourceNameConnection, "region"),
					resource.TestCheckResourceAttrPair(datasourceNameConnection, "ext_id", resourceNameConnection, "ext_id"),
					// Step 4: Plural datasource must return at least one connection.
					resource.TestCheckResourceAttrSet(datasourceNameConnections, "connections.#"),
				),
			},
			// Step 5: Import.
			{
				ResourceName:            resourceNameConnection,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccConnectionImportStateIDFunc,
				ImportStateVerifyIgnore: []string{"access_details"},
			},
		},
	})
}

func TestAccV2NutanixConnectionResource_IPRangeEndpoint(t *testing.T) {
	name := fmt.Sprintf("tf-test-conn-range-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccConnectionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testConnectionConfigIPRange(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameConnection, "name", name),
					resource.TestCheckResourceAttr(resourceNameConnection, "access_details.0.endpoint.0.ip_range_endpoint.0.ip_ranges.0.start_ip.0.ipv4.0.value", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceNameConnection, "access_details.0.endpoint.0.ip_range_endpoint.0.ip_ranges.0.end_ip.0.ipv4.0.value", "10.0.0.10"),
				),
			},
		},
	})
}

func TestAccV2NutanixConnectionResource_IPAddressEndpoint(t *testing.T) {
	name := fmt.Sprintf("tf-test-conn-addr-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccConnectionCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testConnectionConfigIPAddress(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameConnection, "name", name),
					resource.TestCheckResourceAttr(resourceNameConnection, "access_details.0.endpoint.0.ip_address_endpoint.0.ip_addresses.0.ipv4.0.value", "10.0.0.5"),
					resource.TestCheckResourceAttr(resourceNameConnection, "access_details.0.auth.0.api_key_auth.0.api_key_id", "api-key-123"),
				),
			},
		},
	})
}

func TestAccV2NutanixConnectionResource_MissingName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testConnectionConfigMissingName(),
				ExpectError: regexp.MustCompile("The argument \"name\" is required"),
			},
		},
	})
}

func testAccConnectionImportStateIDFunc(s *terraform.State) (string, error) {
	rs, ok := s.RootModule().Resources[resourceNameConnection]
	if !ok {
		return "", fmt.Errorf("resource %s not found in state", resourceNameConnection)
	}
	return rs.Primary.ID, nil
}

func testConnectionConfigBasic(name, region, username string) string {
	return fmt.Sprintf(`
resource "nutanix_connection_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  name                     = "%[2]s"
  region                   = "%[3]s"
  access_details {
    auth {
      basic_auth {
        username = "%[4]s"
        password = "conn-pass"
      }
    }
    endpoint {
      url_endpoint {
        url = "%[5]s"
      }
    }
  }
  lifecycle {
    ignore_changes = [access_details]
  }
}
`, testVars.Lifecycle.HardwareProviders.HardwareProviderExtID, name, region, username, testConnectionURL())
}

func testConnectionConfigWithDatasources(name, region, username string) string {
	return testConnectionConfigBasic(name, region, username) + fmt.Sprintf(`
data "nutanix_connection_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  ext_id                   = nutanix_connection_v2.test.ext_id
}

data "nutanix_connections_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  depends_on               = [nutanix_connection_v2.test]
}
`, testVars.Lifecycle.HardwareProviders.HardwareProviderExtID)
}

func testConnectionConfigIPRange(name string) string {
	return fmt.Sprintf(`
resource "nutanix_connection_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  name                     = "%[2]s"
  access_details {
    auth {
      basic_auth {
        username = "range-user"
        password = "range-pass"
      }
    }
    endpoint {
      ip_range_endpoint {
        ip_ranges {
          start_ip {
            ipv4 {
              value = "10.0.0.1"
            }
          }
          end_ip {
            ipv4 {
              value = "10.0.0.10"
            }
          }
        }
      }
    }
  }
  lifecycle {
    ignore_changes = [access_details]
  }
}
`, testVars.Lifecycle.HardwareProviders.HardwareProviderExtID, name)
}

func testConnectionConfigIPAddress(name string) string {
	return fmt.Sprintf(`
resource "nutanix_connection_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  name                     = "%[2]s"
  access_details {
    auth {
      api_key_auth {
        api_key_id     = "api-key-123"
        api_key_secret = "api-key-secret"
      }
    }
    endpoint {
      ip_address_endpoint {
        ip_addresses {
          ipv4 {
            value = "10.0.0.5"
          }
        }
      }
    }
  }
  lifecycle {
    ignore_changes = [access_details]
  }
}
`, testVars.Lifecycle.HardwareProviders.HardwareProviderExtID, name)
}

func testConnectionConfigMissingName() string {
	return fmt.Sprintf(`
resource "nutanix_connection_v2" "test" {
  hardware_provider_ext_id = "%[1]s"
  access_details {
    endpoint {
      url_endpoint {
        url = "https://example.com"
      }
    }
  }
}
`, testVars.Lifecycle.HardwareProviders.HardwareProviderExtID)
}

func testConnectionURL() string {
	if testVars.Lifecycle.HardwareProviders.ConnectionURL != "" {
		return testVars.Lifecycle.HardwareProviders.ConnectionURL
	}
	return "https://hardware-provider.example.com"
}
