package bot

import (
	"log"
	"os"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
}

func NewBot() (*tgbotapi.BotAPI) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Бот запущен: %s", bot.Self.UserName)

	return bot
}

