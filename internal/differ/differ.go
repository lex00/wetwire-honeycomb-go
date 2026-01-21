// Package differ provides semantic comparison of Honeycomb query configurations.
package differ

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

// HoneycombDiffer implements semantic comparison for Honeycomb Query JSON.
type HoneycombDiffer struct{}

// Compile-time interface check
var _ coredomain.Differ = (*HoneycombDiffer)(nil)

// New creates a new HoneycombDiffer.
func New() *HoneycombDiffer {
	return &HoneycombDiffer{}
}

// Diff compares two Honeycomb query configuration files or directories.
func (d *HoneycombDiffer) Diff(ctx *coredomain.Context, file1, file2 string, opts coredomain.DiffOpts) (*coredomain.DiffResult, error) {
	// Load configurations
	config1, err := loadConfig(file1)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", file1, err)
	}

	config2, err := loadConfig(file2)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", file2, err)
	}

	// Compare configurations
	return compare(config1, config2, opts)
}

// HoneycombConfig represents the structure of Honeycomb configuration output.
type HoneycombConfig struct {
	Queries  map[string]json.RawMessage `json:"queries,omitempty"`
	Boards   map[string]json.RawMessage `json:"boards,omitempty"`
	SLOs     map[string]json.RawMessage `json:"slos,omitempty"`
	Triggers map[string]json.RawMessage `json:"triggers,omitempty"`
}

// loadConfig loads a Honeycomb configuration from a file.
func loadConfig(path string) (*HoneycombConfig, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return loadConfigFromDir(path)
	}

	return loadConfigFromFile(path)
}

// loadConfigFromFile loads configuration from a single JSON file.
func loadConfigFromFile(path string) (*HoneycombConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config HoneycombConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}

	return &config, nil
}

// loadConfigFromDir loads configuration from a directory of JSON files.
func loadConfigFromDir(dir string) (*HoneycombConfig, error) {
	config := &HoneycombConfig{
		Queries:  make(map[string]json.RawMessage),
		Boards:   make(map[string]json.RawMessage),
		SLOs:     make(map[string]json.RawMessage),
		Triggers: make(map[string]json.RawMessage),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var fileConfig HoneycombConfig
		if err := json.Unmarshal(data, &fileConfig); err != nil {
			// Try to parse as a single resource type
			return nil
		}

		// Merge into main config
		for k, v := range fileConfig.Queries {
			config.Queries[k] = v
		}
		for k, v := range fileConfig.Boards {
			config.Boards[k] = v
		}
		for k, v := range fileConfig.SLOs {
			config.SLOs[k] = v
		}
		for k, v := range fileConfig.Triggers {
			config.Triggers[k] = v
		}

		return nil
	})

	return config, err
}

// compare compares two Honeycomb configurations and returns the differences.
func compare(config1, config2 *HoneycombConfig, opts coredomain.DiffOpts) (*coredomain.DiffResult, error) {
	result := &coredomain.DiffResult{
		Entries: []coredomain.DiffEntry{},
		Summary: coredomain.DiffSummary{},
	}

	// Compare queries
	compareResourceMap(config1.Queries, config2.Queries, "query", result, opts)

	// Compare boards
	compareResourceMap(config1.Boards, config2.Boards, "board", result, opts)

	// Compare SLOs
	compareResourceMap(config1.SLOs, config2.SLOs, "slo", result, opts)

	// Compare triggers
	compareResourceMap(config1.Triggers, config2.Triggers, "trigger", result, opts)

	// Calculate total
	result.Summary.Total = result.Summary.Added + result.Summary.Removed + result.Summary.Modified

	return result, nil
}

// compareResourceMap compares two maps of resources.
func compareResourceMap(map1, map2 map[string]json.RawMessage, resourceType string, result *coredomain.DiffResult, opts coredomain.DiffOpts) {
	// Get all keys from both maps
	allKeys := make(map[string]bool)
	for k := range map1 {
		allKeys[k] = true
	}
	for k := range map2 {
		allKeys[k] = true
	}

	// Sort keys for deterministic output
	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		raw1, exists1 := map1[name]
		raw2, exists2 := map2[name]

		if !exists1 && exists2 {
			// Added
			result.Entries = append(result.Entries, coredomain.DiffEntry{
				Resource: name,
				Type:     resourceType,
				Action:   "added",
			})
			result.Summary.Added++
		} else if exists1 && !exists2 {
			// Removed
			result.Entries = append(result.Entries, coredomain.DiffEntry{
				Resource: name,
				Type:     resourceType,
				Action:   "removed",
			})
			result.Summary.Removed++
		} else {
			// Both exist, compare content
			changes := compareJSON(raw1, raw2, opts)
			if len(changes) > 0 {
				result.Entries = append(result.Entries, coredomain.DiffEntry{
					Resource: name,
					Type:     resourceType,
					Action:   "modified",
					Changes:  changes,
				})
				result.Summary.Modified++
			}
		}
	}
}

// compareJSON compares two JSON values and returns a list of changes.
func compareJSON(raw1, raw2 json.RawMessage, opts coredomain.DiffOpts) []string {
	var val1, val2 interface{}
	if err := json.Unmarshal(raw1, &val1); err != nil {
		return []string{"parse error in first value"}
	}
	if err := json.Unmarshal(raw2, &val2); err != nil {
		return []string{"parse error in second value"}
	}

	return compareValues(val1, val2, "", opts)
}

// compareValues recursively compares two values.
func compareValues(v1, v2 interface{}, path string, opts coredomain.DiffOpts) []string {
	var changes []string

	// Handle nil cases
	if v1 == nil && v2 == nil {
		return nil
	}
	if v1 == nil {
		return []string{formatChange(path, "nil", v2)}
	}
	if v2 == nil {
		return []string{formatChange(path, v1, "nil")}
	}

	// Check types
	t1 := reflect.TypeOf(v1)
	t2 := reflect.TypeOf(v2)
	if t1 != t2 {
		return []string{formatChange(path, v1, v2)}
	}

	switch val1 := v1.(type) {
	case map[string]interface{}:
		val2 := v2.(map[string]interface{})
		changes = append(changes, compareMaps(val1, val2, path, opts)...)

	case []interface{}:
		val2 := v2.([]interface{})
		changes = append(changes, compareSlices(val1, val2, path, opts)...)

	default:
		if !reflect.DeepEqual(v1, v2) {
			changes = append(changes, formatChange(path, v1, v2))
		}
	}

	return changes
}

// compareMaps compares two maps.
func compareMaps(m1, m2 map[string]interface{}, path string, opts coredomain.DiffOpts) []string {
	var changes []string

	// Get all keys
	allKeys := make(map[string]bool)
	for k := range m1 {
		allKeys[k] = true
	}
	for k := range m2 {
		allKeys[k] = true
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		keyPath := k
		if path != "" {
			keyPath = path + "." + k
		}

		val1, exists1 := m1[k]
		val2, exists2 := m2[k]

		if !exists1 {
			changes = append(changes, fmt.Sprintf("%s: added", keyPath))
		} else if !exists2 {
			changes = append(changes, fmt.Sprintf("%s: removed", keyPath))
		} else {
			changes = append(changes, compareValues(val1, val2, keyPath, opts)...)
		}
	}

	return changes
}

// compareSlices compares two slices.
func compareSlices(s1, s2 []interface{}, path string, opts coredomain.DiffOpts) []string {
	var changes []string

	if opts.IgnoreOrder {
		// Compare as sets - order doesn't matter
		if !slicesEqualIgnoreOrder(s1, s2) {
			changes = append(changes, fmt.Sprintf("%s: array contents differ", path))
		}
	} else {
		// Compare element by element
		maxLen := len(s1)
		if len(s2) > maxLen {
			maxLen = len(s2)
		}

		for i := 0; i < maxLen; i++ {
			elemPath := fmt.Sprintf("%s[%d]", path, i)
			if i >= len(s1) {
				changes = append(changes, fmt.Sprintf("%s: added", elemPath))
			} else if i >= len(s2) {
				changes = append(changes, fmt.Sprintf("%s: removed", elemPath))
			} else {
				changes = append(changes, compareValues(s1[i], s2[i], elemPath, opts)...)
			}
		}
	}

	return changes
}

// slicesEqualIgnoreOrder checks if two slices have the same elements regardless of order.
func slicesEqualIgnoreOrder(s1, s2 []interface{}) bool {
	if len(s1) != len(s2) {
		return false
	}

	// Simple approach: convert to JSON strings and compare as sets
	set1 := make(map[string]int)
	set2 := make(map[string]int)

	for _, v := range s1 {
		data, _ := json.Marshal(v)
		set1[string(data)]++
	}
	for _, v := range s2 {
		data, _ := json.Marshal(v)
		set2[string(data)]++
	}

	return reflect.DeepEqual(set1, set2)
}

// formatChange formats a change message.
func formatChange(path string, old, new interface{}) string {
	oldStr := formatValue(old)
	newStr := formatValue(new)
	if path == "" {
		return fmt.Sprintf("%s -> %s", oldStr, newStr)
	}
	return fmt.Sprintf("%s: %s -> %s", path, oldStr, newStr)
}

// formatValue formats a value for display.
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case nil:
		return "nil"
	default:
		data, _ := json.Marshal(v)
		s := string(data)
		if len(s) > 50 {
			return s[:47] + "..."
		}
		return s
	}
}
