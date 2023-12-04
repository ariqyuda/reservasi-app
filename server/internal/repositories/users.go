package repositories

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Name     string `db:"nama"`
	Password string `db:"password"`
	Role     string `db:"role"`
	Token    string `db:"token"`
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepositories(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) Register(email, nama, password, nik, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien string) error {

	role := "pasien"
	status := "tidak aktif"
	pasienRepo := NewPasienRepositories(u.db)
	userRepo := NewUserRepositories(u.db)

	//check data input
	if len(email) < 1 || len(nama) < 1 || len(password) < 1 || len(nik) < 1 || len(gender) < 1 || len(tgl_lahir) < 1 ||
		len(tmpt_lahir) < 1 || len(alamat) < 1 || len(no_hp) < 1 || len(ktp_pasien) < 1 {
		return errors.New("data tidak boleh kosong")
	}

	//check nik duplikat
	pasien, _ := pasienRepo.FetchPasienByNIK(nik)
	if pasien.NIK != "" {
		return errors.New("NIK telah terdaftar")
	}

	//check format nik
	checkNIK := regexp.MustCompile(`^[0-9]+$`).MatchString(nik)
	if len(nik) > 16 || len(nik) < 16 || !checkNIK {
		return errors.New("format nik salah")
	}

	// //check format no hp
	// checkNoHP := regexp.MustCompile(`^[0-9]+$`).MatchString(no_hp)
	// if !checkNoHP {
	// 	return errors.New("format no hp salah")
	// }

	err := userRepo.InsertUser(email, nama, role, password, status)
	if err != nil {
		return err
	}

	idUser, _ := userRepo.FetchUserID(email)

	sqlStmt := `INSERT INTO pasien (user_id, nik_pasien, nama, jk_pasien, tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, ktp_pasien, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	timeRepo := NewTimeRepositories(u.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err = u.db.Exec(sqlStmt, idUser, nik, nama, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien, waktuLokal)
	if err != nil {
		return err
	}

	tokenRepo := NewTokenRepository(u.db)
	_ = tokenRepo.SendEmailActivation(email)

	return err
}

func (u *UserRepo) Login(email, password string) (*string, error) {
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
	if passwordHash != nil {
		return nil, errors.New("email atau password yang dimasukkan salah")
	}

	return &user.Email, nil
}

func (u *UserRepo) UbahPassword(id_user int64, password_lama, password_baru string) error {

	// query untuk mengambil data user berdasarkan email dan password
	sqlStmt := `SELECT id, email, password FROM users WHERE id = ?`

	row := u.db.QueryRow(sqlStmt, id_user)

	var user User
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

	if passwordHash != nil {
		return errors.New("password lama tidak sesuai")
	}

	// hash password
	hashPassword, _ := HashPassword(password_baru)

	sqlStmt = `UPDATE users SET password = ? WHERE id = ?`

	_, err = u.db.Exec(sqlStmt, hashPassword, id_user)
	if err != nil {
		return err
	}

	return err
}

func (u *UserRepo) ResetPassword(id_user int64, input_token, password_baru string) error {

	tokenRepo := NewTokenRepository(u.db)
	token, err := tokenRepo.TokenForgetPassword(id_user)

	if err != nil {
		return err
	}

	if input_token != token {
		return errors.New("token tidak sesuai")
	}

	// hash password
	hashPassword, _ := HashPassword(password_baru)

	var sqlStmtUpdate string = `UPDATE users SET password = ? WHERE id = ?`

	_, err = u.db.Exec(sqlStmtUpdate, hashPassword, id_user)
	if err != nil {
		return errors.New("gagal mengubah password")
	}

	_ = tokenRepo.ChangeStatusToken(id_user, token)

	return nil
}

func (u *UserRepo) EmailActivation(id_user int64, input_token string) error {

	tokenRepo := NewTokenRepository(u.db)
	token, err := tokenRepo.TokenEmailActivation(id_user)

	if err != nil {
		return err
	}

	if input_token != token {
		return errors.New("token tidak sesuai")
	}

	var sqlStmtUpdate string = `UPDATE users SET status = "aktif" WHERE id = ?`

	_, err = u.db.Exec(sqlStmtUpdate, id_user)
	if err != nil {
		return errors.New("gagal aktivasi akun")
	}

	_ = tokenRepo.ChangeStatusToken(id_user, token)

	return nil
}

func (u *UserRepo) FetchUserEmail(email string) (User, error) {
	var user User

	var sqlStmt string = `SELECT email FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&user.Email)

	return user, err
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

	// query untuk mengambil status user berdasarkan email
	var sqlStmt string = `SELECT status FROM users WHERE email = ?`

	row := u.db.QueryRow(sqlStmt, email)
	err := row.Scan(&status)

	return &status, err
}

func (u *UserRepo) FetchDataUserByRole(user_role string) ([]User, error) {
	var user []User = make([]User, 0)

	var sqlStmt string = `SELECT id, email, nama FROM users WHERE role = ?`

	rows, err := u.db.Query(sqlStmt, user_role)
	if err != nil {
		return nil, errors.New("gagal menampilkan data user")
	}

	defer rows.Close()

	var dataUser User
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
