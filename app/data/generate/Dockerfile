#
# base image - defines the operating system layer for the build
#
# use this to adjust the version of golang you want a build with
ARG GOLANG_VERSION=1.11.5
# use this to adjust the version of alpine to run for the build
ARG ALPINE_VERSION=3.9
FROM golang:${GOLANG_VERSION}}-alpine${ALPINE_VERSION}} AS base
# allow for passing in of any additional packages you might need
ARG ADDITIONAL_APKS
RUN apk update --no-cache
RUN apk upgrade --no-cache
RUN apk add --no-cache git ca-certificates && update-ca-certificates

#
# development image - where things are actually built
#
FROM base as development
# what should we name our binary? (default indicates "app")
ARG BIN_NAME=app
# any extension we would like for our binary? (default indicates nothing)
ARG BIN_EXT
# which architecture should we build for? (default indicates amd64)
ARG GOARCH=amd64
# which operating system should we build for? (default indicates linux)
ARG GOOS=linux
# should we use static linking? (default indicates yes)
ARG CGO_ENABLED=0
# should we use go modules for the dependencies? (default indicates yes)
ARG GO111MODULE=on
# use something GOPATH/GOROOT friendly - don't anger the gods
WORKDIR /go/src
# copy all we have so far
COPY . /go/src
# do the build
RUN go build \
  -a \
  -ldflags "-extldflags -static" \
  -o /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}
# generate a hash
RUN sha256sum /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} | cut -d " " -f 1 > /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256

#
# production image - the really small image
#
FROM scratch AS production
# what should we name our binary? (default indicates "app")
ARG BIN_NAME=app
# any extension we would like for our binary? (default indicates nothing)
ARG BIN_EXT
# which architecture should we build for? (default indicates amd64)
ARG GOARCH=amd64
# which operating system should we build for? (default indicates linux)
ARG GOOS=linux
# copy everything over from the previous build images
COPY --from=development /etc/ssl/certs /etc/ssl/certs
COPY --from=development /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} /bin/app
COPY --from=development /go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256 /bin/app.sha256
WORKDIR /
# let it start
ENTRYPOINT ["/bin/app"]
# if you're on openshift, you'll need to define this
# EXPOSE 65534
