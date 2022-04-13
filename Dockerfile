# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./goMongo .

EXPOSE 8080

CMD ["./goMongo"]


