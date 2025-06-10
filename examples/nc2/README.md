# Nutanix Cloud Clusters (NC2) Flow Virtual Networking (FVN) Terraform Samples

This folder contains a collection of Terraform samples demonstrating different Flow Networking configurations in Nutanix Cloud Clusters (NC2). These samples show how to set up various networking scenarios using Terraform. Also included are VM image downloads from public repositories and VM creation. 

## Prerequisites

- Nutanix Cloud Clusters (NC2) on AWS cluster
- Prism Central instance
- Terraform installed (version 1.0.0 or later)
- Nutanix Terraform provider (version 2.2.0 or later)



## Nutanix software versions
These samples have been tested and verified on NC2 on AWS with the following Nutanix software versions. 
- Prism Central: pc.2024.3.1.1
- AOS: 10.0.1
- AHV: 7.0.1.5


## Available Samples

**nc2-fvn_nonat-only**
* Creates the "overlay-external-subnet-nonat" subnet in the Flow "transit-vpc"

**nc2-fvn_nat-vpcs**:
* Deploys two Flow VPCs with NAT egress and two subnets each

**nc2-os-images**
* Populates the Prism Central image library with a number of Ubuntu and CentOS images from public repositories

**nc2-fvn_nat-vms-simple**
* Deploys one VPC with NAT egress and two subnets
* Downloads and OS image
* Deploys VMs initiated with cloud-init scripts from the OS image (customize the cloud-init YAML in the "templates" subfolder as desired)
* Creates and associates Floating IPs from the AWS VPC CIDR range to each VM

**nc2-fvn_nonat-vpcs**
* Creates the "overlay-external-subnet-nonat" subnet in the Flow "transit-vpc"
* Creates two Flow VPCs with no-NAT egress
* Creates subnets in each Flow VPC
* Adds the subnets CIDR ranges as Externally Routable Prefixes (ERP) to the Flow VPCs
* Modify the "terraform.tfvars" file to customize the VPCs and subnets

**nc2-fvn_nat-nonat-vms**
* Creates the "overlay-external-subnet-nonat" subnet in the Flow "transit-vpc"
* Creates two Flow VPCs: One with NAT and another with no-NAT egress
* Creates subnets in each Flow VPC
* Downloads an Ubuntu image
* Deploys VMs from the Ubuntu image, including cloud-init scripts (customize them in the "templates" folder)
* Creates and associates Floating IPs from the AWS VPC CIDR range to each VM on the NAT VPC
* Adds the subnets CIDR ranges as Externally Routable Prefixes (ERP) to the Flow no-NAT VPCs


## Usage

Each sample directory contains:
- `main.tf`: Main Terraform configuration
- `variables.tf`: Variable definitions
- `terraform.tfvars`: Variable values
- `set-env.sh`: Environment setup script

To use any sample:

1. Navigate to the desired sample directory
2. Source the environment variables:
   ```bash
   source set-env.sh
   ```
3. Initialize Terraform:
   ```bash
   terraform init
   ```
4. Review the planned changes:
   ```bash
   terraform plan
   ```
5. Apply the configuration:
   ```bash
   terraform apply
   ```

## Environment Variables

Each sample requires the following environment variables (set via `set-env.sh`):
- `TF_VAR_NUTANIX_USERNAME`: Prism Central username
- `TF_VAR_NUTANIX_PASSWORD`: Prism Central password
- `TF_VAR_NUTANIX_ENDPOINT`: Prism Central IP/FQDN
- `TF_VAR_NUTANIX_PORT`: Prism Central port (default: 9440)
- `TF_VAR_NUTANIX_INSECURE`: Set to "true" if using self-signed certificates
- `TF_VAR_SSH_PUBLIC_KEY`: SSH public key for VM access

