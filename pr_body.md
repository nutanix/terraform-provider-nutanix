## Summary

Auto-generated Terraform provider code for the **monitoring** namespace, entity
**event**, based on the inline `sdk_info.json` (see prompt). One PR per code-gen run.

- SDK module: `github.com/nutanix/ntnx-api-golang-clients/monitoring-go-client/v4@v4.2.2`
- API methods covered: 2 (GetEventById, ListEvents)
- Datasources generated: 2 (`nutanix_event_v2`, `nutanix_events_v2`)
- Resources generated: 0 (no Create/Update/Delete methods available for events)

## Files changed

### `nutanix/sdks/v4/monitoring/`
- `monitoring.go` — SDK client wrapper for monitoring EventsApi

### `nutanix/services/monitoringv2/`
- `data_source_nutanix_event_v2.go` — Singular datasource (GetEventById)
- `data_source_nutanix_events_v2.go` — Plural datasource (ListEvents)
- `helpers.go` — Schema helpers and flatten functions
- `data_source_nutanix_event_v2_test.go` — Acceptance test for singular datasource
- `data_source_nutanix_events_v2_test.go` — Acceptance tests for plural datasource

### `examples/monitoring_v2/`
- `event/main.tf`, `event/variables.tf` — Example for singular event datasource
- `events/main.tf`, `events/variables.tf` — Example for plural events datasource

### `website/docs/d/`
- `event_v2.html.markdown` — Documentation for `nutanix_event_v2`
- `events_v2.html.markdown` — Documentation for `nutanix_events_v2`

### Modified files
- `go.mod`, `go.sum` — Added `monitoring-go-client/v4@v4.2.2` dependency
- `nutanix/config.go` — Added `MonitoringAPI` client field and initialization
- `nutanix/provider/provider.go` — Registered `nutanix_event_v2` and `nutanix_events_v2` datasources

## Test execution

| Metric              | Value                                                                                       |
|---------------------|---------------------------------------------------------------------------------------------|
| Command             | `TF_ACC=1 go test ./nutanix/services/monitoringv2 -v -run="TestAccV2Nutanix" -timeout 500m` |
| Total testcases     | 3                                                                                           |
| Passed              | 3                                                                                           |
| Failed              | 0                                                                                           |
| Skipped             | 0                                                                                           |
| Wall-clock duration | 1:44                                                                                        |
| Cluster (PC)        | 10.44.76.91                                                                                 |

### Analysis

All 3 acceptance tests passed on the first run:
- `TestAccV2NutanixEventDataSource_Basic` — 37.68s — Fetches an event by ID (obtained via ListEvents with limit=1)
- `TestAccV2NutanixEventsDataSource_Basic` — 32.57s — Lists all events without filters
- `TestAccV2NutanixEventsDataSource_WithLimit` — 33.72s — Lists events with `$limit=5` query parameter

No failures, no retries needed. Coverage: 51.8% of statements in `monitoringv2` package.

### Test logs (tail)

<details>
<summary>tail -n 500 test_logs.log</summary>

```text
=== RUN   TestAccV2NutanixEventDataSource_Basic
2026/04/29 12:38:04 [DEBUG] config wait_timeout 0
2026-04-29 07:08:04.907Z INFO - OPTIONS https://10.44.76.91:9440/api/monitoring/unversioned/info
2026-04-29 07:08:10.738Z INFO - HTTP/1.1 200 OK
2026-04-29 07:08:10.738Z INFO - Negotiated Version with server : v4.2
2026-04-29 07:08:10.738Z INFO - GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/events?%24limit=1
2026-04-29 07:08:11.052Z INFO - HTTP/1.1 200 OK
2026-04-29 07:08:11.062Z INFO - GET https://10.44.76.91:9440/api/monitoring/v4.2/serviceability/events/663c4e05-02df-4383-a00c-44ae0bc1c497
2026-04-29 07:08:11.377Z INFO - HTTP/1.1 200 OK
--- PASS: TestAccV2NutanixEventDataSource_Basic (37.68s)
=== RUN   TestAccV2NutanixEventsDataSource_Basic
--- PASS: TestAccV2NutanixEventsDataSource_Basic (32.57s)
=== RUN   TestAccV2NutanixEventsDataSource_WithLimit
--- PASS: TestAccV2NutanixEventsDataSource_WithLimit (33.72s)
PASS
coverage: 51.8% of statements
ok  	github.com/terraform-providers/terraform-provider-nutanix/nutanix/services/monitoringv2	104.753s	coverage: 51.8% of statements
```

</details>

## How this was generated

- Triggered via the Nutanix headless code-gen API (`POST /generate`).
- Prompt + inline `sdk_info.json` are stored only in the agent transcript;
  no `sdk_info.json` is committed to this branch.
