# --- Build Stage ---
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build for Linux (Targeting the container environment)
# CGO_ENABLED=0 ensures a static binary with no external dependencies
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o m3tal-backend .

# --- Final Stage ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/m3tal-backend .

# Set environment defaults
ENV STATE_DIR=/docker/control-plane/state
ENV TZ=America/Denver

# Run the backend
CMD ["./m3tal-backend"]
