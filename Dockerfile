FROM golang:1.11

WORKDIR /go/src/github.com/picfight/pfcregtest
COPY . .

RUN apt-get update && apt-get upgrade -y && apt-get install -y rsync

