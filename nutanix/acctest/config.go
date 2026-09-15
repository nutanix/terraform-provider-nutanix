package acctest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

// This file provides a single, shared definition + loader for the acceptance
// test fixtures held in test_config_v2.json. It replaces the per-package
// TestConfig/loadVars/TestMain boilerplate that used to live in every
// nutanix/services/*v2 package.
//
// Usage from a test:
//
//	cfg := acctest.Config(t)
//	subnet := cfg.VMM.Subnet
//
// Connection settings (endpoint, port, credentials, insecure) continue to come
// from environment variables (see TestAccPreCheck); this file only models the
// test *data* fixtures.

// Reusable sub-types shared across domains.

// Credential is a generic username/password pair.
type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DomainUsersUsergroups maps test AD users / user-groups to their ext IDs.
type DomainUsersUsergroups struct {
	Users      map[string]string `json:"users"`
	UserGroups map[string]string `json:"user_groups"`
}

// DomainConfig models an Active Directory / domain definition. It is a superset
// of the fields used by the various consumers (primary vs secondary AD).
type DomainConfig struct {
	Name string `json:"name"`
	// DirectoryServiceName is the directory-service label (no dots; the IAM API
	// rejects dotted names). Name keeps the dotted domain (used by vmmv2 as the
	// Windows AD domain), so tests that create a directory service use this.
	DirectoryServiceName  string                `json:"directory_service_name"`
	Username              string                `json:"username"`
	Password              string                `json:"password"`
	DNS                   string                `json:"dns"`
	URL                   string                `json:"url"`
	ExtID                 string                `json:"ext_id"`
	WhiteListedGroups     []string              `json:"white_listed_groups"`
	DomainUsersUsergroups DomainUsersUsergroups `json:"domain_users_usergroups"`
}

// SubnetConfig models a reusable subnet fixture.
type SubnetConfig struct {
	Name         string `json:"name"`
	VlanID       int    `json:"vlan_id"`
	PrefixLength int    `json:"prefix_length"`
	GatewayIP    string `json:"gateway_ip"`
	StartIP      string `json:"start_ip"`
	EndIP        string `json:"end_ip"`
}

// TestConfig is the complete, typed view of test_config_v2.json. It is the
// union of every per-package config struct that previously existed.
type TestConfig struct {
	Images struct {
		UbuntuImage     string `json:"ubuntu_image"`
		UbuntuImageURL  string `json:"ubuntu_image_url"`
		WindowsImage    string `json:"windows_image"`
		WindowsImageURL string `json:"windows_image_url"`
		CentosImage     string `json:"centos_image"`
		CentosImageURL  string `json:"centos_image_url"`
		NgtImage        string `json:"ngt_image"`
		NgtImageURL     string `json:"ngt_image_url"`
		ISOImageURL     string `json:"iso_image_url"`
		ISOImageSHA1    string `json:"iso_image_sha1"`
		ISOImageSHA256  string `json:"iso_image_sha256"`
	} `json:"images"`
	UsernameForTest string `json:"username_for_test"`
	PasswordForTest string `json:"password_for_test"`
	SshPcUsername   string `json:"ssh_pc_username"`
	SshPcPassword   string `json:"ssh_pc_password"`
	SshPeUsername   string `json:"ssh_pe_username"`
	SshPePassword   string `json:"ssh_pe_password"`
	PeIP            string `json:"pe_ip"`

	// Consumed by passwordmanagerv2 (keys may be absent from the file).
	PCUsername string `json:"pc_username"`
	PCPassword string `json:"pc_password"`

	// Top-level Prism Element credentials (consumed by clustersv2 HCL/asserts).
	PEUsername string `json:"pe_username"`
	PEPassword string `json:"pe_password"`

	// Top-level DNS / NTP servers (consumed by clustersv2 + prismv2 HCL/asserts).
	DNSServers []string `json:"dns_servers"`
	NTPServers []string `json:"ntp_servers"`

	AvailabilityZone struct {
		PcExtID      string `json:"pc_ext_id"`
		ClusterExtID string `json:"cluster_ext_id"`
		RemotePcIP   string `json:"remote_pc_ip"`
	} `json:"availability_zone"`

	Iam        IAMConfig        `json:"iam"`
	Networking NetworkingConfig `json:"networking"`
	NicProfile NicProfileConfig `json:"nic_profile"`
	VMM        VMMConfig        `json:"vmm"`
	Clusters   ClustersConfig   `json:"clusters"`
	Prism      PrismConfig      `json:"prism"`

	DataProtection struct {
		LocalClusterPE   string `json:"local_cluster_pe"`
		LocalClusterVIP  string `json:"local_cluster_vip"`
		RemoteClusterPE  string `json:"remote_cluster_pe"`
		RemoteClusterVIP string `json:"remote_cluster_vip"`
	} `json:"data_protection"`

	Lcm struct {
		EntityModel        string `json:"entity_model"`
		EntityModelVersion string `json:"entity_model_version"`
	} `json:"lcm"`

	ObjectStore struct {
		SubnetName          string   `json:"subnet_name"`
		BucketName          string   `json:"bucket_name"`
		Domain              string   `json:"domain"`
		PublicNetworkIPs    []string `json:"public_network_ips"`
		StorageNetworkDNSIP []string `json:"storage_network_dns_ip"`
		StorageNetworkVip   []string `json:"storage_network_vip"`
	} `json:"object_store"`

	Security struct {
		KMS struct {
			EndpointURL  string `json:"endpoint_url"`
			TenantID     string `json:"tenant_id"`
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			KeyID        string `json:"key_id"`
		} `json:"kms"`
	} `json:"security"`

	Volumes struct {
		VolumeGroupExtIDWithCategory string `json:"vg_ext_id_with_category"`
	} `json:"volumes"`
}

// IAMConfig models the "iam" domain (consumed by iamv2, vmmv2, microsegv2).
type IAMConfig struct {
	DirectoryServicesMain struct {
		PrimaryAD   DomainConfig `json:"primary_ad"`
		SecondaryAD DomainConfig `json:"secondary_ad"`
	} `json:"directory_services_main"`
	Users struct {
		Name                        string `json:"name"`
		IdpID                       string `json:"idp_id"`
		DirectoryServiceID          string `json:"directory_service_id"`
		DirectoryServiceUsername    string `json:"directory_service_username"`
		EmailID                     string `json:"email_id"`
		Locale                      string `json:"locale"`
		Region                      string `json:"region"`
		Password                    string `json:"password"`
		IsForceResetPasswordEnabled bool   `json:"force_reset_password"`
		ExtID                       string `json:"ext_id"`
	} `json:"users"`
	UserGroups struct {
		Name              string `json:"name"`
		SAMLName          string `json:"saml_name"`
		DistinguishedName string `json:"distinguished_name"`
		ExtID             string `json:"ext_id"`
	} `json:"user_groups"`
	IdentityProviders struct {
		IdpMetadataURL string `json:"idp_metadata_url"`
		IdpMetadata    struct {
			EntityID           string `json:"entity_id"`
			LoginURL           string `json:"login_url"`
			LogoutURL          string `json:"logout_url"`
			ErrorURL           string `json:"error_url"`
			Certificate        string `json:"certificate"`
			NameIDPolicyFormat string `json:"name_id_policy_format"`
		} `json:"idp_metadata"`
		IdpMetadataXML          string   `json:"idp_metadata_xml"`
		Name                    string   `json:"name"`
		UsernameAttribute       string   `json:"username_attr"`
		EmailAttribute          string   `json:"email_attr"`
		GroupsAttribute         string   `json:"groups_attr"`
		GroupsDelim             string   `json:"groups_delim"`
		EntityIssuer            string   `json:"entity_issuer"`
		CustomAttributes        []string `json:"custom_attributes"`
		IsSignedAuthnReqEnabled bool     `json:"is_signed_authn_req_enabled"`
	} `json:"identity_providers"`
}

// NetworkingConfig models the "networking" domain.
type NetworkingConfig struct {
	FloatingIP struct {
		VMNicReference string `json:"vm_nic_reference"`
	} `json:"floating_ip"`
	Subnets struct {
		ProjectID     string `json:"project_id"`
		VlanID        int    `json:"vlan_id"`
		NetworkIP     string `json:"network_ip"`
		NetworkPrefix int    `json:"network_prefix"`
		GatewayIP     string `json:"gateway_ip"`
		DHCP          struct {
			StartIP string `json:"start_ip"`
			EndIP   string `json:"end_ip"`
		} `json:"dhcp"`
	} `json:"subnets"`
	BridgeName string       `json:"bridge_name"`
	GcSubnet   SubnetConfig `json:"gc_subnet"`
	// KubernetesClusterExtID is the extId of an NKP/Karbon Kubernetes cluster
	// already registered with Prism Central
	KubernetesClusterExtID string `json:"kubernetes_cluster_ext_id"`
}

// NicProfileConfig models the (currently top-level) "nic_profile" fixture.
type NicProfileConfig struct {
	ExtID          string   `json:"ext_id"`
	Name           string   `json:"name"`
	CapabilityType string   `json:"capability_type"`
	NicFamily      string   `json:"nic_family"`
	HostNicExtIDs  []string `json:"host_nic_ext_ids"`
}

// VMMConfig models the "vmm" domain.
type VMMConfig struct {
	IntegrationVM string `json:"integration_vm"`
	AssignedIP    string `json:"assigned_ip"`
	UnattendXML   string `json:"unattend_xml"`
	Subnet        struct {
		NetworkID int    `json:"network_id"`
		IP        string `json:"ip"`
		Prefix    int    `json:"prefix"`
		GatewayIP string `json:"gateway_ip"`
		StartIP   string `json:"start_ip"`
		EndIP     string `json:"end_ip"`
	} `json:"subnet"`
	NGT struct {
		NgtUpgradeVMName string     `json:"ngt_upgrade_vm_name"`
		Credential       Credential `json:"credential"`
	} `json:"ngt"`
	NgtVM struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"ngt_vm"`
	GPUS []struct {
		DeviceID int    `json:"device_id"`
		Mode     string `json:"mode"`
		Vendor   string `json:"vendor"`
	} `json:"gpus"`
	OvaURL    string `json:"ova_url"`
	GcProfile struct {
		VMName               string `json:"vm_name"`
		DefaultImageUsername string `json:"default_image_username"`
		DefaultImagePassword string `json:"default_image_password"`
	} `json:"gc_profile"`
}

// ClustersConfig models the "clusters" domain.
type ClustersConfig struct {
	Nodes   []string `json:"nodes"`
	Network struct {
		VirtualIP string `json:"virtual_ip"`
		IscsiIP   string `json:"iscsi_ip"`
	} `json:"network"`
	SSLCertificate struct {
		Passphrase        string `json:"passphrase"`
		PrivateKey        string `json:"private_key"`
		PublicCertificate string `json:"public_certificate"`
		CaChain           string `json:"ca_chain"`
	} `json:"ssl_certificate"`
}

// PrismConfig models the "prism" domain.
type PrismConfig struct {
	DeployPC struct {
		PeIP           string `json:"pe_ip"`
		Version        string `json:"version"`
		SubnetName     string `json:"subnet_name"`
		DefaultGateway string `json:"default_gateway"`
		SubnetMask     string `json:"subnet_mask"`
		IPRange        struct {
			Begin string `json:"begin"`
			End   string `json:"end"`
		} `json:"ip_range"`
	} `json:"deploy_pc"`
	Bucket struct {
		Name      string `json:"name"`
		Region    string `json:"region"`
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
	} `json:"bucket"`
	RestoreSource struct {
		PeIP string `json:"pe_ip"`
	} `json:"restore_source"`
	Unregister struct {
		PcExtID string `json:"pc_ext_id"`
	} `json:"unregister"`
	PCRestore struct {
		SkipPCRestoreTest bool `json:"skip_pc_restore_test"`
	} `json:"pc_restore"`
}

var (
	loadedConfig TestConfig
	loadOnce     sync.Once
	loadErr      error
)

// ConfigPath resolves the test config file path. It honors the TEST_CONFIG
// environment variable and otherwise defaults to test_config_v2.json at the
// repository root (resolved relative to this source file, so it is independent
// of the per-package working directory).
func ConfigPath() string {
	if p := os.Getenv("TEST_CONFIG"); p != "" {
		return p
	}
	_, thisFile, _, _ := runtime.Caller(0)
	// this file lives at <repo>/nutanix/acctest/config.go
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "test_config_v2.json")
}

func loadConfig() {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		loadErr = fmt.Errorf("reading test config %q: %w", path, err)
		return
	}
	if err := json.Unmarshal(data, &loadedConfig); err != nil {
		loadErr = fmt.Errorf("parsing test config %q: %w", path, err)
		return
	}
}

// Config returns the parsed test fixtures, loaded once and cached. It fails the
// test with a clear message if the file cannot be read or parsed.
func Config(t *testing.T) TestConfig {
	t.Helper()
	loadOnce.Do(loadConfig)
	if loadErr != nil {
		t.Fatalf("acctest: %v", loadErr)
	}
	return loadedConfig
}

// MustConfig is like Config but does not require a *testing.T. It is intended
// for package-level initialization (e.g. `var testVars = acctest.MustConfig()`)
// and panics if the config cannot be loaded — mirroring the previous
// loadVars/os.Exit(1) behavior in each package's TestMain.
func MustConfig() TestConfig {
	loadOnce.Do(loadConfig)
	if loadErr != nil {
		panic(fmt.Sprintf("acctest: %v", loadErr))
	}
	return loadedConfig
}
