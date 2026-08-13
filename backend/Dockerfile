FROM golang:1.24-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /calculator-api ./cmd/api

FROM alpine:3.22
RUN adduser -D -H appuser
COPY --from=build /calculator-api /usr/local/bin/calculator-api
USER appuser
EXPOSE 8081
CMD ["calculator-api"]
