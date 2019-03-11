## 
## base image - defines the operating system layer for the build
## -------------------------------------------------------------
## use this to adjust the version of golang you want a build with
ARG GOLANG_VERSION=1.11.5
## use this to adjust the version of alpine to run for the build
ARG ALPINE_VERSION=3.9
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS base
## allow for passing in of any additional packages you might need
ARG ADDITIONAL_APKS
## due diligence
RUN apk update --no-cache
RUN apk upgrade --no-cache
## go modules dependencies
RUN apk add --no-cache git
## without these ssl/tls will not work
RUN apk add --no-cache ca-certificates && update-ca-certificates
## other development tooling
RUN apk add --no-cache make

##
## development image - where things are actually built
## ---------------------------------------------------
FROM base as development
## what should we name our binary? (default indicates "app")
ARG BIN_NAME=app
## any extension we would like for our binary? (default indicates nothing)
ARG BIN_EXT
## relative path to the binary from the working directory
ARG BIN_PATH=bin
## which architecture should we build for? (default indicates amd64)
ARG GOARCH=amd64
## which operating system should we build for? (default indicates linux)
ARG GOOS=linux
## should we use static linking? (default indicates yes)
ARG CGO_ENABLED=0
## should we use go modules for the dependencies? (default indicates yes)
ARG GO111MODULE=on
## use something GOPATH/GOROOT friendly - don't anger the gods
WORKDIR /go/src/${BIN_NAME}
## process dependencies first to take advantage of docker caching
COPY ./Makefile ./Makefile
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN make deps
## process everything else
COPY . /go/src/${BIN_NAME}
RUN make compile.linux
## generate a hash
RUN sha256sum ${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} | cut -d " " -f 1 > ${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256
## move things to where they should be
RUN mv /go/src/${BIN_NAME}/${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} /${BIN_NAME}
RUN mv /go/src/${BIN_NAME}/${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256 /${BIN_NAME}.sha256
RUN ln -s /${BIN_NAME} /_
RUN chmod +x /_
## let it start
ENTRYPOINT ["/_"]

##
# production image - the really small image
# -----------------------------------------
FROM scratch AS production
## what should we name our binary? (default indicates "app")
ARG BIN_NAME=app
WORKDIR /
## copy everything over from the development build image
COPY --from=base /etc/ssl/certs /etc/ssl/certs
COPY --from=development /${BIN_NAME} /${BIN_NAME}
COPY --from=development /${BIN_NAME}.sha256 /${BIN_NAME}.sha256
COPY --from=development /_ /_
## let it start
ENTRYPOINT ["/_"]
## if you're on openshift, you'll need to define this to define your application's ports
# EXPOSE 65534
