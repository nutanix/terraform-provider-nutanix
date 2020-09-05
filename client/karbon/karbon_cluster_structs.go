package karbon

// DSMetadata All api calls that return a list will have this metadata block as input
type DSMetadata struct {

	// The filter in FIQL syntax used for the results.
	Filter *string `json:"filter,omitempty" mapstructure:"filter,omitempty"`

	// The kind name
	Kind *string `json:"kind,omitempty" mapstructure:"kind,omitempty"`

	// The number of records to retrieve relative to the offset
	Length *int64 `json:"length,omitempty" mapstructure:"length,omitempty"`

	// Offset from the start of the entity list
	Offset *int64 `json:"offset,omitempty" mapstructure:"offset,omitempty"`

	// The attribute to perform sort on
	SortAttribute *string `json:"sort_attribute,omitempty" mapstructure:"sort_attribute,omitempty"`

	// The sort order in which results are returned
	SortOrder *string `json:"sort_order,omitempty" mapstructure:"sort_order,omitempty"`
}

//KARBON 2.1

type KarbonClusterListIntentResponse []KarbonClusterIntentResponse
type KarbonClusterIntentResponse struct {
	Name                     *string `json:"name" mapstructure:"name, omitempty"`
	UUID                     *string `json:"uuid" mapstructure:"uuid, omitempty"`
	Status                   *string `json:"status" mapstructure:"status, omitempty"`
	Version                  *string `json:"version" mapstructure:"version, omitempty"`
	KubeApiServerIPv4Address *string `json:"kubeapi_server_ipv4_address" mapstructure:"kubeapi_server_ipv4_address, omitempty"`
	ETCDConfig               struct {
		NodePools []string `json:"node_pools" mapstructure:"node_pools, omitempty"`
	} `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
	MasterConfig struct {
		DeploymentType string   `json:"deployment_type" mapstructure:"deployment_type, omitempty"`
		NodePools      []string `json:"node_pools" mapstructure:"node_pools, omitempty"`
	} `json:"master_config" mapstructure:"master_config, omitempty"`
	WorkerConfig struct {
		NodePools []string `json:"node_pools" mapstructure:"node_pools, omitempty"`
	} `json:"worker_config" mapstructure:"worker_config, omitempty"`
}

type KarbonClusterNodePool struct {
	AHVConfig     *KarbonClusterNodePoolAHVConfig    `json:"ahv_config" mapstructure:"ahv_config, omitempty"`
	Name          *string                            `json:"name" mapstructure:"name, omitempty"`
	NodeOSVersion *string                            `json:"node_os_version" mapstructure:"node_os_version, omitempty"`
	NumInstances  *int64                             `json:"num_instances" mapstructure:"num_instances, omitempty"`
	Nodes         *[]KarbonClusterNodeIntentResponse `json:"nodes" mapstructure:"nodes, omitempty"`
}

type KarbonClusterNodeIntentResponse struct {
	Hostname    *string `json:"hostname" mapstructure:"hostname, omitempty"`
	IPv4Address *string `json:"ipv4_address" mapstructure:"ipv4_address, omitempty"`
}

type KarbonClusterKubeconfigResponse struct {
	KubeConfig string `json:"kube_config" mapstructure:"kube_config, omitempty"`
}

//inputs
type KarbonClusterIntentInput struct {
	Name               string                                     `json:"name" mapstructure:"name, omitempty"`
	Version            string                                     `json:"version" mapstructure:"version, omitempty"`
	CNIConfig          KarbonClusterCNIConfigIntentInput          `json:"cni_config" mapstructure:"cni_config, omitempty"`
	ETCDConfig         KarbonClusterETCDConfigIntentInput         `json:"etcd_config" mapstructure:"etcd_config, omitempty"`
	MastersConfig      KarbonClusterMasterConfigIntentInput       `json:"masters_config" mapstructure:"masters_config, omitempty"`
	Metadata           KarbonClusterMetadataIntentInput           `json:"metadata" mapstructure:"metadata, omitempty"`
	StorageClassConfig KarbonClusterStorageClassConfigIntentInput `json:"storage_class_config" mapstructure:"storage_class_config, omitempty"`
	WorkersConfig      KarbonClusterWorkerConfigIntentInput       `json:"workers_config" mapstructure:"workers_config, omitempty"`
}
type KarbonClusterMetadataIntentInput struct {
	APIVersion string `json:"api_version" mapstructure:"api_version, omitempty"`
}

type KarbonClusterMasterConfigIntentInput struct {
	SingleMasterConfig  *KarbonClusterSingleMasterConfigIntentInput        `json:"single_master_config" mapstructure:"single_master_config, omitempty"`
	ActivePassiveConfig *KarbonClusterActivePassiveMasterConfigIntentInput `json:"active_passive_config" mapstructure:"active_passive_config, omitempty"`
	ExternalLBConfig    *KarbonClusterExternalLBMasterConfigIntentInput    `json:"external_lb_config" mapstructure:"external_lb_config, omitempty"`
	NodePools           []KarbonClusterNodePool                            `json:"node_pools" mapstructure:"node_pools, omitempty"`
}

type KarbonClusterActivePassiveMasterConfigIntentInput struct {
	ExternalIPv4Address string `json:"external_ipv4_address" mapstructure:"external_ipv4_address, omitempty"`
}

type KarbonClusterExternalLBMasterConfigIntentInput struct {
	ExternalIPv4Address string                                           `json:"external_ipv4_address" mapstructure:"external_ipv4_address, omitempty"`
	MasterNodesConfig   []KarbonClusterMasterNodeMasterConfigIntentInput `json:"master_nodes_config" mapstructure:"master_nodes_config, omitempty"`
}

type KarbonClusterMasterNodeMasterConfigIntentInput struct {
	IPv4Address  string `json:"ipv4_address" mapstructure:"ipv4_address, omitempty"`
	NodePoolName string `json:"node_pool_name" mapstructure:"node_pool_name, omitempty"`
}

type KarbonClusterSingleMasterConfigIntentInput struct {
}

type KarbonClusterWorkerConfigIntentInput struct {
	NodePools []KarbonClusterNodePool `json:"node_pools" mapstructure:"node_pools, omitempty"`
}
type KarbonClusterETCDConfigIntentInput struct {
	NodePools []KarbonClusterNodePool `json:"node_pools" mapstructure:"node_pools, omitempty"`
}

type KarbonClusterCNIConfigIntentInput struct {
	NodeCIDRMaskSize int64                                  `json:"node_cidr_mask_size" mapstructure:"node_cidr_mask_size, omitempty"`
	PodIPv4CIDR      string                                 `json:"pod_ipv4_cidr" mapstructure:"pod_ipv4_cidr, omitempty"`
	ServiceIPv4CIDR  string                                 `json:"service_ipv4_cidr" mapstructure:"service_ipv4_cidr, omitempty"`
	FlannelConfig    *KarbonClusterFlannelConfigIntentInput `json:"flannel_config" mapstructure:"flannel_config, omitempty"`
	CalicoConfig     *KarbonClusterCalicoConfigIntentInput  `json:"calico_config" mapstructure:"calico_config, omitempty"`
}

type KarbonClusterCalicoConfigIntentInput struct {
	IpPoolConfigs []KarbonClusterCalicoConfigIpPoolConfigIntentInput `json:"ip_pool_configs" mapstructure:"ip_pool_configs,omitempty"`
}

type KarbonClusterCalicoConfigIpPoolConfigIntentInput struct {
	CIDR string `json:"cidr" mapstructure:"cidr"`
}

type KarbonClusterFlannelConfigIntentInput struct{}

type KarbonClusterNodePoolAHVConfig struct {
	CPU                     int64  `json:"cpu" mapstructure:"cpu, omitempty"`
	DiskMib                 int64  `json:"disk_mib" mapstructure:"disk_mib, omitempty"`
	MemoryMib               int64  `json:"memory_mib" mapstructure:"memory_mib, omitempty"`
	NetworkUUID             string `json:"network_uuid" mapstructure:"network_uuid, omitempty"`
	PrismElementClusterUUID string `json:"prism_element_cluster_uuid" mapstructure:"prism_element_cluster_uuid, omitempty"`
}

type KarbonClusterStorageClassConfigIntentInput struct {
	DefaultStorageClass bool                                  `json:"default_storage_class" mapstructure:"default_storage_class, omitempty"`
	Name                string                                `json:"name" mapstructure:"name, omitempty"`
	ReclaimPolicy       string                                `json:"reclaim_policy" mapstructure:"reclaim_policy, omitempty"`
	VolumesConfig       KarbonClusterVolumesConfigIntentInput `json:"volumes_config" mapstructure:"volumes_config, omitempty"`
}

type KarbonClusterVolumesConfigIntentInput struct {
	FileSystem              string `json:"file_system" mapstructure:"file_system, omitempty"`
	FlashMode               bool   `json:"flash_mode" mapstructure:"flash_mode, omitempty"`
	Password                string `json:"password" mapstructure:"password, omitempty"`
	PrismElementClusterUUID string `json:"prism_element_cluster_uuid" mapstructure:"prism_element_cluster_uuid, omitempty"`
	StorageContainer        string `json:"storage_container" mapstructure:"storage_container, omitempty"`
	Username                string `json:"username" mapstructure:"username, omitempty"`
}

//KARBON shared

type KarbonClusterActionResponse struct {
	ClusterName string `json:"cluster_name" mapstructure:"cluster_name, omitempty"`
	ClusterUUID string `json:"cluster_uuid" mapstructure:"cluster_uuid, omitempty"`
	TaskUUID    string `json:"task_uuid" mapstructure:"task_uuid, omitempty"`
}

type KarbonClusterKubeconfig struct {
	APIVersion string `yaml:"apiVersion" mapstructure:"apiVersion, omitempty"`
	Kind       string `yaml:"kind" mapstructure:"kind, omitempty"`
	Clusters   []struct {
		Name    string `yaml:"name" mapstructure:"name, omitempty"`
		Cluster struct {
			Server                   string `yaml:"server" mapstructure:"server, omitempty"`
			CertificateAuthorityData string `yaml:"certificate-authority-data" mapstructure:"certificate-authority-data, omitempty"`
		} `yaml:"cluster" mapstructure:"cluster, omitempty"`
	} `yaml:"clusters" mapstructure:"clusters, omitempty"`
	Users []struct {
		Name string `yaml:"name" mapstructure:"name, omitempty"`
		User struct {
			Token string `yaml:"token" mapstructure:"token, omitempty"`
		} `yaml:"user" mapstructure:"user, omitempty"`
	} `yaml:"users" mapstructure:"users, omitempty"`
	Contexts []struct {
		Context struct {
			Cluster string `yaml:"cluster" mapstructure:"cluster, omitempty"`
			User    string `yaml:"user" mapstructure:"user, omitempty"`
		} `yaml:"context" mapstructure:"context, omitempty"`
		Name string `yaml:"name" mapstructure:"name, omitempty"`
	} `yaml:"contexts" mapstructure:"contexts, omitempty"`
	CurrentContext string `yaml:"current-context" mapstructure:"current-context, omitempty"`
}

type KarbonClusterSSHconfig struct {
	Certificate string `json:"certificate" mapstructure:"certificate, omitempty"`
	ExpiryTime  string `json:"expiry_time" mapstructure:"expiry_time, omitempty"`
	PrivateKey  string `json:"private_key" mapstructure:"private_key, omitempty"`
	Username    string `json:"username" mapstructure:"username, omitempty"`
}
