package discovery

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// DiscoveredSLO represents a discovered SLO definition with metadata.
type DiscoveredSLO struct {
	// Name is the identifier of the SLO (variable name)
	Name string

	// Package is the package name where the SLO is defined
	Package string

	// File is the absolute path to the file containing the SLO
	File string

	// Line is the line number where the SLO is defined
	Line int

	// SLOName is the SLO.Name field value
	SLOName string

	// Description is the SLO.Description field value
	Description string

	// Dataset is the Honeycomb dataset
	Dataset string

	// TargetPercentage is the target SLO percentage
	TargetPercentage float64

	// TimePeriodDays is the rolling window in days
	TimePeriodDays int

	// GoodEventsQueryRef is the name of the good events query
	GoodEventsQueryRef string

	// TotalEventsQueryRef is the name of the total events query
	TotalEventsQueryRef string

	// BurnAlertCount is the number of burn alerts configured
	BurnAlertCount int
}

// DiscoverSLOs discovers all SLO definitions in the specified directory.
func DiscoverSLOs(dir string) ([]DiscoveredSLO, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", dir)
	}

	var discovered []DiscoveredSLO

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		slos, err := discoverSLOsInFile(path)
		if err != nil {
			return nil
		}

		discovered = append(discovered, slos...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return discovered, nil
}

// discoverSLOsInFile discovers SLOs in a single Go source file.
func discoverSLOsInFile(path string) ([]DiscoveredSLO, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var discovered []DiscoveredSLO

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	packageName := node.Name.Name

	ast.Inspect(node, func(n ast.Node) bool {
		if decl, ok := n.(*ast.GenDecl); ok {
			if decl.Tok == token.VAR || decl.Tok == token.CONST {
				for _, spec := range decl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						slos := extractSLOsFromValueSpec(valueSpec, fset, absPath, packageName)
						discovered = append(discovered, slos...)
					}
				}
			}
		}
		return true
	})

	return discovered, nil
}

// extractSLOsFromValueSpec extracts SLOs from a variable declaration.
func extractSLOsFromValueSpec(spec *ast.ValueSpec, fset *token.FileSet, file string, pkg string) []DiscoveredSLO {
	var discovered []DiscoveredSLO

	name := getIdentifierName(spec)
	if name == "" || !isExportedName(name) {
		return discovered
	}

	for _, value := range spec.Values {
		composites := findSLOComposites(value)
		for _, comp := range composites {
			slo := extractSLOFromComposite(comp, fset, file, pkg, name)
			if slo.Name != "" {
				discovered = append(discovered, slo)
			}
		}
	}

	return discovered
}

// findSLOComposites finds all slo.SLO composite literals in an expression.
func findSLOComposites(expr ast.Expr) []*ast.CompositeLit {
	var result []*ast.CompositeLit

	ast.Inspect(expr, func(n ast.Node) bool {
		if comp, ok := n.(*ast.CompositeLit); ok {
			if isSLOType(comp.Type) {
				result = append(result, comp)
			}
		}
		return true
	})

	return result
}

// isSLOType checks if a type expression refers to slo.SLO.
func isSLOType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "slo" && sel.Sel.Name == "SLO"
		}
	}
	return false
}

// extractSLOFromComposite extracts SLO metadata from a composite literal.
func extractSLOFromComposite(comp *ast.CompositeLit, fset *token.FileSet, file string, pkg string, name string) DiscoveredSLO {
	slo := DiscoveredSLO{
		Name:    name,
		Package: pkg,
		File:    file,
		Line:    fset.Position(comp.Pos()).Line,
	}

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
		case "Name":
			slo.SLOName = extractStringLiteral(kv.Value)
		case "Description":
			slo.Description = extractStringLiteral(kv.Value)
		case "Dataset":
			slo.Dataset = extractStringLiteral(kv.Value)
		case "Target":
			slo.TargetPercentage = extractTargetPercentage(kv.Value)
		case "TimePeriod":
			slo.TimePeriodDays = extractTimePeriodDays(kv.Value)
		case "SLI":
			slo.GoodEventsQueryRef, slo.TotalEventsQueryRef = extractSLIQueryRefs(kv.Value)
		case "BurnAlerts":
			slo.BurnAlertCount = extractBurnAlertCount(kv.Value)
		}
	}

	return slo
}

// extractTargetPercentage extracts the percentage from a Target field.
func extractTargetPercentage(expr ast.Expr) float64 {
	// Handle slo.Percentage(99.9)
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "slo" {
				if sel.Sel.Name == "Percentage" && len(call.Args) > 0 {
					return extractFloatLiteral(call.Args[0])
				}
			}
		}
	}
	return 0
}

// extractTimePeriodDays extracts the days from a TimePeriod field.
func extractTimePeriodDays(expr ast.Expr) int {
	// Handle slo.Days(30)
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "slo" {
				if sel.Sel.Name == "Days" && len(call.Args) > 0 {
					return extractIntLiteral(call.Args[0])
				}
			}
		}
	}
	return 0
}

// extractSLIQueryRefs extracts query references from an SLI field.
func extractSLIQueryRefs(expr ast.Expr) (string, string) {
	var goodRef, totalRef string

	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return "", ""
	}

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
		case "GoodEvents":
			if ident, ok := kv.Value.(*ast.Ident); ok {
				goodRef = ident.Name
			}
		case "TotalEvents":
			if ident, ok := kv.Value.(*ast.Ident); ok {
				totalRef = ident.Name
			}
		}
	}

	return goodRef, totalRef
}

// extractBurnAlertCount counts burn alerts from a BurnAlerts field.
func extractBurnAlertCount(expr ast.Expr) int {
	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return 0
	}
	return len(comp.Elts)
}

// extractFloatLiteral extracts a float64 value from an expression.
func extractFloatLiteral(expr ast.Expr) float64 {
	if lit, ok := expr.(*ast.BasicLit); ok {
		if lit.Kind == token.FLOAT || lit.Kind == token.INT {
			val, _ := strconv.ParseFloat(lit.Value, 64)
			return val
		}
	}
	return 0
}
