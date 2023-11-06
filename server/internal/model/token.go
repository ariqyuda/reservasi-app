package model

import "time"

type Token struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Token     string    `db:"token"`
	Fungsi    string    `db:"fungsi"`
	StatusTkn string    `db:"status_tkn"`
	CreatedAt time.Time `db:"created_at"`
}
