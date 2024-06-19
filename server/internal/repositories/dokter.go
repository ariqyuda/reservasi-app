package repositories

import (
	"database/sql"
	"errors"
)

type Dokter struct {
	ID        int64  `db:"id"`
	UserID    int64  `db:"user_id"`
	Email     string `db:"email"`
	Nama      string `db:"nama_dokter"`
	STRDokter string `db:"str_dokter"`
	SIPDokter string `db:"sip_dokter"`
	Status    string `db:"status"`
	PoliID    int64  `db:"poli_id"`
	PoliNama  string `db:"nama_poli"`
}

type DokterRepo struct {
	db *sql.DB
}

func NewDokterRepositories(db *sql.DB) *DokterRepo {
	return &DokterRepo{db: db}
}

func (d *DokterRepo) FetchDokter() ([]Dokter, error) {
	//inisiasi variabel
	var dokter []Dokter = make([]Dokter, 0)

	var sqlStmt = `SELECT id, nama FROM dokter`

	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataDokter Dokter
	for rows.Next() {
		err := rows.Scan(
			&dataDokter.ID,
			&dataDokter.Nama,
		)

		if err != nil {
			return nil, err
		}

		dokter = append(dokter, dataDokter)
	}

	return dokter, nil
}

func (d *DokterRepo) InsertDokter(email, nama, password, str_dokter, sip_dokter, status_dokter, poli_nama string) error {
	role := "dokter"
	status := "aktif"

	userRepo := NewUserRepositories(d.db)

	err := userRepo.InsertUser(email, nama, role, password, status)
	if err != nil {
		return err
	}

	idUser, _ := userRepo.FetchUserID(email)

	poliRepo := NewPoliRepositories(d.db)
	poliID, _ := poliRepo.FetchPoliIDByNama(poli_nama)

	sqlStmt := `INSERT INTO dokter (user_id, poli_id, nama, str_dokter, sip_dokter, status, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`

	timeRepo := NewTimeRepositories(d.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err = d.db.Exec(sqlStmt, idUser, poliID, nama, str_dokter, sip_dokter, status_dokter, waktuLokal)
	if err != nil {
		return err
	}

	return err
}

func (d *DokterRepo) FetchDokterByPoliNama(slug string) ([]Dokter, error) {
	var dokter []Dokter = make([]Dokter, 0)

	poliRepo := NewPoliRepositories(d.db)

	poliID, _ := poliRepo.FetchPoliIDBySlug(slug)

	var sqlStmt string = `SELECT d.id, d.nama, p.nama 
	from dokter d
	JOIN poli P ON d.poli_id = p.id 
	where poli_id = ?`

	rows, err := d.db.Query(sqlStmt, poliID)
	if err != nil {
		return nil, errors.New("gagal menampilkan dokter")
	}

	defer rows.Close()

	var dataDokter Dokter
	for rows.Next() {
		err := rows.Scan(
			&dataDokter.ID,
			&dataDokter.Nama,
			&dataDokter.PoliNama,
		)

		if err != nil {
			return nil, err
		}

		dokter = append(dokter, dataDokter)
	}

	return dokter, nil
}

func (u *DokterRepo) FetchDataDokter(page int) ([]Dokter, Pagination, error) {
	var user []Dokter = make([]Dokter, 0)
	var pagination Pagination

	offSet := (page - 1) * 10

	var sqlStmt string = `SELECT u.id, u.email, u.nama, d.id as id_dokter, d.str_dokter, d.sip_dokter, d.status, p.nama
	FROM users u
	JOIN dokter d ON u.id = d.user_id
	JOIN poli p ON d.poli_id = p.id
	LIMIT 10 OFFSET ?`

	rows, err := u.db.Query(sqlStmt, offSet)
	if err != nil {
		return nil, pagination, errors.New("gagal menampilkan data dokter")
	}

	defer rows.Close()

	var dataDokter Dokter
	for rows.Next() {
		err := rows.Scan(
			&dataDokter.UserID,
			&dataDokter.Email,
			&dataDokter.Nama,
			&dataDokter.ID,
			&dataDokter.STRDokter,
			&dataDokter.SIPDokter,
			&dataDokter.Status,
			&dataDokter.PoliNama,
		)

		if err != nil {
			return nil, pagination, err
		}

		user = append(user, dataDokter)
	}

	var sqlStmtCount = `SELECT COUNT(*) FROM users u JOIN dokter d ON u.id = d.user_id`

	row := u.db.QueryRow(sqlStmtCount)

	var totalRows int
	err = row.Scan(&totalRows)

	if err != nil {
		return nil, pagination, errors.New("gagal menghitung jumlah data")
	}

	pagination = GetDataPageInfo(page, 10, totalRows)

	return user, pagination, nil
}

func (d *DokterRepo) FetchDokterByID(dokter_id int64) (Dokter, error) {
	var dokter Dokter

	var sqlStmt string = `SELECT id, nama, poli_id FROM dokter WHERE id = ?`

	row := d.db.QueryRow(sqlStmt, dokter_id)
	err := row.Scan(
		&dokter.ID,
		&dokter.Nama,
		&dokter.PoliID,
	)

	return dokter, err
}

func (d *DokterRepo) FetchDokterIDByJadwalID(jadwal_id int64) (int64, error) {
	var id int64

	var sqlStmt string = `SELECT dokter_id FROM jadwal_dokter WHERE id = ?`

	row := d.db.QueryRow(sqlStmt, jadwal_id)
	err := row.Scan(&id)

	return id, err
}

func (d *DokterRepo) UbahStatusDokter(dokter_id int64, status string) error {
	var sqlStmt string = `UPDATE dokter SET status = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(d.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := d.db.Exec(sqlStmt, status, waktuLokal, dokter_id)

	if err != nil {
		return err
	}

	return err
}

func (d *DokterRepo) UbahDataDokter(dokter_id int64, nama, str_dokter, sip_dokter string) error {
	var sqlStmt string = `UPDATE dokter SET nama = ?, str_dokter = ?, sip_dokter = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(d.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := d.db.Exec(sqlStmt, nama, str_dokter, sip_dokter, waktuLokal, dokter_id)

	if err != nil {
		return err
	}

	return err
}
