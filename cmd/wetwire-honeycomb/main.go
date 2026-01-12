package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/lex00/wetwire-core-go/agent/agents"
	"github.com/lex00/wetwire-core-go/agent/orchestrator"
	"github.com/lex00/wetwire-core-go/agent/personas"
	"github.com/lex00/wetwire-core-go/providers"
	"github.com/lex00/wetwire-core-go/providers/anthropic"
	"github.com/lex00/wetwire-honeycomb-go/internal/agent"
	"github.com/lex00/wetwire-honeycomb-go/internal/builder"
	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
	"github.com/lex00/wetwire-honeycomb-go/internal/lint"
	"github.com/lex00/wetwire-honeycomb-go/internal/serialize"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "build":
		os.Exit(buildCmd(os.Args[2:]))
	case "lint":
		os.Exit(lintCmd(os.Args[2:]))
	case "list":
		os.Exit(listCmd(os.Args[2:]))
	case "import":
		os.Exit(importCmd(os.Args[2:]))
	case "validate":
		os.Exit(validateCmd(os.Args[2:]))
	case "init":
		os.Exit(initCmd(os.Args[2:]))
	case "graph":
		os.Exit(graphCmd(os.Args[2:]))
	case "diff":
		os.Exit(diffCmd(os.Args[2:]))
	case "watch":
		os.Exit(watchCmd(os.Args[2:]))
	case "design":
		os.Exit(designCmd(os.Args[2:]))
	case "test":
		os.Exit(testCmd(os.Args[2:]))
	case "version":
		fmt.Printf("wetwire-honeycomb %s\n", version)
		os.Exit(0)
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("wetwire-honeycomb - Honeycomb query synthesis")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wetwire-honeycomb <command> [flags] [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  build     Synthesize queries to Query JSON")
	fmt.Println("  lint      Check queries for issues")
	fmt.Println("  list      List all discovered queries")
	fmt.Println("  import    Convert Query JSON to Go code")
	fmt.Println("  validate  Validate Query JSON against Honeycomb constraints")
	fmt.Println("  init      Initialize a new queries directory")
	fmt.Println("  graph     Show query dependency graph")
	fmt.Println("  diff      Compare generated output vs existing config")
	fmt.Println("  watch     Auto-rebuild on source file changes")
	fmt.Println("  design    AI-assisted query generation (requires API key)")
	fmt.Println("  test      Persona-based testing with scoring")
	fmt.Println("  version   Print version information")
	fmt.Println("  help      Print this help message")
	fmt.Println()
	fmt.Println("Run 'wetwire-honeycomb <command> -h' for command-specific help.")
}

func buildCmd(args []string) int {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	output := fs.String("o", "", "Output file or directory")
	format := fs.String("f", "json", "Output format: json, pretty")
	stdout := fs.Bool("stdout", false, "Write to stdout instead of file")
	verbose := fs.Bool("v", false, "Verbose output")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb build [flags] [packages]")
		fmt.Println()
		fmt.Println("Synthesize queries to Query JSON.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Build queries
	b, err := builder.NewBuilder(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	result, err := b.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
		return 1
	}

	if *verbose {
		fmt.Printf("Discovered %d queries\n", result.QueryCount())
	}

	if result.QueryCount() == 0 {
		if *verbose {
			fmt.Println("No queries found")
		}
		return 0
	}

	// Serialize queries
	queries := result.Queries()
	var jsonData []byte

	if len(queries) == 1 {
		q := discoveredToQuery(queries[0])
		if *format == "pretty" {
			jsonData, err = serialize.ToJSONPretty(q)
		} else {
			jsonData, err = serialize.ToJSON(q)
		}
	} else {
		// Multiple queries - wrap in object
		queryMap := make(map[string]json.RawMessage)
		for _, dq := range queries {
			q := discoveredToQuery(dq)
			var data []byte
			data, err = serialize.ToJSON(q)
			if err != nil {
				break
			}
			queryMap[dq.Name] = data
		}
		if err == nil {
			if *format == "pretty" {
				jsonData, err = json.MarshalIndent(queryMap, "", "  ")
			} else {
				jsonData, err = json.Marshal(queryMap)
			}
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Serialization failed: %v\n", err)
		return 1
	}

	// Output
	if *stdout || *output == "" {
		fmt.Println(string(jsonData))
	} else {
		if err := os.WriteFile(*output, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write output: %v\n", err)
			return 1
		}
		if *verbose {
			fmt.Printf("Wrote %d bytes to %s\n", len(jsonData), *output)
		}
	}

	return 0
}

func lintCmd(args []string) int {
	fs := flag.NewFlagSet("lint", flag.ExitOnError)
	format := fs.String("f", "text", "Output format: text, json")
	fix := fs.Bool("fix", false, "Auto-fix issues where possible")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb lint [flags] [packages]")
		fmt.Println()
		fmt.Println("Check queries for issues.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	_ = fix // TODO: implement auto-fix

	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Discover queries
	queries, err := discovery.DiscoverQueries(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Run lint
	results := lint.LintQueries(queries)

	if len(results) == 0 {
		fmt.Println("No issues found")
		return 0
	}

	// Output results
	if *format == "json" {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	} else {
		for _, r := range results {
			severity := "warning"
			if r.Severity == "error" {
				severity = "error"
			}
			fmt.Printf("%s:%d: %s: [%s] %s (%s)\n",
				r.File, r.Line, severity, r.Rule, r.Message, r.Query)
		}
	}

	// Exit 1 if there are errors
	for _, r := range results {
		if r.Severity == "error" {
			return 1
		}
	}

	return 0
}

func listCmd(args []string) int {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	format := fs.String("f", "text", "Output format: text, json")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb list [flags] [packages]")
		fmt.Println()
		fmt.Println("List all discovered queries.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Discover queries
	queries, err := discovery.DiscoverQueries(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if *format == "json" {
		data, _ := json.MarshalIndent(queries, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Found %d queries:\n", len(queries))
		for _, q := range queries {
			fmt.Printf("  %s (%s:%d) - dataset: %s\n",
				q.Name, filepath.Base(q.File), q.Line, q.Dataset)
		}
	}

	return 0
}

func importCmd(args []string) int {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	output := fs.String("o", "", "Output file (default: stdout)")
	pkg := fs.String("p", "queries", "Package name for generated code")
	name := fs.String("n", "Query", "Variable name for the query")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb import [flags] <file.json>")
		fmt.Println()
		fmt.Println("Convert Query JSON to Go code.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: JSON file required")
		fs.Usage()
		return 1
	}

	inputFile := fs.Arg(0)

	// Read JSON file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return 1
	}

	// Parse JSON
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		return 1
	}

	// Generate Go code
	goCode := generateGoCode(*pkg, *name, raw)

	// Output
	if *output == "" {
		fmt.Print(goCode)
	} else {
		if err := os.WriteFile(*output, []byte(goCode), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			return 1
		}
	}

	return 0
}

// discoveredToQuery converts a DiscoveredQuery to a query.Query
func discoveredToQuery(dq discovery.DiscoveredQuery) query.Query {
	q := query.Query{
		Dataset: dq.Dataset,
		TimeRange: query.TimeRange{
			TimeRange: dq.TimeRange.TimeRange,
			StartTime: dq.TimeRange.StartTime,
			EndTime:   dq.TimeRange.EndTime,
		},
		Breakdowns: dq.Breakdowns,
		Limit:      dq.Limit,
	}

	for _, c := range dq.Calculations {
		q.Calculations = append(q.Calculations, query.Calculation{
			Op:     c.Op,
			Column: c.Column,
		})
	}

	for _, f := range dq.Filters {
		q.Filters = append(q.Filters, query.Filter{
			Column: f.Column,
			Op:     f.Op,
			Value:  f.Value,
		})
	}

	return q
}

// generateGoCode generates Go code from a parsed query JSON
func generateGoCode(pkg, name string, raw map[string]any) string {
	var code string
	code += fmt.Sprintf("package %s\n\n", pkg)
	code += "import \"github.com/lex00/wetwire-honeycomb-go/query\"\n\n"
	code += fmt.Sprintf("var %s = query.Query{\n", name)

	// Time range
	if tr, ok := raw["time_range"].(float64); ok {
		hours := int(tr) / 3600
		if hours*3600 == int(tr) {
			code += fmt.Sprintf("\tTimeRange: query.Hours(%d),\n", hours)
		} else {
			code += fmt.Sprintf("\tTimeRange: query.Seconds(%d),\n", int(tr))
		}
	}

	// Breakdowns
	if breakdowns, ok := raw["breakdowns"].([]any); ok && len(breakdowns) > 0 {
		code += "\tBreakdowns: []string{"
		for i, b := range breakdowns {
			if i > 0 {
				code += ", "
			}
			code += fmt.Sprintf("%q", b)
		}
		code += "},\n"
	}

	// Calculations
	if calcs, ok := raw["calculations"].([]any); ok && len(calcs) > 0 {
		code += "\tCalculations: []query.Calculation{\n"
		for _, c := range calcs {
			cm := c.(map[string]any)
			op := cm["op"].(string)
			col, hasCol := cm["column"].(string)

			if op == "COUNT" && !hasCol {
				code += "\t\tquery.Count(),\n"
			} else if hasCol {
				switch op {
				case "P99":
					code += fmt.Sprintf("\t\tquery.P99(%q),\n", col)
				case "P95":
					code += fmt.Sprintf("\t\tquery.P95(%q),\n", col)
				case "P90":
					code += fmt.Sprintf("\t\tquery.P90(%q),\n", col)
				case "P75":
					code += fmt.Sprintf("\t\tquery.P75(%q),\n", col)
				case "P50":
					code += fmt.Sprintf("\t\tquery.P50(%q),\n", col)
				case "AVG":
					code += fmt.Sprintf("\t\tquery.Avg(%q),\n", col)
				case "SUM":
					code += fmt.Sprintf("\t\tquery.Sum(%q),\n", col)
				case "MIN":
					code += fmt.Sprintf("\t\tquery.Min(%q),\n", col)
				case "MAX":
					code += fmt.Sprintf("\t\tquery.Max(%q),\n", col)
				case "COUNT":
					code += fmt.Sprintf("\t\tquery.CountDistinct(%q),\n", col)
				default:
					code += fmt.Sprintf("\t\t{Op: %q, Column: %q},\n", op, col)
				}
			}
		}
		code += "\t},\n"
	}

	// Filters
	if filters, ok := raw["filters"].([]any); ok && len(filters) > 0 {
		code += "\tFilters: []query.Filter{\n"
		for _, f := range filters {
			fm := f.(map[string]any)
			col := fm["column"].(string)
			op := fm["op"].(string)
			val := fm["value"]

			switch op {
			case "=":
				code += fmt.Sprintf("\t\tquery.Equals(%q, %v),\n", col, formatValue(val))
			case "!=":
				code += fmt.Sprintf("\t\tquery.NotEquals(%q, %v),\n", col, formatValue(val))
			case ">":
				code += fmt.Sprintf("\t\tquery.GT(%q, %v),\n", col, formatValue(val))
			case ">=":
				code += fmt.Sprintf("\t\tquery.GTE(%q, %v),\n", col, formatValue(val))
			case "<":
				code += fmt.Sprintf("\t\tquery.LT(%q, %v),\n", col, formatValue(val))
			case "<=":
				code += fmt.Sprintf("\t\tquery.LTE(%q, %v),\n", col, formatValue(val))
			default:
				code += fmt.Sprintf("\t\t{Column: %q, Op: %q, Value: %v},\n", col, op, formatValue(val))
			}
		}
		code += "\t},\n"
	}

	// Limit
	if limit, ok := raw["limit"].(float64); ok && limit > 0 {
		code += fmt.Sprintf("\tLimit: %d,\n", int(limit))
	}

	code += "}\n"
	return code
}

func formatValue(v any) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case float64:
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func validateCmd(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	dryRun := fs.Bool("dry-run", true, "Validate structure only, no API calls")
	dataset := fs.String("dataset", "", "Target dataset for column validation")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb validate [flags] [files]")
		fmt.Println()
		fmt.Println("Validate Query JSON against Honeycomb constraints.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	_ = dataset // Used for API validation (future)

	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Discover queries
	queries, err := discovery.DiscoverQueries(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if len(queries) == 0 {
		fmt.Println("No queries found to validate")
		return 0
	}

	// Validate each query
	var errors []string
	for _, q := range queries {
		errs := validateQuery(q, *dryRun)
		for _, e := range errs {
			errors = append(errors, fmt.Sprintf("%s:%d [%s]: %s", q.File, q.Line, q.Name, e))
		}
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, e := range errors {
			fmt.Printf("  %s\n", e)
		}
		return 1
	}

	fmt.Printf("Validated %d queries successfully\n", len(queries))
	return 0
}

func validateQuery(q discovery.DiscoveredQuery, dryRun bool) []string {
	var errors []string

	// Honeycomb constraints
	const (
		maxTimeRangeDays  = 7
		maxBreakdowns     = 100
		maxCalculations   = 100
		maxFilters        = 100
		maxLimitValue     = 10000
	)

	// Time range check (7 days max = 604800 seconds)
	if q.TimeRange.TimeRange > maxTimeRangeDays*86400 {
		errors = append(errors, fmt.Sprintf("time_range exceeds %d days maximum", maxTimeRangeDays))
	}

	// Breakdown limit
	if len(q.Breakdowns) > maxBreakdowns {
		errors = append(errors, fmt.Sprintf("too many breakdowns (%d, max %d)", len(q.Breakdowns), maxBreakdowns))
	}

	// Calculation limit
	if len(q.Calculations) > maxCalculations {
		errors = append(errors, fmt.Sprintf("too many calculations (%d, max %d)", len(q.Calculations), maxCalculations))
	}

	// Filter limit
	if len(q.Filters) > maxFilters {
		errors = append(errors, fmt.Sprintf("too many filters (%d, max %d)", len(q.Filters), maxFilters))
	}

	// Limit value check
	if q.Limit > maxLimitValue {
		errors = append(errors, fmt.Sprintf("limit exceeds maximum (%d, max %d)", q.Limit, maxLimitValue))
	}

	// Required fields
	if q.Dataset == "" {
		errors = append(errors, "missing required field: dataset")
	}

	if q.TimeRange.TimeRange == 0 && q.TimeRange.StartTime == 0 {
		errors = append(errors, "missing required field: time_range")
	}

	if len(q.Calculations) == 0 {
		errors = append(errors, "missing required field: calculations")
	}

	return errors
}

func initCmd(args []string) int {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	pkgName := fs.String("p", "queries", "Package name")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb init [flags] [directory]")
		fmt.Println()
		fmt.Println("Initialize a new queries directory with example files.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	dir := "queries"
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
	}

	// Create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		return 1
	}

	// Create example query file
	exampleFile := filepath.Join(dir, "example.go")
	exampleContent := fmt.Sprintf(`package %s

import "github.com/lex00/wetwire-honeycomb-go/query"

// SlowRequests finds requests taking longer than 500ms
var SlowRequests = query.Query{
	Dataset:   "production", // TODO: Change to your dataset
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

// ErrorRate tracks error percentage by service
var ErrorRate = query.Query{
	Dataset:   "production", // TODO: Change to your dataset
	TimeRange: query.Hours(1),
	Breakdowns: []string{"service"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("status_code", 400),
	},
}
`, *pkgName)

	if err := os.WriteFile(exampleFile, []byte(exampleContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing example file: %v\n", err)
		return 1
	}

	fmt.Printf("Initialized queries directory: %s\n", dir)
	fmt.Printf("Created example file: %s\n", exampleFile)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit the example queries or create new ones")
	fmt.Println("  2. Run 'wetwire-honeycomb build ./" + dir + "' to generate JSON")
	fmt.Println("  3. Run 'wetwire-honeycomb lint ./" + dir + "' to check for issues")

	return 0
}

func graphCmd(args []string) int {
	fs := flag.NewFlagSet("graph", flag.ExitOnError)
	format := fs.String("f", "text", "Output format: text, dot")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb graph [flags] [packages]")
		fmt.Println()
		fmt.Println("Show query relationships and dependencies.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Discover queries
	queries, err := discovery.DiscoverQueries(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	if len(queries) == 0 {
		fmt.Println("No queries found")
		return 0
	}

	// Group by dataset
	byDataset := make(map[string][]discovery.DiscoveredQuery)
	for _, q := range queries {
		byDataset[q.Dataset] = append(byDataset[q.Dataset], q)
	}

	if *format == "dot" {
		// DOT format for graphviz
		fmt.Println("digraph queries {")
		fmt.Println("  rankdir=LR;")
		fmt.Println("  node [shape=box];")
		fmt.Println()

		for dataset, dqs := range byDataset {
			fmt.Printf("  subgraph cluster_%s {\n", sanitizeID(dataset))
			fmt.Printf("    label=%q;\n", dataset)
			for _, q := range dqs {
				fmt.Printf("    %s [label=%q];\n", sanitizeID(q.Name), q.Name)
			}
			fmt.Println("  }")
		}

		fmt.Println("}")
	} else {
		// Text format
		fmt.Printf("Query Graph (%d queries, %d datasets)\n\n", len(queries), len(byDataset))

		for dataset, dqs := range byDataset {
			fmt.Printf("Dataset: %s\n", dataset)
			for _, q := range dqs {
				fmt.Printf("  ├── %s (%s:%d)\n", q.Name, filepath.Base(q.File), q.Line)
				if len(q.Breakdowns) > 0 {
					fmt.Printf("  │   └── breakdowns: %v\n", q.Breakdowns)
				}
			}
			fmt.Println()
		}
	}

	return 0
}

func sanitizeID(s string) string {
	result := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result += string(r)
		} else {
			result += "_"
		}
	}
	return result
}

func diffCmd(args []string) int {
	fs := flag.NewFlagSet("diff", flag.ExitOnError)
	output := fs.String("output", "", "JSON file to compare against")
	semantic := fs.Bool("semantic", false, "Compare semantic structure instead of text")
	verbose := fs.Bool("v", false, "Verbose output")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb diff [flags] [packages]")
		fmt.Println()
		fmt.Println("Compare generated output vs existing config.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *output == "" {
		fmt.Fprintln(os.Stderr, "Error: --output flag is required")
		fs.Usage()
		return 2
	}

	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	// Build queries
	b, err := builder.NewBuilder(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	result, err := b.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
		return 2
	}

	if result.QueryCount() == 0 {
		fmt.Fprintln(os.Stderr, "No queries found")
		return 2
	}

	// Generate current output
	queries := result.Queries()
	var currentJSON []byte

	if len(queries) == 1 {
		q := discoveredToQuery(queries[0])
		currentJSON, err = serialize.ToJSONPretty(q)
	} else {
		queryMap := make(map[string]json.RawMessage)
		for _, dq := range queries {
			q := discoveredToQuery(dq)
			data, e := serialize.ToJSON(q)
			if e != nil {
				err = e
				break
			}
			queryMap[dq.Name] = data
		}
		if err == nil {
			currentJSON, err = json.MarshalIndent(queryMap, "", "  ")
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Serialization failed: %v\n", err)
		return 2
	}

	// Read existing file
	existingJSON, err := os.ReadFile(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", *output, err)
		return 2
	}

	// Compare
	if *semantic {
		return semanticDiff(currentJSON, existingJSON, *verbose)
	}
	return textDiff(currentJSON, existingJSON, *output, *verbose)
}

func textDiff(current, existing []byte, filename string, verbose bool) int {
	// Normalize line endings
	current = bytes.ReplaceAll(current, []byte("\r\n"), []byte("\n"))
	existing = bytes.ReplaceAll(existing, []byte("\r\n"), []byte("\n"))

	if bytes.Equal(current, existing) {
		if verbose {
			fmt.Println("Files are identical")
		}
		return 0
	}

	// Show line-by-line diff
	currentLines := strings.Split(string(current), "\n")
	existingLines := strings.Split(string(existing), "\n")

	fmt.Printf("--- %s (existing)\n", filename)
	fmt.Println("+++ generated")

	maxLen := len(currentLines)
	if len(existingLines) > maxLen {
		maxLen = len(existingLines)
	}

	for i := 0; i < maxLen; i++ {
		var currLine, existLine string
		if i < len(currentLines) {
			currLine = currentLines[i]
		}
		if i < len(existingLines) {
			existLine = existingLines[i]
		}

		if currLine != existLine {
			if existLine != "" {
				fmt.Printf("-%s\n", existLine)
			}
			if currLine != "" {
				fmt.Printf("+%s\n", currLine)
			}
		}
	}

	return 1
}

func semanticDiff(current, existing []byte, verbose bool) int {
	var currData, existData interface{}

	if err := json.Unmarshal(current, &currData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing generated JSON: %v\n", err)
		return 2
	}

	if err := json.Unmarshal(existing, &existData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing existing JSON: %v\n", err)
		return 2
	}

	if reflect.DeepEqual(currData, existData) {
		if verbose {
			fmt.Println("Semantically identical")
		}
		return 0
	}

	// Show structural differences
	diffs := compareJSON(currData, existData, "")
	for _, d := range diffs {
		fmt.Println(d)
	}

	return 1
}

func compareJSON(a, b interface{}, path string) []string {
	var diffs []string

	switch aTyped := a.(type) {
	case map[string]interface{}:
		bTyped, ok := b.(map[string]interface{})
		if !ok {
			return []string{fmt.Sprintf("Type mismatch at %s: map vs %T", path, b)}
		}

		// Check keys in a
		for k, av := range aTyped {
			newPath := path + "." + k
			if path == "" {
				newPath = k
			}
			if bv, ok := bTyped[k]; ok {
				diffs = append(diffs, compareJSON(av, bv, newPath)...)
			} else {
				diffs = append(diffs, fmt.Sprintf("Key missing in existing: %s", newPath))
			}
		}

		// Check keys in b not in a
		for k := range bTyped {
			newPath := path + "." + k
			if path == "" {
				newPath = k
			}
			if _, ok := aTyped[k]; !ok {
				diffs = append(diffs, fmt.Sprintf("Extra key in existing: %s", newPath))
			}
		}

	case []interface{}:
		bTyped, ok := b.([]interface{})
		if !ok {
			return []string{fmt.Sprintf("Type mismatch at %s: array vs %T", path, b)}
		}

		if len(aTyped) != len(bTyped) {
			diffs = append(diffs, fmt.Sprintf("Array length mismatch at %s: %d vs %d", path, len(aTyped), len(bTyped)))
		}

		minLen := len(aTyped)
		if len(bTyped) < minLen {
			minLen = len(bTyped)
		}

		for i := 0; i < minLen; i++ {
			diffs = append(diffs, compareJSON(aTyped[i], bTyped[i], fmt.Sprintf("%s[%d]", path, i))...)
		}

	default:
		if !reflect.DeepEqual(a, b) {
			diffs = append(diffs, fmt.Sprintf("Value mismatch at %s: %v vs %v", path, a, b))
		}
	}

	return diffs
}

func watchCmd(args []string) int {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	output := fs.String("output", "", "Output file")
	interval := fs.Int("interval", 2, "Polling interval in seconds")
	verbose := fs.Bool("v", false, "Verbose output")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb watch [flags] [packages]")
		fmt.Println()
		fmt.Println("Auto-rebuild on source file changes.")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	path := "."
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	fmt.Printf("Watching %s for changes (interval: %ds)\n", path, *interval)
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	var lastModTime time.Time
	var lastHash string

	for {
		// Get current modification state
		currentModTime, currentHash, err := getDirectoryState(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking files: %v\n", err)
			time.Sleep(time.Duration(*interval) * time.Second)
			continue
		}

		// Check if anything changed
		if !currentModTime.Equal(lastModTime) || currentHash != lastHash {
			if lastModTime.IsZero() {
				fmt.Printf("[%s] Initial build\n", time.Now().Format("15:04:05"))
			} else {
				fmt.Printf("[%s] Changes detected, rebuilding...\n", time.Now().Format("15:04:05"))
			}

			// Build
			b, err := builder.NewBuilder(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  Error: %v\n", err)
			} else {
				result, err := b.Build()
				if err != nil {
					fmt.Fprintf(os.Stderr, "  Build failed: %v\n", err)
				} else {
					if *verbose {
						fmt.Printf("  Found %d queries\n", result.QueryCount())
					}

					if result.QueryCount() > 0 && *output != "" {
						// Write output
						queries := result.Queries()
						var jsonData []byte

						if len(queries) == 1 {
							q := discoveredToQuery(queries[0])
							jsonData, err = serialize.ToJSONPretty(q)
						} else {
							queryMap := make(map[string]json.RawMessage)
							for _, dq := range queries {
								q := discoveredToQuery(dq)
								data, e := serialize.ToJSON(q)
								if e != nil {
									err = e
									break
								}
								queryMap[dq.Name] = data
							}
							if err == nil {
								jsonData, err = json.MarshalIndent(queryMap, "", "  ")
							}
						}

						if err != nil {
							fmt.Fprintf(os.Stderr, "  Serialization failed: %v\n", err)
						} else {
							if err := os.WriteFile(*output, jsonData, 0644); err != nil {
								fmt.Fprintf(os.Stderr, "  Failed to write output: %v\n", err)
							} else {
								fmt.Printf("  Wrote %s (%d bytes)\n", *output, len(jsonData))
							}
						}
					} else if result.QueryCount() > 0 {
						fmt.Printf("  Build succeeded (%d queries)\n", result.QueryCount())
					} else {
						fmt.Println("  No queries found")
					}
				}
			}

			lastModTime = currentModTime
			lastHash = currentHash
		}

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

func getDirectoryState(dir string) (time.Time, string, error) {
	var latestTime time.Time
	var fileList []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and non-Go files
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
		}

		fileList = append(fileList, fmt.Sprintf("%s:%d", path, info.ModTime().UnixNano()))
		return nil
	})

	if err != nil {
		return latestTime, "", err
	}

	// Create a simple hash of file states
	hash := strings.Join(fileList, "|")
	return latestTime, hash, nil
}

func designCmd(args []string) int {
	fs := flag.NewFlagSet("design", flag.ExitOnError)
	output := fs.String("o", ".", "Output directory for generated files")
	provider := fs.String("provider", "anthropic", "AI provider: anthropic, openai, gemini")
	model := fs.String("model", "claude-sonnet-4-20250514", "Model to use")
	maxLintCycles := fs.Int("max-lint-cycles", 5, "Maximum lint/fix cycles")
	stream := fs.Bool("stream", false, "Stream AI responses")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb design [flags] <prompt>")
		fmt.Println()
		fmt.Println("AI-assisted query generation from natural language.")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-honeycomb design \"show me P99 latency by endpoint for the last 2 hours\"")
		fmt.Println("  wetwire-honeycomb design --provider openai \"find slow database queries\"")
		fmt.Println("  wetwire-honeycomb design -o ./queries \"error rate by service\"")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Validate provider
	if *provider != "anthropic" && *provider != "openai" && *provider != "gemini" {
		fmt.Fprintf(os.Stderr, "Error: unsupported provider %q (supported: anthropic, openai, gemini)\n", *provider)
		return 1
	}

	// Get prompt
	var prompt string
	if fs.NArg() > 0 {
		prompt = strings.Join(fs.Args(), " ")
	} else {
		// Interactive mode - read from stdin
		fmt.Println("Enter your query description (Ctrl+D to finish):")
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		prompt = strings.Join(lines, "\n")
	}

	if strings.TrimSpace(prompt) == "" {
		fmt.Fprintln(os.Stderr, "Error: prompt required")
		fs.Usage()
		return 1
	}

	// Check for API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if *provider == "anthropic" && apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ANTHROPIC_API_KEY environment variable required")
		return 1
	}

	// Create output directory if needed
	if err := os.MkdirAll(*output, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		return 1
	}

	// Create provider
	var prov providers.Provider
	var err error
	switch *provider {
	case "anthropic":
		prov, err = anthropic.New(anthropic.Config{APIKey: apiKey})
	default:
		fmt.Fprintf(os.Stderr, "Error: provider %q not yet implemented\n", *provider)
		return 1
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating provider: %v\n", err)
		return 1
	}

	// Create runner agent with Honeycomb domain
	config := agents.RunnerConfig{
		Provider:      prov,
		Domain:        agent.HoneycombDomain(),
		WorkDir:       *output,
		Model:         *model,
		MaxLintCycles: *maxLintCycles,
		Developer:     newConsoleDeveloper(),
	}

	if *stream {
		config.StreamHandler = func(text string) {
			fmt.Print(text)
		}
	}

	runner, err := agents.NewRunnerAgent(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating agent: %v\n", err)
		return 1
	}

	// Run the agent
	ctx := context.Background()
	fmt.Printf("Designing Honeycomb queries for: %s\n\n", prompt)

	if err := runner.Run(ctx, prompt); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Report summary
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  Generated files: %d\n", len(runner.GetGeneratedFiles()))
	for _, f := range runner.GetGeneratedFiles() {
		fmt.Printf("    - %s\n", f)
	}
	fmt.Printf("  Lint cycles: %d\n", runner.GetLintCycles())
	fmt.Printf("  Lint passed: %t\n", runner.LintPassed())

	if !runner.LintPassed() {
		return 1
	}
	return 0
}

// consoleDeveloper implements orchestrator.Developer for interactive sessions
type consoleDeveloper struct {
	reader *bufio.Reader
}

func newConsoleDeveloper() orchestrator.Developer {
	return &consoleDeveloper{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (d *consoleDeveloper) Respond(_ context.Context, question string) (string, error) {
	fmt.Printf("\n[Agent Question] %s\n> ", question)
	answer, err := d.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}

func testCmd(args []string) int {
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	format := fs.String("f", "text", "Output format: text, json")
	persona := fs.String("persona", "intermediate", "Persona: beginner, intermediate, expert, terse, verbose")
	allPersonas := fs.Bool("all-personas", false, "Run with all personas")
	scenario := fs.String("scenario", "default", "Scenario name for tracking")
	provider := fs.String("provider", "anthropic", "AI provider: anthropic, openai, gemini")
	showScore := fs.Bool("score", false, "Show scoring breakdown")
	listOpts := fs.Bool("list", false, "List available personas and scenarios")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-honeycomb test [flags] <prompt>")
		fmt.Println()
		fmt.Println("Persona-based query testing with scoring.")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-honeycomb test --persona beginner \"Create a query to find slow requests\"")
		fmt.Println("  wetwire-honeycomb test --persona expert \"Build an SLI dashboard query set\"")
		fmt.Println("  wetwire-honeycomb test --all-personas \"Create error tracking queries\"")
		fmt.Println("  wetwire-honeycomb test --list")
		fmt.Println()
		fmt.Println("Personas:")
		fmt.Println("  beginner     - New to Honeycomb, needs guidance")
		fmt.Println("  intermediate - Some experience, knows basics")
		fmt.Println("  expert       - Deep observability knowledge")
		fmt.Println("  terse        - Minimal words, expects inference")
		fmt.Println("  verbose      - Over-explains, buries requirements")
		fmt.Println()
		fmt.Println("Scoring (0-15 points):")
		fmt.Println("  Completeness       - Were all required queries generated?")
		fmt.Println("  Lint Quality       - How many lint cycles needed?")
		fmt.Println("  Code Quality       - Idiomatic patterns used?")
		fmt.Println("  Output Validity    - Valid Query JSON produced?")
		fmt.Println("  Question Efficiency - Appropriate clarifications?")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *listOpts {
		fmt.Println("Available Personas:")
		for _, p := range personas.All() {
			fmt.Printf("  %-15s %s\n", p.Name, p.Description)
		}
		fmt.Println()
		fmt.Println("Available Scenarios:")
		fmt.Println("  default        - General query generation")
		fmt.Println("  latency        - Latency and performance queries")
		fmt.Println("  errors         - Error tracking and debugging")
		fmt.Println("  traces         - Distributed tracing queries")
		return 0
	}

	// Validate provider
	if *provider != "anthropic" && *provider != "openai" && *provider != "gemini" {
		fmt.Fprintf(os.Stderr, "Error: unsupported provider %q (supported: anthropic, openai, gemini)\n", *provider)
		return 1
	}

	// Validate persona
	if !*allPersonas {
		if _, err := personas.Get(*persona); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid persona %q\n", *persona)
			fmt.Fprintln(os.Stderr, "Available: beginner, intermediate, expert, terse, verbose")
			return 1
		}
	}

	// Get prompt
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: prompt required")
		fs.Usage()
		return 1
	}
	prompt := strings.Join(fs.Args(), " ")

	// Check for API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if *provider == "anthropic" && apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ANTHROPIC_API_KEY environment variable required")
		return 1
	}

	// Determine which personas to run
	var personasToRun []personas.Persona
	if *allPersonas {
		personasToRun = personas.All()
	} else {
		p, _ := personas.Get(*persona)
		personasToRun = []personas.Persona{p}
	}

	// Run tests for each persona
	type TestResult struct {
		Persona        string         `json:"persona"`
		Scenario       string         `json:"scenario"`
		GeneratedFiles []string       `json:"generated_files"`
		LintCycles     int            `json:"lint_cycles"`
		LintPassed     bool           `json:"lint_passed"`
		Score          *PersonaScore  `json:"score,omitempty"`
		Error          string         `json:"error,omitempty"`
	}

	var results []TestResult

	for _, p := range personasToRun {
		fmt.Printf("\n=== Testing with persona: %s ===\n", p.Name)
		fmt.Printf("Trait: %s\n\n", p.Description)

		// Create temp directory for this test
		tempDir, err := os.MkdirTemp("", fmt.Sprintf("wetwire-honeycomb-test-%s-", p.Name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating temp dir: %v\n", err)
			results = append(results, TestResult{
				Persona:  p.Name,
				Scenario: *scenario,
				Error:    err.Error(),
			})
			continue
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		// Create provider
		var prov providers.Provider
		switch *provider {
		case "anthropic":
			prov, err = anthropic.New(anthropic.Config{APIKey: apiKey})
		default:
			fmt.Fprintf(os.Stderr, "Error: provider %q not yet implemented\n", *provider)
			results = append(results, TestResult{
				Persona:  p.Name,
				Scenario: *scenario,
				Error:    fmt.Sprintf("provider %q not implemented", *provider),
			})
			continue
		}
		if err != nil {
			results = append(results, TestResult{
				Persona:  p.Name,
				Scenario: *scenario,
				Error:    err.Error(),
			})
			continue
		}

		// Modify prompt based on persona
		personaPrompt := applyPersonaStyle(prompt, p)

		// Create runner agent
		config := agents.RunnerConfig{
			Provider:      prov,
			Domain:        agent.HoneycombDomain(),
			WorkDir:       tempDir,
			MaxLintCycles: 5,
		}

		runner, err := agents.NewRunnerAgent(config)
		if err != nil {
			results = append(results, TestResult{
				Persona:  p.Name,
				Scenario: *scenario,
				Error:    err.Error(),
			})
			continue
		}

		// Run the agent
		ctx := context.Background()
		runErr := runner.Run(ctx, personaPrompt)

		result := TestResult{
			Persona:        p.Name,
			Scenario:       *scenario,
			GeneratedFiles: runner.GetGeneratedFiles(),
			LintCycles:     runner.GetLintCycles(),
			LintPassed:     runner.LintPassed(),
		}

		if runErr != nil {
			result.Error = runErr.Error()
		}

		if *showScore {
			result.Score = calculateScore(runner, runErr)
		}

		results = append(results, result)
	}

	// Output results
	if *format == "json" {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("\n=== Test Results ===")
		for _, r := range results {
			status := "PASS"
			if r.Error != "" || !r.LintPassed {
				status = "FAIL"
			}
			fmt.Printf("\n[%s] Persona: %s\n", status, r.Persona)
			fmt.Printf("  Scenario: %s\n", r.Scenario)
			fmt.Printf("  Generated files: %d\n", len(r.GeneratedFiles))
			fmt.Printf("  Lint cycles: %d\n", r.LintCycles)
			fmt.Printf("  Lint passed: %t\n", r.LintPassed)
			if r.Error != "" {
				fmt.Printf("  Error: %s\n", r.Error)
			}
			if r.Score != nil {
				fmt.Printf("  Score: %d/15\n", r.Score.Total)
				fmt.Printf("    Completeness: %d/3\n", r.Score.Completeness)
				fmt.Printf("    Lint Quality: %d/3\n", r.Score.LintQuality)
				fmt.Printf("    Code Quality: %d/3\n", r.Score.CodeQuality)
				fmt.Printf("    Output Validity: %d/3\n", r.Score.OutputValidity)
				fmt.Printf("    Question Efficiency: %d/3\n", r.Score.QuestionEfficiency)
			}
		}

		// Summary
		passed := 0
		for _, r := range results {
			if r.Error == "" && r.LintPassed {
				passed++
			}
		}
		fmt.Printf("\nSummary: %d/%d passed\n", passed, len(results))
	}

	// Return non-zero if any test failed
	for _, r := range results {
		if r.Error != "" || !r.LintPassed {
			return 1
		}
	}
	return 0
}

// PersonaScore represents the scoring breakdown
type PersonaScore struct {
	Completeness       int `json:"completeness"`
	LintQuality        int `json:"lint_quality"`
	CodeQuality        int `json:"code_quality"`
	OutputValidity     int `json:"output_validity"`
	QuestionEfficiency int `json:"question_efficiency"`
	Total              int `json:"total"`
}

func calculateScore(runner *agents.RunnerAgent, runErr error) *PersonaScore {
	score := &PersonaScore{}

	// Completeness (0-3): Were files generated?
	files := runner.GetGeneratedFiles()
	switch {
	case len(files) >= 3:
		score.Completeness = 3
	case len(files) >= 2:
		score.Completeness = 2
	case len(files) >= 1:
		score.Completeness = 1
	}

	// Lint Quality (0-3): How many cycles needed?
	cycles := runner.GetLintCycles()
	switch {
	case runner.LintPassed() && cycles <= 1:
		score.LintQuality = 3
	case runner.LintPassed() && cycles <= 3:
		score.LintQuality = 2
	case runner.LintPassed():
		score.LintQuality = 1
	}

	// Code Quality (0-3): No error and lint passed
	if runErr == nil && runner.LintPassed() {
		score.CodeQuality = 3
	} else if runErr == nil {
		score.CodeQuality = 1
	}

	// Output Validity (0-3): Template generated successfully
	if runner.GetTemplate() != "" {
		score.OutputValidity = 3
	} else if runner.LintPassed() {
		score.OutputValidity = 1
	}

	// Question Efficiency (0-3): Static analysis, no questions
	score.QuestionEfficiency = 2 // Default for automated testing

	score.Total = score.Completeness + score.LintQuality + score.CodeQuality +
		score.OutputValidity + score.QuestionEfficiency

	return score
}

func applyPersonaStyle(prompt string, p personas.Persona) string {
	switch p.Name {
	case "beginner":
		return fmt.Sprintf("I'm new to Honeycomb queries. Can you help me with this? %s", prompt)
	case "terse":
		// Keep it minimal
		return prompt
	case "verbose":
		return fmt.Sprintf("I need your help creating some Honeycomb queries. Let me explain in detail what I'm looking for. %s. Please make sure to include all the necessary fields and use best practices.", prompt)
	case "expert":
		return fmt.Sprintf("Create optimized Honeycomb queries with proper breakdowns, orders, and limits: %s", prompt)
	default:
		return prompt
	}
}
