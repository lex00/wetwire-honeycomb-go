package serialize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

func TestTriggerToJSON_BasicFields(t *testing.T) {
	tr := trigger.Trigger{
		Name:        "High Latency Alert",
		Description: "Alert when P99 exceeds 500ms",
		Dataset:     "production",
		Disabled:    false,
	}

	data, err := TriggerToJSON(tr)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "High Latency Alert", result["name"])
	assert.Equal(t, "Alert when P99 exceeds 500ms", result["description"])
	assert.Equal(t, "production", result["dataset"])
	assert.Equal(t, false, result["disabled"])
}

func TestTriggerToJSON_WithQuery(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
		},
	}

	tr := trigger.Trigger{
		Name:    "High Latency",
		Dataset: "production",
		Query:   q,
	}

	data, err := TriggerToJSON(tr)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	querySpec, ok := result["query"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, float64(3600), querySpec["time_range"])
}

func TestTriggerToJSON_WithThreshold(t *testing.T) {
	tr := trigger.Trigger{
		Name:      "High Latency",
		Dataset:   "production",
		Threshold: trigger.GreaterThan(500),
	}

	data, err := TriggerToJSON(tr)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	threshold, ok := result["threshold"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, ">", threshold["op"])
	assert.Equal(t, float64(500), threshold["value"])
}

func TestTriggerToJSON_WithFrequency(t *testing.T) {
	tr := trigger.Trigger{
		Name:      "High Latency",
		Dataset:   "production",
		Frequency: trigger.Minutes(5),
	}

	data, err := TriggerToJSON(tr)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, float64(300), result["frequency"])
}

func TestTriggerToJSON_WithRecipients(t *testing.T) {
	tr := trigger.Trigger{
		Name:    "High Latency",
		Dataset: "production",
		Recipients: []trigger.Recipient{
			trigger.SlackChannel("#alerts"),
			trigger.PagerDutyService("api-team"),
			trigger.EmailAddress("team@example.com"),
		},
	}

	data, err := TriggerToJSON(tr)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	recipients, ok := result["recipients"].([]interface{})
	require.True(t, ok)
	require.Len(t, recipients, 3)

	slack := recipients[0].(map[string]interface{})
	assert.Equal(t, "slack", slack["type"])
	assert.Equal(t, "#alerts", slack["target"])
}

func TestTriggerToJSON_Complete(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
		Calculations: []query.Calculation{
			query.P99("duration_ms"),
		},
	}

	tr := trigger.Trigger{
		Name:        "High P99 Latency",
		Description: "Alert when P99 exceeds 500ms",
		Dataset:     "production",
		Query:       q,
		Threshold:   trigger.GreaterThan(500),
		Frequency:   trigger.Minutes(5),
		Recipients: []trigger.Recipient{
			trigger.SlackChannel("#alerts"),
		},
		Disabled: false,
	}

	data, err := TriggerToJSON(tr)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "High P99 Latency", result["name"])
	assert.Equal(t, "production", result["dataset"])
	assert.Equal(t, float64(300), result["frequency"])
	assert.False(t, result["disabled"].(bool))
}

func TestTriggerToJSON_Disabled(t *testing.T) {
	tr := trigger.Trigger{
		Name:     "Disabled Alert",
		Dataset:  "production",
		Disabled: true,
	}

	data, err := TriggerToJSON(tr)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.True(t, result["disabled"].(bool))
}

func TestTriggerToJSONPretty(t *testing.T) {
	tr := trigger.Trigger{
		Name: "Test Trigger",
	}

	data, err := TriggerToJSONPretty(tr)
	require.NoError(t, err)

	assert.Contains(t, string(data), "\n")
	assert.Contains(t, string(data), "  ")
}
