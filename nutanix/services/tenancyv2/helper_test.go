package tenancyv2_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	iamAuthn "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/authn"
	iamUserGroupsReq "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/request/usergroups"
	iamUsersReq "github.com/nutanix/ntnx-api-golang-clients/iam-go-client/v4/models/iam/v4/request/users"
	import1 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/config"
	import3 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/request/membership"

	import4 "github.com/nutanix/ntnx-api-golang-clients/tenancy-go-client/v4/models/tenancy/v4/request/tenants"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func checkAttributeLength(resourceName, attribute string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		attrKey := fmt.Sprintf("%s.#", attribute)
		attr, ok := rs.Primary.Attributes[attrKey]
		if !ok {
			return fmt.Errorf("attribute %s not found", attrKey)
		}
		count, err := strconv.Atoi(attr)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %s", attrKey, err)
		}
		if count < minLength {
			return fmt.Errorf("expected %s to be >= %d, got %d", attrKey, minLength, count)
		}
		return nil
	}
}

// tenantMembershipTestInputs holds the infrastructure identifiers required to
// run the live tenant-membership lifecycle test. These reference SP-tenancy
// entities (an allocated domain, a role, and a principal) that cannot be
// created by this provider, so they are supplied via environment variables.
type tenantMembershipTestInputs struct {
	DomainExtID           string
	RoleExtID             string
	IdentityExtID         string
	IdentityProviderExtID string
	PrincipalValue        string
	IdentityType          string
}

// tenantMembershipInputsFromEnv reads the membership test inputs from the
// environment and skips the test when the mandatory ones are missing.
func tenantMembershipInputsFromEnv(t *testing.T) tenantMembershipTestInputs {
	t.Helper()

	in := tenantMembershipTestInputs{
		DomainExtID:           os.Getenv("NUTANIX_DOMAIN_EXT_ID"),
		RoleExtID:             os.Getenv("NUTANIX_MEMBERSHIP_ROLE_EXT_ID"),
		IdentityExtID:         os.Getenv("NUTANIX_MEMBERSHIP_IDENTITY_EXT_ID"),
		IdentityProviderExtID: os.Getenv("NUTANIX_MEMBERSHIP_IDP_EXT_ID"),
		PrincipalValue:        os.Getenv("NUTANIX_MEMBERSHIP_PRINCIPAL_VALUE"),
		IdentityType:          os.Getenv("NUTANIX_MEMBERSHIP_IDENTITY_TYPE"),
	}

	if in.DomainExtID == "" || in.RoleExtID == "" {
		t.Skip("skipping tenant membership lifecycle test: set NUTANIX_DOMAIN_EXT_ID and NUTANIX_MEMBERSHIP_ROLE_EXT_ID (plus either NUTANIX_MEMBERSHIP_IDENTITY_EXT_ID or the identity_details env vars) to run it")
	}
	if in.IdentityExtID == "" && (in.IdentityProviderExtID == "" || in.PrincipalValue == "") {
		t.Skip("skipping tenant membership lifecycle test: provide NUTANIX_MEMBERSHIP_IDENTITY_EXT_ID or both NUTANIX_MEMBERSHIP_IDP_EXT_ID and NUTANIX_MEMBERSHIP_PRINCIPAL_VALUE")
	}
	if in.IdentityType == "" {
		in.IdentityType = "USER"
	}
	return in
}

func testAccCheckNutanixTenantMembershipV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client).ServiceProviderManagementAPI
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_tenant_membership_v2" {
			continue
		}
		getRequest := import3.GetDomainMembershipByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := conn.TenantMembershipAPIInstance.GetDomainMembershipById(ctx, &getRequest); err == nil {
			return fmt.Errorf("tenant membership %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// testAccCleanupDirectoryServiceIdentities deletes IAM users and user groups
// still associated with the directory service identified by resourceName.
// Membership deletion leaves those identities behind; deleting the directory
// service while they exist returns IAM-21016.
func testAccCleanupDirectoryServiceIdentities(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		dsExtID := rs.Primary.ID
		if v := rs.Primary.Attributes["ext_id"]; v != "" {
			dsExtID = v
		}
		if dsExtID == "" {
			return fmt.Errorf("%s has empty ext_id", resourceName)
		}
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		return cleanupDirectoryServiceIdentities(context.Background(), conn, dsExtID)
	}
}

func cleanupDirectoryServiceIdentities(ctx context.Context, conn *conns.Client, dsExtID string) error {
	filter := fmt.Sprintf("idpId eq '%s'", dsExtID)
	log.Printf("[DEBUG] Cleaning IAM identities for directory service %s", dsExtID)

	if err := deleteIAMUsersByIdp(ctx, conn, filter); err != nil {
		return err
	}
	if err := deleteIAMUserGroupsByIdp(ctx, conn, filter); err != nil {
		return err
	}
	return nil
}

func deleteIAMUsersByIdp(ctx context.Context, conn *conns.Client, filter string) error {
	listReq := iamUsersReq.ListUsersRequest{
		Filter_: utils.StringPtr(filter),
		Limit_:  utils.IntPtr(100),
	}
	resp, err := conn.IamAPI.UsersAPIInstance.ListUsers(ctx, &listReq)
	if err != nil {
		return fmt.Errorf("list users for directory service cleanup: %w", err)
	}
	if resp == nil || resp.Data == nil {
		return nil
	}
	users, ok := resp.Data.GetValue().([]iamAuthn.User)
	if !ok {
		return fmt.Errorf("unexpected ListUsers data type %T", resp.Data.GetValue())
	}
	for _, user := range users {
		if user.ExtId == nil || *user.ExtId == "" {
			continue
		}
		if err := deleteIAMUserByExtID(ctx, conn, *user.ExtId); err != nil {
			return err
		}
	}
	return nil
}

func deleteIAMUserGroupsByIdp(ctx context.Context, conn *conns.Client, filter string) error {
	listReq := iamUserGroupsReq.ListUserGroupsRequest{
		Filter_: utils.StringPtr(filter),
		Limit_:  utils.IntPtr(100),
	}
	resp, err := conn.IamAPI.UserGroupsAPIInstance.ListUserGroups(ctx, &listReq)
	if err != nil {
		return fmt.Errorf("list user groups for directory service cleanup: %w", err)
	}
	if resp == nil || resp.Data == nil {
		return nil
	}
	groups, ok := resp.Data.GetValue().([]iamAuthn.UserGroup)
	if !ok {
		return fmt.Errorf("unexpected ListUserGroups data type %T", resp.Data.GetValue())
	}
	for _, group := range groups {
		if group.ExtId == nil || *group.ExtId == "" {
			continue
		}
		if err := deleteIAMUserGroupByExtID(ctx, conn, *group.ExtId); err != nil {
			return err
		}
	}
	return nil
}

func deleteIAMUserGroupByExtID(ctx context.Context, conn *conns.Client, extID string) error {
	getReq := iamUserGroupsReq.GetUserGroupByIdRequest{ExtId: utils.StringPtr(extID)}
	readResp, err := conn.IamAPI.UserGroupsAPIInstance.GetUserGroupById(ctx, &getReq)
	if err != nil {
		if isIAMNotFoundErr(err) {
			return nil
		}
		return fmt.Errorf("get user group %s before delete: %w", extID, err)
	}
	etag := conn.IamAPI.UserGroupsAPIInstance.ApiClient.GetEtag(readResp)
	headers := map[string]interface{}{"If-Match": utils.StringPtr(etag)}
	delReq := iamUserGroupsReq.DeleteUserGroupByIdRequest{ExtId: utils.StringPtr(extID)}
	if _, err := conn.IamAPI.UserGroupsAPIInstance.DeleteUserGroupById(ctx, &delReq, headers); err != nil {
		if isIAMNotFoundErr(err) {
			return nil
		}
		// Tenant Admin often lacks IAM v4 user-group delete (IAM-20003); fall back to v3.
		if isIAMUnauthorizedErr(err) && conn.API != nil && conn.API.V3 != nil {
			if _, v3err := conn.API.V3.DeleteUserGroup(ctx, extID); v3err == nil || isIAMNotFoundErr(v3err) {
				log.Printf("[DEBUG] Deleted IAM user group %s via v3 fallback", extID)
				return nil
			} else if v3err != nil {
				return fmt.Errorf("delete user group %s (v4: %v; v3 fallback: %w)", extID, err, v3err)
			}
		}
		return fmt.Errorf("delete user group %s: %w", extID, err)
	}
	log.Printf("[DEBUG] Deleted IAM user group %s", extID)
	return nil
}

// deleteIAMUserByExtID issues DELETE /iam/v4.1.b3/authn/users/{extId}. The IAM
// Go client does not expose DeleteUserById, so this uses CallApiWithContext.
func deleteIAMUserByExtID(ctx context.Context, conn *conns.Client, extID string) error {
	getReq := iamUsersReq.GetUserByIdRequest{ExtId: utils.StringPtr(extID)}
	readResp, err := conn.IamAPI.UsersAPIInstance.GetUserById(ctx, &getReq)
	if err != nil {
		if isIAMNotFoundErr(err) {
			return nil
		}
		return fmt.Errorf("get user %s before delete: %w", extID, err)
	}
	etag := conn.IamAPI.UsersAPIInstance.ApiClient.GetEtag(readResp)
	uri := "/api/iam/v4.1.b3/authn/users/" + url.PathEscape(extID)
	headerParams := map[string]string{}
	if etag != "" {
		headerParams["If-Match"] = etag
	}
	authNames := []string{"apiKeyAuthScheme", "basicAuthScheme"}
	_, err = conn.IamAPI.UsersAPIInstance.ApiClient.CallApiWithContext(
		ctx,
		&uri,
		http.MethodDelete,
		nil,
		url.Values{},
		headerParams,
		url.Values{},
		[]string{"application/json"},
		[]string{},
		authNames,
	)
	if err != nil {
		if isIAMNotFoundErr(err) {
			return nil
		}
		return fmt.Errorf("delete user %s: %w", extID, err)
	}
	log.Printf("[DEBUG] Deleted IAM user %s", extID)
	return nil
}

func isIAMNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	s := fmt.Sprint(err)
	return strings.Contains(s, "404") ||
		strings.Contains(s, "ENTITY_NOT_FOUND") ||
		strings.Contains(s, "RESOURCE_NOT_FOUND") ||
		strings.Contains(s, "does not exist") ||
		strings.Contains(s, "not found")
}

func isIAMUnauthorizedErr(err error) bool {
	if err == nil {
		return false
	}
	s := fmt.Sprint(err)
	return strings.Contains(s, "IAM-20003") ||
		strings.Contains(s, "unauthorized") ||
		strings.Contains(s, "403")
}

// testAccCheckNutanixTenantV2Destroy verifies that every managed tenant has been
// removed. Deleting a tenant is a soft delete: GetTenantById keeps returning the
// tenant with its state set to DELETED (this is the terminal state the resource
// delete itself waits for). So the tenant is considered destroyed when either
// the GET fails (not found) or the returned tenant is in the DELETED state.
func testAccCheckNutanixTenantV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_tenant_v2" {
			continue
		}

		getTenantRequest := import4.GetTenantByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		resp, err := conn.ServiceProviderManagementAPI.TenantManagementAPIInstance.GetTenantById(ctx, &getTenantRequest)
		if err != nil {
			// Not found: the tenant has been fully removed.
			continue
		}
		if resp != nil && resp.Data != nil {
			if tenant, ok := resp.Data.GetValue().(import1.Tenant); ok {
				if tenant.StateInfo != nil && tenant.StateInfo.State != nil &&
					tenant.StateInfo.State.GetName() == "DELETED" {
					// Soft-deleted: treat as destroyed.
					continue
				}
			}
		}
		return fmt.Errorf("tenant %s still exists", rs.Primary.ID)
	}
	return nil
}

// testAccCheckNutanixTenantDomainAllocationV2Destroy verifies that every managed
// domain allocation has been removed by calling GetTenantDomainAllocationById.
func testAccCheckNutanixTenantDomainAllocationV2Destroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_tenant_domain_allocation_v2" {
			continue
		}

		getRequest := import4.GetTenantDomainAllocationByIdRequest{
			TenantExtId: utils.StringPtr(rs.Primary.Attributes["tenant_ext_id"]),
			ExtId:       utils.StringPtr(rs.Primary.ID),
		}
		resp, err := conn.ServiceProviderManagementAPI.TenantManagementAPIInstance.GetTenantDomainAllocationById(ctx, &getRequest)
		if err != nil {
			// A hard error (e.g. not found) means the allocation is gone.
			continue
		}
		// The allocation is soft-deleted: the GET still succeeds but the state
		// transitions to DELETED. Treat that as destroyed.
		if resp != nil && resp.Data != nil {
			if allocation, ok := resp.Data.GetValue().(import1.DomainAllocation); ok {
				if allocation.StateInfo != nil && allocation.StateInfo.State != nil &&
					allocation.StateInfo.State.GetName() == "DELETED" {
					continue
				}
			}
		}
		return fmt.Errorf("tenant domain allocation %s still exists", rs.Primary.ID)
	}
	return nil
}
