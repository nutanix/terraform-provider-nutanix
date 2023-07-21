---
layout: "nutanix"
page_title: "NUTANIX: nutanix_karbon_cluster_ssh"
sidebar_current: "docs-nutanix-datasource-karbon-cluster-ssh"
description: |-
 Describes the SSH config from a Karbon Cluster
---

# nutanix_karbon_cluster_ssh

Describes the SSH config from a Karbon Cluster

## Example Usage

```hcl
# Get ssh credentials by cluster UUID
data "nutanix_karbon_cluster_ssh" "sshbyid" {
    karbon_cluster_id = "<YOUR-CLUSTER-ID>"
}

# Get ssh credentials by cluster name
data "nutanix_karbon_cluster_ssh" "sshbyname" {
    karbon_cluster_name = "<YOUR-CLUSTER-NAME>"
}
```

## Argument Reference

The following arguments are supported:

* `karbon_cluster_id`: Represents karbon cluster uuid
* `karbon_cluster_name`: Represents the name of karbon cluster

## Attribute Reference

The following arguments are supported:

* `certificate` Certificate of the user for SSH access.
* `expiry_time` Timestamp of certificate expiry in the ISO 8601 format (YYYY-MM-DDThh:mm:ss.sssZ).
* `private_key` The private key of the user for SSH access.
* `username` The username for which credentials are returned.

See detailed information in [Get Nutanix Karbon Cluster SSH Credentials](https://www.nutanix.dev/api_references/nke/#/ca0911a3655d3-get-ssh-credentials-to-remotely-access-nodes-that-belong-to-the-kubernetes-cluster).
