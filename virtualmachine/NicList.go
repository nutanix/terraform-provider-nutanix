package virtualmachine

// NicList struct
type NicList struct {

IPEndpointList []IPEndpointList `json:"ip_endpoint_list,omitempty"bson:"ip_endpoint_list,omitempty"`
MacAddress string `json:"mac_address,omitempty"bson:"mac_address,omitempty"`
NetworkFunctionChainReference NetworkFunctionChainReference `json:"network_function_chain_reference,omitempty"bson:"network_function_chain_reference,omitempty"`
NetworkFunctionNicType string `json:"network_function_nic_type,omitempty"bson:"network_function_nic_type,omitempty"`
NicType string `json:"nic_type,omitempty"bson:"nic_type,omitempty"`
SubnetReference SubnetReference `json:"subnet_reference,omitempty"bson:"subnet_reference,omitempty"`

}