FROM golang:1.10.1-alpine3.7 as build

WORKDIR /go/src/frister.net/go/droplink
ADD . /go/src/frister.net/go/droplink

RUN go install -v


FROM alpine:3.7
RUN apk update && \
   apk add ca-certificates && \
   update-ca-certificates && \
   rm -rf /var/cache/apk/*

RUN mkdir -p /opt/droplink
WORKDIR /opt/droplink
COPY --from=build /go/bin/droplink .
COPY --from=build /go/src/frister.net/go/droplink/upload.html .
COPY --from=build /go/src/frister.net/go/droplink/media ./media

CMD []
ENTRYPOINT ["/opt/droplink/droplink"]
