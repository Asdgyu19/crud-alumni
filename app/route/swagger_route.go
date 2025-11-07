package route

import (
	"crud-alumni/app/service"

	"github.com/gofiber/fiber/v2"
)

// SwaggerRoutes registers routes to serve Swagger UI and the OpenAPI JSON
func SwaggerRoutes(app fiber.Router) {
	swaggerService := service.NewSwaggerService()

	// Serve the JSON spec
	app.Get("/docs/openapi.json", swaggerService.ServeOpenAPISpec)

	// Serve Swagger UI
	app.Get("/docs", swaggerService.ServeSwaggerUI)
}
