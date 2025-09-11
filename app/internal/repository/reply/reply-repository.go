package reply_repository

import (
	"database/sql"
	"log"
)

type ReplyRepository struct {
	db *sql.DB
}

type Reply struct {
	ID int
	From int
	To int
	Text string
}

func NewRepository(db *sql.DB) *ReplyRepository {
	return &ReplyRepository{
		db: db,
	}
}

func (rr *ReplyRepository) Add(from int64, to int64, text string) {
	_, err := rr.db.Exec("INSERT INTO replies (from_user, to_user, text) VALUES (?, ?, ?)", from, to, text)
	if err != nil {
		log.Print(err)
	}
}

func (rr *ReplyRepository) All() *sql.Rows {
	rows, err := rr.db.Query(`SELECT * FROM replies;`)
	if err != nil {
		log.Println("error:", err)
	}

	return rows
}
