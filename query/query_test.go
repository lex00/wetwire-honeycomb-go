package query

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryStructBasic(t *testing.T) {
	q := Query{
		Dataset:   "production",
		TimeRange: Hours(2),
		Calculations: []Calculation{
			Count(),
		},
	}

	assert.Equal(t, "production", q.Dataset)
	assert.Equal(t, 7200, q.TimeRange.TimeRange)
	assert.Len(t, q.Calculations, 1)
	assert.Equal(t, "COUNT", q.Calculations[0].Op)
}

func TestQueryWithBreakdowns(t *testing.T) {
	q := Query{
		Dataset:      "production",
		TimeRange:    Hours(2),
		Breakdowns:   []string{"service", "endpoint"},
		Calculations: []Calculation{Count()},
	}

	assert.Equal(t, []string{"service", "endpoint"}, q.Breakdowns)
}

func TestQueryWithFilters(t *testing.T) {
	q := Query{
		Dataset:   "production",
		TimeRange: Hours(2),
		Filters: []Filter{
			GT("duration_ms", 500),
			Equals("status", "200"),
		},
		Calculations: []Calculation{Count()},
	}

	assert.Len(t, q.Filters, 2)
	assert.Equal(t, ">", q.Filters[0].Op)
	assert.Equal(t, "=", q.Filters[1].Op)
}

func TestQueryWithMultipleCalculations(t *testing.T) {
	q := Query{
		Dataset:   "production",
		TimeRange: Hours(2),
		Calculations: []Calculation{
			P99("duration_ms"),
			Count(),
			Avg("duration_ms"),
		},
	}

	assert.Len(t, q.Calculations, 3)
	assert.Equal(t, "P99", q.Calculations[0].Op)
	assert.Equal(t, "COUNT", q.Calculations[1].Op)
	assert.Equal(t, "AVG", q.Calculations[2].Op)
}

func TestQueryWithAbsoluteTime(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	q := Query{
		Dataset:      "production",
		TimeRange:    Absolute(start, end),
		Calculations: []Calculation{Count()},
	}

	assert.Equal(t, 0, q.TimeRange.TimeRange)
	assert.Equal(t, int(start.Unix()), q.TimeRange.StartTime)
	assert.Equal(t, int(end.Unix()), q.TimeRange.EndTime)
}

func TestQueryWithOrderAndLimit(t *testing.T) {
	q := Query{
		Dataset:      "production",
		TimeRange:    Hours(2),
		Calculations: []Calculation{Count()},
		Orders: []Order{
			{Column: "COUNT", Op: "COUNT", Order: "descending"},
		},
		Limit: 100,
	}

	assert.Len(t, q.Orders, 1)
	assert.Equal(t, "descending", q.Orders[0].Order)
	assert.Equal(t, 100, q.Limit)
}

func TestQueryWithFilterCombination(t *testing.T) {
	q := Query{
		Dataset:   "production",
		TimeRange: Hours(2),
		Filters: []Filter{
			GT("duration_ms", 500),
			Equals("status", "200"),
		},
		FilterCombination: "OR",
		Calculations:      []Calculation{Count()},
	}

	assert.Equal(t, "OR", q.FilterCombination)
}

func TestQueryJSONSerialization(t *testing.T) {
	q := Query{
		Dataset:   "production",
		TimeRange: Hours(2),
		Breakdowns: []string{"endpoint", "service"},
		Calculations: []Calculation{
			P99("duration_ms"),
			Count(),
		},
		Filters: []Filter{
			GT("duration_ms", 500),
		},
	}

	data, err := json.Marshal(q)
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	// Check basic fields
	assert.Equal(t, "production", parsed["dataset"])

	// Check time range
	timeRange, ok := parsed["time_range"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, float64(7200), timeRange["time_range"])

	// Check breakdowns
	breakdowns, ok := parsed["breakdowns"].([]any)
	require.True(t, ok)
	assert.Len(t, breakdowns, 2)

	// Check calculations
	calculations, ok := parsed["calculations"].([]any)
	require.True(t, ok)
	assert.Len(t, calculations, 2)

	// Check filters
	filters, ok := parsed["filters"].([]any)
	require.True(t, ok)
	assert.Len(t, filters, 1)
}

func TestCompleteQueryExample(t *testing.T) {
	// This is the example from the README
	q := Query{
		Dataset:   "production",
		TimeRange: Hours(2),
		Breakdowns: []string{"endpoint", "service"},
		Calculations: []Calculation{
			P99("duration_ms"),
			Count(),
		},
		Filters: []Filter{
			GT("duration_ms", 500),
		},
	}

	assert.NotNil(t, q)
	assert.Equal(t, "production", q.Dataset)
	assert.Equal(t, 7200, q.TimeRange.TimeRange)
	assert.Len(t, q.Breakdowns, 2)
	assert.Len(t, q.Calculations, 2)
	assert.Len(t, q.Filters, 1)
}
