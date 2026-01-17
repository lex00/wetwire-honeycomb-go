// Package domain provides the HoneycombDomain implementation for wetwire-core-go.
package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
	"github.com/lex00/wetwire-honeycomb-go/internal/lint"
	"github.com/lex00/wetwire-honeycomb-go/internal/serialize"
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
	"github.com/spf13/cobra"
)

// Version is set at build time
var Version = "dev"

// Re-export core types for convenience
type (
	Context      = coredomain.Context
	BuildOpts    = coredomain.BuildOpts
	LintOpts     = coredomain.LintOpts
	InitOpts     = coredomain.InitOpts
	ValidateOpts = coredomain.ValidateOpts
	ListOpts     = coredomain.ListOpts
	GraphOpts    = coredomain.GraphOpts
	Result       = coredomain.Result
	Error        = coredomain.Error
)

var (
	NewResult              = coredomain.NewResult
	NewResultWithData      = coredomain.NewResultWithData
	NewErrorResult         = coredomain.NewErrorResult
	NewErrorResultMultiple = coredomain.NewErrorResultMultiple
)

// HoneycombDomain implements the Domain interface for Honeycomb observability.
type HoneycombDomain struct{}

// Compile-time checks
var (
	_ coredomain.Domain        = (*HoneycombDomain)(nil)
	_ coredomain.ListerDomain  = (*HoneycombDomain)(nil)
	_ coredomain.GrapherDomain = (*HoneycombDomain)(nil)
)

// Name returns "honeycomb"
func (d *HoneycombDomain) Name() string {
	return "honeycomb"
}

// Version returns the current version
func (d *HoneycombDomain) Version() string {
	return Version
}

// Builder returns the Honeycomb builder implementation
func (d *HoneycombDomain) Builder() coredomain.Builder {
	return &honeycombBuilder{}
}

// Linter returns the Honeycomb linter implementation
func (d *HoneycombDomain) Linter() coredomain.Linter {
	return &honeycombLinter{}
}

// Initializer returns the Honeycomb initializer implementation
func (d *HoneycombDomain) Initializer() coredomain.Initializer {
	return &honeycombInitializer{}
}

// Validator returns the Honeycomb validator implementation
func (d *HoneycombDomain) Validator() coredomain.Validator {
	return &honeycombValidator{}
}

// Lister returns the Honeycomb lister implementation
func (d *HoneycombDomain) Lister() coredomain.Lister {
	return &honeycombLister{}
}

// Grapher returns the Honeycomb grapher implementation
func (d *HoneycombDomain) Grapher() coredomain.Grapher {
	return &honeycombGrapher{}
}

// CreateRootCommand creates the root command using the domain interface.
func CreateRootCommand(d coredomain.Domain) *cobra.Command {
	return coredomain.Run(d)
}

// honeycombBuilder implements domain.Builder
type honeycombBuilder struct{}

func (b *honeycombBuilder) Build(ctx *Context, path string, opts BuildOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	if resources.TotalCount() == 0 {
		return NewErrorResult("no resources found", Error{
			Path:    absPath,
			Message: "no queries, boards, SLOs, or triggers found",
		}), nil
	}

	// Build output structure
	outputData := make(map[string]json.RawMessage)

	// Filter by type if specified
	resourceType := opts.Type

	// Serialize queries
	if (resourceType == "" || resourceType == "query" || resourceType == "queries") && len(resources.Queries) > 0 {
		queryMap := make(map[string]json.RawMessage)
		for _, dq := range resources.Queries {
			q := discoveredToQuery(dq)
			data, serr := serialize.ToJSON(q)
			if serr != nil {
				return nil, fmt.Errorf("query serialization failed: %w", serr)
			}
			queryMap[dq.Name] = data
		}
		data, _ := json.Marshal(queryMap)
		outputData["queries"] = data
	}

	// Serialize boards
	if (resourceType == "" || resourceType == "board" || resourceType == "boards") && len(resources.Boards) > 0 {
		boardMap := make(map[string]json.RawMessage)
		for _, db := range resources.Boards {
			b := discoveredToBoard(db)
			data, serr := serialize.BoardToJSON(b)
			if serr != nil {
				return nil, fmt.Errorf("board serialization failed: %w", serr)
			}
			boardMap[db.Name] = data
		}
		data, _ := json.Marshal(boardMap)
		outputData["boards"] = data
	}

	// Serialize SLOs
	if (resourceType == "" || resourceType == "slo" || resourceType == "slos") && len(resources.SLOs) > 0 {
		sloMap := make(map[string]json.RawMessage)
		for _, ds := range resources.SLOs {
			s := discoveredToSLO(ds)
			data, serr := serialize.SLOToJSON(s)
			if serr != nil {
				return nil, fmt.Errorf("SLO serialization failed: %w", serr)
			}
			sloMap[ds.Name] = data
		}
		data, _ := json.Marshal(sloMap)
		outputData["slos"] = data
	}

	// Serialize triggers
	if (resourceType == "" || resourceType == "trigger" || resourceType == "triggers") && len(resources.Triggers) > 0 {
		triggerMap := make(map[string]json.RawMessage)
		for _, dt := range resources.Triggers {
			t := discoveredToTrigger(dt)
			data, serr := serialize.TriggerToJSON(t)
			if serr != nil {
				return nil, fmt.Errorf("trigger serialization failed: %w", serr)
			}
			triggerMap[dt.Name] = data
		}
		data, _ := json.Marshal(triggerMap)
		outputData["triggers"] = data
	}

	// Format output
	var jsonData []byte
	if opts.Format == "pretty" {
		jsonData, err = json.MarshalIndent(outputData, "", "  ")
	} else {
		jsonData, err = json.Marshal(outputData)
	}
	if err != nil {
		return nil, fmt.Errorf("serialization failed: %w", err)
	}

	// Handle output file
	if !opts.DryRun && opts.Output != "" {
		if err := os.WriteFile(opts.Output, jsonData, 0644); err != nil {
			return nil, fmt.Errorf("write output: %w", err)
		}
		return NewResult(fmt.Sprintf("Wrote %s", opts.Output)), nil
	}

	return NewResultWithData("Build completed", string(jsonData)), nil
}

// honeycombLinter implements domain.Linter
type honeycombLinter struct{}

func (l *honeycombLinter) Lint(ctx *Context, path string, opts LintOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Run lint on all resources
	results := lint.LintAll(resources)

	if len(results) == 0 {
		return NewResult("No lint issues found"), nil
	}

	// Convert to domain errors
	errs := make([]Error, 0, len(results))
	for _, r := range results {
		errs = append(errs, Error{
			Path:     r.File,
			Line:     r.Line,
			Severity: r.Severity.String(),
			Message:  r.Message,
			Code:     r.Rule,
		})
	}

	return NewErrorResultMultiple("lint issues found", errs), nil
}

// honeycombInitializer implements domain.Initializer
type honeycombInitializer struct{}

func (i *honeycombInitializer) Init(ctx *Context, path string, opts InitOpts) (*Result, error) {
	// Use opts.Path if provided, otherwise fall back to path argument
	targetPath := opts.Path
	if targetPath == "" || targetPath == "." {
		targetPath = path
	}

	// Handle scenario initialization
	if opts.Scenario {
		return i.initScenario(ctx, targetPath, opts)
	}

	// Basic project initialization
	return i.initProject(ctx, targetPath, opts)
}

// initScenario creates a full scenario structure with prompts and expected outputs
func (i *honeycombInitializer) initScenario(ctx *Context, path string, opts InitOpts) (*Result, error) {
	name := opts.Name
	if name == "" {
		name = filepath.Base(path)
	}

	description := opts.Description
	if description == "" {
		description = "Honeycomb observability scenario"
	}

	// Use core's scenario scaffolding
	scenario := coredomain.ScaffoldScenario(name, description, "honeycomb")
	created, err := coredomain.WriteScenario(path, scenario)
	if err != nil {
		return nil, fmt.Errorf("write scenario: %w", err)
	}

	// Create honeycomb-specific expected directories
	expectedDirs := []string{
		filepath.Join(path, "expected", "queries"),
		filepath.Join(path, "expected", "slos"),
		filepath.Join(path, "expected", "triggers"),
		filepath.Join(path, "expected", "boards"),
	}
	for _, dir := range expectedDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	// Create example query in expected/queries/
	exampleQuery := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

// RequestLatency tracks P99/P95/P50 latency across endpoints.
var RequestLatency = query.Query{
	Dataset:   "your-dataset",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"http.route", "http.method"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.P95("duration_ms"),
		query.P50("duration_ms"),
		query.Count(),
	},
	Orders: []query.Order{
		{Op: "P99", Column: "duration_ms", Order: "descending"},
	},
	Limit: 100,
}
`
	queryPath := filepath.Join(path, "expected", "queries", "queries.go")
	if err := os.WriteFile(queryPath, []byte(exampleQuery), 0644); err != nil {
		return nil, fmt.Errorf("write example query: %w", err)
	}
	created = append(created, "expected/queries/queries.go")

	return NewResultWithData(
		fmt.Sprintf("Created scenario %s with %d files", name, len(created)),
		created,
	), nil
}

// initProject creates a basic project with example queries
func (i *honeycombInitializer) initProject(ctx *Context, path string, opts InitOpts) (*Result, error) {
	// Create directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	// Create go.mod
	name := opts.Name
	if name == "" {
		name = filepath.Base(path)
	}
	goMod := fmt.Sprintf(`module %s

go 1.23

require github.com/lex00/wetwire-honeycomb-go v0.0.0
`, name)
	goModPath := filepath.Join(path, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		return nil, fmt.Errorf("write go.mod: %w", err)
	}

	// Create example query file
	exampleContent := `package queries

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
`
	examplePath := filepath.Join(path, "queries.go")
	if err := os.WriteFile(examplePath, []byte(exampleContent), 0644); err != nil {
		return nil, fmt.Errorf("write example: %w", err)
	}

	return NewResultWithData(
		fmt.Sprintf("Created %s with example queries", path),
		[]string{"go.mod", "queries.go"},
	), nil
}

// honeycombValidator implements domain.Validator
type honeycombValidator struct{}

func (v *honeycombValidator) Validate(ctx *Context, path string, opts ValidateOpts) (*Result, error) {
	// For now, validation is the same as lint
	linter := &honeycombLinter{}
	return linter.Lint(ctx, path, LintOpts{})
}

// honeycombLister implements domain.Lister
type honeycombLister struct{}

func (l *honeycombLister) List(ctx *Context, path string, opts ListOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Build list
	list := make([]map[string]string, 0)
	for _, q := range resources.Queries {
		list = append(list, map[string]string{
			"name": q.Name,
			"type": "query",
			"file": q.File,
		})
	}
	for _, b := range resources.Boards {
		list = append(list, map[string]string{
			"name": b.Name,
			"type": "board",
			"file": b.File,
		})
	}
	for _, s := range resources.SLOs {
		list = append(list, map[string]string{
			"name": s.Name,
			"type": "slo",
			"file": s.File,
		})
	}
	for _, t := range resources.Triggers {
		list = append(list, map[string]string{
			"name": t.Name,
			"type": "trigger",
			"file": t.File,
		})
	}

	return NewResultWithData(fmt.Sprintf("Discovered %d resources", len(list)), list), nil
}

// honeycombGrapher implements domain.Grapher
type honeycombGrapher struct{}

func (g *honeycombGrapher) Graph(ctx *Context, path string, opts GraphOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Generate DOT format graph
	var graph string
	switch opts.Format {
	case "dot", "":
		graph = "digraph G {\n"
		for _, q := range resources.Queries {
			graph += fmt.Sprintf("  %s [shape=box];\n", q.Name)
		}
		for _, b := range resources.Boards {
			graph += fmt.Sprintf("  %s [shape=folder];\n", b.Name)
			// Boards reference queries
			for _, q := range resources.Queries {
				graph += fmt.Sprintf("  %s -> %s;\n", b.Name, q.Name)
			}
		}
		graph += "}"
	case "mermaid":
		graph = "graph TD\n"
		for _, q := range resources.Queries {
			graph += fmt.Sprintf("  %s[%s]\n", q.Name, q.Name)
		}
		for _, b := range resources.Boards {
			graph += fmt.Sprintf("  %s{{%s}}\n", b.Name, b.Name)
		}
	default:
		return nil, fmt.Errorf("unknown format: %s", opts.Format)
	}

	return NewResultWithData("Graph generated", graph), nil
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

// discoveredToBoard converts a DiscoveredBoard to a board.Board
func discoveredToBoard(db discovery.DiscoveredBoard) board.Board {
	return board.Board{
		Name:        db.BoardName,
		Description: db.Description,
	}
}

// discoveredToSLO converts a DiscoveredSLO to an slo.SLO
func discoveredToSLO(ds discovery.DiscoveredSLO) slo.SLO {
	return slo.SLO{
		Name:        ds.SLOName,
		Description: ds.Description,
		Dataset:     ds.Dataset,
		Target:      slo.Percentage(ds.TargetPercentage),
		TimePeriod:  slo.Days(ds.TimePeriodDays),
	}
}

// discoveredToTrigger converts a DiscoveredTrigger to a trigger.Trigger
func discoveredToTrigger(dt discovery.DiscoveredTrigger) trigger.Trigger {
	return trigger.Trigger{
		Name:        dt.TriggerName,
		Description: dt.Description,
		Dataset:     dt.Dataset,
		Frequency:   trigger.Seconds(dt.FrequencySeconds),
		Disabled:    dt.Disabled,
	}
}
