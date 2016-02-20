FROM golang:1.6.0-alpine

WORKDIR /go/src/frister.net/go/droplink
ADD . /go/src/frister.net/go/droplink

RUN go install -v

CMD []
ENTRYPOINT ["/go/bin/droplink"]
