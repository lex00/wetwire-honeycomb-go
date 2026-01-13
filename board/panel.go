package board

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
)

// Panel is the interface for board panel types.
type Panel interface {
	panelType() string
}

// Position represents the position and size of a panel on the board.
type Position struct {
	X      int
	Y      int
	Width  int
	Height int
}

// panelConfig holds configuration common to all panel types.
type panelConfig struct {
	Title    string
	Position Position
}

// PanelOption is a function that configures a panel.
type PanelOption func(*panelConfig)

// WithTitle sets the title of a panel.
func WithTitle(title string) PanelOption {
	return func(c *panelConfig) {
		c.Title = title
	}
}

// WithPosition sets the position and size of a panel.
func WithPosition(x, y, width, height int) PanelOption {
	return func(c *panelConfig) {
		c.Position = Position{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		}
	}
}

// queryPanel represents a panel that displays a query visualization.
type queryPanel struct {
	query  query.Query
	config panelConfig
}

func (p *queryPanel) panelType() string {
	return "query"
}

// Query returns the query associated with this panel.
func (p *queryPanel) Query() query.Query {
	return p.query
}

// Config returns the panel configuration.
func (p *queryPanel) Config() panelConfig {
	return p.config
}

// QueryPanel creates a new query panel with the given query and options.
func QueryPanel(q query.Query, opts ...PanelOption) Panel {
	p := &queryPanel{
		query: q,
	}
	for _, opt := range opts {
		opt(&p.config)
	}
	return p
}

// textPanel represents a panel that displays markdown text.
type textPanel struct {
	content string
	config  panelConfig
}

func (p *textPanel) panelType() string {
	return "text"
}

// Content returns the markdown content of this panel.
func (p *textPanel) Content() string {
	return p.content
}

// Config returns the panel configuration.
func (p *textPanel) Config() panelConfig {
	return p.config
}

// TextPanel creates a new text panel with the given markdown content and options.
func TextPanel(content string, opts ...PanelOption) Panel {
	p := &textPanel{
		content: content,
	}
	for _, opt := range opts {
		opt(&p.config)
	}
	return p
}

// sloPanel represents a panel that displays an SLO.
type sloPanel struct {
	sloID  string
	config panelConfig
}

func (p *sloPanel) panelType() string {
	return "slo"
}

// SLOID returns the SLO ID for this panel.
func (p *sloPanel) SLOID() string {
	return p.sloID
}

// Config returns the panel configuration.
func (p *sloPanel) Config() panelConfig {
	return p.config
}

// SLOPanelByID creates a new SLO panel referencing an SLO by its ID.
// Use this for SLOs managed outside of wetwire (e.g., in Terraform).
func SLOPanelByID(id string, opts ...PanelOption) Panel {
	p := &sloPanel{
		sloID: id,
	}
	for _, opt := range opts {
		opt(&p.config)
	}
	return p
}
