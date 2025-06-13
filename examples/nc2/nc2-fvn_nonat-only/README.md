# Non-NAT Only Example

This example demonstrates how to set up a non-NAT networking environment in Nutanix Cloud Infrastructure (NCI) using Terraform.

## Overview

This configuration creates:
- a no-NAT external overlay subnet in the default transit-vpc

## Network Configuration

- No-NAT addition for the transit-vpc

## Prerequisites

- Existing transit VPC named "transit-vpc" (default when Flow is deployed with the cluster)
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


## Variables

Key variables in `terraform.tfvars`:
- VPC name
- Subnet configurations
- Network settings
- Routing parameters

## Cleanup

To destroy the created resources:
```bash
terraform destroy
``` 