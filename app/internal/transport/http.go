package transport

import (
	"context"
	"encoding/json"
	"go-tg-bot/internal/config"
	"go-tg-bot/internal/services"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type HTTP struct {
	cfg      *config.Config
	services *services.Services
	tg       *Telegram
	log      *slog.Logger
}

func NewHTTP(cfg *config.Config, services *services.Services, tg *Telegram, log *slog.Logger) *HTTP {
	return &HTTP{cfg, services, tg, log}
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
