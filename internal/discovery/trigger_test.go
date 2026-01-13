package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverTriggers_BasicTrigger(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "triggers.go")

	content := `package triggers

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

var SlowRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
}

var HighLatencyAlert = trigger.Trigger{
	Name:        "High P99 Latency",
	Description: "Alert when P99 exceeds 500ms",
	Dataset:     "production",
	Query:       SlowRequests,
	Threshold:   trigger.GreaterThan(500),
	Frequency:   trigger.Minutes(5),
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	triggers, err := DiscoverTriggers(dir)
	require.NoError(t, err)
	require.Len(t, triggers, 1)

	tr := triggers[0]
	assert.Equal(t, "HighLatencyAlert", tr.Name)
	assert.Equal(t, "triggers", tr.Package)
	assert.Equal(t, testFile, tr.File)
	assert.Equal(t, "High P99 Latency", tr.TriggerName)
	assert.Equal(t, "Alert when P99 exceeds 500ms", tr.Description)
	assert.Equal(t, "production", tr.Dataset)
	assert.Equal(t, ">", tr.ThresholdOp)
	assert.Equal(t, 500.0, tr.ThresholdValue)
	assert.Equal(t, 300, tr.FrequencySeconds)
}

func TestDiscoverTriggers_WithQueryRef(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "triggers.go")

	content := `package triggers

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
)

var ErrorRate = query.Query{Dataset: "prod"}

var ErrorAlert = trigger.Trigger{
	Name:    "Error Spike",
	Dataset: "prod",
	Query:   ErrorRate,
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	triggers, err := DiscoverTriggers(dir)
	require.NoError(t, err)
	require.Len(t, triggers, 1)

	tr := triggers[0]
	assert.Equal(t, "ErrorRate", tr.QueryRef)
}

func TestDiscoverTriggers_WithRecipients(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "triggers.go")

	content := `package triggers

import "github.com/lex00/wetwire-honeycomb-go/trigger"

var HighLatencyAlert = trigger.Trigger{
	Name:    "High Latency",
	Dataset: "production",
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#alerts"),
		trigger.PagerDutyService("api-team"),
		trigger.EmailAddress("team@example.com"),
	},
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	triggers, err := DiscoverTriggers(dir)
	require.NoError(t, err)
	require.Len(t, triggers, 1)

	tr := triggers[0]
	assert.Equal(t, 3, tr.RecipientCount)
}

func TestDiscoverTriggers_Disabled(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "triggers.go")

	content := `package triggers

import "github.com/lex00/wetwire-honeycomb-go/trigger"

var DisabledAlert = trigger.Trigger{
	Name:     "Disabled Alert",
	Dataset:  "production",
	Disabled: true,
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	triggers, err := DiscoverTriggers(dir)
	require.NoError(t, err)
	require.Len(t, triggers, 1)

	tr := triggers[0]
	assert.True(t, tr.Disabled)
}

func TestDiscoverTriggers_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	triggers, err := DiscoverTriggers(dir)
	require.NoError(t, err)
	assert.Empty(t, triggers)
}

func TestDiscoverTriggers_NoTriggers(t *testing.T) {
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

	triggers, err := DiscoverTriggers(dir)
	require.NoError(t, err)
	assert.Empty(t, triggers)
}
