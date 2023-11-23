package http

import (
	"net/http"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
	"github.com/ElladanTasartir/buffyverse-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCharacters(ctx *gin.Context) {
	var pagedRequest PagedRequest

	if err := ctx.ShouldBindQuery(&pagedRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid page params",
		})
		return
	}

	pagedData, err := s.charactersRepo.FindCharacters(ctx, storage.PageParams{
		Page:     pagedRequest.Page,
		PageSize: pagedRequest.PageSize,
		Search:   pagedRequest.Search,
		Order:    pagedRequest.Order,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, PagedResponse[entity.Character]{
		Result:   pagedData.Results,
		PageSize: pagedRequest.PageSize,
		Page:     pagedRequest.Page,
		Count:    int32(pagedData.Count[0].Count),
	})
}
