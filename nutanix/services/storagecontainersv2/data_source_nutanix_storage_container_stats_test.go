package storagecontainersv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const datasourceNameStorageStatsInfo = "data.nutanix_storage_container_stats_info_v2.test"

func TestAccV2NutanixStorageStatsInfoDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-container-%d", r)

	// Start time is now
	startTime := time.Now()

	// End time is two hours later
	endTime := startTime.Add(2 * time.Hour)

	// Format the times to RFC3339 format
	startTimeFormatted := startTime.UTC().Format(time.RFC3339)
	endTimeFormatted := endTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainerConfig(name) + testStorageStatsDatasourceV2Config(startTimeFormatted, endTimeFormatted),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageStatsInfo, "container_ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixStorageStatsInfoDataSource_SampleInterval(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-container-%d", r)

	// Start time is now
	startTime := time.Now()

	// End time is two hours later
	endTime := startTime.Add(2 * time.Hour)

	// Format the times to RFC3339 format
	startTimeFormatted := startTime.UTC().Format(time.RFC3339)
	endTimeFormatted := endTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainerConfig(name) + testStorageStatsDatasourceV2SampleInterval(startTimeFormatted, endTimeFormatted, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageStatsInfo, "container_ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixStorageStatsInfoDataSource_StatType(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("tf-test-storage-container-%d", r)

	// Start time is now
	startTime := time.Now()

	// End time is two hours later
	endTime := startTime.Add(2 * time.Hour)

	// Format the times to RFC3339 format
	startTimeFormatted := startTime.UTC().Format(time.RFC3339)
	endTimeFormatted := endTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testStorageContainerConfig(name) + testStorageStatsDatasourceV2StatType(startTimeFormatted, endTimeFormatted, "COUNT"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceNameStorageStatsInfo, "container_ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixStorageStatsInfoDataSource_InvalidSampleInterval(t *testing.T) {
	// Start time is now
	startTime := time.Now()

	// End time is two hours later
	endTime := startTime.Add(2 * time.Hour)

	// Format the times to RFC3339 format
	startTimeFormatted := startTime.UTC().Format(time.RFC3339)
	endTimeFormatted := endTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStorageStatsDatasourceV2InvalidSampleInterval(startTimeFormatted, endTimeFormatted, 0),
				ExpectError: regexp.MustCompile("sampling_interval should be greater than 0"),
			},
		},
	})
}

func TestAccV2NutanixStorageStatsInfoDataSource_InvalidStatType(t *testing.T) {
	// Start time is now
	startTime := time.Now()

	// End time is two hours later
	endTime := startTime.Add(2 * time.Hour)

	// Format the times to RFC3339 format
	startTimeFormatted := startTime.UTC().Format(time.RFC3339)
	endTimeFormatted := endTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStorageStatsDatasourceV2InvalidStatType(startTimeFormatted, endTimeFormatted, "INVALID"),
				ExpectError: regexp.MustCompile("running pre-apply refresh"),
			},
		},
	})
}

func TestAccV2NutanixStorageStatsInfoDataSource_MissingRequiredArgs(t *testing.T) {
	// Start time is now
	startTime := time.Now()

	// End time is two hours later
	endTime := startTime.Add(2 * time.Hour)

	// Format the times to RFC3339 format
	startTimeFormatted := startTime.UTC().Format(time.RFC3339)
	endTimeFormatted := endTime.UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testStorageStatsDatasourceV2MissingExtID(startTimeFormatted, endTimeFormatted, "SUM"),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testStorageStatsDatasourceV2MissingStartTime(endTimeFormatted, "SUM"),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testStorageStatsDatasourceV2MissingEndTime(startTimeFormatted, "SUM"),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
		},
	})
}

func testStorageContainerConfig(name string) string {
	return fmt.Sprintf(`

		data "nutanix_clusters_v2" "clusters" {}

		locals{
			cluster = [
				for cluster in data.nutanix_clusters_v2.clusters.cluster_entities :
				cluster.ext_id if cluster.config[0].cluster_function[0] != "PRISM_CENTRAL"
			][0]
		}

		resource "nutanix_storage_containers_v2" "test" {
			name = "%[1]s"
			cluster_ext_id = local.cluster
			logical_advertised_capacity_bytes = 1073741824000
			logical_explicit_reserved_capacity_bytes = 20
			replication_factor = 1
			nfs_whitelist_addresses {
				ipv4  {
					value = "192.168.14.0"
					prefix_length = 32
				}
			}
			erasure_code = "OFF"
			is_inline_ec_enabled = false
			has_higher_ec_fault_domain_preference = false
			cache_deduplication = "OFF"
			on_disk_dedup = "OFF"
			is_compression_enabled = true
			is_internal = false
			is_software_encryption_enabled = false
		}
		
	`, name)
}

func testStorageStatsDatasourceV2Config(startTime, endTime string) string {
	return fmt.Sprintf(`

		
		data "nutanix_storage_container_stats_info_v2" "test" {
			ext_id = nutanix_storage_containers_v2.test.id 
			start_time = "%s"
			end_time = "%s" 
			depends_on = [nutanix_storage_containers_v2.test]
		}

		
	`, startTime, endTime)
}

func testStorageStatsDatasourceV2SampleInterval(startTime, endTime string, sampleInterval int) string {
	return fmt.Sprintf(`

		
		data "nutanix_storage_container_stats_info_v2" "test" {
			ext_id = nutanix_storage_containers_v2.test.id
			start_time = "%s"
			end_time = "%s" 
			sampling_interval = %d
			depends_on = [nutanix_storage_containers_v2.test]
		}

		
	`, startTime, endTime, sampleInterval)
}

func testStorageStatsDatasourceV2StatType(startTime, endTime, statType string) string {
	return fmt.Sprintf(`
		
		data "nutanix_storage_container_stats_info_v2" "test" {
			ext_id = nutanix_storage_containers_v2.test.id 
			start_time = "%s"
			end_time = "%s" 
			stat_type = "%s"
			depends_on = [nutanix_storage_containers_v2.test]
		}
		
	`, startTime, endTime, statType)
}

func testStorageStatsDatasourceV2InvalidSampleInterval(startTime, endTime string, sampleInterval int) string {
	return fmt.Sprintf(`

		
		data "nutanix_storage_container_stats_info_v2" "test" {
			ext_id = "000000-0000000000-00000000"
			start_time = "%s"
			end_time = "%s" 
			sampling_interval = %d
		}

		
	`, startTime, endTime, sampleInterval)
}

func testStorageStatsDatasourceV2InvalidStatType(startTime, endTime, statType string) string {
	return fmt.Sprintf(`
		
		data "nutanix_storage_container_stats_info_v2" "test" {
			ext_id = "000000-0000000000-00000000"
			start_time = "%s"
			end_time = "%s" 
			stat_type = "%s"
		}
		
	`, startTime, endTime, statType)
}

func testStorageStatsDatasourceV2MissingExtID(startTime, endTime, statType string) string {
	return fmt.Sprintf(`
		
		data "nutanix_storage_container_stats_info_v2" "test" {
			start_time = "%s"
			end_time = "%s" 
			stat_type = "%s"
		}
		
	`, startTime, endTime, statType)
}

func testStorageStatsDatasourceV2MissingStartTime(endTime, statType string) string {
	return fmt.Sprintf(`
		
		data "nutanix_storage_container_stats_info_v2" "test" {
			ext_id = "000000-0000000000-00000000"
			end_time = "%s" 
			stat_type = "%s"
		}
		
	`, endTime, statType)
}

func testStorageStatsDatasourceV2MissingEndTime(startTime, statType string) string {
	return fmt.Sprintf(`
		
		data "nutanix_storage_container_stats_info_v2" "test" {
			ext_id = "000000-0000000000-00000000"
			start_time = "%s"
			stat_type = "%s"
		}
		
	`, startTime, statType)
}
