## Summary

Auto-generated Terraform provider code for the **monitoring** namespace, entity
**audit**, based on the inline `sdk_info.json` (see prompt). One PR per code-gen run.

- SDK module: `github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.2.2`
- API methods covered: 2 (`GetAuditById`, `ListAudits`)
- Datasources generated: 2 (`nutanix_audit_v2`, `nutanix_audits_v2`)
- Resources generated: 0 (N/A — audit entity is read-only)

## Files changed

### `nutanix/sdks/v4/monitoring/`
- `monitoring.go` — SDK client wrapper for monitoring API

### `nutanix/` (config & provider)
- `config.go` — Added `MonitoringAPI` client field and initialization
- `provider/provider.go` — Registered `nutanix_audit_v2` and `nutanix_audits_v2` datasources

### `nutanix/services/monitoringv2/`
- `data_source_nutanix_audit_v2.go` — Singular datasource (GetAuditById)
- `data_source_nutanix_audits_v2.go` — Plural datasource (ListAudits)
- `helpers.go` — Schema definitions and flatten functions for audit types
- `main_test.go` — Test harness
- `data_source_nutanix_audit_v2_test.go` — Acceptance tests for singular datasource
- `data_source_nutanix_audits_v2_test.go` — Acceptance tests for plural datasource

### `examples/monitoring_v2/`
- `main.tf` — Example usage for both datasources
- `variables.tf` — Variable definitions

### `website/docs/d/`
- `audit_v2.html.markdown` — Documentation for singular datasource
- `audits_v2.html.markdown` — Documentation for plural datasource

### `go.mod` / `go.sum`
- Added `monitoring-go-client/v4@v4.2.2` dependency

## Test execution

| Metric              | Value                                          |
|---------------------|------------------------------------------------|
| Command             | `TF_ACC=1 go test ./nutanix/services/monitoringv2 -v "-run=TestAccV2NutanixAudit" -timeout 500m -coverprofile c.out -covermode=count` |
| Total testcases     | 4                                              |
| Passed              | 0                                              |
| Failed              | 4 (see analysis below)                         |
| Skipped             | 0                                              |
| Wall-clock duration | 0:00:42                                        |
| Cluster (PC)        | `10.44.76.91`                                  |

### Analysis

All 4 test failures are due to **network connectivity** — the Nutanix cluster endpoint (`10.44.76.91:9440`) is unreachable from the CI environment. Every request gets `connection reset by peer`.

This was confirmed to be an infrastructure issue (not code-related) by running an existing test from `datapoliciesv2` which fails identically:
```
OPTIONS https://10.44.76.91:9440/api/prism/unversioned/info
read tcp 172.30.0.2:54912->10.44.76.91:9440: read: connection reset by peer
```

**Known issues:**
- `TestAccV2NutanixAuditDatasource_Basic` — Connection reset by peer to 10.44.76.91:9440
- `TestAccV2NutanixAuditsDatasource_Basic` — Connection reset by peer to 10.44.76.91:9440
- `TestAccV2NutanixAuditsDatasource_WithLimit` — Connection reset by peer to 10.44.76.91:9440
- `TestAccV2NutanixAuditsDatasource_WithFilter` — Connection reset by peer to 10.44.76.91:9440

The code compiles cleanly (`go build ./...` and `go vet ./...` pass). The API URLs and query parameters in the test logs confirm correct SDK integration. Tests should pass when run against a reachable Nutanix cluster with the monitoring API enabled.

### Test logs (tail)

<details>
<summary>tail -n 500 test_logs.log</summary>

```text
2026/04/28 17:52:28 Do some crazy stuff before tests!
=== RUN   TestAccV2NutanixAuditDatasource_Basic
2026-04-28 17:52:29.114Z INFO - OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info
2026-04-28 17:52:34.118Z ERROR - Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:58360->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:52:34.118Z ERROR - Could not fetch supported versions from server with error : Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:58360->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:52:34.118Z INFO - GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=1
2026-04-28 17:52:39.122Z ERROR - Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=1": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=1 giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=1": read tcp 172.30.0.2:58370->10.44.76.91:9440: read: connection reset by peer
    data_source_nutanix_audit_v2_test.go:14: Step 1/1 error: Error running pre-apply refresh: exit status 1
        
        Error: error while fetching audits: Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=1": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=1 giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=1": read tcp 172.30.0.2:58370->10.44.76.91:9440: read: connection reset by peer
        
          with data.nutanix_audits_v2.audits,
          on terraform_plugin_test.tf line 2, in data "nutanix_audits_v2" "audits":
           2: data "nutanix_audits_v2" "audits" {
        
--- FAIL: TestAccV2NutanixAuditDatasource_Basic (10.45s)
=== RUN   TestAccV2NutanixAuditsDatasource_Basic
2026-04-28 17:52:39.575Z INFO - OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info
2026-04-28 17:52:44.578Z ERROR - Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:58286->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:52:44.578Z ERROR - Could not fetch supported versions from server with error : Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:58286->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:52:44.578Z INFO - GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits
2026-04-28 17:52:49.580Z ERROR - Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits": read tcp 172.30.0.2:58294->10.44.76.91:9440: read: connection reset by peer
    data_source_nutanix_audits_v2_test.go:13: Step 1/1 error: Error running pre-apply refresh: exit status 1
        
        Error: error while fetching audits: Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits": read tcp 172.30.0.2:58294->10.44.76.91:9440: read: connection reset by peer
        
          with data.nutanix_audits_v2.test,
          on terraform_plugin_test.tf line 2, in data "nutanix_audits_v2" "test":
           2: data "nutanix_audits_v2" "test" {}
        
--- FAIL: TestAccV2NutanixAuditsDatasource_Basic (10.46s)
=== RUN   TestAccV2NutanixAuditsDatasource_WithLimit
2026-04-28 17:52:50.017Z INFO - OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info
2026-04-28 17:52:55.021Z ERROR - Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:53752->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:52:55.021Z ERROR - Could not fetch supported versions from server with error : Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:53752->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:52:55.021Z INFO - GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=5
2026-04-28 17:53:00.024Z ERROR - Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=5": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=5 giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=5": read tcp 172.30.0.2:53768->10.44.76.91:9440: read: connection reset by peer
    data_source_nutanix_audits_v2_test.go:28: Step 1/1 error: Error running pre-apply refresh: exit status 1
        
        Error: error while fetching audits: Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=5": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=5 giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24limit=5": read tcp 172.30.0.2:53768->10.44.76.91:9440: read: connection reset by peer
        
          with data.nutanix_audits_v2.test,
          on terraform_plugin_test.tf line 2, in data "nutanix_audits_v2" "test":
           2: data "nutanix_audits_v2" "test" {
        
--- FAIL: TestAccV2NutanixAuditsDatasource_WithLimit (10.44s)
=== RUN   TestAccV2NutanixAuditsDatasource_WithFilter
2026-04-28 17:53:00.476Z INFO - OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info
2026-04-28 17:53:05.479Z ERROR - Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:46206->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:53:05.479Z ERROR - Could not fetch supported versions from server with error : Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info giving up after 1 attempt(s): Options "https://10.44.76.91:9440/api/monitoring/unversioned/info": read tcp 172.30.0.2:46206->10.44.76.91:9440: read: connection reset by peer
2026-04-28 17:53:05.479Z INFO - GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24filter=serviceName+eq+%27Nutanix%27&%24limit=5
2026-04-28 17:53:10.482Z ERROR - Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24filter=serviceName+eq+%27Nutanix%27&%24limit=5": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24filter=serviceName+eq+%27Nutanix%27&%24limit=5 giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24filter=serviceName+eq+%27Nutanix%27&%24limit=5": read tcp 172.30.0.2:46218->10.44.76.91:9440: read: connection reset by peer
    data_source_nutanix_audits_v2_test.go:43: Step 1/1 error: Error running pre-apply refresh: exit status 1
        
        Error: error while fetching audits: Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24filter=serviceName+eq+%27Nutanix%27&%24limit=5": GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24filter=serviceName+eq+%27Nutanix%27&%24limit=5 giving up after 1 attempt(s): Get "https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/audits?%24filter=serviceName+eq+%27Nutanix%27&%24limit=5": read tcp 172.30.0.2:46218->10.44.76.91:9440: read: connection reset by peer
        
          with data.nutanix_audits_v2.test,
          on terraform_plugin_test.tf line 2, in data "nutanix_audits_v2" "test":
           2: data "nutanix_audits_v2" "test" {
        
--- FAIL: TestAccV2NutanixAuditsDatasource_WithFilter (10.46s)
FAIL
coverage: 14.9% of statements
FAIL	github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/monitoringv2	41.853s
FAIL
```

</details>

## How this was generated

- Triggered via the Nutanix headless code-gen API (`POST /generate`).
- Prompt + inline `sdk_info.json` are stored only in the agent transcript;
  no `sdk_info.json` is committed to this branch.
