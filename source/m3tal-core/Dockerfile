FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire source directory
COPY . .

# Build for Linux from the m3tal-core directory
WORKDIR /app/m3tal-core
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o m3tal-core main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates smartmontools radeontop

WORKDIR /root/
COPY --from=builder /app/m3tal-core/m3tal-core .

# State directory mapping
ENV STATE_DIR=/docker/state
ENV TZ=America/Denver

CMD ["./m3tal-core"]
