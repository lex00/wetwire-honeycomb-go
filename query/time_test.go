package query

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelativeTimeRange(t *testing.T) {
	t.Run("Hours", func(t *testing.T) {
		tr := Hours(2)
		assert.Equal(t, 7200, tr.TimeRange)
		assert.Equal(t, 0, tr.StartTime)
		assert.Equal(t, 0, tr.EndTime)
	})

	t.Run("Minutes", func(t *testing.T) {
		tr := Minutes(30)
		assert.Equal(t, 1800, tr.TimeRange)
		assert.Equal(t, 0, tr.StartTime)
		assert.Equal(t, 0, tr.EndTime)
	})

	t.Run("Days", func(t *testing.T) {
		tr := Days(7)
		assert.Equal(t, 604800, tr.TimeRange)
		assert.Equal(t, 0, tr.StartTime)
		assert.Equal(t, 0, tr.EndTime)
	})

	t.Run("Seconds", func(t *testing.T) {
		tr := Seconds(3600)
		assert.Equal(t, 3600, tr.TimeRange)
		assert.Equal(t, 0, tr.StartTime)
		assert.Equal(t, 0, tr.EndTime)
	})
}

func TestAbsoluteTimeRange(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	tr := Absolute(start, end)
	assert.Equal(t, 0, tr.TimeRange)
	assert.Equal(t, int(start.Unix()), tr.StartTime)
	assert.Equal(t, int(end.Unix()), tr.EndTime)
}

func TestLastNHours(t *testing.T) {
	tr := LastNHours(24)
	assert.Equal(t, 86400, tr.TimeRange)
	assert.Equal(t, 0, tr.StartTime)
	assert.Equal(t, 0, tr.EndTime)
}

func TestLast24Hours(t *testing.T) {
	tr := Last24Hours()
	assert.Equal(t, 86400, tr.TimeRange)
	assert.Equal(t, 0, tr.StartTime)
	assert.Equal(t, 0, tr.EndTime)
}

func TestLast7Days(t *testing.T) {
	tr := Last7Days()
	assert.Equal(t, 604800, tr.TimeRange)
	assert.Equal(t, 0, tr.StartTime)
	assert.Equal(t, 0, tr.EndTime)
}
