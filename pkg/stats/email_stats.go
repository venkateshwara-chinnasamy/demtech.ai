package stats

import (
	"sync"
	"time"

	"mock-ses-api/internal/models"
)

type EmailStats struct {
	mu               sync.RWMutex
	dailyCount       map[string]int
	totalCount       int
	accountCreatedAt time.Time
	warmupLimits     map[int]int
}

func NewEmailStats() *EmailStats {
	return &EmailStats{
		dailyCount:       make(map[string]int),
		accountCreatedAt: time.Now(),
		warmupLimits: map[int]int{
			1:  200,   // First 24 hours
			7:  1000,  // First week
			14: 5000,  // Second week
			30: 10000, // First month
		},
	}
}

func (s *EmailStats) CheckWarmupLimit() *models.SESError {
	s.mu.RLock()
	defer s.mu.RUnlock()

	daysSinceCreation := int(time.Since(s.accountCreatedAt).Hours() / 24)
	todayKey := time.Now().Format("2006-01-02")
	dailyCount := s.dailyCount[todayKey]

	var applicableLimit int
	for days, limit := range s.warmupLimits {
		if daysSinceCreation <= days {
			applicableLimit = limit
			break
		}
	}

	if applicableLimit > 0 && dailyCount >= applicableLimit {
		return &models.SESError{
			Code:    "Daily quota exceeded",
			Message: "Account is still in warm-up period. Daily sending quota exceeded.",
		}
	}

	return nil
}

func (s *EmailStats) IncrementCount() {
	s.mu.Lock()
	defer s.mu.Unlock()

	todayKey := time.Now().Format("2006-01-02")
	s.dailyCount[todayKey]++
	s.totalCount++
}

func (s *EmailStats) GetCurrentStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"dailyCount":       s.dailyCount,
		"totalCount":       s.totalCount,
		"accountCreatedAt": s.accountCreatedAt,
		"warmupLimits":     s.warmupLimits,
	}
}
