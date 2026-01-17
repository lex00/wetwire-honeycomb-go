package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

func TestWHC040SLOMissingName(t *testing.T) {
	rule := WHC040SLOMissingName()

	tests := []struct {
		name      string
		slo       discovery.DiscoveredSLO
		wantCount int
	}{
		{
			name: "missing name",
			slo: discovery.DiscoveredSLO{
				Name:    "MySLO",
				SLOName: "",
				File:    "test.go",
				Line:    10,
			},
			wantCount: 1,
		},
		{
			name: "has name",
			slo: discovery.DiscoveredSLO{
				Name:    "MySLO",
				SLOName: "API Availability",
				File:    "test.go",
				Line:    10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.slo)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC040", results[0].Rule)
				assert.Equal(t, SeverityError, results[0].Severity)
			}
		})
	}
}

func TestWHC044TargetOutOfRange(t *testing.T) {
	rule := WHC044TargetOutOfRange()

	tests := []struct {
		name      string
		slo       discovery.DiscoveredSLO
		wantCount int
	}{
		{
			name: "target over 100",
			slo: discovery.DiscoveredSLO{
				Name:             "MySLO",
				TargetPercentage: 105.0,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 1,
		},
		{
			name: "target negative",
			slo: discovery.DiscoveredSLO{
				Name:             "MySLO",
				TargetPercentage: -5.0,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 1,
		},
		{
			name: "valid target",
			slo: discovery.DiscoveredSLO{
				Name:             "MySLO",
				TargetPercentage: 99.9,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 0,
		},
		{
			name: "zero target (allowed)",
			slo: discovery.DiscoveredSLO{
				Name:             "MySLO",
				TargetPercentage: 0,
				File:             "test.go",
				Line:             10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.slo)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC044", results[0].Rule)
				assert.Equal(t, SeverityError, results[0].Severity)
			}
		})
	}
}

func TestWHC047SLONoBurnAlerts(t *testing.T) {
	rule := WHC047SLONoBurnAlerts()

	tests := []struct {
		name      string
		slo       discovery.DiscoveredSLO
		wantCount int
	}{
		{
			name: "no burn alerts",
			slo: discovery.DiscoveredSLO{
				Name:           "MySLO",
				BurnAlertCount: 0,
				File:           "test.go",
				Line:           10,
			},
			wantCount: 1,
		},
		{
			name: "has burn alerts",
			slo: discovery.DiscoveredSLO{
				Name:           "MySLO",
				BurnAlertCount: 2,
				File:           "test.go",
				Line:           10,
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := rule.Check(tt.slo)
			assert.Len(t, results, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, "WHC047", results[0].Rule)
				assert.Equal(t, SeverityInfo, results[0].Severity)
			}
		})
	}
}

func TestAllSLORules(t *testing.T) {
	rules := AllSLORules()
	assert.GreaterOrEqual(t, len(rules), 3) // At least WHC040, WHC044, WHC047
}
