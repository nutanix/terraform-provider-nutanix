package cluster_managementv2_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRsyslogServer = "nutanix_rsyslog_server_v2.test"
const dataSourceNameRsyslogServer = "data.nutanix_rsyslog_server_v2.test"
const dataSourceNameRsyslogServers = "data.nutanix_rsyslog_servers_v2.test"

func TestAccV2NutanixRsyslogServerResource_Basic(t *testing.T) {
	clusterExtID := testVars.ClusterExtID
	if clusterExtID == "" {
		clusterExtID = os.Getenv("NUTANIX_CLUSTER_EXT_ID")
	}
	if clusterExtID == "" {
		t.Skip("NUTANIX_CLUSTER_EXT_ID must be set for this test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceConfig(clusterExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "cluster_ext_id", clusterExtID),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "server_name", "tf-test-rsyslog-server"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "port", "514"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "network_protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "ip_address.0.ipv4.0.value", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.name", "CASSANDRA"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.log_severity_level", "INFO"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.should_log_monitor_files", "true"),
				),
			},
			{
				Config: testRsyslogServerResourceConfigUpdate(clusterExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "cluster_ext_id", clusterExtID),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "server_name", "tf-test-rsyslog-server"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "port", "1514"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "network_protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "ip_address.0.ipv4.0.value", "10.0.0.2"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.name", "PRISM"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.log_severity_level", "WARNING"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.should_log_monitor_files", "false"),
				),
			},
		},
	})
}

func TestAccV2NutanixRsyslogServerDatasource_Basic(t *testing.T) {
	clusterExtID := testVars.ClusterExtID
	if clusterExtID == "" {
		clusterExtID = os.Getenv("NUTANIX_CLUSTER_EXT_ID")
	}
	if clusterExtID == "" {
		t.Skip("NUTANIX_CLUSTER_EXT_ID must be set for this test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceConfig(clusterExtID) + testRsyslogServerDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameRsyslogServer, "server_name", "tf-test-rsyslog-server"),
					resource.TestCheckResourceAttr(dataSourceNameRsyslogServer, "port", "514"),
					resource.TestCheckResourceAttr(dataSourceNameRsyslogServer, "network_protocol", "UDP"),
					resource.TestCheckResourceAttr(dataSourceNameRsyslogServer, "ip_address.0.ipv4.0.value", "10.0.0.1"),
					resource.TestCheckResourceAttr(dataSourceNameRsyslogServer, "modules.0.name", "CASSANDRA"),
					resource.TestCheckResourceAttr(dataSourceNameRsyslogServer, "modules.0.log_severity_level", "INFO"),
				),
			},
		},
	})
}

func TestAccV2NutanixRsyslogServersDatasource_Basic(t *testing.T) {
	clusterExtID := testVars.ClusterExtID
	if clusterExtID == "" {
		clusterExtID = os.Getenv("NUTANIX_CLUSTER_EXT_ID")
	}
	if clusterExtID == "" {
		t.Skip("NUTANIX_CLUSTER_EXT_ID must be set for this test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceConfig(clusterExtID) + testRsyslogServersDatasourceConfig(clusterExtID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameRsyslogServers, "rsyslog_servers.#"),
				),
			},
		},
	})
}

func testRsyslogServerResourceConfig(clusterExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id   = "%s"
  server_name      = "tf-test-rsyslog-server"
  port             = 514
  network_protocol = "UDP"

  ip_address {
    ipv4 {
      value = "10.0.0.1"
    }
  }

  modules {
    name                    = "CASSANDRA"
    log_severity_level      = "INFO"
    should_log_monitor_files = true
  }
}
`, clusterExtID)
}

func testRsyslogServerResourceConfigUpdate(clusterExtID string) string {
	return fmt.Sprintf(`
resource "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id   = "%s"
  server_name      = "tf-test-rsyslog-server"
  port             = 1514
  network_protocol = "TCP"

  ip_address {
    ipv4 {
      value = "10.0.0.2"
    }
  }

  modules {
    name                    = "PRISM"
    log_severity_level      = "WARNING"
    should_log_monitor_files = false
  }
}
`, clusterExtID)
}

func testRsyslogServerDatasourceConfig() string {
	return `
data "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id = nutanix_rsyslog_server_v2.test.cluster_ext_id
  ext_id         = nutanix_rsyslog_server_v2.test.id
}
`
}

func testRsyslogServersDatasourceConfig(clusterExtID string) string {
	return fmt.Sprintf(`
data "nutanix_rsyslog_servers_v2" "test" {
  cluster_ext_id = "%s"
  depends_on     = [nutanix_rsyslog_server_v2.test]
}
`, clusterExtID)
}
