package main

import (
	"crud-alumni/app/route"
	"crud-alumni/database"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// koneksi database
	database.ConnectDB()
	defer database.DB.Close()

	// inisialisasi Fiber
	app := fiber.New()

	// prefix API
	api := app.Group("/api")

	// route auth
	route.AuthRoutes(api)

	// route alumni & pekerjaan (sudah ada middleware di dalamnya)
	route.AlumniRoutes(api)
	route.PekerjaanRoutes(api)

	// listen server
	log.Fatal(app.Listen(":3000"))
}
