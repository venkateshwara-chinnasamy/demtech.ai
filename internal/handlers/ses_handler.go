package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mock-ses-api/internal/models"
	"mock-ses-api/pkg/stats"
)

type SESHandler struct {
	stats *stats.EmailStats
}

func NewSESHandler(stats *stats.EmailStats) *SESHandler {
	return &SESHandler{
		stats: stats,
	}
}

func (h *SESHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

func (h *SESHandler) SendEmail(c *gin.Context) {
	var req models.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid request format: " + err.Error())
		c.JSON(http.StatusBadRequest, models.SESError{
			Code:    "ValidationError",
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := h.stats.CheckWarmupLimit(); err != nil {
		log.Println("Too many requests: ", err)
		c.JSON(http.StatusTooManyRequests, err)
		return
	}

	// Update statistics
	h.stats.IncrementCount()

	c.JSON(http.StatusOK, models.SendEmailResponse{
		MessageId: "mock-message-id-" + time.Now().Format("20060102150405"),
		RequestId: "mock-request-id-" + time.Now().Format("20060102150405"),
	})
}

func (h *SESHandler) GetStats(c *gin.Context) {
	stats := h.stats.GetCurrentStats()
	c.JSON(http.StatusOK, stats)
}
