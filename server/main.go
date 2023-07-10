package main

import (
	"database/sql"
	"tugas-akhir/api"
	"tugas-akhir/cmd/repositories"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "ariq1:2017@tcp(localhost:3306)/reservasi-app")
	if err != nil {
		panic(err)
	}

	usersRepo := repositories.NewUserRepositories(db)
	adminRepo := repositories.NewAdminRepositories(db)
	pasienRepo := repositories.NewPasienRepositories(db)
	petugasRepo := repositories.NewPetugasRepositories(db)

	mainAPI := api.NewAPI(*usersRepo, *adminRepo, *pasienRepo, *petugasRepo)
	mainAPI.Start()
}
