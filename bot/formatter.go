package bot

import (
	"fmt"
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
	}

	if len(f.result.Beaches) > maxResults {
		messages = append(messages, fmt.Sprintf(maxResultsMessage, maxResults))
		f.result.Beaches = f.result.Beaches[0:maxResults]
	}

	cityBeaches := make(map[string][]scraper.Beach)

	for _, beach := range f.result.Beaches {
		if _, ok := cityBeaches[beach.City.Name]; !ok {
			cityBeaches[beach.City.Name] = make([]scraper.Beach, 0)
		}

		cityBeaches[beach.City.Name] = append(cityBeaches[beach.City.Name], beach)
	}

	for city, beaches := range cityBeaches {
		lines := make([]string, 0)
		manyBeachesFound := len(beaches) > 1

		if manyBeachesFound {
			lines = append(lines, fmt.Sprintf(cityHeaderMessage, city))
			lines = append(lines, "")
		}

		for _, beach := range beaches {
			line := fmt.Sprintf(
				"%s A praia %s em %s está %s para banho!",
				ProperEmojiMapping[beach.Quality],
				strings.Title(strings.ToLower(beach.Name)),
				strings.Title(strings.ToLower(beach.City.Name)),
				strings.ToLower(beach.Quality),
			)

			if manyBeachesFound {
				line = fmt.Sprintf(
					"%s A praia %s está %s para banho!",
					ProperEmojiMapping[beach.Quality],
					strings.Title(strings.ToLower(beach.Name)),
					strings.ToLower(beach.Quality),
				)
			}

			lines = append(lines, line)
		}

		messages = append(messages, strings.Join(lines, "\n"))
	}

	return messages
}
