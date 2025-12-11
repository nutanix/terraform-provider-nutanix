package vmmv2_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

// implement vm check destroy function
func testAccCheckNutanixVmsResourceDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	vmClient := conn.VmmAPI.VMAPIInstance

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_virtual_machine_v2" {
			continue
		}
		vmResponse, err := vmClient.GetVmById(utils.StringPtr(rs.Primary.ID))
		if err == nil {
			etag := vmClient.ApiClient.GetEtag(vmResponse)
			args := make(map[string]interface{})
			args["If-Match"] = utils.StringPtr(etag)
			_, err = vmClient.DeleteVmById(utils.StringPtr(rs.Primary.ID), args)
			if err != nil {
				return fmt.Errorf("error: VM still exists: %v", err)
			}
			return nil
		}
	}

	return nil
}
