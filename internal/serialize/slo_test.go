package serialize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

func TestSLOToJSON_BasicFields(t *testing.T) {
	s := slo.SLO{
		Name:        "API Availability",
		Description: "99.9% of requests succeed",
		Dataset:     "production",
	}

	data, err := SLOToJSON(s)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "API Availability", result["name"])
	assert.Equal(t, "99.9% of requests succeed", result["description"])
	assert.Equal(t, "production", result["dataset"])
}

func TestSLOToJSON_WithTarget(t *testing.T) {
	s := slo.SLO{
		Name:   "API Availability",
		Target: slo.Percentage(99.9),
	}

	data, err := SLOToJSON(s)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Honeycomb API uses target_per_million
	assert.Equal(t, float64(999000), result["target_per_million"])
}

func TestSLOToJSON_WithTimePeriod(t *testing.T) {
	s := slo.SLO{
		Name:       "API Availability",
		TimePeriod: slo.Days(30),
	}

	data, err := SLOToJSON(s)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(30), result["time_period_days"])
}

func TestSLOToJSON_WithSLI(t *testing.T) {
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

	s := slo.SLO{
		Name:    "API Availability",
		Dataset: "production",
		SLI: slo.SLI{
			GoodEvents:  goodEvents,
			TotalEvents: totalEvents,
		},
	}

	data, err := SLOToJSON(s)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	sli, ok := result["sli"].(map[string]interface{})
	require.True(t, ok)

	// Check good events query is serialized
	goodEventsQuery, ok := sli["good_events"].(map[string]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, goodEventsQuery)

	// Check total events query is serialized
	totalEventsQuery, ok := sli["total_events"].(map[string]interface{})
	require.True(t, ok)
	assert.NotEmpty(t, totalEventsQuery)
}

func TestSLOToJSON_WithBurnAlerts(t *testing.T) {
	s := slo.SLO{
		Name:    "API Availability",
		Dataset: "production",
		BurnAlerts: []slo.BurnAlert{
			slo.FastBurn(2.0),
			slo.SlowBurn(5.0),
		},
	}

	data, err := SLOToJSON(s)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	alerts, ok := result["burn_alerts"].([]interface{})
	require.True(t, ok)
	require.Len(t, alerts, 2)

	fastBurn := alerts[0].(map[string]interface{})
	assert.Equal(t, "budget_rate", fastBurn["alert_type"])
	assert.Equal(t, float64(2.0), fastBurn["threshold"])
}

func TestSLOToJSON_Complete(t *testing.T) {
	goodEvents := query.Query{
		Dataset: "production",
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	totalEvents := query.Query{
		Dataset: "production",
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	s := slo.SLO{
		Name:        "API Availability",
		Description: "99.9% of requests succeed",
		Dataset:     "production",
		SLI: slo.SLI{
			GoodEvents:  goodEvents,
			TotalEvents: totalEvents,
		},
		Target:     slo.Percentage(99.9),
		TimePeriod: slo.Days(30),
		BurnAlerts: []slo.BurnAlert{
			slo.FastBurn(2.0),
		},
	}

	data, err := SLOToJSON(s)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "API Availability", result["name"])
	assert.Equal(t, "production", result["dataset"])
	assert.Equal(t, float64(999000), result["target_per_million"])
	assert.Equal(t, float64(30), result["time_period_days"])
}

func TestSLOToJSONPretty(t *testing.T) {
	s := slo.SLO{
		Name: "Test SLO",
	}

	data, err := SLOToJSONPretty(s)
	require.NoError(t, err)

	assert.Contains(t, string(data), "\n")
	assert.Contains(t, string(data), "  ")
}
