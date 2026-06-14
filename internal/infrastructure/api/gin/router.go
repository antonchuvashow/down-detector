package gin

import (
	"github.com/gin-gonic/gin"

	apiinspector "detector/internal/infrastructure/api/inspector"
	apireport "detector/internal/infrastructure/api/report"
	"detector/internal/infrastructure/api/route"
	apisuperset "detector/internal/infrastructure/api/superset"
)

type Handlers struct {
	Route     *apiroute.Handler
	Inspector *apiinspector.Handler
	Report    *apireport.Handler
	Superset  *apisuperset.Handler
}

func SetupRoutes(r *gin.Engine, handlers Handlers) {
	// Health check
	// r.GET("/health", health.Check)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Static("/static", "./web/static")
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/report", "./web/report.html")

	// API v1
	api := r.Group("/api/v1")
	{
		// Register controller routes
		handlers.Route.RegisterRoutes(api)
		handlers.Inspector.RegisterInspectors(api)
		handlers.Report.RegisterReports(api)
		handlers.Superset.RegisterSuperset(api)
	}
}
