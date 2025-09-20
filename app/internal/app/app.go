package app

import (
	"context"
	"database/sql"
	"go-tg-bot/internal/config"
	"go-tg-bot/internal/cron"
	"go-tg-bot/internal/repository"
	"go-tg-bot/internal/services"
	"go-tg-bot/internal/transport"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	cfg  *config.Config
	log  *slog.Logger
	db   *sql.DB
	tg   *transport.Telegram
	http *transport.HTTP
	Scheduler *cron.Scheduler
}

func New(ctx context.Context, cfg *config.Config, log *slog.Logger) (*App, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}

	repos := repository.New(db)
	services := services.New(db, repos)

	tg, err := transport.NewTelegram(cfg.TelegramToken, services, log)
	if err != nil {
		return nil, err
	}

	httpSrv := transport.NewHTTP(cfg, repos, services, tg, log)
	scheduler := cron.NewScheduler(services, repos, tg, log)

	return &App{
		cfg:  cfg,
		log:  log,
		db:   db,
		tg:   tg,
		http: httpSrv,
		Scheduler: scheduler,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	go a.tg.Run()
	return a.http.Run(ctx)
}
