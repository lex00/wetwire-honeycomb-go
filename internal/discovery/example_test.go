package discovery_test

import (
	"fmt"
	"path/filepath"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

func ExampleDiscoverQueries() {
	// Discover all queries in a directory
	queries, err := discovery.DiscoverQueries("../../testdata/queries")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print discovered queries
	for _, q := range queries {
		fmt.Printf("Query: %s\n", q.Name)
		fmt.Printf("  Dataset: %s\n", q.Dataset)
		fmt.Printf("  Package: %s\n", q.Package)
		fmt.Printf("  File: %s\n", filepath.Base(q.File))
		fmt.Printf("  Line: %d\n", q.Line)
		fmt.Printf("  Breakdowns: %d\n", len(q.Breakdowns))
		fmt.Printf("  Calculations: %d\n", len(q.Calculations))
		fmt.Printf("  Filters: %d\n", len(q.Filters))
		fmt.Println()
	}
}

func ExampleGroupByDataset() {
	// Discover queries and group by dataset
	queries, err := discovery.DiscoverQueries("../../testdata/queries")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	grouped := discovery.GroupByDataset(queries)

	for dataset, queries := range grouped {
		fmt.Printf("Dataset: %s (%d queries)\n", dataset, len(queries))
		for _, q := range queries {
			fmt.Printf("  - %s\n", q.Name)
		}
	}
}

func ExampleFilterByDataset() {
	// Discover queries and filter by dataset
	queries, err := discovery.DiscoverQueries("../../testdata/queries")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Get only production queries
	prodQueries := discovery.FilterByDataset(queries, "production")

	fmt.Printf("Production queries: %d\n", len(prodQueries))
	for _, q := range prodQueries {
		fmt.Printf("  - %s\n", q.Name)
	}
}
