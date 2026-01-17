// Command diff provides semantic comparison of Honeycomb Query JSON files.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/lex00/wetwire-honeycomb-go/internal/builder"
	"github.com/lex00/wetwire-honeycomb-go/internal/serialize"
	"github.com/spf13/cobra"
)

func newDiffCmd() *cobra.Command {
	var outputFile string
	var semantic bool
	var verbose bool

	cmd := &cobra.Command{
		Use:   "diff [packages]",
		Short: "Compare generated output vs existing config",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if outputFile == "" {
				return fmt.Errorf("--output flag is required")
			}

			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			// Build queries
			b, err := builder.NewBuilder(path)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}

			result, err := b.Build()
			if err != nil {
				return fmt.Errorf("build failed: %w", err)
			}

			if result.QueryCount() == 0 {
				return fmt.Errorf("no queries found")
			}

			// Generate current output
			queries := result.Queries()
			var currentJSON []byte

			if len(queries) == 1 {
				q := discoveredToQuery(queries[0])
				currentJSON, err = serialize.ToJSONPretty(q)
			} else {
				queryMap := make(map[string]json.RawMessage)
				for _, dq := range queries {
					q := discoveredToQuery(dq)
					data, e := serialize.ToJSON(q)
					if e != nil {
						err = e
						break
					}
					queryMap[dq.Name] = data
				}
				if err == nil {
					currentJSON, err = json.MarshalIndent(queryMap, "", "  ")
				}
			}

			if err != nil {
				return fmt.Errorf("serialization failed: %w", err)
			}

			// Read existing file
			existingJSON, err := os.ReadFile(outputFile)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", outputFile, err)
			}

			// Compare
			if semantic {
				return semanticDiff(currentJSON, existingJSON, verbose)
			}
			return textDiff(currentJSON, existingJSON, outputFile, verbose)
		},
	}

	cmd.Flags().StringVar(&outputFile, "output", "", "JSON file to compare against")
	cmd.Flags().BoolVar(&semantic, "semantic", false, "Compare semantic structure instead of text")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return cmd
}

func textDiff(current, existing []byte, filename string, verbose bool) error {
	// Normalize line endings
	current = bytes.ReplaceAll(current, []byte("\r\n"), []byte("\n"))
	existing = bytes.ReplaceAll(existing, []byte("\r\n"), []byte("\n"))

	if bytes.Equal(current, existing) {
		if verbose {
			fmt.Println("Files are identical")
		}
		return nil
	}

	// Show line-by-line diff
	currentLines := strings.Split(string(current), "\n")
	existingLines := strings.Split(string(existing), "\n")

	fmt.Printf("--- %s (existing)\n", filename)
	fmt.Println("+++ generated")

	maxLen := len(currentLines)
	if len(existingLines) > maxLen {
		maxLen = len(existingLines)
	}

	for i := 0; i < maxLen; i++ {
		var currLine, existLine string
		if i < len(currentLines) {
			currLine = currentLines[i]
		}
		if i < len(existingLines) {
			existLine = existingLines[i]
		}

		if currLine != existLine {
			if existLine != "" {
				fmt.Printf("-%s\n", existLine)
			}
			if currLine != "" {
				fmt.Printf("+%s\n", currLine)
			}
		}
	}

	return fmt.Errorf("files differ")
}

func semanticDiff(current, existing []byte, verbose bool) error {
	var currData, existData interface{}

	if err := json.Unmarshal(current, &currData); err != nil {
		return fmt.Errorf("error parsing generated JSON: %w", err)
	}

	if err := json.Unmarshal(existing, &existData); err != nil {
		return fmt.Errorf("error parsing existing JSON: %w", err)
	}

	if reflect.DeepEqual(currData, existData) {
		if verbose {
			fmt.Println("Semantically identical")
		}
		return nil
	}

	// Show structural differences
	diffs := compareJSON(currData, existData, "")
	for _, d := range diffs {
		fmt.Println(d)
	}

	return fmt.Errorf("semantic differences found")
}

func compareJSON(a, b interface{}, path string) []string {
	var diffs []string

	switch aTyped := a.(type) {
	case map[string]interface{}:
		bTyped, ok := b.(map[string]interface{})
		if !ok {
			return []string{fmt.Sprintf("Type mismatch at %s: map vs %T", path, b)}
		}

		// Check keys in a
		for k, av := range aTyped {
			newPath := path + "." + k
			if path == "" {
				newPath = k
			}
			if bv, ok := bTyped[k]; ok {
				diffs = append(diffs, compareJSON(av, bv, newPath)...)
			} else {
				diffs = append(diffs, fmt.Sprintf("Key missing in existing: %s", newPath))
			}
		}

		// Check keys in b not in a
		for k := range bTyped {
			newPath := path + "." + k
			if path == "" {
				newPath = k
			}
			if _, ok := aTyped[k]; !ok {
				diffs = append(diffs, fmt.Sprintf("Extra key in existing: %s", newPath))
			}
		}

	case []interface{}:
		bTyped, ok := b.([]interface{})
		if !ok {
			return []string{fmt.Sprintf("Type mismatch at %s: array vs %T", path, b)}
		}

		if len(aTyped) != len(bTyped) {
			diffs = append(diffs, fmt.Sprintf("Array length mismatch at %s: %d vs %d", path, len(aTyped), len(bTyped)))
		}

		minLen := len(aTyped)
		if len(bTyped) < minLen {
			minLen = len(bTyped)
		}

		for i := 0; i < minLen; i++ {
			diffs = append(diffs, compareJSON(aTyped[i], bTyped[i], fmt.Sprintf("%s[%d]", path, i))...)
		}

	default:
		if !reflect.DeepEqual(a, b) {
			diffs = append(diffs, fmt.Sprintf("Value mismatch at %s: %v vs %v", path, a, b))
		}
	}

	return diffs
}
