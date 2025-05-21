package handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type WebHookHandler struct {
}

func CreateHandler() func(tgbotapi.Update) {
	h := &WebHookHandler{}

	return h.Handle
}

func (wh *WebHookHandler) Handle(u tgbotapi.Update) {
	log.Println(u.Message)
}
