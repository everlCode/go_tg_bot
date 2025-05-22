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
	var id int = u.Message.From.ID
	var name string = u.Message.From.FirstName
	log.Println(u.Message.From.FirstName)
	log.Println(u.Message.From.LastName)
	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row, er := db.Query("select * from users where telegram_id = ?", id)
	if er != nil {
		log.Fatal(er)
	}

	if row.Next() {
		log.Println("Пользователь найден")
	} else {
		_, e := db.Exec("INSERT INTO users (name, telegram_id) VALUES(?,?)", name, id)
		if e != nil {
			log.Print(e )
			log.Fatal("here")
		}
	}
}
