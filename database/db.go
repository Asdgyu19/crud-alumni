package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	// ambil dari env jika perlu
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		// contoh default, ganti password/dbname sesuai kamu
		dsn = "host=localhost user=postgres password=12345678 dbname=alumni_db port=5432 sslmode=disable"
	}
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak bisa di-ping:", err)
	}
	fmt.Println("✅ Berhasil konek ke PostgreSQL")
}
