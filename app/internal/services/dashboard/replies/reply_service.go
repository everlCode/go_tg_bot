package reply_service

import (
	reply_repository "go-tg-bot/internal/repository/reply"
	user_repository "go-tg-bot/internal/repository/user"
	"log"

	"gopkg.in/telebot.v4"
)

type ReplyService struct {
	rr reply_repository.ReplyRepository
	ur user_repository.UserRepository
}

func NewService(rr reply_repository.ReplyRepository, ur user_repository.UserRepository) ReplyService {
	return ReplyService{
		rr: rr,
		ur: ur,
	}
}

func (rs *ReplyService) Handle(c telebot.Context) {
	msg := c.Message()

	replyToId := msg.ReplyTo.Sender.ID
	text := msg.Text
	fromId := msg.Sender.ID

	rs.rr.Add(fromId, replyToId, text)

	user := rs.ur.UserByTelegramId(fromId)
	
	if user == nil || user.Action < 1 || replyToId == fromId {
		return
	}

	if text == "+" || text == "-" {
		rs.ChangeRespect(replyToId, text)
		rs.DecreaseAction(fromId)

	}
}

func (rs *ReplyService) ChangeRespect(id int64, text string) {
	var add int
	if text == "+" {
		add = 1
	}
	if text == "-" {
		add = -1
	}
	rs.ur.AddRespect(id, add)
}

func (rs *ReplyService) DecreaseAction(id int64) {
	rs.ur.DecreaseAction(id)
}
