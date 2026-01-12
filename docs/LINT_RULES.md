# Lint Rules

This document describes all lint rules for wetwire-honeycomb-go.

---

## Overview

wetwire-honeycomb linter checks query declarations for:
- Best practices alignment with Honeycomb patterns
- Type safety enforcement
- Code clarity and maintainability
- Common mistakes and anti-patterns

### Rule Naming

Rules follow the format `WHC<NNN>` where:
- `W` = wetwire
- `HC` = Honeycomb (domain code)
- `<NNN>` = Rule number (001, 002, etc.)

### Severity Levels

| Level | Description | Exit Code |
|-------|-------------|-----------|
| **error** | Must fix - breaks compilation or runtime | 1 |
| **warning** | Should fix - violates best practices | 1 |
| **info** | Consider fixing - suggestions for improvement | 0 |

---

## Rule Index

| Rule | Description | Severity | Auto-fix |
|------|-------------|----------|----------|
| [WHC001](#whc001-use-typed-calculation-functions) | Use typed calculation functions | error | Yes |
| [WHC002](#whc002-use-typed-filter-functions) | Use typed filter functions | error | Yes |
| [WHC003](#whc003-validate-dataset-references) | Validate dataset references | warning | No |
| [WHC004](#whc004-validate-time-range) | Validate time range values | error | No |
| [WHC005](#whc005-unique-calculation-names) | Ensure calculation names are unique | warning | No |
| [WHC006](#whc006-avoid-excessive-breakdowns) | Avoid excessive breakdowns | warning | No |
| [WHC007](#whc007-prefer-direct-field-references) | Prefer direct field references | info | Yes |
| [WHC008](#whc008-validate-filter-values) | Validate filter value types | error | No |
| [WHC009](#whc009-check-calculation-field-compatibility) | Check calculation-field compatibility | warning | No |
| [WHC010](#whc010-limit-query-complexity) | Limit query complexity | info | No |

---

## Rules

### WHC001: Use typed calculation functions

**Severity:** error
**Auto-fix:** Yes

**Description:**

Use typed calculation functions (`query.P99()`, `query.Count()`, etc.) instead of raw struct initialization.

**Why:**

- Type safety: Compile-time validation of calculation types
- Readability: Clear, self-documenting code
- Maintainability: Easier to refactor and update

**Bad:**

```go
var MyQuery = query.Query{
    Calculations: []query.Calculation{
        {Type: "P99", Field: "duration_ms"},
        {Type: "COUNT"},
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
        query.Count(),
    },
}
```

**Auto-fix:**

Automatically replaces raw struct initialization with typed function calls.

---

### WHC002: Use typed filter functions

**Severity:** error
**Auto-fix:** Yes

**Description:**

Use typed filter functions (`query.GT()`, `query.Equals()`, etc.) instead of raw map initialization.

**Why:**

- Type safety: Prevents invalid filter operations
- Validation: Function parameters are validated at compile time
- Consistency: Matches Honeycomb filter patterns

**Bad:**

```go
var MyQuery = query.Query{
    Filters: []query.Filter{
        {Column: "duration_ms", Op: ">", Value: 500},
        {Column: "user_id", Op: "exists"},
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Filters: []query.Filter{
        query.GT("duration_ms", 500),
        query.Exists("user_id"),
    },
}
```

**Auto-fix:**

Automatically replaces raw map initialization with typed function calls.

---

### WHC003: Validate dataset references

**Severity:** warning
**Auto-fix:** No

**Description:**

Dataset names should reference defined constants or be validated against known datasets.

**Why:**

- Prevents typos in dataset names
- Makes dataset changes easier to track
- Improves refactoring safety

**Bad:**

```go
var MyQuery = query.Query{
    Dataset: "production",  // Hard-coded string
}
```

**Good:**

```go
const (
    DatasetProduction = "production"
    DatasetStaging    = "staging"
)

var MyQuery = query.Query{
    Dataset: DatasetProduction,
}
```

**Auto-fix:**

Not available. Requires manual definition of dataset constants.

---

### WHC004: Validate time range

**Severity:** error
**Auto-fix:** No

**Description:**

Time range values must be positive and use typed functions.

**Why:**

- Prevents invalid time ranges (negative, zero)
- Ensures Honeycomb API compatibility
- Catches common mistakes early

**Bad:**

```go
var MyQuery = query.Query{
    TimeRange: 0,           // Invalid: zero
}

var AnotherQuery = query.Query{
    TimeRange: -3600,       // Invalid: negative
}
```

**Good:**

```go
var MyQuery = query.Query{
    TimeRange: query.Hours(1),
}

var AnotherQuery = query.Query{
    TimeRange: query.Days(7),
}
```

**Auto-fix:**

Not available. Requires manual correction of time range values.

---

### WHC005: Unique calculation names

**Severity:** warning
**Auto-fix:** No

**Description:**

When using named calculations, ensure names are unique within a query.

**Why:**

- Prevents ambiguous results
- Makes query results easier to parse
- Aligns with Honeycomb best practices

**Bad:**

```go
var MyQuery = query.Query{
    Calculations: []query.Calculation{
        query.Avg("duration_ms").As("avg_duration"),
        query.Avg("response_time").As("avg_duration"),  // Duplicate name
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Calculations: []query.Calculation{
        query.Avg("duration_ms").As("avg_duration_ms"),
        query.Avg("response_time").As("avg_response_time"),
    },
}
```

**Auto-fix:**

Not available. Requires manual renaming of calculations.

---

### WHC006: Avoid excessive breakdowns

**Severity:** warning
**Auto-fix:** No

**Description:**

Queries with more than 5 breakdowns may have performance issues or be difficult to interpret.

**Why:**

- Performance: Too many breakdowns can slow query execution
- Usability: Results become difficult to visualize and interpret
- Best practice: Honeycomb recommends limiting breakdowns

**Bad:**

```go
var MyQuery = query.Query{
    Breakdowns: []string{
        "endpoint", "service", "region",
        "user_id", "device", "browser", "os",  // 7 breakdowns
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Breakdowns: []string{
        "endpoint", "service", "region",  // 3 breakdowns
    },
}
```

**Auto-fix:**

Not available. Requires manual query redesign.

---

### WHC007: Prefer direct field references

**Severity:** info
**Auto-fix:** Yes

**Description:**

When referencing fields in calculations and filters, prefer direct string references over variables unless the field name is reused multiple times.

**Why:**

- Clarity: Direct references are easier to read
- Simplicity: Reduces unnecessary indirection
- Convention: Matches Honeycomb patterns

**Bad:**

```go
const durationField = "duration_ms"

var MyQuery = query.Query{
    Calculations: []query.Calculation{
        query.P99(durationField),  // Only used once
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Calculations: []query.Calculation{
        query.P99("duration_ms"),
    },
}

// OR if field is reused multiple times:
const durationField = "duration_ms"

var ComplexQuery = query.Query{
    Calculations: []query.Calculation{
        query.P99(durationField),
        query.Avg(durationField),
        query.Max(durationField),
    },
    Filters: []query.Filter{
        query.GT(durationField, 500),
    },
}
```

**Auto-fix:**

Automatically inlines field references used only once.

---

### WHC008: Validate filter values

**Severity:** error
**Auto-fix:** No

**Description:**

Filter values must be compatible with the filter operation and field type.

**Why:**

- Type safety: Prevents runtime errors
- Validation: Catches mismatched types early
- Honeycomb compatibility: Ensures filters work as expected

**Bad:**

```go
var MyQuery = query.Query{
    Filters: []query.Filter{
        query.GT("duration_ms", "500"),     // String instead of int
        query.Equals("count", 3.14),        // Float for count field
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Filters: []query.Filter{
        query.GT("duration_ms", 500),
        query.Equals("count", 3),
    },
}
```

**Auto-fix:**

Not available. Requires manual type correction.

---

### WHC009: Check calculation-field compatibility

**Severity:** warning
**Auto-fix:** No

**Description:**

Ensure calculations are compatible with field types (e.g., percentiles on numeric fields only).

**Why:**

- Prevents query errors
- Aligns with Honeycomb field semantics
- Catches common mistakes

**Bad:**

```go
var MyQuery = query.Query{
    Calculations: []query.Calculation{
        query.P99("endpoint"),      // P99 on string field
        query.Sum("error_message"), // Sum on string field
    },
}
```

**Good:**

```go
var MyQuery = query.Query{
    Calculations: []query.Calculation{
        query.P99("duration_ms"),   // P99 on numeric field
        query.Count(),              // Count is always valid
        query.CountDistinct("endpoint"), // CountDistinct on string
    },
}
```

**Auto-fix:**

Not available. Requires manual calculation correction.

**Note:** This rule requires schema information to determine field types. Without schema metadata, this is a best-effort check based on field naming conventions.

---

### WHC010: Limit query complexity

**Severity:** info
**Auto-fix:** No

**Description:**

Queries with many filters, calculations, and breakdowns may be difficult to maintain and understand. Consider splitting into multiple queries.

**Why:**

- Maintainability: Simpler queries are easier to understand
- Performance: Complex queries may be slower
- Debugging: Easier to isolate issues

**Bad:**

```go
var ComplexQuery = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(24),
    Breakdowns: []string{"endpoint", "service", "region", "user_type"},
    Calculations: []query.Calculation{
        query.Count(),
        query.P50("duration_ms"),
        query.P95("duration_ms"),
        query.P99("duration_ms"),
        query.Avg("duration_ms"),
        query.Max("duration_ms"),
        query.Sum("bytes_sent"),
        query.CountDistinct("user_id"),
    },
    Filters: []query.Filter{
        query.GT("duration_ms", 100),
        query.Exists("user_id"),
        query.Contains("endpoint", "/api/"),
        query.NotEquals("status", 200),
        query.LT("error_count", 10),
    },
}
```

**Good:**

```go
// Split into focused queries
var PerformanceMetrics = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(24),
    Breakdowns: []string{"endpoint", "service"},
    Calculations: []query.Calculation{
        query.P95("duration_ms"),
        query.P99("duration_ms"),
    },
    Filters: []query.Filter{
        query.GT("duration_ms", 100),
    },
}

var UserMetrics = query.Query{
    Dataset:   "production",
    TimeRange: query.Hours(24),
    Breakdowns: []string{"user_type", "region"},
    Calculations: []query.Calculation{
        query.Count(),
        query.CountDistinct("user_id"),
    },
    Filters: []query.Filter{
        query.Exists("user_id"),
    },
}
```

**Auto-fix:**

Not available. Requires manual query redesign.

**Threshold:** Query is flagged if total of (filters + calculations + breakdowns) exceeds 15.

---

## Disabling Rules

### Command Line

Disable specific rules with `--disable`:

```bash
wetwire-honeycomb lint --disable WHC003,WHC007 ./queries/...
```

### Configuration File

Disable rules in `.wetwire-honeycomb.yaml`:

```yaml
lint:
  disabled_rules:
    - WHC003
    - WHC007
```

### Inline Comments

Disable rules for specific lines with comments:

```go
var MyQuery = query.Query{
    Dataset: "production", // wetwire:disable WHC003
}
```

Disable all rules for a query:

```go
// wetwire:disable-all
var LegacyQuery = query.Query{
    // ...
}
```

---

## Configuring Severity

Override severity levels in `.wetwire-honeycomb.yaml`:

```yaml
lint:
  severity_overrides:
    WHC003: error    # Promote to error
    WHC007: warning  # Promote to warning
    WHC010: ignore   # Ignore completely
```

---

## Auto-fix Behavior

Rules marked as auto-fixable can be corrected with `--fix`:

```bash
wetwire-honeycomb lint --fix ./queries/...
```

Auto-fix changes:
- **WHC001**: Replaces raw calculation structs with typed functions
- **WHC002**: Replaces raw filter maps with typed functions
- **WHC007**: Inlines single-use field references

Auto-fix preserves:
- Code formatting (runs `gofmt` after changes)
- Comments and documentation
- Import statements (adds/removes as needed)

Always review auto-fix changes before committing.

---

## Future Rules

Planned rules for future releases:

| Rule | Description | Priority |
|------|-------------|----------|
| WHC011 | Validate field names against schema | High |
| WHC012 | Check for unused breakdowns | Medium |
| WHC013 | Suggest query optimizations | Low |
| WHC014 | Detect duplicate queries | Medium |
| WHC015 | Validate HAVING clause usage | High |
| WHC016 | Check ORDER BY field references | Medium |
| WHC017 | Validate LIMIT values | Low |
| WHC018 | Suggest visualization types | Low |
| WHC019 | Check time alignment settings | Medium |
| WHC020 | Validate granularity values | High |

---

## See Also

- [CLI Reference](CLI.md) - Complete command documentation
- [FAQ](FAQ.md) - Common questions
- [Honeycomb Query Best Practices](https://docs.honeycomb.io/working-with-your-data/queries/) - Official guide
