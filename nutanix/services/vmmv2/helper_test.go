package vmmv2_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/vm"
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
		vmResponse, err := vmClient.GetVmById(ctx, &getVmByIdRequest)
		if err == nil {
			etag := vmClient.ApiClient.GetEtag(vmResponse)
			args := make(map[string]interface{})
			args["If-Match"] = utils.StringPtr(etag)
			deleteVmByIdRequest := import1.DeleteVmByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			_, err = vmClient.DeleteVmById(ctx, &deleteVmByIdRequest, args)
			if err != nil {
				return fmt.Errorf("error: VM still exists: %v", err)
			}
			return nil
		}
	}

	return nil
}
