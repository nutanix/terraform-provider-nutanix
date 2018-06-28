# Terraform Nutanix Provider

Terraform provider plugin to integrate with Nutanix Enterprise Cloud

NOTE: terraform-provider-nutanix is currently tech preview as of 9 May 2018. See "Current Development Status" below.

#### Project, Build, Quality Status

[![Waffle.io - Columns and their card count](https://badge.waffle.io/b58a88b8face3407eccdc79e19d39058d28754b65f6df7513ed23293ef387491a5f538b65dffd5b12ca0201c6d4b50ce.svg?columns=all)](https://waffle.io/nutanix/terraform-provider-nutanix) [![Go Report Card](https://goreportcard.com/badge/github.com/nutanix/terraform-provider-nutanix)](https://goreportcard.com/report/github.com/nutanix/terraform-provider-nutanix)

[![Maintainability](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/maintainability)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/8b9e61df450276bbdbdb/test_coverage)](https://codeclimate.com/github/nutanix/terraform-provider-nutanix/test_coverage)

| Master                                                                                                                                                          | Develop                                                                                                                                                           |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [![Build Status: master](https://travis-ci.com/nutanix/terraform-provider-nutanix.svg?branch=master)](https://travis-ci.com/nutanix/terraform-provider-nutanix) | [![Build Status: develop](https://travis-ci.com/nutanix/terraform-provider-nutanix.svg?branch=develop)](https://travis-ci.com/nutanix/terraform-provider-nutanix) |

## Community

Nutanix is taking an inclusive approach to developing this new feature and welcomes customer feedback. Please see our development project on GitHub (you're here!), comment on requirements, design, code, and/or feel free to join us on Slack. Instructions on commenting, contributing, and joining our community Slack channel are all located within our GitHub Readme.

For a slack invite, please contact terraform@nutanix.com from your business email address, and we'll add you.

## Current Development Status

### Completed

* [x] Finished VM resource (VM resource and VM resource test.)
* [x] Finished subnet resource (Subnet resource and Subnet resource test.)
* [x] Finished Image resource (Image resource and image resource test.)
* [x] Finished VM data source (VM data source and VM data source test.)
* [x] Finished subnet data source (Subnet data source and Subnet data source test)
* [x] Finished Image data source (Image data source and image data source test.)
* [x] Cluster data source.
* [x] Clusters data source.
* [x] Virtual Machines data source.
* [x] Category keys resource.
* [x] Category values resource.
* [x] Network security rule resource.
* [x] Network security rule data source.

### Currently working on: (3rd Milestone)

* [x] Network security rules data source.
* [x] Subnets data source.
* [x] Images data source.
* [x] Volume group resource.
* [x] Volume group datasource.
* [x] Volume groups datasource.
* [x] Documentation for Resources.s
* [x] Documentation for Datasources.

### Issues

* See open issues on GitHub issues

## Requirements

### Provider Use

* [Terraform](https://www.terraform.io/downloads.html) 0.11.7+
* [Nutanix](https://portal.nutanix.com/#/page/home) Prism Central 5.6.0+
* Note: Nutanix Community Edition will be supported, when an AOS 5.6 based version is released

### Provider Development

* [Go](https://golang.org/doc/install) 1.9+ (to build the provider plugin)

### Provider Use

The Terraform Nutanix provider is designed to work with Nutanix Prism Central, such that you can manage one or more Prism Element clusters at scale. AOS/PC 5.6.0 or higher is required, as this Provider makes exclusive use of the v3 APIs

## Example Usage

See the Examples folder for a handful of main.tf demos as well as some pre-compiled binaries.

We'll be refreshing these examples and binaries as we work through tech preview.

Long term, once this is upstream, no pre-compiled binaries will be needed, as terraform will automatically download on use.

## Configuration Reference

The following keys can be used to configure the provider.

* **endpoint** - (Required) IP address for the Nutanix Prism Element.
* **username** - (Required) Username for Nutanix Prism Element. Could be local cluster auth (e.g. `auth`) or directory auth.
* **password** - (Required) Password for the provided username.
* **port** - (Optional) Port for the Nutanix Prism Element. Default port is 9440.
* **insecure** - (Optional) Explicitly allow the provider to perform insecure SSL requests. If omitted, default value is false.

## Resources

* nutanix_virtual_machine
* nutanix_image
* nutanix_subnet
* nutanix_category_key
* nutanix_category_value
* nutanix_network_security_rule
* nutanix_volume_group

## Data Sources

* nutanix_virtual_machine
* nutanix_virtual_machines
* nutanix_image
* nutanix_images
* nutanix_subnet
* nutanix_subnets
* nutanix_clusters
* nutanix_network_security_rule
* nutanix_network_security_rules
* nutanix_cluster
* nutanix_volume_group
* nutanix_volume_groups

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

## Roadmap

This provider will be released as Tech Preview at .NEXT New Orleans, and is linked into the HashiCorp community providers page, here: [](https://www.terraform.io/docs/providers/type/community-index.html)

We'll be working with HashiCorp as code stabilizes to upstream this properly, at which time we'll PR this entire plugin to the terraform providers org.

* Complete upstream work with successful pull request
  * Note: Depending on external review timelines from HashiCorp and subsequent code change(s) as required
* Add Volume Group resource and data source support \* This is dependent on the VG v3 API, which is currently not GA (work in progress)
* Investigate data protection workflows (likely scoped snapshots, but this may directly conflict with overall pets v cattle)
* Investigate project as a resource and data source, for SSP integration
* Investigate Calm once API constructs are available
* Investigate AFS administration workflows
* Investigate cluster administration APIs (think foundation, admin settings, etc) to do physical provisioning (think DC as Code)

## Quick Install

### Install Dependencies

* Terraform

### For developing or build from source

#### Golang

[](https://github.com/golang/go)

#### Dep: Go dependency management tool

On MacOS you can install or upgrade to the latest released version with Homebrew:

```sh
$ brew install dep
$ brew upgrade dep
```

On other platforms you can use the `install.sh` script:

```sh
$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

It will install into your `$GOPATH/bin` directory by default or any other directory you specify using the `INSTALL_DIRECTORY` environment variable.

If your platform is not supported, you'll need to build it manually or let the team (go/dep team) know and we'll consider adding your platform
to the release builds.

More Info: https://github.com/golang/dep

## Install from source

## Install from package

## Building/Developing Provider

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-nutanix`

```sh
$ mkdir $GOPATH/src/github.com/terraform-providers && cd "$_" #if you already created it only change directory

$ git clone https://github.com/nutanix/terraform-provider-nutanix.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-nutanix
$ make deps
$ make build
```

This will create a binary file `terraform-provider-nutanix` you can copy to your terraform specific project.

Alternative build: with our demo

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-nutanix
$ make deps
$ go build -o examples/terraform-provider-nutanix
$ cd examples
$ terraform init #to try out our demo
```
