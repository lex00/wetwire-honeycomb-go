package builder

import (
	"fmt"

	"github.com/lex00/wetwire-honeycomb-go/internal/discover"
)

// Registry stores discovered queries and tracks duplicates.
type Registry struct {
	queries     map[string]discovery.DiscoveredQuery
	duplicates  []DuplicateEntry
	namespacing bool
	strictMode  bool
}

// DuplicateEntry represents a duplicate query entry.
type DuplicateEntry struct {
	Name     string
	Original discovery.DiscoveredQuery
	Duplicate discovery.DiscoveredQuery
}

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		queries: make(map[string]discovery.DiscoveredQuery),
	}
}

// EnableNamespacing enables package-based namespacing.
func (r *Registry) EnableNamespacing() {
	r.namespacing = true
}

// EnableStrictMode enables strict mode.
func (r *Registry) EnableStrictMode() {
	r.strictMode = true
}

// Register adds a query to the registry.
func (r *Registry) Register(q discovery.DiscoveredQuery) error {
	name := r.queryName(q)

	if existing, exists := r.queries[name]; exists {
		r.duplicates = append(r.duplicates, DuplicateEntry{
			Name:      name,
			Original:  existing,
			Duplicate: q,
		})
		if r.strictMode {
			return fmt.Errorf("duplicate query %q: first defined in %s:%d, also defined in %s:%d",
				name, existing.File, existing.Line, q.File, q.Line)
		}
		return nil
	}

	r.queries[name] = q
	return nil
}

// queryName returns the name for a query, with optional namespacing.
func (r *Registry) queryName(q discovery.DiscoveredQuery) string {
	if r.namespacing {
		return q.Package + "." + q.Name
	}
	return q.Name
}

// Get returns a query by name, or nil if not found.
func (r *Registry) Get(name string) *discovery.DiscoveredQuery {
	if q, exists := r.queries[name]; exists {
		return &q
	}
	return nil
}

// All returns all queries in the registry.
func (r *Registry) All() []discovery.DiscoveredQuery {
	result := make([]discovery.DiscoveredQuery, 0, len(r.queries))
	for _, q := range r.queries {
		result = append(result, q)
	}
	return result
}

// Count returns the number of queries in the registry.
func (r *Registry) Count() int {
	return len(r.queries)
}

// Duplicates returns any detected duplicate entries.
func (r *Registry) Duplicates() []DuplicateEntry {
	return r.duplicates
}

// HasDuplicates returns true if any duplicates were detected.
func (r *Registry) HasDuplicates() bool {
	return len(r.duplicates) > 0
}
