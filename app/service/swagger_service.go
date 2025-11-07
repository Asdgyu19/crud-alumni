package service

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type SwaggerService struct{}

func NewSwaggerService() *SwaggerService {
	return &SwaggerService{}
}

// ServeOpenAPISpec serves the OpenAPI JSON specification
func (s *SwaggerService) ServeOpenAPISpec(c *fiber.Ctx) error {
	b, err := ioutil.ReadFile("docs/openapi.json")
	if err != nil {
		log.Println("failed to read openapi.json:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	c.Type("json")
	return c.Send(b)
}

// ServeSwaggerUI serves the Swagger UI HTML page
func (s *SwaggerService) ServeSwaggerUI(c *fiber.Ctx) error {
	html := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: window.location.origin + '/api/docs/openapi.json',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: 'BaseLayout'
      })
    }
  </script>
</body>
</html>`
	c.Type("html")
	return c.SendString(html)
}
