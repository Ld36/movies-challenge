FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o movie-service ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/movie-service .
COPY --from=builder /app/movies.json .

EXPOSE 50051

CMD ["./movie-service"]