package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	calc := Count()
	assert.Equal(t, "COUNT", calc.Op)
	assert.Equal(t, "", calc.Column)
}

func TestCountDistinct(t *testing.T) {
	calc := CountDistinct("user_id")
	assert.Equal(t, "COUNT_DISTINCT", calc.Op)
	assert.Equal(t, "user_id", calc.Column)
}

func TestSum(t *testing.T) {
	calc := Sum("amount")
	assert.Equal(t, "SUM", calc.Op)
	assert.Equal(t, "amount", calc.Column)
}

func TestAvg(t *testing.T) {
	calc := Avg("duration_ms")
	assert.Equal(t, "AVG", calc.Op)
	assert.Equal(t, "duration_ms", calc.Column)
}

func TestMax(t *testing.T) {
	calc := Max("latency")
	assert.Equal(t, "MAX", calc.Op)
	assert.Equal(t, "latency", calc.Column)
}

func TestMin(t *testing.T) {
	calc := Min("latency")
	assert.Equal(t, "MIN", calc.Op)
	assert.Equal(t, "latency", calc.Column)
}

func TestPercentiles(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) Calculation
		op       string
		column   string
		expected string
	}{
		{"P50", P50, "P50", "duration_ms", "P50"},
		{"P75", P75, "P75", "duration_ms", "P75"},
		{"P90", P90, "P90", "duration_ms", "P90"},
		{"P95", P95, "P95", "duration_ms", "P95"},
		{"P99", P99, "P99", "duration_ms", "P99"},
		{"P999", P999, "P999", "duration_ms", "P999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := tt.fn(tt.column)
			assert.Equal(t, tt.expected, calc.Op)
			assert.Equal(t, tt.column, calc.Column)
		})
	}
}

func TestHeatmap(t *testing.T) {
	calc := Heatmap("duration_ms")
	assert.Equal(t, "HEATMAP", calc.Op)
	assert.Equal(t, "duration_ms", calc.Column)
}

func TestRate(t *testing.T) {
	calc := Rate("requests")
	assert.Equal(t, "RATE", calc.Op)
	assert.Equal(t, "requests", calc.Column)
}

func TestRateSum(t *testing.T) {
	calc := RateSum("bytes")
	assert.Equal(t, "RATE_SUM", calc.Op)
	assert.Equal(t, "bytes", calc.Column)
}

func TestRateAvg(t *testing.T) {
	calc := RateAvg("duration_ms")
	assert.Equal(t, "RATE_AVG", calc.Op)
	assert.Equal(t, "duration_ms", calc.Column)
}

func TestRateMax(t *testing.T) {
	calc := RateMax("latency")
	assert.Equal(t, "RATE_MAX", calc.Op)
	assert.Equal(t, "latency", calc.Column)
}

func TestConcurrency(t *testing.T) {
	calc := Concurrency()
	assert.Equal(t, "CONCURRENCY", calc.Op)
	assert.Equal(t, "", calc.Column)
}
