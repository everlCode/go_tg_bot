package message_service

import (
	reaction_repository "go-tg-bot/internal/repository"
	message_repository "go-tg-bot/internal/repository/message"
	reply_repository "go-tg-bot/internal/repository/reply"
	user_repository "go-tg-bot/internal/repository/user"
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

type MessageService struct {
	rr                 reply_repository.ReplyRepository
	ur                 user_repository.UserRepository
	mr                 message_repository.MessageRepository
	reactionRepository reaction_repository.ReactionRepository
	PositiveEmodji     []string
	NegativeEmodji     []string
}

func NewService(
	rr reply_repository.ReplyRepository,
	ur user_repository.UserRepository,
	mr *message_repository.MessageRepository,
	reactionRepository *reaction_repository.ReactionRepository,
) MessageService {
	return MessageService{
		rr:                 rr,
		ur:                 ur,
		mr:                 *mr,
		reactionRepository: *reactionRepository,
		PositiveEmodji:     []string{"👍", "🔥", "\u2764"},
		NegativeEmodji:     []string{"👎", "💩"},
	}
}

func (rs *MessageService) Handle(c telebot.Context) {
	msg := c.Message()
	if msg == nil || msg.Sender == nil {
		return
	}

	id := msg.Sender.ID
	msgId := msg.ID
	name := msg.Sender.FirstName
	text := rs.getText(msg)

	userExist := rs.ur.UserExist(id)
	if !userExist {
		rs.ur.CreateUser(id, name, 1)
	}

	rs.mr.Create(msgId, id, text, msg.Unixtime)
	rs.ur.AddUserMessageCount(id)

	if msg.IsReply() {
		rs.HandleReply(msg)
	}
}

func (rs *MessageService) getText(msg *telebot.Message) string {
	if msg.Text != "" {
		return msg.Text
	}
	if msg.VideoNote != nil {
		return "Видео кружок"
	}

	if msg.Photo != nil {
		return "Фото"
	}

	return ""
}

func (rs *MessageService) HandleReply(msg *telebot.Message) {

	replyToId := msg.ReplyTo.Sender.ID
	text := msg.Text
	fromId := msg.Sender.ID

	rs.rr.Add(fromId, replyToId, text)
}

func (rs *MessageService) HandleReaction(reaction *telebot.MessageReaction) {
	if reaction == nil || len(reaction.NewReaction) == 0 {
		return
	}
	react := reaction.NewReaction[0]
	userFromId := reaction.User.ID

	rs.reactionRepository.Add(userFromId, int64(reaction.MessageID), react.Emoji)

	message := rs.mr.GetById(reaction.MessageID)

	if message == nil {
		return
	}

	if message.FromUser == int(userFromId) {
		return
	}

	user := rs.ur.UserByTelegramId(userFromId)

	if user == nil || user.Action < 1 {
		return
	}

	for _, v := range rs.PositiveEmodji {
		if v == react.Emoji {
			rs.ChangeRespect(message.FromUser, 1)
			rs.DecreaseAction(userFromId)
			break
		}
	}
	for _, v := range rs.NegativeEmodji {
		if v == react.Emoji {
			rs.ChangeRespect(message.FromUser, -1)
			rs.DecreaseAction(userFromId)
			break
		}
	}
}

func (rs *MessageService) ChangeRespect(id int, rate int) {
	log.Println("Change respect ")
	rs.ur.AddRespect(id, rate)
}

func (rs *MessageService) DecreaseAction(id int64) {
	rs.ur.DecreaseAction(id)
}

func (service MessageService) FormatMessagesForGigaChat(messages []message_repository.Message) string {
	names := map[int]string{
		853082525:  "Роман",
		1425523987: "Николай",
		6163710268: "Александр",
		987031755:  "Максим",
		1245777783: "Иван",
		1095940542: "Андрей",
		302048098:  "Андрей",
		812431129:  "Илья",
		788769106:  "Макар",
	}

	content := `Прочитай переписку в чате друзей за сегодня. На осонове сообщений составь краткий обзор, уложись в 25 слов."
	Вот сообщения: \n`

	for _, msg := range messages {
		text := service.deleteBadWords(msg.Text)
		time := time.Unix(
			int64(msg.SendAt),
			0,
		).Format("15:04")
		content += names[msg.UserId] + ": " + text + " " + time + "\n"
	}
	return content
}

func (service MessageService) deleteBadWords(message string) string {
	lower := strings.ToLower(message)
	for _, word := range badWords {
		wordLower := strings.ToLower(word)
		for {
			idx := strings.Index(lower, wordLower)
			if idx == -1 {
				break
			}
			// Удаляем найденное слово (сохраняя регистр оригинального текста)
			message = message[:idx] + message[idx+len(word):]
			lower = strings.ToLower(message)
		}
	}

	return message
}

var badWords = []string{
	"сука", "сучара", "сучка",
	"пиздец", "пизда", "пизде", "пизду", "пиздобол", "пиздабол", "пизд",
	"нахуй", "на хуй", "хуй", "хуевый", "хуев", "нихуя", "охуен", "охуеть", "хули", "хуярить", "хуяр", "ахуе", "хуе", "хуи", "хуя",
	"захуярь",
	"еблан", "ебланище", "еблище", "ебись", "заебца", "заеб", "ебло", "ебать", "ебанный", "ебаный", "ебан", "ебат", "ебля", "ебли",
	"ебал", "разьеб", "разьё", "ебуч", "ебу", "вьеб", "въеб", "въебика", "ёб", "еб",
	"пидар", "пидор", "пидарас", "пидрила",
	"бля", "блять", "блядь", "бл ", "бл\n", "бляди", "блядс",
	"гандон", "гондон",
	"тварь", "твар",
	"долбоеб", "долбоя", "даун", "дибил",
	"говн",
}
