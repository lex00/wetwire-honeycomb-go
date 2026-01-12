package discovery

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredQuery represents a discovered query definition with metadata.
type DiscoveredQuery struct {
	// Name is the identifier of the query (variable name or function name)
	Name string

	// Package is the package name where the query is defined
	Package string

	// File is the absolute path to the file containing the query
	File string

	// Line is the line number where the query is defined
	Line int

	// Dataset is the Honeycomb dataset being queried
	Dataset string

	// TimeRange specifies the time window for the query
	TimeRange TimeRange

	// Breakdowns are the fields to group by
	Breakdowns []string

	// Calculations are the aggregations to compute
	Calculations []Calculation

	// Filters restrict the data being queried
	Filters []Filter

	// FilterCombination specifies how filters are combined ("AND" or "OR")
	FilterCombination string

	// Orders specify how results should be sorted
	Orders []Order

	// Granularity is the time bucket size in seconds
	Granularity int

	// Limit restricts the number of results
	Limit int

	// Style contains metadata for style linting
	Style StyleMetadata
}

// TimeRange represents a time window for a query.
// This matches the structure in the query package.
type TimeRange struct {
	// TimeRange is relative time in seconds
	TimeRange int

	// StartTime is absolute start time in Unix epoch seconds
	StartTime int

	// EndTime is absolute end time in Unix epoch seconds
	EndTime int
}

// Calculation represents an aggregation to compute.
type Calculation struct {
	// Op is the calculation operation (e.g., "COUNT", "P99", "AVG")
	Op string

	// Column is the field to operate on (optional for COUNT)
	Column string

	// Alias is an optional name for the calculation result
	Alias string
}

// Filter represents a condition to filter data.
type Filter struct {
	// Column is the field to filter on
	Column string

	// Op is the filter operator (e.g., "=", ">", "<")
	Op string

	// Value is the value to compare against
	Value interface{}
}

// Order represents a sort specification for query results.
type Order struct {
	// Column is the field to sort by (for breakdown columns)
	Column string

	// Op is the calculation operation to sort by (for calculation results)
	Op string

	// Order is the sort direction ("ascending" or "descending")
	Order string
}

// StyleMetadata contains metadata for style linting.
type StyleMetadata struct {
	// InlineCalculationCount is the number of calculations defined inline
	// (as composite literals rather than named variables)
	InlineCalculationCount int

	// InlineFilterCount is the number of filters defined inline
	// (as composite literals rather than named variables)
	InlineFilterCount int

	// HasRawMapLiteral indicates if raw map literals are used instead of typed builders
	HasRawMapLiteral bool

	// MaxNestingDepth is the maximum nesting depth of the query configuration
	MaxNestingDepth int
}

// DiscoverQueries discovers all Query definitions in the specified directory.
// It parses Go source files and extracts query metadata using AST analysis.
func DiscoverQueries(dir string) ([]DiscoveredQuery, error) {
	// Check if directory exists
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", dir)
	}

	var discovered []DiscoveredQuery

	// Walk the directory tree
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and test files
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file
		queries, err := discoverQueriesInFile(path)
		if err != nil {
			// Log but don't fail on individual file errors
			return nil
		}

		discovered = append(discovered, queries...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return discovered, nil
}

// discoverQueriesInFile discovers queries in a single Go source file.
func discoverQueriesInFile(path string) ([]DiscoveredQuery, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var discovered []DiscoveredQuery

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	packageName := node.Name.Name

	// Walk the AST
	ast.Inspect(node, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.GenDecl:
			// Handle package-level variable and constant declarations
			if decl.Tok == token.VAR || decl.Tok == token.CONST {
				for _, spec := range decl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						queries := extractQueriesFromValueSpec(valueSpec, fset, absPath, packageName)
						discovered = append(discovered, queries...)
					}
				}
			}

		case *ast.FuncDecl:
			// Handle function-scoped queries
			if decl.Body != nil {
				queries := extractQueriesFromFunction(decl, fset, absPath, packageName)
				discovered = append(discovered, queries...)
			}
		}

		return true
	})

	return discovered, nil
}

// extractQueriesFromValueSpec extracts queries from a variable or constant declaration.
func extractQueriesFromValueSpec(spec *ast.ValueSpec, fset *token.FileSet, file string, pkg string) []DiscoveredQuery {
	var discovered []DiscoveredQuery

	name := getIdentifierName(spec)
	if name == "" || !isExportedName(name) {
		return discovered
	}

	// Check each value in the spec
	for _, value := range spec.Values {
		// Find all query composites in this value
		composites := findQueryComposites(value)

		for _, comp := range composites {
			query := extractQueryFromComposite(comp, fset, file, pkg, name)
			if query.Name != "" {
				discovered = append(discovered, query)
			}
		}
	}

	return discovered
}

// extractQueriesFromFunction extracts queries from return statements in a function.
func extractQueriesFromFunction(decl *ast.FuncDecl, fset *token.FileSet, file string, pkg string) []DiscoveredQuery {
	var discovered []DiscoveredQuery

	funcName := getFunctionName(decl)
	if funcName == "" || !isExportedName(funcName) {
		return discovered
	}

	// Walk function body looking for return statements
	ast.Inspect(decl.Body, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			for _, result := range ret.Results {
				composites := findQueryComposites(result)
				for _, comp := range composites {
					query := extractQueryFromComposite(comp, fset, file, pkg, funcName)
					if query.Name != "" {
						discovered = append(discovered, query)
					}
				}
			}
		}
		return true
	})

	return discovered
}

// extractQueryFromComposite extracts query metadata from a composite literal.
func extractQueryFromComposite(comp *ast.CompositeLit, fset *token.FileSet, file string, pkg string, name string) DiscoveredQuery {
	query := DiscoveredQuery{
		Name:    name,
		Package: pkg,
		File:    file,
		Line:    fset.Position(comp.Pos()).Line,
	}

	// Extract fields from the composite literal
	for _, elt := range comp.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch key.Name {
		case "Dataset":
			query.Dataset = extractStringLiteral(kv.Value)

		case "TimeRange":
			query.TimeRange = extractTimeRange(kv.Value)

		case "Breakdowns":
			query.Breakdowns = extractStringSlice(kv.Value)

		case "Calculations":
			query.Calculations = extractCalculations(kv.Value)

		case "Filters":
			query.Filters = extractFilters(kv.Value)

		case "FilterCombination":
			query.FilterCombination = extractStringLiteral(kv.Value)

		case "Orders":
			query.Orders = extractOrders(kv.Value)

		case "Granularity":
			query.Granularity = extractIntLiteral(kv.Value)

		case "Limit":
			query.Limit = extractIntLiteral(kv.Value)
		}
	}

	// Extract style metadata for linting
	query.Style = extractStyleMetadata(comp)

	return query
}

// DiscoverQueriesInPackage discovers queries in a specific Go package.
// This is a convenience wrapper around DiscoverQueries for package paths.
func DiscoverQueriesInPackage(pkgPath string) ([]DiscoveredQuery, error) {
	return DiscoverQueries(pkgPath)
}

// FilterByDataset filters discovered queries by dataset name.
func FilterByDataset(queries []DiscoveredQuery, dataset string) []DiscoveredQuery {
	var result []DiscoveredQuery
	for _, q := range queries {
		if q.Dataset == dataset {
			result = append(result, q)
		}
	}
	return result
}

// FilterByPackage filters discovered queries by package name.
func FilterByPackage(queries []DiscoveredQuery, pkg string) []DiscoveredQuery {
	var result []DiscoveredQuery
	for _, q := range queries {
		if q.Package == pkg {
			result = append(result, q)
		}
	}
	return result
}

// GroupByDataset groups discovered queries by their dataset.
func GroupByDataset(queries []DiscoveredQuery) map[string][]DiscoveredQuery {
	result := make(map[string][]DiscoveredQuery)
	for _, q := range queries {
		result[q.Dataset] = append(result[q.Dataset], q)
	}
	return result
}

// GroupByPackage groups discovered queries by their package.
func GroupByPackage(queries []DiscoveredQuery) map[string][]DiscoveredQuery {
	result := make(map[string][]DiscoveredQuery)
	for _, q := range queries {
		result[q.Package] = append(result[q.Package], q)
	}
	return result
}
