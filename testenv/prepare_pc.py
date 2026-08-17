#!/usr/bin/env python3
"""Prepare a Prism Central for the acceptance tests.

This is the Python port of the preEnv bash helpers and enables, in order:

1. Flow Controller (Microsegmentation / SMSP) -- run as the 'flow' step
     - resolve/create the shared managed VLAN subnet (objects.800)
     - PUT /api/prism/v4.4/management/domain-managers/{dm}/products/{id}
       with enablementState=ENABLED + FlowControllerMetadata (v4, async task)
     - then switch the policy model to Rule-Centric (FLEX) by setting
       isFNSFlexModeEnabled=true on the same product, so FLEX network security
       policies (RULETYPE_FLEX) are supported (skip with
       pc_prep.flow.policy_model=app_centric); see prepare_flow /
       switch_flow_policy_model_to_rule_centric near the bottom of this file
   (replaces the old Network Controller + v3 microseg enablement; see
   enable_flow_controller near the bottom of this file)

2. NuCalm (Self-Service) + Disaster Recovery
     - POST /api/nutanix/v3/services/nucalm            (v3)
     - poll GET /api/nutanix/v3/services/nucalm/status until
       service_enablement_status == ENABLED and service_running_status == HEALTHY
       (falls back to a settle-wait if that endpoint is unavailable)
     - POST /api/nutanix/v3/services/disaster_recovery (v3)
   (mirrors preEnv/scripts/enable_nucalm.sh)

3. Object Store (Objects)
     - POST /api/nutanix/v3/blueprints/marketplace_launch (v3)
     - poll  GET /api/nutanix/v3/services/oss/status until ENABLED/HEALTHY
   (mirrors preEnv/scripts/enable_object.sh)

   NOTE (confirmed via internal docs): Objects has *no* dedicated
   `/api/nutanix/v3/services/oss/enable` endpoint and no v4 enablement API yet --
   it is deployed as a Marketplace application via `blueprints/marketplace_launch`,
   then verified through `/api/nutanix/v3/services/oss/status`. The Objects v4
   APIs (`/api/objects/v4.0/config/object-stores`) only manage object-store
   clusters *after* the service is enabled.

4. Policy Engine enablement (deploys the Policy Engine / Policy VM)
     - pick the first free IP after the PC endpoint (or pc_prep.policy_engine.ip)
     - PUT /api/nutanix/v3/features/policy with is_enabled=true + ip_list
     - poll the feature until it reports COMPLETED / migration Succeeded
   (mirrors preEnv/scripts/enable_policy_engine.sh; this creates the
   policy_engine Zookeeper nodes the LCM step below depends on)

5. LCM Policy Engine downgrade + inventory
     - SSH to the CVM, rewrite the policy-engine Zookeeper transient_data to a
       lower version via cshell.py inside the nucalm container (requires sshpass)
     - POST /api/lcm/v4.0.a1/operations/$actions/performInventory and poll the task
     (so the newer Policy Engine version shows up as an available LCM upgrade)

6. Data protection pre-config (for the near-sync/replication tests)
     - ensure the local PC is paired with the remote PC (availability_zone.
       remote_pc_ip) via an availability zone; if not, connect them with a v3
       cloud_trust (ONPREM_CLOUD)
     - SSH to the local and remote *PCs* (ssh.pc_* creds), set the near-sync
       cerebro gflag and restart cerebro (genesis stop cerebro + cluster start)
     - do the same on both *PEs* (ssh.pe_* creds) so the workload PE cerebro
       advertises the entity-centric NearSync DR capability (single-node /
       hybrid QA clusters otherwise never protect RPO<=15min entities)
     - open the replication firewall ports between the two cluster *PEs*, both
       ways, via modify_firewall (ssh.pe_* creds)
     - ensure sync-rep storage-container parity: the source cluster's
       default-container-* names must also exist on the remote PE (SyncRep
       stretches disks into a same-named container), creating any missing
       ones via ncli (ssh.pe_* creds)
   (mirrors preEnv/data_protection.tf; local PC = PC_ENDPOINT, remote PC =
   availability_zone.remote_pc_ip; PE svm_ip/VIP from the dp config block)

7. Prism (PC unregistration) pre-config
     - ensure the local PC is paired with prism.unregister.remote_pc_ip via an
       availability zone; if not, connect them with a v3 cloud_trust
       (ONPREM_CLOUD) so the unregistration test has a connected PC to unregister
   (no-op when prism.unregister.remote_pc_ip is unset)

8. iSCSI-client VM (mirrors preEnv/scripts/create_vm_for_iscsi_clients.sh)
     - ensure the ubuntu cloud image exists, create + power on a VM, SSH in to
       configure open-iscsi, then create a Volume Group and attach the VM's
       iSCSI client so the volumesv2 iSCSI-client tests have a client to read
   (idempotent: no-op when the VM already exists; requires sshpass)

9. NGT-upgrade VM (mirrors preEnv/scripts/create_vm_for_ngt_upgrade.sh)
     - ensure the Rocky8 image exists, create + power on a VM (with a CD-ROM),
       enable NGT, insert the config ISO, install an OLD NGT bundle in-guest and
       verify -- so the vmmv2 NGT-upgrade test (VM named
       vmm.ngt.ngt_upgrade_vm_name) has a VM with an upgrade available
   (idempotent: no-op when the VM already exists; requires sshpass)

10. Self-Service (Calm) project + blueprints
     - guard: verify Self-Service is enabled from the marketplace (poll
       GET /api/nutanix/v3/services/nucalm/status; enable it if not ready)
     - create project + environment + credential + snapshot policy
       (mirrors preEnv/scripts/create_v3_project.sh)
     - import + patch the two blueprints via /blueprints/import_file, then
       launch an app and run its snapshot action
       (mirrors preEnv/scripts/upload_bp.sh)
   (idempotent: existing project/env/policy/blueprints are reused by name;
    blueprint JSONs live under testenv/payloads/)

Config
------
Reads the same nested YAML as testenv/fill_test_config.py (testenv/config.yaml)
for the PC connection (pc.endpoint/port/username/password/insecure) plus an
optional `pc_prep:` block for tunables. Override with --env; process env vars
(PC_*) take precedence. NEVER commit real secrets (config.yaml is gitignored).

Usage
-----
    python3 testenv/prepare_pc.py                     # run all steps
    python3 testenv/prepare_pc.py --dry-run           # print what it would do
    python3 testenv/prepare_pc.py --only nucalm       # one step (nucalm|objects|policy|lcm|dp|prism|iscsi|ngt|nkp|selfservice|flow)
    python3 testenv/prepare_pc.py --only selfservice  # verify Calm enabled, then create the project + blueprints
    python3 testenv/prepare_pc.py --only flow         # create the managed subnet, then enable (deploy) the Flow Controller
    python3 testenv/prepare_pc.py --skip objects      # skip a step (repeatable)
    python3 testenv/prepare_pc.py -v                  # DEBUG + HTTP traces on console

Standard library only. Reuses helpers from testenv/fill_test_config.py.
"""

from __future__ import annotations

import argparse
import base64
import hashlib
import json
import logging
import re
import shutil
import ssl
import subprocess
import sys
import tempfile
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid
from pathlib import Path

SCRIPT_DIR = Path(__file__).resolve().parent
REPO_ROOT = SCRIPT_DIR.parent
sys.path.insert(0, str(SCRIPT_DIR))

# Reuse the config loading / redaction helpers from the fill script so the two
# scripts share one config format and one secret-redaction policy.
from fill_test_config import (  # noqa: E402
    Env,
    _redact,
    _truncate,
    discover_local_cluster,
    discover_remote_pc_cluster,
    load_env,
)

logger = logging.getLogger("prepare_pc")


# --------------------------------------------------------------------------- #
# Logging (mirrors fill_test_config.setup_logging, but for this logger)
# --------------------------------------------------------------------------- #
def setup_logging(log_dir: Path | None, verbose: bool) -> Path | None:
    logger.setLevel(logging.DEBUG)
    for handler in list(logger.handlers):
        logger.removeHandler(handler)
    fmt = logging.Formatter("%(asctime)s %(levelname)-7s %(message)s", "%Y-%m-%d %H:%M:%S")

    console = logging.StreamHandler(sys.stderr)
    console.setLevel(logging.DEBUG if verbose else logging.INFO)
    console.setFormatter(fmt)
    logger.addHandler(console)

    log_file: Path | None = None
    if log_dir is not None:
        log_dir.mkdir(parents=True, exist_ok=True)
        log_file = log_dir / f"prepare_pc_{time.strftime('%Y%m%d_%H%M%S')}.log"
        file_handler = logging.FileHandler(log_file, encoding="utf-8")
        file_handler.setLevel(logging.DEBUG)
        file_handler.setFormatter(fmt)
        logger.addHandler(file_handler)
    return log_file


def _safe_body(raw: str) -> str:
    try:
        return _truncate(json.dumps(_redact(json.loads(raw))))
    except (ValueError, TypeError):
        return _truncate(raw)


# --------------------------------------------------------------------------- #
# Minimal Prism Central REST client (v3 + v4)
# --------------------------------------------------------------------------- #
class PcClient:
    def __init__(self, env: Env, dry_run: bool = False):
        self.endpoint = env.required("PC_ENDPOINT")
        self.port = env.required("PC_PORT")
        self.username = env.required("PC_USERNAME")
        self.password = env.required("PC_PASSWORD")
        self.insecure = env.bool("PC_INSECURE")
        self.dry_run = dry_run
        token = base64.b64encode(f"{self.username}:{self.password}".encode()).decode()
        self.auth_header = f"Basic {token}"
        self.ctx = ssl.create_default_context()
        if self.insecure:
            self.ctx.check_hostname = False
            self.ctx.verify_mode = ssl.CERT_NONE

    def url(self, path: str) -> str:
        return f"https://{self.endpoint}:{self.port}/{path.lstrip('/')}"

    def request(self, method: str, path: str, body: dict | None = None,
                mutating: bool = True, extra_headers: dict | None = None) -> tuple[int, object]:
        """Return (http_status, parsed_body_or_text). Honours --dry-run for
        mutating calls (POST/PUT/DELETE)."""
        full = self.url(path)
        if self.dry_run and mutating and method.upper() != "GET":
            logger.info("[dry-run] would %s %s\n  payload: %s", method, full,
                        _truncate(json.dumps(_redact(body))) if body is not None else "(none)")
            return 0, {}

        data = json.dumps(body).encode() if body is not None else None
        if body is not None:
            logger.debug("HTTP %s %s\n  request: %s", method, full,
                         _truncate(json.dumps(_redact(body))))
        else:
            logger.debug("HTTP %s %s", method, full)

        req = urllib.request.Request(full, data=data, method=method)
        req.add_header("Authorization", self.auth_header)
        req.add_header("Accept", "application/json")
        req.add_header("NTNX-Request-Id", str(uuid.uuid4()))
        if data is not None:
            req.add_header("Content-Type", "application/json")
        for key, value in (extra_headers or {}).items():
            req.add_header(key, value)
        try:
            with urllib.request.urlopen(req, context=self.ctx, timeout=120) as resp:
                raw = resp.read().decode()
                status = getattr(resp, "status", resp.getcode())
        except urllib.error.HTTPError as exc:
            detail = exc.read().decode(errors="replace")
            logger.debug("HTTP %s %s -> %d\n  response: %s", method, full, exc.code,
                         _safe_body(detail))
            return exc.code, _parse(detail)
        except urllib.error.URLError as exc:
            raise RuntimeError(f"{method} {full} -> connection error: {exc.reason}") from None

        logger.debug("HTTP %s %s -> %d\n  response: %s", method, full, status, _safe_body(raw))
        return status, _parse(raw)

    def get_etag(self, path: str) -> str:
        """GET *path* and return its response ETag header (empty string if none).

        The v4 mutating VM/volume actions (power-on, guest-tools, attach-iscsi-
        client) require an ``If-Match: <etag>`` header. ``request`` does not
        expose response headers, so this small GET fetches just the ETag."""
        full = self.url(path)
        req = urllib.request.Request(full, method="GET")
        req.add_header("Authorization", self.auth_header)
        req.add_header("Accept", "application/json")
        try:
            with urllib.request.urlopen(req, context=self.ctx, timeout=60) as resp:
                return resp.headers.get("ETag", "") or ""
        except urllib.error.HTTPError as exc:
            logger.debug("GET %s (etag) -> %d", full, exc.code)
            return ""
        except urllib.error.URLError as exc:
            raise RuntimeError(f"GET {full} -> connection error: {exc.reason}") from None


def _parse(raw: str) -> object:
    try:
        return json.loads(raw)
    except (ValueError, TypeError):
        return raw


def _dig(obj: object, *keys: str):
    """Return the first present value among dotted key paths (e.g. 'data.extId')."""
    for path in keys:
        cur = obj
        ok = True
        for part in path.split("."):
            if isinstance(cur, dict) and part in cur:
                cur = cur[part]
            else:
                ok = False
                break
        if ok and cur not in (None, ""):
            return cur
    return None


# --------------------------------------------------------------------------- #
# v4 task polling (port of preEnv/scripts/pull_task.sh)
# --------------------------------------------------------------------------- #
_ERGON_PREFIX = "ZXJnb24="  # base64("ergon"): required prefix for raw task UUIDs


def poll_task(client: PcClient, task_ext_id: str, *, timeout_seconds: int = 300,
              interval_seconds: int = 5) -> bool:
    """Poll a Prism v4 task until it terminates. Returns True on success.
    Treats 'network controller already deployed' as success (idempotent)."""
    ext_id = task_ext_id
    # A bare 36-char UUID needs the ergon prefix for the v4 tasks endpoint.
    if len(task_ext_id) == 36 and task_ext_id.count("-") == 4 and ":" not in task_ext_id:
        ext_id = f"{_ERGON_PREFIX}:{task_ext_id}"
    encoded = urllib.parse.quote(ext_id, safe="")
    path = f"api/prism/v4.0/config/tasks/{encoded}"

    logger.info("Polling task %s (api extId: %s)", task_ext_id, ext_id)
    deadline = time.monotonic() + timeout_seconds
    attempt = 0
    while time.monotonic() < deadline:
        attempt += 1
        status_code, body = client.request("GET", path, mutating=False)
        status = _dig(body, "data.status", "status")
        logger.info("[%02d] http=%s status=%s", attempt, status_code, status or "(empty)")

        if status == "SUCCEEDED":
            logger.info("Task succeeded")
            return True
        if status in ("FAILED", "CANCELLED"):
            text = json.dumps(body) if isinstance(body, (dict, list)) else str(body)
            if "already deployed" in text or "already exists" in text:
                logger.info("Resource already deployed -- continuing.")
                return True
            logger.error("Task ended with status %s: %s", status, _truncate(text))
            return False
        time.sleep(interval_seconds)

    logger.error("Timed out waiting for task after %ds", timeout_seconds)
    return False


# --------------------------------------------------------------------------- #
# NuCalm + Disaster Recovery
# --------------------------------------------------------------------------- #
def _nucalm_status(client: PcClient) -> tuple[int, str, str]:
    """Return (http_status, service_enablement_status, service_running_status)
    from GET /services/nucalm/status. The two status strings are upper-cased
    (empty when the endpoint is unavailable or the fields are absent)."""
    http, body = client.request(
        "GET", "api/nutanix/v3/services/nucalm/status", mutating=False)
    enablement = str(_dig(body, "service_enablement_status") or "").upper()
    running = str(_dig(body, "service_running_status") or "").upper()
    return http, enablement, running


def _wait_nucalm_ready(client: PcClient, *, timeout_seconds: int,
                       interval_seconds: int) -> bool | None:
    """Poll /services/nucalm/status until Self-Service is ENABLED + HEALTHY.

    Returns True when ready, False on timeout, and None when the status endpoint
    looks unavailable (repeated 4xx/5xx with no status fields) so the caller can
    fall back to a plain settle-wait on older PC builds."""
    deadline = time.monotonic() + timeout_seconds
    attempt = 0
    err_streak = 0
    enablement = running = ""
    while time.monotonic() < deadline:
        attempt += 1
        http, enablement, running = _nucalm_status(client)
        if http and http >= 400 and not enablement and not running:
            err_streak += 1
            logger.warning("[%02d] nucalm/status http=%s (endpoint unavailable? %d/%d)",
                           attempt, http, err_streak, 6)
            if err_streak >= 6:
                return None
        else:
            err_streak = 0
            logger.info("[%02d] nucalm enablement=%s running=%s", attempt,
                        enablement or "(none)", running or "(none)")
            if enablement == "ENABLED" and running == "HEALTHY":
                logger.info("NuCalm is ENABLED and HEALTHY.")
                return True
        time.sleep(interval_seconds)
    logger.error("NuCalm not ready after %ds (last enablement=%s running=%s)",
                 timeout_seconds, enablement or "(none)", running or "(none)")
    return False


def enable_nucalm(client: PcClient, env: Env) -> bool:
    logger.info("=== Enabling NuCalm (Self-Service) ===")
    status_code, body = client.request(
        "POST", "api/nutanix/v3/services/nucalm",
        {"enable_nutanix_apps": True, "state": "ENABLE"})
    logger.info("NuCalm response (http %s): %s", status_code,
                _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))

    settle = int(env.get("PC_PREP_NUCALM_SETTLE_SECONDS", "600") or 600)
    if client.dry_run:
        logger.info("[dry-run] would poll /services/nucalm/status until ENABLED+HEALTHY, "
                    "then enable Disaster Recovery")
        return True

    # Verify enablement from the marketplace by polling the status endpoint,
    # rather than blindly sleeping. Fall back to a settle-wait if the status
    # endpoint is unavailable (older PC builds).
    timeout = int(env.get("PC_PREP_SELF_SERVICE_ENABLEMENT_TIMEOUT", "1800") or 1800)
    interval = int(env.get("PC_PREP_SELF_SERVICE_ENABLEMENT_POLL_SECONDS", "15") or 15)
    ready = _wait_nucalm_ready(client, timeout_seconds=timeout, interval_seconds=interval)
    if ready is None:
        logger.warning("nucalm/status endpoint did not report readiness; "
                       "falling back to a %ds settle wait.", settle)
        _sleep_with_progress(settle, "Waiting for NuCalm to settle")
    elif ready is False:
        logger.error("NuCalm did not reach ENABLED/HEALTHY within %ds.", timeout)
        return False

    logger.info("Enabling Disaster Recovery ...")
    status_code, body = client.request(
        "POST", "api/nutanix/v3/services/disaster_recovery", {"state": "ENABLE"})
    logger.info("Disaster Recovery response (http %s): %s", status_code,
                _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
    return True


def _sleep_with_progress(seconds: int, msg: str) -> None:
    if seconds <= 0:
        return
    logger.info("%s (%ds) ...", msg, seconds)
    step = 10
    waited = 0
    while waited < seconds:
        time.sleep(min(step, seconds - waited))
        waited += step
        logger.debug("  ... %ds / %ds", min(waited, seconds), seconds)
    logger.info("%s: done.", msg)


# --------------------------------------------------------------------------- #
# Object Store (Objects) via Marketplace launch
# --------------------------------------------------------------------------- #
_OBJECTS_MP_NAME = "Objects"


def _objects_app_exists(client: PcClient) -> bool:
    """True if a Calm app named 'Objects' already exists (idempotency guard).

    Mirrors the UI's POST /apps/list {filter: name==Objects} pre-check before it
    launches the marketplace item."""
    _, body = client.request(
        "POST", "api/nutanix/v3/apps/list",
        {"length": 1, "offset": 0, "filter": f"name=={_OBJECTS_MP_NAME}"},
        mutating=False)
    entities = _dig(body, "entities")
    return bool(isinstance(entities, list) and entities)


def _first_group_value(datum: dict):
    """Pull the scalar value out of a v3 /groups entity data item, whose shape is
    {"name": <attr>, "values": [{"time": ..., "values": [<value>]}]}."""
    vals = datum.get("values") or []
    if vals and isinstance(vals[0], dict):
        inner = vals[0].get("values") or []
        if inner:
            return inner[0]
    return None


def _discover_objects_marketplace_item(client: PcClient, env: Env) -> str | None:
    """Return the UUID of the published 'Objects' marketplace item.

    Honours pc_prep.objects.marketplace_item_uuid; otherwise mirrors the Prism UI
    by querying POST /api/nutanix/v3/groups (entity_type=marketplace_item) for
    published APP items and picking the highest-versioned one named 'Objects'. The
    group entity_id is the marketplace item UUID."""
    override = env.get("PC_PREP_OBJECTS_MARKETPLACE_ITEM_UUID")
    if override:
        return override

    # Mirror the Prism marketplace search query (/dm/marketplace?name=Objects).
    # marketplace_item group queries need the AppFamily category clause and
    # grouping by app_group_uuid; without them PC returns nothing.
    query = {
        "entity_type": "marketplace_item",
        "grouping_attribute": "app_group_uuid",
        "group_member_sort_attribute": "version",
        "group_member_sort_order": "DESCENDING",
        "group_count": 12,
        "group_offset": 0,
        "group_member_count": 1,
        "group_member_attributes": [
            {"attribute": "name"}, {"attribute": "version"},
            {"attribute": "app_state"}, {"attribute": "app_group_uuid"},
        ],
        "filter_criteria": (
            f"marketplace_item_type_list==APP;name=={_OBJECTS_MP_NAME};"
            "(app_state==PUBLISHED);"
            "(category_name==AppFamily;"
            "(category_value==Nutanix,category_value==Preferred_Partners))"
        ),
    }
    status_code, body = client.request("POST", "api/nutanix/v3/groups", query,
                                       mutating=False)
    if status_code and status_code >= 400:
        logger.error("Marketplace item query failed (http %s): %s", status_code,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list))
                               else str(body)))
        if status_code in (401, 403):
            logger.error("Authentication was rejected -- check pc.username/pc.password "
                         "in config.yaml (or an account lockout on Prism Central).")
        return None

    best_uuid: str | None = None
    best_ver: tuple = ()
    for grp in (_dig(body, "group_results") or []):
        for ent in (grp.get("entity_results") or []):
            data = {d.get("name"): _first_group_value(d) for d in ent.get("data", [])}
            if data.get("name") != _OBJECTS_MP_NAME:
                continue
            ver = data.get("version") or "0"
            ver_key = tuple(int(p) if str(p).isdigit() else 0 for p in str(ver).split("."))
            item_uuid = ent.get("entity_id")
            if item_uuid and ver_key >= best_ver:
                best_uuid, best_ver = item_uuid, ver_key
    if not best_uuid:
        logger.error("Could not find a published '%s' marketplace item; set "
                     "pc_prep.objects.marketplace_item_uuid in config.yaml.",
                     _OBJECTS_MP_NAME)
    return best_uuid


def _build_objects_launch_payload(client: PcClient, env: Env) -> dict | None:
    """Build the /blueprints/marketplace_launch body for Objects.

    Mirrors the Prism UI flow: GET the Objects marketplace item, pull its blueprint
    template spec (spec.resources.app_blueprint_template.spec -- which already
    carries this PC's substrate address), and wrap it for launch with a unique app
    name. If pc_prep.objects.payload_file is set and exists, that pre-baked payload
    is used verbatim instead (back-compat / manual override)."""
    payload_file = env.get("PC_PREP_OBJECTS_PAYLOAD_FILE")
    if payload_file:
        path = Path(payload_file)
        if not path.is_absolute():
            path = (REPO_ROOT / payload_file).resolve()
        if path.exists():
            try:
                logger.info("Using pre-baked Objects payload %s", path)
                return json.loads(path.read_text())
            except ValueError as exc:
                logger.error("Objects payload %s is not valid JSON: %s", path, exc)
                return None
        logger.info("Objects payload_file %s not found -- discovering the marketplace "
                    "item instead.", path)

    item_uuid = _discover_objects_marketplace_item(client, env)
    if not item_uuid:
        return None
    logger.info("Using Objects marketplace item %s", item_uuid)

    _, item = client.request(
        "GET", f"api/nutanix/v3/calm_marketplace_items/{item_uuid}", mutating=False)
    bp_spec = _dig(item, "spec.resources.app_blueprint_template.spec")
    resources = bp_spec.get("resources") if isinstance(bp_spec, dict) else None
    if not resources:
        logger.error("Objects marketplace item %s has no blueprint template spec.",
                     item_uuid)
        return None

    profile = "DefaultProfile"
    profiles = resources.get("app_profile_list") or []
    if profiles and isinstance(profiles[0], dict) and profiles[0].get("name"):
        profile = profiles[0]["name"]

    return {
        "spec": {
            "description": bp_spec.get("description") or "Objects Blueprint",
            "resources": resources,
            "source_marketplace_name": _dig(item, "spec.name") or _OBJECTS_MP_NAME,
            "source_marketplace_version": _dig(item, "spec.resources.version") or "",
            "app_blueprint_name": f"{_OBJECTS_MP_NAME} {uuid.uuid4().hex[:8]}",
            "environment_profile_pairs": [{"app_profile": {"name": profile}}],
        },
        "api_version": "3.0",
        "metadata": {"kind": "blueprint", "categories": {"AppFamily": "Nutanix"}},
    }


def enable_object_store(client: PcClient, env: Env) -> bool:
    """Enable Objects by launching its Marketplace app (mirrors the Prism UI).

    Sequence (from a captured browser session):
      1. POST /apps/list {name==Objects}                    -- idempotency check
      2. discover + GET /calm_marketplace_items/{uuid}      -- blueprint spec
      3. POST /blueprints/marketplace_launch {spec}         -- launch the app
      4. poll GET /services/oss/status until ENABLED/HEALTHY
    """
    logger.info("=== Enabling Object Store (Objects) ===")

    if client.dry_run:
        logger.info("[dry-run] would discover the Objects marketplace item, POST "
                    "/blueprints/marketplace_launch, and poll /services/oss/status.")
        return True

    if _objects_app_exists(client):
        logger.info("An '%s' app already exists -- skipping marketplace launch.",
                    _OBJECTS_MP_NAME)
    else:
        payload = _build_objects_launch_payload(client, env)
        if payload is None:
            return False
        status_code, body = client.request(
            "POST", "api/nutanix/v3/blueprints/marketplace_launch", payload)
        logger.info("Marketplace launch response (http %s): %s", status_code,
                    _truncate(json.dumps(body) if isinstance(body, (dict, list))
                              else str(body)))
        if status_code and status_code >= 400:
            text = json.dumps(body) if isinstance(body, (dict, list)) else str(body)
            if "already" not in text.lower():
                logger.error("Objects marketplace launch failed (http %s)", status_code)
                return False
            logger.warning("Objects launch reported 'already' (http %s) -- continuing "
                           "to status poll.", status_code)

    if not env.bool("PC_PREP_OBJECTS_POLL", True):
        logger.info("Skipping Objects status poll (pc_prep.objects.poll=false)")
        return True
    timeout = int(env.get("PC_PREP_OBJECTS_POLL_TIMEOUT_SECONDS", "1800") or 1800)
    return _poll_oss_status(client, timeout_seconds=timeout)


def _poll_oss_status(client: PcClient, *, timeout_seconds: int, interval_seconds: int = 30) -> bool:
    logger.info("Polling Objects service status (timeout %ds) ...", timeout_seconds)
    deadline = time.monotonic() + timeout_seconds
    attempt = 0
    while time.monotonic() < deadline:
        attempt += 1
        status_code, body = client.request(
            "GET", "api/nutanix/v3/services/oss/status", mutating=False)
        enablement = _dig(body, "service_enablement_status")
        running = _dig(body, "service_running_status")
        logger.info("[%02d] http=%s enablement=%s running=%s", attempt, status_code,
                    enablement or "(empty)", running or "(empty)")
        if enablement == "ENABLED" and running == "HEALTHY":
            logger.info("Object Store enabled and healthy.")
            return True
        time.sleep(interval_seconds)
    logger.error("Timed out waiting for Objects to become ENABLED/HEALTHY after %ds",
                 timeout_seconds)
    return False


# --------------------------------------------------------------------------- #
# Step 4: Policy Engine enablement (deploys the Policy Engine / Policy VM)
# --------------------------------------------------------------------------- #
_POLICY_FEATURE_PATH = "api/nutanix/v3/features/policy"


def _ping_alive(ip: str, *, timeout_seconds: int = 1) -> bool:
    """Return True if *ip* answers a single ICMP echo. Cross-platform (the -W
    flag differs between macOS/BSD and Linux)."""
    if sys.platform == "darwin":
        # macOS: -W is per-packet timeout in milliseconds.
        argv = ["ping", "-c", "1", "-W", str(timeout_seconds * 1000), ip]
    else:
        # Linux: -W is in seconds.
        argv = ["ping", "-c", "1", "-W", str(timeout_seconds), ip]
    try:
        proc = subprocess.run(argv, capture_output=True, text=True,
                              timeout=timeout_seconds + 3)
        return proc.returncode == 0
    except (subprocess.TimeoutExpired, OSError):
        return False


def _first_free_ip(base_ip: str, *, max_probe: int = 32) -> str | None:
    """Return the first IP after *base_ip* (same /24) that does not answer ping.
    Mirrors the bash loop that picks the next free IP for the Policy Engine VM."""
    parts = base_ip.split(".")
    if len(parts) != 4 or not all(p.isdigit() for p in parts):
        logger.error("Cannot auto-pick Policy Engine IP: %r is not a dotted IPv4", base_ip)
        return None
    o1, o2, o3, o4 = (int(p) for p in parts)
    for _ in range(max_probe):
        o4 += 1
        if o4 > 255:
            logger.error("Ran off the end of the /24 while looking for a free IP after %s",
                         base_ip)
            return None
        candidate = f"{o1}.{o2}.{o3}.{o4}"
        if not _ping_alive(candidate):
            return candidate
        logger.debug("  %s is in use, trying next", candidate)
    logger.error("No free IP found within %d probes after %s", max_probe, base_ip)
    return None


def enable_policy_engine(client: PcClient, env: Env) -> bool:
    """Enable and deploy the Calm Policy Engine (Policy VM).

    Mirrors preEnv/scripts/enable_policy_engine.sh:
      1. Pick the Policy Engine VM IP (explicit pc_prep.policy_engine.ip, else the
         first free IP after the PC endpoint).
      2. Read the current feature spec_version, then PUT the policy feature with
         is_enabled=true and ip_list=[chosen_ip].
      3. Poll the feature until it reports COMPLETED / migration Succeeded.

    This must run before the `lcm` step, which downgrades the deployed engine.
    """
    logger.info("=== Enabling Policy Engine ===")

    engine_ip = env.get("PC_PREP_POLICY_ENGINE_IP")
    if not engine_ip:
        base = env.get("PC_ENDPOINT")
        logger.info("No pc_prep.policy_engine.ip set; probing for the first free IP "
                    "after %s ...", base)
        engine_ip = _first_free_ip(base)
        if not engine_ip:
            return False
    logger.info("Policy Engine VM IP: %s", engine_ip)

    # Read the current feature spec version (defaults to 0 if the feature has
    # never been touched). The PUT echoes this same value back, matching the
    # bash flow.
    status_code, body = client.request("GET", _POLICY_FEATURE_PATH, mutating=False)
    spec_version = _dig(body, "metadata.spec_version")
    try:
        spec_version = int(spec_version) if spec_version is not None else 0
    except (TypeError, ValueError):
        spec_version = 0
    logger.info("Current policy feature spec_version: %s", spec_version)

    payload = {
        "api_version": "3.1",
        "metadata": {
            "last_update_time": "1750334336803716",
            "creation_time": "1750334336803716",
            "spec_version": spec_version,
            "name": "",
            "kind": "calm_feature",
            "uuid": "02276f82-8e73-4bc4-84ed-0928be3f9cad",
        },
        "spec": {
            "feature_status": {
                "is_ignored": False,
                "is_enabled": True,
                "config": {"data": {"ip_list": [engine_ip]}},
            }
        },
    }
    status_code, body = client.request("PUT", _POLICY_FEATURE_PATH, payload)
    logger.info("Policy Engine enable response (http %s): %s", status_code,
                _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))

    if client.dry_run:
        logger.info("[dry-run] would poll %s until COMPLETED/Succeeded", _POLICY_FEATURE_PATH)
        return True

    # "spec version mismatch" means it is already deployed -- treat as success.
    text = json.dumps(body) if isinstance(body, (dict, list)) else str(body)
    if re.search(r"spec version mismatch", text, re.IGNORECASE):
        logger.info("Policy Engine already exists (spec version mismatch) -- continuing.")

    max_attempts = int(env.get("PC_PREP_POLICY_ENGINE_MAX_ATTEMPTS", "160") or 160)
    sleep_seconds = int(env.get("PC_PREP_POLICY_ENGINE_SLEEP_SECONDS", "5") or 5)
    return _poll_policy_engine(client, max_attempts=max_attempts, sleep_seconds=sleep_seconds)


_POLICY_TERMINAL_FAIL = {"FAILED", "FAILURE", "ERROR", "CANCELLED", "ABORTED"}


def _poll_policy_engine(client: PcClient, *, max_attempts: int, sleep_seconds: int) -> bool:
    logger.info("Polling Policy Engine feature (max %d attempts, %ds apart) ...",
                max_attempts, sleep_seconds)
    start = time.monotonic()
    for attempt in range(1, max_attempts + 1):
        _, body = client.request("GET", _POLICY_FEATURE_PATH, mutating=False)
        state = _dig(body, "status.feature_status.config.state") or "(unknown)"
        state_message = _dig(body, "status.feature_status.config.state_message") or "(no message)"
        current_step = _dig(body, "status.feature_status.config.progress.current_step") or "?"
        total_steps = _dig(body, "status.feature_status.config.progress.total_steps") or "?"
        migration_state = _dig(body, "status.feature_status.migration_status.state")
        elapsed = int(time.monotonic() - start)
        logger.info("[%02d] %-10s step %s/%s  +%3ds  %s", attempt, state, current_step,
                    total_steps, elapsed, _truncate(str(state_message)))

        if state == "COMPLETED" or migration_state == "Succeeded":
            logger.info("Policy Engine feature finished successfully.")
            return True
        if state in _POLICY_TERMINAL_FAIL or migration_state == "Failed":
            logger.error("Policy Engine feature failed: state=%s migration=%s (%s)",
                         state, migration_state, _truncate(str(state_message)))
            return False
        time.sleep(sleep_seconds)

    logger.error("Timed out waiting for Policy Engine after %d attempts (%ds)",
                 max_attempts, max_attempts * sleep_seconds)
    return False


# --------------------------------------------------------------------------- #
# Step 5: Downgrade LCM Policy Engine + run inventory
# --------------------------------------------------------------------------- #
def _ssh_control_path(user: str, host: str) -> str:
    """A short, stable master-socket path for connection multiplexing. Kept
    short to stay under the ~104-char UNIX socket path limit on macOS."""
    digest = hashlib.sha1(f"{user}@{host}".encode()).hexdigest()[:12]
    return str(Path(tempfile.gettempdir()) / f"pp_ssh_{digest}")


def _ssh(host: str, user: str, password: str, command: str, *, timeout: int = 120) -> str:
    """Run a command on a remote host over SSH and return stdout.

    Uses ``sshpass`` for password auth (matching the CVM 'nutanix' account, which
    does not support key auth out of the box). Raises RuntimeError on failure.

    To avoid tripping the account lockout, connections are multiplexed: the first
    call authenticates once and opens a master socket (ControlMaster), and every
    later call in the same run rides that socket without re-authenticating. The
    password is only tried once per connection (NumberOfPasswordPrompts=1) so a
    wrong/locked credential fails fast instead of burning several attempts.
    """
    if shutil.which("sshpass") is None:
        raise RuntimeError(
            "sshpass is required for SSH password auth but was not found in PATH. "
            "Install it (macOS: brew install hudochenkov/sshpass/sshpass; "
            "Ubuntu: apt-get install sshpass; RHEL: yum install sshpass).")
    argv = [
        "sshpass", "-p", password,
        "ssh",
        "-o", "StrictHostKeyChecking=no",
        "-o", "UserKnownHostsFile=/dev/null",
        "-o", "PubkeyAuthentication=no",
        "-o", "PreferredAuthentications=password",
        "-o", "NumberOfPasswordPrompts=1",
        "-o", "ControlMaster=auto",
        "-o", "ControlPersist=60",
        "-o", f"ControlPath={_ssh_control_path(user, host)}",
        f"{user}@{host}", command,
    ]
    # Log with the password masked.
    logger.debug("SSH %s@%s: %s", user, host, _truncate(command))
    proc = subprocess.run(argv, capture_output=True, text=True, timeout=timeout)
    if proc.returncode != 0:
        stderr = proc.stderr.strip()
        if proc.returncode == 255 and "permission denied" in stderr.lower():
            raise RuntimeError(
                f"ssh {user}@{host}: authentication rejected. Check ssh.pc_username / "
                f"ssh.pc_password (and any exported SSH_PC_PASSWORD overriding config). "
                f"If the credentials are correct, the '{user}' account may be temporarily "
                f"locked out from repeated logins -- wait a few minutes and retry.")
        raise RuntimeError(
            f"ssh {user}@{host} exited {proc.returncode}: "
            f"{_truncate(stderr or proc.stdout.strip())}")
    return proc.stdout


_LCM_ZK_KEY = "/appliance/logical/policy_engine/transient_data"
_LCM_STATUS_KEY = "/appliance/logical/policy_engine/status"
_CONFIGURE_LCM = "/home/nutanix/cluster/bin/lcm/configure_lcm"


def _read_lcm_url(host: str, ssh_user: str, ssh_pass: str) -> str:
    """Return the currently configured LCM 'url' from `configure_lcm -p`."""
    out = _ssh(host, ssh_user, ssh_pass, f"{_CONFIGURE_LCM} -p")
    for line in out.splitlines():
        if line.strip().lower().startswith("url:"):
            return line.split(":", 1)[1].strip()
    return ""


def _set_lcm_source_url(host: str, ssh_user: str, ssh_pass: str, url: str) -> bool:
    """Point LCM at a framework source URL via `configure_lcm -u` before inventory.

    LCM always tries to sync its own framework at the start of every inventory
    (there is no flag to skip this). When the configured URL advertises a *lower*
    framework than is installed, the pre-inventory framework auto-update runs the
    test_downgrade precheck and fails the whole inventory (LIF-20039). Pointing
    LCM at a URL whose framework version is >= the installed one makes that
    precheck pass.

    Returns True only if `configure_lcm -p` confirms the URL actually changed to
    the requested value (compared ignoring a trailing '/'). configure_lcm -u does
    not validate the URL, so a silent no-op means the change didn't persist -- we
    surface the command output and fail so the caller can skip a doomed inventory.
    """
    want = url.rstrip("/")
    logger.info("Setting LCM source URL to %s before inventory ...", url)
    try:
        out = _ssh(host, ssh_user, ssh_pass, f"{_CONFIGURE_LCM} -u {url}")
        if out.strip():
            logger.info("configure_lcm -u output: %s", _truncate(out.strip()))
        current = _read_lcm_url(host, ssh_user, ssh_pass)
    except RuntimeError as exc:
        logger.error("configure_lcm -u failed: %s", exc)
        return False

    logger.info("LCM url is now: %s", current or "(empty)")
    if current.rstrip("/") != want:
        logger.error("LCM URL did not change to %s (still %s). configure_lcm -u did not "
                     "persist -- check the URL is well-formed and that this PC is allowed "
                     "to set it; then set pc_prep.lcm.framework_url and retry.",
                     url, current or "(empty)")
        return False
    return True


def downgrade_lcm(client: PcClient, env: Env) -> bool:
    """Downgrade the Policy Engine (Calm) LCM entity to a lower version, then run
    an LCM inventory so the newer version shows up as an available upgrade.

    Mirrors the old preEnv downgrade flow but in pure Python:
      1. SSH to the PC CVM and read the policy-engine transient_data from Zookeeper.
      2. Parse the node id + ip, then rewrite that zk node with the target version
         via cshell.py inside the nucalm container.
      3. Verify the zk value now reports the target version.
      3.5. (optional) Point LCM at a framework source URL whose version is >= the
           installed framework so the pre-inventory framework auto-update's
           test_downgrade precheck (LIF-20039) passes.
      4. POST /api/lcm/v4.0.a1/operations/$actions/performInventory and poll the task.

    Config (pc_prep.lcm):
      * downgrade_version - the Policy Engine version to set (required)
      * inventory_timeout - seconds to wait for the inventory task (default 1800)
      * framework_url     - LCM source URL to set before inventory (optional; only
                            needed when the current URL advertises an older
                            framework than is installed)
      * SSH creds come from ssh.pc_username / ssh.pc_password.
    """
    version = env.get("PC_PREP_LCM_DOWNGRADE_VERSION")
    if not version:
        logger.error("Cannot downgrade LCM: pc_prep.lcm.downgrade_version is empty")
        return False

    host = env.get("PC_ENDPOINT")
    ssh_user = env.get("SSH_PC_USERNAME")
    ssh_pass = env.get("SSH_PC_PASSWORD")
    if not (host and ssh_user and ssh_pass):
        logger.error("Cannot downgrade LCM: PC_ENDPOINT / SSH_PC_USERNAME / "
                     "SSH_PC_PASSWORD must be set")
        return False

    logger.info("=== Downgrading LCM Policy Engine to %s ===", version)

    if client.dry_run:
        logger.info("[dry-run] would SSH %s@%s to read %s, rewrite the policy-engine "
                    "zk node to version %s, verify, then POST performInventory",
                    ssh_user, host, _LCM_ZK_KEY, version)
        return True

    # 0. Confirm the Policy Engine is actually deployed. The transient_data node is
    #    only created once the policy container comes up -- enabling NuCalm alone is
    #    not enough; the Policy Engine VM must be deployed (see the status node).
    zkcat = "/usr/local/nutanix/cluster/bin/zkcat"
    try:
        status = _ssh(host, ssh_user, ssh_pass, f"{zkcat} {_LCM_STATUS_KEY}")
        logger.debug("policy_engine status: %s", _truncate(status))
        if '"is_enabled":true' not in status.replace(" ", ""):
            logger.error("Policy Engine is not enabled on this PC (status=%s). Deploy "
                         "the Policy Engine VM before downgrading.", _truncate(status))
            return False
    except RuntimeError as exc:
        if "no node" in str(exc).lower():
            logger.error("Policy Engine is not deployed on this PC: zk node %s does not "
                         "exist. Enable NuCalm and deploy the Policy Engine VM first, "
                         "then retry.", _LCM_STATUS_KEY)
            return False
        raise

    # 1. Read current transient_data from Zookeeper.
    read_cmd = f"{zkcat} {_LCM_ZK_KEY}"
    try:
        current = _ssh(host, ssh_user, ssh_pass, read_cmd)
    except RuntimeError as exc:
        if "no node" in str(exc).lower():
            logger.error("%s does not exist even though Policy Engine is enabled -- the "
                         "policy container may still be starting. Wait for it to come up "
                         "(or restart policy-container.service on the policy VM) and retry.",
                         _LCM_ZK_KEY)
            return False
        raise
    logger.debug("policy_engine transient_data: %s", _truncate(current))

    node_match = re.search(r'"node_map"\s*:\s*{\s*"([^"]+)"', current)
    ip_match = re.search(r'"ip"\s*:\s*"([^"]+)"', current)
    if not node_match or not ip_match:
        logger.error("Could not parse node_id / ip from policy_engine transient_data")
        return False
    node_id, ip = node_match.group(1), ip_match.group(1)
    logger.info("Policy Engine node_id=%s ip=%s", node_id, ip)

    # 2. Rewrite the zk node with the target version via cshell.py in the nucalm
    #    container. The python snippet is base64-encoded to avoid nested-quoting
    #    issues over SSH.
    snippet = (
        f'val = b"""{{"node_map":{{"{node_id}":{{"ip":"{ip}","version":"{version}"}}}}}}"""\n'
        "from calm.pkg.common.zk_session import get_zookeeper_session\n"
        "zk_handle = get_zookeeper_session()\n"
        f'zk_key = "{_LCM_ZK_KEY}"\n'
        "zk_handle.set(zk_key, val)\n"
        "exit\n"
    )
    b64 = base64.b64encode(snippet.encode()).decode()
    set_cmd = ("docker exec -i nucalm bash -c "
               f"'source /home/calm/venv3/bin/activate && echo {b64} | base64 -d | cshell.py'")
    _ssh(host, ssh_user, ssh_pass, set_cmd)

    # 3. Verify the version was applied.
    verify = _ssh(host, ssh_user, ssh_pass, read_cmd)
    if version not in verify:
        logger.error("Downgrade failed: version %s not present in transient_data", version)
        return False
    logger.info("Policy Engine downgraded to %s", version)

    # 3.5. LCM syncs its own framework from the configured URL at the start of
    #      every inventory (there is no flag to skip this). If that URL advertises
    #      a lower framework than is installed, the pre-inventory auto-update runs
    #      the `test_downgrade` precheck and fails the whole inventory (LIF-20039).
    #      When pc_prep.lcm.framework_url is set, point LCM at a build whose
    #      framework version is >= the installed one so the precheck passes.
    # framework_url = env.get("PC_PREP_LCM_FRAMEWORK_URL")
    # if framework_url:
    #     if not _set_lcm_source_url(host, ssh_user, ssh_pass, framework_url):
    #         logger.error("Not running inventory: the LCM framework URL could not be set, so "
    #                      "the pre-inventory framework auto-update would fail test_downgrade "
    #                      "(LIF-20039) again.")
    #         return False

    # 4. Kick off an LCM inventory and wait for it.
    timeout = int(env.get("PC_PREP_LCM_INVENTORY_TIMEOUT", "1800") or 1800)
    return _perform_inventory(client, env, timeout_seconds=timeout)


def _discover_pc_cluster_ext_id(client: PcClient) -> str:
    """Return the extId of the PRISM_CENTRAL cluster entity on this PC."""
    fltr = ("config/clusterFunction/any(t:t eq "
            "Clustermgmt.Config.ClusterFunctionRef'PRISM_CENTRAL')")
    path = "api/clustermgmt/v4.0/config/clusters?$filter=" + urllib.parse.quote(fltr, safe="")
    _, body = client.request("GET", path, mutating=False)
    items = _dig(body, "data") or []
    if isinstance(items, list) and items:
        return items[0].get("extId") or ""
    return ""


def _perform_inventory(client: PcClient, env: Env, *, timeout_seconds: int) -> bool:
    """Trigger an LCM inventory and poll the resulting task.

    Uses the current LCM UI endpoint (api/lifecycle/v4.3/operations/$actions/
    inventory) with an explicit clusterExtId + X-Cluster-Id, falling back to the
    older api/lcm/v4.0.a1/.../performInventory if the v4.3 route is unavailable.
    clusterExtId defaults to the PRISM_CENTRAL cluster (the PC itself, which owns
    the Policy Engine) but can be pinned via pc_prep.lcm.cluster_ext_id.
    """
    cluster_ext_id = env.get("PC_PREP_LCM_CLUSTER_EXT_ID") or _discover_pc_cluster_ext_id(client)
    if cluster_ext_id:
        logger.info("Performing LCM inventory (v4.3, clusterExtId=%s) ...", cluster_ext_id)
        payload = {
            "credentials": [],
            "inventoryType": "FULL",
            "clusterExtId": cluster_ext_id,
            "$objectType": "lifecycle.v4.operations.InventorySpec",
        }
        status_code, body = client.request(
            "POST", "api/lifecycle/v4.3/operations/$actions/inventory", payload,
            extra_headers={"X-Cluster-Id": cluster_ext_id})
        logger.info("inventory response (http %s): %s", status_code,
                    _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
        task_ext_id = _dig(body, "data.extId", "data.extid", "extId")
        if task_ext_id:
            return poll_task(client, task_ext_id, timeout_seconds=timeout_seconds)
        logger.warning("v4.3 inventory returned no task extId (http %s); falling back to "
                       "the v4.0.a1 performInventory endpoint.", status_code)

    logger.info("Performing LCM inventory (v4.0.a1 performInventory) ...")
    status_code, body = client.request(
        "POST", "api/lcm/v4.0.a1/operations/$actions/performInventory", {})
    logger.info("performInventory response (http %s): %s", status_code,
                _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
    task_ext_id = _dig(body, "data.extId", "data.extid", "extId")
    if not task_ext_id:
        logger.error("No task extId in performInventory response")
        return False
    return poll_task(client, task_ext_id, timeout_seconds=timeout_seconds)


# --------------------------------------------------------------------------- #
# Data protection pre-config (mirrors preEnv/data_protection.tf)
# --------------------------------------------------------------------------- #
# Cerebro gflag that relaxes the near-sync precheck for the QA/acceptance tests.
_DP_CEREBRO_GFLAG = "--near_sync_test_hook_override_pre_checks_for_qa_test=true"
# Replication ports opened between the two clusters (cerebro/stargate/etc.).
_DP_FIREWALL_PORTS = "2030,2036,2073,2090,8740"
# Upper bound for `cluster start` after stopping cerebro (rc 124 is treated ok).
_DP_CEREBRO_TIMEOUT = 300
# ncli path on the CVMs, used for the sync-rep storage-container parity step.
_NCLI = "/home/nutanix/prism/cli/ncli"


def _remote_pc_local_az(remote_ip: str, port: str, user: str, passwd: str,
                        insecure: bool) -> tuple[str, str]:
    """Return (cluster_ext_id, display_name) of the remote PC's own 'Local AZ'
    (the availability_zones/list entity whose management_plane_type == Local).

    That cluster_ext_id is exactly what the local PC records as ``management_url``
    for the paired (non-Local) AZ, and display_name is the ``PC_<advertised-ip>``
    label. Both are stable identifiers that survive the connect-IP vs advertised-IP
    mismatch (e.g. we connect to the PC on its <connect-ip> but the PC advertises
    itself as PC_<advertised-ip>). Returns ("", "") on any failure (caller falls
    back to IP matching)."""
    if not (remote_ip and user and passwd):
        return "", ""
    url = f"https://{remote_ip}:{port}/api/nutanix/v3/availability_zones/list"
    token = base64.b64encode(f"{user}:{passwd}".encode()).decode()
    ctx = ssl.create_default_context()
    if insecure:
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    data = json.dumps({"kind": "availability_zone", "length": 200}).encode()
    req = urllib.request.Request(url, data=data, method="POST")
    req.add_header("Authorization", f"Basic {token}")
    req.add_header("Accept", "application/json")
    req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=30) as resp:
            body = json.loads(resp.read().decode() or "{}")
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.debug("remote PC %s AZ identity lookup failed: %s", remote_ip,
                     getattr(exc, "reason", exc))
        return "", ""
    for ent in (body.get("entities") or []):
        res = (ent.get("status") or {}).get("resources") or {}
        if str(res.get("management_plane_type") or "").lower() == "local":
            return (res.get("management_url") or "", res.get("display_name") or "")
    return "", ""


def _az_is_connected(client: PcClient, remote_ip: str, *, remote_ext_id: str = "",
                     remote_display_name: str = "", quiet: bool = False) -> bool:
    """Return True when the local PC has an availability zone pointing at the
    remote PC (i.e. the two PCs are paired).

    Lists availability zones via the v3 API and looks for a non-Local zone that
    matches the remote PC. Matching prefers the canonical *remote_ext_id* (the
    remote PC's cluster ext_id, recorded as the AZ management_url) and
    *remote_display_name*, since the AZ is labelled by the PC's advertised IP --
    which may differ from the connect IP. Falls back to *remote_ip* substring
    matching. Warns (but still returns True) if the matched zone isn't COMPLETE.
    This is a read, so it runs even under --dry-run. With quiet=True the
    'not found' case logs at debug instead of error."""
    status, body = client.request(
        "POST", "api/nutanix/v3/availability_zones/list",
        {"kind": "availability_zone", "length": 200}, mutating=False)
    if status not in (0, 200):
        logger.error("availability_zones/list returned HTTP %s; cannot verify AZ to %s",
                     status, remote_ip)
        return False
    entities = body.get("entities", []) if isinstance(body, dict) else []
    for ent in entities:
        res = (ent.get("status") or {}).get("resources") or {}
        if str(res.get("management_plane_type") or "").lower() == "local":
            continue
        name = (ent.get("status") or {}).get("name") or ""
        state = (ent.get("status") or {}).get("state") or ""
        mgmt_url = str(res.get("management_url") or "")
        disp = str(res.get("display_name") or "")
        matched = (
            (remote_ext_id and mgmt_url == remote_ext_id)
            or (remote_display_name and disp == remote_display_name)
            or (remote_ip and (remote_ip in mgmt_url or remote_ip in disp or remote_ip in name))
        )
        if matched:
            logger.info("Availability zone for remote PC %s is present "
                        "(name=%r, mgmt_url=%s, state=%s)",
                        remote_ip, name, mgmt_url or "?", state or "?")
            if state and state.upper() != "COMPLETE":
                logger.warning("Availability zone for %s is in state %s (expected COMPLETE); "
                               "the pairing may still be settling", remote_ip, state)
            return True
    msg = ("No availability zone on the local PC (%s) points at the remote PC %s",
           client.endpoint, remote_ip)
    if quiet:
        logger.debug(*msg)
    else:
        logger.error(*msg)
    return False


def _az_connect(client: PcClient, remote_ip: str, remote_user: str, remote_pass: str,
                *, remote_ext_id: str = "", remote_display_name: str = "",
                timeout_seconds: int = 300, interval_seconds: int = 10) -> bool:
    """Pair the local PC with the remote PC at *remote_ip* by creating a v3
    cloud_trust (ONPREM_CLOUD), mirroring
    POST /api/nutanix/v3/cloud_trusts. The pairing creates the AvailabilityZone /
    CloudTrust / RemoteConnection entities; we then poll availability_zones/list
    until the AZ shows up. Returns True once connected (or under --dry-run)."""
    if not (remote_user and remote_pass):
        logger.error("Cannot connect AZ to %s: remote PC creds not set "
                     "(az.remote_pc_username / az.remote_pc_password or the local PC creds)",
                     remote_ip)
        return False
    payload = {
        "spec": {
            "name": "",  # PC derives "PC_<ip>"
            "resources": {
                "url": remote_ip,
                "username": remote_user,
                "password": remote_pass,
                "cloud_type": "ONPREM_CLOUD",
            },
            "description": "",
        },
        "metadata": {"kind": "cloud_trust"},
        "api_version": "3.1.0",
    }
    logger.info("Pairing local PC %s with remote PC %s via cloud_trust (ONPREM_CLOUD) ...",
                client.endpoint, remote_ip)
    status, body = client.request("POST", "api/nutanix/v3/cloud_trusts", payload)
    if client.dry_run:
        logger.info("[dry-run] would create a cloud_trust to %s and wait for the AZ", remote_ip)
        return True
    if status not in (200, 201, 202):
        logger.error("cloud_trust create to %s -> HTTP %s: %s", remote_ip, status,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
        return False
    ref = _dig(body, "status.execution_context.task_uuid", "metadata.uuid")
    logger.info("cloud_trust request accepted (ref=%s); waiting up to %ds for the AZ to connect ...",
                ref or "?", timeout_seconds)
    deadline = time.monotonic() + timeout_seconds
    while time.monotonic() < deadline:
        time.sleep(interval_seconds)
        if _az_is_connected(client, remote_ip, remote_ext_id=remote_ext_id,
                            remote_display_name=remote_display_name, quiet=True):
            return True
        logger.info("... AZ to %s not connected yet; retrying", remote_ip)
    logger.error("Timed out after %ds waiting for the AZ to %s to connect", timeout_seconds, remote_ip)
    return False


def _ensure_az_connected(client: PcClient, remote_ip: str, remote_user: str,
                         remote_pass: str) -> bool:
    """Ensure the local PC is AZ-paired with *remote_ip*; create the pairing via a
    v3 cloud_trust (ONPREM_CLOUD) when it isn't. Returns True when connected.

    Resolves the remote PC's canonical identity (cluster ext_id + advertised
    display_name) up front so detection is robust to the connect-IP vs
    advertised-IP mismatch (the paired AZ is labelled by the PC's advertised IP,
    not the IP we connect to)."""
    remote_ext_id, remote_disp = _remote_pc_local_az(
        remote_ip, client.port, remote_user, remote_pass, client.insecure)
    if remote_ext_id or remote_disp:
        logger.info("Remote PC %s identity: cluster_ext_id=%s display_name=%s",
                    remote_ip, remote_ext_id or "<?>", remote_disp or "<?>")
    else:
        logger.debug("Could not resolve remote PC %s identity; falling back to IP matching",
                     remote_ip)
    if _az_is_connected(client, remote_ip, remote_ext_id=remote_ext_id,
                        remote_display_name=remote_disp, quiet=True):
        logger.info("Availability zone to %s already connected.", remote_ip)
        return True
    logger.info("No availability zone to %s yet; connecting via cloud_trust ...", remote_ip)
    return _az_connect(client, remote_ip, remote_user, remote_pass,
                       remote_ext_id=remote_ext_id, remote_display_name=remote_disp)


def _dp_restart_cerebro(host: str, ssh_user: str, ssh_pass: str, timeout_seconds: int) -> None:
    """Set the near-sync cerebro gflag on *host* and restart cerebro.

    Mirrors the cerebro_pc* null_resources: write the gflag, `genesis stop cerebro`
    then `cluster start`. `cluster start` can block for a long time waiting on
    services, so it is wrapped in `timeout`; rc 124 (timed out) is treated as
    success because the gflag write + cerebro stop already happened.
    """
    command = (
        "set -u; "
        f'echo "{_DP_CEREBRO_GFLAG}" > ~/config/cerebro.gflags && '
        "/usr/local/nutanix/cluster/bin/genesis stop cerebro && "
        f"timeout {timeout_seconds} /usr/local/nutanix/cluster/bin/cluster start; "
        "rc=$?; if [ $rc -eq 0 ] || [ $rc -eq 124 ]; then "
        'if [ $rc -eq 124 ]; then echo "[warn] cluster start did not finish in time; '
        'gflag is in place, continuing." >&2; fi; exit 0; else exit $rc; fi'
    )
    # Give SSH more than the inner cluster-start timeout so it can return the rc.
    _ssh(host, ssh_user, ssh_pass, command, timeout=timeout_seconds + 120)


def _dp_modify_firewall(pe: str, ssh_user: str, ssh_pass: str,
                        peer_pe: str, peer_vip: str) -> None:
    """Open the replication firewall ports on *pe* toward the peer cluster's PE +
    VIP (mirrors the modify_firewall_pe* null_resources)."""
    command = ("/usr/local/nutanix/cluster/bin/modify_firewall -f "
               f"-r {peer_pe},{peer_vip} -p {_DP_FIREWALL_PORTS} -i eth0")
    _ssh(pe, ssh_user, ssh_pass, command, timeout=180)


def _ncli_field_values(raw: str, field: str) -> list[str]:
    """Extract every value whose label (left of the first ':') is exactly *field*
    from ncli's aligned "Label : value" output (e.g. the container/pool "Name")."""
    values = []
    for line in raw.splitlines():
        label, sep, value = line.partition(":")
        if sep and label.strip() == field:
            value = value.strip()
            if value:
                values.append(value)
    return values


def _dp_list_container_names(pe: str, ssh_user: str, ssh_pass: str) -> list[str]:
    """Return the storage-container names on *pe* via `ncli container ls`."""
    raw = _ssh(pe, ssh_user, ssh_pass, f"{_NCLI} container ls", timeout=120)
    return _ncli_field_values(raw, "Name")


def _dp_first_storage_pool(pe: str, ssh_user: str, ssh_pass: str) -> str:
    """Return the first storage-pool name on *pe* via `ncli sp ls` ("" if none)."""
    names = _ncli_field_values(
        _ssh(pe, ssh_user, ssh_pass, f"{_NCLI} sp ls", timeout=120), "Name")
    return names[0] if names else ""


def _dp_ensure_sync_containers(local_pe: str, remote_pe: str,
                               ssh_user: str, ssh_pass: str) -> None:
    """Ensure sync-rep storage-container name parity between the two PEs.

    Synchronous replication (RPO=0) stretches a protected VM's disks into a
    container of the *same name* on the remote cluster; if that container is
    missing the EnableStretch task fails with kContainerNotFound and the VM
    never reaches RULE_PROTECTED. The promote-VM test places its VM in the
    source cluster's ``default-container-<id>`` (it selects the first
    ``default-container-*``), so mirror every such container name onto the
    remote PE, creating any that are missing in the remote's storage pool.
    """
    source = [n for n in _dp_list_container_names(local_pe, ssh_user, ssh_pass)
              if n.startswith("default-container-")]
    if not source:
        logger.warning("sync-rep container parity: no 'default-container-*' on local PE %s; "
                       "skipping", local_pe)
        return
    remote = set(_dp_list_container_names(remote_pe, ssh_user, ssh_pass))
    missing = [n for n in source if n not in remote]
    if not missing:
        logger.info("sync-rep container parity: remote PE %s already has %s",
                    remote_pe, ", ".join(source))
        return
    sp = _dp_first_storage_pool(remote_pe, ssh_user, ssh_pass)
    if not sp:
        logger.error("sync-rep container parity: no storage pool on remote PE %s; cannot "
                     "create %s", remote_pe, ", ".join(missing))
        return
    for name in missing:
        logger.info("sync-rep container parity: creating container %r on remote PE %s (sp=%s)",
                    name, remote_pe, sp)
        _ssh(remote_pe, ssh_user, ssh_pass,
             f"{_NCLI} container create name={name} sp-name={sp}", timeout=180)


def prepare_data_protection(client: PcClient, env: Env) -> bool:
    """Pre-config for the data-protection (near-sync/replication) tests, mirroring
    preEnv/data_protection.tf:

      0. Ensure the local PC is connected to the remote PC via an availability
         zone (availability_zone.remote_pc_ip); if the pairing is missing, create
         it with a v3 cloud_trust (ONPREM_CLOUD) and wait for it to come up.
      1. On both the local and remote *PCs*, set the near-sync cerebro test gflag
         and restart cerebro (genesis stop cerebro + cluster start), using the PC
         SSH creds ssh.pc_username / ssh.pc_password.
      1b. Do the same on both *PEs* (ssh.pe_* creds): the entity-centric NearSync
         DR capability is gated by the workload PE cerebro, so single-node /
         hybrid (non-all-flash) QA clusters need this override or a RPO<=15min
         (NearSync) entity never protects (DP-10400).
      2. Open the replication firewall ports (2030,2036,2073,2090,8740) between the
         two cluster *PEs*, in both directions, via modify_firewall, using the PE
         SSH creds ssh.pe_username / ssh.pe_password.

    The local PC is PC_ENDPOINT; the remote PC is availability_zone.remote_pc_ip.
    The local PE comes from the local PC's cluster (discover_local_cluster); the
    remote PE comes from the *remote PC's* cluster (discover_remote_pc_cluster) --
    NOT dp.remote_cluster_pe, which is a separate input for the tests. Requires
    sshpass.
    """
    local_pc = client.endpoint
    remote_pc = env.get("AZ_REMOTE_PC_IP")
    if not remote_pc:
        logger.error("Cannot prepare data protection: availability_zone.remote_pc_ip is not set "
                     "(it is the remote PC the cerebro step runs on and the AZ peer)")
        return False

    # 0. The two PCs must be paired via an availability zone. Connect if missing.
    logger.info("Ensuring availability-zone connection: local PC %s -> remote PC %s ...",
                local_pc, remote_pc)
    az_user = env.get("AZ_REMOTE_PC_USERNAME") or env.get("PC_USERNAME")
    az_pass = env.get("AZ_REMOTE_PC_PASSWORD") or env.get("PC_PASSWORD")
    if not _ensure_az_connected(client, remote_pc, az_user, az_pass):
        return False

    # PE svm_ips + VIPs for the firewall step (peer PE + VIP). The remote peer is
    # the cluster behind the AZ remote PC (not dp.remote_cluster_pe).
    local = discover_local_cluster(env)
    remote = discover_remote_pc_cluster(env)
    local_pe, local_vip = local.get("pe"), local.get("vip")
    remote_pe, remote_vip = remote.get("pe"), remote.get("vip")

    missing = [name for name, val in (
        ("local PE (local PC cluster)", local_pe), ("local VIP", local_vip),
        ("remote PE (remote PC cluster)", remote_pe), ("remote VIP", remote_vip),
    ) if not val]
    if missing:
        logger.error("Cannot prepare data protection: could not resolve %s for the firewall "
                     "step. Make sure the local PC and the remote PC (%s) are reachable so "
                     "their cluster PEs resolve, or pin the values (dp.local_cluster_pe / "
                     "dp.remote_pc_cluster_pe and their VIPs).",
                     ", ".join(missing), remote_pc)
        return False

    pc_user = env.get("SSH_PC_USERNAME")
    pc_pass = env.get("SSH_PC_PASSWORD")
    pe_user = env.get("SSH_PE_USERNAME")
    pe_pass = env.get("SSH_PE_PASSWORD")
    if not (pc_user and pc_pass):
        logger.error("Cannot prepare data protection: ssh.pc_username / ssh.pc_password must be "
                     "set (the cerebro step runs on the PCs)")
        return False
    if not (pe_user and pe_pass):
        logger.error("Cannot prepare data protection: ssh.pe_username / ssh.pe_password must be "
                     "set (the modify_firewall step runs on the PEs)")
        return False

    timeout_seconds = int(env.get("PC_PREP_DP_CEREBRO_TIMEOUT", str(_DP_CEREBRO_TIMEOUT))
                          or _DP_CEREBRO_TIMEOUT)

    logger.info("=== Preparing data protection ===")
    logger.info("local  pc=%s pe=%s vip=%s", local_pc, local_pe, local_vip)
    logger.info("remote pc=%s pe=%s vip=%s", remote_pc, remote_pe, remote_vip)

    if client.dry_run:
        logger.info("[dry-run] would SSH %s@{%s,%s} to set cerebro gflag '%s' + restart cerebro "
                    "(on the PCs), then SSH %s@{%s,%s} to set the same gflag + restart cerebro "
                    "(on the PEs, to enable the NearSync DR capability), then modify_firewall "
                    "opening ports %s in both directions (on the PEs), then ensure the source's "
                    "default-container-* names exist on remote PE %s (sync-rep container parity)",
                    pc_user, local_pc, remote_pc, _DP_CEREBRO_GFLAG,
                    pe_user, local_pe, remote_pe, _DP_FIREWALL_PORTS, remote_pe)
        return True

    # 1. cerebro gflag + restart on both PCs (local first, then remote).
    for label, pc in (("local", local_pc), ("remote", remote_pc)):
        logger.info("Setting near-sync cerebro gflag and restarting cerebro on %s PC %s ...",
                    label, pc)
        _dp_restart_cerebro(pc, pc_user, pc_pass, timeout_seconds)

    # 1b. The entity-centric NearSync DR *capability* is gated by the workload
    #     *PE* cerebro, not the PC. On single-node / hybrid (non-all-flash) QA
    #     clusters the PE otherwise logs "nearsync DR capability is not set" and
    #     never protects RPO<=15min (NearSync) entities, so the v4 protected-
    #     resource API returns DP-10400. Apply the same QA override gflag +
    #     restart cerebro on BOTH PEs (capability is negotiated bidirectionally).
    for label, pe in (("local", local_pe), ("remote", remote_pe)):
        logger.info("Setting near-sync cerebro gflag and restarting cerebro on %s PE %s ...",
                    label, pe)
        _dp_restart_cerebro(pe, pe_user, pe_pass, timeout_seconds)

    # 2. Open replication firewall ports in both directions on the PEs.
    logger.info("Opening replication ports on local PE %s -> peer %s,%s ...",
                local_pe, remote_pe, remote_vip)
    _dp_modify_firewall(local_pe, pe_user, pe_pass, remote_pe, remote_vip)
    logger.info("Opening replication ports on remote PE %s -> peer %s,%s ...",
                remote_pe, local_pe, local_vip)
    _dp_modify_firewall(remote_pe, pe_user, pe_pass, local_pe, local_vip)

    # 3. SyncRep (RPO=0) stretches disks into a same-named container on the remote
    #    cluster; ensure the source's default-container-* names exist on the remote
    #    PE, else EnableStretch fails with kContainerNotFound (VM stays UNPROTECTED).
    logger.info("Ensuring sync-rep storage-container parity: %s -> %s ...", local_pe, remote_pe)
    _dp_ensure_sync_containers(local_pe, remote_pe, pe_user, pe_pass)

    logger.info("Data protection pre-config complete.")
    return True


def prepare_prism(client: PcClient, env: Env) -> bool:
    """Pre-config for the prism (PC unregistration) tests: AZ-pair the local PC
    with prism.unregister.remote_pc_ip so the unregistration test has a connected
    PC to unregister. Creates the pairing via a v3 cloud_trust (ONPREM_CLOUD) when
    it isn't already connected. Remote PC creds default to
    prism.unregister.remote_pc_username / _password, then the local PC creds.
    No-op (success) when prism.unregister.remote_pc_ip is unset.
    """
    remote_pc = env.get("PRISM_UNREGISTER_REMOTE_PC_IP")
    if not remote_pc:
        logger.info("prism: prism.unregister.remote_pc_ip is not set; nothing to connect.")
        return True

    logger.info("=== Preparing prism (unregister) ===")
    logger.info("Ensuring availability-zone connection: local PC %s -> unregister PC %s ...",
                client.endpoint, remote_pc)
    user = env.get("PRISM_UNREGISTER_REMOTE_PC_USERNAME") or env.get("PC_USERNAME")
    passwd = env.get("PRISM_UNREGISTER_REMOTE_PC_PASSWORD") or env.get("PC_PASSWORD")
    return _ensure_az_connected(client, remote_pc, user, passwd)


# --------------------------------------------------------------------------- #
# Steps 8 & 9 shared: VM / image / subnet helpers
# (port of preEnv/scripts/create_vm_for_iscsi_clients.sh and
#  preEnv/scripts/create_vm_for_ngt_upgrade.sh)
# --------------------------------------------------------------------------- #
_VMS_V40 = "api/vmm/v4.0/ahv/config/vms"
_VMS_V41 = "api/vmm/v4.1/ahv/config/vms"
_IMAGES_V40 = "api/vmm/v4.0/content/images"
_IMAGES_V41 = "api/vmm/v4.1/content/images"
_VG_V40 = "api/volumes/v4.0/config/volume-groups"

def _iscsi_cloud_init_b64(user: str, password: str) -> str:
    """Build the cloud-init user-data (base64) that sets <user>:<password> and
    enables SSH password auth for the iSCSI-client VM. The credentials come from
    config (no password literal is baked into this script)."""
    user_data = (
        "#cloud-config\n"
        "chpasswd:\n"
        "  list: |\n"
        f"    {user}:{password}\n"
        "  expire: false\n"
        "disable_root: false\n"
        "ssh_pwauth:   true"
    )
    return base64.b64encode(user_data.encode()).decode()


def _extract_learned_ip(vm_data: object) -> str:
    """Return the first learned (or assigned) IPv4 of a VM's first NIC.

    Handles both the list shape (nics[].networkInfo) and the single-GET shape
    (nics[].nicNetworkInfo) returned by the v4.1 VM APIs."""
    if not isinstance(vm_data, dict):
        return ""
    nics = vm_data.get("nics") or []
    if not (isinstance(nics, list) and nics and isinstance(nics[0], dict)):
        return ""
    net = nics[0].get("nicNetworkInfo") or nics[0].get("networkInfo") or {}
    ipv4 = net.get("ipv4Info") or {}
    for key in ("learnedIpAddresses", "ipAddresses"):
        addrs = ipv4.get(key) or []
        if isinstance(addrs, list) and addrs and isinstance(addrs[0], dict):
            value = addrs[0].get("value")
            if value:
                return value
    return ""


def _get_vm_by_name(client: PcClient, name: str) -> dict | None:
    """Return the first VM whose name matches *name* (v4.1 $filter), or None."""
    fltr = urllib.parse.quote(f"name eq '{name}'", safe="")
    _, body = client.request("GET", f"{_VMS_V41}?$filter={fltr}", mutating=False)
    items = _dig(body, "data")
    if isinstance(items, list) and items and isinstance(items[0], dict):
        return items[0]
    return None


def _find_image_ext_id(client: PcClient, name: str, *, prefix: bool = False,
                       api_path: str = _IMAGES_V41) -> str | None:
    """Return the extId of the image named (or name-prefixed by) *name*."""
    _, body = client.request("GET", api_path, mutating=False)
    for img in (_dig(body, "data") or []):
        if not isinstance(img, dict):
            continue
        img_name = str(img.get("name") or "")
        if (img_name.startswith(name) if prefix else img_name == name):
            return img.get("extId")
    return None


def _ensure_image(client: PcClient, env: Env, name: str, url: str) -> str | None:
    """Ensure an image named *name* exists, creating it from *url* if missing.
    Returns the image extId (or a placeholder under --dry-run)."""
    existing = _find_image_ext_id(client, name)
    if existing:
        logger.info("Image %r already exists (extId=%s).", name, existing)
        return existing

    if client.dry_run:
        logger.info("[dry-run] would create image %r from %s and wait for the task", name, url)
        return "DRY-RUN-IMAGE-EXT-ID"

    logger.info("Creating image %r from %s ...", name, url)
    payload = {
        "name": name,
        "type": "DISK_IMAGE",
        "source": {
            "$objectType": "vmm.v4.content.UrlSource",
            "url": url,
            "shouldAllowInsecureUrl": True,
        },
    }
    status, body = client.request("POST", _IMAGES_V40, payload)
    task_ext_id = _dig(body, "data.extId", "extId")
    if not task_ext_id:
        logger.error("Image create returned no task extId (http %s): %s", status,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
        return None
    timeout = int(env.get("PC_PREP_IMAGE_TASK_TIMEOUT", "1800") or 1800)
    if not poll_task(client, task_ext_id, timeout_seconds=timeout):
        logger.error("Image %r create task did not succeed.", name)
        return None
    ext_id = _find_image_ext_id(client, name)
    if not ext_id:
        logger.error("Image %r not found after its create task completed.", name)
    return ext_id


def _discover_auto_cluster_ext_id(client: PcClient) -> str | None:
    """Return the extId of the first cluster whose name starts with 'auto'."""
    _, body = client.request("GET", "api/clustermgmt/v4.0/config/clusters", mutating=False)
    for cluster in (_dig(body, "data") or []):
        if isinstance(cluster, dict) and str(cluster.get("name") or "").startswith("auto"):
            return cluster.get("extId")
    return None


def _discover_vm_subnet_ext_id(client: PcClient) -> tuple[str | None, str | None]:
    """Return (extId, name) of the VM subnet, preferring 'vlan.800' then 'vlan.0'."""
    _, body = client.request("GET", "api/networking/v4.0/config/subnets", mutating=False)
    by_name = {}
    for subnet in (_dig(body, "data") or []):
        if isinstance(subnet, dict) and subnet.get("name"):
            by_name[subnet["name"]] = subnet.get("extId")
    for name in ("vlan.800", "vlan.0"):
        if by_name.get(name):
            return by_name[name], name
    return None, None


def _build_vm_payload(*, name: str, description: str, cluster_ext_id: str,
                      image_ext_id: str, subnet_ext_id: str,
                      cloud_init_b64: str | None = None,
                      with_cdrom: bool = False) -> dict:
    """Build a v4 AHV VM create payload (single boot disk from *image_ext_id*,
    single NIC on *subnet_ext_id*). Optionally attaches cloud-init user-data
    and/or an empty CD-ROM (for NGT ISO insertion)."""
    fv = {"$fv": "v4.r1"}
    payload: dict = {
        "$objectType": "vmm.v4.ahv.config.Vm",
        "$reserved": fv,
        "$unknownFields": {},
        "name": name,
        "description": description,
        "cluster": {
            "$objectType": "vmm.v4.ahv.config.ClusterReference",
            "$reserved": fv,
            "$unknownFields": {},
            "extId": cluster_ext_id,
        },
        "numSockets": 2,
        "numCoresPerSocket": 2,
        "memorySizeBytes": 4294967296,
        "isAgentVm": False,
        "hardwareClockTimezone": "UTC",
        "isMemoryOvercommitEnabled": False,
        "apcConfig": {
            "$objectType": "vmm.v4.ahv.config.ApcConfig",
            "$reserved": fv,
            "$unknownFields": {},
            "isApcEnabled": False,
        },
        "disks": [
            {
                "$objectType": "vmm.v4.ahv.config.Disk",
                "$reserved": fv,
                "$unknownFields": {},
                "diskAddress": {
                    "$objectType": "vmm.v4.ahv.config.DiskAddress",
                    "$reserved": fv,
                    "$unknownFields": {},
                    "busType": "SCSI",
                    "index": 0,
                },
                "backingInfo": {
                    "$objectType": "vmm.v4.ahv.config.VmDisk",
                    "$reserved": fv,
                    "$unknownFields": {},
                    "dataSource": {
                        "$objectType": "vmm.v4.ahv.config.DataSource",
                        "$reserved": fv,
                        "$unknownFields": {},
                        "reference": {
                            "imageExtId": image_ext_id,
                            "$objectType": "vmm.v4.ahv.config.ImageReference",
                            "$reserved": fv,
                            "$unknownFields": {},
                        },
                    },
                    "diskSizeBytes": 21474836480,
                },
            }
        ],
        "bootConfig": {
            "$objectType": "vmm.v4.ahv.config.UefiBoot",
            "$reserved": fv,
            "$unknownFields": {},
            "isSecureBootEnabled": False,
        },
        "vtpmConfig": {
            "$objectType": "vmm.v4.ahv.config.VtpmConfig",
            "$reserved": fv,
            "$unknownFields": {},
            "isVtpmEnabled": False,
        },
        "nics": [
            {
                "$objectType": "vmm.v4.ahv.config.Nic",
                "$reserved": fv,
                "$unknownFields": {},
                "nicNetworkInfo": {
                    "$objectType": "vmm.v4.ahv.config.VirtualEthernetNicNetworkInfo",
                    "$reserved": fv,
                    "$unknownFields": {},
                    "subnet": {
                        "$objectType": "vmm.v4.ahv.config.SubnetReference",
                        "$reserved": fv,
                        "$unknownFields": {},
                        "extId": subnet_ext_id,
                    },
                    "vlanMode": "ACCESS",
                },
                "nicBackingInfo": {
                    "$objectType": "vmm.v4.ahv.config.VirtualEthernetNic",
                    "$reserved": fv,
                    "$unknownFields": {},
                    "isConnected": True,
                },
            }
        ],
    }
    if cloud_init_b64:
        payload["guestCustomization"] = {
            "$objectType": "vmm.v4.ahv.config.GuestCustomizationParams",
            "$reserved": fv,
            "$unknownFields": {},
            "config": {
                "$objectType": "vmm.v4.ahv.config.CloudInit",
                "$reserved": fv,
                "$unknownFields": {},
                "datasourceType": "CONFIG_DRIVE_V2",
                "cloudInitScript": {
                    "$objectType": "vmm.v4.ahv.config.Userdata",
                    "$reserved": fv,
                    "$unknownFields": {},
                    "value": cloud_init_b64,
                },
            },
        }
    if with_cdrom:
        payload["cdRoms"] = [
            {
                "$objectType": "vmm.v4.ahv.config.Cdrom",
                "diskAddress": {"busType": "SATA", "index": 0},
            }
        ]
    return payload


def _create_vm(client: PcClient, payload: dict) -> bool:
    """POST a VM create payload and wait for the resulting task."""
    status, body = client.request("POST", _VMS_V40, payload)
    if client.dry_run:
        return True
    task_ext_id = _dig(body, "data.extId", "extId")
    if not task_ext_id:
        logger.error("VM create returned no task extId (http %s): %s", status,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
        return False
    return poll_task(client, task_ext_id, timeout_seconds=600)


def _power_on_vm(client: PcClient, vm_uuid: str) -> bool:
    """Power on *vm_uuid* (v4.1, requires If-Match) and wait for the task."""
    etag = client.get_etag(f"{_VMS_V41}/{vm_uuid}")
    if not etag:
        logger.error("Failed to fetch etag for VM %s (cannot power on).", vm_uuid)
        return False
    status, body = client.request(
        "POST", f"{_VMS_V41}/{vm_uuid}/$actions/power-on", None,
        extra_headers={"If-Match": etag})
    task_ext_id = _dig(body, "data.extId", "extId")
    if task_ext_id:
        return poll_task(client, task_ext_id, timeout_seconds=600)
    if status and status >= 400:
        logger.error("Power on VM %s failed (http %s): %s", vm_uuid, status,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
        return False
    return True


def _power_off_vm(client: PcClient, vm_uuid: str) -> bool:
    """Power off *vm_uuid* (v4.1, requires If-Match) and wait for the task.

    Best-effort: an already-powered-off VM (no task produced / 4xx) is treated as
    success so a subsequent delete can proceed."""
    etag = client.get_etag(f"{_VMS_V41}/{vm_uuid}")
    if not etag:
        logger.error("Failed to fetch etag for VM %s (cannot power off).", vm_uuid)
        return False
    status, body = client.request(
        "POST", f"{_VMS_V41}/{vm_uuid}/$actions/power-off", None,
        extra_headers={"If-Match": etag})
    task_ext_id = _dig(body, "data.extId", "extId")
    if task_ext_id:
        return poll_task(client, task_ext_id, timeout_seconds=600)
    if status and status >= 400:
        logger.info("Power off VM %s returned http %s (likely already off); continuing: %s",
                    vm_uuid, status,
                    _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
    return True


def _delete_vm(client: PcClient, vm_uuid: str) -> bool:
    """Power off (best-effort) then DELETE *vm_uuid* (v4.1, requires If-Match) and
    wait for the resulting task."""
    if client.dry_run:
        logger.info("[dry-run] would power off and delete VM %s", vm_uuid)
        return True
    _power_off_vm(client, vm_uuid)
    etag = client.get_etag(f"{_VMS_V41}/{vm_uuid}")
    if not etag:
        logger.error("Failed to fetch etag for VM %s (cannot delete).", vm_uuid)
        return False
    status, body = client.request(
        "DELETE", f"{_VMS_V41}/{vm_uuid}", None,
        extra_headers={"If-Match": etag})
    task_ext_id = _dig(body, "data.extId", "extId")
    if task_ext_id:
        return poll_task(client, task_ext_id, timeout_seconds=600)
    if status and status >= 400:
        logger.error("Delete VM %s failed (http %s): %s", vm_uuid, status,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
        return False
    return True


def _wait_for_vm_ip(client: PcClient, vm_uuid: str, *, timeout_seconds: int = 300,
                    interval_seconds: int = 5) -> str:
    """Poll a VM until its first NIC learns an IPv4 address (or timeout)."""
    logger.info("Waiting up to %ds for VM %s to acquire an IP ...", timeout_seconds, vm_uuid)
    deadline = time.monotonic() + timeout_seconds
    while time.monotonic() < deadline:
        _, body = client.request("GET", f"{_VMS_V41}/{vm_uuid}", mutating=False)
        ip = _extract_learned_ip(_dig(body, "data") or {})
        if ip:
            logger.info("VM %s acquired IP %s", vm_uuid, ip)
            return ip
        logger.debug("  no IP yet for VM %s; waiting ...", vm_uuid)
        time.sleep(interval_seconds)
    logger.warning("Timed out waiting for VM %s to acquire an IP.", vm_uuid)
    return ""


# --------------------------------------------------------------------------- #
# Step 8: iSCSI-client VM (port of create_vm_for_iscsi_clients.sh)
# --------------------------------------------------------------------------- #
_ISCSI_SETUP_SCRIPT = r"""
if sudo cat /etc/iscsi/initiatorname.iscsi 2>/dev/null | grep -q '^iqn\.'; then
    echo "iSCSI initiator already configured:"
    sudo cat /etc/iscsi/initiatorname.iscsi
    exit 0
fi
sudo apt-get update
sudo apt-get install open-iscsi -y
sudo systemctl enable open-iscsi --now
if sudo cat /etc/iscsi/initiatorname.iscsi 2>/dev/null | grep -q '^iqn\.'; then
    echo "iSCSI initiator configured:"
    sudo cat /etc/iscsi/initiatorname.iscsi
    exit 0
fi
sudo iscsi-iname | sudo tee /etc/iscsi/initiatorname.iscsi
sudo systemctl restart iscsid.service
sudo systemctl restart open-iscsi.service
echo "iSCSI initiator configured:"
sudo cat /etc/iscsi/initiatorname.iscsi
"""


def _save_iscsi_client_iqn(iqn: str) -> None:
    """Persist the client IQN to preEnv/results/pc_data.json (parity with the
    bash script; harmless if unused)."""
    results_file = REPO_ROOT / "preEnv" / "results" / "pc_data.json"
    data: dict = {}
    if results_file.exists():
        try:
            data = json.loads(results_file.read_text()) or {}
        except (ValueError, OSError):
            data = {}
    data["iscsi_client_iqn"] = iqn
    results_file.parent.mkdir(parents=True, exist_ok=True)
    results_file.write_text(json.dumps(data, indent=4, sort_keys=True) + "\n")
    logger.info("Saved iSCSI client IQN to %s", results_file)


def create_iscsi_client_vm(client: PcClient, env: Env) -> bool:
    """Ensure a Linux VM exists with a configured iSCSI initiator, and that its
    initiator is attached to a Volume Group -- so the volumesv2 iSCSI-client
    tests (which read the *first* registered iSCSI client) have one to use.

    Mirrors preEnv/scripts/create_vm_for_iscsi_clients.sh:
      1. If the VM already exists, we're done (idempotent).
      2. Ensure the ubuntu cloud image exists (create from URL otherwise).
      3. Create the VM (cloud-init sets the configured user/password), power it
         on, wait for IP.
      4. SSH in and configure open-iscsi; read the generated IQN.
      5. Create a Volume Group and attach the VM's iSCSI client to it.
    """
    logger.info("=== Creating iSCSI-client VM ===")
    vm_name = env.get("PC_PREP_ISCSI_VM_NAME") or "tf-vm-for-iscsi-clients"
    vg_name = env.get("PC_PREP_ISCSI_VOLUME_GROUP_NAME") or "vg_for_iscsi"
    vm_user = env.get("PC_PREP_ISCSI_VM_USERNAME") or "ubuntu"
    vm_pass = env.required("PC_PREP_ISCSI_VM_PASSWORD")
    image_name = env.get("IMAGES_UBUNTU_IMAGE") or "ubuntu-22.04-server-cloudimg-amd64.qcow2"
    image_url = env.get("IMAGES_UBUNTU_IMAGE_URL") or (
        "http://endor.dyn.nutanix.com/GoldImages/NuCalm/AHV-UVM-Images/"
        "ubuntu-22.04-server-cloudimg-amd64.qcow2")

    existing = _get_vm_by_name(client, vm_name)
    if existing:
        logger.info("VM %r already exists (extId=%s, ip=%s) -- iSCSI client already "
                    "provisioned; nothing to do.", vm_name, existing.get("extId"),
                    _extract_learned_ip(existing) or "N/A")
        return True

    image_ext_id = _ensure_image(client, env, image_name, image_url)
    if not image_ext_id:
        return False
    cluster_ext_id = _discover_auto_cluster_ext_id(client)
    if not cluster_ext_id:
        logger.error("No cluster found whose name starts with 'auto'.")
        return False
    subnet_ext_id, subnet_name = _discover_vm_subnet_ext_id(client)
    if not subnet_ext_id:
        logger.error("No subnet named 'vlan.800' or 'vlan.0' found for the VM NIC.")
        return False
    logger.info("Using cluster=%s subnet=%s (%s) image=%s", cluster_ext_id, subnet_ext_id,
                subnet_name, image_ext_id)

    payload = _build_vm_payload(
        name=vm_name, description="TF VM for iSCSI clients",
        cluster_ext_id=cluster_ext_id, image_ext_id=image_ext_id,
        subnet_ext_id=subnet_ext_id,
        cloud_init_b64=_iscsi_cloud_init_b64(vm_user, vm_pass))

    if client.dry_run:
        logger.info("[dry-run] would create VM %r, power it on, SSH in to configure "
                    "open-iscsi, then create VG %r and attach the iSCSI client.",
                    vm_name, vg_name)
        return True

    if not _create_vm(client, payload):
        return False
    vm = _get_vm_by_name(client, vm_name)
    vm_uuid = vm.get("extId") if vm else None
    if not vm_uuid:
        logger.error("Failed to find VM UUID for %r after creation.", vm_name)
        return False
    if not _power_on_vm(client, vm_uuid):
        return False

    vm_ip = _wait_for_vm_ip(client, vm_uuid)
    if not vm_ip:
        logger.error("VM %r never acquired an IP; cannot configure iSCSI.", vm_name)
        return False

    logger.info("Configuring open-iscsi on %s (user=%s) ...", vm_ip, vm_user)
    setup_out = _ssh(vm_ip, vm_user, vm_pass, _ISCSI_SETUP_SCRIPT, timeout=600)
    logger.debug("iSCSI setup output: %s", _truncate(setup_out))

    iqn = _ssh(vm_ip, vm_user, vm_pass,
               r"sudo cat /etc/iscsi/initiatorname.iscsi 2>/dev/null | grep -oP 'iqn\.[^\s]+'").strip()
    if not iqn:
        logger.error("Failed to read the iSCSI IQN from the VM.")
        return False
    logger.info("iSCSI IQN: %s", iqn)

    if not _ensure_volume_group_with_iscsi_client(client, cluster_ext_id, vg_name, iqn):
        return False

    _save_iscsi_client_iqn(iqn)
    logger.info("iSCSI-client VM ready: ip=%s creds=%s:<redacted> iqn=%s", vm_ip, vm_user, iqn)
    return True


def _ensure_volume_group_with_iscsi_client(client: PcClient, cluster_ext_id: str,
                                           vg_name: str, iqn: str) -> bool:
    """Create the Volume Group *vg_name* (if missing) and attach the iSCSI client
    *iqn* to it. Mirrors the VG-create + attach-iscsi-client part of the bash."""
    vg_ext_id = _get_vg_ext_id(client, vg_name)
    if vg_ext_id:
        logger.info("Volume Group %r already exists (extId=%s).", vg_name, vg_ext_id)
    else:
        logger.info("Creating Volume Group %r ...", vg_name)
        vg_payload = {
            "name": vg_name,
            "description": "",
            "shouldLoadBalanceVmAttachments": False,
            "sharingStatus": "NOT_SHARED",
            "iscsiFeatures": {
                "$reserved": {"$fv": "v4.r1"},
                "$objectType": "volumes.v4.config.IscsiFeatures",
                "$unknownFields": {},
            },
            "clusterReference": cluster_ext_id,
            "usageType": "USER",
            "isHidden": False,
            "$reserved": {"$fv": "v4.r1"},
            "$objectType": "volumes.v4.config.VolumeGroup",
            "$unknownFields": {},
        }
        status, body = client.request("POST", _VG_V40, vg_payload)
        task_ext_id = _dig(body, "data.extId", "extId")
        if task_ext_id:
            if not poll_task(client, task_ext_id, timeout_seconds=600):
                logger.error("Volume Group create task did not succeed.")
                return False
        elif status and status >= 400:
            logger.error("Volume Group create failed (http %s): %s", status,
                         _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
            return False
        vg_ext_id = _get_vg_ext_id(client, vg_name)
        if not vg_ext_id:
            logger.error("Failed to find Volume Group ext_id for %r after creation.", vg_name)
            return False
        logger.info("Volume Group %r extId=%s", vg_name, vg_ext_id)

    etag = client.get_etag(f"{_VG_V40}/{vg_ext_id}")
    if not etag:
        logger.error("Failed to fetch etag for Volume Group %s.", vg_ext_id)
        return False
    logger.info("Attaching iSCSI client %s to Volume Group %r ...", iqn, vg_name)
    iscsi_payload = {
        "iscsiInitiatorName": iqn,
        "enabledAuthentications": "NONE",
        "numVirtualTargets": 32,
        "$reserved": {"$fv": "v4.r1"},
        "$objectType": "volumes.v4.config.IscsiClient",
        "$unknownFields": {},
    }
    status, body = client.request(
        "POST", f"{_VG_V40}/{vg_ext_id}/$actions/attach-iscsi-client", iscsi_payload,
        extra_headers={"If-Match": etag})
    task_ext_id = _dig(body, "data.extId", "extId")
    if task_ext_id:
        if not poll_task(client, task_ext_id, timeout_seconds=600):
            logger.error("Attach iSCSI client task did not succeed.")
            return False
    else:
        text = json.dumps(body) if isinstance(body, (dict, list)) else str(body)
        if status and status >= 400 and "already" not in text.lower():
            logger.error("Attach iSCSI client failed (http %s): %s", status, _truncate(text))
            return False
    logger.info("iSCSI client attached to Volume Group %r.", vg_name)
    return True


def _get_vg_ext_id(client: PcClient, name: str) -> str | None:
    fltr = urllib.parse.quote(f"name eq '{name}'", safe="")
    _, body = client.request("GET", f"{_VG_V40}?$filter={fltr}", mutating=False)
    items = _dig(body, "data")
    if isinstance(items, list) and items and isinstance(items[0], dict):
        return items[0].get("extId")
    return None


# --------------------------------------------------------------------------- #
# Step 9: NGT-upgrade VM (port of create_vm_for_ngt_upgrade.sh)
# --------------------------------------------------------------------------- #
_GUEST_TOOLS_V42 = "api/vmm/v4.2/ahv/config/vms"


def _wait_for_ssh(host: str, user: str, password: str, *, timeout_seconds: int = 300,
                  interval_seconds: int = 5) -> bool:
    """Poll until SSH (password auth) succeeds on *host*, or timeout."""
    logger.info("Waiting up to %ds for SSH on %s ...", timeout_seconds, host)
    deadline = time.monotonic() + timeout_seconds
    while time.monotonic() < deadline:
        try:
            out = _ssh(host, user, password, "echo OpenSSH", timeout=15)
            if "OpenSSH" in out:
                logger.info("SSH is available on %s.", host)
                return True
        except (RuntimeError, subprocess.TimeoutExpired):
            pass
        time.sleep(interval_seconds)
    logger.error("Timed out waiting for SSH on %s.", host)
    return False


def _enable_ngt(client: PcClient, vm_uuid: str) -> bool:
    """Enable NGT on the VM via the v4.2 guest-tools API (requires If-Match)."""
    etag = client.get_etag(f"{_GUEST_TOOLS_V42}/{vm_uuid}/guest-tools")
    if not etag:
        logger.error("Failed to fetch guest-tools etag for VM %s.", vm_uuid)
        return False
    status, body = client.request(
        "PUT", f"{_GUEST_TOOLS_V42}/{vm_uuid}/guest-tools",
        {"isEnabled": True, "$objectType": "vmm.v4.ahv.config.GuestTools"},
        extra_headers={"If-Match": etag})
    task_ext_id = _dig(body, "data.extId", "extId")
    if task_ext_id:
        return poll_task(client, task_ext_id, timeout_seconds=600)
    if status and status >= 400:
        logger.error("Enable NGT failed (http %s): %s", status,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
        return False
    return True


def _insert_ngt_iso(client: PcClient, vm_uuid: str) -> bool:
    """Insert the NGT config ISO (v4.2 guest-tools action, config-only)."""
    etag = client.get_etag(f"{_GUEST_TOOLS_V42}/{vm_uuid}/guest-tools")
    if not etag:
        logger.error("Failed to fetch guest-tools etag for insert-iso on VM %s.", vm_uuid)
        return False
    status, body = client.request(
        "POST", f"{_GUEST_TOOLS_V42}/{vm_uuid}/guest-tools/$actions/insert-iso",
        {"isConfigOnly": True, "$objectType": "vmm.v4.ahv.config.GuestToolsInsertConfig"},
        extra_headers={"If-Match": etag})
    task_ext_id = _dig(body, "data.extId", "extId")
    if task_ext_id:
        return poll_task(client, task_ext_id, timeout_seconds=600)
    if status and status >= 400:
        logger.warning("Insert NGT ISO returned http %s: %s", status,
                       _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body)))
    return True


def _install_old_ngt_in_guest(host: str, user: str, password: str, *,
                              installer_url: str, force_reinstall: bool = False) -> bool:
    """Download + install the OLD NGT bundle in-guest so an upgrade is available.

    With *force_reinstall* (used when the VM already has a newer NGT installed and
    we want to downgrade it again), the currently-installed NGT is uninstalled
    first (best-effort) so the older bundle installs cleanly instead of being
    rejected as an equal/older version."""
    python_bin = _ssh(
        host, user, password,
        'PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" '
        "command -v python3 || command -v python || exit 1", timeout=30).strip()
    if not python_bin:
        logger.error("Failed to discover a Python interpreter on %s.", host)
        return False
    logger.info("Discovered Python interpreter on guest: %s", python_bin)

    logger.info("Downloading old NGT installer to guest ...")
    _ssh(host, user, password,
         f"sudo curl -k -L --fail '{installer_url}' -o /tmp/ngt_installer.tar.gz && "
         "sudo chmod 0644 /tmp/ngt_installer.tar.gz", timeout=600)
    logger.info("Extracting NGT installer ...")
    _ssh(host, user, password, "sudo tar -xzf /tmp/ngt_installer.tar.gz -C /tmp", timeout=120)

    logger.info("Installing lvm2 dependency ...")
    _ssh(host, user, password, (
        "if command -v dnf >/dev/null 2>&1; then sudo dnf -y install lvm2; "
        "elif command -v yum >/dev/null 2>&1; then sudo yum -y install lvm2; "
        "elif command -v apt-get >/dev/null 2>&1; then sudo apt-get update && sudo apt-get install -y lvm2; "
        "elif command -v zypper >/dev/null 2>&1; then sudo zypper --non-interactive install lvm2; "
        'else echo "No supported package manager found to install lvm2" >&2; exit 1; fi'), timeout=600)

    if force_reinstall:
        # A newer NGT is already installed; remove it first so the older bundle
        # can be installed (a straight install would otherwise be a no-op/refused).
        logger.info("Uninstalling the currently-installed NGT before downgrading ...")
        out = _ssh(host, user, password,
                   f"sudo {python_bin} /tmp/linux/install_ngt.py --operation uninstall "
                   "2>&1 || true", timeout=600)
        logger.debug("install_ngt.py uninstall output: %s", _truncate(out))

    logger.info("Running NGT install script ...")
    out = _ssh(host, user, password,
               f"sudo {python_bin} /tmp/linux/install_ngt.py --operation install", timeout=600)
    logger.debug("install_ngt.py output: %s", _truncate(out))
    _ssh(host, user, password,
         "sudo systemctl restart ngt_guest_agent.service ngt_self_service_restore.service "
         "2>/dev/null || true", timeout=60)
    return True


def _verify_ngt(client: PcClient, vm_uuid: str, host: str, user: str, password: str, *,
                expected_version: str, require_iso_inserted: bool,
                attempts: int = 12, delay_seconds: int = 20) -> bool:
    """Poll Prism (falling back to guest-tools + rpm) until NGT reports installed,
    reachable, enabled and at *expected_version* (prefix match)."""
    logger.info("Verifying NGT state (expect installed/reachable/enabled, version %s*) ...",
                expected_version)
    for attempt in range(1, attempts + 1):
        _, vm_body = client.request("GET", f"{_VMS_V41}/{vm_uuid}", mutating=False)
        _, gt_body = client.request(
            "GET", f"{_GUEST_TOOLS_V42}/{vm_uuid}/guest-tools", mutating=False)
        vm_data = _dig(vm_body, "data") or {}
        gt_data = _dig(gt_body, "data") or {}
        gt = vm_data.get("guestTools") or vm_data.get("guest_tools") or {}

        def _flag(name_camel: str, name_snake: str) -> bool:
            return bool(gt.get(name_camel, gt.get(name_snake, gt_data.get(name_camel,
                        gt_data.get(name_snake, False)))))

        installed = _flag("isInstalled", "is_installed")
        reachable = _flag("isReachable", "is_reachable")
        enabled = _flag("isEnabled", "is_enabled")
        iso_inserted = _flag("isIsoInserted", "is_iso_inserted")
        version = (gt.get("version") or gt_data.get("version") or "")
        if not version:
            try:
                version = _ssh(host, user, password,
                               "rpm -q --qf '%{VERSION}' nutanix-guest-agent 2>/dev/null || true",
                               timeout=30).strip()
            except (RuntimeError, subprocess.TimeoutExpired):
                version = ""

        if installed and reachable and enabled and version.startswith(expected_version):
            if require_iso_inserted and not iso_inserted:
                logger.info("NGT ISO was inserted earlier but is not attached now; continuing.")
            logger.info("NGT is installed, reachable, enabled and at version %s.", version)
            return True

        logger.info("[%02d/%02d] NGT not ready: installed=%s reachable=%s enabled=%s "
                    "iso=%s version=%s", attempt, attempts, installed, reachable, enabled,
                    iso_inserted, version or "(unknown)")
        time.sleep(delay_seconds)

    logger.error("NGT did not reach the expected state for VM %s.", vm_uuid)
    return False


def create_ngt_upgrade_vm(client: PcClient, env: Env) -> bool:
    """Ensure a VM (named vmm.ngt.ngt_upgrade_vm_name) exists with an OLD NGT
    version installed, so the vmmv2 NGT-upgrade test has a VM to upgrade.

    Mirrors preEnv/scripts/create_vm_for_ngt_upgrade.sh, with one addition:
      1. If the VM already exists, re-install (downgrade) the OLD NGT bundle
         in-guest so a fresh upgrade is available again (uninstalls the current
         NGT first), then verify -- rather than no-op'ing.
      2. Ensure the Rocky8 image exists (create from URL otherwise).
      3. Create the VM (with a CD-ROM for the NGT ISO), power it on, wait for IP.
      4. Enable NGT (v4.2 guest-tools) and insert the config ISO.
      5. SSH in, download + install the OLD NGT bundle, then verify NGT state.
    """
    logger.info("=== Creating NGT-upgrade VM ===")
    vm_name = (env.get("NGT_UPGRADE_VM_NAME") or env.get("VMM_NGT_NGT_UPGRADE_VM_NAME")
               or env.get("PC_PREP_NGT_VM_NAME"))
    image_name = env.get("PC_PREP_NGT_IMAGE_NAME") or env.get("IMAGES_NGT_IMAGE")
    image_url = env.get("PC_PREP_NGT_IMAGE_URL") or env.get("IMAGES_NGT_IMAGE_URL")
    vm_user = env.get("PC_PREP_NGT_VM_USERNAME")
    vm_pass = env.get("PC_PREP_NGT_VM_PASSWORD")
    ngt_version = env.get("PC_PREP_NGT_VERSION")
    installer_url = env.get("PC_PREP_NGT_INSTALLER_URL")

    missing = [label for label, val in (
        ("vmm.ngt.ngt_upgrade_vm_name / pc_prep.ngt.vm_name", vm_name),
        ("pc_prep.ngt.image_name / images.ngt_image", image_name),
        ("pc_prep.ngt.image_url / images.ngt_image_url", image_url),
        ("pc_prep.ngt.vm_username", vm_user),
        ("pc_prep.ngt.vm_password", vm_pass),
        ("pc_prep.ngt.version", ngt_version),
        ("pc_prep.ngt.installer_url", installer_url),
    ) if not val]
    if missing:
        logger.error("Cannot create NGT-upgrade VM: missing config value(s): %s",
                     "; ".join(missing))
        return False

    existing = _get_vm_by_name(client, vm_name)
    if existing:
        vm_uuid = existing.get("extId")
        vm_ip = _extract_learned_ip(existing)
        logger.info("VM %r already exists (extId=%s, ip=%s) -- deleting it and recreating "
                    "from scratch.", vm_name, vm_uuid, vm_ip or "N/A")
        if not vm_uuid:
            logger.error("Existing VM %r has no extId; cannot delete it.", vm_name)
            return False
        if client.dry_run:
            logger.info("[dry-run] would delete VM %r and recreate it (with CD-ROM), power it "
                        "on, enable NGT, insert the config ISO, then install NGT %s in-guest "
                        "and verify.", vm_name, ngt_version)
            return True
        if not _delete_vm(client, vm_uuid):
            logger.error("Failed to delete existing VM %r; cannot recreate it.", vm_name)
            return False
        logger.info("Deleted existing VM %r; recreating it.", vm_name)

    image_ext_id = _ensure_image(client, env, image_name, image_url)
    if not image_ext_id:
        return False
    cluster_ext_id = _discover_auto_cluster_ext_id(client)
    if not cluster_ext_id:
        logger.error("No cluster found whose name starts with 'auto'.")
        return False
    subnet_ext_id, subnet_name = _discover_vm_subnet_ext_id(client)
    if not subnet_ext_id:
        logger.error("No subnet named 'vlan.800' or 'vlan.0' found for the VM NIC.")
        return False
    logger.info("Using cluster=%s subnet=%s (%s) image=%s", cluster_ext_id, subnet_ext_id,
                subnet_name, image_ext_id)

    payload = _build_vm_payload(
        name=vm_name, description="TF VM for NGT upgrade test",
        cluster_ext_id=cluster_ext_id, image_ext_id=image_ext_id,
        subnet_ext_id=subnet_ext_id, with_cdrom=True)

    if client.dry_run:
        logger.info("[dry-run] would create VM %r (with CD-ROM), power it on, enable NGT, "
                    "insert the config ISO, then install NGT %s in-guest and verify.",
                    vm_name, ngt_version)
        return True

    if not _create_vm(client, payload):
        return False
    vm = _get_vm_by_name(client, vm_name)
    vm_uuid = vm.get("extId") if vm else None
    if not vm_uuid:
        logger.error("Failed to find VM UUID for %r after creation.", vm_name)
        return False
    if not _power_on_vm(client, vm_uuid):
        return False

    vm_ip = _wait_for_vm_ip(client, vm_uuid)
    if not vm_ip:
        logger.error("VM %r never acquired an IP; cannot install NGT.", vm_name)
        return False

    if not _enable_ngt(client, vm_uuid):
        return False

    # Only insert the ISO when the VM actually has a CD-ROM (avoids VMM-31100).
    _, vm_body = client.request("GET", f"{_VMS_V41}/{vm_uuid}", mutating=False)
    cdroms = _dig(vm_body, "data.cdRoms") or []
    require_iso_inserted = bool(isinstance(cdroms, list) and cdroms)
    if require_iso_inserted:
        if not _insert_ngt_iso(client, vm_uuid):
            return False
    else:
        logger.info("Skipping NGT ISO insert: VM has no CD-ROM device.")

    if not _wait_for_ssh(vm_ip, vm_user, vm_pass):
        return False
    if not _install_old_ngt_in_guest(vm_ip, vm_user, vm_pass, installer_url=installer_url):
        return False

    logger.info("Waiting 2 minutes for the CVM to detect NGT version %s ...", ngt_version)
    _sleep_with_progress(120, "Waiting for NGT detection")

    if not _verify_ngt(client, vm_uuid, vm_ip, vm_user, vm_pass,
                       expected_version=ngt_version, require_iso_inserted=require_iso_inserted):
        return False

    logger.info("NGT-upgrade VM ready: ip=%s creds=%s:<redacted> version=%s",
                vm_ip, vm_user, ngt_version)
    return True


# --------------------------------------------------------------------------- #
# Step: NKP node OS image
#
# Ensures the NKP node OS image (images.nkp_image) exists on PC, creating it from
# images.nkp_image_url if missing. The name must match pc_prep.nkp.os_image_name
# so testenv/e2e_nkp_deploy.py can boot the NKP control-plane/worker VMs from it.
# --------------------------------------------------------------------------- #
def ensure_nkp_image(client: PcClient, env: Env) -> bool:
    logger.info("=== Ensuring NKP node OS image ===")
    image_name = env.get("IMAGES_NKP_IMAGE") or env.get("PC_PREP_NKP_OS_IMAGE_NAME")
    image_url = env.get("IMAGES_NKP_IMAGE_URL")
    if not image_name:
        logger.error("Cannot ensure NKP image: set images.nkp_image "
                     "(or pc_prep.nkp.os_image_name).")
        return False
    if not image_url:
        existing = _find_image_ext_id(client, image_name)
        if existing:
            logger.info("NKP image %r already exists (extId=%s); no URL needed.",
                        image_name, existing)
            return True
        logger.error("NKP image %r not found and images.nkp_image_url is empty; "
                     "cannot create it.", image_name)
        return False
    return _ensure_image(client, env, image_name, image_url) is not None


# --------------------------------------------------------------------------- #
# Step: Self-Service (Calm) -- project/environment/policy + blueprints
#
# Mirrors preEnv/scripts/create_v3_project.sh (project + environment +
# credential + snapshot policy) and preEnv/scripts/upload_bp.sh (import + patch
# the two blueprints, then launch an app and run its snapshot action) using the
# same v3/calm APIs. Every sub-step is idempotent: existing entities are reused
# by name so the step can be re-run safely.
# --------------------------------------------------------------------------- #
_SS_V3 = "api/nutanix/v3"
_SS_CALM = "api/calm/v3.0"


def _ss_ok(status: int) -> bool:
    return bool(status) and 200 <= status < 300


def _resolve_repo_path(p: str) -> str:
    path = Path(p)
    return str(path if path.is_absolute() else (REPO_ROOT / p).resolve())


def _ss_body_msg(body: object) -> str:
    return _truncate(json.dumps(body) if isinstance(body, (dict, list)) else str(body))


def _ss_has_reason(body: object, *reasons: str) -> bool:
    wanted = {r.upper() for r in reasons}
    for msg in (_dig(body, "message_list") or []):
        if isinstance(msg, dict) and str(msg.get("reason") or "").upper() in wanted:
            return True
    return False


def _ensure_nucalm_ready(client: PcClient, env: Env) -> bool:
    """Guard: make sure Self-Service is enabled from the marketplace before we
    create projects/blueprints. Enables + waits if it is not already ready."""
    http, enablement, running = _nucalm_status(client)
    if enablement == "ENABLED" and running == "HEALTHY":
        logger.info("Self-Service already enabled (ENABLED/HEALTHY).")
        return True
    if client.dry_run:
        logger.info("[dry-run] Self-Service not confirmed ready; would enable and wait.")
        return True
    logger.info("Self-Service not ready (enablement=%s running=%s); enabling now ...",
                enablement or "(none)", running or "(none)")
    if not enable_nucalm(client, env):
        return False
    _, enablement, running = _nucalm_status(client)
    if enablement == "ENABLED" and running == "HEALTHY":
        return True
    if not enablement and not running:
        logger.warning("Could not confirm nucalm status via API; proceeding best-effort.")
        return True
    logger.error("Self-Service still not ready (enablement=%s running=%s).",
                 enablement or "(none)", running or "(none)")
    return False


# ---- infrastructure lookups (cluster / subnet / account / image) ---------- #
def _ss_resolve_cluster_uuid(client: PcClient) -> str | None:
    """First AOS cluster's UUID via clustermgmt v4 (falls back to an 'auto*'
    named cluster)."""
    fltr = urllib.parse.quote(
        "config/clusterFunction/any(t:t eq "
        "Clustermgmt.Config.ClusterFunctionRef'AOS')", safe="")
    _, body = client.request(
        "GET", f"api/clustermgmt/v4.1/config/clusters?$filter={fltr}", mutating=False)
    data = _dig(body, "data")
    if isinstance(data, list) and data:
        ext = _dig(data[0], "extId")
        if ext:
            return ext
    return _discover_auto_cluster_ext_id(client)


def _ss_resolve_subnet_uuid(client: PcClient, name: str) -> str | None:
    fltr = urllib.parse.quote(f"name eq '{name}'", safe="")
    _, body = client.request(
        "GET", f"api/networking/v4.0/config/subnets?$filter={fltr}", mutating=False)
    data = _dig(body, "data")
    if isinstance(data, list) and data:
        return _dig(data[0], "extId")
    return None


def _ss_resolve_account_uuid(client: PcClient) -> str | None:
    """PC-level account uuid for the NTNX_LOCAL_AZ account, used for the
    project/environment account references (mirrors create_v3_project.sh).

    This is the *PC* account (type nutanix_pc). NOTE: the AHV substrate's
    create_spec.resources.account_uuid must NOT use this value -- it needs the
    per-cluster Nutanix account (see _ss_resolve_substrate_account_uuid)."""
    _, body = client.request(
        "POST", f"{_SS_V3}/accounts/list",
        {"length": 250, "filter": "state!=DELETED;state!=DRAFT"})
    for ent in (_dig(body, "entities") or []):
        if not isinstance(ent, dict) or _dig(ent, "metadata.name") != "NTNX_LOCAL_AZ":
            continue
        acct = _dig(ent, "status.resources.data.pc_account_uuid",
                    "status.resources.pc_account_uuid")
        if acct:
            return acct
        lst = _dig(ent, "status.resources.data.cluster_account_reference_list")
        if isinstance(lst, list) and lst:
            acct = _dig(lst[0], "resources.data.pc_account_uuid")
            if acct:
                return acct
        # Last resort: the account entity's own uuid.
        return _dig(ent, "metadata.uuid")
    return None


def _ss_resolve_substrate_account_uuid(client: PcClient,
                                       cluster_uuid: str | None = None) -> str | None:
    """Account uuid for the AHV substrate's create_spec.resources.account_uuid.

    This must be the *per-cluster* Nutanix account entity uuid -- i.e. an entry's
    `.uuid` under NTNX_LOCAL_AZ.status.resources.data.cluster_account_reference_list
    (equivalently the metadata.uuid of the type=nutanix account for that cluster) --
    NOT the PC-level pc_account_uuid. Using the PC account here makes the blueprint
    launch hang in 'running' with application_uuid=null (the launch is accepted but
    the substrate can never be created). When cluster_uuid is given, prefer the
    reference whose cluster matches; otherwise use the first reference."""
    _, body = client.request(
        "POST", f"{_SS_V3}/accounts/list",
        {"length": 250, "filter": "state!=DELETED;state!=DRAFT"})
    fallback = None
    for ent in (_dig(body, "entities") or []):
        if not isinstance(ent, dict) or _dig(ent, "metadata.name") != "NTNX_LOCAL_AZ":
            continue
        for ref in (_dig(ent, "status.resources.data.cluster_account_reference_list") or []):
            if not isinstance(ref, dict):
                continue
            ref_uuid = ref.get("uuid")
            if not ref_uuid:
                continue
            if fallback is None:
                fallback = ref_uuid
            if cluster_uuid and _dig(ref, "resources.data.cluster_uuid") == cluster_uuid:
                return ref_uuid
        break
    return fallback


def _ss_resolve_image_uuid(client: PcClient, name: str) -> str | None:
    """v3 image metadata.uuid by name (the BP data_source_reference needs the v3
    uuid, not the v4 extId)."""
    _, body = client.request("POST", f"{_SS_V3}/images/list",
                             {"kind": "image", "length": 250})
    for ent in (_dig(body, "entities") or []):
        if not isinstance(ent, dict):
            continue
        nm = _dig(ent, "status.name", "spec.name", "metadata.name")
        if nm == name:
            return _dig(ent, "metadata.uuid")
    return None


# ---- project / environment / snapshot policy ------------------------------ #
def _ss_get_project_uuid(client: PcClient, name: str) -> str | None:
    for endpoint in ("projects_internal/list", "projects/list"):
        for payload in ({"kind": "project", "length": 250, "filter": f"name=={name}"},
                        {"kind": "project", "length": 250}):
            _, body = client.request("POST", f"{_SS_V3}/{endpoint}", payload)
            for ent in (_dig(body, "entities") or []):
                if not isinstance(ent, dict):
                    continue
                nm = _dig(ent, "spec.project_detail.name", "spec.name", "status.name")
                if nm == name:
                    uid = _dig(ent, "metadata.uuid")
                    if uid:
                        return uid
    return None


def _ss_wait_project_ready(client: PcClient, project_uuid: str, *,
                           attempts: int = 60, interval: int = 5) -> bool:
    for i in range(attempts):
        _, body = client.request(
            "GET", f"{_SS_V3}/projects_internal/{project_uuid}", mutating=False)
        state = str(_dig(body, "status.state") or "")
        logger.info("  project poll [%d/%d] state=%s", i + 1, attempts, state or "(none)")
        if state == "ERROR":
            logger.error("Project %s entered ERROR state.", project_uuid)
            return False
        if state and state != "PENDING":
            return True
        time.sleep(interval)
    logger.warning("Project %s did not leave PENDING; continuing.", project_uuid)
    return True


def _ss_create_project(client: PcClient, *, name: str, description: str,
                       cluster_uuid: str, subnet_uuid: str, subnet_name: str,
                       account_uuid: str) -> str | None:
    existing = _ss_get_project_uuid(client, name)
    if existing:
        logger.info("Project %r already exists (uuid=%s); reusing.", name, existing)
        return existing
    if client.dry_run:
        logger.info("[dry-run] would create project %r", name)
        return "DRY-RUN-PROJECT-UUID"

    payload = {
        "api_version": "3.0",
        "metadata": {"kind": "project"},
        "spec": {"project_detail": {"name": name, "description": description, "resources": {
            "default_subnet_reference": {"kind": "subnet", "uuid": subnet_uuid},
            "cluster_reference_list": [{"kind": "cluster", "uuid": cluster_uuid}],
            "subnet_reference_list": [{"kind": "subnet", "uuid": subnet_uuid, "name": subnet_name}],
            "external_network_list": [],
            "directory_reference_list": [],
            "account_reference_list": [{"kind": "account", "uuid": account_uuid}],
            "user_reference_list": [],
            "vpc_reference_list": [],
            "tunnel_reference_list": [],
            "identity_providers_reference_list": [],
            "external_user_group_reference_list": [],
            "enable_directory_and_identity_provider_shortlist": False,
            "environment_reference_list": [],
        }}},
    }
    status, body = client.request("POST", f"{_SS_V3}/projects_internal", payload)
    project_uuid = _dig(body, "metadata.uuid")
    if not project_uuid and _ss_has_reason(body, "DUPLICATE_ENTITY"):
        logger.info("Project %r reports DUPLICATE_ENTITY; resolving by name.", name)
        project_uuid = _ss_get_project_uuid(client, name)
    if not project_uuid:
        logger.error("Project create failed (http %s): %s", status, _ss_body_msg(body))
        return None
    if status == 202:
        _ss_wait_project_ready(client, project_uuid)
    return project_uuid


def _ss_get_env_uuid(client: PcClient, name: str, project_uuid: str) -> str | None:
    _, body = client.request("POST", f"{_SS_V3}/environments/list",
                             {"kind": "environment", "filter": f"name=={name}"})
    for ent in (_dig(body, "entities") or []):
        if not isinstance(ent, dict):
            continue
        nm = _dig(ent, "spec.name", "status.name")
        if nm == name and _dig(ent, "metadata.project_reference.uuid") == project_uuid:
            return _dig(ent, "metadata.uuid")
    return None


def _ss_create_environment(client: PcClient, *, name: str, project_uuid: str,
                           account_uuid: str, cluster_uuid: str, subnet_uuid: str,
                           cred_name: str, cred_username: str, cred_password: str,
                           cred_uuid: str) -> str | None:
    existing = _ss_get_env_uuid(client, name, project_uuid)
    if existing:
        logger.info("Environment %r already exists (uuid=%s); reusing.", name, existing)
        return existing
    if client.dry_run:
        logger.info("[dry-run] would create environment %r", name)
        return "DRY-RUN-ENV-UUID"

    payload = {
        "api_version": "3.0",
        "metadata": {"kind": "environment",
                     "project_reference": {"kind": "project", "uuid": project_uuid}},
        "spec": {"name": name, "description": "", "resources": {
            "substrate_definition_list": [],
            "credential_definition_list": [{
                "name": cred_name, "type": "PASSWORD", "cred_class": "static",
                "username": cred_username,
                "secret": {"attrs": {"is_secret_modified": True}, "value": cred_password},
                "uuid": cred_uuid,
            }],
            "infra_inclusion_list": [{
                "account_reference": {"kind": "account", "uuid": account_uuid},
                "vpc_references": [],
                "default_subnet_reference": {},
                "subnet_references": [{"uuid": subnet_uuid}],
                "type": "nutanix_pc",
                "cluster_references": [{"uuid": cluster_uuid}],
            }],
        }},
    }
    status, body = client.request("POST", f"{_SS_V3}/environments", payload)
    env_uuid = _dig(body, "metadata.uuid")
    if not env_uuid and _ss_has_reason(body, "DUPLICATE_ENTITY"):
        logger.info("Environment %r reports DUPLICATE_ENTITY; resolving by name.", name)
        env_uuid = _ss_get_env_uuid(client, name, project_uuid)
    if not env_uuid:
        logger.error("Environment create failed (http %s): %s", status, _ss_body_msg(body))
        return None
    return env_uuid


def _ss_update_project_env(client: PcClient, *, project_uuid: str, project_name: str,
                           description: str, account_uuid: str, env_uuid: str,
                           cluster_uuid: str, subnet_uuid: str, subnet_name: str,
                           owner_name: str = "admin",
                           owner_uuid: str = "00000000-0000-0000-0000-000000000000") -> bool:
    if client.dry_run:
        logger.info("[dry-run] would attach environment %s to project %s", env_uuid, project_uuid)
        return True
    _, proj = client.request(
        "GET", f"{_SS_V3}/projects_internal/{project_uuid}", mutating=False)
    try:
        spec_version = int(_dig(proj, "metadata.spec_version", "spec_version") or 0)
    except (TypeError, ValueError):
        spec_version = 0
    logger.info("Attaching environment (project spec_version=%d) ...", spec_version)
    payload = {
        "api_version": "3.1",
        "metadata": {
            "kind": "project", "uuid": project_uuid,
            "project_reference": {"kind": "project", "name": project_name, "uuid": project_uuid},
            "categories_mapping": {}, "categories": {},
            "spec_version": spec_version,
            "owner_reference": {"kind": "user", "name": owner_name, "uuid": owner_uuid},
        },
        "spec": {"project_detail": {"name": project_name, "description": description, "resources": {
            "default_environment_reference": {"kind": "environment", "uuid": env_uuid},
            "environment_reference_list": [{"kind": "environment", "uuid": env_uuid}],
            "external_network_list": [],
            "directory_reference_list": [],
            "account_reference_list": [{"kind": "account", "uuid": account_uuid}],
            "user_reference_list": [],
            "vpc_reference_list": [],
            "tunnel_reference_list": [],
            "identity_providers_reference_list": [],
            "external_user_group_reference_list": [],
            "default_subnet_reference": {"kind": "subnet", "uuid": subnet_uuid},
            "subnet_reference_list": [{"kind": "subnet", "name": subnet_name, "uuid": subnet_uuid}],
            "cluster_reference_list": [{"kind": "cluster", "uuid": cluster_uuid}],
            "enable_directory_and_identity_provider_shortlist": False,
        }}},
    }
    status, body = client.request(
        "PUT", f"{_SS_V3}/projects_internal/{project_uuid}", payload)
    if _ss_ok(status):
        return True
    logger.error("Project update failed (http %s): %s", status, _ss_body_msg(body))
    return False


def _ss_get_policy_uuid(client: PcClient, name: str, project_uuid: str) -> str | None:
    _, body = client.request("POST", f"{_SS_CALM}/app_protection_policies/list",
                             {"length": 250, "filter": f"name=={name}"})
    for ent in (_dig(body, "entities") or []):
        if not isinstance(ent, dict):
            continue
        nm = _dig(ent, "spec.name", "status.name")
        proj = _dig(ent, "metadata.project_reference.uuid") or ""
        if nm == name and (proj == project_uuid or not proj):
            return _dig(ent, "metadata.uuid")
    return None


def _ss_create_policy(client: PcClient, *, name: str, project_uuid: str,
                      project_name: str, env_uuid: str, account_uuid: str,
                      cluster_uuid: str) -> str | None:
    existing = _ss_get_policy_uuid(client, name, project_uuid)
    if existing:
        logger.info("Snapshot policy %r already exists (uuid=%s); reusing.", name, existing)
        return existing
    if client.dry_run:
        logger.info("[dry-run] would create snapshot policy %r", name)
        return "DRY-RUN-POLICY-UUID"

    policy_uuid = str(uuid.uuid4())
    rule_uuid = str(uuid.uuid4())
    rule_name = f"rule_{uuid.uuid4().hex[:8]}"
    payload = {
        "api_version": "3.0",
        "metadata": {"kind": "app_protection_policy", "uuid": policy_uuid,
                     "project_reference": {"kind": "project", "name": project_name,
                                           "uuid": project_uuid}},
        "spec": {"name": name, "description": "", "resources": {
            "is_default": True,
            "ordered_availability_site_list": [{
                "environment_reference": {"kind": "environment", "uuid": env_uuid},
                "infra_inclusion_list": {
                    "type": "nutanix_pc",
                    "account_reference": {"kind": "account", "uuid": account_uuid},
                    "cluster_references": [{"kind": "cluster", "uuid": cluster_uuid}],
                },
            }],
            "app_protection_rule_list": [{
                "name": rule_name, "enabled": True,
                "local_snapshot_retention_policy": {"snapshot_expiry_policy": {"multiple": 7}},
                "first_availability_site_index": 0,
                "second_availability_site_index": 0,
                "recovery_point_objective_secs": -1,
                "uuid": rule_uuid,
            }],
        }},
    }
    status, body = client.request(
        "POST", f"{_SS_CALM}/app_protection_policies", payload)
    created = _dig(body, "metadata.uuid")
    if not created and _ss_has_reason(body, "DUPLICATE_NAME", "DUPLICATE_ENTITY"):
        logger.info("Snapshot policy %r reports duplicate; resolving by name.", name)
        created = _ss_get_policy_uuid(client, name, project_uuid)
    if not created:
        logger.error("Snapshot policy create failed (http %s): %s", status, _ss_body_msg(body))
        return None
    return created


def _ss_write_state(client: PcClient, state: dict) -> None:
    if client.dry_run:
        return
    out_dir = SCRIPT_DIR / "results"
    out_dir.mkdir(parents=True, exist_ok=True)
    path = out_dir / "selfservice_v3_state.json"
    try:
        path.write_text(json.dumps(state, indent=2))
        logger.info("Wrote Self-Service state: %s", path)
    except OSError as exc:
        logger.warning("Could not write Self-Service state file: %s", exc)


# ---- blueprints (import + patch + launch + snapshot action) --------------- #
def _ss_import_blueprint(client: PcClient, file_path: str, name: str,
                         project_uuid: str) -> tuple[int, object]:
    """POST /blueprints/import_file as multipart/form-data (file/name/project)."""
    boundary = "----prepPC" + uuid.uuid4().hex
    file_bytes = Path(file_path).read_bytes()
    fname = Path(file_path).name
    parts: list[bytes] = []
    for key, value in (("name", name), ("project_uuid", project_uuid)):
        parts.append(
            f'--{boundary}\r\nContent-Disposition: form-data; name="{key}"\r\n\r\n'
            f'{value}\r\n'.encode())
    parts.append(
        f'--{boundary}\r\nContent-Disposition: form-data; name="file"; '
        f'filename="{fname}"\r\nContent-Type: application/json\r\n\r\n'.encode())
    parts.append(file_bytes + b"\r\n")
    parts.append(f"--{boundary}--\r\n".encode())
    data = b"".join(parts)

    full = client.url(f"{_SS_V3}/blueprints/import_file")
    req = urllib.request.Request(full, data=data, method="POST")
    req.add_header("Authorization", client.auth_header)
    req.add_header("Accept", "application/json")
    req.add_header("Content-Type", f"multipart/form-data; boundary={boundary}")
    try:
        with urllib.request.urlopen(req, context=client.ctx, timeout=180) as resp:
            raw = resp.read().decode()
            return getattr(resp, "status", resp.getcode()), _parse(raw)
    except urllib.error.HTTPError as exc:
        return exc.code, _parse(exc.read().decode(errors="replace"))
    except urllib.error.URLError as exc:
        raise RuntimeError(f"POST {full} -> connection error: {exc.reason}") from None


def _ss_lookup_bp_uuid(client: PcClient, name: str) -> str | None:
    _, body = client.request("POST", f"{_SS_V3}/blueprints/list",
                             {"filter": f"name=={name};state!=DELETED"})
    ents = _dig(body, "entities") or []
    if ents and isinstance(ents[0], dict):
        return _dig(ents[0], "metadata.uuid")
    return None


def _ss_fetch_bp(client: PcClient, bp_uuid: str) -> object:
    _, body = client.request("GET", f"{_SS_V3}/blueprints/{bp_uuid}", mutating=False)
    return body


def _ss_patch_blueprint_body(bp_body: dict, *, project_uuid: str, project_name: str,
                             env_uuid: str, cluster_uuid: str, image_uuid: str,
                             subnet_uuid: str, disk_size_mib: int,
                             bp_admin_password: str, policy_uuid: str,
                             policy_name: str) -> dict:
    """Port of the upload_bp.sh jq patch: wire the imported blueprint to the
    resolved project/env/cluster/image/subnet + snapshot-policy references."""
    body = json.loads(json.dumps(bp_body))  # deep copy
    body.pop("status", None)
    meta = body.setdefault("metadata", {})
    meta.pop("last_update_time", None)
    meta.pop("creation_time", None)
    meta["project_reference"] = {"kind": "project", "name": project_name, "uuid": project_uuid}

    res = body.setdefault("spec", {}).setdefault("resources", {})
    profiles = res.get("app_profile_list")
    if profiles:
        profiles[0]["environment_reference_list"] = [env_uuid]

    substrates = res.get("substrate_definition_list")
    if substrates:
        create_spec = substrates[0].setdefault("create_spec", {})
        create_spec["cluster_reference"] = {"kind": "cluster", "uuid": cluster_uuid}
        cs_res = create_spec.setdefault("resources", {})
        disks = cs_res.get("disk_list")
        if disks:
            disks[0]["data_source_reference"] = {"kind": "image", "uuid": image_uuid}
            disks[0]["disk_size_mib"] = int(disk_size_mib)
        nics = cs_res.get("nic_list")
        if nics:
            nics[0]["subnet_reference"] = {"kind": "subnet", "uuid": subnet_uuid}

    for cred in (res.get("credential_definition_list") or []):
        if cred.get("name") == "admin":
            secret = cred.setdefault("secret", {})
            secret["value"] = bp_admin_password
            secret.setdefault("attrs", {})["is_secret_modified"] = True

    if profiles:
        for snap_cfg in (profiles[0].get("snapshot_config_list") or []):
            for attr in (snap_cfg.get("attrs_list") or []):
                attr["app_protection_policy_reference"] = {
                    "kind": "app_protection_policy", "name": policy_name, "uuid": policy_uuid}
                attr.pop("app_protection_rule_reference", None)
    return body


def _ss_create_blueprint(client: PcClient, *, bp_file: str, bp_name: str,
                         project_uuid: str, project_name: str, env_uuid: str,
                         cluster_uuid: str, image_uuid: str, subnet_uuid: str,
                         disk_size_mib: int, bp_admin_password: str,
                         policy_uuid: str, policy_name: str) -> str | None:
    existing = _ss_lookup_bp_uuid(client, bp_name)
    if existing:
        logger.info("Blueprint %r already exists (uuid=%s); reusing and re-patching.",
                    bp_name, existing)
        bp_uuid = existing
        source = _ss_fetch_bp(client, bp_uuid)
    elif client.dry_run:
        logger.info("[dry-run] would import blueprint %r from %s and patch it",
                    bp_name, bp_file)
        return "DRY-RUN-BP-UUID"
    else:
        logger.info("Importing blueprint %r from %s ...", bp_name, Path(bp_file).name)
        status, up = _ss_import_blueprint(client, bp_file, bp_name, project_uuid)
        bp_uuid = _dig(up, "metadata.uuid")
        if not _ss_ok(status) or not bp_uuid:
            logger.error("Blueprint %r import failed (http %s): %s",
                         bp_name, status, _ss_body_msg(up))
            return None
        source = up

    if client.dry_run:
        return bp_uuid
    if not isinstance(source, dict):
        logger.error("Blueprint %r body could not be read for patching.", bp_name)
        return None

    patched = _ss_patch_blueprint_body(
        source, project_uuid=project_uuid, project_name=project_name, env_uuid=env_uuid,
        cluster_uuid=cluster_uuid, image_uuid=image_uuid, subnet_uuid=subnet_uuid,
        disk_size_mib=disk_size_mib, bp_admin_password=bp_admin_password,
        policy_uuid=policy_uuid, policy_name=policy_name)
    status, body = client.request("PUT", f"{_SS_V3}/blueprints/{bp_uuid}", patched)
    if _ss_ok(status):
        logger.info("Blueprint %r patched (uuid=%s state=%s).", bp_name, bp_uuid,
                    _dig(body, "status.state") or "unknown")
        return bp_uuid
    logger.error("Blueprint %r patch failed (http %s): %s", bp_name, status, _ss_body_msg(body))
    return None


def _ss_find_app_uuid(client: PcClient, app_name: str) -> str | None:
    _, body = client.request("POST", f"{_SS_V3}/apps/list", {"kind": "app", "length": 250})
    for ent in (_dig(body, "entities") or []):
        if not isinstance(ent, dict):
            continue
        nm = _dig(ent, "spec.name", "status.name", "metadata.name")
        state = str(_dig(ent, "status.state") or "").lower()
        if nm == app_name and state != "deleted":
            return _dig(ent, "metadata.uuid")
    return None


def _ss_wait_app_running(client: PcClient, app_uuid: str, *,
                         attempts: int = 60, interval: int = 15) -> bool:
    for i in range(attempts):
        _, body = client.request("GET", f"{_SS_V3}/apps/{app_uuid}", mutating=False)
        state = str(_dig(body, "status.state") or "").lower()
        logger.info("  app poll [%d/%d] state=%s", i + 1, attempts, state or "(none)")
        if state == "running":
            logger.info("App %s is running.", app_uuid)
            return True
        if state in ("error", "deleted", "not_deployed", "failed"):
            logger.error("App %s entered %s state.", app_uuid, state)
            return False
        time.sleep(interval)
    logger.error("Timed out waiting for app %s to run.", app_uuid)
    return False


def _ss_launch_app(client: PcClient, bp_body: dict, app_name: str,
                   app_description: str, account_uuid: str, *,
                   substrate_account_uuid: str | None = None,
                   attempts: int = 60, interval: int = 10) -> str | None:
    bp_uuid = _dig(bp_body, "metadata.uuid")
    profiles = _dig(bp_body, "spec.resources.app_profile_list") or [{}]
    app_profile_name = profiles[0].get("name")
    app_profile_uuid = profiles[0].get("uuid")
    if not (bp_uuid and app_profile_name and app_profile_uuid):
        logger.error("Cannot launch: blueprint/app-profile identifiers missing.")
        return None

    payload = json.loads(json.dumps(bp_body))
    payload.pop("status", None)
    meta = payload.setdefault("metadata", {})
    meta.pop("last_update_time", None)
    meta.pop("creation_time", None)
    meta["use_categories_mapping"] = False
    spec = payload.setdefault("spec", {})
    spec.pop("name", None)
    spec["description"] = app_description
    spec["application_name"] = app_name
    substrates = _dig(spec, "resources.substrate_definition_list") or []
    if substrates:
        # The AHV substrate needs the per-cluster Nutanix account entity uuid, not
        # the PC-level pc_account_uuid; otherwise the launch hangs in 'running'.
        substrates[0].setdefault("create_spec", {}).setdefault(
            "resources", {})["account_uuid"] = substrate_account_uuid or account_uuid
    spec["app_profile_reference"] = {"kind": "app_profile", "name": app_profile_name,
                                     "uuid": app_profile_uuid}

    logger.info("Launching app %r from blueprint %s ...", app_name, bp_uuid)
    status, body = client.request("POST", f"{_SS_V3}/blueprints/{bp_uuid}/launch", payload)
    if not _ss_ok(status):
        logger.error("Launch failed (http %s): %s", status, _ss_body_msg(body))
        return None
    request_id = _dig(body, "status.request_id")
    if not request_id:
        logger.error("Launch response missing request_id: %s", _ss_body_msg(body))
        return None
    logger.info("Launch accepted. request_id=%s", request_id)

    app_uuid = None
    for i in range(attempts):
        _, poll = client.request(
            "GET", f"{_SS_V3}/blueprints/{bp_uuid}/pending_launches/{request_id}",
            mutating=False)
        state = str(_dig(poll, "status.state") or "").lower()
        app_uuid = _dig(poll, "status.application_uuid")
        logger.info("  launch poll [%d/%d] state=%s app_uuid=%s", i + 1, attempts,
                    state or "(none)", app_uuid or "(none)")
        if state == "success":
            if not app_uuid:
                app_uuid = _ss_find_app_uuid(client, app_name)
            break
        if state in ("failed", "error", "not_deployed"):
            logger.error("App launch failed: state=%s", state)
            return None
        time.sleep(interval)

    if not app_uuid:
        logger.error("Launch succeeded but app uuid could not be resolved.")
        return None
    logger.info("App created uuid=%s", app_uuid)
    if not _ss_wait_app_running(client, app_uuid):
        return None
    return app_uuid


def _ss_run_snapshot_action(client: PcClient, app_uuid: str, action_name: str,
                            snapshot_name: str, *, attempts: int = 120,
                            interval: int = 5) -> bool:
    _, app = client.request("GET", f"{_SS_V3}/apps/{app_uuid}", mutating=False)
    api_version = (app.get("api_version") if isinstance(app, dict) else None) or "3.0"
    metadata = json.loads(json.dumps(_dig(app, "metadata") or {}))
    metadata.pop("owner_reference", None)
    metadata["uuid"] = str(uuid.uuid4())

    action_uuid = None
    task_uuid = None
    for act in (_dig(app, "status.resources.action_list") or []):
        if act.get("name") != action_name:
            continue
        action_uuid = act.get("uuid")
        for task in (_dig(act, "runbook.task_definition_list") or []):
            if task.get("type") == "CALL_CONFIG":
                task_uuid = task.get("uuid")
                break
        break
    if not (action_uuid and task_uuid):
        logger.error("Action %r or its CALL_CONFIG task not found on app %s.",
                     action_name, app_uuid)
        return False

    payload = {"api_version": api_version, "metadata": metadata, "spec": {
        "target_uuid": app_uuid, "target_kind": "Application",
        "args": [{"name": "snapshot_name", "value": snapshot_name, "task_uuid": task_uuid}],
    }}
    logger.info("Running action %r on app %s ...", action_name, app_uuid)
    status, body = client.request(
        "POST", f"{_SS_V3}/apps/{app_uuid}/actions/{action_uuid}/run", payload)
    if not _ss_ok(status):
        logger.error("Run action failed (http %s): %s", status, _ss_body_msg(body))
        return False
    runlog_uuid = _dig(body, "status.runlog_uuid", "runlog_uuid")
    if not runlog_uuid:
        logger.error("Action run missing runlog_uuid: %s", _ss_body_msg(body))
        return False
    logger.info("Action %r started. runlog=%s", action_name, runlog_uuid)

    for i in range(attempts):
        _, poll = client.request(
            "GET", f"{_SS_V3}/apps/{app_uuid}/app_runlogs/{runlog_uuid}/output",
            mutating=False)
        state = str(_dig(poll, "status.runlog_state") or "").upper()
        logger.info("  action poll [%d/%d] runlog_state=%s", i + 1, attempts, state or "(none)")
        if state in ("SUCCESS", "WARNING"):
            logger.info("Action %r completed: %s", action_name, state)
            return True
        if state in ("FAILURE", "ERROR", "SYS_FAILURE", "SYS_ERROR", "SYS_ABORTED",
                     "TIMEOUT", "APPROVAL_FAILED"):
            logger.error("Action %r failed: %s", action_name, state)
            return False
        time.sleep(interval)
    logger.error("Timed out waiting for action %r.", action_name)
    return False


def prepare_self_service(client: PcClient, env: Env) -> bool:
    logger.info("=== Preparing Self-Service (Calm): project + blueprints ===")

    # 1. Guard: Self-Service must be enabled from the marketplace first.
    if not _ensure_nucalm_ready(client, env):
        logger.error("Self-Service is not enabled/healthy; aborting. "
                     "Run 'prepare_pc.py --only nucalm' first.")
        return False

    def cfg(key: str, default: str = "") -> str:
        return env.get(f"PC_PREP_SELF_SERVICE_{key}", default) or default

    project_name = cfg("PROJECT_NAME", "tf-project-selfservice-v3")
    project_desc = cfg("PROJECT_DESCRIPTION")
    env_name = cfg("ENVIRONMENT_NAME", "tf-test-env")
    cred_name = cfg("CREDENTIAL_NAME", "tf_test_cred")
    cred_user = cfg("CREDENTIAL_USERNAME", "root")
    cred_pass = env.required("PC_PREP_SELF_SERVICE_CREDENTIAL_PASSWORD")
    policy_name = cfg("SNAPSHOT_POLICY_NAME", "test_local_snapshot_policy_local_account")
    subnet_name = cfg("SUBNET_NAME") or env.get("VMM_SUBNET_NAME") or "vlan.800"
    disk_size = int(cfg("BLUEPRINT_DISK_SIZE_MIB", "40960"))
    image_name = cfg("IMAGE_NAME") or env.get("IMAGES_NGT_IMAGE")
    image_url = cfg("IMAGE_URL") or env.get("IMAGES_NGT_IMAGE_URL")
    bp_admin_pass = cfg("BP_ADMIN_PASSWORD") or cred_pass
    launch = env.bool("PC_PREP_SELF_SERVICE_LAUNCH_APPS", True)
    bp1_file = _resolve_repo_path(cfg("BLUEPRINT1_FILE", "testenv/payloads/test_terraform_bp.json"))
    bp1_name = cfg("BLUEPRINT1_NAME", "test_terraform_bp")
    bp2_file = _resolve_repo_path(
        cfg("BLUEPRINT2_FILE", "testenv/payloads/test_terraform_bp_with_snapshot_config.json"))
    bp2_name = cfg("BLUEPRINT2_NAME", "test_terraform_bp_with_snapshot_config")
    bp2_app_name = cfg("BP2_APP_NAME", "test_terraform_snapshot_restore_app")
    bp2_app_desc = cfg("BP2_APP_DESCRIPTION")
    bp2_action = cfg("BP2_ACTION_NAME", "Snapshot_s1")

    for bp_file in (bp1_file, bp2_file):
        if not Path(bp_file).exists():
            logger.error("Blueprint file not found: %s", bp_file)
            return False

    # 2. Resolve infrastructure references.
    cluster_uuid = _ss_resolve_cluster_uuid(client)
    subnet_uuid = _ss_resolve_subnet_uuid(client, subnet_name)
    account_uuid = _ss_resolve_account_uuid(client)
    if client.dry_run:
        cluster_uuid = cluster_uuid or "DRY-RUN-CLUSTER"
        subnet_uuid = subnet_uuid or "DRY-RUN-SUBNET"
        account_uuid = account_uuid or "DRY-RUN-ACCOUNT"
    for label, value in (("cluster", cluster_uuid), ("subnet", subnet_uuid),
                         ("account", account_uuid)):
        if not value:
            logger.error("Could not resolve %s uuid; aborting.", label)
            return False
    logger.info("Resolved cluster=%s subnet(%s)=%s account=%s",
                cluster_uuid, subnet_name, subnet_uuid, account_uuid)

    # 3. Project + environment + snapshot policy.
    project_uuid = _ss_create_project(
        client, name=project_name, description=project_desc, cluster_uuid=cluster_uuid,
        subnet_uuid=subnet_uuid, subnet_name=subnet_name, account_uuid=account_uuid)
    if not project_uuid:
        return False
    cred_uuid = str(uuid.uuid4())
    env_uuid = _ss_create_environment(
        client, name=env_name, project_uuid=project_uuid, account_uuid=account_uuid,
        cluster_uuid=cluster_uuid, subnet_uuid=subnet_uuid, cred_name=cred_name,
        cred_username=cred_user, cred_password=cred_pass, cred_uuid=cred_uuid)
    if not env_uuid:
        return False
    if not _ss_update_project_env(
            client, project_uuid=project_uuid, project_name=project_name,
            description=project_desc, account_uuid=account_uuid, env_uuid=env_uuid,
            cluster_uuid=cluster_uuid, subnet_uuid=subnet_uuid, subnet_name=subnet_name):
        return False
    policy_uuid = _ss_create_policy(
        client, name=policy_name, project_uuid=project_uuid, project_name=project_name,
        env_uuid=env_uuid, account_uuid=account_uuid, cluster_uuid=cluster_uuid)
    if not policy_uuid:
        return False

    _ss_write_state(client, {
        "project": {"name": project_name, "uuid": project_uuid},
        "environment": {"name": env_name, "uuid": env_uuid},
        "credential": {"name": cred_name, "username": cred_user,
                       "password": cred_pass, "uuid": cred_uuid},
        "policy": {"name": policy_name, "uuid": policy_uuid},
        "infrastructure": {"subnet_name": subnet_name, "subnet_uuid": subnet_uuid,
                           "cluster_uuid": cluster_uuid, "account_uuid": account_uuid},
    })

    # 4. Resolve (or create) the guest image the blueprints deploy from.
    image_uuid = _ss_resolve_image_uuid(client, image_name) if image_name else None
    if not image_uuid and image_name and image_url and not client.dry_run:
        logger.info("Image %r not found; creating it from %s ...", image_name, image_url)
        _ensure_image(client, env, image_name, image_url)
        image_uuid = _ss_resolve_image_uuid(client, image_name)
    if client.dry_run:
        image_uuid = image_uuid or "DRY-RUN-IMAGE"
    if not image_uuid:
        logger.error("Could not resolve blueprint image uuid (name=%r). Set "
                     "pc_prep.self_service.image_name/image_url in config.yaml.", image_name)
        return False
    logger.info("Using blueprint image %r uuid=%s", image_name, image_uuid)

    # 5. Import + patch both blueprints.
    ok = True
    bp_common = dict(
        project_uuid=project_uuid, project_name=project_name, env_uuid=env_uuid,
        cluster_uuid=cluster_uuid, image_uuid=image_uuid, subnet_uuid=subnet_uuid,
        disk_size_mib=disk_size, bp_admin_password=bp_admin_pass,
        policy_uuid=policy_uuid, policy_name=policy_name)
    if not _ss_create_blueprint(client, bp_file=bp1_file, bp_name=bp1_name, **bp_common):
        ok = False
    bp2_uuid = _ss_create_blueprint(client, bp_file=bp2_file, bp_name=bp2_name, **bp_common)
    if not bp2_uuid:
        ok = False

    # 6. Launch an app from the snapshot blueprint and run its snapshot action.
    #    The app is launched with the exact bp2_app_name (no timestamp suffix) so
    #    the self-service acceptance tests -- which look the app up by that exact
    #    name (filter name==test_terraform_snapshot_restore_app) -- can find it.
    #    If an app with that name already exists (and is not deleted), reuse it
    #    instead of launching a duplicate.
    if launch and bp2_uuid and not client.dry_run:
        app_name = bp2_app_name
        snapshot_name = f"snapshot-{time.strftime('%Y%m%d-%H%M%S')}"
        app_uuid = _ss_find_app_uuid(client, app_name)
        if app_uuid:
            logger.info("App %r already exists (uuid=%s); reusing it.", app_name, app_uuid)
        else:
            # The substrate's account must be the per-cluster Nutanix account
            # entity (not the PC pc_account_uuid) or the launch hangs forever.
            substrate_account_uuid = (
                _ss_resolve_substrate_account_uuid(client, cluster_uuid) or account_uuid)
            logger.info("Using substrate account_uuid=%s (cluster=%s)",
                        substrate_account_uuid, cluster_uuid)
            bp2_body = _ss_fetch_bp(client, bp2_uuid)
            if isinstance(bp2_body, dict):
                app_uuid = _ss_launch_app(client, bp2_body, app_name, bp2_app_desc,
                                          account_uuid,
                                          substrate_account_uuid=substrate_account_uuid)
            else:
                logger.warning("Could not fetch BP2 body for launch; skipping app launch.")
        if app_uuid:
            if not _ss_run_snapshot_action(client, app_uuid, bp2_action, snapshot_name):
                logger.warning("Snapshot action did not complete successfully.")
        else:
            logger.warning("BP2 app launch did not complete; snapshot action skipped.")
    elif launch and client.dry_run:
        logger.info("[dry-run] would launch an app named %r from %r and run action %r",
                    bp2_app_name, bp2_name, bp2_action)

    if ok:
        logger.info("Self-Service preparation completed.")
    return ok


# --------------------------------------------------------------------------- #
# Flow Controller: enable (deploy) the FLOW_CONTROLLER product
# --------------------------------------------------------------------------- #
_DOMAIN_MGR_V44 = "api/prism/v4.4/management/domain-managers"


def _discover_domain_manager_id(client: PcClient, env: Env) -> str | None:
    """Return the domain manager (Prism Central) extId.

    Resolution order:
      1. pc_prep.flow.domain_manager_id (env PC_PREP_FLOW_DOMAIN_MANAGER_ID).
      2. First entry from GET /management/domain-managers.
      3. Fallback: the clustermgmt cluster whose function is PRISM_CENTRAL
         (its extId is the domain manager extId). Some builds 404 the
         domain-managers collection GET, so this keeps discovery working.
    """
    override = env.get("PC_PREP_FLOW_DOMAIN_MANAGER_ID")
    if override:
        return override

    _, body = client.request("GET", _DOMAIN_MGR_V44, mutating=False)
    managers = _dig(body, "data") or []
    if isinstance(managers, list) and managers and isinstance(managers[0], dict):
        dm_id = managers[0].get("extId")
        if dm_id:
            return dm_id

    logger.info("%s did not return a domain manager -- falling back to the "
                "PRISM_CENTRAL cluster extId.", _DOMAIN_MGR_V44)
    dm_id = _discover_pc_cluster_ext_id(client)
    if not dm_id:
        logger.error("Could not resolve the domain manager (PC) extId. Set "
                     "pc_prep.flow.domain_manager_id in config.yaml.")
    return dm_id


_SUBNETS_V40 = "api/networking/v4.0/config/subnets"


def _discover_pe_cluster_ext_id(client: PcClient) -> str | None:
    """Return the extId of a Prism Element (AHV) cluster to host the subnet.

    The Flow Controller subnet must live on the underlying PE cluster, not on the
    Prism Central cluster. Preference order: a cluster whose name starts with
    'auto' (the testenv convention), then the first cluster whose function is not
    PRISM_CENTRAL."""
    _, body = client.request("GET", "api/clustermgmt/v4.0/config/clusters", mutating=False)
    clusters = [c for c in (_dig(body, "data") or []) if isinstance(c, dict)]
    for cluster in clusters:
        if str(cluster.get("name") or "").startswith("auto"):
            return cluster.get("extId")
    for cluster in clusters:
        functions = _dig(cluster, "config.clusterFunction", "config.clusterFunctions") or []
        if isinstance(functions, str):
            functions = [functions]
        if "PRISM_CENTRAL" not in functions:
            return cluster.get("extId")
    return None


def _find_subnet_ext_id_by_name(client: PcClient, name: str) -> str | None:
    """Return the extId of the subnet named *name*, or None if absent."""
    fltr = urllib.parse.quote(f"name eq '{name}'", safe="")
    _, body = client.request("GET", f"{_SUBNETS_V40}?$filter={fltr}", mutating=False)
    for subnet in (_dig(body, "data") or []):
        if isinstance(subnet, dict) and subnet.get("name") == name:
            return subnet.get("extId")
    return None


def ensure_flow_managed_subnet(client: PcClient, env: Env) -> tuple[str | None, str | None]:
    """Resolve the managed VLAN subnet the Flow Controller (SMSP) deploys into.

    Flow Controller needs at least 9 free IPs handed out via DHCP, so it must sit
    on a managed VLAN subnet. AHV basic networking allows only ONE subnet per VLAN
    id, and the object store already owns VLAN 800 (objects.800, created by
    testenv/terraform/objects.tf), so the Flow Controller REUSES that same subnet.
    pc_prep.flow.subnet.name therefore normally points at objects.800 and this
    function just looks it up by name and returns its extId.

    Configured via pc_prep.flow.subnet.* in config.yaml (env PC_PREP_FLOW_SUBNET_*):
      name (required) and an optional cluster_ext_id override. If the named subnet
      does NOT already exist, it is created from the shared object-store subnet spec
      (object_store.subnet.*: vlan_id, default_gateway_ip, dhcp_server_address,
      ip_subnet, prefix_length, range_ip_pool_start, range_ip_pool_end; DNS from
      dns.servers) -- each flow key falls back to the matching OBJECT_STORE_SUBNET_*
      value. This keeps a single source of truth for the shared VLAN-800 subnet.

    Returns ``(subnet_ext_id, cluster_ext_id)``:
      * ``(extId, clusterExtId)`` on success (found or created);
      * ``("", None)`` when no subnet name is configured -- the caller should fall
        back to the FLOW_CONTROLLER product's existing metadata;
      * ``(None, None)`` on hard failure.
    """
    def cfg(key: str, default: str = "") -> str:
        # Prefer the flow subnet's own config; fall back to the object-store subnet
        # spec, since the Flow Controller shares that VLAN-800 subnet (objects.800)
        # and its create-time spec lives under object_store.subnet.*.
        return (env.get(f"PC_PREP_FLOW_SUBNET_{key}")
                or env.get(f"OBJECT_STORE_SUBNET_{key}", default) or default)

    name = cfg("NAME")
    if not name:
        logger.info("No pc_prep.flow.subnet.name configured -- skipping managed "
                    "subnet creation (will fall back to the product's subnetExtId).")
        return "", None

    cluster_ext_id = cfg("CLUSTER_EXT_ID") or _discover_pe_cluster_ext_id(client)
    if not cluster_ext_id and client.dry_run:
        cluster_ext_id = "DRY-RUN-CLUSTER"

    logger.info("=== Ensuring Flow Controller managed subnet %r ===", name)

    existing = _find_subnet_ext_id_by_name(client, name)
    if existing:
        logger.info("Managed subnet %r already exists (extId=%s) -- reusing it.",
                    name, existing)
        return existing, cluster_ext_id

    vlan_id = cfg("VLAN_ID")
    default_gateway_ip = cfg("DEFAULT_GATEWAY_IP")
    dhcp_server_address = cfg("DHCP_SERVER_ADDRESS")
    ip_subnet = cfg("IP_SUBNET")
    prefix_length = cfg("PREFIX_LENGTH")
    pool_start_ip = cfg("RANGE_IP_POOL_START")
    pool_end_ip = cfg("RANGE_IP_POOL_END")
    missing = [label for label, value in (
        ("vlan_id", vlan_id), ("ip_subnet", ip_subnet),
        ("prefix_length", prefix_length), ("default_gateway_ip", default_gateway_ip),
        ("range_ip_pool_start", pool_start_ip), ("range_ip_pool_end", pool_end_ip),
    ) if not value]
    if missing:
        # The Flow Controller shares the object-store VLAN-800 subnet (objects.800);
        # its create-time spec normally comes from object_store.subnet.*. If both the
        # flow and object-store specs are missing a field we cannot create it here.
        logger.error("Managed subnet %r not found and cannot be created -- missing "
                     "create params (%s). Set them under object_store.subnet.* (or "
                     "pc_prep.flow.subnet.*), or apply testenv/terraform first.", name,
                     ", ".join(missing))
        return None, None

    if not cluster_ext_id:
        logger.error("Could not resolve a PE cluster extId for the managed subnet; "
                     "set pc_prep.flow.subnet.cluster_ext_id.")
        return None, None

    def _ipv4(value: str) -> dict:
        return {"$objectType": "common.v1.config.IPv4Address", "value": value}

    def _ipaddr(value: str) -> dict:
        # IPAddress wrapper (ipv4 nested) -- required by fields typed as
        # common.v1.config.IPAddress (e.g. dhcpOptions.domainNameServers), as
        # opposed to the flat IPv4Address used by gateway/subnet/pool fields.
        return {"$objectType": "common.v1.config.IPAddress", "ipv4": _ipv4(value)}

    ipv4_cfg = {
        "$objectType": "networking.v4.config.IPv4Config",
        "ipSubnet": {
            "$objectType": "networking.v4.config.IPv4Subnet",
            "ip": _ipv4(ip_subnet),
            "prefixLength": int(prefix_length),
        },
        "defaultGatewayIp": _ipv4(default_gateway_ip),
        "poolList": [
            {
                "$objectType": "networking.v4.config.IPv4Pool",
                "startIp": _ipv4(pool_start_ip),
                "endIp": _ipv4(pool_end_ip),
            }
        ],
    }
    if dhcp_server_address:
        ipv4_cfg["dhcpServerAddress"] = _ipv4(dhcp_server_address)

    payload = {
        "$objectType": "networking.v4.config.Subnet",
        "$reserved": {"$fv": "v4.r3"},
        "name": name,
        "description": f"Managed VLAN subnet {name} (object store + Flow Controller)",
        "subnetType": "VLAN",
        "networkId": int(vlan_id),
        "clusterReference": cluster_ext_id,
        "ipConfig": [
            {
                "$objectType": "networking.v4.config.IPConfig",
                "ipv4": ipv4_cfg,
            }
        ],
    }
    dns_servers = env.list("DNS_SERVERS")
    if dns_servers:
        payload["dhcpOptions"] = {
            "$objectType": "networking.v4.config.DhcpOptions",
            "domainNameServers": [_ipaddr(ip) for ip in dns_servers],
        }

    if client.dry_run:
        logger.info("[dry-run] would POST %s to create managed subnet %r "
                    "(vlan=%s subnet=%s/%s gateway=%s pool=%s-%s cluster=%s).",
                    _SUBNETS_V40, name, vlan_id, ip_subnet, prefix_length,
                    default_gateway_ip, pool_start_ip, pool_end_ip, cluster_ext_id)
        return "DRY-RUN-SUBNET", cluster_ext_id

    status_code, body = client.request("POST", _SUBNETS_V40, payload)
    if status_code and status_code >= 400:
        logger.error("Creating managed subnet %r failed (http %s): %s", name, status_code,
                     _truncate(json.dumps(body) if isinstance(body, (dict, list))
                               else str(body)))
        return None, None

    task_ext_id = _dig(body, "data.extId", "extId")
    if task_ext_id:
        timeout = int(env.get("PC_PREP_FLOW_TASK_TIMEOUT", "1800") or 1800)
        if not poll_task(client, task_ext_id, timeout_seconds=timeout):
            logger.error("Managed subnet %r create task did not succeed.", name)
            return None, None
    else:
        logger.warning("No task extId in subnet-create response (http %s) -- assuming "
                       "applied synchronously.", status_code)

    subnet_ext_id = _find_subnet_ext_id_by_name(client, name)
    if not subnet_ext_id:
        logger.error("Managed subnet %r was created but its extId could not be "
                     "resolved by name.", name)
        return None, None

    logger.info("Managed subnet %r created (extId=%s, vlan=%s, pool %s-%s).", name,
                subnet_ext_id, vlan_id, pool_start_ip, pool_end_ip)
    return subnet_ext_id, cluster_ext_id


def enable_flow_controller(client: PcClient, env: Env) -> bool:
    """Enable (deploy) the FLOW_CONTROLLER product on Prism Central.

    This drives the Prism v4.4 product-management enable flow: it flips the
    FLOW_CONTROLLER product's enablementState to ENABLED, referencing the managed
    VLAN subnet the Flow Controller (SMSP) deploys into. It replaces the older
    "flex mode toggle" behaviour -- the goal here is to bring the Flow Controller
    up from scratch, not just tweak an already-deployed one.

    Steps (mirror the manual curl flow):
      0. Network prep: resolve the managed VLAN subnet the Flow Controller needs
         (>= 9 DHCP IPs on PC's VLAN) and capture its extId. The subnet is shared
         with the object store (objects.800) -- see ensure_flow_managed_subnet.
      1. Resolve the domain manager (PC) extId.
      2. Find the FLOW_CONTROLLER product; read its extId (and fall back to its
         metadata for cluster/subnet when the subnet step was skipped).
      3. If it is already ENABLED, no-op (idempotent).
      4. PUT the product with enablementState=ENABLED and FlowControllerMetadata
         (clusterExtId, subnetExtId, isFNSFlexModeEnabled=false,
         isSecurityAnalyticsEnabled=true), using the product's ETag as If-Match,
         then poll the deployment task.
    """
    # Phase 1: network preparation. The Flow Controller cannot deploy onto an
    # unmanaged subnet (it needs 9 DHCP IPs), so resolve the shared managed VLAN
    # subnet (objects.800) and grab its extId to feed into the enable request.
    subnet_ext_id, cluster_ext_id = ensure_flow_managed_subnet(client, env)
    if subnet_ext_id is None:
        return False

    logger.info("=== Enabling Flow Controller ===")

    dm_id = _discover_domain_manager_id(client, env)
    if not dm_id:
        return False
    logger.info("Using domain manager extId=%s", dm_id)

    products_path = f"{_DOMAIN_MGR_V44}/{dm_id}/products"
    _, body = client.request("GET", products_path, mutating=False)
    products = _dig(body, "data") or []
    flow = None
    for prod in products if isinstance(products, list) else []:
        if isinstance(prod, dict) and prod.get("name") == "FLOW_CONTROLLER":
            flow = prod
            break
    if not flow:
        logger.error("FLOW_CONTROLLER product not found under %s.", products_path)
        return False

    product_id = flow.get("extId")
    metadata = flow.get("metadata") or {}
    # Prefer the subnet/cluster from the network-prep step; fall back to whatever
    # the product already carries (e.g. when subnet.name was left unset).
    if not subnet_ext_id:
        subnet_ext_id = metadata.get("subnetExtId")
    if not cluster_ext_id:
        cluster_ext_id = metadata.get("clusterExtId") or _discover_pe_cluster_ext_id(client)
    if client.dry_run and not cluster_ext_id:
        cluster_ext_id = "DRY-RUN-CLUSTER"
    if not product_id or not cluster_ext_id or not subnet_ext_id:
        logger.error("Cannot enable FLOW_CONTROLLER -- missing extId/clusterExtId/"
                     "subnetExtId (extId=%s cluster=%s subnet=%s).", product_id,
                     cluster_ext_id, subnet_ext_id)
        return False

    if flow.get("enablementState") == "ENABLED":
        logger.info("FLOW_CONTROLLER is already ENABLED (product extId=%s) -- nothing "
                    "to do.", product_id)
        return True

    product_path = f"{products_path}/{product_id}"
    put_body = {
        "$objectType": "prism.v4.management.Product",
        "extId": product_id,
        "name": "FLOW_CONTROLLER",
        "enablementState": "ENABLED",
        "metadata": {
            "$objectType": "prism.v4.management.FlowControllerMetadata",
            "clusterExtId": cluster_ext_id,
            "subnetExtId": subnet_ext_id,
            "isFNSFlexModeEnabled": False,
            "isSecurityAnalyticsEnabled": True,
        },
    }
    # Carry the appliance sizing forward: this PUT is a full replacement, so an
    # omitted resourceSpec could reset a pre-sized Flow Controller's CPU/memory.
    if isinstance(flow.get("resourceSpec"), dict):
        put_body["resourceSpec"] = flow["resourceSpec"]

    if client.dry_run:
        logger.info("[dry-run] would PUT %s to enable Flow Controller "
                    "(cluster=%s subnet=%s).", product_path, cluster_ext_id, subnet_ext_id)
        return True

    # The product carries its ETag both as an HTTP header and inline at
    # $reserved.etag; prefer the header (canonical) and fall back to the body.
    etag = client.get_etag(product_path) or _dig(flow, "$reserved.etag")
    if not etag:
        logger.error("Failed to fetch ETag for %s (required for If-Match).", product_path)
        return False

    status_code, put_resp = client.request(
        "PUT", product_path, put_body, extra_headers={"If-Match": etag})
    if status_code and status_code >= 400:
        logger.error("Enabling Flow Controller failed (http %s): %s", status_code,
                     _truncate(json.dumps(put_resp) if isinstance(put_resp, (dict, list))
                               else str(put_resp)))
        return False

    task_ext_id = _dig(put_resp, "data.extId", "extId")
    if task_ext_id:
        # Deploying the Flow Controller (SMSP cluster) is slow -- it can take ~90
        # minutes. Poll for up to 2h by default; override with
        # pc_prep.flow.enable_timeout (env PC_PREP_FLOW_ENABLE_TIMEOUT).
        timeout = int(env.get("PC_PREP_FLOW_ENABLE_TIMEOUT")
                      or env.get("PC_PREP_FLOW_TASK_TIMEOUT") or 7200)
        logger.info("Waiting for Flow Controller deployment (timeout %ds); this "
                    "typically takes ~90 minutes.", timeout)
        if not poll_task(client, task_ext_id, timeout_seconds=timeout,
                         interval_seconds=15):
            logger.error("Flow Controller enable task did not succeed.")
            return False
    else:
        logger.warning("No task extId in enable PUT response (http %s) -- assuming "
                       "applied synchronously.", status_code)

    logger.info("Flow Controller enabled (product extId=%s, subnet=%s).", product_id,
                subnet_ext_id)
    return True


def _flow_flex_mode_enabled(metadata: dict) -> bool:
    """Return True if the FLOW_CONTROLLER metadata already has FLEX (rule-centric)
    mode on, checking both the legacy flag and the newer capabilities form."""
    if not isinstance(metadata, dict):
        return False
    if metadata.get("isFNSFlexModeEnabled") is True:
        return True
    return bool(_dig(metadata, "capabilities.flexMode.isEnabled") is True)


def switch_flow_policy_model_to_rule_centric(client: PcClient, env: Env) -> bool:
    """Switch the Flow Network Security policy model from App-Centric to
    Rule-Centric (a.k.a. FNS FLEX mode / Enhanced Policy Model).

    Rule-Centric mode is what unlocks FLEX rules (microseg RULETYPE_FLEX); in the
    default App-Centric model the FLEX rule type is not handled and policy create
    fails server-side. The switch is a domain-manager product update on the
    already-deployed FLOW_CONTROLLER (there is no separate microseg endpoint for
    it): flip metadata.isFNSFlexModeEnabled (and capabilities.flexMode.isEnabled)
    to true via PUT .../management/domain-managers/{dm}/products/{id}.

    Must run AFTER enable_flow_controller (the controller has to be ENABLED). The
    mode switch also requires Security Analytics (SAI) to be off, so this sets
    isSecurityAnalyticsEnabled=false during the switch (ENG-923204).

    Gated by pc_prep.flow.policy_model (env PC_PREP_FLOW_POLICY_MODEL):
      - "rule_centric" (default) -> ensure FLEX mode is enabled
      - "app_centric"            -> leave App-Centric (no-op; reverting rule->app
                                    requires deleting policies first and is not
                                    automated here)
    """
    model = (env.get("PC_PREP_FLOW_POLICY_MODEL") or "rule_centric").strip().lower()
    if model in ("app_centric", "app-centric", "appcentric"):
        logger.info("pc_prep.flow.policy_model=%s -- leaving Flow in App-Centric "
                    "model (skipping rule-centric switch).", model)
        return True
    if model not in ("rule_centric", "rule-centric", "rulecentric", "flex", ""):
        logger.warning("Unknown pc_prep.flow.policy_model=%r; defaulting to "
                       "rule_centric.", model)

    logger.info("=== Switching Flow policy model to Rule-Centric (FLEX) ===")

    dm_id = _discover_domain_manager_id(client, env)
    if not dm_id:
        return False

    products_path = f"{_DOMAIN_MGR_V44}/{dm_id}/products"
    _, body = client.request("GET", products_path, mutating=False)
    products = _dig(body, "data") or []
    flow = None
    for prod in products if isinstance(products, list) else []:
        if isinstance(prod, dict) and prod.get("name") == "FLOW_CONTROLLER":
            flow = prod
            break
    if not flow:
        logger.error("FLOW_CONTROLLER product not found under %s.", products_path)
        return False

    product_id = flow.get("extId")
    metadata = flow.get("metadata") or {}

    if flow.get("enablementState") != "ENABLED":
        # The controller must be deployed before its policy model can change.
        if client.dry_run:
            logger.info("[dry-run] FLOW_CONTROLLER not ENABLED yet; would switch to "
                        "rule-centric once the enable step completes.")
            return True
        logger.error("Cannot switch policy model: FLOW_CONTROLLER is not ENABLED "
                     "(state=%s). Run the flow enable step first.",
                     flow.get("enablementState"))
        return False

    if _flow_flex_mode_enabled(metadata):
        logger.info("Flow policy model is already Rule-Centric (FLEX mode enabled) "
                    "-- nothing to do.")
        return True

    cluster_ext_id = metadata.get("clusterExtId")
    subnet_ext_id = metadata.get("subnetExtId")
    if not product_id or (not client.dry_run and (not cluster_ext_id or not subnet_ext_id)):
        logger.error("Cannot switch policy model -- missing extId/clusterExtId/"
                     "subnetExtId on FLOW_CONTROLLER (extId=%s cluster=%s subnet=%s).",
                     product_id, cluster_ext_id, subnet_ext_id)
        return False

    # Preserve any existing flex-mode scale (e.g. securityPolicyCount) if the
    # capabilities block is already present.
    flex_scale = _dig(metadata, "capabilities.flexMode.scale")
    flex_mode = {"isEnabled": True}
    if isinstance(flex_scale, dict):
        flex_mode["scale"] = flex_scale

    product_path = f"{products_path}/{product_id}"
    put_body = {
        "$objectType": "prism.v4.management.Product",
        "extId": product_id,
        "name": "FLOW_CONTROLLER",
        "enablementState": "ENABLED",
        "metadata": {
            "$objectType": "prism.v4.management.FlowControllerMetadata",
            "clusterExtId": cluster_ext_id,
            "subnetExtId": subnet_ext_id,
            # Rule-Centric / FLEX policy model, set in both the legacy flag and the
            # newer capabilities form (Guru syncs them, but send both to be safe).
            "isFNSFlexModeEnabled": True,
            # SAI must be off while switching the policy model (ENG-923204).
            "isSecurityAnalyticsEnabled": False,
            "capabilities": {"flexMode": flex_mode},
        },
    }
    # Full-replacement PUT: carry the appliance sizing forward so we don't reset a
    # pre-sized Flow Controller's CPU/memory.
    if isinstance(flow.get("resourceSpec"), dict):
        put_body["resourceSpec"] = flow["resourceSpec"]

    if client.dry_run:
        logger.info("[dry-run] would PUT %s to switch policy model to Rule-Centric "
                    "(isFNSFlexModeEnabled=true, SAI=false).", product_path)
        return True

    etag = client.get_etag(product_path) or _dig(flow, "$reserved.etag")
    if not etag:
        logger.error("Failed to fetch ETag for %s (required for If-Match).", product_path)
        return False

    status_code, put_resp = client.request(
        "PUT", product_path, put_body, extra_headers={"If-Match": etag})
    if status_code and status_code >= 400:
        logger.error("Switching policy model to Rule-Centric failed (http %s): %s",
                     status_code,
                     _truncate(json.dumps(put_resp) if isinstance(put_resp, (dict, list))
                               else str(put_resp)))
        return False

    task_ext_id = _dig(put_resp, "data.extId", "extId")
    if task_ext_id:
        timeout = int(env.get("PC_PREP_FLOW_POLICY_MODEL_TIMEOUT")
                      or env.get("PC_PREP_FLOW_TASK_TIMEOUT") or 1800)
        logger.info("Waiting for policy-model switch (kFlowSmspUpdate) task "
                    "(timeout %ds).", timeout)
        if not poll_task(client, task_ext_id, timeout_seconds=timeout,
                         interval_seconds=15):
            logger.error("Policy-model switch task did not succeed.")
            return False
    else:
        logger.warning("No task extId in policy-model PUT response (http %s) -- "
                       "assuming applied synchronously.", status_code)

    logger.info("Flow policy model switched to Rule-Centric (FLEX) on product "
                "extId=%s.", product_id)
    return True


def prepare_flow(client: PcClient, env: Env) -> bool:
    """'flow' step: enable (deploy) the Flow Controller, then switch its policy
    model to Rule-Centric (FLEX) so FLEX network security policies are supported.

    The policy-model switch is skipped (leaving App-Centric) when
    pc_prep.flow.policy_model=app_centric."""
    if not enable_flow_controller(client, env):
        return False
    return switch_flow_policy_model_to_rule_centric(client, env)


# --------------------------------------------------------------------------- #
# main
# --------------------------------------------------------------------------- #
STEPS = ("nucalm", "objects", "policy", "lcm", "dp", "prism", "iscsi", "ngt",
         "nkp", "selfservice", "flow")


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Prepare Prism Central: enable Flow Controller, NuCalm/DR, and Objects")
    default_env = SCRIPT_DIR / "config.yaml"
    if not default_env.exists():
        default_env = SCRIPT_DIR / ".env"
    parser.add_argument("--env", default=str(default_env),
                        help="config file: YAML (.yaml/.yml) or flat .env "
                             "(default: testenv/config.yaml if present, else testenv/.env)")
    parser.add_argument("--only", choices=STEPS, action="append", default=None,
                        help="run only this step (repeatable): "
                             "nucalm|objects|policy|lcm|dp|prism|iscsi|ngt|nkp|selfservice|flow")
    parser.add_argument("--skip", choices=STEPS, action="append", default=[],
                        help="skip this step (repeatable): "
                             "nucalm|objects|policy|lcm|dp|prism|iscsi|ngt|nkp|selfservice|flow")
    parser.add_argument("--dry-run", action="store_true",
                        help="log intended calls; make no mutations")
    parser.add_argument("--log-dir", default=str(SCRIPT_DIR / "logs"),
                        help="directory for the run log file (default: testenv/logs)")
    parser.add_argument("--no-log-file", action="store_true",
                        help="log only to the console; do not write a log file")
    parser.add_argument("--verbose", "-v", action="store_true",
                        help="show DEBUG output (incl. HTTP traces) on the console")
    args = parser.parse_args()

    log_dir = None if args.no_log_file else Path(args.log_dir)
    log_file = setup_logging(log_dir, args.verbose)

    selected = tuple(args.only) if args.only else STEPS
    selected = tuple(s for s in selected if s not in set(args.skip))

    logger.info("=== prepare_pc started ===")
    logger.info("env=%s dry_run=%s steps=%s", args.env, args.dry_run, ",".join(selected) or "(none)")
    if log_file:
        logger.info("Full debug log: %s", log_file)
    if not selected:
        logger.warning("No steps selected -- nothing to do.")
        return 0

    try:
        env = Env(load_env(Path(args.env)))
        client = PcClient(env, dry_run=args.dry_run)

        handlers = {
            "nucalm": enable_nucalm,
            "objects": enable_object_store,
            "policy": enable_policy_engine,
            "lcm": downgrade_lcm,
            "dp": prepare_data_protection,
            "prism": prepare_prism,
            "iscsi": create_iscsi_client_vm,
            "ngt": create_ngt_upgrade_vm,
            "nkp": ensure_nkp_image,
            "selfservice": prepare_self_service,
            "flow": prepare_flow,
        }
        failed = []
        for step in selected:
            try:
                if not handlers[step](client, env):
                    failed.append(step)
            except Exception as exc:  # noqa: BLE001 - surface a clean per-step message
                logger.error("Step '%s' failed: %s", step, exc)
                logger.debug("traceback", exc_info=True)
                failed.append(step)

        if failed:
            logger.error("Completed with failures in: %s", ", ".join(failed))
            return 1
        logger.info("=== prepare_pc completed successfully (%s) ===", ",".join(selected))
        return 0
    except SystemExit:
        raise
    except Exception:  # noqa: BLE001
        logger.exception("Unexpected failure")
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
