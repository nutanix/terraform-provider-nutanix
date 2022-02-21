package foundation

//Node imaging input
type ImageNodesInput struct {
	XsMasterLabel         *string       `json:"xs_master_label,omitempty"`
	LayoutEggUUID         *string       `json:"layout_egg_uuid,omitempty"`
	IpmiPassword          string        `json:"ipmi_password"`
	CvmGateway            string        `json:"cvm_gateway"`
	HypervExternalVnic    *string       `json:"hyperv_external_vnic,omitempty"`
	XenConfigType         *string       `json:"xen_config_type,omitempty"`
	UcsmIP                *string       `json:"ucsm_ip,omitempty"`
	UcsmPassword          *string       `json:"ucsm_password,omitempty"`
	HypervisorIso         HypervisorIso `json:"hypervisor_iso"`
	UncPath               *string       `json:"unc_path,omitempty"`
	HypervisorNetmask     string        `json:"hypervisor_netmask"`
	FcSettings            *FcSettings   `json:"fc_settings,omitempty"`
	XsMasterPassword      *string       `json:"xs_master_password,omitempty"`
	SvmRescueArgs         []*string     `json:"svm_rescue_args,omitempty"`
	CvmNetmask            string        `json:"cvm_netmask"`
	XsMasterIP            *string       `json:"xs_master_ip,omitempty"`
	Clusters              []*Clusters   `json:"clusters,omitempty"`
	HypervExternalVswitch *string       `json:"hyperv_external_vswitch,omitempty"`
	HypervisorNameserver  string        `json:"hypervisor_nameserver"`
	HypervSku             *string       `json:"hyperv_sku,omitempty"`
	EosMetadata           *EosMetadata  `json:"eos_metadata,omitempty"`
	Tests                 *Tests        `json:"tests,omitempty"`
	Blocks                []Block       `json:"blocks"`
	HypervProductKey      *string       `json:"hyperv_product_key,omitempty"`
	UncUsername           *string       `json:"unc_username,omitempty"`
	InstallScript         *string       `json:"install_script,omitempty"`
	IpmiUser              string        `json:"ipmi_user"`
	HypervisorPassword    *string       `json:"hypervisor_password,omitempty"`
	UncPassword           *string       `json:"unc_password,omitempty"`
	XsMasterUsername      *string       `json:"xs_master_username,omitempty"`
	SkipHypervisor        *bool         `json:"skip_hypervisor,omitempty"`
	HypervisorGateway     *string       `json:"hypervisor_gateway"`
	NosPackage            string        `json:"nos_package"` //To-do for cluster creation:null for cluster creation
	UcsmUser              *string       `json:"ucsm_user,omitempty"`
}

//Specific hypervisor defination for imaging
type Hypervisor struct {
	Checksum *string `json:"checksum,omitempty"`
	Filename string  `json:"filename"`
}

//Hypervisor ISO's for various kinds
type HypervisorIso struct {
	Hyperv *Hypervisor `json:"hyperv,omitempty"`
	Kvm    *Hypervisor `json:"kvm,omitempty"`
	Xen    *Hypervisor `json:"xen,omitempty"`
	Esx    *Hypervisor `json:"esx,omitempty"`
}

//Foundation Central Metadata
type FcMetadata struct {
	FcIP   string `json:"fc_ip"`
	APIKey string `json:"api_key"`
}

//Foundaton Central settings
type FcSettings struct {
	FcMetadata        FcMetadata `json:"fc_metadata"`
	FoundationCentral bool       `json:"foundation_central"`
}

//To-do for cluster creation: check for required fields
//Clusters creation related information
type Clusters struct {
	EnableNs              *bool    `json:"enable_ns,omitempty"`
	BackplaneSubnet       *string  `json:"backplane_subnet,omitempty"`
	ClusterInitSuccessful *bool    `json:"cluster_init_successful"`
	BackplaneNetmask      *string  `json:"backplane_netmask,omitempty"`
	RedundancyFactor      int      `json:"redundancy_factor"`
	BackplaneVlan         *string  `json:"backplane_vlan,omitempty"` //type unknown in docs
	ClusterName           string   `json:"cluster_name"`
	ClusterExternalIP     *string  `json:"cluster_external_ip,omitempty"`
	CvmNtpServers         *string  `json:"cvm_ntp_servers,omitempty"`
	SingleNodeCluster     *bool    `json:"single_node_cluster,omitempty"`
	ClusterMembers        []string `json:"cluster_members"`
	CvmDNSServers         *string  `json:"cvm_dns_servers,omitempty"`
	ClusterInitNow        bool     `json:"cluster_init_now"`
	HypervisorNtpServers  *string  `json:"hypervisor_ntp_servers,omitempty"`
}

type EosMetadata struct {
	ConfigID    string   `json:"config_id"`
	AccountName []string `json:"account_name"`
	Email       string   `json:"email"`
}

type Tests struct {
	RunSyscheck bool `json:"run_syscheck"`
	RunNcc      bool `json:"run_ncc"`
}

type UcsmParams struct {
	NativeVlan       bool   `json:"native_vlan"`
	KeepUcsmSettings bool   `json:"keep_ucsm_settings"`
	MacPool          string `json:"mac_pool"`
	VlanName         string `json:"vlan_name"`
}

type Vswitches struct {
	Lacp        string   `json:"lacp"`
	BondMode    string   `json:"bond_mode"`
	Name        string   `json:"name"`
	Uplinks     []string `json:"uplinks"`
	OtherConfig []string `json:"other_config"`
	Mtu         int      `json:"mtu"`
}

//Single node defination
type Node struct {
	Ipv6Address             *string      `json:"ipv6_address,omitempty"`
	NodePosition            *string      `json:"node_position,omitempty"`
	ImageDelay              int          `json:"image_delay,omitempty"`
	UcsmParams              UcsmParams   `json:"ucsm_params,omitempty"`
	HypervisorHostname      *string      `json:"hypervisor_hostname"`
	CvmGbRAM                *int         `json:"cvm_gb_ram,omitempty"`
	DeviceHint              *string      `json:"device_hint,omitempty"`
	BondMode                *string      `json:"bond_mode,omitempty"`
	RdmaPassthrough         *bool        `json:"rdma_passthrough,omitempty"`
	ClusterID               *string      `json:"cluster_id,omitempty"` //type unknown in docs
	UcsmNodeSerial          *string      `json:"ucsm_node_serial,omitempty"`
	HypervisorIP            *string      `json:"hypervisor_ip"`
	NodeSerial              *string      `json:"node_serial,omitempty"`
	IpmiConfigureNow        *bool        `json:"ipmi_configure_now,omitempty"`
	ImageSuccessful         *bool        `json:"image_successful,omitempty"`
	Ipv6Interface           *string      `json:"ipv6_interface,omitempty"`
	CvmNumVcpus             *int         `json:"cvm_num_vcpus,omitempty"`
	IpmiMac                 *string      `json:"ipmi_mac,omitempty"`
	RdmaMacAddr             *string      `json:"rdma_mac_addr,omitempty"`
	BondUplinks             []*string    `json:"bond_uplinks,omitempty"`
	CurrentNetworkInterface *string      `json:"current_network_interface,omitempty"`
	Hypervisor              *string      `json:"hypervisor"`
	Vswitches               []*Vswitches `json:"vswitches,omitempty"`
	BondLacpRate            *string      `json:"bond_lacp_rate,omitempty"`
	ImageNow                *bool        `json:"image_now"`
	UcsmManagedMode         *string      `json:"ucsm_managed_mode,omitempty"`
	IpmiIP                  *string      `json:"ipmi_ip"`
	CurrentCvmVlanTag       *string      `json:"current_cvm_vlan_tag,omitempty"` //type unknown in docs
	CvmIP                   *string      `json:"cvm_ip"`
	ExludeBootSerial        *string      `json:"exlude_boot_serial,omitempty"`
	MitigateLowBootSpace    *bool        `json:"mitigate_low_boot_space,omitempty"`
}

//Block containing multiple nodes
type Block struct {
	Nodes   []Node  `json:"nodes"`
	BlockID *string `json:"block_id,omitempty"`
}

//Response from /image_nodes API call
type ImageNodesAPIResponse struct {
	SessionID string `mapstructure:"session_id"`
}

//Node Imaging progress response
type ProgressResponse struct {
	AbortSession             bool              `mapstructure:"abort_session"`
	Results                  []string          `mapstructure:"results"`
	SessionID                string            `mapstructure:"session_id"`
	ImagingStopped           bool              `mapstructure:"imaging_stopped"`
	AggregatePercentComplete int               `mapstructure:"aggregate_percent_complete"`
	Action                   string            `mapstructure:"action"`
	Clusters                 []ClusterProgress `mapstructure:"clusters"`
	Nodes                    []NodeProgress    `mapstructure:"nodes"`
}
type ClusterProgress struct {
	Category        []string `mapstructure:"category,omitempty"`
	Status          string   `mapstructure:"status"`
	Messages        []string `mapstructure:"messages"`
	ClusterName     string   `mapstructure:"cluster_name"`
	TimeElapsed     int      `mapstructure:"time_elapsed"`
	ClusterMembers  []string `mapstructure:"cluster_members"`
	PercentComplete int      `mapstructure:"percent_complete"`
	TimeTotal       int      `mapstructure:"time_total"`
}
type NodeProgress struct {
	Category        []string `mapstructure:"category,omitempty"`
	Status          string   `mapstructure:"status"`
	Messages        []string `mapstructure:"messages"`
	TimeElapsed     int      `mapstructure:"time_elapsed"`
	CvmIP           string   `mapstructure:"cvm_ip"`
	PercentComplete int      `mapstructure:"percent_complete"`
	HypervisorIP    string   `mapstructure:"hypervisor_ip"`
	TimeTotal       int      `mapstructure:"time_total"`
}

//Response from /enumerate_nos_packages api
type ListNOSPackagesResponse []string

//Reference to hypervisor for ListHypervisorISOsResponse
type HypervisorISOReference struct {
	Supported bool   `mapstructure:"supported"`
	Filename  string `mapstructure:"filename"`
}

//Response from /enumerate_hypervisor_isos api
type ListHypervisorISOsResponse struct {
	Hyperv []HypervisorISOReference `mapstructure:"hyperv"`
	Kvm    []HypervisorISOReference `mapstructure:"kvm"`
	Esx    []HypervisorISOReference `mapstructure:"esx"`
	Linux  []HypervisorISOReference `mapstructure:"linux"`
}
