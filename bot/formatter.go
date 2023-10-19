package bot

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/mendelgusmao/cetesb-telegram-bot/scraper"
	"github.com/mendelgusmao/cetesb-telegram-bot/store"
)

type Formatter struct {
	query  string
	result store.QueryResult
}

func NewFormatter(query string, result store.QueryResult) *Formatter {
	return &Formatter{
		query:  query,
		result: result,
	}
}

func (f *Formatter) format() (messages []string) {
	messages = make([]string, 0)

	if len(f.result.Beaches) == 0 {
		messages = append(messages, fmt.Sprintf(notFoundMessage, f.query))
		return
	}

	if len(f.result.Beaches) > maxResults {
		beaches := f.result.Beaches

		rand.Shuffle(len(beaches), func(a, b int) {
			beaches[a], beaches[b] = beaches[b], beaches[a]
		})

		f.result.Beaches = beaches[0:maxResults]
		messages = append(messages, maxResultsMessage)
	}

	cityBeaches := make(map[string][]scraper.Beach)
	samplingPeriod := f.result.Beaches[0].Sampling

	for _, beach := range f.result.Beaches {
		if _, ok := cityBeaches[beach.City.Name]; !ok {
			cityBeaches[beach.City.Name] = make([]scraper.Beach, 0)
		}

		cityBeaches[beach.City.Name] = append(cityBeaches[beach.City.Name], beach)
	}

	if len(cityBeaches) > 1 {
		lines := make([]string, 0)

		for _, beaches := range cityBeaches {
			for _, beach := range beaches {
				line := fmt.Sprintf(
					"%s A praia %s em %s está %s para banho!",
					ProperEmojiMapping[beach.Quality],
					strings.Title(strings.ToLower(beach.Name)),
					strings.Title(strings.ToLower(beach.City.Name)),
					strings.ToLower(beach.Quality),
				)

				lines = append(lines, line)
			}
		}

		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf(samplingPeriodMessage, samplingPeriod.StartDate, samplingPeriod.EndDate))

		return append(messages, strings.Join(lines, "\n"))
	}

	for city, beaches := range cityBeaches {
		lines := make([]string, 0)
		manyBeachesFound := len(beaches) > 1

		if manyBeachesFound {
			lines = append(lines, fmt.Sprintf(cityHeaderMessage, city))
			lines = append(lines, "")
		}

		for _, beach := range beaches {
			var line string

			if manyBeachesFound {
				line = fmt.Sprintf(
					"%s A praia %s está %s para banho!",
					ProperEmojiMapping[beach.Quality],
					strings.Title(strings.ToLower(beach.Name)),
					strings.ToLower(beach.Quality),
				)
			} else {
				line = fmt.Sprintf(
					"%s A praia %s em %s está %s para banho!",
					ProperEmojiMapping[beach.Quality],
					strings.Title(strings.ToLower(beach.Name)),
					strings.Title(strings.ToLower(beach.City.Name)),
					strings.ToLower(beach.Quality),
				)
			}

			lines = append(lines, line)
		}

		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf(samplingPeriodMessage, samplingPeriod.StartDate, samplingPeriod.EndDate))

		messages = append(messages, strings.Join(lines, "\n"))
	}

	return messages
}
