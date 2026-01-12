package query

// Calculation represents a Honeycomb calculation/aggregation operation.
type Calculation struct {
	// Op is the calculation operation (COUNT, SUM, AVG, P99, etc.)
	Op string `json:"op"`

	// Column is the field to aggregate (empty for COUNT and CONCURRENCY)
	Column string `json:"column,omitempty"`
}

// Count returns the total number of events.
func Count() Calculation {
	return Calculation{
		Op: "COUNT",
	}
}

// CountDistinct returns the count of unique values in a column.
func CountDistinct(column string) Calculation {
	return Calculation{
		Op:     "COUNT_DISTINCT",
		Column: column,
	}
}

// Sum returns the sum of values in a column.
func Sum(column string) Calculation {
	return Calculation{
		Op:     "SUM",
		Column: column,
	}
}

// Avg returns the average of values in a column.
func Avg(column string) Calculation {
	return Calculation{
		Op:     "AVG",
		Column: column,
	}
}

// Max returns the maximum value in a column.
func Max(column string) Calculation {
	return Calculation{
		Op:     "MAX",
		Column: column,
	}
}

// Min returns the minimum value in a column.
func Min(column string) Calculation {
	return Calculation{
		Op:     "MIN",
		Column: column,
	}
}

// P50 returns the 50th percentile (median) of values in a column.
func P50(column string) Calculation {
	return Calculation{
		Op:     "P50",
		Column: column,
	}
}

// P75 returns the 75th percentile of values in a column.
func P75(column string) Calculation {
	return Calculation{
		Op:     "P75",
		Column: column,
	}
}

// P90 returns the 90th percentile of values in a column.
func P90(column string) Calculation {
	return Calculation{
		Op:     "P90",
		Column: column,
	}
}

// P95 returns the 95th percentile of values in a column.
func P95(column string) Calculation {
	return Calculation{
		Op:     "P95",
		Column: column,
	}
}

// P99 returns the 99th percentile of values in a column.
func P99(column string) Calculation {
	return Calculation{
		Op:     "P99",
		Column: column,
	}
}

// P999 returns the 99.9th percentile of values in a column.
func P999(column string) Calculation {
	return Calculation{
		Op:     "P999",
		Column: column,
	}
}

// Heatmap generates a heatmap visualization of value distribution.
func Heatmap(column string) Calculation {
	return Calculation{
		Op:     "HEATMAP",
		Column: column,
	}
}

// Rate returns the rate per second of events.
func Rate(column string) Calculation {
	return Calculation{
		Op:     "RATE",
		Column: column,
	}
}

// RateSum returns the rate per second of summed values.
func RateSum(column string) Calculation {
	return Calculation{
		Op:     "RATE_SUM",
		Column: column,
	}
}

// RateAvg returns the rate per second of average values.
func RateAvg(column string) Calculation {
	return Calculation{
		Op:     "RATE_AVG",
		Column: column,
	}
}

// RateMax returns the rate per second of maximum values.
func RateMax(column string) Calculation {
	return Calculation{
		Op:     "RATE_MAX",
		Column: column,
	}
}

// Concurrency returns the number of concurrent events.
func Concurrency() Calculation {
	return Calculation{
		Op: "CONCURRENCY",
	}
}
