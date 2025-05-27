---
layout: "nutanix"
page_title: "NUTANIX: nutanix_pc_registration_v2 "
sidebar_current: "docs-nutanix-resource-pc-registration-v2"
description: |-
  This operation Registers a domain manager (Prism Central) instance to other entities like PE and PC. This process is asynchronous, creating a registration task and returning its UUID.


---

# nutanix_pc_registration_v2

Provides a resource to Registers a domain manager (Prism Central) instance to other entities like PE and PC. This process is asynchronous, creating a registration task and returning its UUID.


## Example Usage

```hcl


// DomainManagerRemoteClusterSpec
resource "nutanix_pc_registration_v2 " "pc-domain-manager"{
  pc_ext_id = "00000000-0000-0000-0000-000000000000"
  remote_cluster {
    domain_manager_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = "0.0.0.0"
          }
        }
        credentials {
          authentication {
            username = "example"
            password = "example.123"
          }
        }
      }
      cloud_type = "NUTANIX_HOSTED_CLOUD"
    }
  }
}

// AOSRemoteClusterSpec
resource "nutanix_pc_registration_v2 " "pc-aos"{
  pc_ext_id = "00000000-0000-0000-0000-000000000000"
  remote_cluster {
    aos_remote_cluster_spec {
      remote_cluster {
        address {
          ipv4 {
            value = "0.0.0.0"
          }
        }
        credentials {
          authentication {
            username = "example"
            password = "example.123"
          }
        }
      }
    }
  }
}

// ClusterReference
resource "nutanix_pc_registration_v2 " "pc-cluster-reference"{
  pc_ext_id = "00000000-0000-0000-0000-000000000000"
  remote_cluster {
    cluster_reference {
      ext_id = "11111111-1111-1111-1111-111111111111"
    }
  }
}
```

## Argument Reference
The following arguments are supported:


* `pc_ext_id`: -(Required) The display name for the Role.
* `remote_cluster`: -(Required)  The registration request consists of the remote cluster details. Credentials must be of domain manager (Prism Central) role.
  The remote cluster details are different based on the object type. The object type is used to determine the type of remote cluster. The object type can be one of the following:
  * `prism.v4.management.DomainManagerRemoteClusterSpec`
  * `prism.v4.management.AOSRemoteClusterSpec`
  * `prism.v4.management.ClusterReference`


### Remote Cluster
The remote_cluster argument supports the following, depending on the object type:

> only one of domain_manager_remote_cluster_spec, aos_remote_cluster_spec, cluster_reference can be used at a time.


* `domain_manager_remote_cluster_spec`: - The registration request consists of the remote cluster details. and cloud type.
* `aos_remote_cluster_spec`: - The registration request consists of the remote cluster details.
* `cluster_reference`: - The registration request consists of the remote cluster details. using the cluster reference.

### DomainManagerRemoteClusterSpec
the `domain_manager_remote_cluster_spec` argument supports the following:
* `remote_cluster`: -(Required)  Address configuration of a remote cluster. It requires the address of the remote, that is an IP or domain name along with the basic authentication credentials.
* `cloud_type`: -(Required)  Enum denoting whether the domain manager (Prism Central) instance is reachable with its physical address or reachable through the My Nutanix portal. Based on the above description, the allowed enum values are:
  * `NUTANIX_HOSTED_CLOUD` : Domain manager (Prism Central) reachable through My Nutanix portal.
  * `ONPREM_CLOUD`: Domain manager (Prism Central) reachable on it's physical address.

### AOSRemoteClusterSpec
The `aos_remote_cluster_spec` argument supports the following:
* `remote_cluster`: -(Required)  Address configuration of a remote cluster. It requires the address of the remote, that is an IP or domain name along with the basic authentication credentials.

### ClusterReference
The `cluster_reference` argument supports the following:
* `ext_id`: -(Required)  Cluster UUID of a remote cluster.


#### Remote Cluster Details
The remote_cluster argument for `prism.v4.management.DomainManagerRemoteClusterSpec` and `prism.v4.management.AOSRemoteClusterSpec` supports the following:

* `address`: -(Required)  An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format or a Fully Qualified Domain Name.
* `credentials`: -(Required)  Credentials to connect to a remote cluster.

##### Address
The address argument supports the following:

* `ipv4`: - An unique address that identifies a device on the internet or a local network in IPv4 format.
* `ipv6`: - An unique address that identifies a device on the internet or a local network in IPv6 format.
* `fqdn`: - A fully qualified domain name that specifies its exact location in the tree hierarchy of the Domain Name System.

###### IPV4, IPV6

The ipv4 attribute supports the following:

* `value`: - An unique address that identifies a device on the internet or a local network in IPv4/IPv6 format.
* `prefix_length`: - The prefix length of the network to which this host IPv4/IPv6 address belongs.

###### FQDN

The fqdn attribute supports the following:

* `value`: - The fully qualified domain name of the host.


##### Credentials
The credentials argument supports the following:

* `authentication`: -(Required)  An authentication scheme that requires the client to present a username and password. The server will service the request only if it can validate the user-ID and password for the protection space of the Request-URI.

###### Authentication
The authentication argument supports the following:

* `username`: -(Required)  Username required for the basic auth scheme. As per RFC 2617 usernames might be case sensitive.
* `password`: -(Required)  Password required for the basic auth scheme.

See detailed information in [Nutanix Register a PC Docs](https://developers.nutanix.com/api-reference?namespace=prism&version=v4.0#tag/DomainManager/operation/register).
