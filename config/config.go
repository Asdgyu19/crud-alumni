package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// PostgreSQL Config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// MongoDB Config
	MongoURI    string
	MongoDBName string
	UseMongoDB  bool

	// App Config
	JWTSecret string
	AppPort   string
}

var Cfg *Config

func LoadConfig() {
	// load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Tidak menemukan file .env, pakai default env")
	}

	Cfg = &Config{
		// PostgreSQL Config
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "crud_alumni"),

		// MongoDB Config
		MongoURI:    getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGODB_DATABASE", "alumni_db"),
		UseMongoDB:  getEnv("USE_MONGODB", "") == "true",

		// App Config
		JWTSecret: getEnv("JWT_SECRET", "mysecretkey"),
		AppPort:   getEnv("APP_PORT", "3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
