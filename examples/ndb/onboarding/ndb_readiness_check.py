#!/usr/bin/env python3
"""
NDB storage-readiness precheck for onboarding.

What it does:
1) Find existing PE cluster (by name/IP), or register it if missing.
2) Poll /era/v0.9/clusters/i/{cloudId}/storage-containers.
3) Print PASS/FAIL diagnosis before Terraform apply.
"""

import argparse
import json
import sys
import time
from typing import Any, Dict, List, Optional

import requests
import urllib3


def normalize_endpoint(endpoint: str) -> str:
    raw = endpoint.strip()
    if raw.startswith("https://"):
        raw = raw[len("https://") :]
    if raw.startswith("http://"):
        raw = raw[len("http://") :]
    if ":" not in raw:
        raw = f"{raw}:8443"
    return raw


def parse_clusters(payload: Any) -> List[Dict[str, Any]]:
    if isinstance(payload, list):
        return payload
    if isinstance(payload, dict):
        entities = payload.get("entities")
        if isinstance(entities, list):
            return entities
    return []


def first_ip(cluster: Dict[str, Any]) -> str:
    ips = cluster.get("ipAddresses") or cluster.get("ipaddresses") or []
    if isinstance(ips, list) and ips:
        return str(ips[0])
    return ""


def find_cluster(clusters: List[Dict[str, Any]], pe_name: str, pe_ip: str) -> Optional[Dict[str, Any]]:
    for c in clusters:
        if c.get("name") == pe_name or first_ip(c) == pe_ip:
            return c
    return None


def register_cluster(session: requests.Session, base: str, pe_name: str, pe_ip: str, pe_user: str, pe_pass: str) -> requests.Response:
    payload_ladder = [
        {
            "name": pe_name,
            "cloudType": "NTNX",
            "version": "v2",
            "description": "",
            "ipAddresses": [pe_ip],
            "credentials": {"username": pe_user, "password": pe_pass},
            "status": "UP",
        },
        {
            "name": pe_name,
            "cloudType": "NTNX",
            "version": "v2",
            "description": "",
            "ipAddresses": [pe_ip],
            "username": pe_user,
            "password": pe_pass,
            "status": "UP",
        },
        {
            "clusterName": pe_name,
            "cloudType": "NTNX",
            "version": "v2",
            "clusterIP": pe_ip,
            "credentialsInfo": {"username": pe_user, "password": pe_pass},
        },
    ]

    last_resp: Optional[requests.Response] = None
    for idx, payload in enumerate(payload_ladder, start=1):
        resp = session.post(f"{base}/clusters", json=payload, timeout=60)
        body = (resp.text or "").replace("\n", " ")[:220]
        print(f"[register {idx}] status={resp.status_code} body={body}")
        last_resp = resp
        if resp.status_code in (200, 201, 202):
            return resp
    assert last_resp is not None
    return last_resp


def main() -> int:
    parser = argparse.ArgumentParser(description="Check NDB storage readiness before Terraform onboarding.")
    parser.add_argument("--ndb-endpoint", required=True, help="NDB endpoint (IP[:port] or URL)")
    parser.add_argument("--ndb-user", default="admin")
    parser.add_argument("--ndb-pass", required=True)
    parser.add_argument("--pe-name", required=True)
    parser.add_argument("--pe-ip", required=True)
    parser.add_argument("--pe-user", required=True)
    parser.add_argument("--pe-pass", required=True)
    parser.add_argument("--attempts", type=int, default=18, help="Storage probe attempts")
    parser.add_argument("--interval-seconds", type=int, default=10, help="Delay between attempts")
    parser.add_argument("--register-if-missing", action="store_true", default=True)
    args = parser.parse_args()

    urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
    endpoint = normalize_endpoint(args.ndb_endpoint)
    base = f"https://{endpoint}/era/v0.9"

    session = requests.Session()
    session.verify = False
    session.auth = (args.ndb_user, args.ndb_pass)
    session.headers.update({"Accept": "application/json"})

    print(f"NDB base: {base}")
    print("== step1: list/find cluster ==")
    list_resp = session.get(f"{base}/clusters", timeout=30)
    print(f"clusters status={list_resp.status_code}")
    if list_resp.status_code != 200:
        print("FAIL: unable to list clusters")
        return 2

    clusters = parse_clusters(list_resp.json())
    cluster = find_cluster(clusters, args.pe_name, args.pe_ip)

    if cluster is None and args.register_if_missing:
        print("Cluster not found; registering PE cluster...")
        reg_resp = register_cluster(session, base, args.pe_name, args.pe_ip, args.pe_user, args.pe_pass)
        if reg_resp.status_code not in (200, 201, 202):
            print("FAIL: registration did not succeed.")
            return 2
        # Re-list after registration.
        list_resp = session.get(f"{base}/clusters", timeout=30)
        clusters = parse_clusters(list_resp.json())
        cluster = find_cluster(clusters, args.pe_name, args.pe_ip)

    if cluster is None:
        print("FAIL: PE cluster not found and not registered.")
        return 2

    cluster_id = cluster.get("id")
    print(
        "Cluster:",
        json.dumps(
            {
                "id": cluster_id,
                "name": cluster.get("name"),
                "ip": first_ip(cluster),
                "status": cluster.get("status"),
            }
        ),
    )
    if not cluster_id:
        print("FAIL: cluster id missing")
        return 2

    print("== step2: storage readiness poll ==")
    url = f"{base}/clusters/i/{cluster_id}/storage-containers"
    last_body = ""
    for idx in range(1, args.attempts + 1):
        resp = session.get(url, timeout=30)
        last_body = (resp.text or "").replace("\n", " ")[:260]
        print(f"[probe {idx:02d}] status={resp.status_code} body={last_body}")
        if resp.status_code == 200:
            print("PASS: storage endpoint is ready.")
            print(json.dumps({"result": "PASS", "cluster_id": cluster_id, "endpoint": url}))
            return 0
        if idx < args.attempts:
            time.sleep(args.interval_seconds)

    print("FAIL: storage endpoint not ready in probe window.")
    print(json.dumps({"result": "FAIL", "cluster_id": cluster_id, "endpoint": url, "last_body": last_body}))
    return 3


if __name__ == "__main__":
    sys.exit(main())
