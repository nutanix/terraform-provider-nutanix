package dataprotectionv2_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
	"time"
)

func waitForVmToBeProtected(resourceName, attributeName, desiredValue string, maxRetries int, retryInterval, sleepTime time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var lastValue string
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.VmmAPI.VMAPIInstance

		for i := 0; i < maxRetries; i++ {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource not found: %s", resourceName)
			}

			vmResp, err := client.GetVmById(utils.StringPtr(rs.Primary.ID))
			if err != nil {
				return fmt.Errorf("error getting vm by id: %v", err)
			}

			// read the attribute value from the response
			vm := vmResp.Data.GetValue().(config.Vm)
			lastValue = config.ProtectionType.GetName(*vm.ProtectionType)
			if lastValue == desiredValue {
				time.Sleep(sleepTime)
				fmt.Printf("[DEBUG] VM is %s\n", lastValue)
				return nil // Desired value reached
			}

			fmt.Printf("[DEBUG] Waiting for vm to be protected:  attribute %q to be %q. Current value: %q\n", attributeName, desiredValue, lastValue)
			// Wait before retrying
			time.Sleep(retryInterval)
		}

		return fmt.Errorf("failed to reach desired value for attribute %q: expected %q, got %q after %d retries", attributeName, desiredValue, lastValue, maxRetries)
	}
}

func testCheckDestroyPromoteProtectedResource(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	vmClient := conn.VmmAPI.VMAPIInstance
	ppClient := conn.DataPoliciesAPI.ProtectionPolicies

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_virtual_machine_v2" {
			//aJSON, _ := json.MarshalIndent(rs, "", " ")
			//fmt.Printf("resource: %s", aJSON)
			fmt.Printf("resource name: %s", rs.Primary.ID)
			_, err := vmClient.GetVmById(utils.StringPtr(rs.Primary.ID))
			if err == nil {
				fmt.Printf("VM still exists")
				_, err = vmClient.DeleteVmById(utils.StringPtr(rs.Primary.ID))
				if err != nil {
					return fmt.Errorf("error: VM still exists: %v", err)
				}
				return nil
			}
		}

		if rs.Type == "nutanix_protection_policy_v2" {
			_, err := ppClient.GetProtectionPolicyById(utils.StringPtr(rs.Primary.ID))
			if err == nil {
				fmt.Printf("Protection Policy still exists")
				_, err = ppClient.DeleteProtectionPolicyById(utils.StringPtr(rs.Primary.ID))
				if err != nil {
					return fmt.Errorf("error: Protection Policy still exists : %v", err)
				}
				return nil
			}
		}

	}

	return nil
}
