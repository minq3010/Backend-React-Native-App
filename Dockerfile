FROM golang:1.22.2-alpine3.19

WORKDIR /src/app

RUN apk add --no-cache git

RUN go install github.com/cosmtrek/air@v1.40.4

COPY . .

RUN go mod tidy