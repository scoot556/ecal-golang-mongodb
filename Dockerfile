FROM ubuntu:20.04
FROM mongo:3
FROM golang:latest

RUN apt-get update && apt-get install -y golang

RUN apt-get install -y git

ENV GOPATH=/go

ENV GOPROXY=https://proxy.golang.org,direct


#RUN go get -u github.com/gin-gonic/gin
#RUN go install gopkg.in/mgo.v3
#RUN go install gopkg.in/mgo.v3/bson
#RUN go install github.com/jwuensche/trommel

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . .
RUN go mod download
RUN go install

RUN wget http://archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.1f-1ubuntu2_amd64.deb

RUN dpkg -i libssl1.1_1.1.1f-1ubuntu2_amd64.deb

RUN apt-get update && \
    apt-get install -y gnupg wget && \
    wget -qO - https://www.mongodb.org/static/pgp/server-4.4.asc | apt-key add -

# Add the MongoDB repository for Ubuntu 20.04
RUN echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/4.4 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-4.4.list

# Install MongoDB
RUN apt-get update && \
    apt-get install -y mongodb-org


RUN mkdir -p /data/db

COPY movies.json /movies.json
COPY comments.json /comments.json

RUN go build -o main

EXPOSE 6000

#CMD ["mongod"]

CMD mongod --fork --logpath /var/log/mongodb.log && \
    mongosh "use ecal;" && \
    mongosh "db.createCollection('movies'); db.movies.insertMany([]);" && \
    mongosh " db.createCollection('comments'); db.comments.insertMany([]);" && \
    mongoimport --host localhost --db ecal --collection movies --file movies.json --type json && \
    mongoimport --host localhost --db ecal --collection comments --file comments.json --type json && \
    ./main

#CMD ["./main"]