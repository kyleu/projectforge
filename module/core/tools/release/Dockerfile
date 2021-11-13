FROM golang:alpine

RUN apk add --update --no-cache ca-certificates libc6-compat

# main http port
EXPOSE {{{ .Port }}}{{{$inc := .}}}{{{ range $k, $v := .PortOffsets }}}
# {{{ $k }}} port
EXPOSE {{{ $inc.PortIncremented $v }}}{{{ end }}}

ENTRYPOINT ["/{{{ .Exec }}}", "-a", "0.0.0.0"]

COPY {{{ .Exec }}} /
