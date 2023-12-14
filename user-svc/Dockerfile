FROM golang:latest

RUN go version

ENV GOPATH=/

COPY ./ ./

RUN go mod download

RUN go build -o user-svc ./cmd/server.go