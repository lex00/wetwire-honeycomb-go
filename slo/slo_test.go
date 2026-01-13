package slo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/query"
)

func TestSLO_BasicFields(t *testing.T) {
	s := SLO{
		Name:        "API Availability",
		Description: "99.9% of requests succeed",
		Dataset:     "production",
	}

	assert.Equal(t, "API Availability", s.Name)
	assert.Equal(t, "99.9% of requests succeed", s.Description)
	assert.Equal(t, "production", s.Dataset)
}

func TestSLO_WithSLI(t *testing.T) {
	goodEvents := query.Query{
		Dataset:   "production",
		TimeRange: query.Days(30),
		Calculations: []query.Calculation{
			query.Count(),
		},
		Filters: []query.Filter{
			query.LT("http.status_code", 500),
		},
	}

	totalEvents := query.Query{
		Dataset:   "production",
		TimeRange: query.Days(30),
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	s := SLO{
		Name:    "API Availability",
		Dataset: "production",
		SLI: SLI{
			GoodEvents:  goodEvents,
			TotalEvents: totalEvents,
		},
	}

	assert.Equal(t, "production", s.SLI.GoodEvents.Dataset)
	assert.Equal(t, "production", s.SLI.TotalEvents.Dataset)
	require.Len(t, s.SLI.GoodEvents.Filters, 1)
}

func TestSLO_WithTargetAndTimePeriod(t *testing.T) {
	s := SLO{
		Name:       "API Availability",
		Dataset:    "production",
		Target:     Percentage(99.9),
		TimePeriod: Days(30),
	}

	assert.Equal(t, 99.9, s.Target.Percentage)
	assert.Equal(t, 30, s.TimePeriod.Days)
}

func TestPercentage_Builder(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"99.9%", 99.9, 99.9},
		{"99.99%", 99.99, 99.99},
		{"95%", 95.0, 95.0},
		{"100%", 100.0, 100.0},
		{"0%", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := Percentage(tt.input)
			assert.Equal(t, tt.expected, target.Percentage)
		})
	}
}

func TestDays_Builder(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"7 days", 7, 7},
		{"14 days", 14, 14},
		{"30 days", 30, 30},
		{"90 days", 90, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			period := Days(tt.input)
			assert.Equal(t, tt.expected, period.Days)
		})
	}
}

func TestSLO_WithBurnAlerts(t *testing.T) {
	s := SLO{
		Name:    "API Availability",
		Dataset: "production",
		Target:  Percentage(99.9),
		BurnAlerts: []BurnAlert{
			FastBurn(2.0),
			SlowBurn(5.0),
		},
	}

	require.Len(t, s.BurnAlerts, 2)
}

func TestSLO_Complete(t *testing.T) {
	goodEvents := query.Query{
		Dataset:   "production",
		TimeRange: query.Days(30),
		Calculations: []query.Calculation{
			query.Count(),
		},
		Filters: []query.Filter{
			query.LT("http.status_code", 500),
		},
	}

	totalEvents := query.Query{
		Dataset:   "production",
		TimeRange: query.Days(30),
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	s := SLO{
		Name:        "API Availability",
		Description: "99.9% of requests succeed",
		Dataset:     "production",
		SLI: SLI{
			GoodEvents:  goodEvents,
			TotalEvents: totalEvents,
		},
		Target:     Percentage(99.9),
		TimePeriod: Days(30),
		BurnAlerts: []BurnAlert{
			FastBurn(2.0),
			SlowBurn(5.0),
		},
	}

	assert.Equal(t, "API Availability", s.Name)
	assert.Equal(t, 99.9, s.Target.Percentage)
	assert.Equal(t, 30, s.TimePeriod.Days)
	require.Len(t, s.BurnAlerts, 2)
}
