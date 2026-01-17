// MCP server implementation using domain.BuildMCPServer().
package main

import (
	"context"

	"github.com/lex00/wetwire-honeycomb-go/domain"
	"github.com/spf13/cobra"

	coredomain "github.com/lex00/wetwire-core-go/domain"
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

// runMCPServer starts the MCP server on stdio transport using domain.BuildMCPServer().
func runMCPServer() error {
	server := coredomain.BuildMCPServer(&domain.HoneycombDomain{})
	return server.Start(context.Background())
}
