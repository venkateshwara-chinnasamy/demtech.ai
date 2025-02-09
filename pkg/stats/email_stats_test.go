package stats

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewEmailStats(t *testing.T) {
	stats := NewEmailStats()
	assert.NotNil(t, stats)
	assert.NotNil(t, stats.dailyCount)
	assert.NotNil(t, stats.warmupLimits)
	assert.NotZero(t, stats.accountCreatedAt)
}

func TestCheckWarmupLimit(t *testing.T) {
	stats := NewEmailStats()

	// Test within limits
	stats.CheckWarmupLimit()

	// Test exceeding limits
	for i := 0; i < 201; i++ {
		stats.IncrementCount()
	}

	err := stats.CheckWarmupLimit()
	assert.Error(t, errors.New(err.Message))
}

func TestGetCurrentStats(t *testing.T) {
	stats := NewEmailStats()
	stats.IncrementCount()

	currentStats := stats.GetCurrentStats()
	assert.NotNil(t, currentStats)
	assert.Equal(t, 1, currentStats["totalCount"])
}
