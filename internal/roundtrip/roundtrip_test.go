package roundtrip

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
	"github.com/lex00/wetwire-honeycomb-go/internal/serialize"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// TestRoundTrip_BasicQuery tests round-trip conversion for a basic query.
// This verifies that: JSON -> Go code -> Parse -> Build -> JSON produces semantically equivalent output.
func TestRoundTrip_BasicQuery(t *testing.T) {
	testRoundTrip(t, "basic_query.json")
}

// TestRoundTrip_ComplexQuery tests round-trip conversion for a complex query with all features.
// This verifies that: JSON -> Go code -> Parse -> Build -> JSON produces semantically equivalent output.
func TestRoundTrip_ComplexQuery(t *testing.T) {
	testRoundTrip(t, "complex_query.json")
}

// testRoundTrip performs the complete round-trip test for a given fixture file.
func testRoundTrip(t *testing.T, fixtureFile string) {
	t.Helper()

	// 1. Read original JSON fixture
	fixturePath := filepath.Join("../../testdata/roundtrip", fixtureFile)
	originalJSON, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("Failed to read fixture %s: %v", fixtureFile, err)
	}

	// Parse original JSON to get the expected structure
	var originalData map[string]interface{}
	if err := json.Unmarshal(originalJSON, &originalData); err != nil {
		t.Fatalf("Failed to parse original JSON: %v", err)
	}

	// 2. Import JSON to Go code (simulating the import command)
	goCode := generateGoCode("testpkg", "TestQuery", originalData)

	// 3. Write Go code to temporary file
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "query.go")
	if err := os.WriteFile(goFile, []byte(goCode), 0644); err != nil {
		t.Fatalf("Failed to write temporary Go file: %v", err)
	}

	// 4. Discover the query from the generated Go file
	queries, err := discovery.DiscoverQueries(tmpDir)
	if err != nil {
		t.Fatalf("Failed to discover queries: %v", err)
	}

	if len(queries) != 1 {
		t.Fatalf("Expected 1 query, found %d", len(queries))
	}

	discoveredQuery := queries[0]

	// 5. Convert discovered query to query.Query
	q := discoveredToQuery(discoveredQuery)

	// 6. Build back to JSON
	generatedJSON, err := serialize.ToJSON(q)
	if err != nil {
		t.Fatalf("Failed to serialize query to JSON: %v", err)
	}

	// Parse generated JSON
	var generatedData map[string]interface{}
	if err := json.Unmarshal(generatedJSON, &generatedData); err != nil {
		t.Fatalf("Failed to parse generated JSON: %v", err)
	}

	// 7. Compare original and generated JSON semantically
	if err := compareJSON(originalData, generatedData); err != nil {
		t.Errorf("Round-trip comparison failed for %s:\n%v", fixtureFile, err)
		t.Logf("Original JSON:\n%s", string(originalJSON))
		t.Logf("Generated JSON:\n%s", string(generatedJSON))
	}
}

// compareJSON performs semantic comparison of two JSON objects.
// It ignores whitespace, key ordering, and float vs int differences.
func compareJSON(expected, actual map[string]interface{}) error {
	// Normalize both maps
	normalizeMap(expected)
	normalizeMap(actual)

	// Compare each key in expected
	for key, expectedVal := range expected {
		actualVal, exists := actual[key]
		if !exists {
			return fmt.Errorf("missing key %q in actual JSON", key)
		}

		if err := compareValues(key, expectedVal, actualVal); err != nil {
			return err
		}
	}

	// Check for extra keys in actual
	for key := range actual {
		if _, exists := expected[key]; !exists {
			return fmt.Errorf("unexpected key %q in actual JSON", key)
		}
	}

	return nil
}

// compareValues compares two values semantically.
func compareValues(path string, expected, actual interface{}) error {
	// Normalize values
	expected = normalizeValue(expected)
	actual = normalizeValue(actual)

	// Handle nil
	if expected == nil && actual == nil {
		return nil
	}
	if expected == nil || actual == nil {
		return fmt.Errorf("at %q: expected %v, got %v", path, expected, actual)
	}

	// Get types
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	// Handle numeric comparisons (allow float64 vs int)
	if isNumeric(expected) && isNumeric(actual) {
		if !numericEqual(expected, actual) {
			return fmt.Errorf("at %q: expected %v, got %v", path, expected, actual)
		}
		return nil
	}

	// Types must match for non-numeric values
	if expectedType != actualType {
		return fmt.Errorf("at %q: type mismatch - expected %T, got %T", path, expected, actual)
	}

	// Compare based on type
	switch exp := expected.(type) {
	case map[string]interface{}:
		act := actual.(map[string]interface{})
		for key, val := range exp {
			actVal, exists := act[key]
			if !exists {
				return fmt.Errorf("at %q: missing key %q", path, key)
			}
			if err := compareValues(path+"."+key, val, actVal); err != nil {
				return err
			}
		}
		// Check for extra keys
		for key := range act {
			if _, exists := exp[key]; !exists {
				return fmt.Errorf("at %q: unexpected key %q", path, key)
			}
		}

	case []interface{}:
		act := actual.([]interface{})
		if len(exp) != len(act) {
			return fmt.Errorf("at %q: array length mismatch - expected %d, got %d", path, len(exp), len(act))
		}
		for i := range exp {
			if err := compareValues(fmt.Sprintf("%s[%d]", path, i), exp[i], act[i]); err != nil {
				return err
			}
		}

	case string:
		if exp != actual.(string) {
			return fmt.Errorf("at %q: expected %q, got %q", path, exp, actual)
		}

	case bool:
		if exp != actual.(bool) {
			return fmt.Errorf("at %q: expected %v, got %v", path, exp, actual)
		}

	default:
		if !reflect.DeepEqual(expected, actual) {
			return fmt.Errorf("at %q: expected %v, got %v", path, expected, actual)
		}
	}

	return nil
}

// normalizeMap normalizes a map by removing omitempty fields with zero values.
func normalizeMap(m map[string]interface{}) {
	for key, val := range m {
		switch v := val.(type) {
		case map[string]interface{}:
			normalizeMap(v)
			// Remove empty maps
			if len(v) == 0 {
				delete(m, key)
			}
		case []interface{}:
			if len(v) == 0 {
				delete(m, key)
			} else {
				for _, item := range v {
					if itemMap, ok := item.(map[string]interface{}); ok {
						normalizeMap(itemMap)
					}
				}
			}
		case nil:
			delete(m, key)
		}
	}
}

// normalizeValue normalizes a value for comparison.
func normalizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case float64:
		// If the float is actually an integer, keep it as float64 for consistent comparison
		return val
	case map[string]interface{}:
		normalizeMap(val)
		return val
	default:
		return val
	}
}

// isNumeric checks if a value is numeric.
func isNumeric(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	}
	return false
}

// numericEqual compares two numeric values.
func numericEqual(a, b interface{}) bool {
	// Convert both to float64 for comparison
	var aFloat, bFloat float64

	switch v := a.(type) {
	case int:
		aFloat = float64(v)
	case int64:
		aFloat = float64(v)
	case float64:
		aFloat = v
	}

	switch v := b.(type) {
	case int:
		bFloat = float64(v)
	case int64:
		bFloat = float64(v)
	case float64:
		bFloat = v
	}

	return aFloat == bFloat
}

// discoveredToQuery converts a DiscoveredQuery to a query.Query.
func discoveredToQuery(dq discovery.DiscoveredQuery) query.Query {
	q := query.Query{
		Dataset: dq.Dataset,
		TimeRange: query.TimeRange{
			TimeRange: dq.TimeRange.TimeRange,
			StartTime: dq.TimeRange.StartTime,
			EndTime:   dq.TimeRange.EndTime,
		},
		Breakdowns:        dq.Breakdowns,
		FilterCombination: dq.FilterCombination,
		Limit:             dq.Limit,
		Granularity:       dq.Granularity,
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

	for _, o := range dq.Orders {
		q.Orders = append(q.Orders, query.Order{
			Column: o.Column,
			Op:     o.Op,
			Order:  o.Order,
		})
	}

	return q
}

// generateGoCode generates Go code from a parsed query JSON.
// This is a copy of the logic from cmd/wetwire-honeycomb/main.go.
func generateGoCode(pkg, name string, raw map[string]interface{}) string {
	var code string
	code += fmt.Sprintf("package %s\n\n", pkg)
	code += "import \"github.com/lex00/wetwire-honeycomb-go/query\"\n\n"
	code += fmt.Sprintf("var %s = query.Query{\n", name)

	// Dataset (not present in fixture, but would be required in real queries)
	code += "\tDataset: \"test-dataset\",\n"

	// Time range
	if tr, ok := raw["time_range"].(float64); ok {
		hours := int(tr) / 3600
		if hours*3600 == int(tr) {
			code += fmt.Sprintf("\tTimeRange: query.Hours(%d),\n", hours)
		} else {
			code += fmt.Sprintf("\tTimeRange: query.Seconds(%d),\n", int(tr))
		}
	}

	// Breakdowns
	if breakdowns, ok := raw["breakdowns"].([]interface{}); ok && len(breakdowns) > 0 {
		code += "\tBreakdowns: []string{"
		for i, b := range breakdowns {
			if i > 0 {
				code += ", "
			}
			code += fmt.Sprintf("%q", b)
		}
		code += "},\n"
	}

	// Calculations
	if calcs, ok := raw["calculations"].([]interface{}); ok && len(calcs) > 0 {
		code += "\tCalculations: []query.Calculation{\n"
		for _, c := range calcs {
			cm := c.(map[string]interface{})
			op := cm["op"].(string)
			col, hasCol := cm["column"].(string)

			if op == "COUNT" && !hasCol {
				code += "\t\tquery.Count(),\n"
			} else if hasCol {
				switch op {
				case "P99":
					code += fmt.Sprintf("\t\tquery.P99(%q),\n", col)
				case "P95":
					code += fmt.Sprintf("\t\tquery.P95(%q),\n", col)
				case "P90":
					code += fmt.Sprintf("\t\tquery.P90(%q),\n", col)
				case "P75":
					code += fmt.Sprintf("\t\tquery.P75(%q),\n", col)
				case "P50":
					code += fmt.Sprintf("\t\tquery.P50(%q),\n", col)
				case "AVG":
					code += fmt.Sprintf("\t\tquery.Avg(%q),\n", col)
				case "SUM":
					code += fmt.Sprintf("\t\tquery.Sum(%q),\n", col)
				case "MIN":
					code += fmt.Sprintf("\t\tquery.Min(%q),\n", col)
				case "MAX":
					code += fmt.Sprintf("\t\tquery.Max(%q),\n", col)
				case "COUNT_DISTINCT":
					code += fmt.Sprintf("\t\tquery.CountDistinct(%q),\n", col)
				default:
					code += fmt.Sprintf("\t\t{Op: %q, Column: %q},\n", op, col)
				}
			}
		}
		code += "\t},\n"
	}

	// Filters
	if filters, ok := raw["filters"].([]interface{}); ok && len(filters) > 0 {
		code += "\tFilters: []query.Filter{\n"
		for _, f := range filters {
			fm := f.(map[string]interface{})
			col := fm["column"].(string)
			op := fm["op"].(string)
			val := fm["value"]

			switch op {
			case "=":
				code += fmt.Sprintf("\t\tquery.Equals(%q, %v),\n", col, formatValue(val))
			case "!=":
				code += fmt.Sprintf("\t\tquery.NotEquals(%q, %v),\n", col, formatValue(val))
			case ">":
				code += fmt.Sprintf("\t\tquery.GT(%q, %v),\n", col, formatValue(val))
			case ">=":
				code += fmt.Sprintf("\t\tquery.GTE(%q, %v),\n", col, formatValue(val))
			case "<":
				code += fmt.Sprintf("\t\tquery.LT(%q, %v),\n", col, formatValue(val))
			case "<=":
				code += fmt.Sprintf("\t\tquery.LTE(%q, %v),\n", col, formatValue(val))
			case "contains":
				code += fmt.Sprintf("\t\tquery.Contains(%q, %v),\n", col, formatValue(val))
			case "does-not-contain":
				code += fmt.Sprintf("\t\tquery.DoesNotContain(%q, %v),\n", col, formatValue(val))
			case "exists":
				code += fmt.Sprintf("\t\tquery.Exists(%q),\n", col)
			case "does-not-exist":
				code += fmt.Sprintf("\t\tquery.DoesNotExist(%q),\n", col)
			default:
				code += fmt.Sprintf("\t\t{Column: %q, Op: %q, Value: %v},\n", col, op, formatValue(val))
			}
		}
		code += "\t},\n"
	}

	// Filter combination
	if fc, ok := raw["filter_combination"].(string); ok && fc != "" {
		code += fmt.Sprintf("\tFilterCombination: %q,\n", fc)
	}

	// Orders
	if orders, ok := raw["orders"].([]interface{}); ok && len(orders) > 0 {
		code += "\tOrders: []query.Order{\n"
		for _, o := range orders {
			om := o.(map[string]interface{})
			order := om["order"].(string)

			if col, hasCol := om["column"].(string); hasCol {
				code += fmt.Sprintf("\t\t{Column: %q, Order: %q},\n", col, order)
			} else if op, hasOp := om["op"].(string); hasOp {
				code += fmt.Sprintf("\t\t{Op: %q, Order: %q},\n", op, order)
			}
		}
		code += "\t},\n"
	}

	// Limit
	if limit, ok := raw["limit"].(float64); ok && limit > 0 {
		code += fmt.Sprintf("\tLimit: %d,\n", int(limit))
	}

	// Granularity
	if granularity, ok := raw["granularity"].(float64); ok && granularity > 0 {
		code += fmt.Sprintf("\tGranularity: %d,\n", int(granularity))
	}

	code += "}\n"
	return code
}

// formatValue formats a value for Go code generation.
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case float64:
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
