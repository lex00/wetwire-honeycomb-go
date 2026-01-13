package slo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBurnAlert_BasicFields(t *testing.T) {
	alert := BurnAlert{
		Name:      "Fast Burn Alert",
		AlertType: ExhaustionTime,
		Threshold: 2.0,
		Window:    Days(1),
	}

	assert.Equal(t, "Fast Burn Alert", alert.Name)
	assert.Equal(t, ExhaustionTime, alert.AlertType)
	assert.Equal(t, 2.0, alert.Threshold)
	assert.Equal(t, 1, alert.Window.Days)
}

func TestAlertType_Constants(t *testing.T) {
	assert.Equal(t, AlertType("exhaustion_time"), ExhaustionTime)
	assert.Equal(t, AlertType("budget_rate"), BudgetRate)
}

func TestFastBurn_Builder(t *testing.T) {
	alert := FastBurn(2.0)

	assert.Equal(t, BudgetRate, alert.AlertType)
	assert.Equal(t, 2.0, alert.Threshold)
	// FastBurn uses 1 hour window by default
	assert.Equal(t, 1, alert.Window.Hours)
}

func TestSlowBurn_Builder(t *testing.T) {
	alert := SlowBurn(5.0)

	assert.Equal(t, BudgetRate, alert.AlertType)
	assert.Equal(t, 5.0, alert.Threshold)
	// SlowBurn uses 24 hour window by default
	assert.Equal(t, 24, alert.Window.Hours)
}

func TestFastBurn_WithDifferentThresholds(t *testing.T) {
	tests := []struct {
		name      string
		threshold float64
	}{
		{"1%", 1.0},
		{"2%", 2.0},
		{"5%", 5.0},
		{"10%", 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert := FastBurn(tt.threshold)
			assert.Equal(t, tt.threshold, alert.Threshold)
		})
	}
}

func TestSlowBurn_WithDifferentThresholds(t *testing.T) {
	tests := []struct {
		name      string
		threshold float64
	}{
		{"5%", 5.0},
		{"10%", 10.0},
		{"15%", 15.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alert := SlowBurn(tt.threshold)
			assert.Equal(t, tt.threshold, alert.Threshold)
		})
	}
}

func TestBurnAlert_WithRecipients(t *testing.T) {
	alert := BurnAlert{
		Name:      "Critical Alert",
		AlertType: ExhaustionTime,
		Threshold: 1.0,
		Recipients: []Recipient{
			{Type: "slack", Target: "#alerts"},
			{Type: "pagerduty", Target: "api-team"},
		},
	}

	assert.Len(t, alert.Recipients, 2)
	assert.Equal(t, "slack", alert.Recipients[0].Type)
	assert.Equal(t, "#alerts", alert.Recipients[0].Target)
}

func TestRecipient_Types(t *testing.T) {
	recipients := []Recipient{
		{Type: "slack", Target: "#alerts"},
		{Type: "pagerduty", Target: "service-123"},
		{Type: "email", Target: "team@example.com"},
		{Type: "webhook", Target: "https://example.com/webhook"},
	}

	expectedTypes := []string{"slack", "pagerduty", "email", "webhook"}
	for i, r := range recipients {
		assert.Equal(t, expectedTypes[i], r.Type)
	}
}
