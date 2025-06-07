package message_service

import (
	message_repository "go-tg-bot/internal/repository/message"
	reply_repository "go-tg-bot/internal/repository/reply"
	user_repository "go-tg-bot/internal/repository/user"

	"gopkg.in/telebot.v4"
)

type ReplyService struct {
	rr reply_repository.ReplyRepository
	ur user_repository.UserRepository
	mr message_repository.MessageRepository
}

func NewService(rr reply_repository.ReplyRepository, ur user_repository.UserRepository, mr message_repository.MessageRepository) ReplyService {
	return ReplyService{
		rr: rr,
		ur: ur,
		mr: mr,
	}
}

func (rs *ReplyService) Handle(c telebot.Context) {
	msg := c.Message()
	if msg == nil || msg.Sender == nil {
		return
	}

	id := msg.Sender.ID
	msgId := msg.ID
	name := msg.Sender.FirstName
	text := msg.Text

	userExist := rs.ur.UserExist(id)
	if !userExist {
		rs.ur.CreateUser(id, name, 1)
	}

	rs.mr.Create(msgId, id, text, msg.Unixtime)
	rs.ur.AddUserMessageCount(id)

	if msg.IsReply() {
		rs.HandleReply(msg)
	}

	return
}

func (rs *ReplyService) HandleReply(msg *telebot.Message) {

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
