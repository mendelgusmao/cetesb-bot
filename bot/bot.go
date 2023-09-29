package bot

import (
	"log"

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
			b.handleUpdate(update)
		}
	}()
}

func (b *Bot) handleUpdate(update telegram.Update) {
	if update.Message != nil {
		query := update.Message.Text
		log.Printf("@%s: %s", update.Message.From.UserName, query)

		result, err := b.store.Query(query)

		if err != nil {
			log.Printf("[ERROR] @%s: %s => %v", update.Message.From.UserName, query, err)
			msg := telegram.NewMessage(update.Message.Chat.ID, unknownErrorMessage)
			b.telegram.Send(msg)
		}

		formatter := NewFormatter(query, result)

		for _, message := range formatter.format() {
			msg := telegram.NewMessage(update.Message.Chat.ID, message)
			b.telegram.Send(msg)
		}
	}
}
