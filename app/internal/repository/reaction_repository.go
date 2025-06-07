package reaction_repository

import (
	"database/sql"
	"log"
)

type ReactionRepository struct {
	db *sql.DB
}

type Reaction struct {
	ID        int
	UserID    int
	MessageId int
	Text      string
}

func NewRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{
		db: db,
	}
}

func (rr *ReactionRepository) Add(user_id int64, message_id int64, text string) {
	log.Println("CREATE REACTION")
	log.Println(user_id)
	log.Println(message_id)
	log.Println(text)
	_, err := rr.db.Exec("INSERT INTO reactions (user_id, message_id, reaction) VALUES (?, ?, ?)", user_id, message_id, text)
	if err != nil {
		log.Print(err)
	}
}

func (rr *ReactionRepository) All() *sql.Rows {
	rows, err := rr.db.Query(`SELECT * FROM reactions;`)
	if err != nil {
		log.Println("error:", err)
	}

	return rows
}
