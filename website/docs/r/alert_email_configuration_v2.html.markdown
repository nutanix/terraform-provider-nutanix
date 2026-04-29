---
layout: "nutanix"
page_title: "NUTANIX: nutanix_alert_email_configuration_v2"
sidebar_current: "docs-nutanix-resource-alert-email-configuration-v2"
description: |-
  Updates the configuration that is used to send alert emails.
---

# nutanix_alert_email_configuration_v2

Updates the configuration that is used to send alert emails. This is a singleton resource — the alert email configuration always exists on the cluster. Creating this resource applies an update; deleting it merely removes it from state.

## Example Usage

```hcl
resource "nutanix_alert_email_configuration_v2" "example" {
  is_enabled              = true
  is_email_digest_enabled = true
  email_contact_list      = ["admin@example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `is_enabled`: - (Optional) Indicates whether alert emails are enabled or not.
* `has_default_nutanix_email`: - (Optional) Indicates whether alert emails are enabled or not on default Nutanix email ID.
* `default_nutanix_email`: - (Optional) The default Nutanix email ID to which alert emails are sent.
* `is_email_digest_enabled`: - (Optional) Indicates whether alert email digest is enabled or not.
* `is_empty_alert_email_digest_skipped`: - (Optional) Send alert email digest only if there are one or more alerts.
* `alert_email_digest_send_time`: - (Optional) Time in HH:mm format when the alert email digest is sent daily.
* `alert_email_digest_send_timezone`: - (Optional) Timezone for the time at which the alert email digest is sent daily.
* `email_contact_list`: - (Optional) List of email contacts.
* `email_config_rules`: - (Optional) Rules for email configuration.
* `email_template`: - (Optional) Email template configuration.

### email_config_rules
* `cluster_uuids`: (Optional) Cluster UUIDs to which this rule applies.
* `has_global_email_contact_list`: (Optional) Indicates whether to include a global email contact list.
* `impact_types`: (Optional) Impact types for the rule.
* `is_enabled`: (Optional) Indicates whether the configuration rule is enabled or not.
* `match_phrases`: (Optional) List of phrases to match the alert.
* `recipients`: (Optional) List of recipients who will receive emails.
* `severities`: (Optional) Severity levels for the rule.

### email_template
* `body_suffix`: (Optional) Suffix for the email body.
* `subject_prefix`: (Optional) Prefix for the email subject.

## Attribute Reference

The following attributes are exported:

* `ext_id`: A globally unique identifier of an instance that is suitable for external consumption.
* `tenant_id`: A globally unique identifier that represents the tenant that owns this entity.
* `links`: A HATEOAS style link for the response.
* `tunnel_details`: Details of the remote tunnel configuration.

### tunnel_details
* `connection_status`: Connection status details.
* `http_proxy`: HTTP proxy configuration.
* `service_center`: Service center details.
* `transport_status`: Transport status details.

### links
* `href`: The URL at which the entity described by the link can be accessed.
* `rel`: A name that identifies the relationship of the link to the object that is returned by the URL.

See detailed information in [Nutanix Monitoring v4 Alert Email Configuration](https://developers.nutanix.com/api-reference?namespace=monitoring&version=v4.0).
