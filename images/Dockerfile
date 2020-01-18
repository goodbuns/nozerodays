FROM golang:1.14beta1-alpine3.11 AS builder

RUN apk update && \
    apk add --no-cache git=2.24.1-r0 ca-certificates && \
    update-ca-certificates

RUN adduser -D -g '' app

WORKDIR $HOME/src
COPY . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -o /go/bin/nozerodays github.com/goodbuns/nozerodays/cmd

FROM alpine:3.11.2

RUN apk add --no-cache tzdata
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/nozerodays /go/bin/nozerodays

USER app

ENTRYPOINT ["/go/bin/nozerodays"]
