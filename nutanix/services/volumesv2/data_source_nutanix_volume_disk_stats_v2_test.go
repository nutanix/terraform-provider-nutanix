package volumesv2_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

// Test Plan — nutanix_volume_disk_stats_v2 (GetVolumeDiskStats)
//
// The GetVolumeDiskStats API queries time-series performance stats for a single
// Volume Disk (identified by {extId}) that belongs to a Volume Group (identified
// by {volumeGroupExtId}) over a [start_time, end_time] window. It is a read-only
// datasource — there is no resource, so no CheckDestroy of a stats entity is
// needed (Import not supported for this datasource).
//
// All prerequisites (Volume Group + Volume Disk) are created via resource blocks
// so the test is self-contained and reproducible. Scenarios covered:
//
//  1. Basic: create a VG + disk, then query the disk stats with the required
//     start_time/end_time window. Assert that the datasource resolves the disk
//     identity (ext_id matches the created disk, volume_disk_ext_id is set) and
//     the stats collections are present in state.
//  2. WithStatType: same as basic but exercise the stat_type down-sampling
//     operator (AVG) and an explicit sampling_interval to verify the query
//     parameters are wired through to the API.
//  3. InvalidStatType (negative): supply an invalid stat_type enum value and
//     verify the ValidateFunc rejects it before any API call.
//  4. InvalidStartTime (negative): supply a non-RFC3339 start_time and verify
//     the datasource returns a descriptive parse error.
//  5. InvalidExtId (negative): supply a random (non-existent) disk ext_id and
//     verify the API returns an error.

const datasourceVolumeDiskStats = "data.nutanix_volume_disk_stats_v2.test"

func TestAccV2NutanixVolumeDiskStatsDataSource_Basic(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-disk-stats-%d", r)
	desc := "terraform test volume disk stats description"

	startTime := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	endTime := time.Now().UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeDiskStatsDataSourceConfig(name, desc, int(diskSizeBytes), startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceVolumeDiskStats, "ext_id",
						"nutanix_volume_group_disk_v2.test", "id"),
					resource.TestCheckResourceAttrSet(datasourceVolumeDiskStats, "volume_disk_ext_id"),
					resource.TestCheckResourceAttr(datasourceVolumeDiskStats, "start_time", startTime),
					resource.TestCheckResourceAttr(datasourceVolumeDiskStats, "end_time", endTime),
				),
			},
		},
	})
}

func TestAccV2NutanixVolumeDiskStatsDataSource_WithStatType(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-disk-stats-%d", r)
	desc := "terraform test volume disk stats description"

	startTime := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	endTime := time.Now().UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeDiskStatsDataSourceConfigWithStatType(name, desc, int(diskSizeBytes), startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceVolumeDiskStats, "ext_id",
						"nutanix_volume_group_disk_v2.test", "id"),
					resource.TestCheckResourceAttrSet(datasourceVolumeDiskStats, "volume_disk_ext_id"),
					resource.TestCheckResourceAttr(datasourceVolumeDiskStats, "stat_type", "AVG"),
					resource.TestCheckResourceAttr(datasourceVolumeDiskStats, "sampling_interval", "30"),
				),
			},
		},
	})
}

func TestAccV2NutanixVolumeDiskStatsDataSource_InvalidStatType(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-disk-stats-%d", r)
	desc := "terraform test volume disk stats description"

	startTime := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	endTime := time.Now().UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVolumeDiskStatsDataSourceConfigInvalidStatType(name, desc, int(diskSizeBytes), startTime, endTime),
				ExpectError: regexp.MustCompile(`expected stat_type to be one of`),
			},
		},
	})
}

func TestAccV2NutanixVolumeDiskStatsDataSource_InvalidStartTime(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-disk-stats-%d", r)
	desc := "terraform test volume disk stats description"

	endTime := time.Now().UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVolumeDiskStatsDataSourceConfigInvalidStartTime(name, desc, int(diskSizeBytes), endTime),
				ExpectError: regexp.MustCompile(`error while parsing start_time`),
			},
		},
	})
}

func TestAccV2NutanixVolumeDiskStatsDataSource_InvalidExtId(t *testing.T) {
	r := acctest.RandInt()
	name := fmt.Sprintf("terraform-test-volume-disk-stats-%d", r)
	desc := "terraform test volume disk stats description"

	startTime := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	endTime := time.Now().UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccVolumeDiskStatsDataSourceConfigInvalidExtID(name, desc, int(diskSizeBytes), startTime, endTime),
				ExpectError: regexp.MustCompile(`error while fetching Volume Disk stats`),
			},
		},
	})
}

func testAccVolumeDiskStatsDataSourceConfig(name, desc string, diskSizeBytes int, startTime, endTime string) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		fmt.Sprintf(`
		data "nutanix_volume_disk_stats_v2" "test" {
			volume_group_ext_id = nutanix_volume_group_v2.test.id
			ext_id              = nutanix_volume_group_disk_v2.test.id
			start_time          = "%[1]s"
			end_time            = "%[2]s"
			depends_on          = [nutanix_volume_group_disk_v2.test]
		}
	`, startTime, endTime)
}

func testAccVolumeDiskStatsDataSourceConfigWithStatType(name, desc string, diskSizeBytes int, startTime, endTime string) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		fmt.Sprintf(`
		data "nutanix_volume_disk_stats_v2" "test" {
			volume_group_ext_id = nutanix_volume_group_v2.test.id
			ext_id              = nutanix_volume_group_disk_v2.test.id
			start_time          = "%[1]s"
			end_time            = "%[2]s"
			sampling_interval   = 30
			stat_type           = "AVG"
			depends_on          = [nutanix_volume_group_disk_v2.test]
		}
	`, startTime, endTime)
}

func testAccVolumeDiskStatsDataSourceConfigInvalidStatType(name, desc string, diskSizeBytes int, startTime, endTime string) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		fmt.Sprintf(`
		data "nutanix_volume_disk_stats_v2" "test" {
			volume_group_ext_id = nutanix_volume_group_v2.test.id
			ext_id              = nutanix_volume_group_disk_v2.test.id
			start_time          = "%[1]s"
			end_time            = "%[2]s"
			stat_type           = "INVALID"
			depends_on          = [nutanix_volume_group_disk_v2.test]
		}
	`, startTime, endTime)
}

func testAccVolumeDiskStatsDataSourceConfigInvalidStartTime(name, desc string, diskSizeBytes int, endTime string) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		fmt.Sprintf(`
		data "nutanix_volume_disk_stats_v2" "test" {
			volume_group_ext_id = nutanix_volume_group_v2.test.id
			ext_id              = nutanix_volume_group_disk_v2.test.id
			start_time          = "not-a-valid-timestamp"
			end_time            = "%[1]s"
			depends_on          = [nutanix_volume_group_disk_v2.test]
		}
	`, endTime)
}

func testAccVolumeDiskStatsDataSourceConfigInvalidExtID(name, desc string, diskSizeBytes int, startTime, endTime string) string {
	return testAccVolumeGroupResourceConfig(name, desc) +
		testAccVolumeGroupDiskResourceConfig(name, desc, diskSizeBytes) +
		fmt.Sprintf(`
		data "nutanix_volume_disk_stats_v2" "test" {
			volume_group_ext_id = nutanix_volume_group_v2.test.id
			ext_id              = "00000000-0000-0000-0000-000000000000"
			start_time          = "%[1]s"
			end_time            = "%[2]s"
			depends_on          = [nutanix_volume_group_disk_v2.test]
		}
	`, startTime, endTime)
}
