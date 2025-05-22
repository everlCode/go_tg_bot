package bot

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler func(tgbotapi.Update)
}

func NewBot(h func(update tgbotapi.Update)) *Bot {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	b, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot := &Bot{
		api:     b,
		handler: h,
	}

	return bot
}

func (b *Bot) HandleWebHook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var update tgbotapi.Update
	json.Unmarshal(body, &update)

	go b.handler(update) // async
	w.WriteHeader(http.StatusOK)
}
