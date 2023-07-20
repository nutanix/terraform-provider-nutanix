---
layout: "nutanix"
page_title: "NUTANIX: nutanix_foundation_central_cluster_details"
sidebar_current: "docs-nutanix-datasource-foundation-central-cluster-details"
description: |-
 Details of the cluster with the UUID. 
---

# nutanix_foundation_central_cluster_details

Get a cluster details created using Foundation Central.

## Example Usage

```hcl
data "nutanix_foundation_central_cluster_details" "imaged_cluster_details" {
    imaged_cluster_uuid = "<CLUSTER-UUID>"
}
```

## Argument Reference

*`imaged_cluster_uuid`: UUID of the cluster whose details need to be fetched.

## Attribute Reference

The following attributes are exported:

### imaged_clusters
* `current_time`:Current time of Foundation Central.
* `archived`: True if the cluster creation request is archived, False otherwise
* `cluster_external_ip`: External management ip of the cluster.
* `imaged_node_uuid_list`: List of UUIDs of imaged nodes.
* `common_network_settings`: Common network settings across the nodes in the cluster.
* `storage_node_count`: Number of storage only nodes in the cluster. AHV iso for storage node will be taken from aos package.
* `redundancy_factor`: Redundancy factor of the cluster.
* `foundation_init_node_uuid`: UUID of the first node coordinating cluster creation.
* `workflow_type`: If imaging and cluster creation is coordinated by Foundation, value will be FOUNDATION_WF. If the nodes are in phoenix, value will be PHOENIX_WF.
* `cluster_name`: Name of the cluster.
* `foundation_init_config`: Json config used by Foundation to create the cluster.
* `cluster_status`: Details of cluster creation process.
* `cluster_size`: Number of nodes in the cluster.
* `destroyed`: True if the cluster is destroyed, False otherwise
* `created_timestamp`: Time when the cluster creation request was received in Foundation Central.
* `imaged_cluster_uuid`: UUID of the cluster.


### common network settings
* `cvm_dns_servers`: List of dns servers for the cvms in the cluster.
* `hypervisor_dns_servers`: List of dns servers for the hypervisors in the cluster.
* `cvm_ntp_servers`: List of ntp servers for the cvms in the cluster.
* `hypervisor_ntp_servers`: List of ntp servers for the hypervisors in the cluster.

### cluster status
* `cluster_creation_started`: Denotes whether cluster creation has started in a phoenix workflow. For foundation workflows, this field will be same as intent_picked_up.
* `intent_picked_up`: Denotes whether remote node has picked up the cluster creation intent.
* `imaging_stopped`: Describes whether foundation imaging and cluster creation has stopped. True indicates that process has stopped. False indicates that process is still going on.
* `node_progress_details`: List of progress details of each node.
* `aggregate_percent_complete`: Overall progress percentage including imaging and cluster creation.
* `current_foundation_ip`: Current IP address of the coordinating foundation node.
* `cluster_progress_details`: Denotes the progress status of cluster creation.
* `foundation_session_id`: Foundation session id for cluster creation.

### node progress details
* `status`: Current status of the node imaging process.
* `imaged_node_uuid`: UUID of the node.
* `imaging_stopped`: Describes whether imaging has stopped. True indicates that process has stopped. False indicates that process is still going on. This field will only be used by phoenix nodes to update FC.
* `intent_picked_up`: Denotes whether the remote nodes has picked up the cluster creation intent.
* `percent_complete`: Percent completion of the node imaging process
* `message_list`: List of messages for the client based on process state.

### cluster_progress_details
* `cluster_name`: Cluster name.
* `status`: Current status of cluster creation process.
* `percent_complete`: Percent completion of cluster creation process.
* `message_list`: List of messages for the client based on process state.


See detailed information in [Nutanix Foundation Central Get the details of a cluster](https://www.nutanix.dev/api_references/foundation-central/#/52a237a955f44-get-the-details-of-a-cluster).