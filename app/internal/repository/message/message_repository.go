package message_repository

import (
	"database/sql"
	"log"
)

type MessageRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (mr *MessageRepository) Create(message_id int, from_user int64, text string, send_at int64) {
	_, err := mr.db.Exec("INSERT INTO messages (message_id, from, text, send_at) VALUES(?, ?, ?, ?)", message_id, from_user, text, send_at)
	if err != nil {
		log.Println(err)
	}
}
