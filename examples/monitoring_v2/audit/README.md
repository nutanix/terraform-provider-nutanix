# Nutanix Monitoring v2 - Audit Data Sources Example

This example demonstrates how to use the Nutanix Monitoring v2 audit data sources to fetch audit information.

## Features

- Fetch all audits using `nutanix_audits_v2` data source
- Fetch a specific audit by ext_id using `nutanix_audit_v2` data source
- Display audit details including type, service name, operation type, and status

## Usage

1. Update the `variables.tf` file with your Nutanix Prism Central credentials:

```hcl
nutanix_username = "admin"
nutanix_password = "your-password"
nutanix_endpoint = "your-prism-central-ip"
nutanix_port     = "9440"
nutanix_insecure = true
```

2. Initialize Terraform:

```bash
terraform init
```

3. Review the plan:

```bash
terraform plan
```

4. Apply the configuration:

```bash
terraform apply
```

## Data Sources Used

- `nutanix_audits_v2` - Fetches a list of all audits
- `nutanix_audit_v2` - Fetches details of a specific audit by ext_id

## Outputs

- `first_audit_ext_id` - The ext_id of the first audit (if available)
- `audit_details` - Complete details of the fetched audit including:
  - ext_id
  - audit_type
  - service_name
  - creation_time
  - operation_type
  - status
  - message
