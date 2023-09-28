package bot

import (
	"fmt"
	"log"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mendelgusmao/cetesb-telegram-bot/store"
)

func New(store *store.Store, telegramAPIToken string) *Bot {
	telegram, err := telegram.NewBotAPI(telegramAPIToken)

	if err != nil {
		log.Panic(err)
	}

	return &Bot{store: store, telegram: telegram}
}

func (b *Bot) Work() {
	go func() {
		u := telegram.NewUpdate(0)
		u.Timeout = 60

		updates := b.telegram.GetUpdatesChan(u)

		for update := range updates {
			if update.Message != nil {
				log.Printf("@%s: %s", update.Message.From.UserName, update.Message.Text)

				result, err := b.store.Query(update.Message.Text)

				if err != nil {
					log.Printf("[ERROR] @%s: %s => %v", update.Message.From.UserName, update.Message.Text, err)
					msg := telegram.NewMessage(update.Message.Chat.ID, unknownErrorMessage)
					b.telegram.Send(msg)
				}

				lines := make([]string, len(result.Beaches))

				for index, result := range result.Beaches {
					lines[index] = fmt.Sprintf(
						"%s A praia %s na cidade de %s está %s para banho!",
						ProperEmojiMapping[result.Proper],
						strings.Title(strings.ToLower(result.Name)),
						strings.Title(result.City.Name),
						ProperTextMapping[result.Proper],
					)
				}

				msg := telegram.NewMessage(update.Message.Chat.ID, strings.Join(lines, "\n"))
				b.telegram.Send(msg)
			}
		}
	}()
}
