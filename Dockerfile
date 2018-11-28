FROM golang:1.11-alpine3.8

ARG CGO_ENABLED=0
ARG GOARCH=amd64
ARG GOOS=linux

RUN apk update --no-cache
RUN apk add --no-cache bash inotify-tools curl git gcc g++ vim

COPY --chown=1000:0 ./scripts /scripts
ENV PATH="/scripts:${PATH}"
RUN chmod +x -R /scripts

WORKDIR /go/src/app

RUN printf -- "alias ll=\"ls -al\";\nexport PATH=\"${PATH}\";" > /.profile

ENV GO111MODULE auto
ENV CGO_ENABLED ${CGO_ENABLED}
ENV GOARCH ${GOARCH}
ENV GOOS ${GOOS}

RUN mkdir -p /.cache/go-build && chmod 777 -R /.cache

ENTRYPOINT [ "entrypoint" ]