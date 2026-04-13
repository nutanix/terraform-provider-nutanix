package multidomainv2_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/projects"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/multidomain-go-client/v17/models/multidomain/v4/request/resourcegroups"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameProjectV2        = "nutanix_project_v2.test"
const resourceNameResourceGroupV2  = "nutanix_resource_group_v2.test"

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
