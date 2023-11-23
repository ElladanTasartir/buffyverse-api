package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ElladanTasartir/buffyverse-api/internal/config"
	"github.com/ElladanTasartir/buffyverse-api/internal/storage"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer     *gin.Engine
	config         *config.Config
	storage        *storage.Storage
	charactersRepo storage.CharactersStorage
}

type PagedResponse[T interface{}] struct {
	Result   []T   `json:"results"`
	PageSize int32 `json:"pageSize"`
	Page     int32 `json:"page"`
	Count    int32 `json:"count"`
}

type PagedRequest struct {
	PageSize int32  `form:"pageSize"`
	Page     int32  `form:"page"`
	Search   string `form:"search"`
	Order    string `form:"order"`
}

func NewServer(config *config.Config, db *storage.Storage) (*Server, error) {
	engine := gin.New()

	charactersRepo, err := storage.NewCharactersRepository(db)
	if err != nil {
		return nil, err
	}

	engine.Use(gin.Logger(), gin.Recovery())

	return &Server{
		httpServer:     engine,
		config:         config,
		storage:        db,
		charactersRepo: charactersRepo,
	}, nil
}

func (s *Server) Start() error {
	s.loadRoutes()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.httpServer,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start http server. err = %v", err)
	}

	return nil
}

func (s *Server) notFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"error": "resource not found",
	})
}

func (s *Server) timeout(ctx *gin.Context) {
	ctx.Next()
}

func (s *Server) timeoutResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusGatewayTimeout, gin.H{
		"error": "request timed out",
	})
}

func (s *Server) loadRoutes() {
	s.httpServer.Use(timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(s.timeout),
		timeout.WithResponse(s.timeoutResponse),
	))

	s.httpServer.NoRoute(s.notFound)

	s.httpServer.GET("/", s.healthCheck)
	s.httpServer.POST("/scrape/characters", s.scrapeCharacters)
	s.httpServer.GET("/characters", s.getCharacters)
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
