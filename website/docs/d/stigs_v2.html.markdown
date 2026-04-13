---
layout: "nutanix"
page_title: "NUTANIX: nutanix_stigs_v2"
sidebar_current: "docs-nutanix-datasource-stigs-v2"
description: |-
   Provides a Nutanix STIGS datasource to Get the STIG controls details. A Security Technical Implementation Guide (STIG) is a cybersecurity methodology for standardizing security protocols within networks, servers, computers, and logical designs to enhance overall security. These guides, when implemented, enhance security for software, hardware, and physical and logical architectures to further reduce vulnerabilities. This datasource retrieves Security Technical Implementation Guide (STIG) control details for each cluster.
---

# nutanix_stigs_v2

A Security Technical Implementation Guide (STIG) is a cybersecurity methodology for standardizing security protocols within networks, servers, computers, and logical designs to enhance overall security. These guides, when implemented, enhance security for software, hardware, and physical and logical architectures to further reduce vulnerabilities.

This datasource retrieves Security Technical Implementation Guide (STIG) control details for each cluster. Each STIG record represents a specific rule or control evaluated against one or more clusters, containing metadata such as rule ID, severity, compliance status, and remediation guidance.

This datasource uses Prism Central (PC) v4 APIs based SDKs.

## Example

```hcl

data "nutanix_stigs_v2" "all" {}

data "nutanix_stigs_v2" "filtered-status"{
  filter = "status eq Security.Report.StigStatus'APPLICABLE'"
}

data "nutanix_stigs_v2" "filtered-severity"{
  filter = "severity eq Security.Report.Severity'HIGH'"
}

data "nutanix_stigs_v2" "limited"{
  limit = 2
}

data "nutanix_stigs_v2" "select-example"{
  select = "stigVersion,status"
}
```

## Argument Reference

The following arguments are supported:

- `page`: - A URL query parameter that specifies the page number of the result set. It must be a positive integer between 0 and the maximum number of pages that are available for that resource. Any number out of this range might lead to no results.
- `limit` : A URL query parameter that specifies the total number of records returned in the result set. Must be a positive integer between 1 and 100. Any number out of this range will lead to a validation error. If the limit is not provided, a default value of 50 records will be returned in the result set.
- `filter` :A URL query parameter that allows clients to filter a collection of resources. The expression specified with \$filter is evaluated for each resource in the collection, and only items where the expression evaluates to true are included in the response. Expression specified with the \$filter must conform to the OData V4.01 URL conventions. For example, filter '\$filter=name eq 'karbon-ntnx-1.0' would filter the result on cluster name 'karbon-ntnx1.0', filter '\$filter=startswith(name, 'C')' would filter on cluster name starting with 'C'. The filter can be applied to the following fields:

  - <details>
    <summary>severity</summary>

    **Example:**

    ```
    filter = "severity eq Security.Report.Severity'HIGH'"
    ```

    </details>

  - <details>
    <summary>status</summary>

    **Example:**

    ```
    filter = "status eq Security.Report.StigStatus'APPLICABLE'"
    ```

    </details>

* `orderby` : A URL query parameter that allows clients to specify the sort criteria for the returned list of objects. Resources can be sorted in ascending order using asc or descending order using desc. If asc or desc are not specified, the resources will be sorted in ascending order by default. For example, '\$orderby=templateName desc' would get all templates sorted by templateName in descending order. The orderby can be applied to the following fields:

  - <details>
    <summary>severity</summary>

    **Example:**

    ```
    filter = "severity"
    ```

    </details>

  - <details>
    <summary>stigVersion</summary>

    **Example:**

    ```
    order_by = "stigVersion desc"
    ```

    </details>

* `select` : A URL query parameter that allows clients to request a specific set of properties for each entity or complex type. Expression specified with the \$select must conform to the OData V4.01 URL conventions. If a \$select expression consists of a single select item that is an asterisk (i.e., \*), then all properties on the matching resource will be returned. following fields:

  - <details>
      <summary>affectedClusters</summary>

    **Example:**

    ```
    select = "affectedClusters"
    ```

    </details>

  - <details>
    <summary>benchmarkId</summary>

    **Example:**

    ```
    select = "benchmarkId"
    ```

    </details>

  - <details>
    <summary>comments</summary>

    **Example:**

    ```
    select = "comments"
    ```

    </details>

  - <details>
    <summary>extId</summary>

    **Example:**

    ```
    select = "extId"
    ```

    </details>

  - <details>
    <summary>fixText</summary>

    **Example:**

    ```
    select = "fixText"
    ```

    </details>

  - <details>
    <summary>identifiers</summary>

    **Example:**

    ```
    select = "identifiers"
    ```

    </details>

  - <details>
    <summary>links</summary>

    **Example:**

    ```
    select = "links"
    ```

    </details>

  - <details>
    <summary>ruleId</summary>

    **Example:**

    ```
    select = "ruleId"
    ```

    </details>

  - <details>
    <summary>severity</summary>

    **Example:**

    ```
    select = "severity"
    ```

    </details>

  - <details>
    <summary>status</summary>

    **Example:**

    ```
    select = "status"
    ```

    </details>

  - <details>
    <summary>stigVersion</summary>

    **Example:**

    ```
    select = "stigVersion"
    ```

    </details>

  - <details>
    <summary>tenantId</summary>

    **Example:**

    ```
    select = "tenantId"
    ```

    </details>

  - <details>
    <summary>title</summary>

    **Example:**

    ```
    select = "title"
    ```

    </details>

## Attribute Reference

The Following attributes are exported:

- `stigs`: -List of STIG controls details for STIG rules on each cluster.

### stigs

The `stigs` exports the following:

- `tenantId`: - A globally unique identifier that represents the tenant that owns this entity. The system automatically assigns it, and it and is immutable from an API consumer perspective (some use cases may cause this ID to change - For instance, a use case may require the transfer of ownership of the entity, but these cases are handled automatically on the server).
- `ext_id`: - A globally unique identifier of an instance that is suitable for external consumption.
- `links`: - A HATEOAS style link for the response. Each link contains a user-friendly name identifying the link and an address for retrieving the particular resource.
- `title`: - Title of the STIG control.
- `rule_id`: - Rule ID of the STIG control.
- `stig_version`: - STIG ID of the control.
- `identifiers`: - Additional identifiers used to describe this control.
- `affected_clusters`: - List of clusters that failed the STIG control.
- `severity`: - Contains possible values for the severity level of a vulnerability.
  | Enum | Description |
  |--------------|--------------------|
  | `HIGH` | Severity level high. |
  | `MEDIUM` | Severity level medium. |
  | `LOW` | Severity level low. |
  | `UNKNOWN` | Unknown value. |
  | `CRITICAL` | Severity level critical. |
  | `REDACTED` | Redacted value. |

- `status`: - Current status of the STIG rule.
  | Enum | Description |
  |--------------|--------------------|
  | `NOT_APPLICABLE` | STIG is not applicable. |
  | `NEEDS_REVIEW` | STIG application needs a review. |
  | `APPLICABLE` | STIG is applicable. |
  | `UNKNOWN` | Unknown value. |
  | `REDACTED` | Redacted value. |
- `comments`: - The comments to explain why a STIG rule applies or does not apply to the cluster.
- `fix_text`: - The command/steps to fix the STIG rule failure.
- `benchmark_id`: - Benchmark ID of the STIG rules.

See detailed information in [Nutanix Get the STIG controls details V4](https://developers.nutanix.com/api-reference?namespace=security&version=v4.0#tag/STIGs/operation/listStigs)
