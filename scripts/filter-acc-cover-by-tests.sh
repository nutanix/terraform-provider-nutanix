#!/usr/bin/env bash
# Filter an Acc cover profile to only resource/datasource files relevant to the
# -run test pattern (not the whole package). Used for single/few Acc test runs.
#
# Usage:
#   scripts/filter-acc-cover-by-tests.sh <in.out> <out.out> <run_pattern> [package_path]
#
# run_pattern examples:
#   TestAccV2NutanixVolumeGroupResource_Basic
#   TestAccV2NutanixVolumeGroupResource_Basic|TestAccV2NutanixVolumeGroupDiskResource_Basic
#   .                  → no filter (copy through)
#   TestAccV2Nutanix*  → no filter (broad suite)
#
# Matching: longest CamelCase stem from production *.go basenames that appears
# in the test name (e.g. volume_group_disk → VolumeGroupDisk).

set -euo pipefail

IN="${1:?Usage: $0 <in.out> <out.out> <run_pattern> [package_path]}"
OUT="${2:?}"
RUN_PATTERN="${3:?}"
PACKAGE_PATH="${4:-}"

# Broad patterns → keep full profile (package / v3 / v4 style runs).
if [[ "$RUN_PATTERN" == "." || "$RUN_PATTERN" == ".*" || "$RUN_PATTERN" == "*" ]]; then
  cp "$IN" "$OUT"
  echo "cover-filter: broad -run=.; keeping full profile"
  exit 0
fi
if [[ "$RUN_PATTERN" == *\** ]]; then
  # Wildcards like TestAccV2Nutanix* / TestAccNutanix*
  cp "$IN" "$OUT"
  echo "cover-filter: wildcard -run=$RUN_PATTERN; keeping full profile"
  exit 0
fi

python3 - "$IN" "$OUT" "$RUN_PATTERN" "$PACKAGE_PATH" <<'PY'
import re
import sys
from pathlib import Path

in_path, out_path, run_pattern, package_path = sys.argv[1:5]

def snake_to_camel(s: str) -> str:
    return "".join(p.title() for p in s.split("_") if p)

def file_stem_key(path: str):
    name = Path(path).name
    if not name.endswith(".go") or name.endswith("_test.go"):
        return None
    base = name[:-3]
    for prefix in ("resource_nutanix_", "data_source_nutanix_", "resource_", "data_source_"):
        if base.startswith(prefix):
            base = base[len(prefix) :]
            break
    else:
        # helpers / misc — skip for test-scoped targeting
        return None
    if base.endswith("_v2"):
        base = base[:-3]
    return base  # e.g. volume_group_disk

def test_names(pattern: str) -> list[str]:
    # Split Go -run alternatives; ignore empty.
    return [p for p in pattern.split("|") if p]

def match_files_for_test(test: str, stems: dict) -> list:
    # Strip common Acc prefixes for matching.
    core = test
    for pref in (
        "TestAccV2Nutanix",
        "TestAccNutanix",
        "TestAccV2",
        "TestAccFC",
        "TestAccEra",
        "TestAccFoundation",
        "TestAccKarbon",
        "TestAcc",
    ):
        if core.startswith(pref):
            core = core[len(pref) :]
            break
    # Drop scenario suffix: _Basic, _WithLimit, ...
    core = core.split("_", 1)[0]

    want_resource = "Resource" in core and "DataSource" not in core
    want_datasource = "DataSource" in core

    best_len = 0
    best_files = []
    for snake, files in stems.items():
        camel = snake_to_camel(snake)
        if not camel:
            continue
        if camel not in core:
            continue
        ranked = []
        for f in files:
            base = Path(f).name
            if want_resource and not base.startswith("resource_"):
                continue
            if want_datasource and not base.startswith("data_source_"):
                continue
            ranked.append(f)
        # If kind filter removed everything (unusual naming), fall back to all stem files.
        candidates = ranked if ranked else list(files)
        if len(camel) > best_len:
            best_len = len(camel)
            best_files = candidates
        elif len(camel) == best_len:
            best_files.extend(candidates)

    seen = set()
    out = []
    for f in best_files:
        if f not in seen:
            seen.add(f)
            out.append(f)
    return out

# Collect candidate files from package dir and/or profile.
stems: dict[str, list[str]] = {}

def add_file(path: str) -> None:
    key = file_stem_key(path)
    if not key:
        return
    stems.setdefault(key, [])
    # normalize to path as it appears in profile when possible
    if path not in stems[key]:
        stems[key].append(path)

pkg = package_path.strip()
if pkg.startswith("./"):
    pkg_dir = Path(pkg[2:])
elif pkg and pkg != "./...":
    pkg_dir = Path(pkg)
else:
    pkg_dir = None

if pkg_dir and pkg_dir.is_dir():
    for p in sorted(pkg_dir.glob("*.go")):
        if p.name.endswith("_test.go"):
            continue
        add_file(str(p).replace("\\", "/"))

# Also index paths from the profile (module-qualified).
mode = None
blocks: list[str] = []
with open(in_path) as f:
    mode = f.readline()
    for line in f:
        blocks.append(line)
        loc = line.rsplit(" ", 2)[0]
        file_path = loc.split(":", 1)[0]
        # Map to repo-relative if possible
        marker = "terraform-provider-nutanix/"
        rel = file_path
        if marker in file_path:
            rel = file_path.split(marker, 1)[1]
        add_file(rel)
        add_file(file_path)

wanted_rels: set[str] = set()
wanted_suffixes: set[str] = set()
for t in test_names(run_pattern):
    matched = match_files_for_test(t, stems)
    for m in matched:
        wanted_rels.add(m)
        wanted_suffixes.add(Path(m).name)

if not wanted_rels:
    # Fallback: keep only files that had any hit (still better than whole package zeros).
    print(
        f"cover-filter: no file match for -run={run_pattern}; "
        "falling back to files with executed statements",
        file=sys.stderr,
    )
    keep_files = set()
    for line in blocks:
        loc, num_stmt, count = line.rsplit(" ", 2)
        if int(count) > 0:
            keep_files.add(loc.split(":", 1)[0])
    with open(out_path, "w") as out:
        out.write(mode if mode else "mode: atomic\n")
        for line in blocks:
            loc = line.rsplit(" ", 2)[0]
            if loc.split(":", 1)[0] in keep_files:
                out.write(line)
    print(f"cover-filter: wrote {out_path} ({len(keep_files)} touched files)")
    sys.exit(0)

def keep_line(line: str) -> bool:
    loc = line.rsplit(" ", 2)[0]
    file_path = loc.split(":", 1)[0]
    base = Path(file_path).name
    if base in wanted_suffixes:
        return True
    for w in wanted_rels:
        if file_path.endswith(w) or file_path.endswith("/" + w) or file_path == w:
            return True
        if w.endswith(file_path) or file_path.endswith(Path(w).name):
            return True
    return False

kept = 0
with open(out_path, "w") as out:
    out.write(mode if mode else "mode: atomic\n")
    for line in blocks:
        if keep_line(line):
            out.write(line)
            kept += 1

print(
    "cover-filter: tests → files: "
    + ", ".join(sorted(wanted_suffixes))
)
print(f"cover-filter: wrote {out_path} ({kept} blocks)")
PY
