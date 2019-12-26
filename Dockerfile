FROM golang:1.13.5-alpine3.10 as build-env

WORKDIR /go/src/app
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wasme


FROM alpine:3.10

COPY --from=build-env /go/src/app/wasme /usr/local/bin/wasme
ENTRYPOINT ["/usr/local/bin/wasme"]