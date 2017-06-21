# Terraform AHV Provider Plugin
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

- **Create**: This create the new vm. This takes the nested configuration of vm from main.tf and send the POST request on prism v3 api for creating vm. If the http response status something else than 200 - 208 then it gives error. Otherwise it keeps polling till the vm gets created by taking the status state from GET Api call response. If the vm is POWERED_ON and there is atleast one network adapter than it keep polling till the vm gets assigned an ip. Then it sets the ip_address with the ip and
  id of the resource with the vm's uuid.
- **Update**:  This is called to update the properties of the existing vm. For updating the memory and cpu we have to first update power_state to POWER_OFF and then update the memory. With this updates ip_address of the vm also get recomputed.
- **Destroy**: This called to delete the vm. It takes the uuid from the id of the resource and then call DELETE on that uuid.

Environment variable **HTTP_LOG** can be set to define the path of file from which HTTP request logs can be accessed.

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ cd $GOPATH/github.com/ideadevice/terraform-ahv-provider-plugin/cmd
$ make build
...
$ $GOPATH/bin/terraform-provider-nutanix
...
```
In order to test the provider, you can simply run `make test`.

```sh
$ cd $GOPATH/github.com/ideadevice/terraform-ahv-provider-plugin/
$ go test --username=username --password=password --endpoint=1.1.1.1 --insecure=true
```
Following flags are defined for the testcases :- 

```sh
  -apiVersion string
        api_version (default "3.0")
  -diskDeviceType1 string
        disk_device_type_1 (default "DISK")
  -diskDeviceType2 string
        disk_device_type_2 (default "DISK")
  -diskKind1 string
        disk_kind_1 (default "image")
  -diskKind2 string
        disk_kind_2 (default "image")
  -diskName1 string
        disk_name_1 (default "Centos7")
  -diskName2 string
        disk_name_2 (default "Centos7")
  -diskNo string
        disk_No (default "2")
  -diskSize1 string
        disk_size_1 (default "1")
  -diskSize2 string
        disk_size_2 (default "1")
  -diskUUID1 string
        disk_uuid_1 (default "9eabbb39-1baf-4872-beaf-adedcb612a0b")
  -diskUUID2 string
        disk_uuid_2 (default "9eabbb39-1baf-4872-beaf-adedcb612a0b")
  -endpoint string
        endpoint must be set
  -http-log string
        path to file where http request and response headers must be stored
  -insecure
        insecure flag must set true to allow provider to perform insecure SSL requests.
  -kind string
        kind (default "vm")
  -memorySize string
        memory_size_mb (default "1024")
  -name string
        name (default "vm_test1")
  -numSockets string
        num_sockets (default "1")
  -numVCPUs string
        num_vcpus (default "1")
  -password string
        password for api call
  -port string
        port for api call (default "9440")
  -powerState string
        power_state (default "POWERED_ON")
  -specVersion string
        spec_version
  -updateMemorySize string
        update_memory_size_name (default "2048")
  -updateName string
        update_name (default "vm_test2")
  -username string
        username for api call
```
