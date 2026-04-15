#!/usr/bin/env bash
# Run acceptance tests locally (same as /ok-to-test on GitHub).
# Must be run from the repository root. Requires NUTANIX_* env vars and test config files.
#
# Usage:
#   ./scripts/run-acceptance-test.sh [-p package] <test_pattern>
#
# Examples (equivalent to /ok-to-test on GitHub):
#   ./scripts/run-acceptance-test.sh -p vmmv2 TestAccV2NutanixOvaVmDeployResource_DeployVMFromOva
#   ./scripts/run-acceptance-test.sh -p vmmv2 TestAccV2NutanixNGTInsertIsoResource_InsertNGTIsoIntoVmHaveNGTTest
#   ./scripts/run-acceptance-test.sh v4                                    # all TestAccV2Nutanix* (runs from repo root)
#   ./scripts/run-acceptance-test.sh TestAccV2NutanixOvaVmDeployResource_DeployVMFromOva  # single test (needs -p for vmmv2)
#
# Options:
#   -p package   Package under nutanix/services/ (e.g. vmmv2). Required for vmmv2 so test_config_v2.json is found.

set -e

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

PACKAGE=""
RUN_PATTERN=""

while getopts "p:" opt; do
  case "$opt" in
    p) PACKAGE="$OPTARG" ;;
    *) exit 1 ;;
  esac
done
shift "$((OPTIND - 1))"
RUN_PATTERN="${1:?Usage: $0 [-p package] <test_pattern>}"

# Map provider names to -run patterns (same as GitHub workflow)
case "$RUN_PATTERN" in
  foundation)           RUN_PATTERN="TestAccFoundation*" ;;
  foundation_central)   RUN_PATTERN="TestAccFC*" ;;
  karbon)               RUN_PATTERN="TestAccKarbon*" ;;
  v3)                   RUN_PATTERN="TestAccNutanix*" ;;
  v4)                   RUN_PATTERN="TestAccV2Nutanix*" ;;
  lcm)                  RUN_PATTERN="TestAccV2NutanixLcm*" ;;
  era)                  RUN_PATTERN="TestAccEra*" ;;
esac

export TF_ACC=1
export TF_LOG="${TF_LOG:-ERROR}"
export GOTRACEBACK=all

# Pre-check: NUTANIX_* required by TestAccPreCheck
if [[ -z "${NUTANIX_USERNAME:-}" || -z "${NUTANIX_PASSWORD:-}" || -z "${NUTANIX_ENDPOINT:-}" ]]; then
  echo "Error: NUTANIX_USERNAME, NUTANIX_PASSWORD, NUTANIX_ENDPOINT (and NUTANIX_INSECURE, NUTANIX_PORT, NUTANIX_STORAGE_CONTAINER) must be set for acceptance tests."
  echo "Copy from your GitHub Actions secrets or set in .env and source it."
  exit 1
fi

if [[ -n "$PACKAGE" ]]; then
  # Run from package dir so paths like ../../../test_config_v2.json resolve to repo root
  PKG_DIR="$REPO_ROOT/nutanix/services/$PACKAGE"
  if [[ ! -d "$PKG_DIR" ]]; then
    echo "Error: Package directory not found: $PKG_DIR"
    exit 1
  fi
  if [[ "$PACKAGE" == "vmmv2" && ! -f "$REPO_ROOT/test_config_v2.json" ]]; then
    echo "Warning: test_config_v2.json not found at repo root; vmmv2 tests may fail."
  fi
  echo "==> Running acceptance tests in package: $PACKAGE (pattern: $RUN_PATTERN)"
  (cd "$PKG_DIR" && go test . -v -run="$RUN_PATTERN" -timeout 500m -count=1)
else
  # Run all packages matching the pattern from repo root
  echo "==> Running acceptance tests in ./... (pattern: $RUN_PATTERN)"
  go test ./... -v -run="$RUN_PATTERN" -timeout 500m -count=1
fi
