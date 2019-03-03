# for instances where we wish to use an internally built golang
ARG IMAGE_NAME=golang
ARG IMAGE_TAG=1.11.5-alpine3.9

# define the base image
FROM ${IMAGE_NAME}:${IMAGE_TAG}

# set standard variables (this has been tested for linux only)
# for a full list of GOOS/GOARCH, see https://github.com/golang/go/blob/master/src/go/build/syslist.go
## turns the CGO off (allows for static linking so we can use the `scratch` image)
ARG CGO_ENABLED=0
## sets the CPU architecture to amd64 (windows, use: 386)
ARG GOARCH=amd64
## sets the target system to linux (windows, use: windows, mac, use: darwin)
ARG GOOS=linux

# update the packages and add useful build tools
RUN apk update --no-cache
RUN apk add --no-cache bash inotify-tools curl git gcc g++ vim

# for developer happiness when debugging
RUN printf -- "alias ll=\"ls -al\";\nexport PATH=\"${PATH}\";" > /.profile

# set the environment variables for the build process
ENV GO111MODULE auto
ENV CGO_ENABLED ${CGO_ENABLED}
ENV GOARCH ${GOARCH}
ENV GOOS ${GOOS}

# create cache directories so this works even if consumer doesn't map volumes
RUN mkdir -p /.cache/go-build && chmod 777 -R /.cache

# sets scripts to be runnable by any user belonging to the `root` group
COPY --chown=1000:0 ./scripts /scripts
ENV PATH="/scripts:${PATH}"
RUN chmod +x -R /scripts

# famous last workds
WORKDIR /go/src/app
ENTRYPOINT [ "entrypoint" ]
