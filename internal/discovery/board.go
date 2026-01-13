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

// DiscoveredBoard represents a discovered board definition with metadata.
type DiscoveredBoard struct {
	// Name is the identifier of the board (variable name)
	Name string

	// Package is the package name where the board is defined
	Package string

	// File is the absolute path to the file containing the board
	File string

	// Line is the line number where the board is defined
	Line int

	// BoardName is the Board.Name field value
	BoardName string

	// Description is the Board.Description field value
	Description string

	// PanelCount is the number of panels in the board
	PanelCount int

	// QueryRefs are the names of queries referenced by QueryPanel
	QueryRefs []string

	// SLORefs are the IDs of SLOs referenced by SLOPanelByID
	SLORefs []string

	// IsTemplate indicates if the board is generated from a function call
	IsTemplate bool
}

// DiscoverBoards discovers all Board definitions in the specified directory.
func DiscoverBoards(dir string) ([]DiscoveredBoard, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", dir)
	}

	var discovered []DiscoveredBoard

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		boards, err := discoverBoardsInFile(path)
		if err != nil {
			return nil
		}

		discovered = append(discovered, boards...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return discovered, nil
}

// discoverBoardsInFile discovers boards in a single Go source file.
func discoverBoardsInFile(path string) ([]DiscoveredBoard, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var discovered []DiscoveredBoard

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
						boards := extractBoardsFromValueSpec(valueSpec, fset, absPath, packageName)
						discovered = append(discovered, boards...)
					}
				}
			}
		}
		return true
	})

	return discovered, nil
}

// extractBoardsFromValueSpec extracts boards from a variable declaration.
func extractBoardsFromValueSpec(spec *ast.ValueSpec, fset *token.FileSet, file string, pkg string) []DiscoveredBoard {
	var discovered []DiscoveredBoard

	name := getIdentifierName(spec)
	if name == "" || !isExportedName(name) {
		return discovered
	}

	for _, value := range spec.Values {
		composites := findBoardComposites(value)
		for _, comp := range composites {
			board := extractBoardFromComposite(comp, fset, file, pkg, name)
			if board.Name != "" {
				discovered = append(discovered, board)
			}
		}
	}

	return discovered
}

// findBoardComposites finds all board.Board composite literals in an expression.
func findBoardComposites(expr ast.Expr) []*ast.CompositeLit {
	var result []*ast.CompositeLit

	ast.Inspect(expr, func(n ast.Node) bool {
		if comp, ok := n.(*ast.CompositeLit); ok {
			if isBoardType(comp.Type) {
				result = append(result, comp)
			}
		}
		return true
	})

	return result
}

// isBoardType checks if a type expression refers to board.Board.
func isBoardType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "board" && sel.Sel.Name == "Board"
		}
	}
	return false
}

// extractBoardFromComposite extracts board metadata from a composite literal.
func extractBoardFromComposite(comp *ast.CompositeLit, fset *token.FileSet, file string, pkg string, name string) DiscoveredBoard {
	board := DiscoveredBoard{
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
			board.BoardName = extractStringLiteral(kv.Value)
		case "Description":
			board.Description = extractStringLiteral(kv.Value)
		case "Panels":
			board.PanelCount, board.QueryRefs, board.SLORefs = extractPanelInfo(kv.Value)
		}
	}

	return board
}

// extractPanelInfo extracts panel count and references from a Panels field.
func extractPanelInfo(expr ast.Expr) (int, []string, []string) {
	var panelCount int
	var queryRefs []string
	var sloRefs []string

	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return 0, nil, nil
	}

	for _, elt := range comp.Elts {
		panelCount++

		// Check for board.QueryPanel(SomeQuery) or board.SLOPanelByID("id")
		if call, ok := elt.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "board" {
					switch sel.Sel.Name {
					case "QueryPanel":
						// Extract query reference
						if len(call.Args) > 0 {
							if ident, ok := call.Args[0].(*ast.Ident); ok {
								queryRefs = append(queryRefs, ident.Name)
							}
						}
					case "SLOPanelByID":
						// Extract SLO ID
						if len(call.Args) > 0 {
							if id := extractStringLiteral(call.Args[0]); id != "" {
								sloRefs = append(sloRefs, id)
							}
						}
					}
				}
			}
		}
	}

	return panelCount, queryRefs, sloRefs
}
