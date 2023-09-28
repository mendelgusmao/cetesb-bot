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
	extraItems := page.MustElement("body").MustEval(beachExtraInfoExtractor).Arr()
	beaches := make([]Beach, len(items))

	cityNameHeader := extraItems[0].Str()
	currentDateHeader := extraItems[1].Str()
	samplingDatesHeader := extraItems[2].Str()

	cityName := cityRE.FindStringSubmatch(cityNameHeader)[1]
	currentDate := currentDateRE.FindStringSubmatch(currentDateHeader)[1]
	samplingDates := samplingDatesRE.FindStringSubmatch(samplingDatesHeader)

	for index, item := range items {
		city.Name = cityName

		beaches[index] = Beach{
			City:   city,
			Name:   item.Arr()[1].Str(),
			Proper: item.Arr()[0].Bool(),
			Sampling: Sampling{
				CurrentDate: currentDate,
				StartDate:   samplingDates[1],
				EndDate:     samplingDates[2],
			},
		}
	}

	return beaches
}

func (s *Scraper) Finish() {
	s.browser.MustClose()
}
