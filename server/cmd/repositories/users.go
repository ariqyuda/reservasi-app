package repositories

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUserRepositories(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (u *UsersRepo) FetchUserEmail(email string) (User, error) {
	var sqlStmt string
	var user User

	sqlStmt = `SELECT email FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&user.Email)

	return user, err
}

func (u *UsersRepo) FetchPasienNIK(nik string) (Pasien, error) {
	var sqlStmt string
	var pasien Pasien

	sqlStmt = `SELECT nik_pasien FROM pasien WHERE nik_pasien = ?`

	row := u.db.QueryRow(sqlStmt, nik)
	err := row.Scan(&pasien.NIK)

	return pasien, err
}

func (u *UsersRepo) FetchUserRole(email string) (*string, error) {
	var sqlStmt string
	var role string

	// query untuk mengambil role user berdasarkan email
	sqlStmt = `SELECT role FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&role)

	return &role, err
}

func (u *UsersRepo) FetchUserID(email string) (*int64, error) {
	var sqlStmt string
	var id int64

	// query untuk mengambil id user berdasarkan email
	sqlStmt = `SELECT id FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&id)

	return &id, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UsersRepo) FetchPoliID(nama string) (*int64, error) {
	var sqlStmt string
	var id int64

	sqlStmt = `SELECT id FROM poli WHERE nama = ?`

	row := u.db.QueryRow(sqlStmt, nama)
	err := row.Scan(&id)

	return &id, err
}

func (u *UsersRepo) InsertUser(email, nama, role, password string) error {
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

func (u *UsersRepo) Register(email, nama, password, nik, gender, tgl_lahir, tmpt_lahir, alamat, no_hp string) error {

	role := "pasien"

	pasien, _ := u.FetchPasienNIK(nik)
	if pasien.NIK != "" {
		return errors.New("NIK telah terdaftar")
	}

	err := u.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	idUser, _ := u.FetchUserID(email)

	sqlStmt := `INSERT INTO pasien (user_id, nik_pasien, nama, jk_pasien, tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = u.db.Exec(sqlStmt, idUser, nik, nama, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, time.Now())
	if err != nil {
		return err
	}

	return err
}

func (u *UsersRepo) Login(email, password string) (*string, error) {
	// query untuk mengambil data user berdasarkan email dan password
	sqlStmt := `SELECT id, email, nama, password, role FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)

	var user User
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
