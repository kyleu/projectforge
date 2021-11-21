FROM golang:alpine

LABEL "org.opencontainers.image.authors"="Kyle U"
LABEL "org.opencontainers.image.source"="https://github.com/kyleu/projectforge"
LABEL "org.opencontainers.image.vendor"="kyleu"
LABEL "org.opencontainers.image.title"="Project Forge"
LABEL "org.opencontainers.image.description"="Build and maintain feature-rich applications using Golang"

RUN apk add --update --no-cache ca-certificates libc6-compat

# main http port
EXPOSE 40000
# marketing port
EXPOSE 40001

ENTRYPOINT ["/projectforge", "-a", "0.0.0.0"]

COPY projectforge /
