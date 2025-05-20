package main

import (
	"log"
	"net/http"
	"os"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing without it")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	// webhookURL := os.Getenv("WEBHOOK_URL")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Бот запущен: %s", bot.Self.UserName)

	// wh := tgbotapi.NewWebhook(webhookURL)
	// if err != nil {
	// 	log.Fatalf("Ошибка создания вебхука: %v", err)
	// }

	// _, err = bot.SetWebhook(wh)
	if err != nil {
		log.Fatalf("Ошибка установки вебхука: %v", err)
	}

	http.HandleFunc("/bot", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
