FROM golang:latest AS builder 

RUN apt-get update && apt-get install -y \
    gcc \
    musl-tools \
    librdkafka-dev \
    pkg-config 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY . . 
RUN go build -o main ./cmd/main.go 


FROM alpine:latest
RUN apk add --no-cache \
    librdkafka \
    libc6-compat
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]

