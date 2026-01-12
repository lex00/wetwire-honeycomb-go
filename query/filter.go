package query

// Filter represents a Honeycomb filter condition.
type Filter struct {
	// Column is the field to filter on
	Column string `json:"column"`

	// Op is the filter operator (=, !=, >, <, contains, exists, etc.)
	Op string `json:"op"`

	// Value is the value to compare against (nil for exists/does-not-exist)
	Value any `json:"value,omitempty"`
}

// Equals creates a filter for exact equality (=).
func Equals(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     "=",
		Value:  value,
	}
}

// NotEquals creates a filter for inequality (!=).
func NotEquals(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     "!=",
		Value:  value,
	}
}

// GreaterThan creates a filter for greater than (>).
func GreaterThan(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     ">",
		Value:  value,
	}
}

// GreaterThanOrEqual creates a filter for greater than or equal (>=).
func GreaterThanOrEqual(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     ">=",
		Value:  value,
	}
}

// LessThan creates a filter for less than (<).
func LessThan(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     "<",
		Value:  value,
	}
}

// LessThanOrEqual creates a filter for less than or equal (<=).
func LessThanOrEqual(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     "<=",
		Value:  value,
	}
}

// Contains creates a filter for substring matching.
func Contains(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     "contains",
		Value:  value,
	}
}

// DoesNotContain creates a filter for negative substring matching.
func DoesNotContain(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     "does-not-contain",
		Value:  value,
	}
}

// Exists creates a filter checking if a field exists.
func Exists(column string) Filter {
	return Filter{
		Column: column,
		Op:     "exists",
		Value:  nil,
	}
}

// DoesNotExist creates a filter checking if a field does not exist.
func DoesNotExist(column string) Filter {
	return Filter{
		Column: column,
		Op:     "does-not-exist",
		Value:  nil,
	}
}

// StartsWith creates a filter for prefix matching.
func StartsWith(column string, value any) Filter {
	return Filter{
		Column: column,
		Op:     "starts-with",
		Value:  value,
	}
}

// In creates a filter checking if a value is in a list.
func In(column string, values []any) Filter {
	return Filter{
		Column: column,
		Op:     "in",
		Value:  values,
	}
}

// NotIn creates a filter checking if a value is not in a list.
func NotIn(column string, values []any) Filter {
	return Filter{
		Column: column,
		Op:     "not-in",
		Value:  values,
	}
}

// Convenience aliases for common operators

// GT is an alias for GreaterThan.
func GT(column string, value any) Filter {
	return GreaterThan(column, value)
}

// GTE is an alias for GreaterThanOrEqual.
func GTE(column string, value any) Filter {
	return GreaterThanOrEqual(column, value)
}

// LT is an alias for LessThan.
func LT(column string, value any) Filter {
	return LessThan(column, value)
}

// LTE is an alias for LessThanOrEqual.
func LTE(column string, value any) Filter {
	return LessThanOrEqual(column, value)
}

// Eq is an alias for Equals.
func Eq(column string, value any) Filter {
	return Equals(column, value)
}

// Ne is an alias for NotEquals.
func Ne(column string, value any) Filter {
	return NotEquals(column, value)
}
