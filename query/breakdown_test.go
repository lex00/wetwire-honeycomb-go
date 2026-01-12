package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBreakdown(t *testing.T) {
	t.Run("Single field", func(t *testing.T) {
		bd := Breakdown("service")
		assert.Equal(t, []string{"service"}, bd)
	})

	t.Run("Multiple fields", func(t *testing.T) {
		bd := Breakdown("service", "endpoint", "user_id")
		assert.Equal(t, []string{"service", "endpoint", "user_id"}, bd)
	})

	t.Run("Empty", func(t *testing.T) {
		bd := Breakdown()
		assert.Equal(t, []string{}, bd)
	})
}
