FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
# Skip go get/mod download since we'll use vendor or it'll download on build
COPY . .
RUN go build -o m3tal-core main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates smartmontools radeontop
WORKDIR /root/
COPY --from=builder /app/m3tal-core .
# State directory mapping
ENV STATE_DIR=/docker/state
CMD ["./m3tal-core"]
