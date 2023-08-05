package user

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"
	"tugas-akhir/internal/repositories/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepositories(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) FetchUserEmail(email string) (model.User, error) {
	var sqlStmt string
	var user model.User

	sqlStmt = `SELECT email FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&user.Email)

	return user, err
}

func (u *UserRepo) FetchPasienByNIK(nik string) (model.Pasien, error) {
	var sqlStmt string
	var pasien model.Pasien

	sqlStmt = `SELECT nik_pasien FROM pasien WHERE nik_pasien = ?`

	row := u.db.QueryRow(sqlStmt, nik)
	err := row.Scan(&pasien.NIK)

	return pasien, err
}

func (u *UserRepo) FetchUserRole(email string) (*string, error) {
	var sqlStmt string
	var role string

	// query untuk mengambil role user berdasarkan email
	sqlStmt = `SELECT role FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&role)

	return &role, err
}

func (u *UserRepo) FetchUserID(email string) (*int64, error) {
	var sqlStmt string
	var id int64

	// query untuk mengambil id user berdasarkan email
	sqlStmt = `SELECT id FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&id)

	return &id, err
}

func (u *UserRepo) FetchDataUserByRole(user_role string) ([]model.User, error) {
	var user []model.User = make([]model.User, 0)

	var sqlStmt string = `SELECT id, email, nama FROM users WHERE role = ?`

	rows, err := u.db.Query(sqlStmt, user_role)
	if err != nil {
		return nil, errors.New("gagal menampilkan data user")
	}

	defer rows.Close()

	var dataUser model.User
	for rows.Next() {
		err := rows.Scan(
			&dataUser.ID,
			&dataUser.Email,
			&dataUser.Name,
		)

		if err != nil {
			return nil, err
		}

		user = append(user, dataUser)
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UserRepo) FetchPoliID(nama string) (*int64, error) {
	var sqlStmt string
	var id int64

	sqlStmt = `SELECT id FROM poli WHERE nama = ?`

	row := u.db.QueryRow(sqlStmt, nama)
	err := row.Scan(&id)

	return &id, err
}

func (u *UserRepo) InsertUser(email, nama, role, password string) error {
	var sqlStmt string

	// email checking
	// check format email
	checkEmailFormat := strings.Contains(email, "@gmail.com")
	if !checkEmailFormat {
		return errors.New("format email tidak sesuai")
	}

	checkEmailSpace := strings.Contains(email, " ")
	if checkEmailSpace {
		return errors.New("email tidak boleh mengandung spasi")
	}

	// check email duplikat
	user, _ := u.FetchUserEmail(email)
	if user.Email != "" {
		return errors.New("email telah terdaftar")
	}

	// password checking
	// check panjang password
	if len(password) < 6 || len(password) > 12 {
		return errors.New("password harus memilik 6-12 karakter")
	}

	// check password mengandung spasi
	checkSpacePass := strings.Contains(password, " ")
	if checkSpacePass {
		return errors.New("password tidak boleh mengandung spasi")
	}

	// check password mengandung hanya huruf dan angka
	checkAlphaNumericPass := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(password)
	if !checkAlphaNumericPass {
		return errors.New("password hanya boleh mengandung huruf dan angka saja")
	}

	newFormatPassword := "%+%" + password + "%+%"

	// hash password
	hashPassword, _ := HashPassword(newFormatPassword)

	// insert to database
	// query untuk insert user
	sqlStmt = `INSERT INTO USERS (email, nama, password, role, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err := u.db.Exec(sqlStmt, email, nama, hashPassword, role, time.Now())
	if err != nil {
		return err
	}

	return err
}
