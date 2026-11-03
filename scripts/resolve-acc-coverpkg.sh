#!/usr/bin/env bash
# Resolve the -coverpkg denominator for an acceptance run.
#
# The percentage is executed statements divided by statements in this set.
# A named test is divided by the service package that contains it, not by the
# whole module and not by a single resource file.
#
# Usage:
#   scripts/resolve-acc-coverpkg.sh <package_path> <cover_scope> <run_pattern>
# stdout: two lines — coverpkg, scope_label
set -euo pipefail

package_path="${1:-./...}"
cover_scope="${2:-}"
run_pattern="${3:-}"

if [[ "$package_path" != "./..." ]]; then
  printf '%s\n' "$package_path"
  printf '%s\n' "package $package_path"
  exit 0
fi

join_dirs() {
  local coverpkg="" d
  for d in "$@"; do
    [[ -d "$d" ]] || continue
    coverpkg="${coverpkg:+$coverpkg,}./$d"
  done
  printf '%s' "$coverpkg"
}

case "$cover_scope" in
  v4)
    dirs=()
    for d in nutanix/services/*v2; do
      [[ -d "$d" ]] || continue
      dirs+=("$d")
    done
    if ((${#dirs[@]})); then
      printf '%s\n' "$(join_dirs "${dirs[@]}")"
    else
      printf '\n'
    fi
    printf '%s\n' "overall v4 v2 service packages"
    exit 0
    ;;
  v3)
    dirs=()
    for d in nutanix/services/*; do
      [[ -d "$d" ]] || continue
      case "$d" in *v2) continue ;; esac
      dirs+=("$d")
    done
    if ((${#dirs[@]})); then
      printf '%s\n' "$(join_dirs "${dirs[@]}")"
    else
      printf '\n'
    fi
    printf '%s\n' "overall v3 non-v2 service packages"
    exit 0
    ;;
  lcm)
    printf '%s\n' "./nutanix/services/lcmv2"
    printf '%s\n' "package ./nutanix/services/lcmv2"
    exit 0
    ;;
  era)
    printf '%s\n' "./nutanix/services/ndb"
    printf '%s\n' "package ./nutanix/services/ndb"
    exit 0
    ;;
  foundation)
    printf '%s\n' "./nutanix/services/foundation"
    printf '%s\n' "package ./nutanix/services/foundation"
    exit 0
    ;;
  foundation_central)
    printf '%s\n' "./nutanix/services/foundationCentral"
    printf '%s\n' "package ./nutanix/services/foundationCentral"
    exit 0
    ;;
  karbon)
    printf '%s\n' "./nutanix/services/nke"
    printf '%s\n' "package ./nutanix/services/nke"
    exit 0
    ;;
esac

# Plain test names (no wildcard): denominator is each package that declares them.
if [[ -n "$run_pattern" && "$run_pattern" != "." && "$run_pattern" != ".*" && "$run_pattern" != "*" && "$run_pattern" != *\** ]]; then
  coverpkg=""
  IFS='|' read -ra names <<< "$run_pattern"
  for name in "${names[@]}"; do
    [[ -z "$name" ]] && continue
    while IFS= read -r file; do
      [[ -n "$file" ]] || continue
      rel="./$(dirname "$file")"
      case ",$coverpkg," in
        *",$rel,"*) ;;
        *) coverpkg="${coverpkg:+$coverpkg,}$rel" ;;
      esac
    done < <(grep -l -R --include='*_test.go' -F "func ${name}" nutanix/services 2>/dev/null || true)
  done
  if [[ -n "$coverpkg" ]]; then
    printf '%s\n' "$coverpkg"
    printf '%s\n' "package(s) containing [$run_pattern]"
    exit 0
  fi
fi

printf '%s\n' "./..."
printf '%s\n' "all packages"
