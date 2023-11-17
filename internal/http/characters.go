package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getCharacters(ctx *gin.Context) {
	characters, err := s.charactersRepo.FindCharacters(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": characters,
	})
}
