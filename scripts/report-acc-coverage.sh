#!/usr/bin/env bash
# Report statement coverage from an Acc cover profile and optionally append to a log.
# Usage:
#   scripts/report-acc-coverage.sh <profile> [logfile] [scope_label] [run_pattern] [package_path]
#
# When run_pattern names specific Acc tests (not '.' / wildcards), the profile is
# filtered to only the matching resource/datasource source files before reporting.
#
# Soft-fails on missing profile (prints a warning). Hard-fails only on invalid parse
# when the profile exists but cannot be read.
# When GITHUB_ENV is set (Actions), also writes CODE_COVERAGE_OUTPUT for PR comments.
#
# Also writes a local HTML/text report under COVERAGE_REPORT_DIR (default: coverage-report/):
#   coverage-report/coverage.html
#   coverage-report/coverage.txt
#   coverage-report/c.out  (profile used for the report; may be filtered)

set -euo pipefail

PROFILE="${1:?Usage: $0 <profile> [logfile] [scope_label] [run_pattern] [package_path]}"
LOGFILE="${2:-}"
SCOPE_LABEL="${3:-acceptance tests}"
RUN_PATTERN="${4:-}"
PACKAGE_PATH="${5:-}"
THRESHOLD="${TESTCOV_THRESHOLD:-85}"
REPORT_DIR="${COVERAGE_REPORT_DIR:-coverage-report}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

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

REPORT_PROFILE="$PROFILE"
FILTERED_PROFILE=""
if [[ -n "$RUN_PATTERN" ]]; then
  FILTERED_PROFILE="$(mktemp)"
  if bash "$SCRIPT_DIR/filter-acc-cover-by-tests.sh" "$PROFILE" "$FILTERED_PROFILE" "$RUN_PATTERN" "$PACKAGE_PATH"; then
    if [[ -s "$FILTERED_PROFILE" ]] && grep -q '^mode:' "$FILTERED_PROFILE"; then
      # Only use filtered profile if it still has block lines.
      if [[ "$(wc -l < "$FILTERED_PROFILE" | tr -d ' ')" -gt 1 ]]; then
        REPORT_PROFILE="$FILTERED_PROFILE"
        if [[ "$RUN_PATTERN" != "." && "$RUN_PATTERN" != *\** ]]; then
          SCOPE_LABEL="tests [$RUN_PATTERN] → matched resource/datasource files"
        fi
      fi
    fi
  fi
fi

cleanup() {
  [[ -n "$FILTERED_PROFILE" && -f "$FILTERED_PROFILE" ]] && rm -f "$FILTERED_PROFILE"
}
trap cleanup EXIT

total_line=$(
  go tool cover -func="$REPORT_PROFILE" |
    awk '$1 == "total:" { print; exit }'
)

if [[ -z "${total_line}" ]]; then
  report "==> Coverage: ERROR could not find total line in $REPORT_PROFILE"
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
report "Profile: $REPORT_PROFILE"
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
write_html_report "$REPORT_PROFILE"
report "================================================================================================================================================"
