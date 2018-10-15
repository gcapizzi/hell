FROM golang:latest

RUN go get -u github.com/onsi/ginkgo/ginkgo
RUN go get -u github.com/onsi/gomega/...

WORKDIR /go/src/github.com/gcapizzi/hell
CMD bash
