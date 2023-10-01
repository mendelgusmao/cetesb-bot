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
	startMessage      = "Olá, %s! Eu sou o CETESB Praia bot, um bot não oficial que usa dados da CETESB para te informar sobre a qualidade da água das praias do litoral paulista.\n\n" +
		"Para fazer uma consulta, digite o nome da praia, a cidade, ou até melhor, o nome da praia com o nome da cidade. Em instantes eu vou te dizer como está a qualidade da água de lá.\n\n" +
		"Boa praia!"
	samplingPeriodMessage = "Período de amostragem: de %s a %s"
	maxResults = 5
)

var (
	ProperEmojiMapping = map[string]string{
		"Própria":   "\U0001F7E2",
		"Imprópria": "\U0001F534",
	}
)

type Bot struct {
	store    *store.Store
	telegram *telegram.BotAPI
}
