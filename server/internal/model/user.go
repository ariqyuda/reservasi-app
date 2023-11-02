package model

import "time"

type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"nama"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	Token     string    `db:"token"`
}
