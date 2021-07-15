# Modified from github.com/neilotoole/xcgo
ARG OSX_SDK="MacOSX10.15.sdk"
ARG OSX_CODENAME="catalina"
ARG OSX_VERSION_MIN="10.10"
ARG OSX_SDK_BASEURL="https://github.com/neilotoole/xcgo/releases/download/v0.1"
ARG OSX_SDK_SUM="d97054a0aaf60cb8e9224ec524315904f0309fbbbac763eb7736bdfbdad6efc8"
ARG OSX_CROSS_COMMIT="bee9df60f169abdbe88d8529dbcc1ec57acf656d"
ARG LIBTOOL_VERSION="2.4.6_1"
ARG LIBTOOL_BASEURL="https://github.com/neilotoole/xcgo/releases/download/v0.1"

FROM ubuntu:bionic AS golangcore
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash curl git jq lsb-core software-properties-common

ENV GOPATH="/go"
RUN mkdir -p "${GOPATH}/src"

RUN add-apt-repository -y ppa:longsleep/golang-backports
RUN apt update && apt install -y golang-1.16
RUN ln -s /usr/lib/go-1.16 /usr/lib/go
RUN ln -s /usr/lib/go/bin/go /usr/bin/go
RUN ln -s /usr/lib/go/bin/gofmt /usr/bin/gofmt

RUN go version

FROM golangcore AS devtools
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential clang cmake file gcc-mingw-w64 gcc-mingw-w64-i686 gcc-mingw-w64-x86-64 less \
    libc6-dev libc6-dev-i386 libc++-dev libltdl-dev libsqlite3-dev libssl-dev libxml2-dev \
    llvm man parallel patch sqlite3 tree vim webkit2gtk-4.0 wget xz-utils zlib1g-dev zsh

#### Modified from github.com/tpoechtrager/osxcross
FROM devtools AS gotools
ARG OSX_SDK
ARG OSX_CODENAME
ARG OSX_SDK_BASEURL
ARG OSX_SDK_SUM
ARG OSX_CROSS_COMMIT
ARG OSX_VERSION_MIN
ARG LIBTOOL_VERSION
ARG LIBTOOL_BASEURL
ENV OSX_CROSS_PATH=/osxcross

WORKDIR "${OSX_CROSS_PATH}"
RUN git clone https://github.com/tpoechtrager/osxcross.git . && git checkout -q "${OSX_CROSS_COMMIT}" && rm -rf ./.git

RUN curl -fsSL "${OSX_SDK_BASEURL}/${OSX_SDK}.tar.xz" -o "${OSX_CROSS_PATH}/tarballs/${OSX_SDK}.tar.xz"
RUN echo "${OSX_SDK_SUM}"  "${OSX_CROSS_PATH}/tarballs/${OSX_SDK}.tar.xz" | sha256sum -c -

RUN UNATTENDED=yes OSX_VERSION_MIN=${OSX_VERSION_MIN} ./build.sh

RUN mkdir -p "${OSX_CROSS_PATH}/target/SDK/${OSX_SDK}/usr/"
RUN curl -fsSL "${LIBTOOL_BASEURL}/libtool-${LIBTOOL_VERSION}.${OSX_CODENAME}.bottle.tar.gz" \
	| gzip -dc \
	| tar xf - -C "${OSX_CROSS_PATH}/target/SDK/${OSX_SDK}/usr/" --strip-components=2 "libtool/${LIBTOOL_VERSION}/include/" "libtool/${LIBTOOL_VERSION}/lib/"

WORKDIR /root

#### app container
FROM gotools AS builder
LABEL maintainer="kyle@kyleu.com"
ENV PATH=${OSX_CROSS_PATH}/target/bin:$PATH:${GOPATH}/bin
ENV CGO_ENABLED=1

WORKDIR /src

RUN git init

RUN go get -u github.com/pyros2097/go-embed

ADD "./go.mod" "/src/go.mod"
ADD "./go.sum" "/src/go.sum"

RUN go mod download

ADD . /src/
RUN /src/tools/desktop/package.sh

FROM ubuntu:bionic
COPY --from=builder /src/dist /dist
