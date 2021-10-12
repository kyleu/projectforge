FROM golang:alpine

RUN apk add --update --no-cache ca-certificates

ENTRYPOINT ["/projectforge", "-a", "0.0.0.0"]
EXPOSE 40000
COPY projectforge /
