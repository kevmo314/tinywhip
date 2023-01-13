FROM golang:alpine AS builder

COPY . /go/src/github.com/kevmo314/tinywhip

WORKDIR /go/src/github.com/kevmo314/tinywhip

RUN go build -o /go/bin/tinywhip cmd/main.go

FROM alpine:latest

COPY --from=builder /go/bin/tinywhip /tinywhip

ENTRYPOINT ["/tinywhip"]
