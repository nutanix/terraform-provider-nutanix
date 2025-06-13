# OS Images Management Example

This example demonstrates how to manage OS images in Nutanix Cloud Infrastructure (NCI) using Terraform.

## Overview

This configuration:
- Manages OS images in Prism Central
- Handles image uploads and configurations
- Sets up image properties and metadata
- Manages image lifecycle

## Prerequisites

- Prism Central access
- OS image files available
- Sufficient storage space in Prism Central
- Understanding of image requirements

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
- Manages OS images in Prism Central
- Configures image properties
- Sets up image metadata
- Handles image lifecycle
- Manages image versions

## Variables

Key variables in `terraform.tfvars`:
- Image names
- Image source locations
- Image properties
- Version information

## Cleanup

To destroy the created resources:
```bash
terraform destroy
```

Note: Be careful when destroying images as this will remove them from Prism Central. 