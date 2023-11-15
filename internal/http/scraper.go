package http

import (
	"net/http"

	"github.com/ElladanTasartir/buffyverse-api/internal/scraper"
	"github.com/gin-gonic/gin"
)

func (s *Server) scrapeEpisode(ctx *gin.Context) {}

func (s *Server) scrapeSeason(ctx *gin.Context) {}

func (s *Server) scrapeCharacters(ctx *gin.Context) {
	scraper, err := scraper.NewBuffyScraper(s.config.ScraperURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	characters, err := scraper.ScrapeCharacters()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, characters)
}
