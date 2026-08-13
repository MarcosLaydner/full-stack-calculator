# Full-stack Calculator

A technical-assignment-ready calculator with a React and TypeScript frontend and
a Go REST API. The application supports addition, subtraction, multiplication,
division, exponentiation, square root, and percentage calculations.

## Architecture

```text
frontend/             React, TypeScript, Vite, Vitest
  src/api/            HTTP client and API types
  src/components/     Calculator UI
backend/              Go HTTP microservice
  cmd/api/            Application entry point
  internal/calculator Domain logic
  internal/httpapi/   HTTP routing and validation
```

The backend separates arithmetic from transport concerns, which keeps the core
logic reusable and easy to unit test. The frontend similarly isolates API access
from display components. In Docker, Nginx serves the compiled frontend and proxies
`/api` and `/health` to the backend, so the browser only talks to one origin.

## Run with Docker (*Recomended*)

Docker and Docker Compose are the only prerequisites.

```bash
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080). Stop the stack with:

```bash
docker compose down
```

## Run locally

Prerequisites: Node.js 22+, npm 10+, and Go 1.24+.

Install frontend dependencies and start both development servers in separate
terminals:

```bash
npm install
npm run dev
```

```bash
cd backend
go run ./cmd/api
```

The frontend is available at `http://localhost:5173`; Vite proxies API requests
to the Go service at `http://localhost:8081`.

## Interactive API documentation

Swagger UI is available at [http://localhost:8080/docs](http://localhost:8080/docs)
when using Docker, or [http://localhost:8081/docs](http://localhost:8081/docs)
when running the backend directly. The OpenAPI 3.1 document is embedded directly
into that page. Swagger UI includes editable examples and supports sending requests
to every operation with **Try it out**.

## API

### Operations

Each operation has its own endpoint and handler:

| Operation | Endpoint | Request body |
| --- | --- | --- |
| Addition | `POST /api/add` | `{"a": number, "b": number}` |
| Subtraction | `POST /api/subtract` | `{"a": number, "b": number}` |
| Multiplication | `POST /api/multiply` | `{"a": number, "b": number}` |
| Division | `POST /api/divide` | `{"a": number, "b": number}` |
| Exponentiation | `POST /api/power` | `{"a": number, "b": number}` |
| Square root | `POST /api/square-root` | `{"a": number}` |
| Percentage | `POST /api/percentage` | `{"a": number, "b": number}` |

Binary operation example:

```bash
curl -X POST http://localhost:8081/api/add \
  -H 'Content-Type: application/json' \
  -d '{"a":12,"b":3}'
```

```json
{"result":15}
```

Percentage calculates `a%` of `b`. Every route validates its own request shape;
unknown fields and missing operands return `400 Bad Request`.

Errors use a consistent JSON response and an appropriate HTTP status:

```json
{"error":"division by zero"}
```

### Health check

```bash
curl http://localhost:8081/health
```

## Tests and coverage

Run all tests from the repository root:

```bash
npm test
```

Run each layer with coverage reports:

```bash
npm run test:frontend -- --coverage
cd backend && go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The tests cover arithmetic behavior and edge cases, HTTP validation and status
codes, and the primary frontend calculation and error workflows.
