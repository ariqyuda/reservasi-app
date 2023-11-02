package repositories

import (
	"database/sql"
	"errors"
	"strings"
	"tugas-akhir/internal/model"
)

type PoliRepo struct {
	db *sql.DB
}

func NewPoliRepositories(db *sql.DB) *PoliRepo {
	return &PoliRepo{db: db}
}

func (p *PoliRepo) FetchPoliIDByNama(nama_poli string) (*int64, error) {
	var sqlStmt string
	var id int64

	sqlStmt = `SELECT id FROM poli WHERE nama = ?`

	row := p.db.QueryRow(sqlStmt, nama_poli)
	err := row.Scan(&id)

	return &id, err
}

func (p *PoliRepo) FetchPoliIDBySlug(slug string) (*int64, error) {
	var sqlStmt string
	var id int64

	sqlStmt = `SELECT id FROM poli WHERE slug = ?`

	row := p.db.QueryRow(sqlStmt, slug)
	err := row.Scan(&id)

	return &id, err
}

func (p *PoliRepo) FetchPoli() ([]model.Poli, error) {
	var poli []model.Poli = make([]model.Poli, 0)

	var sqlStmt string = `SELECT id, nama, slug FROM poli`

	rows, err := p.db.Query(sqlStmt)
	if err != nil {
		return nil, errors.New("gagal menampilkan poli")
	}

	defer rows.Close()

	var dataPoli model.Poli
	for rows.Next() {
		err := rows.Scan(
			&dataPoli.ID,
			&dataPoli.Name,
			&dataPoli.Slug,
		)

		if err != nil {
			return nil, err
		}

		poli = append(poli, dataPoli)
	}

	return poli, nil
}

func (p *PoliRepo) FetchPoliIDByDokterID(dokter_id int64) (*int64, error) {
	var id int64

	var sqlStmt string = `SELECT poli_id FROM dokter WHERE id = ?`

	row := p.db.QueryRow(sqlStmt, dokter_id)
	err := row.Scan(&id)

	return &id, err
}

func (p *PoliRepo) FetchPoliNameBySlug(slug string) (string, error) {
	var poliName string

	var sqlStmt string = `SELECT nama FROM poli WHERE slug = ?`

	row := p.db.QueryRow(sqlStmt, slug)
	err := row.Scan(&poliName)

	return poliName, err
}

func (p *PoliRepo) FetchPoliByNama(nama_poli string) (model.Poli, error) {
	var poli model.Poli

	var sqlStmt string = `SELECT nama FROM poli WHERE nama LIKE ?`

	row := p.db.QueryRow(sqlStmt, nama_poli)
	err := row.Scan(&poli.Name)

	return poli, err
}

func (p *PoliRepo) SetSlugByName(nama_poli string) (string, error) {
	// membuat kata menjadi huruf kecil
	words := strings.ToLower(nama_poli)

	// menghilangkan tanda koma
	words = strings.Replace(words, ",", "", -1)

	// menggantikan spasi dengan tanda strip
	slug := strings.Replace(words, " ", "-", -1)

	return slug, nil
}

func (p *PoliRepo) InsertPoli(nama_poli string) error {

	poli, _ := p.FetchPoliByNama(nama_poli)
	if poli.Name != "" {
		return errors.New("poli sudah terdaftar")
	}

	slug, _ := p.SetSlugByName(nama_poli)

	sqlStmt := `INSERT INTO POLI (nama, slug, created_at) VALUES (?, ?, ?)`

	timeRepo := NewTimeRepositories(p.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := p.db.Exec(sqlStmt, nama_poli, slug, waktuLokal)
	if err != nil {
		return err
	}

	return err
}

func (p *PoliRepo) UbahPoli(poli_id int64, nama string) error {

	slug, _ := p.SetSlugByName(nama)

	var sqlStmt string = `UPDATE poli SET nama = ?, slug = ?, updated_at = ? WHERE id = ?`

	timeRepo := NewTimeRepositories(p.db)
	waktuLokal, _ := timeRepo.SetLocalTime()

	_, err := p.db.Exec(sqlStmt, nama, slug, waktuLokal, poli_id)

	if err != nil {
		return err
	}

	return err
}
