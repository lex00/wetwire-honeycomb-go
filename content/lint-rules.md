---
title: "Lint Rules"
---

This document describes all lint rules for wetwire-honeycomb-go.

---

## Overview

wetwire-honeycomb linter checks query declarations for:
- Required fields (dataset, time range, calculations)
- Honeycomb API constraints and limits
- Best practices for query performance
- Common mistakes and anti-patterns

### Rule Naming

Rules follow the format `WHC<NNN>` where:
- `W` = wetwire
- `HC` = Honeycomb (domain code)
- `<NNN>` = Rule number (001, 002, etc.)

### Severity Levels

| Level | Description | Exit Code |
|-------|-------------|-----------|
| **error** | Must fix - violates Honeycomb constraints or missing required fields | 1 |
| **warning** | Should fix - potential performance or usability issue | 1 |

---

## Rule Index

| Rule | Description | Severity |
|------|-------------|----------|
| WHC001 | Missing dataset | error |
| WHC002 | Missing time range | error |
| WHC003 | Empty calculations | error |
| WHC004 | Breakdown without order | warning |
| WHC005 | High cardinality breakdown | warning |
| WHC006 | Invalid calculation for column type | error |
| WHC007 | Invalid filter operator | error |
| WHC008 | Missing limit with breakdowns | warning |
| WHC009 | Time range exceeds 7 days | error |
| WHC010 | Excessive filter count | warning |
| WHC011 | Circular dependency | error |
| WHC012 | Secret in filter | error |
| WHC013 | Sensitive column exposure | warning |
| WHC014 | Hardcoded credentials | error |
| WHC020 | Inline calculation definition | warning |
| WHC021 | Inline filter definition | warning |
| WHC022 | Raw map literal | warning |
| WHC023 | Deeply nested configuration | warning |
| **Board Rules** | | |
| WHC030 | Board has no panels | error |
| WHC034 | Board exceeds panel limit | warning |
| **SLO Rules** | | |
| WHC040 | SLO missing name | error |
| WHC044 | Target out of range | error |
| WHC047 | SLO no burn alerts | warning |
| **Trigger Rules** | | |
| WHC050 | Trigger missing name | error |
| WHC053 | Trigger no recipients | warning |
| WHC054 | Trigger frequency under 1 minute | warning |
| WHC056 | Trigger is disabled | info |

---

## Query Rules

### WHC001: Missing dataset

**Severity:** error

Every query must specify a dataset. This is a required field for the Honeycomb Query API.

**Bad:**

```go
var MyQuery = query.Query{
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
    // Missing Dataset field
}
```

**Good:**

```go
var MyQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(1),
    Calculations: []query.Calculation{
        query.Count(),
    },
}
```

---

### WHC002: Missing time range

**Severity:** error

Every query must specify a time range, either as a relative duration or absolute start/end times.

---

### WHC003: Empty calculations

**Severity:** error

Every query must have at least one calculation. A query without calculations would return no useful data.

---

### WHC009: Time range exceeds 7 days

**Severity:** error

Honeycomb has a maximum time range of 7 days for queries. Queries exceeding this limit will fail.

**Bad:**

```go
TimeRange: query.Days(30) // 30 days exceeds limit
```

**Good:**

```go
TimeRange: query.Days(7) // Maximum allowed
```

---

### WHC012: Secret in filter

**Severity:** error

Detects potential secrets, tokens, or credentials in filter values that should not be hardcoded.

**Detected Patterns:**
- OpenAI keys (`sk-...`)
- Generic tokens (values containing `token`, `bearer`, `secret`)
- API keys (values containing `api_key`, `apikey`)
- Passwords (values containing `password`)

---

## Board Rules

### WHC030: Board has no panels

**Severity:** error

Every board must have at least one panel.

### WHC034: Board exceeds panel limit

**Severity:** warning

Boards with more than 20 panels may have performance issues.

---

## SLO Rules

### WHC040: SLO missing name

**Severity:** error

Every SLO must have a name for identification.

### WHC044: Target out of range

**Severity:** error

SLO target percentage must be between 0 and 100.

### WHC047: SLO no burn alerts

**Severity:** warning

SLOs without burn alerts won't notify you when the error budget is being consumed too quickly.

---

## Trigger Rules

### WHC050: Trigger missing name

**Severity:** error

Every trigger must have a name for identification in alerts.

### WHC053: Trigger no recipients

**Severity:** warning

Triggers without recipients won't notify anyone when they fire.

### WHC054: Trigger frequency under 1 minute

**Severity:** warning

Trigger frequencies under 1 minute may cause excessive alerting.

### WHC056: Trigger is disabled

**Severity:** info

The trigger is explicitly disabled and will not fire alerts. This is informational, not necessarily a problem.

---

## Disabling Rules

### Command Line

Check specific rules only:

```bash
wetwire-honeycomb lint --rules WHC001,WHC002 ./queries/...
```

Disable specific rules:

```bash
wetwire-honeycomb lint --disable WHC003 ./queries/...
```

---

## See Also

- [CLI Reference](../cli/) - Complete command documentation
- [FAQ](../faq/) - Common questions
- [Honeycomb Query Best Practices](https://docs.honeycomb.io/working-with-your-data/queries/) - Official guide
