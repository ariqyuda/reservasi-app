package repositories

import (
	"database/sql"
	"errors"
)

type PetugasRepo struct {
	db *sql.DB
}

func NewPetugasRepositories(db *sql.DB) *PetugasRepo {
	return &PetugasRepo{db: db}
}

func (prs *PetugasRepo) InsertPetugas(email, nama, password string) error {

	// check input
	if email == "" || nama == "" || password == "" {
		return errors.New("input tidak boleh kosong")
	}

	role := "petugas"
	status := "aktif"

	userRepo := NewUserRepositories(prs.db)

	err := userRepo.InsertUser(email, nama, role, password, status)
	if err != nil {
		return err
	}

	return err
}

func (prs *PetugasRepo) FetchDataPetugas(page int) ([]User, Pagination, error) {
	var user []User = make([]User, 0)
	var pagination Pagination

	offSet := (page - 1) * 10

	var sqlStmt string = `SELECT id, email, nama FROM users WHERE role = 'petugas' LIMIT 10 OFFSET ?`

	rows, err := prs.db.Query(sqlStmt, offSet)
	if err != nil {
		return nil, pagination, errors.New("gagal menampilkan data petugas")
	}

	defer rows.Close()

	var dataUser User
	for rows.Next() {
		err := rows.Scan(
			&dataUser.ID,
			&dataUser.Email,
			&dataUser.Name,
		)

		if err != nil {
			return nil, pagination, err
		}

		user = append(user, dataUser)
	}

	var sqlStmtCount = `SELECT COUNT(*) FROM users WHERE role = 'petugas'`

	row := prs.db.QueryRow(sqlStmtCount)

	var totalRows int
	err = row.Scan(&totalRows)

	if err != nil {
		return nil, pagination, err
	}

	pagination = GetDataPageInfo(page, 10, totalRows)

	return user, pagination, nil
}
