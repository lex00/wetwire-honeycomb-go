package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

func TestWHC050TriggerMissingName(t *testing.T) {
	rule := WHC050TriggerMissingName()

	tests := []struct {
		name      string
		trigger   discovery.DiscoveredTrigger
		wantCount int
	}{
		{
			name: "missing name",
			trigger: discovery.DiscoveredTrigger{
				Name:        "MyTrigger",
				TriggerName: "",
				File:        "test.go",
				Line:        10,
			},
			wantCount: 1,
		},
		{
			name: "has name",
			trigger: discovery.DiscoveredTrigger{
				Name:        "MyTrigger",
				TriggerName: "High Latency Alert",
				File:        "test.go",
				Line:        10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.trigger)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC050", results[0].Rule)
				assert.Equal(t, "error", results[0].Severity)
			}
		})
	}
}

func TestWHC053TriggerNoRecipients(t *testing.T) {
	rule := WHC053TriggerNoRecipients()

	tests := []struct {
		name      string
		trigger   discovery.DiscoveredTrigger
		wantCount int
	}{
		{
			name: "no recipients",
			trigger: discovery.DiscoveredTrigger{
				Name:           "MyTrigger",
				RecipientCount: 0,
				File:           "test.go",
				Line:           10,
			},
			wantCount: 1,
		},
		{
			name: "has recipients",
			trigger: discovery.DiscoveredTrigger{
				Name:           "MyTrigger",
				RecipientCount: 2,
				File:           "test.go",
				Line:           10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.trigger)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC053", results[0].Rule)
				assert.Equal(t, "error", results[0].Severity)
			}
		})
	}
}

func TestWHC054TriggerFrequencyUnder1Minute(t *testing.T) {
	rule := WHC054TriggerFrequencyUnder1Minute()

	tests := []struct {
		name      string
		trigger   discovery.DiscoveredTrigger
		wantCount int
	}{
		{
			name: "frequency under 60 seconds",
			trigger: discovery.DiscoveredTrigger{
				Name:             "MyTrigger",
				FrequencySeconds: 30,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 1,
		},
		{
			name: "frequency at 60 seconds",
			trigger: discovery.DiscoveredTrigger{
				Name:             "MyTrigger",
				FrequencySeconds: 60,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 0,
		},
		{
			name: "frequency over 60 seconds",
			trigger: discovery.DiscoveredTrigger{
				Name:             "MyTrigger",
				FrequencySeconds: 300,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 0,
		},
		{
			name: "no frequency set",
			trigger: discovery.DiscoveredTrigger{
				Name:             "MyTrigger",
				FrequencySeconds: 0,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.trigger)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC054", results[0].Rule)
				assert.Equal(t, "warning", results[0].Severity)
			}
		})
	}
}

func TestWHC056TriggerIsDisabled(t *testing.T) {
	rule := WHC056TriggerIsDisabled()

	tests := []struct {
		name      string
		trigger   discovery.DiscoveredTrigger
		wantCount int
	}{
		{
			name: "disabled trigger",
			trigger: discovery.DiscoveredTrigger{
				Name:     "MyTrigger",
				Disabled: true,
				File:     "test.go",
				Line:     10,
			},
			wantCount: 1,
		},
		{
			name: "enabled trigger",
			trigger: discovery.DiscoveredTrigger{
				Name:     "MyTrigger",
				Disabled: false,
				File:     "test.go",
				Line:     10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.trigger)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC056", results[0].Rule)
				assert.Equal(t, "info", results[0].Severity)
			}
		})
	}
}

func TestAllTriggerRules(t *testing.T) {
	rules := AllTriggerRules()
	assert.GreaterOrEqual(t, len(rules), 4) // At least WHC050, WHC053, WHC054, WHC056
}
