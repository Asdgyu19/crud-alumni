package main

import (
	"crud-alumni/app/route"
	"crud-alumni/database"
	"log"

	"github.com/gofiber/fiber/v2"
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
	route.SetupFileRoutes(api)      // Tambah route untuk upload file
	route.MongoPekerjaanRoutes(api) // Tambah route untuk MongoDB pekerjaan

	// listen server
	log.Fatal(app.Listen(":3000"))
}
