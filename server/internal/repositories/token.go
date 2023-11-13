package repositories

import (
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"tugas-akhir/internal/model"

	"gopkg.in/gomail.v2"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func GenerateToken() string {

	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	token := make([]byte, 6)
	for i := range token {
		token[i] = letters[rand.Intn(len(letters))]
	}

	return string(token)
}

func (tkn *TokenRepo) CreateToken(id_user int64, fungsi string) (string, error) {

	var token = GenerateToken()

	var status = "not used"

	timeRepo := NewTimeRepositories(tkn.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	var sqlStmt = `INSERT INTO token (user_id, token, fungsi, status_tkn, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err := tkn.db.Exec(sqlStmt, id_user, token, fungsi, status, waktuLokal)
	if err != nil {
		return "", err
	}

	return token, err
}

func (tkn *TokenRepo) SendTokenForgetPassword(email string) error {

	userRepo := NewUserRepositories(tkn.db)

	user_id, err := userRepo.FetchUserID(email)

	if err != nil {
		return err
	}

	token, err := tkn.CreateToken(int64(*user_id), "forget password")

	if err != nil {
		return err
	}

	userId := strconv.FormatInt(int64(*user_id), 10)

	// kirim email token ke user
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "API Test")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Token Reset Password")

	// kirim token reset password
	mailer.SetBody("text/html", "Klik link berikut untuk reset password anda <a href='http://localhost:3000/reset-password/?userid="+userId+"&token="+token+"'>Reset Password</a>")

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"apitest481@gmail.com",
		"ewiv ryzn xevw jwgr",
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}

func (tkn *TokenRepo) TokenForgetPassword(id_user int64) (string, error) {

	var sqlStmt = `SELECT t.token, t.status_tkn
	FROM token t
	JOIN users u ON u.id = t.user_id
	WHERE t.user_id = ? AND t.fungsi = "forget password"
	ORDER BY t.created_at DESC`

	row := tkn.db.QueryRow(sqlStmt, id_user)

	var token string
	var tokenModel model.Token
	err := row.Scan(
		&tokenModel.Token,
		&tokenModel.StatusTkn,
	)

	if err != nil {
		return token, err
	}

	if tokenModel.StatusTkn == "used" {
		return token, errors.New("tidak ada token yang tersedia")
	}

	token = tokenModel.Token

	return token, nil
}

func (tkn *TokenRepo) SendEmailActivation(email string) error {

	userRepo := NewUserRepositories(tkn.db)

	user_id, err := userRepo.FetchUserID(email)

	if err != nil {
		return err
	}

	token, err := tkn.CreateToken(int64(*user_id), "email activation")

	if err != nil {
		return err
	}

	// int64 to string
	userId := strconv.FormatInt(int64(*user_id), 10)

	// kirim email token ke user
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "apitest481@gmail.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Token Aktivasi Email")

	// kirim link aktivasi email
	mailer.SetBody("text/html", "Klik link berikut untuk aktivasi email anda <a href='http://localhost:3000/verify/?userid="+userId+"&token="+token+"'>Aktivasi Email</a>")

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"apitest481@gmail.com",
		"ewiv ryzn xevw jwgr",
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}

func (tkn *TokenRepo) TokenEmailActivation(id_user int64) (string, error) {

	var sqlStmt = `SELECT t.token, t.status_tkn
	FROM token t
	JOIN users u ON u.id = t.user_id
	WHERE t.user_id = ? AND t.fungsi = "email activation"
	ORDER BY t.created_at DESC`

	row := tkn.db.QueryRow(sqlStmt, id_user)

	var token string
	var tokenModel model.Token
	err := row.Scan(
		&tokenModel.Token,
		&tokenModel.StatusTkn,
	)

	if err != nil {
		return token, err
	}

	if tokenModel.StatusTkn == "used" {
		return token, errors.New("tidak ada token yang tersedia")
	}

	token = tokenModel.Token

	return token, nil
}

func (tkn *TokenRepo) ChangeStatusToken(id_user int64, token string) error {

	var sqlStmt = `UPDATE token SET status_tkn = "used" WHERE user_id = ? AND token = ?`

	_, err := tkn.db.Exec(sqlStmt, id_user, token)
	if err != nil {
		return err
	}

	return nil
}
