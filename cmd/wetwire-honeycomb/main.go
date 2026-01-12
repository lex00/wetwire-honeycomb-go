// Command wetwire-honeycomb generates Honeycomb Query JSON from Go query declarations.
//
// Usage:
//
//	wetwire-honeycomb build ./queries/...   Generate Query JSON
//	wetwire-honeycomb lint ./queries/...    Check for issues
//	wetwire-honeycomb validate ./queries/...Validate queries
//	wetwire-honeycomb list ./queries/...    List discovered queries
//	wetwire-honeycomb graph ./queries/...   Generate dependency graph
//	wetwire-honeycomb init myqueries        Create new queries directory
//	wetwire-honeycomb import query.json     Import Query JSON to Go
//	wetwire-honeycomb design "prompt"       AI-assisted query design
//	wetwire-honeycomb test "prompt"         Run persona-based testing
//	wetwire-honeycomb diff old.json new.json Compare two query files
//	wetwire-honeycomb watch ./queries/...   Auto-rebuild on file changes
//	wetwire-honeycomb version               Show version
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/lex00/wetwire-honeycomb-go/internal/builder"
	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
	"github.com/lex00/wetwire-honeycomb-go/internal/lint"
	"github.com/lex00/wetwire-honeycomb-go/internal/serialize"
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/spf13/cobra"
)

const version = "0.2.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "wetwire-honeycomb",
		Short: "Generate Honeycomb Query JSON from Go",
		Long: `wetwire-honeycomb generates Honeycomb Query JSON from Go query declarations.

Define your queries using native Go syntax:

    var SlowRequests = query.Query{
        Dataset:   "production",
        TimeRange: query.Hours(2),
        Breakdowns: []string{"endpoint"},
        Calculations: []query.Calculation{
            query.P99("duration_ms"),
        },
    }

Then generate Query JSON:

    wetwire-honeycomb build ./queries/...`,
		Version: version,
	}

	rootCmd.AddCommand(
		newBuildCmd(),
		newLintCmd(),
		newListCmd(),
		newImportCmd(),
		newValidateCmd(),
		newInitCmd(),
		newGraphCmd(),
		newDiffCmd(),
		newWatchCmd(),
		newDesignCmd(),
		newTestCmd(),
		newVersionCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("wetwire-honeycomb %s\n", version)
		},
	}
}

func newBuildCmd() *cobra.Command {
	var output string
	var format string
	var stdout bool
	var verbose bool

	cmd := &cobra.Command{
		Use:   "build [packages]",
		Short: "Synthesize queries to Query JSON",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Build queries
			b, err := builder.NewBuilder(path)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}

			result, err := b.Build()
			if err != nil {
				return fmt.Errorf("build failed: %w", err)
			}

			if verbose {
				fmt.Printf("Discovered %d queries\n", result.QueryCount())
			}

			if result.QueryCount() == 0 {
				if verbose {
					fmt.Println("No queries found")
				}
				return nil
			}

			// Serialize queries
			queries := result.Queries()
			var jsonData []byte

			if len(queries) == 1 {
				q := discoveredToQuery(queries[0])
				if format == "pretty" {
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
					if format == "pretty" {
						jsonData, err = json.MarshalIndent(queryMap, "", "  ")
					} else {
						jsonData, err = json.Marshal(queryMap)
					}
				}
			}

			if err != nil {
				return fmt.Errorf("serialization failed: %w", err)
			}

			// Output
			if stdout || output == "" {
				fmt.Println(string(jsonData))
			} else {
				if err := os.WriteFile(output, jsonData, 0644); err != nil {
					return fmt.Errorf("failed to write output: %w", err)
				}
				if verbose {
					fmt.Printf("Wrote %d bytes to %s\n", len(jsonData), output)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file or directory")
	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format: json, pretty")
	cmd.Flags().BoolVar(&stdout, "stdout", false, "Write to stdout instead of file")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func newLintCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "lint [packages]",
		Short: "Check queries for issues",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Discover queries
			queries, err := discovery.DiscoverQueries(path)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}

			// Run lint
			results := lint.LintQueries(queries)

			if len(results) == 0 {
				fmt.Println("No issues found")
				return nil
			}

			// Output results
			if format == "json" {
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
					return fmt.Errorf("lint errors found")
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")

	return cmd
}

func newListCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "list [packages]",
		Short: "List all discovered queries",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Discover queries
			queries, err := discovery.DiscoverQueries(path)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}

			if format == "json" {
				data, _ := json.MarshalIndent(queries, "", "  ")
				fmt.Println(string(data))
			} else {
				fmt.Printf("Found %d queries:\n", len(queries))
				for _, q := range queries {
					fmt.Printf("  %s (%s:%d) - dataset: %s\n",
						q.Name, filepath.Base(q.File), q.Line, q.Dataset)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, json")

	return cmd
}

func newImportCmd() *cobra.Command {
	var output string
	var pkg string
	var name string

	cmd := &cobra.Command{
		Use:   "import <file.json>",
		Short: "Convert Query JSON to Go code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile := args[0]

			// Read JSON file
			data, err := os.ReadFile(inputFile)
			if err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}

			// Parse JSON
			var raw map[string]any
			if err := json.Unmarshal(data, &raw); err != nil {
				return fmt.Errorf("error parsing JSON: %w", err)
			}

			// Generate Go code
			goCode := generateGoCode(pkg, name, raw)

			// Output
			if output == "" {
				fmt.Print(goCode)
			} else {
				if err := os.WriteFile(output, []byte(goCode), 0644); err != nil {
					return fmt.Errorf("error writing file: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file (default: stdout)")
	cmd.Flags().StringVarP(&pkg, "package", "p", "queries", "Package name for generated code")
	cmd.Flags().StringVarP(&name, "name", "n", "Query", "Variable name for the query")

	return cmd
}

func newValidateCmd() *cobra.Command {
	var dryRun bool
	var dataset string

	cmd := &cobra.Command{
		Use:   "validate [files]",
		Short: "Validate Query JSON against Honeycomb constraints",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			_ = dataset // Used for API validation (future)

			// Discover queries
			queries, err := discovery.DiscoverQueries(path)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}

			if len(queries) == 0 {
				fmt.Println("No queries found to validate")
				return nil
			}

			// Validate each query
			var errors []string
			for _, q := range queries {
				errs := validateQuery(q, dryRun)
				for _, e := range errs {
					errors = append(errors, fmt.Sprintf("%s:%d [%s]: %s", q.File, q.Line, q.Name, e))
				}
			}

			if len(errors) > 0 {
				fmt.Println("Validation errors:")
				for _, e := range errors {
					fmt.Printf("  %s\n", e)
				}
				return fmt.Errorf("validation failed")
			}

			fmt.Printf("Validated %d queries successfully\n", len(queries))
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", true, "Validate structure only, no API calls")
	cmd.Flags().StringVar(&dataset, "dataset", "", "Target dataset for column validation")

	return cmd
}

func newInitCmd() *cobra.Command {
	var pkgName string

	cmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new queries directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "queries"
			if len(args) > 0 {
				dir = args[0]
			}

			// Create directory
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating directory: %w", err)
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
`, pkgName)

			if err := os.WriteFile(exampleFile, []byte(exampleContent), 0644); err != nil {
				return fmt.Errorf("error writing example file: %w", err)
			}

			fmt.Printf("Initialized queries directory: %s\n", dir)
			fmt.Printf("Created example file: %s\n", exampleFile)
			fmt.Println()
			fmt.Println("Next steps:")
			fmt.Println("  1. Edit the example queries or create new ones")
			fmt.Println("  2. Run 'wetwire-honeycomb build ./" + dir + "' to generate JSON")
			fmt.Println("  3. Run 'wetwire-honeycomb lint ./" + dir + "' to check for issues")

			return nil
		},
	}

	cmd.Flags().StringVarP(&pkgName, "package", "p", "queries", "Package name")

	return cmd
}

func newGraphCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "graph [packages]",
		Short: "Show query relationships and dependencies",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Discover queries
			queries, err := discovery.DiscoverQueries(path)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}

			if len(queries) == 0 {
				fmt.Println("No queries found")
				return nil
			}

			// Group by dataset
			byDataset := make(map[string][]discovery.DiscoveredQuery)
			for _, q := range queries {
				byDataset[q.Dataset] = append(byDataset[q.Dataset], q)
			}

			if format == "dot" {
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
						fmt.Printf("  +-- %s (%s:%d)\n", q.Name, filepath.Base(q.File), q.Line)
						if len(q.Breakdowns) > 0 {
							fmt.Printf("  |   +-- breakdowns: %v\n", q.Breakdowns)
						}
					}
					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, dot")

	return cmd
}

func newDiffCmd() *cobra.Command {
	var outputFile string
	var semantic bool
	var verbose bool

	cmd := &cobra.Command{
		Use:   "diff [packages]",
		Short: "Compare generated output vs existing config",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputFile == "" {
				return fmt.Errorf("--output flag is required")
			}

			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Build queries
			b, err := builder.NewBuilder(path)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}

			result, err := b.Build()
			if err != nil {
				return fmt.Errorf("build failed: %w", err)
			}

			if result.QueryCount() == 0 {
				return fmt.Errorf("no queries found")
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
				return fmt.Errorf("serialization failed: %w", err)
			}

			// Read existing file
			existingJSON, err := os.ReadFile(outputFile)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", outputFile, err)
			}

			// Compare
			if semantic {
				return semanticDiff(currentJSON, existingJSON, verbose)
			}
			return textDiff(currentJSON, existingJSON, outputFile, verbose)
		},
	}

	cmd.Flags().StringVar(&outputFile, "output", "", "JSON file to compare against")
	cmd.Flags().BoolVar(&semantic, "semantic", false, "Compare semantic structure instead of text")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func newWatchCmd() *cobra.Command {
	var outputFile string
	var interval int
	var verbose bool

	cmd := &cobra.Command{
		Use:   "watch [packages]",
		Short: "Auto-rebuild on source file changes",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			fmt.Printf("Watching %s for changes (interval: %ds)\n", path, interval)
			fmt.Println("Press Ctrl+C to stop")
			fmt.Println()

			var lastModTime time.Time
			var lastHash string

			for {
				// Get current modification state
				currentModTime, currentHash, err := getDirectoryState(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error checking files: %v\n", err)
					time.Sleep(time.Duration(interval) * time.Second)
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
							if verbose {
								fmt.Printf("  Found %d queries\n", result.QueryCount())
							}

							if result.QueryCount() > 0 && outputFile != "" {
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
									if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
										fmt.Fprintf(os.Stderr, "  Failed to write output: %v\n", err)
									} else {
										fmt.Printf("  Wrote %s (%d bytes)\n", outputFile, len(jsonData))
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

				time.Sleep(time.Duration(interval) * time.Second)
			}
		},
	}

	cmd.Flags().StringVar(&outputFile, "output", "", "Output file")
	cmd.Flags().IntVar(&interval, "interval", 2, "Polling interval in seconds")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

// Helper functions

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

func validateQuery(q discovery.DiscoveredQuery, dryRun bool) []string {
	var errors []string

	// Honeycomb constraints
	const (
		maxTimeRangeDays = 7
		maxBreakdowns    = 100
		maxCalculations  = 100
		maxFilters       = 100
		maxLimitValue    = 10000
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

func textDiff(current, existing []byte, filename string, verbose bool) error {
	// Normalize line endings
	current = bytes.ReplaceAll(current, []byte("\r\n"), []byte("\n"))
	existing = bytes.ReplaceAll(existing, []byte("\r\n"), []byte("\n"))

	if bytes.Equal(current, existing) {
		if verbose {
			fmt.Println("Files are identical")
		}
		return nil
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

	return fmt.Errorf("files differ")
}

func semanticDiff(current, existing []byte, verbose bool) error {
	var currData, existData interface{}

	if err := json.Unmarshal(current, &currData); err != nil {
		return fmt.Errorf("error parsing generated JSON: %w", err)
	}

	if err := json.Unmarshal(existing, &existData); err != nil {
		return fmt.Errorf("error parsing existing JSON: %w", err)
	}

	if reflect.DeepEqual(currData, existData) {
		if verbose {
			fmt.Println("Semantically identical")
		}
		return nil
	}

	// Show structural differences
	diffs := compareJSON(currData, existData, "")
	for _, d := range diffs {
		fmt.Println(d)
	}

	return fmt.Errorf("semantic differences found")
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
