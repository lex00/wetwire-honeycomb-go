// Package builder provides the query build pipeline.
package builder

import (
	"fmt"

	"github.com/lex00/wetwire-honeycomb-go/internal/discover"
)

// Builder orchestrates the query discovery and build pipeline.
type Builder struct {
	path string

	// Options
	namespacing   bool
	strictMode    bool
	packageFilter string
	datasetFilter string
}

// NewBuilder creates a new Builder for the given path.
func NewBuilder(path string) (*Builder, error) {
	return &Builder{
		path: path,
	}, nil
}

// WithNamespacing enables package-based namespacing for query names.
func (b *Builder) WithNamespacing(enabled bool) *Builder {
	b.namespacing = enabled
	return b
}

// WithStrictMode enables strict mode which fails on duplicates.
func (b *Builder) WithStrictMode(enabled bool) *Builder {
	b.strictMode = enabled
	return b
}

// WithPackageFilter filters queries to only include those from the specified package.
func (b *Builder) WithPackageFilter(pkg string) *Builder {
	b.packageFilter = pkg
	return b
}

// WithDatasetFilter filters queries to only include those targeting the specified dataset.
func (b *Builder) WithDatasetFilter(dataset string) *Builder {
	b.datasetFilter = dataset
	return b
}

// Build discovers queries and builds the registry.
func (b *Builder) Build() (*BuildResult, error) {
	// Discover queries
	queries, err := discovery.DiscoverQueries(b.path)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Create registry
	registry := NewRegistry()
	if b.namespacing {
		registry.EnableNamespacing()
	}
	if b.strictMode {
		registry.EnableStrictMode()
	}

	// Apply filters and register queries
	for _, q := range queries {
		// Apply package filter
		if b.packageFilter != "" && q.Package != b.packageFilter {
			continue
		}

		// Apply dataset filter
		if b.datasetFilter != "" && q.Dataset != b.datasetFilter {
			continue
		}

		if err := registry.Register(q); err != nil {
			if b.strictMode {
				return nil, err
			}
			// In non-strict mode, continue on duplicates
		}
	}

	return &BuildResult{
		registry: registry,
	}, nil
}

// BuildResult contains the result of a build operation.
type BuildResult struct {
	registry *Registry
}

// QueryCount returns the number of queries in the result.
func (r *BuildResult) QueryCount() int {
	return r.registry.Count()
}

// Queries returns all queries in the result.
func (r *BuildResult) Queries() []discovery.DiscoveredQuery {
	return r.registry.All()
}

// Query returns a query by name, or nil if not found.
func (r *BuildResult) Query(name string) *discovery.DiscoveredQuery {
	return r.registry.Get(name)
}

// Duplicates returns any duplicate queries detected.
func (r *BuildResult) Duplicates() []DuplicateEntry {
	return r.registry.Duplicates()
}
