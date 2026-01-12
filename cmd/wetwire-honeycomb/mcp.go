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

	"github.com/lex00/wetwire-honeycomb-go/internal/builder"
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
	server.RegisterToolWithSchema("wetwire_build", "Generate Query JSON from wetwire declarations", mcpHandleBuild, map[string]any{
		"type": "object",
		"properties": map[string]any{
			"package": map[string]any{
				"type":        "string",
				"description": "Package path to discover queries from",
			},
			"format": map[string]any{
				"type":        "string",
				"enum":        []string{"json", "pretty"},
				"description": "Output format (default: pretty)",
			},
			"dry_run": map[string]any{
				"type":        "boolean",
				"description": "Return content without writing files",
			},
		},
	})

	// wetwire_lint tool
	server.RegisterToolWithSchema("wetwire_lint", "Check code quality and style (domain lint rules)", mcpHandleLint, map[string]any{
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
	server.RegisterToolWithSchema("wetwire_list", "List all discovered queries", mcpHandleList, map[string]any{
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

	// Discover queries
	queries, err := discovery.DiscoverQueries(absPath)
	if err != nil {
		result.Issues = append(result.Issues, MCPLintIssue{
			Severity: "error",
			Message:  fmt.Sprintf("discovery failed: %v", err),
			RuleID:   "internal",
		})
		return mcpJSONResult(result)
	}

	// Run lint
	lintResults := lint.LintQueries(queries)

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

	// Build queries
	b, err := builder.NewBuilder(absPath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("builder error: %v", err))
		return mcpJSONResult(result)
	}

	buildResult, err := b.Build()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("build failed: %v", err))
		return mcpJSONResult(result)
	}

	if buildResult.QueryCount() == 0 {
		result.Errors = append(result.Errors, "no queries found")
		return mcpJSONResult(result)
	}

	// Serialize queries
	queries := buildResult.Queries()
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
		result.Errors = append(result.Errors, fmt.Sprintf("serialization failed: %v", err))
		return mcpJSONResult(result)
	}

	// Build success result
	for _, q := range queries {
		result.Queries = append(result.Queries, q.Name)
	}

	result.Success = true
	result.Output = string(jsonData)
	return mcpJSONResult(result)
}

// mcpHandleList implements the wetwire_list tool.
func mcpHandleList(_ context.Context, args map[string]any) (string, error) {
	pkg, _ := args["package"].(string)
	format, _ := args["format"].(string)

	result := MCPListResult{}

	if pkg == "" {
		result.Message = "package is required"
		return mcpJSONResult(result)
	}

	// Discover queries
	queries, err := discovery.DiscoverQueries(pkg)
	if err != nil {
		result.Message = fmt.Sprintf("discovery failed: %v", err)
		return mcpJSONResult(result)
	}

	// Build query list
	for _, q := range queries {
		result.Queries = append(result.Queries, MCPQueryInfo{
			Name:    q.Name,
			File:    q.File,
			Line:    q.Line,
			Dataset: q.Dataset,
		})
	}

	result.Success = true
	result.Count = len(queries)

	if format == "text" {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Found %d queries:\n", len(queries)))
		for _, q := range queries {
			sb.WriteString(fmt.Sprintf("  %s (%s:%d) - dataset: %s\n",
				q.Name, filepath.Base(q.File), q.Line, q.Dataset))
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

	// Discover queries
	queries, err := discovery.DiscoverQueries(pkg)
	if err != nil {
		result.Message = fmt.Sprintf("discovery failed: %v", err)
		return mcpJSONResult(result)
	}

	if len(queries) == 0 {
		result.Message = "no queries found"
		return mcpJSONResult(result)
	}

	// Group by dataset
	byDataset := make(map[string][]discovery.DiscoveredQuery)
	for _, q := range queries {
		byDataset[q.Dataset] = append(byDataset[q.Dataset], q)
	}

	if format == "dot" {
		// DOT format for graphviz
		var sb strings.Builder
		sb.WriteString("digraph queries {\n")
		sb.WriteString("  rankdir=LR;\n")
		sb.WriteString("  node [shape=box];\n\n")

		for dataset, dqs := range byDataset {
			sb.WriteString(fmt.Sprintf("  subgraph cluster_%s {\n", mcpSanitizeID(dataset)))
			sb.WriteString(fmt.Sprintf("    label=%q;\n", dataset))
			for _, q := range dqs {
				sb.WriteString(fmt.Sprintf("    %s [label=%q];\n", mcpSanitizeID(q.Name), q.Name))
			}
			sb.WriteString("  }\n")
		}

		sb.WriteString("}\n")
		result.Graph = sb.String()
	} else {
		// Text format
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Query Graph (%d queries, %d datasets)\n\n", len(queries), len(byDataset)))

		for dataset, dqs := range byDataset {
			sb.WriteString(fmt.Sprintf("Dataset: %s\n", dataset))
			for _, q := range dqs {
				sb.WriteString(fmt.Sprintf("  +-- %s (%s:%d)\n", q.Name, filepath.Base(q.File), q.Line))
				if len(q.Breakdowns) > 0 {
					sb.WriteString(fmt.Sprintf("  |   +-- breakdowns: %v\n", q.Breakdowns))
				}
			}
			sb.WriteString("\n")
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
	Success bool     `json:"success"`
	Output  string   `json:"output,omitempty"`
	Queries []string `json:"queries,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// MCPListResult is the result of the wetwire_list tool.
type MCPListResult struct {
	Success bool           `json:"success"`
	Count   int            `json:"count"`
	Message string         `json:"message,omitempty"`
	Queries []MCPQueryInfo `json:"queries,omitempty"`
}

// MCPQueryInfo represents query metadata.
type MCPQueryInfo struct {
	Name    string `json:"name"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Dataset string `json:"dataset"`
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
