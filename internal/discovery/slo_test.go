package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverSLOs_BasicSLO(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "slos.go")

	content := `package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

var GoodEvents = query.Query{
	Dataset:   "production",
	TimeRange: query.Days(30),
}

var TotalEvents = query.Query{
	Dataset:   "production",
	TimeRange: query.Days(30),
}

var APIAvailability = slo.SLO{
	Name:        "API Availability",
	Description: "99.9% of requests succeed",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents:  GoodEvents,
		TotalEvents: TotalEvents,
	},
	Target:     slo.Percentage(99.9),
	TimePeriod: slo.Days(30),
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	slos, err := DiscoverSLOs(dir)
	require.NoError(t, err)
	require.Len(t, slos, 1)

	s := slos[0]
	assert.Equal(t, "APIAvailability", s.Name)
	assert.Equal(t, "slos", s.Package)
	assert.Equal(t, testFile, s.File)
	assert.Equal(t, "API Availability", s.SLOName)
	assert.Equal(t, "99.9% of requests succeed", s.Description)
	assert.Equal(t, "production", s.Dataset)
	assert.Equal(t, 99.9, s.TargetPercentage)
	assert.Equal(t, 30, s.TimePeriodDays)
}

func TestDiscoverSLOs_WithQueryRefs(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "slos.go")

	content := `package slos

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
)

var SuccessRequests = query.Query{Dataset: "prod"}
var AllRequests = query.Query{Dataset: "prod"}

var LatencySLO = slo.SLO{
	Name:    "Latency SLO",
	Dataset: "prod",
	SLI: slo.SLI{
		GoodEvents:  SuccessRequests,
		TotalEvents: AllRequests,
	},
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	slos, err := DiscoverSLOs(dir)
	require.NoError(t, err)
	require.Len(t, slos, 1)

	s := slos[0]
	assert.Contains(t, s.GoodEventsQueryRef, "SuccessRequests")
	assert.Contains(t, s.TotalEventsQueryRef, "AllRequests")
}

func TestDiscoverSLOs_WithBurnAlerts(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "slos.go")

	content := `package slos

import "github.com/lex00/wetwire-honeycomb-go/slo"

var APIAvailability = slo.SLO{
	Name:    "API Availability",
	Dataset: "production",
	BurnAlerts: []slo.BurnAlert{
		slo.FastBurn(2.0),
		slo.SlowBurn(5.0),
	},
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	slos, err := DiscoverSLOs(dir)
	require.NoError(t, err)
	require.Len(t, slos, 1)

	s := slos[0]
	assert.Equal(t, 2, s.BurnAlertCount)
}

func TestDiscoverSLOs_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	slos, err := DiscoverSLOs(dir)
	require.NoError(t, err)
	assert.Empty(t, slos)
}

func TestDiscoverSLOs_NoSLOs(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "queries.go")

	content := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
	Dataset: "production",
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	slos, err := DiscoverSLOs(dir)
	require.NoError(t, err)
	assert.Empty(t, slos)
}
