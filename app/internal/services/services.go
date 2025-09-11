package services

import (
	"database/sql"
	"go-tg-bot/internal/repository"
	dashboard_service "go-tg-bot/internal/services/dashboard"
	"go-tg-bot/internal/services/gigachad"
	message_service "go-tg-bot/internal/services/messages"
	stat_service "go-tg-bot/internal/services/stat"
)

type Services struct {
	Dashboard *dashboard_service.DashboardService
	Messages  *message_service.MessageService
	Stats     *stat_service.StatService
	GigaChat *gigachad.GigaChatApi
}

func New(db *sql.DB, r *repository.Repositories) *Services {
	api, _ := gigachad.NewApi()
	return &Services{
		Dashboard: dashboard_service.NewService(*r.User),
		Messages:  message_service.NewService(*r.Reply, *r.User, *r.Message, *r.Reaction),
		Stats:     stat_service.NewService(db, r.Message, r.User, r.Reaction),
		GigaChat: api,
	}
}
