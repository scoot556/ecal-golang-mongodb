FROM ubuntu:20.04
FROM mongo:3
FROM golang:latest

RUN apt-get update && apt-get install -y golang

RUN apt-get install -y git

ENV GOPATH=/go

ENV GOPROXY=https://proxy.golang.org,direct

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . .
RUN go mod download
RUN go install

#RUN wget http://archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.1f-1ubuntu2_amd64.deb
#
#RUN dpkg -i libssl1.1_1.1.1f-1ubuntu2_amd64.deb
#
#RUN apt-get update && \
#    apt-get install -y gnupg wget && \
#    wget -qO - https://www.mongodb.org/static/pgp/server-4.4.asc | apt-key add -
#
#RUN echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/4.4 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-4.4.list
#
#RUN apt-get update && \
#    apt-get install -y mongodb-org
#
#
#RUN mkdir -p /data/db

#COPY movies.json /movies.json
#COPY comments.json /comments.json

COPY movies.json .
COPY comments.json .

#RUN mongoimport --collection='comments' --file='comments.json' --jsonArray && \
#    mongoimport --collection='movies' --file='movies.json' --jsonArray




RUN go build -o main

EXPOSE 8080


CMD ["./main"]