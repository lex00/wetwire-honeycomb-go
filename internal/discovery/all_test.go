package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverAll_AllResourceTypes(t *testing.T) {
	dir := t.TempDir()

	// Create queries file
	queriesFile := filepath.Join(dir, "queries.go")
	queriesContent := `package observability

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
}

var ErrorRate = query.Query{
	Dataset: "production",
}
`
	err := os.WriteFile(queriesFile, []byte(queriesContent), 0644)
	require.NoError(t, err)

	// Create SLOs file
	slosFile := filepath.Join(dir, "slos.go")
	slosContent := `package observability

import "github.com/lex00/wetwire-honeycomb-go/slo"

var APIAvailability = slo.SLO{
	Name:    "API Availability",
	Dataset: "production",
}
`
	err = os.WriteFile(slosFile, []byte(slosContent), 0644)
	require.NoError(t, err)

	// Create triggers file
	triggersFile := filepath.Join(dir, "triggers.go")
	triggersContent := `package observability

import "github.com/lex00/wetwire-honeycomb-go/trigger"

var HighLatencyAlert = trigger.Trigger{
	Name:    "High Latency",
	Dataset: "production",
}
`
	err = os.WriteFile(triggersFile, []byte(triggersContent), 0644)
	require.NoError(t, err)

	// Create boards file
	boardsFile := filepath.Join(dir, "boards.go")
	boardsContent := `package observability

import "github.com/lex00/wetwire-honeycomb-go/board"

var PerformanceBoard = board.Board{
	Name: "Performance",
}
`
	err = os.WriteFile(boardsFile, []byte(boardsContent), 0644)
	require.NoError(t, err)

	// Discover all resources
	resources, err := DiscoverAll(dir)
	require.NoError(t, err)

	assert.Len(t, resources.Queries, 2)
	assert.Len(t, resources.SLOs, 1)
	assert.Len(t, resources.Triggers, 1)
	assert.Len(t, resources.Boards, 1)
}

func TestDiscoverAll_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	resources, err := DiscoverAll(dir)
	require.NoError(t, err)

	assert.Empty(t, resources.Queries)
	assert.Empty(t, resources.SLOs)
	assert.Empty(t, resources.Triggers)
	assert.Empty(t, resources.Boards)
}

func TestDiscoveredResources_TotalCount(t *testing.T) {
	resources := &DiscoveredResources{
		Queries:  make([]DiscoveredQuery, 3),
		SLOs:     make([]DiscoveredSLO, 2),
		Triggers: make([]DiscoveredTrigger, 1),
		Boards:   make([]DiscoveredBoard, 4),
	}

	assert.Equal(t, 10, resources.TotalCount())
}
