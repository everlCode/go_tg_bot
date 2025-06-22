package message_repository

import (
	"database/sql"
	"log"
)

type MessageRepository struct {
	db *sql.DB
}

type Message struct {
	ID        int
	MessageId int
	FromUser  int
	Text      string
	SendAt    int

	//addtitional params
	UserName string
	UserId   int
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
		ID:        id,
		MessageId: message_id,
		FromUser:  from,
		Text:      text,
		SendAt:    send_at,
	}
}

func (mr *MessageRepository) GetMessagesForToday() []Message {
	rows, error := mr.db.Query(`SELECT u.name, u.telegram_id, m.text, m.send_at FROM messages m
		JOIN users u on u.telegram_id = m.from_user
		WHERE send_at >= strftime('%s', 'now', 'start of day')
	  AND send_at < strftime('%s', 'now', 'start of day', '+1 day');`)

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
			SendAt:   send_at,
			Text:     text,
			UserName: name,
			UserId:   user_id,
		}
		messages = append(messages, msg)
	}

	return messages
}

func (rep MessageRepository) MessageCountForWeek() map[int]int {
	rows, err := rep.db.Query(`
		SELECT from_user, count(*)
		FROM messages m
		WHERE send_at >= STRFTIME("%s", "now", "-7 days")
		GROUP BY from_user
		ORDER BY COUNT(*) DESC;
	`)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()

	stat := make(map[int]int)
	for rows.Next() {
		var userId, count int
		if err := rows.Scan(&userId, &count); err != nil {
			log.Println("scan error:", err)
			continue
		}
		stat[userId] = count
	}

	return stat
}
