package slo

// BurnAlert represents an error budget burn alert configuration.
type BurnAlert struct {
	// Name is the display name of the burn alert
	Name string

	// AlertType specifies the burn rate calculation method
	AlertType AlertType

	// Threshold is the trigger threshold value
	Threshold float64

	// Window is the time window for burn rate calculation
	Window TimePeriod

	// Recipients are the notification targets for this alert
	Recipients []Recipient
}

// AlertType represents the type of burn alert calculation.
type AlertType string

const (
	// ExhaustionTime alerts when error budget will be exhausted within a time threshold
	ExhaustionTime AlertType = "exhaustion_time"

	// BudgetRate alerts when error budget consumption rate exceeds a threshold
	BudgetRate AlertType = "budget_rate"
)

// Recipient represents a notification target for burn alerts.
type Recipient struct {
	// Type is the recipient type (slack, pagerduty, email, webhook)
	Type string

	// Target is the destination (channel, service ID, email address, URL)
	Target string
}

// FastBurn creates a BurnAlert configured for fast burn detection.
// Uses BudgetRate alert type with a 1 hour window.
func FastBurn(budgetPercent float64) BurnAlert {
	return BurnAlert{
		AlertType: BudgetRate,
		Threshold: budgetPercent,
		Window:    TimePeriod{Hours: 1},
	}
}

// SlowBurn creates a BurnAlert configured for slow burn detection.
// Uses BudgetRate alert type with a 24 hour window.
func SlowBurn(budgetPercent float64) BurnAlert {
	return BurnAlert{
		AlertType: BudgetRate,
		Threshold: budgetPercent,
		Window:    TimePeriod{Hours: 24},
	}
}
