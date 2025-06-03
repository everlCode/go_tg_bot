package user_repository

import (
	"database/sql"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	MessageCount int    `json:"message_count"`
	Percent 	 float64    `json:"percent"`
}

func NewRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (ur *UserRepository) UserExist(telegram_id int64) bool {
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

func (ur *UserRepository) CreateUser(telegram_id int64, name string, message_count int) {
	_, err := ur.db.Exec("INSERT INTO users (name, telegram_id, message_count) VALUES (?, ?, ?)", name, telegram_id, message_count)
	if err != nil {
		log.Print(err)
	}
}

func (ur *UserRepository) AddUserMessageCount(telegram_id int64) {
	_, err := ur.db.Exec("UPDATE users SET message_count = message_count + 1 WHERE telegram_id = ?", telegram_id)
	if err != nil {
		log.Print(err)
	}
}

func (ur *UserRepository) GetTopUsers() *sql.Rows {
	rows, err := ur.db.Query(`
		SELECT id, name, message_count, ROUND((message_count * 100.0) / total.total_messages, 2) AS percent
		FROM users,
    	(SELECT SUM(message_count) AS total_messages FROM users) AS total
		ORDER BY message_count DESC;
	`)
	if err != nil {
		log.Println("error:", err)
	}

	return rows
}



