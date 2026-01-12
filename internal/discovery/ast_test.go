package discovery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNormalizeCalculationOp(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Count", "COUNT"},
		{"CountDistinct", "COUNT_DISTINCT"},
		{"Sum", "SUM"},
		{"Avg", "AVG"},
		{"Max", "MAX"},
		{"Min", "MIN"},
		{"P50", "P50"},
		{"P75", "P75"},
		{"P90", "P90"},
		{"P95", "P95"},
		{"P99", "P99"},
		{"P999", "P999"},
		{"Heatmap", "HEATMAP"},
		{"Rate", "RATE"},
		{"RateSum", "RATE_SUM"},
		{"RateAvg", "RATE_AVG"},
		{"RateMax", "RATE_MAX"},
		{"Concurrency", "CONCURRENCY"},
		// Unknown ops should be uppercased
		{"UnknownOp", "UNKNOWNOP"},
		{"custom", "CUSTOM"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := normalizeCalculationOp(tc.input)
			if result != tc.expected {
				t.Errorf("normalizeCalculationOp(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestMapFilterFuncToOp(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"GT", ">"},
		{"GTE", ">="},
		{"LT", "<"},
		{"LTE", "<="},
		{"Equals", "="},
		{"Eq", "="},
		{"NotEquals", "!="},
		{"Ne", "!="},
		{"GreaterThan", ">"},
		{"GreaterThanOrEqual", ">="},
		{"LessThan", "<"},
		{"LessThanOrEqual", "<="},
		{"Contains", "contains"},
		{"DoesNotContain", "does-not-contain"},
		{"Exists", "exists"},
		{"DoesNotExist", "does-not-exist"},
		{"StartsWith", "starts-with"},
		{"In", "in"},
		{"NotIn", "not-in"},
		// Unknown funcs should be lowercased
		{"UnknownFunc", "unknownfunc"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := mapFilterFuncToOp(tc.input)
			if result != tc.expected {
				t.Errorf("mapFilterFuncToOp(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsExportedName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"ExportedName", true},
		{"unexportedName", false},
		{"", false},
		{"A", true},
		{"a", false},
		{"URL", true},
		{"urlParser", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := isExportedName(tc.input)
			if result != tc.expected {
				t.Errorf("isExportedName(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestExtractStringLiteral(t *testing.T) {
	// Test with actual AST node
	src := `package test
var s = "hello"`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Find the string literal
	var strLit ast.Expr
	ast.Inspect(file, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			strLit = lit
			return false
		}
		return true
	})

	if strLit == nil {
		t.Fatal("String literal not found")
	}

	result := extractStringLiteral(strLit)
	if result != "hello" {
		t.Errorf("extractStringLiteral() = %q, want %q", result, "hello")
	}

	// Test with non-string expression
	intSrc := `package test
var i = 42`
	intFile, _ := parser.ParseFile(fset, "test2.go", intSrc, 0)
	var intLit ast.Expr
	ast.Inspect(intFile, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.INT {
			intLit = lit
			return false
		}
		return true
	})

	result = extractStringLiteral(intLit)
	if result != "" {
		t.Errorf("extractStringLiteral(int) = %q, want empty string", result)
	}
}

func TestExtractIntLiteral(t *testing.T) {
	src := `package test
var i = 42`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var intLit ast.Expr
	ast.Inspect(file, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.INT {
			intLit = lit
			return false
		}
		return true
	})

	if intLit == nil {
		t.Fatal("Int literal not found")
	}

	result := extractIntLiteral(intLit)
	if result != 42 {
		t.Errorf("extractIntLiteral() = %d, want 42", result)
	}

	// Test with non-int expression
	strSrc := `package test
var s = "hello"`
	strFile, _ := parser.ParseFile(fset, "test2.go", strSrc, 0)
	var strLit ast.Expr
	ast.Inspect(strFile, func(n ast.Node) bool {
		if lit, ok := n.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			strLit = lit
			return false
		}
		return true
	})

	result = extractIntLiteral(strLit)
	if result != 0 {
		t.Errorf("extractIntLiteral(string) = %d, want 0", result)
	}
}

func TestIsQueryType(t *testing.T) {
	// Test query.Query selector
	src := `package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var q = query.Query{}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var compLit *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit = cl
			return false
		}
		return true
	})

	if compLit == nil || compLit.Type == nil {
		t.Fatal("CompositeLit not found")
	}

	if !isQueryType(compLit.Type) {
		t.Error("isQueryType should return true for query.Query")
	}

	// Test non-query type
	src2 := `package test
type Other struct{}
var o = Other{}`
	file2, _ := parser.ParseFile(fset, "test2.go", src2, 0)
	var compLit2 *ast.CompositeLit
	ast.Inspect(file2, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit2 = cl
			return false
		}
		return true
	})

	if compLit2 != nil && compLit2.Type != nil {
		if isQueryType(compLit2.Type) {
			t.Error("isQueryType should return false for Other type")
		}
	}

	// Test direct Query ident (dot import)
	src3 := `package test
type Query struct{}
var q = Query{}`
	file3, _ := parser.ParseFile(fset, "test3.go", src3, 0)
	var compLit3 *ast.CompositeLit
	ast.Inspect(file3, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit3 = cl
			return false
		}
		return true
	})

	if compLit3 != nil && compLit3.Type != nil {
		if !isQueryType(compLit3.Type) {
			t.Error("isQueryType should return true for Query ident")
		}
	}
}

func TestExtractStringSlice(t *testing.T) {
	src := `package test
var s = []string{"a", "b", "c"}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var compLit *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit = cl
			return false
		}
		return true
	})

	if compLit == nil {
		t.Fatal("CompositeLit not found")
	}

	result := extractStringSlice(compLit)
	expected := []string{"a", "b", "c"}
	if len(result) != len(expected) {
		t.Errorf("extractStringSlice() len = %d, want %d", len(result), len(expected))
	}
	for i, v := range expected {
		if i >= len(result) || result[i] != v {
			t.Errorf("extractStringSlice()[%d] = %q, want %q", i, result[i], v)
		}
	}

	// Test with non-composite expression
	result = extractStringSlice(&ast.BasicLit{})
	if len(result) != 0 {
		t.Error("extractStringSlice on non-composite should return empty slice")
	}
}

func TestExtractTimeRange_CompositeLiteral(t *testing.T) {
	src := `package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.TimeRange{TimeRange: 3600, StartTime: 1000, EndTime: 2000}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var compLit *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit = cl
			return false
		}
		return true
	})

	if compLit == nil {
		t.Fatal("CompositeLit not found")
	}

	result := extractTimeRange(compLit)
	if result.TimeRange != 3600 {
		t.Errorf("TimeRange = %d, want 3600", result.TimeRange)
	}
	if result.StartTime != 1000 {
		t.Errorf("StartTime = %d, want 1000", result.StartTime)
	}
	if result.EndTime != 2000 {
		t.Errorf("EndTime = %d, want 2000", result.EndTime)
	}
}

func TestExtractTimeRange_FunctionCalls(t *testing.T) {
	testCases := []struct {
		src      string
		expected int
	}{
		{`package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.Hours(2)`, 7200},
		{`package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.Days(1)`, 86400},
		{`package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.Minutes(30)`, 1800},
		{`package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.Seconds(60)`, 60},
		{`package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.Last24Hours()`, 86400},
		{`package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.Last7Days()`, 604800},
		{`package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var tr = query.LastNHours(3)`, 10800},
	}

	fset := token.NewFileSet()
	for _, tc := range testCases {
		file, err := parser.ParseFile(fset, "test.go", tc.src, 0)
		if err != nil {
			t.Fatalf("Failed to parse: %v", err)
		}

		var callExpr *ast.CallExpr
		ast.Inspect(file, func(n ast.Node) bool {
			if ce, ok := n.(*ast.CallExpr); ok {
				callExpr = ce
				return false
			}
			return true
		})

		if callExpr == nil {
			t.Fatalf("CallExpr not found for %s", tc.src)
		}

		result := extractTimeRange(callExpr)
		if result.TimeRange != tc.expected {
			t.Errorf("extractTimeRange() = %d, want %d for %s", result.TimeRange, tc.expected, tc.src)
		}
	}
}

func TestExtractCalculation_CompositeLiteral(t *testing.T) {
	src := `package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var c = query.Calculation{Op: "P99", Column: "duration_ms", Alias: "p99_latency"}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var compLit *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit = cl
			return false
		}
		return true
	})

	if compLit == nil {
		t.Fatal("CompositeLit not found")
	}

	result := extractCalculation(compLit)
	if result.Op != "P99" {
		t.Errorf("Op = %q, want P99", result.Op)
	}
	if result.Column != "duration_ms" {
		t.Errorf("Column = %q, want duration_ms", result.Column)
	}
	if result.Alias != "p99_latency" {
		t.Errorf("Alias = %q, want p99_latency", result.Alias)
	}
}

func TestExtractFilter_CompositeLiteral(t *testing.T) {
	src := `package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var f = query.Filter{Column: "status", Op: "=", Value: "error"}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var compLit *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit = cl
			return false
		}
		return true
	})

	if compLit == nil {
		t.Fatal("CompositeLit not found")
	}

	result := extractFilter(compLit)
	if result.Column != "status" {
		t.Errorf("Column = %q, want status", result.Column)
	}
	if result.Op != "=" {
		t.Errorf("Op = %q, want =", result.Op)
	}
	if result.Value != "error" {
		t.Errorf("Value = %v, want error", result.Value)
	}
}

func TestExtractFilter_IntValue(t *testing.T) {
	src := `package test
import "github.com/lex00/wetwire-honeycomb-go/query"
var f = query.Filter{Column: "count", Op: ">", Value: 100}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var compLit *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit = cl
			return false
		}
		return true
	})

	if compLit == nil {
		t.Fatal("CompositeLit not found")
	}

	result := extractFilter(compLit)
	if result.Value != 100 {
		t.Errorf("Value = %v, want 100", result.Value)
	}
}

func TestGetIdentifierName(t *testing.T) {
	src := `package test
var TestVar = 1`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var valueSpec *ast.ValueSpec
	ast.Inspect(file, func(n ast.Node) bool {
		if vs, ok := n.(*ast.ValueSpec); ok {
			valueSpec = vs
			return false
		}
		return true
	})

	if valueSpec == nil {
		t.Fatal("ValueSpec not found")
	}

	result := getIdentifierName(valueSpec)
	if result != "TestVar" {
		t.Errorf("getIdentifierName() = %q, want TestVar", result)
	}

	// Test empty names
	emptySpec := &ast.ValueSpec{Names: []*ast.Ident{}}
	result = getIdentifierName(emptySpec)
	if result != "" {
		t.Errorf("getIdentifierName(empty) = %q, want empty", result)
	}
}

func TestGetFunctionName(t *testing.T) {
	src := `package test
func TestFunc() {}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var funcDecl *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			funcDecl = fd
			return false
		}
		return true
	})

	if funcDecl == nil {
		t.Fatal("FuncDecl not found")
	}

	result := getFunctionName(funcDecl)
	if result != "TestFunc" {
		t.Errorf("getFunctionName() = %q, want TestFunc", result)
	}

	// Test nil name
	nilNameDecl := &ast.FuncDecl{}
	result = getFunctionName(nilNameDecl)
	if result != "" {
		t.Errorf("getFunctionName(nil) = %q, want empty", result)
	}
}

func TestExtractOrders_Empty(t *testing.T) {
	// Test with non-composite expression
	result := extractOrders(&ast.BasicLit{})
	if len(result) != 0 {
		t.Error("extractOrders on non-composite should return empty slice")
	}
}

func TestExtractOrder_Empty(t *testing.T) {
	// Test with non-composite expression
	result := extractOrder(&ast.BasicLit{})
	if result.Column != "" || result.Op != "" || result.Order != "" {
		t.Error("extractOrder on non-composite should return empty Order")
	}
}

func TestExtractFieldValue(t *testing.T) {
	src := `package test
type T struct { Field string }
var t = T{Field: "value"}`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var compLit *ast.CompositeLit
	ast.Inspect(file, func(n ast.Node) bool {
		if cl, ok := n.(*ast.CompositeLit); ok {
			compLit = cl
			return false
		}
		return true
	})

	if compLit == nil {
		t.Fatal("CompositeLit not found")
	}

	// Test existing field
	result := extractFieldValue(compLit, "Field")
	if result == nil {
		t.Error("extractFieldValue should find Field")
	}

	// Test non-existing field
	result = extractFieldValue(compLit, "NonExistent")
	if result != nil {
		t.Error("extractFieldValue should return nil for non-existent field")
	}
}

func TestIsQueryCompositeLit_NilType(t *testing.T) {
	comp := &ast.CompositeLit{Type: nil}
	if isQueryCompositeLit(comp) {
		t.Error("isQueryCompositeLit should return false for nil type")
	}
}

func TestExtractFilters_Empty(t *testing.T) {
	// Test with non-composite expression
	result := extractFilters(&ast.BasicLit{})
	if len(result) != 0 {
		t.Error("extractFilters on non-composite should return empty slice")
	}
}

func TestExtractCalculations_Empty(t *testing.T) {
	// Test with non-composite expression
	result := extractCalculations(&ast.BasicLit{})
	if len(result) != 0 {
		t.Error("extractCalculations on non-composite should return empty slice")
	}
}

func TestQualifyTypeName(t *testing.T) {
	// Test SelectorExpr (pkg.Type)
	src := `package test
import "fmt"
var x fmt.Stringer`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	var selExpr *ast.SelectorExpr
	ast.Inspect(file, func(n ast.Node) bool {
		if se, ok := n.(*ast.SelectorExpr); ok {
			selExpr = se
			return false
		}
		return true
	})

	if selExpr != nil {
		result := qualifyTypeName(selExpr)
		if result != "fmt.Stringer" {
			t.Errorf("qualifyTypeName(SelectorExpr) = %q, want fmt.Stringer", result)
		}
	}

	// Test Ident
	src2 := `package test
type MyType int
var x MyType`
	file2, _ := parser.ParseFile(fset, "test2.go", src2, 0)
	var ident *ast.Ident
	ast.Inspect(file2, func(n ast.Node) bool {
		if vs, ok := n.(*ast.ValueSpec); ok && vs.Type != nil {
			if id, ok := vs.Type.(*ast.Ident); ok {
				ident = id
				return false
			}
		}
		return true
	})

	if ident != nil {
		result := qualifyTypeName(ident)
		if result != "MyType" {
			t.Errorf("qualifyTypeName(Ident) = %q, want MyType", result)
		}
	}

	// Test unsupported type
	result := qualifyTypeName(&ast.StarExpr{})
	if result != "" {
		t.Errorf("qualifyTypeName(StarExpr) = %q, want empty", result)
	}
}
