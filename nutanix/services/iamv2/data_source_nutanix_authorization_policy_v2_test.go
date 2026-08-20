package iamv2_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameAuthorizationPolicy = "data.nutanix_authorization_policy_v2.test"

func TestAccV2NutanixAuthorizationPolicyDatasource_Basic(t *testing.T) {
	acpDisplayName := fmt.Sprintf("tf-test-acp-%d", acctest.RandInt())
	roleDisplayName := fmt.Sprintf("tf-test-role-display-name-%d", acctest.RandInt())
	roleDescription := "tf test role description"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAuthorizationPolicyDatasourceV2Config(acpDisplayName, roleDisplayName, roleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameAuthorizationPolicy, "ext_id"),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "display_name", acpDisplayName),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "description", authPolicyDescription),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "authorization_policy_type", authPolicyType),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "identities.#", strconv.Itoa(len(authPolicyIdentities))),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "identities.0.reserved", authPolicyIdentities[0]),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "entities.#", strconv.Itoa(len(authPolicyEntities))),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "entities.0.reserved", authPolicyEntities[0]),
					resource.TestCheckResourceAttr(datasourceNameAuthorizationPolicy, "entities.1.reserved", authPolicyEntities[1]),
				),
			},
		},
	})
}

func testAuthorizationPolicyDatasourceV2Config(acpDisplayName, roleDisplayName, roleDescription string) string {
	return fmt.Sprintf(`
		%s

		data "nutanix_authorization_policy_v2" "test" {
			ext_id = nutanix_authorization_policy_v2.auth_policy_test.id
			depends_on = [nutanix_authorization_policy_v2.auth_policy_test]
		}
	`, authPolicyConfig(acpDisplayName, roleDisplayName, roleDescription))
}
