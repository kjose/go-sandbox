#FROM docker/whalesay:latest
#RUN apt-get -y update && apt-get install -y fortunes
#CMD /usr/games/fortune -a | cowsay

#FROM ubuntu:latest
#RUN apt-get -y update && apt-get install -y curl

FROM golang:1.15
WORKDIR /go/src
COPY . .
RUN go get -d -v ./...
RUN go build -o ../bin/app

ENTRYPOINT /go/bin/app

MAINTAINER kevin.jose.92@gmail.com

# To build this Dockerfile
# docker build -t go-server .
# 
