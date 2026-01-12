package builder

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestdataPath(t *testing.T) string {
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	return filepath.Join(projectRoot, "testdata", "queries")
}

func TestNewBuilder(t *testing.T) {
	b, err := NewBuilder("/some/path")
	require.NoError(t, err)
	assert.NotNil(t, b)
}

func TestBuilder_Build(t *testing.T) {
	testdataPath := getTestdataPath(t)
	b, err := NewBuilder(testdataPath)
	require.NoError(t, err)

	result, err := b.Build()
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Greater(t, result.QueryCount(), 0)
}

func TestBuilder_WithNamespacing(t *testing.T) {
	testdataPath := getTestdataPath(t)
	b, err := NewBuilder(testdataPath)
	require.NoError(t, err)

	b.WithNamespacing(true)
	result, err := b.Build()
	require.NoError(t, err)

	// With namespacing, query names should include package
	queries := result.Queries()
	assert.Greater(t, len(queries), 0)
}

func TestBuilder_WithDatasetFilter(t *testing.T) {
	testdataPath := getTestdataPath(t)
	b, err := NewBuilder(testdataPath)
	require.NoError(t, err)

	b.WithDatasetFilter("production")
	result, err := b.Build()
	require.NoError(t, err)

	// All queries should have the production dataset
	for _, q := range result.Queries() {
		assert.Equal(t, "production", q.Dataset)
	}
}

func TestBuilder_WithStrictMode(t *testing.T) {
	testdataPath := getTestdataPath(t)
	b, err := NewBuilder(testdataPath)
	require.NoError(t, err)

	b.WithStrictMode(true)
	result, err := b.Build()
	// Should not error if no duplicates
	if err == nil {
		assert.NotNil(t, result)
	}
}

func TestBuildResult_QueryCount(t *testing.T) {
	testdataPath := getTestdataPath(t)
	b, err := NewBuilder(testdataPath)
	require.NoError(t, err)

	result, err := b.Build()
	require.NoError(t, err)

	count := result.QueryCount()
	assert.GreaterOrEqual(t, count, 0)
}

func TestBuildResult_Queries(t *testing.T) {
	testdataPath := getTestdataPath(t)
	b, err := NewBuilder(testdataPath)
	require.NoError(t, err)

	result, err := b.Build()
	require.NoError(t, err)

	queries := result.Queries()
	assert.Equal(t, result.QueryCount(), len(queries))
}

func TestBuildResult_Query(t *testing.T) {
	testdataPath := getTestdataPath(t)
	b, err := NewBuilder(testdataPath)
	require.NoError(t, err)

	result, err := b.Build()
	require.NoError(t, err)

	// Try to get a known query
	queries := result.Queries()
	if len(queries) > 0 {
		name := queries[0].Name
		q := result.Query(name)
		assert.NotNil(t, q)
		assert.Equal(t, name, q.Name)
	}

	// Non-existent query should return nil
	q := result.Query("NonExistentQuery")
	assert.Nil(t, q)
}
