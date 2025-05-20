package main

import (
	"database/sql"
	"go-tg-bot/internal/bot"
	"io"
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

	bot.NewBot()

	// webhookURL := os.Getenv("WEBHOOK_URL")
	// wh := tgbotapi.NewWebhook(webhookURL)
	// if err != nil {
	// 	log.Fatalf("Ошибка создания вебхука: %v", err)
	// }

	// _, err = bot.SetWebhook(wh)
	// if err != nil {
	// 	log.Fatalf("Ошибка установки вебхука: %v", err)
	// }

	http.HandleFunc("/bot", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Ошибка чтения тела запроса: %v", err)
			http.Error(w, "Ошибка", http.StatusInternalServerError)
			return
		}
		// Логируем тело как строку
		log.Printf("Тело запроса: %s", string(body))
		w.Write([]byte("Hello, world"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
