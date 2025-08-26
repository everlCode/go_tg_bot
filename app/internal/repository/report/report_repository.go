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

func (mr *ReportRepository) GetLast() (*Report, error) {
	row := mr.db.QueryRow("SELECT id, text, created_at FROM reports ORDER BY created_at DESC LIMIT 1")
	var report Report
	err := row.Scan(&report.ID, &report.Text, &report.CreatedAt)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &report, nil
}
