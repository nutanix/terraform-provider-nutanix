package clustersv2_test

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const dataSourceNameSSLCertificate = "data.nutanix_ssl_certificate_v2.test"
const resourceNameSSLCertificate = "nutanix_ssl_certificate_v2.test"

func TestAccV2NutanixSSLCertificateV2_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// update the ssl certificate
			{
				Config: testAccSSLCertificateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSSLCertificate, "public_certificate", testVars.Clusters.SSLCertificate.PublicCertificate),
					resource.TestCheckResourceAttr(resourceNameSSLCertificate, "ca_chain", testVars.Clusters.SSLCertificate.CaChain),
					resource.TestCheckResourceAttr(resourceNameSSLCertificate, "private_key_algorithm", "RSA_2048"),
				),
			},
			// read the ssl certificate
			{
				Config: testAccSSLCertificateConfig() + `
				
				data "nutanix_ssl_certificate_v2" "test" {
					cluster_ext_id = local.clusterUUID
					depends_on = [nutanix_ssl_certificate_v2.test]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					// Decode base64 certificate from config and compare with formatted certificate from API
					checkPublicCertificateDecoded(dataSourceNameSSLCertificate, "public_certificate", testVars.Clusters.SSLCertificate.PublicCertificate),
					resource.TestCheckResourceAttr(dataSourceNameSSLCertificate, "private_key_algorithm", "RSA_2048"),
				),
			},
		},
	})
}

func TestAccV2NutanixSSLCertificateV2_Regenerate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			// regenerate the ssl certificate
			{
				Config: `
				data "nutanix_clusters_v2" "clusters" {
					filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
				}
				locals {
					clusterUUID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
				}
				# regenerate the ssl certificate
				resource "nutanix_ssl_certificate_v2" "test" {
					cluster_ext_id = local.clusterUUID
					private_key_algorithm = "RSA_2048"
				}
				`,
			},
			// read the ssl certificate
			{
				Config: `
				data "nutanix_clusters_v2" "clusters" {
					filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
				}
				locals {
					clusterUUID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
				}
				# regenerate the ssl certificate
				resource "nutanix_ssl_certificate_v2" "test" {
					cluster_ext_id = local.clusterUUID
					private_key_algorithm = "RSA_2048"
				}
				# read the ssl certificate
				data "nutanix_ssl_certificate_v2" "test" {
					cluster_ext_id = local.clusterUUID
					depends_on = [nutanix_ssl_certificate_v2.test]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameSSLCertificate, "public_certificate"),
					resource.TestCheckResourceAttr(dataSourceNameSSLCertificate, "private_key_algorithm", "RSA_2048"),
				),
			},
		},
	})
}

func testAccSSLCertificateConfig() string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {
			filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
		}

		locals {
			clusterUUID = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
			config = jsondecode(file("%[1]s"))
			ssl_certificate = local.config.clusters.ssl_certificate
		}
		# update the ssl certificate
		resource "nutanix_ssl_certificate_v2" "test" {
			cluster_ext_id = local.clusterUUID
			passphrase = local.ssl_certificate.passphrase
			private_key = local.ssl_certificate.private_key
			public_certificate = local.ssl_certificate.public_certificate
			ca_chain = local.ssl_certificate.ca_chain
			private_key_algorithm = "RSA_2048"
			lifecycle {
				ignore_changes = [ private_key, passphrase, ca_chain ]
			}
		}
	`, filepath)
}

// checkPublicCertificateDecoded decodes the base64 value from config and compares it with
// the formatted certificate value from the API (resource/data source)
// It normalizes whitespace to handle differences in trailing newlines
func checkPublicCertificateDecoded(resourceName, attrName, expectedBase64 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// Get the actual value from the resource/data source (formatted certificate from API)
		actualValue := rs.Primary.Attributes[attrName]

		// Decode the expected base64 value to get the formatted certificate
		decodedExpected, err := base64.StdEncoding.DecodeString(expectedBase64)
		if err != nil {
			return fmt.Errorf("failed to decode base64 %s: %w", attrName, err)
		}
		expectedValue := string(decodedExpected)

		// Normalize whitespace by trimming all leading/trailing whitespace
		// This handles cases where the API might return certificates with/without trailing newlines
		normalizedExpected := strings.TrimSpace(expectedValue)
		normalizedActual := strings.TrimSpace(actualValue)

		// Compare the normalized values
		if normalizedActual != normalizedExpected {
			return fmt.Errorf("%s.%s: expected decoded value (normalized):\n%q\n\ngot (normalized):\n%q\n\nOriginal expected:\n%s\n\nOriginal got:\n%s",
				resourceName, attrName, normalizedExpected, normalizedActual, expectedValue, actualValue)
		}

		return nil
	}
}
