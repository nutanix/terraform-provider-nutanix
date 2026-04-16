# NAT-Enabled VMs Simple Example

This example demonstrates how to deploy virtual machines in a NAT-enabled VPC environment using Terraform in Nutanix Cloud Infrastructure (NCI).

## Overview

This configuration:
- Creates VMs in a NAT-enabled VPC
- Configures networking for the VMs
- Sets up cloud-init for VM initialization
- Manages VM lifecycle

## Prerequisites

- Existing NAT-enabled VPC
- OS image available in Prism Central
- Prism Central access
- SSH public key for VM access

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

The configuration:
- Deploys VMs in a NAT-enabled environment
- Uses cloud-init for VM initialization
- Configures networking and security
- Manages VM resources and lifecycle

## Templates

The `templates` directory contains:
- Cloud-init configuration templates
- VM initialization scripts
- Network configuration templates

## Variables

Key variables in `terraform.tfvars`:
- VM names and configurations
- Network settings
- Resource specifications

## Cleanup

To destroy the created resources:
```bash
terraform destroy
``` 