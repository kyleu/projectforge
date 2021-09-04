FROM golang:alpine

RUN apk add --update --no-cache ca-certificates

ENTRYPOINT ["/{{{ .Exec }}}", "-a", "0.0.0.0"]
EXPOSE {{{ .Port }}}
COPY {{{ .Exec }}} /
