package dashboard_service

import (
	user_repository "go-tg-bot/internal/repository/user"
	"log"
)

type DashboardService struct {
	ur user_repository.UserRepository
}

func NewService(ur user_repository.UserRepository) DashboardService {
	return DashboardService{ur}
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

// func (ds DashboardService) UsersTop() []user_repository.User {
// 	users := ds.DashboardData()

// }
