package bot

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mendelgusmao/cetesb-telegram-bot/store"
)

const (
	unknownErrorMessage = "Ocorreu um erro na consulta. Tente novamente mais tarde."
	notFoundMessage     = "Não encontrei informações sobre %s. Modifique sua consulta e tente novamente."
	maxResultsMessage   = "Encontrei muitas praias e vou te mandar informações de %d delas. " +
		"Caso não encontre a que quer, tente fazer uma consulta mais específica."
	cityHeaderMessage = "Encontrei informações das seguintes praias da cidade de %s:"
	maxResults        = 5
)

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
