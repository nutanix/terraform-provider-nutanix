# ClusterNetwork

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DomainServer** | [**ClusterDomainServer**](cluster_domain_server.md) |  | [optional] [default to null]
**ExternalDataServicesIp** | **string** | The cluster IP address that provides external entities access to various cluster data services.  | [optional] [default to null]
**ExternalIp** | **string** | The local IP of cluster visible externally. | [optional] [default to null]
**ExternalSubnet** | **string** | External subnet for cross server communication. The format is IP/netmask.  | [optional] [default to null]
**HttpProxyList** | [**[]ClusterNetworkEntity**](cluster_network_entity.md) | List of proxies to connect to the service centers. | [optional] [default to null]
**HttpProxyWhitelist** | [**[]HttpProxyWhitelist**](http_proxy_whitelist.md) | List of HTTP proxy whitelist. | [optional] [default to null]
**InternalSubnet** | **string** | The internal subnet is local to every server - its not visible outside.iSCSI requests generated internally within the appliance (by user VMs or VMFS) are sent to the internal subnet. The format is IP/netmask.  | [optional] [default to null]
**NameServerIpList** | **[]string** | The list of IP addresses of the name servers. | [optional] [default to null]
**NfsSubnetWhitelist** | **[]string** | Comma separated list of subnets (of the form &#39;a.b.c.d/l.m.n.o&#39;) that are allowed to send NFS requests to this container. If not specified, the global NFS whitelist will be looked up for access permission. The internal subnet is always automatically considered part of the whitelist, even if the field below does not explicitly specify it. Similarly, all the hypervisor IPs are considered part of the whitelist. Finally, to permit debugging, all of the SVMs local IPs are considered to be implicitly part of the whitelist.  | [optional] [default to null]
**NtpServerIpList** | **[]string** | The list of IP addresses or FQDNs of the NTP servers. | [optional] [default to null]
**SmtpServer** | [**SmtpServer**](smtp_server.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


