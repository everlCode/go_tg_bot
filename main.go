package main

import (
	"database/sql"
	"go-tg-bot/internal/bot"
	"go-tg-bot/internal/handler"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf(db.Ping().Error())
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing without it")
	}
	db.Close()

	wh := handler.CreateHandler()
	b := bot.NewBot(wh)

	// webhookURL := os.Getenv("WEBHOOK_URL")
	// wh := tgbotapi.NewWebhook(webhookURL)
	// if err != nil {
	// 	log.Fatalf("Ошибка создания вебхука: %v", err)
	// }

	// _, err = bot.SetWebhook(wh)
	// if err != nil {
	// 	log.Fatalf("Ошибка установки вебхука: %v", err)
	// }

	http.HandleFunc("/bot", b.HandleWebHook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
