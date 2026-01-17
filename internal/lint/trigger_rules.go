package lint

import (
	"fmt"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

// TriggerRule represents a lint rule for triggers.
type TriggerRule struct {
	Code     string
	Severity Severity
	Message  string
	Check    func(trigger discovery.DiscoveredTrigger) []LintResult
}

// AllTriggerRules returns all available trigger lint rules.
func AllTriggerRules() []TriggerRule {
	return []TriggerRule{
		WHC050TriggerMissingName(),
		WHC053TriggerNoRecipients(),
		WHC054TriggerFrequencyUnder1Minute(),
		WHC056TriggerIsDisabled(),
	}
}

// WHC050TriggerMissingName checks if a trigger is missing a name.
func WHC050TriggerMissingName() TriggerRule {
	return TriggerRule{
		Code:     "WHC050",
		Severity: SeverityError,
		Message:  "Trigger missing name",
		Check: func(trigger discovery.DiscoveredTrigger) []LintResult {
			if trigger.TriggerName == "" {
				return []LintResult{
					{
						Rule:     "WHC050",
						Severity: SeverityError,
						Message:  "Trigger missing name",
						File:     trigger.File,
						Line:     trigger.Line,
						Query:    trigger.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC053TriggerNoRecipients checks if a trigger has no recipients configured.
func WHC053TriggerNoRecipients() TriggerRule {
	return TriggerRule{
		Code:     "WHC053",
		Severity: SeverityError,
		Message:  "Trigger has no recipients",
		Check: func(trigger discovery.DiscoveredTrigger) []LintResult {
			if trigger.RecipientCount == 0 {
				return []LintResult{
					{
						Rule:     "WHC053",
						Severity: SeverityError,
						Message:  "Trigger has no recipients - alerts won't be delivered",
						File:     trigger.File,
						Line:     trigger.Line,
						Query:    trigger.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC054TriggerFrequencyUnder1Minute warns when trigger frequency is under 1 minute.
func WHC054TriggerFrequencyUnder1Minute() TriggerRule {
	return TriggerRule{
		Code:     "WHC054",
		Severity: SeverityWarning,
		Message:  "Trigger frequency under 1 minute",
		Check: func(trigger discovery.DiscoveredTrigger) []LintResult {
			if trigger.FrequencySeconds > 0 && trigger.FrequencySeconds < 60 {
				return []LintResult{
					{
						Rule:     "WHC054",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Trigger frequency under 1 minute (%d seconds) may cause excessive alerts", trigger.FrequencySeconds),
						File:     trigger.File,
						Line:     trigger.Line,
						Query:    trigger.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC056TriggerIsDisabled provides info when a trigger is disabled.
func WHC056TriggerIsDisabled() TriggerRule {
	return TriggerRule{
		Code:     "WHC056",
		Severity: SeverityInfo,
		Message:  "Trigger is disabled",
		Check: func(trigger discovery.DiscoveredTrigger) []LintResult {
			if trigger.Disabled {
				return []LintResult{
					{
						Rule:     "WHC056",
						Severity: SeverityInfo,
						Message:  "Trigger is disabled",
						File:     trigger.File,
						Line:     trigger.Line,
						Query:    trigger.Name,
					},
				}
			}
			return nil
		},
	}
}
