FROM golang:1.17-alpine

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod . 
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o example main.go
CMD ["./example"]