package scraper

import (
	"context"
	"fmt"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
	"github.com/gocolly/colly"
)

type Scraper interface {
	ScrapeCharacter(address string) (entity.Character, error)
	ScrapeCharacters(ctx context.Context) ([]entity.Character, error)
	ScrapeEpisodes(address string) ([]entity.Episode, error)
}

type BuffyScraper struct {
	collector *colly.Collector
	address   string
}

func NewBuffyScraper(address string) (Scraper, error) {
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.CacheDir("./tmp"),
	)

	return &BuffyScraper{
		collector: collector,
		address:   address,
	}, nil
}

func (b *BuffyScraper) buildAddress(endpoint string) string {
	return fmt.Sprintf("%s%s", b.address, endpoint)
}
