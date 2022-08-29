package nutanix

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/spf13/cast"
)

func TestAccNutanixProject_basic(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"
	limit := cast.ToString(acctest.RandIntRange(2, 4))
	rsType := "STORAGE"

	updateName := acctest.RandomWithPrefix("test-project-name-dou")
	updateDescription := acctest.RandomWithPrefix("test-project-desc-dou")
	updateCategoryName := "Environment"
	updateCategoryVal := "Production"
	updateLimit := cast.ToString(acctest.RandIntRange(4, 8))
	updateRSType := "MEMORY"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectConfig(subnetName, name, description, categoryName, categoryVal, limit, rsType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.limit", limit),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.resource_type", rsType),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
			{
				Config: testAccNutanixProjectConfig(
					subnetName, updateName, updateDescription, updateCategoryName, updateCategoryVal, updateLimit, updateRSType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNutanixProjectExists(&resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.limit", updateLimit),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.resource_type", updateRSType),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
				),
			},
		},
	})
}

func TestAccNutanixProject_importBasic(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"
	limit := cast.ToString(acctest.RandIntRange(2, 4))
	rsType := "STORAGE"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNutanixProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectConfig(subnetName, name, description, categoryName, categoryVal, limit, rsType),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckNutanixProjectImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNutanixProject_withInternal(t *testing.T) {
	resourceName := "nutanix_project.project_test"

	subnetName := acctest.RandomWithPrefix("test-subnateName")
	name := acctest.RandomWithPrefix("test-project-name-dou")
	description := acctest.RandomWithPrefix("test-project-desc-dou")
	categoryName := "Environment"
	categoryVal := "Staging"
	limit := cast.ToString(acctest.RandIntRange(2, 4))
	rsType := "STORAGE"

	// updateName := acctest.RandomWithPrefix("test-project-name-dou")
	// updateDescription := acctest.RandomWithPrefix("test-project-desc-dou")
	// updateCategoryName := "Environment"
	// updateCategoryVal := "Production"
	// updateLimit := cast.ToString(acctest.RandIntRange(4, 8))
	// updateRSType := "MEMORY"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixProjectInternalConfig(subnetName, name, description, categoryName, categoryVal, limit, rsType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.limit", limit),
					resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.resource_type", rsType),
					resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_reference_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cluster_reference_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpc_reference_list.#", "1"),
				),
			},
			// {
			// 	Config: testAccNutanixProjectInternalConfig(
			// 		subnetName, updateName, updateDescription, updateCategoryName, updateCategoryVal, updateLimit, updateRSType),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		resource.TestCheckResourceAttr(resourceName, "name", updateName),
			// 		resource.TestCheckResourceAttr(resourceName, "description", updateDescription),
			// 		resource.TestCheckResourceAttr(resourceName, "resource_domain.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.limit", updateLimit),
			// 		resource.TestCheckResourceAttr(resourceName, "resource_domain.0.resources.0.resource_type", updateRSType),
			// 		resource.TestCheckResourceAttr(resourceName, "categories.#", "1"),
			// 		resource.TestCheckResourceAttr(resourceName, "api_version", "3.1"),
			// 	),
			// },
		},
	})
}

func testAccCheckNutanixProjectImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccCheckNutanixProjectExists(resourceName *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[*resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", *resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckNutanixProjectDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_project" {
			continue
		}
		for {
			_, err := conn.API.V3.GetProject(rs.Primary.ID)
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
					return nil
				}
				return err
			}
			time.Sleep(3000 * time.Millisecond)
		}
	}
	return nil
}

func testAccNutanixProjectConfig(subnetName, name, description, categoryName, categoryVal, limit, rsType string) string {
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

			resource_domain {
				resources {
					limit         = %s
					resource_type = "%s"
				}
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}
			subnet_reference_list{
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			api_version = "3.1"
		}
	`, subnetName, name, description, categoryName, categoryVal, limit, rsType)
}

func testAccNutanixProjectInternalConfig(subnetName, name, description, categoryName, categoryVal, limit, rsType string) string {
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


		resource "nutanix_subnet" "overlay-subnet" {
			cluster_uuid = local.cluster1
			name        = "acctest-subnet-updated"
			description = "Description of my unit test VLAN"
			vlan_id     = 876
			subnet_type = "VLAN"
			subnet_ip          = "10.250.144.0"
		  default_gateway_ip = "10.250.144.1"
		  prefix_length = 24
		  is_external = true
		  ip_config_pool_list_ranges = ["10.250.144.10 10.250.144.20"]
		}


		resource "nutanix_vpc" "acctest-managed" {
			depends_on = [
				resource.nutanix_subnet.overlay-subnet
			]
			name = "acctest-managed-vpc"
		  
			external_subnet_reference_name = [
			  "acctest-subnet-updated"
			]
		  
			common_domain_name_server_ip_list{
					ip = "8.8.8.9"
			}
		  
			externally_routable_prefix_list{
			  ip=  "172.30.0.0"
			  prefix_length= 16
			}
			externally_routable_prefix_list{
				ip=  "172.34.0.0"
				prefix_length= 16
			  }
		  }

		resource "nutanix_project" "project_test" {
			name        = "%s"
			description = "%s"

			categories {
				name  = "%s"
				value = "%s"
			}

			resource_domain {
				resources {
					limit         = %s
					resource_type = "%s"
				}
			}

			default_subnet_reference {
				uuid = nutanix_subnet.subnet.metadata.uuid
			}

			use_project_internal = true

			api_version = "3.1"

			subnet_reference_list{
				kind="subnet"
				uuid=nutanix_subnet.subnet.metadata.uuid
			}
			subnet_reference_list{
				kind="subnet"
				uuid=nutanix_subnet.overlay-subnet.id
			}
			cluster_reference_list{
				kind="cluster"
				uuid=local.cluster1
			}
			vpc_reference_list{
				kind="vpc"
				uuid= nutanix_vpc.acctest-managed.id
			}
		}
	`, subnetName, name, description, categoryName, categoryVal, limit, rsType)
}
