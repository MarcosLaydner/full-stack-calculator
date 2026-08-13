package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDocsEndpointServesSwaggerUI(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/docs", nil)
	recorder := httptest.NewRecorder()
	NewHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q, want text/html; charset=utf-8", contentType)
	}
	if !strings.Contains(recorder.Body.String(), "SwaggerUIBundle") {
		t.Error("response does not contain Swagger UI")
	}
	if !strings.Contains(recorder.Body.String(), `"openapi": "3.1.0"`) {
		t.Error("response does not contain the embedded OpenAPI specification")
	}
	for _, path := range []string{"/api/add", "/api/subtract", "/api/multiply", "/api/divide", "/api/power", "/api/square-root", "/api/percentage"} {
		if !strings.Contains(recorder.Body.String(), `"`+path+`"`) {
			t.Errorf("embedded OpenAPI specification is missing %s", path)
		}
	}
}

func TestOpenAPISpecIsNotExposedSeparately(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil)
	recorder := httptest.NewRecorder()
	NewHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}
