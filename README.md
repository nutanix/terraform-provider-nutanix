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
$ go test -v $(glide novendor) --username=username --password=password --endpoint=1.1.1.1 --insecure=true
```
Following flags are defined for the testcases :- 

```sh
 -api-version string
        This is api_version for testcase vm. (default "3.0")
  -disk-device-type-1 string
        This is device type for the first disk. (default "DISK")
  -disk-device-type-2 string
        This is device type for the second disk. (default "DISK")
  -disk-kind-1 string
        This is Kind field for the first disk. (default "image")
  -disk-kind-2 string
        This is Kind field for the second disk. (default "image")
  -disk-name-1 string
        This is disk name of first disk. (default "Centos7")
  -disk-name-2 string
        This is disk name of second disk. (default "Centos7")
  -disk-size-1 string
        This is size of the first disk (default "1")
  -disk-size-2 string
        This is size of the second disk (default "1")
  -disk-uuid-1 string
        This is UUID of first disk. (default "9eabbb39-1baf-4872-beaf-adedcb612a0b")
  -disk-uuid-2 string
        This is UUID of second disk. (default "9eabbb39-1baf-4872-beaf-adedcb612a0b")
  -diskNo string
        This is the number of disks attached to the disktestcase vm. (default "2")
  -endpoint string
        endpoint must be set
  -http-log string
        path to file where http request and response headers must be stored
  -insecure
        insecure flag must set true to allow provider to perform insecure SSL requests.
  -kind string
        This is kind for testcase vm. (default "vm")
  -memory-size string
        This is the memory_size_mb for testcase vm. (default "1024")
  -name string
        This is the name of the vm. (default "vm_test1")
  -num-sockets string
        This is num_sockets for the testcase vm. (default "1")
  -num-vcpus string
        This is num_vcpus for the testcase vm. (default "1")
  -password string
        password for api call
  -port string
        port for api call (default "9440")
  -power-state string
        This is power_state for testcase vm. (default "POWERED_ON")
  -spec-version string
        This is spec_version for testcase vm.
  -update-memory-size string
        This is the memory size to which vm gets upgraded in updateMemory testcase. (default "2048")
  -update-name string
        This is the updated name of the vm in updateName testcase. (default "vm_test2")
  -username string
        username for api call.
```

Instead of flags values can also be passed through the [Environment Variables](https://en.wikipedia.org/wiki/Environment_variable). A flag with name "x-y" can be set in CLI as <binary> --x-y, if the same flag has to be set in ENV, it has to be set as X_Y. For example a flag abc-xyz can be omitted by setting environment variable ABC_XYZ.
conflicts & resolution order in the descending order of precedence
    flag
    env
