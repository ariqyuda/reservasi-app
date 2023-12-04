package api

import (
	"context"
	"encoding/json"

	// "errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func (api *API) AllowOrigin(w http.ResponseWriter, req *http.Request) {
	// localhost:8080 origin mendapat ijin akses
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080, http://localhost:3000")
	// localhost:8080 origin mendapat ijin akses
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// semua method diperbolehkan masuk
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
	// semua header diperbolehkan untuk disisipkan
	w.Header().Set("Access-Control-Allow-Headers", "*")
	// allow cookie
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

func (api *API) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			encoder.Encode(AuthErrorResponse{Error: "Token not valid"})
			return
		}
		tokenSplit := strings.Split(token, " ")
		jwtToken := tokenSplit[1]

		claims := &CLaims{}

		//parse JWT token ke dalam claim
		tkn, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// return unauthorized ketika signature invalid
				w.WriteHeader(http.StatusUnauthorized)
				encoder.Encode(AuthErrorResponse{Error: err.Error()})
				return
			}
			// return bad request ketika field token tidak ada
			w.WriteHeader(http.StatusBadRequest)
			encoder.Encode(AuthErrorResponse{Error: err.Error()})
			return
		}

		//return unauthorized ketika token sudah tidak valid (biasanya karna token expired)
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			encoder.Encode(AuthErrorResponse{Error: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "email", claims.Email)
		ctx = context.WithValue(ctx, "role", claims.Role)
		ctx = context.WithValue(ctx, "status", claims.Status)
		ctx = context.WithValue(ctx, "id", claims.ID)
		ctx = context.WithValue(ctx, "props", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) StatusAKunMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		status := r.Context().Value("status")
		if status != "aktif" {
			w.WriteHeader(http.StatusForbidden)
			encoder.Encode(AuthErrorResponse{Error: "akun tidak aktif"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		role := r.Context().Value("role")
		if role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			encoder.Encode(AuthErrorResponse{Error: "forbidden access"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) DokterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		role := r.Context().Value("role")
		if role != "dokter" {
			w.WriteHeader(http.StatusForbidden)
			encoder.Encode(AuthErrorResponse{Error: "forbidden access"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) PetugasMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		role := r.Context().Value("role")
		if role != "petugas" {
			w.WriteHeader(http.StatusForbidden)
			encoder.Encode(AuthErrorResponse{Error: "forbidden access"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) GET(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			encoder.Encode(AuthErrorResponse{Error: "Need GET Method!"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) POST(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			encoder.Encode(AuthErrorResponse{Error: "Need POST Method!"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) PUT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.AllowOrigin(w, r)
		encoder := json.NewEncoder(w)
		if r.Method != http.MethodPut {
			w.WriteHeader(http.StatusMethodNotAllowed)
			encoder.Encode(AuthErrorResponse{Error: "Need PUT Method!"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
