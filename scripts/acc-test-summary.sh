#!/usr/bin/env bash
# Appends test summary (last result per test name) to the given log file.
# Usage: scripts/acc-test-summary.sh <test_output.log>

set -e

LOGFILE="${1:?Usage: $0 <test_output.log>}"

if [[ ! -f "$LOGFILE" ]] || ! grep -qE '^--- (PASS|FAIL|SKIP):' "$LOGFILE" 2>/dev/null; then
  exit 0
fi

TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT

awk '
  /^--- (PASS|FAIL|SKIP): / {
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
' "$LOGFILE" | sort -n > "$TMP"

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
} >> "$LOGFILE"

echo ""
echo "==> Test summary appended to $LOGFILE"
