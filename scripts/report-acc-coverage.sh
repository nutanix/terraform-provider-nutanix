#!/usr/bin/env bash
# Report statement coverage from an Acc cover profile and optionally append to a log.
# Usage:
#   scripts/report-acc-coverage.sh <profile> [logfile] [scope_label]
#
# The total is go tool cover's statement percentage of the profile as written.
# That profile must be -covermode=atomic with -coverpkg set to the intended
# denominator (see scripts/resolve-acc-coverpkg.sh). This script does not drop
# files: dropping them changes the denominator and inflates the percentage.
#
# Soft-fails on missing profile (prints a warning). Hard-fails only on invalid parse
# when the profile exists but cannot be read.
# When GITHUB_ENV is set (Actions), also writes CODE_COVERAGE_OUTPUT for PR comments.
#
# Also writes a local HTML/text report under COVERAGE_REPORT_DIR (default: coverage-report/):
#   coverage-report/coverage.html
#   coverage-report/coverage.txt
#   coverage-report/c.out

set -euo pipefail

PROFILE="${1:?Usage: $0 <profile> [logfile] [scope_label]}"
LOGFILE="${2:-}"
SCOPE_LABEL="${3:-acceptance tests}"
THRESHOLD="${TESTCOV_THRESHOLD:-85}"
REPORT_DIR="${COVERAGE_REPORT_DIR:-coverage-report}"

report() {
  local line="$1"
  echo "$line"
  if [[ -n "$LOGFILE" ]]; then
    echo "$line" >> "$LOGFILE"
  fi
}

set_github_env() {
  local value="$1"
  if [[ -n "${GITHUB_ENV:-}" ]]; then
    echo "CODE_COVERAGE_OUTPUT=${value}" >> "$GITHUB_ENV"
  fi
}

write_html_report() {
  local profile="$1"
  mkdir -p "$REPORT_DIR"
  go tool cover -html="$profile" -o "$REPORT_DIR/coverage.html"
  go tool cover -func="$profile" > "$REPORT_DIR/coverage.txt"
  cp "$profile" "$REPORT_DIR/c.out"
  report "HTML report: $REPORT_DIR/coverage.html"
  report "Func report: $REPORT_DIR/coverage.txt"
  report "Open with: open $REPORT_DIR/coverage.html"
}

if [[ ! -s "$PROFILE" ]]; then
  report "==> Coverage: profile missing or empty ($PROFILE); skipped"
  set_github_env "Coverage profile missing or empty; could not compute coverage."
  exit 0
fi

total_line=$(
  go tool cover -func="$PROFILE" |
    awk '$1 == "total:" { print; exit }'
)

if [[ -z "${total_line}" ]]; then
  report "==> Coverage: ERROR could not find total line in $PROFILE"
  set_github_env "Could not parse total coverage from profile."
  exit 1
fi

coverage=$(
  echo "$total_line" |
    awk '{ gsub("%", "", $NF); print $NF }'
)

if ! [[ "$coverage" =~ ^[0-9]+([.][0-9]+)?$ ]]; then
  report "==> Coverage: ERROR invalid coverage value: ${coverage:-<empty>}"
  set_github_env "Invalid coverage value parsed from profile."
  exit 1
fi

report ""
report "================================================== 📊 CODE COVERAGE 📊 ================================================================================"
report "Scope: $SCOPE_LABEL"
report "Profile: $PROFILE"
report "Result: $total_line"
report "Statement coverage: ${coverage}%"
report "Advisory threshold: ${THRESHOLD}%"
if awk -v actual="$coverage" -v threshold="$THRESHOLD" 'BEGIN { exit !(actual+0 >= threshold+0) }'; then
  report "Status: meets advisory threshold"
  set_github_env "Statement coverage is ${coverage}% (scope: ${SCOPE_LABEL}; threshold ${THRESHOLD}%, advisory)."
else
  report "Status: below advisory threshold (not failing)"
  set_github_env "Statement coverage ${coverage}% is below advisory threshold (${THRESHOLD}%). Scope: ${SCOPE_LABEL}."
fi
write_html_report "$PROFILE"
report "================================================================================================================================================"
