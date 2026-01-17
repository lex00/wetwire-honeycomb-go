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

// DiscoveredTrigger represents a discovered trigger definition with metadata.
type DiscoveredTrigger struct {
	// Name is the identifier of the trigger (variable name)
	Name string

	// Package is the package name where the trigger is defined
	Package string

	// File is the absolute path to the file containing the trigger
	File string

	// Line is the line number where the trigger is defined
	Line int

	// TriggerName is the Trigger.Name field value
	TriggerName string

	// Description is the Trigger.Description field value
	Description string

	// Dataset is the Honeycomb dataset
	Dataset string

	// QueryRef is the name of the referenced query
	QueryRef string

	// ThresholdOp is the threshold operator (>, >=, <, <=)
	ThresholdOp string

	// ThresholdValue is the threshold value
	ThresholdValue float64

	// FrequencySeconds is the evaluation frequency in seconds
	FrequencySeconds int

	// RecipientCount is the number of recipients configured
	RecipientCount int

	// Disabled indicates if the trigger is disabled
	Disabled bool
}

// DiscoverTriggers discovers all Trigger definitions in the specified directory.
func DiscoverTriggers(dir string) ([]DiscoveredTrigger, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", dir)
	}

	var discovered []DiscoveredTrigger

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		triggers, err := discoverTriggersInFile(path)
		if err != nil {
			return nil
		}

		discovered = append(discovered, triggers...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return discovered, nil
}

// discoverTriggersInFile discovers triggers in a single Go source file.
func discoverTriggersInFile(path string) ([]DiscoveredTrigger, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var discovered []DiscoveredTrigger

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
						triggers := extractTriggersFromValueSpec(valueSpec, fset, absPath, packageName)
						discovered = append(discovered, triggers...)
					}
				}
			}
		}
		return true
	})

	return discovered, nil
}

// extractTriggersFromValueSpec extracts triggers from a variable declaration.
func extractTriggersFromValueSpec(spec *ast.ValueSpec, fset *token.FileSet, file string, pkg string) []DiscoveredTrigger {
	var discovered []DiscoveredTrigger

	name := getIdentifierName(spec)
	if name == "" || !isExportedName(name) {
		return discovered
	}

	for _, value := range spec.Values {
		composites := findTriggerComposites(value)
		for _, comp := range composites {
			trigger := extractTriggerFromComposite(comp, fset, file, pkg, name)
			if trigger.Name != "" {
				discovered = append(discovered, trigger)
			}
		}
	}

	return discovered
}

// findTriggerComposites finds all trigger.Trigger composite literals in an expression.
func findTriggerComposites(expr ast.Expr) []*ast.CompositeLit {
	var result []*ast.CompositeLit

	ast.Inspect(expr, func(n ast.Node) bool {
		if comp, ok := n.(*ast.CompositeLit); ok {
			if isTriggerType(comp.Type) {
				result = append(result, comp)
			}
		}
		return true
	})

	return result
}

// isTriggerType checks if a type expression refers to trigger.Trigger.
func isTriggerType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "trigger" && sel.Sel.Name == "Trigger"
		}
	}
	return false
}

// extractTriggerFromComposite extracts trigger metadata from a composite literal.
func extractTriggerFromComposite(comp *ast.CompositeLit, fset *token.FileSet, file string, pkg string, name string) DiscoveredTrigger {
	trigger := DiscoveredTrigger{
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
			trigger.TriggerName = extractStringLiteral(kv.Value)
		case "Description":
			trigger.Description = extractStringLiteral(kv.Value)
		case "Dataset":
			trigger.Dataset = extractStringLiteral(kv.Value)
		case "Query":
			if ident, ok := kv.Value.(*ast.Ident); ok {
				trigger.QueryRef = ident.Name
			}
		case "Threshold":
			trigger.ThresholdOp, trigger.ThresholdValue = extractThreshold(kv.Value)
		case "Frequency":
			trigger.FrequencySeconds = extractFrequencySeconds(kv.Value)
		case "Recipients":
			trigger.RecipientCount = extractRecipientCount(kv.Value)
		case "Disabled":
			trigger.Disabled = extractBoolLiteral(kv.Value)
		}
	}

	return trigger
}

// extractThreshold extracts operator and value from a Threshold field.
func extractThreshold(expr ast.Expr) (string, float64) {
	// Handle trigger.GreaterThan(500), trigger.LessThan(10), etc.
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "trigger" {
				var op string
				switch sel.Sel.Name {
				case "GreaterThan":
					op = ">"
				case "GreaterThanOrEqual":
					op = ">="
				case "LessThan":
					op = "<"
				case "LessThanOrEqual":
					op = "<="
				}
				if op != "" && len(call.Args) > 0 {
					return op, extractFloatLiteral(call.Args[0])
				}
			}
		}
	}
	return "", 0
}

// extractFrequencySeconds extracts seconds from a Frequency field.
func extractFrequencySeconds(expr ast.Expr) int {
	// Handle trigger.Minutes(5), trigger.Seconds(30)
	if call, ok := expr.(*ast.CallExpr); ok {
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "trigger" {
				if len(call.Args) > 0 {
					n := extractIntLiteral(call.Args[0])
					switch sel.Sel.Name {
					case "Minutes":
						return n * 60
					case "Seconds":
						return n
					}
				}
			}
		}
	}
	return 0
}

// extractRecipientCount counts recipients from a Recipients field.
func extractRecipientCount(expr ast.Expr) int {
	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return 0
	}
	return len(comp.Elts)
}

// extractBoolLiteral extracts a bool value from an expression.
func extractBoolLiteral(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "true"
	}
	return false
}
