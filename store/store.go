package store

import (
	"fmt"
	"log"
	"time"

	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/scoredb/lib/database"
)

func New(database *database.Database, scraper *scraper.Scraper) *Store {
	return &Store{
		database: database,
		scraper:  scraper,
	}
}

func (s *Store) ScrapeAndStore() error {
	cities, beaches := s.Scrape()

	if err := s.Store("cities", cities); err != nil {
		return err
	}

	return s.Store("beaches", beaches)
}

func (s *Store) Scrape() (cities []database.Document, beaches []database.Document) {
	scrapedCities := s.scraper.ScrapeCities()
	cities = make([]database.Document, len(scrapedCities))
	beaches = make([]database.Document, 0)

	log.Printf("[store.Scrape] Found %d cities\n", len(scrapedCities))

	for cityIndex, city := range scrapedCities {
		log.Printf("[store.Scrape] Scraping %s beaches\n", city.Name)

		scrapedBeaches := s.scraper.ScrapeBeaches(city)

		log.Printf("[store.Scrape] Found %d beaches in %s\n", len(scrapedBeaches), city.Name)

		for _, beach := range scrapedBeaches {
			beachDocument := database.Document{
				Keys: []string{
					beach.Name,
				},
				ExactKeys: []string{
					fmt.Sprintf("%s %s", beach.City.Name, beach.Name),
				},
				Content: beach,
			}

			beaches = append(beaches, beachDocument)
		}

		cities[cityIndex] = database.Document{
			Keys:      []string{city.Name},
			ExactKeys: []string{city.Name},
			Content:   scrapedBeaches,
		}
	}

	return
}

func (s *Store) Store(collection string, documents []database.Document) error {
	var err error

	if !s.database.CollectionExists(collection) {
		err = s.database.CreateCollection(collection, databaseConfiguration, documents)
	} else {
		err = s.database.UpdateCollection(collection, documents)
	}

	if err != nil {
		return fmt.Errorf("[store.Store] %v", err)
	}

	return nil
}

func (s *Store) Query(key string) (QueryResult, error) {
	cityMatches, _ := s.database.Query("cities", key)
	citiesQueryResult := newQueryResult("cities", cityMatches)

	if citiesQueryResult.HasPerfectMatches {
		return citiesQueryResult, nil
	}

	beachMatches, _ := s.database.Query("beaches", key)

	return newQueryResult("beaches", beachMatches), nil
}

func (s *Store) Work() {
	if err := s.ScrapeAndStore(); err != nil {
		log.Printf("[store.Work (ticker)] %v", err)
	}

	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.ScrapeAndStore(); err != nil {
					log.Printf("[store.Work (ticker)] %v", err)
				}
			}
		}
	}()
}
