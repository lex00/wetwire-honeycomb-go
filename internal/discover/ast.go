package discovery

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// extractStringLiteral extracts a string value from an expression.
func extractStringLiteral(expr ast.Expr) string {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		val, _ := strconv.Unquote(lit.Value)
		return val
	}
	return ""
}

// extractIntLiteral extracts an int value from an expression.
func extractIntLiteral(expr ast.Expr) int {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.INT {
		val, _ := strconv.Atoi(lit.Value)
		return val
	}
	return 0
}

// extractStringSlice extracts a slice of strings from a composite literal.
func extractStringSlice(expr ast.Expr) []string {
	var result []string

	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return result
	}

	for _, elt := range comp.Elts {
		if s := extractStringLiteral(elt); s != "" {
			result = append(result, s)
		}
	}

	return result
}

// isQueryType checks if a type expression refers to query.Query.
func isQueryType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		// Check for query.Query
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name == "query" && t.Sel.Name == "Query"
		}
	case *ast.Ident:
		// Direct reference to Query (if imported as .)
		return t.Name == "Query"
	}
	return false
}

// isQueryCompositeLit checks if a composite literal is a Query type.
func isQueryCompositeLit(comp *ast.CompositeLit) bool {
	if comp.Type == nil {
		return false
	}
	return isQueryType(comp.Type)
}

// extractFieldValue extracts the value expression for a named field in a composite literal.
func extractFieldValue(comp *ast.CompositeLit, fieldName string) ast.Expr {
	for _, elt := range comp.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		if key.Name == fieldName {
			return kv.Value
		}
	}
	return nil
}

// extractCalculations extracts calculation information from a composite literal.
func extractCalculations(expr ast.Expr) []Calculation {
	var result []Calculation

	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return result
	}

	for _, elt := range comp.Elts {
		if calc := extractCalculation(elt); calc.Op != "" {
			result = append(result, calc)
		}
	}

	return result
}

// extractCalculation extracts a single calculation from an expression.
func extractCalculation(expr ast.Expr) Calculation {
	var calc Calculation

	// Handle query.P99("column"), query.Count(), etc.
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "query" {
				calc.Op = normalizeCalculationOp(sel.Sel.Name)

				// Extract column argument if present
				if len(call.Args) > 0 {
					calc.Column = extractStringLiteral(call.Args[0])
				}
			}
		}
	}

	// Handle composite literal: query.Calculation{Op: "P99", Column: "duration"}
	if comp, ok := expr.(*ast.CompositeLit); ok {
		if op := extractFieldValue(comp, "Op"); op != nil {
			calc.Op = normalizeCalculationOp(extractStringLiteral(op))
		}
		if col := extractFieldValue(comp, "Column"); col != nil {
			calc.Column = extractStringLiteral(col)
		}
		if alias := extractFieldValue(comp, "Alias"); alias != nil {
			calc.Alias = extractStringLiteral(alias)
		}
	}

	return calc
}

// normalizeCalculationOp normalizes calculation operation names to uppercase Honeycomb format.
func normalizeCalculationOp(funcName string) string {
	mapping := map[string]string{
		"Count":         "COUNT",
		"CountDistinct": "COUNT_DISTINCT",
		"Sum":           "SUM",
		"Avg":           "AVG",
		"Max":           "MAX",
		"Min":           "MIN",
		"P50":           "P50",
		"P75":           "P75",
		"P90":           "P90",
		"P95":           "P95",
		"P99":           "P99",
		"P999":          "P999",
		"Heatmap":       "HEATMAP",
		"Rate":          "RATE",
		"RateSum":       "RATE_SUM",
		"RateAvg":       "RATE_AVG",
		"RateMax":       "RATE_MAX",
		"Concurrency":   "CONCURRENCY",
	}
	if op, ok := mapping[funcName]; ok {
		return op
	}
	return strings.ToUpper(funcName)
}

// extractTimeRange extracts time range information from an expression.
func extractTimeRange(expr ast.Expr) TimeRange {
	var tr TimeRange

	// Handle query.Hours(n), query.Minutes(n), query.Days(n) function calls
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "query" {
				funcName := sel.Sel.Name
				if len(call.Args) > 0 {
					n := extractIntLiteral(call.Args[0])
					switch funcName {
					case "Hours", "LastNHours":
						tr.TimeRange = n * 3600
					case "Days":
						tr.TimeRange = n * 86400
					case "Minutes":
						tr.TimeRange = n * 60
					case "Seconds":
						tr.TimeRange = n
					case "Last24Hours":
						tr.TimeRange = 24 * 3600
					case "Last7Days":
						tr.TimeRange = 7 * 86400
					}
				} else if funcName == "Last24Hours" {
					tr.TimeRange = 24 * 3600
				} else if funcName == "Last7Days" {
					tr.TimeRange = 7 * 86400
				}
			}
		}
	}

	// Handle composite literal: query.TimeRange{TimeRange: 7200}
	if comp, ok := expr.(*ast.CompositeLit); ok {
		if timeRange := extractFieldValue(comp, "TimeRange"); timeRange != nil {
			tr.TimeRange = extractIntLiteral(timeRange)
		}
		if startTime := extractFieldValue(comp, "StartTime"); startTime != nil {
			tr.StartTime = extractIntLiteral(startTime)
		}
		if endTime := extractFieldValue(comp, "EndTime"); endTime != nil {
			tr.EndTime = extractIntLiteral(endTime)
		}
	}

	return tr
}

// getIdentifierName extracts the name from a variable/const declaration.
func getIdentifierName(spec *ast.ValueSpec) string {
	if len(spec.Names) > 0 {
		return spec.Names[0].Name
	}
	return ""
}

// getFunctionName extracts the function name from a function declaration.
func getFunctionName(decl *ast.FuncDecl) string {
	if decl.Name != nil {
		return decl.Name.Name
	}
	return ""
}

// findQueryComposites recursively finds all query.Query composite literals in an expression.
func findQueryComposites(expr ast.Expr) []*ast.CompositeLit {
	var result []*ast.CompositeLit

	ast.Inspect(expr, func(n ast.Node) bool {
		if comp, ok := n.(*ast.CompositeLit); ok {
			if isQueryCompositeLit(comp) {
				result = append(result, comp)
			}
			// Check if this is a composite that might contain an embedded Query
			if comp.Type != nil {
				// Look for Query field in the composite
				for _, elt := range comp.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						if key, ok := kv.Key.(*ast.Ident); ok && key.Name == "Query" {
							// This is an embedded Query field
							if innerComp, ok := kv.Value.(*ast.CompositeLit); ok {
								if isQueryCompositeLit(innerComp) {
									result = append(result, innerComp)
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	return result
}

// isExportedName checks if a name is exported (starts with capital letter).
func isExportedName(name string) bool {
	if name == "" {
		return false
	}
	r := []rune(name)
	return len(r) > 0 && r[0] >= 'A' && r[0] <= 'Z'
}

// qualifyTypeName returns a qualified type name from an expression.
// Currently unused but may be needed for future type qualification features.
// nolint:unused
func qualifyTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name + "." + t.Sel.Name
		}
	case *ast.Ident:
		return t.Name
	}
	return ""
}

// extractFilters extracts filter information from a composite literal.
func extractFilters(expr ast.Expr) []Filter {
	var result []Filter

	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return result
	}

	for _, elt := range comp.Elts {
		if filter := extractFilter(elt); filter.Column != "" {
			result = append(result, filter)
		}
	}

	return result
}

// extractFilter extracts a single filter from an expression.
func extractFilter(expr ast.Expr) Filter {
	var filter Filter

	// Handle query.GT("duration_ms", 500), query.Equals("status", "error"), etc.
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "query" {
				funcName := sel.Sel.Name

				// Extract column argument (first arg)
				if len(call.Args) > 0 {
					filter.Column = extractStringLiteral(call.Args[0])
				}

				// Map function name to operator
				filter.Op = mapFilterFuncToOp(funcName)

				// Extract value argument (second arg) - can be string or int
				if len(call.Args) > 1 {
					if s := extractStringLiteral(call.Args[1]); s != "" {
						filter.Value = s
					} else {
						filter.Value = extractIntLiteral(call.Args[1])
					}
				}
			}
		}
	}

	// Handle composite literal: query.Filter{Column: "duration_ms", Op: ">", Value: 500}
	if comp, ok := expr.(*ast.CompositeLit); ok {
		if col := extractFieldValue(comp, "Column"); col != nil {
			filter.Column = extractStringLiteral(col)
		}
		if op := extractFieldValue(comp, "Op"); op != nil {
			filter.Op = extractStringLiteral(op)
		}
		if val := extractFieldValue(comp, "Value"); val != nil {
			if s := extractStringLiteral(val); s != "" {
				filter.Value = s
			} else {
				filter.Value = extractIntLiteral(val)
			}
		}
	}

	return filter
}

// mapFilterFuncToOp maps filter function names to operators.
func mapFilterFuncToOp(funcName string) string {
	mapping := map[string]string{
		"GT":                 ">",
		"GTE":                ">=",
		"LT":                 "<",
		"LTE":                "<=",
		"Equals":             "=",
		"Eq":                 "=",
		"NotEquals":          "!=",
		"Ne":                 "!=",
		"GreaterThan":        ">",
		"GreaterThanOrEqual": ">=",
		"LessThan":           "<",
		"LessThanOrEqual":    "<=",
		"Contains":           "contains",
		"DoesNotContain":     "does-not-contain",
		"Exists":             "exists",
		"DoesNotExist":       "does-not-exist",
		"StartsWith":         "starts-with",
		"In":                 "in",
		"NotIn":              "not-in",
	}
	if op, ok := mapping[funcName]; ok {
		return op
	}
	return strings.ToLower(funcName)
}

// extractOrders extracts order specifications from a composite literal.
func extractOrders(expr ast.Expr) []Order {
	var result []Order

	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return result
	}

	for _, elt := range comp.Elts {
		if order := extractOrder(elt); order.Column != "" || order.Op != "" {
			result = append(result, order)
		}
	}

	return result
}

// extractOrder extracts a single order specification from an expression.
func extractOrder(expr ast.Expr) Order {
	var order Order

	// Handle composite literal: query.Order{Column: "endpoint", Order: "descending"}
	if comp, ok := expr.(*ast.CompositeLit); ok {
		if col := extractFieldValue(comp, "Column"); col != nil {
			order.Column = extractStringLiteral(col)
		}
		if op := extractFieldValue(comp, "Op"); op != nil {
			order.Op = extractStringLiteral(op)
		}
		if ord := extractFieldValue(comp, "Order"); ord != nil {
			order.Order = extractStringLiteral(ord)
		}
	}

	return order
}

// extractStyleMetadata extracts style metadata from a query composite literal.
func extractStyleMetadata(comp *ast.CompositeLit) StyleMetadata {
	var meta StyleMetadata

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
		case "Calculations":
			meta.InlineCalculationCount = countInlineDefinitions(kv.Value)
		case "Filters":
			meta.InlineFilterCount = countInlineDefinitions(kv.Value)
		}

		// Check for raw map literals in any field
		if hasRawMapLiteral(kv.Value) {
			meta.HasRawMapLiteral = true
		}
	}

	// Calculate nesting depth
	meta.MaxNestingDepth = calculateNestingDepth(comp)

	return meta
}

// countInlineDefinitions counts the number of inline composite literal definitions
// in a slice expression. An inline definition is a composite literal (e.g., query.Calculation{...})
// rather than a reference to a named variable.
func countInlineDefinitions(expr ast.Expr) int {
	count := 0

	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return 0
	}

	for _, elt := range comp.Elts {
		switch e := elt.(type) {
		case *ast.CompositeLit:
			// Composite literal (inline definition)
			count++
		case *ast.CallExpr:
			// Function call like query.P99("col") is considered inline
			// but helper functions like query.Count() are acceptable
			if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "query" {
					// This is a typed builder call, which is fine
					continue
				}
			}
			count++
		case *ast.Ident:
			// Reference to a named variable - not inline
			continue
		default:
			// Unknown, count as inline to be safe
			count++
		}
	}

	return count
}

// hasRawMapLiteral checks if an expression contains raw map literals
// instead of typed structs.
func hasRawMapLiteral(expr ast.Expr) bool {
	found := false

	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}

		if comp, ok := n.(*ast.CompositeLit); ok {
			// Check if this is a map type
			if mapType, ok := comp.Type.(*ast.MapType); ok {
				// This is a map literal - check if it's map[string]interface{} or similar
				if _, ok := mapType.Key.(*ast.Ident); ok {
					found = true
					return false
				}
			}
		}

		return true
	})

	return found
}

// calculateNestingDepth calculates the maximum nesting depth of composite literals.
func calculateNestingDepth(expr ast.Expr) int {
	maxDepth := 0

	var walk func(n ast.Node, depth int)
	walk = func(n ast.Node, depth int) {
		if n == nil {
			return
		}

		if _, ok := n.(*ast.CompositeLit); ok {
			if depth > maxDepth {
				maxDepth = depth
			}
			depth++
		}

		ast.Inspect(n, func(child ast.Node) bool {
			if child == n {
				return true
			}
			if child != nil {
				walk(child, depth)
			}
			return false
		})
	}

	walk(expr, 0)

	return maxDepth
}
