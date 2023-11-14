package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *gin.Engine
}

func NewServer(port int) (*Server, error) {
	engine := gin.New()

	return &Server{
		httpServer: engine,
	}, nil
}

func (s *Server) Start() error {
	err := s.httpServer.Run()
	if err != nil {
		return fmt.Errorf("failed to start http server. err = %v", err)
	}

	return nil
}
