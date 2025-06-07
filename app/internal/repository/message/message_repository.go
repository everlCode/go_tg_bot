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
