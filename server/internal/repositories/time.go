package repositories

import (
	"database/sql"
	"time"
)

type TimeRepo struct {
	db *sql.DB
}

func NewTimeRepositories(db *sql.DB) *TimeRepo {
	return &TimeRepo{db: db}
}

func (u *TimeRepo) SetLocalTime() (string, error) {
	// set lokasi
	location, _ := time.LoadLocation("Asia/Jakarta")

	// get waktu lokal
	localTime := time.Now().In(location).Format("2006-01-02 15:04:05")

	return localTime, nil
}
