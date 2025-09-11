package transport

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"

	"go-tg-bot/internal/services"

	"gopkg.in/telebot.v4"
)

type Telegram struct {
	bot      *telebot.Bot
	services *services.Services
	log      *slog.Logger
}

func NewTelegram(token string, services *services.Services, log *slog.Logger) (*Telegram, error) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		Client: &http.Client{},
	})
	if err != nil {
		return nil, err
	}

	tg := &Telegram{
		bot:      bot,
		services: services,
		log:      log,
	}

	tg.registerHandlers()

	return tg, nil
}

func (t *Telegram) registerHandlers() {
	// Команда /start
	t.bot.Handle("/start", func(c telebot.Context) error {
		return c.Send(`👋 Добро пожаловать!
		
Этот бот позволяет оценивать ответы пользователей и вести рейтинг активности.
		
📌 Основной функционал:
— Счётчик сообщений.
— Реакции на чужие сообщения (❤️, 👍, 🔥, 🥰, ❤️‍🔥, 👎, 💩).
— Лимит 10 действий в день.
— Веб-панель статистики: https://www.everl.ru/dashboard
`)
	})

	// Текстовые сообщения
	t.bot.Handle(telebot.OnText, func(c telebot.Context) error {
		t.services.Messages.Handle(c)
		return nil
	})

	// Медиа
	t.bot.Handle(telebot.OnMedia, func(c telebot.Context) error {
		t.services.Messages.Handle(c)
		return nil
	})

	// Топ пользователей
	t.bot.Handle("/top", func(c telebot.Context) error {
		top := t.services.Dashboard.UsersTop()
		return c.Send(top, telebot.ModeHTML)
	})
}

// Webhook HTTP endpoint (для mux.HandleFunc("/bot", ...))
func (t *Telegram) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var update telebot.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "55555", http.StatusBadRequest)
		t.log.Error("telegram webhook decode error", "err", err)
		return
	}
	log.Println("333333333")
	t.bot.ProcessUpdate(update)

	// костыль: в либе нет обработки реакций
	if update.MessageReaction != nil {
		t.services.Messages.HandleReaction(update.MessageReaction)
	}
}

// Запуск (для long polling, если без webhook)
func (t *Telegram) Run() {
	t.log.Info("starting telegram bot (long polling)")
	t.bot.Start()
}

// Доступ к bot (например, для отправки сообщений из cron)
func (t *Telegram) Bot() *telebot.Bot {
	return t.bot
}
