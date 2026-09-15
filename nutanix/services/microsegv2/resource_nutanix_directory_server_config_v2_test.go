package microsegv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// ---------------------------------------------------------------------------
// Directory Server Config — Create / Update / Delete
// ---------------------------------------------------------------------------

func TestAccV2DirectoryServerConfig_CreateContainsMatchType(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServerConfigV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "directory_service_reference", dsRef),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.#", "1"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_entity", "VM"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_field", "NAME"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_type", "CONTAINS"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.criteria", "DeveloperVM"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "domain_controllers.#", "1"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "domain_controllers.0.ipv4.0.value", dcIP),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "is_default_category_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login", "false"),
				),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_CreateAllMatchType(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_All(dsRef, dcIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameDirectoryServerConfigV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "is_default_category_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login", "false"),
				),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_UpdateContainsToAll(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_type", "CONTAINS"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.criteria", "DeveloperVM"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "is_default_category_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login", "false"),
				),
			},
			{
				Config: testAccDSCConfig_All(dsRef, dcIP),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_type", "ALL"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "is_default_category_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login", "false"),
				),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_UpdateCriteria(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.criteria", "DeveloperVM"),
				),
			},
			{
				Config: testAccDSCConfig_Contains(dsRef, dcIP, "ProductionVM", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "matching_criterias.0.criteria", "ProductionVM"),
				),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_UpdateDefaultCategoryFlags(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "is_default_category_enabled", "false"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login", "false"),
				),
			},
			{
				Config: testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", true, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "is_default_category_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login", "false"),
				),
			},
			{
				Config: testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "is_default_category_enabled", "true"),
					resource.TestCheckResourceAttr(resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login", "true"),
				),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_ImportState(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", false, false),
			},
			{
				ResourceName:      resourceNameDirectoryServerConfigV2,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_Datasources_Contains(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	singularDS := "data.nutanix_directory_server_config_v2.test"
	pluralDS := "data.nutanix_directory_server_configs_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_WithDatasources_Contains(dsRef, dcIP, "DeveloperVM", true, false),
				Check: resource.ComposeTestCheckFunc(
					// --- Singular datasource: pair every field with the resource ---
					resource.TestCheckResourceAttrPair(singularDS, "ext_id", resourceNameDirectoryServerConfigV2, "ext_id"),
					resource.TestCheckResourceAttrPair(singularDS, "directory_service_reference", resourceNameDirectoryServerConfigV2, "directory_service_reference"),
					resource.TestCheckResourceAttrPair(singularDS, "is_default_category_enabled", resourceNameDirectoryServerConfigV2, "is_default_category_enabled"),
					resource.TestCheckResourceAttrPair(singularDS, "should_keep_default_category_on_login", resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.#", resourceNameDirectoryServerConfigV2, "matching_criterias.#"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.0.match_entity", resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_entity"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.0.match_field", resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_field"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.0.match_type", resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_type"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.0.criteria", resourceNameDirectoryServerConfigV2, "matching_criterias.0.criteria"),
					resource.TestCheckResourceAttrPair(singularDS, "domain_controllers.#", resourceNameDirectoryServerConfigV2, "domain_controllers.#"),
					resource.TestCheckResourceAttrPair(singularDS, "domain_controllers.0.ipv4.0.value", resourceNameDirectoryServerConfigV2, "domain_controllers.0.ipv4.0.value"),

					// --- Singular: explicit value checks for CONTAINS ---
					resource.TestCheckResourceAttr(singularDS, "directory_service_reference", dsRef),
					resource.TestCheckResourceAttr(singularDS, "matching_criterias.0.match_entity", "VM"),
					resource.TestCheckResourceAttr(singularDS, "matching_criterias.0.match_field", "NAME"),
					resource.TestCheckResourceAttr(singularDS, "matching_criterias.0.match_type", "CONTAINS"),
					resource.TestCheckResourceAttr(singularDS, "matching_criterias.0.criteria", "DeveloperVM"),
					resource.TestCheckResourceAttr(singularDS, "domain_controllers.0.ipv4.0.value", dcIP),
					resource.TestCheckResourceAttr(singularDS, "is_default_category_enabled", "true"),
					resource.TestCheckResourceAttr(singularDS, "should_keep_default_category_on_login", "false"),

					// --- Plural datasource: list-level checks ---
					resource.TestCheckResourceAttrSet(pluralDS, "directory_server_configs.#"),
					resource.TestCheckResourceAttrSet(pluralDS, "directory_server_configs.0.ext_id"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.directory_service_reference", dsRef),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.is_default_category_enabled", "true"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.should_keep_default_category_on_login", "false"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.#", "1"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.0.match_entity", "VM"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.0.match_field", "NAME"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.0.match_type", "CONTAINS"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.0.criteria", "DeveloperVM"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.domain_controllers.#", "1"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.domain_controllers.0.ipv4.0.value", dcIP),
				),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_Datasources_All(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	singularDS := "data.nutanix_directory_server_config_v2.test"
	pluralDS := "data.nutanix_directory_server_configs_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testDirectoryServerConfigV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccDSCConfig_WithDatasources_All(dsRef, dcIP),
				Check: resource.ComposeTestCheckFunc(
					// --- Singular datasource: pair every field with the resource ---
					resource.TestCheckResourceAttrPair(singularDS, "ext_id", resourceNameDirectoryServerConfigV2, "ext_id"),
					resource.TestCheckResourceAttrPair(singularDS, "directory_service_reference", resourceNameDirectoryServerConfigV2, "directory_service_reference"),
					resource.TestCheckResourceAttrPair(singularDS, "is_default_category_enabled", resourceNameDirectoryServerConfigV2, "is_default_category_enabled"),
					resource.TestCheckResourceAttrPair(singularDS, "should_keep_default_category_on_login", resourceNameDirectoryServerConfigV2, "should_keep_default_category_on_login"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.0.match_entity", resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_entity"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.0.match_field", resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_field"),
					resource.TestCheckResourceAttrPair(singularDS, "matching_criterias.0.match_type", resourceNameDirectoryServerConfigV2, "matching_criterias.0.match_type"),
					resource.TestCheckResourceAttrPair(singularDS, "domain_controllers.0.ipv4.0.value", resourceNameDirectoryServerConfigV2, "domain_controllers.0.ipv4.0.value"),

					// --- Singular: explicit value checks for ALL ---
					resource.TestCheckResourceAttr(singularDS, "directory_service_reference", dsRef),
					resource.TestCheckResourceAttr(singularDS, "matching_criterias.0.match_entity", "VM"),
					resource.TestCheckResourceAttr(singularDS, "matching_criterias.0.match_field", "NAME"),
					resource.TestCheckResourceAttr(singularDS, "matching_criterias.0.match_type", "ALL"),
					resource.TestCheckResourceAttr(singularDS, "domain_controllers.0.ipv4.0.value", dcIP),
					resource.TestCheckResourceAttr(singularDS, "is_default_category_enabled", "false"),
					resource.TestCheckResourceAttr(singularDS, "should_keep_default_category_on_login", "false"),

					// --- Plural datasource: list-level checks ---
					resource.TestCheckResourceAttrSet(pluralDS, "directory_server_configs.#"),
					resource.TestCheckResourceAttrSet(pluralDS, "directory_server_configs.0.ext_id"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.directory_service_reference", dsRef),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.is_default_category_enabled", "false"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.should_keep_default_category_on_login", "false"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.#", "1"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.0.match_entity", "VM"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.0.match_field", "NAME"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.matching_criterias.0.match_type", "ALL"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.domain_controllers.#", "1"),
					resource.TestCheckResourceAttr(pluralDS, "directory_server_configs.0.domain_controllers.0.ipv4.0.value", dcIP),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Directory Server Config — Negative tests
// ---------------------------------------------------------------------------

func TestAccV2DirectoryServerConfig_InvalidMatchType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDSCConfig_InvalidMatchType(),
				ExpectError: regexp.MustCompile(`expected matching_criterias\.0\.match_type to be one of`),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_CriteriaWithAllMatchType(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDSCConfig_CriteriaWithAll(dsRef, dcIP),
				ExpectError: regexp.MustCompile(`'criteria' must not be set when match_type is "ALL"`),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_DefaultCategoryEnabledWithAllMatchType(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDSCConfig_DefaultCategoryEnabledWithAll(dsRef, dcIP),
				ExpectError: regexp.MustCompile(`'is_default_category_enabled' must be false when match_type is "ALL"`),
			},
		},
	})
}

func TestAccV2DirectoryServerConfig_KeepOnLoginWithoutDefaultCategoryEnabled(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", false, true),
				ExpectError: regexp.MustCompile(`'should_keep_default_category_on_login' can only be true when 'is_default_category_enabled' is also true`),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Category Mapping — Create / Update / Delete (depends on Directory Server Config)
// ---------------------------------------------------------------------------

func TestAccV2CategoryMapping_CreateAndUpdate(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCategoryMappingV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccCategoryMappingConfig_Create(dsRef, dcIP, "dnd_approval_group_3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameCategoryMappingV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameCategoryMappingV2, "name", "dnd_approval_group_3"),
					resource.TestCheckResourceAttr(resourceNameCategoryMappingV2, "category_name", "ADGroup"),
					resource.TestCheckResourceAttr(resourceNameCategoryMappingV2, "category_value", "dnd_approval_group_3"),
					resource.TestCheckResourceAttr(resourceNameCategoryMappingV2, "ad_info.#", "1"),
					resource.TestCheckResourceAttr(resourceNameCategoryMappingV2, "ad_info.0.directory_service_reference", dsRef),
					resource.TestCheckResourceAttrSet(resourceNameCategoryMappingV2, "ad_info.0.object_identifier"),
				),
			},
			{
				Config: testAccCategoryMappingConfig_Update(dsRef, dcIP, "dnd_approval_group_3", "dnd_approval_group_3_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameCategoryMappingV2, "category_value", "dnd_approval_group_3_updated"),
				),
			},
		},
	})
}

func TestAccV2CategoryMapping_MultipleMappings(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCategoryMappingV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config: testAccCategoryMappingConfig_Multiple(dsRef, dcIP,
					"dnd_approval_group_3", "dnd_approval_group_4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameCategoryMappingV2, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameCategoryMappingV2, "name", "dnd_approval_group_3"),
					resource.TestCheckResourceAttrSet("nutanix_ad_group_category_mapping_v2.test2", "ext_id"),
					resource.TestCheckResourceAttr("nutanix_ad_group_category_mapping_v2.test2", "name", "dnd_approval_group_4"),
				),
			},
		},
	})
}

func TestAccV2CategoryMapping_ImportState(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCategoryMappingV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccCategoryMappingConfig_Create(dsRef, dcIP, "dnd_approval_group_3"),
			},
			{
				ResourceName:      resourceNameCategoryMappingV2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ad_info.0.object_path",
				},
				Config: testAccCategoryMappingImportConfig(dsRef, dcIP, "dnd_approval_group_3"),
			},
		},
	})
}

func TestAccV2CategoryMapping_Datasources(t *testing.T) {
	dsRef := testVars.Iam.DirectoryServicesMain.SecondaryAD.ExtID
	dcIP := testVars.Iam.DirectoryServicesMain.SecondaryAD.DNS
	groupName := "dnd_approval_group_3"

	singularDS := "data.nutanix_ad_group_category_mapping_v2.test"
	pluralDS := "data.nutanix_ad_group_category_mappings_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testCategoryMappingV2CheckDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: tearDownAll,
				Config:    testAccCategoryMappingConfig_WithDatasources(dsRef, dcIP, groupName),
				Check: resource.ComposeTestCheckFunc(
					// --- Singular datasource: pair every field with the resource ---
					resource.TestCheckResourceAttrPair(singularDS, "ext_id", resourceNameCategoryMappingV2, "ext_id"),
					resource.TestCheckResourceAttrPair(singularDS, "name", resourceNameCategoryMappingV2, "name"),
					resource.TestCheckResourceAttrPair(singularDS, "category_name", resourceNameCategoryMappingV2, "category_name"),
					resource.TestCheckResourceAttrPair(singularDS, "category_value", resourceNameCategoryMappingV2, "category_value"),
					resource.TestCheckResourceAttrPair(singularDS, "ad_info.#", resourceNameCategoryMappingV2, "ad_info.#"),
					resource.TestCheckResourceAttrPair(singularDS, "ad_info.0.directory_service_reference", resourceNameCategoryMappingV2, "ad_info.0.directory_service_reference"),
					resource.TestCheckResourceAttrPair(singularDS, "ad_info.0.object_identifier", resourceNameCategoryMappingV2, "ad_info.0.object_identifier"),

					// --- Singular: explicit value checks ---
					resource.TestCheckResourceAttr(singularDS, "name", groupName),
					resource.TestCheckResourceAttr(singularDS, "category_name", "ADGroup"),
					resource.TestCheckResourceAttr(singularDS, "category_value", groupName),
					resource.TestCheckResourceAttr(singularDS, "ad_info.0.directory_service_reference", dsRef),
					resource.TestCheckResourceAttrSet(singularDS, "ad_info.0.object_identifier"),
					resource.TestCheckResourceAttrSet(singularDS, "ad_info.0.object_path"),

					// --- Plural datasource: list-level checks ---
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.#"),
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.0.ext_id"),
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.0.name"),
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.0.category_name"),
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.0.category_value"),
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.0.ad_info.0.directory_service_reference"),
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.0.ad_info.0.object_identifier"),
					resource.TestCheckResourceAttrSet(pluralDS, "category_mappings.0.ad_info.0.object_path"),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// Category Mapping — Negative tests
// ---------------------------------------------------------------------------

func TestAccV2CategoryMapping_MissingCategoryValue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCategoryMappingConfig_MissingCategoryValue(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2CategoryMapping_MissingAdInfo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCategoryMappingConfig_MissingAdInfo(),
				ExpectError: regexp.MustCompile("Insufficient ad_info blocks"),
			},
		},
	})
}

func TestAccV2CategoryMapping_MissingName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCategoryMappingConfig_MissingName(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

// ===========================================================================
// Config generators — Directory Server Config
// ===========================================================================

func testAccDSCConfig_Contains(dsRef, dcIP, criteria string, defaultEnabled, keepOnLogin bool) string {
	return fmt.Sprintf(`
resource "nutanix_directory_server_config_v2" "test" {
  directory_service_reference           = "%[1]s"
  is_default_category_enabled           = %[4]t
  should_keep_default_category_on_login = %[5]t

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "CONTAINS"
    criteria     = "%[3]s"
  }

  domain_controllers {
    ipv4 {
      value = "%[2]s"
    }
  }
}
`, dsRef, dcIP, criteria, defaultEnabled, keepOnLogin)
}

func testAccDSCConfig_All(dsRef, dcIP string) string {
	return fmt.Sprintf(`
resource "nutanix_directory_server_config_v2" "test" {
  directory_service_reference           = "%[1]s"
  is_default_category_enabled           = false
  should_keep_default_category_on_login = false

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "ALL"
  }

  domain_controllers {
    ipv4 {
      value = "%[2]s"
    }
  }
}
`, dsRef, dcIP)
}

func testAccDSCDatasourceSuffix() string {
	return `
data "nutanix_directory_server_config_v2" "test" {
  ext_id     = nutanix_directory_server_config_v2.test.ext_id
  depends_on = [nutanix_directory_server_config_v2.test]
}

data "nutanix_directory_server_configs_v2" "test" {
  depends_on = [nutanix_directory_server_config_v2.test]
}
`
}

func testAccDSCConfig_WithDatasources_Contains(dsRef, dcIP, criteria string, defaultEnabled, keepOnLogin bool) string {
	return testAccDSCConfig_Contains(dsRef, dcIP, criteria, defaultEnabled, keepOnLogin) + testAccDSCDatasourceSuffix()
}

func testAccDSCConfig_WithDatasources_All(dsRef, dcIP string) string {
	return testAccDSCConfig_All(dsRef, dcIP) + testAccDSCDatasourceSuffix()
}

func testAccDSCConfig_InvalidMatchType() string {
	return `
resource "nutanix_directory_server_config_v2" "test" {
  directory_service_reference = "00000000-0000-0000-0000-000000000000"

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "INVALID"
  }

  domain_controllers {
    ipv4 {
      value = "10.0.0.1"
    }
  }
}
`
}

func testAccDSCConfig_CriteriaWithAll(dsRef, dcIP string) string {
	return fmt.Sprintf(`
resource "nutanix_directory_server_config_v2" "test" {
  directory_service_reference           = "%[1]s"
  is_default_category_enabled           = false
  should_keep_default_category_on_login = false

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "ALL"
    criteria     = "should-not-be-set"
  }

  domain_controllers {
    ipv4 {
      value = "%[2]s"
    }
  }
}
`, dsRef, dcIP)
}

func testAccDSCConfig_DefaultCategoryEnabledWithAll(dsRef, dcIP string) string {
	return fmt.Sprintf(`
resource "nutanix_directory_server_config_v2" "test" {
  directory_service_reference           = "%[1]s"
  is_default_category_enabled           = true
  should_keep_default_category_on_login = false

  matching_criterias {
    match_entity = "VM"
    match_field  = "NAME"
    match_type   = "ALL"
  }

  domain_controllers {
    ipv4 {
      value = "%[2]s"
    }
  }
}
`, dsRef, dcIP)
}

// ===========================================================================
// Config generators — Category Mapping (depends on Directory Server Config)
// ===========================================================================

func testAccCategoryMappingBaseConfig(dsRef, dcIP string) string {
	return testAccDSCConfig_Contains(dsRef, dcIP, "DeveloperVM", false, false)
}

func testAccCategoryMappingSearchBlock(dsRef, groupName string) string {
	return fmt.Sprintf(`
data "nutanix_directory_service_users_search_v2" "%[2]s" {
  directory_service_ext_id = "%[1]s"
  query                    = "%[2]s"
  is_wildcard_search       = true
}

locals {
  %[2]s_object_guid = one(flatten([
    for result in data.nutanix_directory_service_users_search_v2.%[2]s.search_results : [
      for attr in result.attributes :
      attr.values[0] if attr.name == "objectGUID"
    ]
  ]))
}
`, dsRef, groupName)
}

// testAccCategoryMappingImportConfig provides a minimal config for the import
// step without data sources. The SDK's ImportStateVerify iterates all resources
// in the post-import state and tries to match them against the pre-import state
// but skips data-source entries only in the old state, not the new one. When the
// full config's data sources are resolved during import they end up in the new
// state and verification fails. Using a stripped config avoids this.
func testAccCategoryMappingImportConfig(dsRef, dcIP, groupName string) string {
	return testAccCategoryMappingBaseConfig(dsRef, dcIP) + fmt.Sprintf(`
resource "nutanix_ad_group_category_mapping_v2" "test" {
  name           = "%[1]s"
  category_value = "%[1]s"

  ad_info {
    directory_service_reference = "%[2]s"
    object_identifier           = "placeholder"
  }

  depends_on = [nutanix_directory_server_config_v2.test]
}
`, groupName, dsRef)
}

func testAccCategoryMappingConfig_Create(dsRef, dcIP, groupName string) string {
	return testAccCategoryMappingBaseConfig(dsRef, dcIP) +
		testAccCategoryMappingSearchBlock(dsRef, groupName) +
		fmt.Sprintf(`
resource "nutanix_ad_group_category_mapping_v2" "test" {
  name           = "%[1]s"
  category_value = "%[1]s"

  ad_info {
    directory_service_reference = "%[2]s"
    object_identifier           = local.%[1]s_object_guid
  }

  depends_on = [nutanix_directory_server_config_v2.test]
}
`, groupName, dsRef)
}

func testAccCategoryMappingConfig_Update(dsRef, dcIP, groupName, newCategoryValue string) string {
	return testAccCategoryMappingBaseConfig(dsRef, dcIP) +
		testAccCategoryMappingSearchBlock(dsRef, groupName) +
		fmt.Sprintf(`
resource "nutanix_ad_group_category_mapping_v2" "test" {
  name           = "%[1]s"
  category_value = "%[4]s"

  ad_info {
    directory_service_reference = "%[2]s"
    object_identifier           = local.%[1]s_object_guid
  }

  depends_on = [nutanix_directory_server_config_v2.test]
}
`, groupName, dsRef, groupName, newCategoryValue)
}

func testAccCategoryMappingConfig_Multiple(dsRef, dcIP, group1, group2 string) string {
	return testAccCategoryMappingBaseConfig(dsRef, dcIP) +
		testAccCategoryMappingSearchBlock(dsRef, group1) +
		testAccCategoryMappingSearchBlock(dsRef, group2) +
		fmt.Sprintf(`
resource "nutanix_ad_group_category_mapping_v2" "test" {
  name           = "%[1]s"
  category_value = "%[1]s"

  ad_info {
    directory_service_reference = "%[3]s"
    object_identifier           = local.%[1]s_object_guid
  }

  depends_on = [nutanix_directory_server_config_v2.test]
}

resource "nutanix_ad_group_category_mapping_v2" "test2" {
  name           = "%[2]s"
  category_value = "%[2]s"

  ad_info {
    directory_service_reference = "%[3]s"
    object_identifier           = local.%[2]s_object_guid
  }

  depends_on = [nutanix_directory_server_config_v2.test]
}
`, group1, group2, dsRef)
}

func testAccCategoryMappingConfig_WithDatasources(dsRef, dcIP, groupName string) string {
	return testAccCategoryMappingConfig_Create(dsRef, dcIP, groupName) + `
data "nutanix_ad_group_category_mapping_v2" "test" {
  ext_id     = nutanix_ad_group_category_mapping_v2.test.ext_id
  depends_on = [nutanix_ad_group_category_mapping_v2.test]
}

data "nutanix_ad_group_category_mappings_v2" "test" {
  depends_on = [nutanix_ad_group_category_mapping_v2.test]
}
`
}

func testAccCategoryMappingConfig_MissingCategoryValue() string {
	return `
resource "nutanix_ad_group_category_mapping_v2" "test" {
  name = "missing-category-value"

  ad_info {
    directory_service_reference = "00000000-0000-0000-0000-000000000000"
    object_identifier           = "11111111-1111-1111-1111-111111111111"
  }
}
`
}

func testAccCategoryMappingConfig_MissingAdInfo() string {
	return `
resource "nutanix_ad_group_category_mapping_v2" "test" {
  name           = "missing-ad-info"
  category_value = "test-value"
}
`
}

func testAccCategoryMappingConfig_MissingName() string {
	return `
resource "nutanix_ad_group_category_mapping_v2" "test" {
  category_value = "test-value"

  ad_info {
    directory_service_reference = "00000000-0000-0000-0000-000000000000"
    object_identifier           = "11111111-1111-1111-1111-111111111111"
  }
}
`
}
