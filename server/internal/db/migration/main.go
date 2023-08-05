package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "ariq1:2017@tcp(localhost:3306)/reservasi-app")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		nama VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL
	);
	`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS pasien (
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		nik_pasien VARCHAR(255) NOT NULL,
		nama VARCHAR(255) NOT NULL,
		jk_pasien VARCHAR(255) NOT NULL,
		tgl_lahir_pasien VARCHAR(255) NOT NULL,
		tmpt_lahir_pasien VARCHAR(255) NOT NULL,
		alamat_pasien VARCHAR(255) NOT NULL,
		no_hp_pasien VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS poli (
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		nama VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL
	)
	`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS dokter (
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		poli_id INTEGER NOT NULL,
		nama VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (poli_id) REFERENCES poli(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS jadwal_dokter (
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		dokter_id INTEGER NOT NULL,
		jadwal_hari VARCHAR(255) NOT NULL,
		jadwal_waktu VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (dokter_id) REFERENCES dokter(id)
	);
	`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS reservasi (
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		dokter_id INTEGER NOT NULL,
		poli_id INTEGER NOT NULL,
		nik_pasien VARCHAR(255) NOT NULL,
		nama VARCHAR(255) NOT NULL,
		jk_pasien VARCHAR(255) NOT NULL,
		tgl_lahir_pasien VARCHAR(255) NOT NULL,
		tmpt_lahir_pasien VARCHAR(255) NOT NULL,
		alamat_pasien VARCHAR(255) NOT NULL,
		no_hp_pasien VARCHAR(255) NOT NULL,
		jadwal_tanggal VARCHAR(255) NOT NULL,
		jadwal_hari VARCHAR(255) NOT NULL,
		jadwal_waktu VARCHAR(255) NOT NULL,
		tipe VARCHAR(255) NOT NULL,
		status VARCHAR (255) NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (dokter_id) REFERENCES dokter(id),
		FOREIGN KEY (poli_id) REFERENCES poli(id)
	);
	`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS ganti_password (
		id INTEGER AUTO_INCREMENT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		token STRING VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (dokter_id) REFERENCES dokter(id),
		FOREIGN KEY (poli_id) REFERENCES poli(id)
	);
	`)

	if err != nil {
		panic(err)
	}
}
