package serialize

import (
	"bytes"
	"encoding/json"

	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// boardJSON is the internal representation for board JSON serialization.
type boardJSON struct {
	Name          string            `json:"name"`
	Description   string            `json:"description,omitempty"`
	Panels        []panelJSON       `json:"panels,omitempty"`
	PresetFilters []boardFilterJSON `json:"preset_filters,omitempty"`
}

type panelJSON struct {
	Type     string        `json:"type"`
	Title    string        `json:"title,omitempty"`
	Position *positionJSON `json:"position,omitempty"`
	// For query panels
	Query *queryJSON `json:"query,omitempty"`
	// For text panels
	Content string `json:"content,omitempty"`
	// For SLO panels
	SLOID string `json:"slo_id,omitempty"`
}

type positionJSON struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type boardFilterJSON struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Value  any    `json:"value,omitempty"`
}

// Panel accessor interfaces for type assertions
type queryPanelAccessor interface {
	Query() query.Query
	Config() board.PanelConfig
}

type textPanelAccessor interface {
	Content() string
	Config() board.PanelConfig
}

type sloPanelAccessor interface {
	SLOID() string
	Config() board.PanelConfig
}

// BoardToJSON serializes a Board to Honeycomb Board JSON format.
func BoardToJSON(b board.Board) ([]byte, error) {
	jb := toBoardJSON(b)
	return json.Marshal(jb)
}

// BoardToJSONPretty serializes a Board to indented JSON format.
func BoardToJSONPretty(b board.Board) ([]byte, error) {
	jb := toBoardJSON(b)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(jb); err != nil {
		return nil, err
	}
	result := buf.Bytes()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

func toBoardJSON(b board.Board) boardJSON {
	jb := boardJSON{
		Name:        b.Name,
		Description: b.Description,
	}

	// Convert panels
	if len(b.Panels) > 0 {
		jb.Panels = make([]panelJSON, len(b.Panels))
		for i, p := range b.Panels {
			jb.Panels[i] = toPanelJSON(p)
		}
	}

	// Convert preset filters
	if len(b.PresetFilters) > 0 {
		jb.PresetFilters = make([]boardFilterJSON, len(b.PresetFilters))
		for i, f := range b.PresetFilters {
			jb.PresetFilters[i] = boardFilterJSON{
				Column: f.Column,
				Op:     f.Operation,
				Value:  f.Value,
			}
		}
	}

	return jb
}

func toPanelJSON(p board.Panel) panelJSON {
	jp := panelJSON{}

	// Use type assertions to access panel-specific data
	switch panel := p.(type) {
	case queryPanelAccessor:
		jp.Type = "query"
		config := panel.Config()
		jp.Title = config.Title
		if config.Position.Width > 0 || config.Position.Height > 0 {
			jp.Position = &positionJSON{
				X:      config.Position.X,
				Y:      config.Position.Y,
				Width:  config.Position.Width,
				Height: config.Position.Height,
			}
		}
		q := panel.Query()
		qj := toQueryJSON(q)
		jp.Query = &qj

	case textPanelAccessor:
		jp.Type = "text"
		config := panel.Config()
		jp.Title = config.Title
		jp.Content = panel.Content()
		if config.Position.Width > 0 || config.Position.Height > 0 {
			jp.Position = &positionJSON{
				X:      config.Position.X,
				Y:      config.Position.Y,
				Width:  config.Position.Width,
				Height: config.Position.Height,
			}
		}

	case sloPanelAccessor:
		jp.Type = "slo"
		config := panel.Config()
		jp.Title = config.Title
		jp.SLOID = panel.SLOID()
		if config.Position.Width > 0 || config.Position.Height > 0 {
			jp.Position = &positionJSON{
				X:      config.Position.X,
				Y:      config.Position.Y,
				Width:  config.Position.Width,
				Height: config.Position.Height,
			}
		}
	}

	return jp
}
