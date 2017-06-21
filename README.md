# Terraform AHV Provider Plugin
=================================
Provider plugin to integrate with AHV APIs

- Website: https://www.terraform.io

Requirments
------------

-   [Terraform](https://www.terraform.io/downloads.html) 0.9.x
-   [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin`

```sh
$ mkdir -p $GOPATH/src/github.com/ideadevice; cd $GOPATH/src/github.com/ideadevice
$ git clone git@github.com:ideadevice/terraform-ahv-provider-plugin
```
Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin
$ glide install
$ cd cmd
$ make build
```
Using the  provider
-------------------

The Nutanix Provider is used to interact with AHV APIs. The Provider needs to be configured with the proper credentials before it can be used.

## Example Usage 

```sh
// Configure the Nutanix Provider
provider "nutanix"{
    username  = "username"
    password  = "password"
    endpoint = 1.1.1.1
    insecure = false
    port = 12345
 }

// Create a new instance
resource "nutanix_virtual_machine" "my-machine"{
    ...
}
```
## Configuration Reference 
The following keys can be used to configure the provider.

- **endpoint** - (Required) IP address for the Nutanix Prism Element.
- **username** - (Required) Username for Nutanix Prism Element. Could be local cluster auth (e.g. `auth`) or directory auth.
- **password** - (Required) Password for the provided username.
- **port**     - (Optional) Port for the Nutanix Prism Element. Default port is 9440.
- **insecure** - (Optional) Explicitly allow the provider to perform insecure SSL requests. If omitted, default value is false.

Environment variable HTTP_LOG can be set to define the path of file from which HTTP request logs can be accessed.

Resources
---------

- nutanix_virtual_machine 
-------------------------

Creates, Updates and Destroy virtual machine resource using Prism Element APIs. Example of usage is given at  `$GOPATH/src/github.com/ideadevice/terraform-ahv-provider-plugin/examples/main.tf`

Following are the required arguments :- 

```sh
resource "nutanix_virtual_machine" "my-machine"{
    name = "testname"    
    spec {
        resources = {
            num_vcpus_per_socket = 1
            num_sockets = 1
            memory_size_mib = 1024
            power_state = "POWERED_ON"
        }
    }
    api_version = "3.0"
    metadata = {
        kind = "vm"    
        spec_version = 0
        name = "testname"
    }
}

output "ip" {
    value = "${nutanix_virtual_machine.my-machine.ip_address}"    
}
```
Features :
----------

- Create: 
- Update:
- Destroy:
