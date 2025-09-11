package main

import (
	"context"
	"go-tg-bot/internal/app"
	"go-tg-bot/internal/config"
	"go-tg-bot/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	application, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Error("failed to init app", "err", err)
		os.Exit(1)
	}

	if err := application.Run(ctx); err != nil {
		log.Error("app stopped with error", "err", err)
	}
}
