package scraper

import (
	"strings"

	"github.com/go-rod/rod"
)

func New() *Scraper {
	browser := rod.New().MustConnect()

	return &Scraper{browser: browser}
}

func (s *Scraper) ScrapeCities() []City {
	page := s.browser.MustPage(qualityMapURL).MustWaitStable()
	items := page.MustElement("map").MustEval(citiesExtractor).Arr()
	cities := make([]City, len(items))

	for index, item := range items {
		url := item.Str()
		parts := strings.Split(url, "/")
		name := strings.Replace(parts[len(parts)-1], ".phtml", "", -1)

		cities[index] = City{Name: name, URL: url}
	}

	return cities
}

func (s *Scraper) ScrapeBeaches(city City) []Beach {
	page := s.browser.MustPage(city.URL).MustWaitStable()
	items := page.MustElement("body").MustEval(beachExtractor).Arr()
	beaches := make([]Beach, len(items))

	for index, item := range items {
		beaches[index] = Beach{
			City:   city,
			Name:   item.Arr()[1].Str(),
			Proper: item.Arr()[0].Bool(),
		}
	}

	return beaches
}

func (s *Scraper) Finish() {
	s.browser.MustClose()
}
