FROM golang:1.9
WORKDIR /go/src/github.com/OpenBazaar/feeproxy
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o /feeproxy .

FROM alpine:3.6
RUN apk --update --no-cache add ca-certificates
COPY --from=0 /feeproxy /feeproxy
ENTRYPOINT ["/feeproxy"]