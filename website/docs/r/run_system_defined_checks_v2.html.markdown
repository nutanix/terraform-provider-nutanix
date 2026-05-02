---
layout: "nutanix"
page_title: "NUTANIX: nutanix_run_system_defined_checks_v2"
sidebar_current: "docs-nutanix-resource-run-system-defined-checks-v2"
description: |-
  Run System-Defined Checks on a cluster.
---

# nutanix_run_system_defined_checks_v2

Run System-Defined Checks on a cluster. This is an action resource that triggers a run of System-Defined Checks. It creates a task and waits for completion.

## Example

```hcl
resource "nutanix_run_system_defined_checks_v2" "example" {
  cluster_ext_id                              = "00000000-0000-0000-0000-000000000000"
  should_run_all_checks                       = true
  should_send_report_to_configured_recipients = false
}
```

## Argument Reference

The following arguments are supported:

- `cluster_ext_id`: (Required) Unique identifier for the cluster for which run System-Defined Checks is requested.
- `additional_recipients`: (Optional) A list of additional email addresses for sending the run summary. Either this should be set or should_send_report_to_configured_recipients should be true. If both are set then email would be sent to all the recipients.
- `node_ips`: (Optional) List of node IP addresses where the Check will run. This field will be ignored if the check scope is a cluster.
- `sda_ext_ids`: (Optional) List of Check IDs to be executed. This field cannot be set simultaneously with should_run_all_checks; only one of them should be specified.
- `should_anonymize`: (Optional) Indicates whether to mask sensitive data in the check run summary.
- `should_run_all_checks`: (Optional) Indicates whether all System-Defined Checks applicable to the specified cluster should be executed. This field is mutually exclusive with the sda_ext_ids parameter.
- `should_send_report_to_configured_recipients`: (Optional) Determines if the run summary should be sent to the configured email address associated with the cluster.

### node_ips

- `prefix_length`: (Optional) The prefix length of the network to which this host IPv4 address belongs.
- `value`: (Optional) The IPv4 address of the host.

## Attribute Reference

The following attributes are exported:

- `task_ext_id`: A globally unique identifier for the task created by this action.
