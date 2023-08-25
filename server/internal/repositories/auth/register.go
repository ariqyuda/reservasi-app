package auth

import (
	"errors"
	"tugas-akhir/internal/repositories/user"
)

func (auth *AuthRepo) Register(email, nama, password, nik, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien string) error {

	role := "pasien"
	userRepo := user.NewUserRepositories(auth.db)

	err := userRepo.CheckUserInput(email, nama, password, nik, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien)
	if err != nil {
		return err
	}

	pasien, _ := userRepo.FetchPasienByNIK(nik)
	if pasien.NIK != "" {
		return errors.New("NIK telah terdaftar")
	}

	err = userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	idUser, _ := userRepo.FetchUserID(email)

	sqlStmt := `INSERT INTO pasien (user_id, nik_pasien, nama, jk_pasien, tgl_lahir_pasien, tmpt_lahir_pasien, alamat_pasien, no_hp_pasien, ktp_pasien, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	waktuLokal, _ := userRepo.SetLocalTime()

	_, err = auth.db.Exec(sqlStmt, idUser, nik, nama, gender, tgl_lahir, tmpt_lahir, alamat, no_hp, ktp_pasien, waktuLokal)
	if err != nil {
		return err
	}

	return err
}
