# discover-nodes-network-details

This module is used to discover nodes and get node network details of nodes which are not a part of any cluster.

Note : This module was created for internal use in aos and dos based modules.

## Resources & DataSources used

1. nutanix_foundation_discover_nodes
2. nutanix_foundation_node_network_details

## Usage

Basic example of usage. 

```hcl
module "discovered_nodes_network_details" {
    // 
    source = "<local-path-to-nutanix-terraform-provider-repo>/terraform-provider-nutanix/modules/foundation/discover-nodes-network-details/"
}

```

Check output.tf for exported fields from this module