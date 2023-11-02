package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Email    string `json:"email"`
	Name     string `json:"nama"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Pasien struct {
	Email       string `json:"email"`
	Nama        string `json:"nama"`
	Password    string `json:"password"`
	NIK         string `json:"nik_pasien"`
	Gender      string `json:"jk_pasien"`
	BornDate    string `json:"tgl_lahir_pasien"`
	BornPlace   string `json:"tmpt_lahir_pasien"`
	Adress      string `json:"alamat_pasien"`
	PhoneNumber string `json:"no_hp_pasien"`
	KTP         string `json:"ktp_pasien"`
}

type Password struct {
	PasswordBaru string `json:"password_baru"`
	PasswordLama string `json:"password_lama"`
}

type RegisterSuccessResponse struct {
	Message string `json:"message"`
	Data    Pasien `json:"data_pasien"`
}

type LoginSuccessResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Role    string `json:"role"`
}

type AuthErrorResponse struct {
	Error string `json:"error"`
}

type PasswordResponse struct {
	Message string `json:"message"`
}

// Jwt key untuk membuat signature
var jwtKey = []byte("key")

type CLaims struct {
	ID    int64
	Email string
	Role  string
	jwt.RegisteredClaims
}

func (api *API) register(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)
	var pasien Pasien
	err := json.NewDecoder(req.Body).Decode(&pasien)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.authRepo.Register(pasien.Email, pasien.Nama, pasien.Password, pasien.NIK, pasien.Gender, pasien.BornDate, pasien.BornPlace, pasien.Adress, pasien.PhoneNumber, pasien.KTP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	registerPasien := Pasien{
		Email:       pasien.Email,
		Nama:        pasien.Nama,
		NIK:         pasien.NIK,
		Gender:      pasien.Gender,
		BornDate:    pasien.BornDate,
		BornPlace:   pasien.BornPlace,
		Adress:      pasien.Adress,
		PhoneNumber: pasien.PhoneNumber,
		KTP:         pasien.KTP,
	}

	registerResponse := RegisterSuccessResponse{
		Message: "register success",
		Data:    registerPasien,
	}

	json.NewEncoder(w).Encode(registerResponse)
}

func (api *API) login(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := api.authRepo.Login(user.Email, user.Password)

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	userRole, _ := api.usersRepo.FetchUserRole(*res)
	userId, _ := api.usersRepo.FetchUserID(*res)

	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &CLaims{
		ID:    *userId,
		Email: *res,
		Role:  *userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	})

	loginResponse := LoginSuccessResponse{
		Message: "login success",
		Token:   tokenString,
		Role:    *userRole,
	}

	json.NewEncoder(w).Encode(loginResponse)
}

func (api *API) logout(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	token, err := req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// return unauthorized ketika token kosong
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// return bad request ketika field token tidak ada
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if token.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	c := http.Cookie{
		Name:   "token",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logged out"))
}

func (api *API) ubahPassword(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	var password Password
	err := json.NewDecoder(req.Body).Decode(&password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.authRepo.UbahPassword(userID, password.PasswordLama, password.PasswordBaru)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	authResponse := PasswordResponse{
		Message: "Berhasil mengganti password",
	}

	json.NewEncoder(w).Encode(authResponse)
}

func (api *API) resetPassword(w http.ResponseWriter, req *http.Request) {
	api.AllowOrigin(w, req)

	userID := req.Context().Value("id").(int64)

	var password Password
	err := json.NewDecoder(req.Body).Decode(&password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encoder := json.NewEncoder(w)
	err = api.authRepo.ResetPassword(userID, password.PasswordBaru)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(AuthErrorResponse{Error: err.Error()})
		return
	}

	authResponse := PasswordResponse{
		Message: "Berhasil mengganti password",
	}

	json.NewEncoder(w).Encode(authResponse)
}
