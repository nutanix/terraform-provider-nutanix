package prism_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

func TestAccNutanixProjectDataSourceByID_basic(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"

	updateName := acctest.RandomWithPrefix("test-project-name-dou")
	updateDescription := acctest.RandomWithPrefix("test-project-desc-dou")
	updateCategoryName := "Environment"
	updateCategoryVal := "Production"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceByIDConfig(subnetName, name, description, categoryName, categoryVal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
			{
				Config: testAccProjectDataSourceByIDConfig(subnetName, updateName, updateDescription, updateCategoryName, updateCategoryVal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
		},
	})
}

func TestAccNutanixProjectDataSourceByName_basic(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"

	updateName := acctest.RandomWithPrefix("test-project-name-dou")
	updateDescription := acctest.RandomWithPrefix("test-project-desc-dou")
	updateCategoryName := "Environment"
	updateCategoryVal := "Production"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckNutanixProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectDataSourceByNameConfig(subnetName, name, description, categoryName, categoryVal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
			{
				Config: testAccProjectDataSourceByNameConfig(subnetName, updateName, updateDescription, updateCategoryName, updateCategoryVal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
		},
	})
}

func testAccProjectDataSourceByIDConfig(subnetName, name, description, categoryName, categoryVal string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "subnet" {
			cluster_uuid       = local.cluster1
			name               = "%s"
			description        = "Description of my unit test VLAN"
			vlan_id            = 31
			subnet_type        = "VLAN"
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24

			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}

			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		}

		resource "nutanix_project" "project_test" {
			name        = "%s"
			description = "%s"

			categories {
				name  = "%s"
				value = "%s"
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			subnet_reference_list {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			api_version = "3.1"
		}

		data "nutanix_project" "test" {
			project_id = nutanix_project.project_test.id
		}
	`, subnetName, name, description, categoryName, categoryVal)
}

func testAccProjectDataSourceByNameConfig(subnetName, name, description, categoryName, categoryVal string) string {
	return fmt.Sprintf(`
		data "nutanix_clusters" "clusters" {}

		locals {
			cluster1 = [
				for cluster in data.nutanix_clusters.clusters.entities :
				cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_subnet" "subnet" {
			cluster_uuid       = local.cluster1
			name               = "%s"
			description        = "Description of my unit test VLAN"
			vlan_id            = 31
			subnet_type        = "VLAN"
			subnet_ip          = "10.250.140.0"
			default_gateway_ip = "10.250.140.1"
			prefix_length      = 24

			dhcp_options = {
				boot_file_name   = "bootfile"
				domain_name      = "nutanix"
				tftp_server_name = "10.250.140.200"
			}

			dhcp_domain_name_server_list = ["8.8.8.8", "4.2.2.2"]
			dhcp_domain_search_list      = ["terraform.nutanix.com", "terraform.unit.test.com"]
		}

		resource "nutanix_project" "project_test" {
			name        = "%s"
			description = "%s"

			categories {
				name  = "%s"
				value = "%s"
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			subnet_reference_list {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			api_version = "3.1"
		}

		data "nutanix_project" "test" {
			project_name = nutanix_project.project_test.name
		}
	`, subnetName, name, description, categoryName, categoryVal)
}
