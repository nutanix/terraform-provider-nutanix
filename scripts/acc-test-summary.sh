#!/usr/bin/env bash
# Appends test summary (last result per test name) to the given log file.
# Usage: scripts/acc-test-summary.sh <test_output.log>

set -e

LOGFILE="${1:?Usage: $0 <test_output.log>}"

# grep -a treats the file as text even if DEBUG HTTP dumps contain binary (e.g. ISO uploads).
# awk on the raw file can fail with "multibyte conversion failure".
if [[ ! -f "$LOGFILE" ]] || ! grep -a -qE '^--- (PASS|FAIL|SKIP):' "$LOGFILE" 2>/dev/null; then
  exit 0
fi

TMP=$(mktemp)
SUMMARY_TMP=$(mktemp)
trap 'rm -f "$TMP" "$SUMMARY_TMP"' EXIT

grep -aE '^--- (PASS|FAIL|SKIP): ' "$LOGFILE" | LC_ALL=C awk '
  {
    result = $2; sub(/:$/, "", result);
    name = $3;
    last_result[name] = result;
    last_line[name] = NR;
  }
  END {
    for (n in last_result) {
      print last_line[n], last_result[n], n
    }
  }
' | sort -n > "$TMP"

TOTAL_PASSED=$(awk '$2 == "PASS"' "$TMP" | wc -l | tr -d ' ')
TOTAL_FAILED=$(awk '$2 == "FAIL"' "$TMP" | wc -l | tr -d ' ')
TOTAL_SKIPPED=$(awk '$2 == "SKIP"' "$TMP" | wc -l | tr -d ' ')
UNIQUE_TESTS=$((TOTAL_PASSED + TOTAL_FAILED + TOTAL_SKIPPED))

PASS_PERCENT=0
FAIL_PERCENT=0
SKIP_PERCENT=0
if [[ $UNIQUE_TESTS -gt 0 ]]; then
  [[ $TOTAL_PASSED -gt 0 ]] && PASS_PERCENT=$((TOTAL_PASSED * 100 / UNIQUE_TESTS))
  [[ $TOTAL_FAILED -gt 0 ]] && FAIL_PERCENT=$((TOTAL_FAILED * 100 / UNIQUE_TESTS))
  [[ $TOTAL_SKIPPED -gt 0 ]] && SKIP_PERCENT=$((TOTAL_SKIPPED * 100 / UNIQUE_TESTS))
fi

# Fetch PC / AOS versions via ncli when SSH credentials are available
# (same approach as nutanix.ansible ok-to-test-command).
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CONFIG_FILE="${NUTANIX_TEST_CONFIG:-$REPO_ROOT/test_config_v2.json}"
NCLI="/home/nutanix/prism/cli/ncli"
VERSION_CMD="$NCLI cluster info | grep -i Version | grep -vi NCC"
SSH_OPTS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=10"

PC_VERSION_INFO="Failed to fetch PC version"
AOS_VERSION_INFO="Failed to fetch AOS version"

if [[ -f "$CONFIG_FILE" ]] && command -v sshpass >/dev/null 2>&1 && command -v jq >/dev/null 2>&1; then
  PC_SSH_USER=$(jq -r '.ssh_pc_username // "nutanix"' "$CONFIG_FILE")
  PC_SSH_PASS=$(jq -r '.ssh_pc_password // empty' "$CONFIG_FILE")
  PE_SSH_USER=$(jq -r '.ssh_pe_username // "nutanix"' "$CONFIG_FILE")
  PE_SSH_PASS=$(jq -r '.ssh_pe_password // empty' "$CONFIG_FILE")
  PE_IP=$(jq -r '.data_protection.local_cluster_pe // empty' "$CONFIG_FILE")
  PC_IP="${NUTANIX_ENDPOINT:-}"

  if [[ -n "$PC_IP" && -n "$PC_SSH_PASS" ]]; then
    PC_VERSION_INFO=$(sshpass -p "$PC_SSH_PASS" \
      ssh $SSH_OPTS "$PC_SSH_USER@$PC_IP" \
      "$VERSION_CMD" 2>/dev/null || echo "Failed to fetch PC version")
  elif [[ -z "$PC_IP" ]]; then
    PC_VERSION_INFO="NUTANIX_ENDPOINT not set"
  fi

  if [[ -n "$PE_IP" && -n "$PE_SSH_PASS" ]]; then
    AOS_VERSION_INFO=$(sshpass -p "$PE_SSH_PASS" \
      ssh $SSH_OPTS "$PE_SSH_USER@$PE_IP" \
      "$VERSION_CMD" 2>/dev/null || echo "Failed to fetch AOS version")
  else
    AOS_VERSION_INFO="PE IP not configured"
  fi
elif [[ ! -f "$CONFIG_FILE" ]]; then
  PC_VERSION_INFO="test_config_v2.json not found"
  AOS_VERSION_INFO="test_config_v2.json not found"
elif ! command -v sshpass >/dev/null 2>&1; then
  PC_VERSION_INFO="sshpass not installed"
  AOS_VERSION_INFO="sshpass not installed"
fi

{
  echo ""
  echo "================================================== 🧪 TEST SUMMARY 🧪 ================================================================================="
  echo "Total Test Cases Run 🚀: $UNIQUE_TESTS"
  echo "Total Test Cases Passed ✅: $TOTAL_PASSED ($PASS_PERCENT %)"
  echo "Total Test Cases Failed ❌: $TOTAL_FAILED ($FAIL_PERCENT %)"
  echo "Total Test Cases Skipped ⚠️: $TOTAL_SKIPPED ($SKIP_PERCENT %)"
  echo "================================================================================================================================================"
  echo ""
  echo "================================================== TESTS SUCCEEDED ✅ ============================================================================="
  if [[ $TOTAL_PASSED -gt 0 ]]; then
    awk '$2 == "PASS" { print "✅ " $3 }' "$TMP"
  else
    echo "No tests passed 😞❗"
  fi
  echo "================================================================================================================================================"
  echo ""
  echo "================================================== TESTS FAILED ❌ ================================================================================"
  if [[ $TOTAL_FAILED -gt 0 ]]; then
    awk '$2 == "FAIL" { print "❌ " $3 }' "$TMP"
  else
    echo "🎉💃 No tests failed 🕺🎉"
  fi
  echo "================================================================================================================================================"
  echo ""
  echo "================================================== TESTS SKIPPED ⚠️ =============================================================================="
  if [[ $TOTAL_SKIPPED -gt 0 ]]; then
    awk '$2 == "SKIP" { print "⚠️ " $3 "   : Reason: ➡️ :  See log for details" }' "$TMP"
  else
    echo "🎉💃 No tests skipped 🕺🎉"
  fi
  echo "================================================================================================================================================"
  echo ""
  echo "PC Version:"
  echo "$PC_VERSION_INFO"
  echo ""
  echo "AOS Version:"
  echo "$AOS_VERSION_INFO"
} > "$SUMMARY_TMP"

cat "$SUMMARY_TMP" >> "$LOGFILE"
echo ""
echo "==> Test summary appended to $LOGFILE"
echo ""
cat "$SUMMARY_TMP"
