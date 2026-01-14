# Adoption Guide

A practical guide for teams migrating to wetwire-honeycomb-go for Honeycomb query management.

---

## Table of Contents

1. [Why Adopt wetwire-honeycomb-go](#why-adopt-wetwire-honeycomb-go)
2. [Migration Paths](#migration-paths)
   - [From Manual Honeycomb UI](#from-manual-honeycomb-ui)
   - [From JSON/YAML Config Files](#from-jsonyaml-config-files)
   - [From Other Query Management Tools](#from-other-query-management-tools)
3. [Step-by-Step Adoption Plan](#step-by-step-adoption-plan)
   - [Phase 1: Install and Try with One Query](#phase-1-install-and-try-with-one-query)
   - [Phase 2: Import Existing Queries](#phase-2-import-existing-queries)
   - [Phase 3: Add SLOs and Triggers](#phase-3-add-slos-and-triggers)
   - [Phase 4: Create Boards and Full Stack](#phase-4-create-boards-and-full-stack)
4. [Team Onboarding Checklist](#team-onboarding-checklist)
5. [CI/CD Integration](#cicd-integration)
6. [Common Challenges and Solutions](#common-challenges-and-solutions)
7. [Success Metrics to Track](#success-metrics-to-track)

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

### Comparison: Before vs After

**Before (Manual/Ad-hoc):**
- Queries created in Honeycomb UI, difficult to track changes
- No review process for query modifications
- Copy-paste errors when sharing queries
- No visibility into who created or modified queries
- Queries scattered across teams with no central repository

**After (wetwire-honeycomb-go):**
- Queries defined in code, versioned in Git
- Pull request reviews for all query changes
- Import and share queries as Go packages
- Full audit trail of query modifications
- Central repository with organized structure

### When to Adopt

wetwire-honeycomb-go is ideal for teams that:

- Have more than 10-20 queries to manage
- Need to share queries across multiple services or teams
- Want infrastructure-as-code patterns for observability
- Require audit trails and review processes
- Are comfortable with Go development workflows

Consider alternatives if:

- You only need a few ad-hoc queries
- Your team doesn't use Go
- You prefer GUI-based workflows exclusively

---

## Migration Paths

### From Manual Honeycomb UI

If your team creates queries directly in the Honeycomb web interface:

#### Step 1: Export Query JSON from Honeycomb

In the Honeycomb UI:
1. Open your query
2. Click the "JSON" tab or export option
3. Copy the query JSON

Example exported JSON:
```json
{
  "time_range": 7200,
  "breakdowns": ["endpoint", "service"],
  "calculations": [
    {"op": "P99", "column": "duration_ms"},
    {"op": "COUNT"}
  ],
  "filters": [
    {"column": "duration_ms", "op": ">", "value": 500}
  ]
}
```

#### Step 2: Convert to Go Using the Import Command

```bash
wetwire-honeycomb import query.json -o queries/performance.go --name SlowRequests
```

This generates:
```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(2),
    Breakdowns: []string{"endpoint", "service"},
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
        query.Count(),
    },
    Filters: []query.Filter{
        query.GT("duration_ms", 500),
    },
}
```

#### Step 3: Validate and Iterate

```bash
# Lint for issues
wetwire-honeycomb lint ./queries/...

# Build to verify JSON output matches original
wetwire-honeycomb build --pretty ./queries/...
```

### From JSON/YAML Config Files

If you manage queries in JSON or YAML configuration files:

#### Step 1: Inventory Existing Files

```bash
# Find all query config files
find . -name "*.json" -o -name "*.yaml" | xargs grep -l "calculations\|time_range"
```

#### Step 2: Batch Convert Using Import

Convert each file:
```bash
# Single file
wetwire-honeycomb import ./config/latency-query.json -o queries/latency.go --name LatencyQuery

# Multiple files (script)
for f in ./config/*.json; do
  name=$(basename "$f" .json | sed 's/-/_/g' | sed 's/\b\(.\)/\u\1/g')
  wetwire-honeycomb import "$f" -o "queries/$(basename "$f" .json).go" --name "$name"
done
```

#### Step 3: Organize and Consolidate

Move related queries into logical files:
```
queries/
├── performance/
│   ├── latency.go       # P99, P95 queries
│   └── throughput.go    # Request rate queries
├── errors/
│   ├── rates.go         # Error rate queries
│   └── breakdown.go     # Error breakdown queries
└── slos/
    └── availability.go  # SLI queries
```

#### Step 4: Update Deployment Pipeline

Replace JSON file reads with wetwire-honeycomb build output:

```bash
# Old: copy JSON files
cp ./config/*.json ./deploy/

# New: generate from Go
wetwire-honeycomb build -o ./deploy/queries.json ./queries/...
```

### From Other Query Management Tools

If migrating from another tool (Terraform, custom scripts, etc.):

#### Step 1: Export to Intermediate Format

Most tools can export to JSON. Generate JSON exports of your existing queries.

#### Step 2: Follow JSON Migration Path

Use the import command to convert JSON to Go as described above.

#### Step 3: Map Tool-Specific Features

| Feature | wetwire-honeycomb Equivalent |
|---------|------------------------------|
| Query templates | Go functions returning `query.Query` |
| Variable substitution | Go variables and composition |
| Conditional queries | Go `if` statements in helper functions |
| Query grouping | Package organization |
| Deployment automation | CI/CD with `build` command |

---

## Step-by-Step Adoption Plan

### Phase 1: Install and Try with One Query

**Duration:** 1-2 hours

**Goal:** Get a single query working end-to-end

#### 1.1 Install wetwire-honeycomb

```bash
go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest
```

#### 1.2 Initialize a Project

```bash
mkdir honeycomb-queries
cd honeycomb-queries
wetwire-honeycomb init queries
```

This creates a basic structure:
```
honeycomb-queries/
├── queries/
│   └── example.go
└── go.mod
```

#### 1.3 Create Your First Query

Edit `queries/example.go`:

```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// SlowRequests identifies API requests with high latency
var SlowRequests = query.Query{
    Dataset:   "production",  // Your actual dataset name
    TimeRange: query.Hours(2),
    Breakdowns: []string{"endpoint"},
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
        query.Count(),
    },
    Filters: []query.Filter{
        query.GT("duration_ms", 500),
    },
}
```

#### 1.4 Validate and Build

```bash
# Check for issues
wetwire-honeycomb lint ./queries/...

# Generate JSON
wetwire-honeycomb build --pretty ./queries/...
```

#### 1.5 Verify in Honeycomb

Copy the generated JSON and paste into Honeycomb's Query Builder JSON view to verify it works.

**Phase 1 Success Criteria:**
- [ ] wetwire-honeycomb installed successfully
- [ ] One query defined in Go
- [ ] Query builds to valid JSON
- [ ] JSON works when tested in Honeycomb UI

### Phase 2: Import Existing Queries

**Duration:** 1-2 days (depending on query count)

**Goal:** Convert critical existing queries to wetwire-honeycomb

#### 2.1 Identify High-Priority Queries

List your most important queries:
- Dashboard queries
- SLO indicator queries
- Alert trigger queries
- Frequently-used debugging queries

#### 2.2 Export and Convert

For each query:

```bash
# Save JSON from Honeycomb UI to file
# Then import:
wetwire-honeycomb import exported_query.json \
  -o queries/performance/latency.go \
  --name EndpointLatency \
  --package performance
```

#### 2.3 Organize by Domain

Create a clear file structure:

```
queries/
├── performance/
│   ├── latency.go
│   └── throughput.go
├── errors/
│   ├── rates.go
│   └── top_errors.go
├── sli/
│   ├── availability.go
│   └── latency_budget.go
└── debugging/
    └── trace_analysis.go
```

#### 2.4 Add Documentation

Add Go comments to explain each query:

```go
// EndpointLatency tracks P99 latency across all API endpoints.
// Used in the API Performance Dashboard and LatencySLO.
// Threshold: Alert when P99 > 1000ms
var EndpointLatency = query.Query{
    // ...
}
```

#### 2.5 Lint and Fix Issues

```bash
# Check all queries
wetwire-honeycomb lint ./queries/...

# Auto-fix where possible
wetwire-honeycomb lint --fix ./queries/...
```

**Phase 2 Success Criteria:**
- [ ] All critical queries converted to Go
- [ ] Queries organized into logical packages
- [ ] All queries pass linting
- [ ] Generated JSON matches original behavior

### Phase 3: Add SLOs and Triggers

**Duration:** 2-3 days

**Goal:** Define SLOs and triggers that reference your queries

#### 3.1 Create SLOs

Define Service Level Objectives using your query definitions:

```go
package slos

import (
    "github.com/lex00/wetwire-honeycomb-go/query"
    "github.com/lex00/wetwire-honeycomb-go/slo"
    "yourorg/queries/sli"
)

// APIAvailability tracks 99.9% availability over 30 days
var APIAvailability = slo.SLO{
    Name:        "API Availability",
    Description: "99.9% of API requests succeed (status < 500)",
    Dataset:     "production",
    SLI: slo.SLI{
        GoodEvents:  sli.SuccessfulRequests,  // Reference existing query
        TotalEvents: sli.AllRequests,          // Reference existing query
    },
    Target:     slo.Percentage(99.9),
    TimePeriod: slo.Days(30),
    BurnAlerts: []slo.BurnAlert{
        slo.FastBurn(2.0),
        slo.SlowBurn(5.0),
    },
}
```

See [SLOS.md](SLOS.md) for complete SLO documentation.

#### 3.2 Create Triggers

Define alerts that monitor query results:

```go
package triggers

import (
    "github.com/lex00/wetwire-honeycomb-go/trigger"
    "yourorg/queries/performance"
)

// HighLatencyAlert fires when P99 exceeds 1 second
var HighLatencyAlert = trigger.Trigger{
    Name:        "High P99 Latency",
    Description: "Alert when P99 latency exceeds 1000ms",
    Dataset:     "production",
    Query:       performance.EndpointLatency,  // Reference existing query
    Threshold:   trigger.GreaterThan(1000),
    Frequency:   trigger.Minutes(5),
    Recipients: []trigger.Recipient{
        trigger.SlackChannel("#alerts"),
        trigger.PagerDutyService("api-oncall"),
    },
}
```

See [TRIGGERS.md](TRIGGERS.md) for complete trigger documentation.

#### 3.3 Validate Connections

Verify queries, SLOs, and triggers connect correctly:

```bash
# List all resources
wetwire-honeycomb list ./...

# Build all resources
wetwire-honeycomb build -o all_resources.json ./...
```

**Phase 3 Success Criteria:**
- [ ] SLOs defined for critical services
- [ ] Triggers defined for key metrics
- [ ] SLOs and triggers reference shared queries
- [ ] All resources build successfully

### Phase 4: Create Boards and Full Stack

**Duration:** 2-3 days

**Goal:** Complete monitoring stack with dashboards

#### 4.1 Create Boards

Build dashboards that combine queries, SLOs, and documentation:

```go
package boards

import (
    "github.com/lex00/wetwire-honeycomb-go/board"
    "yourorg/queries/performance"
    "yourorg/queries/errors"
)

// APIPerformance is the primary API monitoring dashboard
var APIPerformance = board.Board{
    Name:        "API Performance Dashboard",
    Description: "Real-time monitoring of API latency, throughput, and errors",
    Panels: []board.Panel{
        board.QueryPanel(performance.EndpointLatency,
            board.WithTitle("P99 Latency by Endpoint"),
            board.WithPosition(0, 0, 6, 4),
        ),
        board.QueryPanel(errors.ErrorRate,
            board.WithTitle("Error Rate"),
            board.WithPosition(6, 0, 6, 4),
        ),
        board.TextPanel(`
## Runbook
1. Check error rate trends
2. Identify affected endpoints
3. Review recent deployments

**On-call:** #api-oncall
        `,
            board.WithTitle("Investigation Guide"),
            board.WithPosition(0, 4, 12, 3),
        ),
    },
    PresetFilters: []board.Filter{
        {Column: "environment", Operation: "=", Value: "production"},
    },
    Tags: []board.Tag{
        {Key: "team", Value: "api"},
        {Key: "tier", Value: "critical"},
    },
}
```

See [BOARDS.md](BOARDS.md) for complete board documentation.

#### 4.2 Reference the Full Stack Example

Study the complete example in `examples/full_stack/`:

```
examples/full_stack/
├── queries.go    # Base metric queries
├── slos.go       # SLOs referencing queries
├── triggers.go   # Triggers referencing queries
└── boards.go     # Boards combining all resources
```

#### 4.3 Build Complete Output

Generate all resources:

```bash
# Build queries
wetwire-honeycomb build -o output/queries.json ./queries/...

# Build SLOs
wetwire-honeycomb build -o output/slos.json ./slos/...

# Build triggers
wetwire-honeycomb build -o output/triggers.json ./triggers/...

# Build boards
wetwire-honeycomb build -o output/boards.json ./boards/...

# Or build everything at once
wetwire-honeycomb build -o output/all.json ./...
```

**Phase 4 Success Criteria:**
- [ ] Boards created for key use cases
- [ ] Full stack builds without errors
- [ ] Dashboard layout reviewed with team
- [ ] JSON output validated in Honeycomb

---

## Team Onboarding Checklist

Use this checklist when onboarding new team members to wetwire-honeycomb-go.

### Prerequisites

- [ ] Go 1.21+ installed
- [ ] wetwire-honeycomb CLI installed
- [ ] Git repository access
- [ ] Honeycomb account with API access

### Documentation Review

- [ ] Read [QUICK_START.md](QUICK_START.md) - Basic usage
- [ ] Read [CLI.md](CLI.md) - Command reference
- [ ] Read [LINT_RULES.md](LINT_RULES.md) - Understanding lint rules
- [ ] Review [examples/](../examples/) directory

### Hands-On Practice

- [ ] Create a simple query using `wetwire-honeycomb init`
- [ ] Run `wetwire-honeycomb lint` and fix any issues
- [ ] Run `wetwire-honeycomb build` and review JSON output
- [ ] Test generated JSON in Honeycomb UI
- [ ] Import an existing query using `wetwire-honeycomb import`

### Team Conventions

- [ ] Understand file/package organization
- [ ] Learn query naming conventions
- [ ] Review PR process for query changes
- [ ] Know who to ask for help

### Development Workflow

- [ ] Clone the queries repository
- [ ] Create a branch for changes
- [ ] Make changes and run lint/build locally
- [ ] Submit PR for review
- [ ] Merge after approval

---

## CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/honeycomb-queries.yml
name: Honeycomb Queries

on:
  push:
    branches: [main]
    paths:
      - 'queries/**'
      - 'slos/**'
      - 'triggers/**'
      - 'boards/**'
  pull_request:
    paths:
      - 'queries/**'
      - 'slos/**'
      - 'triggers/**'
      - 'boards/**'

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
        run: wetwire-honeycomb lint --format json ./... > lint-results.json

      - name: Check for lint errors
        run: |
          if [ $(jq '.summary.errors' lint-results.json) -gt 0 ]; then
            echo "Lint errors found:"
            jq '.issues[] | select(.severity == "error")' lint-results.json
            exit 1
          fi

      - name: Build all resources
        run: wetwire-honeycomb build -o output.json ./...

      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: honeycomb-resources
          path: output.json

  deploy:
    needs: validate
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/download-artifact@v3
        with:
          name: honeycomb-resources

      - name: Deploy to Honeycomb
        env:
          HONEYCOMB_API_KEY: ${{ secrets.HONEYCOMB_API_KEY }}
        run: |
          # Your deployment script here
          # This could push to Honeycomb API or store in a central location
          echo "Deploying resources..."
```

### GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - validate
  - deploy

variables:
  GO_VERSION: "1.21"

validate:
  stage: validate
  image: golang:${GO_VERSION}
  script:
    - go install github.com/lex00/wetwire-honeycomb-go/cmd/wetwire-honeycomb@latest
    - wetwire-honeycomb lint ./...
    - wetwire-honeycomb build -o output.json ./...
  artifacts:
    paths:
      - output.json
    expire_in: 1 week
  rules:
    - changes:
        - queries/**/*
        - slos/**/*
        - triggers/**/*
        - boards/**/*

deploy:
  stage: deploy
  image: alpine:latest
  script:
    - echo "Deploying resources..."
    # Add deployment commands here
  dependencies:
    - validate
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
      changes:
        - queries/**/*
        - slos/**/*
        - triggers/**/*
        - boards/**/*
```

### Pre-commit Hook

Add a pre-commit hook to validate before pushing:

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running wetwire-honeycomb lint..."

if ! wetwire-honeycomb lint ./queries/... ./slos/... ./triggers/... ./boards/...; then
    echo ""
    echo "Lint errors found. Please fix before committing."
    echo "Run: wetwire-honeycomb lint --fix ./..."
    exit 1
fi

echo "Lint passed!"
```

Make it executable:
```bash
chmod +x .git/hooks/pre-commit
```

---

## Common Challenges and Solutions

### Challenge: Team Resistance to Code-Based Queries

**Problem:** Team members prefer the Honeycomb UI and resist learning Go syntax.

**Solutions:**
- Start with power users who are comfortable with code
- Show productivity benefits (version control, review process)
- Use `wetwire-honeycomb import` to reduce manual conversion work
- Pair program on first few queries
- Create templates for common query patterns

### Challenge: Large Backlog of Existing Queries

**Problem:** Hundreds of queries to migrate is overwhelming.

**Solutions:**
- Prioritize by usage (start with dashboard queries)
- Migrate incrementally (phase approach)
- Use batch import scripts
- Set a deadline for new queries (new in code, legacy in UI)
- Track migration progress with a checklist

### Challenge: Maintaining Consistency Across Teams

**Problem:** Different teams use different patterns and naming conventions.

**Solutions:**
- Establish organization-wide conventions (naming, file structure)
- Use lint rules to enforce patterns
- Create shared packages for common queries
- Regular code reviews for query changes
- Document conventions in a CONTRIBUTING.md

### Challenge: Debugging JSON Output Issues

**Problem:** Generated JSON doesn't work as expected in Honeycomb.

**Solutions:**
- Use `wetwire-honeycomb build --pretty` for readable output
- Compare generated JSON with working UI queries
- Enable verbose mode: `wetwire-honeycomb build -v`
- Check dataset and column names match exactly
- Test in Honeycomb Query Builder JSON tab

### Challenge: CI/CD Integration Complexity

**Problem:** Integrating with existing deployment pipelines is difficult.

**Solutions:**
- Start with lint-only CI (just validation)
- Add build step once lint is stable
- Separate validation from deployment
- Use artifacts to pass JSON between stages
- Document the pipeline for team reference

### Challenge: Query Dependencies and Ordering

**Problem:** SLOs depend on queries, triggers depend on queries, order matters.

**Solutions:**
- Build resources in dependency order (queries first)
- Use separate output files per resource type
- Organize packages to reflect dependencies
- Document dependency chains in comments

---

## Success Metrics to Track

### Adoption Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| **Query Coverage** | Percentage of queries managed in wetwire | 80%+ within 6 months |
| **Team Onboarding** | Time for new team member to create first query | < 2 hours |
| **Active Contributors** | Number of team members with merged query PRs | > 50% of observability users |
| **Import Velocity** | Queries converted per week during migration | 10-20 queries/week |

### Quality Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| **Lint Pass Rate** | Percentage of PRs passing lint on first try | > 90% |
| **Build Success Rate** | Percentage of builds completing without errors | > 99% |
| **Query Reuse** | Average number of SLOs/triggers per query | > 1.5 |
| **Documentation Coverage** | Percentage of queries with meaningful comments | > 80% |

### Operational Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| **Time to Change** | Time from query modification to production | < 1 hour |
| **Change Review Time** | Average time for query PR review | < 24 hours |
| **Incident Detection** | Time to identify query-related issues | < 5 minutes |
| **Rollback Frequency** | How often query changes need to be reverted | < 5% of changes |

### Tracking Progress

Create a simple tracking document:

```markdown
# wetwire-honeycomb Adoption Progress

## Migration Status
- Total existing queries: 150
- Queries migrated: 45 (30%)
- Target completion: Q2 2024

## Weekly Stats
| Week | Queries Added | PRs Merged | Lint Issues |
|------|--------------|------------|-------------|
| W1   | 12           | 3          | 8           |
| W2   | 15           | 4          | 3           |
| W3   | 18           | 5          | 2           |

## Blockers
- [ ] Need Honeycomb API access for CI/CD
- [ ] Training session scheduled for Team B
```

---

## Next Steps

After completing adoption:

1. **Expand Coverage** - Continue migrating remaining queries
2. **Share Learnings** - Document patterns that worked well
3. **Contribute Back** - Report issues or contribute improvements to wetwire-honeycomb-go
4. **Optimize Workflow** - Refine CI/CD based on experience
5. **Train Others** - Help other teams adopt the same approach

---

## See Also

- [QUICK_START.md](QUICK_START.md) - Getting started guide
- [CLI.md](CLI.md) - Complete command reference
- [SLOS.md](SLOS.md) - SLO documentation
- [TRIGGERS.md](TRIGGERS.md) - Trigger documentation
- [BOARDS.md](BOARDS.md) - Board documentation
- [LINT_RULES.md](LINT_RULES.md) - Lint rule reference
- [FAQ.md](FAQ.md) - Frequently asked questions
- [examples/full_stack/](../examples/full_stack/) - Complete working example
