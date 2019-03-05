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
ARG GOARCH=amd64
ARG GOOS=linux
ARG VERSION=0.0.0
ENV VERSION=${VERSION}
ARG COMMIT=0000000
ENV COMMIT=${COMMIT}
WORKDIR /go/build
COPY . /go/build
RUN make compile

FROM base AS production
COPY --from=build /go/build/bin/godev /bin/godev
