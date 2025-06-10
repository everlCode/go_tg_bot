package report_repository

import (
	"database/sql"
	"log"
)

type ReportRepository struct {
	db *sql.DB
}

type Report struct {
	ID        int
	Text      string
	CreatedAt int
}

func NewRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{
		db: db,
	}
}

func (mr *ReportRepository) Create(text string) {
	_, err := mr.db.Exec("INSERT INTO reports (text) VALUES(?)", text)
	if err != nil {
		log.Println(err)
	}
}
