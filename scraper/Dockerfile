# Stage 1: Build the binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the scraper binary
RUN CGO_ENABLED=0 GOOS=linux go build -o scraper-bin ./scraper/main.go

# Stage 2: Create a tiny runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/scraper-bin .
COPY --from=builder /app/aniscale-key.json .

CMD ["./scraper-bin"]