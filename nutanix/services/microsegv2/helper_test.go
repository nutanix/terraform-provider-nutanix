package microsegv2_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	import2 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/config"
	dscReq "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/directoryserverconfigs"
	import1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/entitygroups"
	nspRequest "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/microseg/v4/request/networksecuritypolicies"
	prismMicroseg "github.com/nutanix-core/ntnx-api-golang-sdk-internal/microseg-go-client/v17/models/prism/v4/config"
	prismConfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/config"
	categoryRequest "github.com/nutanix-core/ntnx-api-golang-sdk-internal/prism-go-client/v17/models/prism/v4/request/categories"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/microseg"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v4/prism"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const resourceNameEntityGroupV2 = "nutanix_entity_group_v2.test"
const resourceNameNetworkSecurityPolicyImportV2 = "nutanix_network_security_policy_import_v2.test"
const resourceNameNetworkSecurityPolicyExportV2 = "nutanix_network_security_policy_export_v2.test"

func testEntityGroupV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	api := conn.MicroSegAPI.EntityGroupsAPIInstance

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_entity_group_v2" {
			getEntityGroupByIDRequest := import1.GetEntityGroupByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			_, err := api.GetEntityGroupById(context.Background(), &getEntityGroupByIDRequest)
			if err == nil {
				return fmt.Errorf("entity group v2 still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

// Import not supported for the network security policy import/export resources — they are
// action resources backed by async tasks with no GetById equivalent.

// testNetworkSecurityPolicyImportV2CheckDestroy verifies that destroying the configuration
// cleaned up every entity touched by the import flow: the policies created by the import
// action (tracked in imported_policy_ext_ids and deleted on destroy), any directly
// provisioned policies, and finally the categories they reference.
func testNetworkSecurityPolicyImportV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	nspClient := conn.MicroSegAPI.NetworkingSecurityInstance
	categoryClient := conn.PrismAPI.CategoriesAPIInstance

	// 1. Policies created by the import action must be destroyed.
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_network_security_policy_import_v2" {
			continue
		}
		count, _ := strconv.Atoi(rs.Primary.Attributes["imported_policy_ext_ids.#"])
		for i := 0; i < count; i++ {
			extID := rs.Primary.Attributes[fmt.Sprintf("imported_policy_ext_ids.%d", i)]
			if extID == "" {
				continue
			}
			getReq := nspRequest.GetNetworkSecurityPolicyByIdRequest{
				ExtId: utils.StringPtr(extID),
			}
			if _, err := nspClient.GetNetworkSecurityPolicyById(ctx, &getReq); err == nil {
				return fmt.Errorf("imported network security policy %s still exists", extID)
			} else if !isEntityNotFoundErr(err) {
				return fmt.Errorf("error checking if imported network security policy %s exists: %v", extID, err)
			}
		}
	}

	// 2. Any directly provisioned policies must be destroyed.
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_network_security_policy_v2" {
			continue
		}
		getReq := nspRequest.GetNetworkSecurityPolicyByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := nspClient.GetNetworkSecurityPolicyById(ctx, &getReq); err == nil {
			return fmt.Errorf("network security policy %s still exists", rs.Primary.ID)
		} else if !isEntityNotFoundErr(err) {
			return fmt.Errorf("error checking if network security policy %s exists: %v", rs.Primary.ID, err)
		}
	}

	// 3. Categories must be destroyed (checked last, after the policies that reference them).
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_category_v2" {
			continue
		}
		getReq := categoryRequest.GetCategoryByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := categoryClient.GetCategoryById(ctx, &getReq); err == nil {
			return fmt.Errorf("category %s still exists", rs.Primary.ID)
		} else if !isEntityNotFoundErr(err) {
			return fmt.Errorf("error checking if category %s exists: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

// testNetworkSecurityPolicyExportV2CheckDestroy verifies that destroying the configuration
// cleaned up the underlying entities. The export resource itself is a no-op on destroy, but
// the test configuration also provisions network security policies and categories. Policies
// are checked first; categories are checked last because they can only be removed once the
// policies referencing them are gone.
func testNetworkSecurityPolicyExportV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	ctx := context.Background()

	nspClient := conn.MicroSegAPI.NetworkingSecurityInstance
	categoryClient := conn.PrismAPI.CategoriesAPIInstance

	// 1. Network security policies must be destroyed.
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_network_security_policy_v2" {
			continue
		}
		getReq := nspRequest.GetNetworkSecurityPolicyByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := nspClient.GetNetworkSecurityPolicyById(ctx, &getReq); err == nil {
			return fmt.Errorf("network security policy %s still exists", rs.Primary.ID)
		} else if !isEntityNotFoundErr(err) {
			return fmt.Errorf("error checking if network security policy %s exists: %v", rs.Primary.ID, err)
		}
	}

	// 2. Categories must be destroyed (checked last, after the policies that reference them).
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "nutanix_category_v2" {
			continue
		}
		getReq := categoryRequest.GetCategoryByIdRequest{
			ExtId: utils.StringPtr(rs.Primary.ID),
		}
		if _, err := categoryClient.GetCategoryById(ctx, &getReq); err == nil {
			return fmt.Errorf("category %s still exists", rs.Primary.ID)
		} else if !isEntityNotFoundErr(err) {
			return fmt.Errorf("error checking if category %s exists: %v", rs.Primary.ID, err)
		}
	}

	return nil
}

// isEntityNotFoundErr reports whether an API error indicates the entity no longer exists
// (the expected outcome after a successful destroy).
func isEntityNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "not found") ||
		strings.Contains(s, "does not exist") ||
		strings.Contains(s, "entity_not_found") ||
		strings.Contains(s, "could not be fetched") ||
		strings.Contains(s, "valid metadata or attributes could not be fetched") ||
		strings.Contains(s, "rbac_authorization_error") ||
		strings.Contains(s, "unauthorized") ||
		strings.Contains(s, "plat-10007") ||
		strings.Contains(s, "mic-10013") ||
		strings.Contains(s, "failed to authorize")
}

// testAccCaptureImportedPolicyIDs records the imported policy ext_ids from state so a
// later teardown step can verify they were removed from the cluster.
func testAccCaptureImportedPolicyIDs(importRes string, policyIDs *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		imp, ok := s.RootModule().Resources[importRes]
		if !ok {
			return fmt.Errorf("resource not found: %s", importRes)
		}
		count, _ := strconv.Atoi(imp.Primary.Attributes["imported_policy_ext_ids.#"])
		ids := make([]string, 0, count)
		for i := 0; i < count; i++ {
			ids = append(ids, imp.Primary.Attributes[fmt.Sprintf("imported_policy_ext_ids.%d", i)])
		}
		*policyIDs = ids
		return nil
	}
}

// testAccCheckImportEntitiesNotCreated asserts that a dry-run import created nothing on the
// cluster: no category exists for the given key (the import would have created one had it not
// been a dry run).
func testAccCheckImportEntitiesNotCreated(categoryKey string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		ctx := context.Background()
		categoryClient := conn.PrismAPI.CategoriesAPIInstance

		listReq := categoryRequest.ListCategoriesRequest{
			Filter_: utils.StringPtr(fmt.Sprintf("key eq '%s'", categoryKey)),
		}
		listResp, err := categoryClient.ListCategories(ctx, &listReq)
		if err != nil {
			return fmt.Errorf("error listing categories with key %s: %v", categoryKey, err)
		}
		if listResp.Data == nil {
			return nil
		}
		if cats := listResp.Data.GetValue().([]prismConfig.Category); len(cats) > 0 {
			return fmt.Errorf("dry-run import created %d category(ies) with key %q; expected none", len(cats), categoryKey)
		}
		return nil
	}
}

// testAccCheckImportedEntitiesDestroyed verifies the imported policies are gone, then removes
// the category that the import recreated from the exported file. That category is created by
// the import action itself (not by Terraform), so it is never in Terraform state and nothing
// destroys it automatically. We locate it by name (key) via a filtered list, then delete it so
// it does not linger on the cluster after the test.
func testAccCheckImportedEntitiesDestroyed(policyIDs *[]string, categoryKey string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		conn := acc.TestAccProvider.Meta().(*conns.Client)
		ctx := context.Background()

		nspClient := conn.MicroSegAPI.NetworkingSecurityInstance
		categoryClient := conn.PrismAPI.CategoriesAPIInstance

		for _, id := range *policyIDs {
			if id == "" {
				continue
			}
			getReq := nspRequest.GetNetworkSecurityPolicyByIdRequest{
				ExtId: utils.StringPtr(id),
			}
			if _, err := nspClient.GetNetworkSecurityPolicyById(ctx, &getReq); err == nil {
				return fmt.Errorf("imported network security policy %s still exists after teardown", id)
			} else if !isEntityNotFoundErr(err) {
				return fmt.Errorf("error checking if imported network security policy %s exists: %v", id, err)
			}
		}

		if categoryKey == "" {
			return nil
		}

		// Search and filter on the category name (key) to find the id(s) created by the import.
		listReq := categoryRequest.ListCategoriesRequest{
			Filter_: utils.StringPtr(fmt.Sprintf("key eq '%s'", categoryKey)),
		}
		listResp, err := categoryClient.ListCategories(ctx, &listReq)
		if err != nil {
			return fmt.Errorf("error listing categories with key %s: %v", categoryKey, err)
		}
		if listResp.Data == nil {
			return nil
		}

		categories := listResp.Data.GetValue().([]prismConfig.Category)
		for _, c := range categories {
			if c.ExtId == nil {
				continue
			}
			// A category that is still shared with a project cannot be deleted (and would
			// also block the project's own deletion), so unshare it from every project first.
			if err := unshareCategoryFromAllProjects(ctx, conn, *c.ExtId); err != nil {
				return err
			}
			delReq := categoryRequest.DeleteCategoryByIdRequest{
				ExtId: c.ExtId,
			}
			if _, err := categoryClient.DeleteCategoryById(ctx, &delReq); err != nil && !isEntityNotFoundErr(err) {
				return fmt.Errorf("error deleting category %s (key %s) left behind by import: %v", *c.ExtId, categoryKey, err)
			}
		}

		return nil
	}
}

// unshareCategoryFromAllProjects removes every project association from a category so it can be
// deleted (and so the projects it was shared with can themselves be deleted). It waits for each
// unshare task to complete and tolerates a category that has already been removed.
func unshareCategoryFromAllProjects(ctx context.Context, conn *conns.Client, categoryExtID string) error {
	categoryClient := conn.PrismAPI.CategoriesAPIInstance

	getReq := categoryRequest.GetCategoryByIdRequest{ExtId: utils.StringPtr(categoryExtID)}
	getResp, err := categoryClient.GetCategoryById(ctx, &getReq)
	if err != nil {
		if isEntityNotFoundErr(err) {
			return nil
		}
		return fmt.Errorf("error fetching category %s before unshare: %v", categoryExtID, err)
	}

	category, ok := getResp.Data.GetValue().(prismConfig.Category)
	if !ok {
		return nil
	}

	for _, projectID := range category.SharedWithProjects {
		// Re-fetch to obtain a fresh ETag for the If-Match precondition each iteration.
		latest, err := categoryClient.GetCategoryById(ctx, &getReq)
		if err != nil {
			if isEntityNotFoundErr(err) {
				return nil
			}
			return fmt.Errorf("error fetching category %s before unshare: %v", categoryExtID, err)
		}
		headers := map[string]interface{}{
			"If-Match": utils.StringPtr(categoryClient.ApiClient.GetEtag(latest)),
		}

		unshareReq := categoryRequest.UnshareCategoryRequest{
			CategoryExtId: utils.StringPtr(categoryExtID),
			Body: &prismConfig.UnshareCategoryRequest{
				ProjectExtId: utils.StringPtr(projectID),
			},
		}
		unshareResp, err := categoryClient.UnshareCategory(ctx, &unshareReq, headers)
		if err != nil {
			if isEntityNotFoundErr(err) {
				continue
			}
			return fmt.Errorf("error unsharing category %s from project %s: %v", categoryExtID, projectID, err)
		}

		taskRef, ok := unshareResp.Data.GetValue().(prismConfig.TaskReference)
		if !ok || taskRef.ExtId == nil {
			continue
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING", "RUNNING", "QUEUED"},
			Target:  []string{"SUCCEEDED"},
			Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, conn.PrismAPI, utils.StringValue(taskRef.ExtId)),
			Timeout: 5 * time.Minute,
		}
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("error waiting for unshare task of category %s (project %s): %v", categoryExtID, projectID, err)
		}
	}

	return nil
}

// testAccCheckExportFileValid asserts that the export was downloaded to the file_path
// recorded in state: the path attribute is set, the file exists on disk, and it is non-empty.
func testAccCheckExportFileValid(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		statePath := rs.Primary.Attributes["file_path"]
		if statePath == "" {
			return fmt.Errorf("file_path is empty in state for %s", resourceName)
		}

		info, err := os.Stat(statePath)
		if err != nil {
			return fmt.Errorf("export file was not downloaded to the provided file path (%s): %w", statePath, err)
		}
		if info.Size() == 0 {
			return fmt.Errorf("export file (%s) is empty", statePath)
		}

		return nil
	}
}

const (
	resourceNameDirectoryServerConfigV2 = "nutanix_directory_server_config_v2.test"
	resourceNameCategoryMappingV2       = "nutanix_ad_group_category_mapping_v2.test"
)

func testDirectoryServerConfigV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	api := conn.MicroSegAPI.DirectoryServerConfigsAPIInstance

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_directory_server_config_v2" {
			getRequest := dscReq.GetDirectoryServerConfigByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			_, err := api.GetDirectoryServerConfigById(context.Background(), &getRequest)
			if err == nil {
				return fmt.Errorf("directory server config v2 still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testCategoryMappingV2CheckDestroy(state *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)
	api := conn.MicroSegAPI.DirectoryServerConfigsAPIInstance

	for _, rs := range state.RootModule().Resources {
		if rs.Type == "nutanix_ad_group_category_mapping_v2" {
			getRequest := dscReq.GetDsCategoryMappingByIdRequest{
				ExtId: utils.StringPtr(rs.Primary.ID),
			}
			_, err := api.GetDsCategoryMappingById(context.Background(), &getRequest)
			if err == nil {
				return fmt.Errorf("category mapping v2 still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

// getOrCreateConn returns the provider's Client if available, otherwise
// creates a standalone client from environment variables. This allows
// teardown to work even before the provider has been configured (e.g. the
// very first test in a run).
func getOrCreateConn() (*conns.Client, error) {
	if acc.TestAccProvider.Meta() != nil {
		return acc.TestAccProvider.Meta().(*conns.Client), nil
	}

	creds := client.Credentials{
		URL:      fmt.Sprintf("%s:%s", os.Getenv("NUTANIX_ENDPOINT"), os.Getenv("NUTANIX_PORT")),
		Endpoint: os.Getenv("NUTANIX_ENDPOINT"),
		Username: os.Getenv("NUTANIX_USERNAME"),
		Password: os.Getenv("NUTANIX_PASSWORD"),
		Port:     os.Getenv("NUTANIX_PORT"),
		Insecure: os.Getenv("NUTANIX_INSECURE") == "true",
	}
	microsegClient, err := microseg.NewMicrosegClient(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to create microseg client: %v", err)
	}
	prismClient, err := prism.NewPrismClient(creds)
	if err != nil {
		return nil, fmt.Errorf("failed to create prism client: %v", err)
	}
	return &conns.Client{
		MicroSegAPI: microsegClient,
		PrismAPI:    prismClient,
	}, nil
}

// tearDownCategoryMappings deletes all existing category mappings.
// Call this before tests to ensure a clean slate.
func tearDownCategoryMappings() {
	conn, err := getOrCreateConn()
	if err != nil {
		log.Printf("[TEARDOWN] warning: %v", err)
		return
	}
	api := conn.MicroSegAPI.DirectoryServerConfigsAPIInstance
	ctx := context.Background()

	listResp, err := api.ListCategoryMappings(ctx, &dscReq.ListCategoryMappingsRequest{})
	if err != nil {
		log.Printf("[TEARDOWN] warning: failed to list category mappings: %v", err)
		return
	}
	if listResp.Data == nil {
		return
	}
	mappings, ok := listResp.Data.GetValue().([]import2.CategoryMapping)
	if !ok || len(mappings) == 0 {
		return
	}

	for _, m := range mappings {
		extID := utils.StringValue(m.ExtId)
		log.Printf("[TEARDOWN] deleting category mapping %s", extID)

		getReq := dscReq.GetDsCategoryMappingByIdRequest{ExtId: utils.StringPtr(extID)}
		getResp, err := api.GetDsCategoryMappingById(ctx, &getReq)
		if err != nil {
			log.Printf("[TEARDOWN] warning: failed to get category mapping %s for etag: %v", extID, err)
			continue
		}
		etag := api.ApiClient.GetEtag(getResp)
		args := map[string]interface{}{"If-Match": utils.StringPtr(etag)}

		delReq := dscReq.DeleteDsCategoryMappingByIdRequest{ExtId: utils.StringPtr(extID)}
		resp, err := api.DeleteDsCategoryMappingById(ctx, &delReq, args)
		if err != nil {
			log.Printf("[TEARDOWN] warning: failed to delete category mapping %s: %v", extID, err)
			continue
		}
		waitForTask(resp.Data, conn)
	}
}

// tearDownDirectoryServerConfigs deletes all existing directory server configs.
// Call this before tests to ensure a clean slate.
func tearDownDirectoryServerConfigs() {
	conn, err := getOrCreateConn()
	if err != nil {
		log.Printf("[TEARDOWN] warning: %v", err)
		return
	}
	api := conn.MicroSegAPI.DirectoryServerConfigsAPIInstance
	ctx := context.Background()

	listResp, err := api.ListDirectoryServerConfigs(ctx, &dscReq.ListDirectoryServerConfigsRequest{})
	if err != nil {
		log.Printf("[TEARDOWN] warning: failed to list directory server configs: %v", err)
		return
	}
	if listResp.Data == nil {
		return
	}
	configs, ok := listResp.Data.GetValue().([]import2.DirectoryServerConfig)
	if !ok || len(configs) == 0 {
		return
	}

	for _, c := range configs {
		extID := utils.StringValue(c.ExtId)
		log.Printf("[TEARDOWN] deleting directory server config %s", extID)

		delReq := dscReq.DeleteDirectoryServerConfigByIdRequest{ExtId: utils.StringPtr(extID)}
		resp, err := api.DeleteDirectoryServerConfigById(ctx, &delReq)
		if err != nil {
			log.Printf("[TEARDOWN] warning: failed to delete directory server config %s: %v", extID, err)
			continue
		}
		waitForTask(resp.Data, conn)
	}
}

// tearDownAll cleans up both category mappings and directory server configs.
// Category mappings must be deleted first since they depend on the config.
func tearDownAll() {
	tearDownCategoryMappings()
	tearDownDirectoryServerConfigs()
}

func waitForTask(data interface{ GetValue() interface{} }, conn *conns.Client) {
	if data == nil {
		return
	}
	taskRef, ok := data.GetValue().(prismMicroseg.TaskReference)
	if !ok {
		return
	}
	ctx := context.Background()
	taskConn := conn.PrismAPI
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING", "QUEUED"},
		Target:  []string{"SUCCEEDED"},
		Refresh: common.TaskStateRefreshPrismTaskGroupFunc(ctx, taskConn, utils.StringValue(taskRef.ExtId)),
		Timeout: 5 * time.Minute,
	}
	_, _ = stateConf.WaitForStateContext(ctx)
}
