---
title: "Adoption"
---

A practical guide for teams migrating to wetwire-honeycomb-go for Honeycomb query management.

---

## Why Adopt wetwire-honeycomb-go

### Key Benefits

| Benefit | Description |
|---------|-------------|
| **Type Safety** | Compile-time validation catches errors before deployment |
| **Version Control** | Queries live in Git with full history, diffs, and reviews |
| **Code Reuse** | Share queries across teams via Go modules |
| **Linting** | Built-in best practices enforcement (WHC rules) |
| **Documentation** | Self-documenting code with Go comments |
| **Consistency** | Standardized query patterns across the organization |
| **Automation** | CI/CD integration for validation and deployment |

### When to Adopt

wetwire-honeycomb-go is ideal for teams that:

- Have more than 10-20 queries to manage
- Need to share queries across multiple services or teams
- Want infrastructure-as-code patterns for observability
- Require audit trails and review processes
- Are comfortable with Go development workflows

---

## Migration Paths

### From Manual Honeycomb UI

#### Step 1: Export Query JSON from Honeycomb

In the Honeycomb UI:
1. Open your query
2. Click the "JSON" tab or export option
3. Copy the query JSON

#### Step 2: Convert to Go Using the Import Command

```bash
wetwire-honeycomb import query.json -o queries/performance.go --name SlowRequests
```

#### Step 3: Validate and Iterate

```bash
wetwire-honeycomb lint ./queries/...
wetwire-honeycomb build --pretty ./queries/...
```

### From JSON/YAML Config Files

#### Batch Convert Using Import

```bash
for f in ./config/*.json; do
  name=$(basename "$f" .json | sed 's/-/_/g' | sed 's/\b\(.\)/\u\1/g')
  wetwire-honeycomb import "$f" -o "queries/$(basename "$f" .json).go" --name "$name"
done
```

---

## Step-by-Step Adoption Plan

### Phase 1: Install and Try with One Query

**Duration:** 1-2 hours

```bash
go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest
mkdir honeycomb-queries
cd honeycomb-queries
wetwire-honeycomb init queries
```

### Phase 2: Import Existing Queries

**Duration:** 1-2 days

1. Identify high-priority queries (dashboards, SLOs, alerts)
2. Export and convert each query
3. Organize by domain
4. Add documentation
5. Lint and fix issues

### Phase 3: Add SLOs and Triggers

**Duration:** 2-3 days

Define SLOs and triggers that reference your queries. See [SLOs](../slos/) and [Triggers](../triggers/) documentation.

### Phase 4: Create Boards and Full Stack

**Duration:** 2-3 days

Build dashboards that combine queries, SLOs, and documentation. See [Boards](../boards/) documentation.

---

## Team Onboarding Checklist

### Prerequisites

- [ ] Go 1.21+ installed
- [ ] wetwire-honeycomb CLI installed
- [ ] Git repository access
- [ ] Honeycomb account with API access

### Documentation Review

- [ ] Read [Quick Start](../quick-start/)
- [ ] Read [CLI](../cli/)
- [ ] Read [Lint Rules](../lint-rules/)
- [ ] Review examples directory

### Hands-On Practice

- [ ] Create a simple query using `wetwire-honeycomb init`
- [ ] Run `wetwire-honeycomb lint` and fix any issues
- [ ] Run `wetwire-honeycomb build` and review JSON output
- [ ] Test generated JSON in Honeycomb UI
- [ ] Import an existing query using `wetwire-honeycomb import`

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Honeycomb Queries

on:
  push:
    branches: [main]
    paths:
      - 'queries/**'
  pull_request:
    paths:
      - 'queries/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install wetwire-honeycomb
        run: go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest

      - name: Lint queries
        run: wetwire-honeycomb lint ./...

      - name: Build all resources
        run: wetwire-honeycomb build -o output.json ./...
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

if ! wetwire-honeycomb lint ./queries/... ./slos/... ./triggers/... ./boards/...; then
    echo "Lint errors found. Please fix before committing."
    exit 1
fi
```

---

## Common Challenges and Solutions

### Team Resistance to Code-Based Queries

**Solutions:**
- Start with power users who are comfortable with code
- Show productivity benefits (version control, review process)
- Use `wetwire-honeycomb import` to reduce manual conversion work
- Pair program on first few queries

### Large Backlog of Existing Queries

**Solutions:**
- Prioritize by usage (start with dashboard queries)
- Migrate incrementally (phase approach)
- Use batch import scripts
- Set a deadline for new queries (new in code, legacy in UI)

### Maintaining Consistency Across Teams

**Solutions:**
- Establish organization-wide conventions
- Use lint rules to enforce patterns
- Create shared packages for common queries
- Regular code reviews for query changes

---

## Success Metrics to Track

### Adoption Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| **Query Coverage** | Percentage of queries managed in wetwire | 80%+ within 6 months |
| **Team Onboarding** | Time for new team member to create first query | < 2 hours |
| **Active Contributors** | Team members with merged query PRs | > 50% of users |

### Quality Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| **Lint Pass Rate** | PRs passing lint on first try | > 90% |
| **Build Success Rate** | Builds completing without errors | > 99% |
| **Query Reuse** | Average number of SLOs/triggers per query | > 1.5 |

---

## See Also

- [Quick Start](../quick-start/) - Getting started guide
- [CLI](../cli/) - Complete command reference
- [SLOs](../slos/) - SLO documentation
- [Triggers](../triggers/) - Trigger documentation
- [Boards](../boards/) - Board documentation
- [Lint Rules](../lint-rules/) - Lint rule reference
- [FAQ](../faq/) - Frequently asked questions
