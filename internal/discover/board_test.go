package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscoverBoards_BasicBoard(t *testing.T) {
	// Create a temporary directory with a test file
	dir := t.TempDir()
	testFile := filepath.Join(dir, "boards.go")

	content := `package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

var PerformanceBoard = board.Board{
	Name:        "Service Performance",
	Description: "Latency and error tracking",
	Panels: []board.Panel{
		board.QueryPanel(SlowRequests),
		board.TextPanel("## Notes"),
	},
}

var SlowRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	boards, err := DiscoverBoards(dir)
	require.NoError(t, err)
	require.Len(t, boards, 1)

	board := boards[0]
	assert.Equal(t, "PerformanceBoard", board.Name)
	assert.Equal(t, "boards", board.Package)
	assert.Equal(t, testFile, board.File)
	assert.Equal(t, "Service Performance", board.BoardName)
	assert.Equal(t, "Latency and error tracking", board.Description)
	assert.Equal(t, 2, board.PanelCount)
}

func TestDiscoverBoards_WithQueryRefs(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "boards.go")

	content := `package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
	"github.com/lex00/wetwire-honeycomb-go/query"
)

var SlowRequests = query.Query{
	Dataset: "production",
}

var ErrorRates = query.Query{
	Dataset: "production",
}

var Dashboard = board.Board{
	Name: "Dashboard",
	Panels: []board.Panel{
		board.QueryPanel(SlowRequests),
		board.QueryPanel(ErrorRates),
	},
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	boards, err := DiscoverBoards(dir)
	require.NoError(t, err)
	require.Len(t, boards, 1)

	board := boards[0]
	assert.Equal(t, "Dashboard", board.Name)
	assert.Equal(t, 2, board.PanelCount)
	assert.Contains(t, board.QueryRefs, "SlowRequests")
	assert.Contains(t, board.QueryRefs, "ErrorRates")
}

func TestDiscoverBoards_WithSLORefs(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "boards.go")

	content := `package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
)

var Dashboard = board.Board{
	Name: "SLO Dashboard",
	Panels: []board.Panel{
		board.SLOPanelByID("api-availability"),
		board.SLOPanelByID("latency-slo"),
	},
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	boards, err := DiscoverBoards(dir)
	require.NoError(t, err)
	require.Len(t, boards, 1)

	board := boards[0]
	assert.Equal(t, 2, board.PanelCount)
	assert.Contains(t, board.SLORefs, "api-availability")
	assert.Contains(t, board.SLORefs, "latency-slo")
}

func TestDiscoverBoards_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	boards, err := DiscoverBoards(dir)
	require.NoError(t, err)
	assert.Empty(t, boards)
}

func TestDiscoverBoards_NoBoards(t *testing.T) {
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

	boards, err := DiscoverBoards(dir)
	require.NoError(t, err)
	assert.Empty(t, boards)
}

func TestDiscoverBoards_SkipsTestFiles(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "boards_test.go")

	content := `package boards

import (
	"github.com/lex00/wetwire-honeycomb-go/board"
)

var TestBoard = board.Board{
	Name: "Test Board",
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	boards, err := DiscoverBoards(dir)
	require.NoError(t, err)
	assert.Empty(t, boards)
}
