FROM golang:latest

RUN apt-get update && apt-get install -y golang

ENV GOPATH=/go

ENV GOPROXY=https://proxy.golang.org,direct

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . .
COPY . ./
RUN go mod download

RUN go build -o main .
CMD ["./main"]