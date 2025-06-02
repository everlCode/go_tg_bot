package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go-tg-bot/internal/bot"
	"go-tg-bot/internal/handler"
	user_repository "go-tg-bot/internal/repository"
	dashboard_service "go-tg-bot/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing without it")
	}

	// Подключение к БД
	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация репозитория и бота
	userRepository := user_repository.NewRepository(db)
	webHook := handler.CreateHandler(&userRepository)
	b := bot.NewBot(webHook)

	// Роутинг
	http.HandleFunc("/bot", b.HandleWebHook)
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	})
	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		dashboardService := dashboard_service.NewService(userRepository)
		users := dashboardService.DashboardData()
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	// Определяем окружение
	env := os.Getenv("ENV") // "production" или "local"
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if env == "production" {
		// Пути к сертификатам
		certFile := os.Getenv("TLS_CERT")
		keyFile := os.Getenv("TLS_KEY")

		if certFile == "" || keyFile == "" {
			log.Fatal("TLS_CERT and TLS_KEY must be set in production")
		}

		log.Println("Starting HTTPS server on port 443")
		log.Fatal(http.ListenAndServeTLS(":443", certFile, keyFile, nil))
	} else {
		// Локальный режим — обычный HTTP
		log.Println("Running in local mode on http://localhost:" + port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
