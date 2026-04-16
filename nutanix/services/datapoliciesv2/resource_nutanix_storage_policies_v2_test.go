package datapoliciesv2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameStoragePolicy = "nutanix_storage_policy_v2.test"
const dataSourceNameStoragePolicy = "data.nutanix_storage_policy_v2.test"
const dataSourceNameStoragePolicies = "data.nutanix_storage_policies_v2.test"

func TestAccV2NutanixStoragePolicyResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "compression_spec.0.compression_state", "POSTPROCESS"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "encryption_spec.0.encryption_state", "ENABLED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "1000"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "THREE"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "policy_type", "USER"),
				),
			},
			{
				Config: testStoragePolicyResourceConfig(name) + testStoragePolicyDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicy, "compression_spec.0.compression_state", "POSTPROCESS"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicy, "encryption_spec.0.encryption_state", "ENABLED"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicy, "qos_spec.0.throttled_iops", "1000"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "THREE"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicy, "policy_type", "USER"),
				),
			},
			{
				Config: testStoragePolicyResourceConfig(name) + testStoragePoliciesDatasourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameStoragePolicies, "storage_policies.#"),
					checkAttributeLength(dataSourceNameStoragePolicies, "storage_policies", 1),
					resource.TestCheckResourceAttrSet(dataSourceNameStoragePolicies, "storage_policies.0.ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicies, "storage_policies.0.name", name),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicies, "storage_policies.0.compression_spec.0.compression_state", "POSTPROCESS"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicies, "storage_policies.0.encryption_spec.0.encryption_state", "ENABLED"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicies, "storage_policies.0.qos_spec.0.throttled_iops", "1000"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicies, "storage_policies.0.fault_tolerance_spec.0.replication_factor", "THREE"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicies, "storage_policies.0.policy_type", "USER"),
					resource.TestCheckResourceAttr(dataSourceNameStoragePolicies, "total_available_results", "1"),
				),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_WithSystemDerived(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-sys-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceWithSystemDerivedConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "compression_spec.0.compression_state", "SYSTEM_DERIVED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "encryption_spec.0.encryption_state", "SYSTEM_DERIVED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "100"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "SYSTEM_DERIVED"),
				),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_WithInlineCompression(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-inline-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceWithInlineCompressionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "compression_spec.0.compression_state", "INLINE"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "encryption_spec.0.encryption_state", "ENABLED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "5000"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "TWO"),
				),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_WithDisabledCompression(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-disabled-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceWithDisabledCompressionConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "compression_spec.0.compression_state", "DISABLED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "encryption_spec.0.encryption_state", "SYSTEM_DERIVED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "2000"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "SYSTEM_DERIVED"),
				),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_Update(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-update-%d", r)
	updatedName := fmt.Sprintf("tf-test-storage-policy-updated-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "compression_spec.0.compression_state", "POSTPROCESS"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "encryption_spec.0.encryption_state", "ENABLED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "1000"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "THREE"),
				),
			},
			{
				Config: testStoragePolicyResourceUpdateConfig(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", updatedName),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "compression_spec.0.compression_state", "INLINE"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "2000"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "TWO"),
				),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_UpdateEncryptionState(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-update-encryption-%d", r)
	updatedName := fmt.Sprintf("tf-test-storage-policy-updated-encryption-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "encryption_spec.0.encryption_state", "ENABLED"),
				),
			},
			{
				Config:      testStoragePolicyResourceUpdateEncryptionStateConfig(updatedName, "SYSTEM_DERIVED"),
				ExpectError: regexp.MustCompile("Encryption value cannot be changed once enabled because it is not supported."),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_WithCategories(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-cat-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceWithCategoriesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "category_ext_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "100"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "compression_spec.0.compression_state", "SYSTEM_DERIVED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "encryption_spec.0.encryption_state", "SYSTEM_DERIVED"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "fault_tolerance_spec.0.replication_factor", "SYSTEM_DERIVED"),
				),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_MinimalConfig(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-min-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStoragePolicyResourceMinimalConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameStoragePolicy, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "name", name),
					resource.TestCheckResourceAttr(resourceNameStoragePolicy, "qos_spec.0.throttled_iops", "100"),
				),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_InvalidWithoutName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStoragePolicyResourceWithoutNameConfig(),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_InvalidWithoutQosSpec(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-min-%d", r)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStoragePolicyResourceWithoutQosSpecConfig(name),
				ExpectError: regexp.MustCompile("qos_spec must be provided when compression_state, encryption_state, and replication_factor are all SYSTEM_DERIVED"),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_InvalidCompressionState(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-invalid-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStoragePolicyResourceWithInvalidCompressionStateConfig(name),
				ExpectError: regexp.MustCompile("got INVALID_STATE"),
			},
		},
	})
}

func TestAccV2NutanixStoragePolicyResource_InvalidThrottledIops(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-policy-iops-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStoragePolicyResourceWithInvalidThrottledIopsConfig(name),
				ExpectError: regexp.MustCompile("Numeric instance is lower than the required minimum"),
			},
		},
	})
}

// Test configuration functions

func testStoragePolicyResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  compression_spec {
    compression_state = "POSTPROCESS"
  }

  encryption_spec {
    encryption_state = "ENABLED"
  }

  qos_spec {
    throttled_iops = 1000
  }

  fault_tolerance_spec {
    replication_factor = "THREE"
  }
}`, name)
}

func testStoragePolicyResourceWithSystemDerivedConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  compression_spec {
    compression_state = "SYSTEM_DERIVED"
  }

  encryption_spec {
    encryption_state = "SYSTEM_DERIVED"
  }

  qos_spec {
    throttled_iops = 100
  }

  fault_tolerance_spec {
    replication_factor = "SYSTEM_DERIVED"
  }
}`, name)
}

func testStoragePolicyResourceWithInlineCompressionConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  compression_spec {
    compression_state = "INLINE"
  }

  encryption_spec {
    encryption_state = "ENABLED"
  }

  qos_spec {
    throttled_iops = 5000
  }

  fault_tolerance_spec {
    replication_factor = "TWO"
  }
}`, name)
}

func testStoragePolicyResourceWithDisabledCompressionConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  compression_spec {
    compression_state = "DISABLED"
  }

  qos_spec {
    throttled_iops = 2000
  }
}`, name)
}

func testStoragePolicyResourceUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  compression_spec {
    compression_state = "INLINE"
  }

  encryption_spec {
    encryption_state = "ENABLED"
  }

  qos_spec {
    throttled_iops = 2000
  }

  fault_tolerance_spec {
    replication_factor = "TWO"
  }
}`, name)
}

func testStoragePolicyResourceWithCategoriesConfig(name string) string {
	return fmt.Sprintf(`
data "nutanix_categories_v2" "category-list" {}
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"
  category_ext_ids = [data.nutanix_categories_v2.category-list.categories.0.ext_id, data.nutanix_categories_v2.category-list.categories.1.ext_id]
  qos_spec {
    throttled_iops = 100
  }
}`, name)
}

func testStoragePolicyResourceMinimalConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  qos_spec {
    throttled_iops = 100
  }
}`, name)
}

func testStoragePolicyResourceWithoutNameConfig() string {
	return `
resource "nutanix_storage_policy_v2" "test" {
  qos_spec {
    throttled_iops = 100
  }
}`
}

func testStoragePolicyResourceWithInvalidCompressionStateConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  compression_spec {
    compression_state = "INVALID_STATE"
  }

  qos_spec {
    throttled_iops = 100
  }
}`, name)
}

func testStoragePolicyResourceWithInvalidThrottledIopsConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"

  qos_spec {
    throttled_iops = 50
  }
}`, name)
}

func testStoragePolicyResourceWithoutQosSpecConfig(name string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"
}`, name)
}

func testStoragePolicyResourceUpdateEncryptionStateConfig(name, encryptionState string) string {
	return fmt.Sprintf(`
resource "nutanix_storage_policy_v2" "test" {
  name = "%s"
	compression_spec {
    compression_state = "POSTPROCESS"
  }

  encryption_spec {
    encryption_state = "%s"
  }

  qos_spec {
    throttled_iops = 1000
  }

  fault_tolerance_spec {
    replication_factor = "THREE"
  }
}`, name, encryptionState)
}

func testStoragePolicyDatasourceConfig() string {
	return `

data "nutanix_storage_policy_v2" "test" {
	ext_id = nutanix_storage_policy_v2.test.ext_id
	depends_on = [nutanix_storage_policy_v2.test]
}

`
}

func testStoragePoliciesDatasourceConfig(name string) string {
	return fmt.Sprintf(`

data "nutanix_storage_policies_v2" "test" {
	filter = "name eq '%s'"
	depends_on = [nutanix_storage_policy_v2.test]
}

`, name)
}
