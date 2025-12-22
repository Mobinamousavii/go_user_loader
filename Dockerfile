FROM docker.arvancloud.ir/golang:1.23-alpine AS builder

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go 



FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app .
CMD ["./app", "--file" ,"users.csv" ,"--port","8080"]
