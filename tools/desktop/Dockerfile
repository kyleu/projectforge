FROM ghcr.io/gythialy/golang-cross:v1.17.1-1 AS builder

LABEL maintainer="kyle@kyleu.com"
ENV PATH=${OSX_CROSS_PATH}/target/bin:$PATH:${GOPATH}/bin
ENV CGO_ENABLED=1

RUN apt-get update && apt-get install -y --no-install-recommends sqlite3 webkit2gtk-4.0

WORKDIR /src

ADD "./go.mod" "/src/go.mod"
ADD "./go.sum" "/src/go.sum"

RUN go mod download

ADD . /src/
RUN /src/tools/desktop/package.sh

FROM ubuntu:bionic
COPY --from=builder /src/dist /dist
