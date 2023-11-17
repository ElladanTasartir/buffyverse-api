package http

import (
	"fmt"

	"github.com/ElladanTasartir/buffyverse-api/internal/config"
	"github.com/ElladanTasartir/buffyverse-api/internal/storage"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer     *gin.Engine
	config         *config.Config
	storage        *storage.Storage
	charactersRepo storage.CharactersStorage
}

func NewServer(config *config.Config, db *storage.Storage) (*Server, error) {
	engine := gin.New()

	charactersRepo, err := storage.NewCharactersRepository(db)
	if err != nil {
		return nil, err
	}

	return &Server{
		httpServer:     engine,
		config:         config,
		storage:        db,
		charactersRepo: charactersRepo,
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
	s.httpServer.POST("/scrape/characters", s.scrapeCharacters)
	s.httpServer.GET("/characters", s.getCharacters)
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "OK",
	})
}
