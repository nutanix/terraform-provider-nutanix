package vmmv2_test

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/request/vm"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// implement vm check destroy function
func testAccCheckNutanixVmsResourceDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	vmClient := conn.VmmAPI.VMAPIInstance
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine_v2" {
			continue
		}
		getVmByIdRequest := import1.GetVmByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := vmClient.GetVmById(ctx, &getVmByIdRequest)
		if err == nil {
			deleteVmByIdRequest := import1.DeleteVmByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			_, err = vmClient.DeleteVmById(ctx, &deleteVmByIdRequest)
			if err != nil {
				return fmt.Errorf("error: VM still exists: %v", err)
			}
			return nil
		}
	}

	return nil
}

// testAccCheckGcDeployVMsDestroy deletes VMs created by nutanix_deploy_templates_v2
// resources. The deploy resource's Delete is a no-op, so we clean up here.
func testAccCheckGcDeployVMsDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	vmClient := conn.VmmAPI.VMAPIInstance
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_deploy_templates_v2" {
			continue
		}
		vmID := rs.Primary.ID
		if vmID == "" {
			continue
		}

		getReq := import1.GetVmByIdRequest{ExtId: utils.StringPtr(vmID)}
		if _, err := vmClient.GetVmById(ctx, &getReq); err != nil {
			log.Printf("[DEBUG] Deploy VM %s already gone: %v", vmID, err)
			continue
		}

		log.Printf("[DEBUG] Cleaning up deploy VM %s", vmID)
		delReq := import1.DeleteVmByIdRequest{ExtId: utils.StringPtr(vmID)}
		if _, err := vmClient.DeleteVmById(ctx, &delReq); err != nil {
			return fmt.Errorf("failed to delete deploy VM %s: %v", vmID, err)
		}
	}

	return nil
}
