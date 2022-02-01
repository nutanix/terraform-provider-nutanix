package fc

import (
	"time"
)

type ErrorResponse struct{
	Code *int32
	MessageList []*string
}


// Metadata for List Operations Input
type ListMetadataInput struct {
	Length *int64 `json:"length,omitempty"`
	Offset *int64 `json:"offset,omitempty"`
}

// Metadata for List Operations Output
type ListMetadataOutput struct {
	TotalMatches *int64 `json:"total_matches,omitempty"`
	Length *int64 `json:"length,omitempty"`
	Offset *int64 `json:"offset,omitempty"`
}

// CommonNetworkSetting ...
type CommonNetworkSetting struct {
	CvmDnsServers []*string `json:"cvm_dns_servers,omitempty"`
	HypervisorDnsServers []*string `json:"hypervisor_dns_servers,omitempty"`
	CvmNtpServers[]*string `json:"cvm_ntp_servers,omitempty"`
	HypervisorNtpServers []*string `json:"hypervisor_ntp_servers,omitempty"`
}

type ImagedNodeListFilter struct {
	NodeState *string `json:"node_state,omitempty"`
}

type HardwareAttribute struct {

}

// ImagedNodeDetails ...
type ImagedNodeDetails struct{
	CvmVlanID *int64 `json:"cvm_vlan_id,omitempty"`
	NodeType *string `json:"node_type,omitempty"`
	CreatedTimestamp *string `json:"created_timestamp,omitempty"`
	Ipv6Interface *string `json:"ipv6_interface,omitempty"`
	APIKeyUUID *string `json:"api_key_uuid,omitempty"`
	FoundationVersion *string `json:"foundation_version,omitempty"`
	CurrentTime *string `json:"current_time,omitempty"`
	NodePosition *string `json:"node_position,omitempty"`
	CvmNetmask *string `json:"cvm_netmask,omitempty"`
	IpmiIP *string `json:"ipmi_ip,omitempty"`
	CvmUUID *string `json:"cvm_uuid,omitempty"`
	CvmIpv6 *string `json:"cvm_ipv6,omitempty"`
	ImagedClusterUuid *string `json:"imaged_cluster_uuid,omitempty"`
	CvmUp *bool `json:"cvm_up,omitempty"`
	Available *bool `json:"available,omitempty"`
	ObjectVersion *int64 `json:"object_version,omitempty"`
	IpmiNetmask *string `json:"ipmi_netmask,omitempty"`
	HypervisorHostname *string `json:"hypervisor_hostname,omitempty"`
	NodeState *string `json:"node_state,omitempty"`
	HypervisorVersion *string `json:"hypervisor_version,omitempty"`
	HypervisorIP *string `json:"hypervisor_ip,omitempty"`
	Model *string `json:"model,omitempty"`
	IpmiGateway *string `json:"ipmi_gateway,omitempty"`
	HardwareAttributes *HardwareAttribute `json:"hardware_attributes,omitempty"`
	CvmGateway *string `json:"cvm_gateway,omitempty"`
	NodeSerial *string `json:"node_serial,omitempty"`
	ImagedNodeUUID *string `json:"imaged_node_uuid,omitempty"`
	BlockSerial *string `json:"block_serial,omitempty"`
	HypervisorType *string `json:"hypervisor_type,omitempty"`
	LatestHbTsList []*string `json:"latest_hb_ts_list,omitempty"`
	HypervisorNetmask *string `json:"hypervisor_netmask,omitempty"`
	HypervisorGateway *string `json:"hypervisor_gateway,omitempty"`
	CvmIP *string `json:"cvm_ip,omitempty"`
	AosVersion *string `json:"aos_version,omitempty"`
}

// ImagedNodesInput ...
type ImagedNodesInput struct {
	CvmVlanID *int64 `json:"cvm_vlan_id,omitempty"`
	NodeType *string `json:"node_type,omitempty"`
	Ipv6Interface *string `json:"ipv6_interface,omitempty"`
	FoundationVersion *string `json:"foundation_version,omitempty"`
	IpmiNetmask *string `json:"ipmi_netmask,omitempty"`
	CvmNetmask *string `json:"cvm_netmask,omitempty"`
	IpmiIP *string `json:"ipmi_ip,omitempty"`
	CvmUUID *string `json:"cvm_uuid,omitempty"`
	CvmIpv6 *string `json:"cvm_ipv6,omitempty"`
	CvmUp *bool `json:"cvm_up,omitempty"`
	NodePosition *string `json:"node_position,omitempty"`
	HypervisorHostname *string `json:"hypervisor_hostname,omitempty"`
	HypervisorVersion *string `json:"hypervisor_version,omitempty"`
	HypervisorIP *string `json:"hypervisor_ip,omitempty"`
	CvmIP *string `json:"cvm_ip,omitempty"`
	IpmiGateway *string `json:"ipmi_gateway,omitempty"`
	HardwareAttributes *HardwareAttribute `json:"hardware_attributes,omitempty"`
	CvmGateway *string `json:"cvm_gateway,omitempty"`
	NodeSerial *string `json:"node_serial,omitempty"`
	BlockSerial *string `json:"block_serial,omitempty"`
	HypervisorType *string `json:"hypervisor_type,omitempty"`
	HypervisorNetmask *string `json:"hypervisor_netmask,omitempty"`
	HypervisorGateway *string `json:"hypervisor_gateway,omitempty"`
	Model *string `json:"model,omitempty"`
	AosVersion *string `json:"aos_version,omitempty"`
}

// ImagedNodeResponse ...
type ImagedNodesResponse struct {
	ObjectVersion *int64 `json:"object_version,omitempty"`
	ImagedNodeUUID *string  `json:"imaged_node_uuid,omitempty"`
}

type ImagedNodesListInput struct {
	Length *int64 `json:"object_version,omitempty"`
	Filters *ImagedNodeListFilter `json:"object_version,omitempty"`
	Offset *int64 `json:"object_version,omitempty"`
}

type ImagedNodesListResponse struct {
	Metadata *ListMetadataOutput `json:"metadata,omitempty"`
	ImagedNodes []*ImagedNodeDetails `json:"imaged_nodes,omitempty"`
}