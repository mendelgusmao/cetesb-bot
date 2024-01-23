package scraper

import (
	"fmt"

	"github.com/go-rod/rod"
)

func New() *Scraper {
	browser := rod.New().MustConnect()

	return &Scraper{browser: browser}
}

func (s *Scraper) Scrape() (cityBeaches map[string][]Beach, e error) {
	cityBeaches = make(map[string][]Beach)

	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("Scraper.Scrape: %v", err)
		}
	}()

	page := s.browser.MustPage(beachQualityURL).MustWaitStable()
	extractedCityBeaches := page.MustElement("body").MustEval(beachExtractor).Map()
	extractedSamplingDates := page.MustElement("body").MustEval(samplingDatesExtractor).Str()
	samplingDates := samplingDatesRE.FindStringSubmatch(extractedSamplingDates)
	page.MustClose()

	for cityName, extractedBeaches := range extractedCityBeaches {
		beaches := make([]Beach, len(extractedBeaches.Arr()))

		for index, extractedBeach := range extractedBeaches.Arr() {
			beaches[index] = Beach{
				City:    City{Name: cityName},
				Name:    extractedBeach.Arr()[0].Str(),
				Quality: extractedBeach.Arr()[2].Str(),
				Sampling: Sampling{
					StartDate: samplingDates[1],
					EndDate:   samplingDates[2],
				},
			}
		}

		cityBeaches[cityName] = append(cityBeaches[cityName], beaches...)
	}

	return cityBeaches, nil
}
