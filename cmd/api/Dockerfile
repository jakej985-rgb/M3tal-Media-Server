# --- Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire root directory (since pkg is now in root)
COPY . .

# Build for Linux from the cmd/api directory
WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o m3tal-backend .

# Use a stable alpine version
FROM alpine:3.20

RUN echo "https://dl-cdn.alpinelinux.org/alpine/v3.20/community" >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache ca-certificates radeontop smartmontools pciutils || true

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cmd/api/m3tal-backend .

# Set environment defaults
ENV STATE_DIR=/docker/control-plane/state
ENV TZ=America/Denver

# Run the backend
CMD ["./m3tal-backend"]
