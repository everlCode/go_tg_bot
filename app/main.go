package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	reaction_repository "go-tg-bot/internal/repository"
	message_repository "go-tg-bot/internal/repository/message"
	reply_repository "go-tg-bot/internal/repository/reply"
	user_repository "go-tg-bot/internal/repository/user"
	dashboard_service "go-tg-bot/internal/services/dashboard"
	message_service "go-tg-bot/internal/services/dashboard/messages"
	"go-tg-bot/internal/services/gigachad"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v4"
)

func main() {
	err := godotenv.Load(".env") // укажи путь до .env
	if err != nil {
		log.Printf("Ошибка загрузки .env файла: %v", err)
	}
	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := user_repository.NewRepository(db)
	replyRepository := reply_repository.NewRepository(db)
	messageRepository := message_repository.NewRepository(db)
	reactionRepository := reaction_repository.NewRepository(db)
	messageService := message_service.NewService(replyRepository, userRepository, messageRepository, reactionRepository)
	// Загружаем переменные окружения
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Client: &http.Client{},
	})
	if err != nil {
		log.Fatal(err)
	}

	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		_, err := db.Exec("UPDATE users SET action = 10")
		if err != nil {
			log.Println(err)
		}
		log.Println("Cron: задание выполнено в", time.Now())
	})

	c.AddFunc("0 22 * * *", func() {
		messages := messageRepository.GetMessagesForToday()

		content := messageService.FormatMessagesForGigaChat(messages)

		gigaChatApi, _ := gigachad.NewApi()
		result := gigaChatApi.Send(content)

		bot.Send(telebot.ChatID(-4204971428), result.Choices[0].Message.Content)
	})
	c.Start()

	// Регистрируем хендлеры
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send(`👋 Добро пожаловать!
	
		Этот бот позволяет оценивать ответы пользователей в беседе и вести рейтинг активности.
		
		📌 Основной функционал:
		— Бот учитывает количество сообщений каждого пользователя.
		— Вы можете реагировать на чужие сообщения с помощью эмодзи ❤️, 👍, 🔥 или 👎, 💩, чтобы изменить "уважение" (рейтинг) автора.
		— Каждому пользователю доступно до 10 действий (оценок) в день.
		— Вся активность отображается на дашборде (веб-интерфейсе).
		
		🌐 Панель статистики: https://www.everl.ru/dashboard
		`)
	})

	mux := http.NewServeMux()

	// Telegram Webhook: используем bot.HandleUpdate
	mux.HandleFunc("/bot", func(w http.ResponseWriter, r *http.Request) {
		log.Println("BOT")
		bodyBytes, _ := io.ReadAll(r.Body)
		log.Println("Raw body:", string(bodyBytes))
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		var update telebot.Update

		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "invalid update: ", http.StatusBadRequest)
			log.Println(err)		log.Println(messages)
			return
		}

		bot.ProcessUpdate(update)
		//костыль, влибе нет обработки реакций
		if update.MessageReaction != nil {
			messageService.HandleReaction(update.MessageReaction)
		}

	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		u := c.Update()
		log.Println(u)
		messageService.Handle(c)

		return nil
	})

	bot.Handle(telebot.OnVideoNote, func(c telebot.Context) error {
		messageService.Handle(c)
		return nil
	})

	// Свой API маршрут
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	})

	mux.HandleFunc("/gigachat", func(w http.ResponseWriter, r *http.Request) {

		messages := messageRepository.GetMessagesForToday()
		if len(messages) == 0 {
			w.Write([]byte("No messages"))
		}

		content := messageService.FormatMessagesForGigaChat(messages)
		log.Println(content)
		gigaChatApi, err := gigachad.NewApi()
		result := gigaChatApi.Send(content)

		if err != nil {
			http.Error(w, "Failed to get API token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		//bot.Send(telebot.ChatID(-4204971428), result.Choices[0].Message.Content)
		// bot.Send(telebot.ChatID(-4204971428), "TEST")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		//Подключение к БД
		dashboardService := dashboard_service.NewService(userRepository)
		users := dashboardService.DashboardData()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	mux.HandleFunc("/api/replies", func(w http.ResponseWriter, r *http.Request) {
		rows := replyRepository.All()
		defer rows.Close()

		var replies []reply_repository.Reply
		for rows.Next() {
			var id, from, to int
			var text string
			if err := rows.Scan(&id, &from, &to, &text); err != nil {
				log.Println("Row scan error:", err)
				return
			}
			replies = append(replies, reply_repository.Reply{
				ID:   id,
				From: from,
				To:   to,
				Text: text,
			})

		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(replies)
	})

	mux.HandleFunc("/setwebhook", func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем webhook URL
		publicURL := os.Getenv("WEBHOOK_URL")
		err = bot.SetWebhook(&telebot.Webhook{
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: publicURL,
			},
			AllowedUpdates: []string{"message", "message_reaction"},
		})
		if err != nil {
			log.Fatal(err)
		}
	})

	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	if env == "production" {
		// Пути к сертификатам
		certFile := os.Getenv("TLS_CERT")
		keyFile := os.Getenv("TLS_KEY")

		if certFile == "" || keyFile == "" {
			log.Fatal("TLS_CERT and TLS_KEY must be set in production")
		}

		log.Println("Starting HTTPS server on port 443")
		log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, mux))
	} else {
		// Локальный режим — обычный HTTP
		log.Println("Running in local mode on http://localhost:" + port)
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}

}
