package transport

import (
	"context"
	"encoding/json"
	"go-tg-bot/internal/config"
	"go-tg-bot/internal/services"
	"log"
	"log/slog"
	"net/http"
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

	srv := &http.Server{
		Addr:    ":" + h.cfg.Port,
		Handler: mux,
	}

	h.log.Info("starting HTTP server", "port", h.cfg.Port)
	return srv.ListenAndServe()
}
