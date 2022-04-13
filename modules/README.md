# Terraform Nutanix Modules

Terraform nutanix modules are used to create various resources on nutanix. They are one step ahead automation around various data sources and resources of terraform nutanix provider.

## Steps to use 

Currently these modules can be used by following process :

Pull modules code from nutanix/terraform-provider-nutanix github repository

```sh
$ git clone https://github.com/nutanix/terraform-provider-nutanix.git
```

Create main.tf file (Doesn't need to be in same directory as repository)

Now, modules can be accessed by mentioning file path to the directory of particular module you want to use in `source` parameter. Modules are under `modules` folder in terraform-provider-nutanix repository.

For example, if wanted to use any foundation module

```hcl
module "image_nodes" {
    source = "../../files/terraform-provider-nutanix/modules/foundation/aos-based-node-imaging/node-serials-filter"
    .
    .
    .
    .
}
```

Example to use foundation central module
```hcl
module "image_nodes"{
    source = "../../files/terraform-provider-nutanix/modules/foundationCentral/aos-based-node-imaging/
    node-serial-filter"
    .
    .
    .
}
```

## Foundation Modules

The modules based on foundation data sources and resources are given below : 

| Name | Description | Required Version|
|------|-------------|-----------------|
| aos-based-node-imaging/nod-serial-filter| This module can be used to image nutanix imaged node having cvm running by giving node serials and some required information. Internally it uses node network details and discover nodes data sources to discover & get network details of those node serial corresponding nodes, which can be used as imaging input. | >1.4.1 |
|dos-based-node-imaging/nod-serial-filter| This module can be used to image nodes having discovery os running by giving node serials and some required information. Internally it uses node network details and discover nodes data sources to discover & get network details of those node serial corresponding nodes, which can be used as imaging input.| >1.4.1|
| manual-node-imaging | This module can be used to image nodes , which cannot be discovered, by providing defaults and node specific information. | >1.4.1 |
| discover-nodes-network-details | This module can be used to node discovery and get network details of nodes which are not part of cluster. This module is used by other foundatio modules as well | >1.4.1 |

Note : `Required Version` denotes required terraform nutanix provider version.

Check terraform-provider-nutanix/modules/foundation/examples for example configuration

## Foundation Central Modules

The modules  based on foundation central datasources and resources  are given below:

| Name | Description | Required Version |
|------|-------------|------------------|
|aos-based-node-imaging/node-serial-filter| This modules can be used to image nutanix nodes by giving node serial numbers and other required inputs. Internally it uses datasources to fetch the corresponding node information and later used as node imaging input. | >1.5.0-beta |
|manual-node-imaging/node-serial-filter| This modules is used to image nodes given required input information and node specific details. |>1.5.0-beta |  

Note : `Required Version` denotes required terraform nutanix provider version.

Check terraform-provider-nutanix/modules/foundationCentral/examples for example configuration