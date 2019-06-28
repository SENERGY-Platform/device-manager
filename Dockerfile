FROM golang:1.12


COPY . /go/src/device-manager
WORKDIR /go/src/device-manager

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./device-manager