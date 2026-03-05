package dataprotectionv2_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/ahv/config"
	volumesClient "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v17/models/volumes/v4/config"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/request/vm"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/datapolicies-go-client/v17/models/datapolicies/v4/request/protectionpolicies"
	import3 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/categories"
	import4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/volumes-go-client/v17/models/volumes/v4/request/volumegroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func waitForVMToBeProtected(resourceName, attributeName, desiredValue string, maxRetries int, retryInterval, sleepTime time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var lastValue string
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		client := conn.VmmAPI.VMAPIInstance
		ctx := context.Background()

		for i := 0; i < maxRetries; i++ {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("resource not found: %s", resourceName)
			}

			getVmByIdRequest := import1.GetVmByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			vmResp, err := client.GetVmById(ctx, &getVmByIdRequest)
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
	ctx := context.Background()

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
		getVmByIdRequest := import1.GetVmByIdRequest{
			ExtId: utils.StringPtr(vmExtID),
		}
		readResp, err := vmClient.GetVmById(ctx, &getVmByIdRequest)
		if err == nil {
			args := make(map[string]interface{})
			etag := vmClient.ApiClient.GetEtag(readResp)
			args["If-Match"] = utils.StringPtr(etag)
			deleteVmByIdRequest := import1.DeleteVmByIdRequest{
				ExtId: utils.StringPtr(vmExtID),
			}
			_, err = vmClient.DeleteVmById(ctx, &deleteVmByIdRequest, args)
			if err != nil {
				return fmt.Errorf("error: VM still exists: %v", err)
			}
			log.Printf("[DEBUG] VM deleted\n")
		}
	}

	// delete protection policy
	if ppExtID != "" {
		getProtectionPolicyByIdRequest := import2.GetProtectionPolicyByIdRequest{
			ExtId: utils.StringPtr(ppExtID),
		}
		_, err := ppClient.GetProtectionPolicyById(ctx, &getProtectionPolicyByIdRequest)
		if err == nil {
			log.Printf("[DEBUG] Protection Policy still exists")
			deleteProtectionPolicyByIdRequest := import2.DeleteProtectionPolicyByIdRequest{
				ExtId: utils.StringPtr(ppExtID),
			}
			_, err = ppClient.DeleteProtectionPolicyById(ctx, &deleteProtectionPolicyByIdRequest)
			if err != nil {
				return fmt.Errorf("error: Protection Policy still exists : %v", err)
			}
			log.Printf("[DEBUG] Protection Policy deleted\n")
		}
	}

	// delete category
	if categoryExtID != "" {
		getCategoryByIdRequest := import3.GetCategoryByIdRequest{
			ExtId: utils.StringPtr(categoryExtID),
		}
		_, err := categoryClient.GetCategoryById(ctx, &getCategoryByIdRequest)
		if err == nil {
			log.Printf("[DEBUG] Category still exists")
			deleteCategoryByIdRequest := import3.DeleteCategoryByIdRequest{
				ExtId: utils.StringPtr(categoryExtID),
			}
			_, err = categoryClient.DeleteCategoryById(ctx, &deleteCategoryByIdRequest)
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
			ctx := context.Background()

			getVmByIdRequest := import1.GetVmByIdRequest{
				ExtId: utils.StringPtr(vmExtID),
			}
			readResp, err := client.GetVmById(ctx, &getVmByIdRequest)

			if err == nil {
				args := make(map[string]interface{})
				etag := client.ApiClient.GetEtag(readResp)
				args["If-Match"] = utils.StringPtr(etag)

				deleteVmByIdRequest := import1.DeleteVmByIdRequest{
					ExtId: utils.StringPtr(vmExtID),
				}
				_, err = client.DeleteVmById(ctx, &deleteVmByIdRequest, args)
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
		ctx := context.Background()

		filter := fmt.Sprintf("startswith(name, '%s')", vmName)

		listVmsRequest := import1.ListVmsRequest{
			Filter_: utils.StringPtr(filter),
		}
		resp, err := client.ListVms(ctx, &listVmsRequest)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if resp.Data == nil {
			return fmt.Errorf("no data returned from list vms on Remote site")
		}
		vms := resp.Data.GetValue().([]config.Vm)

		vm := vms[0]
		getVmByIdRequest := import1.GetVmByIdRequest{
			ExtId: vm.ExtId,
		}
		readResp, err := client.GetVmById(ctx, &getVmByIdRequest)
		if err == nil {
			args := make(map[string]interface{})
			etag := client.ApiClient.GetEtag(readResp)
			args["If-Match"] = utils.StringPtr(etag)
			deleteVmByIdRequest := import1.DeleteVmByIdRequest{
				ExtId: vm.ExtId,
			}
			_, err = client.DeleteVmById(ctx, &deleteVmByIdRequest, args)
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
		ctx := context.Background()
		filter := fmt.Sprintf("startswith(name, '%s')", vgName)

		listVolumeGroupsRequest := import4.ListVolumeGroupsRequest{
			Filter_: utils.StringPtr(filter),
		}
		resp, err := client.ListVolumeGroups(ctx, &listVolumeGroupsRequest)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if resp.Data == nil {
			return fmt.Errorf("no data returned from list Volume Groups on Remote site")
		}
		vgs := resp.Data.GetValue().([]volumesClient.VolumeGroup)

		vg := vgs[0]
		if err == nil {
			deleteVolumeGroupByIdRequest := import4.DeleteVolumeGroupByIdRequest{
				ExtId: vg.ExtId,
			}
			_, err = client.DeleteVolumeGroupById(ctx, &deleteVolumeGroupByIdRequest)
			if err != nil {
				return fmt.Errorf("error: Restored Volume Group still exists: %v", err)
			}
			log.Printf("[DEBUG] Restored Volume Group deleted\n")
			return nil
		}

		return nil
	}
}
