# AI-Assisted Terraform Provider Code Generation – VP Demo Script

**Presenter:** Nins Bisht
**Duration:** ~5–7 minutes
**Audience:** VPs / Leadership
**Goal:** Show how Cursor + our code-gen pipeline produces production-ready Terraform provider code for any new Nutanix API entity in a single request.

---

## 1. Opening (30 sec)

> "Good morning. I'm Nins Bisht. In the next few minutes I'll show how we use **AI-assisted code generation** to build the complete Terraform provider support for a new Nutanix entity – schema, resource, data sources, tests, docs and examples – end-to-end.
>
> For today's demo, the entity is **Scope Template** from the IAM namespace."

---

## 2. The Problem We're Solving (30 sec)

> "Every new API entity needs the same set of deliverables – SDK wiring, GetById and List data sources, full CRUD resource, flatten/expand helpers, acceptance tests, docs, examples and provider registration.
>
> Manually this is **2–3 days of repetitive, error-prone work**. We've automated 70–80% of it."

---

## 3. The Pipeline – 3 Steps (45 sec)

> "The pipeline has three steps:
>
> 1. **Extract SDK info** – a Go tool that reads the SDK module cache and produces a JSON describing every API method, request and response struct.
> 2. **Generate code** – a Python driver that feeds that JSON into Cursor with our codified rules and repo conventions.
> 3. **Developer review** – we review, build, run `terraform apply`, and ship."

---

## 4. Live Demo (3–4 min)

### Step 1 – Extract the SDK info (~20 sec)

> "I'm running the extractor on the IAM SDK package, scoped to the `ScopeTemplate` keyword."

```bash
go run code_gen/extract_sdk_info.go \
  -package=".../iam-go-client/v17" \
  -keyword="ScopeTemplate"
```

> "It produced `sdk_info.json`. Scope Template exposes three APIs – `CreateScopeTemplate`, `GetScopeTemplate`, `ListScopeTemplates` – and you can see all three captured here, with full request and response detail."

### Step 2 – Trigger code generation (~30 sec)

```bash
cd code_gen && python3 cursor_auto_generate.py sdk_extract_output/sdk_info.json
```

> "This invokes Cursor with the **claude-4.6-opus-high** agent. It already understands our framework, schema patterns and repo layout. Watch – it lays out the task list, then generates every file we need."

*(Wait while generation runs – ~1–2 minutes.)*

> "Done. In a single request it generated:
> - the **data source** and **resource** Go files
> - the **acceptance test** files
> - the **HCL examples**
> - the **documentation**
> - and registered both the resource and data source in `provider.go`."

### Step 3 – Review and validate (~2 min)

> "Now let's review what was produced. Open the list data source for Scope Templates."

> "Notice it added the standard filter, order-by, and limit options, and the `scope_templates` attribute with `display_name`, `description`, `entities` (a list with `entity_filter`), `created_by` and `created_time` – five attributes total."

> "Cross-checking the API documentation – tenant ID, external ID, display name, entities, created_by, created_time – every field is in the schema. Nothing missed."

> "Now let's compile the provider:"

```bash
make build
```

> "Build is successful."

> "I'll grab the registered data source name from `provider.go` and run a quick `terraform plan` and `apply` against PC `10.x.x.x` which has one Scope Template – `project-scope-template`, scoped to the `project` entity."

```bash
terraform apply
```

> "Apply succeeds. Looking at the state file – display name `project-scope-template`, description `scope template for projects 2.0` – matches the API response exactly."

---

## 5. The Numbers (45 sec)

> "On the **Role Membership** entity, this same pipeline produced:
> - **17 files**
> - **~1,475 lines** of production-ready code
> - in **a single Cursor request** (~1.2M tokens, claude-4.6-opus-high)
>
> Manual effort: **12–18 hours**. AI-assisted: **2–4 hours**. That's a **~10x speedup** – same-day delivery instead of multi-day effort.
>
> PR is live: github.com/nutanix/terraform-provider-nutanix/pull/1128."

---

## 6. Why It Matters to the Business (30 sec)

> "Three concrete wins:
> 1. **Speed** – ~10x faster per entity, faster time-to-market.
> 2. **Consistency** – every resource follows the exact same patterns; reviewers focus on business logic, not boilerplate.
> 3. **Completeness** – tests, docs, examples and provider registration ship together, every time.
>
> And this runs on our existing Cursor Business license – **zero incremental cost**."

---

## 7. What's Next (20 sec)

> "Two next steps:
> 1. Move the pipeline out of local-only into a shared, repo-hosted workflow.
> 2. Package it as a reusable **Cursor Agent skill** so any team member can invoke it on any new entity.
>
> The same approach already applies to **Ansible** with comparable value.
>
> Thank you."

---

## Appendix – Quick FAQ Cheat Sheet

| Question | Answer |
|---|---|
| Is it fully autonomous? | No – ~70–80% automated. Developers still review, dev-test and finalize tests. |
| What about new SDKs? | Pipeline reads any Go SDK in the module cache – no per-SDK wiring. |
| Tooling cost? | One-time ~1,300 lines of Python + Go. Already complete and operational. |
| Model used? | claude-4.6-opus-high inside Cursor. |
| Token cost? | ~1.2M tokens per entity, covered by existing Cursor Business license. |
