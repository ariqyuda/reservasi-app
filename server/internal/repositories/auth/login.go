package auth

import (
	"errors"
	"tugas-akhir/internal/repositories/model"

	"golang.org/x/crypto/bcrypt"
)

func (u *AuthRepo) Login(email, password string) (*string, error) {
	// query untuk mengambil data user berdasarkan email dan password
	sqlStmt := `SELECT id, email, nama, password, role FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return nil, errors.New("tidak ada user dengan email tersebut")
	}

	newFormatPassword := "%+%" + password + "%+%"

	// check kesesuaian input dan stored password
	hashedPassword := []byte(user.Password)
	pass := []byte(newFormatPassword)
	passwordHash := bcrypt.CompareHashAndPassword(hashedPassword, pass)
	if passwordHash == nil {
		return &user.Email, nil
	}

	return nil, errors.New("email atau password yang dimasukkan salah")
}
