ARG BASE_IMAGE=golang
ARG BASE_TAG=1.11.5-alpine3.9
FROM ${BASE_IMAGE}:${BASE_TAG} as base
# due diligence
RUN apk update --no-cache && apk upgrade --no-cache
# add common development tools
RUN apk add git make bash curl jq g++ ca-certificates
# add ssl capabilities
RUN update-ca-certificates

FROM base AS build
WORKDIR /go/build
COPY . /go/build
RUN make compile.linux

FROM base AS production
COPY --from=build /go/build/bin/godev-linux-amd64 /bin/godev-linux-amd64
RUN ln -s /bin/godev /bin/godev-linux-amd64