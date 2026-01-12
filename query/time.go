package query

import "time"

// TimeRange represents time parameters for a Honeycomb query.
// Use either relative (TimeRange in seconds) or absolute (StartTime/EndTime).
type TimeRange struct {
	// TimeRange is relative time in seconds (e.g., 7200 = last 2 hours)
	TimeRange int `json:"time_range,omitempty"`

	// StartTime is absolute start time in Unix epoch seconds
	StartTime int `json:"start_time,omitempty"`

	// EndTime is absolute end time in Unix epoch seconds
	EndTime int `json:"end_time,omitempty"`
}

// Hours creates a relative time range for the last N hours.
func Hours(n int) TimeRange {
	return TimeRange{
		TimeRange: n * 3600,
	}
}

// Minutes creates a relative time range for the last N minutes.
func Minutes(n int) TimeRange {
	return TimeRange{
		TimeRange: n * 60,
	}
}

// Days creates a relative time range for the last N days.
func Days(n int) TimeRange {
	return TimeRange{
		TimeRange: n * 86400,
	}
}

// Seconds creates a relative time range for the last N seconds.
func Seconds(n int) TimeRange {
	return TimeRange{
		TimeRange: n,
	}
}

// Absolute creates an absolute time range from start to end.
func Absolute(start, end time.Time) TimeRange {
	return TimeRange{
		StartTime: int(start.Unix()),
		EndTime:   int(end.Unix()),
	}
}

// LastNHours is a convenience function for Hours.
func LastNHours(n int) TimeRange {
	return Hours(n)
}

// Last24Hours creates a time range for the last 24 hours.
func Last24Hours() TimeRange {
	return Hours(24)
}

// Last7Days creates a time range for the last 7 days.
func Last7Days() TimeRange {
	return Days(7)
}
