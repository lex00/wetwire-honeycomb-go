package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

func TestWHC030BoardHasNoPanels(t *testing.T) {
	rule := WHC030BoardHasNoPanels()

	tests := []struct {
		name      string
		board     discovery.DiscoveredBoard
		wantCount int
	}{
		{
			name: "no panels",
			board: discovery.DiscoveredBoard{
				Name:       "EmptyBoard",
				PanelCount: 0,
				File:       "test.go",
				Line:       10,
			},
			wantCount: 1,
		},
		{
			name: "has panels",
			board: discovery.DiscoveredBoard{
				Name:       "Dashboard",
				PanelCount: 3,
				File:       "test.go",
				Line:       10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.board)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC030", results[0].Rule)
				assert.Equal(t, SeverityError, results[0].Severity)
			}
		})
	}
}

func TestWHC034BoardExceedsPanelLimit(t *testing.T) {
	rule := WHC034BoardExceedsPanelLimit()

	tests := []struct {
		name      string
		board     discovery.DiscoveredBoard
		wantCount int
	}{
		{
			name: "exceeds 20 panels",
			board: discovery.DiscoveredBoard{
				Name:       "LargeBoard",
				PanelCount: 25,
				File:       "test.go",
				Line:       10,
			},
			wantCount: 1,
		},
		{
			name: "exactly 20 panels",
			board: discovery.DiscoveredBoard{
				Name:       "MaxBoard",
				PanelCount: 20,
				File:       "test.go",
				Line:       10,
			},
			wantCount: 0,
		},
		{
			name: "under limit",
			board: discovery.DiscoveredBoard{
				Name:       "SmallBoard",
				PanelCount: 5,
				File:       "test.go",
				Line:       10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.board)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC034", results[0].Rule)
				assert.Equal(t, SeverityWarning, results[0].Severity)
			}
		})
	}
}

func TestAllBoardRules(t *testing.T) {
	rules := AllBoardRules()
	assert.GreaterOrEqual(t, len(rules), 2) // At least WHC030 and WHC034
}
