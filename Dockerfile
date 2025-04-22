# Build stage
FROM golang:1.24.2-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/server

# Final image
FROM debian:bullseye-slim

# ✅ Instala certificados
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080
ENTRYPOINT ["./app"]
