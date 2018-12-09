FROM golang:1.11-alpine3.8

# set standard variables (this works for linux)
## turns the CGO off (allows for static linking - we can use the `scratch` image)
ARG CGO_ENABLED=0
## sets the CPU architecture to amd64
ARG GOARCH=amd64
## sets the target system to linux
ARG GOOS=linux

# update the packages and add useful build tools
RUN apk update --no-cache
RUN apk add --no-cache bash inotify-tools curl git gcc g++ vim

# sets scripts to be runnable by any user belonging to the `root` group
COPY --chown=1000:0 ./scripts /scripts
ENV PATH="/scripts:${PATH}"
RUN chmod +x -R /scripts
WORKDIR /go/src/app

# for developer happiness
RUN printf -- "alias ll=\"ls -al\";\nexport PATH=\"${PATH}\";" > /.profile

# set the environment variables for the building
ENV GO111MODULE auto
ENV CGO_ENABLED ${CGO_ENABLED}
ENV GOARCH ${GOARCH}
ENV GOOS ${GOOS}

RUN mkdir -p /.cache/go-build && chmod 777 -R /.cache

COPY --chown=1000:0 ./scripts /scripts
ENV PATH="/scripts:${PATH}"
RUN chmod +x -R /scripts

ENTRYPOINT [ "entrypoint" ]
