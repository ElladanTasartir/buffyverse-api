package scraper

import (
	"fmt"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
	"github.com/gocolly/colly"
)

type Scraper interface {
	ScrapeCharacter(address string) (*entity.Character, error)
	ScrapeCharacters() ([]*entity.Character, error)
}

type BuffyScraper struct {
	collector *colly.Collector
	address   string
}

func NewBuffyScraper(address string) (Scraper, error) {
	collector := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	return &BuffyScraper{
		collector: collector,
		address:   address,
	}, nil
}

func (b *BuffyScraper) ScrapeCharacter(address string) (*entity.Character, error) {
	var character entity.Character

	b.collector.OnHTML("aside[role=region]", func(e *colly.HTMLElement) {
		if href, exists := e.DOM.Find("figure > .image").Attr("href"); exists {
			character.Image = href
		}

		character.Name = e.DOM.Find("div[data-source=Name] > div").Text()
		character.Status = e.DOM.Find("div[data-source=Status] > div").Text()
		_ = e.DOM.Find("div[data-source=Born] > div > a").Remove()
		character.Birth = e.DOM.Find("div[data-source=Born] > div").Text()
	})

	err := b.collector.Visit(address)
	if err != nil {
		return &character, fmt.Errorf("failed to scrape character. err = %v", err)
	}

	return &character, nil
}

func (b *BuffyScraper) ScrapeCharacters() ([]*entity.Character, error) {
	searchPages := map[string]string{
		"Scoobies": "Scooby_Gang",
		"Angel":    "Angel_Investigations",
		"Vamps":    "Category:Vampires",
		"Slayers":  "Category:Slayers",
	}

	pageMethods := map[string]func(address string) ([]*entity.Character, error){
		"Scoobies": b.scrapeScoobyMembers,
	}

	var characters []*entity.Character

	for key, page := range searchPages {
		endpoint := b.buildAddress(page)

		if _, ok := pageMethods[key]; !ok {
			continue
		}

		scrapedCharacters, err := pageMethods[key](endpoint)
		if err != nil {
			return scrapedCharacters, fmt.Errorf("failed to scrape characters. key = %s/ err = %v", key, err)
		}

		characters = append(characters, scrapedCharacters...)
	}

	return characters, nil
}

func (b *BuffyScraper) scrapeScoobyMembers(address string) ([]*entity.Character, error) {
	return nil, nil
}

func (b *BuffyScraper) buildAddress(endpoint string) string {
	return fmt.Sprintf("%s/%s", b.address, endpoint)
}
