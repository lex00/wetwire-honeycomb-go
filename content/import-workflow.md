---
title: "Import Workflow"
---

This document describes the workflow for importing existing Honeycomb Query JSON into wetwire-honeycomb-go code.

---

## Overview

The `import` command converts Honeycomb Query JSON files into type-safe Go code. This enables migration from existing JSON-based query definitions to the wetwire-honeycomb approach, bringing benefits like compile-time validation, linting, and version control.

```bash
wetwire-honeycomb import <file.json>
```

---

## Prerequisites

Before importing, you need:

1. **Honeycomb Query JSON** - Export from Honeycomb UI or existing configuration
2. **wetwire-honeycomb CLI** - Installed and available in PATH
3. **Target Go package** - A directory where imported queries will live

### Obtaining Query JSON from Honeycomb

**From the Honeycomb UI:**
1. Open your query in the Honeycomb query builder
2. Click the "Share" or "Export" button
3. Select "Copy as JSON" or "Download JSON"
4. Save the JSON to a file

**From the Honeycomb API:**
```bash
curl -X GET "https://api.honeycomb.io/1/queries/YOUR_DATASET/QUERY_ID" \
  -H "X-Honeycomb-Team: $HONEYCOMB_API_KEY" \
  -o query.json
```

---

## Step-by-Step Workflow

### Step 1: Export Query from Honeycomb

Save your query JSON to a file:

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
  ],
  "limit": 100
}
```

### Step 2: Run the Import Command

```bash
# Output to stdout
wetwire-honeycomb import slow_requests.json

# Output to file
wetwire-honeycomb import -o queries/slow_requests.go slow_requests.json

# Specify package and variable name
wetwire-honeycomb import -p queries -n SlowRequests -o queries/slow_requests.go slow_requests.json
```

### Step 3: Review Generated Code

The import command generates:

```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
	TimeRange: query.Hours(2),
	Breakdowns: []string{"endpoint", "service"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.GT("duration_ms", 500),
	},
	Limit: 100,
}
```

### Step 4: Customize Generated Code

After import, you should:

1. **Add the Dataset field** - Import does not include dataset information
2. **Verify time range conversion** - Check that the time range matches your intent
3. **Review filter values** - Ensure values are correctly typed
4. **Add documentation** - Add comments explaining the query purpose
5. **Run lint** - Validate the generated code

### Step 5: Validate with Lint and Build

```bash
wetwire-honeycomb lint ./queries/...
wetwire-honeycomb build ./queries/...
```

---

## Command Reference

### import

```bash
wetwire-honeycomb import [OPTIONS] <file.json>
```

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output FILE` | Write generated Go code to FILE | stdout |
| `-p, --package NAME` | Package name for generated code | `queries` |
| `-n, --name NAME` | Variable name for the query | `Query` |

---

## Supported JSON Fields

| JSON Field | Go Field | Notes |
|------------|----------|-------|
| `time_range` | `TimeRange` | Converted to `query.Hours()` or `query.Seconds()` |
| `breakdowns` | `Breakdowns` | Array of column names |
| `calculations` | `Calculations` | Converted to typed functions |
| `filters` | `Filters` | Converted to typed filter functions |
| `limit` | `Limit` | Integer limit value |

### Calculation Conversions

| JSON op | Go Function |
|---------|-------------|
| `COUNT` (no column) | `query.Count()` |
| `COUNT` (with column) | `query.CountDistinct(column)` |
| `P50` | `query.P50(column)` |
| `P95` | `query.P95(column)` |
| `P99` | `query.P99(column)` |
| `AVG` | `query.Avg(column)` |
| `SUM` | `query.Sum(column)` |
| `MIN` | `query.Min(column)` |
| `MAX` | `query.Max(column)` |

### Filter Conversions

| JSON op | Go Function |
|---------|-------------|
| `=` | `query.Equals(column, value)` |
| `!=` | `query.NotEquals(column, value)` |
| `>` | `query.GT(column, value)` |
| `>=` | `query.GTE(column, value)` |
| `<` | `query.LT(column, value)` |
| `<=` | `query.LTE(column, value)` |

---

## Limitations

The import command has the following limitations that may require manual editing:

### Not Imported

| Feature | Manual Action |
|---------|---------------|
| Dataset | Add `Dataset` field manually |
| Granularity | Add `Granularity` if needed |
| Orders | Add `Orders` field if needed |
| Havings | Add `Havings` field if needed |
| Filter combination | Set `FilterCombination` if needed |

---

## Best Practices

### Organizing Imported Queries

```
queries/
├── imported/       # Recently imported, needs review
│   ├── legacy_query1.go
│   └── legacy_query2.go
└── reviewed/       # Reviewed and customized
    ├── performance.go
    └── errors.go
```

### Post-Import Checklist

- [ ] Add `Dataset` field
- [ ] Add documentation comment explaining the query purpose
- [ ] Verify time range matches original intent
- [ ] Check filter values are correctly typed
- [ ] Run `wetwire-honeycomb lint` to check for issues
- [ ] Run `wetwire-honeycomb build` and compare output to original JSON
- [ ] Move from `imported/` to reviewed location
- [ ] Commit with descriptive message

### Batch Import Multiple Queries

```bash
for f in exported/*.json; do
  name=$(basename "$f" .json | sed 's/_\([a-z]\)/\U\1/g' | sed 's/^./\U&/')
  wetwire-honeycomb import -p queries -n "$name" -o "queries/$(basename $f .json).go" "$f"
done
```

---

## Troubleshooting

### "Error reading file"

Verify the file path is correct and the file exists.

### "Error parsing JSON"

Validate the JSON file format:

```bash
jq . query.json
```

### Generated code has raw structs

The operation is not in the import mapping. Manually convert to an appropriate typed function or keep the raw struct if the operation is custom.

### Missing fields after import

The import command only handles standard Query JSON fields. Add missing fields manually.

---

## See Also

- [CLI Reference](../cli/) - Complete command documentation
- [FAQ](../faq/) - Common questions
- [Lint Rules](../lint-rules/) - All WHC rules
- [Honeycomb Query Specification](https://docs.honeycomb.io/api/query-specification/) - API reference
