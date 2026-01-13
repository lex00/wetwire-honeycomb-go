package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lex00/wetwire-honeycomb-go/query"
)

func TestQueryPanel_Basic(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(2),
		Calculations: []query.Calculation{
			query.Count(),
			query.P99("duration_ms"),
		},
	}

	panel := QueryPanel(q)
	require.NotNil(t, panel)
	assert.Equal(t, "query", panel.panelType())
}

func TestQueryPanel_WithTitle(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
	}

	panel := QueryPanel(q, WithTitle("Slow Requests"))
	require.NotNil(t, panel)

	qp, ok := panel.(*queryPanel)
	require.True(t, ok)
	assert.Equal(t, "Slow Requests", qp.config.Title)
}

func TestQueryPanel_WithPosition(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
	}

	panel := QueryPanel(q, WithPosition(0, 0, 12, 4))
	require.NotNil(t, panel)

	qp, ok := panel.(*queryPanel)
	require.True(t, ok)
	assert.Equal(t, 0, qp.config.Position.X)
	assert.Equal(t, 0, qp.config.Position.Y)
	assert.Equal(t, 12, qp.config.Position.Width)
	assert.Equal(t, 4, qp.config.Position.Height)
}

func TestQueryPanel_WithMultipleOptions(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
	}

	panel := QueryPanel(q,
		WithTitle("Error Rates"),
		WithPosition(0, 4, 6, 4),
	)
	require.NotNil(t, panel)

	qp, ok := panel.(*queryPanel)
	require.True(t, ok)
	assert.Equal(t, "Error Rates", qp.config.Title)
	assert.Equal(t, 6, qp.config.Position.Width)
}

func TestTextPanel_Basic(t *testing.T) {
	panel := TextPanel("## Notes\nMonitor during peak hours")
	require.NotNil(t, panel)
	assert.Equal(t, "text", panel.panelType())
}

func TestTextPanel_WithOptions(t *testing.T) {
	panel := TextPanel("## Dashboard Notes",
		WithTitle("Notes"),
		WithPosition(0, 8, 12, 2),
	)
	require.NotNil(t, panel)

	tp, ok := panel.(*textPanel)
	require.True(t, ok)
	assert.Equal(t, "## Dashboard Notes", tp.content)
	assert.Equal(t, "Notes", tp.config.Title)
}

func TestSLOPanelByID_Basic(t *testing.T) {
	panel := SLOPanelByID("api-availability-slo")
	require.NotNil(t, panel)
	assert.Equal(t, "slo", panel.panelType())
}

func TestSLOPanelByID_WithOptions(t *testing.T) {
	panel := SLOPanelByID("api-availability-slo",
		WithTitle("API SLO"),
		WithPosition(6, 0, 6, 4),
	)
	require.NotNil(t, panel)

	sp, ok := panel.(*sloPanel)
	require.True(t, ok)
	assert.Equal(t, "api-availability-slo", sp.sloID)
	assert.Equal(t, "API SLO", sp.config.Title)
}

func TestPanel_Interface(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
	}

	// All panel types implement Panel interface
	var panels []Panel
	panels = append(panels, QueryPanel(q))
	panels = append(panels, TextPanel("notes"))
	panels = append(panels, SLOPanelByID("slo-id"))

	require.Len(t, panels, 3)

	expectedTypes := []string{"query", "text", "slo"}
	for i, panel := range panels {
		assert.Equal(t, expectedTypes[i], panel.panelType())
	}
}

func TestPosition_Struct(t *testing.T) {
	pos := Position{
		X:      0,
		Y:      4,
		Width:  6,
		Height: 4,
	}

	assert.Equal(t, 0, pos.X)
	assert.Equal(t, 4, pos.Y)
	assert.Equal(t, 6, pos.Width)
	assert.Equal(t, 4, pos.Height)
}

func TestPanelConfig_Defaults(t *testing.T) {
	q := query.Query{
		Dataset:   "production",
		TimeRange: query.Hours(1),
	}

	// Panel without options should have empty config
	panel := QueryPanel(q)
	qp, ok := panel.(*queryPanel)
	require.True(t, ok)
	assert.Equal(t, "", qp.config.Title)
	assert.Equal(t, 0, qp.config.Position.X)
}
