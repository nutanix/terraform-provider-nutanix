package clustersv2_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"regexp"
	"testing"
)

// Logic covered in create cluster resource test
func TestAccNutanixPePcRegistrationV2Resource_ValidationsDomainManagerRemoteClusterSpec(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutPcExtId(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutAuthenticationPassword(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutAuthenticationUsername(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutCloudType(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixPePcRegistrationV2Resource_ValidationsAOSRemoteClusterSpec(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{

			{
				Config:      testAccClusterResourceAOSRemoteClusterSpecInvalidConfigWithoutRemoteCluster(),
				ExpectError: regexp.MustCompile("Insufficient remote_cluster blocks"),
			},
			{
				Config:      testAccClusterResourceAOSRemoteClusterSpecInvalidConfigWithoutAuthenticationPassword(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccClusterResourceAOSRemoteClusterSpecInvalidConfigWithoutAuthenticationUsername(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccNutanixPePcRegistrationV2Resource_ValidationsClusterReference(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterResourceClusterReferenceInvalidConfigWithoutClusterExtId(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// Invalid Configs for DomainManagerRemoteClusterSpec
func testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutPcExtId() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
	  remote_cluster {
		domain_manager_remote_cluster_spec {
		  remote_cluster {
			address {
			  ipv4 {
				value = "0.0.0.0"
			  }
			}
			credentials {
			  authentication {
				username = "test"
				password = "test"
			  }
			}
		  }
		  cloud_type = "NUTANIX_HOSTED_CLOUD"
		}
	  }
	}`
}

func testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutAuthenticationPassword() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
		pc_ext_id = "00000000-0000-0000-0000-000000000000"
        remote_cluster {
			domain_manager_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = "0.0.0.0"
				  }
				}
				credentials {
				  authentication {
					username = "test"
				  }
				}
			  }
			  cloud_type = "NUTANIX_HOSTED_CLOUD"
			}
		}				
	}`
}

func testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutAuthenticationUsername() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
		pc_ext_id = "00000000-0000-0000-0000-000000000000"
        remote_cluster {
			domain_manager_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = "0.0.0.0"
				  }
				}
				credentials {
				  authentication {
					password = "test"
				  }
				}
			  }
			  cloud_type = "NUTANIX_HOSTED_CLOUD"
			}
		}				
	}`
}

func testAccClusterResourceDomainManagerRemoteClusterSpecInvalidConfigWithoutCloudType() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
		pc_ext_id = "00000000-0000-0000-0000-000000000000"
        remote_cluster {
			domain_manager_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = "0.0.0.0"
				  }
				}
				credentials {
				  authentication {
					username = "test"
					password = "test"
				  }
				}
			  }
			}
		}				
	}`
}

func testAccClusterResourceAOSRemoteClusterSpecInvalidConfigWithoutAuthenticationPassword() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
		pc_ext_id = "00000000-0000-0000-0000-000000000000"
		remote_cluster {
			aos_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = "0.0.0.0"
				  }
				}
				credentials {
				  authentication {
					username = "test"
				  }
				}
			  }
			}
		}		
	}`
}

func testAccClusterResourceAOSRemoteClusterSpecInvalidConfigWithoutRemoteCluster() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
		pc_ext_id = "00000000-0000-0000-0000-000000000000"
		remote_cluster {
			aos_remote_cluster_spec {			  
			}
		}		
	}`
}

func testAccClusterResourceAOSRemoteClusterSpecInvalidConfigWithoutAuthenticationUsername() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
		pc_ext_id = "00000000-0000-0000-0000-000000000000"
		remote_cluster {
			aos_remote_cluster_spec {
			  remote_cluster {
				address {
				  ipv4 {
					value = "0.0.0.0"
				  }
				}
				credentials {
				  authentication {
					password = "test"
				  }
				}
			  }
			}
		}		
	}`
}

// Invalid Configs for ClusterReference

func testAccClusterResourceClusterReferenceInvalidConfigWithoutClusterExtId() string {
	return `
	resource "nutanix_pc_registration_v2" "test" {
	  remote_cluster {
		cluster_reference {
		  ext_id = "11111111-1111-1111-1111-111111111111"
		}
	  }
	}`
}
