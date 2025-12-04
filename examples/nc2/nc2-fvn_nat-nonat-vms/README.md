# Mixed NAT/Non-NAT VMs Example

This example demonstrates how to deploy virtual machines in a mixed NAT and no-NAT environment using Terraform in Nutanix Cloud Clusters (NC2).

## Overview

This configuration:
- Creates VMs in both NAT and no-NAT VPCs
- Configures networking for different VM types
- Sets up cloud-init for VM initialization

## Prerequisites

- Existing NAT and non-NAT VPCs
- OS image available in Prism Central
- Prism Central access
- SSH public key for VM access
- Understanding of mixed networking requirements

## Usage

1. Set up environment variables:
   ```bash
   source set-env.sh
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Review the planned changes:
   ```bash
   terraform plan
   ```

4. Apply the configuration:
   ```bash
   terraform apply
   ```

## NOTE 

The CIDR ranges of the no-NAT subnets must be added as Externally Routable Prefixes (ERPs) to both each no-NAT Flow VPC and to the transit-vpc. This Terraform sample will add the ERP entries to the Flow VPCs created by the sript but not to the "transit-vpc", as this will also automatically update the AWS VPC route table. As such, after the Terraform deployment is done, add the CIDR ranges of the subnets as ERPs to the "transit-vpc" and the VMs will be able to communicate with the outside world (as well as being accessible from EC2 instances in the VPC NC2 is deployed into). 

For more info on NAT and no-NAT networking on NC2, including how to update the ERP entries, please refer to: https://jonamiki.com/2024/09/09/how-do-set-up-nat-and-no-nat-networking-with-nc2-on-aws/


## Configuration Details

The configuration:
- Deploys VMs in both NAT and no-NAT environments
- Uses cloud-init for VM initialization
- Configures networking and security for each environment
- Handles mixed networking scenarios

## Templates

The `templates` directory contains:
- Cloud-init configuration templates
- VM initialization scripts
- Network configuration templates for both NAT and no-NAT environments

## Variables

Key variables in `terraform.tfvars`:
- VM names and configurations
- Network settings for both environments
- Resource specifications
- Environment-specific parameters

## Cleanup

To destroy the created resources:
```bash
terraform destroy
``` 
