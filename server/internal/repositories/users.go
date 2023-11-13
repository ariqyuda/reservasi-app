package repositories

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"tugas-akhir/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepositories(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) FetchUserEmail(email string) (model.User, error) {
	var user model.User

	var sqlStmt string = `SELECT email FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&user.Email)

	return user, err
}

func (u *UserRepo) FetchPasienByNIK(nik string) (model.Pasien, error) {
	var pasien model.Pasien

	var sqlStmt string = `SELECT nik_pasien FROM pasien WHERE nik_pasien = ?`

	row := u.db.QueryRow(sqlStmt, nik)
	err := row.Scan(&pasien.NIK)

	return pasien, err
}

func (u *UserRepo) FetchUserRole(email string) (*string, error) {
	var role string

	// query untuk mengambil role user berdasarkan email
	var sqlStmt string = `SELECT role FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&role)

	return &role, err
}

func (u *UserRepo) FetchUserID(email string) (*int64, error) {
	var id int64

	// query untuk mengambil id user berdasarkan email
	var sqlStmt string = `SELECT id FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&id)

	return &id, err
}

func (u *UserRepo) FetchAccountStatus(email string) (*string, error) {
	var status string

	// query untuk mengambil role user berdasarkan email
	var sqlStmt string = `SELECT status FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&status)

	return &status, err
}

func (u *UserRepo) FetchPasienID(user_id int64) (*int64, error) {
	var sqlStmt string = `SELECT id FROM pasien WHERE user_id = ?`
	var id int64

	row := u.db.QueryRow(sqlStmt, user_id)
	err := row.Scan(&id)

	return &id, err
}

func (u *UserRepo) FetchDataDokter() ([]model.Dokter, error) {
	var user []model.Dokter = make([]model.Dokter, 0)

	var sqlStmt string = `SELECT u.id, u.email, u.nama, d.str_dokter, d.sip_dokter ,p.nama
	FROM users u
	JOIN dokter d ON u.id = d.user_id
	JOIN poli p ON d.poli_id = p.id`

	rows, err := u.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan data user")
	}

	defer rows.Close()

	var dataDokter model.Dokter
	for rows.Next() {
		err := rows.Scan(
			&dataDokter.UserID,
			&dataDokter.Email,
			&dataDokter.Nama,
			&dataDokter.STRDokter,
			&dataDokter.SIPDokter,
			&dataDokter.PoliName,
		)

		if err != nil {
			return nil, err
		}

		user = append(user, dataDokter)
	}

	return user, nil
}

func (u *UserRepo) FetchDataPasien() ([]model.Pasien, error) {
	var user []model.Pasien = make([]model.Pasien, 0)

	var sqlStmt string = `SELECT u.id, u.email, u.nama, p.nik_pasien
	FROM users u
	JOIN pasien p ON u.id = p.user_id`

	rows, err := u.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan data user")
	}

	defer rows.Close()

	var dataUser model.Pasien
	for rows.Next() {
		err := rows.Scan(
			&dataUser.ID,
			&dataUser.Email,
			&dataUser.Nama,
			&dataUser.NIK,
		)

		if err != nil {
			return nil, err
		}

		user = append(user, dataUser)
	}

	return user, nil
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

	newFormatPassword := "%+%" + password + "%+%"

	bytes, err := bcrypt.GenerateFromPassword([]byte(newFormatPassword), 14)
	return string(bytes), err
}

// func (u *UserRepo) CheckUserInput(email, nama, password, nik, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien string) error {

// 	//check data input
// 	if len(email) < 1 || len(nama) < 1 || len(password) < 1 || len(nik) < 1 || len(gender) < 1 || len(tgl_lahir) < 1 ||
// 		len(tmpt_lahir) < 1 || len(alamat) < 1 || len(no_hp) < 1 || len(ktp_pasien) < 1 {
// 		return errors.New("data tidak boleh kosong")
// 	}

// 	//check format nik
// 	if len(nik) > 16 {
// 		return errors.New("nik tidak boleh lebih dari 16 karakter")
// 	}

// 	checkNIK := regexp.MustCompile(`^[0-9]+$`).MatchString(nik)
// 	if !checkNIK {
// 		return errors.New("nik hanya boleh mengandung angka")
// 	}

// 	//check format no hp
// 	checkNoHP := regexp.MustCompile(`^[0-9]+$`).MatchString(no_hp)
// 	if !checkNoHP {
// 		return errors.New("format no hp salah")
// 	}

// 	return nil
// }

func (u *UserRepo) InsertUser(email, nama, role, password, status string) error {
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

	// hash password
	hashPassword, _ := HashPassword(password)

	// insert to database
	// query untuk insert user
	var sqlStmt = `INSERT INTO USERS (email, nama, password, role, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`

	// set waktu lokal
	timeRepo := NewTimeRepositories(u.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := u.db.Exec(sqlStmt, email, nama, hashPassword, role, status, waktuLokal)
	if err != nil {
		return err
	}

	return err
}
