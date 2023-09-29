package api

import (
	"fmt"
	"net/http"
	"tugas-akhir/internal/repositories/admin"
	"tugas-akhir/internal/repositories/auth"
	"tugas-akhir/internal/repositories/dokter"
	"tugas-akhir/internal/repositories/pasien"
	"tugas-akhir/internal/repositories/petugas"
	"tugas-akhir/internal/repositories/user"
)

type API struct {
	usersRepo   user.UserRepo
	authRepo    auth.AuthRepo
	adminRepo   admin.AdminRepo
	pasienRepo  pasien.PasienRepo
	petugasRepo petugas.PetugasRepo
	dokterRepo  dokter.DokterRepo

	mux *http.ServeMux
}

func NewAPI(usersRepo user.UserRepo, authRepo auth.AuthRepo, adminRepo admin.AdminRepo,
	pasienRepo pasien.PasienRepo, petugasRepo petugas.PetugasRepo, dokterRepo dokter.DokterRepo) API {

	mux := http.NewServeMux()
	api := API{
		usersRepo, authRepo, adminRepo, pasienRepo, petugasRepo, dokterRepo, mux,
	}

	// API without middleware
	mux.Handle("/api/register", api.POST(http.HandlerFunc(api.register)))
	mux.Handle("/api/login", api.POST(http.HandlerFunc(api.login)))
	mux.Handle("/api/logout", api.POST(http.HandlerFunc(api.logout)))

	// API dokter
	mux.Handle("/api/dokter/lihat/jadwal", api.GET(api.AuthMiddleWare(api.DokterMiddleware(http.HandlerFunc(api.lihatJadwalReservasi)))))

	// API pasien with middleware
	mux.Handle("/api/pasien/lihat/poli", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatPoli))))
	mux.Handle("/api/pasien/lihat/poli/dokter", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatDokter))))
	mux.Handle("/api/pasien/lihat/poli/dokter/jadwal", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatJadwalDokter))))
	mux.Handle("/api/pasien/reservasi/pribadi", api.POST(api.AuthMiddleWare(http.HandlerFunc(api.reservasiPribadi))))
	mux.Handle("/api/pasien/reservasi/riwayat", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatReservasi))))
	mux.Handle("/api/pasien/reservasi/hasil", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatHasilReservasi))))
	mux.Handle("/api/pasien/profile", api.GET(api.AuthMiddleWare(http.HandlerFunc(api.lihatDataDiri))))
	mux.Handle("/api/pasien/profile/ubah", api.POST(api.AuthMiddleWare(http.HandlerFunc(api.ubahDataDiri))))

	// API petugas with middleware
	mux.Handle("/api/petugas/insert/poli", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.insertPoli)))))
	mux.Handle("/api/petugas/insert/dokter", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.insertDokter)))))
	mux.Handle("/api/petugas/insert/jadwal", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.insertJadwalDokter)))))
	mux.Handle("/api/petugas/lihat/poli", api.GET(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.lihatPoli)))))
	mux.Handle("/api/petugas/lihat/poli/dokter", api.GET(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.lihatDokter)))))
	mux.Handle("/api/petugas/lihat/poli/dokter/jadwal", api.GET(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.lihatJadwalDokter)))))
	mux.Handle("/api/petugas/reservasi/pasien", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.reservasiPasien)))))
	mux.Handle("/api/petugas/lihat/data/reservasi", api.GET(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.lihatReservasiUser)))))
	mux.Handle("/api/petugas/verifikasi/reservasi", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.verifikasiReservasi)))))
	mux.Handle("/api/petugas/ubah/poli", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.ubahPoli)))))
	mux.Handle("/api/petugas/ubah/jadwal", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.ubahJadwalDokter)))))
	mux.Handle("/api/petugas/laporan", api.POST(api.AuthMiddleWare(api.PetugasMiddleware(http.HandlerFunc(api.kirimDataLaporan)))))

	// API admin with middleware
	mux.Handle("/api/admin/lihat/data/user", api.GET(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.lihatDataUser)))))
	mux.Handle("/api/admin/lihat/data/dokter", api.GET(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.lihatDataDokter)))))
	mux.Handle("/api/admin/insert/petugas", api.POST(api.AuthMiddleWare(api.AdminMiddleware(http.HandlerFunc(api.insertPetugas)))))

	return api
}

func (api *API) Handler() *http.ServeMux {
	return api.mux
}

func (api *API) Start() {
	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", api.Handler())
}
