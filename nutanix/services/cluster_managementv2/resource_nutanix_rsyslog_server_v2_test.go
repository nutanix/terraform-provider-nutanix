package cluster_managementv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
)

const resourceNameRsyslogServer = "nutanix_rsyslog_server_v2.test"
const dataSourceNameRsyslogServer = "data.nutanix_rsyslog_server_v2.test"
const dataSourceNameRsyslogServers = "data.nutanix_rsyslog_servers_v2.test"

func TestAccV2NutanixRsyslogServerResource_Basic(t *testing.T) {
	r := acctest.RandInt()
	serverName := fmt.Sprintf("tf-test-rsyslog-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceConfig(serverName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "server_name", serverName),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "port", "514"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "network_protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "ip_address.0.ipv4.0.value", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.name", "CASSANDRA"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.log_severity_level", "INFO"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.should_log_monitor_files", "true"),
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "cluster_ext_id"),
				),
			},
		},
	})
}

func TestAccV2NutanixRsyslogServerResource_WithDatasource(t *testing.T) {
	r := acctest.RandInt()
	serverName := fmt.Sprintf("tf-test-rsyslog-ds-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceConfig(serverName) + testRsyslogServerDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(dataSourceNameRsyslogServer, "server_name", serverName),
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

func TestAccV2NutanixRsyslogServerResource_WithListDatasource(t *testing.T) {
	r := acctest.RandInt()
	serverName := fmt.Sprintf("tf-test-rsyslog-list-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceConfig(serverName) + testRsyslogServersDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceNameRsyslogServers, "rsyslog_servers.#"),
					checkAttributeLength(dataSourceNameRsyslogServers, "rsyslog_servers", 1),
				),
			},
		},
	})
}

func TestAccV2NutanixRsyslogServerResource_Update(t *testing.T) {
	r := acctest.RandInt()
	serverName := fmt.Sprintf("tf-test-rsyslog-upd-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceConfig(serverName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "server_name", serverName),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "port", "514"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "network_protocol", "UDP"),
				),
			},
			{
				Config: testRsyslogServerResourceUpdateConfig(serverName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "server_name", serverName),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "port", "1514"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "network_protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.name", "STARGATE"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.0.log_severity_level", "WARNING"),
				),
			},
		},
	})
}

func TestAccV2NutanixRsyslogServerResource_TCPProtocol(t *testing.T) {
	r := acctest.RandInt()
	serverName := fmt.Sprintf("tf-test-rsyslog-tcp-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceTCPConfig(serverName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "server_name", serverName),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "port", "6514"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "network_protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "ip_address.0.ipv4.0.value", "10.0.0.2"),
				),
			},
		},
	})
}

func TestAccV2NutanixRsyslogServerResource_MultipleModules(t *testing.T) {
	r := acctest.RandInt()
	serverName := fmt.Sprintf("tf-test-rsyslog-multi-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testRsyslogServerResourceMultiModulesConfig(serverName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceNameRsyslogServer, "ext_id"),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "server_name", serverName),
					resource.TestCheckResourceAttr(resourceNameRsyslogServer, "modules.#", "3"),
				),
			},
		},
	})
}

func testRsyslogServerResourceConfig(serverName string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
}

resource "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id = local.cluster_ext_id
  server_name    = "%s"
  port           = 514
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
}`, serverName)
}

func testRsyslogServerResourceUpdateConfig(serverName string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
}

resource "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id = local.cluster_ext_id
  server_name    = "%s"
  port           = 1514
  network_protocol = "TCP"

  ip_address {
    ipv4 {
      value = "10.0.0.1"
    }
  }

  modules {
    name                    = "STARGATE"
    log_severity_level      = "WARNING"
    should_log_monitor_files = false
  }
}`, serverName)
}

func testRsyslogServerResourceTCPConfig(serverName string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
}

resource "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id = local.cluster_ext_id
  server_name    = "%s"
  port           = 6514
  network_protocol = "TCP"

  ip_address {
    ipv4 {
      value = "10.0.0.2"
    }
  }
}`, serverName)
}

func testRsyslogServerResourceMultiModulesConfig(serverName string) string {
	return fmt.Sprintf(`
data "nutanix_clusters_v2" "clusters" {
  filter = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
}

locals {
  cluster_ext_id = data.nutanix_clusters_v2.clusters.cluster_entities.0.ext_id
}

resource "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id = local.cluster_ext_id
  server_name    = "%s"
  port           = 514
  network_protocol = "UDP"

  ip_address {
    ipv4 {
      value = "10.0.0.3"
    }
  }

  modules {
    name                    = "CASSANDRA"
    log_severity_level      = "INFO"
    should_log_monitor_files = true
  }

  modules {
    name                    = "STARGATE"
    log_severity_level      = "WARNING"
    should_log_monitor_files = false
  }

  modules {
    name                    = "GENESIS"
    log_severity_level      = "ERROR"
    should_log_monitor_files = true
  }
}`, serverName)
}

func testRsyslogServerDatasourceConfig() string {
	return `
data "nutanix_rsyslog_server_v2" "test" {
  cluster_ext_id = nutanix_rsyslog_server_v2.test.cluster_ext_id
  ext_id         = nutanix_rsyslog_server_v2.test.ext_id
  depends_on     = [nutanix_rsyslog_server_v2.test]
}
`
}

func testRsyslogServersDatasourceConfig() string {
	return `
data "nutanix_rsyslog_servers_v2" "test" {
  cluster_ext_id = nutanix_rsyslog_server_v2.test.cluster_ext_id
  depends_on     = [nutanix_rsyslog_server_v2.test]
}
`
}
