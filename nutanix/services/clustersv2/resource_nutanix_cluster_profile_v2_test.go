package clustersv2_test

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	conns "github.com/terraform-providers/terraform-provider-nutanix/nutanix"
	acc "github.com/terraform-providers/terraform-provider-nutanix/nutanix/acctest"
	"github.com/terraform-providers/terraform-provider-nutanix/nutanix/common"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func TestAccNutanixClusterProfileV2_basic(t *testing.T) {
	resourceName1 := "nutanix_cluster_profile_v2.tf_first"
	resourceName2 := "nutanix_cluster_profile_v2.tf_second"

	dataSourceNameList := "data.nutanix_cluster_profiles_v2.all_profiles"
	dataSourceNameFilter := "data.nutanix_cluster_profiles_v2.filtered_profiles"
	dataSourceNameLimit := "data.nutanix_cluster_profiles_v2.limited_profiles"
	dataSourceNameSingle := "data.nutanix_cluster_profile_v2.first_profile"

	profileName1 := fmt.Sprintf("tf-test-cluster-profile1-%d", acc.RandIntBetween(1, 5000))
	profileName2 := fmt.Sprintf("tf-test-cluster-profile2-%d", acc.RandIntBetween(5001, 10000))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckClusterProfileDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create both cluster profiles
			{
				Config: testAccClusterProfilesConfig(profileName1, profileName2),
				Check: resource.ComposeTestCheckFunc(
					// First profile checks
					resource.TestCheckResourceAttrSet(resourceName1, "id"),
					resource.TestCheckResourceAttrSet(resourceName1, "ext_id"),
					resource.TestCheckResourceAttr(resourceName1, "name", profileName1),
					resource.TestCheckResourceAttr(resourceName1, "description", "Example First Cluster Profile created via Terraform"),
					resource.TestCheckResourceAttr(resourceName1, "allowed_overrides.#", "2"),
					resource.TestCheckResourceAttr(resourceName1, "allowed_overrides.0", "NTP_SERVER_CONFIG"),
					resource.TestCheckResourceAttr(resourceName1, "allowed_overrides.1", "SNMP_SERVER_CONFIG"),
					resource.TestCheckResourceAttr(resourceName1, "name_server_ip_list.0.ipv4.0.value", "240.29.254.180"),
					resource.TestCheckResourceAttr(resourceName1, "name_server_ip_list.0.ipv6.0.value", "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"),
					resource.TestCheckResourceAttr(resourceName1, "ntp_server_ip_list.0.ipv4.0.value", "240.29.254.180"),
					resource.TestCheckResourceAttr(resourceName1, "ntp_server_ip_list.0.ipv6.0.value", "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"),
					resource.TestCheckResourceAttr(resourceName1, "ntp_server_ip_list.0.fqdn.0.value", "ntp.example.com"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.email_address", "email@example.com"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.type", "SSL"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.ip_address.0.ipv4.0.value", "240.29.254.180"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.ip_address.0.ipv6.0.value", "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.ip_address.0.fqdn.0.value", "smtp.example.com"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.port", "587"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.username", "example_user"),
					resource.TestCheckResourceAttr(resourceName1, "nfs_subnet_white_list.0", "10.110.106.45/255.255.255.255"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.users.0.username", "snmpuser1"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.users.0.auth_type", "MD5"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.users.0.priv_type", "DES"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.transports.0.protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.transports.0.port", "21"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.address.0.ipv4.0.value", "240.29.254.180"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.address.0.ipv4.0.prefix_length", "24"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.address.0.ipv6.0.value", "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.username", "trapuser"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.port", "59"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.engine_id", "0x1234567890abcdef12"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.version", "V2"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.receiver_name", "trap-receiver"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.traps.0.community_string", "snmp-server community public RO 192.168.1.0 255.255.255.0"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.server_name", "testServer1"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.port", "29"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.network_protocol", "UDP"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.ip_address.0.ipv4.0.value", "240.29.254.180"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.ip_address.0.ipv6.0.value", "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.modules.0.name", "CASSANDRA"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.modules.0.log_severity_level", "EMERGENCY"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.modules.0.should_log_monitor_files", "true"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.modules.1.name", "CURATOR"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.modules.1.log_severity_level", "ERROR"),
					resource.TestCheckResourceAttr(resourceName1, "rsyslog_server_list.0.modules.1.should_log_monitor_files", "false"),
					resource.TestCheckResourceAttr(resourceName1, "pulse_status.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName1, "pulse_status.0.pii_scrubbing_level", "DEFAULT"),

					// Second profile checks
					resource.TestCheckResourceAttrSet(resourceName2, "id"),
					resource.TestCheckResourceAttrSet(resourceName2, "ext_id"),
					resource.TestCheckResourceAttr(resourceName2, "name", profileName2),
					resource.TestCheckResourceAttr(resourceName2, "description", "Example Second Cluster Profile created via Terraform"),
					resource.TestCheckResourceAttr(resourceName2, "allowed_overrides.#", "2"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.email_address", "email2@example.com"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.type", "STARTTLS"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.ip_address.0.ipv4.0.value", "240.29.254.190"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.ip_address.0.ipv6.0.value", "1c89:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.ip_address.0.fqdn.0.value", "smtp2.example.com"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.port", "468"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.username", "smtp2-user"),
					resource.TestCheckResourceAttr(resourceName2, "snmp_config.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "snmp_config.0.users.0.username", "snmpuser2"),
					resource.TestCheckResourceAttr(resourceName2, "snmp_config.0.users.0.auth_type", "SHA"),
					resource.TestCheckResourceAttr(resourceName2, "pulse_status.0.is_enabled", "true"),

					// Data source checks
					// List checks
					common.CheckAttributeLength(dataSourceNameList, "cluster_profiles", 2),
					resource.TestCheckResourceAttrSet(dataSourceNameList, "cluster_profiles.0.ext_id"),
					resource.TestCheckResourceAttrSet(dataSourceNameList, "cluster_profiles.1.ext_id"),
					// Filter checks
					common.CheckAttributeLengthEqual(dataSourceNameFilter, "cluster_profiles", 1),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.ext_id", resourceName1, "ext_id"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.ext_id", resourceName1, "id"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.name", resourceName1, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.description", resourceName1, "description"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.allowed_overrides.#", resourceName1, "allowed_overrides.#"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.allowed_overrides.0", resourceName1, "allowed_overrides.0"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.allowed_overrides.1", resourceName1, "allowed_overrides.1"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.name_server_ip_list.0.ipv4.0.value", resourceName1, "name_server_ip_list.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.name_server_ip_list.0.ipv6.0.value", resourceName1, "name_server_ip_list.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.ntp_server_ip_list.0.ipv4.0.value", resourceName1, "ntp_server_ip_list.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.ntp_server_ip_list.0.ipv6.0.value", resourceName1, "ntp_server_ip_list.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.ntp_server_ip_list.0.fqdn.0.value", resourceName1, "ntp_server_ip_list.0.fqdn.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.smtp_server.0.email_address", resourceName1, "smtp_server.0.email_address"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.smtp_server.0.type", resourceName1, "smtp_server.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.smtp_server.0.server.0.ip_address.0.ipv4.0.value", resourceName1, "smtp_server.0.server.0.ip_address.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.smtp_server.0.server.0.ip_address.0.ipv6.0.value", resourceName1, "smtp_server.0.server.0.ip_address.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.smtp_server.0.server.0.ip_address.0.fqdn.0.value", resourceName1, "smtp_server.0.server.0.ip_address.0.fqdn.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.smtp_server.0.server.0.port", resourceName1, "smtp_server.0.server.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.smtp_server.0.server.0.username", resourceName1, "smtp_server.0.server.0.username"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.nfs_subnet_white_list.0", resourceName1, "nfs_subnet_white_list.0"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.is_enabled", resourceName1, "snmp_config.0.is_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.users.0.username", resourceName1, "snmp_config.0.users.0.username"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.users.0.auth_type", resourceName1, "snmp_config.0.users.0.auth_type"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.users.0.priv_type", resourceName1, "snmp_config.0.users.0.priv_type"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.transports.0.protocol", resourceName1, "snmp_config.0.transports.0.protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.transports.0.port", resourceName1, "snmp_config.0.transports.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.address.0.ipv4.0.value", resourceName1, "snmp_config.0.traps.0.address.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.address.0.ipv4.0.prefix_length", resourceName1, "snmp_config.0.traps.0.address.0.ipv4.0.prefix_length"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.address.0.ipv6.0.value", resourceName1, "snmp_config.0.traps.0.address.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.username", resourceName1, "snmp_config.0.traps.0.username"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.protocol", resourceName1, "snmp_config.0.traps.0.protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.port", resourceName1, "snmp_config.0.traps.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.engine_id", resourceName1, "snmp_config.0.traps.0.engine_id"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.version", resourceName1, "snmp_config.0.traps.0.version"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.receiver_name", resourceName1, "snmp_config.0.traps.0.receiver_name"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.snmp_config.0.traps.0.community_string", resourceName1, "snmp_config.0.traps.0.community_string"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.server_name", resourceName1, "rsyslog_server_list.0.server_name"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.port", resourceName1, "rsyslog_server_list.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.network_protocol", resourceName1, "rsyslog_server_list.0.network_protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.ip_address.0.ipv4.0.value", resourceName1, "rsyslog_server_list.0.ip_address.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.ip_address.0.ipv6.0.value", resourceName1, "rsyslog_server_list.0.ip_address.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.modules.0.name", resourceName1, "rsyslog_server_list.0.modules.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.modules.0.log_severity_level", resourceName1, "rsyslog_server_list.0.modules.0.log_severity_level"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.modules.0.should_log_monitor_files", resourceName1, "rsyslog_server_list.0.modules.0.should_log_monitor_files"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.modules.1.name", resourceName1, "rsyslog_server_list.0.modules.1.name"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.modules.1.log_severity_level", resourceName1, "rsyslog_server_list.0.modules.1.log_severity_level"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.rsyslog_server_list.0.modules.1.should_log_monitor_files", resourceName1, "rsyslog_server_list.0.modules.1.should_log_monitor_files"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.pulse_status.0.is_enabled", resourceName1, "pulse_status.0.is_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceNameFilter, "cluster_profiles.0.pulse_status.0.pii_scrubbing_level", resourceName1, "pulse_status.0.pii_scrubbing_level"),
					// Limit checks
					common.CheckAttributeLengthEqual(dataSourceNameLimit, "cluster_profiles", 1),
					// Single profile checks
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "ext_id", resourceName1, "ext_id"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "name", resourceName1, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "description", resourceName1, "description"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "allowed_overrides.#", resourceName1, "allowed_overrides.#"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "allowed_overrides.0", resourceName1, "allowed_overrides.0"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "allowed_overrides.1", resourceName1, "allowed_overrides.1"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "name_server_ip_list.0.ipv4.0.value", resourceName1, "name_server_ip_list.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "name_server_ip_list.0.ipv6.0.value", resourceName1, "name_server_ip_list.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "ntp_server_ip_list.0.ipv4.0.value", resourceName1, "ntp_server_ip_list.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "ntp_server_ip_list.0.ipv6.0.value", resourceName1, "ntp_server_ip_list.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "ntp_server_ip_list.0.fqdn.0.value", resourceName1, "ntp_server_ip_list.0.fqdn.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "smtp_server.0.email_address", resourceName1, "smtp_server.0.email_address"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "smtp_server.0.type", resourceName1, "smtp_server.0.type"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "smtp_server.0.server.0.ip_address.0.ipv4.0.value", resourceName1, "smtp_server.0.server.0.ip_address.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "smtp_server.0.server.0.ip_address.0.ipv6.0.value", resourceName1, "smtp_server.0.server.0.ip_address.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "smtp_server.0.server.0.ip_address.0.fqdn.0.value", resourceName1, "smtp_server.0.server.0.ip_address.0.fqdn.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "smtp_server.0.server.0.port", resourceName1, "smtp_server.0.server.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "smtp_server.0.server.0.username", resourceName1, "smtp_server.0.server.0.username"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "nfs_subnet_white_list.0", resourceName1, "nfs_subnet_white_list.0"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.is_enabled", resourceName1, "snmp_config.0.is_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.users.0.username", resourceName1, "snmp_config.0.users.0.username"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.users.0.auth_type", resourceName1, "snmp_config.0.users.0.auth_type"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.users.0.priv_type", resourceName1, "snmp_config.0.users.0.priv_type"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.transports.0.protocol", resourceName1, "snmp_config.0.transports.0.protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.transports.0.port", resourceName1, "snmp_config.0.transports.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.address.0.ipv4.0.value", resourceName1, "snmp_config.0.traps.0.address.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.address.0.ipv4.0.prefix_length", resourceName1, "snmp_config.0.traps.0.address.0.ipv4.0.prefix_length"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.address.0.ipv6.0.value", resourceName1, "snmp_config.0.traps.0.address.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.username", resourceName1, "snmp_config.0.traps.0.username"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.protocol", resourceName1, "snmp_config.0.traps.0.protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.port", resourceName1, "snmp_config.0.traps.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.engine_id", resourceName1, "snmp_config.0.traps.0.engine_id"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.version", resourceName1, "snmp_config.0.traps.0.version"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.receiver_name", resourceName1, "snmp_config.0.traps.0.receiver_name"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "snmp_config.0.traps.0.community_string", resourceName1, "snmp_config.0.traps.0.community_string"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.server_name", resourceName1, "rsyslog_server_list.0.server_name"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.port", resourceName1, "rsyslog_server_list.0.port"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.network_protocol", resourceName1, "rsyslog_server_list.0.network_protocol"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.ip_address.0.ipv4.0.value", resourceName1, "rsyslog_server_list.0.ip_address.0.ipv4.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.ip_address.0.ipv6.0.value", resourceName1, "rsyslog_server_list.0.ip_address.0.ipv6.0.value"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.modules.0.name", resourceName1, "rsyslog_server_list.0.modules.0.name"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.modules.0.log_severity_level", resourceName1, "rsyslog_server_list.0.modules.0.log_severity_level"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.modules.0.should_log_monitor_files", resourceName1, "rsyslog_server_list.0.modules.0.should_log_monitor_files"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.modules.1.name", resourceName1, "rsyslog_server_list.0.modules.1.name"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.modules.1.log_severity_level", resourceName1, "rsyslog_server_list.0.modules.1.log_severity_level"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "rsyslog_server_list.0.modules.1.should_log_monitor_files", resourceName1, "rsyslog_server_list.0.modules.1.should_log_monitor_files"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "pulse_status.0.is_enabled", resourceName1, "pulse_status.0.is_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceNameSingle, "pulse_status.0.pii_scrubbing_level", resourceName1, "pulse_status.0.pii_scrubbing_level"),
				),
			},
			// Step 2: Update first profile's description and pulse_status
			{
				Config: testAccClusterProfileFullUpdateConfig(profileName1+"_updated", profileName2),
				Check: resource.ComposeTestCheckFunc(
					// First profile full updated checks
					resource.TestCheckResourceAttrSet(resourceName1, "id"),
					resource.TestCheckResourceAttrSet(resourceName1, "ext_id"),
					resource.TestCheckResourceAttr(resourceName1, "name", profileName1+"_updated"),
					resource.TestCheckResourceAttr(resourceName1, "description", "Fully Updated First Cluster Profile"),
					resource.TestCheckResourceAttr(resourceName1, "allowed_overrides.#", "1"),
					resource.TestCheckResourceAttr(resourceName1, "allowed_overrides.0", "PULSE_CONFIG"),
					resource.TestCheckResourceAttr(resourceName1, "name_server_ip_list.0.ipv4.0.value", "10.1.1.1"),
					resource.TestCheckResourceAttr(resourceName1, "name_server_ip_list.0.ipv6.0.value", "fd00::1"),
					resource.TestCheckResourceAttr(resourceName1, "nfs_subnet_white_list.0", "192.168.1.0/255.255.255.255"),
					resource.TestCheckResourceAttr(resourceName1, "ntp_server_ip_list.0.ipv4.0.value", "10.1.1.2"),
					resource.TestCheckResourceAttr(resourceName1, "ntp_server_ip_list.0.ipv6.0.value", "fd00::2"),
					resource.TestCheckResourceAttr(resourceName1, "ntp_server_ip_list.0.fqdn.0.value", "ntp-updated.example.com"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.email_address", "updated@example.com"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.type", "STARTTLS"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.ip_address.0.ipv4.0.value", "10.1.1.3"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.ip_address.0.ipv6.0.value", "fd00::3"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.ip_address.0.fqdn.0.value", "smtp-updated.example.com"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.port", "2525"),
					resource.TestCheckResourceAttr(resourceName1, "smtp_server.0.server.0.username", "updated_user"),
					resource.TestCheckResourceAttr(resourceName1, "snmp_config.0.is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName1, "pulse_status.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName1, "pulse_status.0.pii_scrubbing_level", "ALL"),

					// Second profile remains unchanged
					resource.TestCheckResourceAttrSet(resourceName2, "id"),
					resource.TestCheckResourceAttrSet(resourceName2, "ext_id"),
					resource.TestCheckResourceAttr(resourceName2, "name", profileName2),
					resource.TestCheckResourceAttr(resourceName2, "description", "Example Second Cluster Profile created via Terraform"),
					resource.TestCheckResourceAttr(resourceName2, "allowed_overrides.#", "2"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.email_address", "email2@example.com"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.type", "STARTTLS"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.ip_address.0.ipv4.0.value", "240.29.254.190"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.ip_address.0.ipv6.0.value", "1c89:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.ip_address.0.fqdn.0.value", "smtp2.example.com"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.port", "468"),
					resource.TestCheckResourceAttr(resourceName2, "smtp_server.0.server.0.username", "smtp2-user"),
					resource.TestCheckResourceAttr(resourceName2, "snmp_config.0.is_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName2, "snmp_config.0.users.0.username", "snmpuser2"),
					resource.TestCheckResourceAttr(resourceName2, "snmp_config.0.users.0.auth_type", "SHA"),
					resource.TestCheckResourceAttr(resourceName2, "pulse_status.0.is_enabled", "true"),
				),
			},
		},
	})
}

func TestAccNutanixClusterProfileV2_duplicate(t *testing.T) {
	profileName1 := "tf-test-cluster-profile"
	profileName2 := "tf-test-cluster-profile"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckClusterProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterProfilesConfig(profileName1, profileName2),
				ExpectError: regexp.MustCompile(`profile name already exists on another profile`),
			},
		},
	})
}

func TestAccNutanixClusterProfileV2_fetchCPWrongExtID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckClusterProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
					data "nutanix_cluster_profile_v2" "test" {
						ext_id = "00000000-0000-0000-0000-000000000000"
					}
				`,
				ExpectError: regexp.MustCompile(`Error - profile 00000000-0000-0000-0000-000000000000 not found`),
			},
		},
	})
}

func testAccClusterProfilesConfig(profile1, profile2 string) string {
	return fmt.Sprintf(`
resource "nutanix_cluster_profile_v2" "tf_first" {
  name = "%s"
  description = "Example First Cluster Profile created via Terraform"
  allowed_overrides = ["NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"]

  name_server_ip_list {
    ipv4 { value = "240.29.254.180" }
    ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
  }

  ntp_server_ip_list {
    ipv4 { value = "240.29.254.180" }
    ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
    fqdn { value = "ntp.example.com" }
  }

  smtp_server {
    email_address = "email@example.com"
    type = "SSL"
    server {
      ip_address {
        ipv4 { value = "240.29.254.180" }
        ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
        fqdn { value = "smtp.example.com" }
      }
      port     = 587
      username = "example_user"
      password = "example_password"
    }
  }

  nfs_subnet_white_list = ["10.110.106.45/255.255.255.255"]

  snmp_config {
    is_enabled = true
    users {
      username  = "snmpuser1"
      auth_type = "MD5"
      auth_key  = "Test_SNMP_user_authentication_key"
      priv_type = "DES"
      priv_key  = "Test_SNMP_user_encryption_key"
    }
    transports {
      protocol = "UDP"
      port     = 21
    }
    traps {
      address {
        ipv4 {
					value         = "240.29.254.180"
					prefix_length = 24
				}
        ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
      }
      username         = "trapuser"
      protocol         = "UDP"
      port             = 59
      should_inform    = false
      engine_id        = "0x1234567890abcdef12"
      version          = "V2"
      receiver_name    = "trap-receiver"
      community_string = "snmp-server community public RO 192.168.1.0 255.255.255.0"
    }
  }

  rsyslog_server_list {
    server_name      = "testServer1"
    port             = 29
    network_protocol = "UDP"
    ip_address {
      ipv4 { value = "240.29.254.180" }
      ipv6 { value = "1a7d:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
    }
    modules {
      name                     = "CASSANDRA"
      log_severity_level       = "EMERGENCY"
      should_log_monitor_files = true
    }
    modules {
      name                     = "CURATOR"
      log_severity_level       = "ERROR"
      should_log_monitor_files = false
    }
  }

  pulse_status {
    is_enabled          = false
    pii_scrubbing_level = "DEFAULT"
  }

  lifecycle {
    ignore_changes = [
      smtp_server.0.server.0.password,
      snmp_config.0.users.0.auth_key,
      snmp_config.0.users.0.priv_key
    ]
  }
}

resource "nutanix_cluster_profile_v2" "tf_second" {
  name = "%s"
  description = "Example Second Cluster Profile created via Terraform"
  allowed_overrides = ["NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"]

  smtp_server {
    email_address = "email2@example.com"
    type = "STARTTLS"
    server {
      ip_address {
        ipv4 { value = "240.29.254.190" }
        ipv6 { value = "1c89:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
        fqdn { value = "smtp2.example.com" }
      }
      port     = 468
      username = "smtp2-user"
      password = "smtp2-password"
    }
  }

  snmp_config {
    is_enabled = true
    users {
      username  = "snmpuser2"
      auth_type = "SHA"
      auth_key  = "Test_SNMP_user_authentication_key2"
    }
  }

  pulse_status {
    is_enabled = true
  }

  lifecycle {
    ignore_changes = [
      smtp_server.0.server.0.password,
      snmp_config.0.users.0.auth_key
    ]
  }
}

# list all cluster profiles
data "nutanix_cluster_profiles_v2" "all_profiles" {
	depends_on = [
		nutanix_cluster_profile_v2.tf_first,
		nutanix_cluster_profile_v2.tf_second,
	]
}

# List cluster profile with filter
data "nutanix_cluster_profiles_v2" "filtered_profiles" {
	filter = "name eq '${nutanix_cluster_profile_v2.tf_first.name}'"
}

# List cluster profile with limit
data "nutanix_cluster_profiles_v2" "limited_profiles" {
	limit = 1
	depends_on = [
		nutanix_cluster_profile_v2.tf_first,
		nutanix_cluster_profile_v2.tf_second,
	]
}

# Get single cluster profile by ext_id
data "nutanix_cluster_profile_v2" "first_profile" {
	ext_id = nutanix_cluster_profile_v2.tf_first.id
}

`, profile1, profile2)
}

// ----------------------
// Full update config for first profile
// ----------------------
func testAccClusterProfileFullUpdateConfig(profile1, profile2 string) string {
	return fmt.Sprintf(`
resource "nutanix_cluster_profile_v2" "tf_first" {
  name = "%s"
  description = "Fully Updated First Cluster Profile"
  allowed_overrides = ["PULSE_CONFIG"]

  name_server_ip_list {
    ipv4 { value = "10.1.1.1" }
    ipv6 { value = "fd00::1" }
  }

  ntp_server_ip_list {
    ipv4 { value = "10.1.1.2" }
    ipv6 { value = "fd00::2" }
    fqdn { value = "ntp-updated.example.com" }
  }

  smtp_server {
		email_address = "updated@example.com"
		type = "STARTTLS"

		server {
			ip_address {
				ipv4 { value = "10.1.1.3" }
				ipv6 { value = "fd00::3" }
				fqdn { value = "smtp-updated.example.com" }
			}
			port     = 2525
			username = "updated_user"
			password = "updated_password"
		}
	}


  nfs_subnet_white_list = ["192.168.1.0/255.255.255.255"]

  snmp_config {
    is_enabled = false
  }

  pulse_status {
    is_enabled = true
    pii_scrubbing_level = "ALL"
  }

  lifecycle {
    ignore_changes = [smtp_server.0.server.0.password]
  }
}

resource "nutanix_cluster_profile_v2" "tf_second" {
  name = "%s"
  description = "Example Second Cluster Profile created via Terraform"
  allowed_overrides = ["NTP_SERVER_CONFIG", "SNMP_SERVER_CONFIG"]

  smtp_server {
    email_address = "email2@example.com"
    type = "STARTTLS"
    server {
      ip_address {
        ipv4 { value = "240.29.254.190" }
        ipv6 { value = "1c89:9a64:df8d:dfd8:39c6:c4ea:e35c:0ba4" }
        fqdn { value = "smtp2.example.com" }
      }
      port     = 468
      username = "smtp2-user"
      password = "smtp2-password"
    }
  }

  snmp_config {
    is_enabled = true
    users {
      username  = "snmpuser2"
      auth_type = "SHA"
      auth_key  = "Test_SNMP_user_authentication_key2"
    }
  }

  pulse_status {
    is_enabled = true
  }

  lifecycle {
    ignore_changes = [
      smtp_server.0.server.0.password,
      snmp_config.0.users.0.auth_key
    ]
  }
}
`, profile1, profile2)
}

// Check destroy function
func testAccCheckClusterProfileDestroy(s *terraform.State) error {
	conn := acc.TestAccProvider.Meta().(*conns.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "nutanix_cluster_profile_v2" {
			continue
		}
		// Check API if resource exists
		_, errRead := conn.ClusterAPI.ClusterProfilesAPI.GetClusterProfileById(utils.StringPtr(rs.Primary.ID))
		if errRead != nil {
			if strings.Contains(fmt.Sprint(errRead), "profile "+rs.Primary.ID+" not found") {
				return nil
			}
			return errRead
		}
		log.Printf("[DEBUG] Cluster Profile %s still exists, destroying...", rs.Primary.ID)
		_, err := conn.ClusterAPI.ClusterProfilesAPI.DeleteClusterProfileById(utils.StringPtr(rs.Primary.ID))
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Cluster Profile %s destroyed successfully", rs.Primary.ID)
	}
	return nil
}
