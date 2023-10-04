FROM mongo:3

RUN apt-get update && apt-get install -y golang

RUN go get github.com/gin-gonic/gin
RUN go get gopkg.in/mgo.v2
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/jwuensche/trommel

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

COPY . .

RUN go build -o main

EXPOSE 8080

CMD ["./main"]