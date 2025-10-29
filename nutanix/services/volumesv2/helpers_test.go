package volumesv2_test

import (
	"fmt"
	"log"
	"strings"

	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var diskSizeBytes int64 = 5368709120
var updatedDiskSizeBytes int64 = 10737418240

func testAccCheckResourceAttrListNotEmpty(resourceName, attrName, subAttr string) resource.TestCheckFunc {
	log.Printf("[DEBUG] testAccCheckResourceAttrNotEmpty ###############################")
	return func(s *terraform.State) error {
		resourceInstance := s.RootModule().Resources[resourceName]

		if resourceInstance == nil {
			return fmt.Errorf("resource %s not found", resourceName)
		}

		prefix := attrName + "."
		subAttrPrefix := prefix + "%d." + subAttr
		log.Printf("[DEBUG] Attributes : %s", resourceInstance.Primary.Attributes)
		for i := 0; ; i++ {
			attr := fmt.Sprintf(subAttrPrefix, i)
			if _, ok := resourceInstance.Primary.Attributes[attr]; !ok {
				// No more items in the list
				break
			}
			log.Printf("[DEBUG]  Attribute : %s", attr)
			if resourceInstance.Primary.Attributes[attr] == "" {
				return fmt.Errorf("%s attribute %s is empty", resourceName, attr)
			}
		}
		return nil
	}
}

func testAccCheckNutanixVolumeGroupV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_volume_group_v2" {
			continue
		}
		if _, err := conn.VolumeAPI.VolumeAPIInstance.DeleteVolumeGroupById(utils.StringPtr(rs.Primary.ID)); err != nil {
			if strings.Contains(fmt.Sprint(err), "VOLUME_UNKNOWN_ENTITY_ERROR") {
				return nil
			}
			return err
		}
	}
	return nil
}

func resourceNutanixVolumeGroupV2Exists(conn *conns.Client, name string) (*string, error) {
	var vgUUID *string

	filter := fmt.Sprintf("name==%s", name)
	vgList, err := conn.VolumeAPI.VolumeAPIInstance.ListVolumeGroups(nil, nil, &filter, nil, nil, nil)

	log.Printf("Volume Group List: %v", vgList)

	if err != nil {
		return nil, err
	}

	for _, vg := range vgList.Data.GetValue().([]volumesClient.VolumeGroup) {
		if utils.StringValue(vg.Name) == name {
			vgUUID = vg.ExtId
		}
	}
	log.Printf("Volume Group UUID: %v", vgUUID)
	return vgUUID, nil
}

// Helper Functions
func testAndCheckComputedValues(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}

		return nil
	}
}

// VolumeGroup Resource

func testAccVolumeGroupResourceConfig(name, desc string) string {
	return fmt.Sprintf(`
	data "nutanix_clusters_v2" "clusters" {}

	locals {
		cluster1 =  [
			  for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
			  cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
		][0]		
	}

	resource "nutanix_volume_group_v2" "test" {
		name                               = "%[2]s"
		description                        = "%[3]s"
		should_load_balance_vm_attachments = false
		sharing_status                     = "SHARED"		
		created_by 						   = "admin"
		cluster_reference                  = local.cluster1
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
	`, filepath, name, desc)
}

func testAccVolumeGroupDiskResourceConfig(name, desc string, diskSizeBytes int) string {
	return fmt.Sprintf(`	  

      data "nutanix_storage_containers_v2" "test" {
		  filter = "clusterExtId eq '${local.cluster1}'"
		  limit  = 1
	  }
	  resource "nutanix_volume_group_disk_v2" "test" {
		volume_group_ext_id = resource.nutanix_volume_group_v2.test.id
		index               = 1
		description         = "%[3]s"
		disk_size_bytes     = %[4]d
	  
		disk_data_source_reference {
		  name        = "vg-disk-%[2]s"
		  ext_id      = data.nutanix_storage_containers_v2.test.storage_containers[0].ext_id
		  entity_type = "STORAGE_CONTAINER"
		  uris        = ["uri1","uri2"]
		}
	  
		disk_storage_features {
		  flash_mode {
			is_enabled = false
		  }
		}
	  
		lifecycle {
		  ignore_changes = [
			disk_data_source_reference
		  ]
		}
	  
		depends_on = [resource.nutanix_volume_group_v2.test]
	  }	  
	`, filepath, name, desc, diskSizeBytes)
}
