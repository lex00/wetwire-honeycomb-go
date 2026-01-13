package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/query"
)

func TestBoard_BasicFields(t *testing.T) {
	b := Board{
		Name:        "Test Board",
		Description: "A test board",
	}

	assert.Equal(t, "Test Board", b.Name)
	assert.Equal(t, "A test board", b.Description)
}

func TestBoard_WithPanels(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
	}

	b := Board{
		Name: "Performance Board",
		Panels: []Panel{
			QueryPanel(q),
			TextPanel("## Notes"),
		},
	}

	assert.Equal(t, "Performance Board", b.Name)
	require.Len(t, b.Panels, 2)
}

func TestBoard_WithPresetFilters(t *testing.T) {
	b := Board{
		Name: "Filtered Board",
		PresetFilters: []Filter{
			{Column: "service.name", Operation: "=", Value: "api"},
			{Column: "status_code", Operation: ">=", Value: 400},
		},
	}

	require.Len(t, b.PresetFilters, 2)
	assert.Equal(t, "service.name", b.PresetFilters[0].Column)
	assert.Equal(t, "=", b.PresetFilters[0].Operation)
	assert.Equal(t, "api", b.PresetFilters[0].Value)
}

func TestBoard_WithTags(t *testing.T) {
	b := Board{
		Name: "Tagged Board",
		Tags: []Tag{
			{Key: "team", Value: "platform"},
			{Key: "environment", Value: "production"},
		},
	}

	require.Len(t, b.Tags, 2)
	assert.Equal(t, "team", b.Tags[0].Key)
	assert.Equal(t, "platform", b.Tags[0].Value)
}

func TestFilter_AllOperations(t *testing.T) {
	operations := []string{"=", "!=", ">", ">=", "<", "<=", "contains", "does-not-contain"}

	for _, op := range operations {
		f := Filter{
			Column:    "test_column",
			Operation: op,
			Value:     "test_value",
		}
		assert.Equal(t, op, f.Operation)
	}
}
