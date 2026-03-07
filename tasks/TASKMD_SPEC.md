# taskmd Specification

**Version:** 1.2
**Last Updated:** 2026-02-15

> **See also:** [Operations Specification](./taskmd_operations.md) — defines behavioral contracts for scanning, filtering, validation, ranking, dependency resolution, graph construction, and search.

## Quick Reference

Each task is a `.md` file with YAML frontmatter and a markdown body.

```yaml
---
id: "001"
title: "Task title"
status: pending
---

# Task Title

Description and subtasks go here.
```

### Field Summary

| Field | Type | Required | Values / Format |
|-------|------|----------|-----------------|
| `id` | string | **Yes** | Unique identifier (e.g., `"001"`, `"42"`, `"cli-049"`) |
| `title` | string | **Yes** | Brief, descriptive text |
| `status` | enum | Recommended | `pending`, `in-progress`, `completed`, `in-review`, `blocked`, `cancelled` |
| `priority` | enum | No | `low`, `medium`, `high`, `critical` |
| `effort` | enum | No | `small`, `medium`, `large` |
| `type` | enum | No | `feature`, `bug`, `improvement`, `chore`, `docs` |
| `dependencies` | array | No | List of task ID strings (e.g., `["001", "015"]`) |
| `tags` | array | No | Lowercase, hyphen-separated strings |
| `group` | string | No | Logical grouping (derived from directory if omitted) |
| `owner` | string | No | Free-form assignee name or identifier |
| `touches` | array | No | Abstract scope identifiers (e.g., `["cli/graph", "cli/output"]`) |
| `context` | array | No | Explicit file paths relevant to the task (e.g., `["docs/api.md"]`) |
| `parent` | string | No | Single task ID (e.g., `"045"`) |
| `created` | date | No | `YYYY-MM-DD` |
| `verify` | array | No | List of typed verification checks (see below) |
| `pr` | array | No | List of pull request URLs |
| `external_id` | string | No | Identifier from an external system (e.g., `"PROJ-123"`, `"42"`) |

## Frontmatter Schema

<!-- Unknown frontmatter fields are silently ignored by the parser and preserved as-is in the file. -->

### Required Fields

**`id`** — Unique identifier for the task. Any non-empty string is valid (e.g., `"001"`, `"42"`, `"cli-049"`). Must be unique across all tasks in the project.

**`title`** — Brief, action-oriented description of the task.

### Optional Fields

**`status`** — Current state of the task (recommended for all tasks):

| Status | Meaning |
|--------|---------|
| `pending` | Not started (initial state) |
| `in-progress` | Currently being worked on |
| `in-review` | Submitted for review (PR open) |
| `completed` | Finished and verified |
| `blocked` | Cannot proceed due to a blocker |
| `cancelled` | Will not be completed |

```
pending → in-progress → in-review → completed
   ↓            ↓            ↓            ↓
   ↓         blocked         ↓            ↓
   ↓            ↓            ↓            ↓
   └──→ cancelled ←──────────┴────────────┘
```

**`priority`** — Importance level:

| Priority | Use Case |
|----------|----------|
| `low` | Nice to have, can be deferred |
| `medium` | Standard work items (default) |
| `high` | Important for project success |
| `critical` | Urgent, must address immediately |

**`effort`** — Estimated complexity:

| Effort | Typical Duration |
|--------|------------------|
| `small` | < 2 hours |
| `medium` | 2–8 hours |
| `large` | > 8 hours / multi-day |

**`type`** — Classification of the work item:

| Type | Meaning |
|------|---------|
| `feature` | New functionality |
| `bug` | Defect fix |
| `improvement` | Enhancement to existing functionality |
| `chore` | Maintenance or housekeeping |
| `docs` | Documentation-only change |

**`dependencies`** — List of task IDs that must be completed before this task can start. Always reference by ID, always use array format:

```yaml
dependencies: ["001", "015"]
```

**`tags`** — Labels for categorization and filtering. Use lowercase, hyphen-separated strings:

```yaml
tags:
  - core
  - api
```

**`group`** — Logical grouping. If omitted, derived from the parent directory name. Root-level tasks have no group.

**`owner`** — Free-form string for assigning a task to a person or team. Used for filtering and display; no validation is applied.

**`touches`** — List of abstract scope identifiers declaring which code areas a task modifies. Used by the `tracks` command to detect spatial overlap and assign tasks to parallel work tracks. Two tasks that share a scope should not be worked on simultaneously (risk of merge conflicts).

```yaml
touches:
  - cli/graph
  - cli/output
```

Scopes are user-defined identifiers. Concrete scope-to-path mappings can be configured in `.taskmd.yaml`:

```yaml
# .taskmd.yaml
scopes:
  cli/graph:
    description: "Graph visualization and dependency rendering"
    paths:
      - "apps/cli/internal/graph/"
      - "apps/cli/internal/cli/graph.go"
  cli/output:
    paths:
      - "apps/cli/internal/cli/format.go"
```

The optional `description` field provides a human-readable explanation of what a scope covers. When present, it is included in validation error messages for easier debugging.

When scopes are configured, `touches` values not found in the config produce a warning. When no scopes config exists, all values are accepted silently.

**`context`** — List of explicit file paths relevant to the task. Use this for files that fall outside scope mappings, such as test fixtures, documentation, or configuration files. Paths are relative to the project root.

```yaml
context:
  - "docs/api-design.md"
  - "apps/cli/internal/web/handlers_test.go"
```

The `context` command merges files from both `touches` (via scope resolution) and `context` (explicit paths), deduplicating by path. Each entry is tagged with its source (`scope:<name>` or `explicit`) and checked for existence. Non-existent files are not errors — the task may create them.

**`parent`** — Task ID of the parent task for hierarchical grouping. A task can have at most one parent. Children are computed dynamically (not stored in frontmatter) by finding all tasks whose `parent` matches a given ID.

- Purely organizational — does not imply blocking or dependency
- No status cascading — completing all children does not auto-complete the parent
- Must reference an existing task ID; self-references and cycles are flagged by validation

```yaml
parent: "045"
```

**`created`** — Date when the task was created, in `YYYY-MM-DD` format.

**`verify`** — List of typed acceptance checks for validating task completion. Each entry is a map with a `type` field that determines the check kind. Run checks with `taskmd verify <ID>`.

| Type | Fields | Behavior |
|------|--------|----------|
| `bash` | `run` (required), `dir` (optional) | Runs `run` in a shell subprocess; pass if exit code 0, fail otherwise |
| `assert` | `check` (required) | Displays `check` text for an agent to evaluate (not executed) |

- `dir` is relative to the project root (where `.taskmd.yaml` lives); defaults to `.`
- Unknown types are preserved in the file but produce a warning and are skipped during execution

```yaml
verify:
  - type: bash
    run: "go test ./internal/api/... -run TestPagination"
    dir: "apps/cli"
  - type: bash
    run: "npm test"
    dir: "apps/web"
  - type: assert
    check: "Pagination links appear in the API response headers"
  - type: assert
    check: "Page size defaults to 20 when not specified"
```

**`pr`** — List of pull request URLs associated with this task. Used in `pr-review` workflow mode to track open PRs. Managed via `taskmd set --add-pr <url>` and `taskmd set --remove-pr <url>`.

```yaml
pr: ["https://github.com/owner/repo/pull/42"]
```

**`external_id`** — Identifier from an external system (e.g., a GitHub issue number or Jira issue key). Used to trace synced tasks back to their source. Written by the sync engine; not typically set manually.

```yaml
external_id: "PROJ-123"
```

## Workflow Modes

The `workflow` key in `.taskmd.yaml` controls how tasks transition to completion.

| Mode | Behavior |
|------|----------|
| `solo` (default) | Agent sets task status directly to `completed` |
| `pr-review` | Agent opens a PR, sets status to `in-review` with `--add-pr`, and stops. Task moves to `completed` on PR merge (via CI or manually). |

```yaml
# .taskmd.yaml
workflow: pr-review
```

When using `--done` in `pr-review` mode, the CLI sets status to `in-review` instead of `completed`.

## ID Generation

The `id` section in `.taskmd.yaml` configures how task IDs are generated.

### Configuration

```yaml
# .taskmd.yaml
id:
  strategy: sequential  # "sequential", "prefixed", "random", or "ulid"
  prefix: ""            # required when strategy is "prefixed"
  length: 6             # ID length (used by random and ulid strategies)
  padding: 3            # zero-padding width (used by sequential strategy)
```

### Strategies

| Strategy | Format | Example Filename | Derived ID |
|----------|--------|-----------------|------------|
| `sequential` (default) | Zero-padded number | `009-add-feature.md` | `009` |
| `prefixed` | Prefix + number | `dr-001-fix-login.md` | `dr-001` |
| `random` | Alphanumeric string | `a3f9x2-slug-title.md` | `a3f9x2` |
| `ulid` | ULID (timestamp + random) | `01h5a3mpk2-fix-bug.md` | `01h5a3mpk2` |

### Defaults

When the `id` section is omitted, the following defaults apply:

| Field | Default |
|-------|---------|
| `strategy` | `sequential` |
| `length` | `6` |
| `padding` | `3` |

### Filename Patterns

The parser automatically derives task IDs from filenames based on these patterns:

- **Sequential**: Filename starts with digits — `001-slug.md` → ID `001`
- **Prefixed**: Lowercase alpha prefix, hyphen, digits — `dr-001-slug.md` → ID `dr-001`
- **Random**: 3-8 lowercase alphanumeric chars with at least one digit — `a3f9x2-slug.md` → ID `a3f9x2`
- **ULID**: Crockford Base32 string (timestamp + random) — `01h5a3mpk2-slug.md` → ID `01h5a3mpk2`

## File Organization

### File Naming

Task files follow this pattern:

```
ID-descriptive-title.md
```

Where `ID` is the task ID and `descriptive-title` is a lowercase hyphen-separated slug. Examples:

- `001-project-scaffolding.md`
- `042-implement-user-auth.md`
- `cli-049-add-graph-command.md`

The ID prefix may be omitted if the `id` field in frontmatter is the sole identifier.

### Directory Structure

Tasks can be organized into subdirectories for grouping:

```
tasks/
├── 001-taskmd-specification.md     # No group
├── web/                             # Group: "web"
│   ├── 001-project-scaffolding.md
│   └── 002-typescript-types.md
└── cli/                             # Group: "cli"
    ├── 015-go-cli-scaffolding.md
    └── 016-task-model-parsing.md
```

Group resolution priority:
1. Explicit `group` in frontmatter
2. Parent directory name
3. No group (root-level tasks)

## Validation

A valid taskmd file **must**:

1. Have YAML frontmatter enclosed in `---` delimiters
2. Include required fields: `id`, `title`
3. Use valid enum values for `status`, `priority`, `effort`, `type`
4. Have unique IDs across the project
5. Reference only existing tasks in `dependencies`
6. Have no circular dependency chains
7. Reference an existing task in `parent` (if set), with no self-reference or parent cycles

A valid taskmd file **should**:

1. Follow the `NNN-task-name.md` naming pattern
2. Include a creation date
3. Have a descriptive markdown body

## Examples

### Minimal Task

```markdown
---
id: "001"
title: "Fix login button alignment"
status: pending
---

# Fix Login Button Alignment

The login button on the homepage is misaligned. Update the CSS to center it.
```

### Full Task

```markdown
---
id: "015"
title: "Implement user authentication"
status: in-progress
priority: high
effort: large
type: feature
dependencies: ["012", "013"]
parent: "012"
tags:
  - auth
  - security
  - api
created: 2026-02-08
---

# Implement User Authentication

## Objective

Add JWT-based authentication to the API.

## Tasks

- [x] Design authentication flow
- [x] Implement JWT signing and verification
- [ ] Create login endpoint
- [ ] Create logout endpoint
- [ ] Add authentication middleware
- [ ] Write integration tests

## Acceptance Criteria

- Users can log in with email and password
- JWT tokens expire after 24 hours
- Protected routes require valid JWT
- All endpoints have > 90% test coverage
```
