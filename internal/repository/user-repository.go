package user_repository

import (
	"database/sql"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur *UserRepository) UserExist(telegram_id int) bool {
	var exist bool
	err := ur.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = ?)", telegram_id).Scan(&exist)
	if err != nil {
		log.Fatal(err)
		return false
	}

	if exist {
		return true
	}

	return false
}

func (ur *UserRepository) CreateUser(telegram_id int, name string, message_count int) {
	_, err := ur.db.Exec("INSERT INTO users (name, telegram_id, message_count) VALUES (?, ?, ?)", telegram_id, name, message_count)
	if err != nil {
		log.Fatal(err)
	}
}

func (ur *UserRepository) AddUserMessageCount(telegram_id int) {
	_, err := ur.db.Exec("UPDATE users SET message_count = message_count + 1 WHERE telegram_id = ?", telegram_id)
	if err != nil {
		log.Fatal(err)
	}
}
