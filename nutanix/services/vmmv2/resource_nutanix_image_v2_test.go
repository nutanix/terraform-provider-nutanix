package vmmv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameImage = "nutanix_images_v2.test"

func TestAccV2NutanixImagesResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagesV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "ISO_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "placement_policy_status.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesResource_WithUpdate(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	updatedName := fmt.Sprintf("test-image-updated-%d", r)
	desc := "test image description"
	updatedDesc := "test image description updated"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagesV2Config(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "ISO_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "placement_policy_status.#"),
				),
			},
			{
				Config: testImagesV2Config(updatedName, updatedDesc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "ISO_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", updatedDesc),
					resource.TestCheckResourceAttrSet(resourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "placement_policy_status.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesResource_WithDisk(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagesV2ConfigWithDisk(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "DISK_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "placement_policy_status.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesResource_WithVMDiskSource(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagesV2ConfigWithVMDiskSource(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "DISK_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "placement_policy_status.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesResource_WithClusterExts(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-%d", r)
	desc := "test image description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagesV2ConfigWithDisk(name, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "DISK_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", desc),
					resource.TestCheckResourceAttrSet(resourceNameImage, "create_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "last_update_time"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "owner_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "size_bytes"),
					resource.TestCheckResourceAttrSet(resourceNameImage, "placement_policy_status.#"),
				),
			},
		},
	})
}

func TestAccV2NutanixImagesResource_WithMoreThanOneSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testImagesV2ConfigWithMoreThanOneSource(),
				ExpectError: regexp.MustCompile("only one of url_source, vm_disk_source, or object_lite_source can be specified in source"),
			},
		},
	})
}

func TestAccV2NutanixImagesResource_WithChecksum(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("test-image-checksum-%d", r)
	desc := "test image checksum description"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testImagesV2ConfigWithChecksum(name, desc, testVars.Images.ISOImageSHA1, "sha1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "ISO_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", desc),
					resource.TestCheckResourceAttr(resourceNameImage, "checksum.#", "1"),
					resource.TestCheckResourceAttr(resourceNameImage, "checksum.0.hex_digest", testVars.Images.ISOImageSHA1),
					resource.TestCheckResourceAttr(resourceNameImage, "checksum.0.object_type", "sha1"),
				),
			},
			{
				// A checksum is verified/stored at image-create time and cannot be
				// changed by an in-place update (the platform keeps the create-time
				// digest). Taint the image so this step destroys and recreates it,
				// exercising a real sha256 create+verify instead of a no-op update.
				Taint:  []string{resourceNameImage},
				Config: testImagesV2ConfigWithChecksum(name, desc, testVars.Images.ISOImageSHA256, "sha256"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameImage, "name", name),
					resource.TestCheckResourceAttr(resourceNameImage, "type", "ISO_IMAGE"),
					resource.TestCheckResourceAttr(resourceNameImage, "description", desc),
					resource.TestCheckResourceAttr(resourceNameImage, "checksum.#", "1"),
					resource.TestCheckResourceAttr(resourceNameImage, "checksum.0.hex_digest", testVars.Images.ISOImageSHA256),
					resource.TestCheckResourceAttr(resourceNameImage, "checksum.0.object_type", "sha256"),
				),
			},
		},
	})
}

func testImagesV2Config(name, desc string) string {
	return fmt.Sprintf(`
		locals {
			config = jsondecode(file("%[3]s"))
		}
		resource "nutanix_images_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			type = "ISO_IMAGE"
			source{
				url_source{
					url = local.config.images.iso_image_url
				}
			}
		}
`, name, desc, filepath)
}

func testImagesV2ConfigWithChecksum(name, desc, hexDigest, objectType string) string {
	return fmt.Sprintf(`
		locals {
			config = jsondecode(file("%[5]s"))
		}
		resource "nutanix_images_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			type = "ISO_IMAGE"
			checksum {
				hex_digest = "%[3]s"
				object_type = "%[4]s"
			}
			source{
				url_source{
					url = local.config.images.iso_image_url
				}
			}
		}
`, name, desc, hexDigest, objectType, filepath)
}

func testImagesV2ConfigWithDisk(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
		config = jsondecode(file("%[3]s"))
		cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_images_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			type = "DISK_IMAGE"
			source{
				url_source{
					url = local.config.images.iso_image_url
				}
			}
			cluster_location_ext_ids = [
				local.cluster0
			]
		}
`, name, desc, filepath)
}

func testImagesV2ConfigWithVMDiskSource(name, desc string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters_v2" "clusters" {}

		locals {
			cluster0 = [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
			config = jsondecode(file("%[3]s"))
			vmm = local.config.vmm
		}

		data "nutanix_storage_containers_v2" "ngt-sc" {
		  filter = "clusterExtId eq '${local.cluster0}' and startswith(name,'default-container-')"
		  limit = 1
		}

		resource "nutanix_virtual_machine_v2" "test"{
			name= "tf-test-vm-disk"
			description =  "desc vm"
			num_cores_per_socket = 1
			num_sockets = 1
			cluster {
				ext_id = local.cluster0
			}
			disks{
				disk_address{
					bus_type = "SCSI"
					index = 0
				}
				backing_info{
					vm_disk{
						disk_size_bytes = 1073741824
						storage_container{
							ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
						}
					}
				}
			}
			power_state = "OFF"
		}

		resource "nutanix_images_v2" "test" {
			name = "%[1]s"
			description = "%[2]s"
			type = "DISK_IMAGE"
			source{
				vm_disk_source{
					ext_id = resource.nutanix_virtual_machine_v2.test.disks.0.ext_id
				}
			}
			cluster_location_ext_ids = [
				local.cluster0
			]
			depends_on = [nutanix_virtual_machine_v2.test]
		}
`, name, desc, filepath)
}

func testImagesV2ConfigWithMoreThanOneSource() string {
	return `
resource "nutanix_images_v2" "test" {
	name = "img-with-two sources"
	description = "%[2]s"
	type = "DISK_IMAGE"
	source{
		url_source{
			url = "http://invalid-url.com"
		}
		vm_disk_source{
			ext_id = "796cef72-ceb9-4d23-9146-af16eec1345f"
		}
	}
	cluster_location_ext_ids = [
		"6eebcfc0-acdc-4c2c-a367-496df04acaea"
	]
}
`
}
