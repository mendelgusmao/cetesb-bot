package bot

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mendelgusmao/cetesb-telegram-bot/store"
)

const unknownErrorMessage = "Ocorreu um erro na consulta. Tente novamente mais tarde."

var (
	ProperEmojiMapping = map[bool]string{
		true:  "\U0001F7E2",
		false: "\U0001F534",
	}
	ProperTextMapping = map[bool]string{
		true:  "própria",
		false: "imprópria",
	}
)

type Bot struct {
	store    *store.Store
	telegram *telegram.BotAPI
}
