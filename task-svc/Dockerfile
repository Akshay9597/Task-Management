FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

RUN go mod download

RUN go build -o task-svc ./cmd/server.go