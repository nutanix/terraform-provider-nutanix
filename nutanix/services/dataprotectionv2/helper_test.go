package dataprotectionv2_test

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nutanix/ntnx-api-golang-clients/vmm-go-client/v4/models/vmm/v4/ahv/config"
	volumesClient "github.com/nutanix/ntnx-api-golang-clients/volumes-go-client/v4/models/volumes/v4/config"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func waitForVMToBeProtected(resourceName, attributeName, desiredValue string, maxRetries int, retryInterval, sleepTime time.Duration) resource.TestCheckFunc {
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
				log.Printf("[DEBUG] VM is %s\n", lastValue)
				time.Sleep(sleepTime)
				return nil // Desired value reached
			}

			log.Printf("[DEBUG] Waiting for vm to be protected:  attribute %q to be %q. Current value: %q\n", attributeName, desiredValue, lastValue)
			// Wait before retrying
			time.Sleep(retryInterval)
		}

		return fmt.Errorf("VM: failed to reach desired value for attribute %q: expected %q, got %q after %d retries", attributeName, desiredValue, lastValue, maxRetries)
	}
}

func testCheckDestroyProtectedResourceAndCleanup(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	vmClient := conn.VmmAPI.VMAPIInstance
	categoryClient := conn.PrismAPI.CategoriesAPIInstance
	ppClient := conn.DataPoliciesAPI.ProtectionPolicies

	vmExtID := ""
	ppExtID := ""
	categoryExtID := ""

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_virtual_machine_v2" {
			log.Printf("[DEBUG] Checking if VM still exists\n")
			vmExtID = rs.Primary.ID
		}

		if rs.Type == "nutanix_protection_policy_v2" {
			log.Printf("[DEBUG] Checking if Protection Policy still exists\n")
			ppExtID = rs.Primary.ID
		}

		if rs.Type == "nutanix_category_v2" {
			log.Printf("[DEBUG] Checking if Category still exists\n")
			categoryExtID = rs.Primary.ID
		}
	}

	// delete vm
	if vmExtID != "" {
		readResp, err := vmClient.GetVmById(utils.StringPtr(vmExtID))
		if err == nil {
			args := make(map[string]interface{})
			etag := vmClient.ApiClient.GetEtag(readResp)
			args["If-Match"] = utils.StringPtr(etag)
			_, err = vmClient.DeleteVmById(utils.StringPtr(vmExtID), args)
			if err != nil {
				return fmt.Errorf("error: VM still exists: %v", err)
			}
			log.Printf("[DEBUG] VM deleted\n")
		}
	}

	// delete protection policy
	if ppExtID != "" {
		_, err := ppClient.GetProtectionPolicyById(utils.StringPtr(ppExtID))
		if err == nil {
			log.Printf("[DEBUG] Protection Policy still exists")
			_, err = ppClient.DeleteProtectionPolicyById(utils.StringPtr(ppExtID))
			if err != nil {
				return fmt.Errorf("error: Protection Policy still exists : %v", err)
			}
			log.Printf("[DEBUG] Protection Policy deleted\n")
		}
	}

	// delete category
	if categoryExtID != "" {
		_, err := categoryClient.GetCategoryById(utils.StringPtr(categoryExtID), nil)
		if err == nil {
			log.Printf("[DEBUG] Category still exists")

			_, err = categoryClient.DeleteCategoryById(utils.StringPtr(categoryExtID))
			if err != nil {
				return fmt.Errorf("error: Category still exists : %v", err)
			}
			log.Printf("[DEBUG] Category deleted\n")
		}
	}

	return nil
}

func testCheckDestroyProtectedResourceAndCleanupForPromoteVM(state *terraform.State) error {
	// check on protection policy, category and vm on local site
	testCheckDestroyProtectedResourceAndCleanup(state)

	// delete promoted vm
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_virtual_machines_v2" {
			log.Printf("[DEBUG] Checking if VM still exists\n")
			vmExtID := rs.Primary.Attributes["vms.0.ext_id"]

			connRemote := acc.TestAccProvider2.Meta().(*conns.Client)
			client := connRemote.VmmAPI.VMAPIInstance

			readResp, err := client.GetVmById(utils.StringPtr(vmExtID))

			if err == nil {
				args := make(map[string]interface{})
				etag := client.ApiClient.GetEtag(readResp)
				args["If-Match"] = utils.StringPtr(etag)

				_, err = client.DeleteVmById(utils.StringPtr(vmExtID), args)
				if err != nil {
					return fmt.Errorf("error: Promoted VM still exists: %v", err)
				}
				log.Printf("[DEBUG] Promoted VM deleted\n")
				return nil
			}
			return fmt.Errorf("promoted VM not found")
		}
	}

	return nil
}

func deleteRestoredVM(vmName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider2.Meta().(*conns.Client)
		client := conn.VmmAPI.VMAPIInstance

		filter := fmt.Sprintf("startswith(name, '%s')", vmName)

		resp, err := client.ListVms(nil, nil, utils.StringPtr(filter), nil, nil)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if resp.Data == nil {
			return fmt.Errorf("no data returned from list vms on Remote site")
		}
		vms := resp.Data.GetValue().([]config.Vm)

		vm := vms[0]
		readResp, err := client.GetVmById(vm.ExtId)
		if err == nil {
			args := make(map[string]interface{})
			etag := client.ApiClient.GetEtag(readResp)
			args["If-Match"] = utils.StringPtr(etag)
			_, err = client.DeleteVmById(vm.ExtId, args)
			if err != nil {
				return fmt.Errorf("error: Restored VM still exists: %v", err)
			}
			log.Printf("[DEBUG] Restored VM deleted\n")
			return nil
		}

		return nil
	}
}

func deleteRestoredVg(vgName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestAccProvider2.Meta().(*conns.Client)
		client := conn.VolumeAPI.VolumeAPIInstance
		filter := fmt.Sprintf("startswith(name, '%s')", vgName)

		resp, err := client.ListVolumeGroups(nil, nil, utils.StringPtr(filter), nil, nil, nil)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if resp.Data == nil {
			return fmt.Errorf("no data returned from list Volume Groups on Remote site")
		}
		vgs := resp.Data.GetValue().([]volumesClient.VolumeGroup)

		vg := vgs[0]
		if err == nil {
			_, err = client.DeleteVolumeGroupById(vg.ExtId)
			if err != nil {
				return fmt.Errorf("error: Restored Volume Group still exists: %v", err)
			}
			log.Printf("[DEBUG] Restored Volume Group deleted\n")
			return nil
		}

		return nil
	}
}
