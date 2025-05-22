package handler

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

type WebHookHandler struct {
}

func CreateHandler() func(tgbotapi.Update) {
	h := &WebHookHandler{}

	return h.Handle
}

func (wh *WebHookHandler) Handle(u tgbotapi.Update) {
	log.Println(u.Message.From.ID)
	id := u.Message.From.ID
	name := u.Message.From.FirstName
	log.Println(name, u.Message.From.LastName)

	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = ?)", id).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		log.Print("exist")
		_, err = db.Exec("UPDATE users SET message_count = message_count + 1 WHERE telegram_id = ?", id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err = db.Exec("INSERT INTO users (name, telegram_id) VALUES (?, ?)", name, id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
