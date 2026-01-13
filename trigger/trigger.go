// Package trigger provides type-safe Honeycomb trigger declarations.
package trigger

import "github.com/lex00/wetwire-honeycomb-go/query"

// Trigger represents a complete Honeycomb trigger specification.
type Trigger struct {
	// Name is the display name of the trigger
	Name string

	// Description provides additional context about the trigger
	Description string

	// Dataset is the Honeycomb dataset this trigger monitors
	Dataset string

	// Query is the query that defines the metric to monitor
	Query query.Query

	// Threshold defines the condition that fires the trigger
	Threshold Threshold

	// Frequency is how often the trigger evaluates
	Frequency Frequency

	// Recipients are the notification targets when the trigger fires
	Recipients []Recipient

	// Disabled indicates whether the trigger is active
	Disabled bool
}

// Threshold represents a trigger threshold condition.
type Threshold struct {
	// Op is the comparison operator
	Op Op

	// Value is the threshold value
	Value float64
}

// Op represents a comparison operator for thresholds.
type Op string

const (
	// GT represents greater than (>)
	GT Op = ">"

	// GTE represents greater than or equal (>=)
	GTE Op = ">="

	// LT represents less than (<)
	LT Op = "<"

	// LTE represents less than or equal (<=)
	LTE Op = "<="
)

// Frequency represents how often a trigger evaluates.
type Frequency struct {
	// Seconds is the evaluation interval in seconds
	Seconds int
}

// GreaterThan creates a Threshold with the > operator.
func GreaterThan(value float64) Threshold {
	return Threshold{Op: GT, Value: value}
}

// GreaterThanOrEqual creates a Threshold with the >= operator.
func GreaterThanOrEqual(value float64) Threshold {
	return Threshold{Op: GTE, Value: value}
}

// LessThan creates a Threshold with the < operator.
func LessThan(value float64) Threshold {
	return Threshold{Op: LT, Value: value}
}

// LessThanOrEqual creates a Threshold with the <= operator.
func LessThanOrEqual(value float64) Threshold {
	return Threshold{Op: LTE, Value: value}
}

// Minutes creates a Frequency with the specified number of minutes.
func Minutes(m int) Frequency {
	return Frequency{Seconds: m * 60}
}

// Seconds creates a Frequency with the specified number of seconds.
func Seconds(s int) Frequency {
	return Frequency{Seconds: s}
}
