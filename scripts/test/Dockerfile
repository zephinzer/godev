FROM zephinzer/golang-dev:latest as development
COPY . /go/src/app

FROM development as build
RUN build

FROM scratch as production
COPY --from=build /go/src/app/app /app
ENTRYPOINT [ "/app" ]
