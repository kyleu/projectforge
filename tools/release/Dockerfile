FROM golang:alpine

RUN apk add --update --no-cache ca-certificates

# main http port
EXPOSE 40000
# marketing port
EXPOSE 40001

ENTRYPOINT ["/projectforge", "-a", "0.0.0.0"]

COPY projectforge /
