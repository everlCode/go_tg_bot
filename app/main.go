package main

import (
	"database/sql"
	"encoding/json"
	reply_repository "go-tg-bot/internal/repository/reply"
	user_repository "go-tg-bot/internal/repository/user"
	dashboard_service "go-tg-bot/internal/services/dashboard"
	reply_service "go-tg-bot/internal/services/dashboard/replies"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v4"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	c := cron.New()
	c.AddFunc("0 0 * * *", func() {
		_, err := db.Exec("UPDATE users SET action = 3")
		if err != nil {
			log.Println(err)
		}
		log.Println("Cron: задание выполнено в", time.Now())
	})
	c.Start()

	userRepository := user_repository.NewRepository(db)
	replyRepository := reply_repository.NewRepository(db)
	replyService := reply_service.NewService(replyRepository, userRepository)
	// Загружаем переменные окружения
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Client: &http.Client{},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Устанавливаем webhook URL
	publicURL := os.Getenv("WEBHOOK_URL")
	err = bot.SetWebhook(&telebot.Webhook{
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: publicURL,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Регистрируем хендлеры
	bot.Handle("/start", func(c telebot.Context) error {
		return c.Send(`👋 Добро пожаловать!
	
		Этот бот позволяет оценивать ответы пользователей в беседе и вести рейтинг активности.
		
		📌 Основной функционал:
		— Бот учитывает количество сообщений каждого пользователя.
		— Вы можете отвечать на чужие сообщения знаком **+** или **-**, чтобы изменить "уважение" (рейтинг) автора.
		— Каждому пользователю доступно до 3 действий (оценок) в день.
		— Вся активность отображается на дашборде (веб-интерфейсе).
		
		🌐 Панель статистики: https://www.everl.ru/dashboard
		
		Просто начните общение или поставьте + / - в ответ на сообщение, чтобы участвовать!
		`)
	})

	mux := http.NewServeMux()

	// Telegram Webhook: используем bot.HandleUpdate
	mux.HandleFunc("/bot", func(w http.ResponseWriter, r *http.Request) {
		var update telebot.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "invalid update", http.StatusBadRequest)
			return
		}

		bot.ProcessUpdate(update)
	})

	bot.Handle(telebot.OnText, func(c telebot.Context) error {
		log.Println("text")
		msg := c.Message()
		if msg == nil || msg.Sender == nil {
			log.Print(msg)
			return nil
		}

		id := msg.Sender.ID
		name := msg.Sender.FirstName
		userExist := userRepository.UserExist(id)

		if userExist {
			userRepository.AddUserMessageCount(id)
		} else {
			userRepository.CreateUser(id, name, 1)
		}

		return nil
	})

	bot.Handle(telebot.OnReply, func(c telebot.Context) error {
		replyService.Handle(c)

		return nil
	})

	// Свой API маршрут
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		//Подключение к БД

		dashboardService := dashboard_service.NewService(userRepository)
		users := dashboardService.DashboardData()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	mux.HandleFunc("/api/replies", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT * FROM replies;`)
		if err != nil {
			log.Println("error:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
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
