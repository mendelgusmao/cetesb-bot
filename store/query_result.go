package store

import (
	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/scoredb/lib/fuzzymap"
)

type QueryResult struct {
	Kind              string
	Beaches           []scraper.Beach
	HasPerfectMatches bool
}

func newQueryResult(kind string, matches []fuzzymap.Match[any]) QueryResult {
	perfectMatches := hasPerfectMatches(matches)
	beaches := make([]scraper.Beach, 0)

	for _, match := range matches {
		if perfectMatches && match.Score < 100 {
			continue
		}

		if matchedBeaches, ok := match.Content.([]scraper.Beach); ok {
			beaches = append(beaches, matchedBeaches...)
		} else {
			beaches = append(beaches, match.Content.(scraper.Beach))
		}
	}

	return QueryResult{
		Kind:              kind,
		Beaches:           beaches,
		HasPerfectMatches: perfectMatches,
	}
}

func hasPerfectMatches(matches []fuzzymap.Match[any]) bool {
	if len(matches) == 0 {
		return false
	}

	perfects := 0

	for _, match := range matches {
		if match.Score == 100 {
			perfects += 1
		}
	}

	return perfects > 0
}
