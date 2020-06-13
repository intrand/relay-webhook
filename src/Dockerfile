FROM golang:alpine as builder
RUN apk add --no-cache git ca-certificates
RUN adduser -D -g '' app
COPY . $GOPATH/src/mypackage/myapp/
WORKDIR $GOPATH/src/mypackage/myapp/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/relay-webhook

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/relay-webhook /go/bin/relay-webhook

EXPOSE 8000
USER app

ENTRYPOINT ["/go/bin/relay-webhook"]
