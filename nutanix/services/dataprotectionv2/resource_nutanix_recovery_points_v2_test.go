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

const resourceNameRecoveryPoints = "nutanix_recovery_points_v2.test"

// this test cases must be updated after merge it on temp design branch, since it need to create a vm and volume group within the test case

var (
	path, _  = os.Getwd()
	filepath = path + "/../../../test_config_v2.json"
)

func TestAccNutanixRecoveryPointsV2Resource_VmRecoveryPoints2(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)
	vmName := fmt.Sprintf("tf-test-rp-vm-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmConfig(vmName) + testRecoveryPointsResourceConfigWithVmRecoveryPoints(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.0.vm_ext_id"),
					//resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "vm_recovery_points.0.status", "COMPLETE"),
					//resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "vm_recovery_points.0.expiration_time", expirationTimeFormatted),
					//resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "vm_recovery_points.0.recovery_point_type", "APPLICATION_CONSISTENT"),
					//resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "vm_recovery_points.0.name", "tf-test-rp-vm-name"),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPointsV2Resource_VmRecoveryPointsWithAppConsProps(t *testing.T) {
	t.Skip("Skipping this test case as it is failing due to missing app consistent properties in get request")
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)
	vmName := fmt.Sprintf("tf-test-rp-vm-%d", r)
	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmConfig(vmName) + testRecoveryPointsResourceConfigWithVmRecoveryPointsWithAppConsProps(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.0.vm_ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "vm_recovery_points.0.application_consistent_properties.0.backup_type", "FULL_BACKUP"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "vm_recovery_points.0.application_consistent_properties.0.should_include_writers", "true"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "vm_recovery_points.0.application_consistent_properties.0.should_store_vss_metadata", "true"),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPointsV2Resource_VmRecoveryPointsWithMultipleVms(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)
	vmName := fmt.Sprintf("tf-test-rp-vm-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmConfig(vmName) + testVm2Config(vmName) +
					testRecoveryPointsResourceConfigWithVmRecoveryPointsWithMultipleVms(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.0.vm_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.1.vm_ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPointsV2Resource_VolumeGroupRecoveryPoints(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPoints(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "volume_group_recovery_points.0.volume_group_ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPointsV2Resource_VolumeGroupRecoveryPointsWithMultipleVGs(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPointsWithMultipleVGs(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "volume_group_recovery_points.0.volume_group_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "volume_group_recovery_points.1.volume_group_ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPointsV2Resource_RecoveryPointWithMultipleVmAndVGs(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)
	vmName := fmt.Sprintf("tf-test-rp-vm-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmConfig(vmName) + testVm2Config(vmName) +
					testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPointsWithMultipleVmAndVGs(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "CRASH_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.0.vm_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.1.vm_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "volume_group_recovery_points.0.volume_group_ext_id"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "volume_group_recovery_points.1.volume_group_ext_id"),
				),
			},
		},
	})
}

func TestAccNutanixRecoveryPointsV2Resource_UpdateExpirationTime(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-recovery-point-%d", r)
	vmName := fmt.Sprintf("tf-test-rp-vm-%d", r)

	// End time is two week later
	expirationTime := time.Now().Add(14 * 24 * time.Hour)
	// End time is one month later
	expirationTimeUpdate := time.Now().Add(30 * 24 * time.Hour)

	expirationTimeFormatted := expirationTime.UTC().Format(time.RFC3339)
	expirationTimeUpdateFormatted := expirationTimeUpdate.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccFoundationPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testVmConfig(vmName) + testRecoveryPointsResourceConfigWithVmRecoveryPoints(name, expirationTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.0.vm_ext_id"),
				),
			},
			{
				Config: testVmConfig(vmName) + testRecoveryPointsResourceConfigWithVmRecoveryPoints(name, expirationTimeUpdateFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "name", name),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "status", "COMPLETE"),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "expiration_time", expirationTimeUpdateFormatted),
					resource.TestCheckResourceAttr(resourceNameRecoveryPoints, "recovery_point_type", "APPLICATION_CONSISTENT"),
					resource.TestCheckResourceAttrSet(resourceNameRecoveryPoints, "vm_recovery_points.0.vm_ext_id"),
				),
			},
		},
	})
}

func testRecoveryPointsResourceConfigWithVmRecoveryPoints(name, expirationTime string) string {
	return fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			name = "tf-test-rp-vm-name"
			expiration_time = "%[2]s"
			recovery_point_type = "APPLICATION_CONSISTENT"
			status = "COMPLETE"
			vm_ext_id = nutanix_virtual_machine_v2.test-1.id 
		}
	}`, name, expirationTime)
}

func testRecoveryPointsResourceConfigWithVmRecoveryPointsWithMultipleVms(name, expirationTime string) string {
	return fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "CRASH_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-1.id 
		}
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-2.id 
		}
	}`, name, expirationTime)
}

func testRecoveryPointsResourceConfigWithVmRecoveryPointsWithAppConsProps(name, expirationTime string) string {
	return fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "APPLICATION_CONSISTENT"
		vm_recovery_points {
			vm_ext_id = vm_ext_id = nutanix_virtual_machine_v2.test-1.id 
			application_consistent_properties {
				  backup_type               = "FULL_BACKUP"
				  should_include_writers    = true
				  writers                   = ["0f95b402-67aa-431c-9eab-bf0907a99345", "0f95b402-67aa-431c-9eab-bf0907a99346"]
				  should_store_vss_metadata = true
				  object_type = "dataprotection.v4.common.VssProperties"
			}  
		}
	}`, name, expirationTime)
}

func testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPoints(name, expirationTime string) string {
	vg := testAccVolumeGroup1ResourceConfig("vg-"+name, "test volume group description")
	return vg + fmt.Sprintf(`
	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "CRASH_CONSISTENT"
		volume_group_recovery_points {
			volume_group_ext_id = nutanix_volume_group_v2.test-1.id
		}			
	}`, name, expirationTime)
}

func testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPointsWithMultipleVGs(name, expirationTime string) string {
	vg1 := testAccVolumeGroup1ResourceConfig("vg-1-"+name, "test volume group description")
	vg2 := testAccVolumeGroup2ResourceConfig("vg-2-"+name, "test volume group description")
	return vg1 + vg2 + fmt.Sprintf(`
	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "CRASH_CONSISTENT"
		volume_group_recovery_points {
			volume_group_ext_id = nutanix_volume_group_v2.test-1.id
		}	
		volume_group_recovery_points {
			volume_group_ext_id = nutanix_volume_group_v2.test-2.id
		}			
	}`, name, expirationTime)
}

func testRecoveryPointsResourceConfigWithVolumeGroupRecoveryPointsWithMultipleVmAndVGs(name, expirationTime string) string {
	vg1 := testAccVolumeGroup1ResourceConfig("vg-1-"+name, "test volume group description")
	vg2 := testAccVolumeGroup2ResourceConfig("vg-2-"+name, "test volume group description")
	return vg1 + vg2 + fmt.Sprintf(`

	resource "nutanix_recovery_points_v2" "test" {
		name                = "%[1]s"
		expiration_time     = "%[2]s"
		status              = "COMPLETE"
		recovery_point_type = "CRASH_CONSISTENT"
        vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-1.id 
		}
		vm_recovery_points {
			vm_ext_id = nutanix_virtual_machine_v2.test-2.id 
		}
		volume_group_recovery_points {
			volume_group_ext_id = nutanix_volume_group_v2.test-1.id
		}	
		volume_group_recovery_points {
			volume_group_ext_id = nutanix_volume_group_v2.test-2.id
		}			
	}`, name, expirationTime)
}

func testAccVolumeGroup1ResourceConfig(name, desc string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "cluster-list" {}

	locals {
		cluster_id = [
		  for cluster in data.nutanix_clusters_v2.cluster-list.cluster_entities :
		  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]
	}

	resource "nutanix_volume_group_v2" "test-1" {
		name                               = "%[1]s"
		description                        = "%[2]s"
		should_load_balance_vm_attachments = false
		sharing_status                     = "SHARED"		
		created_by 						   = "admin"
		cluster_reference                  = local.cluster_id
		iscsi_features {
			target_secret			 = "1234567891011"
			enabled_authentications  = "CHAP"
		}
		storage_features {
		  flash_mode {
			is_enabled = true
		  }
		}
		usage_type = "USER"
		is_hidden = false
		lifecycle {
			ignore_changes = [
			  iscsi_features[0].target_secret
			]
		}
	  }	  
	`, name, desc)
}

func testAccVolumeGroup2ResourceConfig(name, desc string) string {
	return fmt.Sprintf(`

	resource "nutanix_volume_group_v2" "test-2" {
		name                               = "%[1]s"
		description                        = "%[2]s"
		should_load_balance_vm_attachments = false
		sharing_status                     = "SHARED"		
		created_by 						   = "admin"
		cluster_reference                  = local.cluster_id
		iscsi_features {
			target_secret			 = "1234567891011"
			enabled_authentications  = "CHAP"
		}
		storage_features {
		  flash_mode {
			is_enabled = true
		  }
		}
		usage_type = "USER"
		is_hidden = false
		lifecycle {
			ignore_changes = [
			  iscsi_features[0].target_secret
			]
		}
	  }	  
	`, name, desc)
}
