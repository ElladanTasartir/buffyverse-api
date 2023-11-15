package http

import (
	"context"
	"log"
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

	go func() {
		ctx := context.Background()

		characters, err := scraper.ScrapeCharacters(ctx)
		if err != nil {
			log.Printf("failed to scrape characters. err = %v", err)
			return
		}

		err = s.charactersRepo.CreateCharacters(ctx, characters)
		if err != nil {
			log.Printf("failed to scrape characters. err = %v", err)
			return
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"status": "processing started",
	})
}
