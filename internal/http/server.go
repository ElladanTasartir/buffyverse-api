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

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      engine,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := &Server{
		httpServer:     httpServer,
		engine:         engine,
		config:         config,
		storage:        db,
		charactersRepo: charactersRepo,
	}

	engine.Use(gin.Logger(), gin.CustomRecovery(server.errorWrapper))

	return server, nil
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

func (s *Server) errorWrapper(c *gin.Context, err any) {
	log.Printf("panic err in http server. err = %v\n", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": "Internal Server Error",
	})
}

func (s *Server) shutdown(c context.Context) error {
	return s.httpServer.Shutdown(c)
}

func (s *Server) notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "resource not found",
	})
}

func (s *Server) timeout(c *gin.Context) {
	c.Next()
}

func (s *Server) timeoutResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"error": "request timed out",
	})
}

func (s *Server) loadRoutes() {
	s.engine.Use(s.cors)

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

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (s *Server) cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, hx-*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
