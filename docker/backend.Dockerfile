FROM golang:1.22-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest
