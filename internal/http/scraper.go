package http

import (
	"context"
	"log"
	"net/http"

	"github.com/ElladanTasartir/buffyverse-api/internal/scraper"
	"github.com/gin-gonic/gin"
)

func (s *Server) scrapeEpisode(c *gin.Context) {}

func (s *Server) scrapeSeason(c *gin.Context) {}

func (s *Server) scrapeCharacters(c *gin.Context) {
	scraper, err := scraper.NewBuffyScraper(s.config.ScraperURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	c.JSON(http.StatusOK, gin.H{
		"status": "processing started",
	})
}
