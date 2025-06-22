package stat_service

import (
	"database/sql"
	message_repository "go-tg-bot/internal/repository/message"
	"sort"
)

type StatService struct {
	db                *sql.DB
	messageRepository message_repository.MessageRepository
}

type WeekStat struct {
	Stats []UserStat
}

type UserStat struct {
	UserId       int
	MessageCount int
}

func NewService(db *sql.DB, messageRepository message_repository.MessageRepository) StatService {
	return StatService{
		db:                db,
		messageRepository: messageRepository,
	}
}

func (service StatService) WeekStat() WeekStat {
	messageCount := service.messageRepository.MessageCountForWeek()

	var stat []UserStat

	for i, v := range messageCount {
		userStat := UserStat{}
		userStat.UserId = i
		userStat.MessageCount = v
		stat = append(stat, userStat)
	}

	sort.Slice(stat, func(i, j int) bool {
		return stat[i].MessageCount > stat[j].MessageCount
	})

	return WeekStat{
		Stats: stat,
	}
}
