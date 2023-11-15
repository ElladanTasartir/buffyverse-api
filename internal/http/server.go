package http

import (
	"fmt"

	"github.com/ElladanTasartir/buffyverse-api/internal/config"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *gin.Engine
	config     *config.Config
}

func NewServer(config *config.Config) (*Server, error) {
	engine := gin.New()

	return &Server{
		httpServer: engine,
		config:     config,
	}, nil
}

func (s *Server) Start() error {
	s.loadRoutes()

	err := s.httpServer.Run(fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to start http server. err = %v", err)
	}

	return nil
}

func (s *Server) loadRoutes() {
	s.httpServer.GET("/", s.healthCheck)
	s.httpServer.POST("/scrape/character", s.scrapeCharacters)
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "OK",
	})
}
