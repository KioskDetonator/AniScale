# Stage 1: Build the binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the Go module files from the root
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the notifier binary
RUN CGO_ENABLED=0 GOOS=linux go build -o notifier-bin ./notifier/main.go

# Stage 2: Create a tiny runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/notifier-bin .
# Copy the key file so the app can authenticate with GCP
COPY --from=builder /app/aniscale-key.json .

CMD ["./notifier-bin"]