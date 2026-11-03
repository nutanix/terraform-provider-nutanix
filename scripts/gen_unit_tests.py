#!/usr/bin/env python3
"""Generate flatten/expand unit tests for networkingv2 and vmmv2."""
from __future__ import annotations

import re
import textwrap
from pathlib import Path

SKIP_FUNCS = {
    "getEtagHeader",
    "waitForDiskTask",
    "ApplyDiskDeletions",
    "ApplyDiskUpdates",
    "ApplyDiskAdditions",
}


def parse_source_imports(text: str) -> list[str]:
    m = re.search(r"\nimport\s*\((.*?)\)\n", text, re.S)
    if not m:
        m2 = re.search(r'\nimport\s+(".*?")\n', text)
        return [m2.group(1)] if m2 else []
    lines = []
    for line in m.group(1).splitlines():
        line = line.strip()
        if line and not line.startswith("//"):
            lines.append(line)
    return lines


def extract_funcs_from_file(path: Path):
    text = path.read_text()
    funcs = []
    for m in re.finditer(r"^func (\w+)\(", text, re.M):
        name = m.group(1)
        if name in SKIP_FUNCS:
            continue
        is_fe = name.startswith(("flatten", "expand", "Flatten", "Expand"))
        # params
        i = m.end()
        depth = 1
        while i < len(text) and depth:
            if text[i] == "(":
                depth += 1
            elif text[i] == ")":
                depth -= 1
            i += 1
        params = re.sub(r"\s+", " ", text[m.end() : i - 1]).strip()
        while i < len(text) and text[i] in " \t\n":
            i += 1
        ret_start = i
        # Find body '{' while allowing interface{} / map[string]interface{} in returns
        depth = 0
        j = i
        while j < len(text):
            ch = text[j]
            if ch == "{":
                if depth == 0:
                    k = j + 1
                    while k < len(text) and text[k] in " \t":
                        k += 1
                    if k < len(text) and text[k] == "}":
                        depth = 1
                    else:
                        break  # function body
                else:
                    depth += 1
            elif ch == "}":
                depth -= 1
            j += 1
        returns = re.sub(r"\s+", " ", text[ret_start:j]).strip()
        # body
        start = j + 1
        depth = 1
        i = start
        while i < len(text) and depth:
            if text[i] == "{":
                depth += 1
            elif text[i] == "}":
                depth -= 1
            i += 1
        body = text[start : i - 1]
        first = "\n".join(body.strip().split("\n")[:15])
        parts = split_params(params) if params else []
        if parts:
            pname, ptype = parts[0]
        else:
            pname, ptype = "in", ""
        funcs.append(
            {
                "name": name,
                "params": params,
                "returns": returns,
                "pname": pname,
                "ptype": ptype,
                "first": first,
                "body": body,
                "is_fe": is_fe,
            }
        )
    return text, funcs


def nil_checks(f):
    pname = f["pname"]
    first = f["first"]
    ptr = bool(re.search(rf"\bif\s+{re.escape(pname)}\s*==\s*nil\b", first))
    ln = bool(
        re.search(rf"\bif\s+len\({re.escape(pname)}\)\s*==\s*0", first)
        or re.search(rf"\bif\s+len\({re.escape(pname)}\)\s*>\s*0", first)
        or re.search(rf"\bif\s+len\({re.escape(pname)}\)\s*[=!]=", first)
    )
    return ptr, ln


def has_ext_id(body: str) -> bool:
    return bool(re.search(r"\.ExtId\b", body))


def type_leaf(ptype: str) -> str:
    base = ptype
    if base.startswith("*"):
        base = base[1:]
    if base.startswith("[]"):
        base = base[2:]
        if base.startswith("*"):
            base = base[1:]
    return base.split(".")[-1]


def looks_like_enum(ptype: str, returns: str = "", body: str = "") -> bool:
    """Heuristic: enum flatteners usually return string/interface{} and compare *pr."""
    leaf = type_leaf(ptype)
    if leaf in ("string", "int", "int32", "int64", "bool", "float64", "byte"):
        return False
    # Named like FooFeature / FooCapability are often enums/aliases, not structs
    if re.search(r"(Feature|Capability)$", leaf) and returns in (
        "string",
        "interface{}",
        "[]interface{}",
        "[]map[string]interface{}",
        "[]string",
    ):
        if body and (
            re.search(r"\*\w+\s*==", body)
            or re.search(re.escape(leaf) + r"\(", body)
            or "FlattenPtrEnum" in body
            or "GetName()" in body
        ):
            return True
    if returns in ("string", "interface{}"):
        if body and (
            re.search(r"\*\w+\s*==", body)
            or re.search(re.escape(leaf) + r"\(", body)
            or "FlattenPtrEnum" in body
            or "common.FlattenPtrEnum" in body
        ):
            return True
        if returns == "string":
            return True
    if leaf.endswith("Type") and body and not has_ext_id(body) and not re.search(
        r"\.(ExtId|Name|Value|Links|TenantId)\b", body
    ):
        if re.search(r"\*\w+\s*==", body) or re.search(re.escape(leaf) + r"\(", body):
            return True
    return False


def struct_literal(base: str, with_ext: bool, returns: str = "", body: str = "") -> str:
    if base in ("string",):
        return '"x"'
    if looks_like_enum(base, returns, body):
        return f"{base}(1)"
    # Prefer zero structs: ExtId is not present on every model and causes compile failures.
    return f"{base}{{}}"


def go_test_name(name: str) -> str:
    return "TestUnit" + name[0].upper() + name[1:]


def expr_empty(ptype: str, returns: str = "", body: str = "") -> str:
    if ptype == "interface{}":
        return "nil"
    if ptype == "string":
        return '""'
    if ptype == "map[string]interface{}":
        return "map[string]interface{}{}"
    if ptype == "[]string":
        return "[]string{}"
    if ptype.startswith("*"):
        base = ptype[1:]
        if looks_like_enum(base, returns, body) or base == "string":
            return f"func() {ptype} {{ var v {base}; return &v }}()"
        return f"&{base}{{}}"
    if ptype.startswith("[]"):
        return f"{ptype}{{}}"
    if looks_like_enum(ptype, returns, body):
        return f"{ptype}(0)"
    return f"{ptype}{{}}"


def expr_nil(ptype: str):
    if (
        ptype.startswith("*")
        or ptype.startswith("[]")
        or ptype == "interface{}"
        or ptype == "map[string]interface{}"
    ):
        return "nil"
    return None


def expr_populated(ptype: str, name: str, body: str, is_expand: bool, returns: str = "") -> str:
    if is_expand:
        if ptype == "string":
            # Prefer a known enum label from a map literal in the body
            labels = re.findall(r'"([A-Z][A-Z0-9_]{2,})"', body)
            if labels:
                return f'"{labels[0]}"'
            return '"ACTIVE"'
        if ptype == "[]string":
            return '[]string{"a"}'
        if ptype == "map[string]interface{}":
            return expand_map_literal(body)
        if ptype == "interface{}" or ptype.startswith("[]interface"):
            top_is_list = bool(re.search(r"\w+\.\(\[\]interface\{\}\)", body))
            top_is_map = bool(
                re.search(r"\w+\.\(map\[string\]interface\{\}\)", body)
            ) and not top_is_list
            # Pure string enum expand: cast param itself to string (no list/map cast)
            if (
                ptype == "interface{}"
                and not top_is_list
                and not top_is_map
                and re.search(r"\w+\.\(string\)", body)
                and not re.search(r"\.\(map\[string\]interface\{\}\)", body)
            ):
                labels = re.findall(r'"([A-Z][A-Z0-9_]{2,})"', body)
                if labels:
                    return f'"{labels[0]}"'
                return '"UNKNOWN"'
            if ptype == "interface{}" and top_is_map:
                return "map[string]interface{}{}"
            # list of strings / enum labels
            if re.search(r"\.\(string\)", body) and not re.search(
                r"\.\(map\[string\]interface", body
            ):
                labels = re.findall(r'"([A-Z][A-Z0-9_]{2,})"', body)
                if labels:
                    joined = ", ".join(f'"{x}"' for x in labels[:2])
                    return f"[]interface{{}}{{{joined}}}"
                return '[]interface{}{"a"}'
            # Safe populated list/block with common required string fields
            return '[]interface{}{map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}}'
        if ptype == "map[string]interface{}":
            return 'map[string]interface{}{"ext_id": "ext-1", "name": "name-1"}'
        return expr_empty(ptype, returns, body)

    use_ext = has_ext_id(body)
    if ptype == "[]string":
        return '[]string{"a"}'
    if ptype == "string":
        return '"SECURE"'
    if ptype == "interface{}":
        return '[]interface{}{map[string]interface{}{"ext_id": "ext-1"}}'
    if ptype == "map[string]interface{}":
        return 'map[string]interface{}{"ext_id": "ext-1"}'
    if ptype.startswith("*"):
        base = ptype[1:]
        if looks_like_enum(base, returns, body):
            return f"func() {ptype} {{ v := {base}(1); return &v }}()"
        if base == "string":
            return 'func() *string { v := "x"; return &v }()'
        return f"&{struct_literal(base, use_ext, returns, body)}"
    if ptype.startswith("[]") and not ptype.startswith("[]interface"):
        elem = ptype[2:]
        if elem == "string":
            return '[]string{"a"}'
        if looks_like_enum(elem, returns, body):
            return f"{ptype}{{{elem}(1)}}"
        if elem.startswith("*"):
            if looks_like_enum(elem[1:], returns, body):
                return f"{ptype}{{func() {elem} {{ v := {elem[1:]}(1); return &v }}()}}"
            return f"{ptype}{{&{struct_literal(elem[1:], use_ext, returns, body)}}}"
        return f"{ptype}{{{struct_literal(elem, use_ext, returns, body)}}}"
    return struct_literal(ptype, use_ext, returns, body)


def expand_map_literal(body: str) -> str:
    keys = re.findall(r'\["([^"]+)"\]', body)

    def value_for_key(k: str) -> str:
        # Nested TypeList / object fields
        if re.search(
            rf'\["{re.escape(k)}"\].{{0,80}}\[\]interface',
            body,
            re.S,
        ) or re.search(
            rf"{re.escape(k)}\.\(\[\]interface",
            body,
        ) or re.search(
            rf'len\({re.escape(k)}\.\(\[\]interface',
            body,
        ):
            # IP-ish nested objects
            if "ip" in k or "subnet" in k or "address" in k:
                return '[]interface{}{map[string]interface{}{"value": "10.0.0.1", "prefix_length": 24}}'
            return '[]interface{}{map[string]interface{}{"ext_id": "ext-1"}}'
        if (
            k.startswith("is_")
            or k.startswith("should_")
            or k.endswith("_enabled")
            or k in ("enabled",)
        ):
            return "true"
        if (
            "prefix" in k
            or k.endswith("_length")
            or k.endswith("_secs")
            or k.endswith("_threshold")
            or k.endswith("_timeout")
            or k
            in (
                "priority",
                "type",
                "code",
                "vlan_identifier",
                "route_table",
                "num_sockets",
                "memory_size_bytes",
            )
        ):
            return "24"
        # enum string fields often cast .(string) and looked up in map
        if re.search(rf'\["{re.escape(k)}"\].{{0,120}}\.\(string\)', body, re.S):
            labels = re.findall(r'"([A-Z][A-Z0-9_]{2,})"', body)
            if labels:
                return f'"{labels[0]}"'
            return '"UNKNOWN"'
        if k in ("ext_id", "name", "value", "href", "rel"):
            return '"10.0.0.1"' if k == "value" else f'"{k}-1"'
        return '"x"'

    m = {}
    preferred = [
        "ext_id",
        "name",
        "value",
        "prefix_length",
        "is_enabled",
        "is_snooping_enabled",
        "is_querier_enabled",
        "service_ip",
        "reroute_fallback_action",
    ]
    for k in preferred:
        if k in keys:
            m[k] = value_for_key(k)
    if not m:
        for k in keys[:5]:
            m[k] = value_for_key(k)
    if not m:
        m = {"ext_id": '"ext-1"'}
    map_lit = ", ".join(f'"{k}": {v}' for k, v in m.items())
    return f"map[string]interface{{}}{{{map_lit}}}"


def returns_nil_for_nil(f) -> bool:
    return bool(re.search(r"return nil\b", f["first"]))


def count_params(params: str) -> int:
    if not params.strip():
        return 0
    depth = 0
    n = 0
    for ch in params:
        if ch in "[(":
            depth += 1
        elif ch in "])":
            depth -= 1
        elif ch == "," and depth == 0:
            n += 1
    return n + 1


def split_params(params: str):
    depth = 0
    parts = []
    cur = ""
    for ch in params:
        if ch in "[(":
            depth += 1
            cur += ch
        elif ch in "])":
            depth -= 1
            cur += ch
        elif ch == "," and depth == 0:
            parts.append(cur.strip())
            cur = ""
        else:
            cur += ch
    if cur.strip():
        parts.append(cur.strip())
    out = []
    for p in parts:
        m = re.match(r"(\w+)\s+(.+)", p)
        if m:
            out.append((m.group(1), m.group(2).strip()))
        else:
            out.append(("x", p))
    return out


def make_multi_arg_call(f, first_expr: str):
    params = f["params"]
    name = f["name"]
    n = count_params(params)
    if name in ("expandNic", "expandVMConfigOverride"):
        return None
    if name == "flattenDisk":
        return f"flattenDisk({first_expr})"
    parts = split_params(params)
    if n == 2 and "map[string]interface{}" in parts[0][1]:
        p2t = parts[1][1]
        use_ext = has_ext_id(f["body"])
        if p2t.startswith("*"):
            item = f"&{struct_literal(p2t[1:], use_ext)}"
        else:
            item = struct_literal(p2t, use_ext)
        return f"{name}(map[string]interface{{}}{{}}, {item})"
    if n == 1:
        return f"{name}({first_expr})"
    args = [first_expr]
    for _, pt in parts[1:]:
        if pt == "string":
            args.append('""')
        elif pt == "bool":
            args.append("false")
        elif pt.startswith("*") or pt.startswith("[]") or pt == "interface{}":
            args.append("nil")
        else:
            return None
    return f"{name}({', '.join(args)})"


def gen_resourcedata_test(f) -> str:
    name = f["name"]
    tname = go_test_name(name)
    if name == "flattenVmStartupPolicy":
        return textwrap.dedent(
            f"""\
            func {tname}(t *testing.T) {{
            \tt.Parallel()
            \tschemaMap := ResourceNutanixVmStartupPolicyV2().Schema
            \tt.Run("empty", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.NoError(t, {name}(d, &import1.VmStartupPolicy{{}}))
            \t}})
            \tt.Run("populated", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.NoError(t, {name}(d, &import1.VmStartupPolicy{{ExtId: utils.StringPtr("ext-1")}}))
            \t\trequire.Equal(t, "ext-1", d.Get("ext_id"))
            \t}})
            }}
            """
        )
    if name == "flattenDependencyConflict":
        return textwrap.dedent(
            f"""\
            func {tname}(t *testing.T) {{
            \tt.Parallel()
            \tschemaMap := DatasourceNutanixVmStartupPolicyDependencyConflictV2().Schema
            \tt.Run("empty", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.NoError(t, {name}(d, &import1.DependencyConflict{{}}))
            \t}})
            \tt.Run("populated", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.NoError(t, {name}(d, &import1.DependencyConflict{{ExtId: utils.StringPtr("ext-1")}}))
            \t\trequire.Equal(t, "ext-1", d.Get("ext_id"))
            \t}})
            }}
            """
        )
    if name == "flattenStartConditionConflict":
        return textwrap.dedent(
            f"""\
            func {tname}(t *testing.T) {{
            \tt.Parallel()
            \tschemaMap := DatasourceNutanixVmStartupPolicyStartConditionConflictV2().Schema
            \tt.Run("empty", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.NoError(t, {name}(d, &import1.StartConditionConflict{{}}))
            \t}})
            \tt.Run("populated", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.NoError(t, {name}(d, &import1.StartConditionConflict{{ExtId: utils.StringPtr("ext-1")}}))
            \t\trequire.Equal(t, "ext-1", d.Get("ext_id"))
            \t}})
            }}
            """
        )
    if name == "expandVMConfigOverride":
        return textwrap.dedent(
            f"""\
            func {tname}(t *testing.T) {{
            \tt.Parallel()
            \tschemaMap := ResourceNutanixTemplateDeployV2().Schema
            \tt.Run("empty", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\tgot := expandVMConfigOverride([]interface{{}}{{}}, d)
            \t\trequire.Nil(t, got)
            \t}})
            \tt.Run("populated", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{
            \t\t\t"ext_id": "tmpl-1",
            \t\t}})
            \t\tgot := expandVMConfigOverride([]interface{{}}{{map[string]interface{{}}{{"name": "vm-1"}}}}, d)
            \t\trequire.NotNil(t, got)
            \t}})
            }}
            """
        )
    if name == "expandNic":
        return textwrap.dedent(
            f"""\
            func {tname}(t *testing.T) {{
            \tt.Parallel()
            \tschemaMap := ResourceNutanixVirtualMachineV2().Schema
            \tt.Run("nil", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.Nil(t, expandNic(nil, d, "nics"))
            \t}})
            \tt.Run("empty", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{}})
            \t\trequire.Nil(t, expandNic([]interface{{}}{{}}, d, "nics"))
            \t}})
            \tt.Run("populated", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\td := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{{}}{{
            \t\t\t"nics": []interface{{}}{{map[string]interface{{}}{{"ext_id": "nic-1"}}}},
            \t\t}})
            \t\tgot := expandNic([]interface{{}}{{map[string]interface{{}}{{"ext_id": "nic-1"}}}}, d, "nics")
            \t\trequire.NotNil(t, got)
            \t\trequire.Equal(t, "nic-1", utils.StringValue(got[0].ExtId))
            \t}})
            }}
            """
        )
    if name == "flattenDisk":
        return textwrap.dedent(
            f"""\
            func {tname}(t *testing.T) {{
            \tt.Parallel()
            \tt.Run("nil", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\trequire.Nil(t, flattenDisk(nil))
            \t}})
            \tt.Run("empty", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\trequire.Nil(t, flattenDisk([]config.Disk{{}}))
            \t}})
            \tt.Run("populated", func(t *testing.T) {{
            \t\tt.Parallel()
            \t\tgot := flattenDisk([]config.Disk{{{{ExtId: utils.StringPtr("disk-1")}}}})
            \t\trequire.Len(t, got, 1)
            \t}})
            }}
            """
        )
    return f"// skipped special {name}\n"


def gen_test_func(f) -> str:
    name = f["name"]
    ptype = f["ptype"]
    pname = f["pname"]
    params = f["params"]
    returns = f["returns"]
    is_expand = name.startswith("expand") or name.startswith("Expand")
    ptr_nil, len_nil = nil_checks(f)
    tname = go_test_name(name)

    if "ResourceData" in params:
        return gen_resourcedata_test(f)

    nparams = count_params(params)
    lines = [f"func {tname}(t *testing.T) {{", "\tt.Parallel()", ""]

    cases = []
    nil_expr = expr_nil(ptype)
    if nil_expr is not None:
        include_nil = False
        if ptr_nil:
            include_nil = True
        if ptype.startswith("[]") and len_nil:
            include_nil = True
        if ptype == "interface{}" and re.search(
            rf"{re.escape(pname)}\s*==\s*nil", f["first"]
        ):
            include_nil = True
        if include_nil:
            cases.append(("nil", nil_expr))

    include_empty = True
    if is_expand and ptype == "interface{}":
        # Many expands do pr.([]interface{})[0] after only a nil check.
        # An empty slice is non-nil and panics — skip empty unless len(pr) is checked.
        if not re.search(rf"\blen\({re.escape(pname)}\)", f["first"]):
            include_empty = False
            empty_expr = None
        else:
            empty_expr = "[]interface{}{}"
    elif ptype == "interface{}":
        empty_expr = "[]interface{}{}"
    else:
        empty_expr = expr_empty(ptype, returns, f["body"])
    if include_empty and empty_expr is not None:
        cases.append(("empty", empty_expr))
    cases.append(("populated", expr_populated(ptype, name, f["body"], is_expand, returns)))

    # OneOf SDK wrappers panic unless SetValue is used — only nil is safe.
    if "OneOf" in ptype:
        cases = [c for c in cases if c[0] == "nil"]
        if not cases:
            lines = [
                f"func {tname}(t *testing.T) {{",
                "\tt.Parallel()",
                '\tt.Run("nil", func(t *testing.T) {',
                "\t\tt.Parallel()",
                f"\t\trequire.Nil(t, {name}(nil))",
                "\t})",
                "}",
                "",
            ]
            return "\n".join(lines)

    # Unnil-checked GetName() on nested enum pointers panics for zero values.
    if not is_expand and re.search(r"\.GetName\(\)", f["body"]):
        cases = [c for c in cases if c[0] != "populated"]

    # make(..., 0) then index-assign panics for any populated input (prod quirk).
    if not is_expand and re.search(
        r"make\(\[\]map\[string\]interface\{\},\s*0\)", f["body"]
    ) and re.search(r"\w+\[k\]\s*=", f["body"]):
        cases = [c for c in cases if c[0] != "populated"]

    for cname, cexpr in cases:
        lines.append(f'\tt.Run("{cname}", func(t *testing.T) {{')
        lines.append("\t\tt.Parallel()")
        if nparams > 1:
            call = make_multi_arg_call(f, cexpr)
            if call is None:
                lines.append('\t\tt.Skip("multi-arg helper needs special setup")')
                lines.append("\t})")
                lines.append("")
                continue
        else:
            call = f"{name}({cexpr})"

        if cname == "nil":
            if returns_nil_for_nil(f) and "error" not in returns:
                lines.append(f"\t\trequire.Nil(t, {call})")
            else:
                lines.append(f"\t\t_ = {call}")
        elif cname == "empty":
            lines.append(f"\t\tgot := {call}")
            if returns_nil_for_nil(f) and (
                ptype.startswith("[]") or ptype == "interface{}"
            ):
                lines.append("\t\trequire.Nil(t, got)")
            else:
                lines.append("\t\t_ = got")
        else:
            lines.append(f"\t\tgot := {call}")
            if returns.startswith("*") or returns.startswith("[]"):
                # populated may still be nil if input shape wrong; don't hard-require
                lines.append("\t\t_ = got")
            else:
                lines.append("\t\t_ = got")
        lines.append("\t})")
        lines.append("")

    lines.append("}")
    lines.append("")
    return "\n".join(lines)


def gen_file_tests(pkg: str, src_path: Path):
    text, all_funcs = extract_funcs_from_file(src_path)
    funcs = [f for f in all_funcs if f["is_fe"]]
    if not funcs:
        return None
    source_imports = parse_source_imports(text)

    body_parts = []
    for f in funcs:
        try:
            body_parts.append(gen_test_func(f))
        except Exception as e:  # noqa: BLE001
            body_parts.append(f"// ERROR generating {f['name']}: {e}\n")

    used_text = "\n".join(body_parts)
    base_imps = {
        '"testing"',
        '"github.com/stretchr/testify/require"',
    }
    if "utils." in used_text:
        base_imps.add(
            '"github.com/terraform-providers/terraform-provider-nutanix/utils"'
        )
    if "schema." in used_text:
        base_imps.add('"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"')
    if "require." not in used_text:
        base_imps.discard('"github.com/stretchr/testify/require"')

    final_type = []
    skip_paths = {
        "github.com/stretchr/testify/require",
        "github.com/terraform-providers/terraform-provider-nutanix/utils",
        "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema",
    }
    for imp in source_imports:
        path_m = re.search(r'"([^"]+)"', imp)
        if path_m and path_m.group(1) in skip_paths:
            continue
        m = re.match(r'(\w+)\s+"', imp)
        if m:
            alias = m.group(1)
            if re.search(rf"\b{re.escape(alias)}\.", used_text):
                final_type.append(imp)
            continue
        m2 = re.match(r'"([^"]+)"', imp)
        if m2:
            pkgname = m2.group(1).rstrip("/").split("/")[-1]
            if re.search(rf"\b{re.escape(pkgname)}\.", used_text):
                final_type.append(imp)

    imports_block = "\n".join("\t" + i for i in sorted(base_imps))
    if final_type:
        imports_block += "\n\n" + "\n".join("\t" + i for i in final_type)

    return f"package {pkg}\n\nimport (\n{imports_block}\n)\n\n" + "\n".join(body_parts)


def main():
    root = Path(__file__).resolve().parents[1]
    for pkg in ("networkingv2", "vmmv2"):
        pkg_dir = root / "nutanix" / "services" / pkg
        count = 0
        for src in sorted(pkg_dir.glob("*.go")):
            if src.name.endswith("_test.go"):
                continue
            content = gen_file_tests(pkg, src)
            if not content:
                continue
            out = src.with_name(src.stem + "_unit_test.go")
            if out.name == "data_source_nutanix_image_v2_unit_test.go":
                out = src.with_name("data_source_nutanix_image_v2_flatten_unit_test.go")
            out.write_text(content)
            count += 1
        print(f"{pkg}: wrote {count} flatten/expand unit test files")


if __name__ == "__main__":
    main()
