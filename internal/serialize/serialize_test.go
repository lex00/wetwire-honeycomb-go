package serialize

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/query"
)

func TestToJSON_BasicQuery(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(2),
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(7200), result["time_range"])
	assert.NotNil(t, result["calculations"])
}

func TestToJSON_WithBreakdowns(t *testing.T) {
	q := query.Query{
		Dataset:    "production",
		TimeRange:  query.Hours(1),
		Breakdowns: []string{"endpoint", "service"},
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
		},
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	breakdowns := result["breakdowns"].([]any)
	assert.Equal(t, "endpoint", breakdowns[0])
	assert.Equal(t, "service", breakdowns[1])
}

func TestToJSON_WithFilters(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
		Calculations: []query.Calculation{
			query.Count(),
		},
		Filters: []query.Filter{
			query.GT("duration_ms", 500),
			query.Equals("status", "error"),
		},
		FilterCombination: "AND",
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	filters := result["filters"].([]any)
	assert.Len(t, filters, 2)

	filter1 := filters[0].(map[string]any)
	assert.Equal(t, "duration_ms", filter1["column"])
	assert.Equal(t, ">", filter1["op"])
	assert.Equal(t, float64(500), filter1["value"])

	assert.Equal(t, "AND", result["filter_combination"])
}

func TestToJSON_WithOrders(t *testing.T) {
	q := query.Query{
		Dataset:    "production",
		TimeRange:  query.Hours(1),
		Breakdowns: []string{"endpoint"},
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
		},
		Orders: []query.Order{
			{Op: "P99", Column: "duration_ms", Order: "descending"},
		},
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	orders := result["orders"].([]any)
	assert.Len(t, orders, 1)

	order := orders[0].(map[string]any)
	assert.Equal(t, "P99", order["op"])
	assert.Equal(t, "duration_ms", order["column"])
	assert.Equal(t, "descending", order["order"])
}

func TestToJSON_OmitsEmptyFields(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// These fields should be omitted when empty/zero
	_, hasFilters := result["filters"]
	assert.False(t, hasFilters, "empty filters should be omitted")

	_, hasBreakdowns := result["breakdowns"]
	assert.False(t, hasBreakdowns, "empty breakdowns should be omitted")

	_, hasLimit := result["limit"]
	assert.False(t, hasLimit, "zero limit should be omitted")
}

func TestToJSON_AllCalculationTypes(t *testing.T) {
	calculations := []query.Calculation{
		query.Count(),
		query.Sum("bytes"),
		query.Avg("duration_ms"),
		query.Min("duration_ms"),
		query.Max("duration_ms"),
		query.P50("duration_ms"),
		query.P75("duration_ms"),
		query.P90("duration_ms"),
		query.P95("duration_ms"),
		query.P99("duration_ms"),
		query.Heatmap("duration_ms"),
	}

	q := query.Query{
		Dataset:      "production",
		TimeRange:   query.Hours(1),
		Calculations: calculations,
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	calcs := result["calculations"].([]any)
	assert.Len(t, calcs, len(calculations))
}

func TestToJSONPretty(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(2),
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	data, err := ToJSONPretty(q)
	require.NoError(t, err)

	// Should be indented
	assert.Contains(t, string(data), "\n")
	assert.Contains(t, string(data), "  ")

	// Should still be valid JSON
	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)
}

func TestToJSON_AbsoluteTimeRange(t *testing.T) {
	start := time.Unix(1704067200, 0)
	end := time.Unix(1704153600, 0)

	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Absolute(start, end),
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(1704067200), result["start_time"])
	assert.Equal(t, float64(1704153600), result["end_time"])
}

func TestToJSON_ComplexQuery(t *testing.T) {
	q := query.Query{
		Dataset:           "production",
		TimeRange:         query.Hours(24),
		Breakdowns:        []string{"service", "environment"},
		FilterCombination: "AND",
		Calculations: []query.Calculation{
			query.Count(),
			query.P99("duration_ms"),
			query.Avg("response_size"),
		},
		Filters: []query.Filter{
			query.GT("duration_ms", 1000),
			query.Equals("environment", "production"),
		},
		Orders: []query.Order{
			{Op: "P99", Column: "duration_ms", Order: "descending"},
		},
		Limit:       100,
		Granularity: 60,
	}

	data, err := ToJSON(q)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(86400), result["time_range"])
	assert.Equal(t, float64(100), result["limit"])
	assert.Equal(t, float64(60), result["granularity"])
	assert.Equal(t, "AND", result["filter_combination"])
}
