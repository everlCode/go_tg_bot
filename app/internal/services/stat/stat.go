package stat_service

import (
	"database/sql"
	"fmt"
	reaction_repository "go-tg-bot/internal/repository"
	message_repository "go-tg-bot/internal/repository/message"
	user_repository "go-tg-bot/internal/repository/user"
	"sort"
	"strings"
)

type StatService struct {
	db                 *sql.DB
	messageRepository  message_repository.MessageRepository
	userRepository     user_repository.UserRepository
	reactionRepository reaction_repository.ReactionRepository
}

type WeekStat struct {
	Stats []UserStat
}

type UserStat struct {
	UserId       int
	UserName     string
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
		db:                 db,
		messageRepository:  messageRepository,
		userRepository:     userRepository,
		reactionRepository: reactionRepository,
	}
}

func (service StatService) WeekStat() WeekStat {
	users := service.userRepository.All()
	messageCount := service.messageRepository.MessageCountForWeek()
	reactionStat := service.reactionRepository.ReactionStat()

	var stat []UserStat

	for _, user := range users {
		count, ok := messageCount[user.ID]
		if !ok {
			continue
		}
		userStat := UserStat{
			UserId:       user.ID,
			UserName:     user.Name,
			MessageCount: count,
			ReactionStat: reactionStat[user.ID],
		}
		stat = append(stat, userStat)
	}

	sort.Slice(stat, func(i, j int) bool {
		return stat[i].MessageCount > stat[j].MessageCount
	})

	return WeekStat{
		Stats: stat,
	}
}

func (service StatService) FormatDataForWeekReport(data WeekStat) string {
	if len(data.Stats) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("✨ <b>Итоги недели</b> ✨\n\n")
	sb.WriteString("🏆 <b>Топ сообщений</b>\n\n")

	for i, stat := range data.Stats {
		medal := ""
		switch i {
		case 0:
			medal = "🥇"
		case 1:
			medal = "🥈"
		case 2:
			medal = "🥉"
		default:
			medal = "🔹"
		}
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> — <b>%d</b>\n", medal, stat.UserName, stat.MessageCount))
	}

	sort.Slice(data.Stats, func(i, j int) bool {
		return data.Stats[i].ReactionStat.GetReactionCount > data.Stats[j].ReactionStat.GetReactionCount
	})

	sb.WriteString("\n")
	sb.WriteString("🏆 <b>Топ полученных реакций</b>\n\n")
	for i, stat := range data.Stats {
		medal := ""
		switch i {
		case 0:
			medal = "🥇"
		case 1:
			medal = "🥈"
		case 2:
			medal = "🥉"
		default:
			medal = "🔹"
		}
		
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> — <b>%d</b>\n", medal, stat.UserName, stat.ReactionStat.GetReactionCount))
	}

	sort.Slice(data.Stats, func(i, j int) bool {
		return data.Stats[i].ReactionStat.MadeReactionCount > data.Stats[j].ReactionStat.MadeReactionCount
	})

	sb.WriteString("\n")
	sb.WriteString("🏆 <b>Топ оставленных реакций</b>\n\n")
	for i, stat := range data.Stats {
		medal := ""
		switch i {
		case 0:
			medal = "🥇"
		case 1:
			medal = "🥈"
		case 2:
			medal = "🥉"
		default:
			medal = "🔹"
		}
		
		sb.WriteString(fmt.Sprintf("%s <b>%s</b> — <b>%d</b>\n", medal, stat.UserName, stat.ReactionStat.MadeReactionCount))
	}

	return sb.String()
}
