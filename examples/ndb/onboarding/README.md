# NDB onboarding example

This folder shows how to run `nutanix_ndb_onboarding` with optional Prism Central details.

## Notes

- `prism_element_info` is required.
- `prism_central_info` is optional and can be toggled by setting `pc_ip`.
- `ndb_config` is optional. If omitted, the provider uses existing NDB DNS/NTP values discovered from ERA server config.
- `storage.container_name` is optional in `auto` mode; the provider picks the first discovered container when available.

## Readiness script

`ndb_readiness_check.py` is an optional support utility for preflight checks. It is not required for provider execution.
