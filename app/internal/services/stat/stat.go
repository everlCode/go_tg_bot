package stat_service

import (
	"database/sql"
	reaction_repository "go-tg-bot/internal/repository"
	message_repository "go-tg-bot/internal/repository/message"
	user_repository "go-tg-bot/internal/repository/user"
	"sort"
)

type StatService struct {
	db                *sql.DB
	messageRepository message_repository.MessageRepository
	userRepository user_repository.UserRepository
	reactionRepository reaction_repository.ReactionRepository
}

type WeekStat struct {
	Stats []UserStat
}

type UserStat struct {
	UserId       int
	UserName string
	MessageCount int
	ReactionStat reaction_repository.ReactionStat
}

func NewService(
	db *sql.DB,
	messageRepository message_repository.MessageRepository,
	userRepository user_repository.UserRepository,
	reactionRepository reaction_repository.ReactionRepository,
	) StatService {
	return StatService{
		db:                db,
		messageRepository: messageRepository,
		userRepository: userRepository,
		reactionRepository: reactionRepository,
	}
}

func (service StatService) WeekStat() WeekStat {
	users := service.userRepository.All()
	messageCount := service.messageRepository.MessageCountForWeek()
	reactionStat := service.reactionRepository.GetReactionsCountByMessageAuthor()

	var stat []UserStat

	for _, user := range users {
		count, ok := messageCount[user.ID]; if !ok { continue }
		userStat := UserStat{
			UserId: user.ID,
			UserName: user.Name,
			MessageCount: count,
			ReactionStat: reactionStat[user.ID],
		}
		stat = append(stat, userStat)
	}

	sort.Slice(stat, func(i, j int) bool {
		return stat[i].MessageCount > stat[j].MessageCount
	})

	return WeekStat {
		Stats: stat,
	}
}
