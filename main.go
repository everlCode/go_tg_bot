package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	_ "github.com/mattn/go-sqlite3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/db.db") // Файл будет внутри контейнера
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf(db.Ping().Error())
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing without it")
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram_id INTEGER UNIQUE,
		username TEXT);`)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Бот запущен: %s", bot.Self.UserName)

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
