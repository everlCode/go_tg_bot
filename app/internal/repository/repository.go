package repository

import (
	"database/sql"

	message_repository "go-tg-bot/internal/repository/message"
	reaction_repository "go-tg-bot/internal/repository/reaction"
	reply_repository "go-tg-bot/internal/repository/reply"
	report_repository "go-tg-bot/internal/repository/report"
	user_repository "go-tg-bot/internal/repository/user"
)

type Repositories struct {
	User     *user_repository.Repository
	Message  *message_repository.MessageRepository
	Reply    *reply_repository.ReplyRepository
	Reaction *reaction_repository.ReactionRepository
	Report   *report_repository.ReportRepository
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		User:     user_repository.NewRepository(db),
		Message:  message_repository.NewRepository(db),
		Reply:    reply_repository.NewRepository(db),
		Reaction: reaction_repository.NewRepository(db),
		Report:   report_repository.NewRepository(db),
	}
}
