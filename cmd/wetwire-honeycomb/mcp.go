// MCP server implementation for embedded design mode.
//
// When design --mcp-server is called, this runs the MCP protocol over stdio,
// providing wetwire_init, wetwire_lint, and wetwire_build tools.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	coremcp "github.com/lex00/wetwire-core-go/mcp"
	"github.com/spf13/cobra"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
	"github.com/lex00/wetwire-honeycomb-go/internal/lint"
	"github.com/lex00/wetwire-honeycomb-go/internal/serialize"
)

// newMCPCmd creates the "mcp" subcommand that runs the MCP server.
func newMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Run MCP server on stdio",
		Long: `Run the Model Context Protocol (MCP) server on stdio transport.

This command starts an MCP server that exposes wetwire-honeycomb tools
for AI assistants to use. The server provides the following tools:
  - wetwire_init: Initialize a new wetwire-honeycomb project
  - wetwire_lint: Lint Go packages for wetwire-honeycomb issues
  - wetwire_build: Generate Query JSON from Go packages
  - wetwire_list: List discovered queries
  - wetwire_graph: Generate dependency graph (DOT/Mermaid)

This is typically used by AI tools and should not be called directly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer()
		},
	}
}

// runMCPServer starts the MCP server on stdio transport.
func runMCPServer() error {
	server := coremcp.NewServer(coremcp.Config{
		Name:    "wetwire-honeycomb",
		Version: "1.0.0",
	})

	// Register tools using core MCP infrastructure
	mcpRegisterStandardTools(server)

	// Run on stdio transport
	return server.Start(context.Background())
}

// mcpRegisterStandardTools registers all standard wetwire tools with the MCP server.
func mcpRegisterStandardTools(server *coremcp.Server) {
	// wetwire_init tool
	server.RegisterToolWithSchema("wetwire_init", "Initialize a new wetwire-honeycomb project with example code", mcpHandleInit, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Project name",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Output directory (default: current directory)",
			},
		},
	})

	// wetwire_build tool
	server.RegisterToolWithSchema("wetwire_build", "Generate JSON from all wetwire resources (queries, boards, SLOs, triggers)", mcpHandleBuild, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"package": map[string]any{
				"type":        "string",
				"description": "Package path to discover resources from",
			},
			"format": map[string]any{
				"type":        "string",
				"enum":        []string{"json", "pretty"},
				"description": "Output format (default: pretty)",
			},
			"type": map[string]any{
				"type":        "string",
				"enum":        []string{"query", "board", "slo", "trigger"},
				"description": "Build only specific resource type (default: all)",
			},
		},
	})

	// wetwire_lint tool
	server.RegisterToolWithSchema("wetwire_lint", "Lint all resources (queries, boards, SLOs, triggers) with WHC rules", mcpHandleLint, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"package": map[string]any{
				"type":        "string",
				"description": "Package path to lint",
			},
			"format": map[string]any{
				"type":        "string",
				"enum":        []string{"text", "json"},
				"description": "Output format (default: text)",
			},
		},
	})

	// wetwire_list tool
	server.RegisterToolWithSchema("wetwire_list", "List all discovered resources (queries, boards, SLOs, triggers)", mcpHandleList, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"package": map[string]any{
				"type":        "string",
				"description": "Package path to analyze",
			},
			"format": map[string]any{
				"type":        "string",
				"enum":        []string{"text", "json"},
				"description": "Output format (default: text)",
			},
			"type": map[string]any{
				"type":        "string",
				"enum":        []string{"query", "board", "slo", "trigger"},
				"description": "Filter by resource type (default: all)",
			},
		},
	})

	// wetwire_graph tool
	server.RegisterToolWithSchema("wetwire_graph", "Visualize query relationships (DOT/Mermaid)", mcpHandleGraph, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"package": map[string]any{
				"type":        "string",
				"description": "Package path to analyze",
			},
			"format": map[string]any{
				"type":        "string",
				"enum":        []string{"text", "dot"},
				"description": "Output format (default: text)",
			},
		},
	})
}

// mcpHandleInit implements the wetwire_init tool.
func mcpHandleInit(_ context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	name, _ := args["name"].(string)

	result := MCPInitResult{Path: path}

	if path == "" {
		result.Error = "path is required"
		return mcpJSONResult(result)
	}

	// Create project directory
	if err := os.MkdirAll(path, 0755); err != nil {
		result.Error = fmt.Sprintf("creating project directory: %v", err)
		return mcpJSONResult(result)
	}

	// Get the module name
	moduleName := name
	if moduleName == "" {
		moduleName = filepath.Base(path)
	}

	// Write go.mod
	goMod := fmt.Sprintf(`module %s

go 1.23

require github.com/lex00/wetwire-honeycomb-go v0.0.0
`, moduleName)

	goModPath := filepath.Join(path, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		result.Error = fmt.Sprintf("writing go.mod: %v", err)
		return mcpJSONResult(result)
	}
	result.Files = append(result.Files, "go.mod")

	// Write queries.go (example queries)
	queriesGo := `package queries

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

	queriesGoPath := filepath.Join(path, "queries.go")
	if err := os.WriteFile(queriesGoPath, []byte(queriesGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing queries.go: %v", err)
		return mcpJSONResult(result)
	}
	result.Files = append(result.Files, "queries.go")

	result.Success = true
	return mcpJSONResult(result)
}

// mcpHandleLint implements the wetwire_lint tool.
func mcpHandleLint(_ context.Context, args map[string]any) (string, error) {
	pkg, _ := args["package"].(string)
	_ = args["format"] // format parameter reserved for future use

	result := MCPLintResult{}

	if pkg == "" {
		result.Issues = append(result.Issues, MCPLintIssue{
			Severity: "error",
			Message:  "package is required",
			RuleID:   "internal",
		})
		return mcpJSONResult(result)
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(pkg)
	if err != nil {
		result.Issues = append(result.Issues, MCPLintIssue{
			Severity: "error",
			Message:  fmt.Sprintf("invalid path: %v", err),
			RuleID:   "internal",
		})
		return mcpJSONResult(result)
	}

	// Verify path exists
	if _, err := os.Stat(absPath); err != nil {
		result.Issues = append(result.Issues, MCPLintIssue{
			Severity: "error",
			Message:  fmt.Sprintf("path not found: %s", pkg),
			RuleID:   "internal",
		})
		return mcpJSONResult(result)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(absPath)
	if err != nil {
		result.Issues = append(result.Issues, MCPLintIssue{
			Severity: "error",
			Message:  fmt.Sprintf("discovery failed: %v", err),
			RuleID:   "internal",
		})
		return mcpJSONResult(result)
	}

	// Run lint on all resources
	lintResults := lint.LintAll(resources)

	// Convert lint results to our format
	for _, issue := range lintResults {
		result.Issues = append(result.Issues, MCPLintIssue{
			RuleID:   issue.Rule,
			Message:  issue.Message,
			File:     issue.File,
			Line:     issue.Line,
			Severity: issue.Severity,
		})
	}

	result.Success = len(result.Issues) == 0
	return mcpJSONResult(result)
}

// mcpHandleBuild implements the wetwire_build tool.
func mcpHandleBuild(_ context.Context, args map[string]any) (string, error) {
	pkg, _ := args["package"].(string)
	format, _ := args["format"].(string)
	resourceType, _ := args["type"].(string)

	result := MCPBuildResult{}

	if pkg == "" {
		result.Errors = append(result.Errors, "package is required")
		return mcpJSONResult(result)
	}

	if format == "" {
		format = "pretty"
	}

	// Convert path to absolute
	absPath, err := filepath.Abs(pkg)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("invalid path: %v", err))
		return mcpJSONResult(result)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(absPath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return mcpJSONResult(result)
	}

	if resources.TotalCount() == 0 {
		result.Errors = append(result.Errors, "no resources found")
		return mcpJSONResult(result)
	}

	// Build output structure
	outputData := make(map[string]json.RawMessage)

	// Serialize queries
	if (resourceType == "" || resourceType == "query" || resourceType == "queries") && len(resources.Queries) > 0 {
		queryMap := make(map[string]json.RawMessage)
		for _, dq := range resources.Queries {
			q := discoveredToQuery(dq)
			data, serr := serialize.ToJSON(q)
			if serr != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("query serialization failed: %v", serr))
				return mcpJSONResult(result)
			}
			queryMap[dq.Name] = data
			result.Queries = append(result.Queries, dq.Name)
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
				result.Errors = append(result.Errors, fmt.Sprintf("board serialization failed: %v", serr))
				return mcpJSONResult(result)
			}
			boardMap[db.Name] = data
			result.Boards = append(result.Boards, db.Name)
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
				result.Errors = append(result.Errors, fmt.Sprintf("SLO serialization failed: %v", serr))
				return mcpJSONResult(result)
			}
			sloMap[ds.Name] = data
			result.SLOs = append(result.SLOs, ds.Name)
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
				result.Errors = append(result.Errors, fmt.Sprintf("trigger serialization failed: %v", serr))
				return mcpJSONResult(result)
			}
			triggerMap[dt.Name] = data
			result.Triggers = append(result.Triggers, dt.Name)
		}
		data, _ := json.Marshal(triggerMap)
		outputData["triggers"] = data
	}

	// Final JSON output
	var jsonData []byte
	if format == "pretty" {
		jsonData, err = json.MarshalIndent(outputData, "", "  ")
	} else {
		jsonData, err = json.Marshal(outputData)
	}
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("serialization failed: %v", err))
		return mcpJSONResult(result)
	}

	result.Success = true
	result.Output = string(jsonData)
	return mcpJSONResult(result)
}

// mcpHandleList implements the wetwire_list tool.
func mcpHandleList(_ context.Context, args map[string]any) (string, error) {
	pkg, _ := args["package"].(string)
	format, _ := args["format"].(string)
	resourceType, _ := args["type"].(string)

	result := MCPListResult{}

	if pkg == "" {
		result.Message = "package is required"
		return mcpJSONResult(result)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(pkg)
	if err != nil {
		result.Message = fmt.Sprintf("discovery failed: %v", err)
		return mcpJSONResult(result)
	}

	// Build query list
	if resourceType == "" || resourceType == "query" || resourceType == "queries" {
		for _, q := range resources.Queries {
			result.Queries = append(result.Queries, MCPQueryInfo{
				Name:    q.Name,
				File:    q.File,
				Line:    q.Line,
				Dataset: q.Dataset,
			})
		}
	}

	// Build board list
	if resourceType == "" || resourceType == "board" || resourceType == "boards" {
		for _, b := range resources.Boards {
			result.Boards = append(result.Boards, MCPBoardInfo{
				Name:       b.Name,
				File:       b.File,
				Line:       b.Line,
				BoardName:  b.BoardName,
				PanelCount: b.PanelCount,
				QueryRefs:  b.QueryRefs,
				SLORefs:    b.SLORefs,
			})
		}
	}

	// Build SLO list
	if resourceType == "" || resourceType == "slo" || resourceType == "slos" {
		for _, s := range resources.SLOs {
			result.SLOs = append(result.SLOs, MCPSLOInfo{
				Name:             s.Name,
				File:             s.File,
				Line:             s.Line,
				SLOName:          s.SLOName,
				Dataset:          s.Dataset,
				TargetPercentage: s.TargetPercentage,
				BurnAlertCount:   s.BurnAlertCount,
			})
		}
	}

	// Build trigger list
	if resourceType == "" || resourceType == "trigger" || resourceType == "triggers" {
		for _, t := range resources.Triggers {
			result.Triggers = append(result.Triggers, MCPTriggerInfo{
				Name:             t.Name,
				File:             t.File,
				Line:             t.Line,
				TriggerName:      t.TriggerName,
				Dataset:          t.Dataset,
				FrequencySeconds: t.FrequencySeconds,
				RecipientCount:   t.RecipientCount,
				Disabled:         t.Disabled,
			})
		}
	}

	result.Success = true
	result.Count = resources.TotalCount()

	if format == "text" {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d resources:\n\n", result.Count))

		if len(result.Queries) > 0 {
			sb.WriteString(fmt.Sprintf("Queries (%d):\n", len(result.Queries)))
			for _, q := range result.Queries {
				sb.WriteString(fmt.Sprintf("  %s (%s:%d) - dataset: %s\n",
					q.Name, filepath.Base(q.File), q.Line, q.Dataset))
			}
			sb.WriteString("\n")
		}

		if len(result.Boards) > 0 {
			sb.WriteString(fmt.Sprintf("Boards (%d):\n", len(result.Boards)))
			for _, b := range result.Boards {
				sb.WriteString(fmt.Sprintf("  %s (%s:%d) - %d panels\n",
					b.Name, filepath.Base(b.File), b.Line, b.PanelCount))
			}
			sb.WriteString("\n")
		}

		if len(result.SLOs) > 0 {
			sb.WriteString(fmt.Sprintf("SLOs (%d):\n", len(result.SLOs)))
			for _, s := range result.SLOs {
				sb.WriteString(fmt.Sprintf("  %s (%s:%d) - %.1f%%\n",
					s.Name, filepath.Base(s.File), s.Line, s.TargetPercentage))
			}
			sb.WriteString("\n")
		}

		if len(result.Triggers) > 0 {
			sb.WriteString(fmt.Sprintf("Triggers (%d):\n", len(result.Triggers)))
			for _, t := range result.Triggers {
				status := "enabled"
				if t.Disabled {
					status = "disabled"
				}
				sb.WriteString(fmt.Sprintf("  %s (%s:%d) - %s\n",
					t.Name, filepath.Base(t.File), t.Line, status))
			}
		}
		result.Message = sb.String()
	}

	return mcpJSONResult(result)
}

// mcpHandleGraph implements the wetwire_graph tool.
func mcpHandleGraph(_ context.Context, args map[string]any) (string, error) {
	pkg, _ := args["package"].(string)
	format, _ := args["format"].(string)

	result := MCPGraphResult{}

	if pkg == "" {
		result.Message = "package is required"
		return mcpJSONResult(result)
	}

	if format == "" {
		format = "text"
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(pkg)
	if err != nil {
		result.Message = fmt.Sprintf("discovery failed: %v", err)
		return mcpJSONResult(result)
	}

	if resources.TotalCount() == 0 {
		result.Message = "no resources found"
		return mcpJSONResult(result)
	}

	if format == "dot" {
		// DOT format for graphviz with full dependency chain
		var sb strings.Builder
		sb.WriteString("digraph resources {\n")
		sb.WriteString("  rankdir=LR;\n")
		sb.WriteString("  node [shape=box];\n\n")

		// Queries
		if len(resources.Queries) > 0 {
			sb.WriteString("  subgraph cluster_queries {\n")
			sb.WriteString("    label=\"Queries\";\n")
			sb.WriteString("    style=filled;\n")
			sb.WriteString("    fillcolor=\"#e3f2fd\";\n")
			for _, q := range resources.Queries {
				sb.WriteString(fmt.Sprintf("    query_%s [label=%q];\n", mcpSanitizeID(q.Name), q.Name))
			}
			sb.WriteString("  }\n\n")
		}

		// SLOs
		if len(resources.SLOs) > 0 {
			sb.WriteString("  subgraph cluster_slos {\n")
			sb.WriteString("    label=\"SLOs\";\n")
			sb.WriteString("    style=filled;\n")
			sb.WriteString("    fillcolor=\"#e8f5e9\";\n")
			for _, s := range resources.SLOs {
				sb.WriteString(fmt.Sprintf("    slo_%s [label=%q];\n", mcpSanitizeID(s.Name), s.SLOName))
			}
			sb.WriteString("  }\n\n")
		}

		// Triggers
		if len(resources.Triggers) > 0 {
			sb.WriteString("  subgraph cluster_triggers {\n")
			sb.WriteString("    label=\"Triggers\";\n")
			sb.WriteString("    style=filled;\n")
			sb.WriteString("    fillcolor=\"#fff3e0\";\n")
			for _, t := range resources.Triggers {
				sb.WriteString(fmt.Sprintf("    trigger_%s [label=%q];\n", mcpSanitizeID(t.Name), t.TriggerName))
			}
			sb.WriteString("  }\n\n")
		}

		// Boards
		if len(resources.Boards) > 0 {
			sb.WriteString("  subgraph cluster_boards {\n")
			sb.WriteString("    label=\"Boards\";\n")
			sb.WriteString("    style=filled;\n")
			sb.WriteString("    fillcolor=\"#fce4ec\";\n")
			for _, b := range resources.Boards {
				sb.WriteString(fmt.Sprintf("    board_%s [label=%q];\n", mcpSanitizeID(b.Name), b.BoardName))
			}
			sb.WriteString("  }\n\n")
		}

		// Edges: Board -> Query refs
		for _, b := range resources.Boards {
			for _, qref := range b.QueryRefs {
				sb.WriteString(fmt.Sprintf("  board_%s -> query_%s;\n", mcpSanitizeID(b.Name), mcpSanitizeID(qref)))
			}
			for _, sref := range b.SLORefs {
				sb.WriteString(fmt.Sprintf("  board_%s -> slo_%s;\n", mcpSanitizeID(b.Name), mcpSanitizeID(sref)))
			}
		}

		sb.WriteString("}\n")
		result.Graph = sb.String()
	} else {
		// Text format with dependency chain
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Resource Graph (%d total)\n", resources.TotalCount()))
		sb.WriteString(fmt.Sprintf("  Queries: %d, SLOs: %d, Triggers: %d, Boards: %d\n\n",
			len(resources.Queries), len(resources.SLOs), len(resources.Triggers), len(resources.Boards)))

		// Group queries by dataset
		byDataset := make(map[string][]discovery.DiscoveredQuery)
		for _, q := range resources.Queries {
			byDataset[q.Dataset] = append(byDataset[q.Dataset], q)
		}

		for dataset, dqs := range byDataset {
			sb.WriteString(fmt.Sprintf("Dataset: %s\n", dataset))
			for _, q := range dqs {
				sb.WriteString(fmt.Sprintf("  +-- Query: %s\n", q.Name))
			}
			sb.WriteString("\n")
		}

		// SLOs
		if len(resources.SLOs) > 0 {
			sb.WriteString("SLOs:\n")
			for _, s := range resources.SLOs {
				sb.WriteString(fmt.Sprintf("  +-- %s (%.1f%%, %d burn alerts)\n", s.SLOName, s.TargetPercentage, s.BurnAlertCount))
			}
			sb.WriteString("\n")
		}

		// Triggers
		if len(resources.Triggers) > 0 {
			sb.WriteString("Triggers:\n")
			for _, t := range resources.Triggers {
				status := "enabled"
				if t.Disabled {
					status = "disabled"
				}
				sb.WriteString(fmt.Sprintf("  +-- %s (%s)\n", t.TriggerName, status))
			}
			sb.WriteString("\n")
		}

		// Boards with dependencies
		if len(resources.Boards) > 0 {
			sb.WriteString("Boards:\n")
			for _, b := range resources.Boards {
				sb.WriteString(fmt.Sprintf("  +-- %s (%d panels)\n", b.BoardName, b.PanelCount))
				if len(b.QueryRefs) > 0 {
					sb.WriteString(fmt.Sprintf("      +-- queries: %v\n", b.QueryRefs))
				}
				if len(b.SLORefs) > 0 {
					sb.WriteString(fmt.Sprintf("      +-- slos: %v\n", b.SLORefs))
				}
			}
		}

		result.Graph = sb.String()
	}

	result.Success = true
	result.Format = format
	return mcpJSONResult(result)
}

// MCP Result types

// MCPInitResult is the result of the wetwire_init tool.
type MCPInitResult struct {
	Success bool     `json:"success"`
	Path    string   `json:"path"`
	Files   []string `json:"files"`
	Error   string   `json:"error,omitempty"`
}

// MCPLintResult is the result of the wetwire_lint tool.
type MCPLintResult struct {
	Success bool           `json:"success"`
	Issues  []MCPLintIssue `json:"issues,omitempty"`
}

// MCPLintIssue represents a single lint issue.
type MCPLintIssue struct {
	RuleID   string `json:"rule_id"`
	Message  string `json:"message"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Severity string `json:"severity"`
}

// MCPBuildResult is the result of the wetwire_build tool.
type MCPBuildResult struct {
	Success  bool     `json:"success"`
	Output   string   `json:"output,omitempty"`
	Queries  []string `json:"queries,omitempty"`
	Boards   []string `json:"boards,omitempty"`
	SLOs     []string `json:"slos,omitempty"`
	Triggers []string `json:"triggers,omitempty"`
	Errors   []string `json:"errors,omitempty"`
}

// MCPListResult is the result of the wetwire_list tool.
type MCPListResult struct {
	Success  bool              `json:"success"`
	Count    int               `json:"count"`
	Message  string            `json:"message,omitempty"`
	Queries  []MCPQueryInfo    `json:"queries,omitempty"`
	Boards   []MCPBoardInfo    `json:"boards,omitempty"`
	SLOs     []MCPSLOInfo      `json:"slos,omitempty"`
	Triggers []MCPTriggerInfo  `json:"triggers,omitempty"`
}

// MCPQueryInfo represents query metadata.
type MCPQueryInfo struct {
	Name    string `json:"name"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Dataset string `json:"dataset"`
}

// MCPBoardInfo represents board metadata.
type MCPBoardInfo struct {
	Name       string   `json:"name"`
	File       string   `json:"file"`
	Line       int      `json:"line"`
	BoardName  string   `json:"board_name,omitempty"`
	PanelCount int      `json:"panel_count"`
	QueryRefs  []string `json:"query_refs,omitempty"`
	SLORefs    []string `json:"slo_refs,omitempty"`
}

// MCPSLOInfo represents SLO metadata.
type MCPSLOInfo struct {
	Name             string  `json:"name"`
	File             string  `json:"file"`
	Line             int     `json:"line"`
	SLOName          string  `json:"slo_name,omitempty"`
	Dataset          string  `json:"dataset,omitempty"`
	TargetPercentage float64 `json:"target_percentage,omitempty"`
	BurnAlertCount   int     `json:"burn_alert_count"`
}

// MCPTriggerInfo represents trigger metadata.
type MCPTriggerInfo struct {
	Name             string `json:"name"`
	File             string `json:"file"`
	Line             int    `json:"line"`
	TriggerName      string `json:"trigger_name,omitempty"`
	Dataset          string `json:"dataset,omitempty"`
	FrequencySeconds int    `json:"frequency_seconds,omitempty"`
	RecipientCount   int    `json:"recipient_count"`
	Disabled         bool   `json:"disabled"`
}

// MCPGraphResult is the result of the wetwire_graph tool.
type MCPGraphResult struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Format  string `json:"format"`
	Graph   string `json:"graph,omitempty"`
}

// Helper functions

// mcpJSONResult creates a JSON string from any value.
func mcpJSONResult(v any) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling result: %w", err)
	}
	return string(data), nil
}

func mcpSanitizeID(s string) string {
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
