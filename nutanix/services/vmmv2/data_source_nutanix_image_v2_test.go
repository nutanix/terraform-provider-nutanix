package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameImage = "data.nutanix_image_v2.test"

func TestAccV2NutanixImageDatasource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceConfigV2(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameImage, "name", name),
					resource.TestCheckResourceAttr(datasourceNameImage, "type", "ISO_IMAGE"),
					resource.TestCheckResourceAttr(datasourceNameImage, "description", desc),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "placement_policy_status.#"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "cluster_location_ext_ids.#"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "source.#"),
				),
			},
		},
	})
}

func testAccImageDataSourceConfigV2(name, desc string) string {
	return fmt.Sprintf(`
		resource "nutanix_images_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			type = "ISO_IMAGE"
			source{
				url_source{
					url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
				}
			}
		}

		data "nutanix_image_v2" "test"{
			ext_id = resource.nutanix_images_v2.test.id
		}
`, name, desc)
}

// TestAccV2NutanixImageDatasource_WithPlacementPolicy verifies that when an image
// matches an image placement policy (via category filter), the image data source
// returns a populated `placement_policy_status` block whose enum fields
// (`compliance_status`, `enforcement_mode`) are serialized as strings.
//
// This is the regression test for the bug where `flattenImagePlacementStatus`
// assigned the typed SDK enum pointer directly to a `schema.TypeString` field
// and apply failed with:
//
//	Error: placement_policy_status.0.compliance_status: '' expected type 'string',
//	got unconvertible type 'content.ComplianceStatus'
//
// PC will only stamp `placement_policy_status` on an image when BOTH halves of
// the policy match: the image has a category that matches `image_entity_filter`
// AND there is at least one cluster whose categories match `cluster_entity_filter`.
// To satisfy the cluster side from a test (where we cannot create a cluster), we
// associate the test category to an existing AOS cluster via the PC v4 API using
// a `null_resource` + `local-exec` provisioner. The destroy provisioner
// disassociates the category before the resources are torn down.
//
// Required env (already set by the standard acceptance test runner):
//
//	NUTANIX_ENDPOINT, NUTANIX_PORT, NUTANIX_USERNAME, NUTANIX_PASSWORD
func TestAccV2NutanixImageDatasource_WithPlacementPolicy(t *testing.T) {
	r := acctest.RandInt()
	imageName := fmt.Sprintf("test-image-ipp-%d", r)
	catSuffix := fmt.Sprintf("ipp_%d", r)

	const ippResourceName = "nutanix_image_placement_policy_v2.test-ipp"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {Source: "hashicorp/null"},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccImageDataSourceConfigV2WithPlacementPolicy(imageName, catSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceNameImage, "name", imageName),
					resource.TestCheckResourceAttr(datasourceNameImage, "type", "DISK_IMAGE"),
					resource.TestCheckResourceAttrSet(datasourceNameImage, "ext_id"),

					// The image matches the policy's image_entity_filter and the
					// cluster matches its cluster_entity_filter, so PC must stamp
					// at least one placement_policy_status entry on the image.
					resource.TestCheckResourceAttr(datasourceNameImage, "placement_policy_status.#", "1"),

					// The status entry must point back to the policy we created.
					resource.TestCheckResourceAttrPair(
						datasourceNameImage, "placement_policy_status.0.placement_policy_ext_id",
						ippResourceName, "id",
					),

					// Core regression assertion: compliance_status must be a real
					// string enum name (e.g. COMPLIANT / NON_COMPLIANT / PENDING),
					// not the typed-pointer value that the SDK couldn't coerce.
					resource.TestMatchResourceAttr(
						datasourceNameImage,
						"placement_policy_status.0.compliance_status",
						regexp.MustCompile(`^[A-Z_]+$`),
					),
					// enforcement_mode is also flattened via common.FlattenPtrEnum.
					// PC may leave it empty initially, so we only require the key
					// to be present in state.
					resource.TestCheckResourceAttrSet(
						datasourceNameImage, "placement_policy_status.0.enforcement_mode",
					),
				),
			},
		},
	})
}

func testAccImageDataSourceConfigV2WithPlacementPolicy(name string, categoryName string) string {
	categoryKey := fmt.Sprintf("category_key_%s", categoryName)
	categoryValue := fmt.Sprintf("category_value_%s", categoryName)

	return fmt.Sprintf(`
		# Pick any AOS cluster (skip PC) — its only role here is to be the cluster
		# whose categories we modify so the placement policy has a compliant target.
		data "nutanix_clusters_v2" "aos" {
			filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
		}

		resource "nutanix_category_v2" "test-category" {
			key         = "%[2]s"
			value       = "%[3]s"
			description = "category for image placement policy IT test"
		}

		resource "nutanix_cluster_category_associations_v2" "attach_category_to_cluster" {
			cluster_ext_id = data.nutanix_clusters_v2.aos.cluster_entities.0.ext_id
			categories = [nutanix_category_v2.test-category.id]
		}

		resource "nutanix_image_placement_policy_v2" "test-ipp" {
			name           = "ipp-%[1]s"
			description    = "image placement policy for image data source IT test"
			placement_type = "SOFT"

			cluster_entity_filter {
				category_ext_ids = [nutanix_category_v2.test-category.id]
				type             = "CATEGORIES_MATCH_ALL"
			}
			image_entity_filter {
				category_ext_ids = [nutanix_category_v2.test-category.id]
				type             = "CATEGORIES_MATCH_ALL"
			}

			lifecycle {
				ignore_changes = [
					cluster_entity_filter,
					image_entity_filter,
				]
			}
		}

		resource "nutanix_images_v2" "test" {
			depends_on = [nutanix_image_placement_policy_v2.test-ipp]

			name             = "%[1]s"
			type             = "DISK_IMAGE"
			category_ext_ids = [nutanix_category_v2.test-category.id]

			source {
				url_source {
					url = "http://archive.ubuntu.com/ubuntu/dists/bionic/main/installer-amd64/current/images/netboot/mini.iso"
				}
			}

			# Give PC a moment to evaluate the placement policy against the new
			# image so placement_policy_status is populated before the data
			# source read happens.
			provisioner "local-exec" {
				command = "sleep 15"
			}
		}

		data "nutanix_image_v2" "test" {
			ext_id = nutanix_images_v2.test.id
		}
	`, name, categoryKey, categoryValue)
}
