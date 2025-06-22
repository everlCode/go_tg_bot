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

type ReactionStat struct {
	UserID            int
	GetReactionCount  int
	MadeReactionCount int
}

func NewRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{
		db: db,
	}
}

func (rr *ReactionRepository) Add(user_id int64, message_id int64, text string) {
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

func (repo *ReactionRepository) ReactionStat() map[int]ReactionStat {
	rows, err := repo.db.Query(`
        SELECT u.telegram_id, COUNT(*) AS reaction_count
		FROM reactions r
		INNER JOIN messages m ON m.message_id = r.message_id
		INNER JOIN users u ON u.telegram_id = m.from_user
		WHERE m.send_at >= STRFTIME("%s", "now", "-7 days")
		GROUP BY u.telegram_id
		ORDER BY reaction_count DESC;
    `)
	if err != nil {
		log.Println("error:", err)
	}

	stats := make(map[int]ReactionStat)
	for rows.Next() {
		stat := ReactionStat{}
		rows.Scan(&stat.UserID, &stat.GetReactionCount)
		stats[stat.UserID] = stat
	}
	rows.Close()

	rowws, err := repo.db.Query(`
        SELECT u.telegram_id, COUNT(*) AS reaction_count
		FROM reactions r
		INNER JOIN messages m ON m.message_id = r.message_id
		INNER JOIN users u ON u.telegram_id = r.user_id
		WHERE m.send_at >= STRFTIME("%s", "now", "-7 days")
		GROUP BY u.telegram_id
		ORDER BY reaction_count DESC;
    `)
	if err != nil {
		log.Println("error:", err)
	}
	defer rowws.Close()

	for rowws.Next() {
		var id, count int
		rowws.Scan(&id, &count)

		stat, ok := stats[id]
		if !ok {
			continue
		}
		stat.MadeReactionCount = count
		stats[id] = stat
	}

	return stats
}
