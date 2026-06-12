package gin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port     int
	Mode     string
	Handlers Handlers
}

type Server struct {
	engine *gin.Engine
	server *http.Server
}

func NewServer(cfg Config) *Server {
	gin.SetMode(cfg.Mode)

	engine := gin.New()

	// Global middleware
	engine.Use(gin.Recovery())

	// Setup routes
	SetupRoutes(engine, cfg.Handlers)

	return &Server{
		engine: engine,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", cfg.Port),
			Handler:        engine,
			ReadTimeout:    15 * time.Second,
			WriteTimeout:   15 * time.Second,
			IdleTimeout:    60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.Background())
}
