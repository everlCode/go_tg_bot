package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"go-tg-bot/internal/config"
	"go-tg-bot/internal/repository"
	"go-tg-bot/internal/services"
	"log"
	"log/slog"
	"net/http"
	"os"

	"gopkg.in/telebot.v4"
)

type HTTP struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *services.Services
	tg           *Telegram
	log          *slog.Logger
}

func NewHTTP(cfg *config.Config, repositories *repository.Repositories, services *services.Services, tg *Telegram, log *slog.Logger) *HTTP {
	return &HTTP{cfg, repositories, services, tg, log}
}

func (h *HTTP) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	log.Println("register /bot")
	mux.HandleFunc("/bot", h.tg.HandleWebhook)

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	})

	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		users := h.services.Dashboard.DashboardData()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		messages := h.repositories.Message.GetMessagesForToday()
		if len(messages) == 0 {
			log.Println("Сегодня сообщений нет")
			return
		}
		content := h.services.Messages.FormatMessagesForGigaChat(messages)

		result := h.services.GigaChat.Send(content)

		txt := result.Choices[0].Message.Content

		h.repositories.Report.Create(txt)

		imgData, e := h.services.GigaChat.GenerateImage("Нарисуй изображение по отчету из нашено чата: " + txt)

		bot := h.tg.Bot()
		if e != nil || len(imgData) == 0 {
			// Если нет изображения — отправляем только текст
			_, e := bot.Send(telebot.ChatID(-4204971428), txt)
			if e != nil {
				log.Println("Ошибка при отправке текста:", e)
			}
			log.Println("DONE (только текст)")
			return
		}

		photo := &telebot.Photo{
			File:    telebot.FromReader(bytes.NewReader(imgData)),
			Caption: txt,
		}
		log.Println("DONE!!!!!!!!!")
		_, er := bot.Send(telebot.ChatID(-4204971428), photo)
		if er != nil {
			log.Fatal("Ошибка получения последнего отчета:", er)
		}
	})

	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	if env == "production" {
		certFile := os.Getenv("TLS_CERT")
		keyFile := os.Getenv("TLS_KEY")

		if certFile == "" || keyFile == "" {
			log.Fatal("TLS_CERT and TLS_KEY must be set in production")
		}

		log.Println("Starting HTTPS server on port 443")
		return http.ListenAndServeTLS(":443", certFile, keyFile, mux)
	} else {
		// Локальный режим — обычный HTTP
		log.Println("Running in local mode on http://localhost:" + port)
		return http.ListenAndServe(":"+port, mux)
	}
}
