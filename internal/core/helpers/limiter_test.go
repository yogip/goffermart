package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStopLimiter(t *testing.T) {
	t.Run("Test NewStopLimiter", func(t *testing.T) {
		limiter := NewStopLimiter()
		assert.NotNil(t, limiter)
		assert.NotNil(t, limiter.lock)
	})

	t.Run("Test CanWork with future stopUnill", func(t *testing.T) {
		limiter := NewStopLimiter()
		limiter.stopUnill = time.Now().Add(1 * time.Second)
		assert.False(t, limiter.CanWork())
	})

	t.Run("Test CanWork with past stopUnill", func(t *testing.T) {
		limiter := NewStopLimiter()
		limiter.stopUnill = time.Now().Add(-1 * time.Second)
		assert.True(t, limiter.CanWork())
	})

	t.Run("Test EnsureLimit with future stopUnill", func(t *testing.T) {
		limiter := NewStopLimiter()
		limiter.stopUnill = time.Now().Add(1 * time.Second)
		start := time.Now()
		limiter.EnsureLimit()
		assert.InDelta(t, 1*time.Second, time.Since(start), float64(100*time.Millisecond))
	})

	t.Run("Test EnsureLimit with past stopUnill", func(t *testing.T) {
		limiter := NewStopLimiter()
		limiter.stopUnill = time.Now().Add(-1 * time.Second)
		start := time.Now()
		limiter.EnsureLimit()
		assert.InDelta(t, 0, time.Since(start), float64(100*time.Millisecond))
	})

	t.Run("Test StopUntil with zero duration", func(t *testing.T) {
		limiter := NewStopLimiter()
		limiter.StopUntil(0)
		assert.InDelta(t, 0, time.Since(limiter.stopUnill), float64(100*time.Millisecond))
	})

	t.Run("Test StopUntil with negative duration", func(t *testing.T) {
		limiter := NewStopLimiter()
		limiter.StopUntil(-1 * time.Second)
		start := time.Now()
		limiter.EnsureLimit()
		assert.InDelta(t, 0*time.Second, time.Since(start), float64(100*time.Millisecond))
	})

	t.Run("Test StopUntil with positive duration", func(t *testing.T) {
		limiter := NewStopLimiter()
		limiter.StopUntil(1 * time.Second)
		start := time.Now()
		limiter.EnsureLimit()
		assert.InDelta(t, 1*time.Second, time.Since(start), float64(100*time.Millisecond))
	})
}
