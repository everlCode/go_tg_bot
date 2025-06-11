package message_repository

import (
	"database/sql"
	"log"
)

type MessageRepository struct {
	db *sql.DB
}

type Message struct {
	ID int
	MessageId int
	FromUser int
	Text string
	SendAt int

	//addtitional params
	UserName string
	UserId int
}

func NewRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (mr *MessageRepository) Create(message_id int, from_user int64, text string, send_at int64) {
	_, err := mr.db.Exec("INSERT INTO messages (message_id, from_user, text, send_at) VALUES(?, ?, ?, ?)", message_id, from_user, text, send_at)
	if err != nil {
		log.Println(err)
	}
}

func (mr *MessageRepository) GetById(message_id int) *Message {
	msg := mr.db.QueryRow("SELECT * FROM messages WHERE message_id = ?", message_id)
	var id, from, send_at int
	var text string
	
	msg.Scan(&id, &message_id, &from, &text, &send_at)

	return &Message{
		ID: id,
		MessageId: message_id,
		FromUser: from,
		Text: text,
		SendAt: send_at,
	}
}

func (mr *MessageRepository) GetMessagesForToday() []Message {
	rows, error := mr.db.Query(`SELECT u.name, u.telegram_id, m.text, m.send_at FROM messages m
		JOIN users u on u.telegram_id = m.from_user;`)

	if error != nil {
		log.Println(error)
		return nil
	}

	var messages []Message

	for rows.Next() {
		var send_at, user_id int
		var text, name string

		rows.Scan(&name, &user_id, &text, &send_at)
		msg := Message{
			SendAt: send_at,
			Text: text,
			UserName: name,
			UserId: user_id,
		}
		messages = append(messages, msg)
	}

	return messages
}