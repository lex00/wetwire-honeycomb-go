package lint

import (
	"fmt"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

// BoardRule represents a lint rule for boards.
type BoardRule struct {
	Code     string
	Severity Severity
	Message  string
	Check    func(board discovery.DiscoveredBoard) []LintResult
}

// AllBoardRules returns all available board lint rules.
func AllBoardRules() []BoardRule {
	return []BoardRule{
		WHC030BoardHasNoPanels(),
		WHC034BoardExceedsPanelLimit(),
	}
}

// WHC030BoardHasNoPanels checks if a board has no panels.
func WHC030BoardHasNoPanels() BoardRule {
	return BoardRule{
		Code:     "WHC030",
		Severity: SeverityError,
		Message:  "Board has no panels",
		Check: func(board discovery.DiscoveredBoard) []LintResult {
			if board.PanelCount == 0 {
				return []LintResult{
					{
						Rule:     "WHC030",
						Severity: SeverityError,
						Message:  "Board has no panels",
						File:     board.File,
						Line:     board.Line,
						Query:    board.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC034BoardExceedsPanelLimit checks if a board exceeds the recommended panel limit.
func WHC034BoardExceedsPanelLimit() BoardRule {
	return BoardRule{
		Code:     "WHC034",
		Severity: SeverityWarning,
		Message:  "Board exceeds 20 panels",
		Check: func(board discovery.DiscoveredBoard) []LintResult {
			const maxPanels = 20
			if board.PanelCount > maxPanels {
				return []LintResult{
					{
						Rule:     "WHC034",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Board exceeds %d panels (has %d)", maxPanels, board.PanelCount),
						File:     board.File,
						Line:     board.Line,
						Query:    board.Name,
					},
				}
			}
			return nil
		},
	}
}
