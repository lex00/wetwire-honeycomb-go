package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
