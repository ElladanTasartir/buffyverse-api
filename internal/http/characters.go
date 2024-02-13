package http

import (
	"net/http"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
	"github.com/ElladanTasartir/buffyverse-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCharacters(c *gin.Context) {
	var pagedRequest PagedRequest

	if err := c.ShouldBindQuery(&pagedRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid page params",
		})
		return
	}

	pagedData, err := s.charactersRepo.FindCharacters(c, storage.PageParams{
		Page:     pagedRequest.Page,
		PageSize: pagedRequest.PageSize,
		Search:   pagedRequest.Search,
		Order:    pagedRequest.Order,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, PagedResponse[entity.Character]{
		Result:   pagedData.Results,
		PageSize: pagedRequest.PageSize,
		Page:     pagedRequest.Page,
		Count:    int32(pagedData.Count[0].Count),
	})
}
