# Golang Development Image
This repository contains a Docker image used for an easier Golang 1.x development experience

# Requirements
You'll need:

1. Docker
2. A sense of adventure

# Opinions forced on you

1. You should be using Git
1. Usage of `go mod` is the only way to manage dependencies
1. Having a vendor library is good for your health
1. Live-reload is happiness on a CLI
1. Live-reload should also live-install dependencies added to your `*.go`

# Usage

Run the following in your directory of choice:

```sh
docker run -it -u ${UID} -v "$(pwd):/go/src/app" zephinzer/golang:latest init;
```

Setup a `Dockerfile` in the directory of your choice with:

```
FROM zephinzer/golang:latest as development
COPY . /go/src/app
ENTRYPOINT [ "start" ]

FROM development as build
RUN build
ENTRYPOINT [ "/go/src/app/app" ]

FROM scratch as production
COPY --from=build /go/src/app/app /
ENTRYPOINT [ "/app" ]
```

Add a `main.go` to the same directory:

```go
package main

import (
  "fmt"
)

func main() {
  fmt.Println("hello world")
}
```

Run the following to start your application in live-reload:

```sh
docker run -it \
  -u ${UID} \
  -v "$(pwd):/go/src/app" \
  zephinzer/golang:latest start;
```

If you need to expose ports (3000 in the example):

```sh
docker run -it \
  -u ${UID} \
  -v "$(pwd):/go/src/app" \
  -p "3000:3000" \
  zephinzer/golang:latest start;
```

Run the following to build your application:

```sh
docker run -it -u ${UID} -v "$(pwd):/go/src/app" zephinzer/golang:latest build;
```

You should now see a binary named `app` in that directory which you can redistribute.

Alternatively, to package it into a nice Docker image:

```sh
docker build \
  -t yourusername/imagename:latest \
  .
```

Run `docker images` and you should see your image. This image is built from `scratch`, which is super minimal.

Try to run it:

```sh
docker run yourusername/imagename:latest;
```

# Development

## Testing
Tests are contained [in the `./test` directory](./test) but you can run it from the root using `make test`.

