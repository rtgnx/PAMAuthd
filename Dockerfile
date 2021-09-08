FROM golang:1.16.4-alpine3.13 as builder

RUN apk update && apk add alpine-sdk linux-pam linux-pam-dev

RUN mkdir -p /go/src/github.com/Reverse-Labs/pamauthd
WORKDIR /go/src/github.com/Reverse-Labs/pamauthd
COPY . .
RUN go fmt
RUN addgroup testuser
RUN adduser -G testuser -D testuser
RUN passwd -d testuser
RUN go test

FROM alpine:3.13

COPY --from=builder /usr/bin/pamauthd /usr/bin/pamauthd

RUN go build -o /usr/bin/pamauthd 

CMD ["/usr/bin/pamauthd", "serve"]