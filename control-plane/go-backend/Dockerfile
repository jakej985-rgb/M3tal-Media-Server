# --- Build Stage ---
FROM golang:1.24-alpine AS builder

# NO NETWORK REQUIRED: Using local vendor folder
WORKDIR /app

# Copy all source and dependencies
COPY . .

# Build for Linux using the vendor folder
# -mod=vendor tells Go to use the local vendor directory instead of the network
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o m3tal-backend .

# --- Final Stage ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates radeontop smartmontools pciutils

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/m3tal-backend .

# Set environment defaults
ENV STATE_DIR=/docker/control-plane/state
ENV TZ=America/Denver

# Run the backend
CMD ["./m3tal-backend"]
