package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddEndpoint(t *testing.T) {
	assertEndpointResult(t, "/api/add", `{"a":6,"b":7}`, `{"result":13}`)
}

func TestSubtractEndpoint(t *testing.T) {
	assertEndpointResult(t, "/api/subtract", `{"a":10,"b":4}`, `{"result":6}`)
}

func TestMultiplyEndpoint(t *testing.T) {
	assertEndpointResult(t, "/api/multiply", `{"a":6,"b":7}`, `{"result":42}`)
}

func TestDivideEndpoint(t *testing.T) {
	assertEndpointResult(t, "/api/divide", `{"a":10,"b":4}`, `{"result":2.5}`)
}

func TestDivideEndpointRejectsZeroDivisor(t *testing.T) {
	assertEndpointError(t, "/api/divide", `{"a":1,"b":0}`, http.StatusUnprocessableEntity, "division by zero")
}

func TestPowerEndpoint(t *testing.T) {
	assertEndpointResult(t, "/api/power", `{"a":2,"b":10}`, `{"result":1024}`)
}

func TestSquareRootEndpoint(t *testing.T) {
	assertEndpointResult(t, "/api/square-root", `{"a":81}`, `{"result":9}`)
}

func TestSquareRootEndpointRejectsNegativeNumber(t *testing.T) {
	assertEndpointError(t, "/api/square-root", `{"a":-1}`, http.StatusUnprocessableEntity, "cannot calculate the square root of a negative number")
}

func TestPercentageEndpoint(t *testing.T) {
	assertEndpointResult(t, "/api/percentage", `{"a":20,"b":80}`, `{"result":16}`)
}

func TestPercentageEndpointAvoidsIntermediateOverflow(t *testing.T) {
	assertEndpointResult(t, "/api/percentage", `{"a":1e308,"b":100}`, `{"result":1e+308}`)
}

func TestEndpointValidation(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		body    string
		message string
	}{
		{name: "empty body", path: "/api/add", message: "request body is required"},
		{name: "malformed JSON", path: "/api/add", body: `{`, message: "request body contains malformed JSON"},
		{name: "missing a", path: "/api/add", body: `{"b":1}`, message: "a is required"},
		{name: "missing b", path: "/api/add", body: `{"a":1}`, message: "b is required"},
		{name: "missing unary operand", path: "/api/square-root", body: `{}`, message: "a is required"},
		{name: "wrong operand type", path: "/api/add", body: `{"a":"one","b":2}`, message: "a must be a number"},
		{name: "unknown field", path: "/api/add", body: `{"a":1,"b":2,"c":3}`, message: `request body contains unknown field \"c\"`},
		{name: "multiple values", path: "/api/add", body: `{"a":1,"b":2} {}`, message: "request body must contain one JSON object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertEndpointError(t, tt.path, tt.body, http.StatusBadRequest, tt.message)
		})
	}
}

func TestEndpointRejectsOversizedBody(t *testing.T) {
	body := `{"a":1,"b":2,"padding":"` + strings.Repeat("x", 1<<20) + `"}`
	assertEndpointError(t, "/api/add", body, http.StatusBadRequest, "request body must not exceed 1 MB")
}

func TestHealthEndpoint(t *testing.T) {
	assertEndpointResult(t, "/health", "", `{"status":"ok"}`)
}

func assertEndpointResult(t *testing.T, path, body, want string) {
	t.Helper()
	method := http.MethodPost
	if path == "/health" {
		method = http.MethodGet
	}
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	recorder := httptest.NewRecorder()
	NewHandler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("%s status = %d, want %d; body = %s", path, recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if got := strings.TrimSpace(recorder.Body.String()); got != want {
		t.Errorf("%s body = %s, want %s", path, got, want)
	}
}

func assertEndpointError(t *testing.T, path, body string, status int, message string) {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	recorder := httptest.NewRecorder()
	NewHandler().ServeHTTP(recorder, request)
	if recorder.Code != status {
		t.Fatalf("%s status = %d, want %d; body = %s", path, recorder.Code, status, recorder.Body.String())
	}
	if message != "" && !strings.Contains(recorder.Body.String(), message) {
		t.Errorf("%s body = %s, want error containing %q", path, recorder.Body.String(), message)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}
}
