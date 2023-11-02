package model

import "time"

type Poli struct {
	ID        int64     `db:"id"`
	Name      string    `db:"nama_poli"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
}
