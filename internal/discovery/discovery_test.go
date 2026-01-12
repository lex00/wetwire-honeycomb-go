package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverQueries_SimplePackageLevel(t *testing.T) {
	// Test discovering simple package-level query variables
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should find SlowRequests and ErrorRate from simple.go
	if len(discovered) < 2 {
		t.Errorf("Expected at least 2 queries, got %d", len(discovered))
	}

	// Check SlowRequests query
	slowRequests := findQuery(discovered, "SlowRequests")
	if slowRequests == nil {
		t.Fatal("SlowRequests query not found")
	}

	if slowRequests.Dataset != "production" {
		t.Errorf("Expected dataset 'production', got %q", slowRequests.Dataset)
	}

	if slowRequests.Package != "queries" {
		t.Errorf("Expected package 'queries', got %q", slowRequests.Package)
	}

	if slowRequests.File == "" {
		t.Error("Expected non-empty file path")
	}

	if slowRequests.Line <= 0 {
		t.Error("Expected positive line number")
	}

	if len(slowRequests.Breakdowns) != 2 {
		t.Errorf("Expected 2 breakdowns, got %d", len(slowRequests.Breakdowns))
	}

	if len(slowRequests.Calculations) != 2 {
		t.Errorf("Expected 2 calculations, got %d", len(slowRequests.Calculations))
	}
}

func TestDiscoverQueries_EmbeddedTypes(t *testing.T) {
	// Test discovering queries embedded in custom types
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should find DatabaseQueries from embedded.go
	dbQuery := findQuery(discovered, "DatabaseQueries")
	if dbQuery == nil {
		t.Fatal("DatabaseQueries query not found")
	}

	if dbQuery.Dataset != "database" {
		t.Errorf("Expected dataset 'database', got %q", dbQuery.Dataset)
	}

	if len(dbQuery.Breakdowns) != 1 {
		t.Errorf("Expected 1 breakdown, got %d", len(dbQuery.Breakdowns))
	}
}

func TestDiscoverQueries_FunctionScoped(t *testing.T) {
	// Test discovering queries defined within functions
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should find the query in GetLatencyQuery function
	latencyQuery := findQuery(discovered, "GetLatencyQuery")
	if latencyQuery == nil {
		t.Fatal("GetLatencyQuery query not found")
	}

	if latencyQuery.Dataset != "api" {
		t.Errorf("Expected dataset 'api', got %q", latencyQuery.Dataset)
	}

	if len(latencyQuery.Calculations) != 1 {
		t.Errorf("Expected 1 calculation, got %d", len(latencyQuery.Calculations))
	}
}

func TestDiscoverQueries_EmptyDirectory(t *testing.T) {
	// Test discovering queries in an empty directory
	tempDir := t.TempDir()

	discovered, err := DiscoverQueries(tempDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	if len(discovered) != 0 {
		t.Errorf("Expected 0 queries in empty directory, got %d", len(discovered))
	}
}

func TestDiscoverQueries_NonExistentDirectory(t *testing.T) {
	// Test handling of non-existent directory
	_, err := DiscoverQueries("/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestDiscoverQueries_TimeRangeExtraction(t *testing.T) {
	// Test that time range information is correctly extracted
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	slowRequests := findQuery(discovered, "SlowRequests")
	if slowRequests == nil {
		t.Fatal("SlowRequests query not found")
	}

	// Verify time range was extracted (2 hours = 7200 seconds)
	expectedTimeRange := 2 * 3600
	if slowRequests.TimeRange.TimeRange != expectedTimeRange {
		t.Errorf("Expected time range %d, got %d", expectedTimeRange, slowRequests.TimeRange.TimeRange)
	}
}

func TestDiscoverQueries_MultipleFiles(t *testing.T) {
	// Test discovering queries across multiple files
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should find queries from both simple.go and embedded.go
	fileSet := make(map[string]bool)
	for _, q := range discovered {
		fileSet[filepath.Base(q.File)] = true
	}

	if !fileSet["simple.go"] {
		t.Error("Expected to find queries from simple.go")
	}

	if !fileSet["embedded.go"] {
		t.Error("Expected to find queries from embedded.go")
	}
}

func TestDiscoverQueries_AdvancedFeatures(t *testing.T) {
	// Test discovering queries with advanced features (multiple calculations, filters, etc.)
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should find AdvancedQuery from advanced.go
	advQuery := findQuery(discovered, "AdvancedQuery")
	if advQuery == nil {
		t.Fatal("AdvancedQuery query not found")
	}

	if advQuery.Dataset != "logs" {
		t.Errorf("Expected dataset 'logs', got %q", advQuery.Dataset)
	}

	if len(advQuery.Breakdowns) != 3 {
		t.Errorf("Expected 3 breakdowns, got %d", len(advQuery.Breakdowns))
	}

	if len(advQuery.Calculations) != 4 {
		t.Errorf("Expected 4 calculations, got %d", len(advQuery.Calculations))
	}

	if len(advQuery.Filters) != 3 {
		t.Errorf("Expected 3 filters, got %d", len(advQuery.Filters))
	}

	if advQuery.Limit != 100 {
		t.Errorf("Expected limit 100, got %d", advQuery.Limit)
	}

	// Verify MetricsQuery
	metricsQuery := findQuery(discovered, "MetricsQuery")
	if metricsQuery == nil {
		t.Fatal("MetricsQuery query not found")
	}

	if metricsQuery.Dataset != "metrics" {
		t.Errorf("Expected dataset 'metrics', got %q", metricsQuery.Dataset)
	}

	// Verify ErrorTracking
	errorTracking := findQuery(discovered, "ErrorTracking")
	if errorTracking == nil {
		t.Fatal("ErrorTracking query not found")
	}

	if errorTracking.Dataset != "errors" {
		t.Errorf("Expected dataset 'errors', got %q", errorTracking.Dataset)
	}
}

func TestGroupByDataset(t *testing.T) {
	// Test grouping queries by dataset
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	grouped := GroupByDataset(discovered)

	// Should have multiple datasets
	if len(grouped) < 3 {
		t.Errorf("Expected at least 3 datasets, got %d", len(grouped))
	}

	// Verify production dataset has queries
	if prodQueries, ok := grouped["production"]; ok {
		if len(prodQueries) < 2 {
			t.Errorf("Expected at least 2 queries in production dataset, got %d", len(prodQueries))
		}
	} else {
		t.Error("Expected production dataset to be present")
	}
}

func TestFilterByDataset(t *testing.T) {
	// Test filtering queries by dataset
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Filter by production dataset
	prodQueries := FilterByDataset(discovered, "production")
	if len(prodQueries) < 2 {
		t.Errorf("Expected at least 2 queries in production dataset, got %d", len(prodQueries))
	}

	// Verify all returned queries are from production
	for _, q := range prodQueries {
		if q.Dataset != "production" {
			t.Errorf("Expected query from production dataset, got %q", q.Dataset)
		}
	}
}

func TestDiscoverQueriesInPackage(t *testing.T) {
	// Test DiscoverQueriesInPackage wrapper function
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueriesInPackage(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueriesInPackage failed: %v", err)
	}

	// Should find the same queries as DiscoverQueries
	if len(discovered) < 2 {
		t.Errorf("Expected at least 2 queries, got %d", len(discovered))
	}
}

func TestFilterByPackage(t *testing.T) {
	// Test filtering queries by package name
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Filter by queries package
	queriesPackage := FilterByPackage(discovered, "queries")
	if len(queriesPackage) == 0 {
		t.Error("Expected at least 1 query in 'queries' package")
	}

	// Verify all returned queries are from queries package
	for _, q := range queriesPackage {
		if q.Package != "queries" {
			t.Errorf("Expected query from queries package, got %q", q.Package)
		}
	}

	// Filter by non-existent package
	nonExistent := FilterByPackage(discovered, "nonexistent")
	if len(nonExistent) != 0 {
		t.Errorf("Expected 0 queries in non-existent package, got %d", len(nonExistent))
	}
}

func TestGroupByPackage(t *testing.T) {
	// Test grouping queries by package
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")

	discovered, err := DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	grouped := GroupByPackage(discovered)

	// Should have at least queries package
	if len(grouped) == 0 {
		t.Error("Expected at least 1 package")
	}

	// Verify queries package exists
	if queriesPackage, ok := grouped["queries"]; ok {
		if len(queriesPackage) == 0 {
			t.Error("Expected at least 1 query in queries package")
		}
	} else {
		t.Error("Expected queries package to be present")
	}

	// Verify all queries in each group have matching package
	for pkg, queries := range grouped {
		for _, q := range queries {
			if q.Package != pkg {
				t.Errorf("Query %s has package %q but is in group %q", q.Name, q.Package, pkg)
			}
		}
	}
}

func TestDiscoverQueries_InvalidFile(t *testing.T) {
	// Test discovering queries in a directory with a file (not directory)
	tempDir := t.TempDir()

	// Create a regular file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := DiscoverQueries(testFile)
	if err == nil {
		t.Error("Expected error for path that is not a directory, got nil")
	}
}

func TestDiscoverQueries_MalformedGoFile(t *testing.T) {
	// Test discovering queries in a directory with malformed Go file
	tempDir := t.TempDir()

	// Create a malformed Go file
	malformedFile := filepath.Join(tempDir, "malformed.go")
	if err := os.WriteFile(malformedFile, []byte("package test\n\nthis is not valid go syntax {{{"), 0644); err != nil {
		t.Fatalf("Failed to create malformed file: %v", err)
	}

	// Should not fail on malformed files, just skip them
	discovered, err := DiscoverQueries(tempDir)
	if err != nil {
		t.Fatalf("DiscoverQueries should not fail on malformed files: %v", err)
	}

	// Should return empty results since file is malformed
	if len(discovered) != 0 {
		t.Errorf("Expected 0 queries from malformed file, got %d", len(discovered))
	}
}

func TestExtractQueriesFromValueSpec_Unexported(t *testing.T) {
	// Test that unexported queries are skipped
	tempDir := t.TempDir()

	// Create a file with unexported query
	testFile := filepath.Join(tempDir, "unexported.go")
	fileContent := `package test

import "github.com/lex00/wetwire-honeycomb-go/query"

var unexportedQuery = query.Query{
	Dataset: "test",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{query.Count()},
}
`
	if err := os.WriteFile(testFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	discovered, err := DiscoverQueries(tempDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should not find unexported query
	for _, q := range discovered {
		if q.Name == "unexportedQuery" {
			t.Error("Should not discover unexported queries")
		}
	}
}

func TestExtractQueriesFromFunction_Unexported(t *testing.T) {
	// Test that unexported functions are skipped
	tempDir := t.TempDir()

	// Create a file with unexported function
	testFile := filepath.Join(tempDir, "unexported_func.go")
	fileContent := `package test

import "github.com/lex00/wetwire-honeycomb-go/query"

func unexportedFunc() query.Query {
	return query.Query{
		Dataset: "test",
		TimeRange: query.Hours(1),
		Calculations: []query.Calculation{query.Count()},
	}
}
`
	if err := os.WriteFile(testFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	discovered, err := DiscoverQueries(tempDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should not find unexported function
	for _, q := range discovered {
		if q.Name == "unexportedFunc" {
			t.Error("Should not discover queries from unexported functions")
		}
	}
}

func TestExtractQueryFromComposite_MissingFields(t *testing.T) {
	// Test extracting query from composite with missing fields
	tempDir := t.TempDir()

	// Create a file with minimal query (only dataset)
	testFile := filepath.Join(tempDir, "minimal.go")
	fileContent := `package test

import "github.com/lex00/wetwire-honeycomb-go/query"

var MinimalQuery = query.Query{
	Dataset: "test",
}
`
	if err := os.WriteFile(testFile, []byte(fileContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	discovered, err := DiscoverQueries(tempDir)
	if err != nil {
		t.Fatalf("DiscoverQueries failed: %v", err)
	}

	// Should find the query even with missing fields
	found := false
	for _, q := range discovered {
		if q.Name == "MinimalQuery" {
			found = true
			if q.Dataset != "test" {
				t.Errorf("Expected dataset 'test', got %q", q.Dataset)
			}
			// Other fields should be empty/zero
			if len(q.Calculations) != 0 {
				t.Error("Expected no calculations")
			}
		}
	}

	if !found {
		t.Error("Expected to find MinimalQuery")
	}
}

// Helper functions

func findQuery(queries []DiscoveredQuery, name string) *DiscoveredQuery {
	for i, q := range queries {
		if q.Name == name {
			return &queries[i]
		}
	}
	return nil
}

func getRepoRoot(t *testing.T) string {
	// Walk up from current directory to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Could not find repository root (go.mod not found)")
		}
		dir = parent
	}
}
