ARG BASE_IMAGE=golang
ARG BASE_TAG=1.12.0-alpine3.9
FROM ${BASE_IMAGE}:${BASE_TAG} as base
# due diligence
RUN apk update --no-cache && apk upgrade --no-cache
# add common development tools
RUN apk add git make bash curl jq g++ ca-certificates
# add ssl capabilities
RUN update-ca-certificates

FROM base AS build
ARG VERSION=0.0.0
ENV VERSION=${VERSION}
ARG COMMIT=0000000
ENV COMMIT=${COMMIT}
WORKDIR /go/build
COPY . /go/build
RUN make compile
RUN chmod +x -R /go/build/bin/*
RUN ln -s /go/build/bin/godev-${VERSION}-linux-amd64 /bin/godev
RUN chmod +x /bin/godev

FROM base AS production
ARG VERSION=0.0.0
ENV VERSION=${VERSION}
COPY --from=build /go/build/bin/godev-${VERSION}-linux-amd64 /bin/godev
