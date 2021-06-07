FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod downloadd

ENTRYPOINT go run cmd/videoserver/main.go