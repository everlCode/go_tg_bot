package main

import (
	"database/sql"
	"encoding/json"
	"go-tg-bot/internal/bot"
	"go-tg-bot/internal/handler"
	"log"
	"net/http"
	"os"

	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	MessageCount int    `json:"message_count"`
}

func main() {

	// log.Printf(db.Ping().Error())
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing without it")
	}

	wh := handler.CreateHandler()
	b := bot.NewBot(wh)

	// webhookURL := os.Getenv("WEBHOOK_URL")
	// wh := tgbotapi.NewWebhook(webhookURL)
	// if err != nil {
	// 	log.Fatalf("Ошибка создания вебхука: %v", err)
	// }

	// _, err = bot.SetWebhook(wh)
	// if err != nil {
	// 	log.Fatalf("Ошибка установки вебхука: %v", err)
	// }

	http.HandleFunc("/bot", b.HandleWebHook)
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "./db/db.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("select id, name, message_count from users ORDER BY message_count DESC")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var id int
		var name string
		var message_count int

		var users []User
		for rows.Next() {
			if err := rows.Scan(&id, &name, &message_count); err != nil {
				http.Error(w, "Row scan error", http.StatusInternalServerError)
				log.Println("Row scan error:", err)
				return
			}

			users = append(users, User{
				ID:           id,
				Name:         name,
				MessageCount: message_count,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))

}
