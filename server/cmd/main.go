package main

import (
	"database/sql"
	"tugas-akhir/api"
	"tugas-akhir/internal/repositories/admin"
	"tugas-akhir/internal/repositories/auth"
	"tugas-akhir/internal/repositories/dokter"
	"tugas-akhir/internal/repositories/pasien"
	"tugas-akhir/internal/repositories/petugas"
	"tugas-akhir/internal/repositories/user"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "ariq1:2017@tcp(localhost:3306)/reservasi-app")
	if err != nil {
		panic(err)
	}

	usersRepo := user.NewUserRepositories(db)
	authRepo := auth.NewAuthRepositories(db)
	adminRepo := admin.NewAdminRepositories(db)
	pasienRepo := pasien.NewPasienRepositories(db)
	petugasRepo := petugas.NewPetugasRepositories(db)
	dokterRepo := dokter.NewDokterRepositories(db)

	mainAPI := api.NewAPI(*usersRepo, *authRepo, *adminRepo, *pasienRepo, *petugasRepo, *dokterRepo)
	mainAPI.Start()
}
