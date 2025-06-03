package main

import (
	"database/sql"
	"encoding/json"
	user_repository "go-tg-bot/internal/repository"
	dashboard_service "go-tg-bot/internal/services"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/telebot.v4"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := user_repository.NewRepository(db)
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
		return c.Send("👋 Hello from telebot.v4 webhook!")
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
