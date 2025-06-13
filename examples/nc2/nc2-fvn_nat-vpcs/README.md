# NAT-Enabled VPCs Example

This example demonstrates how to create NAT-enabled VPCs with multiple subnets in Nutanix Cloud Infrastructure (NCI) using Terraform.

## Overview

This configuration creates:
- Two NAT-enabled VPCs
- Two subnets per VPC (Subnet-A and Subnet-B)
- Integration with an existing transit VPC
- Connection to an external NAT subnet

## Network Configuration

- VPC 1:
  - Subnet-A: 192.168.10.0/24
  - Subnet-B: 192.168.30.0/24
- VPC 2:
  - Subnet-A: 192.168.20.0/24
  - Subnet-B: 192.168.40.0/24

## Prerequisites

- Existing transit VPC named "transit-vpc"
- External NAT subnet named "overlay-external-subnet-nat"
- Prism Central access

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

## Configuration Details

The configuration creates:
- Two VPCs with NAT capabilities
- Each VPC has two subnets (A and B)
- Subnets are configured with:
  - DHCP enabled
  - DNS server configuration
  - IP address pools for VM deployment
  - Default gateway configuration

## Variables

Key variables in `terraform.tfvars`:
- VPC names
- Subnet names
- Network configurations

## Cleanup

To destroy the created resources:
```bash
terraform destroy
``` 