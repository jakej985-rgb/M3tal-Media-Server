FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the entire root directory
COPY . .

# Build for Linux from the cmd/m3tal directory
WORKDIR /app/cmd/m3tal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o m3tal-core .

FROM alpine:latest
RUN apk add --no-cache ca-certificates smartmontools radeontop

WORKDIR /root/
COPY --from=builder /app/cmd/m3tal/m3tal-core .

# State directory mapping
ENV STATE_DIR=/docker/state
ENV TZ=America/Denver

CMD ["./m3tal-core"]
