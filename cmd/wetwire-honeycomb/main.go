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
	"fmt"
	"os"

	"github.com/lex00/wetwire-honeycomb-go/domain"
	"github.com/lex00/wetwire-honeycomb-go/internal/discover"
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/spf13/cobra"
)

// Version information set via ldflags
var version = "dev"

func main() {
	// Set domain version from ldflags
	domain.Version = version

	// Use domain interface for auto-generated commands
	rootCmd := domain.CreateRootCommand(&domain.HoneycombDomain{})

	// Add domain-specific commands
	addDomainSpecificCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// addDomainSpecificCommands adds Honeycomb-specific commands to the root command.
func addDomainSpecificCommands(rootCmd *cobra.Command) {
	// Add custom commands not covered by domain interface
	rootCmd.AddCommand(
		newDiffCmd(),
		newWatchCmd(),
		newDesignCmd(),
		newTestCmd(),
		newMCPCmd(),
	)
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
