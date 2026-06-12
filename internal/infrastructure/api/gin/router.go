package gin

import (
	"github.com/gin-gonic/gin"

	"detector/internal/infrastructure/api/route"
)

type Handlers struct {
	Route *apiroute.Handler
	// Assignment *assignment.Handler
	// Report     *report.Handler
}

func SetupRoutes(r *gin.Engine, handlers Handlers) {
	// Health check
	// r.GET("/health", health.Check)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API v1
	api := r.Group("/api/v1")
	{
		// Register controller routes
		handlers.Route.RegisterRoutes(api)
		// handlers.Assignment.RegisterRoutes(api)
		// handlers.Report.RegisterRoutes(api)
	}
}
