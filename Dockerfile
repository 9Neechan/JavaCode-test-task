# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main main.go

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xz
RUN mv migrate.linux-amd64 /usr/bin/migrate

# Run stage
FROM alpine
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /usr/bin/migrate /usr/bin/migrate
COPY db/migration ./migration

COPY config.env . 
COPY start.sh .
COPY wait-for.sh .

RUN chmod +x /app/start.sh
RUN chmod +x /app/wait-for.sh

CMD ["./main"]
ENTRYPOINT ["/app/start.sh"]
