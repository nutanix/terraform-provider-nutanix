---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_cluster_kubeconfig"
sidebar_current: "docs-nutanix-datasource-karbon-cluster-kubeconfig"
description: |-
 Describes the SSH config from a Karbon Cluster
---

# nutanix_karbon_cluster_kubeconfig

Describes the SSH config from a Karbon Cluster

## Example Usage

```hcl
# Get kubeconfig by cluster UUID
data "nutanix_karbon_cluster_kubeconfig" "configbyid" {
    karbon_cluster_id = "<YOUR-CLUSTER-ID>"
}


# Get Kubeconfig by cluster name
data "nutanix_karbon_cluster_kubeconfig" "configbyname" {
    karbon_cluster_name = "<YOUR-CLUSTER-NAME>"
}
```

## Argument Reference

The following arguments are supported:

* `karbon_cluster_id`: Represents karbon cluster uuid
* `karbon_cluster_name`: Represents the name of karbon cluster

## Attribute Reference

The following arguments are supported:

* `name` 
* `access_token` 
* `cluster_ca_certificate` 
* `cluster_url`

See detailed information in [Get Nutanix Karbon Cluster KubeConfig](https://www.nutanix.dev/api_references/nke/#/1afde071b42be-get-the-kubeconfig-to-access-the-kubernetes-cluster).
