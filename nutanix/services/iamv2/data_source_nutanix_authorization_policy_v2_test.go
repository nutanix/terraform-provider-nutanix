package iamv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameAuthorizationPolicy = "data.nutanix_authorization_policy_v2.test"

func TestAccV2NutanixAuthorizationPolicyDatasource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPolicyDatasourceV2Config(filepath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameAuthorizationPolicy, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "display_name", testVars.Iam.AuthPolicies.DisplayName),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "description", testVars.Iam.AuthPolicies.Description),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "authorization_policy_type", testVars.Iam.AuthPolicies.AuthPolicyType),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "identities.#", strconv.Itoa(len(testVars.Iam.AuthPolicies.Identities))),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "identities.0.reserved", testVars.Iam.AuthPolicies.Identities[0]),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "entities.#", strconv.Itoa(len(testVars.Iam.AuthPolicies.Entities))),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "entities.0.reserved", testVars.Iam.AuthPolicies.Entities[0]),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "entities.1.reserved", testVars.Iam.AuthPolicies.Entities[1]),
				),
			},
		},
	})
}

func testAuthorizationPolicyDatasourceV2Config(filepath string) string {
	return fmt.Sprintf(`

		locals{
			config = (jsondecode(file("%s")))
			auth_policies = local.config.iam.auth_policies
			roles = local.config.iam.roles
		}

		%s
			
		data "nutanix_authorization_policy_v2" "test" {
			ext_id = nutanix_authorization_policy_v2.auth_policy_test.id
			depends_on = [nutanix_authorization_policy_v2.auth_policy_test]
		}

		
	`, filepath, authPolicy)
}
