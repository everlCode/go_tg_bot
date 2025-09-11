package dashboard_service

import (
	"fmt"
	user_repository "go-tg-bot/internal/repository/user"
	"log"
	"strings"
)

type DashboardService struct {
	ur user_repository.Repository
}

func NewService(ur user_repository.Repository) *DashboardService {
	return &DashboardService{ur}
}

func (ds DashboardService) DashboardData() []user_repository.User {
	rows := ds.ur.GetTopUsers()
	defer rows.Close()

	var users []user_repository.User
	for rows.Next() {
		var id, message_count, respect, action int
		var percent float64
		var name string
		if err := rows.Scan(&id, &name, &message_count, &percent, &respect, &action); err != nil {
			log.Println("Row scan error:", err)
			return nil
		}
		users = append(users, user_repository.User{
			ID:           id,
			Name:         name,
			MessageCount: message_count,
			Percent:      percent,
			Respect:      respect,
			Action:       action,
		})
	}

	return users
}

func (ds DashboardService) UsersTop() string {
	users := ds.DashboardData()

	var b strings.Builder
	for i, u := range users {
		b.WriteString(fmt.Sprintf(
			"%d. %s — ✉️ %d | 🏅 %d\n",
			i+1, u.Name, u.MessageCount, u.Respect,
		))
	}

	return b.String()
}
