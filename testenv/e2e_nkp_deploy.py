#!/usr/bin/env python3
"""
End-to-end Nutanix Kubernetes Platform (NKP) + Flow CNI deployment orchestrator.

This is the single, consolidated script (it replaces the former trio
automate_nkp.py / automate_nkp_vpc.py / the SSH-driven e2e_nkp_deploy.py).

Execution model: the whole workflow must run ON THE BASTION VM (the machine that
has docker/kind and where the nkp CLI bootstraps a management cluster). You can
launch this script from anywhere -- when it is not already on the bastion it
SSHes into the bastion (pc_prep.nkp.bastion_*), copies itself + config.yaml over,
and re-executes itself there, streaming the output back. Pass --local (or set
NKP_ON_BASTION=1) to force it to run the workflow on the current machine.

Once on the bastion it mechanises the whole NKP_FLOW_CNI_RUNBOOK.md path end to
end:

  1. Preflight the bastion tools (nkp, kubectl, helm, docker) + SSH keypair.
  2. `nkp create bootstrap`                       -> local management (kind) cluster.
  3. Resolve a free control-plane VIP + MetalLB LB range (config or auto-discover).
  4. `nkp create cluster nutanix --dry-run -o yaml`, then REMOVE addons.cni
     (BYO-CNI) so NKP installs no CNI, and `kubectl apply` + wait ControlPlaneReady.
  5. Fetch the workload kubeconfig; pre-apply the Flow CRDs and the harbor pull
     secret on the workload cluster.
  6. Install Flow CNI with a direct `helm upgrade --install` (stable release name,
     pre-stamped namespaces, --wait --atomic) against the validated priyankar
     (harbor) build, with version-matched image overrides; wait nodes Ready.
     NOTE: a caaph HelmChartProxy is intentionally NOT used -- caaph pre-creates the
     release namespace as a plain namespace, which collides with the flow-cni-system
     namespace the chart templates, and it renames the release on every retry so it
     deadlocks on "cannot be imported ... missing managed-by".
  7. Read the Konnector-registered cluster extId from Prism Central, then
     Flow-activate the cluster (PATCH its kubeconfig to PC) so the extId becomes a
     valid VPC `kubernetesClusters` reference.

Why BYO-CNI instead of NKP's native Flow provider: the native
`addons.cni.provider: Flow` renders placeholder images `nutanix/flow-*:1.0.0`
that exist in no reachable registry, and its managed addon ignores value
overrides. See NKP_FLOW_CNI_RUNBOOK.md for the full rationale/troubleshooting.

Single source of truth for all settings: testenv/config.yaml
  - Prism Central connection:   top-level  pc.{endpoint,port,username,password}
  - NKP cluster + Flow spec:     pc_prep.nkp.*  (incl. pc_prep.nkp.flow_cni.*)
Override the config path with the NKP_CONFIG environment variable if needed.

Local dependencies:  pip3 install requests pyyaml
"""

from __future__ import annotations

import base64
import getpass
import ipaddress
import json
import os
import platform
import shlex
import shutil
import subprocess
import sys
import tempfile
import time
from datetime import datetime

import requests
import urllib3

try:
    import yaml
except ImportError:
    print("[-] Missing dependency 'PyYAML'. Install it with: pip3 install pyyaml")
    sys.exit(1)

# Self-signed Prism Central certs -> silence the insecure-HTTPS noise.
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


# =============================================================================
# CONFIGURATION  (single source of truth: testenv/config.yaml)
# =============================================================================

_CONFIG_PATH = os.environ.get(
    "NKP_CONFIG",
    os.path.join(os.path.dirname(os.path.abspath(__file__)), "config.yaml"),
)


def _load_config(path: str) -> dict:
    try:
        with open(path) as fh:
            return yaml.safe_load(fh) or {}
    except FileNotFoundError:
        print(f"[-] Config file not found: {path}")
        sys.exit(1)
    except yaml.YAMLError as exc:
        print(f"[-] Failed to parse {path}: {exc}")
        sys.exit(1)


def _require(value, label: str):
    if value in (None, ""):
        print(f"[-] Missing required config value: {label} (in {_CONFIG_PATH})")
        sys.exit(1)
    return value


_CFG = _load_config(_CONFIG_PATH)
_PC = _CFG.get("pc", {}) or {}
_NKP = (_CFG.get("pc_prep", {}) or {}).get("nkp", {}) or {}

# --- Prism Central (pc.*) ---
PC_IP = str(_require(_PC.get("endpoint"), "pc.endpoint"))
PC_PORT = str(_PC.get("port", "9440"))
PC_USERNAME = str(_require(_PC.get("username"), "pc.username"))
PC_PASSWORD = str(_require(_PC.get("password"), "pc.password"))

# --- NKP cluster specs (pc_prep.nkp.*) ---
CLUSTER_NAME = str(_require(_NKP.get("cluster_name"), "pc_prep.nkp.cluster_name"))
NKP_OS_IMAGE_NAME = str(_require(_NKP.get("os_image_name"), "pc_prep.nkp.os_image_name"))
NKP_SUBNET_NAME = str(_require(_NKP.get("subnet_name"), "pc_prep.nkp.subnet_name"))
# PE_CLUSTER_NAME and STORAGE_CONTAINER are discovered from Prism Central at
# runtime via the clustermgmt v4 API (see discover_pe_cluster / discover_storage_
# container); they are intentionally NOT read from config.

# Control-plane VIP + MetalLB LB range. Empty => auto-discover free static IPs.
CONTROL_PLANE_VIP = str(_NKP.get("control_plane_vip", "") or "")
LB_IP_RANGE = str(_NKP.get("lb_ip_range", "") or "")
LB_IP_COUNT = int(_NKP.get("lb_ip_count", 2))
_STATIC_WINDOW = _NKP.get("static_ip_window", {}) or {}
STATIC_IP_WINDOW_START = str(_STATIC_WINDOW.get("start", "") or "")
STATIC_IP_WINDOW_END = str(_STATIC_WINDOW.get("end", "") or "")

CP_REPLICAS = int(_NKP.get("control_plane_replicas", 1))
CP_VCPUS = int(_NKP.get("control_plane_vcpus", 4))
CP_MEMORY_GIB = int(_NKP.get("control_plane_memory_gib", 16))
WORKER_REPLICAS = int(_NKP.get("worker_replicas", 3))
WORKER_VCPUS = int(_NKP.get("worker_vcpus", 8))
WORKER_MEMORY_GIB = int(_NKP.get("worker_memory_gib", 16))

# --- Flow CNI (pc_prep.nkp.flow_cni.*) ---
_FLOW = _NKP.get("flow_cni", {}) or {}
FLOW_CNI_ENABLED = bool(_FLOW.get("enabled", False))
FLOW_HARBOR = str(_FLOW.get("harbor", "harbor.eng.nutanix.com"))
FLOW_CHART_REPO = str(_FLOW.get("chart_repo", f"oci://{FLOW_HARBOR}/priyankar"))
FLOW_CHART_VERSION = str(_FLOW.get("chart_version", "1.1.0-937"))
FLOW_K8S_CNI_REPO = str(_FLOW.get("flow_k8s_cni_repo", f"{FLOW_HARBOR}/priyankar/flow-k8s-cni"))
FLOW_K8S_CNI_TAG = str(_FLOW.get("flow_k8s_cni_tag", "937"))
FLOW_OVN_REPO = str(_FLOW.get("ovn_repo", f"{FLOW_HARBOR}/priyankar/flow-ovn-kubernetes"))
FLOW_OVN_TAG = str(_FLOW.get("ovn_tag", "222"))
FLOW_REGISTRY_USER = str(_FLOW.get("registry_username", "svc.flow_nw"))
FLOW_DOCKER_AUTH_B64 = str(_FLOW.get("docker_auth_b64", "") or "")
# CAPI Cluster namespace; empty => auto-detect after create.
CAPI_NAMESPACE_CFG = str(_FLOW.get("capi_namespace", "") or "")

# --- SSH public key injected into the nodes (generated on the bastion if missing) ---
SSH_PRIVATE_KEY = os.path.expanduser("~/.ssh/id_rsa")
SSH_PUBLIC_KEY = SSH_PRIVATE_KEY + ".pub"

# --- Bastion VM (the script copies itself here and runs the whole workflow on it) ---
# All heavy lifting (nkp/docker/kind/kubectl/helm) must happen ON the bastion, so
# when launched anywhere else this script SSHes into the bastion and re-executes
# itself there. bastion_host may be an IP/DNS name; if empty it is discovered from
# Prism Central by bastion_vm_name.
BASTION_VM_NAME = str(_NKP.get("bastion_vm_name", "") or "")
BASTION_HOST = str(_NKP.get("bastion_host", "") or "")
BASTION_SSH_USER = str(_NKP.get("bastion_ssh_username", "ubuntu") or "ubuntu")
BASTION_SSH_PASSWORD = str(_NKP.get("bastion_ssh_password", "") or "")
BASTION_SSH_KEY_PATH = os.path.expanduser(str(_NKP.get("bastion_ssh_key_path", "") or "")) \
    if _NKP.get("bastion_ssh_key_path") else ""
# Remote working directory the script + config are copied into on the bastion.
BASTION_REMOTE_DIR = str(_NKP.get("bastion_remote_dir", "/tmp/nkp-deploy") or "/tmp/nkp-deploy")
# NKP CLI tarball URL (auto-installed on the bastion when the nkp CLI is missing).
NKP_CLI_DOWNLOAD_URL = str(_NKP.get("cli_download_url", "") or "")
# Sentinel: set on the bastion re-exec so we run the workflow instead of re-remoting.
ON_BASTION = os.environ.get("NKP_ON_BASTION") == "1"

# --- Work directory for generated manifests + kubeconfig ---
WORKDIR = os.environ.get("NKP_WORKDIR", os.path.expanduser("~"))
CLUSTER_MANIFEST = os.path.join(WORKDIR, f"cluster-{CLUSTER_NAME}.yaml")
FLOW_VALUES = os.path.join(WORKDIR, "flow-cni-values.yaml")
WORKLOAD_KUBECONFIG = os.path.join(WORKDIR, f"{CLUSTER_NAME}.conf")

# Stable Helm release name for the Flow CNI chart. We install it DIRECTLY with
# `helm upgrade --install` (not via a caaph HelmChartProxy): caaph pre-creates the
# release namespace as a PLAIN namespace, which collides with the flow-cni-system
# Namespace the chart templates itself ("cannot be imported ... missing managed-by"),
# and caaph renames the release on every retry so it never recovers. A direct
# install with a fixed release name + pre-stamped namespaces is deterministic.
FLOW_RELEASE = "nutanix-flow-cni"
# Namespaces the Flow chart owns/templates. We pre-stamp any that already exist with
# Helm ownership for FLOW_RELEASE so `helm install` adopts them cleanly on re-runs
# (and after a prior failed attempt) instead of erroring on ownership metadata.
_FLOW_NAMESPACES = ("flow-cni-system", "ovn-kubernetes", "flow-cns-system",
                    "ovn-host-network")

# --- Timeouts / polling ---
REGISTRATION_RETRIES = 12
REGISTRATION_DELAY_S = 30
POST_CREATE_GRACE_S = 60
CONTROL_PLANE_TIMEOUT = "30m"
FLOW_READY_RETRIES = 40
FLOW_READY_DELAY_S = 15

V4_BASE = f"https://{PC_IP}:{PC_PORT}"
_SUBNETS_V42 = "api/networking/v4.2/config/subnets"
_CLUSTERS_V40 = "api/clustermgmt/v4.0/config/clusters"
_STORAGE_CONTAINERS_V40 = "api/clustermgmt/v4.0/config/storage-containers"
# OData filter for the AOS-function (Prism Element) cluster, mirroring
# terraform/main.tf's aosFilter; excludes the PRISM_CENTRAL cluster itself.
_AOS_CLUSTER_FILTER = ("config/clusterFunction/any(t:t eq "
                       "Clustermgmt.Config.ClusterFunctionRef'AOS')")
_STORAGE_CONTAINER_PREFIX = "default-container-"


# =============================================================================
# Logging + subprocess helpers
# =============================================================================


def log(msg: str, level: str = "*") -> None:
    ts = datetime.now().strftime("%H:%M:%S")
    print(f"[{ts}] [{level}] {msg}", flush=True)


def die(msg: str) -> None:
    log(msg, level="-")
    sys.exit(1)


def _nkp_env() -> dict:
    """nkp/CAPX read PC creds from the environment, not the flags. Without these
    it aborts with: expected environment var "NUTANIX_USER" not found."""
    env = os.environ.copy()
    env["NUTANIX_USER"] = PC_USERNAME
    env["NUTANIX_PASSWORD"] = PC_PASSWORD
    env["NUTANIX_ENDPOINT"] = PC_IP
    env["NUTANIX_PORT"] = PC_PORT
    env["NUTANIX_INSECURE"] = "true"
    return env


def sh(cmd: list[str], *, check: bool = True, capture: bool = False,
       stdout_path: str | None = None, quiet: bool = False) -> subprocess.CompletedProcess:
    """Run a command locally. Streams to the console unless capture/stdout_path
    is set. Uses the nkp env (PC creds) for every call -- harmless for kubectl/helm."""
    if not quiet:
        log("$ " + " ".join(cmd))
    kwargs: dict = {"env": _nkp_env(), "text": True}
    if capture:
        kwargs["stdout"] = subprocess.PIPE
        kwargs["stderr"] = subprocess.PIPE
    elif stdout_path:
        kwargs["stderr"] = subprocess.PIPE
    try:
        if stdout_path:
            with open(stdout_path, "w") as out:
                kwargs["stdout"] = out
                proc = subprocess.run(cmd, **kwargs)
        else:
            proc = subprocess.run(cmd, **kwargs)
    except FileNotFoundError:
        die(f"Command not found: {cmd[0]}. Install it on the bastion and retry.")
    if check and proc.returncode != 0:
        detail = (proc.stderr or "").strip() if (capture or stdout_path) else ""
        die(f"Command failed (exit {proc.returncode}): {' '.join(cmd)}\n{detail}")
    return proc


def kubectl(args: list[str], *, workload: bool = False, **kw) -> subprocess.CompletedProcess:
    """kubectl against the management (kind) cluster by default, or the workload
    cluster when workload=True (uses the fetched kubeconfig)."""
    base = ["kubectl"]
    if workload:
        base += ["--kubeconfig", WORKLOAD_KUBECONFIG]
    return sh(base + args, **kw)


# =============================================================================
# Prism Central v4 REST helpers
# =============================================================================


def _v4_session() -> requests.Session:
    s = requests.Session()
    s.auth = (PC_USERNAME, PC_PASSWORD)
    s.verify = False
    s.headers.update({"Content-Type": "application/json", "Accept": "application/json"})
    return s


def v4_get(session: requests.Session, path: str) -> dict:
    resp = session.get(f"{V4_BASE}/{path.lstrip('/')}", timeout=30)
    resp.raise_for_status()
    return resp.json()


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


# =============================================================================
# Bastion remoting: run the whole workflow ON the bastion VM
# =============================================================================


_SSH_OPTS = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null",
             "-o", "LogLevel=ERROR", "-o", "ConnectTimeout=20",
             "-o", "ServerAliveInterval=30"]


def _sshpass_prefix() -> list[str]:
    """sshpass prefix for password auth (only when no key path is configured)."""
    if BASTION_SSH_KEY_PATH:
        return []
    if BASTION_SSH_PASSWORD:
        if not shutil.which("sshpass"):
            die("pc_prep.nkp.bastion_ssh_password is set but 'sshpass' is not installed "
                "locally.\n      Install it (brew install hudochenkov/sshpass/sshpass | "
                "apt-get install sshpass) or set pc_prep.nkp.bastion_ssh_key_path.")
        return ["sshpass", "-p", BASTION_SSH_PASSWORD]
    return []


def _key_opts() -> list[str]:
    return ["-i", BASTION_SSH_KEY_PATH] if BASTION_SSH_KEY_PATH else []


def _discover_bastion_host(session: requests.Session) -> str:
    """Resolve the bastion VM's IP from Prism Central by name (v3 vms/list)."""
    if not BASTION_VM_NAME:
        die("Cannot reach the bastion: set pc_prep.nkp.bastion_host (IP/DNS) or "
            "pc_prep.nkp.bastion_vm_name so it can be discovered from Prism Central.")
    log(f"Discovering bastion VM {BASTION_VM_NAME!r} IP from Prism Central...")
    try:
        resp = session.post(f"{V4_BASE}/api/nutanix/v3/vms/list",
                            json={"kind": "vm", "filter": f"vm_name=={BASTION_VM_NAME}"},
                            timeout=30)
        resp.raise_for_status()
        for entity in resp.json().get("entities", []):
            nics = _dig(entity, "status.resources.nic_list") or []
            for nic in nics:
                for ep in nic.get("ip_endpoint_list", []):
                    ip = ep.get("ip")
                    if ip:
                        return ip
    except requests.RequestException as exc:
        die(f"Failed to query Prism Central for the bastion IP: {exc}")
    die(f"Bastion VM {BASTION_VM_NAME!r} has no learned IP yet. Set "
        "pc_prep.nkp.bastion_host explicitly, or wait for the VM to get an IP.")
    return ""  # unreachable (die exits)


def run_on_bastion() -> int:
    """Copy this script + config to the bastion and re-execute it there, so the
    entire NKP/Flow workflow (docker/kind/nkp/kubectl/helm) runs on the bastion."""
    host = BASTION_HOST or _discover_bastion_host(_v4_session())
    target = f"{BASTION_SSH_USER}@{host}"
    log(f"Running the deployment on the bastion {target} (dir {BASTION_REMOTE_DIR}).")

    script_path = os.path.abspath(__file__)
    remote_script = f"{BASTION_REMOTE_DIR}/{os.path.basename(script_path)}"
    remote_config = f"{BASTION_REMOTE_DIR}/config.yaml"

    ssh_cmd = _sshpass_prefix() + ["ssh"] + _SSH_OPTS + _key_opts()
    scp_cmd = _sshpass_prefix() + ["scp"] + _SSH_OPTS + _key_opts()

    # 1. Ensure the remote dir exists, then copy the script + config over.
    sh(ssh_cmd + [target, f"mkdir -p {BASTION_REMOTE_DIR}"])
    sh(scp_cmd + [script_path, _CONFIG_PATH, f"{target}:{BASTION_REMOTE_DIR}/"])
    # scp lands config under its original basename; normalise to config.yaml.
    src_cfg = f"{BASTION_REMOTE_DIR}/{os.path.basename(_CONFIG_PATH)}"
    if os.path.basename(_CONFIG_PATH) != "config.yaml":
        sh(ssh_cmd + [target, f"mv -f {src_cfg} {remote_config}"])

    # 2. Make sure the bastion's python3 has the script's deps (requests, pyyaml).
    #    Prefer Ubuntu's apt packages (no pip needed); fall back to bootstrapping
    #    pip (apt python3-pip -> ensurepip -> get-pip.py) only if that's not enough.
    log("Ensuring python3 + requests/pyyaml on the bastion...")
    prep = (
        'set -u; export NKP_SUDO_PASS=' + shlex.quote(BASTION_SSH_PASSWORD) + '; '
        'sudo_do() { if [ "$(id -u)" = 0 ]; then "$@"; '
        'elif [ -n "$NKP_SUDO_PASS" ]; then printf "%s\\n" "$NKP_SUDO_PASS" | sudo -S -p "" "$@"; '
        'else sudo -n "$@"; fi; }; '
        'have() { python3 -c "import $1" >/dev/null 2>&1; }; '
        'command -v python3 >/dev/null || { sudo_do apt-get update -y && sudo_do apt-get install -y python3; }; '
        # First choice: distro packages (universe for python3-requests).
        'if ! have yaml || ! have requests; then '
        '  sudo_do apt-get update -y || true; '
        '  sudo_do add-apt-repository -y universe >/dev/null 2>&1 || true; '
        '  sudo_do apt-get install -y python3-yaml python3-requests || true; '
        'fi; '
        # Fallback: bootstrap pip and pip-install into the user site.
        'if ! have yaml || ! have requests; then '
        '  python3 -m pip --version >/dev/null 2>&1 || sudo_do apt-get install -y python3-pip || true; '
        '  python3 -m pip --version >/dev/null 2>&1 || python3 -m ensurepip --upgrade || true; '
        '  python3 -m pip --version >/dev/null 2>&1 || (curl -fsSL https://bootstrap.pypa.io/get-pip.py | python3 - --user) || true; '
        '  python3 -m pip install --user --quiet --disable-pip-version-check requests pyyaml || true; '
        'fi; '
        'have yaml && have requests || '
        '{ echo "[-] Could not install requests/pyyaml on the bastion" >&2; exit 1; }'
    )
    sh(ssh_cmd + [target, prep])

    # 3. Re-exec on the bastion with the sentinel + config path set. Stream output.
    remote_cmd = (
        f"cd {BASTION_REMOTE_DIR} && "
        f"NKP_ON_BASTION=1 NKP_CONFIG={remote_config} "
        f"NKP_SUDO_PASS={shlex.quote(BASTION_SSH_PASSWORD)} "
        f"python3 {shlex.quote(remote_script)}"
    )
    log("Streaming the bastion run (this is the full deployment)...")
    proc = subprocess.run(ssh_cmd + ["-tt", target, remote_cmd])
    return proc.returncode


# =============================================================================
# STEP 1: Preflight tools + SSH key
# =============================================================================


def _os_arch() -> tuple[str, str] | None:
    """Return (os, arch) as the k8s/helm download conventions, or None."""
    system = platform.system().lower()  # 'linux' | 'darwin'
    arch = {"x86_64": "amd64", "amd64": "amd64",
            "aarch64": "arm64", "arm64": "arm64"}.get(platform.machine().lower())
    if system not in ("linux", "darwin") or not arch:
        return None
    return system, arch


def _sudo(args: list[str], *, check: bool = False) -> subprocess.CompletedProcess:
    """Run *args* as root (password-aware). Uses NKP_SUDO_PASS via `sudo -S` when
    set, else assumes passwordless sudo (typical on cloud-init VMs)."""
    if os.geteuid() == 0 or not shutil.which("sudo"):
        return sh(args, check=check)
    sudo_pass = os.environ.get("NKP_SUDO_PASS")
    if sudo_pass:
        cmd = ("printf '%s\\n' \"$NKP_SUDO_PASS\" | sudo -S -p '' "
               + " ".join(shlex.quote(a) for a in args))
        return sh(["bash", "-c", cmd], check=check)
    return sh(["sudo", "-n"] + args, check=check)


def _sudo_install(src: str, dest: str) -> bool:
    """`install -m 0755 src dest`, using sudo (password-aware) when not root."""
    return _sudo(["install", "-m", "0755", src, dest]).returncode == 0


def _install_kubectl() -> bool:
    """Install kubectl into /usr/local/bin (runbook section 1)."""
    oa = _os_arch()
    if not oa:
        log(f"Cannot auto-install kubectl for {platform.system()}/{platform.machine()}.",
            level="-")
        return False
    system, arch = oa
    log(f"Installing kubectl ({system}/{arch})...")
    try:
        stable = subprocess.run(
            ["curl", "-fsSL", "https://dl.k8s.io/release/stable.txt"],
            check=True, capture_output=True, text=True).stdout.strip()
        with tempfile.TemporaryDirectory() as tmp:
            binpath = os.path.join(tmp, "kubectl")
            url = f"https://dl.k8s.io/release/{stable}/bin/{system}/{arch}/kubectl"
            sh(["curl", "-fsSLo", binpath, url])
            if not _sudo_install(binpath, "/usr/local/bin/kubectl"):
                return False
    except (subprocess.CalledProcessError, FileNotFoundError) as exc:
        log(f"kubectl install failed: {exc}", level="-")
        return False
    return shutil.which("kubectl") is not None


def _latest_helm_version() -> str:
    """Latest helm tag (e.g. 'v3.16.3'), falling back to a known-good pin."""
    try:
        data = subprocess.run(
            ["curl", "-fsSL", "https://api.github.com/repos/helm/helm/releases/latest"],
            check=True, capture_output=True, text=True).stdout
        tag = json.loads(data).get("tag_name")
        if tag:
            return tag
    except (subprocess.CalledProcessError, FileNotFoundError, json.JSONDecodeError):
        pass
    return "v3.16.3"


def _install_helm() -> bool:
    """Install helm into /usr/local/bin from the official release tarball."""
    oa = _os_arch()
    if not oa:
        log(f"Cannot auto-install helm for {platform.system()}/{platform.machine()}.",
            level="-")
        return False
    system, arch = oa
    version = _latest_helm_version()
    log(f"Installing helm {version} ({system}/{arch})...")
    try:
        url = f"https://get.helm.sh/helm-{version}-{system}-{arch}.tar.gz"
        with tempfile.TemporaryDirectory() as tmp:
            tgz = os.path.join(tmp, "helm.tar.gz")
            sh(["curl", "-fsSLo", tgz, url])
            sh(["tar", "-xzf", tgz, "-C", tmp])
            binpath = os.path.join(tmp, f"{system}-{arch}", "helm")
            if not _sudo_install(binpath, "/usr/local/bin/helm"):
                return False
    except (subprocess.CalledProcessError, FileNotFoundError) as exc:
        log(f"helm install failed: {exc}", level="-")
        return False
    return shutil.which("helm") is not None


def _install_nkp() -> bool:
    """Install the nkp CLI into /usr/local/bin from pc_prep.nkp.cli_download_url."""
    if not NKP_CLI_DOWNLOAD_URL:
        log("Cannot auto-install nkp: pc_prep.nkp.cli_download_url is empty.", level="-")
        return False
    log("Installing nkp CLI from cli_download_url...")
    try:
        with tempfile.TemporaryDirectory() as tmp:
            tgz = os.path.join(tmp, "nkp.tar.gz")
            sh(["curl", "-fsSLo", tgz, NKP_CLI_DOWNLOAD_URL])
            sh(["tar", "-xzf", tgz, "-C", tmp])
            binpath = os.path.join(tmp, "nkp")
            if not os.path.exists(binpath):  # locate it if nested in the tarball
                for root, _dirs, files in os.walk(tmp):
                    if "nkp" in files:
                        binpath = os.path.join(root, "nkp")
                        break
            if not os.path.exists(binpath):
                log("nkp binary not found in the downloaded tarball.", level="-")
                return False
            if not _sudo_install(binpath, "/usr/local/bin/nkp"):
                return False
    except (subprocess.CalledProcessError, FileNotFoundError) as exc:
        log(f"nkp install failed: {exc}", level="-")
        return False
    return shutil.which("nkp") is not None


_INSTALLERS = {"kubectl": _install_kubectl, "helm": _install_helm, "nkp": _install_nkp}


def _docker_ready() -> bool:
    return subprocess.run(["docker", "info"], stdout=subprocess.DEVNULL,
                          stderr=subprocess.DEVNULL).returncode == 0


def _install_docker() -> bool:
    """Install + start docker.io on a Debian/Ubuntu bastion and make the socket
    usable by the current user for THIS session (group changes need a re-login,
    which the single SSH re-exec doesn't get)."""
    if platform.system().lower() != "linux" or not shutil.which("apt-get"):
        log("Cannot auto-install docker (need a Debian/Ubuntu apt-based bastion).",
            level="-")
        return False
    log("Installing docker.io ...")
    try:
        _sudo(["apt-get", "update", "-y"])
        _sudo(["env", "DEBIAN_FRONTEND=noninteractive", "apt-get", "install", "-y",
               "docker.io"], check=True)
        _sudo(["systemctl", "enable", "--now", "docker"], check=True)
        user = os.environ.get("SUDO_USER") or getpass.getuser() or "ubuntu"
        _sudo(["usermod", "-aG", "docker", user])
        # Group membership won't apply to this SSH session; open the socket so the
        # current user can talk to dockerd right now (test bastion, resets on restart).
        if not _docker_ready():
            _sudo(["chmod", "666", "/var/run/docker.sock"])
    except (subprocess.CalledProcessError, FileNotFoundError) as exc:
        log(f"docker install failed: {exc}", level="-")
        return False
    return _docker_ready()


def preflight_tools() -> None:
    log("Preflight: checking bastion tools (nkp, kubectl, helm, docker)...")
    missing = [t for t in ("nkp", "kubectl", "helm") if not shutil.which(t)]

    # Auto-install what we can (kubectl/helm from upstream, nkp from cli_download_url).
    installable = [t for t in missing if t in _INSTALLERS]
    if installable:
        log("Missing CLI(s): " + ", ".join(installable) + " -- attempting auto-install...")
        for tool in installable:
            if _INSTALLERS[tool]():
                log(f"Installed {tool}.", level="+")
            else:
                die(f"Failed to auto-install {tool}. Install it per "
                    "NKP_FLOW_CNI_RUNBOOK.md section 1 and retry"
                    + (" (set pc_prep.nkp.cli_download_url for nkp)." if tool == "nkp"
                       else "."))

    still_missing = [t for t in missing if not shutil.which(t)]
    if still_missing:
        die("Missing required CLI(s) on the bastion PATH: " + ", ".join(still_missing) +
            ".\n      Install them per NKP_FLOW_CNI_RUNBOOK.md section 1.")

    engine = shutil.which("docker") or shutil.which("podman")
    if not engine:
        log("No container engine (docker/podman) found -- attempting to install docker...")
        if _install_docker():
            engine = shutil.which("docker")
            log("Installed docker.", level="+")
        else:
            die("No container engine (docker/podman) found and auto-install failed. "
                "NKP bootstraps a local kind cluster and needs one. Install docker "
                "and `systemctl enable --now docker`.")
    if engine and engine.endswith("docker") and not _docker_ready():
        # Daemon up but not reachable as this user (group not applied this session).
        _sudo(["systemctl", "enable", "--now", "docker"], check=False)
        if not _docker_ready():
            _sudo(["chmod", "666", "/var/run/docker.sock"], check=False)
        if not _docker_ready():
            die("Docker is installed but the daemon is not reachable. Start it and "
                "ensure your user is in the docker group (re-login after usermod).")
    log(f"Preflight OK (container engine: {engine}).", level="+")


def ensure_ssh_key() -> None:
    if os.path.isfile(SSH_PUBLIC_KEY):
        log(f"Using existing SSH public key: {SSH_PUBLIC_KEY}")
        return
    log(f"Generating SSH keypair at {SSH_PRIVATE_KEY}...")
    os.makedirs(os.path.dirname(SSH_PRIVATE_KEY), exist_ok=True)
    sh(["ssh-keygen", "-t", "rsa", "-b", "4096", "-N", "", "-f", SSH_PRIVATE_KEY])


# =============================================================================
# STEP 2: Bootstrap the management (kind) cluster
# =============================================================================


def _mgmt_capi_ready() -> bool:
    """True if the local management cluster already has the CAPI/CAAPH CRDs
    (i.e. a bootstrap already ran); used to make bootstrap idempotent."""
    proc = kubectl(["get", "crd", "helmchartproxies.addons.cluster.x-k8s.io"],
                   check=False, capture=True, quiet=True)
    return proc.returncode == 0


def create_bootstrap() -> None:
    if _mgmt_capi_ready():
        log("Management cluster already bootstrapped (CAPI CRDs present); skipping.",
            level="+")
        return
    log("Bootstrapping the local management (kind) cluster...")
    sh(["nkp", "create", "bootstrap", "--kubeconfig",
        os.path.expanduser("~/.kube/config")])
    kubectl(["get", "crd", "helmchartproxies.addons.cluster.x-k8s.io"], check=False)


# =============================================================================
# STEP 3: Resolve control-plane VIP + LB range (config or auto-discover)
# =============================================================================


def _subnet_ipam(session: requests.Session, name: str):
    """(found, net, gateway, pools) for subnet *name*; net/pools empty for an
    unmanaged VLAN (external DHCP; PC exposes no CIDR)."""
    from urllib.parse import quote
    fltr = quote(f"name eq '{name}'", safe="")
    body = v4_get(session, f"{_SUBNETS_V42}?$filter={fltr}")
    data = _dig(body, "data")
    if not (isinstance(data, list) and data and isinstance(data[0], dict)):
        return False, None, None, None
    for entry in (data[0].get("ipConfig") or []):
        v4 = entry.get("ipv4") or {}
        subnet = v4.get("ipSubnet") or {}
        ip_val = _dig(subnet, "ip.value")
        prefix = subnet.get("prefixLength")
        if not ip_val or prefix is None:
            continue
        net = ipaddress.ip_network(f"{ip_val}/{prefix}", strict=False)
        gateway = _dig(v4, "defaultGatewayIp.value")
        pools = []
        for pool in (v4.get("poolList") or []):
            start = _dig(pool, "startIp.value")
            end = _dig(pool, "endIp.value")
            if start and end:
                pools.append((int(ipaddress.ip_address(start)),
                              int(ipaddress.ip_address(end))))
        return True, net, gateway, pools
    return True, None, None, []


def _static_candidates(net, gateway, pc_ip, pools, window):
    reserved = set()
    if gateway:
        reserved.add(int(ipaddress.ip_address(gateway)))
    if pc_ip:
        reserved.add(int(ipaddress.ip_address(pc_ip)))
    if net is not None:
        reserved.update({int(net.network_address), int(net.broadcast_address)})
        lo, hi = int(net.network_address) + 1, int(net.broadcast_address) - 1
    else:
        lo = hi = None
    if window:
        w_lo, w_hi = int(ipaddress.ip_address(window[0])), int(ipaddress.ip_address(window[1]))
        lo = w_lo if lo is None else max(lo, w_lo)
        hi = w_hi if hi is None else min(hi, w_hi)
    if lo is None or hi is None:
        return
    for i in range(lo, hi + 1):
        if i in reserved or any(a <= i <= b for a, b in pools):
            continue
        yield i


def _ip_free_local(ip: str) -> bool:
    """Ping *ip* from the bastion (on the node subnet). Non-zero exit -> free."""
    return subprocess.run(["ping", "-c", "1", "-W", "1", ip],
                          stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL).returncode != 0


def _contiguous_run(values: list[int], count: int):
    run: list[int] = []
    for x in values:
        run = run + [x] if (run and x == run[-1] + 1) else [x]
        if len(run) >= count:
            return run[:count]
    return None


def resolve_nkp_ips(session: requests.Session, probe_max: int = 64) -> tuple[str, str]:
    """Return (control_plane_vip, lb_ip_range) from config, or auto-allocate free
    static IPs outside the subnet's IPAM pools / within the static window."""
    vip_cfg, lb_cfg = CONTROL_PLANE_VIP.strip(), LB_IP_RANGE.strip()
    if vip_cfg and lb_cfg:
        log(f"Using explicit VIP={vip_cfg} and LB range={lb_cfg} from config.")
        return vip_cfg, lb_cfg

    log(f"Auto-discovering control-plane VIP + LB range on subnet '{NKP_SUBNET_NAME}'...")
    found, net, gateway, pools = _subnet_ipam(session, NKP_SUBNET_NAME)
    if not found:
        die(f"Subnet '{NKP_SUBNET_NAME}' not found via the networking v4 API.")

    window = None
    if pools:
        log(f"Subnet {net} is IPAM-managed; excluding pools {pools} and picking free IPs.")
    else:
        if not (STATIC_IP_WINDOW_START and STATIC_IP_WINDOW_END):
            die(f"Subnet '{NKP_SUBNET_NAME}' is not IPAM-managed; set static_ip_window "
                f"(or explicit control_plane_vip + lb_ip_range) in config.yaml.")
        window = (STATIC_IP_WINDOW_START, STATIC_IP_WINDOW_END)
        log(f"Subnet not IPAM-managed; picking free IPs in static window "
            f"{window[0]}-{window[1]} (assumed outside DHCP).")

    free: list[int] = []
    vip: int | None = None
    lb_run = None
    probes = 0
    for cand_int in _static_candidates(net, gateway, PC_IP, pools, window):
        if probes >= probe_max:
            break
        probes += 1
        cand = str(ipaddress.ip_address(cand_int))
        if not _ip_free_local(cand):
            log(f"  {cand} is in use, skipping.")
            continue
        free.append(cand_int)
        if vip is None:
            vip = free[0]
        lb_run = _contiguous_run([i for i in free if i != vip], LB_IP_COUNT)
        if vip is not None and lb_run:
            break

    if vip is None or not lb_run:
        die(f"Could not find a free VIP + {LB_IP_COUNT} contiguous free LB IPs "
            f"within {probe_max} probes on '{NKP_SUBNET_NAME}'.")

    vip_ip = str(ipaddress.ip_address(vip))
    lb_range = f"{ipaddress.ip_address(lb_run[0])}-{ipaddress.ip_address(lb_run[-1])}"
    log(f"Selected control-plane VIP={vip_ip}, LB range={lb_range}", level="+")
    return vip_ip, lb_range


# =============================================================================
# STEP 4: Generate the NKPCluster manifest, strip CNI (BYO), apply, wait
# =============================================================================


def discover_pe_cluster(session: requests.Session) -> str:
    """Discover the Prism Element (AOS-function) cluster NAME via clustermgmt v4,
    using the same aosFilter as terraform/main.tf. This excludes the Prism
    Central cluster itself, so we get the PE that hosts the NKP VMs."""
    from urllib.parse import quote
    fltr = quote(_AOS_CLUSTER_FILTER, safe="")
    body = v4_get(session, f"{_CLUSTERS_V40}?$filter={fltr}")
    for cluster in (_dig(body, "data") or []):
        name = _dig(cluster, "name")
        if name:
            log(f"Discovered PE cluster '{name}' (AOS function) via clustermgmt v4.",
                level="+")
            return name
    die("Could not discover an AOS (Prism Element) cluster via clustermgmt v4 "
        f"({_CLUSTERS_V40}?$filter=<aosFilter>).")
    return ""  # unreachable


def discover_storage_container(session: requests.Session) -> str:
    """Discover the CSI storage container NAME via clustermgmt v4 -- the one whose
    name starts with 'default-container-' (mirrors terraform/main.tf). Falls back
    to an unfiltered list + client-side prefix match if the server rejects the
    startswith() filter."""
    from urllib.parse import quote

    def _first_matching(body: object) -> str:
        for sc in (_dig(body, "data") or []):
            name = _dig(sc, "name")
            if name and name.startswith(_STORAGE_CONTAINER_PREFIX):
                return name
        return ""

    fltr = quote(f"startswith(name,'{_STORAGE_CONTAINER_PREFIX}')", safe="")
    try:
        name = _first_matching(v4_get(session, f"{_STORAGE_CONTAINERS_V40}?$filter={fltr}"))
    except requests.RequestException:
        name = ""
    if not name:  # fallback: list all and match client-side
        name = _first_matching(v4_get(session, _STORAGE_CONTAINERS_V40))
    if name:
        log(f"Discovered storage container '{name}' via clustermgmt v4.", level="+")
        return name
    die(f"Could not discover a storage container named '{_STORAGE_CONTAINER_PREFIX}*' "
        f"via clustermgmt v4 ({_STORAGE_CONTAINERS_V40}).")
    return ""  # unreachable


def generate_cluster_manifest(vip: str, lb_range: str, pe_cluster: str,
                              storage_container: str) -> None:
    log("Generating the NKPCluster manifest (nkp create cluster --dry-run)...")
    cmd = [
        "nkp", "create", "cluster", "nutanix",
        "--cluster-name", CLUSTER_NAME,
        "--endpoint", f"https://{PC_IP}:{PC_PORT}",
        "--control-plane-endpoint-ip", vip,
        "--control-plane-prism-element-cluster", pe_cluster,
        "--control-plane-subnets", NKP_SUBNET_NAME,
        "--control-plane-vm-image", NKP_OS_IMAGE_NAME,
        "--control-plane-memory", str(CP_MEMORY_GIB),
        "--control-plane-vcpus", str(CP_VCPUS),
        "--control-plane-replicas", str(CP_REPLICAS),
        "--worker-prism-element-cluster", pe_cluster,
        "--worker-subnets", NKP_SUBNET_NAME,
        "--worker-vm-image", NKP_OS_IMAGE_NAME,
        "--worker-memory", str(WORKER_MEMORY_GIB),
        "--worker-vcpus", str(WORKER_VCPUS),
        "--worker-replicas", str(WORKER_REPLICAS),
        "--csi-storage-container", storage_container,
        "--kubernetes-service-load-balancer-ip-range", lb_range,
        "--ssh-public-key-file", SSH_PUBLIC_KEY,
        "--insecure",
        "--dry-run", "-o", "yaml",
    ]
    sh(cmd, stdout_path=CLUSTER_MANIFEST)
    docs = [d for d in yaml.safe_load_all(open(CLUSTER_MANIFEST)) if d]
    if not any(d.get("kind") for d in docs):
        die(f"Generated manifest {CLUSTER_MANIFEST} has no resources -- check the "
            "nkp preflight errors above (PE cluster / storage / image / creds).")
    log(f"Manifest written to {CLUSTER_MANIFEST} ({len(docs)} resources).", level="+")


def strip_cni_addon() -> None:
    """Remove addons.cni from the NKPCluster so NKP installs NO CNI (BYO-CNI).
    Nodes will come up NotReady until Flow CNI is applied in step 6."""
    log("Removing the managed CNI addon from the manifest (BYO-CNI)...")
    docs = list(yaml.safe_load_all(open(CLUSTER_MANIFEST)))
    removed = 0

    def walk(o):
        nonlocal removed
        if isinstance(o, dict):
            addons = o.get("addons")
            if isinstance(addons, dict) and "cni" in addons:
                addons.pop("cni", None)
                removed += 1
            for v in o.values():
                walk(v)
        elif isinstance(o, list):
            for v in o:
                walk(v)

    for d in docs:
        walk(d)
    with open(CLUSTER_MANIFEST, "w") as fh:
        yaml.dump_all(docs, fh, default_flow_style=False)
    if removed == 0:
        log("No addons.cni block found to remove (schema may differ) -- continuing.",
            level="-")
    else:
        log(f"Removed addons.cni ({removed} block(s)); NKP will install no CNI.", level="+")


def _timeout_seconds(spec: str) -> int:
    """Parse a kubectl-style duration ("30m", "600s", "1h") into seconds."""
    spec = spec.strip()
    try:
        if spec.endswith("h"):
            return int(float(spec[:-1]) * 3600)
        if spec.endswith("m"):
            return int(float(spec[:-1]) * 60)
        if spec.endswith("s"):
            return int(float(spec[:-1]))
        return int(spec)
    except ValueError:
        return 1800


def wait_for_workload_api(ns: str) -> None:
    """Wait until the workload API server answers -- WITHOUT requiring a Ready
    node. kube-vip serves the control-plane VIP before any CNI is installed, so
    the API is reachable while nodes are still NotReady. In the BYO-CNI flow we
    must gate on this (not ControlPlaneReady): ControlPlaneReady needs a Ready
    node, which needs a CNI, which we install in the NEXT step -- waiting on it
    here would deadlock (cluster never ready => Flow CNI never installed)."""
    log(f"Waiting up to {CONTROL_PLANE_TIMEOUT} for the workload API server to "
        "respond (CNI not required)...")
    deadline = time.time() + _timeout_seconds(CONTROL_PLANE_TIMEOUT)
    while time.time() < deadline:
        secret = kubectl(["-n", ns, "get", f"secret/{CLUSTER_NAME}-kubeconfig"],
                         check=False, capture=True, quiet=True)
        if secret.returncode == 0:
            fetch_workload_kubeconfig(ns)
            probe = kubectl(["get", "--raw=/readyz"], workload=True,
                            check=False, capture=True, quiet=True)
            if probe.returncode == 0:
                log("Workload API server is reachable.", level="+")
                return
        time.sleep(10)
    die(f"Workload API server for '{CLUSTER_NAME}' did not respond within "
        f"{CONTROL_PLANE_TIMEOUT}.")


def apply_cluster_and_wait() -> str:
    log("Applying the NKPCluster manifest to the management cluster...")
    kubectl(["apply", "-f", CLUSTER_MANIFEST])
    ns = detect_capi_namespace()
    if FLOW_CNI_ENABLED:
        # BYO-CNI: the CNI is stripped, so the control plane will never report
        # ControlPlaneReady until Flow CNI is installed (next step). Gate only on
        # API reachability here to avoid deadlocking; install_flow_cni() then waits
        # for nodes to go Ready.
        wait_for_workload_api(ns)
    else:
        log(f"Waiting up to {CONTROL_PLANE_TIMEOUT} for ControlPlaneReady in ns '{ns}'...")
        kubectl(["-n", ns, "wait", "--for=condition=ControlPlaneReady",
                 f"cluster/{CLUSTER_NAME}", f"--timeout={CONTROL_PLANE_TIMEOUT}"])
    return ns


def detect_capi_namespace() -> str:
    if CAPI_NAMESPACE_CFG:
        return CAPI_NAMESPACE_CFG
    # Poll briefly: the CAPI Cluster object appears a few seconds after apply.
    for _ in range(30):
        proc = kubectl(
            ["get", "clusters.cluster.x-k8s.io", "-A",
             "-o", "jsonpath={range .items[*]}{.metadata.namespace}/{.metadata.name}{'\\n'}{end}"],
            check=False, capture=True, quiet=True,
        )
        for line in (proc.stdout or "").splitlines():
            if line.strip().endswith(f"/{CLUSTER_NAME}"):
                ns = line.split("/", 1)[0]
                log(f"CAPI Cluster '{CLUSTER_NAME}' is in namespace '{ns}'.", level="+")
                return ns
        time.sleep(2)
    die(f"Could not find CAPI Cluster '{CLUSTER_NAME}' in any namespace.")
    return ""  # unreachable


def fetch_workload_kubeconfig(ns: str) -> None:
    log("Fetching the workload cluster kubeconfig...")
    tmp = WORKLOAD_KUBECONFIG + ".tmp"
    sh(["nkp", "get", "kubeconfig", "-c", CLUSTER_NAME, "-n", ns], stdout_path=tmp)
    os.replace(tmp, WORKLOAD_KUBECONFIG)
    kubectl(["get", "nodes"], workload=True, check=False)


def cluster_cidrs(ns: str) -> tuple[str, str]:
    """Return (podCIDR, serviceCIDR) of the created cluster for the HCP values."""
    proc = kubectl(
        ["get", "cluster", CLUSTER_NAME, "-n", ns, "-o",
         "jsonpath={.spec.clusterNetwork.pods.cidrBlocks[0]} "
         "{.spec.clusterNetwork.services.cidrBlocks[0]}"],
        capture=True,
    )
    parts = (proc.stdout or "").split()
    pod = parts[0] if len(parts) > 0 and parts[0] else "192.168.0.0/16"
    svc = parts[1] if len(parts) > 1 and parts[1] else "10.96.0.0/12"
    log(f"Cluster CIDRs: pods={pod} services={svc}", level="+")
    return pod, svc


# =============================================================================
# STEP 5: Pre-apply the Flow subchart CRDs on the workload cluster
# =============================================================================


def preapply_crds() -> None:
    """Pre-apply the Flow chart's CRDs. Helm installs crds/ only from the TOP
    chart, never from subcharts, so the OVN + container-security subchart CRDs
    must be applied manually here.

    Namespaces and the cni-secret are intentionally NOT pre-created: the chart
    OWNS them (parent templates the flow-cni-system Namespace + its secret; the
    OVN subchart creates the ovn-kubernetes secret and tolerates a pre-existing
    namespace) via dockerConfigSecret.create=true, and multus pulls from public
    ghcr.io. Pre-creating any of them only triggers Helm 'cannot be imported ...
    invalid ownership metadata' errors."""
    chart_dir = "/tmp/flowcni"
    log("Pulling the Flow chart and pre-applying its (subchart) CRDs...")
    shutil.rmtree(chart_dir, ignore_errors=True)
    sh(["helm", "pull", f"{FLOW_CHART_REPO}/nutanix-flow-cni",
        "--version", FLOW_CHART_VERSION, "--untar", "--untardir", chart_dir])
    base = os.path.join(chart_dir, "nutanix-flow-cni")
    crd_dirs = [os.path.join(base, "crds")]
    for sub in ("nutanix-core-flow-ovn-kubernetes", "nutanix-core-flow-container-security"):
        crd_dirs.append(os.path.join(base, "charts", sub, "crds"))
    for d in crd_dirs:
        if os.path.isdir(d):
            kubectl(["apply", "--server-side", "--force-conflicts", "-f", d],
                    workload=True, check=False)


# =============================================================================
# STEP 6: Install Flow CNI directly with `helm upgrade --install`, wait Ready
# =============================================================================


def write_flow_values(vip: str, pod_cidr: str, svc_cidr: str) -> None:
    """Write the Helm values for the Flow CNI chart (installed directly, not via
    a caaph HelmChartProxy)."""
    values = f"""crdUpgrade:
  enabled: false
image:
  repository: {FLOW_K8S_CNI_REPO}
  tag: "{FLOW_K8S_CNI_TAG}"
nutanix-core-flow-ovn-kubernetes:
  crdUpgrade:
    enabled: false
  k8sAPIServer: "https://{vip}:6443"
  podNetwork: "{pod_cidr}"
  serviceNetwork: "{svc_cidr}"
nutanix-core-flow-container-security:
  flowCns:
    enabled: false
  flowProcessor:
    enabled: false
global:
  enableEgressIp: true
  enableEgressService: true
  image:
    repository: {FLOW_OVN_REPO}
    tag: "{FLOW_OVN_TAG}"
  dockerConfigSecret:
    registry: {FLOW_HARBOR}
    auth: "{FLOW_DOCKER_AUTH_B64}"
    # The chart creates + owns 'cni-secret' in the namespaces it manages
    # (flow-cni-system via the parent, ovn-kubernetes via the OVN subchart).
    create: true
  imagePullSecretName: "cni-secret"
"""
    with open(FLOW_VALUES, "w") as fh:
        fh.write(values)
    log(f"Flow CNI values written to {FLOW_VALUES}.")


def _adopt_flow_namespaces() -> None:
    """Stamp Helm ownership metadata onto the Flow namespaces so `helm install`
    ADOPTS them instead of failing with 'exists and cannot be imported ...'.

    The chart templates its own release namespace (flow-cni-system) plus the
    OVN/container-security namespaces. Helm refuses to take over a namespace it
    did not create unless it carries the managed-by=Helm label + release
    annotations. We create flow-cni-system (the release namespace must exist) and
    stamp every Flow namespace that is already present (e.g. leftovers from a
    previous attempt) for our stable release name."""
    kubectl(["create", "namespace", "flow-cni-system"], workload=True,
            check=False, quiet=True)
    for ns in _FLOW_NAMESPACES:
        exists = kubectl(["get", "namespace", ns], workload=True,
                         check=False, capture=True, quiet=True)
        if exists.returncode != 0:
            continue
        kubectl(["label", "namespace", ns, "app.kubernetes.io/managed-by=Helm",
                 "--overwrite"], workload=True, check=False, quiet=True)
        kubectl(["annotate", "namespace", ns,
                 f"meta.helm.sh/release-name={FLOW_RELEASE}",
                 "meta.helm.sh/release-namespace=flow-cni-system",
                 "--overwrite"], workload=True, check=False, quiet=True)


def install_flow_cni(ns: str) -> None:
    """Install the Flow CNI chart directly with Helm on the workload cluster.

    We deliberately do NOT use a caaph HelmChartProxy: caaph pre-creates the
    release namespace as a plain (unowned) namespace, which permanently collides
    with the flow-cni-system Namespace the chart templates, and caaph renames the
    Helm release on every retry so it never recovers. A direct `helm upgrade
    --install` with a fixed release name + pre-stamped namespaces is idempotent
    and deterministic."""
    _adopt_flow_namespaces()

    log("Installing the Flow CNI chart directly with Helm (--wait --atomic)...")
    # --wait blocks until every workload (incl. the OVN daemonsets that flip the
    # nodes Ready) is up; --atomic rolls a failed attempt fully back so a re-run
    # starts clean. No --create-namespace: the chart owns flow-cni-system and we
    # already pre-stamped it for adoption above.
    sh(["helm", "upgrade", "--install", FLOW_RELEASE,
        f"{FLOW_CHART_REPO}/nutanix-flow-cni", "--version", FLOW_CHART_VERSION,
        "--namespace", "flow-cni-system", "--values", FLOW_VALUES,
        "--wait", "--atomic", "--timeout", "20m",
        "--kubeconfig", WORKLOAD_KUBECONFIG], check=False)

    # Confirm the real success signal: all nodes Ready once OVN is up.
    for attempt in range(1, FLOW_READY_RETRIES + 1):
        nodes = kubectl(["get", "nodes", "--no-headers"], workload=True,
                        check=False, capture=True, quiet=True)
        lines = [l for l in (nodes.stdout or "").splitlines() if l.strip()]
        ready_count = sum(1 for l in lines if _node_ready(l))
        if lines and ready_count == len(lines):
            log(f"All {len(lines)} nodes are Ready -- Flow CNI is up.", level="+")
            return
        log(f"Flow CNI progressing: {ready_count}/{len(lines) or '?'} nodes Ready "
            f"(attempt {attempt}/{FLOW_READY_RETRIES})...")
        time.sleep(FLOW_READY_DELAY_S)

    status = sh(["helm", "status", FLOW_RELEASE, "--namespace", "flow-cni-system",
                 "--kubeconfig", WORKLOAD_KUBECONFIG],
                check=False, capture=True, quiet=True)
    log("Nodes did not all become Ready in time. Inspect the Flow release + OVN "
        f"pods on the workload cluster:\n{status.stdout or status.stderr}", level="-")


def _node_ready(line: str) -> bool:
    parts = line.split()
    return len(parts) >= 2 and parts[1] == "Ready"


# =============================================================================
# STEP 7: registration extId + Flow activation
# =============================================================================


def _registration_records(payload: object) -> list:
    if isinstance(payload, list):
        return [c for c in payload if isinstance(c, dict)]
    if isinstance(payload, dict):
        for key in ("data", "clusters"):
            items = payload.get(key)
            if isinstance(items, list):
                return [c for c in items if isinstance(c, dict)]
    return []


def _record_ext_id(cluster: dict) -> str | None:
    for key in ("extId", "uuid", "KubernetesClusterUUID", "kubernetesClusterUuid"):
        val = cluster.get(key)
        if val:
            return val
    return None


def get_registered_cluster_extid(retries: int | None = None, *, quiet: bool = False) -> str | None:
    retries = retries or REGISTRATION_RETRIES
    if not quiet:
        log(f"Querying Prism Central for registered cluster '{CLUSTER_NAME}'...")
    endpoints = [
        f"{V4_BASE}/api/nke/v4.0.b1/config/cluster-registrations",
        f"{V4_BASE}/api/nke/v4.0.a1/config/cluster-registrations",
        f"{V4_BASE}/karbon/v1-alpha.1/k8s/cluster-registrations",
    ]
    headers = {"Accept": "application/json"}
    for attempt in range(1, retries + 1):
        for endpoint in endpoints:
            try:
                resp = requests.get(endpoint, auth=(PC_USERNAME, PC_PASSWORD),
                                    headers=headers, verify=False, timeout=10)
                if resp.status_code == 404:
                    continue
                resp.raise_for_status()
                for cluster in _registration_records(resp.json()):
                    if cluster.get("name") == CLUSTER_NAME:
                        ext_id = _record_ext_id(cluster)
                        log(f"Found '{CLUSTER_NAME}' -> extId={ext_id} "
                            f"status={cluster.get('status', 'UNKNOWN')}", level="+")
                        return ext_id
            except requests.exceptions.RequestException as e:
                if not quiet:
                    log(f"API request error on {endpoint}: {e}", level="-")
        if attempt < retries:
            log(f"Not registered yet. Waiting {REGISTRATION_DELAY_S}s... "
                f"(attempt {attempt}/{retries})")
            time.sleep(REGISTRATION_DELAY_S)
    if not quiet:
        log("Timed out waiting for the cluster to register with Prism Central.", level="-")
    return None


def flow_activated_ext_ids(session: requests.Session) -> set:
    """Cluster UUIDs known to the Flow subsystem (flow_kube_cluster_config). The
    ONLY extIds valid for a VPC kubernetesClusters reference. Empty on error."""
    body = {
        "entity_type": "flow_kube_cluster_config",
        "group_member_count": 100,
        "group_member_attributes": [{"attribute": "cluster_uuid"}, {"attribute": "name"}],
    }
    try:
        resp = session.post(f"{V4_BASE}/api/nutanix/v3/groups", json=body, timeout=30)
        resp.raise_for_status()
    except requests.exceptions.RequestException as e:
        log(f"flow_kube_cluster_config groups query failed: {e}", level="-")
        return set()
    out = set()
    for group in resp.json().get("group_results", []) or []:
        for entity in group.get("entity_results", []) or []:
            uuid = ""
            for field in entity.get("data", []) or []:
                if field.get("name") == "cluster_uuid":
                    vals = (field.get("values") or [{}])[0].get("values") or []
                    if vals:
                        uuid = vals[0]
            uuid = uuid or entity.get("entity_id", "")
            if uuid:
                out.add(uuid)
    return out


def activate_flow_capabilities(session: requests.Session, ext_id: str) -> bool:
    """Flow-activate the cluster by uploading its kubeconfig to PC. REQUIRED
    before a VPC can reference it. Flow CNI must already be installed (step 6)."""
    if ext_id in flow_activated_ext_ids(session):
        log(f"Cluster {ext_id} is already Flow-activated.", level="+")
        return True
    try:
        with open(WORKLOAD_KUBECONFIG) as fh:
            kubeconfig = fh.read()
    except OSError as e:
        log(f"Could not read {WORKLOAD_KUBECONFIG} ({e}); skipping Flow activation.", level="-")
        return False

    kubeconfig_b64 = base64.b64encode(kubeconfig.encode()).decode()
    url = f"{V4_BASE}/karbon/v1-alpha.1/k8s/cluster-registrations/{ext_id}"
    log("Activating Flow capabilities (PATCH kubeconfig to PC)...")
    try:
        resp = session.patch(url, json={"kubeconfig": kubeconfig_b64}, timeout=60)
        if resp.status_code >= 400:
            log(f"Flow activation PATCH -> HTTP {resp.status_code}: {resp.text[:300]}", level="-")
    except requests.exceptions.RequestException as e:
        log(f"Flow activation PATCH failed: {e}", level="-")
        return False

    for attempt in range(1, REGISTRATION_RETRIES + 1):
        if ext_id in flow_activated_ext_ids(session):
            log(f"Cluster {ext_id} is now Flow-activated (VPC-ready).", level="+")
            return True
        log(f"Waiting for Flow activation... (attempt {attempt}/{REGISTRATION_RETRIES})")
        time.sleep(REGISTRATION_DELAY_S)

    log("Cluster did not become Flow-activated in time. Ensure Flow CNI is Ready, "
        "then activate it (Kubernetes Clusters -> Actions -> Activate Nutanix Flow "
        "Capabilities).", level="-")
    return False


# =============================================================================
# MAIN
# =============================================================================


def run_workflow() -> int:
    print("=" * 64)
    print("   Nutanix NKP + Flow CNI End-to-End Deployment (on bastion)")
    print("=" * 64)

    session = _v4_session()

    # Resume guard: if CLUSTER_NAME is already registered, deploy already ran
    # (nkp create is not idempotent). Skip straight to Flow activation.
    log("Checking whether the cluster is already registered (resume guard)...")
    ext_id = get_registered_cluster_extid(retries=1, quiet=True)

    if ext_id:
        log(f"Cluster '{CLUSTER_NAME}' is already registered (extId={ext_id}); "
            "skipping bootstrap + deploy.", level="+")
    else:
        preflight_tools()
        ensure_ssh_key()
        create_bootstrap()
        vip, lb_range = resolve_nkp_ips(session)
        pe_cluster = discover_pe_cluster(session)
        storage_container = discover_storage_container(session)
        generate_cluster_manifest(vip, lb_range, pe_cluster, storage_container)
        if FLOW_CNI_ENABLED:
            strip_cni_addon()
        ns = apply_cluster_and_wait()
        fetch_workload_kubeconfig(ns)

        if FLOW_CNI_ENABLED:
            pod_cidr, svc_cidr = cluster_cidrs(ns)
            preapply_crds()
            write_flow_values(vip, pod_cidr, svc_cidr)
            install_flow_cni(ns)
        else:
            log("flow_cni.enabled is false; leaving NKP's default CNI in place.")

        log(f"Waiting {POST_CREATE_GRACE_S}s for the Konnector agent to register the cluster...")
        time.sleep(POST_CREATE_GRACE_S)
        ext_id = get_registered_cluster_extid()

    if not ext_id:
        log("Workflow finished, but cluster registration could not be verified.", level="-")
        return 1

    vpc_ready = activate_flow_capabilities(session, ext_id) if FLOW_CNI_ENABLED else False

    print("\n" + "=" * 64)
    print(" SUCCESS!")
    print("=" * 64)
    if vpc_ready:
        print(" Cluster is Flow-activated and ready for VPC association.")
    elif FLOW_CNI_ENABLED:
        print(" NOTE: cluster is registered but NOT yet Flow-activated; a VPC")
        print(" kubernetesClusters reference will fail until Flow activation succeeds")
        print(" (needs Flow CNI Ready on the cluster).")
    else:
        print(" NOTE: Flow CNI was not deployed (flow_cni.enabled=false); this")
        print(" cluster cannot be referenced by a VPC until it runs Flow + is activated.")
    print("\n testenv/fill_test_config.py will discover this extId for")
    print(" networking.kubernetes_cluster_ext_id once the cluster is Flow-activated.\n")
    print(" VPC v4 API kubernetesClusters payload:\n")
    print(json.dumps({"kubernetesClusters": [{"extId": ext_id}]}, indent=2))
    return 0


def main() -> int:
    """Dispatch: run the workflow when already on the bastion (or when --local is
    passed), otherwise SSH into the bastion and run the whole workflow there."""
    if ON_BASTION or "--local" in sys.argv:
        return run_workflow()
    return run_on_bastion()


if __name__ == "__main__":
    try:
        sys.exit(main())
    except KeyboardInterrupt:
        die("Interrupted by user.")
