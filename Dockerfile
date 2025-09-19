FROM golang:1.24.0-alpine AS builder 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY . . 

RUN go build -o main ./cmd/main.go 

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]

