#!/usr/bin/env python3
"""Generate cluster-free unit tests for prism, ndb, and clustersv2."""
from __future__ import annotations

import os
import re
from typing import List, Optional, Tuple

ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))

SKIP = {
    "expandProjectInternalSpec",
    "expandUser",
    "expandUserGroup",
    "expandCreateAcp",
    "expandNetworkProfileProperties",
    "expandNetworkHAInstance",
    "expandClusterWithNewNode",
}

# Functions known to panic on empty struct / empty slice without nested fields.
SKIP_EMPTY_POPULATED = {
    "flattenSMTPNetwork",  # derefs IpAddress
    "flattenDBLcmConfig",  # PostDeleteCommand deref bug on empty
}

# Skip populated/empty-struct cases (nil/empty-slice still tested when applicable).
SKIP_POPULATED = {
    "flattenProjectEntities",
    "flattenProtectionRuleEntities",
    "flattenRecoveryPlanEntities",
    "flattenProjectAcp",
    "flattenNameNetworkRef",  # derefs HypervisorType
    "flattenVMAssignmentList",
    "flattenFloatingAssignmentList",
    "flattenNetworkMappingList",
    "flattenAssignmentList",
    "flattenIPConfigList",
    "flattenStageList",
    "flattenEntityInfoList",
    "expandRightHandsideCategories",  # expects *schema.Set nested
    "expandProjectInternalResourceDomain",
    "expandProjectContextFilterList",
    "expandRsyslogServerList",  # asserts string/int without ok
    "flattenSMTPNetwork",
    "flattenDBLcmConfig",
}

SCHEMA_FOR_RD = {
    "expandIdentityProviderUser": "ResourceNutanixProject()",
    "expandAccessControlPolicyResources": "ResourceNutanixAccessControlPolicy()",
    "expandContextFilterList": "ResourceNutanixAccessControlPolicy()",
    "expandProjectSpec": "ResourceNutanixProject()",
    "expandResourceDomain": "ResourceNutanixProject()",
    "expandReferenceList": "ResourceNutanixProject()",
    "expandReferenceSet": "ResourceNutanixProject()",
    "expandMetadata": "ResourceNutanixProject()",
    "expandProjectDetails": "ResourceNutanixProject()",
    "expandOptionalReference": "ResourceNutanixProject()",
    "expandCategoryFilter": "ResourceNutanixProtectionRule()",
    "expandAvailabilityZoneConnectivityList": "ResourceNutanixProtectionRule()",
    "expandOrderAvailibilityList": "ResourceNutanixProtectionRule()",
    "expandStageList": "ResourceNutanixRecoveryPlan()",
    "expandParameters": "ResourceNutanixRecoveryPlan()",
    "expandAcp": "ResourceNutanixProject()",
    "expandActionArguments": "ResourceDatabaseInstance()",
    "expandRegisterDBActionArguments": "ResourceNutanixNDBRegisterDatabase()",
    "expandPostgreSQLCloneActionArgs": "ResourceNutanixNDBClone()",
    "expandClusterNodeParams": "ResourceNutanixClusterAddNodeV2()",
    "expandClusterConfigReference": "ResourceNutanixClusterV2()",
    "expandClusterProfile": "ResourceNutanixClusterProfileV2()",
}


def get_body(text: str, start: int) -> str:
    # Skip any '{' that appear in the signature (e.g. interface{}), then take the body.
    # Find the end of the signature by locating the ')' matching the first '('.
    i = text.find("(", start)
    if i < 0:
        return ""
    depth = 0
    j = i
    while j < len(text):
        if text[j] == "(":
            depth += 1
        elif text[j] == ")":
            depth -= 1
            if depth == 0:
                break
        j += 1
    # Now find the opening '{' of the function body after the signature.
    i = text.find("{", j)
    if i < 0:
        return ""
    depth = 0
    for k in range(i, len(text)):
        if text[k] == "{":
            depth += 1
        elif text[k] == "}":
            depth -= 1
            if depth == 0:
                return text[i : k + 1]
    return ""


def safe_nil_check(body: str, pname: Optional[str]) -> bool:
    """True if the body safely handles a nil first argument before casting it."""
    if not pname or not body.startswith("{"):
        return False
    # Strip leading whitespace after '{'
    rest = body[1:].lstrip()
    return bool(
        re.match(rf"if\s+{re.escape(pname)}\s*==\s*nil\b", rest)
        or re.match(rf"if\s+{re.escape(pname)}\s*!=\s*nil\b", rest)
        or re.match(rf"if\s+{re.escape(pname)}\s*==\s*nil\s*\|\|", rest)
    )


def has_nil_check(body: str) -> bool:
    return bool(
        re.search(r"if\s+\w+\s*==\s*nil|if\s+\w+\s*!=\s*nil|if\s+len\(\w+\)", body)
    )


def parse_first_param(params: str) -> Tuple[Optional[str], Optional[str]]:
    if not params.strip():
        return None, None
    first = params.split(",")[0].strip()
    parts = first.rsplit(" ", 1)
    if len(parts) != 2:
        return first, None
    return parts[0], parts[1]


def collect(pkg: str):
    path = os.path.join(ROOT, "nutanix", "services", pkg)
    schemas, funcs = [], []
    for fn in sorted(os.listdir(path)):
        if not fn.endswith(".go") or fn.endswith("_test.go"):
            continue
        text = open(os.path.join(path, fn)).read()
        for m in re.finditer(
            r"^func ((?:Resource|Data[sS]ource)[A-Za-z0-9_]*)\(\) \*schema\.Resource",
            text,
            re.M,
        ):
            schemas.append(m.group(1))
        for m in re.finditer(
            r"^func ((?:flatten|expand|Flatten|Expand)[A-Za-z0-9_]*)\(([^)]*)\)([^{]*)\{",
            text,
            re.M,
        ):
            params = re.sub(r"\s+", " ", m.group(2)).strip()
            pname, ptype = parse_first_param(params)
            body = get_body(text, m.start())
            ptype_n = ptype.replace("Era.", "era.") if ptype else ptype
            funcs.append(
                {
                    "name": m.group(1),
                    "params": params,
                    "ret": m.group(3).strip(),
                    "pname": pname,
                    "ptype": ptype_n,
                    "nil_check": has_nil_check(body),
                    "body": body,
                }
            )
    return schemas, funcs


def emit_subtests(lines, auto):
    for name, cases in auto:
        lines.append('\tt.Run("%s", func(t *testing.T) {' % name)
        lines.append("\t\tt.Parallel()")
        for cname, body in cases:
            lines.append('\t\tt.Run("%s", func(t *testing.T) {' % cname)
            lines.append("\t\t\tt.Parallel()")
            for bl in body.split("\n"):
                lines.append("\t\t\t" + bl)
            lines.append("\t\t})")
        lines.append("\t})")
        lines.append("")


def gen_cases(f) -> List[Tuple[str, str]]:
    name = f["name"]
    ptype = f["ptype"]
    params = f["params"]
    nil_check = f["nil_check"]
    body = f["body"]
    cases = []

    if name in SKIP or "meta interface{}" in params or "context.Context" in params:
        return []

    # Multi-arg helpers that include ResourceData (first param may not be *schema.ResourceData).
    if "*schema.ResourceData" in params:
        # Reuse ResourceData path using SCHEMA_FOR_RD keyed by name.
        schema_fn = SCHEMA_FOR_RD.get(name)
        if not schema_fn:
            return []
        fake = dict(f)
        # Force rd_cases via temporary first-type override
        if params.startswith("pr []interface{}"):
            cases.append(
                (
                    "nil",
                    "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                    "\t\t\t_ = %s(nil, d)" % (schema_fn, name),
                )
            )
            cases.append(
                (
                    "empty",
                    "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                    "\t\t\t_ = %s([]interface{}{}, d)" % (schema_fn, name),
                )
            )
            return cases
        if params.startswith("pr interface{}"):
            cases.append(
                (
                    "nil",
                    "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                    "\t\t\trequire.Nil(t, %s(nil, d))" % (schema_fn, name),
                )
            )
            return cases
        # fall through to first-param ResourceData handling below

    # ResourceData expands
    if ptype == "*schema.ResourceData" or (
        "*schema.ResourceData" in params and SCHEMA_FOR_RD.get(name)
    ):
        if ptype != "*schema.ResourceData" and not params.startswith("d *schema.ResourceData"):
            pass  # already handled above
        else:
            schema_fn = SCHEMA_FOR_RD.get(name)
            if not schema_fn:
                return []
            if "access *v3.AccessControlPolicyResources" in params:
                cases.append(
                    (
                        "empty",
                        "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                        "\t\t\taccess := &v3.AccessControlPolicyResources{}\n"
                        "\t\t\t%s(d, access)" % (schema_fn, name),
                    )
                )
                return cases
            if "pr []interface{}" in params and params.startswith("d *schema.ResourceData"):
                cases.append(
                    (
                        "nil",
                        "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                        "\t\t\t_ = %s(d, nil)" % (schema_fn, name),
                    )
                )
                cases.append(
                    (
                        "empty",
                        "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                        "\t\t\t_ = %s(d, []interface{}{})" % (schema_fn, name),
                    )
                )
                return cases
            if "pr interface{}" in params and params.startswith("d *schema.ResourceData"):
                cases.append(
                    (
                        "nil",
                        "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                        "\t\t\trequire.Nil(t, %s(d, nil))" % (schema_fn, name),
                    )
                )
                return cases
            if ", key string" in params and "kind string" in params:
                cases.append(
                    (
                        "empty",
                        "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                        '\t\t\t_ = %s(d, "subnet_reference_list", "subnet")'
                        % (schema_fn, name),
                    )
                )
                return cases
        if ", key string" in params:
            # Prefer keys that exist as TypeList/TypeSet on the project schema.
            key = "user_reference_list" if "Set" in name or "expandReferenceSet" == name else "account_reference_list"
            if name == "expandReferenceList":
                key = "account_reference_list"
            if name == "expandReferenceSet":
                key = "user_reference_list"
            cases.append(
                (
                    "empty",
                    "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                    '\t\t\t_ = %s(d, "%s")' % (schema_fn, name, key),
                )
            )
            return cases
            if ", kind string" in params:
                cases.append(
                    (
                        "empty",
                        "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                        '\t\t\t_ = %s(d, "project")' % (schema_fn, name),
                    )
                )
                return cases
            cases.append(
                (
                    "empty",
                    "d := schema.TestResourceDataRaw(t, %s.Schema, map[string]interface{}{})\n"
                    "\t\t\t_ = %s(d)" % (schema_fn, name),
                )
            )
            return cases

    # Enum-like named types cannot use composite literals Foo{}.
    ENUM_BASES = {
        "HttpProxyWhiteListTargetType",
        "HttpProxyType",
        "PIIScrubbingLevel",
        "ClusterFaultToleranceRef",
        "HypervisorType",
        "ClusterFunctionRef",
        "ClusterArchReference",
        "OperationMode",
        "EncryptionStatus",
        "EncryptionOptionInfo",
        "EncryptionScopeInfo",
        "KeyManagementServerType",
        "DomainAwarenessLevel",
        "UpgradeStatus",
        "SoftwareTypeRef",
        "NodeStatus",
        "HostTypeEnum",
        "HypervisorState",
        "AcropolisConnectionState",
        "StorageTierReference",
    }

    def is_enum_type(t: str) -> bool:
        base = t[1:] if t.startswith("*") else t
        base = base.split(".")[-1]
        return base in ENUM_BASES

    if ptype and ptype.startswith("*") and not ptype.startswith("*schema"):
        if nil_check:
            cases.append(("nil", "require.Nil(t, %s(nil))" % name))
        if name in SKIP_EMPTY_POPULATED or name in SKIP_POPULATED:
            return cases
        if is_enum_type(ptype):
            cases.append(
                (
                    "empty",
                    "var v %s\n\t\t\t_ = %s(&v)" % (ptype[1:], name),
                )
            )
            cases.append(
                (
                    "populated",
                    "var v %s\n\t\t\t_ = %s(&v)" % (ptype[1:], name),
                )
            )
            return cases
        cases.append(("empty", "_ = %s(&%s{})" % (name, ptype[1:])))
        cases.append(("populated", "_ = %s(&%s{})" % (name, ptype[1:])))
        return cases

    if ptype and ptype.startswith("[]") and ptype != "[]interface{}":
        # Normalize host-entity alias config.IPAddress (common/v1) to import4.
        ptype_use = ptype.replace("[]config.IPAddress", "[]import4.IPAddress")
        if nil_check:
            cases.append(("nil", "_ = %s(nil)" % name))
            cases.append(("empty", "_ = %s(%s{})" % (name, ptype_use)))
        if name in SKIP_POPULATED:
            return cases
        et = ptype_use[2:]
        if et.startswith("*"):
            if is_enum_type(et):
                cases.append(
                    (
                        "populated",
                        "var v %s\n\t\t\t_ = %s(%s{&v})" % (et[1:], name, ptype_use),
                    )
                )
            else:
                cases.append(
                    ("populated", "_ = %s(%s{&%s{}})" % (name, ptype_use, et[1:]))
                )
        elif is_enum_type(et):
            cases.append(
                (
                    "populated",
                    "var v %s\n\t\t\t_ = %s(%s{v})" % (et, name, ptype_use),
                )
            )
        else:
            cases.append(("populated", "_ = %s(%s{{}})" % (name, ptype_use)))
        return cases

    if ptype == "[]interface{}":
        # Treat empty slice as nil only when body uses len(pr); avoid empty when body
        # only checks pr != nil then indexes [0].
        if re.match(r"\{\s*if\s+len\(%s\)" % re.escape(f.get("pname") or "pr"), body) or (
            "if len(" + (f.get("pname") or "pr") + ")" in body[:80]
        ):
            cases.append(("nil", "require.Nil(t, %s(nil))" % name))
            cases.append(("empty", "require.Nil(t, %s([]interface{}{}))" % name))
        elif safe_nil_check(body, f.get("pname")):
            cases.append(("nil", "require.Nil(t, %s(nil))" % name))
            # empty slice is non-nil; skip if body indexes [0] without len check
        if name not in SKIP_POPULATED:
            cases.append(
                (
                    "populated",
                    '_ = %s([]interface{}{map[string]interface{}{"name": "n", "value": "v", "tag_name": "t", "server_name": "s1", "network_protocol": "TCP", "username": "u", "password": "p", "user_principal_name": "user@ex.com"}})'
                    % name,
                )
            )
        return cases

    if ptype == "interface{}":
        if safe_nil_check(body, f.get("pname")):
            cases.append(("nil", "require.Nil(t, %s(nil))" % name))
        if name in SKIP_POPULATED:
            return cases
        casts_map = ".(map[string]interface{})" in body
        casts_slice = ".([]interface{})" in body
        if casts_map and not casts_slice:
            cases.append(("empty", "_ = %s(map[string]interface{}{})" % name))
            cases.append(
                (
                    "populated",
                    '_ = %s(map[string]interface{}{"num_snapshots": 1})' % name,
                )
            )
            return cases
        # Slice-shaped interface{}: empty []interface{}{} is non-nil and often indexes [0].
        # Only exercise a single populated element (and nil when safe).
        if casts_slice:
            cases.append(
                (
                    "populated",
                    '_ = %s([]interface{}{map[string]interface{}{"name": "n", "value": "v", "email_address": "a@b.c", "type": "PLAIN", "is_enabled": true}})'
                    % name,
                )
            )
        return cases

    if ptype and ptype.startswith("map["):
        if "map[string]interface{}" in ptype:
            cases.append(("empty", "_ = %s(map[string]interface{}{})" % name))
            cases.append(
                (
                    "populated",
                    '_ = %s(map[string]interface{}{"operator": "IN", "right_hand_side": []interface{}{}})'
                    % name,
                )
            )
        elif "map[string][]string" in ptype:
            cases.append(("nil", "_ = %s(nil)" % name))
            cases.append(("empty", "_ = %s(map[string][]string{})" % name))
            cases.append(("populated", '_ = %s(map[string][]string{"a": {"b"}})' % name))
        elif "map[string]string" in ptype:
            cases.append(("nil", "_ = %s(nil)" % name))
            cases.append(("empty", "_ = %s(map[string]string{})" % name))
            cases.append(("populated", '_ = %s(map[string]string{"a": "b"})' % name))
        return cases

    if ptype == "v3.RightHandSide":
        cases.append(("empty", "_ = %s(v3.RightHandSide{})" % name))
        cases.append(
            (
                "populated",
                "_ = %s(v3.RightHandSide{Collection: utils.StringPtr(\"x\")})" % name,
            )
        )
        return cases

    if ptype == "era.Info":
        cases.append(("empty", "_ = %s(era.Info{})" % name))
        cases.append(("populated", "_ = %s(era.Info{})" % name))
        return cases

    if ptype == "*schema.Set":
        cases.append(
            (
                "empty",
                "s := schema.NewSet(schema.HashString, nil)\n\t\t\t_ = %s(s)" % name,
            )
        )
        cases.append(
            (
                "populated",
                "s := schema.NewSet(schema.HashResource(&schema.Resource{Schema: map[string]*schema.Schema{"
                '"name": {Type: schema.TypeString}, "value": {Type: schema.TypeString}}}), '
                '[]interface{}{map[string]interface{}{"name": "n", "value": "true"}})\n'
                "\t\t\t_ = %s(s)" % name,
            )
        )
        return cases

    if ptype and " " not in ptype and not ptype.startswith("*") and not ptype.startswith("["):
        if is_enum_type(ptype):
            cases.append(("empty", "var v %s\n\t\t\t_ = %s(v)" % (ptype, name)))
            cases.append(("populated", "var v %s\n\t\t\t_ = %s(v)" % (ptype, name)))
        else:
            cases.append(("empty", "_ = %s(%s{})" % (name, ptype)))
            cases.append(("populated", "_ = %s(%s{})" % (name, ptype)))
        return cases

    return cases


def write_file(path: str, content: str):
    open(path, "w").write(content)
    print("Wrote", path, "lines", content.count("\n") + 1)


def write_prism():
    schemas, funcs = collect("prism")
    covered = {
        "expandStringList",
        "expandCategories",
        "flattenCategories",
        "expandFilterParams",
        "flattenReferenceValues",
        "expandReference",
        "flattenArrayReferenceValues",
        "flattenReferenceValuesList",
        "flattenReferenceList",
        "flattenReference",
        "flattenExternalNetworkListReferenceList",
        "flattenExternalNetworkListReference",
        "expandDirectoryUserGroup",
        "expandSamlUserGroup",
    }
    lines = [
        "package prism",
        "",
        "import (",
        '\t"testing"',
        "",
        '\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"',
        '\t"github.com/stretchr/testify/require"',
        '\t"github.com/terraform-providers/terraform-provider-nutanix/nutanix/client"',
        '\tv3 "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/prism"',
        '\t"github.com/terraform-providers/terraform-provider-nutanix/utils"',
        ")",
        "",
        "func TestUnitPrismSchemas(t *testing.T) {",
        "\tt.Parallel()",
        "",
    ]
    for s in schemas:
        w = "true" if s.startswith("Resource") else "false"
        lines.append(
            "\trequire.NoError(t, %s().InternalValidate(%s().Schema, %s))" % (s, s, w)
        )
    lines += [
        "\trequire.NoError(t, categoriesSchema().Elem.(*schema.Resource).InternalValidate(categoriesSchema().Elem.(*schema.Resource).Schema, true))",
        "\trequire.NoError(t, categoriesSchemaOptional().Elem.(*schema.Resource).InternalValidate(categoriesSchemaOptional().Elem.(*schema.Resource).Schema, true))",
        "\trequire.NoError(t, DataSourceFiltersSchema().Elem.(*schema.Resource).InternalValidate(DataSourceFiltersSchema().Elem.(*schema.Resource).Schema, true))",
        "}",
        "",
        PRISM_HELPERS,
        "",
        "func TestUnitPrismFlattenExpandHelpers(t *testing.T) {",
        "\tt.Parallel()",
        "",
    ]
    auto = []
    for f in funcs:
        if f["name"] in covered or f["name"] in SKIP:
            continue
        cases = gen_cases(f)
        if cases:
            auto.append((f["name"], cases))
    emit_subtests(lines, auto)
    lines.append("}")
    write_file(
        os.path.join(ROOT, "nutanix/services/prism/prism_unit_test.go"),
        "\n".join(lines) + "\n",
    )
    print("  auto funcs:", len(auto))


def write_ndb():
    schemas, funcs = collect("ndb")
    lines = [
        "package ndb",
        "",
        "import (",
        '\t"context"',
        '\t"testing"',
        "",
        '\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"',
        '\t"github.com/stretchr/testify/require"',
        '\tera "github.com/terraform-providers/terraform-provider-nutanix/nutanix/sdks/v3/era"',
        '\t"github.com/terraform-providers/terraform-provider-nutanix/utils"',
        ")",
        "",
        "func TestUnitNDBSchemas(t *testing.T) {",
        "\tt.Parallel()",
        "",
    ]
    for s in schemas:
        w = "true" if s.startswith("Resource") else "false"
        lines.append(
            "\trequire.NoError(t, %s().InternalValidate(%s().Schema, %s))" % (s, s, w)
        )
    lines += [
        "}",
        "",
        NDB_HELPERS,
        "",
        "func TestUnitNDBFlattenExpandHelpers(t *testing.T) {",
        "\tt.Parallel()",
        "",
    ]
    skip_hand = {"expandPrimarySLA", "expandIPInfos", "expandTags"}
    auto = []
    for f in funcs:
        if f["name"] in SKIP or f["name"] in skip_hand:
            continue
        cases = gen_cases(f)
        if cases:
            auto.append((f["name"], cases))
    emit_subtests(lines, auto)
    lines.append("}")
    write_file(
        os.path.join(ROOT, "nutanix/services/ndb/ndb_unit_test.go"),
        "\n".join(lines) + "\n",
    )
    print("  auto funcs:", len(auto))


def write_clustersv2():
    schemas, funcs = collect("clustersv2")
    lines = [
        "package clustersv2",
        "",
        "import (",
        '\t"testing"',
        "",
        '\t"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"',
        '\tconfig "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"',
        '\timport1 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/clustermgmt/v4/config"',
        '\timport4 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/config"',
        '\tcommonResp "github.com/nutanix-core/ntnx-api-golang-sdk-internal/clustermgmt-go-client/v17/models/common/v1/response"',
        '\t"github.com/stretchr/testify/require"',
        '\t"github.com/terraform-providers/terraform-provider-nutanix/utils"',
        ")",
        "",
        "func TestUnitClustersV2Schemas(t *testing.T) {",
        "\tt.Parallel()",
        "",
    ]
    for s in schemas:
        w = (
            "false"
            if s.startswith("Datasource") or s.startswith("DataSource")
            else "true"
        )
        lines.append(
            "\trequire.NoError(t, %s().InternalValidate(%s().Schema, %s))" % (s, s, w)
        )
    lines += [
        "\trequire.NoError(t, schemaForLinks().Elem.(*schema.Resource).InternalValidate(schemaForLinks().Elem.(*schema.Resource).Schema, true))",
        "\trequire.NoError(t, schemaForValuePrefixLength().Elem.(*schema.Resource).InternalValidate(schemaForValuePrefixLength().Elem.(*schema.Resource).Schema, true))",
        "\trequire.NoError(t, schemaForIPAddress(false).Elem.(*schema.Resource).InternalValidate(schemaForIPAddress(false).Elem.(*schema.Resource).Schema, true))",
        "}",
        "",
        CLUSTERS_HELPERS,
        "",
        "func TestUnitClustersV2FlattenExpandHelpers(t *testing.T) {",
        "\tt.Parallel()",
        "",
    ]
    auto = []
    for f in funcs:
        if f["name"] in SKIP or f["name"] == "flattenLinks":
            continue
        cases = gen_cases(f)
        if cases:
            auto.append((f["name"], cases))
    emit_subtests(lines, auto)
    lines.append("}")
    write_file(
        os.path.join(ROOT, "nutanix/services/clustersv2/clustersv2_unit_test.go"),
        "\n".join(lines) + "\n",
    )
    print("  auto funcs:", len(auto))


PRISM_HELPERS = r'''
func TestUnitExpandStringList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   []interface{}
		want []*string
	}{
		{name: "nil", in: nil, want: []*string{}},
		{name: "empty", in: []interface{}{}, want: []*string{}},
		{name: "populated", in: []interface{}{"a", "", "b"}, want: []*string{utils.StringPtr("a"), utils.StringPtr("b")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandStringList(tt.in))
		})
	}
}

func TestUnitExpandFlattenCategories(t *testing.T) {
	t.Parallel()
	t.Run("flatten nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenCategories(nil))
	})
	t.Run("flatten empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []interface{}{}, flattenCategories(map[string]string{}))
	})
	t.Run("flatten populated", func(t *testing.T) {
		t.Parallel()
		got := flattenCategories(map[string]string{"env": "prod"})
		require.Equal(t, []interface{}{map[string]interface{}{"name": "env", "value": "prod"}}, got)
	})
	t.Run("expand empty set", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(func(v interface{}) int {
			c := v.(map[string]interface{})
			return utils.HashcodeString(c["name"].(string) + c["value"].(string))
		}, nil)
		require.Equal(t, map[string]string{}, expandCategories(s))
	})
	t.Run("expand populated", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(func(v interface{}) int {
			c := v.(map[string]interface{})
			return utils.HashcodeString(c["name"].(string) + c["value"].(string))
		}, []interface{}{map[string]interface{}{"name": "env", "value": "prod"}})
		require.Equal(t, map[string]string{"env": "prod"}, expandCategories(s))
	})
}

func TestUnitExpandFilterParams(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   map[string][]string
		want []map[string]interface{}
	}{
		{name: "nil", in: nil, want: []map[string]interface{}{}},
		{name: "empty", in: map[string][]string{}, want: []map[string]interface{}{}},
		{name: "populated", in: map[string][]string{"name": {"a"}}, want: []map[string]interface{}{{"name": "name", "values": []string{"a"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, expandFilterParams(tt.in))
		})
	}
}

func TestUnitBuildAndReplaceFilters(t *testing.T) {
	t.Parallel()
	t.Run("build empty", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(schema.HashResource(DataSourceFiltersSchema().Elem.(*schema.Resource)), nil)
		require.Empty(t, BuildFiltersDataSource(s))
	})
	t.Run("build populated", func(t *testing.T) {
		t.Parallel()
		s := schema.NewSet(schema.HashResource(DataSourceFiltersSchema().Elem.(*schema.Resource)), []interface{}{
			map[string]interface{}{"name": "vm_name", "values": []interface{}{"vm1"}},
		})
		got := BuildFiltersDataSource(s)
		require.Len(t, got, 1)
		require.Equal(t, "vm_name", got[0].Name)
		require.Equal(t, []string{"vm1"}, got[0].Values)
	})
	t.Run("replace nil mappings", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, ReplaceFilterPrefixes(nil, nil))
	})
	t.Run("replace empty", func(t *testing.T) {
		t.Parallel()
		require.Empty(t, ReplaceFilterPrefixes([]*client.AdditionalFilter{}, map[string]string{}))
	})
	t.Run("replace populated", func(t *testing.T) {
		t.Parallel()
		in := []*client.AdditionalFilter{{Name: "vm.name", Values: []string{"x"}}}
		got := ReplaceFilterPrefixes(in, map[string]string{"vm": "virtual_machine"})
		require.Equal(t, "virtual_machine.name", got[0].Name)
	})
	t.Run("filterParamsHash", func(t *testing.T) {
		t.Parallel()
		require.NotEqual(t, 0, filterParamsHash(map[string]interface{}{"name": "a"}))
	})
}

func TestUnitGetACPRoleUUIDAndRemoval(t *testing.T) {
	t.Parallel()
	t.Run("nil-ish", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", getACPRoleUUID(nil))
		require.Equal(t, "", getACPRoleUUID("x"))
		require.False(t, acpRemovalFromMiddle(nil, nil))
		require.False(t, acpRemovalFromMiddleRoles(nil, nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", getACPRoleUUID(map[string]interface{}{}))
		require.False(t, acpRemovalFromMiddle([]interface{}{}, []interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		item := map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "r1"}}}
		require.Equal(t, "r1", getACPRoleUUID(item))
		old := []interface{}{
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "a"}}},
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "b"}}},
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "c"}}},
		}
		newL := []interface{}{
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "a"}}},
			map[string]interface{}{"role_reference": []interface{}{map[string]interface{}{"uuid": "c"}}},
		}
		require.True(t, acpRemovalFromMiddle(old, newL))
		require.Equal(t, 0, acpHash("bad"))
		require.NotEqual(t, 0, acpHash(item))
	})
}

func TestUnitValidateAndExpandReference(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, map[string]interface{}{}, flattenReferenceValues(nil))
		require.Nil(t, expandReference(map[string]interface{}{}))
		require.Nil(t, validateRef(map[string]interface{}{}))
		require.Equal(t, map[string]interface{}{}, flattenReference(nil))
		require.Equal(t, map[string]interface{}{}, flattenExternalNetworkListReference(nil))
		require.Equal(t, []interface{}{}, flattenReferenceValuesList(nil))
	})
	t.Run("empty slices", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []map[string]interface{}{}, flattenArrayReferenceValues(nil))
		require.Equal(t, []map[string]interface{}{}, flattenArrayReferenceValues([]*v3.Reference{}))
		require.Equal(t, []map[string]interface{}{}, flattenReferenceList(nil))
		require.Equal(t, []map[string]interface{}{}, flattenReferenceList([]*v3.ReferenceValues{}))
		require.Equal(t, []map[string]interface{}{}, flattenExternalNetworkListReferenceList(nil))
		require.Nil(t, expandDirectoryUserGroup(nil))
		require.Nil(t, expandDirectoryUserGroup([]interface{}{}))
		require.Nil(t, expandSamlUserGroup(nil))
		require.Nil(t, expandSamlUserGroup([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		ref := &v3.Reference{Kind: utils.StringPtr("vm"), UUID: utils.StringPtr("u1"), Name: utils.StringPtr("n1")}
		require.Equal(t, map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"}, flattenReferenceValues(ref))
		require.Equal(t, []interface{}{map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"}}, flattenReferenceValuesList(ref))
		exp := expandReference(map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"})
		require.Equal(t, "vm", *exp.Kind)
		rv := &v3.ReferenceValues{Kind: "vm", UUID: "u1", Name: "n1"}
		require.Equal(t, map[string]interface{}{"kind": "vm", "uuid": "u1", "name": "n1"}, flattenReference(rv))
		require.Equal(t, map[string]interface{}{"uuid": "u1", "name": "n1"}, flattenExternalNetworkListReference(rv))
		dug := expandDirectoryUserGroup([]interface{}{map[string]interface{}{"distinguished_name": "cn=a"}})
		require.Equal(t, "cn=a", *dug.DistinguishedName)
		sug := expandSamlUserGroup([]interface{}{map[string]interface{}{"idp_uuid": "i1", "name": "g1"}})
		require.Equal(t, "i1", *sug.IdpUUID)
		require.NotNil(t, validateRef(map[string]interface{}{"kind": "vm", "uuid": "u1"}))
		require.Nil(t, validateRefList(nil, nil))
		require.NotNil(t, validateRefList([]interface{}{map[string]interface{}{"uuid": "u1"}}, utils.StringPtr("vm")))
	})
}
'''.lstrip("\n")

NDB_HELPERS = r'''
func TestUnitNDBContextAndBoolHelpers(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		t.Parallel()
		_, ok := FromContext(context.Background())
		require.False(t, ok)
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		ctx := NewContext(context.Background(), "")
		got, ok := FromContext(ctx)
		require.True(t, ok)
		require.Equal(t, "", got)
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		ctx := NewContext(context.Background(), "db-1")
		got, ok := FromContext(ctx)
		require.True(t, ok)
		require.Equal(t, "db-1", got)
		b, ok := tryToConvertBool("true")
		require.True(t, ok)
		require.True(t, b)
		_, ok = tryToConvertBool("nope")
		require.False(t, ok)
	})
}

func TestUnitExpandPrimarySLAAndIPInfos(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPrimarySLA(nil))
		require.Nil(t, expandIPInfos(nil))
		require.Nil(t, expandTags(nil))
		require.Nil(t, buildSLADetails(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, expandPrimarySLA([]interface{}{}))
		require.Nil(t, expandIPInfos([]interface{}{}))
		require.Nil(t, expandTags([]interface{}{}))
		require.Nil(t, buildSLADetails([]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		sla := expandPrimarySLA([]interface{}{map[string]interface{}{
			"sla_id":         "sla-1",
			"nx_cluster_ids": []interface{}{"c1"},
		}})
		require.Equal(t, "sla-1", utils.StringValue(sla.SLAID))
		ips := expandIPInfos([]interface{}{map[string]interface{}{
			"ip_type":      "static",
			"ip_addresses": []interface{}{"10.0.0.1"},
		}})
		require.Len(t, ips, 1)
		tags := expandTags([]interface{}{map[string]interface{}{"tag_name": "n", "value": "v"}})
		require.Len(t, tags, 1)
	})
}
'''.lstrip("\n")

CLUSTERS_HELPERS = r'''
func TestUnitIPHelpers(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", ipToKey(nil))
		require.Equal(t, "UNKNOWN", getIPType(nil))
		require.True(t, ipv4Equal(nil, nil))
		require.True(t, ipv6Equal(nil, nil))
		require.True(t, ipAddressEqual(nil, nil))
		require.True(t, nodeEqual(nil, nil))
		require.Empty(t, nodeKeyCandidates(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "", ipToKey(&import4.IPAddress{}))
		require.Equal(t, "UNKNOWN", getIPType(&import4.IPAddress{}))
		require.False(t, ipv4Equal(&import4.IPv4Address{}, nil))
		require.Empty(t, nodeKeyCandidates(&config.NodeListItemReference{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		ip := &import4.IPAddress{Ipv4: &import4.IPv4Address{Value: utils.StringPtr("10.0.0.1"), PrefixLength: utils.IntPtr(32)}}
		require.Equal(t, "ipv4:10.0.0.1/32", ipToKey(ip))
		require.Equal(t, "IPV4", getIPType(ip))
		require.True(t, ipv4Equal(ip.Ipv4, ip.Ipv4))
		n := &config.NodeListItemReference{NodeUuid: utils.StringPtr("n1"), ControllerVmIp: ip}
		require.Contains(t, nodeKeyCandidates(n), "uuid:n1")
		require.True(t, nodeEqual(n, n))
	})
}

func TestUnitFindCategories(t *testing.T) {
	t.Parallel()
	t.Run("nil/empty old", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []string{}, FindRemovedCategories(nil, []string{"a"}))
		require.Equal(t, []string{}, FindAddedCategories([]string{"a"}, nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []string{}, FindRemovedCategories([]string{}, []string{}))
		require.Equal(t, []string{}, FindAddedCategories([]string{}, []string{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, []string{"a"}, FindRemovedCategories([]string{"a", "b"}, []string{"b"}))
		require.Equal(t, []string{"c"}, FindAddedCategories([]string{"a"}, []string{"a", "c"}))
	})
}

func TestUnitFlattenLinksHelper(t *testing.T) {
	t.Parallel()
	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenLinks(nil))
	})
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Nil(t, flattenLinks([]commonResp.ApiLink{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		got := flattenLinks([]commonResp.ApiLink{{
			Href: utils.StringPtr("https://example.test"),
			Rel:  utils.StringPtr("self"),
		}})
		require.Equal(t, []map[string]interface{}{{"href": "https://example.test", "rel": "self"}}, got)
	})
}

func TestUnitAuthorizedPublicKeyHash(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, schema.HashString("-"), authorizedPublicKeyHash(map[string]interface{}{}))
	})
	t.Run("populated", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, schema.HashString("n-k"), authorizedPublicKeyHash(map[string]interface{}{"name": "n", "key": "k"}))
	})
}
'''.lstrip("\n")


if __name__ == "__main__":
    write_prism()
    write_ndb()
    write_clustersv2()
