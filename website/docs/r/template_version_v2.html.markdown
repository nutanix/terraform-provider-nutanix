---
layout: "nutanix"
page_title: "NUTANIX: nutanix_template_version_v4"
sidebar_current: "docs-nutanix-resource-template-version-v4"
description: |-
  Create a Template from the given VM identifier. A Template stores the VM configuration and disks from the source VM.
---

# nutanix_template_version_v4

Create a Template from the given VM identifier. A Template stores the VM configuration and disks from the source VM.

## Example 

```hcl
    resource "nutanix_template_version_v4" "test" {
        template_ext_id = resource.nutanix_template_v2.test.id
        template_version_spec{
            version_name = "Second temp"
            version_description = "second desc"
            version_source{
                template_version_reference{
                    version_id= nutanix_template_v2.test.template_version_spec.0.ext_id
                    override_vm_config{
                        num_sockets=1
                        num_threads_per_core=2
                        memory_size_bytes= 1073741824
                        num_cores_per_socket = 1    
                    }
                }
            }
            is_active_version = true
        }
    }
```


## Argument Reference

The following arguments are supported:

* `template_ext_id`: (Required) The identifier of a Template    .
* `template_version_spec`: (Required) A model that represents an object instance that is accessible through an API endpoint. Instances of this type get an extId field that contains the globally unique identifier for that instance. Externally accessible instances are always tenant aware and, therefore, extend the TenantAwareModel

### template_version_spec

* `version_source`: (Required) Source of the created Template Version. The source can either be a VM when creating a new Template Version or an existing Version within a Template when creating a new Version. Either `template_vm_reference` or `template_version_reference` . 


* `version_source.template_version_reference`: (Optional) Template Version Reference


### version_source.template_version_reference

* `version_id`: (Required) The identifier of a Template Version.
* `override_vm_config`: (Required) Overrides specification for VM create from a Template.


### override_vm_config

* `name`: (Optional) VM name.
* `num_sockets`: (Optional) Number of vCPU sockets.
* `num_cores_per_socket`: (Optional) Number of cores per socket.
* `num_threads_per_core`: (Optional) Number of threads per core.
* `memory_size_bytes`: (Optional) Memory size in bytes.
* `nics`: (Optional) NICs attached to the VM.
* `guest_customization`: (Optional) Stage a Sysprep or cloud-init configuration file to be used by the guest for the next boot. Note that the Sysprep command must be used to generalize the Windows VMs before triggering this API call.


See detailed information in [Nutanix Template](https://developers.nutanix.com/api-reference?namespace=vmm&version=v4.0.b1).