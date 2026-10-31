# NDB onboarding example

This folder shows how to run `nutanix_ndb_onboarding` with optional Prism Central details.

## Notes

- `prism_element_info` is required.
- `prism_central_info` is optional and can be toggled by setting `pc_ip`.
- `ndb_config` is optional. If omitted, the provider uses existing NDB DNS/NTP values discovered from ERA server config.
- `storage.container_name` is optional in `auto` mode; the provider picks the first discovered container when available.

## Manual acceptance checklist

Run these cases against supported NDB builds before marking the onboarding workflow release-ready:

- Safe mode onboarding: set `enable_full_onboarding = false` and verify PE registration completes without triggering storage, network, or setup.
- Full onboarding without Prism Central: omit `pc_ip` and verify PE, DNS/NTP, storage, network, and setup complete.
- Full onboarding with Prism Central: set `pc_ip`, `pc_username`, and `pc_password` and verify the PC validation/onboarding path succeeds.
- Explicit storage selection: set `storage_container` and verify that container is applied.
- Auto storage selection: leave `storage_container` empty and verify the provider selects a discovered container.
- Strict storage validation: set `selection_mode = "strict"` with an invalid container and verify apply fails before setup.
- Network skipped: set `network_details.skip = true` and verify setup proceeds without a custom network selection.
- Existing network selected: set `network_details.skip = false` with a valid `existing_network_name` and verify it is applied.
- Setup disabled: set `setup.trigger = false` and verify the resource stops before setup execution.
- Re-apply behavior: run `terraform apply` again and verify there are no unexpected duplicate onboarding side effects.

For `/clusters` payload compatibility, test at least the NDB build where the payload mismatch was observed and the latest NDB build supported by this provider branch.
