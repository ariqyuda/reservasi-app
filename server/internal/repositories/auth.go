package repositories

import (
	"database/sql"
	"errors"
	"tugas-akhir/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepositories(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (auth *AuthRepo) Register(email, nama, password, nik, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien string) error {

	role := "pasien"
	userRepo := NewUserRepositories(auth.db)

	err := userRepo.CheckUserInput(email, nama, password, nik, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien)
	if err != nil {
		return err
	}

	pasien, _ := userRepo.FetchPasienByNIK(nik)
	if pasien.NIK != "" {
		return errors.New("NIK telah terdaftar")
	}

	err = userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	idUser, _ := userRepo.FetchUserID(email)

	sqlStmt := `INSERT INTO pasien (user_id, nik_pasien, nama, jk_pasien, tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, ktp_pasien, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	timeRepo := NewTimeRepositories(auth.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err = auth.db.Exec(sqlStmt, idUser, nik, nama, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien, waktuLokal)
	if err != nil {
		return err
	}

	return err
}

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

func (u *AuthRepo) UbahPassword(id_user int64, password_lama, password_baru string) error {

	// query untuk mengambil data user berdasarkan email dan password
	sqlStmt := `SELECT id, email, password FROM users WHERE id = ?`

	row := u.db.QueryRow(sqlStmt, id_user)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil
	}

	newFormatPassword := "%+%" + password_lama + "%+%"

	// check kesesuaian input dan stored password
	hashedPassword := []byte(user.Password)
	pass := []byte(newFormatPassword)
	passwordHash := bcrypt.CompareHashAndPassword(hashedPassword, pass)

	if passwordHash == nil {
		// hash password
		hashPassword, _ := HashPassword(password_baru)

		var sqlStmt string = `UPDATE users SET password = ? WHERE id = ?`

		_, err = u.db.Exec(sqlStmt, hashPassword, id_user)
		if err != nil {
			return err
		}
	}

	return errors.New("password lama tidak sesuai")
}

func (u *AuthRepo) ResetPassword(id_user int64, password_baru string) error {

	// hash password
	hashPassword, _ := HashPassword(password_baru)

	var sqlStmtUpdate string = `UPDATE users SET password = ? WHERE id = ?`

	_, err := u.db.Exec(sqlStmtUpdate, hashPassword, id_user)
	if err != nil {
		return errors.New("gagal mengubah password")
	}

	return nil
}
