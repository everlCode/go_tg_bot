package handler

import (
	user_repository "go-tg-bot/internal/repository"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

type WebHookHandler struct {
	userRepository *user_repository.UserRepository
}

func CreateHandler(ur *user_repository.UserRepository) func(tgbotapi.Update) {
	h := &WebHookHandler{
		userRepository: ur,
	}

	return h.Handle
}

func (wh *WebHookHandler) Handle(u tgbotapi.Update) {
	if u.Message == nil || u.Message.From == nil {
		log.Print(u.Message)
		return
	}
	id := u.Message.From.ID
	name := u.Message.From.FirstName
	userExist := wh.userRepository.UserExist(id)

	if userExist {
		wh.userRepository.AddUserMessageCount(id)
	} else {
		wh.userRepository.CreateUser(id, name, 1)
	}
}
