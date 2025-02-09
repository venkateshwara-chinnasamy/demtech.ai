package routes

import (
	"github.com/gin-gonic/gin"

	"mock-ses-api/internal/config"
	"mock-ses-api/internal/handlers"
	"mock-ses-api/pkg/stats"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Initialize dependencies
	emailStats := stats.NewEmailStats()
	sesHandler := handlers.NewSESHandler(emailStats)

	route := router.Group("/api/v1/emails")

	// Health check
	route.GET("/health", sesHandler.HealthCheck)

	// SES API endpoints
	route.POST("/outbound-emails", sesHandler.SendEmail)

	// Stats endpoint (for testing)
	route.GET("/stats", sesHandler.GetStats)

	return router
}
