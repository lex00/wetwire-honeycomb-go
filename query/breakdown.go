package query

// Breakdown creates a list of field names to group by.
// These correspond to the "breakdowns" array in a Honeycomb query.
func Breakdown(fields ...string) []string {
	if fields == nil {
		return []string{}
	}
	return fields
}
