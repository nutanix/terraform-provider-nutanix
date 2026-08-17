#!/usr/bin/env python3
"""Populate test_config_v2.json from a local .env file plus live IAM provisioning.

This script is SAFE TO COMMIT: it contains no secrets. All sensitive / environment
specific values are read from a gitignored ``.env`` file (see
``testenv/test_config.env.example`` for the full list of keys).

What it does
------------
1. Reads ``testenv/.env`` (override with ``--env``).
2. Fills the top-level test_config_v2.json fields that the acceptance tests
   consume: pc/pe credentials, ssh credentials, dns/ntp servers, the ``clusters``
   block (nodes, network, ssl_certificate) and ``data_protection``. When
   ``cluster.ssl.generate=true`` it generates a fresh rsa2048 leaf cert signed by a
   throwaway single-tier CA (created on the fly via the openssl CLI) and
   base64-encodes it into clusters.ssl_certificate. The cluster ``virtual_ip``
   and ``iscsi_ip`` are
   discovered live from the first node's Prism Element REST API (using the PE
   credentials) so they don't have to be hand-maintained in .env; the
   CLUSTER_VIRTUAL_IP / CLUSTER_ISCSI_IP env vars are only used as a fallback when
   the cluster doesn't return them (or with --skip-iam).
3. For the ``iam`` block it talks to the Prism Central IAM v4 API to
   provision (or reuse, by name) the prerequisite objects and capture their
   ext_ids:
     * a SAML identity provider          -> iam.identity_providers
       (created from IAM_IDP_METADATA_XML_FILE, e.g. testenv/federationmetadata.xml;
        IAM_IDP_METADATA_TEST_FILE is copied to the repo-root test_idp_metadata.txt
        that the iamv2 SAML tests read via file(...))
     * two ACTIVE_DIRECTORY services      -> iam.directory_services_main.{primary,secondary}_ad
       (and searches each directory for the configured users/groups to fill
        domain_users_usergroups ext_ids)
     * a user                             -> iam.users
     * a user group                       -> iam.user_groups
       (when IAM_OPENLDAP_NAME is set, an OPEN_LDAP directory is provisioned and a
        real user + group are discovered from it -- like preEnv/iam.tf -- so the
        LDAP user group always exists in the directory; otherwise the hardcoded
        IAM_USER_* / IAM_GROUP_* values are used)
4. Merges everything into the existing test_config_v2.json (preserving every
   other section) and writes it back in place. The only file the script writes
   is test_config_v2.json (plus a debug log under testenv/logs/).

Usage
-----
    cp testenv/test_config.env.example testenv/.env   # then edit testenv/.env
    python3 testenv/fill_test_config.py --dry-run     # preview, no writes/mutations
    python3 testenv/fill_test_config.py               # fill everything
    python3 testenv/fill_test_config.py --skip-iam    # offline: only top-level fields
    python3 testenv/fill_test_config.py -v            # verbose (HTTP traces on console)

Logging
-------
Every run writes a full-detail, timestamped debug log to testenv/logs/ (override
with --log-dir, disable with --no-log-file). The console shows INFO by default
(-v / --verbose for DEBUG). Secrets (passwords, keys, auth headers) are redacted
from the logs. Log files match *.log and are gitignored.

Standard library only -- no pip install required.
"""

from __future__ import annotations

import argparse
import base64
import concurrent.futures
import json
import logging
import os
import re
import shutil
import socket
import ssl
import subprocess
import sys
import tempfile
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path

SCRIPT_DIR = Path(__file__).resolve().parent
REPO_ROOT = SCRIPT_DIR.parent

logger = logging.getLogger("fill_test_config")


# --------------------------------------------------------------------------- #
# .env loading + small typed accessors
# --------------------------------------------------------------------------- #
def load_env(path: Path) -> dict:
    """Load configuration from a YAML file (nested + anchors) or a flat KEY=VALUE
    .env file, returning the flat KEY=VALUE dict the rest of the script consumes.
    Process env vars take precedence over file values."""
    if path.suffix.lower() in (".yaml", ".yml"):
        values = load_yaml_config(path)
    elif path.exists():
        values = _parse_dotenv(path)
    else:
        log(f"WARNING: config file not found at {path} (relying on process env only)")
        values = {}
    # Process environment overrides file values.
    for key in list(values) + [k for k in os.environ if k.startswith(("PC_", "PE_", "CFG_", "SSH_", "IAM_", "DNS_", "NTP_", "CLUSTER_", "SSL_", "DP_"))]:
        if key in os.environ:
            values[key] = os.environ[key]
    return values


def _parse_dotenv(path: Path) -> dict:
    """Parse a simple KEY=VALUE .env file."""
    values: dict = {}
    for raw in path.read_text().splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, _, val = line.partition("=")
        values[key.strip()] = val.strip().strip('"').strip("'")
    return values


def load_yaml_config(path: Path) -> dict:
    """Load a nested YAML config and flatten it to the flat KEY=VALUE names the
    script uses. Nested maps join with '_' and upper-case (pc.endpoint ->
    PC_ENDPOINT, iam.idp.metadata.login_url -> IAM_IDP_METADATA_LOGIN_URL); lists
    become comma-joined strings. YAML anchors/aliases (&name / *name) deduplicate
    shared values, and top-level keys starting with '_' or '.' are treated as
    anchor-only and skipped."""
    try:
        import yaml
    except ImportError:
        raise SystemExit(
            "ERROR: PyYAML is required to read a YAML config. Install it with\n"
            "    python3 -m pip install pyyaml\n"
            "or point --env at a flat .env file instead.")
    if not path.exists():
        log(f"WARNING: config file not found at {path} (relying on process env only)")
        return {}
    data = yaml.safe_load(path.read_text()) or {}
    if not isinstance(data, dict):
        raise SystemExit(f"ERROR: {path} must contain a YAML mapping at the top level")
    return _flatten_config(data)


def _flatten_config(data: dict, prefix: str = "") -> dict:
    out: dict = {}
    for key, value in data.items():
        if prefix == "" and str(key).startswith(("_", ".")):
            continue  # anchor-only / hidden section
        flat = f"{prefix}_{key}" if prefix else str(key)
        if isinstance(value, dict):
            out.update(_flatten_config(value, flat))
        elif isinstance(value, list):
            out[flat.upper()] = ",".join("" if v is None else str(v) for v in value)
        elif isinstance(value, bool):
            out[flat.upper()] = "true" if value else "false"
        elif value is None:
            out[flat.upper()] = ""
        else:
            out[flat.upper()] = str(value)
    return out


class Env:
    def __init__(self, data: dict):
        self.data = data

    def get(self, key: str, default: str = "") -> str:
        return self.data.get(key, default)

    def required(self, key: str) -> str:
        val = self.data.get(key, "")
        if val == "":
            raise SystemExit(f"ERROR: required env var '{key}' is empty/missing")
        return val

    def list(self, key: str) -> list:
        raw = self.data.get(key, "")
        return [item.strip() for item in raw.split(",") if item.strip()]

    def bool(self, key: str, default: bool = False) -> bool:
        raw = self.data.get(key, "")
        if raw == "":
            return default
        return raw.lower() in ("1", "true", "yes", "on")

    def file_or_inline(self, file_key: str, inline_key: str) -> str:
        """Return file contents if *_FILE is set, else the inline value."""
        file_path = self.data.get(file_key, "")
        if file_path:
            p = Path(file_path)
            if not p.is_absolute():
                p = (REPO_ROOT / file_path).resolve()
            if not p.exists():
                raise SystemExit(f"ERROR: file referenced by {file_key} not found: {p}")
            return p.read_text()
        return self.data.get(inline_key, "")


def log(msg: str) -> None:
    """Backwards-compatible info-level log used throughout the script."""
    logger.info(msg)


# Keys whose values must never be written to the log files.
_SECRET_HINTS = ("password", "secret", "passphrase", "private_key", "privatekey",
                 "access_key", "client_secret", "secret_key", "token", "auth")


def _redact(obj):
    """Deep-copy a JSON-ish value, masking anything that looks like a secret."""
    if isinstance(obj, dict):
        out = {}
        for key, val in obj.items():
            if any(hint in str(key).lower() for hint in _SECRET_HINTS):
                out[key] = "***REDACTED***"
            else:
                out[key] = _redact(val)
        return out
    if isinstance(obj, list):
        return [_redact(item) for item in obj]
    return obj


def _truncate(text, limit: int = 2000) -> str:
    text = str(text)
    if len(text) <= limit:
        return text
    return f"{text[:limit]}...<truncated {len(text) - limit} chars>"


def _safe_json(raw: str) -> str:
    """Redact + truncate a (possibly JSON) response body for logging."""
    try:
        return _truncate(json.dumps(_redact(json.loads(raw))))
    except (ValueError, TypeError):
        return _truncate(raw)


def setup_logging(log_dir: Path | None, verbose: bool) -> Path | None:
    """Configure console (INFO, or DEBUG with --verbose) + full-detail file logging."""
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
        log_file = log_dir / f"fill_test_config_{time.strftime('%Y%m%d_%H%M%S')}.log"
        file_handler = logging.FileHandler(log_file, encoding="utf-8")
        file_handler.setLevel(logging.DEBUG)  # always capture full detail to file
        file_handler.setFormatter(fmt)
        logger.addHandler(file_handler)
        logger.debug("Logging to %s", log_file)
    return log_file


# --------------------------------------------------------------------------- #
# IAM v4 REST client (get-or-create + directory search)
# --------------------------------------------------------------------------- #
class IamClient:
    def __init__(self, env: Env, dry_run: bool = False):
        self.endpoint = env.required("PC_ENDPOINT")
        self.port = env.required("PC_PORT")
        self.username = env.required("PC_USERNAME")
        self.password = env.required("PC_PASSWORD")
        self.insecure = env.bool("PC_INSECURE")
        self.version = env.required("IAM_API_VERSION")
        self.fv = env.required("IAM_FV")
        # Candidate API versions for the directory-service "share with all
        # projects" action, newest first. It lives under the config namespace
        # introduced in v4.1.b3 (pc.2024.3+/AOS 7.6); older PCs (e.g. AOS 7.5)
        # top out at v4.1.b2 and may not offer it at all -- so we probe these in
        # order and skip gracefully if none are supported.
        self.share_api_versions = ["v4.1.b3", "v4.1.b2"]
        self.dry_run = dry_run
        token = base64.b64encode(f"{self.username}:{self.password}".encode()).decode()
        self.auth_header = f"Basic {token}"
        self.ctx = ssl.create_default_context()
        if self.insecure:
            self.ctx.check_hostname = False
            self.ctx.verify_mode = ssl.CERT_NONE

    def base(self) -> str:
        return f"https://{self.endpoint}:{self.port}/api/iam/{self.version}/authn"

    def _request(self, method: str, url: str, body: dict | None = None) -> dict:
        data = json.dumps(body).encode() if body is not None else None
        redacted_body = _truncate(json.dumps(_redact(body))) if body is not None else "<none>"
        if body is not None:
            logger.debug("HTTP %s %s\n  request: %s", method, url,
                         redacted_body)
        else:
            logger.debug("HTTP %s %s", method, url)
        req = urllib.request.Request(url, data=data, method=method)
        req.add_header("Authorization", self.auth_header)
        req.add_header("Accept", "application/json")
        if data is not None:
            req.add_header("Content-Type", "application/json")
        started = time.monotonic()
        try:
            with urllib.request.urlopen(req, context=self.ctx, timeout=60) as resp:
                raw = resp.read().decode()
                status = getattr(resp, "status", resp.getcode())
        except urllib.error.HTTPError as exc:
            elapsed = time.monotonic() - started
            detail = exc.read().decode(errors="replace")
            logger.error("HTTP %s %s -> %d after %.1fs\n  request: %s\n  response: %s",
                         method, url, exc.code, elapsed, redacted_body, _truncate(detail))
            raise RuntimeError(f"{method} {url} -> HTTP {exc.code}: {detail}") from None
        except (socket.timeout, TimeoutError) as exc:
            elapsed = time.monotonic() - started
            logger.error("HTTP %s %s -> timeout after %.1fs (timeout=60s): %s\n  request: %s",
                         method, url, elapsed, exc, redacted_body)
            raise RuntimeError(f"{method} {url} -> timeout after {elapsed:.1f}s: {exc}") from None
        except urllib.error.URLError as exc:
            elapsed = time.monotonic() - started
            if isinstance(exc.reason, (socket.timeout, TimeoutError)):
                logger.error("HTTP %s %s -> timeout after %.1fs (timeout=60s): %s\n  request: %s",
                             method, url, elapsed, exc.reason, redacted_body)
                raise RuntimeError(
                    f"{method} {url} -> timeout after {elapsed:.1f}s: {exc.reason}") from None
            logger.error("HTTP %s %s -> connection error after %.1fs: %s\n  request: %s",
                         method, url, elapsed, exc.reason, redacted_body)
            raise RuntimeError(f"{method} {url} -> connection error: {exc.reason}") from None
        elapsed = time.monotonic() - started
        logger.debug("HTTP %s %s -> %s in %.1fs (%d bytes)\n  response: %s", method, url,
                     status, elapsed, len(raw), _safe_json(raw) if raw.strip() else "<empty>")
        return json.loads(raw) if raw.strip() else {}

    def _send(self, method: str, url: str, body: dict | None = None,
              headers: dict | None = None, quiet: bool = False):
        """Like _request but supports extra request headers and returns
        (parsed_body, response_headers). response_headers is the case-insensitive
        http.client message object. Pass quiet=True to log expected failures
        (e.g. version probes) at debug level instead of error."""
        data = json.dumps(body).encode() if body is not None else None
        redacted_body = _truncate(json.dumps(_redact(body))) if body is not None else "<none>"
        logger.debug("HTTP %s %s%s", method, url,
                     "" if body is None else f"\n  request: {redacted_body}")
        req = urllib.request.Request(url, data=data, method=method)
        req.add_header("Authorization", self.auth_header)
        req.add_header("Accept", "application/json")
        if data is not None:
            req.add_header("Content-Type", "application/json")
        for key, value in (headers or {}).items():
            if value is not None:
                req.add_header(key, value)
        started = time.monotonic()
        try:
            with urllib.request.urlopen(req, context=self.ctx, timeout=60) as resp:
                raw = resp.read().decode()
                resp_headers = resp.headers
        except urllib.error.HTTPError as exc:
            elapsed = time.monotonic() - started
            detail = exc.read().decode(errors="replace")
            (logger.debug if quiet else logger.error)(
                "HTTP %s %s -> %d after %.1fs\n  request: %s\n  response: %s",
                method, url, exc.code, elapsed, redacted_body, _truncate(detail))
            raise RuntimeError(f"{method} {url} -> HTTP {exc.code}: {detail}") from None
        except (socket.timeout, TimeoutError) as exc:
            elapsed = time.monotonic() - started
            (logger.debug if quiet else logger.error)(
                "HTTP %s %s -> timeout after %.1fs (timeout=60s): %s\n  request: %s",
                method, url, elapsed, exc, redacted_body)
            raise RuntimeError(f"{method} {url} -> timeout after {elapsed:.1f}s: {exc}") from None
        except urllib.error.URLError as exc:
            elapsed = time.monotonic() - started
            if isinstance(exc.reason, (socket.timeout, TimeoutError)):
                (logger.debug if quiet else logger.error)(
                    "HTTP %s %s -> timeout after %.1fs (timeout=60s): %s\n  request: %s",
                    method, url, elapsed, exc.reason, redacted_body)
                raise RuntimeError(
                    f"{method} {url} -> timeout after {elapsed:.1f}s: {exc.reason}") from None
            (logger.debug if quiet else logger.error)(
                "HTTP %s %s -> connection error after %.1fs: %s\n  request: %s",
                method, url, elapsed, exc.reason, redacted_body)
            raise RuntimeError(f"{method} {url} -> connection error: {exc.reason}") from None
        logger.debug("HTTP %s %s -> completed in %.1fs (%d bytes)", method, url,
                     time.monotonic() - started, len(raw))
        return (json.loads(raw) if raw.strip() else {}), resp_headers

    def share_directory_service_with_all_projects(self, ext_id: str) -> None:
        """Share a directory service (IDP) with all projects. Required before a
        role membership can be created for users/groups from that IDP (IAM-20027).
        The share endpoints live under the v4.1.b3 ``config`` namespace and need
        an If-Match (ETag) header taken from a prior GET of the directory service.
        """
        if not ext_id:
            return
        if self.dry_run:
            log(f"  - [dry-run] would share directory service {ext_id} with all projects")
            return
        # The share-all action lives under a newer IAM namespace (v4.1.b3 on PC
        # 7.6+). Older PCs (e.g. 7.5) only serve up to v4.1.b2 and may lack the
        # endpoint entirely -- try newest first, fall back to v4.1.b2, and skip
        # with a warning if unsupported so IAM prep doesn't abort on older releases.
        last_exc: Exception | None = None
        for ver in self.share_api_versions:
            base_v41 = f"https://{self.endpoint}:{self.port}/api/iam/{ver}/authn"
            try:
                get_body, get_headers = self._send(
                    "GET", f"{base_v41}/directory-services/{ext_id}", quiet=True)
            except RuntimeError as exc:
                last_exc = exc
                if self._is_version_unsupported(exc):
                    continue
                raise
            etag = get_headers.get("ETag")
            if not etag:
                entity = get_body.get("data") if isinstance(get_body.get("data"), dict) else get_body
                etag = ((entity or {}).get("$reserved") or {}).get("ETag")
            share_url = f"{base_v41}/config/directory-service/{ext_id}/$actions/share-all"
            try:
                self._send("POST", share_url, headers={"If-Match": etag}, quiet=True)
                log(f"  - shared directory service {ext_id} with all projects (iam {ver})")
                return
            except RuntimeError as exc:
                last_exc = exc
                if "already" in str(exc).lower():
                    log(f"  - directory service {ext_id} is already shared with all projects")
                    return
                if self._is_version_unsupported(exc):
                    continue
                raise
        logger.warning(
            "  - share-all endpoint unsupported on this PC (tried v4.1.b3/v4.1.b2); "
            "skipping share of directory service %s. IDP-user role memberships may "
            "fail (IAM-20027). Last error: %s", ext_id, last_exc)

    @staticmethod
    def _is_version_unsupported(exc: Exception) -> bool:
        """True when a request failed because the PC does not serve that API
        version / endpoint (older releases): the v4 gateway returns HTTP 400 with
        'Invalid API version passed in the request', or a plain 404."""
        s = str(exc)
        return "Invalid API version" in s or "-> HTTP 404" in s

    @staticmethod
    def _data(resp: dict):
        return resp.get("data", resp)

    def _envelope(self, object_type: str, payload: dict) -> dict:
        body = {"$objectType": object_type, "$reserved": {"$fv": self.fv}}
        body.update(payload)
        return body

    def _list(self, collection: str) -> list:
        """Return all entities from a list endpoint (best-effort pagination)."""
        items: list = []
        page = 0
        while True:
            url = f"{self.base()}/{collection}?$page={page}&$limit=100"
            data = self._data(self._request("GET", url))
            if not isinstance(data, list):
                data = data or []
            items.extend(data)
            if len(data) < 100:
                break
            page += 1
            if page > 50:  # safety valve
                break
        return items

    def find_by(self, collection: str, field: str, value: str) -> dict | None:
        """Look up an entity by an attribute. Tries an OData $filter first (as
        create_idps.sh does) and falls back to scanning the full list."""
        if not value:
            return None
        flt = urllib.parse.quote(f"{field} eq '{value}'")
        try:
            data = self._data(self._request("GET", f"{self.base()}/{collection}?$filter={flt}"))
            if isinstance(data, list):
                for item in data:
                    if item.get(field) == value:
                        return item
                if data:
                    return data[0]
        except RuntimeError:
            pass  # some fields/PC versions reject $filter; fall back to listing
        for item in self._list(collection):
            if item.get(field) == value:
                return item
        return None

    @staticmethod
    def _already_exists(resp: dict) -> bool:
        errors = (resp.get("data") or {}).get("error") if isinstance(resp.get("data"), dict) else None
        for err in errors or []:
            blob = json.dumps(err).lower()
            if "already exist" in blob:
                return True
        return False

    def get_or_create(self, collection: str, match_field: str, match_value: str,
                      payload: dict, label: str, resolvers: list | None = None) -> dict:
        """Create an entity, reusing an existing one (by match_field, then by any
        extra ``resolvers`` of (field, value)) when it already exists. Payloads are
        sent as plain JSON -- matching the working preEnv/scripts/create_idps.sh."""
        lookups = [(match_field, match_value)] + list(resolvers or [])

        for field, value in lookups:
            existing = self.find_by(collection, field, value)
            if existing:
                log(f"  - reusing existing {label} '{match_value}' (extId={existing.get('extId')})")
                return existing

        if self.dry_run:
            log(f"  - [dry-run] would CREATE {label} '{match_value}'")
            return {"extId": f"<dry-run-{label}-ext-id>", **payload}

        log(f"  - creating {label} '{match_value}'")
        exc_text = ""
        try:
            resp = self._request("POST", f"{self.base()}/{collection}", payload)
        except RuntimeError as exc:
            # A 4xx "already exists" is expected when re-running -> resolve instead.
            # Any other error (bad creds, unreachable LDAP, ...) is a real failure
            # and must propagate so the run fails.
            exc_text = str(exc)
            if "already exist" not in exc_text.lower():
                raise
            resp = {}

        created = self._data(resp) if resp else {}
        if isinstance(created, dict) and created.get("extId"):
            return created

        # Create returned no extId (already-exists error, or async task) -> resolve.
        for field, value in lookups:
            resolved = self.find_by(collection, field, value)
            if resolved:
                log(f"  - resolved existing {label} '{match_value}' (extId={resolved.get('extId')})")
                return resolved

        # LDAP users aren't returned by the list endpoint, but the "already exists"
        # error carries the uuid directly (e.g. "... already exists with uuid <id>").
        m = re.search(r"already exists with uuid\s+([0-9a-fA-F-]{36})", exc_text)
        if m:
            ext_id = m.group(1)
            log(f"  - resolved existing {label} '{match_value}' from error (extId={ext_id})")
            return {"extId": ext_id}
        raise RuntimeError(f"failed to create or resolve {label} '{match_value}'")

    # Attribute names differ by directory type: Active Directory exposes accounts
    # via sAMAccountName/userPrincipalName, while OpenLDAP uses cn/uid. Searching
    # the wrong set returns an empty result list (see IAM v4 search API / ldap.go).
    _AD_SEARCHED_ATTRS = ["sAMAccountName", "userPrincipalName", "name", "distinguishedName"]
    _AD_RETURNED_ATTRS = ["sAMAccountName", "userPrincipalName", "name",
                          "distinguishedName", "memberOf", "uuid"]
    _LDAP_SEARCHED_ATTRS = ["cn", "dn", "name", "uid"]
    _LDAP_RETURNED_ATTRS = ["cn", "dn", "uid", "memberUid", "uuid"]

    def search_directory(self, ext_id: str, query: str,
                         searched_attributes: list = None,
                         returned_attributes: list = None) -> list:
        """Run a directory-service search and return the raw entity list."""
        if self.dry_run or not ext_id or ext_id.startswith("<dry-run"):
            return []
        url = f"{self.base()}/directory-services/{ext_id}/$actions/search"
        body = self._envelope("iam.v4.authn.DirectoryServiceSearchQuery", {
            "query": query,
            "searchedAttributes": searched_attributes or self._LDAP_SEARCHED_ATTRS,
            "returnedAttributes": returned_attributes or self._LDAP_RETURNED_ATTRS,
            "isWildcardSearch": True,
        })
        data = self._data(self._request("POST", url, body))
        if isinstance(data, dict):
            return data.get("searchResults") or data.get("entities") or []
        return data or []

def _prune_empty(obj):
    """Recursively drop empty strings/None/empty containers (keeps booleans)."""
    if isinstance(obj, dict):
        return {k: _prune_empty(v) for k, v in obj.items()
                if _prune_empty(v) not in ("", None, {}, [])}
    if isinstance(obj, list):
        return [_prune_empty(v) for v in obj]
    return obj


def _rdn_value(dn: str) -> str:
    """Return the value of the first RDN of a distinguished name, e.g.
    'CN=dnd_approval_group_1,CN=Users,DC=qa,DC=nucalm,DC=io' ->
    'dnd_approval_group_1'. Returns the input unchanged when it isn't a DN."""
    if not dn:
        return ""
    first = dn.split(",")[0].strip()
    return first.split("=", 1)[1] if "=" in first else first


def attr(entity: dict, name: str) -> str:
    """Pull the first value of a named attribute out of a search entity."""
    for a in entity.get("attributes", []) or []:
        if a.get("name") == name:
            vals = a.get("values") or []
            if vals:
                return vals[0]
    return ""


def search_ext_id_for(client: IamClient, ds_ext_id: str, wanted_name: str,
                      entity_type: str = "", directory_type: str = "ACTIVE_DIRECTORY") -> str:
    """Find a directory entity whose name matches wanted_name; return its ext_id.

    The directory-search `query` must be the search term itself (the name), not a
    fixed literal, and the searched attributes must match the directory type:
    Active Directory keys on sAMAccountName/userPrincipalName, OpenLDAP on cn/uid
    (see the IAM v4 search API / ldap.go)."""
    target = wanted_name.split("@")[0].lower()
    if directory_type == "ACTIVE_DIRECTORY":
        searched = IamClient._AD_SEARCHED_ATTRS
        returned = IamClient._AD_RETURNED_ATTRS
    else:
        searched = IamClient._LDAP_SEARCHED_ATTRS
        returned = IamClient._LDAP_RETURNED_ATTRS
    for entity in client.search_directory(ds_ext_id, target, searched, returned):
        etype = (entity.get("entityType") or "").lower()
        if entity_type and etype and etype != entity_type:
            continue
        candidates = {
            (entity.get("name") or "").lower(),
            attr(entity, "cn").lower(),
            attr(entity, "uid").lower(),
            attr(entity, "sAMAccountName").lower(),
            attr(entity, "userPrincipalName").split("@")[0].lower(),
        }
        if target in candidates:
            return entity.get("identityExtId") or attr(entity, "uuid") or ""
    return ""


def materialize_ad_user(client: IamClient, ds_ext_id: str, wanted: str) -> str:
    """Create (or reuse) an IAM LDAP-user reference for an AD identity and return
    its IAM extId.

    A directory *search* only yields the directory identity id, which role
    membership rejects with IAM-20027 ('user ... does not exist'). The identity
    must exist as an IAM user entity first, exactly like nutanix_users_v2
    (user_type=LDAP) in preEnv/iam.tf. We look the account up in the directory to
    learn the username the AD actually exposes (the configured key may use a
    different UPN suffix than the real domain), then create the user with it."""
    if not ds_ext_id or client.dry_run:
        return ""
    target = wanted.split("@")[0].lower()
    username = ""
    for entity in client.search_directory(ds_ext_id, target,
                                          IamClient._AD_SEARCHED_ATTRS,
                                          IamClient._AD_RETURNED_ATTRS):
        etype = (entity.get("entityType") or "").lower()
        if etype and etype != "person":
            continue
        upn = attr(entity, "userPrincipalName")
        sam = attr(entity, "sAMAccountName")
        candidates = {
            (entity.get("name") or "").lower(),
            sam.lower(),
            upn.split("@")[0].lower(),
        }
        if target in candidates:
            username = upn or sam or (entity.get("name") or "")
            break
    if not username:
        raise RuntimeError(f"AD user '{wanted}' not found in directory {ds_ext_id}; "
                           f"cannot create the user reference for role membership")
    payload = {"username": username, "userType": "LDAP", "idpId": ds_ext_id}
    entity = client.get_or_create("users", "username", username, payload,
                                  f"AD user '{username}'")
    return entity.get("extId", "")


def discover_ad_user_identity(client: IamClient, ds_ext_id: str,
                              wanted: str) -> tuple:
    """Look an AD account up in the directory and return its real
    (userPrincipalName, displayName).

    *wanted* is a sAMAccountName / short name (e.g. 'test1'); the configured
    value may not carry the correct UPN suffix, so the actual UPN + display name
    are read from the directory (never hardcoded). Returns ("", "") when the
    account isn't found. Adds ``displayName`` to the returned attributes since the
    default AD set does not include it."""
    if not ds_ext_id or client.dry_run:
        return "", ""
    target = wanted.split("@")[0].lower()
    returned = IamClient._AD_RETURNED_ATTRS + ["displayName"]
    for entity in client.search_directory(ds_ext_id, target,
                                          IamClient._AD_SEARCHED_ATTRS, returned):
        etype = (entity.get("entityType") or "").lower()
        if etype and etype != "person":
            continue
        upn = attr(entity, "userPrincipalName")
        sam = attr(entity, "sAMAccountName")
        candidates = {
            (entity.get("name") or "").lower(),
            sam.lower(),
            upn.split("@")[0].lower(),
        }
        if target in candidates:
            display = (attr(entity, "displayName") or attr(entity, "name")
                       or (entity.get("name") or ""))
            return (upn or sam or (entity.get("name") or "")), display
    return "", ""


def materialize_ad_group(client: IamClient, ds_ext_id: str, wanted: str) -> str:
    """Create (or reuse) an IAM LDAP user-group reference for an AD group and
    return its IAM extId.

    Mirrors nutanix_user_groups_v2 (group_type=LDAP) in preEnv/iam.tf: the
    group's distinguishedName is discovered from the directory (IAM rejects a
    group it cannot resolve, IAM-21807) and the group is created so role
    membership can target it (otherwise IAM-20027 'group ... does not exist')."""
    if not ds_ext_id or client.dry_run:
        return ""
    target = wanted.split("@")[0].lower()
    dn = ""
    for entity in client.search_directory(ds_ext_id, target,
                                          IamClient._AD_SEARCHED_ATTRS,
                                          IamClient._AD_RETURNED_ATTRS):
        etype = (entity.get("entityType") or "").lower()
        if etype and etype != "group":
            continue
        # For AD groups the top-level "name" is the full distinguishedName, so we
        # match on sAMAccountName/cn and the first RDN value, never on "name".
        candidates = {
            attr(entity, "sAMAccountName").lower(),
            attr(entity, "cn").lower(),
            _rdn_value(entity.get("name") or "").lower(),
            _rdn_value(attr(entity, "distinguishedName")).lower(),
        }
        if target in candidates:
            dn = (attr(entity, "distinguishedName")
                  or entity.get("distinguishedName", "")
                  or (entity.get("name") or ""))
            break
    if not dn:
        raise RuntimeError(f"AD group '{wanted}' not found in directory {ds_ext_id}; "
                           f"cannot create the user-group reference for role membership")
    # The group's IAM `name` is the configured simple name (the Go tests look the
    # group up by that key); the distinguishedName from the directory is what IAM
    # validates and uniquely identifies the group (IAM-20006 'already exists with
    # same DN'), so we also resolve an existing group by it.
    name, dn = target, dn.lower()
    payload = {"groupType": "LDAP", "idpId": ds_ext_id, "name": name,
               "distinguishedName": dn}
    entity = client.get_or_create("user-groups", "name", name, payload,
                                  f"AD group '{name}'",
                                  resolvers=[("distinguishedName", dn)])
    return entity.get("extId", "")


# --------------------------------------------------------------------------- #
# IAM block builder
# --------------------------------------------------------------------------- #
def build_identity_provider(client: IamClient, env: Env) -> tuple[dict, str]:
    name = env.get("IAM_IDP_NAME")
    metadata_url = env.get("IAM_IDP_METADATA_URL")

    # The IdP named IAM_IDP_NAME is NOT created here: every SAML *IdP* acceptance
    # test (resource + data source) creates its own nutanix_saml_identity_providers_v2
    # with this name and destroys it, so pre-creating it makes them fail with a 409
    # (IAM-20006, name already exists). The emitted iam.identity_providers block only
    # carries the values those tests read.
    #
    # The SAML *user* / *user-group* tests are different: they reference an existing
    # IdP via iam.users.idp_id and never delete it (the provider's user delete is a
    # no-op, so a test-managed IdP can't be torn down -- IAM-21016). So we create a
    # dedicated, persistent IdP (distinct name + entity issuer so it coexists with
    # the CRUD tests' IdP) and return its ext_id for iam.users.idp_id.
    log(f"IAM: identity provider '{name}' (values only, not created)")
    persistent_ext_id = _create_persistent_saml_idp(client, env)

    # Inline PEM in .env keeps newlines as literal "\n"; restore real newlines so
    # the written config matches the PEM stored previously.
    cert = env.get("IAM_IDP_METADATA_CERTIFICATE").replace("\\n", "\n")
    block = {
        "custom_attributes": env.list("IAM_IDP_CUSTOM_ATTRIBUTES"),
        "email_attr": env.get("IAM_IDP_EMAIL_ATTR"),
        "entity_issuer": env.get("IAM_IDP_ENTITY_ISSUER"),
        "ext_id": "",
        "groups_attr": env.get("IAM_IDP_GROUPS_ATTR"),
        "groups_delim": env.get("IAM_IDP_GROUPS_DELIM"),
        "idp_metadata": {
            "certificate": cert,
            "entity_id": env.get("IAM_IDP_METADATA_ENTITY_ID"),
            "error_url": env.get("IAM_IDP_METADATA_ERROR_URL"),
            "login_url": env.get("IAM_IDP_METADATA_LOGIN_URL"),
            "logout_url": env.get("IAM_IDP_METADATA_LOGOUT_URL"),
            "name_id_policy_format": env.get("IAM_IDP_METADATA_NAME_ID_POLICY_FORMAT"),
        },
        "idp_metadata_url": metadata_url,
        # Left as the inline XML (usually empty): the create payload uses
        # IAM_IDP_METADATA_XML_FILE, while the tests read the metadata via
        # file(test_idp_metadata.txt) which place_test_idp_metadata() puts in place.
        "idp_metadata_xml": env.get("IAM_IDP_METADATA_XML"),
        "is_signed_authn_req_enabled": env.bool("IAM_IDP_IS_SIGNED_AUTHN"),
        "name": name,
        "username_attr": env.get("IAM_IDP_USERNAME_ATTR"),
    }
    return block, persistent_ext_id


def _create_persistent_saml_idp(client: IamClient, env: Env) -> str:
    """Create (or reuse) the persistent SAML IdP that the SAML user / user-group
    tests reference via iam.users.idp_id. It must outlive the tests (the provider
    cannot delete SAML users, so a test-managed IdP can't be destroyed), and must
    use a different name + entity issuer than IAM_IDP_NAME so it coexists with the
    IdP the SAML-IdP CRUD tests create and destroy."""
    base_name = env.get("IAM_IDP_NAME")
    name = env.get("IAM_USER_IDP_NAME") or (f"{base_name}_persistent" if base_name else "")
    if not name:
        return ""
    base_issuer = env.get("IAM_IDP_ENTITY_ISSUER")
    entity_issuer = env.get("IAM_USER_IDP_ENTITY_ISSUER") or (
        f"{base_issuer}_persistent" if base_issuer else "")
    log(f"IAM: persistent identity provider '{name}' (for iam.users.idp_id)")
    payload = {
        "name": name,
        "usernameAttribute": env.get("IAM_IDP_USERNAME_ATTR"),
        "emailAttribute": env.get("IAM_IDP_EMAIL_ATTR"),
        "groupsAttribute": env.get("IAM_IDP_GROUPS_ATTR"),
        "groupsDelim": env.get("IAM_IDP_GROUPS_DELIM"),
        "entityIssuer": entity_issuer,
        "customAttributes": env.list("IAM_IDP_CUSTOM_ATTRIBUTES"),
        "isSignedAuthnReqEnabled": env.bool("IAM_IDP_IS_SIGNED_AUTHN"),
    }
    # Prefer the structured idpMetadata (entityId + URLs + certificate) -- the same
    # source the SAML-IdP CRUD tests use. A bare idpMetadataXml/Url pointing at an
    # unreachable file leaves the IdP "incorrectly configured", so SAML *user*
    # creation against it fails with IAM-21417.
    cert = (env.get("IAM_IDP_METADATA_CERTIFICATE") or "").replace("\\n", "\n")
    # IAM derives the IdP ext_id from the metadata entityId, so the persistent IdP
    # must use a different entityId than the CRUD-test IdP (IAM_IDP_METADATA_ENTITY_ID)
    # -- otherwise a CRUD test's teardown deletes this IdP too and the SAML user test
    # then fails with IAM-21417 (IdP no longer exists).
    base_entity_id = env.get("IAM_IDP_METADATA_ENTITY_ID")
    entity_id = env.get("IAM_USER_IDP_METADATA_ENTITY_ID") or (
        f"{base_entity_id}/persistent" if base_entity_id else "")
    idp_metadata = {
        "certificate": cert,
        "entityId": entity_id,
        "loginUrl": env.get("IAM_IDP_METADATA_LOGIN_URL"),
        "logoutUrl": env.get("IAM_IDP_METADATA_LOGOUT_URL"),
        "errorUrl": env.get("IAM_IDP_METADATA_ERROR_URL"),
        "nameIdPolicyFormat": env.get("IAM_IDP_METADATA_NAME_ID_POLICY_FORMAT"),
    }
    idp_metadata = {k: v for k, v in idp_metadata.items() if v not in ("", None)}
    if idp_metadata.get("certificate") and idp_metadata.get("entityId") and \
            idp_metadata.get("loginUrl"):
        payload["idpMetadata"] = idp_metadata
    else:
        xml = env.file_or_inline("IAM_IDP_METADATA_XML_FILE", "IAM_IDP_METADATA_XML")
        if xml:
            payload["idpMetadataXml"] = xml
        elif env.get("IAM_IDP_METADATA_URL"):
            payload["idpMetadataUrl"] = env.get("IAM_IDP_METADATA_URL")
    payload = {k: v for k, v in payload.items() if v not in ("", None, [])}
    entity = client.get_or_create("saml-identity-providers", "name", name,
                                  payload, "persistent SAML IdP")
    return entity.get("extId", "")


def _sanitize_ds_name(name: str) -> str:
    """Directory-service `name` rejects dots/special chars (IAM-21011); keep only
    alphanumerics, '-' and '_', mapping everything else to '-'. e.g.
    'qa.nutanix.com' -> 'qa-nutanix-com'."""
    cleaned = "".join(c if (c.isalnum() or c in "-_") else "-" for c in name)
    while "--" in cleaned:
        cleaned = cleaned.replace("--", "-")
    return cleaned.strip("-") or "tf-directory-service"


def build_directory_service(client: IamClient, env: Env, prefix: str,
                            create: bool = True) -> tuple[dict, str]:
    name = env.get(f"{prefix}_NAME")
    if not name:
        return {}, ""
    domain_name = env.get(f"{prefix}_DOMAIN_NAME") or name
    # The config's `name` keeps the dotted domain (vmmv2 uses it as the Windows
    # AD domain), but the directory-service resource `name` must not contain dots.
    service_name = env.get(f"{prefix}_SERVICE_NAME") or _sanitize_ds_name(name)

    # When create=False we only emit the connection values (the directory-service
    # acceptance tests create/destroy this AD themselves, so pre-registering it
    # would make their "create" step fail with "already exists"). Its ext_id is
    # not consumed anywhere, so leaving it empty is safe.
    if not create:
        log(f"IAM: directory service '{service_name}' (values only, not registered)")
        ext_id = ""
        users_map, groups_map = {}, {}
    else:
        log(f"IAM: directory service '{service_name}' (domain '{domain_name}')")
        payload = {
            "name": service_name,
            "url": env.get(f"{prefix}_URL"),
            "directoryType": "ACTIVE_DIRECTORY",
            "domainName": domain_name,
            "serviceAccount": {
                "username": env.get(f"{prefix}_USERNAME"),
                "password": env.get(f"{prefix}_PASSWORD"),
            },
        }
        payload = {k: v for k, v in payload.items() if v not in ("", None)}
        try:
            entity = client.get_or_create("directory-services", "name", service_name, payload,
                                          "directory service", resolvers=[("domainName", domain_name)])
            ext_id = entity.get("extId", "")
        except RuntimeError as exc:
            logger.warning(
                "IAM: directory service '%s' could not be registered; continuing with "
                "values only. Dependent AD user/group ext_ids will be left empty. Error: %s",
                service_name, exc)
            ext_id = ""

        # Role membership targets these AD identities by IAM extId, so they must
        # exist as IAM user/user-group entities (a directory search alone yields
        # an id the role-membership API rejects with IAM-20027). Create the
        # references and record the resulting IAM extIds.
        users_map = {}
        if ext_id:
            for user in env.list(f"{prefix}_USERS"):
                try:
                    users_map[user] = materialize_ad_user(client, ext_id, user)
                except RuntimeError as exc:
                    logger.warning(
                        "IAM: AD user '%s' in directory service '%s' could not be "
                        "materialized; leaving its ext_id empty. Error: %s",
                        user, service_name, exc)
                    users_map[user] = ""
        groups_map = {}
        if ext_id:
            for group in env.list(f"{prefix}_GROUPS"):
                try:
                    groups_map[group] = materialize_ad_group(client, ext_id, group)
                except RuntimeError as exc:
                    logger.warning(
                        "IAM: AD group '%s' in directory service '%s' could not be "
                        "materialized; leaving its ext_id empty. Error: %s",
                        group, service_name, exc)
                    groups_map[group] = ""

    block = {
        "name": name,
        # Directory-service label (no dots; IAM-21011). `name` keeps the dotted
        # domain for vmmv2's Windows domain join.
        "directory_service_name": service_name,
        "username": env.get(f"{prefix}_USERNAME"),
        "password": env.get(f"{prefix}_PASSWORD"),
        "dns": env.get(f"{prefix}_DNS"),
        "ext_id": ext_id,
    }
    if env.get(f"{prefix}_URL"):
        block["url"] = env.get(f"{prefix}_URL")
    if users_map or groups_map:
        block["domain_users_usergroups"] = {"users": users_map, "user_groups": groups_map}
    return block, ext_id


def build_user(client: IamClient, env: Env, idp_ext_id: str,
               primary_ext_id: str, secondary_ext_id: str) -> dict:
    name = env.get("IAM_USER_NAME")
    if not name:
        log("IAM: skipping users block (IAM_USER_NAME not set)")
        return {}
    log(f"IAM: user '{name}'")
    user_type = env.get("IAM_USER_TYPE").upper()
    sources = {"idp": idp_ext_id, "primary": primary_ext_id, "secondary": secondary_ext_id}
    idp_id = sources.get(env.get("IAM_USER_IDP_SOURCE"), idp_ext_id)
    directory_service_id = sources.get(env.get("IAM_USER_DIRECTORY_SOURCE"), secondary_ext_id)

    payload = {"username": name, "userType": user_type}
    if user_type == "LOCAL":
        payload.update({
            # firstName/lastName are required by the IAM v4 API for LOCAL users.
            "firstName": env.get("IAM_USER_FIRST_NAME"),
            "lastName": env.get("IAM_USER_LAST_NAME"),
            "emailId": env.get("IAM_USER_EMAIL"),
            "password": env.get("IAM_USER_PASSWORD"),
            "locale": env.get("IAM_USER_LOCALE"),
            "region": env.get("IAM_USER_REGION"),
            "isForceResetPasswordEnabled": env.bool("IAM_USER_FORCE_RESET"),
        })
    else:
        payload["idpId"] = directory_service_id
    payload = {k: v for k, v in payload.items() if v not in ("", None)}

    try:
        entity = client.get_or_create("users", "username", name, payload, "user")
        ext_id = entity.get("extId", "")
    except RuntimeError as exc:
        logger.warning(
            "IAM: user '%s' could not be created/resolved; continuing with values "
            "only and leaving ext_id empty. Error: %s", name, exc)
        ext_id = ""
    return {
        "name": name,
        "idp_id": idp_id,
        "directory_service_id": directory_service_id,
        "directory_service_username": env.get("IAM_USER_DIRECTORY_SERVICE_USERNAME"),
        "email_id": env.get("IAM_USER_EMAIL"),
        "locale": env.get("IAM_USER_LOCALE"),
        "region": env.get("IAM_USER_REGION"),
        "password": env.get("IAM_USER_PASSWORD"),
        "force_reset_password": env.bool("IAM_USER_FORCE_RESET"),
        "ext_id": ext_id,
    }


def build_user_group(client: IamClient, env: Env,
                     primary_ext_id: str, secondary_ext_id: str, idp_ext_id: str) -> dict:
    name = env.get("IAM_GROUP_NAME")
    if not name:
        log("IAM: skipping user_groups block (IAM_GROUP_NAME not set)")
        return {}
    log(f"IAM: user group '{name}'")
    sources = {"idp": idp_ext_id, "primary": primary_ext_id, "secondary": secondary_ext_id}
    idp_id = sources.get(env.get("IAM_GROUP_IDP_SOURCE"), secondary_ext_id)
    distinguished_name = env.get("IAM_GROUP_DISTINGUISHED_NAME")

    payload = {
        "groupType": env.get("IAM_GROUP_TYPE").upper(),
        "idpId": idp_id,
        "name": name,
    }
    if distinguished_name:
        payload["distinguishedName"] = distinguished_name
    payload = {k: v for k, v in payload.items() if v not in ("", None)}

    try:
        entity = client.get_or_create("user-groups", "name", name, payload, "user group")
        ext_id = entity.get("extId", "")
    except RuntimeError as exc:
        logger.warning(
            "IAM: user group '%s' could not be created/resolved; continuing with "
            "values only and leaving ext_id empty. Error: %s", name, exc)
        ext_id = ""
    return {
        "name": name,
        "saml_name": env.get("IAM_GROUP_SAML_NAME") or name,
        "distinguished_name": distinguished_name,
        "ext_id": ext_id,
    }


def build_openldap_directory(client: IamClient, env: Env) -> str:
    """Provision (or reuse) an OPEN_LDAP directory service and return its ext_id.
    Returns '' when IAM_OPENLDAP_NAME is unset (feature disabled). Mirrors the
    open_ldap directory in preEnv/iam.tf."""
    name = env.get("IAM_OPENLDAP_NAME")
    if not name:
        return ""
    domain = env.get("IAM_OPENLDAP_DOMAIN_NAME") or name
    service_name = env.get("IAM_OPENLDAP_SERVICE_NAME") or _sanitize_ds_name(name)
    log(f"IAM: OpenLDAP directory '{service_name}' (domain '{domain}')")
    payload = {
        "name": service_name,
        "url": env.get("IAM_OPENLDAP_URL"),
        "directoryType": "OPEN_LDAP",
        "domainName": domain,
        "serviceAccount": {
            "username": env.get("IAM_OPENLDAP_USERNAME"),
            "password": env.get("IAM_OPENLDAP_PASSWORD"),
        },
        "openLdapConfiguration": {
            "userConfiguration": {
                "userObjectClass": env.get("IAM_OPENLDAP_USER_OBJECT_CLASS"),
                "userSearchBase": env.get("IAM_OPENLDAP_USER_SEARCH_BASE"),
                "usernameAttribute": env.get("IAM_OPENLDAP_USERNAME_ATTR"),
            },
            "userGroupConfiguration": {
                "groupObjectClass": env.get("IAM_OPENLDAP_GROUP_OBJECT_CLASS"),
                "groupSearchBase": env.get("IAM_OPENLDAP_GROUP_SEARCH_BASE"),
                "groupMemberAttribute": env.get("IAM_OPENLDAP_GROUP_MEMBER_ATTR"),
                "groupMemberAttributeValue": env.get("IAM_OPENLDAP_GROUP_MEMBER_ATTR_VALUE"),
            },
        },
        "groupSearchType": env.get("IAM_OPENLDAP_GROUP_SEARCH_TYPE"),
    }
    payload = _prune_empty(payload)
    try:
        entity = client.get_or_create("directory-services", "name", service_name, payload,
                                      "OpenLDAP directory", resolvers=[("domainName", domain)])
        return entity.get("extId", "")
    except RuntimeError as exc:
        logger.warning(
            "IAM: OpenLDAP directory '%s' could not be registered; continuing without "
            "OpenLDAP-discovered user/group values. Error: %s", service_name, exc)
        return ""


def discover_ldap_group(client: IamClient, ds_ext_id: str) -> dict | None:
    """Find the first real group in a directory (cn + dn). Mirrors the group
    search in preEnv/iam.tf so the created user-group always exists in the dir.
    The literal "Tenant.Group" is the wildcard token OpenLDAP expects here."""
    for entity in client.search_directory(ds_ext_id, "Tenant.Group"):
        etype = (entity.get("entityType") or "").lower()
        if etype and etype != "group":
            continue
        cn, dn = attr(entity, "cn"), attr(entity, "dn")
        if cn and dn:
            return {"name": cn.lower(), "distinguished_name": dn.lower()}
    return None


def discover_ldap_user(client: IamClient, ds_ext_id: str) -> dict | None:
    """Find the first real user in a directory (uid/cn).
    The literal "Tenant.Group.User" is the wildcard token OpenLDAP expects here."""
    for entity in client.search_directory(ds_ext_id, "Tenant.Group.User"):
        etype = (entity.get("entityType") or "").lower()
        if etype and etype != "person":
            continue
        uid = attr(entity, "cn") or attr(entity, "uid")
        if uid:
            return {"uid": uid}
    return None


def build_openldap_user(client: IamClient, env: Env, openldap_ext_id: str,
                        idp_ext_id: str, user: dict) -> dict:
    """Fill the iam.users block from a user discovered in the directory.

    The user is NOT created in IAM: the acceptance tests create (and destroy)
    it themselves, so pre-creating it here would make those tests fail with a
    409 conflict. ext_id is left empty because no test consumes it.
    """
    uid = user["uid"]
    domain = env.get("IAM_OPENLDAP_DOMAIN_NAME") or env.get("IAM_OPENLDAP_NAME")
    log(f"IAM: OpenLDAP user '{uid}' (values only, not created)")
    username_at_domain = f"{uid}@{domain}" if domain else uid
    return {
        "name": uid.lower(),
        "idp_id": idp_ext_id,
        "directory_service_id": openldap_ext_id,
        "directory_service_username": username_at_domain,
        "email_id": username_at_domain,
        "locale": env.get("IAM_USER_LOCALE"),
        "region": env.get("IAM_USER_REGION"),
        "password": env.get("IAM_USER_PASSWORD"),
        "force_reset_password": env.bool("IAM_USER_FORCE_RESET"),
        "ext_id": "",
    }


def build_openldap_group(client: IamClient, env: Env, openldap_ext_id: str,
                         group: dict) -> dict:
    """Fill the iam.user_groups block from a group discovered in the directory.

    The group is NOT created in IAM: the acceptance tests create (and destroy)
    it themselves, so pre-creating it here would make those tests fail with a
    409 conflict (group already exists with same DN). ext_id is left empty
    because no test consumes it.
    """
    name = group["name"]
    dn = group["distinguished_name"]
    log(f"IAM: OpenLDAP user group '{name}' (values only, not created)")
    return {
        "name": name,
        "saml_name": name,
        "distinguished_name": dn,
        "ext_id": "",
    }


def build_iam(env: Env, dry_run: bool) -> dict:
    client = IamClient(env, dry_run=dry_run)
    log(f"Connecting to IAM API at {client.base()} (insecure={client.insecure})")

    idp_block, idp_ext_id = build_identity_provider(client, env)
    # The primary AD doubles as the standalone iam.directory_services entity that
    # the directory-service CRUD tests create/destroy, so don't pre-register it
    # (its ext_id is unused). The secondary AD must be registered: the role
    # membership test consumes secondary_ad.ext_id.
    primary_block, primary_ext_id = build_directory_service(client, env, "IAM_PRIMARY_AD",
                                                            create=False)
    secondary_block, secondary_ext_id = build_directory_service(client, env, "IAM_SECONDARY_AD")
    # Extra AD kept in the config for convenience only; nothing consumes it yet, so
    # emit values-only (create=False) and don't register it in IAM.
    tertiary_block, _tertiary_ext_id = build_directory_service(client, env, "IAM_TERTIARY_AD",
                                                               create=False)
    # The role membership test creates memberships for secondary-AD users/groups
    # against a throwaway project. IAM requires the IDP (directory service) to be
    # shared with the project first (IAM-20027); sharing with all projects makes
    # any project usable.
    if secondary_ext_id:
        client.share_directory_service_with_all_projects(secondary_ext_id)

    # If an OpenLDAP directory is configured, discover a real user/group from it
    # (read-only search; mirrors preEnv/iam.tf) and fill their values into the
    # config. They are NOT created in IAM here -- the acceptance tests create and
    # destroy them, so pre-creating would cause 409 conflicts. Otherwise fall
    # back to the IAM_USER_* / IAM_GROUP_* values.
    openldap_ext_id = build_openldap_directory(client, env)
    if openldap_ext_id and not client.dry_run:
        group = discover_ldap_group(client, openldap_ext_id)
        user = discover_ldap_user(client, openldap_ext_id)
        if not group:
            raise RuntimeError(f"OpenLDAP directory {openldap_ext_id} returned no groups to "
                               f"fill the user group from (check the LDAP server / search base)")
        if not user:
            raise RuntimeError(f"OpenLDAP directory {openldap_ext_id} returned no users to "
                               f"fill the user from (check the LDAP server / search base)")
        user_block = build_openldap_user(client, env, openldap_ext_id, idp_ext_id, user)
        group_block = build_openldap_group(client, env, openldap_ext_id, group)
    elif openldap_ext_id and client.dry_run:
        log("  - [dry-run] would discover an OpenLDAP user/group and fill their values")
        user_block, group_block = {}, {}
    else:
        user_block = build_user(client, env, idp_ext_id, primary_ext_id, secondary_ext_id)
        group_block = build_user_group(client, env, primary_ext_id, secondary_ext_id, idp_ext_id)

    iam: dict = {}
    dsm: dict = {}
    if primary_block:
        dsm["primary_ad"] = primary_block
    if secondary_block:
        dsm["secondary_ad"] = secondary_block
    if tertiary_block:
        dsm["tertiary_ad"] = tertiary_block
    if dsm:
        iam["directory_services_main"] = dsm
    if user_block:
        iam["users"] = user_block
    if group_block:
        iam["user_groups"] = group_block
    if idp_block:
        iam["identity_providers"] = idp_block
    return iam


def place_test_idp_metadata(env: Env, dry_run: bool) -> None:
    """Ensure the IdP metadata file the acceptance tests read via file(...) is in
    place. The iamv2 tests load it from the repo root as test_idp_metadata.txt
    (see nutanix/services/iamv2/main_test.go); copy IAM_IDP_METADATA_TEST_FILE
    there so the SAML IdP tests have their metadata."""
    src_raw = env.get("IAM_IDP_METADATA_TEST_FILE")
    if not src_raw:
        return
    src = Path(src_raw)
    if not src.is_absolute():
        src = (REPO_ROOT / src_raw).resolve()
    dest = REPO_ROOT / "test_idp_metadata.txt"
    if not src.exists():
        logger.warning("IAM_IDP_METADATA_TEST_FILE not found: %s (skipping)", src)
        return
    if src.resolve() == dest.resolve():
        return
    if dry_run:
        log(f"  - [dry-run] would copy IdP test metadata {src} -> {dest}")
        return
    shutil.copyfile(src, dest)
    log(f"Placed IdP test metadata at {dest} (from {src})")


# --------------------------------------------------------------------------- #
# SSL certificate generation (self-contained; requires the openssl CLI)
# --------------------------------------------------------------------------- #
def _b64_file(path: Path) -> str:
    return base64.b64encode(path.read_bytes()).decode()


_IPV4_RE = re.compile(r"^\d{1,3}(\.\d{1,3}){3}$")


def _openssl(args: list, cwd: str) -> None:
    """Run an openssl subcommand, raising with captured output on failure."""
    proc = subprocess.run(["openssl", *args], cwd=cwd, capture_output=True, text=True)
    logger.debug("openssl %s -> exit=%s\n  stderr: %s", args[0], proc.returncode,
                 _truncate(proc.stderr))
    if proc.returncode != 0:
        raise RuntimeError(f"openssl {args[0]} failed (exit {proc.returncode}):\n"
                           f"{proc.stdout}\n{proc.stderr}")


def generate_ssl_certificate(env: Env, dry_run: bool) -> dict:
    """Generate an RSA-2048 server cert and return the base64-encoded material for
    the clusters.ssl_certificate block.

    Self-contained: creates a throwaway single-tier CA and signs a leaf cert for
    the cluster VIP/FQDN with it, entirely inside a temp directory using the
    ``openssl`` CLI. The leaf cert is uploaded to the cluster together with the
    generated CA cert as the ``ca_chain`` -- the cluster only needs a key + cert +
    chain that verify against each other, so a fresh per-run CA is sufficient and
    avoids depending on any committed CA material.
    """
    nodes = env.list("CLUSTER_NODES")
    fqdn = env.get("CLUSTER_SSL_FQDN") or env.get("CLUSTER_VIRTUAL_IP") or (nodes[0] if nodes else "")
    passphrase = env.get("CLUSTER_SSL_PASSPHRASE")

    if not fqdn:
        raise SystemExit("ERROR: SSL generation needs a FQDN/IP (set cluster.ssl.fqdn or cluster.virtual_ip)")
    if dry_run:
        log(f"  - [dry-run] would generate a fresh rsa2048 cert + CA for '{fqdn}' via openssl")
        return {"passphrase": passphrase, "private_key": "", "public_certificate": "", "ca_chain": ""}
    if shutil.which("openssl") is None:
        raise SystemExit("ERROR: SSL generation requires the 'openssl' CLI, which was not "
                         "found in PATH.")

    log(f"SSL: generating rsa2048 cert for '{fqdn}' with a throwaway CA")
    subj_base = "/C=US/ST=California/L=San Jose/O=Nutanix/OU=Test"
    san = f"{'IP' if _IPV4_RE.match(fqdn) else 'DNS'}:{fqdn}"

    with tempfile.TemporaryDirectory(prefix="ssl_gen_") as tmp:
        ca_key = "ca_key.pem"
        ca_crt = "ca_cert.pem"
        leaf_key = "leaf_key.pem"
        leaf_csr = "leaf.csr"
        leaf_crt = "leaf_cert.pem"

        # 1. Self-signed CA (acts as the single-tier issuer / chain).
        _openssl(["req", "-x509", "-newkey", "rsa:2048", "-nodes",
                  "-keyout", ca_key, "-out", ca_crt, "-days", "3650", "-sha256",
                  "-subj", f"{subj_base}/CN=Nutanix Test CA"], cwd=tmp)

        # 2. Leaf key + CSR for the cluster VIP/FQDN.
        _openssl(["req", "-newkey", "rsa:2048", "-nodes",
                  "-keyout", leaf_key, "-out", leaf_csr,
                  "-subj", f"{subj_base}/CN={fqdn}"], cwd=tmp)

        # 3. Sign the leaf with the CA, adding SAN + server-auth extensions.
        ext_path = Path(tmp) / "leaf_ext.cnf"
        ext_path.write_text(
            "basicConstraints=CA:FALSE\n"
            "keyUsage=digitalSignature,keyEncipherment\n"
            "extendedKeyUsage=serverAuth\n"
            f"subjectAltName={san}\n"
        )
        _openssl(["x509", "-req", "-in", leaf_csr, "-CA", ca_crt, "-CAkey", ca_key,
                  "-CAcreateserial", "-out", leaf_crt, "-days", "825", "-sha256",
                  "-extfile", "leaf_ext.cnf"], cwd=tmp)

        return {
            "passphrase": passphrase,
            "private_key": _b64_file(Path(tmp) / leaf_key),
            "public_certificate": _b64_file(Path(tmp) / leaf_crt),
            "ca_chain": _b64_file(Path(tmp) / ca_crt),
        }


# --------------------------------------------------------------------------- #
# Cluster network discovery (Prism Element VIP + iSCSI data services IP)
# --------------------------------------------------------------------------- #
def _coerce_ip(value) -> str:
    """Reduce a VIP/iSCSI field to a plain IPv4 string.

    The v2.0 cluster API may return the address as a string, a nested
    ``{"ipv4": ...}`` object, or a single-element list of such objects. The test
    config and Go struct expect a plain string.
    """
    if not value:
        return ""
    if isinstance(value, str):
        return value
    if isinstance(value, list):
        return _coerce_ip(value[0]) if value else ""
    if isinstance(value, dict):
        return value.get("ipv4") or value.get("ipv6") or ""
    return str(value)


_MAC_RE = re.compile(r"(?:[0-9a-fA-F]{1,2}[:-]){5}[0-9a-fA-F]{1,2}")


def _ping(ip: str, timeout: int = 3) -> bool:
    """Return True if *ip* answers a single ICMP echo. Cross-platform (the -W
    reply-wait flag is milliseconds on macOS, seconds on Linux)."""
    if sys.platform == "darwin":
        argv = ["ping", "-c", "1", "-n", "-W", "1500", ip]
    else:
        argv = ["ping", "-c", "1", "-n", "-w", "2", ip]
    try:
        proc = subprocess.run(argv, capture_output=True, timeout=timeout)
        return proc.returncode == 0
    except (subprocess.TimeoutExpired, OSError):
        return False


def _arp_has_entry(ip: str, timeout: int = 3) -> bool:
    """Return True if the local ARP table has a resolved MAC for *ip* (i.e. the
    address is live on this L2 segment). 'incomplete'/'no entry' -> False. Only
    meaningful when this host shares the broadcast domain with the target."""
    try:
        proc = subprocess.run(["arp", "-n", ip], capture_output=True, text=True, timeout=timeout)
    except (subprocess.TimeoutExpired, OSError):
        return False
    out = proc.stdout or ""
    if "incomplete" in out.lower():
        return False
    return bool(_MAC_RE.search(out))


def _ip_is_free(ip: str) -> bool:
    """An IP is considered free when it neither answers ping nor has a resolved
    ARP entry. (ARP is best-effort: on a routed/VPN path there will be no ARP
    entries, so freeness rests on the ping result.)"""
    if _ping(ip):
        logger.debug("candidate %s answers ping -> in use", ip)
        return False
    if _arp_has_entry(ip):
        logger.debug("candidate %s has a resolved ARP entry -> in use", ip)
        return False
    logger.debug("candidate %s is free (no ping reply, no ARP entry)", ip)
    return True


def _iter_ips_after(anchor_ip: str):
    """Yield host IPs in *anchor_ip*'s /24, starting just above it and wrapping to
    the addresses below it, skipping the network (.0) and broadcast (.255)."""
    parts = anchor_ip.split(".")
    if len(parts) != 4 or not all(p.isdigit() for p in parts):
        return
    base = ".".join(parts[:3])
    last = int(parts[3])
    for octet in list(range(last + 1, 255)) + list(range(1, last)):
        yield f"{base}.{octet}"


def _find_free_ips(anchor_ip: str, count: int, exclude: set, label: str) -> list:
    """Return up to *count* free IPs walking upward from *anchor_ip* within its
    /24 (see _iter_ips_after), skipping anything in *exclude* and adding each pick
    to *exclude* so callers don't hand out duplicates. Free is judged by ping+ARP.

    A quick reachability sanity check pings the anchor first: if even the PE IP
    doesn't answer, free detection can't be trusted (nothing is reachable), so we
    warn -- the picks may be wrong until the host/subnet is reachable.
    """
    if count <= 0:
        return []
    if not _ping(anchor_ip):
        logger.warning("%s: anchor %s does not answer ping; free-IP detection may be "
                       "unreliable (subnet not reachable from here)", label, anchor_ip)
    found = []
    for cand in _iter_ips_after(anchor_ip):
        if cand in exclude:
            continue
        if _ip_is_free(cand):
            found.append(cand)
            exclude.add(cand)
            if len(found) >= count:
                break
    if len(found) < count:
        logger.warning("%s: only found %d/%d free IP(s) after %s in its /24",
                       label, len(found), count, anchor_ip)
    return found


def discover_cluster_network(env: Env) -> dict:
    """Pick the cluster external (virtual) IP and external data-services (iSCSI)
    IP as the first two *free* IPs after the PE IP (CLUSTER_NODES[0]) within its
    /24 -- freeness judged by ping + ARP (see _find_free_ips). Explicit
    CLUSTER_VIRTUAL_IP / CLUSTER_ISCSI_IP values win per field. Returns
    {"virtual_ip": ..., "iscsi_ip": ...}.
    """
    vip = env.get("CLUSTER_VIRTUAL_IP")
    dsip = env.get("CLUSTER_ISCSI_IP")
    if vip and dsip:
        log(f"clusters.network: using configured virtual_ip={vip} iscsi_ip={dsip}")
        return {"virtual_ip": vip, "iscsi_ip": dsip}

    nodes = [n for n in env.list("CLUSTER_NODES") if n]
    if not nodes:
        logger.warning("clusters.network: CLUSTER_NODES is empty; using configured values")
        return {"virtual_ip": vip, "iscsi_ip": dsip}
    pe = nodes[0]

    exclude = set(nodes)
    needed = 2 - sum(1 for v in (vip, dsip) if v)
    log(f"clusters.network: finding {needed} free IP(s) after PE {pe} (virtual_ip, iscsi_ip)")
    free = _find_free_ips(pe, needed, exclude, "clusters.network")
    free_iter = iter(free)
    if not vip:
        vip = next(free_iter, "")
    if not dsip:
        dsip = next(free_iter, "")

    if not vip:
        logger.warning("clusters.network: could not pick a free virtual_ip; left empty")
    if not dsip:
        logger.warning("clusters.network: could not pick a free iscsi_ip; left empty")
    log(f"clusters.network: virtual_ip={vip or '<unset>'} iscsi_ip={dsip or '<unset>'}")
    return {"virtual_ip": vip, "iscsi_ip": dsip}


def _group_attr(entity: dict, name: str) -> str:
    """Extract a single attribute value from a v3 groups entity_results item."""
    for field in entity.get("data") or []:
        if field.get("name") != name:
            continue
        for v in field.get("values") or []:
            vals = v.get("values") or []
            if vals:
                return vals[0]
    return ""


def _flow_kube_clusters(host: str, port: str, token: str, ctx) -> list:
    """Return the Kubernetes clusters known to the networking/Flow subsystem via
    the v3 groups API over the flow_kube_cluster_config IDF entity. Only clusters
    ACTIVATED for Flow networking appear here, and only their cluster_uuid is a
    valid VPC kubernetesClusters.extId. Returns [{'name':..., 'uuid':...}, ...]
    (empty on any error or when none are activated)."""
    body = {
        "entity_type": "flow_kube_cluster_config",
        "group_member_count": 100,
        "group_member_attributes": [
            {"attribute": "cluster_uuid"},
            {"attribute": "name"},
        ],
    }
    try:
        status, resp = _pc_v3_post(host, port, token, ctx, "api/nutanix/v3/groups", body)
    except Exception as exc:  # noqa: BLE001 - connection/parse issues are non-fatal
        logger.warning("flow_kube_cluster_config groups query failed (%s)", exc)
        return []
    if status // 100 != 2:
        logger.warning("flow_kube_cluster_config groups query -> HTTP %s", status)
        return []
    out = []
    for group in resp.get("group_results") or []:
        for entity in group.get("entity_results") or []:
            uuid = _group_attr(entity, "cluster_uuid") or entity.get("entity_id", "")
            if uuid:
                out.append({"name": _group_attr(entity, "name"), "uuid": uuid})
    return out


def discover_kubernetes_cluster_ext_id(env: Env) -> str:
    """Return the extId of a Kubernetes cluster registered with Prism Central,
    for networking.kubernetes_cluster_ext_id (consumed by the VPC v2
    kubernetes_clusters association test).

    An explicit NETWORKING_KUBERNETES_CLUSTER_EXT_ID always wins. Otherwise the
    cluster is discovered from PC's registration API (the same one served after
    an NKP/Konnector self-registration, e.g. via testenv/e2e_nkp_deploy.py). When
    NETWORKING_KUBERNETES_CLUSTER_NAME is set, the cluster with that name is
    selected; otherwise the first registered cluster is used.

    Returns "" (and logs a warning) when PC creds are missing or no cluster is
    registered -- the caller omits the key so an existing value is preserved and
    the VPC test self-skips.
    """
    override = env.get("NETWORKING_KUBERNETES_CLUSTER_EXT_ID")
    if override:
        log(f"networking.kubernetes_cluster_ext_id: using configured {override}")
        return override

    host = env.get("PC_ENDPOINT")
    port = env.get("PC_PORT")
    username = env.get("PC_USERNAME")
    password = env.get("PC_PASSWORD")
    if not (host and port and username and password):
        logger.warning("networking.kubernetes_cluster_ext_id: PC_ENDPOINT/PORT/"
                       "USERNAME/PASSWORD not set; leaving it to the existing value")
        return ""

    wanted = env.get("NETWORKING_KUBERNETES_CLUSTER_NAME")
    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    token = base64.b64encode(f"{username}:{password}".encode()).decode()

    # AUTHORITATIVE source first: the flow_kube_cluster_config IDF entity is what
    # the networking/Flow subsystem validates VPC kubernetesClusters against. A
    # cluster only appears here once it has been ACTIVATED for Flow networking
    # (Kubernetes Clusters page -> Actions -> Activate Nutanix Flow Capabilities).
    # A cluster that is merely PC-registered (visible in the k8s UI) but not
    # Flow-activated is NOT here, and referencing it makes VPC create fail with
    # "K8s Cluster <uuid> not found". So prefer these extIds.
    flow_clusters = _flow_kube_clusters(host, port, token, ctx)
    if flow_clusters:
        chosen = None
        if wanted:
            chosen = next((c for c in flow_clusters if c.get("name") == wanted), None)
        chosen = chosen or flow_clusters[0]
        if chosen.get("uuid"):
            log(f"networking.kubernetes_cluster_ext_id: {chosen.get('name', '<unnamed>')} "
                f"-> {chosen['uuid']} (Flow-activated, from flow_kube_cluster_config)")
            return chosen["uuid"]
    else:
        logger.warning("networking.kubernetes_cluster_ext_id: no Flow-activated clusters "
                       "found (flow_kube_cluster_config empty); falling back to the "
                       "registration API -- note such an extId may be rejected by VPC "
                       "create until the cluster is activated for Flow networking")

    # Karbon (NKE) endpoint first -- it serves Konnector self-registrations today
    # (returns 'uuid'); the v4 NKE paths are forward-compat fallbacks ('extId').
    endpoints = [
        "/karbon/v1-alpha.1/k8s/cluster-registrations",
        "/api/nke/v4.0.b1/config/cluster-registrations",
        "/api/nke/v4.0.a1/config/cluster-registrations",
    ]
    for path in endpoints:
        try:
            body = _pc_get_json(host, port, path, token, ctx)
        except urllib.error.HTTPError as exc:
            if exc.code == 404:
                continue
            logger.warning("networking.kubernetes_cluster_ext_id: %s -> HTTP %s",
                           path, exc.code)
            continue
        except Exception as exc:  # noqa: BLE001 - connection/parse issues are non-fatal
            logger.warning("networking.kubernetes_cluster_ext_id: %s failed (%s)", path, exc)
            continue

        if isinstance(body, list):
            records = [c for c in body if isinstance(c, dict)]
        else:
            records = body.get("data") or body.get("clusters") or []
            records = [c for c in records if isinstance(c, dict)]
        if not records:
            continue

        chosen = None
        if wanted:
            chosen = next((c for c in records if c.get("name") == wanted), None)
            if chosen is None:
                logger.warning("networking.kubernetes_cluster_ext_id: cluster "
                               "'%s' not found via %s", wanted, path)
                continue
        else:
            chosen = records[0]

        for key in ("extId", "uuid", "KubernetesClusterUUID", "kubernetesClusterUuid"):
            ext_id = chosen.get(key)
            if ext_id:
                log(f"networking.kubernetes_cluster_ext_id: "
                    f"{chosen.get('name', '<unnamed>')} -> {ext_id}")
                return ext_id

    logger.warning("networking.kubernetes_cluster_ext_id: no registered Kubernetes "
                   "cluster found; leaving it to the existing value")
    return ""


def _pc_get_json(host: str, port: str, path: str, token: str, ctx,
                 timeout: int = 30) -> dict:
    """GET a JSON document from a PC v4 API. Returns {} on an empty body."""
    url = f"https://{host}:{port}{path}"
    req = urllib.request.Request(url, method="GET")
    req.add_header("Authorization", f"Basic {token}")
    req.add_header("Accept", "application/json")
    logger.debug("HTTP GET %s", url)
    with urllib.request.urlopen(req, context=ctx, timeout=timeout) as resp:
        raw = resp.read().decode()
    logger.debug("response %s: %s", path, _safe_json(raw))
    return json.loads(raw) if raw.strip() else {}


def _discover_cluster_ext_id(host: str, port: str, token: str, ctx, cluster_function: str) -> str:
    """Return the extId of the first cluster on *host* whose clusterFunction
    matches *cluster_function* (e.g. 'AOS' or 'PRISM_CENTRAL').

    Mirrors the preEnv clusters.tf data sources, which filter
    nutanix_clusters_v2 by clusterFunction and take cluster_entities[0].ext_id.
    """
    fltr = (f"config/clusterFunction/any(t:t eq "
            f"Clustermgmt.Config.ClusterFunctionRef'{cluster_function}')")
    path = "/api/clustermgmt/v4.0/config/clusters?$filter=" + urllib.parse.quote(fltr, safe="")
    data = _pc_get_json(host, port, path, token, ctx)
    items = data.get("data") or []
    if not items:
        return ""
    return items[0].get("extId") or ""


def _pc_v3_post(host: str, port: str, token: str, ctx, path: str, body: dict,
                timeout: int = 60) -> tuple:
    """POST a v3 JSON request to a PC. Returns (status_code, parsed_body). A
    non-2xx HTTPError is captured and returned as (code, {...}) rather than
    raised; connection/parse errors propagate to the caller."""
    url = f"https://{host}:{port}/{path}"
    data = json.dumps(body).encode()
    req = urllib.request.Request(url, data=data, method="POST")
    req.add_header("Authorization", f"Basic {token}")
    req.add_header("Accept", "application/json")
    req.add_header("Content-Type", "application/json")
    logger.debug("HTTP POST %s", url)
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=timeout) as resp:
            raw = resp.read().decode()
            status = getattr(resp, "status", resp.getcode())
    except urllib.error.HTTPError as exc:
        detail = exc.read().decode(errors="replace")
        logger.debug("HTTP POST %s -> %d: %s", url, exc.code, _truncate(detail))
        return exc.code, {}
    logger.debug("response %s: %s", path, _safe_json(raw) if raw.strip() else "<empty>")
    return status, (json.loads(raw) if raw.strip() else {})


def _remote_pc_az_identity(remote_ip: str, port: str, token: str, ctx) -> tuple:
    """Return (cluster_ext_id, display_name) of the remote PC's own 'Local AZ'
    (the availability_zones/list entity whose management_plane_type == Local).

    That cluster_ext_id is what the local PC records as ``management_url`` for the
    paired AZ, and display_name is the ``PC_<advertised-ip>`` label. Both survive
    the connect-IP vs advertised-IP mismatch. Returns ("", "") on any failure."""
    try:
        status, body = _pc_v3_post(remote_ip, port, token, ctx,
                                   "api/nutanix/v3/availability_zones/list",
                                   {"kind": "availability_zone", "length": 200})
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.debug("remote PC %s AZ identity lookup failed: %s", remote_ip,
                     getattr(exc, "reason", exc))
        return "", ""
    if status not in (0, 200) or not isinstance(body, dict):
        return "", ""
    for ent in (body.get("entities") or []):
        res = (ent.get("status") or {}).get("resources") or {}
        if str(res.get("management_plane_type") or "").lower() == "local":
            return res.get("management_url") or "", res.get("display_name") or ""
    return "", ""


def _az_zone_present(body: dict, remote_ip: str, remote_ext_id: str,
                     remote_disp: str) -> tuple:
    """Scan an availability_zones/list body for a non-Local zone matching the
    remote PC. Returns (matched, state). Matching prefers the canonical
    remote_ext_id / display_name, falling back to remote_ip substring matching."""
    for ent in (body.get("entities") or []) if isinstance(body, dict) else []:
        res = (ent.get("status") or {}).get("resources") or {}
        if str(res.get("management_plane_type") or "").lower() == "local":
            continue
        name = (ent.get("status") or {}).get("name") or ""
        state = (ent.get("status") or {}).get("state") or ""
        mgmt_url = str(res.get("management_url") or "")
        disp = str(res.get("display_name") or "")
        if ((remote_ext_id and mgmt_url == remote_ext_id)
                or (remote_disp and disp == remote_disp)
                or (remote_ip and (remote_ip in mgmt_url or remote_ip in disp
                                   or remote_ip in name))):
            return True, state
    return False, ""


def _ensure_az_connected(env: Env, remote_ip: str, remote_user: str,
                         remote_pass: str, label: str = "prism.unregister") -> bool:
    """Ensure the local PC (PC_ENDPOINT) is availability-zone paired with
    *remote_ip*; create the pairing via a v3 cloud_trust (ONPREM_CLOUD) when it
    isn't, then poll until it connects.

    Mirrors prepare_pc.prepare_prism so tests that need a connected remote PC (the
    unregistration test, and the protection_policy destination AZ) have one.
    Non-fatal: logs and returns False on any failure (discovery must not abort the
    rest of the fill). *label* only prefixes the log lines."""
    local_ep = env.get("PC_ENDPOINT")
    local_user = env.get("PC_USERNAME")
    local_pass = env.get("PC_PASSWORD")
    port = env.get("PC_PORT")
    if not (local_ep and local_user and local_pass and port):
        logger.warning("%s: PC_ENDPOINT/USERNAME/PASSWORD/PORT not set; "
                       "cannot verify the availability-zone connection to %s", label, remote_ip)
        return False

    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    local_token = base64.b64encode(f"{local_user}:{local_pass}".encode()).decode()
    remote_token = base64.b64encode(f"{remote_user}:{remote_pass}".encode()).decode()

    # Resolve the remote PC's canonical identity so detection survives the
    # connect-IP vs advertised-IP mismatch (the paired AZ is labelled by the
    # advertised IP, not the IP we connect to).
    remote_ext_id, remote_disp = _remote_pc_az_identity(remote_ip, port, remote_token, ctx)

    def _connected() -> bool:
        try:
            status, body = _pc_v3_post(local_ep, port, local_token, ctx,
                                       "api/nutanix/v3/availability_zones/list",
                                       {"kind": "availability_zone", "length": 200})
        except (urllib.error.URLError, ValueError, OSError) as exc:
            logger.warning("%s: availability_zones/list on %s failed (%s)",
                           label, local_ep, getattr(exc, "reason", exc))
            return False
        if status not in (0, 200):
            logger.warning("%s: availability_zones/list on %s -> HTTP %s",
                           label, local_ep, status)
            return False
        matched, state = _az_zone_present(body, remote_ip, remote_ext_id, remote_disp)
        if matched and state and state.upper() != "COMPLETE":
            logger.warning("%s: AZ to %s present but state=%s (expected "
                           "COMPLETE); the pairing may still be settling", label, remote_ip, state)
        return matched

    if _connected():
        log(f"{label}: availability zone local PC {local_ep} -> {remote_ip} "
            f"already connected")
        return True

    if not (remote_user and remote_pass):
        logger.warning("%s: no availability zone to %s and remote PC creds "
                       "unavailable; cannot connect", label, remote_ip)
        return False

    log(f"{label}: no availability zone to {remote_ip}; connecting via "
        f"cloud_trust (ONPREM_CLOUD) ...")
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
    try:
        status, _ = _pc_v3_post(local_ep, port, local_token, ctx,
                                "api/nutanix/v3/cloud_trusts", payload)
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.warning("%s: cloud_trust create to %s failed (%s)",
                       label, remote_ip, getattr(exc, "reason", exc))
        return False
    if status not in (200, 201, 202):
        logger.warning("%s: cloud_trust create to %s -> HTTP %s", label, remote_ip, status)
        return False

    deadline = time.monotonic() + 300
    while time.monotonic() < deadline:
        time.sleep(10)
        if _connected():
            log(f"{label}: availability zone to {remote_ip} connected")
            return True
        logger.info("%s: AZ to %s not connected yet; retrying", label, remote_ip)
    logger.warning("%s: timed out waiting for the AZ to %s to connect", label, remote_ip)
    return False


def discover_availability_zone(env: Env) -> dict:
    """Fill the availability_zone block, which points at a *remote* PC used as a
    replication target by the data-protection / data-policies tests.

    The only required input is the remote PC IP (AZ_REMOTE_PC_IP). We connect to
    that PC (PC credentials, or the AZ_REMOTE_PC_USERNAME/PASSWORD overrides when
    the remote PC uses different ones) and discover, via the v4 clusters API --
    exactly like preEnv/clusters.tf:
      * pc_ext_id      - extId of the PRISM_CENTRAL cluster entity on that PC
      * cluster_ext_id - extId of an AOS (PE) cluster registered to that PC
    Any value already provided via AZ_PC_EXT_ID / AZ_CLUSTER_EXT_ID is used as a
    fallback whenever discovery is unavailable or fails.
    """
    fallback = {
        "pc_ext_id": env.get("AZ_PC_EXT_ID"),
        "cluster_ext_id": env.get("AZ_CLUSTER_EXT_ID"),
        "remote_pc_ip": env.get("AZ_REMOTE_PC_IP"),
    }
    remote_ip = env.get("AZ_REMOTE_PC_IP")
    if not remote_ip:
        logger.warning("Cannot discover availability zone: AZ_REMOTE_PC_IP is empty; using .env values")
        return fallback
    username = env.get("AZ_REMOTE_PC_USERNAME") or env.get("PC_USERNAME")
    password = env.get("AZ_REMOTE_PC_PASSWORD") or env.get("PC_PASSWORD")
    port = env.get("PC_PORT")
    if not password or not port:
        logger.warning("Cannot discover availability zone: PC_USERNAME/PC_PORT not set; using .env values")
        return fallback

    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    token = base64.b64encode(f"{username}:{password}".encode()).decode()

    out = dict(fallback)
    log(f"Discovering availability zone from remote PC {remote_ip}")

    # PRISM_CENTRAL cluster entity extId -> pc_ext_id.
    try:
        pc_ext_id = _discover_cluster_ext_id(remote_ip, port, token, ctx, "PRISM_CENTRAL")
        if pc_ext_id:
            out["pc_ext_id"] = pc_ext_id
        else:
            logger.warning("No PRISM_CENTRAL cluster on %s; pc_ext_id left as .env value", remote_ip)
    except (urllib.error.URLError, ValueError) as exc:
        logger.warning("pc_ext_id discovery failed (%s); using .env value", getattr(exc, "reason", exc))

    # AOS (PE) cluster extId -> cluster_ext_id.
    try:
        cluster_ext_id = _discover_cluster_ext_id(remote_ip, port, token, ctx, "AOS")
        if cluster_ext_id:
            out["cluster_ext_id"] = cluster_ext_id
        else:
            logger.warning("No AOS cluster on %s; cluster_ext_id left as .env value", remote_ip)
    except (urllib.error.URLError, ValueError) as exc:
        logger.warning("cluster_ext_id discovery failed (%s); using .env value", getattr(exc, "reason", exc))

    log(f"Availability zone: pc_ext_id={out['pc_ext_id'] or '<unset>'} "
        f"cluster_ext_id={out['cluster_ext_id'] or '<unset>'} remote_pc_ip={remote_ip}")
    return out


def discover_unregister_pc_ext_id(env: Env) -> str:
    """Return the extId of the *remote* PC to unregister (prism.unregister).

    Mirrors preEnv/prism.tf, which connects to prism_unregister_pc_ip and reads
    that PC's PRISM_CENTRAL cluster entity extId. Only PRISM_UNREGISTER_REMOTE_PC_IP
    is required; PC credentials are reused (override with
    PRISM_UNREGISTER_REMOTE_PC_USERNAME / _PASSWORD when the remote PC differs).
    Falls back to PRISM_UNREGISTER_PC_EXT_ID when discovery is unavailable/fails.
    """
    fallback = env.get("PRISM_UNREGISTER_PC_EXT_ID")
    remote_ip = env.get("PRISM_UNREGISTER_REMOTE_PC_IP")
    if not remote_ip:
        logger.warning("prism.unregister: remote_pc_ip empty; using configured pc_ext_id")
        return fallback
    username = env.get("PRISM_UNREGISTER_REMOTE_PC_USERNAME") or env.get("PC_USERNAME")
    password = env.get("PRISM_UNREGISTER_REMOTE_PC_PASSWORD") or env.get("PC_PASSWORD")
    port = env.get("PC_PORT")
    if not (username and password and port):
        logger.warning("prism.unregister: PC_USERNAME/PASSWORD/PORT not set; using configured pc_ext_id")
        return fallback

    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    token = base64.b64encode(f"{username}:{password}".encode()).decode()

    # Ensure the local PC is availability-zone paired with the remote PC (connect
    # via cloud_trust when it isn't), so the unregistration test has a connected PC
    # to unregister -- mirrors prepare_pc.prepare_prism. Best-effort; a failure here
    # does not stop pc_ext_id discovery below.
    _ensure_az_connected(env, remote_ip, username, password, label="prism.unregister")

    log(f"Discovering prism.unregister pc_ext_id from remote PC {remote_ip}")
    try:
        ext_id = _discover_cluster_ext_id(remote_ip, port, token, ctx, "PRISM_CENTRAL")
    except (urllib.error.URLError, ValueError) as exc:
        logger.warning("prism.unregister pc_ext_id discovery failed (%s); using configured value",
                       getattr(exc, "reason", exc))
        return fallback
    if ext_id:
        log(f"prism.unregister: discovered pc_ext_id={ext_id}")
        return ext_id
    logger.warning("No PRISM_CENTRAL cluster on %s; using configured pc_ext_id", remote_ip)
    return fallback


def _discover_aos_cluster_name(host: str, port: str, token: str, ctx,
                               prefer_auto: bool = True) -> str:
    """Return the name of an AOS cluster on *host*, mirroring the cluster_fetch
    module's nutanix_clusters_v2 query. When prefer_auto is set the first cluster
    whose name starts with 'auto' is preferred; if none match (or prefer_auto is
    False) it falls back to the first AOS cluster."""
    base = "config/clusterFunction/any(t:t eq Clustermgmt.Config.ClusterFunctionRef'AOS')"
    filters = ([base + " and startswith(name, 'auto')"] if prefer_auto else []) + [base]
    for fltr in filters:
        path = ("/api/clustermgmt/v4.0/config/clusters?$limit=1&$filter="
                + urllib.parse.quote(fltr, safe=""))
        data = _pc_get_json(host, port, path, token, ctx)
        items = data.get("data") or []
        if items:
            return items[0].get("name") or ""
    return ""


def _jarvis_parse_custom(obj: dict) -> dict:
    """Return the ``custom`` map of a Jarvis node/cluster document. Jarvis stores
    ``custom`` as a JSON *string* (e.g. '{"v1_ip": "...", "v2_ip": "..."}'), so
    decode it when needed; return {} when it is missing or malformed."""
    custom = (obj or {}).get("custom")
    if isinstance(custom, str) and custom.strip():
        try:
            custom = json.loads(custom)
        except ValueError:
            return {}
    return custom if isinstance(custom, dict) else {}


def _jarvis_cluster_doc(cluster_name: str) -> dict:
    """GET the Jarvis cluster inventory document for *cluster_name* (GET
    https://jarvis.eng.nutanix.com/api/v1/clusters/<name>, no auth) and return the
    inner payload, unwrapping the top-level ``data`` envelope. Jarvis holds the
    reserved IPs / network layout even while the cluster is being (re)built, which
    is why we prefer it over the live PrismGateway API. Network/parse errors
    propagate to the caller."""
    jctx = ssl.create_default_context()
    jctx.check_hostname = False
    jctx.verify_mode = ssl.CERT_NONE
    url = ("https://jarvis.eng.nutanix.com/api/v1/clusters/"
           + urllib.parse.quote(cluster_name, safe=""))
    req = urllib.request.Request(url, method="GET")
    req.add_header("Accept", "application/json")
    logger.debug("HTTP GET %s", url)
    with urllib.request.urlopen(req, context=jctx, timeout=30) as resp:
        raw = resp.read().decode()
    logger.debug("jarvis cluster response: %s", _safe_json(raw))
    data = json.loads(raw) if raw.strip() else {}
    return data.get("data") if isinstance(data.get("data"), dict) else data


def _jarvis_nodes(payload: dict) -> list:
    """Return the node list of a Jarvis cluster document (top-level ``nodes`` or
    ``additional_data.nodes``)."""
    return (payload.get("nodes")
            or (payload.get("additional_data") or {}).get("nodes") or [])


def _jarvis_cluster_ips(cluster_name: str) -> dict:
    """Return {"svm_ip", "virtual_ip", "iscsi_ip"} for *cluster_name* from the
    Jarvis cluster document:
      svm_ip     = nodes[0].svm_ip
      virtual_ip = custom.v1_ip  (the cluster external / VIP)
      iscsi_ip   = custom.v2_ip  (the external data-services / iSCSI IP)
    where ``custom`` is a JSON string at the cluster level (falling back to the
    first node's ``custom``). Network/parse errors propagate to the caller.
    """
    payload = _jarvis_cluster_doc(cluster_name)
    nodes = _jarvis_nodes(payload)
    svm_ip = (nodes[0] or {}).get("svm_ip") if nodes else ""
    custom = _jarvis_parse_custom(payload)
    if not custom and nodes:
        custom = _jarvis_parse_custom(nodes[0])
    return {
        "svm_ip": svm_ip or "",
        "virtual_ip": _coerce_ip(custom.get("v1_ip")),
        "iscsi_ip": _coerce_ip(custom.get("v2_ip")),
    }


def _jarvis_cluster_network(cluster_name: str) -> dict:
    """Return the deploy-PC network fields for *cluster_name* from the Jarvis
    cluster document:
      default_gateway = nodes[0].network.default_gw
      subnet_mask     = nodes[0].network.svm_subnet_mask (else host_subnet_mask)
      network         = nodes[0].network.network  (the subnet/network address)
      ip              = custom.v1_ip  (the reserved external IP; used as ip_range)
    Missing fields come back as empty strings. Network/parse errors propagate to
    the caller.
    """
    payload = _jarvis_cluster_doc(cluster_name)
    nodes = _jarvis_nodes(payload)
    net = (nodes[0].get("network") if nodes and isinstance(nodes[0], dict) else {}) or {}
    custom = _jarvis_parse_custom(payload)
    if not custom and nodes:
        custom = _jarvis_parse_custom(nodes[0])
    return {
        "default_gateway": (net.get("default_gw") or "").strip(),
        "subnet_mask": (net.get("svm_subnet_mask") or net.get("host_subnet_mask") or "").strip(),
        "network": (net.get("network") or "").strip(),
        "ip": _coerce_ip(custom.get("v1_ip")),
    }


def _ssh_cluster_name(host: str, user: str, password: str, timeout: int = 45) -> str:
    """Return the AOS cluster name read over SSH from a PE CVM at *host*.

    Runs ``~/prism/cli/ncli cluster info`` ("Cluster Name : ...") on the CVM and,
    as a backstop, ``zeus_config_printer`` (Zeus config, field ``cluster_name``).
    ncli is invoked by its full path (it is not on the login PATH) and without
    ``-h true`` (which suppressed the output). Uses ``sshpass`` for the 'nutanix'
    account (no key auth); the connection tries the password once
    (NumberOfPasswordPrompts=1) so a locked/wrong credential fails fast. Returns
    "" when sshpass is missing, the host is unreachable, or the name can't be
    parsed.
    """
    if not (host and user and password):
        return ""
    if shutil.which("sshpass") is None:
        logger.warning("Cannot read cluster name over SSH: sshpass not found in PATH "
                       "(install it or pin the Jarvis name in config)")
        return ""
    remote = ("~/prism/cli/ncli cluster info 2>/dev/null | grep -i 'Cluster Name'; "
              "zeus_config_printer 2>/dev/null | grep -m1 cluster_name")
    argv = [
        "sshpass", "-p", password, "ssh",
        "-o", "StrictHostKeyChecking=no",
        "-o", "UserKnownHostsFile=/dev/null",
        "-o", "PubkeyAuthentication=no",
        "-o", "PreferredAuthentications=password",
        "-o", "NumberOfPasswordPrompts=1",
        "-o", "ConnectTimeout=15",
        f"{user}@{host}", remote,
    ]
    logger.debug("SSH %s@%s: reading cluster_name via `%s`", user, host, remote)
    try:
        proc = subprocess.run(argv, capture_output=True, text=True, timeout=timeout)
    except subprocess.TimeoutExpired:
        logger.warning("cluster name SSH to %s timed out after %ss "
                       "(host unreachable or slow)", host, timeout)
        return ""
    except OSError as exc:
        logger.warning("cluster name SSH to %s failed to launch (%s)", host, exc)
        return ""
    out = (proc.stdout or "").strip()
    err = (proc.stderr or "").strip()
    logger.debug("cluster name SSH to %s: rc=%s\n--- stdout ---\n%s\n--- stderr ---\n%s",
                 host, proc.returncode, out or "<empty>", err or "<empty>")
    m = re.search(r'cluster_name\s*:\s*"?([A-Za-z0-9._-]+)', out)
    if not m:
        m = re.search(r'Cluster Name\s*:\s*([A-Za-z0-9._-]+)', out)
    name = m.group(1) if m else ""
    if name:
        logger.debug("cluster name on %s = %s", host, name)
        return name
    # No name parsed. Explain the most likely cause from what we got back.
    if proc.returncode == 255 or "permission denied" in err.lower():
        logger.warning("cluster name SSH to %s: auth/connection failed (rc=%s): %s",
                       host, proc.returncode, err[:200] or "<no stderr>")
    elif "nutanix controller vm" in out.lower() or (not out and proc.returncode != 0):
        # sshd is up (we got the CVM login banner) but zeus_config_printer/ncli
        # produced nothing -> the AOS cluster services are down/not configured.
        logger.warning("cluster name SSH to %s: connected but cluster services returned "
                       "no cluster_name (rc=%s) -- AOS services are likely down; pin a "
                       "*_jarvis_name in config to resolve IPs while the cluster is down",
                       host, proc.returncode)
    else:
        logger.warning("cluster name SSH to %s: could not parse cluster_name from output "
                       "(rc=%s); run with -v to see the raw stdout/stderr", host, proc.returncode)
    return name


# ROBO / small-cluster hardware flags. Added to each PE CVM's
# /etc/nutanix/hardware_config.json ("hardware_attributes" section) so single /
# two-node and mixed-hypervisor clusters are accepted. Standalone twin of this
# logic lives in temp/scripts/enable_robo_flags.sh.
ROBO_HARDWARE_FLAGS = {
    "one_node_cluster": True,
    "two_node_cluster": True,
    "robo_mixed_hypervisor": True,
}

# Full path -- genesis is not on PATH for non-interactive ssh sessions.
GENESIS_BIN = "/usr/local/nutanix/cluster/bin/genesis"

# Python run on the CVM (under sudo) to inject the flags idempotently, with a
# timestamped backup and re-validation before the atomic swap.
_ROBO_FLAGS_PY = r'''
import json, os, sys, time, shutil

path = "/etc/nutanix/hardware_config.json"
flags = %s

with open(path) as f:
    data = json.load(f)

# genesis/foundation read these ROBO flags from the per-node section
# (node.hardware_attributes), so append them into that existing block -- do NOT
# create a separate top-level hardware_attributes object.
node = data.get("node")
if not isinstance(node, dict):
    print("ERROR: no 'node' object in hardware_config.json; cannot apply flags")
    sys.exit(1)
attrs = node.get("hardware_attributes")
if not isinstance(attrs, dict):
    attrs = {}
    node["hardware_attributes"] = attrs

changed = False
for k, v in flags.items():
    if attrs.get(k) != v:
        attrs[k] = v
        changed = True

# Remove the misplaced top-level hardware_attributes (genesis/foundation ignore
# it; earlier runs of this workaround wrongly wrote the flags here).
if "hardware_attributes" in data:
    del data["hardware_attributes"]
    changed = True
    print("removed misplaced top-level hardware_attributes")

if not changed:
    print("flags already present; no change needed")
    sys.exit(0)

backup = "%%s.bak.%%s" %% (path, time.strftime("%%Y%%m%%d-%%H%%M%%S"))
shutil.copy2(path, backup)
print("backup written to %%s" %% backup)

tmp = path + ".tmp"
with open(tmp, "w") as f:
    json.dump(data, f, indent=2)
    f.write("\n")
with open(tmp) as f:
    json.load(f)
os.rename(tmp, path)
print("flags added to hardware_attributes")
''' % repr(ROBO_HARDWARE_FLAGS)


def _ssh_apply_robo_flags(host: str, user: str, password: str,
                          restart: bool = True, timeout: int = 300) -> bool:
    """SSH into the PE CVM at *host* and add the ROBO small-cluster flags to
    /etc/nutanix/hardware_config.json, then (when *restart*) stop foundation and
    restart genesis. The edit is idempotent (skips when the flags are already
    present) and keeps a timestamped backup. Returns True on success.
    """
    if not (host and user and password):
        return False
    if shutil.which("sshpass") is None:
        logger.warning("robo flags: sshpass not found in PATH; cannot edit %s "
                       "(install it or apply the flags manually)", host)
        return False
    b64 = base64.b64encode(_ROBO_FLAGS_PY.encode()).decode()
    remote = ("echo %s | base64 -d | sudo bash -c '"
              "command -v python3 >/dev/null 2>&1 && PY=python3 || PY=python; $PY -'" % b64)
    if restart:
        remote += " && %s stop foundation && %s restart" % (GENESIS_BIN, GENESIS_BIN)
    argv = [
        "sshpass", "-p", password, "ssh",
        "-o", "StrictHostKeyChecking=no",
        "-o", "UserKnownHostsFile=/dev/null",
        "-o", "PubkeyAuthentication=no",
        "-o", "PreferredAuthentications=password",
        "-o", "NumberOfPasswordPrompts=1",
        "-o", "ConnectTimeout=15",
        f"{user}@{host}", remote,
    ]
    logger.info("robo flags: applying to %s ...", host)
    try:
        proc = subprocess.run(argv, capture_output=True, text=True, timeout=timeout)
    except subprocess.TimeoutExpired:
        logger.warning("robo flags: SSH to %s timed out after %ss "
                       "(host unreachable or genesis restart slow)", host, timeout)
        return False
    except OSError as exc:
        logger.warning("robo flags: SSH to %s failed to launch (%s)", host, exc)
        return False
    out = (proc.stdout or "").strip()
    err = (proc.stderr or "").strip()
    logger.debug("robo flags SSH to %s: rc=%s\n--- stdout ---\n%s\n--- stderr ---\n%s",
                 host, proc.returncode, out or "<empty>", err or "<empty>")
    if proc.returncode != 0:
        logger.warning("robo flags: %s FAILED (rc=%s): %s", host, proc.returncode,
                       (err or out)[:300] or "<no output>")
        return False
    for line in out.splitlines():
        low = line.lower()
        if "flag" in low or "backup" in low or "genesis started" in low:
            logger.info("robo flags: %s: %s", host, line.strip())
    logger.info("robo flags: %s done", host)
    return True


def apply_robo_flags(env: "Env", hosts, dry_run: bool = False) -> None:
    """Apply the ROBO hardware flags (+ genesis restart) to each PE node in
    *hosts* in parallel, using the ssh.pe_* credentials. De-dups hosts and skips
    silently when there is nothing to do. Failures are logged, not fatal.
    """
    hosts = list(dict.fromkeys(h for h in hosts if h))  # de-dup, preserve order
    if not hosts:
        logger.info("robo flags: no target PE nodes resolved; skipping")
        return
    user = env.get("SSH_PE_USERNAME")
    password = env.get("SSH_PE_PASSWORD")
    if not (user and password):
        logger.warning("robo flags: ssh.pe_username/pe_password not set; skipping %s",
                       ", ".join(hosts))
        return
    if dry_run:
        logger.info("robo flags: [dry-run] would edit hardware_config.json + restart "
                    "genesis on: %s", ", ".join(hosts))
        return
    logger.info("robo flags: applying to %d PE node(s) in parallel: %s",
                len(hosts), ", ".join(hosts))
    results = {}
    with concurrent.futures.ThreadPoolExecutor(max_workers=len(hosts)) as ex:
        futs = {ex.submit(_ssh_apply_robo_flags, h, user, password): h for h in hosts}
        for fut in concurrent.futures.as_completed(futs):
            h = futs[fut]
            try:
                results[h] = fut.result()
            except Exception as exc:  # noqa: BLE001
                logger.warning("robo flags: %s raised %s", h, exc)
                results[h] = False
    failed = [h for h, ok in results.items() if not ok]
    if failed:
        logger.warning("robo flags: completed with failures on: %s", ", ".join(failed))
    else:
        logger.info("robo flags: all %d node(s) updated successfully", len(hosts))


def _resolve_cluster_name(host: str, port: str, token: str, ctx, label: str,
                          explicit_name: str = "", prefer_auto: bool = True,
                          ssh_user: str = "", ssh_pass: str = "") -> tuple:
    """Resolve the AOS cluster name for the cluster at *host*, returning
    (name, source). The name is resolved in order:
      1. *explicit_name* (a pinned ``*_jarvis_name`` from config),
      2. SSH to *host* (zeus_config_printer / ncli) when ssh_user/ssh_pass given --
         works even when the cluster's HTTP APIs are down but sshd is up,
      3. *host*'s clusters v4 API (needs the API to be reachable).
    Returns ("", "") when nothing resolves. The name-resolution steps log their
    own failures and are swallowed (a down clusters v4 API is not fatal)."""
    name, source = "", ""
    if explicit_name:
        name, source = explicit_name, "pinned config *_jarvis_name"
        logger.debug("%s: using pinned Jarvis cluster name %r", label, name)
    if not name and ssh_user and ssh_pass:
        logger.debug("%s: resolving cluster name over SSH from %s", label, host)
        name = _ssh_cluster_name(host, ssh_user, ssh_pass)
        if name:
            source = "SSH"
    elif not name:
        logger.debug("%s: skipping SSH name lookup (ssh.pe_username/pe_password not set)", label)
    if not name:
        logger.debug("%s: resolving cluster name via %s clusters v4 API", label, host)
        try:
            name = _discover_aos_cluster_name(host, port, token, ctx, prefer_auto)
            if name:
                source = "clusters v4 API"
        except (urllib.error.URLError, ValueError) as exc:
            logger.warning("%s: clusters v4 name lookup on %s failed (%s)",
                           label, host, getattr(exc, "reason", exc))
    if not name:
        logger.warning("%s: could not resolve a cluster name for %s (SSH and clusters v4 "
                       "both unavailable); pin a *_jarvis_name in config to resolve while "
                       "the cluster is down", label, host)
    return name, source


def _resolve_cluster_ips(host: str, port: str, token: str, ctx, label: str,
                         explicit_name: str = "", prefer_auto: bool = True,
                         ssh_user: str = "", ssh_pass: str = "") -> dict:
    """Resolve {"svm_ip", "virtual_ip", "iscsi_ip", "name"} for the AOS cluster at
    *host* using the Jarvis inventory (which holds the reserved IPs even while the
    cluster is down). The cluster name is resolved via _resolve_cluster_name and
    the node IPs are then looked up through Jarvis. Returns empty strings when no
    name can be resolved; only Jarvis errors propagate to the caller."""
    empty = {"svm_ip": "", "virtual_ip": "", "iscsi_ip": "", "name": ""}
    name, source = _resolve_cluster_name(host, port, token, ctx, label,
                                         explicit_name=explicit_name,
                                         prefer_auto=prefer_auto,
                                         ssh_user=ssh_user, ssh_pass=ssh_pass)
    if not name:
        return empty

    logger.info("%s: resolved cluster name %r via %s; querying Jarvis", label, name, source)
    ips = _jarvis_cluster_ips(name)
    ips["name"] = name
    logger.debug("%s: Jarvis returned svm_ip=%s virtual_ip=%s iscsi_ip=%s for %r",
                 label, ips["svm_ip"] or "<unset>", ips["virtual_ip"] or "<unset>",
                 ips["iscsi_ip"] or "<unset>", name)
    return ips


def _discover_pe_svm_from_pc(env: Env, endpoint: str, username: str, password: str,
                             label: str, jarvis_name: str = "") -> str:
    """Return the PE svm_ip of the AOS cluster registered to the PC at *endpoint*
    (via that PC's clusters v4 API -> Jarvis). Returns "" on missing creds or
    discovery failure."""
    port = env.get("PC_PORT")
    if not (endpoint and username and password and port):
        logger.warning("%s: PC endpoint/creds/port not set; PE left empty", label)
        return ""
    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    token = base64.b64encode(f"{username}:{password}".encode()).decode()
    try:
        info = _resolve_cluster_ips(endpoint, port, token, ctx, label,
                                    explicit_name=jarvis_name, prefer_auto=True)
        return info["svm_ip"]
    except (urllib.error.URLError, ValueError) as exc:
        logger.warning("%s PE discovery failed (%s); PE left empty",
                       label, getattr(exc, "reason", exc))
        return ""


def discover_local_cluster(env: Env) -> dict:
    """Resolve data_protection.local_cluster_pe / _vip (also reused for
    prism.restore_source.pe_ip).

    The PE svm_ip is discovered from the AOS cluster registered to the *local* PC
    (PC_ENDPOINT; DP_LOCAL_CLUSTER_PE overrides it). The VIP is then chosen as the
    first *free* IP after that PE within its /24 (ping + ARP; see _find_free_ips),
    unless DP_LOCAL_CLUSTER_VIP is set explicitly.
    """
    pe = env.get("DP_LOCAL_CLUSTER_PE")
    vip = env.get("DP_LOCAL_CLUSTER_VIP")
    if pe and vip:
        log(f"data_protection.local_cluster: using configured pe={pe} vip={vip} (skipping discovery)")
        return {"pe": pe, "vip": vip}

    if not pe:
        pe = _discover_pe_svm_from_pc(
            env, env.get("PC_ENDPOINT"), env.get("PC_USERNAME"), env.get("PC_PASSWORD"),
            "data_protection.local_cluster", env.get("DP_LOCAL_CLUSTER_JARVIS_NAME"))

    # VIP = first free IP after the local PE within its /24.
    if pe and not vip:
        log(f"data_protection.local_cluster: finding 1 free IP after PE {pe} (local_cluster_vip)")
        free = _find_free_ips(pe, 1, {pe}, "data_protection.local_cluster")
        vip = free[0] if free else ""

    log(f"data_protection.local_cluster: pe={pe or '<unset>'} vip={vip or '<unset>'}")
    return {"pe": pe, "vip": vip}


def discover_remote_cluster(env: Env) -> dict:
    """Resolve data_protection.remote_cluster_pe / _vip for the test config.

    remote_cluster_pe is an explicit input (DP_REMOTE_CLUSTER_PE) -- the remote
    cluster the DP acceptance tests target. The VIP is the first *free* IP after
    it within its /24 (ping + ARP; see _find_free_ips), unless
    DP_REMOTE_CLUSTER_VIP is set explicitly.

    NOTE: this is distinct from the AZ remote PC used by the prepare_pc firewall
    step -- see discover_remote_pc_cluster.
    """
    pe = env.get("DP_REMOTE_CLUSTER_PE")
    vip = env.get("DP_REMOTE_CLUSTER_VIP")
    if not pe:
        logger.warning("data_protection.remote_cluster: DP_REMOTE_CLUSTER_PE empty")
        return {"pe": pe, "vip": vip}
    if not vip:
        log(f"data_protection.remote_cluster: finding 1 free IP after PE {pe} (remote_cluster_vip)")
        free = _find_free_ips(pe, 1, {pe}, "data_protection.remote_cluster")
        vip = free[0] if free else ""
    log(f"data_protection.remote_cluster: pe={pe or '<unset>'} vip={vip or '<unset>'}")
    return {"pe": pe, "vip": vip}


def discover_remote_pc_cluster(env: Env) -> dict:
    """Resolve the PE svm_ip + VIP of the AOS cluster behind the *remote PC*
    connected via the availability zone (AZ_REMOTE_PC_IP).

    Used only by the prepare_pc data-protection firewall step (the remote peer of
    the modify_firewall pair) -- NOT the test config's data_protection.remote_*
    (see discover_remote_cluster). The PE svm_ip is discovered from the remote
    PC's AOS cluster (remote PC creds default to AZ_REMOTE_PC_USERNAME / _PASSWORD,
    then the local PC creds); the VIP is the first *free* IP after it. Explicit
    DP_REMOTE_PC_CLUSTER_PE / _VIP override discovery.
    """
    pe = env.get("DP_REMOTE_PC_CLUSTER_PE")
    vip = env.get("DP_REMOTE_PC_CLUSTER_VIP")
    if pe and vip:
        log(f"data_protection.remote_pc_cluster: using configured pe={pe} vip={vip}")
        return {"pe": pe, "vip": vip}

    remote_pc = env.get("AZ_REMOTE_PC_IP")
    if not pe:
        if not remote_pc:
            logger.warning("data_protection.remote_pc_cluster: az.remote_pc_ip not set; PE left empty")
        else:
            pe = _discover_pe_svm_from_pc(
                env, remote_pc,
                env.get("AZ_REMOTE_PC_USERNAME") or env.get("PC_USERNAME"),
                env.get("AZ_REMOTE_PC_PASSWORD") or env.get("PC_PASSWORD"),
                "data_protection.remote_pc_cluster", env.get("DP_REMOTE_PC_CLUSTER_JARVIS_NAME"))

    if pe and not vip:
        log(f"data_protection.remote_pc_cluster: finding 1 free IP after PE {pe} (vip)")
        free = _find_free_ips(pe, 1, {pe}, "data_protection.remote_pc_cluster")
        vip = free[0] if free else ""

    log(f"data_protection.remote_pc_cluster: pe={pe or '<unset>'} vip={vip or '<unset>'} "
        f"(remote PC {remote_pc or '<unset>'})")
    return {"pe": pe, "vip": vip}


def discover_remote_cluster_vip(env: Env) -> str:
    """Thin wrapper returning just data_protection.remote_cluster_vip."""
    return discover_remote_cluster(env).get("vip") or ""


def discover_pc_domain(env: Env) -> str:
    """Return the PC's MSP/DNS domain used as the Objects store ``domain``.

    Fetches GET /api/nutanix/v3/prism_central on the local PC and reads
    ``resources.cmsp_config.pc_domain_name`` (e.g. ``msp.pc-xxxx.nutanix.com``).
    Returns "" whenever discovery is unavailable/fails.
    """
    endpoint = env.get("PC_ENDPOINT")
    username = env.get("PC_USERNAME")
    password = env.get("PC_PASSWORD")
    port = env.get("PC_PORT")
    if not (endpoint and username and password and port):
        logger.warning("Cannot discover PC domain: PC_ENDPOINT/USERNAME/PASSWORD/PORT "
                       "not set; leaving domain empty")
        return ""

    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    token = base64.b64encode(f"{username}:{password}".encode()).decode()

    try:
        data = _pc_get_json(endpoint, port, "/api/nutanix/v3/prism_central", token, ctx)
    except (urllib.error.URLError, ValueError) as exc:
        logger.warning("PC domain discovery failed (%s); leaving domain empty",
                       getattr(exc, "reason", exc))
        return ""

    domain = (((data.get("resources") or {}).get("cmsp_config") or {})
              .get("pc_domain_name")) or ""
    if domain:
        log(f"Discovered PC domain: {domain}")
    else:
        logger.warning("pc_domain_name not present in /api/nutanix/v3/prism_central; "
                       "leaving domain empty")
    return domain


# --------------------------------------------------------------------------- #
# Legacy test_config.json ("v3") builder
#
# Fills the flat test_config.json (mirrors preEnv/create_json_v3.tf). Everything
# except the two static config.yaml inputs (subnet_name, ad_rule_target) is
# discovered from the local PC + the AZ-connected remote PC.
# --------------------------------------------------------------------------- #
def _discover_default_container(host: str, port: str, token: str, ctx) -> dict:
    """Return the first storage container whose name starts with
    'default-container-' as its raw v4 object (``{}`` if none / on error).
    Mirrors preEnv/create_json_v3.tf's data.nutanix_storage_containers_v2 filter
    (take [0]). The v4 clustermgmt object carries both ``name`` and
    ``containerExtId`` (the UUID)."""
    fltr = "startswith(name,'default-container-')"
    path = ("/api/clustermgmt/v4.0/config/storage-containers?$limit=1&$filter="
            + urllib.parse.quote(fltr, safe=""))
    try:
        data = _pc_get_json(host, port, path, token, ctx)
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.warning("default storage-container discovery failed (%s)",
                       getattr(exc, "reason", exc))
        return {}
    items = data.get("data") or []
    if not items:
        logger.warning("no storage container matching 'default-container-*' found")
        return {}
    return items[0] or {}


def _discover_default_container_name(host: str, port: str, token: str, ctx) -> str:
    """Return the name of the default storage container (see
    _discover_default_container)."""
    return _discover_default_container(host, port, token, ctx).get("name") or ""


def _discover_pc_account_uuid(host: str, port: str, token: str, ctx) -> str:
    """Return the local PC's Calm/self-service account uuid.

    Mirrors preEnv/scripts/get_pc_account_uuid.sh: POST v3 accounts/list, pick the
    entity named 'NTNX_LOCAL_AZ' and read its
    status.resources.data.cluster_account_reference_list[0].resources.data
    .pc_account_uuid. Returns "" on any failure."""
    try:
        status, body = _pc_v3_post(host, port, token, ctx,
                                   "api/nutanix/v3/accounts/list",
                                   {"length": 250, "filter": "state!=DELETED;state!=DRAFT"})
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.warning("account_uuid discovery failed (%s)", getattr(exc, "reason", exc))
        return ""
    if status not in (0, 200) or not isinstance(body, dict):
        logger.warning("account_uuid: accounts/list on %s -> HTTP %s", host, status)
        return ""
    for ent in body.get("entities") or []:
        if ((ent.get("metadata") or {}).get("name")) != "NTNX_LOCAL_AZ":
            continue
        data = (((ent.get("status") or {}).get("resources") or {}).get("data") or {})
        refs = data.get("cluster_account_reference_list") or []
        if refs:
            uuid = (((refs[0].get("resources") or {}).get("data") or {})
                    .get("pc_account_uuid")) or ""
            if uuid:
                return uuid
    logger.warning("account_uuid: no NTNX_LOCAL_AZ account with a pc_account_uuid found")
    return ""


def discover_ad_group_dn(client: IamClient, ds_ext_id: str, wanted: str) -> str:
    """Look an AD group up in the directory and return its distinguishedName
    (read-only; does not create anything in IAM). Returns "" when not found."""
    if not ds_ext_id or client.dry_run:
        return ""
    target = wanted.split("@")[0].lower()
    for entity in client.search_directory(ds_ext_id, target,
                                          IamClient._AD_SEARCHED_ATTRS,
                                          IamClient._AD_RETURNED_ATTRS):
        etype = (entity.get("entityType") or "").lower()
        if etype and etype != "group":
            continue
        candidates = {
            attr(entity, "sAMAccountName").lower(),
            attr(entity, "cn").lower(),
            _rdn_value(entity.get("name") or "").lower(),
            _rdn_value(attr(entity, "distinguishedName")).lower(),
        }
        if target in candidates:
            return (attr(entity, "distinguishedName")
                    or entity.get("distinguishedName", "")
                    or (entity.get("name") or ""))
    return ""


def _lookup_v3_user_group_by_dn(host: str, port: str, token: str, ctx, dn: str) -> str:
    """Return the uuid of an existing v3 user_group with the given
    distinguished_name (case-insensitive), or "" when none matches."""
    try:
        status, body = _pc_v3_post(host, port, token, ctx,
                                   "api/nutanix/v3/user_groups/list",
                                   {"kind": "user_group", "length": 250})
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.warning("user_groups/list on %s failed (%s)", host, getattr(exc, "reason", exc))
        return ""
    if status not in (0, 200) or not isinstance(body, dict):
        return ""
    for ent in body.get("entities") or []:
        res = (ent.get("status") or {}).get("resources") or {}
        ent_dn = ((res.get("directory_service_user_group") or {}).get("distinguished_name")) or ""
        if ent_dn.lower() == dn.lower():
            return ((ent.get("metadata") or {}).get("uuid")) or ""
    return ""


def _ensure_v3_user_group(host: str, port: str, token: str, ctx, dn: str) -> str:
    """Ensure a DIRECTORY_SERVICE user_group with *dn* exists in PC, creating it
    via the v3 user_groups API when absent. Returns its uuid ("" on failure).

    This is what makes the user-groups "duplicate entity" acceptance test
    deterministic: the group must already exist so the test's own create collides
    with a bad-request. Idempotent -- a create that comes back as an
    already-exists error falls through to a list lookup."""
    body = {
        "metadata": {"kind": "user_group"},
        "spec": {"resources": {"directory_service_user_group": {"distinguished_name": dn}}},
    }
    try:
        status, resp = _pc_v3_post(host, port, token, ctx,
                                   "api/nutanix/v3/user_groups", body)
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.warning("user_groups create on %s failed (%s)", host, getattr(exc, "reason", exc))
        return ""
    if status in (200, 201, 202) and isinstance(resp, dict):
        uuid = ((resp.get("metadata") or {}).get("uuid")) or ""
        log(f"test_config.json: pre-created duplicate user_group '{dn}' (uuid={uuid})")
        return uuid
    # Non-2xx (e.g. already exists) -> resolve the existing group's uuid.
    uuid = _lookup_v3_user_group_by_dn(host, port, token, ctx, dn)
    if uuid:
        log(f"test_config.json: duplicate user_group '{dn}' already exists (uuid={uuid})")
    else:
        logger.warning("test_config.json: could not create or find user_group '%s'", dn)
    return uuid


def build_v3_test_config(env: Env, existing: dict, dry_run: bool = False,
                         discover: bool = True) -> dict:
    """Build the fields for the legacy test_config.json (create_json_v3.tf).

    subnet_name / ad_rule_target / self_service are static config.yaml inputs.
    The rest is discovered from the local PC (PC_ENDPOINT) and the AZ remote PC
    (AZ_REMOTE_PC_IP): default_container_name, account_uuid, the users'
    directory_service_uuid, and protection_policy.{local,destination}_az. Fields
    that cannot be discovered are omitted so the existing file value is kept.
    """
    out: dict = {}

    # --- static inputs from config.yaml ---
    subnet_name = env.get("TEST_CONFIG_SUBNET_NAME") or env.get("VMM_SUBNET_NAME")
    if subnet_name:
        out["subnet_name"] = subnet_name
    ad_name = env.get("TEST_CONFIG_AD_RULE_TARGET_NAME")
    ad_values = env.get("TEST_CONFIG_AD_RULE_TARGET_VALUES")
    if ad_name or ad_values:
        out["ad_rule_target"] = {"name": ad_name, "values": ad_values}
    self_service = {
        "bp_name_with_snapshot_config": (
            env.get("TEST_CONFIG_SELF_SERVICE_BP_NAME_WITH_SNAPSHOT_CONFIG")
            or env.get("PC_PREP_SELF_SERVICE_BLUEPRINT2_NAME")
        ),
        "bp_name": (
            env.get("TEST_CONFIG_SELF_SERVICE_BP_NAME")
            or env.get("PC_PREP_SELF_SERVICE_BLUEPRINT1_NAME")
        ),
        "app_name_with_snapshot_config": (
            env.get("TEST_CONFIG_SELF_SERVICE_APP_NAME_WITH_SNAPSHOT_CONFIG")
            or env.get("PC_PREP_SELF_SERVICE_BP2_APP_NAME")
        ),
    }
    self_service = {k: v for k, v in self_service.items() if v}
    if self_service:
        out["self_service"] = self_service

    if not discover:
        log("test_config.json: filling static fields only (discovery disabled)")
        return out

    endpoint = env.get("PC_ENDPOINT")
    username = env.get("PC_USERNAME")
    password = env.get("PC_PASSWORD")
    port = env.get("PC_PORT")
    if not (endpoint and username and password and port):
        logger.warning("test_config.json: PC_ENDPOINT/USERNAME/PASSWORD/PORT not set; "
                       "skipping discovery")
        return out

    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    token = base64.b64encode(f"{username}:{password}".encode()).decode()

    # --- default_container_name (local PC storage container) ---
    container = _discover_default_container_name(endpoint, port, token, ctx)
    if container:
        log(f"test_config.json: default_container_name={container}")
        out["default_container_name"] = container

    # --- account_uuid (NTNX_LOCAL_AZ Calm account) ---
    account_uuid = _discover_pc_account_uuid(endpoint, port, token, ctx)
    if account_uuid:
        log(f"test_config.json: account_uuid={account_uuid}")
        out["account_uuid"] = account_uuid

    # --- protection_policy.local_az (local PC + its AOS cluster) ---
    protection_policy: dict = {}
    try:
        local_pc_ext_id = _discover_cluster_ext_id(endpoint, port, token, ctx, "PRISM_CENTRAL")
        local_cluster_ext_id = _discover_cluster_ext_id(endpoint, port, token, ctx, "AOS")
    except (urllib.error.URLError, ValueError, OSError) as exc:
        logger.warning("protection_policy.local_az discovery failed (%s)",
                       getattr(exc, "reason", exc))
        local_pc_ext_id = local_cluster_ext_id = ""
    if local_pc_ext_id or local_cluster_ext_id:
        protection_policy["local_az"] = {
            "uuid": local_pc_ext_id,
            "cluster_uuid": local_cluster_ext_id,
        }
        log(f"test_config.json: protection_policy.local_az uuid={local_pc_ext_id or '<unset>'} "
            f"cluster_uuid={local_cluster_ext_id or '<unset>'}")

    # --- protection_policy.destination_az (AZ-connected remote PC) ---
    remote_ip = env.get("AZ_REMOTE_PC_IP")
    if remote_ip:
        remote_user = env.get("AZ_REMOTE_PC_USERNAME") or username
        remote_pass = env.get("AZ_REMOTE_PC_PASSWORD") or password
        # Ensure the local PC is AZ-paired with the destination PC (cloud_trust);
        # a mutation, so skip it on dry runs. Non-fatal.
        if not dry_run:
            _ensure_az_connected(env, remote_ip, remote_user, remote_pass,
                                 label="protection_policy.destination_az")
        az = discover_availability_zone(env)
        dest_pc = az.get("pc_ext_id") or ""
        dest_cluster = az.get("cluster_ext_id") or ""
        if dest_pc or dest_cluster:
            protection_policy["destination_az"] = {
                "uuid": dest_pc,
                "cluster_uuid": dest_cluster,
            }
            log(f"test_config.json: protection_policy.destination_az uuid={dest_pc or '<unset>'} "
                f"cluster_uuid={dest_cluster or '<unset>'}")
    else:
        logger.warning("test_config.json: az.remote_pc_ip not set; "
                       "protection_policy.destination_az left as-is")
    if protection_policy:
        out["protection_policy"] = protection_policy

    # --- resolve the qa-nucalm-io directory service once (shared by users[] and
    #     the duplicate user-group) ---
    ds_name = (env.get("TEST_CONFIG_USERS_DIRECTORY_SERVICE_NAME")
               or _sanitize_ds_name(env.get("IAM_SECONDARY_AD_NAME")))
    client = None
    ds_ext_id = ""
    if ds_name and not dry_run:
        try:
            client = IamClient(env, dry_run=dry_run)
            ds = client.find_by("directory-services", "name", ds_name)
            ds_ext_id = (ds or {}).get("extId") if ds else ""
        except (RuntimeError, SystemExit) as exc:
            logger.warning("test_config.json: directory-service '%s' lookup failed (%s)",
                           ds_name, exc)
        if ds_ext_id:
            log(f"test_config.json: directory service '{ds_name}' extId={ds_ext_id}")
        else:
            logger.warning("test_config.json: directory service '%s' not found; users[] "
                           "and the duplicate user-group left as-is", ds_name)

    # --- users[] (principal_name + display_name + directory_service_uuid all
    #     fetched from the directory; nothing hardcoded) ---
    existing_users = existing.get("users") if isinstance(existing, dict) else None
    # The AD accounts to look up: TEST_CONFIG_USERS, falling back to the short
    # names of the users already in the file.
    wanted = env.list("TEST_CONFIG_USERS") or [
        (u.get("principal_name") or "").split("@")[0]
        for u in (existing_users or []) if isinstance(u, dict)
    ]
    wanted = [w for w in wanted if w]
    if client and ds_ext_id and wanted:
        discovered = []
        for name in wanted:
            try:
                principal, display = discover_ad_user_identity(client, ds_ext_id, name)
            except RuntimeError as exc:
                logger.warning("test_config.json: user '%s' lookup failed (%s)", name, exc)
                principal, display = "", ""
            if principal:
                log(f"test_config.json: user '{name}' -> principal={principal} "
                    f"display_name={display!r}")
                discovered.append({
                    "principal_name": principal,
                    "expected_display_name": display,
                    "directory_service_uuid": ds_ext_id,
                })
            else:
                logger.warning("test_config.json: user '%s' not found in directory '%s'",
                               name, ds_name)
        # Only replace the block when every configured user resolved, so a
        # partial/failed search never truncates the file's users[].
        if discovered and len(discovered) == len(wanted):
            out["users"] = discovered
        elif discovered:
            logger.warning("test_config.json: only %d/%d users resolved; users[] "
                           "left as-is", len(discovered), len(wanted))

    # --- user_group_with_distinguished_name: back the user-groups tests with
    #     groups from the CONNECTED directory (never a stale/disconnected one).
    #     DNs are stored LOWER-CASED because PC normalises directory DNs to lower
    #     case, so the resource/data-source reads must match.
    #       [0] duplicate-entity + data-source tests -> discovered AND pre-created
    #           in PC (v3) so the test's own create collides with a bad-request.
    #       [1] basic (directory_service_user_group) and [2] with-org-unit
    #           (directory_service_ou) -> DN discovered only; the tests create and
    #           destroy these groups themselves, so they are NOT pre-created here.
    existing_groups = (existing.get("user_group_with_distinguished_name")
                       if isinstance(existing, dict) else None) or []
    groups_out = [dict(g) if isinstance(g, dict) else g for g in existing_groups]

    def _set_group(idx: int, entry: dict) -> None:
        while len(groups_out) <= idx:
            groups_out.append({})
        groups_out[idx] = entry

    dup_name = env.get("TEST_CONFIG_DUPLICATE_USER_GROUP_NAME")
    if not dup_name and existing_groups and isinstance(existing_groups[0], dict):
        dup_name = _rdn_value(existing_groups[0].get("distinguished_name") or "")
    if client and ds_ext_id and dup_name:
        dn = discover_ad_group_dn(client, ds_ext_id, dup_name)
        if not dn and existing_groups and isinstance(existing_groups[0], dict):
            dn = existing_groups[0].get("distinguished_name") or ""
        if dn:
            uuid = _ensure_v3_user_group(endpoint, port, token, ctx, dn)
            if not uuid and existing_groups and isinstance(existing_groups[0], dict):
                uuid = existing_groups[0].get("uuid") or ""
            _set_group(0, {
                "distinguished_name": dn.lower(),
                "display_name": _rdn_value(dn),
                "uuid": uuid,
            })
            log(f"test_config.json: user_group_with_distinguished_name[0] "
                f"dn={dn.lower()} uuid={uuid}")
        else:
            logger.warning("test_config.json: duplicate user-group '%s' not found in "
                           "directory '%s'; user_group_with_distinguished_name[0] left as-is",
                           dup_name, ds_name)

    # [1] basic and [2] with-org-unit groups (TEST_CONFIG_USER_GROUP_NAMES in
    # order, mapped to indices 1, 2, ...) discovered from the same directory.
    extra_names = env.list("TEST_CONFIG_USER_GROUP_NAMES")
    if client and ds_ext_id and extra_names:
        for offset, gname in enumerate(extra_names):
            idx = offset + 1
            dn = discover_ad_group_dn(client, ds_ext_id, gname)
            if dn:
                _set_group(idx, {"distinguished_name": dn.lower()})
                log(f"test_config.json: user_group_with_distinguished_name[{idx}] "
                    f"dn={dn.lower()}")
            else:
                logger.warning("test_config.json: user-group '%s' not found in directory "
                               "'%s'; user_group_with_distinguished_name[%d] left as-is",
                               gname, ds_name, idx)

    if groups_out:
        out["user_group_with_distinguished_name"] = groups_out

    return out


# --------------------------------------------------------------------------- #
# Top-level (offline) builder
# --------------------------------------------------------------------------- #
def pad_nodes(nodes: list, size: int = 4) -> list:
    nodes = list(nodes)
    return (nodes + [""] * size)[:max(size, len(nodes))]


def build_ssl_certificate(env: Env, dry_run: bool) -> dict:
    if env.bool("CLUSTER_SSL_GENERATE"):
        return generate_ssl_certificate(env, dry_run)
    return {
        "passphrase": env.get("CLUSTER_SSL_PASSPHRASE"),
        "private_key": env.file_or_inline("CLUSTER_SSL_PRIVATE_KEY_FILE", "CLUSTER_SSL_PRIVATE_KEY"),
        "public_certificate": env.file_or_inline("CLUSTER_SSL_PUBLIC_CERTIFICATE_FILE", "CLUSTER_SSL_PUBLIC_CERTIFICATE"),
        "ca_chain": env.file_or_inline("CLUSTER_SSL_CA_CHAIN_FILE", "CLUSTER_SSL_CA_CHAIN"),
    }


def _to_int(value, default: int = 0) -> int:
    """Coerce an env string to int, returning default when empty/non-numeric."""
    try:
        return int(str(value).strip())
    except (TypeError, ValueError):
        return default


def build_networking(env: Env, discover_ips: bool = True) -> dict:
    """Build the networking block used by the networkingv2 subnet/route tests.

    subnets.* describe an external VLAN to create test subnets on, so they are
    environment inputs (no discovery). gc_subnet is a full subnet fixture (the
    vmmv2 GC deploy test reads its start_ip/end_ip to pick free VM IPs).
    external_nat_subnet / overlay_subnet are the *names* of subnets provisioned
    by testenv/terraform that the tests reference by name.

    kubernetes_cluster_ext_id is discovered from PC (when discover_ips is set)
    and only included when found, so an existing value is never clobbered.

    project_id is optional: the route tests now create their own project inline,
    so it is only kept for any config that still references a pre-existing one.
    """
    out = {
        "subnets": {
            "vlan_id": _to_int(env.get("NETWORKING_SUBNETS_VLAN_ID")),
            "network_ip": env.get("NETWORKING_SUBNETS_NETWORK_IP"),
            "network_prefix": _to_int(env.get("NETWORKING_SUBNETS_NETWORK_PREFIX")),
            "gateway_ip": env.get("NETWORKING_SUBNETS_GATEWAY_IP"),
            "dhcp": {
                "start_ip": env.get("NETWORKING_SUBNETS_DHCP_START_IP"),
                "end_ip": env.get("NETWORKING_SUBNETS_DHCP_END_IP"),
            },
        },
        "gc_subnet": {
            "name": env.get("NETWORKING_GC_SUBNET_NAME"),
            "vlan_id": _to_int(env.get("NETWORKING_GC_SUBNET_VLAN_ID")),
            "prefix_length": _to_int(env.get("NETWORKING_GC_SUBNET_PREFIX_LENGTH")),
            "gateway_ip": env.get("NETWORKING_GC_SUBNET_GATEWAY_IP"),
            "start_ip": env.get("NETWORKING_GC_SUBNET_START_IP"),
            "end_ip": env.get("NETWORKING_GC_SUBNET_END_IP"),
        },
        "external_nat_subnet": env.get("NETWORKING_EXTERNAL_NAT_SUBNET"),
        "overlay_subnet": env.get("NETWORKING_OVERLAY_SUBNET"),
    }
    # Only set kubernetes_cluster_ext_id when we resolve a non-empty value, so a
    # failed/offline discovery preserves whatever is already in test_config_v2.json.
    k8s_ext_id = discover_kubernetes_cluster_ext_id(env) if discover_ips else \
        env.get("NETWORKING_KUBERNETES_CLUSTER_EXT_ID")
    if k8s_ext_id:
        out["kubernetes_cluster_ext_id"] = k8s_ext_id
    return out


def discover_deploy_pc(env: Env) -> dict:
    """Resolve the prism.deploy_pc network fields from the PE that hosts the
    cluster (PRISM_DEPLOY_PC_PE_IP).

    The cluster name is resolved from that PE (PRISM_DEPLOY_PC_JARVIS_NAME override
    -> SSH via ssh.pe_* -> the PE's clusters v4 API), then the Jarvis inventory
    supplies:
      default_gateway    <- nodes[0].network.default_gw
      subnet_mask        <- nodes[0].network.svm_subnet_mask
      ip_range.begin/end <- custom.v1_ip  (the reserved external IP), but only when
                            that IP is actually free (ping + ARP); otherwise the
                            first free IP in the PE's /24 is picked instead so the
                            deploy does not fail with "UVM IP already in use".
    Explicit PRISM_DEPLOY_PC_* values win per field. ``version`` stays an explicit
    config input; the deploy network is no longer taken from config -- the deploy
    test creates the external subnet on the PE and uses its ext_id (see
    resource_nutanix_pc_deploy_v2_test.go, mirroring preEnv/prism.tf).
    Returns {"default_gateway", "subnet_mask", "ip_range": {"begin", "end"}}.
    """
    # pe_ip is the anchor for the whole deploy-PC block: with no PE to deploy from,
    # there is nothing to discover, so skip filling the network fields entirely and
    # emit blanks (the deploy test self-skips on empty pe_ip).
    pe_ip = env.get("PRISM_DEPLOY_PC_PE_IP")
    if not pe_ip:
        logger.warning("prism.deploy_pc: pe_ip is empty; skipping deploy-PC network "
                       "discovery and leaving default_gateway/subnet_mask/ip_range blank")
        return {"default_gateway": "", "subnet_mask": "",
                "ip_range": {"begin": "", "end": ""}}

    gw = env.get("PRISM_DEPLOY_PC_DEFAULT_GATEWAY")
    mask = env.get("PRISM_DEPLOY_PC_SUBNET_MASK")
    begin = env.get("PRISM_DEPLOY_PC_IP_RANGE_BEGIN")
    end = env.get("PRISM_DEPLOY_PC_IP_RANGE_END")

    def _result():
        return {"default_gateway": gw, "subnet_mask": mask,
                "ip_range": {"begin": begin, "end": end}}

    if gw and mask and begin and end:
        log(f"prism.deploy_pc: using configured gateway={gw} subnet_mask={mask} "
            f"ip_range={begin}-{end} (skipping discovery)")
        return _result()

    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    pe_user = env.get("CFG_PE_USERNAME") or env.get("PC_USERNAME")
    pe_pass = env.get("CFG_PE_PASSWORD") or env.get("PC_PASSWORD")
    token = base64.b64encode(f"{pe_user}:{pe_pass}".encode()).decode()
    port = env.get("PE_PORT") or env.get("PC_PORT") or "9440"

    name, source = _resolve_cluster_name(
        pe_ip, port, token, ctx, "prism.deploy_pc",
        explicit_name=env.get("PRISM_DEPLOY_PC_JARVIS_NAME"), prefer_auto=True,
        ssh_user=env.get("SSH_PE_USERNAME"), ssh_pass=env.get("SSH_PE_PASSWORD"))
    if not name:
        logger.warning("prism.deploy_pc: could not resolve a cluster name from PE %s; "
                       "leaving network fields as configured", pe_ip)
        return _result()

    logger.info("prism.deploy_pc: resolved cluster name %r via %s; querying Jarvis",
                name, source)
    try:
        net = _jarvis_cluster_network(name)
    except (urllib.error.URLError, ValueError) as exc:
        logger.warning("prism.deploy_pc: Jarvis lookup for %r failed (%s); leaving network "
                       "fields as configured", name, getattr(exc, "reason", exc))
        return _result()

    gw = gw or net["default_gateway"]
    mask = mask or net["subnet_mask"]

    # Jarvis returns a *reserved* external IP (custom.v1_ip) for the deploy target,
    # but nothing guarantees it is actually unused -- PC deploy then fails hard with
    # "UVM IP(s) is/are already in use". So treat the Jarvis IP as the preferred
    # candidate, verify it is free (ping + ARP), and otherwise fall back to the
    # first free IP in the same /24. The scan is anchored at pe_ip (which is
    # reachable, so _find_free_ips' reachability check is meaningful) rather than
    # the possibly-dead reserved IP. An explicit PRISM_DEPLOY_PC_IP_RANGE_* still
    # wins (begin is already set in that case, so we skip discovery here).
    if not begin:
        candidate = net["ip"]
        if candidate and _ip_is_free(candidate):
            begin = end = candidate
        else:
            if candidate:
                logger.warning("prism.deploy_pc: Jarvis reserved IP %s is in use; "
                               "scanning %s's /24 for a free deploy IP",
                               candidate, pe_ip)
            free = _find_free_ips(pe_ip, 1, {pe_ip, candidate}, "prism.deploy_pc")
            begin = end = (free[0] if free else candidate)
    end = end or begin
    log(f"prism.deploy_pc: cluster={name} gateway={gw or '<unset>'} "
        f"subnet_mask={mask or '<unset>'} ip_range={begin or '<unset>'}-{end or '<unset>'}")
    return _result()


def build_prism(env: Env, discover: bool = True, local_cluster_pe: str = "") -> dict:
    """Build the prism block used by the prismv2 tests (PC deploy / unregister /
    backup-target / restore-source / restore-PC).

    deploy_pc.default_gateway / subnet_mask / ip_range are discovered from the PE
    that hosts the cluster (PRISM_DEPLOY_PC_PE_IP) via the Jarvis inventory (see
    discover_deploy_pc); pe_ip / version stay explicit inputs. The deploy network
    is not in config -- the deploy test creates the external subnet on the PE and
    uses its ext_id, and reads name/NTP servers from the top-level dns_servers /
    ntp_servers (as
    preEnv/create_json.tf does). unregister.pc_ext_id is discovered from the
    remote unregister PC. restore_source.pe_ip is the local cluster PE (reuses the
    value discovered for data_protection.local_cluster_pe; pass it via
    local_cluster_pe). bucket.* carries the AWS S3 secrets (sourced from the
    gitignored config, never committed).
    """
    unregister_ext_id = (discover_unregister_pc_ext_id(env) if discover
                         else env.get("PRISM_UNREGISTER_PC_EXT_ID"))
    # restore_source.pe_ip is the local cluster PE (== data_protection.local_cluster_pe
    # in preEnv/create_json.tf): an explicit override wins, else reuse the value
    # already discovered for the local cluster.
    restore_pe_ip = env.get("PRISM_RESTORE_SOURCE_PE_IP") or (local_cluster_pe if discover else "")
    deploy = discover_deploy_pc(env) if discover else {
        "default_gateway": env.get("PRISM_DEPLOY_PC_DEFAULT_GATEWAY"),
        "subnet_mask": env.get("PRISM_DEPLOY_PC_SUBNET_MASK"),
        "ip_range": {
            "begin": env.get("PRISM_DEPLOY_PC_IP_RANGE_BEGIN"),
            "end": env.get("PRISM_DEPLOY_PC_IP_RANGE_END"),
        },
    }
    return {
        "deploy_pc": {
            "pe_ip": env.get("PRISM_DEPLOY_PC_PE_IP"),
            "version": env.get("PRISM_DEPLOY_PC_VERSION"),
            # Name of the external subnet pre-created on the deploy PE by
            # testenv/terraform (networking.tf); the deploy test looks it up by
            # this name instead of creating its own subnet.
            "subnet_name": env.get("PRISM_DEPLOY_PC_SUBNET_NAME"),
            "default_gateway": deploy.get("default_gateway"),
            "subnet_mask": deploy.get("subnet_mask"),
            "ip_range": deploy.get("ip_range"),
        },
        "unregister": {
            "pc_ext_id": unregister_ext_id,
        },
        "bucket": {
            "name": env.get("PRISM_BUCKET_NAME"),
            "region": env.get("PRISM_BUCKET_REGION"),
            "access_key": env.get("PRISM_BUCKET_ACCESS_KEY"),
            "secret_key": env.get("PRISM_BUCKET_SECRET_KEY"),
        },
        "restore_source": {
            "pe_ip": restore_pe_ip,
        },
        "pc_restore": {
            # Admin creds come from the top-level pc_username / pc_password.
            "skip_pc_restore_test": env.bool("PRISM_PC_RESTORE_SKIP_PC_RESTORE_TEST", True),
        },
    }


# Top-level test_config_v2.json sections that --only understands. "cfg" groups the
# shared login/ssh/dns/ntp scalars that are not their own nested block. "iam" is
# handled separately in main() (it is a different builder), listed here so --only
# validation accepts it.
TOP_LEVEL_SECTIONS = (
    "images", "cfg", "clusters", "data_protection", "lcm", "availability_zone",
    "networking", "vmm", "object_store", "prism", "security", "iam",
)


def build_top_level(env: Env, dry_run: bool = False, discover_ips: bool = True,
                    only: "set[str] | None" = None) -> dict:
    """Build the top-level test_config_v2.json blocks.

    When *only* is a non-empty set of section names (e.g. {"prism"}), just those
    sections are built -- and only their discovery is run -- so the rest of the
    document is left untouched by the caller's deep-merge. Recognised sections are
    in TOP_LEVEL_SECTIONS ("cfg" covers the shared login/ssh/dns/ntp scalars).
    """
    only = set(only or ())

    def want(section: str) -> bool:
        return not only or section in only

    # local_cluster feeds both data_protection and prism.restore_source.pe_ip;
    # resolve it at most once and only when a section that needs it is built.
    _local_cluster: "dict | None" = None

    def local_cluster() -> dict:
        nonlocal _local_cluster
        if _local_cluster is None:
            if discover_ips:
                _local_cluster = discover_local_cluster(env)
            else:
                _local_cluster = {
                    "pe": env.get("DP_LOCAL_CLUSTER_PE"),
                    "vip": env.get("DP_LOCAL_CLUSTER_VIP"),
                }
        return _local_cluster

    out: dict = {}

    if want("images"):
        out["images"] = {
            "ubuntu_image": env.get("IMAGES_UBUNTU_IMAGE"),
            "ubuntu_image_url": env.get("IMAGES_UBUNTU_IMAGE_URL"),
            "windows_image": env.get("IMAGES_WINDOWS_IMAGE"),
            "windows_image_url": env.get("IMAGES_WINDOWS_IMAGE_URL"),
            "centos_image": env.get("IMAGES_CENTOS_IMAGE"),
            "centos_image_url": env.get("IMAGES_CENTOS_IMAGE_URL"),
            "ngt_image": env.get("IMAGES_NGT_IMAGE"),
            "ngt_image_url": env.get("IMAGES_NGT_IMAGE_URL"),
            "iso_image_url": env.get("IMAGES_ISO_IMAGE_URL"),
            "iso_image_sha1": env.get("IMAGES_ISO_IMAGE_SHA1"),
            "iso_image_sha256": env.get("IMAGES_ISO_IMAGE_SHA256"),
        }

    if want("cfg"):
        out.update({
            "username_for_test": env.get("CFG_USERNAME_FOR_TEST"),
            "password_for_test": env.get("CFG_PASSWORD_FOR_TEST"),
            "pc_username": env.get("CFG_PC_USERNAME"),
            "pc_password": env.get("CFG_PC_PASSWORD"),
            "pe_username": env.get("CFG_PE_USERNAME"),
            "pe_password": env.get("CFG_PE_PASSWORD"),
            "ssh_pc_username": env.get("SSH_PC_USERNAME"),
            "ssh_pc_password": env.get("SSH_PC_PASSWORD"),
            "ssh_pe_username": env.get("SSH_PE_USERNAME"),
            "ssh_pe_password": env.get("SSH_PE_PASSWORD"),
            "dns_servers": env.list("DNS_SERVERS"),
            "ntp_servers": env.list("NTP_SERVERS"),
        })

    if want("clusters"):
        if discover_ips:
            network = discover_cluster_network(env)
        else:
            network = {
                "virtual_ip": env.get("CLUSTER_VIRTUAL_IP"),
                "iscsi_ip": env.get("CLUSTER_ISCSI_IP"),
            }
        out["clusters"] = {
            "nodes": pad_nodes(env.list("CLUSTER_NODES")),
            "network": network,
            "ssl_certificate": build_ssl_certificate(env, dry_run),
        }

    if want("data_protection"):
        lc = local_cluster()
        remote_cluster_vip = (discover_remote_cluster_vip(env) if discover_ips
                              else env.get("DP_REMOTE_CLUSTER_VIP"))
        out["data_protection"] = {
            "local_cluster_pe": lc["pe"],
            "local_cluster_vip": lc["vip"],
            "remote_cluster_pe": env.get("DP_REMOTE_CLUSTER_PE"),
            "remote_cluster_vip": remote_cluster_vip,
        }

    if want("lcm"):
        out["lcm"] = {
            "entity_model": env.get("LCM_ENTITY_MODEL"),
            "entity_model_version": env.get("LCM_ENTITY_MODEL_VERSION"),
        }

    if want("availability_zone"):
        if discover_ips:
            out["availability_zone"] = discover_availability_zone(env)
        else:
            out["availability_zone"] = {
                "pc_ext_id": env.get("AZ_PC_EXT_ID"),
                "cluster_ext_id": env.get("AZ_CLUSTER_EXT_ID"),
                "remote_pc_ip": env.get("AZ_REMOTE_PC_IP"),
            }

    if want("networking"):
        out["networking"] = build_networking(env, discover_ips)

    if want("vmm"):
        out["vmm"] = {
            "integration_vm": env.get("VMM_INTEGRATION_VM"),
            "subnet_name": env.get("VMM_SUBNET_NAME"),
            "assigned_ip": env.get("VMM_ASSIGNED_IP"),
            "unattend_xml": env.get("VMM_UNATTEND_XML"),
            "subnet": {
                "network_id": _to_int(env.get("VMM_SUBNET_NETWORK_ID")),
                "ip": env.get("VMM_SUBNET_IP"),
                "prefix": _to_int(env.get("VMM_SUBNET_PREFIX")),
                "gateway_ip": env.get("VMM_SUBNET_GATEWAY_IP"),
                "start_ip": env.get("VMM_SUBNET_START_IP"),
                "end_ip": env.get("VMM_SUBNET_END_IP"),
            },
            "ngt": {
                "ngt_upgrade_vm_name": env.get("VMM_NGT_NGT_UPGRADE_VM_NAME"),
                "credential": {
                    "username": env.get("VMM_NGT_CREDENTIAL_USERNAME"),
                    "password": env.get("VMM_NGT_CREDENTIAL_PASSWORD"),
                },
            },
            "ngt_vm": {
                "name": env.get("VMM_NGT_VM_NAME"),
                "username": env.get("VMM_NGT_VM_USERNAME"),
                "password": env.get("VMM_NGT_VM_PASSWORD"),
            },
            "gpus": [
                {
                    "device_id": _to_int(env.get("VMM_GPU_DEVICE_ID")),
                    "mode": env.get("VMM_GPU_MODE"),
                    "vendor": env.get("VMM_GPU_VENDOR"),
                }
            ],
            "ova_url": env.get("VMM_OVA_URL"),
            "gc_profile": {
                "vm_name": env.get("VMM_GC_PROFILE_VM_NAME"),
                "default_image_username": env.get("VMM_GC_PROFILE_DEFAULT_IMAGE_USERNAME"),
                "default_image_password": env.get("VMM_GC_PROFILE_DEFAULT_IMAGE_PASSWORD"),
            },
        }

    if want("object_store"):
        object_store_domain = discover_pc_domain(env) if discover_ips else ""
        out["object_store"] = {
            "subnet_name": env.get("OBJECT_STORE_SUBNET_NAME"),
            "bucket_name": env.get("OBJECT_STORE_BUCKET_NAME"),
            "domain": object_store_domain,
            "public_network_ips": env.list("OBJECT_STORE_PUBLIC_NETWORK_IPS"),
            "storage_network_dns_ip": env.list("OBJECT_STORE_STORAGE_NETWORK_DNS_IP"),
            "storage_network_vip": env.list("OBJECT_STORE_STORAGE_NETWORK_VIP"),
        }

    if want("prism"):
        out["prism"] = build_prism(env, discover_ips, local_cluster()["pe"])

    if want("security"):
        out["security"] = {
            "kms": {
                "endpoint_url": env.get("SECURITY_KMS_ENDPOINT_URL"),
                "key_id": env.get("SECURITY_KMS_KEY_ID"),
                "tenant_id": env.get("SECURITY_KMS_TENANT_ID"),
                "client_id": env.get("SECURITY_KMS_CLIENT_ID"),
                "client_secret": env.get("SECURITY_KMS_CLIENT_SECRET"),
            },
        }

    return out


# --------------------------------------------------------------------------- #
# Merge + write
# --------------------------------------------------------------------------- #
def deep_merge(base: dict, override: dict) -> dict:
    """Recursively merge override into base (override wins for scalars/lists)."""
    for key, val in override.items():
        if isinstance(val, dict) and isinstance(base.get(key), dict):
            deep_merge(base[key], val)
        else:
            base[key] = val
    return base


def _set_env_file_var(path: Path, key: str, value: str) -> bool:
    """Update ``KEY=...`` in a flat .env file in place, preserving any inline
    trailing comment and appending the line if the key is absent. Returns True
    when the file was written (i.e. the value actually changed)."""
    new_line_value = f'"{value}"'
    pattern = re.compile(rf"^(\s*(?:export\s+)?{re.escape(key)}\s*=)(.*)$")

    if path.exists():
        lines = path.read_text().splitlines()
    else:
        logger.warning("%s: not found; creating it to set %s", path, key)
        lines = []

    for i, line in enumerate(lines):
        m = pattern.match(line)
        if not m:
            continue
        rest = m.group(2)
        # Keep a trailing inline comment (e.g. `VALUE  # note`) if present.
        hash_idx = rest.find("#")
        comment = f"  {rest[hash_idx:].rstrip()}" if hash_idx != -1 else ""
        updated = f"{m.group(1)}{new_line_value}{comment}"
        if updated == line:
            log(f"{path}: {key} already up to date ({value})")
            return False
        lines[i] = updated
        path.write_text("\n".join(lines) + "\n")
        return True

    # Key not present -> append it.
    lines.append(f"{key}={new_line_value}")
    path.write_text("\n".join(lines) + "\n")
    return True


def update_env_storage_container(env: Env, env_file: Path, dry_run: bool) -> None:
    """Discover the local PC's default storage container and write its ext_id
    (UUID) into env_file's NUTANIX_STORAGE_CONTAINER -- the value the Go
    acceptance tests read via TestAccPreCheck. No-op when the PC connection is
    not configured or the container cannot be found."""
    endpoint = env.get("PC_ENDPOINT")
    username = env.get("PC_USERNAME")
    password = env.get("PC_PASSWORD")
    port = env.get("PC_PORT")
    if not (endpoint and username and password and port):
        logger.warning("NUTANIX_STORAGE_CONTAINER: PC_ENDPOINT/USERNAME/PASSWORD/PORT "
                       "not set; leaving %s untouched", env_file)
        return

    ctx = ssl.create_default_context()
    if env.bool("PC_INSECURE"):
        ctx.check_hostname = False
        ctx.verify_mode = ssl.CERT_NONE
    token = base64.b64encode(f"{username}:{password}".encode()).decode()

    container = _discover_default_container(endpoint, port, token, ctx)
    ext_id = container.get("containerExtId") or container.get("extId") or ""
    if not ext_id:
        logger.warning("NUTANIX_STORAGE_CONTAINER: default container ext_id not found; "
                       "leaving %s untouched", env_file)
        return

    name = container.get("name") or ""
    if dry_run:
        log(f"  - [dry-run] would set NUTANIX_STORAGE_CONTAINER={ext_id} "
            f"(default container {name!r}) in {env_file}")
        return

    if _set_env_file_var(env_file, "NUTANIX_STORAGE_CONTAINER", ext_id):
        log(f"{env_file}: NUTANIX_STORAGE_CONTAINER={ext_id} (default container {name!r})")


def main() -> int:
    parser = argparse.ArgumentParser(description="Fill test_config_v2.json from .env + IAM API")
    # Prefer the nested YAML config when present, fall back to the flat .env.
    default_env = SCRIPT_DIR / "config.yaml"
    if not default_env.exists():
        default_env = SCRIPT_DIR / ".env"
    parser.add_argument("--env", default=str(default_env),
                        help="path to the config file: YAML (.yaml/.yml) or flat .env "
                             "(default: testenv/config.yaml if present, else testenv/.env)")
    parser.add_argument("--config", default=str(REPO_ROOT / "test_config_v2.json"),
                        help="path to test_config_v2.json to update")
    parser.add_argument("--config-v3", default=str(REPO_ROOT / "test_config.json"),
                        help="path to the legacy test_config.json to update "
                             "(mirrors preEnv/create_json_v3.tf)")
    parser.add_argument("--skip-v3", action="store_true",
                        help="do not touch the legacy test_config.json")
    parser.add_argument("--env-file", default=str(REPO_ROOT / ".env"),
                        help="path to the root .env whose NUTANIX_STORAGE_CONTAINER "
                             "is updated with the discovered default-container ext_id "
                             "(default: <repo>/.env)")
    parser.add_argument("--skip-env-update", action="store_true",
                        help="do not update NUTANIX_STORAGE_CONTAINER in the root .env")
    parser.add_argument("--skip-robo-flags", action="store_true",
                        help="do not SSH into the PE nodes to add the ROBO small-cluster "
                             "flags to /etc/nutanix/hardware_config.json (+ genesis restart)")
    parser.add_argument("--dry-run", action="store_true",
                        help="print the merged result; make no writes and no API mutations")
    parser.add_argument("--skip-iam", action="store_true",
                        help="only fill the offline top-level fields (no PC connection)")
    parser.add_argument("--only", metavar="SECTIONS",
                        help="build/discover ONLY these comma-separated test_config_v2.json "
                             "sections and leave the rest untouched, e.g. --only prism. "
                             "Implies skipping the legacy test_config.json, the "
                             "NUTANIX_STORAGE_CONTAINER .env update and the ROBO flags "
                             "(unless a listed section needs them). Recognised sections: "
                             + ", ".join(TOP_LEVEL_SECTIONS))
    parser.add_argument("--log-dir", default=str(SCRIPT_DIR / "logs"),
                        help="directory for the run log file (default: testenv/logs)")
    parser.add_argument("--no-log-file", action="store_true",
                        help="log only to the console; do not write a log file")
    parser.add_argument("--verbose", "-v", action="store_true",
                        help="show DEBUG output (incl. HTTP traces) on the console too")
    args = parser.parse_args()

    log_dir = None if args.no_log_file else Path(args.log_dir)
    log_file = setup_logging(log_dir, args.verbose)

    # --only: restrict work to specific top-level sections of test_config_v2.json.
    only = {s.strip() for s in (args.only or "").split(",") if s.strip()}
    unknown = only - set(TOP_LEVEL_SECTIONS)
    if unknown:
        logger.error("--only: unknown section(s): %s. Recognised: %s",
                     ", ".join(sorted(unknown)), ", ".join(TOP_LEVEL_SECTIONS))
        return 2

    logger.info("=== fill_test_config started ===")
    logger.info("env=%s config=%s dry_run=%s skip_iam=%s only=%s", args.env, args.config,
                args.dry_run, args.skip_iam, ",".join(sorted(only)) or "<all>")
    if log_file:
        logger.info("Full debug log: %s", log_file)

    try:
        env = Env(load_env(Path(args.env)))

        config_path = Path(args.config)
        if config_path.exists():
            existing = json.loads(config_path.read_text())
            logger.debug("Loaded existing config (%d top-level keys)", len(existing))
        else:
            logger.info("NOTE: %s does not exist; starting from an empty document", config_path)
            existing = {}

        # --only decouples discovery from IAM: in normal runs discovery is tied to
        # the PC connection (disabled by --skip-iam), but when targeting specific
        # sections we still want their discovery (e.g. prism deploy-IP resolution)
        # while skipping IAM unless it was explicitly requested.
        if only:
            run_iam = "iam" in only
            discover = True
        else:
            run_iam = not args.skip_iam
            discover = not args.skip_iam
        only_top = only - {"iam"}
        build_top = (not only) or bool(only_top)
        # --only targets test_config_v2.json only; skip the unrelated side effects
        # (legacy v3 config, storage-container .env update, ROBO node flags).
        skip_v3 = args.skip_v3 or bool(only)
        skip_env_update = args.skip_env_update or bool(only)
        skip_robo_flags = args.skip_robo_flags or bool(only)

        # The top-level work (cluster VIP/iSCSI discovery + SSL generation) and the
        # IAM provisioning hit different endpoints and share no state, so run them
        # concurrently to overlap the network round-trips.
        with concurrent.futures.ThreadPoolExecutor(max_workers=2) as executor:
            top_future = (executor.submit(build_top_level, env, args.dry_run, discover,
                                          only_top) if build_top else None)
            iam_future = executor.submit(build_iam, env, args.dry_run) if run_iam else None
            if not run_iam:
                logger.info("Skipping IAM provisioning (%s)",
                            "--only without iam" if only else "--skip-iam")

            update = top_future.result() if top_future is not None else {}

            if iam_future is not None:
                try:
                    iam = iam_future.result()
                    if iam:
                        update["iam"] = iam
                    place_test_idp_metadata(env, dry_run=args.dry_run)
                except Exception as exc:  # noqa: BLE001 - surface a clean message
                    logger.error("ERROR while provisioning IAM: %s", exc)
                    logger.debug("IAM provisioning traceback", exc_info=True)
                    logger.error("Tip: re-run with --skip-iam to fill only the offline fields, "
                                 "or fix the connection / IAM_* values in your .env")
                    return 1

        merged = deep_merge(existing, update)
        rendered = json.dumps(merged, indent=2, ensure_ascii=False) + "\n"
        logger.debug("Merged config keys: %s", ", ".join(sorted(merged)))

        # --- legacy test_config.json ("v3") ---
        v3_path = Path(args.config_v3)
        rendered_v3 = None
        if not skip_v3:
            if v3_path.exists():
                existing_v3 = json.loads(v3_path.read_text())
                logger.debug("Loaded existing test_config.json (%d top-level keys)",
                             len(existing_v3))
            else:
                logger.info("NOTE: %s does not exist; starting from an empty document", v3_path)
                existing_v3 = {}
            v3_update = build_v3_test_config(env, existing_v3, dry_run=args.dry_run,
                                             discover=run_iam)
            merged_v3 = deep_merge(existing_v3, v3_update)
            rendered_v3 = json.dumps(merged_v3, indent=2, ensure_ascii=False) + "\n"
        else:
            logger.info("Skipping the legacy test_config.json (%s)",
                        "--only" if only else "--skip-v3")

        # --- root .env: NUTANIX_STORAGE_CONTAINER = default container ext_id ---
        if skip_env_update:
            logger.info("Skipping NUTANIX_STORAGE_CONTAINER update (%s)",
                        "--only" if only and not args.skip_env_update else "--skip-env-update")
        elif run_iam:
            update_env_storage_container(env, Path(args.env_file), dry_run=args.dry_run)
        else:
            logger.info("Skipping NUTANIX_STORAGE_CONTAINER update (discovery disabled)")

        # ROBO / small-cluster workaround: add the hardware_config.json flags and
        # restart genesis on the cluster nodes (cluster.nodes) plus the DP remote
        # PE (data_protection.remote_cluster_pe). These IPs change every run, so
        # they are always taken from the current config.
        robo_hosts = []
        if not skip_robo_flags:
            robo_hosts = [n for n in env.list("CLUSTER_NODES") if n]
            remote_pe = env.get("DP_REMOTE_CLUSTER_PE")
            if remote_pe:
                robo_hosts.append(remote_pe)

        if args.dry_run:
            logger.info("--- dry run: merged test_config_v2.json (printed to stdout) ---")
            print(rendered)
            if rendered_v3 is not None:
                logger.info("--- dry run: merged test_config.json (printed to stdout) ---")
                print(rendered_v3)
            if not skip_robo_flags:
                apply_robo_flags(env, robo_hosts, dry_run=True)
            return 0

        config_path.write_text(rendered)
        logger.info("Wrote %s", config_path)
        if rendered_v3 is not None:
            v3_path.write_text(rendered_v3)
            logger.info("Wrote %s", v3_path)

        if skip_robo_flags:
            logger.info("Skipping ROBO hardware_config.json flags (%s)",
                        "--only" if only and not args.skip_robo_flags else "--skip-robo-flags")
        else:
            apply_robo_flags(env, robo_hosts, dry_run=False)
        return 0
    except SystemExit:
        raise
    except Exception:  # noqa: BLE001 - log full traceback before exiting non-zero
        logger.exception("Unexpected failure")
        return 1
    finally:
        logger.info("=== fill_test_config finished ===")


if __name__ == "__main__":
    raise SystemExit(main())
