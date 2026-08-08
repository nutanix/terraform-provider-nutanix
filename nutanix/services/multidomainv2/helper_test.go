package multidomainv2_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/management"
	localdomainReq "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/localdomain"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/projects"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/resourcegroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameProjectV2 = "nutanix_project_v2.test"
const resourceNameResourceGroupV2 = "nutanix_resource_group_v2.test"
const resourceNameLocalDomainV2 = "nutanix_local_domain_v2.test"

// testLocalDomainV2CheckDestroy verifies the local domain is unregistered after
// the test. LocalDomain is a singleton — after unregister its registration state
// must no longer be REGISTERED.
func testLocalDomainV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	client := conn.MultidomainAPI.LocalDomain
	ctx := context.Background()

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_local_domain_v2" {
			continue
		}
		resp, err := client.GetLocalDomain(ctx, &localdomainReq.GetLocalDomainRequest{})
		if err != nil {
			// A fetch error is acceptable — the domain is gone.
			return nil
		}
		domain, ok := resp.Data.GetValue().(management.LocalDomain)
		if !ok {
			return nil
		}
		if domain.RegistrationStatus != nil && domain.RegistrationStatus.State != nil {
			if domain.RegistrationStatus.State.GetName() == "REGISTERED" {
				return fmt.Errorf("local domain still registered: %s", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testProjectV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	client := conn.MultidomainAPI.Projects
	ctx := context.Background()

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_project_v2" {
			continue
		}
		getReq := import1.GetProjectByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := client.GetProjectById(ctx, &getReq)
		if err == nil {
			return fmt.Errorf("project still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testResourceGroupV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	client := conn.MultidomainAPI.ResourceGroups
	ctx := context.Background()

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_resource_group_v2" {
			continue
		}
		getReq := import2.GetResourceGroupByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		_, err := client.GetResourceGroupById(ctx, &getReq)
		if err == nil {
			return fmt.Errorf("resource group still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func printAllStateAttributes(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		fmt.Printf("\n=== State attributes for %s ===\n", resourceName)
		for key, value := range rs.Primary.Attributes {
			fmt.Printf("  %s = %s\n", key, value)
		}
		fmt.Printf("=== End state attributes for %s ===\n\n", resourceName)
		return nil
	}
}

func checkAttributeLength(resourceName, attribute string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		attrKey := fmt.Sprintf("%s.#", attribute)
		countStr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}
		c, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", attrKey, err)
		}
		if c < minLength {
			return fmt.Errorf("expected %s >= %d, got %d", attrKey, minLength, c)
		}
		return nil
	}
}
