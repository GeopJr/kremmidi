FROM golang:1.19-buster

MAINTAINER GeopJr <docker@geopjr.dev>

WORKDIR /go/src/app
COPY . .

RUN make

CMD ["/go/src/app/kremmidi"]