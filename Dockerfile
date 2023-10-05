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

RUN apt-get update && \
    apt-get install -y gnupg && \
    apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10

RUN echo "deb http://repo.mongodb.org/apt/debian stretch/mongodb-org/3.6 main" | tee /etc/apt/sources.list.d/mongodb-org-3.6.list && \
    apt-get update && \
    apt-get install -y mongodb-org

RUN mkdir -p /data/db

COPY movies.json /movies.json
COPY comments.json /comments.json

RUN go build -o main

EXPOSE 6000

CMD mongosh --fork --logpath /var/log/mongodb.log && \
    mongo --host localhost --eval "db = db.getSiblingDB('ecal'); db.createCollection('movies'); db.movies.insertMany([]);" && \
    mongo --host localhost --eval "db = db.getSiblingDB('ecal'); db.createCollection('comments'); db.comments.insertMany([]);" && \
    mongoimport --host localhost --db reach-engine --collection movies --type json --file /movies.json --jsonArray && \
    mongoimport --host localhost --db reach-engine --collection comments --type json --file /comments.json --jsonArray && \
    ./main

#CMD ["./main"]