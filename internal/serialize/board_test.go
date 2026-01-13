package serialize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

func TestBoardToJSON_BasicFields(t *testing.T) {
	b := board.Board{
		Name:        "Service Performance",
		Description: "Latency and error tracking",
	}

	data, err := BoardToJSON(b)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "Service Performance", result["name"])
	assert.Equal(t, "Latency and error tracking", result["description"])
}

func TestBoardToJSON_WithQueryPanel(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(2),
		Calculations: []query.Calculation{
			query.Count(),
		},
	}

	b := board.Board{
		Name: "Dashboard",
		Panels: []board.Panel{
			board.QueryPanel(q, board.WithTitle("Requests")),
		},
	}

	data, err := BoardToJSON(b)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	panels, ok := result["panels"].([]interface{})
	require.True(t, ok)
	require.Len(t, panels, 1)

	panel := panels[0].(map[string]interface{})
	assert.Equal(t, "query", panel["type"])
	assert.Equal(t, "Requests", panel["title"])
}

func TestBoardToJSON_WithTextPanel(t *testing.T) {
	b := board.Board{
		Name: "Dashboard",
		Panels: []board.Panel{
			board.TextPanel("## Notes\nMonitor during peak hours"),
		},
	}

	data, err := BoardToJSON(b)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	panels, ok := result["panels"].([]interface{})
	require.True(t, ok)
	require.Len(t, panels, 1)

	panel := panels[0].(map[string]interface{})
	assert.Equal(t, "text", panel["type"])
	assert.Equal(t, "## Notes\nMonitor during peak hours", panel["content"])
}

func TestBoardToJSON_WithSLOPanel(t *testing.T) {
	b := board.Board{
		Name: "SLO Dashboard",
		Panels: []board.Panel{
			board.SLOPanelByID("api-availability", board.WithTitle("API SLO")),
		},
	}

	data, err := BoardToJSON(b)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	panels, ok := result["panels"].([]interface{})
	require.True(t, ok)
	require.Len(t, panels, 1)

	panel := panels[0].(map[string]interface{})
	assert.Equal(t, "slo", panel["type"])
	assert.Equal(t, "api-availability", panel["slo_id"])
	assert.Equal(t, "API SLO", panel["title"])
}

func TestBoardToJSON_WithPosition(t *testing.T) {
	b := board.Board{
		Name: "Dashboard",
		Panels: []board.Panel{
			board.TextPanel("notes", board.WithPosition(0, 0, 12, 4)),
		},
	}

	data, err := BoardToJSON(b)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	panels := result["panels"].([]interface{})
	panel := panels[0].(map[string]interface{})
	position := panel["position"].(map[string]interface{})

	assert.Equal(t, float64(0), position["x"])
	assert.Equal(t, float64(0), position["y"])
	assert.Equal(t, float64(12), position["width"])
	assert.Equal(t, float64(4), position["height"])
}

func TestBoardToJSON_WithPresetFilters(t *testing.T) {
	b := board.Board{
		Name: "Filtered Board",
		PresetFilters: []board.Filter{
			{Column: "service.name", Operation: "=", Value: "api"},
		},
	}

	data, err := BoardToJSON(b)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	filters, ok := result["preset_filters"].([]interface{})
	require.True(t, ok)
	require.Len(t, filters, 1)

	filter := filters[0].(map[string]interface{})
	assert.Equal(t, "service.name", filter["column"])
	assert.Equal(t, "=", filter["op"])
	assert.Equal(t, "api", filter["value"])
}

func TestBoardToJSONPretty(t *testing.T) {
	b := board.Board{
		Name: "Test Board",
	}

	data, err := BoardToJSONPretty(b)
	require.NoError(t, err)

	// Should be formatted with indentation
	assert.Contains(t, string(data), "\n")
	assert.Contains(t, string(data), "  ")
}
