package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ElladanTasartir/buffyverse-api/internal/config"
	"github.com/ElladanTasartir/buffyverse-api/internal/storage"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer     *http.Server
	engine         *gin.Engine
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
	if config.Environment != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	charactersRepo, err := storage.NewCharactersRepository(db)
	if err != nil {
		return nil, err
	}

	engine.Use(gin.Logger(), gin.Recovery())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		httpServer:     server,
		engine:         engine,
		config:         config,
		storage:        db,
		charactersRepo: charactersRepo,
	}, nil
}

func (s *Server) Start() error {
	s.loadRoutes()

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server err = %v", err)
	}

	return nil
}

func (s *Server) GracefulShutdown() error {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error. err = %v", err)
	}

	log.Println("server has been shutdown gracefully")

	return nil
}

func (s *Server) shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
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
	s.engine.Use(timeout.New(
		timeout.WithTimeout(5*time.Second),
		timeout.WithHandler(s.timeout),
		timeout.WithResponse(s.timeoutResponse),
	))

	s.engine.NoRoute(s.notFound)

	s.engine.GET("/", s.healthCheck)
	s.engine.POST("/scrape/characters", s.scrapeCharacters)
	s.engine.GET("/characters", s.getCharacters)
}

func (s *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
