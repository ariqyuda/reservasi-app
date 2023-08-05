package admin

import "tugas-akhir/internal/repositories/user"

func (a *AdminRepo) InsertPetugas(email, nama, password string) error {

	role := "petugas"

	userRepo := user.NewUserRepositories(a.db)

	err := userRepo.InsertUser(email, nama, role, password)
	if err != nil {
		return err
	}

	return err
}
