package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing without it")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	webhookURL := os.Getenv("WEBHOOK_URL")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Бот запущен: %s", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatalf("Ошибка создания вебхука: %v", err)
	}
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("Ошибка установки вебхука: %v", err)
	}

	http.HandleFunc("/telegram_webhook", func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update

		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Printf("Ошибка декодирования update: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if update.Message == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("Сообщение от @%s: %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы сказали: "+update.Message.Text)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Ошибка отправки сообщения: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Запуск сервера на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
