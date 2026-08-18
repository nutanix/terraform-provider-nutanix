package dataprotectionv2_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNamePromoteProtectedResource = "nutanix_promote_protected_resource_v2.test"

const maxRetries = 90
const retryInterval = 10 * time.Second
const sleepTime = 5 * time.Minute

func TestAccV2NutanixPromoteProtectedResourceResource_PromoteVm(t *testing.T) {
	// if the test is running using NUTANIX_API_KEY, skip the test
	if os.Getenv("NUTANIX_API_KEY") != "" {
		t.Skip("Skipping test as it not supported using NUTANIX_API_KEY")
	}
	r := acctest.RandIntRange(1, 99)
	vmName := fmt.Sprintf("tf-test-protected-vm-promote-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-vm-%d", r)
	description := "create a new protected vm and promote it"

	vmResourceName := "nutanix_virtual_machine_v2.test"
	datasourceNamePromotedVM := "data.nutanix_virtual_machines_v2.promoted-vm"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		ProtoV5ProviderFactories: acc.TestAccProtoV5ProviderFactories,
		CheckDestroy:             testCheckDestroyProtectedResourceAndCleanupForPromoteVM,
		Steps: []resource.TestStep{
			// create protection policy and protected vm
			{
				Config: testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(vmResourceName, "id"),
					resource.TestCheckResourceAttr(vmResourceName, "name", vmName),
					waitForVMToBeProtected(vmResourceName, "protection_type", "RULE_PROTECTED", maxRetries, retryInterval, sleepTime),
				),
			},
			//promote protected vm
			{
				PreConfig: func() {
					fmt.Println("Step 2: Promote Protected Resource")
				},
				Config: testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description, r) +
					testPromoteProtectedResourceVMConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNamePromoteProtectedResource, "id"),
					resource.TestCheckResourceAttrSet(resourceNamePromoteProtectedResource, "ext_id"),
					// check if the promoted vm is created
					resource.TestCheckResourceAttrPair(datasourceNamePromotedVM, "vms.0.name", vmResourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceNamePromotedVM, "vms.0.num_cores_per_socket", vmResourceName, "num_cores_per_socket"),
					resource.TestCheckResourceAttrPair(datasourceNamePromotedVM, "vms.0.num_sockets", vmResourceName, "num_sockets"),
				),
			},
		},
	})
}

func TestAccV2NutanixPromoteProtectedResourceResource_PromoteVG(t *testing.T) {
	r := acctest.RandIntRange(1, 99)

	// variables for the test
	vgName := fmt.Sprintf("tf-test-protected-vg-promote-%d", r)
	ppName := fmt.Sprintf("tf-test-protected-policy-promote-vg-%d", r)
	clusterName := fmt.Sprintf("tf-test-cluster-pp-%d", r)
	description := "create a new protected VG and promote it"
	categoryKey := fmt.Sprintf("tf-test-category-pp-promote-vg-%d", r)
	categoryValue := fmt.Sprintf("tf_test_category_pp_promote_vg_%d", r)

	// resource names for the test
	newClusterDataSourceName := "data.nutanix_clusters_v2.new-cls"
	categoryResourceName := "nutanix_category_v2.test"
	protectionPolicyResourceName := "nutanix_protection_policy_v2.test"
	vgResourceName := "nutanix_volume_group_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {Source: "hashicorp/null"},
		},
		CheckDestroy: testCheckDestroyProtectedResourceAndCleanup,
		Steps: []resource.TestStep{
			//// create protection policy and protected vg
			{
				Config: testPromoteProtectedResourceVGAndProtectionPolicyConfig(clusterName, vgName, ppName, description, categoryKey, categoryValue),
				Check: resource.ComposeTestCheckFunc(
					// cluster check
					resource.TestCheckResourceAttrSet(newClusterDataSourceName, "cluster_entities.0.ext_id"),
					resource.TestCheckResourceAttr(newClusterDataSourceName, "cluster_entities.0.name", clusterName),
					// category check
					resource.TestCheckResourceAttrSet(categoryResourceName, "id"),
					resource.TestCheckResourceAttr(categoryResourceName, "key", categoryKey),
					resource.TestCheckResourceAttr(categoryResourceName, "value", categoryValue),
					// protection policy check
					resource.TestCheckResourceAttrSet(protectionPolicyResourceName, "id"),
					resource.TestCheckResourceAttr(protectionPolicyResourceName, "name", ppName),
					resource.TestCheckResourceAttr(protectionPolicyResourceName, "description", description),
					// volume group check
					resource.TestCheckResourceAttrSet(vgResourceName, "id"),
					resource.TestCheckResourceAttr(vgResourceName, "name", vgName),
					resource.TestCheckResourceAttr(vgResourceName, "description", description),
				),
			},
		},
	})
}

func testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description string, r int) string {
	return fmt.Sprintf(`
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

# list Clusters
data "nutanix_clusters_v2" "clusters" {}

locals {
  clusterExtId = [
    for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
    cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
  ][0]
  config = jsondecode(file("%[1]s"))
  availability_zone = local.config.availability_zone
}

# Create Category
resource "nutanix_category_v2" "test" {
  key   = "tf-synchronous-pp-%[5]d"
  value = "tf_synchronous_pp_%[5]d"
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[4]s"
  description = "%[3]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = local.availability_zone.pc_ext_id
    label                 = "target"
    is_primary            = false
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.availability_zone.cluster_ext_id]
      }
    }
  }

  category_ids = [nutanix_category_v2.test.id]
}

data "nutanix_storage_containers_v2" "ngt-sc" {
	filter = "clusterExtId eq '${local.clusterExtId}' and startswith(name,'default-container-')"
	limit = 1
}

resource "nutanix_virtual_machine_v2" "test" {
  name                 = "%[2]s"
  description          = "%[3]s"
  num_cores_per_socket = 1
  num_sockets          = 1
  cluster {
    ext_id = local.clusterExtId
  }
  categories {
    ext_id = nutanix_category_v2.test.id
  }
power_state = "ON"
  disks {
    disk_address{
      bus_type = "SCSI"
      index = 0
    }
    backing_info{
      vm_disk{
        disk_size_bytes = "1073741824" # 10 GB
        storage_container{
          ext_id = data.nutanix_storage_containers_v2.ngt-sc.storage_containers[0].ext_id
        }
      }
    }
  }
  cd_roms {
    disk_address {
      bus_type = "IDE"
      index    = 0
    }
  }
  depends_on = [nutanix_protection_policy_v2.test]
}


	`, filepath, vmName, description, ppName, r)
}

func testPromoteProtectedResourceVMConfig() string {
	remotePcIP := testVars.AvailabilityZone.RemotePcIP

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

`, username, password, remotePcIP, insecure, port)

	return fmt.Sprintf(
		`


%[1]s

resource "nutanix_promote_protected_resource_v2" "test" {
  provider = nutanix-2
  ext_id = nutanix_virtual_machine_v2.test.id
  provisioner "local-exec" {
    command = "sleep 10" # sleep for 10 seconds after promoting the resource, to read the promoted resource
  }
}

data "nutanix_virtual_machines_v2" "promoted-vm" {
  provider = nutanix-2
  filter = "name eq '${nutanix_virtual_machine_v2.test.name}'"
  limit = 1
  depends_on = [nutanix_promote_protected_resource_v2.test]
}

`, remoteHostProviderConfig)
}

// testPromoteProtectedResourceVGAndProtectionPolicyConfig returns the configuration for promoting a protected VG
// Steps:
// 1. Create a new cluster
// 2. Register the cluster to PC
// 3. Modify the firewall rules between the new cluster and the PC cluster
// 4. Create a category
// 5. Create a protection policy
// 6. Create a volume group
// 7. Associate the category to the volume group
// 8. Promote the protected resource (VG)
func testPromoteProtectedResourceVGAndProtectionPolicyConfig(clusterName, vgName, ppName, description, categoryKey, categoryValue string) string {
	return fmt.Sprintf(`

# list Clusters
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}
# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}


locals {
  clusterExtId      = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
  pcExtId           = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
  config            = (jsondecode(file("%[1]s")))
  clusters          = local.config.clusters
  availability_zone = local.config.availability_zone
  data_protection   = local.config.data_protection

  randomNum = 99

  # Commands to reset the cluster password
  localClusterIP   = local.data_protection.local_cluster_pe
  localClusterVIP  = local.data_protection.local_cluster_vip
  remoteClusterIP  = local.data_protection.remote_cluster_pe
  remoteClusterVIP = local.data_protection.remote_cluster_vip
  username         = local.config.pe_username
  password         = local.config.pe_password

  resetClusterPassword = "/home/nutanix/prism/cli/ncli user reset-password user-name=${local.username} password=${local.password}"

  remoteClusterSSHCommand = "sshpass -p '${local.config.ssh_pe_password}' ssh -o StrictHostKeyChecking=no ${local.config.ssh_pe_username}@${local.remoteClusterIP}"
  localClusterSSHCommand  = "sshpass -p '${local.config.ssh_pe_password}' ssh -o StrictHostKeyChecking=no ${local.config.ssh_pe_username}@${local.localClusterIP}"

  resetClusterPasswordCommand = "${local.remoteClusterSSHCommand} '${local.resetClusterPassword}'"

  # Commands to modify the firewall rules between the clusters
  modifyFirewallRulesCommand       = "/usr/local/nutanix/cluster/bin/modify_firewall -f -r"
  modifyLocalClusterFirewallRules  = "${local.localClusterSSHCommand} '${local.modifyFirewallRulesCommand} ${local.remoteClusterIP},${local.remoteClusterVIP} -p 2030,2036,2073,2090,8740 -i eth0'"
  modifyRemoteClusterFirewallRules = "${local.remoteClusterSSHCommand} '${local.modifyFirewallRulesCommand} ${local.localClusterIP},${local.localClusterVIP}  -p 2030,2036,2073,2090,8740 -i eth0'"

  # KB 12749: a PE cannot be unregistered while it still has synchronously
  # protected entities. Deleting the policy/removing the category above only
  # *starts* unprotection; the stretch params take time (~10 min) to clear from
  # the target PE. Poll the target PE until stretch_params_printer reports no
  # remaining stretch params (bounded ~15 min) before the cluster/registration is
  # destroyed, instead of a fixed sleep that can under-wait and trigger
  # "error_code 50 / Unregistration failure".
  waitForStretchClearedOnTarget = "${local.remoteClusterSSHCommand} 'source /etc/profile >/dev/null 2>&1; for i in $(seq 1 60); do o=$(stretch_params_printer 2>/dev/null); if [ -z \"$o\" ]; then echo \"stretch params cleared on target PE\"; exit 0; fi; echo \"waiting for stretch params to clear ($i/60)\"; sleep 15; done; echo \"WARN: timed out (~15m) waiting for stretch params to clear\"; exit 0'"
}

# check if the nodes is un configured or not
resource "nutanix_clusters_discover_unconfigured_nodes_v2" "test-discover-cluster-node" {
  ext_id       = local.pcExtId
  address_type = "IPV4"
  ip_filter_list {
    ipv4 {
      value = local.remoteClusterIP
    }
  }

  ## check if the node is  un configured or not
  lifecycle {
    postcondition {
      condition     = length(self.unconfigured_nodes) == 1
      error_message = "The node ${local.remoteClusterIP} are not unconfigured"
    }
  }

  depends_on = [data.nutanix_clusters_v2.clusters]
}

# create a new cluster
resource "nutanix_cluster_v2" "test" {
  name = "%[2]s"
  nodes {
    node_list {
      controller_vm_ip {
        ipv4 {
          value = local.remoteClusterIP
        }
      }
    }
  }
  config {
    cluster_function = ["AOS"]
    cluster_arch     = "X86_64"
    fault_tolerance_state {
      domain_awareness_level = "DISK"
    }
    redundancy_factor = 1
  }
  network {
    external_address {
      ipv4 {
        value = local.remoteClusterVIP
      }
    }
  }

  # Reset the cluster password
  provisioner "local-exec" {
    command    = local.resetClusterPasswordCommand
    on_failure = continue
  }
  # Set lifecycle to ignore changes
  lifecycle {
    ignore_changes = [network.0.smtp_server.0.server.0.password, links, categories, config.0.cluster_function]
  }
  depends_on = [nutanix_clusters_discover_unconfigured_nodes_v2.test-discover-cluster-node]
}


# register the cluster to pc
resource "nutanix_pc_registration_v2" "node-registration" {
  pc_ext_id = local.pcExtId
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = local.remoteClusterIP
          }
        }
        credentials {
          authentication {
            username = local.username
            password = local.password
          }
        }
      }
    }
  }
  # Modify the firewall rules on Remote cluster
  provisioner "local-exec" {
    command    = local.modifyRemoteClusterFirewallRules
    when       = create
    on_failure = continue
  }
  depends_on = [nutanix_cluster_v2.test]
}

# After registering the remote PE to PC, Prism Central takes a few seconds to
# list the newly registered cluster. Poll the clusters API until it appears so
# the data.nutanix_clusters_v2.new-cls lookup below doesn't read an empty list
# (which surfaced as "Invalid index" on cluster_entities.0.ext_id).
resource "null_resource" "wait_cluster_registered" {
  triggers = {
    cluster_name = nutanix_cluster_v2.test.name
  }

  provisioner "local-exec" {
    command = <<-EOT
      for i in $(seq 1 60); do
        if curl -ksS -u "$NUTANIX_USERNAME:$NUTANIX_PASSWORD" \
             "https://$NUTANIX_ENDPOINT:$NUTANIX_PORT/api/clustermgmt/v4.0/config/clusters" \
             | grep -q "\"${self.triggers.cluster_name}\""; then
          echo "cluster ${self.triggers.cluster_name} is now listed on PC"
          exit 0
        fi
        echo "waiting for cluster ${self.triggers.cluster_name} to appear on PC (attempt $i)"
        sleep 5
      done
      echo "ERROR: timed out waiting for cluster ${self.triggers.cluster_name} to appear on PC"
      exit 1
    EOT
  }

  depends_on = [nutanix_pc_registration_v2.node-registration]
}

# create a category, protection policy, volume group and associate it to the volume group
# list Clusters
data "nutanix_clusters_v2" "new-cls" {
  filter = "name eq '${nutanix_cluster_v2.test.name}'"
  depends_on = [null_resource.wait_cluster_registered]
}

locals {
  newClusterExtId = data.nutanix_clusters_v2.new-cls.cluster_entities.0.ext_id
}

# Create Category
resource "nutanix_category_v2" "test" {
  key   = "%[6]s"
  value = "%[7]s"

  # Modify the firewall rules on Local cluster
  provisioner "local-exec" {
    command    = local.modifyLocalClusterFirewallRules
    on_failure = continue
  }
  depends_on = [nutanix_pc_registration_v2.node-registration, nutanix_cluster_v2.test]
}

# Before the cluster/registration is destroyed, wait until the target PE has no
# synchronously protected entities left (stretch params cleared) so it can be
# unregistered cleanly (KB 12749). A destroy-time provisioner may only reference
# "self", so the poll command is captured in triggers and read via self.triggers.
# depends_on ties this to the cluster/registration (so it is destroyed before
# them), and the protection policy depends on it (so it is destroyed after the
# policy, i.e. only once unprotection has been triggered by the policy delete).
resource "null_resource" "wait_stretch_cleared" {
  triggers = {
    wait_cmd = local.waitForStretchClearedOnTarget
  }

  provisioner "local-exec" {
    command    = self.triggers.wait_cmd
    when       = destroy
    on_failure = continue
  }

  depends_on = [nutanix_pc_registration_v2.node-registration, nutanix_cluster_v2.test]
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[4]s"
  description = "%[5]s"

  replication_configurations {
    source_location_label = "source"
    remote_location_label = "target"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }
  replication_configurations {
    source_location_label = "target"
    remote_location_label = "source"
    schedule {
      recovery_point_type                           = "CRASH_CONSISTENT"
      recovery_point_objective_time_seconds         = 0
      sync_replication_auto_suspend_timeout_seconds = 10
    }
  }

  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "source"
    is_primary            = true
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.clusterExtId]
      }
    }
  }
  replication_locations {
    domain_manager_ext_id = data.nutanix_pcs_v2.pcs-list.pcs[0].ext_id
    label                 = "target"
    is_primary            = false
    replication_sub_location {
      cluster_ext_ids {
        cluster_ext_ids = [local.newClusterExtId]
      }
    }
  }

  category_ids = [nutanix_category_v2.test.id]

  # Ensure the stretch-cleared wait is destroyed *after* this policy, so the poll
  # runs only once the policy delete has triggered unprotection.
  depends_on = [null_resource.wait_stretch_cleared]
}

resource "nutanix_volume_group_v2" "test" {
  name              = "%[3]s"
  description       = "%[5]s"
  cluster_reference = local.clusterExtId
  lifecycle {
    ignore_changes = [cluster_reference]
  }
  depends_on = [nutanix_protection_policy_v2.test]
}


resource "nutanix_associate_category_to_volume_group_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  categories {
    ext_id = nutanix_category_v2.test.id
  }
  provisioner "local-exec" {
    # sleep 9 min to wait for the vg to be protected
    command    = "sleep 540"
  }
}

resource "nutanix_promote_protected_resource_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  depends_on = [nutanix_associate_category_to_volume_group_v2.test]
}

`, filepath, clusterName, vgName, ppName, description, categoryKey, categoryValue)
}
