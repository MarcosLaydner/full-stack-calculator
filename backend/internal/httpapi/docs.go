package httpapi

import (
	"bytes"
	_ "embed"
	"log"
	"net/http"
)

//go:embed openapi.json
var openAPI []byte

const swaggerUI = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Calculator API documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      spec: __OPENAPI_SPEC__,
      dom_id: '#swagger-ui',
      deepLinking: true,
      displayRequestDuration: true,
      tryItOutEnabled: true
    })
  </script>
</body>
</html>`

func docs(w http.ResponseWriter, _ *http.Request) {
	page := bytes.Replace([]byte(swaggerUI), []byte("__OPENAPI_SPEC__"), openAPI, 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(page); err != nil {
		log.Printf("write Swagger UI response: %v", err)
	}
}
