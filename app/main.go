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

	"github.com/joho/godotenv"
	"golang.org/x/crypto/acme/autocert"
)

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	MessageCount int    `json:"message_count"`
}

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
		rows := userRepository.GetTopUsers()
		defer rows.Close()

		var users []User
		for rows.Next() {
			var id int
			var name string
			var message_count int
			if err := rows.Scan(&id, &name, &message_count); err != nil {
				http.Error(w, "Row scan error", http.StatusInternalServerError)
				log.Println("Row scan error:", err)
				return
			}
			users = append(users, User{ID: id, Name: name, MessageCount: message_count})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		log.Printf("Returned %d users\n", len(users))
	})

	// Определяем окружение
	env := os.Getenv("ENV") // "production" или "local"
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if env == "production" {
		domain := os.Getenv("DOMAIN")
		if domain == "" {
			log.Fatal("DOMAIN must be set in production")
		}

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
			Cache:      autocert.DirCache("certs"), // папка для хранения сертификатов
		}

		server := &http.Server{
			Addr:      ":443",
			Handler:   nil,
			TLSConfig: certManager.TLSConfig(),
		}

		// Запускаем HTTPS сервер
		go http.ListenAndServe(":80", certManager.HTTPHandler(nil)) // HTTP для получения сертификата
		log.Println("Running in production mode on https://", domain)
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		// Локальный режим — обычный HTTP
		log.Println("Running in local mode on http://localhost:" + port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}
}
