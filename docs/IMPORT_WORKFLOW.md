<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

This document describes the workflow for importing existing Honeycomb Query JSON into wetwire-honeycomb-go code.

---

## Overview

The `import` command converts Honeycomb Query JSON files into type-safe Go code. This enables migration from existing JSON-based query definitions to the wetwire-honeycomb approach, bringing benefits like compile-time validation, linting, and version control.

```bash
wetwire-honeycomb import <file.json>
```

The command reads a Query JSON file and generates Go code that declares an equivalent `query.Query` struct.

---

## Prerequisites

Before importing, you need:

1. **Honeycomb Query JSON** - Export from Honeycomb UI or existing configuration
2. **wetwire-honeycomb CLI** - Installed and available in PATH
3. **Target Go package** - A directory where imported queries will live

### Obtaining Query JSON from Honeycomb

There are several ways to get Query JSON from Honeycomb:

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

**From existing configuration files:**
If you have Query JSON in configuration management systems, export it to a standalone JSON file.

---

## Step-by-Step Workflow

### Step 1: Export Query from Honeycomb

Save your query JSON to a file. The JSON should follow the Honeycomb Query Specification format:

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

Save this as `slow_requests.json`.

### Step 2: Run the Import Command

Use `wetwire-honeycomb import` to generate Go code:

```bash
# Output to stdout
wetwire-honeycomb import slow_requests.json

# Output to file
wetwire-honeycomb import -o queries/slow_requests.go slow_requests.json

# Specify package and variable name
wetwire-honeycomb import -p queries -n SlowRequests -o queries/slow_requests.go slow_requests.json
```

### Step 3: Review Generated Code

The import command generates Go code like:

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

1. **Add the Dataset field** - Import does not include dataset information from the JSON
2. **Verify time range conversion** - Check that the time range matches your intent
3. **Review filter values** - Ensure values are correctly typed
4. **Add documentation** - Add comments explaining the query purpose
5. **Run lint** - Validate the generated code

Example of customized output:

```go
package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// SlowRequests finds requests taking longer than 500ms.
// Used for performance monitoring dashboards.
var SlowRequests = query.Query{
	Dataset:   "production",  // Added manually
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

### Step 5: Validate with Lint and Build

```bash
# Lint the imported query
wetwire-honeycomb lint ./queries/...

# Build to verify JSON output matches original
wetwire-honeycomb build ./queries/...
```

---

## Command Reference

### import

Convert Query JSON to Go code.

```bash
wetwire-honeycomb import [OPTIONS] <file.json>
```

**Arguments:**

| Argument | Description | Required |
|----------|-------------|----------|
| `file.json` | Path to the Query JSON file to import | Yes |

**Options:**

| Flag | Description | Default |
|------|-------------|---------|
| `-o, --output FILE` | Write generated Go code to FILE | stdout |
| `-p, --package NAME` | Package name for generated code | `queries` |
| `-n, --name NAME` | Variable name for the query | `Query` |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error reading or parsing JSON file |

---

## Examples

### Basic Import

Import a query and output to stdout:

```bash
wetwire-honeycomb import query.json
```

### Import with Custom Package

Generate code for a specific package:

```bash
wetwire-honeycomb import -p performance -n LatencyMetrics query.json
```

Output:

```go
package performance

import "github.com/lex00/wetwire-honeycomb-go/query"

var LatencyMetrics = query.Query{
	// ...
}
```

### Import to File

Write directly to a Go source file:

```bash
wetwire-honeycomb import -o queries/api_errors.go -n APIErrors api_errors.json
```

### Batch Import Multiple Queries

Import multiple JSON files using a shell loop:

```bash
for f in exported/*.json; do
  name=$(basename "$f" .json | sed 's/_\([a-z]\)/\U\1/g' | sed 's/^./\U&/')
  wetwire-honeycomb import -p queries -n "$name" -o "queries/$(basename $f .json).go" "$f"
done
```

### Import and Immediately Lint

Chain import with lint to catch issues:

```bash
wetwire-honeycomb import -o queries/new_query.go -n NewQuery query.json && \
  wetwire-honeycomb lint ./queries/new_query.go
```

---

## Supported JSON Fields

The import command converts the following Query JSON fields:

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
| `P75` | `query.P75(column)` |
| `P90` | `query.P90(column)` |
| `P95` | `query.P95(column)` |
| `P99` | `query.P99(column)` |
| `AVG` | `query.Avg(column)` |
| `SUM` | `query.Sum(column)` |
| `MIN` | `query.Min(column)` |
| `MAX` | `query.Max(column)` |
| Other | Raw struct `{Op: "...", Column: "..."}` |

### Filter Conversions

| JSON op | Go Function |
|---------|-------------|
| `=` | `query.Equals(column, value)` |
| `!=` | `query.NotEquals(column, value)` |
| `>` | `query.GT(column, value)` |
| `>=` | `query.GTE(column, value)` |
| `<` | `query.LT(column, value)` |
| `<=` | `query.LTE(column, value)` |
| Other | Raw struct `{Column: "...", Op: "...", Value: ...}` |

---

## Limitations

The import command has the following limitations that may require manual editing:

### Not Imported

| Feature | Reason | Manual Action |
|---------|--------|---------------|
| Dataset | Not in Query JSON | Add `Dataset` field manually |
| Granularity | Complex time logic | Add `Granularity` if needed |
| Orders | Ordering requires context | Add `Orders` field if needed |
| Havings | Complex clause structure | Add `Havings` field if needed |
| Filter combination | AND/OR logic | Set `FilterCombination` if needed |
| Start/end time | Absolute time handling | Use `query.Absolute()` if needed |

### Partial Conversions

| Feature | Issue | Manual Action |
|---------|-------|---------------|
| Unknown calculations | Ops not in mapping | Convert raw struct to typed function |
| Unknown filters | Ops not in mapping | Convert raw struct to typed function |
| Complex filters | Nested OR logic | Restructure using `query.Or()` |
| Derived columns | Not supported | Define separately if needed |

### Example: Manual Additions

After import, you may need to add:

```go
var MyQuery = query.Query{
	Dataset:   "production",     // Add dataset
	TimeRange: query.Hours(2),   // From import
	Granularity: 60,             // Add if needed
	Breakdowns: []string{"endpoint"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Orders: []query.Order{       // Add if needed
		{Op: "COUNT", Order: "descending"},
	},
	FilterCombination: "AND",    // Add if needed
	Limit: 100,
}
```

---

## Best Practices

### Organizing Imported Queries

1. **Create a dedicated directory for imports**

   ```
   queries/
   ├── imported/       # Recently imported, needs review
   │   ├── legacy_query1.go
   │   └── legacy_query2.go
   └── reviewed/       # Reviewed and customized
       ├── performance.go
       └── errors.go
   ```

2. **Use meaningful variable names**

   ```bash
   # Bad: Generic name
   wetwire-honeycomb import -n Query query.json

   # Good: Descriptive name
   wetwire-honeycomb import -n SlowAPIRequests query.json
   ```

3. **Group related queries in the same file**

   Import related queries, then manually combine them into a single file with shared constants.

### Post-Import Checklist

After importing each query:

- [ ] Add `Dataset` field
- [ ] Add documentation comment explaining the query purpose
- [ ] Verify time range matches original intent
- [ ] Check filter values are correctly typed
- [ ] Run `wetwire-honeycomb lint` to check for issues
- [ ] Run `wetwire-honeycomb build` and compare output to original JSON
- [ ] Move from `imported/` to reviewed location
- [ ] Commit with descriptive message

### Handling Large Migrations

For migrating many queries:

1. **Export all queries first**
   ```bash
   mkdir -p exported imported reviewed
   # Export queries from Honeycomb to exported/
   ```

2. **Batch import**
   ```bash
   for f in exported/*.json; do
     name=$(basename "$f" .json)
     wetwire-honeycomb import -o "imported/${name}.go" -n "${name^}" "$f"
   done
   ```

3. **Review incrementally**
   - Move queries to `reviewed/` as you validate them
   - Track progress with a checklist

4. **Validate the migration**
   ```bash
   wetwire-honeycomb lint ./reviewed/...
   wetwire-honeycomb build ./reviewed/...
   ```

### Maintaining Consistency

When importing multiple queries:

- Use consistent package naming (`package queries`)
- Use consistent variable naming conventions (PascalCase)
- Group queries by domain (performance, errors, auth, etc.)
- Extract common filters and calculations to shared variables

---

## Troubleshooting

### "Error reading file"

```
error reading file: open query.json: no such file or directory
```

**Solution:** Verify the file path is correct and the file exists.

### "Error parsing JSON"

```
error parsing JSON: invalid character ',' looking for beginning of value
```

**Solution:** Validate the JSON file format. Use a JSON validator or `jq`:

```bash
jq . query.json
```

### Generated code has raw structs

If the output contains raw structs instead of typed functions:

```go
Calculations: []query.Calculation{
	{Op: "CUSTOM_OP", Column: "field"},  // Raw struct
},
```

**Solution:** The operation is not in the import mapping. Manually convert to an appropriate typed function or keep the raw struct if the operation is custom.

### Missing fields after import

If critical fields are missing from the generated code:

**Solution:** The import command only handles standard Query JSON fields. Add missing fields manually:
- `Dataset`
- `Granularity`
- `Orders`
- `Havings`
- `FilterCombination`

---

## See Also

- [CLI Reference](CLI.md) - Complete command documentation
- [FAQ](FAQ.md) - Common questions
- [Lint Rules](LINT_RULES.md) - All WHC rules
- [Honeycomb Query Specification](https://docs.honeycomb.io/api/query-specification/) - API reference
