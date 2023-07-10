package api

import (
	"fmt"
	"net/http"
	"tugas-akhir/cmd/repositories"
)

type API struct {
	usersRepo   repositories.UsersRepo
	adminRepo   repositories.AdminRepo
	pasienRepo  repositories.PasienRepo
	petugasRepo repositories.PetugasRepo

	mux *http.ServeMux
}

func NewAPI(usersRepo repositories.UsersRepo, adminRepo repositories.AdminRepo,
	pasienRepo repositories.PasienRepo, petugasRepo repositories.PetugasRepo) API {

	mux := http.NewServeMux()
	api := API{
		usersRepo, adminRepo, pasienRepo, petugasRepo, mux,
	}

	// API without middleware
	mux.Handle("/api/register", api.POST(http.HandlerFunc(api.register)))
	mux.Handle("/api/login", api.POST(http.HandlerFunc(api.login)))
	mux.Handle("/api/logout", api.POST(http.HandlerFunc(api.logout)))
	// mux.Handle("/api/insert", api.POST(http.HandlerFunc(api.insertUser)))

	// API pasien with middleware
	mux.Handle("/api/pasien/reservasi/pribadi", api.POST(api.AuthMiddleWare(http.HandlerFunc(api.reservasiPribadi))))
	mux.Handle("/api/pasien/reservasi/riwayat", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatReservasi))))
	mux.Handle("/api/pasien/reservasi/hasil", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatHasilReservasi))))

	// API admin with middleware
	mux.Handle("/api/admin/insert/dokter", api.POST(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.insertDokter)))))
	mux.Handle("/api/admin/insert/petugas", api.POST(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.insertPetugas)))))
	mux.Handle("/api/admin/insert/admin", api.POST(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.insertAdmin)))))
	mux.Handle("/api/admin/insert/poli", api.POST(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.insertPoli)))))
	mux.Handle("/api/admin/insert/dokter/jadwal", api.POST(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.insertJadwalDokter)))))
	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", api.Handler())
}
