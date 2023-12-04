package main

import (
	"database/sql"
	"tugas-akhir/api"
	"tugas-akhir/internal/repositories"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "ariq1:2017@tcp(localhost:3306)/reservasi-app")
	if err != nil {
		panic(err)
	}

	usersRepo := repositories.NewUserRepositories(db)
	dokterRepo := repositories.NewDokterRepositories(db)
	pasienRepo := repositories.NewPasienRepositories(db)
	petugasRepo := repositories.NewPetugasRepositories(db)
	timeRepo := repositories.NewTimeRepositories(db)
	jadwalRepo := repositories.NewJadwalRepositories(db)
	poliRepo := repositories.NewPoliRepositories(db)
	reservasiRepo := repositories.NewReservasiRepositories(db)
	tokenRepo := repositories.NewTokenRepository(db)

	mainAPI := api.NewAPI(*usersRepo, *dokterRepo, *pasienRepo, *petugasRepo, *timeRepo, *jadwalRepo, *poliRepo, *reservasiRepo, *tokenRepo)
	mainAPI.Start()
}
