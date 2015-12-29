FROM golang:1.5

ENV GO15VENDOREXPERIMENT 1

WORKDIR /go/src/frister.net/go/droplink
ADD . /go/src/frister.net/go/droplink

RUN go install -v

CMD []
ENTRYPOINT ["/go/bin/droplink"]
