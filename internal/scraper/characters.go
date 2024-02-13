package scraper

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ElladanTasartir/buffyverse-api/internal/entity"
	"github.com/gocolly/colly"
)

func (b *BuffyScraper) ScrapeCharacter(address string) (entity.Character, error) {
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
		return character, fmt.Errorf("failed to scrape character. err = %v", err)
	}

	return character, nil
}

func (b *BuffyScraper) ScrapeCharacters(ctx context.Context) ([]entity.Character, error) {
	searchPages := map[string]string{
		"Scoobies": "/wiki/Scooby_Gang",
		"Angel":    "/wiki/Angel_Investigations",
		"Vamps":    "/wiki/Category:Vampires",
		"Slayers":  "/wiki/Category:Slayers",
	}

	pageMethods := map[string]func(address string) ([]entity.Character, error){
		"Scoobies": b.scrapeScoobyMembers,
		"Angel":    b.scrapeAngelMembers,
		"Vamps":    b.scrapeVamps,
		"Slayers":  b.scrapeSlayers,
	}

	var characters []entity.Character

	charactersChan := make(chan []entity.Character, 4)
	errChan := make(chan error, 4)

	wg := sync.WaitGroup{}

	for key, page := range searchPages {
		wg.Add(1)
		go func(key, page string) {
			defer wg.Done()

			if _, ok := pageMethods[key]; !ok {
				return
			}

			scrapedCharacters, err := pageMethods[key](page)
			if err != nil {
				errChan <- fmt.Errorf("failed to scrape characters. key = %s/ err = %v", key, err)
				return
			}

			charactersChan <- scrapedCharacters
		}(key, page)
	}

	wg.Wait()

	close(charactersChan)
	close(errChan)

	for err := range errChan {
		return characters, err
	}

	for scrapedCharacters := range charactersChan {
		characters = append(characters, scrapedCharacters...)
	}

	return characters, nil
}

func (b *BuffyScraper) scrapeScoobyMembers(address string) ([]entity.Character, error) {
	var characters []entity.Character

	b.collector.OnHTML("aside[role=region]", func(e *colly.HTMLElement) {
		var membersHref []string

		for _, node := range e.DOM.Find("div[data-source=Leader] > div > ul").Contents().Nodes {
			for _, attr := range node.FirstChild.Attr {
				if attr.Key == "href" {
					membersHref = append(membersHref, attr.Val)
					break
				}
			}
		}

		for _, node := range e.DOM.Find("div[data-source=Members] > div > ul").Contents().Nodes {
			for _, attr := range node.FirstChild.Attr {
				if attr.Key == "href" {
					membersHref = append(membersHref, attr.Val)
					break
				}
			}
		}

		for _, href := range membersHref {
			character, err := b.ScrapeCharacter(b.buildAddress(href))
			if err != nil {
				log.Printf("Failed to scrape character: %s. err = %v", href, err)
				continue
			}

			characters = append(characters, character)
		}
	})

	b.collector.OnRequest(func(f *colly.Request) {
		log.Printf("scraping scoobies...")
	})

	err := b.collector.Visit(b.buildAddress(address))
	if err != nil {
		return characters, fmt.Errorf("failed to scrape scoobies. err = %v", err)
	}

	return characters, nil
}

func (b *BuffyScraper) scrapeAngelMembers(address string) ([]entity.Character, error) {
	var characters []entity.Character

	b.collector.OnHTML("aside[role=region]", func(e *colly.HTMLElement) {
		var membersHref []string

		for _, node := range e.DOM.Find("div[data-source=Members] > div > ul").Contents().Nodes {
			for _, attr := range node.FirstChild.Attr {
				if attr.Key == "href" {
					membersHref = append(membersHref, attr.Val)
					break
				}
			}
		}

		for _, href := range membersHref {
			character, err := b.ScrapeCharacter(b.buildAddress(href))
			if err != nil {
				log.Printf("Failed to scrape character: %s. err = %v", href, err)
				continue
			}

			characters = append(characters, character)
		}
	})

	b.collector.OnRequest(func(f *colly.Request) {
		log.Printf("scraping angel investigations...")
	})

	err := b.collector.Visit(b.buildAddress(address))
	if err != nil {
		return characters, fmt.Errorf("faild to scrape angel investigations. err = %v", err)
	}

	return characters, nil
}

func (b *BuffyScraper) scrapeVamps(address string) ([]entity.Character, error) {
	var characters []entity.Character

	const nextPage = "?from=Singing+vampire"

	b.collector.OnHTML(".category-page__member-link", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		character, err := b.ScrapeCharacter(b.buildAddress(href))
		if err != nil {
			log.Printf("Failed to scrape character: %s. err = %v", href, err)
			return
		}

		if character.Name == "" {
			return
		}

		characters = append(characters, character)
	})

	b.collector.OnRequest(func(f *colly.Request) {
		log.Printf("scraping vamps...")
	})

	err := b.collector.Visit(b.buildAddress(address))
	if err != nil {
		return characters, fmt.Errorf("faild to scrape vamps. err = %v", err)
	}

	err = b.collector.Visit(b.buildAddress(fmt.Sprintf("%s%s", address, nextPage)))
	if err != nil {
		return characters, fmt.Errorf("faild to scrape vamps. err = %v", err)
	}

	return characters, nil
}

func (b *BuffyScraper) scrapeSlayers(address string) ([]entity.Character, error) {
	var characters []entity.Character

	const nextPage = "?from=Unidentified+Slayer+%281590s%29"

	b.collector.OnHTML(".category-page__member-link", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		character, err := b.ScrapeCharacter(b.buildAddress(href))
		if err != nil {
			log.Printf("Failed to scrape character: %s. err = %v", href, err)
			return
		}

		if character.Name == "" {
			return
		}

		characters = append(characters, character)
	})

	b.collector.OnRequest(func(f *colly.Request) {
		log.Printf("scraping slayers...")
	})

	err := b.collector.Visit(b.buildAddress(address))
	if err != nil {
		return characters, fmt.Errorf("faild to scrape slayers. err = %v", err)
	}

	err = b.collector.Visit(b.buildAddress(fmt.Sprintf("%s%s", address, nextPage)))
	if err != nil {
		return characters, fmt.Errorf("faild to scrape slayers. err = %v", err)
	}

	return characters, nil
}
