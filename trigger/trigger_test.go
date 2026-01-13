package trigger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/query"
)

func TestTrigger_BasicFields(t *testing.T) {
	tr := Trigger{
		Name:        "High Latency Alert",
		Description: "Alert when P99 exceeds 500ms",
		Dataset:     "production",
		Disabled:    false,
	}

	assert.Equal(t, "High Latency Alert", tr.Name)
	assert.Equal(t, "Alert when P99 exceeds 500ms", tr.Description)
	assert.Equal(t, "production", tr.Dataset)
	assert.False(t, tr.Disabled)
}

func TestTrigger_WithQuery(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
		},
	}

	tr := Trigger{
		Name:    "High Latency Alert",
		Dataset: "production",
		Query:   q,
	}

	assert.Equal(t, "production", tr.Query.Dataset)
	require.Len(t, tr.Query.Calculations, 1)
}

func TestTrigger_WithThreshold(t *testing.T) {
	tr := Trigger{
		Name:      "High Latency Alert",
		Dataset:   "production",
		Threshold: GreaterThan(500),
	}

	assert.Equal(t, GT, tr.Threshold.Op)
	assert.Equal(t, 500.0, tr.Threshold.Value)
}

func TestTrigger_WithFrequency(t *testing.T) {
	tr := Trigger{
		Name:      "High Latency Alert",
		Dataset:   "production",
		Frequency: Minutes(5),
	}

	assert.Equal(t, 300, tr.Frequency.Seconds)
}

func TestOp_Constants(t *testing.T) {
	assert.Equal(t, Op(">"), GT)
	assert.Equal(t, Op(">="), GTE)
	assert.Equal(t, Op("<"), LT)
	assert.Equal(t, Op("<="), LTE)
}

func TestGreaterThan_Builder(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"500", 500, 500},
		{"100.5", 100.5, 100.5},
		{"0", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := GreaterThan(tt.value)
			assert.Equal(t, GT, th.Op)
			assert.Equal(t, tt.expected, th.Value)
		})
	}
}

func TestGreaterThanOrEqual_Builder(t *testing.T) {
	th := GreaterThanOrEqual(100)
	assert.Equal(t, GTE, th.Op)
	assert.Equal(t, 100.0, th.Value)
}

func TestLessThan_Builder(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"10", 10, 10},
		{"0.5", 0.5, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := LessThan(tt.value)
			assert.Equal(t, LT, th.Op)
			assert.Equal(t, tt.expected, th.Value)
		})
	}
}

func TestLessThanOrEqual_Builder(t *testing.T) {
	th := LessThanOrEqual(50)
	assert.Equal(t, LTE, th.Op)
	assert.Equal(t, 50.0, th.Value)
}

func TestMinutes_Builder(t *testing.T) {
	tests := []struct {
		name     string
		minutes  int
		expected int
	}{
		{"1 minute", 1, 60},
		{"5 minutes", 5, 300},
		{"10 minutes", 10, 600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			freq := Minutes(tt.minutes)
			assert.Equal(t, tt.expected, freq.Seconds)
		})
	}
}

func TestSeconds_Builder(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected int
	}{
		{"30 seconds", 30, 30},
		{"60 seconds", 60, 60},
		{"120 seconds", 120, 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			freq := Seconds(tt.seconds)
			assert.Equal(t, tt.expected, freq.Seconds)
		})
	}
}

func TestTrigger_Complete(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
		},
		Filters: []query.Filter{
			query.GT("duration_ms", 100),
		},
	}

	tr := Trigger{
		Name:        "High P99 Latency",
		Description: "Alert when P99 exceeds 500ms",
		Dataset:     "production",
		Query:       q,
		Threshold:   GreaterThan(500),
		Frequency:   Minutes(5),
		Recipients: []Recipient{
			SlackChannel("#alerts"),
			PagerDutyService("api-team"),
		},
		Disabled: false,
	}

	assert.Equal(t, "High P99 Latency", tr.Name)
	assert.Equal(t, "production", tr.Dataset)
	assert.Equal(t, GT, tr.Threshold.Op)
	assert.Equal(t, 500.0, tr.Threshold.Value)
	assert.Equal(t, 300, tr.Frequency.Seconds)
	require.Len(t, tr.Recipients, 2)
	assert.False(t, tr.Disabled)
}

func TestTrigger_Disabled(t *testing.T) {
	tr := Trigger{
		Name:     "Disabled Alert",
		Dataset:  "production",
		Disabled: true,
	}

	assert.True(t, tr.Disabled)
}
