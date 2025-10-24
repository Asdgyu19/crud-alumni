package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB          *sql.DB
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
)

func ConnectDB() {
	// PostgreSQL Connection
	var err error
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=12345678 dbname=alumni_db port=5432 sslmode=disable"
	}
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi database PostgreSQL:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Database PostgreSQL tidak bisa di-ping:", err)
	}
	fmt.Println("✅ Berhasil konek ke PostgreSQL")

	// MongoDB Connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Gagal koneksi ke MongoDB:", err)
	}

	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB tidak bisa di-ping:", err)
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "alumni_db"
	}
	MongoDB = MongoClient.Database(dbName)
	fmt.Println("✅ Berhasil konek ke MongoDB")
}
