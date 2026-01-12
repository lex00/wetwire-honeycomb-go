package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	filter := Equals("status", "200")
	assert.Equal(t, "status", filter.Column)
	assert.Equal(t, "=", filter.Op)
	assert.Equal(t, "200", filter.Value)
}

func TestNotEquals(t *testing.T) {
	filter := NotEquals("status", "500")
	assert.Equal(t, "status", filter.Column)
	assert.Equal(t, "!=", filter.Op)
	assert.Equal(t, "500", filter.Value)
}

func TestGreaterThan(t *testing.T) {
	filter := GreaterThan("duration_ms", 1000)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, ">", filter.Op)
	assert.Equal(t, 1000, filter.Value)
}

func TestGreaterThanOrEqual(t *testing.T) {
	filter := GreaterThanOrEqual("duration_ms", 1000)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, ">=", filter.Op)
	assert.Equal(t, 1000, filter.Value)
}

func TestLessThan(t *testing.T) {
	filter := LessThan("duration_ms", 100)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, "<", filter.Op)
	assert.Equal(t, 100, filter.Value)
}

func TestLessThanOrEqual(t *testing.T) {
	filter := LessThanOrEqual("duration_ms", 100)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, "<=", filter.Op)
	assert.Equal(t, 100, filter.Value)
}

func TestContains(t *testing.T) {
	filter := Contains("message", "error")
	assert.Equal(t, "message", filter.Column)
	assert.Equal(t, "contains", filter.Op)
	assert.Equal(t, "error", filter.Value)
}

func TestDoesNotContain(t *testing.T) {
	filter := DoesNotContain("message", "debug")
	assert.Equal(t, "message", filter.Column)
	assert.Equal(t, "does-not-contain", filter.Op)
	assert.Equal(t, "debug", filter.Value)
}

func TestExists(t *testing.T) {
	filter := Exists("user_id")
	assert.Equal(t, "user_id", filter.Column)
	assert.Equal(t, "exists", filter.Op)
	assert.Nil(t, filter.Value)
}

func TestDoesNotExist(t *testing.T) {
	filter := DoesNotExist("error")
	assert.Equal(t, "error", filter.Column)
	assert.Equal(t, "does-not-exist", filter.Op)
	assert.Nil(t, filter.Value)
}

func TestStartsWith(t *testing.T) {
	filter := StartsWith("path", "/api")
	assert.Equal(t, "path", filter.Column)
	assert.Equal(t, "starts-with", filter.Op)
	assert.Equal(t, "/api", filter.Value)
}

func TestIn(t *testing.T) {
	filter := In("status", []any{"200", "201", "204"})
	assert.Equal(t, "status", filter.Column)
	assert.Equal(t, "in", filter.Op)
	assert.Equal(t, []any{"200", "201", "204"}, filter.Value)
}

func TestNotIn(t *testing.T) {
	filter := NotIn("status", []any{"500", "502", "503"})
	assert.Equal(t, "status", filter.Column)
	assert.Equal(t, "not-in", filter.Op)
	assert.Equal(t, []any{"500", "502", "503"}, filter.Value)
}

// Convenience aliases
func TestGT(t *testing.T) {
	filter := GT("duration_ms", 500)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, ">", filter.Op)
	assert.Equal(t, 500, filter.Value)
}

func TestGTE(t *testing.T) {
	filter := GTE("duration_ms", 500)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, ">=", filter.Op)
	assert.Equal(t, 500, filter.Value)
}

func TestLT(t *testing.T) {
	filter := LT("duration_ms", 500)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, "<", filter.Op)
	assert.Equal(t, 500, filter.Value)
}

func TestLTE(t *testing.T) {
	filter := LTE("duration_ms", 500)
	assert.Equal(t, "duration_ms", filter.Column)
	assert.Equal(t, "<=", filter.Op)
	assert.Equal(t, 500, filter.Value)
}

func TestEq(t *testing.T) {
	filter := Eq("status", "200")
	assert.Equal(t, "status", filter.Column)
	assert.Equal(t, "=", filter.Op)
	assert.Equal(t, "200", filter.Value)
}

func TestNe(t *testing.T) {
	filter := Ne("status", "500")
	assert.Equal(t, "status", filter.Column)
	assert.Equal(t, "!=", filter.Op)
	assert.Equal(t, "500", filter.Value)
}
