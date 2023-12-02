package store

import (
	"fmt"
	"hash/crc32"
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

func (s *Store) CreateOrUpdateCollection(collectionName string, config database.Configuration, documents []database.Document) error {
	if s.database.CollectionExists(collectionName) {
		return s.database.UpdateCollection(collectionName, documents)
	}

	return s.database.CreateCollection(collectionName, config, documents)
}

func (s *Store) ScrapeAndStore() error {
	cities, beaches := s.ScrapeAndTransform()

	if err := s.CreateOrUpdateCollection("cities", databaseConfiguration, cities); err != nil {
		return err
	}

	return s.CreateOrUpdateCollection("beaches", databaseConfiguration, beaches)
}

func (s *Store) ScrapeAndTransform() (cities []database.Document, beaches []database.Document) {
	scrapedCityBeaches := s.scraper.Scrape()
	cities = make([]database.Document, 0)
	beaches = make([]database.Document, 0)

	checksum := crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v|%v", cities, beaches)))

	if checksum == s.lastChecksum {
		return
	}

	s.lastChecksum = checksum
	log.Println("[store.Scrape] Change detected! Going to update the database.")

	for cityName, cityBeaches := range scrapedCityBeaches {
		for _, beach := range cityBeaches {
			beachDocument := database.Document{
				Keys: []string{
					beach.Name,
					fmt.Sprintf("%s %s", beach.Name, beach.City.Name),
				},
				ExactKeys: []string{
					fmt.Sprintf("%s %s", beach.City.Name, beach.Name),
				},
				Content: beach,
			}

			beaches = append(beaches, beachDocument)
		}

		cities = append(cities, database.Document{
			Keys:      []string{cityName},
			ExactKeys: []string{cityName},
			Content:   cityBeaches,
		})
	}

	log.Printf("[store.Scrape] Found %d cities\n", len(scrapedCityBeaches))
	log.Printf("[store.Scrape] Found %d beaches\n", len(beaches))

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
