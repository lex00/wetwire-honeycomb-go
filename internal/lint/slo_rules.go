package lint

import (
	"fmt"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

// SLORule represents a lint rule for SLOs.
type SLORule struct {
	Code     string
	Severity Severity
	Message  string
	Check    func(slo discovery.DiscoveredSLO) []LintResult
}

// AllSLORules returns all available SLO lint rules.
func AllSLORules() []SLORule {
	return []SLORule{
		WHC040SLOMissingName(),
		WHC044TargetOutOfRange(),
		WHC047SLONoBurnAlerts(),
	}
}

// WHC040SLOMissingName checks if an SLO is missing a name.
func WHC040SLOMissingName() SLORule {
	return SLORule{
		Code:     "WHC040",
		Severity: SeverityError,
		Message:  "SLO missing name",
		Check: func(slo discovery.DiscoveredSLO) []LintResult {
			if slo.SLOName == "" {
				return []LintResult{
					{
						Rule:     "WHC040",
						Severity: SeverityError,
						Message:  "SLO missing name",
						File:     slo.File,
						Line:     slo.Line,
						Query:    slo.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC044TargetOutOfRange checks if the SLO target percentage is out of valid range (0-100).
func WHC044TargetOutOfRange() SLORule {
	return SLORule{
		Code:     "WHC044",
		Severity: SeverityError,
		Message:  "Target percentage out of range (0-100)",
		Check: func(slo discovery.DiscoveredSLO) []LintResult {
			if slo.TargetPercentage < 0 || slo.TargetPercentage > 100 {
				return []LintResult{
					{
						Rule:     "WHC044",
						Severity: SeverityError,
						Message:  fmt.Sprintf("Target percentage out of range: %.2f (must be 0-100)", slo.TargetPercentage),
						File:     slo.File,
						Line:     slo.Line,
						Query:    slo.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC047SLONoBurnAlerts provides info when an SLO has no burn alerts configured.
func WHC047SLONoBurnAlerts() SLORule {
	return SLORule{
		Code:     "WHC047",
		Severity: SeverityInfo,
		Message:  "SLO has no burn alerts configured",
		Check: func(slo discovery.DiscoveredSLO) []LintResult {
			if slo.BurnAlertCount == 0 {
				return []LintResult{
					{
						Rule:     "WHC047",
						Severity: SeverityInfo,
						Message:  "SLO has no burn alerts configured - consider adding fast and slow burn alerts",
						File:     slo.File,
						Line:     slo.Line,
						Query:    slo.Name,
					},
				}
			}
			return nil
		},
	}
}
