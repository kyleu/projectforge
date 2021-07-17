# Build image
FROM golang:alpine AS builder

ENV GOFLAGS="-mod=readonly"

RUN apk add --update --no-cache bash ca-certificates make git curl build-base

RUN mkdir /app

WORKDIR /app

RUN go get -u github.com/pyros2097/go-embed
RUN go get -u github.com/valyala/quicktemplate
RUN go get -u github.com/valyala/quicktemplate/qtc

ADD ./go.mod        /app/go.mod
ADD ./go.sum        /app/go.sum

RUN go mod download

ADD ./app           /app/app
ADD ./bin           /app/bin
ADD ./main.go       /app/main.go
ADD ./Makefile      /app/Makefile
ADD ./views         /app/views

RUN go mod download

RUN set -xe && bash -c 'make build-release-ci'

RUN mv build/release /build

# Final image
FROM alpine

RUN apk add --update --no-cache ca-certificates tzdata bash curl

SHELL ["/bin/bash", "-c"]

COPY --from=builder /build/* /usr/local/bin/

EXPOSE 10101
CMD ["projectforge"]
