---
layout: "nutanix"
page_title: "NUTANIX: nutanix_ndb_onboarding"
sidebar_current: "docs-nutanix-resource-ndb-onboarding"
description: |-
  Runs the Nutanix Database Service (NDB) onboarding workflow for a Prism Element cluster.
---

# nutanix_ndb_onboarding

Provides a workflow-style resource to onboard a Prism Element cluster into Nutanix Database Service (NDB). The resource supports optional Prism Central validation, NDB DNS/NTP configuration, storage selection, network selection, and setup execution.

## Example Usage

```hcl
resource "nutanix_ndb_onboarding" "wizard" {
  enable_full_onboarding = true
  selection_mode         = "auto"

  prism_central_info {
    name       = "pc-10-102-144-166"
    ip_address = "10.102.144.166"
    username   = var.pc_username
    password   = var.pc_password
  }

  prism_element_info {
    name       = var.pe_name
    cluster_ip = var.pe_cluster_ip
    username   = var.pe_username
    password   = var.pe_password
  }

  ndb_config {
    dns_servers = ["10.22.64.16", "10.40.64.16"]
    ntp_servers = ["pool.ntp.org"]
    timezone    = "UTC"
  }

  storage {
    container_name = "SelfServiceContainer"
  }

  network_details {
    skip                  = false
    existing_network_name = "vlan.0"
  }

  setup {
    trigger         = true
    timeout_minutes = 120
  }
}
```

## Argument Reference

* `enable_full_onboarding`: (Optional) When `true`, runs the full onboarding workflow after Prism Element registration. Defaults to `false`.
* `selection_mode`: (Optional) Selection behavior for optional discovered values. Supported values are `auto` and `strict`. Defaults to `auto`.
* `prism_central_info`: (Optional) Prism Central validation details.
* `prism_element_info`: (Required) Prism Element details to onboard.
* `ndb_config`: (Optional) NDB DNS, NTP, timezone, and SMTP configuration.
* `storage`: (Optional) Storage container selection.
* `network_details`: (Optional) Network selection for setup.
* `setup`: (Optional) Setup trigger and timeout settings.

### prism_central_info

* `name`: (Optional) Prism Central display name.
* `description`: (Optional) Prism Central description.
* `ip_address`: (Required) Prism Central IP address.
* `port`: (Optional) Prism Central port. Defaults to `9440`.
* `username`: (Required) Prism Central username. This field is sensitive.
* `password`: (Required) Prism Central password. This field is sensitive.

### prism_element_info

* `name`: (Required) Prism Element cluster name.
* `description`: (Optional) Prism Element description.
* `cluster_ip`: (Required) Prism Element cluster IP address.
* `username`: (Required) Prism Element username. This field is sensitive.
* `password`: (Required) Prism Element password. This field is sensitive.
* `version`: (Optional) Prism Element API version. Defaults to `v2`.
* `cloud_type`: (Optional) Cluster cloud type. Defaults to `NTNX`.

### ndb_config

* `dns_servers`: (Optional) List of DNS servers.
* `ntp_servers`: (Optional) List of NTP servers.
* `timezone`: (Optional) NDB timezone. Defaults to `UTC`.
* `smtp_server_ip_port`: (Optional) SMTP server and port.
* `smtp_username`: (Optional) SMTP username.
* `smtp_password`: (Optional) SMTP password. This field is sensitive.
* `email_from_address`: (Optional) SMTP sender email address.
* `smtp_tls_enabled`: (Optional) Whether SMTP TLS is enabled.
* `smtp_unsecured`: (Optional) Whether unsecured SMTP is enabled.
* `apply_smtp_even_empty`: (Optional) Apply SMTP settings even when empty. Defaults to `true`.

### storage

* `container_name`: (Optional) Storage container name. In `auto` mode, the provider selects the first discovered container when this is empty.

### network_details

* `skip`: (Optional) Skip custom network selection. Defaults to `true`.
* `existing_network_name`: (Optional) Existing network name to select.
* `vlan_name`: (Optional) VLAN name.
* `static_ip`: (Optional) Static IP address.
* `gateway`: (Optional) Network gateway.
* `subnet_mask`: (Optional) Network subnet mask.

### setup

* `trigger`: (Optional) Trigger NDB setup after earlier onboarding steps. Defaults to `true`.
* `timeout_minutes`: (Optional) Timeout for setup operation polling. Defaults to `90`.

## Attributes Reference

The following attributes are exported:

* `id`: Resource ID. This is the onboarded NDB cluster ID.
* `cluster_id`: Onboarded NDB cluster ID.
* `operation_id`: Prism Element registration operation ID.
* `current_step`: Current or final onboarding step.
* `completed_steps`: List of completed onboarding steps.
* `status`: Onboarding status.
* `available_storage_containers`: Discovered storage containers.
* `available_dns_servers`: Discovered DNS servers.
* `available_ntp_servers`: Discovered NTP servers.
* `available_network_names`: Discovered network names.
* `effective_storage_container`: Storage container applied by the workflow.
* `effective_dns_servers`: DNS servers applied by the workflow.
* `effective_ntp_servers`: NTP servers applied by the workflow.
* `effective_network_name`: Network selected by the workflow.
* `setup_operation_id`: Setup operation ID.
* `setup_progress_percent`: Setup progress percentage.
* `setup_current_step`: Current setup step.

## Workflow Behavior

This is a workflow-style resource. `Create` executes onboarding, `Read` refreshes the recorded state, `Update` is intentionally unsupported, and `Delete` only removes the resource from Terraform state without deregistering or tearing down the remote cluster.
