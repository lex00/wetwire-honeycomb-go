// Package board provides type-safe Honeycomb board declarations.
package board

// Board represents a complete Honeycomb board specification.
type Board struct {
	// Name is the display name of the board
	Name string

	// Description provides additional context about the board
	Description string

	// Panels are the visual components of the board
	Panels []Panel

	// PresetFilters are board-level filters applied to all query panels
	PresetFilters []Filter

	// Tags are key-value metadata for organizing boards
	Tags []Tag
}

// Filter represents a board-level preset filter.
type Filter struct {
	// Column is the field to filter on
	Column string

	// Operation is the filter operator ("=", "!=", ">", ">=", "<", "<=", "contains", etc.)
	Operation string

	// Value is the value to compare against
	Value any
}

// Tag represents a key-value metadata pair for boards.
type Tag struct {
	// Key is the tag name
	Key string

	// Value is the tag value
	Value string
}
