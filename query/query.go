// Package query provides type-safe Honeycomb query declarations.
package query

// Query represents a complete Honeycomb query specification.
type Query struct {
	// Dataset is the name of the Honeycomb dataset to query
	Dataset string `json:"dataset"`

	// TimeRange specifies the time window for the query
	TimeRange TimeRange `json:"time_range"`

	// Calculations are the aggregation operations to perform
	Calculations []Calculation `json:"calculations"`

	// Filters restrict which events are included in the query
	Filters []Filter `json:"filters,omitempty"`

	// FilterCombination specifies how multiple filters are combined ("AND" or "OR")
	// Defaults to "AND" if not specified
	FilterCombination string `json:"filter_combination,omitempty"`

	// Breakdowns are the fields to group results by
	Breakdowns []string `json:"breakdowns,omitempty"`

	// Orders specify how to sort the results
	Orders []Order `json:"orders,omitempty"`

	// Limit restricts the number of results returned
	Limit int `json:"limit,omitempty"`

	// Granularity is the time bucket size for time series queries (in seconds)
	Granularity int `json:"granularity,omitempty"`
}

// Order specifies how query results should be sorted.
type Order struct {
	// Column is the field to sort by
	Column string `json:"column,omitempty"`

	// Op is the calculation operation if sorting by a calculated value
	Op string `json:"op,omitempty"`

	// Order is "ascending" or "descending"
	Order string `json:"order"`
}
