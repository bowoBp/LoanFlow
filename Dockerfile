# Menggunakan image Go untuk build
FROM golang:1.23 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory di dalam container
WORKDIR /app

# Copy semua file kode ke dalam container
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build aplikasi
RUN go build -o main ./cmd/api/main.go

# Menggunakan scratch sebagai base image untuk aplikasi runtime
FROM alpine:latest

# Copy file binary dari stage builder
COPY --from=builder /app/main /app/main

# Copy file .env jika ada
COPY --from=builder /app/.env /app/.env

# Set working directory di dalam container
WORKDIR /app

# Ekspose port default aplikasi
EXPOSE 8080

# Perintah untuk menjalankan aplikasi
CMD ["./main"]
