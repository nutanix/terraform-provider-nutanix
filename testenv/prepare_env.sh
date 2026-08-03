#!/usr/bin/env bash
#
# prepare_env.sh -- orchestrate the three PC-preparation stages, in order:
#
#   1. pc      python3 testenv/prepare_pc.py       (enable PC services + prereq VMs)
#   2. tf      terraform apply in testenv/terraform (provision test infra resources)
#   3. config  python3 testenv/fill_test_config.py (generate test_config_v2.json + IAM)
#
# Ordering matters: stage 1 enables the Network Controller (Flow) that stage 2's
# VPC/overlay subnet require, and enables the Object Store service used by the
# object-store tests. Stage 3 is independent of stage 2 but is run last so the
# live cluster discovery it does reflects the fully provisioned environment.
#
# All stages read testenv/config.yaml (terraform additionally reads
# testenv/terraform/terraform.tfvars).
#
# Usage
# -----
#   testenv/prepare_env.sh                 # run all three stages, in order
#   testenv/prepare_env.sh --dry-run       # preview only (pc/config --dry-run, tf plan)
#   testenv/prepare_env.sh --only pc       # run a single stage (repeatable)
#   testenv/prepare_env.sh --skip tf       # skip a stage (repeatable)
#   testenv/prepare_env.sh --yes           # auto-approve `terraform apply`
#   testenv/prepare_env.sh --pc-args "--skip objects --skip dp"   # pass args to prepare_pc.py
#   testenv/prepare_env.sh --tf-args "-target=..."                # pass args to terraform apply
#   testenv/prepare_env.sh --config-args "--skip-iam"             # pass args to fill_test_config.py
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
TF_DIR="${SCRIPT_DIR}/terraform"

ALL_STAGES=(pc tf config)

DRY_RUN=false
AUTO_APPROVE=false
ONLY_STAGES=()
SKIP_STAGES=()
PC_ARGS=""
TF_ARGS=""
CONFIG_ARGS=""

die() { echo "ERROR: $*" >&2; exit 1; }
log() { echo; echo "==> $*"; }

usage() { awk 'NR>1 && /^#/{sub(/^# ?/,""); print; next} NR>1{exit}' "${BASH_SOURCE[0]}"; exit "${1:-0}"; }

is_valid_stage() {
  local s="$1"
  for v in "${ALL_STAGES[@]}"; do [[ "$s" == "$v" ]] && return 0; done
  return 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)      DRY_RUN=true; shift ;;
    --yes|-y)       AUTO_APPROVE=true; shift ;;
    --only)         is_valid_stage "${2:-}" || die "--only expects one of: ${ALL_STAGES[*]}"; ONLY_STAGES+=("$2"); shift 2 ;;
    --skip)         is_valid_stage "${2:-}" || die "--skip expects one of: ${ALL_STAGES[*]}"; SKIP_STAGES+=("$2"); shift 2 ;;
    --pc-args)      PC_ARGS="${2:-}"; shift 2 ;;
    --tf-args)      TF_ARGS="${2:-}"; shift 2 ;;
    --config-args)  CONFIG_ARGS="${2:-}"; shift 2 ;;
    -h|--help)      usage 0 ;;
    *)              die "unknown argument: $1 (see --help)" ;;
  esac
done

# Resolve the set of stages to run (preserve canonical order).
stage_selected() {
  local s="$1"
  if [[ ${#ONLY_STAGES[@]} -gt 0 ]]; then
    for v in "${ONLY_STAGES[@]}"; do [[ "$s" == "$v" ]] && return 0; done
    return 1
  fi
  for v in "${SKIP_STAGES[@]}"; do [[ "$s" == "$v" ]] && return 1; done
  return 0
}

PY="${PYTHON:-python3}"

run_pc() {
  command -v "$PY" >/dev/null 2>&1 || die "python3 not found (set \$PYTHON to override)"
  [[ -f "${SCRIPT_DIR}/config.yaml" || -f "${SCRIPT_DIR}/.env" ]] \
    || die "neither testenv/config.yaml nor testenv/.env exists"
  local extra="$PC_ARGS"
  $DRY_RUN && extra="--dry-run $extra"
  log "Stage 1/3 [pc]: prepare_pc.py $extra"
  # shellcheck disable=SC2086
  "$PY" "${SCRIPT_DIR}/prepare_pc.py" $extra
}

run_tf() {
  command -v terraform >/dev/null 2>&1 || die "terraform not found in PATH"
  [[ -d "$TF_DIR" ]] || die "terraform dir not found: $TF_DIR"
  [[ -f "${TF_DIR}/terraform.tfvars" ]] \
    || die "missing ${TF_DIR}/terraform.tfvars (copy/fill it before provisioning)"
  log "Stage 2/3 [tf]: terraform init"
  terraform -chdir="$TF_DIR" init -input=false
  if $DRY_RUN; then
    log "Stage 2/3 [tf]: terraform plan $TF_ARGS"
    # shellcheck disable=SC2086
    terraform -chdir="$TF_DIR" plan $TF_ARGS
  else
    local approve=""
    $AUTO_APPROVE && approve="-auto-approve"
    log "Stage 2/3 [tf]: terraform apply $approve $TF_ARGS"
    # shellcheck disable=SC2086
    terraform -chdir="$TF_DIR" apply $approve $TF_ARGS
  fi
}

run_config() {
  command -v "$PY" >/dev/null 2>&1 || die "python3 not found (set \$PYTHON to override)"
  local extra="$CONFIG_ARGS"
  $DRY_RUN && extra="--dry-run $extra"
  log "Stage 3/3 [config]: fill_test_config.py $extra"
  # shellcheck disable=SC2086
  "$PY" "${SCRIPT_DIR}/fill_test_config.py" $extra
}

log "prepare_env: repo=${REPO_ROOT} dry_run=${DRY_RUN} auto_approve=${AUTO_APPROVE}"
for stage in "${ALL_STAGES[@]}"; do
  if stage_selected "$stage"; then
    case "$stage" in
      pc)     run_pc ;;
      tf)     run_tf ;;
      config) run_config ;;
    esac
  else
    log "Skipping stage [$stage]"
  fi
done

log "prepare_env: done."
