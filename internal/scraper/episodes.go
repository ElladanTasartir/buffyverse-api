package scraper

import (
	"fmt"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
	"github.com/gocolly/colly"
)

func (b *BuffyScraper) ScrapeEpisodes(address string) ([]entity.Episode, error) {
	pagesMap := map[string]string{
		"Buffy": "/wiki/List_of_Buffy_the_Vampire_Slayer_episodes",
		"Angel": "/wiki/List_of_Angel_episodes",
	}

	pageMethods := map[string]func(address string) ([]entity.Episode, error){
		"Buffy": b.scrapeBuffyEpisodes,
	}

	var episodes []entity.Episode

	for key, page := range pagesMap {
		if _, ok := pageMethods[key]; !ok {
			continue
		}

		scrapedEpisodes, err := pageMethods[key](page)
		if err != nil {
			return scrapedEpisodes, fmt.Errorf("failed to scrape episodes. key = %s/ err = %v", key, err)
		}

		episodes = append(episodes, scrapedEpisodes...)
	}

	return episodes, nil
}

func (b *BuffyScraper) scrapeBuffyEpisodes(address string) ([]entity.Episode, error) {
	var episodes []entity.Episode

	b.collector.OnHTML(".wikitable > tbody", func(e *colly.HTMLElement) {

	})

	err := b.collector.Visit(b.buildAddress(address))
	if err != nil {
		return episodes, fmt.Errorf("failed to scrape angel episodes. err = %v", err)
	}

	return nil, nil
}
