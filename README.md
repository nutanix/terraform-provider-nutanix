# Terraform Nutanix Provider

Terraform provider plugin to integrate with Nutanix Enterprise Cloud

NOTE: The latest version of the Nutanix provider is [v1.8.0-beta.1](https://github.com/nutanix/terraform-provider-nutanix/releases/tag/v1.8.0-beta.1)

Modules based on Terraform Nutanix Provider can be found here : [Modules](https://github.com/nutanix/terraform-provider-nutanix/tree/master/modules)
## Build, Quality Status

 [![Go Report Card](https://goreportcard.com/badge/github.com/nutanix/terraform-provider-nutanix)](https://goreportcard.com/report/github.com/nutanix/terraform-provider-nutanix)
<!-- [![Maintainability](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/maintainability)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/test_coverage)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/test_coverage) -->

| Master                                                                                                                                                          | Develop                                                                                                                                                           |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [![Build Status](https://travis-ci.org/nutanix/terraform-provider-nutanix.svg?branch=master)](https://travis-ci.org/nutanix/terraform-provider-nutanix) | [![Build Status](https://travis-ci.org/nutanix/terraform-provider-nutanix.svg?branch=develop)](https://travis-ci.org/nutanix/terraform-provider-nutanix) |

## Community

Nutanix is taking an inclusive approach to developing this new feature and welcomes customer feedback. Please see our development project on GitHub (you're here!), comment on requirements, design, code, and/or feel free to join us on Slack. Instructions on commenting, contributing, and joining our community Slack channel are all located within our GitHub Readme.

For a slack invite, please contact terraform@nutanix.com from your business email address, and we'll add you.

### Provider Development
* [Terraform](https://www.terraform.io/downloads.html) 0.12+
* [Go](https://golang.org/doc/install) 1.17+ (to build the provider plugin)
* This provider uses [SDKv2](https://www.terraform.io/plugin/sdkv2/sdkv2-intro) from release 1.3.0

### Provider Use

The Terraform Nutanix provider is designed to work with Nutanix Prism Central and Standalone Foundation, such that you can manage one or more Prism Element clusters at scale. AOS/PC 5.6.0 or higher is required, as this Provider makes exclusive use of the v3 APIs. It also consists components to work with Foundation to performing node imaging and related activities.

> For the 1.2.0 release of the provider it will have an N-1 compatibility with the Prism Central APIs. This provider was tested against Prism Central versions 2020.9 and 2020.11, as well as AOS version 5.18 and 5.19


> For the 1.3.0 release of the provider it will have N-2 compatibility with the Prism Central APIs. This release was tested against Prism Central versions pc.2021.9.0.4, pc.2021.8.0.1 and pc.2021.7.


> For the 1.4.0 & 1.4.1 release of the provider it will have N-2 compatibility with the Prism Central APIs. This release was tested against Prism Central versions pc2022.1 pc.2021.9.0.4 and pc.2021.8.0.1.  

> For the 1.5.0 release of the provider it will have N-2 compatibility with the Prism Central APIs. This release was tested against Prism Central versions pc2022.1.0.2 pc.2021.9.0.4 and pc.2021.8.0.1.

> For the 1.6.1 release of the provider it will have N-2 compatibility with the Prism Central APIs. This release was tested against Prism Central versions pc2022.4 pc2022.1.0.2 and pc2021.9.0.4.

> For the 1.7.0 release of the provider it will have N-2 compatibility with the Prism Central APIs. This release was tested against Prism Central versions pc2022.6, pc2022.4 and pc2022.1.0.2.

> For the 1.7.1 release of the provider it will have N-2 compatibility with the Prism Central APIs. This release was tested against Prism Central versions pc2022.6, pc2022.4.0.1 and pc2022.1.0.2.

### note
With v1.6.1 release of flow networking feature in provider, IAMv2 setups would be mandate. 
Also, there is known issue for access_control_policies resource where update would be failing. We are continuously tracking the issue internally.

with v1.7.0 release of user groups feature in provider, pc version should be minimum 2022.1 to support organisational and saml user group. 

With v1.7.1 release of project internal  in provider is supported. Note to use this, set "use_project_internal" to true. It also enables the ACP mapping with projects. 

## Foundation
> For the 1.5.0-beta release of the provider it will have N-1 compatibility with the Foundation. This release was tested against Foundation versions v5.2 and v5.1.1

> For the 1.5.0 release of the provider it will have N-1 compatibility with the Foundation. This release was tested against Foundation versions v5.2 and v5.1.1

Foundation based examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/foundation/

Foundation based modules & examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/modules/foundation/

## Foundation Central
> For the 1.5.0-beta.2 release of the provider it will have N-1 compatibility with the Foundation Central. This release was tested with v1.2 and v1.3 Foundation Central versions.

> For the 1.5.0 release of the provider it will have N-1 compatibility with the Foundation Central. This release was tested with v1.2 and v1.3 Foundation Central versions.

Foundation Central based examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/foundationCentral/

Foundation Central based modules and examples : Foundation based modules & examples : https://github.com/nutanix/terraform-provider-nutanix/blob/master/modules/foundationCentral/

## Nutanix Database Service
> For the 1.8.0-beta.1 release of the provider, it will have N-1 compatibility with the Nutanix database service. This release was tested with v2.4 and v2.4.1 versions.

Note: For 1.8.0-beta.1 release, only postgress database type is qualified and officially supported.

Checkout example : https://github.com/nutanix/terraform-provider-nutanix/blob/master/examples/ndb/database_instance

## Example Usage

See the Examples folder for a handful of main.tf demos as well as some pre-compiled binaries.

We'll be refreshing these examples and binaries as we work through tech preview.

Long term, once this is upstream, no pre-compiled binaries will be needed, as terraform will automatically download on use.

## Configuration Reference

The following keys can be used to configure the provider.

* **endpoint** - (Required) IP address for the Nutanix Prism Central.
* **username** - (Required) Username for Nutanix Prism Central. Could be local cluster auth (e.g. `auth`) or directory auth.
* **password** - (Required) Password for the provided username.
* **port** - (Optional) Port for the Nutanix Prism Central. Default port is 9440.
* **insecure** - (Optional) Explicitly allow the provider to perform insecure SSL requests. If omitted, default value is false.
* **wait_timeout** - (optional) Set if you know that the creation o update of a resource may take long time (minutes).

```hcl
provider "nutanix" {
  username     = "admin"
  password     = "myPassword"
  port         = 9440
  endpoint     = "10.36.7.201"
  insecure     = true
  wait_timeout = 10
}
```

## From terraform-provider-nutanix v1.5.0-beta :

The following keys can be used to configure the provider.

* **endpoint** - (Optional) IP address for the Nutanix Prism Central.
* **username** - (Optional) Username for Nutanix Prism Central. Could be local cluster auth (e.g. `auth`) or directory auth.
* **password** - (Optional) Password for the provided username.
* **port** - (Optional) Port for the Nutanix Prism Central. Default port is 9440.
* **insecure** - (Optional) Explicitly allow the provider to perform insecure SSL requests. If omitted, default value is false.
* **wait_timeout** - (optional) Set if you know that the creation or update of a resource may take long time (minutes).
* **foundation_endpoint** - (optional) IP address of foundation vm.
* **foundation_port** - (optional) Port of foundation vm. Default port is 8000.

```hcl
provider "nutanix" {
  username            = "admin"
  password            = "myPassword"
  port                = 9440
  endpoint            = "10.36.7.201"
  insecure            = true
  wait_timeout        = 10
  foundation_endpoint = "10.xx.xx.xx"
  foundation_port     = 8000
}
```

## Additional fields for using Nutanix Database Service:

* **ndb_username** - (Optional) Username of Nutanix Database Service server
* **ndb_password** - (Optional) Password of Nutanix Database Service server
* **ndb_endpoint** - (Optional) IP of Nutanix Database Service server

```hcl
provider "nutanix" {
  ndb_username = var.ndb_username
  ndb_password = var.ndb_password
  ndb_endpoint = var.ndb_endpoint
}
```

### Provider Configuration Requirements & Warnings
From foundation getting released in 1.5.0-beta, provider configuration will accomodate prism central and foundation apis connection details. **It will show warnings for disabled api connections as per the attributes given in provider configuration in above mentioned format**. The below are the required attributes for corresponding provider componenets :
* endpoint, username and password are required fields for using Prism Central & Karbon based resources and data sources
* foundation_endpoint is required field for using Foundation based resources and data sources
* ndb_username, ndb_password and ndb_endpoint are required fields for using NDB based resources and data sources
## Resources

* nutanix_access_control_policy
* nutanix_category_key
* nutanix_category_value
* nutanix_image
* nutanix_karbon_cluster
* nutanix_karbon_private_registry
* nutanix_network_security_rule
* nutanix_project
* nutanix_protection_rule
* nutanix_recovery_plan
* nutanix_role
* nutanix_subnet
* nutanix_user
* nutanix_virtual_machine
* nutanix_service_group
* nutanix_address_group
* nutanix_foundation_image_nodes
* nutanix_foundation_ipmi_config
* nutanix_foundation_image
* nutanix_foundation_central_api_keys
* nutanix_foundation_central_image_cluster
* nutanix_vpc
* nutanix_pbr
* nutanix_static_routes
* nutanix_floating_ip
* nutanix_user_groups
* nutanix_ndb_database

## Data Sources

* nutanix_access_control_policies
* nutanix_access_control_policy
* nutanix_category_key
* nutanix_cluster
* nutanix_clusters
* nutanix_host
* nutanix_hosts
* nutanix_image
* nutanix_karbon_cluster_kubeconfig
* nutanix_karbon_cluster_ssh
* nutanix_karbon_cluster
* nutanix_karbon_clusters
* nutanix_karbon_private_registries
* nutanix_karbon_private_registry
* nutanix_network_security_rule
* nutanix_permission
* nutanix_permissions
* nutanix_project
* nutanix_projects
* nutanix_role
* nutanix_roles
* nutanix_subnet
* nutanix_subnets
* nutanix_user_group
* nutanix_user_groups
* nutanix_user
* nutanix_users
* nutanix_virtual_machine
* nutanix_protection_rule
* nutanix_protection_rules
* nutanix_recovery_plan
* nutanix_recovery_plans
* nutanix_address_groups
* nutanix_address_group
* nutanix_foundation_discover_nodes
* nutanix_foundation_node_network_details
* nutanix_foundation_nos_packages
* nutanix_foundation_hypervisor_isos
* nutanix_foundation_central_api_keys
* nutanix_foundation_central_list_api_keys
* nutanix_foundation_central_imaged_nodes_list
* nutanix_foundation_central_imaged_clusters_list
* nutanix_foundation_central_cluster_details
* nutanix_foundation_central_imaged_node_details
* nutanix_vpc
* nutanix_vpcs
* nutanix_pbr
* nutanix_pbrs
* nutanix_floating_ip
* nutanix_floating_ips
* nutanix_static_routes
* nutanix_ndb_cluster
* nutanix_ndb_clusters
* nutanix_ndb_database
* nutanix_ndb_databases
* nutanix_ndb_profile
* nutanix_ndb_profiles
* nutanix_ndb_sla
* nutanix_ndb_slas

## Quick Install

### Install Dependencies

* [Terraform](https://www.terraform.io/downloads.html) 0.12+

### For developing or build from source


* [Go](https://golang.org/doc/install) 1.12+ (to build the provider plugin)


### Building/Developing Provider

We recomment to use Go 1.12+ to be able to use `go modules`

```sh
$ git clone https://github.com/nutanix/terraform-provider-nutanix.git
```

Enter the provider directory and build the provider

```sh
$ make tools
$ make build
```

This will create a binary file `terraform-provider-nutanix` you can copy to your terraform specific project.

Alternative build: with our demo

```sh
$ make tools
$ go build -o examples/terraform-provider-nutanix
$ cd examples
$ terraform init #to try out our demo
```

If you need multi-OS binaries such as Linux, macOS, Windows. Run the following command.

```sh
$ make tools
$ make cibuild
```

This command will create a `pkg/` directory with all the binaries for the most popular OS.

### Running tests of provider

For running unit tests:
```sh
make test
```

For running integration tests:

1. Add environment variables for setup related details:
```ssh
export NUTANIX_USERNAME="<username>"
export NUTANIX_PASSWORD="<password>"
export NUTANIX_INSECURE=true
export NUTANIX_PORT=9440
export NUTANIX_ENDPOINT="<pc-ip>"
export NUTANIX_STORAGE_CONTAINER="<storage-container-uuid-for-vm-tests>"
export FOUNDATION_ENDPOINT="<foundation-vm-ip-for-foundation-related-tests>"
export FOUNDATION_PORT=8000
export NOS_IMAGE_TEST_URL="<test-image-url>"
```

2. Some tests need setup related constants for resource creation. So add/replace details in test_config.json (for pc tests) and test_foundation_config.json (for foundation and foundation central tests)

3. To run all tests:
```ssh
make testacc
```

4. To run specific tests:
```ssh 
export TESTARGS='-run=TestAccNutanixPbr_WithSourceExternalDestinationNetwork'
make testacc
```

5. To run collection of tests:
``` ssh
export TESTARGS='-run=TestAccNutanixPbr*'
make testacc
```

### Common Issues using the development binary.

Terraform download the released binary instead developent one.

Just follow this steps to get the development binary:

1. Copy the development terraform binary in the root folder of the project (i.e. where your main.tf is), this should be named `terraform-provider-nutanix`
2. Remove the entire “.terraform” directory.
    ```sh
    rm -rf .terraform/
    ```

3. Run the following command in the same folder where you have copied the development terraform binary.
    ```sh
    terraform init -upgrade
    terraform providers -version
    ```

4. You should see version as “nutanix (unversioned)”
5. Then run your main.tf

## Release it

1. Install `goreleaser` tool:

    ```bash
    go get -v github.com/goreleaser/goreleaser
    cd $GOPATH/src/github.com/goreleaser/goreleaser
    go install
    ```

    Alternatively you can download a latest release from [goreleaser Releases Page](https://github.com/goreleaser/goreleaser/releases)

1. Clean up folder `(builds)` if exists

1. Make sure that the repository state is clean:

    ```bash
    git status
    ```

1. Tag the release:

    ```bash
    git tag v1.1.0
    ```

1. Run `goreleaser`:

    ```bash
    cd (TODO: go dir)
    goreleaser --skip-publish v1.1.0
    ```

1. Check builds inside `(TODO: build dir)` directory.

1. Publish release tag to GitHub:

    ```bash
    git push origin v1.1.0
    ```

## Additional Resources

We've got a handful of resources outside of this repository that will help users understand the interactions between terraform and Nutanix

* YouTube
  _ Overview Video: [](https://www.youtube.com/watch?v=V8_Lu1mxV6g)
  _ Working with images: [](https://www.youtube.com/watch?v=IW0eQevZ73I)
* Nutanix GitHub
  _ [](https://github.com/nutanix/terraform-provider-nutanix)
  _ Private repo until code goes upstream
* Jon’s GitHub
  _ [](https://github.com/JonKohler/ThisOldCloud/tree/master/Terraform-Nutanix)
  _ Contains sample TF’s and PDFs from the youtube videos
* Slack channel \* User community slack channel is available on nutanix.slack.com. Email terraform@nutanix.com to gain entry.
