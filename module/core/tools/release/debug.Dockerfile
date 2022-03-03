FROM golang:alpine

LABEL "org.opencontainers.image.authors"="{{{ .Info.AuthorName }}}"
LABEL "org.opencontainers.image.source"="{{{ .Info.Sourcecode }}}"
LABEL "org.opencontainers.image.vendor"="{{{ .Info.Org }}}"
LABEL "org.opencontainers.image.title"="{{{ .Title }}}"
LABEL "org.opencontainers.image.description"="{{{ .Info.Summary }}}"

RUN apk add --update --no-cache ca-certificates tzdata bash curl htop libc6-compat

RUN apk add --no-cache ca-certificates dpkg gcc git musl-dev \
    && mkdir -p "$GOPATH/src" "$GOPATH/bin" \
    && chmod -R 777 "$GOPATH" \
    && go get github.com/go-delve/delve/cmd/dlv

SHELL ["/bin/bash", "-c"]

# main http port
EXPOSE {{{ .Port }}}{{{$inc := .}}}{{{ range $k, $v := .PortOffsets }}}
# {{{ $k }}} port
EXPOSE {{{ $inc.PortIncremented $v }}}{{{ end }}}

WORKDIR /

ENTRYPOINT ["/{{{ .Exec }}}", "-a", "0.0.0.0"]

COPY {{{ .Exec }}} /{{{ .ExtraFilesDocker }}}
