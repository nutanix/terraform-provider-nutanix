package networking_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePbr = "nutanix_pbr.acctest-managed"

func TestAccNutanixPbr_basic(t *testing.T) {
	r := randIntBetween(221, 230)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
				),
			},
		},
	})
}

func TestAccNutanixPbr_WithSourceExternalDestinationNetwork(t *testing.T) {
	r := randIntBetween(231, 240)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.address_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "ALL"),
				),
			},
			{
				Config: testAccNutanixPbrConfigUpdateWithSourceExternalDestinationNetwork(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", fmt.Sprintf("acctest-managed-%d-updated", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "DENY"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.prefix_length", "24"),
				),
			},
		},
	})
}

func TestAccNutanixPbr_WithTCP(t *testing.T) {
	r := randIntBetween(241, 250)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfigWithSourceNetworkDestinationExternalWithTCP(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "TCP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.source_port_range_list.0.start_port", "50"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.destination_port_range_list.0.end_port", "40"),
				),
			},
			{
				Config: testAccNutanixPbrConfigUpdateWithSourceExternalDestinationNetwork(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", fmt.Sprintf("acctest-managed-%d-updated", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "DENY"),
				),
			},
		},
	})
}

func TestAccNutanixPbr_WithProtocolUDP(t *testing.T) {
	r := randIntBetween(251, 260)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfigWithSourceNetworkDestinationExternalWithTCP(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "TCP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.source_port_range_list.0.start_port", "50"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.destination_port_range_list.0.end_port", "40"),
				),
			},
			{
				Config: testAccNutanixPbrConfigWithSourceExternalDestinationAnyWithUDP(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", fmt.Sprintf("acctest-managed-%d-updated", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "UDP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "REROUTE"),
					resource.TestCheckResourceAttr(resourceNamePbr, "service_ip_list.0", "10.2.2.34"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.udp.0.source_port_range_list.0.start_port", "50"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.udp.0.destination_port_range_list.0.end_port", "40"),
				),
			},
		},
	})
}

func TestAccNutanixPbr_WithProtocolICMP(t *testing.T) {
	r := randIntBetween(261, 270)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfigWithSourceNetworkDestinationExternalWithTCP(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "TCP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.source_port_range_list.0.start_port", "50"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.destination_port_range_list.0.end_port", "40"),
				),
			},
			{
				Config: testAccNutanixPbrConfigWithSourceAnyDestinationExternalWithICMP(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", fmt.Sprintf("acctest-managed-%d-updated", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "ICMP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.address_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.icmp.0.icmp_code", "20"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.icmp.0.icmp_type", "2"),
				),
			},
		},
	})
}

func TestAccNutanixPbr_WithProtocolNumber(t *testing.T) {
	r := randIntBetween(271, 280)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfigWithSourceNetworkDestinationExternalWithTCP(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "TCP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "INTERNET"),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "PERMIT"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.source_port_range_list.0.start_port", "50"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.destination_port_range_list.0.end_port", "40"),
				),
			},
			{
				Config: testAccNutanixPbrConfigUpdateWithSourceAnyDestinationAnyWithProtocolNumber(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", fmt.Sprintf("acctest-managed-%d-updated", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "PROTOCOL_NUMBER"),
					resource.TestCheckResourceAttr(resourceNamePbr, "action", "DENY"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "source.0.address_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "destination.0.address_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.protocol_number", "50"),
				),
			},
		},
	})
}

func TestAccNutanixPbr_WithVPCName(t *testing.T) {
	r := randIntBetween(281, 290)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfigWithVpcName(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
				),
			},
		},
	})
}

func TestAccNutanixPbr_WithVPCNameAndBidirectional(t *testing.T) {
	r := randIntBetween(291, 300)
	pbrName := fmt.Sprintf("acctest-managed-%d", r)
	vpcName := fmt.Sprintf("acctest-vpc-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNutanixPbrConfigWithVpcNameAndBidirectional(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNamePbr, "name", pbrName),
					resource.TestCheckResourceAttr(resourceNamePbr, "vpc_name", vpcName),
					resource.TestCheckResourceAttr(resourceNamePbr, "is_bidirectional", "true"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_type", "TCP"),
					resource.TestCheckResourceAttr(resourceNamePbr, "priority", fmt.Sprintf("%d", r)),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.destination_port_range_list.#", "2"),
					resource.TestCheckResourceAttr(resourceNamePbr, "protocol_parameters.0.tcp.0.source_port_range_list.#", "2"),
				),
			},
		},
	})
}

func testAccNutanixPbrConfig(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.42.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d"
		priority = %[1]d
		protocol_type = "ALL"
		action = "PERMIT"
		vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id
		source{
		  address_type = "ALL"
		}
		destination{
		  address_type = "ALL"
		}
	}
	`, r)
}

func testAccNutanixPbrConfigUpdateWithSourceExternalDestinationNetwork(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.43.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d-updated"
		priority = %[1]d
		protocol_type = "ALL"
		action = "DENY"
		vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id
		source{
		  address_type = "INTERNET"
		}
		destination{
			subnet_ip=  "1.2.2.0"
			prefix_length=  24
		}
	}
	`, r)
}

func testAccNutanixPbrConfigWithSourceNetworkDestinationExternalWithTCP(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.44.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d"
		priority = %[1]d
		action = "PERMIT"
		vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id
		destination{
		  address_type = "INTERNET"
		}
		source{
			subnet_ip=  "1.2.2.0"
			prefix_length=  24
		}
		protocol_type = "TCP"
		protocol_parameters{
			tcp{
				source_port_range_list{
					end_port  = 50
					start_port = 50
				}
				destination_port_range_list{
					end_port  = 40
					start_port = 40
				}
			}
		}

	}
	`, r)
}

func testAccNutanixPbrConfigWithSourceExternalDestinationAnyWithUDP(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.45.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d-updated"
		priority = %[1]d
		vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id
		destination{
		  address_type = "ALL"
		}
		source{
			address_type = "INTERNET"
		}
		protocol_type = "UDP"
		protocol_parameters{
			udp{
				source_port_range_list{
					end_port  = 50
					start_port = 50
				}
				destination_port_range_list{
					end_port  = 40
					start_port = 40
				}
			}
		}
		action = "REROUTE"
		service_ip_list = ["10.2.2.34"]

	}
	`, r)
}

func testAccNutanixPbrConfigWithSourceAnyDestinationExternalWithICMP(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.46.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d-updated"
		priority = %[1]d
		vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id
		source{
		  address_type = "ALL"
		}
		destination{
			address_type = "INTERNET"
		}
		protocol_type = "ICMP"
		protocol_parameters{
			icmp {
				icmp_type = 2
				icmp_code = 20
			}
		}
		action = "PERMIT"
	}
	`, r)
}

func testAccNutanixPbrConfigUpdateWithSourceAnyDestinationAnyWithProtocolNumber(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.47.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d-updated"
		priority = %[1]d
		vpc_reference_uuid = resource.nutanix_vpc.test-vpc.id
		destination{
		  address_type = "ALL"
		}
		source{
			address_type = "ALL"
		}
		protocol_type = "PROTOCOL_NUMBER"
		protocol_parameters{
			protocol_number= "50"
		}
		action = "DENY"
	}
	`, r)
}

func testAccNutanixPbrConfigWithVpcName(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.48.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d"
		priority = %[1]d
		protocol_type = "ALL"
		action = "PERMIT"
		vpc_name = resource.nutanix_vpc.test-vpc.name
		source{
		  address_type = "ALL"
		}
		destination{
		  address_type = "ALL"
		}
	}
	`, r)
}

func testAccNutanixPbrConfigWithVpcNameAndBidirectional(r int) string {
	return fmt.Sprintf(`

	data "nutanix_clusters" "clusters" {}

	locals {
		cluster1 = [
		for cluster in data.nutanix_clusters.clusters.entities :
		cluster.metadata.uuid if cluster.service_list[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_subnet" "sub-test" {
		cluster_uuid = local.cluster1
		name        = "acctest-managed-%[1]d"
		description = "Description of my unit test VLAN"
		vlan_id     = %[1]d
		subnet_type = "VLAN"
		subnet_ip          = "10.250.140.0"
	  default_gateway_ip = "10.250.140.1"
	  prefix_length = 24
	  is_external = true
	  ip_config_pool_list_ranges = ["10.250.140.10 10.250.140.20"]
	}

	resource "nutanix_vpc" "test-vpc" {
		name = "acctest-vpc-%[1]d"


		external_subnet_reference_uuid = [
		  resource.nutanix_subnet.sub-test.id
		]

		common_domain_name_server_ip_list{
				ip = "8.8.8.9"
		}

		externally_routable_prefix_list{
		  ip=  "172.49.0.0"
		  prefix_length= 16
		}
	  }

	resource "nutanix_pbr" "acctest-managed" {
		name = "acctest-managed-%[1]d"
		priority = %[1]d
		protocol_type = "TCP"
		action = "PERMIT"
		vpc_name = resource.nutanix_vpc.test-vpc.name
		source{
		  address_type = "ALL"
		}
		destination{
		  address_type = "ALL"
		}
		protocol_parameters{
			tcp{
				source_port_range_list{
					end_port  = 40
					start_port = 30
				}
				destination_port_range_list{
					end_port  = 60
					start_port = 50
				}
				source_port_range_list{
					end_port  = 70
					start_port = 65
				}
				destination_port_range_list{
					end_port  = 80
					start_port = 75
				}
			}
		}
		is_bidirectional=true
	}
	`, r)
}
