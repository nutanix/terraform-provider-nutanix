package prismv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameDeployPC = "nutanix_pc_deploy_v2.test"

func TestAccV2NutanixDeployPcResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-deploy-pc-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployPCConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDeployPC, "id"),
					resource.TestCheckResourceAttrSet(resourceNameDeployPC, "network.0.external_networks.0.network_ext_id"),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "config.0.name", name),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "config.0.size", "STARTER"),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "config.0.build_info.0.version", testVars.Prism.DeployPC.Version),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "network.0.external_networks.0.ip_ranges.0.begin.0.ipv4.0.value", testVars.Prism.DeployPC.IPRange.Begin),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "network.0.external_networks.0.default_gateway.0.ipv4.0.value", testVars.Prism.DeployPC.DefaultGateway),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "network.0.external_networks.0.subnet_mask.0.ipv4.0.value", testVars.Prism.DeployPC.SubnetMask),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "network.0.name_servers.0.ipv4.0.value", testVars.Prism.DeployPC.NameServers[0]),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "network.0.name_servers.1.ipv4.0.value", testVars.Prism.DeployPC.NameServers[1]),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "network.0.ntp_servers.0.fqdn.0.value", testVars.Prism.DeployPC.NtpServers[0]),
					resource.TestCheckResourceAttr(resourceNameDeployPC, "network.0.ntp_servers.1.fqdn.0.value", testVars.Prism.DeployPC.NtpServers[1]),
				),
			},
		},
	})
}

func testAccDeployPCConfig(name string) string {
	username := os.Getenv("NUTANIX_USERNAME")
	password := os.Getenv("NUTANIX_PASSWORD")
	port, _ := strconv.Atoi(os.Getenv("NUTANIX_PORT"))
	insecure, _ := strconv.ParseBool(os.Getenv("NUTANIX_INSECURE"))
	remoteHostProviderConfig := fmt.Sprintf(`
provider "nutanix-2" {
  username = "%[1]s"
  password = "%[2]s"
  endpoint = "%[3]s"
  insecure = %[4]t
  port     = %[5]d
}

`, username, password, testVars.Prism.DeployPC.PeIP, insecure, port)

	return fmt.Sprintf(`

%[2]s

locals {
  config = jsondecode(file("%[1]s"))
  deploy_pc = local.config.prism.deploy_pc
}

resource "nutanix_pc_deploy_v2" "test" {
  provider = nutanix-2
  timeouts {
    create = "120m"
  }
  config {
    build_info {
      version = local.deploy_pc.version
    }
    size = "STARTER"
    name = "%[3]s"
  }
  network {
    external_networks {
      network_ext_id = local.deploy_pc.network_id
      default_gateway {
        ipv4 {
          value = local.deploy_pc.default_gateway
        }
      }
      subnet_mask {
        ipv4 {
          value = local.deploy_pc.subnet_mask
        }
      }
      ip_ranges {
        begin {
          ipv4 {
            value = local.deploy_pc.ip_range.begin
          }
        }
        end {
          ipv4 {
            value = local.deploy_pc.ip_range.end
          }
        }
      }
    }
    name_servers {
      ipv4 {
        value = local.deploy_pc.name_servers[0]
      }
    }
    name_servers {
      ipv4 {
        value = local.deploy_pc.name_servers[1]
      }
    }
    ntp_servers {
      fqdn {
        value = local.deploy_pc.ntp_servers[0]
      }
    }
    ntp_servers {
      fqdn {
        value = local.deploy_pc.ntp_servers[1]
      }
    }
    ntp_servers {
      fqdn {
        value = local.deploy_pc.ntp_servers[2]
      }
    }
    ntp_servers {
      fqdn {
        value = local.deploy_pc.ntp_servers[3]
      }
    }
  }
}

 `, filepath, remoteHostProviderConfig, name)
}
