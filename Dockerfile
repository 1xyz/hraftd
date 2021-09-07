FROM golang:alpine AS builder

RUN apk update && apk add make git build-base && \
     rm -rf /var/cache/apk/*

ADD . /go/src/github.com/1xyz/hraftd
WORKDIR /go/src/github.com/1xyz/hraftd
RUN make release/linux

###

FROM alpine:latest AS hraftd

RUN apk update && apk add ca-certificates bash
WORKDIR /root/
COPY --from=builder /go/src/github.com/1xyz/hraftd/bin/linux/hraftd .