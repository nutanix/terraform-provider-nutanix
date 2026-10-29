package dataprotectionv2_test

import (
	"fmt"
	"os"
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
					// Sync promote removes the source VM (or races to). Drop it from
					// state so post-test destroy does not fail with "VM not found".
					dropResourceFromState(vmResourceName),
					// Clean up the promoted VM on the remote PC.
					deleteRestoredVM(vmName),
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
	description := "create a new protected VG and promote it"
	categoryKey := fmt.Sprintf("tf-test-category-pp-promote-vg-%d", r)
	categoryValue := fmt.Sprintf("tf_test_category_pp_promote_vg_%d", r)

	// resource names for the test
	categoryResourceName := "nutanix_category_v2.test"
	protectionPolicyResourceName := "nutanix_protection_policy_v2.test"
	vgResourceName := "nutanix_volume_group_v2.test"
	registrationResourceName := "nutanix_pc_registration_v2.node-registration"

	// Always tear down PE registration / blacklist state after the test,
	// including when apply fails before CheckDestroy can capture ext IDs
	// (e.g. "remote cluster is blacklisted for registration").
	t.Cleanup(func() {
		if err := ensurePromoteVGRemotePEClean(); err != nil {
			t.Logf("PromoteVG post-test PE cleanup: %v", err)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			// Seed PC/PE UUIDs and clear leftover blacklist/registration from a
			// prior run so this attempt can register the remote PE again.
			seedPromoteVGCleanupIDs()
			if err := ensurePromoteVGRemotePEClean(); err != nil {
				t.Logf("PromoteVG pre-test PE cleanup: %v", err)
			}
			seedPromoteVGCleanupIDs()
		},
		Providers: acc.TestAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {Source: "hashicorp/null"},
		},
		CheckDestroy: testCheckDestroyPromoteVGAndUnregister,
		Steps: []resource.TestStep{
			// register existing remote PE, sync PP, VG + category; wait until protected-resource is ready
			{
				Config: testPromoteProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description, categoryKey, categoryValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(registrationResourceName, "id"),
					capturePromoteVGUnregisterIDs(),
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
					// poll GET protected-resources until ready, then settle (no blind sleep)
					waitForProtectedResourceAvailable(vgResourceName, maxRetries, retryInterval, sleepTime),
				),
			},
			// promote protected VG after readiness
			{
				PreConfig: func() {
					fmt.Println("Step 2: Promote Protected VG")
				},
				Config: testPromoteProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description, categoryKey, categoryValue) +
					testPromoteProtectedResourceVGConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNamePromoteProtectedResource, "id"),
					resource.TestCheckResourceAttrSet(resourceNamePromoteProtectedResource, "ext_id"),
				),
			},
		},
	})
}

func testPromoteProtectedResourceVMAndProtectionPolicyConfig(vmName, ppName, description string, r int) string {
	return fmt.Sprintf(`

locals {
  pcFilter            = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')"
  aosFilter           = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

# List domain Managers
data "nutanix_pcs_v2" "pcs-list" {}

# list Clusters
data "nutanix_clusters_v2" "clusters" {
  filter = local.aosFilter
}

locals {
  clusterExtId = data.nutanix_clusters_v2.clusters.cluster_entities[0].ext_id
  config = jsondecode(file("%[1]s"))
  availability_zone = local.config.availability_zone
}

data "nutanix_storage_containers_v2" "ssc" {
  filter = "clusterExtId eq '${local.clusterExtId}' and name eq 'SelfServiceContainer'"
  limit  = 1
}

locals {
  sourceStorageContainerExtId = data.nutanix_storage_containers_v2.ssc.storage_containers[0].ext_id
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
  disks {
    disk_address {
      bus_type = "SCSI"
      index    = 0
    }
    backing_info {
      vm_disk {
        disk_size_bytes = 1073741824 # 1 GiB
        storage_container {
          ext_id = local.sourceStorageContainerExtId
        }
      }
    }
  }
  power_state = "ON"

  depends_on = [nutanix_protection_policy_v2.test]
}


	`, filepath, vmName, description, ppName, r)
}

func testPromoteProtectedResourceVMConfig() string {
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

`, remotePCProviderConfig(testVars.AvailabilityZone.RemotePcIP))
}

// testPromoteProtectedResourceVGAndProtectionPolicyConfig returns the setup for a sync-protected VG.
// Promote is applied in a second test step via testPromoteProtectedResourceVGConfig after
// waitForProtectedResourceAvailable confirms the protected-resource API is ready.
// The remote PE already has a cluster; this config only registers it to PC (and
// CheckDestroy unregisters it after stretch params clear).
// Steps:
// 1. Register the existing remote PE cluster to PC
// 2. Modify the firewall rules between the remote PE and the local PE
// 3. Create a category
// 4. Create a sync protection policy (RPO=0)
// 5. Create a volume group
// 6. Associate the category to the volume group
func testPromoteProtectedResourceVGAndProtectionPolicyConfig(vgName, ppName, description, categoryKey, categoryValue string) string {
	// Resolve PC host once for both direct-PC and NC setups so local-exec
	// polls hit the same endpoint the provider uses (see acc.PCEndpointHost).
	pcHost := acc.PCEndpointHost()
	pcPort := acc.PCEndpointPort()
	// Use acceptance-test credentials (env), not test_config pc_username —
	// on NC, local admin from JSON is not valid against the domain FQDN.
	pcUser := os.Getenv("NUTANIX_USERNAME")
	pcPass := os.Getenv("NUTANIX_PASSWORD")

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

  # PC API host/creds for local-exec polls: same host + auth as the provider
  # (direct PC or NC <domain>.domains.<central>).
  pc_host     = "%[7]s"
  pc_port     = "%[8]s"
  pc_username = "%[9]s"
  pc_password = "%[10]s"

  localClusterIP   = local.data_protection.local_cluster_pe
  localClusterVIP  = local.data_protection.local_cluster_vip
  remoteClusterIP  = local.data_protection.remote_cluster_pe
  remoteClusterVIP = local.data_protection.remote_cluster_vip
  username         = local.config.pe_username
  password         = local.config.pe_password

  # UserKnownHostsFile=/dev/null is required: StrictHostKeyChecking=no alone still
  # fails when known_hosts has a stale key for a reimaged PE (CI runners reuse hosts).
  remoteClusterSSHCommand = "sshpass -p '${local.config.ssh_pe_password}' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o PreferredAuthentications=password -o PubkeyAuthentication=no ${local.config.ssh_pe_username}@${local.remoteClusterIP}"
  localClusterSSHCommand  = "sshpass -p '${local.config.ssh_pe_password}' ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o PreferredAuthentications=password -o PubkeyAuthentication=no ${local.config.ssh_pe_username}@${local.localClusterIP}"

  # Commands to modify the firewall rules between the clusters
  modifyFirewallRulesCommand       = "/usr/local/nutanix/cluster/bin/modify_firewall -f -r"
  modifyLocalClusterFirewallRules  = "${local.localClusterSSHCommand} '${local.modifyFirewallRulesCommand} ${local.remoteClusterIP},${local.remoteClusterVIP} -p 2030,2036,2073,2090,8740 -i eth0'"
  modifyRemoteClusterFirewallRules = "${local.remoteClusterSSHCommand} '${local.modifyFirewallRulesCommand} ${local.localClusterIP},${local.localClusterVIP}  -p 2030,2036,2073,2090,8740 -i eth0'"

  # KB 12749: a PE cannot be unregistered while it still has synchronously
  # protected entities. Deleting the policy/removing the category above only
  # *starts* unprotection; the stretch params take time (~10 min) to clear from
  # the target PE. Poll the target PE until stretch_params_printer reports no
  # remaining stretch params (bounded ~15 min) before CheckDestroy unregisters
  # the PE, instead of a fixed sleep that can under-wait and trigger
  # "error_code 50 / Unregistration failure".
  waitForStretchClearedOnTarget = "${local.remoteClusterSSHCommand} 'source /etc/profile >/dev/null 2>&1; for i in $(seq 1 60); do o=$(stretch_params_printer 2>/dev/null); if [ -z \"$o\" ]; then echo \"stretch params cleared on target PE\"; exit 0; fi; echo \"waiting for stretch params to clear ($i/60)\"; sleep 15; done; echo \"WARN: timed out (~15m) waiting for stretch params to clear\"; exit 0'"
}

# Register the existing remote PE cluster to PC (node already has a cluster).
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
  # Modify the firewall rules on Remote cluster (fail apply if SSH/firewall fails)
  provisioner "local-exec" {
    command = local.modifyRemoteClusterFirewallRules
    when    = create
  }
}

# After registering the remote PE to PC, Prism Central takes a few seconds to
# list the newly registered cluster. Poll until the remote PE IP appears so the
# new-cls lookup below can resolve local.newClusterExtId.
resource "null_resource" "wait_cluster_registered" {
  triggers = {
    remote_pe = local.remoteClusterIP
  }

  provisioner "local-exec" {
    # Single-line command (same pattern as cluster entity tests): use $(seq)
    # with one $, and $$i so Terraform leaves a shell $i. Heredoc + $$(seq)
    # was passed through literally and broke /bin/sh.
    command = "for i in $(seq 1 60); do if curl -ksSf -u '${local.pc_username}:${local.pc_password}' 'https://${local.pc_host}:${local.pc_port}/api/clustermgmt/v4.3/config/clusters' | grep -q '${self.triggers.remote_pe}'; then echo 'remote PE ${self.triggers.remote_pe} is now listed on PC'; exit 0; fi; echo \"waiting for remote PE ${self.triggers.remote_pe} to appear on PC (attempt $$i)\"; sleep 5; done; echo 'ERROR: timed out waiting for remote PE ${self.triggers.remote_pe} to appear on PC'; exit 1"
  }

  depends_on = [nutanix_pc_registration_v2.node-registration]
}

# Resolve the registered remote PE by CVM IP / VIP (cluster name is unknown
# because we did not create it in this test).
data "nutanix_clusters_v2" "new-cls" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
  depends_on = [null_resource.wait_cluster_registered]
}

locals {
  remote_cluster_matches = [
    for c in data.nutanix_clusters_v2.new-cls.cluster_entities : c.ext_id
    if(
      length([
        for n in try(c.nodes[0].node_list, []) : n
        if try(n.controller_vm_ip[0].ipv4[0].value, "") == local.remoteClusterIP
      ]) > 0
      || try(c.network[0].external_address[0].ipv4[0].value, "") == local.remoteClusterVIP
    )
  ]
  newClusterExtId = local.remote_cluster_matches[0]
}

# Sync VG protection needs at least one disk; empty VGs never get a dormant
# secondary (Castor: "Dormant VG config not found" on promote).
data "nutanix_storage_containers_v2" "ssc" {
  filter = "clusterExtId eq '${local.clusterExtId}' and name eq 'SelfServiceContainer'"
  limit  = 1
}

locals {
  sourceStorageContainerExtId = data.nutanix_storage_containers_v2.ssc.storage_containers[0].ext_id
}

# Create Category
resource "nutanix_category_v2" "test" {
  key   = "%[5]s"
  value = "%[6]s"

  # Modify the firewall rules on Local cluster (fail apply if SSH/firewall fails)
  provisioner "local-exec" {
    command = local.modifyLocalClusterFirewallRules
  }
  depends_on = [nutanix_pc_registration_v2.node-registration]
}

# Before CheckDestroy unregisters the PE, wait until the target has no
# synchronously protected entities left (stretch params cleared) so unregister
# succeeds (KB 12749). A destroy-time provisioner may only reference "self", so
# the poll command and PE/PC ext IDs are captured in triggers.
# depends_on ties this to the registration (destroyed before it), and the
# protection policy depends on it (destroyed after the policy, i.e. only once
# unprotection has been triggered by the policy delete).
resource "null_resource" "wait_stretch_cleared" {
  triggers = {
    wait_cmd              = local.waitForStretchClearedOnTarget
    pc_ext_id             = local.pcExtId
    remote_cluster_ext_id = local.newClusterExtId
  }

  provisioner "local-exec" {
    command    = self.triggers.wait_cmd
    when       = destroy
    on_failure = continue
  }

  depends_on = [nutanix_pc_registration_v2.node-registration, null_resource.wait_cluster_registered]
}

resource "nutanix_protection_policy_v2" "test" {
  name        = "%[3]s"
  description = "%[4]s"

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
  name              = "%[2]s"
  description       = "%[4]s"
  cluster_reference = local.clusterExtId
  disks {
    disk_size_bytes = 1073741824 # 1 GiB
    index           = 1
    disk_data_source_reference {
      name        = "vg-disk-%[2]s"
      ext_id      = local.sourceStorageContainerExtId
      entity_type = "STORAGE_CONTAINER"
    }
  }
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
}

`, filepath, vgName, ppName, description, categoryKey, categoryValue, pcHost, pcPort, pcUser, pcPass)
}

func testPromoteProtectedResourceVGConfig() string {
	return `

resource "nutanix_promote_protected_resource_v2" "test" {
  ext_id = nutanix_volume_group_v2.test.id
  depends_on = [nutanix_associate_category_to_volume_group_v2.test]
}
`
}
