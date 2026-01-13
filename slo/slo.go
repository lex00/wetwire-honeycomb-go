// Package slo provides type-safe Honeycomb SLO declarations.
package slo

import "github.com/lex00/wetwire-honeycomb-go/query"

// SLO represents a complete Honeycomb SLO specification.
type SLO struct {
	// Name is the display name of the SLO
	Name string

	// Description provides additional context about the SLO
	Description string

	// Dataset is the Honeycomb dataset this SLO measures
	Dataset string

	// SLI defines the Service Level Indicator (good/total events)
	SLI SLI

	// Target is the SLO target percentage
	Target Target

	// TimePeriod is the rolling window for SLO calculation
	TimePeriod TimePeriod

	// BurnAlerts are alerts triggered by error budget consumption
	BurnAlerts []BurnAlert
}

// SLI represents a Service Level Indicator definition.
type SLI struct {
	// GoodEvents is the query that counts successful events
	GoodEvents query.Query

	// TotalEvents is the query that counts all events
	TotalEvents query.Query
}

// Target represents an SLO target percentage.
type Target struct {
	// Percentage is the target value (e.g., 99.9 for 99.9%)
	Percentage float64
}

// TimePeriod represents a rolling time window for SLO calculation.
type TimePeriod struct {
	// Days is the number of days in the time period
	Days int

	// Hours is the number of hours in the time period
	Hours int
}

// Percentage creates a Target with the specified percentage.
func Percentage(p float64) Target {
	return Target{Percentage: p}
}

// Days creates a TimePeriod with the specified number of days.
func Days(d int) TimePeriod {
	return TimePeriod{Days: d}
}
